// cope-replay generates prose under each arm of the experiment, offline.
//
// The live crossover can only consume prose as fast as it gets written — about
// 165 hits a day. This replays real prompts out of past transcripts, calls the
// API once per arm with the reminder swapped, and scores both replies with the
// same rules the hook uses. Same question, hours instead of days.
//
// What it cannot do: the replies come from a reconstructed context with no real
// tool use, so if the card only bites during real work this will miss it. It
// complements the live experiment rather than replacing it.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

const model = "claude-opus-5"

// The arms, in the order a case runs them. Held first so a failure part-way
// through still leaves a usable reference arm.
var arms = []string{"hold", "inject", "positive"}

// framing tells the model what setting it is writing in. Without it the
// replayed context reads as a transcript to comment on rather than a
// conversation to continue, and the prose comes out in the wrong register.
const framing = `You are Claude Code, Anthropic's official CLI for Claude, in an interactive session with a software engineer. What follows is the conversation so far. Reply to the final message as you normally would.`

func main() {
	var (
		transcript = flag.String("transcript", "", "Claude Code session JSONL to replay")
		n          = flag.Int("n", 8, "how many prompts to sample")
		historyN   = flag.Int("history", 6, "messages of prior conversation to carry")
		minReply   = flag.Int("min-reply", 400, "skip prompts whose real reply was shorter than this")
		out        = flag.String("out", "replay-turns.jsonl", "where to write the scored turns")
		gate       = flag.String("gate", "cope-gate", "path to the cope-gate binary")
		maxTok     = flag.Int64("max-tokens", 8000, "output ceiling; thinking counts against it")
		dryRun     = flag.Bool("dry-run", false, "assemble every request and report size, without calling")
		conc       = flag.Int("concurrency", 3, "requests in flight")

		pairsMode   = flag.Bool("pairs", false, "generate a blind pairwise preference test instead of scoring")
		transcripts = flag.String("transcripts", "", "glob of transcripts to draw cases from (-pairs)")
		nTest       = flag.Int("test-pairs", 20, "pairs of card vs CLAUDE.md (-pairs)")
		nControl    = flag.Int("control-pairs", 10, "pairs of card vs bare framing, the positive control (-pairs)")
		seed        = flag.Int("seed", 1, "shuffle seed, so a run is reproducible (-pairs)")
		maxFiles    = flag.Int("max-files", 40, "transcripts to read cases from, after shuffling (-pairs)")
		outDir      = flag.String("outdir", "pairs", "where to write pairs.html and pairs.json (-pairs)")
		repageFrom  = flag.String("repage", "", "rebuild pairs.html from an existing pairs.json, without calling the API")
		appendTo    = flag.Bool("append", false, "merge this run into the pairs.json already in -outdir (-pairs)")
		exclude     = flag.String("exclude", "", "comma-separated pairs.json paths whose prompts this run must not re-sample (-pairs)")
		rules       = flag.String("rules", "", "card JSON to test instead of the one built into the gate (-pairs)")
		baseline    = flag.String("baseline", "", "file of prose guidance the card is tested on top of, e.g. the Response style section of a CLAUDE.md (-pairs)")
		rival       = flag.String("rival", "", "second card JSON to test the card head to head against, instead of stacking it on -baseline (-pairs)")

		discriminate = flag.Bool("discriminate", false, "ask which reply carries a named voice instead of which reads better (-pairs)")
		judge        = flag.Bool("judge", false, "answer every pair with the model too, on the same materials, so its agreement with you can be read (-discriminate)")
	)
	flag.Parse()

	if *repageFrom != "" {
		if err := repage(*repageFrom, *seed); err != nil {
			fmt.Fprintf(os.Stderr, "cope-replay: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *pairsMode {
		if *transcripts == "" {
			fmt.Fprintln(os.Stderr, "cope-replay: -pairs needs -transcripts")
			os.Exit(2)
		}
		if *judge && !*discriminate {
			fmt.Fprintln(os.Stderr, "cope-replay: -judge needs -discriminate; there is no key to score a preference against")
			os.Exit(2)
		}
		if err := runPairs(*gate, *transcripts, *outDir, *exclude, *rules, *baseline, *rival, *nTest, *nControl, *historyN, *seed, *maxFiles, *maxTok, *conc, *dryRun, *appendTo, *discriminate, *judge); err != nil {
			fmt.Fprintf(os.Stderr, "cope-replay: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *transcript == "" {
		fmt.Fprintln(os.Stderr, "cope-replay: -transcript is required")
		os.Exit(2)
	}

	cases, err := readCases(*transcript, *historyN, *minReply)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-replay: %v\n", err)
		os.Exit(1)
	}
	if len(cases) == 0 {
		fmt.Fprintln(os.Stderr, "cope-replay: no replayable prompts in that transcript")
		os.Exit(1)
	}
	cases = sample(cases, *n)

	card, err := run(*gate, "--inject")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-replay: reading the card: %v\n", err)
		os.Exit(1)
	}
	reminders, err := loadReminders(*gate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-replay: rendering the arms: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d prompts x %d arms = %d calls, model %s\n", len(cases), len(arms), len(cases)*len(arms), model)

	if *dryRun {
		reportDryRun(cases, card, reminders)
		return
	}

	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-replay: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	client := anthropic.NewClient()
	var (
		mu    sync.Mutex
		enc   = json.NewEncoder(f)
		usage totals
		sem   = make(chan struct{}, *conc)
		wg    sync.WaitGroup
	)

	for i, c := range cases {
		for _, arm := range arms {
			wg.Add(1)
			sem <- struct{}{}
			go func(i int, c Case, arm string) {
				defer wg.Done()
				defer func() { <-sem }()

				rec, u, err := replayOne(client, c, arm, card, reminders[arm], *maxTok, *gate, i)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					fmt.Fprintf(os.Stderr, "  case %d %s: %v\n", i, arm, err)
					return
				}
				usage.add(u)
				_ = enc.Encode(rec)
				fmt.Printf("  case %-3d %-9s %5d ch  %v\n", i, arm, rec.Chars, rec.Rules)
			}(i, c, arm)
		}
	}
	wg.Wait()

	fmt.Printf("\nwrote %s\n", filepath.Clean(*out))
	usage.report()
	fmt.Printf("\nread it with: %s --ab-report %s\n", *gate, filepath.Clean(*out))
	fmt.Println("ignore that report's drift column. Every call here is independent," +
		"\nso there is no session for a voice to drift across — replay can measure" +
		"\nthe per-turn rate and nothing about the trajectory.")
}

// turnRecord matches what the hook writes to turns.jsonl, so --ab-report reads
// a replay and a live run through exactly the same code.
type turnRecord struct {
	Time        string   `json:"time"`
	Session     string   `json:"session"`
	Key         string   `json:"key"`
	Arm         string   `json:"arm"`
	Chars       int      `json:"chars"`
	Rules       []string `json:"rules"`
	InjectChars int      `json:"inject_chars,omitempty"`
}

func replayOne(client anthropic.Client, c Case, arm, card, reminder string, maxTok int64, gate string, i int) (turnRecord, anthropic.Usage, error) {
	msgs := make([]anthropic.MessageParam, 0, len(c.History)+1)
	for j, m := range c.History {
		block := anthropic.NewTextBlock(m.Text)
		// The history is byte-identical across the three arms of one case, so
		// marking its last block means arms two and three read it back at a
		// tenth of the price instead of paying for it again.
		if j == len(c.History)-1 {
			block.OfText.CacheControl = anthropic.NewCacheControlEphemeralParam()
		}
		if m.Role == "assistant" {
			msgs = append(msgs, anthropic.NewAssistantMessage(block))
		} else {
			msgs = append(msgs, anthropic.NewUserMessage(block))
		}
	}
	// The hook prepends its reminder to the prompt, so the replay does too.
	final := c.Prompt
	if reminder != "" {
		final = reminder + "\n" + c.Prompt
	}
	msgs = append(msgs, anthropic.NewUserMessage(anthropic.NewTextBlock(final)))

	resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     model,
		MaxTokens: maxTok,
		System: []anthropic.TextBlockParam{
			{Text: framing},
			// The card is identical in every call of the whole run, so this is
			// the one breakpoint that pays off across cases as well as arms.
			{Text: card, CacheControl: anthropic.NewCacheControlEphemeralParam()},
		},
		Messages: msgs,
	})
	if err != nil {
		return turnRecord{}, anthropic.Usage{}, err
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
		return turnRecord{}, resp.Usage, fmt.Errorf("no text in reply (stop_reason %s)", resp.StopReason)
	}

	rules, err := score(gate, text)
	if err != nil {
		return turnRecord{}, resp.Usage, err
	}
	return turnRecord{
		Time:    time.Now().UTC().Format(time.RFC3339),
		Session: "replay-" + filepath.Base(strings.TrimSuffix(model, ".")),
		// The arm belongs in the key. A live turn is one reply scored twice and
		// deduped on this; a replayed case is three different replies, and
		// keying them alike would collapse the run to its last arm.
		Key:         fmt.Sprintf("case-%03d-%s", i, arm),
		Arm:         arm,
		Chars:       len(text),
		Rules:       rules,
		InjectChars: len(reminder),
	}, resp.Usage, nil
}

