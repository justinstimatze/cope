package main

// The mid-session experiment.
//
// The question is narrow and it is the one the whole project turns on: does the
// sliced mid-session injection do anything a static file does not? Both arms
// carry the SessionStart card, so the card is held constant and what varies is
// only the reminder chosen from measured output.
//
// Alternate refresher windows are held back. A turn is attributed to whichever
// arm was live when it was written, and the two arms are compared on hits per
// thousand characters — the rate MEASUREMENTS.md already argues for, because
// the share of turns hit tracks how long the turns were.
//
// What this cannot answer: whether the prose got better. It counts the moves
// the card names. A reply scoring zero may still be worse to read than one
// scoring three, and nothing here measures that.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/justinstimatze/cope/internal/state"
)

type armTotals struct {
	turns      int
	chars      int
	hits       int
	injChars   int
	injWindows int
	byRule     map[string]int

	// The same totals split by distance from the window boundary. This is where
	// the largest effect in the data lives, and it is not the arm: the rate
	// climbs with distance under injection and under hold alike, more steeply
	// under hold. Whatever it is, a reminder no arm received cannot be causing
	// it.
	byDist [distBuckets]halfTotals

	// The same totals split by where the turn sat in its session. Persona
	// evaluation is two measurements, not one: how a single reply scores, and
	// whether the voice holds across the session. A per-turn rate can stay flat
	// while the voice erodes, because every turn clears the per-turn bar and
	// the drift only shows up in the trajectory. Comparing the arms on the mean
	// alone measures the first and is blind to the second, which is the one a
	// mid-session reminder is supposed to act on.
	early, late halfTotals
}

type halfTotals struct {
	turns int
	chars int
	hits  int
}

// Distance from the boundary of the window a turn was written in. Zero is the
// turn whose prompt carried the injection; the rest are downstream of it.
//
// The arm is attributed per window, not per injection, so a window's later turns
// are counted as injected while having seen the reminder only at second hand. If
// the reminder did anything, distance 0 is where it would show.
const (
	distNear = iota // the turn that carried the reminder
	distMid         // one or two turns after it
	distFar         // three or more
	distBuckets
)

var distLabels = [distBuckets]string{"d0", "d1-2", "d3+"}

func distBucket(d int) int {
	switch {
	case d == 0:
		return distNear
	case d <= 2:
		return distMid
	default:
		return distFar
	}
}

// What a rule detects, which turns out to matter more than which arm it fired
// under. The distance gradient is almost entirely `formatting`, and those rules
// track what kind of reply the turn called for: a status report uses labels and
// even blocks, and reporting clusters in the dense stretches of a session. The
// `rhetoric` rules are the ones that read as voice rather than layout, and they
// show no gradient at all.
//
// Grouping them is a judgement about what each rule is sensitive to, not a fact
// the card records. It is here rather than in `card/rules.json` because only two
// of these are POSTPROC patterns; the rest are shape rules in internal/scan.
// `bold_label` stays in the map after leaving the card on 2026-07-30. The logs
// this reads are historical, and dropping the id would silently reclassify every
// turn it ever fired on as `other` — changing what past runs say rather than
// what future ones measure.
func ruleClass(id string) string {
	switch id {
	case "bold_label", "labelled_opening", "paragraph_uniformity", "echoed_heading":
		return "formatting"
	case "flip", "clause_symmetry", "fragment_run":
		return "rhetoric"
	case "dangling_end", "buried_decision", "ask_not_last", "forked_end":
		return "ending"
	case "apology", "self_postmortem", "announced_length":
		return "deference"
	case "cross_turn_repeat", "repeated_opening":
		return "repetition"
	case "unverified_done", "loop_ask":
		return "handoff"
	}
	return "other"
}

// deference is its own class rather than part of rhetoric because the three
// rules in it came from outside this repo — the recurring complaints in Zvi
// Mowshowitz's July 2026 Opus 5 reaction roundup — and the point of separating
// them is to be able to ask whether they move independently of the rules mined
// from these transcripts. Folding them into rhetoric would answer that question
// by construction.
// repetition is alone in its class because it is the only rule that reads the
// window rather than the reply. Every other id here can fire on a turn read in
// isolation; this one cannot fire at all until the session has a history, so
// its rate is a function of session length in a way none of the others are.
// Averaged into any of them it would drag that class's rate up over a long run
// and look like the class getting worse.
//
// handoff is the loop lane's replacement for ending. The two never appear on
// the same turn by construction — CheckLane drops one set and adds the other —
// so folding them together would produce a class whose rate is really a
// measure of how much of the run was unattended.
var classOrder = []string{"formatting", "rhetoric", "ending", "deference", "repetition", "handoff", "other"}

