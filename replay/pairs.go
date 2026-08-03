package main

// Blind pairwise preference: the measurement the rule counts were standing in for.
//
// Every other cope measurement scores prose with regexes, and four fifths of
// what those regexes catch turned out to track what a turn was for rather than
// how it was written. This asks the only question the project actually cares
// about — is the prose better — by generating two replies to the same real
// prompt under different system prompts and handing them over unlabelled.
//
// Three conditions, because a null is only worth reading if the instrument can
// register a hit:
//
//	bare      the framing alone
//	baseline  framing + whatever prose guidance is already deployed (-baseline)
//	card      baseline + the voice card under test (-rules, or the built-in one)
//
// baseline vs card is the question: what does the card add on top of what is
// already in place. bare vs card is the positive control, the widest gap
// available. If the control is picked reliably and the test is not, the card
// does nothing beyond the baseline. If neither is picked reliably, the
// instrument is blind and no measurement it has produced means anything.
//
// Nothing about the two documents is baked in. -rules takes any card the gate
// can render and -baseline takes any file, so this tests someone else's card
// against someone else's CLAUDE.md without touching the code — which is also
// why the maintainer's own style section is not a const in here.
//
// Cases come only from transcripts with no <voice-card> in them. Prior
// assistant turns anchor register hard, and history written under the card
// would seed both arms with the thing being tested.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

type condition struct {
	ID     string
	Blocks []string
}

// conditions builds the arms from the documents supplied at run time. The test
// arm is always "card" against whatever "rival" resolves to, and the two ways it
// resolves ask different questions.
//
// With a second card, it is one card against the other and neither arm carries
// the baseline, because the question is which card wins rather than what either
// adds. With a baseline document, the card is stacked on top of it, because the
// question is what the card adds to guidance you already have.
//
// Neither one is a real configuration for the test arm, and that is not a
// missing argument: it means the card is being compared to nothing, which is
// what the control arm already does.
func conditions(card, baseline, rival string) map[string]condition {
	c := map[string]condition{
		"bare": {"bare", []string{framing}},
		"card": {"card", []string{framing, card}},
	}
	switch {
	case rival != "":
		c["rival"] = condition{"rival", []string{framing, rival}}
	case baseline != "":
		c["rival"] = condition{"baseline", []string{framing, baseline}}
		c["card"] = condition{"card", []string{framing, baseline, card}}
	}
	return c
}

// Pair is one screen of the blind test: a prompt and two replies to it.
type Pair struct {
	Index   int    `json:"index"`
	Prompt  string `json:"prompt"`
	Context string `json:"context"`

	Left      string `json:"left"`
	Right     string `json:"right"`
	LeftCond  string `json:"left_cond"`
	RightCond string `json:"right_cond"`
	Set       string `json:"set"` // "test" or "control"

	// Target is the condition whose voice description the reader is shown, on
	// a discrimination run. Empty on a preference run, where no description is
	// shown at all — which is what keeps one pairs.json readable by both pages.
	Target string `json:"target,omitempty"`
	// Judge is the model's answer to the same question the reader gets:
	// "left", "right", or "unsure". Empty when -judge did not run.
	Judge string `json:"judge,omitempty"`
}

