package main

// Blind discrimination: did the voicing arrive at all.
//
// The preference instrument next door asks which of two replies reads better.
// That question cannot answer "make it sound like this instead", and reading it
// as though it could is how a null on n=11 became a claim that voice guidance
// does nothing. Preference and identity are different axes: a voice can land
// completely and leave the reader indifferent about whether it was worth
// landing.
//
// So this asks the other question. The reader is shown one voice — the aim and
// the register, taken out of the card itself — and two replies to the same real
// prompt, one written under that voice and one not. Which is which?
//
//	above chance   the voice is on the page and a reader can find it
//	at chance      whatever the card did, it did not reach the prose
//
// Chance is 50%, so the arithmetic is the same exact binomial the preference
// page uses. What differs is that a null here is a much stronger statement: two
// replies can be equally good in a hundred ways, but if they cannot be told
// apart under a stated target then no voicing happened, whatever the card said.
//
// Two sets, same structure as the preference run:
//
//	control  card vs bare framing — the widest gap there is, and if the reader
//	         cannot pick the card out against no guidance at all, the
//	         instrument is blind and nothing else it says means anything
//	test     card vs a second card — two deliberate voices told apart, which is
//	         the question a framework for arbitrary voicing has to answer
//
// The test set alternates which of the two cards is named as the target, so a
// reader cannot settle into recognising one voice and answering everything by
// its absence, and so both cards get measured rather than one.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

// balanceTargets names the voice each pair asks about — half to each card
// within the test set, always the card in the control set — and then puts that
// voice on the left in half of them.
//
// Both halves are balanced rather than flipped, for the same reason sides are
// in the preference run: at thirty pairs a coin lands lopsided often enough to
// matter, and any asymmetry would ride on it.
//
// The second half is the one that is easy to leave out, and leaving it out is
// not harmless. balanceSides balances where the *card* sits, which is the
// answer on a preference run; here the answer is where the *target* sits, and
// the target is the rival in half the pairs. Assigning targets over
// card-balanced sides put the answer on the right in five pairs of six on the
// first live run — free marks for any reader with a side habit, in the one
// direction blinding cannot cover.
//
// The control set has no choice of target to make. bare is the framing alone
// and describes no voice, so the only nameable target is the card — which is
// what makes that set a positive control rather than a second test.
func balanceTargets(pairs []Pair, descs map[string]string, rng *rand.Rand) {
	bySet := map[string][]int{}
	for i, p := range pairs {
		bySet[p.Set] = append(bySet[p.Set], i)
	}
	for set, idx := range bySet {
		rng.Shuffle(len(idx), func(a, b int) { idx[a], idx[b] = idx[b], idx[a] })
		for k, i := range idx {
			p := &pairs[i]
			other := p.LeftCond
			if other == "card" {
				other = p.RightCond
			}
			// Half the test set asks about the rival instead, when there is a
			// rival with a voice to ask about.
			p.Target = "card"
			if set == "test" && descs[other] != "" && k >= len(idx)/2 {
				p.Target = other
			}
		}
		// Re-shuffle before the side pass so which pairs get the target on the
		// left is independent of which ones were given the rival as target.
		rng.Shuffle(len(idx), func(a, b int) { idx[a], idx[b] = idx[b], idx[a] })
		for k, i := range idx {
			p := &pairs[i]
			if wantLeft := k < len(idx)/2; (p.LeftCond == p.Target) != wantLeft {
				p.Left, p.Right = p.Right, p.Left
				p.LeftCond, p.RightCond = p.RightCond, p.LeftCond
			}
		}
	}
}

// correct reports which side the pair's target was actually written under.
func (p Pair) correct() string {
	if p.LeftCond == p.Target {
		return "left"
	}
	return "right"
}

// The instruction may not say the other reply had no voice. On a test-set pair
// both were written under a card and only one was written under this one, so
// "the other was not" points the judge at whichever reply reads as more voiced
// — toward a loud rival and away from a quiet card. It survived the first two
// runs because the bias ran against the shipped card and it scored anyway, but
// a number is not worth having twice from a question asked wrong.
const judgeInstruction = `You will be shown a description of a writing voice, then two replies to the same message. One of them was written by a model working under that voice description. The other was written under different instructions, which may include a different voice.

Say which one. Answer with exactly one word: "one", "two", or "unsure".

Judge only how the reply is written — its register, sentence shapes, how it opens and closes, what it does with detail. Do not judge which reply is better, more helpful, or more correct. Answer "unsure" only if the two are genuinely indistinguishable in register; a weak signal is still a signal and should be called.`