func (h *halfTotals) add(r turnRecord) {
	h.turns++
	h.chars += r.Chars
	h.hits += len(r.Rules)
}

func newArmTotals() *armTotals { return &armTotals{byRule: map[string]int{}} }

// add takes the turn's position in its session as a fraction, so the totals can
// be split into the opening and closing halves of a session, and its distance
// from the boundary of the window it was written in.
func (a *armTotals) add(r turnRecord, frac float64, dist int) {
	a.byDist[distBucket(dist)].add(r)
	a.addTotals(r, frac)
}

func (a *armTotals) addTotals(r turnRecord, frac float64) {
	a.turns++
	a.chars += r.Chars
	a.hits += len(r.Rules)
	if r.InjectChars > 0 {
		a.injChars += r.InjectChars
		a.injWindows++
	}
	for _, id := range r.Rules {
		a.byRule[id]++
	}
	if frac < 0.5 {
		a.early.add(r)
	} else {
		a.late.add(r)
	}
}

// drift is how much this arm's rate climbs from the opening half of a session
// to the closing half, in hits per thousand characters.
//
// A positive number is the voice loosening as the session runs on. The claim a
// mid-session reminder makes is that this number is smaller under injection
// than under hold — which is a different claim from the mean being lower, and
// the one the rate ratio cannot see.
func (a *armTotals) drift() (float64, bool) {
	if a.early.chars == 0 || a.late.chars == 0 {
		return 0, false
	}
	return per1k(a.late.hits, a.late.chars) - per1k(a.early.hits, a.early.chars), true
}

// meanInject is the average size of what this arm put in front of the writer,
// over the turns that carried one. Reported because arms differ in length as
// well as in content, and a reader comparing two rates deserves to see whether
// they also differ by a factor of three in how much text they cost.
func (a *armTotals) meanInject() int {
	if a.injWindows == 0 {
		return 0
	}
	return a.injChars / a.injWindows
}

// per1k is the rate the arms are compared on. Zero characters reads as zero
// rather than as an error: an arm with no turns yet is the normal state on the
// first day of collection.
func per1k(hits, chars int) float64 {
	if chars == 0 {
		return 0
	}
	return float64(hits) / float64(chars) * 1000
}

// readTurnLog returns the last record for each turn. A turn is scored twice —
// at Stop against a partial transcript, and again on the next prompt once the
// reply is complete — and only the second reading saw the whole thing.
//
// Records with no key cannot be matched to a later rescore and are all kept.
// That is the honest reading: without an id there is nothing to say two records
// describe one turn.
func readTurnLog(path string) ([]turnRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var order []string
	latest := map[string]turnRecord{}
	var unkeyed []turnRecord
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 64<<10), 4<<20)
	for s.Scan() {
		var r turnRecord
		if err := json.Unmarshal(s.Bytes(), &r); err != nil {
			continue // a half-written line costs one turn, not the report
		}
		if r.Key == "" {
			unkeyed = append(unkeyed, r)
			continue
		}
		id := r.Session + "\x00" + r.Key
		if _, seen := latest[id]; !seen {
			order = append(order, id)
		}
		latest[id] = r
	}
	out := make([]turnRecord, 0, len(order)+len(unkeyed))
	for _, id := range order {
		out = append(out, latest[id])
	}
	return append(out, unkeyed...), s.Err()
}

// rateRatio compares the two arms as a Poisson rate ratio with a 95% interval.
//
// Hits are counts over a character exposure, which is what a Poisson rate is
// for. The interval is the standard log-scale one: an interval containing 1
// means the data does not separate the arms, and saying so is the point of
// running this rather than reading two numbers and picking a story.
func rateRatio(hitsA, charsA, hitsB, charsB int) (rr, lo, hi float64, ok bool) {
	if hitsA == 0 || hitsB == 0 || charsA == 0 || charsB == 0 {
		return 0, 0, 0, false
	}
	rr = per1k(hitsA, charsA) / per1k(hitsB, charsB)
	se := math.Sqrt(1/float64(hitsA) + 1/float64(hitsB))
	return rr, rr * math.Exp(-1.96*se), rr * math.Exp(1.96*se), true
}

