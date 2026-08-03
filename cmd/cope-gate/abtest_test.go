package main

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justinstimatze/cope/internal/state"
)

func writeTurnLog(t *testing.T, records []turnRecord) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "turns.jsonl")
	var b strings.Builder
	for _, r := range records {
		raw, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		b.Write(raw)
		b.WriteByte('\n')
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// Every turn is written twice — at Stop against a partial transcript, and again
// on the next prompt once the reply is complete. Counting both would double the
// exposure and mix a truncated reading into the totals.
func TestReadTurnLogKeepsTheLastReadingOfEachTurn(t *testing.T) {
	path := writeTurnLog(t, []turnRecord{
		{Session: "s1", Key: "k1", Arm: state.ArmInject, Chars: 900, Rules: []string{"bold_label"}},
		{Session: "s1", Key: "k1", Arm: state.ArmInject, Chars: 1400, Rules: []string{"bold_label", "dangling_end"}},
		{Session: "s1", Key: "k2", Arm: state.ArmHold, Chars: 500, Rules: nil},
		// Same key, different session: a different turn.
		{Session: "s2", Key: "k1", Arm: state.ArmHold, Chars: 200, Rules: []string{"flip"}},
	})

	got, err := readTurnLog(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d turns, want 3: %+v", len(got), got)
	}
	if got[0].Chars != 1400 || len(got[0].Rules) != 2 {
		t.Errorf("the rescore should win: %+v", got[0])
	}
	if got[1].Key != "k2" || got[2].Session != "s2" {
		t.Errorf("turns came back out of order: %+v", got)
	}
}

// A turn with no key cannot be matched to a later rescore, so nothing may be
// dropped on the assumption that two of them describe one turn.
func TestReadTurnLogKeepsEveryUnkeyedTurn(t *testing.T) {
	path := writeTurnLog(t, []turnRecord{
		{Session: "s1", Arm: state.ArmInject, Chars: 100},
		{Session: "s1", Arm: state.ArmInject, Chars: 200},
	})
	got, err := readTurnLog(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("got %d turns, want both kept: %+v", len(got), got)
	}
}

// A half-written line costs one turn, not the report.
func TestReadTurnLogSurvivesATornLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "turns.jsonl")
	body := `{"session":"s1","key":"k1","arm":"inject","chars":100,"rules":["flip"]}` + "\n" +
		`{"session":"s1","key":"k2","arm":"hold","cha` + "\n" +
		`{"session":"s1","key":"k3","arm":"hold","chars":300,"rules":[]}` + "\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := readTurnLog(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("got %d turns, want the two well-formed ones: %+v", len(got), got)
	}
}

func TestRateRatio(t *testing.T) {
	// Half the hits over the same exposure is a ratio of 0.5.
	rr, lo, hi, ok := rateRatio(50, 100_000, 100, 100_000)
	if !ok {
		t.Fatal("a populated pair of arms should compare")
	}
	if math.Abs(rr-0.5) > 1e-9 {
		t.Errorf("rate ratio is %.4f, want 0.5", rr)
	}
	if !(lo < rr && rr < hi) {
		t.Errorf("the interval %.3f–%.3f does not bracket %.3f", lo, hi, rr)
	}
	if hi >= 1 {
		t.Errorf("a halving over 150 hits should clear 1, got upper bound %.3f", hi)
	}

	// The same rate in both arms cannot separate them, however much data.
	_, lo, hi, _ = rateRatio(500, 1_000_000, 500, 1_000_000)
	if !(lo < 1 && hi > 1) {
		t.Errorf("equal rates should give an interval containing 1, got %.3f–%.3f", lo, hi)
	}

	// Thin data must widen rather than assert.
	_, thinLo, thinHi, _ := rateRatio(2, 4_000, 4, 4_000)
	if thinLo > 1 || thinHi < 1 {
		t.Errorf("six hits should not separate the arms, got %.3f–%.3f", thinLo, thinHi)
	}

	if _, _, _, ok := rateRatio(0, 1000, 5, 1000); ok {
		t.Error("an arm with no hits has no ratio to report")
	}
}

// The report discovers its arms from the log, so adding a third one must not
// need a flag — and hold has to stay the reference every other arm is measured
// against.
func TestArmOrderPutsTheReferenceFirst(t *testing.T) {
	got := armOrder(map[string]*armTotals{
		state.ArmPositive: newArmTotals(),
		state.ArmInject:   newArmTotals(),
		state.ArmHold:     newArmTotals(),
	})
	want := []string{state.ArmHold, state.ArmInject, state.ArmPositive}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("arm order is %v, want %v", got, want)
		}
	}
}

// The positive arm is shorter than the paired one by construction, so the
// report has to carry that length or a reader cannot tell a difference in
// content from a difference in size.
func TestMeanInjectIgnoresHeldWindows(t *testing.T) {
	a := newArmTotals()
	a.add(turnRecord{Chars: 100, InjectChars: 1700}, 0, 0)
	a.add(turnRecord{Chars: 100, InjectChars: 1500}, 0.5, 1)
	a.add(turnRecord{Chars: 100}, 0.9, 2) // a turn later in the same window
	if got := a.meanInject(); got != 1600 {
		t.Errorf("mean injection is %d, want 1600 — turns carrying none should not dilute it", got)
	}
	if got := newArmTotals().meanInject(); got != 0 {
		t.Errorf("an arm that injects nothing should report 0, got %d", got)
	}
}