// score runs the reply through the same binary the hook uses, so the replay
// cannot drift from the live scoring. --check writes one JSON record per
// violation to the log, which is the machine-readable surface.
func score(gate, text string) ([]string, error) {
	tmp, err := os.CreateTemp("", "cope-replay-*.jsonl")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	cmd := exec.Command(gate, "--check", "-", "--log", tmp.Name())
	cmd.Stdin = strings.NewReader(text)
	cmd.Stdout = nil
	if err := cmd.Run(); err != nil {
		// --check exits non-zero only under --block, which is off here.
		return nil, fmt.Errorf("scoring: %w", err)
	}

	f, err := os.Open(tmp.Name())
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rules := []string{}
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 64<<10), 4<<20)
	for s.Scan() {
		var v struct {
			RuleID string `json:"rule_id"`
		}
		if json.Unmarshal(s.Bytes(), &v) == nil && v.RuleID != "" {
			rules = append(rules, v.RuleID)
		}
	}
	return rules, s.Err()
}

// loadReminders renders each arm once. The rules named are the two that
// dominate every measured session, which is what the live refresher would be
// speaking to; hold renders nothing by definition.
func loadReminders(gate string) (map[string]string, error) {
	out := map[string]string{"hold": ""}
	for _, arm := range []string{"inject", "positive"} {
		t, err := run(gate, "--render-arm", arm, "--render-for", "flip,labelled_opening")
		if err != nil {
			return nil, err
		}
		out[arm] = t
	}
	return out, nil
}