func runABReport(path string) {
	if path == "-" {
		dir := stateDir()
		if dir == "" {
			fmt.Fprintln(os.Stderr, "cope-gate: no state directory")
			os.Exit(1)
		}
		path = turnLogPath(dir)
	}
	records, err := readTurnLog(path)
	if os.IsNotExist(err) {
		// Not an error. The log appears when the first window is claimed, so
		// this is what the report says on the day the experiment starts.
		fmt.Printf("%s does not exist yet.\nIt is written once --ab claims its first refresher window.\n", filepath.Clean(path))
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
		os.Exit(1)
	}

	// Arms are discovered from the log rather than fixed here, so a run that
	// added an arm reports it without the reader passing a flag to say so.
	// Where each turn sat in its session, needed before any of it is totalled:
	// the drift split is by position, and position is only knowable once the
	// session's full length is.
	seen := map[string]int{}
	for _, r := range records {
		if r.Arm != "" {
			seen[r.Session]++
		}
	}

	arms := map[string]*armTotals{}
	sessions := map[string]bool{}
	at := map[string]int{}
	builds := map[string]int{}
	var buildOrder []string
	// Window boundaries are recovered from the log rather than recorded: the arm
	// changes exactly when a window is claimed, and `catchUp` rescores the
	// previous turn before the claim, so the first record bearing a new arm is
	// the turn whose prompt carried the reminder.
	lastArm := map[string]string{}
	dist := map[string]int{}
	// Every class shares the same character exposure — the denominator is all the
	// prose written, not just the prose that rule class happened to fire on — so
	// they are all present from the start.
	classDist := map[string]*[distBuckets]halfTotals{}
	for _, c := range classOrder {
		classDist[c] = &[distBuckets]halfTotals{}
	}
	total, skipped := 0, 0
	for _, r := range records {
		if r.Arm == "" {
			skipped++
			continue
		}
		if prev, seen := lastArm[r.Session]; !seen || prev != r.Arm {
			dist[r.Session] = 0
		} else {
			dist[r.Session]++
		}
		lastArm[r.Session] = r.Arm
		b := distBucket(dist[r.Session])
		for _, id := range r.Rules {
			classDist[ruleClass(id)][b].hits++
		}
		for _, c := range classOrder {
			classDist[c][b].chars += r.Chars
		}
		if _, known := arms[r.Arm]; !known {
			arms[r.Arm] = newArmTotals()
		}
		v := r.Version
		if v == "" {
			v = "unrecorded"
		}
		if builds[v] == 0 {
			buildOrder = append(buildOrder, v)
		}
		builds[v]++
		frac := float64(at[r.Session]) / float64(seen[r.Session])
		at[r.Session]++
		arms[r.Arm].add(r, frac, dist[r.Session])
		sessions[r.Session] = true
		total++
	}
	if len(arms) == 0 {
		fmt.Printf("%s\nno turns are tagged with an arm yet.\n", filepath.Clean(path))
		return
	}

	fmt.Printf("%s\n%d turns across %d session(s)", filepath.Clean(path), total, len(sessions))
	if skipped > 0 {
		fmt.Printf(", %d before the experiment started", skipped)
	}
	fmt.Print("\n")
	reportBuilds(buildOrder, builds)
	fmt.Print("\n")

	// ArmHold is the reference every other arm is measured against, and it
	// sorts first for that reason. The rest follow in a fixed order so two
	// runs of the report line up.
	order := armOrder(arms)
	fmt.Printf("%-10s %7s %11s %7s %9s %10s\n", "arm", "turns", "chars", "hits", "per 1k", "injected")
	for _, name := range order {
		a := arms[name]
		inj := "—"
		if a.meanInject() > 0 {
			inj = fmt.Sprintf("%d ch", a.meanInject())
		}
		fmt.Printf("%-10s %7d %11d %7d %9.3f %10s\n", name, a.turns, a.chars, a.hits, per1k(a.hits, a.chars), inj)
	}

	ref, hasRef := arms[state.ArmHold]
	fmt.Println()
	if !hasRef {
		fmt.Println("no held-back windows yet, so there is nothing to compare against")
	}
	for _, name := range order {
		if name == state.ArmHold || !hasRef {
			continue
		}
		a := arms[name]
		rr, lo, hi, ok := rateRatio(a.hits, a.chars, ref.hits, ref.chars)
		if !ok {
			fmt.Printf("%s/hold: not enough data in one of the arms yet\n", name)
			continue
		}
		verdict := "the interval contains 1: this data does not separate them"
		switch {
		case hi < 1:
			verdict = "scores lower than hold, and the interval clears 1"
		case lo > 1:
			verdict = "scores higher than hold, and the interval clears 1"
		}
		fmt.Printf("%s/hold: %.3f (95%% CI %.3f–%.3f) — %s\n", name, rr, lo, hi, verdict)
	}

	// The trajectory half of the question. A reminder that holds the voice
	// steady shows up here as a smaller climb, whether or not it moved the mean.
	fmt.Printf("\n%-10s %9s %9s %9s   (hits per 1k, by position in session)\n", "arm", "early", "late", "drift")
	for _, name := range order {
		a := arms[name]
		d, ok := a.drift()
		if !ok {
			fmt.Printf("%-10s %9s %9s %9s\n", name, "—", "—", "—")
			continue
		}
		fmt.Printf("%-10s %9.3f %9.3f %+9.3f\n", name,
			per1k(a.early.hits, a.early.chars), per1k(a.late.hits, a.late.chars), d)
	}

	// Distance from the window boundary. Read this before the arm comparison
	// above: it is the largest effect in the log and it is present in every arm,
	// including the one that receives nothing, so it is not the reminder.
	fmt.Printf("\n%-10s %9s %9s %9s   (hits per 1k, by distance from the window boundary)\n",
		"arm", distLabels[0], distLabels[1], distLabels[2])
	for _, name := range order {
		fmt.Printf("%-10s", name)
		for i := range arms[name].byDist {
			d := arms[name].byDist[i]
			if d.chars == 0 {
				fmt.Printf(" %9s", "—")
				continue
			}
			fmt.Printf(" %9.3f", per1k(d.hits, d.chars))
		}
		fmt.Println()
	}

	// And what the gradient is made of. It is nearly all `formatting`, which
	// tracks the kind of reply the turn called for rather than the voice it was
	// written in — a status report uses labels and even blocks, and reporting
	// clusters in the dense stretches of a session. `rhetoric` is the subset that
	// reads as voice, and it stays flat.
	hitsIn := func(c string) int {
		n := 0
		for i := range classDist[c] {
			n += classDist[c][i].hits
		}
		return n
	}
	allHits := 0
	for _, c := range classOrder {
		allHits += hitsIn(c)
	}
	if allHits > 0 {
		fmt.Printf("\n%-12s %9s %9s %9s %7s %7s   (all arms pooled)\n",
			"rule class", distLabels[0], distLabels[1], distLabels[2], "hits", "share")
		for _, c := range classOrder {
			n := hitsIn(c)
			if n == 0 {
				continue
			}
			fmt.Printf("%-12s", c)
			for i := range classDist[c] {
				d := classDist[c][i]
				if d.chars == 0 {
					fmt.Printf(" %9s", "—")
					continue
				}
				fmt.Printf(" %9.3f", per1k(d.hits, d.chars))
			}
			fmt.Printf(" %7d %6.0f%%\n", n, float64(n)/float64(allHits)*100)
		}
	}

	ids := map[string]bool{}
	for _, a := range arms {
		for id := range a.byRule {
			ids[id] = true
		}
	}
	if len(ids) == 0 {
		return
	}
	names := make([]string, 0, len(ids))
	for id := range ids {
		names = append(names, id)
	}
	totalFor := func(id string) int {
		n := 0
		for _, a := range arms {
			n += a.byRule[id]
		}
		return n
	}
	sort.Slice(names, func(i, j int) bool { return totalFor(names[i]) > totalFor(names[j]) })

	fmt.Printf("\n%-22s", "rule")
	for _, name := range order {
		fmt.Printf(" %9s", name)
	}
	fmt.Println("   (hits per 1k chars)")
	for _, id := range names {
		fmt.Printf("%-22s", id)
		for _, name := range order {
			fmt.Printf(" %9.3f", per1k(arms[name].byRule[id], arms[name].chars))
		}
		fmt.Println()
	}
}

