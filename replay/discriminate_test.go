package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func discPairs(n int, set string) []Pair {
	out := make([]Pair, n)
	for i := range out {
		other := "rival"
		if set == "control" {
			other = "bare"
		}
		out[i] = Pair{
			Index: i, Set: set,
			Left: "a", LeftCond: "card",
			Right: "b", RightCond: other,
		}
	}
	return out
}

// The named voice alternates so a reader cannot answer everything by
// recognising one card and calling its absence, and so both cards are measured
// rather than one being the permanent subject and the other the permanent foil.
func TestTargetsAlternateAcrossTheTestSet(t *testing.T) {
	pairs := discPairs(20, "test")
	descs := map[string]string{"card": "one voice", "rival": "another voice"}
	balanceTargets(pairs, descs, rand.New(rand.NewPCG(7, 0xc09e)))

	count := map[string]int{}
	for _, p := range pairs {
		count[p.Target]++
	}
	if count["card"] != 10 || count["rival"] != 10 {
		t.Errorf("targets split %v, want 10 each", count)
	}
}

// bare is the framing alone and describes no voice, so it can never be the
// thing a reader is asked to find. That is what makes the control set a control
// rather than a second test.
func TestControlAlwaysAsksAboutTheCard(t *testing.T) {
	pairs := discPairs(10, "control")
	balanceTargets(pairs, map[string]string{"card": "one voice"}, rand.New(rand.NewPCG(3, 0xc09e)))
	for _, p := range pairs {
		if p.Target != "card" {
			t.Fatalf("control pair targets %q, and %q has no voice to recognise", p.Target, p.Target)
		}
	}
}

// A baseline run has a second arm with no VOICE block, which is the same shape
// as the control: describable on one side only.
func TestATargetIsNeverAVoicelessArm(t *testing.T) {
	pairs := discPairs(12, "test")
	for i := range pairs {
		pairs[i].RightCond = "rival" // a baseline arm keys the same way
	}
	// No description for the rival: -baseline supplies prose guidance, not a card.
	balanceTargets(pairs, map[string]string{"card": "one voice"}, rand.New(rand.NewPCG(11, 0xc09e)))
	for _, p := range pairs {
		if p.Target != "card" {
			t.Errorf("asked the reader to find %q, which has no description", p.Target)
		}
	}
}

// The answer on this page is the target's side, and the target is the rival in
// half the pairs — so balancing where the card sits leaves the answer free to
// pile up on one side. It did, five of six, on the first live run.
func TestTheAnswerIsBalancedAcrossSides(t *testing.T) {
	descs := map[string]string{"card": "one voice", "rival": "another voice"}
	for _, seed := range []uint64{1, 5, 9, 23} {
		pairs := discPairs(20, "test")
		balanceTargets(pairs, descs, rand.New(rand.NewPCG(seed, 0xc09e)))
		left := 0
		for _, p := range pairs {
			if p.correct() == "left" {
				left++
			}
		}
		if left != 10 {
			t.Errorf("seed %d: the answer is on the left in %d of 20 pairs, want 10", seed, left)
		}
	}
}

// Swapping sides moves the prose, not the label. A pair whose sides were
// swapped to balance the answer must still carry the reply that condition
// actually wrote.
func TestSideSwapsKeepEachReplyWithItsCondition(t *testing.T) {
	pairs := discPairs(8, "test")
	for i := range pairs {
		pairs[i].Left, pairs[i].Right = "written by card", "written by rival"
	}
	balanceTargets(pairs, map[string]string{"card": "a", "rival": "b"}, rand.New(rand.NewPCG(4, 0xc09e)))
	for _, p := range pairs {
		if (p.LeftCond == "card") != (p.Left == "written by card") {
			t.Fatalf("pair %d: %s is labelled %q on the left", p.Index, p.Left, p.LeftCond)
		}
		if (p.RightCond == "card") != (p.Right == "written by card") {
			t.Fatalf("pair %d: %s is labelled %q on the right", p.Index, p.Right, p.RightCond)
		}
	}
}

// The key. correct() is what both the page and the judge score against, so a
// swap here would invert every result silently.
func TestCorrectFollowsTheTargetToItsSide(t *testing.T) {
	left := Pair{LeftCond: "card", RightCond: "rival", Target: "card"}
	if got := left.correct(); got != "left" {
		t.Errorf("card on the left, target card: got %q", got)
	}
	right := Pair{LeftCond: "card", RightCond: "rival", Target: "rival"}
	if got := right.correct(); got != "right" {
		t.Errorf("rival on the right, target rival: got %q", got)
	}
}

// A malformed answer is missing data. Scoring it wrong would understate a judge
// that was only disobedient about format.
func TestParseVerdictReadsTheWordAndNothingElse(t *testing.T) {
	for in, want := range map[string]string{
		"one": "left", "One.": "left", " 1 ": "left", "reply one": "left",
		"two": "right", "TWO": "right", "2": "right", "reply 2": "right",
		"unsure": "unsure", "": "unsure",
		"I think the first one, because its sentences are shorter": "unsure",
	} {
		if got := parseVerdict(in); got != want {
			t.Errorf("parseVerdict(%q) = %q, want %q", in, got, want)
		}
	}
}

var atobArg = regexp.MustCompile(`atob\("([A-Za-z0-9+/=]+)"\)`)

func decodeBlock(t *testing.T, page string, n int, into any) {
	t.Helper()
	m := atobArg.FindAllStringSubmatch(page, -1)
	if len(m) <= n {
		t.Fatalf("page carries %d embedded blocks, wanted at least %d", len(m), n+1)
	}
	raw, err := base64.StdEncoding.DecodeString(m[n][1])
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatal(err)
	}
}

// The page shows a voice and asks which reply carries it. A pair whose target
// has no description renders an empty panel and the reader answers a question
// nobody asked — so the pairing of target to description is checked here rather
// than discovered mid-run.
func TestEveryPairOnThePageHasItsVoice(t *testing.T) {
	pairs := discPairs(6, "test")
	descs := map[string]string{
		"card":  "Aim: the phrasing should cost the reader nothing",
		"rival": "Aim: say the fewest words that answer",
	}
	balanceTargets(pairs, descs, rand.New(rand.NewPCG(2, 0xc09e)))

	path := filepath.Join(t.TempDir(), "discriminate.html")
	if err := writeDiscPage(path, pairs, descs); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	page := string(raw)

	var gotPairs []Pair
	var gotDescs map[string]string
	decodeBlock(t, page, 0, &gotPairs)
	decodeBlock(t, page, 1, &gotDescs)

	if len(gotPairs) != len(pairs) {
		t.Fatalf("page carries %d pairs, want %d", len(gotPairs), len(pairs))
	}
	for _, p := range gotPairs {
		if gotDescs[p.Target] == "" {
			t.Errorf("pair %d asks about %q, which the page cannot describe", p.Index, p.Target)
		}
	}

	// The key is the one thing a reader must not read off the page by accident.
	// Base64 is a speed bump rather than a secret, but the serialised form must
	// not also be sitting in the markup where a scroll would land on it.
	for _, leak := range []string{`"left_cond":`, `"target":`, `"judge":`} {
		if strings.Contains(page, leak) {
			t.Errorf("the page inlines the answer key as readable JSON: %s", leak)
		}
	}
}