func run(bin string, args ...string) (string, error) {
	out, err := exec.Command(bin, args...).Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", bin, strings.Join(args, " "), err)
	}
	return string(out), nil
}

// sample spreads the picks evenly through the session rather than taking the
// first n, so a replay is not all opening-of-conversation prompts.
func sample(cases []Case, n int) []Case {
	if n <= 0 || n >= len(cases) {
		return cases
	}
	out := make([]Case, 0, n)
	step := float64(len(cases)) / float64(n)
	for i := range n {
		out = append(out, cases[int(float64(i)*step)])
	}
	return out
}

type totals struct {
	in, out, cacheWrite, cacheRead int64
}

func (t *totals) add(u anthropic.Usage) {
	t.in += u.InputTokens
	t.out += u.OutputTokens
	t.cacheWrite += u.CacheCreationInputTokens
	t.cacheRead += u.CacheReadInputTokens
}

// report prices the run at Opus 5's published rates: $5 per million input,
// $25 per million output, cache reads at a tenth of input and 5-minute cache
// writes at 1.25x.
func (t *totals) report() {
	const inRate, outRate = 5.0 / 1e6, 25.0 / 1e6
	cost := float64(t.in)*inRate +
		float64(t.cacheWrite)*inRate*1.25 +
		float64(t.cacheRead)*inRate*0.1 +
		float64(t.out)*outRate
	fmt.Printf("tokens: %d in, %d cache-write, %d cache-read, %d out\n",
		t.in, t.cacheWrite, t.cacheRead, t.out)
	if t.cacheRead+t.cacheWrite > 0 {
		fmt.Printf("cache hit rate: %.0f%% of cacheable input\n",
			100*float64(t.cacheRead)/float64(t.cacheRead+t.cacheWrite))
	}
	fmt.Printf("cost: $%.2f at $5/$25 per Mtok\n", cost)
}

// reportDryRun assembles what would be sent and estimates it, so a run can be
// priced before it is paid for. Characters over four is a rough token count;
// it is an order-of-magnitude check, not a bill.
func reportDryRun(cases []Case, card string, reminders map[string]string) {
	var shared, variable int
	for _, c := range cases {
		for _, m := range c.History {
			shared += len(m.Text)
		}
		for _, arm := range arms {
			variable += len(c.Prompt) + len(reminders[arm])
		}
	}
	cardTok := len(card) / 4
	fmt.Printf("\ncard: %d ch (~%d tok), cached once and read by every call\n", len(card), cardTok)
	for _, arm := range arms {
		fmt.Printf("  %-9s reminder %d ch\n", arm, len(reminders[arm]))
	}
	fmt.Printf("history across cases: %d ch (~%d tok), each cached once and read %d times\n",
		shared, shared/4, len(arms)-1)
	fmt.Printf("variable tail: %d ch (~%d tok)\n", variable, variable/4)
	fmt.Printf("\nrough input: ~%d tok billed, plus ~%d tok of cache reads at a tenth\n",
		cardTok+shared/4+variable/4, (len(arms)-1)*(cardTok+shared/4))
	fmt.Println("output is the unknown — budget ~600 tok of prose per call, plus thinking.")
}