// reportBuilds says which binaries wrote the data, in the order they first
// appear.
//
// One build is the normal case and prints as a single line. More than one means
// the rules moved underneath the run: `go install` replaces the file the hook
// names, so a rule added or retuned partway through scored only the turns after
// it, and the totals average across a change the numbers cannot show. That is
// not automatically fatal — a version bump that touched no rule is harmless —
// but it is the reader's call to make and they need to be told there is one.
func reportBuilds(order []string, counts map[string]int) {
	switch len(order) {
	case 0:
		return
	case 1:
		if order[0] != "unrecorded" {
			fmt.Printf("written by %s\n", order[0])
		}
		return
	}
	fmt.Printf("written by %d builds — the totals span a change in the scorer:\n", len(order))
	for _, v := range order {
		fmt.Printf("  %-40s %5d turns\n", v, counts[v])
	}
}

// armOrder puts the reference arm first and the rest in a stable order, so the
// columns of two runs of the report can be read against each other.
func armOrder(arms map[string]*armTotals) []string {
	out := make([]string, 0, len(arms))
	for name := range arms {
		out = append(out, name)
	}
	sort.Slice(out, func(i, j int) bool {
		if (out[i] == state.ArmHold) != (out[j] == state.ArmHold) {
			return out[i] == state.ArmHold
		}
		return out[i] < out[j]
	})
	return out
}
