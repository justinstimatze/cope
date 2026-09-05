package scan

import "testing"

// fires reports whether any violation carries the id.
func fires(v []Violation, id string) bool {
	for _, x := range v {
		if x.RuleID == id {
			return true
		}
	}
	return false
}

func TestEchoedHeading(t *testing.T) {
	echo := "## Cache invalidation\n\n" +
		"Cache invalidation is the part nobody measured, and the current path reruns it on every write."
	if !fires(EchoedHeading(echo), "echoed_heading") {
		t.Error("did not fire on a first sentence repeating the whole heading")
	}

	// The heading names the topic and the sentence says something about it.
	fresh := "## Cache invalidation\n\n" +
		"Every write reruns the whole map, which is the cost that never showed up in the profile."
	if fires(EchoedHeading(fresh), "echoed_heading") {
		t.Error("fired on a sentence that adds to its heading")
	}

	// One content word is not an echo: a heading of one topic word and a
	// sentence using it is ordinary prose, and flagging it would flag most
	// documents. The rule needs two.
	thin := "## Caching\n\nCaching is the part nobody measured on this path."
	if fires(EchoedHeading(thin), "echoed_heading") {
		t.Error("fired on a one-word heading")
	}
}

func TestRepeatedOpening(t *testing.T) {
	three := "The gate reads the transcript once per turn. The gate writes nothing when the " +
		"reply comes back clean. The gate never edits what has already been sent."
	if !fires(RepeatedOpening(three), "repeated_opening") {
		t.Error("did not fire on three sentences opening the same way")
	}

	// Two is a rhythm. A rule that fires here would fire on deliberate
	// anaphora, which the source page itself says to leave alone.
	two := "The gate reads the transcript once per turn. The gate writes nothing when the " +
		"reply comes back clean. Nothing else touches the file after that."
	if fires(RepeatedOpening(two), "repeated_opening") {
		t.Error("fired on two")
	}
}

func TestFragmentRun(t *testing.T) {
	run := "The profile came back this morning and it points somewhere unexpected. " +
		"Lock contention. Not algorithmic cost. The same shape twice."
	if !fires(FragmentRun(run), "fragment_run") {
		t.Error("did not fire on three consecutive fragments")
	}

	// Two fragments and a sentence between them is this repo's own register.
	// A rule firing here would fire on the card that ships with it.
	paced := "The profile came back this morning and it points somewhere unexpected. " +
		"Lock contention. The remedy is a per-shard mutex, roughly a day of work. Not algorithmic cost."
	if fires(FragmentRun(paced), "fragment_run") {
		t.Error("fired on fragments a real sentence separates")
	}
}

// Redact leaves a run of spaces where a code span was, so a sentence opening on
// one has no prose opening. Reading across the hole produced a key whose two
// words are not adjacent in the reply, and a quoted match full of blanks.
func TestRepeatedOpeningIgnoresCodeLedSentences(t *testing.T) {
	code := "`cope-gate --refresher` runs on every prompt and injects nothing most of the time. " +
		"`cope-gate --pretool` runs on a Linear write and scores the body it carries. " +
		"`cope-gate --check` runs on a file and prints what fired."
	if v := RepeatedOpening(code); len(v) != 0 {
		t.Errorf("keyed three sentences across their redacted openings: %+v", v)
	}

	// The same three with the code spans written as prose still fire, so the
	// skip is about the hole and not about the sentences.
	prose := "The refresher runs on every prompt and injects nothing most of the time. " +
		"The refresher runs on a Linear write and scores the body it carries. " +
		"The refresher runs on a file and prints what fired."
	if !fires(RepeatedOpening(prose), "repeated_opening") {
		t.Error("the same three openings written without code spans did not fire")
	}
}

// The three run in every lane. They are about a reader who is reading, not one
// who is about to answer, so the loop and external lanes keep all of them —
// which is the first thing either lane has gained rather than lost.
func TestTellsSurviveEveryLane(t *testing.T) {
	text := "The profile came back this morning and it points somewhere unexpected. " +
		"Lock contention. Not algorithmic cost. The same shape twice."
	for _, lane := range []string{LaneInteractive, LaneLoop, LaneExternal} {
		if !fires((&Card{}).CheckLane(text, 0.35, lane), "fragment_run") {
			t.Errorf("lane %s dropped fragment_run", lane)
		}
	}
}

func TestValidLane(t *testing.T) {
	for _, ok := range []string{"", LaneInteractive, LaneLoop, LaneExternal} {
		if err := ValidLane(ok); err != nil {
			t.Errorf("%q rejected: %v", ok, err)
		}
	}
	if ValidLane("extenal") == nil {
		t.Error("a misspelled lane was accepted, which is how a caller scores a lane it did not ask for")
	}
}