// A per-turn rate can hold flat while the voice erodes across a session, which
// is the failure a mid-session reminder exists to prevent. The report has to
// carry that second measurement or the arms are compared on the half of the
// question the reminder is not about.
func TestDriftSplitsOnPositionInSession(t *testing.T) {
	a := newArmTotals()
	// Opening half: clean. Closing half: the same turns, scoring more.
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip"}}, 0.0, 0)
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip"}}, 0.25, 1)
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip", "bold_label", "dangling_end"}}, 0.6, 2)
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip", "bold_label", "dangling_end"}}, 0.9, 3)

	d, ok := a.drift()
	if !ok {
		t.Fatal("both halves are populated; drift should compute")
	}
	if math.Abs(d-2.0) > 1e-9 {
		t.Errorf("drift is %.3f, want +2.000 — 1/1k early against 3/1k late", d)
	}
	// The mean is blind to it: the same four turns average 2 hits per 1k
	// whichever half they landed in.
	if got := per1k(a.hits, a.chars); math.Abs(got-2.0) > 1e-9 {
		t.Errorf("mean is %.3f, want 2.000", got)
	}

	// One-sided data cannot report a climb it did not observe.
	only := newArmTotals()
	only.add(turnRecord{Chars: 1000, Rules: []string{"flip"}}, 0.1, 0)
	if _, ok := only.drift(); ok {
		t.Error("an arm with no closing-half turns has no drift to report")
	}
}

// The hook names an absolute path and `go install` replaces the file underneath
// it, so a run can span several builds. The report has to say so: the totals
// average across whatever changed in the scorer, and only the reader can judge
// whether that mattered.
func TestReportBuildsNamesEveryBinaryThatWrote(t *testing.T) {
	out := capture(t, func() {
		reportBuilds([]string{"v0.2.0", "v0.2.1"}, map[string]int{"v0.2.0": 40, "v0.2.1": 198})
	})
	for _, want := range []string{"2 builds", "v0.2.0", "40", "v0.2.1", "198"} {
		if !strings.Contains(out, want) {
			t.Errorf("report does not mention %q:\n%s", want, out)
		}
	}

	// One build is the ordinary case and says so in one line.
	if got := capture(t, func() { reportBuilds([]string{"v0.2.1"}, map[string]int{"v0.2.1": 238}) }); !strings.Contains(got, "written by v0.2.1") || strings.Contains(got, "builds") {
		t.Errorf("a single build should print one plain line, got %q", got)
	}

	// Turns written before the version was recorded have nothing to report.
	if got := capture(t, func() { reportBuilds([]string{"unrecorded"}, map[string]int{"unrecorded": 5}) }); got != "" {
		t.Errorf("unversioned turns should print nothing, got %q", got)
	}
}

func capture(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var b strings.Builder
	if _, err := io.Copy(&b, r); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// The distance gradient is the largest effect in the log and it is not the arm,
// so the report has to carry it. Splitting on distance from the window boundary
// is what separates "the reminder worked" from "the turn that carried it was
// already different".
func TestDistanceSplitsOnPositionInWindow(t *testing.T) {
	a := newArmTotals()
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip"}}, 0, 0)                 // d0
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip", "flip"}}, 0, 2)         // d1-2
	a.add(turnRecord{Chars: 1000, Rules: []string{"flip", "flip", "flip"}}, 0, 7) // d3+

	for i, want := range []float64{1, 2, 3} {
		d := a.byDist[i]
		if got := per1k(d.hits, d.chars); got != want {
			t.Errorf("%s is %.1f per 1k, want %.1f", distLabels[i], got, want)
		}
	}

	// A turn three or more past the boundary shares a bucket with one ten past;
	// the reminder is equally second-hand to both.
	b := newArmTotals()
	b.add(turnRecord{Chars: 1000}, 0, 3)
	b.add(turnRecord{Chars: 1000}, 0, 40)
	if b.byDist[distFar].turns != 2 {
		t.Errorf("d3 and d40 belong in the same bucket, got %+v", b.byDist)
	}
}

// The gradient is nearly all formatting rules, which track what kind of reply
// the turn called for rather than the voice it was written in. Separating them
// is the difference between measuring drift and measuring how much of the
// session was spent reporting.
func TestRuleClassSeparatesLayoutFromVoice(t *testing.T) {
	for id, want := range map[string]string{
		"bold_label":           "formatting",
		"labelled_opening":     "formatting",
		"paragraph_uniformity": "formatting",
		"flip":                 "rhetoric",
		"clause_symmetry":      "rhetoric",
		"dangling_end":         "ending",
		"buried_decision":      "ending",
		"ask_not_last":         "ending",
		"something_new":        "other",
	} {
		if got := ruleClass(id); got != want {
			t.Errorf("ruleClass(%q) = %q, want %q", id, got, want)
		}
	}

	// Every class has to appear in classOrder or its hits vanish from the report.
	seen := map[string]bool{}
	for _, c := range classOrder {
		seen[c] = true
	}
	for _, id := range []string{"bold_label", "flip", "dangling_end", "something_new"} {
		if !seen[ruleClass(id)] {
			t.Errorf("%q classifies as %q, which classOrder does not list", id, ruleClass(id))
		}
	}
}