func runPairs(gate, glob, outDir, exclude, rules, baselinePath, rivalRules string, nTest, nControl, historyN, seed, maxFiles int, maxTok int64, conc int, dry, append_, discriminate, judge bool) error {
	ask := func(rulesPath, mode string) (string, error) {
		argv := []string{mode}
		if rulesPath != "" {
			argv = []string{"--rules", rulesPath, mode}
		}
		return run(gate, argv...)
	}
	render := func(rulesPath string) (string, error) { return ask(rulesPath, "--inject") }

	card, err := render(rules)
	if err != nil {
		return fmt.Errorf("reading the card: %w", err)
	}

	// On a discrimination run the reader is shown a voice and asked which reply
	// was written under it, so every arm that can be a target needs a
	// description. bare has none by construction and is never a target: with no
	// guidance at all it is the thing the target is being told apart from.
	descs := map[string]string{}
	if discriminate {
		for id, path := range map[string]string{"card": rules, "rival": rivalRules} {
			if id == "rival" && rivalRules == "" {
				continue
			}
			d, err := ask(path, "--describe")
			if err != nil {
				return fmt.Errorf("describing %s: %w", id, err)
			}
			if strings.TrimSpace(d) == "" {
				return fmt.Errorf("the %s card describes no voice: -discriminate shows the reader a "+
					"VOICE block and asks which reply carries it, and a card with an empty one has "+
					"nothing to recognise", id)
			}
			descs[id] = strings.TrimSpace(d)
		}
	}

	var rival string
	if rivalRules != "" {
		if rival, err = render(rivalRules); err != nil {
			return fmt.Errorf("reading -rival: %w", err)
		}
	}

	var baseline string
	if baselinePath != "" {
		raw, err := os.ReadFile(baselinePath)
		if err != nil {
			return fmt.Errorf("reading -baseline: %w", err)
		}
		baseline = strings.TrimSpace(string(raw))
	}
	if baseline == "" && rival == "" && nTest > 0 {
		return fmt.Errorf("-test-pairs needs -baseline or -rival: the test arm compares the card " +
			"either against guidance you already have or against a second card, and with neither " +
			"it is the same comparison as -control-pairs. Pass one of them, or run control pairs only")
	}
	conds := conditions(card, baseline, rival)

	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	rng := rand.New(rand.NewPCG(uint64(seed), 0xc09e))
	rng.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	// The corpus is 11G across hundreds of sessions and reading all of it to
	// pick thirty prompts is most of an hour. Shuffle first, then take the
	// first maxFiles that qualify, so the sample still spans the corpus.
	var clean []string
	for _, f := range files {
		if len(clean) >= maxFiles {
			break
		}
		// Subagent transcripts are structured reports to a parent process, a
		// different genre from a reply to a person, and they are all card-free
		// because subagents get no SessionStart hook. Sampling them would fill
		// the test with prose neither condition was written for.
		if strings.HasPrefix(filepath.Base(f), "agent-") {
			continue
		}
		st, err := os.Stat(f)
		if err != nil || st.Size() < 150<<10 || st.Size() > 4<<20 {
			continue
		}
		carded, err := hasCard(f)
		if err != nil || carded {
			continue
		}
		clean = append(clean, f)
	}
	if len(clean) == 0 {
		return fmt.Errorf("no card-free transcripts match %s", glob)
	}

	// Prompts already judged in an earlier run. Re-sampling one carries the
	// reader's memory of how he called it last time into the new pair, and that
	// is the single form of unblinding a shuffled side assignment cannot undo.
	read := map[string]bool{}
	for _, path := range strings.Split(exclude, ",") {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		prior, err := loadPairs(path)
		if err != nil {
			return fmt.Errorf("reading -exclude %s: %w", path, err)
		}
		for _, p := range prior {
			read[p.Prompt] = true
		}
	}
	if len(read) > 0 {
		fmt.Printf("excluding %d prompts already read\n", len(read))
	}

	var cases []Case
	for _, f := range clean {
		cs, err := readCases(f, historyN, 400)
		if err != nil {
			continue
		}
		for _, c := range cs {
			if read[c.Prompt] {
				continue
			}
			// An upper bound as well as a lower one: these get read by a
			// person, one pair at a time, and a wall of text is how a blind
			// test turns into an abandoned one.
			if len(c.Reply) > 2500 || !humanTyped(c.Prompt) || carriesSecret(c) || !intactHistory(c) {
				continue
			}
			cases = append(cases, c)
		}
	}
	want := nTest + nControl
	if len(cases) < want {
		return fmt.Errorf("only %d usable cases across %d transcripts, need %d", len(cases), len(clean), want)
	}

	pool := len(cases)
	rng.Shuffle(len(cases), func(i, j int) { cases[i], cases[j] = cases[j], cases[i] })
	cases = cases[:want]

	fmt.Printf("%d card-free transcripts, %d usable cases, %d sampled, %d calls\n",
		len(clean), pool, want, want*2)
	if dry {
		for i, c := range cases {
			set := "test"
			if i >= nTest {
				set = "control"
			}
			fmt.Printf("  %-8s %s\n", set, c)
		}
		return nil
	}

	client := anthropic.NewClient()
	pairs := make([]Pair, want)
	var (
		mu    sync.Mutex
		usage totals
		sem   = make(chan struct{}, conc)
		wg    sync.WaitGroup
	)

	for i, c := range cases {
		against, set := "rival", "test"
		if i >= nTest {
			against, set = "bare", "control"
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, c Case, against, set string) {
			defer wg.Done()
			defer func() { <-sem }()

			a, ua, err := generate(client, c, conds["card"], maxTok)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  case %d card: %v\n", i, err)
				return
			}
			b, ub, err := generate(client, c, conds[against], maxTok)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  case %d %s: %v\n", i, against, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			usage.add(ua)
			usage.add(ub)
			// Sides are assigned after the run, as a balanced permutation
			// rather than a coin flip per pair. A flip is unbiased in
			// expectation and lopsided in practice at this n, and a reader with
			// any left/right habit would turn that lopsidedness into an effect.
			pairs[i] = Pair{
				Index:     i,
				Prompt:    c.Prompt,
				Context:   contextOf(c),
				Set:       set,
				Left:      a,
				LeftCond:  "card",
				Right:     b,
				RightCond: against,
			}
			fmt.Printf("  case %-3d %-8s card %4d ch, %s %4d ch\n", i, set, len(a), against, len(b))
		}(i, c, against, set)
	}
	wg.Wait()

	done := pairs[:0]
	for _, p := range pairs {
		if p.Left != "" && p.Right != "" {
			done = append(done, p)
		}
	}
	// Renumber so the page's indices are contiguous after any dropped call.
	for i := range done {
		done[i].Index = i
	}
	if len(done) == 0 {
		return fmt.Errorf("every call failed")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	dataPath := filepath.Join(outDir, "pairs.json")

	if append_ {
		prior, err := loadPairs(dataPath)
		if err != nil {
			return fmt.Errorf("appending to %s: %w", dataPath, err)
		}
		seen := map[string]bool{}
		for _, p := range prior {
			seen[p.Prompt] = true
		}
		kept := 0
		for _, p := range done {
			if !seen[p.Prompt] {
				prior = append(prior, p)
				kept++
			}
		}
		fmt.Printf("appended %d of %d new pairs to %d existing\n", kept, len(done), len(prior)-kept)
		done = prior
	}
	done = arrange(done, rng)
	if discriminate {
		balanceTargets(done, descs, rng)
		if judge {
			runJudge(done, descs, conc, &usage)
		}
	}
	raw, err := json.MarshalIndent(done, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(dataPath, raw, 0o644); err != nil {
		return err
	}
	page, write := "pairs.html", writePage
	if discriminate {
		page, write = "discriminate.html", func(path string, pairs []Pair) error {
			return writeDiscPage(path, pairs, descs)
		}
	}
	pagePath := filepath.Join(outDir, page)
	if err := write(pagePath, done); err != nil {
		return err
	}

	fmt.Printf("\n%d pairs written\n", len(done))
	usage.report()
	fmt.Printf("\nblind test:  %s\nanswer key:  %s\n", pagePath, dataPath)
	return nil
}

// arrange fixes both things about presentation that a partial run would
// otherwise get wrong: which side the card renders on, and what a reader who
// stops early ends up having judged.
// dedupe drops what should never have been sampled and what got sampled twice.
// It runs on every arrange, including -repage, so a filter written after a run
// still cleans that run's data without paying for new prose.
func dedupe(pairs []Pair) []Pair {
	seen := map[string]bool{}
	out := pairs[:0]
	for _, p := range pairs {
		if seen[p.Prompt] || !humanTyped(p.Prompt) {
			continue
		}
		seen[p.Prompt] = true
		out = append(out, p)
	}
	return out
}

func arrange(pairs []Pair, rng *rand.Rand) []Pair {
	pairs = dedupe(pairs)
	balanceSides(pairs, rng)
	// Interleave the sets. Written in set order, the control sits entirely
	// behind the test, so anyone who stops at pair 20 has judged no control at
	// all — and the control is the only thing that makes a null on the test
	// readable. Shuffling the order means every prefix samples both.
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	for i := range pairs {
		pairs[i].Index = i
	}
	return pairs
}

func loadPairs(path string) ([]Pair, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pairs []Pair
	return pairs, json.Unmarshal(raw, &pairs)
}

// balanceSides puts the card on the left in exactly half of each set, then
// shuffles which half. Presentation order is the one nuisance variable a
// reader cannot be blinded to, so it gets balanced rather than randomised.
func balanceSides(pairs []Pair, rng *rand.Rand) {
	bySet := map[string][]int{}
	for i, p := range pairs {
		bySet[p.Set] = append(bySet[p.Set], i)
	}
	for _, idx := range bySet {
		rng.Shuffle(len(idx), func(a, b int) { idx[a], idx[b] = idx[b], idx[a] })
		for k, i := range idx {
			wantLeft := k < len(idx)/2
			p := &pairs[i]
			if (p.LeftCond == "card") != wantLeft {
				p.Left, p.Right = p.Right, p.Left
				p.LeftCond, p.RightCond = p.RightCond, p.LeftCond
			}
		}
	}
}

// repage rebuilds the page from an existing run. Reading thirty pairs is the
// expensive half of this experiment for the model and the cheap half for the
// wallet; this makes the page revisable without paying for the prose twice.
func repage(jsonPath string, seed int) error {
	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	var pairs []Pair
	if err := json.Unmarshal(raw, &pairs); err != nil {
		return err
	}
	pairs = arrange(pairs, rand.New(rand.NewPCG(uint64(seed), 0xc09e)))
	out, err := json.MarshalIndent(pairs, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, out, 0o644); err != nil {
		return err
	}
	page := filepath.Join(filepath.Dir(jsonPath), "pairs.html")
	if err := writePage(page, pairs); err != nil {
		return err
	}
	fmt.Printf("%d pairs repaged\n%s\n", len(pairs), page)
	return nil
}

// generate writes one reply under one condition. The last system block carries
// the cache breakpoint: within a condition it is identical across every case,
// so it is paid for once and read back at a tenth for the rest of the run.
func generate(client anthropic.Client, c Case, cond condition, maxTok int64) (string, anthropic.Usage, error) {
	system := make([]anthropic.TextBlockParam, 0, len(cond.Blocks))
	for i, b := range cond.Blocks {
		p := anthropic.TextBlockParam{Text: b}
		if i == len(cond.Blocks)-1 {
			p.CacheControl = anthropic.NewCacheControlEphemeralParam()
		}
		system = append(system, p)
	}

	msgs := make([]anthropic.MessageParam, 0, len(c.History)+1)
	for j, m := range c.History {
		block := anthropic.NewTextBlock(m.Text)
		if j == len(c.History)-1 {
			block.OfText.CacheControl = anthropic.NewCacheControlEphemeralParam()
		}
		if m.Role == "assistant" {
			msgs = append(msgs, anthropic.NewAssistantMessage(block))
		} else {
			msgs = append(msgs, anthropic.NewUserMessage(block))
		}
	}
	msgs = append(msgs, anthropic.NewUserMessage(anthropic.NewTextBlock(c.Prompt)))

	// A 529 drops the whole pair, because a pair needs both arms and there is no
	// half of one worth showing. Losing pairs to server load is losing power on
	// the only question this tool asks, so overload and rate limiting are worth
	// waiting out rather than reporting. Nothing else retries: a 400 is a bug in
	// the request and retrying it just pays twice for the same error.
	var (
		resp *anthropic.Message
		err  error
	)
	for attempt := 0; ; attempt++ {
		resp, err = client.Messages.New(context.Background(), anthropic.MessageNewParams{
			Model:     model,
			MaxTokens: maxTok,
			System:    system,
			Messages:  msgs,
		})
		if err == nil {
			break
		}
		var apiErr *anthropic.Error
		if attempt == 4 || !errors.As(err, &apiErr) || (apiErr.StatusCode != 529 && apiErr.StatusCode != 429 && apiErr.StatusCode < 500) {
			return "", anthropic.Usage{}, err
		}
		time.Sleep(time.Duration(1<<attempt) * 2 * time.Second)
	}
	var reply strings.Builder
	for _, b := range resp.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			if reply.Len() > 0 {
				reply.WriteString("\n\n")
			}
			reply.WriteString(t.Text)
		}
	}
	text := strings.TrimSpace(reply.String())
	if text == "" {
		return "", resp.Usage, fmt.Errorf("no text (stop_reason %s)", resp.StopReason)
	}
	return text, resp.Usage, nil
}

