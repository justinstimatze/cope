package scan

import (
	"strings"
	"testing"
)

func phraseTexts(ps []Phrase) []string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Text)
	}
	return out
}

func hasPhrase(ps []Phrase, want string) bool {
	for _, p := range ps {
		if p.Text == want {
			return true
		}
	}
	return false
}

// The whole rule turns on this distinction. A construction is made of
// connective tissue and recurs because the writer reaches for it; a subject
// recurs because the conversation is about it, and flagging that would make the
// rule fire hardest on the most focused sessions.
func TestPhrasesReadConstructionsAndNotSubjectMatter(t *testing.T) {
	construction := Phrases(
		"That is not the whole point of it, and the whole point of it is " +
			"something else again once you have looked at what came before.")
	if !hasPhrase(construction, "that is not the whole") {
		t.Errorf("missed a construction made entirely of common words: %q", phraseTexts(construction))
	}

	// One content word is still a construction: it is the frame that repeats,
	// with the topic riding along inside it.
	oneContent := Phrases(
		"The whole point of cope is that a linter reads one reply and this " +
			"reads the session, which is a different question altogether.")
	if !hasPhrase(oneContent, "the whole point of cope") {
		t.Errorf("a four-scaffold shingle should count: %q", phraseTexts(oneContent))
	}

	subject := Phrases(
		"The cross turn repetition detector fingerprints five word shingles " +
			"drawn from redacted prose paragraphs, then discards duplicates.")
	if len(subject) != 0 {
		t.Errorf("subject-matter spans were fingerprinted, so a focused session "+
			"would report itself for staying on topic: %q", phraseTexts(subject))
	}
}

// A reply that says something twice is one turn saying it. Counting it twice
// would let a single paragraph clear the threshold on its own.
func TestPhrasesAreDedupedWithinATurn(t *testing.T) {
	ps := Phrases(
		"That is not the whole point at all here. Something intervenes and " +
			"then, later on, that is not the whole point either.")
	n := 0
	for _, p := range ps {
		if p.Text == "that is not the whole" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("the same span appears %d times in one turn's phrases, want 1", n)
	}
}

// Fingerprints are what reaches disk, so equal text must give equal sums and
// different text must not. Without the first half the window never matches;
// without the second it matches everything.
func TestFingerprintIgnoresCaseAndNotWords(t *testing.T) {
	a := Phrases("That Is Not The Whole point of it, said with the capitals left in place.")
	b := Phrases("that is not the whole point of it, said with the capitals taken out.")
	if len(a) == 0 || len(b) == 0 {
		t.Fatal("no phrases to compare")
	}
	if a[0].Sum != b[0].Sum {
		t.Error("case changed the fingerprint, so a sentence-initial tic never matches itself")
	}
	c := Phrases("That is not the point of it, and one word is all that differs here.")
	if len(c) > 0 && c[0].Sum == a[0].Sum {
		t.Error("different spans share a fingerprint")
	}
}

// Two earlier turns, not one. A phrase used twice in a long session is a
// coincidence; the complaint is about the third and fourth airing.
func TestCrossTurnRepeatWaitsForTheThirdAiring(t *testing.T) {
	text := "That is not the whole point of it, and there is more to say " +
		"about what the earlier turns were already doing here."
	ps := Phrases(text)
	if len(ps) == 0 {
		t.Fatal("no phrases extracted")
	}
	sum := ps[0].Sum

	if v := CrossTurnRepeat(ps, map[uint32]int{sum: 1}); len(v) != 0 {
		t.Errorf("fired on one earlier turn: %+v", v)
	}
	v := CrossTurnRepeat(ps, map[uint32]int{sum: 2})
	if len(v) != 1 {
		t.Fatalf("two earlier turns should fire once, got %d", len(v))
	}
	if v[0].RuleID != "cross_turn_repeat" {
		t.Errorf("rule id %q", v[0].RuleID)
	}
	if v[0].Matched != ps[0].Text {
		t.Errorf("matched %q, want the repeated phrase %q", v[0].Matched, ps[0].Text)
	}
	if !strings.Contains(v[0].Why, "2 earlier turns") {
		t.Errorf("the reason should say how many turns: %q", v[0].Why)
	}
}

// One violation per reply, and it names the worst offender. Reporting every
// repeated span would hand the writer a list to work through, which is the flat
// card again with a different heading.
func TestCrossTurnRepeatReportsTheBusiestPhraseOnce(t *testing.T) {
	ps := []Phrase{
		{Sum: 1, Text: "that is not the whole"},
		{Sum: 2, Text: "which is the whole point"},
		{Sum: 3, Text: "and this one is fresh"},
	}
	v := CrossTurnRepeat(ps, map[uint32]int{1: 2, 2: 5})
	if len(v) != 1 {
		t.Fatalf("got %d violations, want 1", len(v))
	}
	if v[0].Matched != "which is the whole point" {
		t.Errorf("matched %q, want the phrase repeated across the most turns", v[0].Matched)
	}
	if !strings.Contains(v[0].Context, "that is not the whole") {
		t.Errorf("the runner-up should still be named: %q", v[0].Context)
	}
}

func TestSumsMatchThePhrases(t *testing.T) {
	ps := []Phrase{{Sum: 7, Text: "a"}, {Sum: 9, Text: "b"}}
	got := Sums(ps)
	if len(got) != 2 || got[0] != 7 || got[1] != 9 {
		t.Errorf("Sums = %v", got)
	}
}

// A long reply can produce hundreds of shingles and twenty of those live in one
// session file. The cap is what keeps the state directory small.
func TestPhrasesAreCapped(t *testing.T) {
	var b strings.Builder
	for i := range 2000 {
		b.WriteString("that is not the whole point of it and neither was the one before number ")
		b.WriteString(strings.Repeat("x", 1+i%5))
		b.WriteString(". ")
	}
	if n := len(Phrases(b.String())); n > maxPhrases {
		t.Errorf("%d phrases from one turn, want at most %d", n, maxPhrases)
	}
}