// runJudge answers every pair with the model, on the same materials and the
// same question the reader gets.
//
// It exists because reading is the bottleneck. A person can call thirty pairs
// in an evening, and thirty pairs resolves a large effect and nothing else — so
// every card revision costs another evening, and a framework meant to carry
// arbitrary voices needs the loop to be shorter than that.
//
// The judge is not trusted on its own. It answers the pairs the reader also
// answers, so the page can report agreement, and until that agreement has been
// looked at the judge's rate is a number about a model rather than about a
// voice. That check is the whole reason it runs on the same pairs instead of a
// cheaper separate sample.
func runJudge(pairs []Pair, descs map[string]string, conc int, usage *totals) {
	client := anthropic.NewClient()
	var (
		mu  sync.Mutex
		sem = make(chan struct{}, conc)
		wg  sync.WaitGroup
	)
	fmt.Printf("\njudging %d pairs\n", len(pairs))
	for i := range pairs {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			verdict, u, err := judgeOne(client, pairs[i], descs[pairs[i].Target])
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fmt.Fprintf(os.Stderr, "  judge %d: %v\n", i, err)
				return
			}
			usage.add(u)
			pairs[i].Judge = verdict
		}(i)
	}
	wg.Wait()

	hit, n := 0, 0
	for _, p := range pairs {
		if p.Judge == "" || p.Judge == "unsure" {
			continue
		}
		n++
		if p.Judge == p.correct() {
			hit++
		}
	}
	if n > 0 {
		fmt.Printf("judge: %d/%d correct (%.0f%%)\n", hit, n, 100*float64(hit)/float64(n))
	}
}

func judgeOne(client anthropic.Client, p Pair, desc string) (string, anthropic.Usage, error) {
	if desc == "" {
		return "", anthropic.Usage{}, fmt.Errorf("no description for target %q", p.Target)
	}
	// The instruction and the target description are identical across every
	// pair asking about the same voice, and they are the long part. Marking the
	// end of the description caches everything above the two replies, so a
	// hundred-pair run pays for this prefix once.
	system := []anthropic.TextBlockParam{
		{Text: judgeInstruction},
		{Text: "The voice:\n\n" + desc, CacheControl: anthropic.NewCacheControlEphemeralParam()},
	}
	ask := fmt.Sprintf("The message both replies are answering:\n\n%s\n\n"+
		"--- reply one ---\n\n%s\n\n--- reply two ---\n\n%s\n\n"+
		"Which reply was written under the voice? One word: one, two, or unsure.",
		p.Prompt, p.Left, p.Right)

	var (
		resp *anthropic.Message
		err  error
	)
	for attempt := 0; ; attempt++ {
		resp, err = client.Messages.New(context.Background(), anthropic.MessageNewParams{
			Model: model,
			// One word of answer needs almost none of this. The room is for the
			// thinking block Opus 5 emits ahead of the text whether or not it
			// was asked for: at a ceiling sized for the answer alone, thinking
			// spends the whole budget and no text block arrives at all. Six
			// pairs came back "unsure" that way, which reads as a judge that
			// cannot tell two voices apart rather than as a truncated reply.
			MaxTokens: 1024,
			System:    system,
			Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(ask))},
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
	var said strings.Builder
	for _, b := range resp.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			said.WriteString(t.Text)
		}
	}
	// No text block is missing data, not an answer of "unsure". The page counts
	// the two differently — an empty verdict is dropped from the run, an unsure
	// is a call the judge made and declined to resolve — and collapsing them
	// would turn a truncated reply into a finding about the voices.
	if strings.TrimSpace(said.String()) == "" {
		return "", resp.Usage, fmt.Errorf("no answer (stop_reason %s)", resp.StopReason)
	}
	return parseVerdict(said.String()), resp.Usage, nil
}

// parseVerdict reads the judge's word. Anything it cannot place counts as
// unsure and is dropped from the rate rather than guessed at — a malformed
// answer is missing data, and scoring it as wrong would understate a judge that
// was merely disobedient about format.
func parseVerdict(s string) string {
	switch strings.ToLower(strings.Trim(strings.TrimSpace(s), ".\"'`")) {
	case "one", "1", "reply one", "reply 1":
		return "left"
	case "two", "2", "reply two", "reply 2":
		return "right"
	default:
		return "unsure"
	}
}

// writeDiscPage renders the discrimination test. Like the preference page the
// data is base64'd, which hides the key from a glance at the source without
// pretending to hide it from anyone who wants it.
func writeDiscPage(path string, pairs []Pair, descs map[string]string) error {
	raw, err := json.Marshal(pairs)
	if err != nil {
		return err
	}
	dsc, err := json.Marshal(descs)
	if err != nil {
		return err
	}
	page := strings.Replace(discHTML, "__DATA__", base64.StdEncoding.EncodeToString(raw), 1)
	page = strings.Replace(page, "__DESCS__", base64.StdEncoding.EncodeToString(dsc), 1)
	page = strings.Replace(page, "__COUNT__", fmt.Sprint(len(pairs)), 1)
	return os.WriteFile(path, []byte(page), 0o644)
}