// contextOf renders the conversation above the prompt for the reader, who
// needs enough of it to judge whether a reply answers the question.
func contextOf(c Case) string {
	var b strings.Builder
	for _, m := range c.History {
		who := "you"
		if m.Role == "assistant" {
			who = "claude"
		}
		t := m.Text
		if len(t) > 700 {
			t = t[:700] + "\n…"
		}
		fmt.Fprintf(&b, "**%s:** %s\n\n", who, t)
	}
	return strings.TrimSpace(b.String())
}

// humanTyped rejects machinery that arrives in the user role. realPrompt
// already drops the hooks and slash commands the scoring replay cares about;
// a blind test also needs to drop anything a person did not write, because the
// reader is being asked which reply reads better and a reply to a task
// notification is not prose either arm was aimed at.
func humanTyped(t string) bool {
	if strings.HasPrefix(t, "<") {
		return false
	}
	for _, s := range []string{"<task-notification>", "[SYSTEM NOTIFICATION", "<system-reminder>",
		"<command-message>", "Stop hook feedback", "hook feedback:", "PreToolUse", "UserPromptSubmit"} {
		if strings.Contains(t, s) {
			return false
		}
	}
	// A bare path or a bare date is a fragment of something, not a message.
	// Three words is the cheapest test that catches both without guessing at
	// what the fragment came from.
	if len(strings.Fields(t)) < 3 {
		return false
	}
	return true
}

// intactHistory drops a case whose context contains a turn that was all tool
// calls and no prose. Those serialize to an empty text block, which the API
// rejects outright, and the alternative — substituting a placeholder — would
// put text in the context that nobody wrote. There are more cases than the
// test needs, so the picky option is free.
func intactHistory(c Case) bool {
	for _, m := range c.History {
		if strings.TrimSpace(m.Text) == "" {
			return false
		}
	}
	return true
}

// carriesSecret drops a case whose prompt or history quotes a credential.
// These transcripts were written across every project on the machine, the
// pairs file lands on disk unencrypted, and the page renders whatever it is
// handed. Cheap to skip a case; not cheap to un-write a key.
func carriesSecret(c Case) bool {
	markers := []string{"sk-ant-", "ghp_", "github_pat_", "AKIA", "xoxb-", "-----BEGIN", "Bearer ey", "API_KEY=", "_TOKEN=", "_SECRET="}
	check := func(s string) bool {
		for _, m := range markers {
			if strings.Contains(s, m) {
				return true
			}
		}
		return false
	}
	if check(c.Prompt) {
		return true
	}
	for _, m := range c.History {
		if check(m.Text) {
			return true
		}
	}
	return false
}

func hasCard(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	buf := make([]byte, 128<<10)
	var tail string
	for {
		n, err := f.Read(buf)
		if n > 0 {
			chunk := tail + string(buf[:n])
			if strings.Contains(chunk, "<voice-card>") {
				return true, nil
			}
			if len(chunk) > 32 {
				tail = chunk[len(chunk)-32:]
			} else {
				tail = chunk
			}
		}
		if err != nil {
			return false, nil
		}
	}
}

// writePage renders the blind test. The data is base64'd rather than inlined
// as readable JSON so that a glance at the page source does not give away
// which side is which — it is a speed bump for the honest, not a secret.
func writePage(path string, pairs []Pair) error {
	raw, err := json.Marshal(pairs)
	if err != nil {
		return err
	}
	enc := base64.StdEncoding.EncodeToString(raw)
	page := strings.Replace(pageHTML, "__DATA__", enc, 1)
	page = strings.Replace(page, "__COUNT__", fmt.Sprint(len(pairs)), 1)
	return os.WriteFile(path, []byte(page), 0o644)
}

var pageHTML = `<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>cope — blind pairs</title>
<style>
:root { color-scheme: light dark; --fg:#111; --dim:#666; --bg:#fff; --panel:#f6f6f4; --line:#ddd; --accent:#8a5a00; }
@media (prefers-color-scheme: dark) {
  :root { --fg:#e6e4e0; --dim:#8f8d88; --bg:#16161a; --panel:#1e1e23; --line:#33333a; --accent:#d0a050; }
}
* { box-sizing: border-box; }
body { margin:0; background:var(--bg); color:var(--fg); font:16px/1.6 -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif; }
header { position:sticky; top:0; background:var(--bg); border-bottom:1px solid var(--line); padding:.7rem 1.2rem; display:flex; gap:1rem; align-items:center; flex-wrap:wrap; z-index:5; }
header h1 { font-size:.95rem; font-weight:600; margin:0; letter-spacing:.02em; }
#bar { flex:1; height:5px; background:var(--panel); border-radius:3px; overflow:hidden; min-width:120px; }
#fill { height:100%; width:0; background:var(--accent); transition:width .2s; }
#count { font-variant-numeric:tabular-nums; color:var(--dim); font-size:.85rem; }
main { max-width:1400px; margin:0 auto; padding:1.2rem; }
.ctx { background:var(--panel); border:1px solid var(--line); border-radius:8px; padding:.9rem 1.1rem; margin-bottom:1rem; font-size:.9rem; }
.ctx summary { cursor:pointer; color:var(--dim); font-size:.82rem; text-transform:uppercase; letter-spacing:.06em; }
.ctx .body { margin-top:.8rem; color:var(--dim); max-height:340px; overflow:auto; }
.prompt { border-left:3px solid var(--accent); padding:.2rem 0 .2rem 1rem; margin-bottom:1.4rem; }
.prompt .lbl { font-size:.72rem; text-transform:uppercase; letter-spacing:.08em; color:var(--dim); }
.cols { display:grid; grid-template-columns:1fr 1fr; gap:1.2rem; }
@media (max-width:900px) { .cols { grid-template-columns:1fr; } }
.card { border:1px solid var(--line); border-radius:8px; padding:1rem 1.2rem; background:var(--panel); }
.card h2 { font-size:.75rem; text-transform:uppercase; letter-spacing:.1em; color:var(--dim); margin:0 0 .8rem; font-weight:600; }
.card .md { overflow-wrap:anywhere; }
.md p { margin:0 0 .8em; } .md ul,.md ol { margin:0 0 .8em; padding-left:1.3em; } .md li { margin:.2em 0; }
.md h1,.md h2,.md h3 { font-size:1em; font-weight:600; margin:1.2em 0 .4em; }
.md code { background:rgba(128,128,128,.16); padding:.1em .35em; border-radius:3px; font-size:.88em; }
.md pre { background:rgba(128,128,128,.14); padding:.7em .9em; border-radius:6px; overflow-x:auto; }
.md pre code { background:none; padding:0; }
.ask { margin:1.6rem auto .4rem; max-width:900px; }
.ask .q { text-align:center; color:var(--dim); font-size:.82rem; margin-bottom:.5rem; }
.pick { display:flex; gap:.7rem; justify-content:center; flex-wrap:wrap; margin-bottom:1.1rem; }
button { font:inherit; padding:.55rem 1.1rem; border:1px solid var(--line); background:var(--bg); color:var(--fg); border-radius:6px; cursor:pointer; }
button:hover { border-color:var(--accent); }
button.on { background:var(--accent); border-color:var(--accent); color:#fff; }
button kbd { font:inherit; font-size:.78em; opacity:.6; margin-right:.4em; }
.hint { text-align:center; color:var(--dim); font-size:.8rem; }
#results { max-width:760px; margin:0 auto; }
table { border-collapse:collapse; width:100%; margin:1rem 0; font-variant-numeric:tabular-nums; }
th,td { text-align:left; padding:.45rem .7rem; border-bottom:1px solid var(--line); font-size:.92rem; }
th { color:var(--dim); font-weight:600; font-size:.78rem; text-transform:uppercase; letter-spacing:.05em; }
.verdict { border-left:3px solid var(--accent); padding:.6rem 1rem; margin:1.2rem 0; }
</style></head><body>
<header>
  <h1>blind pairs</h1>
  <div id="bar"><div id="fill"></div></div>
  <span id="count"></span>
  <button id="reveal">reveal</button>
</header>
<main id="app"></main>
<script>
var DATA = JSON.parse(decodeURIComponent(escape(atob("__DATA__"))));
var TOTAL = __COUNT__;
var answers = {};
var at = 0;

function esc(s){ var d=document.createElement("div"); d.textContent=s; return d.innerHTML; }

// Enough markdown to read a Claude reply the way the terminal renders one.
function md(src){
  var out=[], lines=src.split("\n"), i=0, para=[], list=null;
  function flushPara(){ if(para.length){ out.push("<p>"+inline(para.join(" "))+"</p>"); para=[]; } }
  function flushList(){ if(list){ out.push("<"+list.tag+">"+list.items.map(function(x){return "<li>"+inline(x)+"</li>";}).join("")+"</"+list.tag+">"); list=null; } }
  function flush(){ flushPara(); flushList(); }
  while(i<lines.length){
    var ln=lines[i];
    if(/^` + "```" + `/.test(ln)){
      flush(); i++; var buf=[];
      while(i<lines.length && !/^` + "```" + `/.test(lines[i])) buf.push(lines[i++]);
      i++; out.push("<pre><code>"+esc(buf.join("\n"))+"</code></pre>"); continue;
    }
    var h=ln.match(/^(#{1,6})\s+(.*)$/);
    if(h){ flush(); out.push("<h3>"+inline(h[2])+"</h3>"); i++; continue; }
    var ul=ln.match(/^\s*[-*+]\s+(.*)$/);
    if(ul){ flushPara(); if(!list||list.tag!=="ul"){ flushList(); list={tag:"ul",items:[]}; } list.items.push(ul[1]); i++; continue; }
    var ol=ln.match(/^\s*\d+[.)]\s+(.*)$/);
    if(ol){ flushPara(); if(!list||list.tag!=="ol"){ flushList(); list={tag:"ol",items:[]}; } list.items.push(ol[1]); i++; continue; }
    if(/^\s*$/.test(ln)){ flush(); i++; continue; }
    flushList(); para.push(ln.trim()); i++;
  }
  flush();
  return out.join("");
}
function inline(s){
  s = esc(s);
  s = s.replace(/` + "`" + `([^` + "`" + `]+)` + "`" + `/g, "<code>$1</code>");
  s = s.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
  s = s.replace(/(^|[^*])\*([^*]+)\*/g, "$1<em>$2</em>");
  return s;
}

function progress(){
  var n=Object.keys(answers).length;
  document.getElementById("fill").style.width=(100*n/TOTAL)+"%";
  document.getElementById("count").textContent=n+" / "+TOTAL;
}

function render(){
  if(at>=DATA.length){ return results(); }
  var p=DATA[at];
  var app=document.getElementById("app");
  app.innerHTML =
    (p.context ? '<details class="ctx"><summary>conversation so far</summary><div class="body md">'+md(p.context)+'</div></details>' : '') +
    '<div class="prompt"><div class="lbl">the message being answered</div><div class="md">'+md(p.prompt)+'</div></div>' +
    '<div class="cols">' +
      '<div class="card"><h2>reply 1</h2><div class="md">'+md(p.left)+'</div></div>' +
      '<div class="card"><h2>reply 2</h2><div class="md">'+md(p.right)+'</div></div>' +
    '</div>' +
    '<div class="ask"><div class="q">which one reads better?</div><div class="pick">' +
      '<button data-q="left"><kbd>1</kbd>reply 1</button>' +
      '<button data-q="tie"><kbd>3</kbd>no difference</button>' +
      '<button data-q="right"><kbd>2</kbd>reply 2</button>' +
    '</div></div>' +
    '<div class="hint">1 / 2 / 3 to answer, → or enter to skip, ← to go back</div>';
  var cur=answers[p.index]||{};
  Array.prototype.forEach.call(app.querySelectorAll("[data-q]"), function(b){
    var v=b.getAttribute("data-q");
    if(cur.q===v) b.className="on";
    b.onclick=function(){ setQ(v); };
  });
  window.scrollTo(0,0);
  progress();
}

function slot(){
  var k=DATA[at].index;
  if(!answers[k]) answers[k]={};
  return answers[k];
}
// A second axis lived here, asking which reply the reader would have to push
// back on. Across 52 pairs it was answered "neither" 52 times: a reply read on
// its own carries no evidence about the judgement behind it, so the question
// had no honest answer and cost half the reading time to not answer. One axis,
// one keypress, and a pair now costs about fifteen seconds.
function setQ(v){ slot().q=v; step(); }
function step(){ at++; render(); }

document.addEventListener("keydown", function(e){
  if(at>=DATA.length) return;
  var q={"1":"left","2":"right","3":"tie"}[e.key];
  if(q) setQ(q);
  else if(e.key==="ArrowRight"||e.key==="Enter") step();
  else if(e.key==="ArrowLeft" && at>0){ at--; render(); }
});
document.getElementById("reveal").onclick=results;

// Two-sided exact binomial against p=0.5, ties dropped.
function sign(k,n){
  if(!n) return 1;
  function C(n,k){ var r=1; for(var i=0;i<k;i++) r=r*(n-i)/(i+1); return r; }
  var lo=Math.min(k,n-k), s=0;
  for(var i=0;i<=lo;i++) s+=C(n,i);
  return Math.min(1, 2*s/Math.pow(2,n));
}

// tally counts how often the card side was chosen on one axis. On "q" that is
// the card winning; on "c" it is the card being the one that needs pushing
// back on, so the same count reads in the opposite direction.
function tally(set, axis){
  var card=0, other=0, tie=0;
  DATA.forEach(function(p){
    if(p.set!==set) return;
    var a=answers[p.index]; if(!a) return;
    var v=a[axis]; if(!v) return;
    if(v==="tie"||v==="neither"||v==="both"){ tie++; return; }
    var picked = v==="left" ? p.left_cond : p.right_cond;
    if(picked==="card") card++; else other++;
  });
  return {card:card, other:other, tie:tie};
}

// The comparator for a set is read out of the data rather than written into
// the page, so a run against someone else's card and someone else's baseline
// labels its own results. Nothing here knows what the two documents were.
function against(set){
  for(var i=0;i<DATA.length;i++){
    if(DATA[i].set!==set) continue;
    return DATA[i].left_cond==="card" ? DATA[i].right_cond : DATA[i].left_cond;
  }
  return "";
}
function has(set){ return DATA.some(function(p){ return p.set===set; }); }

function row(label, t, vs){
  var n=t.card+t.other;
  var p=sign(t.card,n);
  return "<tr><td>"+label+"</td><td>"+t.card+" / "+n+"</td><td>"+vs+" "+t.other+"</td><td>"+t.tie+"</td><td>"+
    (n ? (100*t.card/n).toFixed(0)+"%" : "—")+"</td><td>"+(n ? p.toFixed(3) : "—")+"</td></tr>";
}

function table(caption, axis, cardCol, tieCol){
  var rows="";
  ["test","control"].forEach(function(set){
    if(!has(set)) return;
    var a=against(set);
    rows += row(set+" — card vs "+a, tally(set,axis), a);
  });
  return "<h3>"+caption+"</h3><table><tr><th>set</th><th>"+cardCol+"</th><th>vs</th><th>"+tieCol+
    "</th><th>rate</th><th>p</th></tr>" + rows + "</table>";
}

// A run may carry no control at all. The control's job is to show the axis can
// register a hit; once that is established, pairs spent on it cost more power
// on the question than the second reading buys back.
function verdictFor(axis){
  var test=tally("test",axis), ctrl=tally("control",axis);
  var nT=test.card+test.other, nC=ctrl.card+ctrl.other;
  var vs=against("test")||"the baseline";
  if(nT<6 && nC<6) return "Too few calls to read.";
  var out="";
  if(nT>=6){
    var pT=sign(test.card,nT), rate=(100*test.card/nT).toFixed(0)+"%";
    if(pT<0.05){
      out = test.card>test.other
        ? "The card wins at "+rate+", p="+pT.toFixed(3)+". That is a real difference against "+vs+" on its own."
        : "The card LOSES at "+rate+", p="+pT.toFixed(3)+". It is making the prose worse than "+vs+" alone.";
    } else {
      out = "No separation: "+rate+", p="+pT.toFixed(3)+". At this n only a large move is readable, so this rules out a large effect and says nothing about a small one.";
    }
  }
  if(nC>=6){
    var pC=sign(ctrl.card,nC), cr=(100*ctrl.card/nC).toFixed(0)+"%";
    out += (out?" ":"") + (pC<0.05
      ? "The control separates at "+cr+", n="+nC+", so the axis can register a hit."
      : "The control is "+cr+" at n="+nC+" and does not separate. Until it does, a flat test result is as likely to be a blind instrument as a null.");
  }
  return out;
}

function results(){
  document.getElementById("app").innerHTML =
    '<div id="results"><h2>results</h2>' +
    table("which reads better", "q", "card picked", "no difference") +
    '<div class="verdict">'+verdictFor("q")+'</div>' +
    '<p style="color:var(--dim);font-size:.88rem">p is a two-sided exact binomial against a coin flip, ' +
    'ties dropped. Reload to start over; nothing is saved.</p>' +
    '<p><button id="back">back to the pairs</button></p></div>';
  // Reveal is meant to be looked at mid-run, so it has to be leavable. Without
  // this the only way back is the arrow keys, and reloading to escape costs
  // every answer given so far.
  document.getElementById("back").onclick=function(){
    if(at>=DATA.length) at=DATA.length-1;
    render();
  };
  progress();
}

render();
</script></body></html>
`
