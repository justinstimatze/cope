package scan

import (
	"strings"
	"testing"
)

func TestGateParsesRuleAndReason(t *testing.T) {
	c := &Card{Gate: []string{"dangling_end off — VOICE holds the conclusion back"}}
	got, err := c.Gates()
	if err != nil {
		t.Fatalf("Gates: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d directives, want 1", len(got))
	}
	if got[0].Rule != "dangling_end" {
		t.Errorf("rule = %q, want dangling_end", got[0].Rule)
	}
	if got[0].Why != "VOICE holds the conclusion back" {
		t.Errorf("why = %q", got[0].Why)
	}
}

func TestGateAcceptsHyphenAndNoReason(t *testing.T) {
	for _, line := range []string{
		"apology off - the register never apologises",
		"apology off -- the register never apologises",
		"apology off",
	} {
		c := &Card{Gate: []string{line}}
		got, err := c.Gates()
		if err != nil {
			t.Fatalf("%q: %v", line, err)
		}
		if len(got) != 1 || got[0].Rule != "apology" {
			t.Fatalf("%q parsed to %+v", line, got)
		}
	}
}

func TestGateRejectsWhatWouldSilentlyDoNothing(t *testing.T) {
	// Each of these looks like a working exemption in the card file. The rule
	// they name keeps firing, so a card believing itself exempt is the whole
	// failure mode — every one has to be loud.
	for name, line := range map[string]string{
		"unknown rule": "dangling_ends off — typo in the id",
		"no state":     "dangling_end",
		"state is on":  "dangling_end on — rules are on by default",
		"too my words": "dangling_end off please — extra word before the reason",
	} {
		c := &Card{Gate: []string{line}}
		if _, err := c.Gates(); err == nil {
			t.Errorf("%s: %q was accepted", name, line)
		}
	}
}

func TestGateRejectsTheSameRuleTwice(t *testing.T) {
	c := &Card{Gate: []string{"apology off — one", "apology off — two"}}
	if _, err := c.Gates(); err == nil || !strings.Contains(err.Error(), "twice") {
		t.Errorf("err = %v, want a complaint about the repeat", err)
	}
}

func TestExemptRuleIsDroppedFromTheScore(t *testing.T) {
	// A reply naming an open problem in its last block with no question after
	// it, which is what dangling_end is for.
	text := "The migration ran and moved 41 rows.\n\n" +
		"Two rows did not move and the reason is still open."

	plain := &Card{}
	if got := plain.Check(text, 0.35); len(got) == 0 {
		t.Fatal("fixture no longer trips any rule; the test needs a new one")
	} else if !hasRule(got, "dangling_end") {
		t.Fatalf("fixture trips %v, wanted dangling_end", ids(got))
	}

	exempt := &Card{Gate: []string{"dangling_end off — the card asks for this ending"}}
	if got := exempt.Check(text, 0.35); hasRule(got, "dangling_end") {
		t.Errorf("dangling_end survived the exemption: %v", ids(got))
	}
}

func TestExemptionDoesNotLeakToAnotherCard(t *testing.T) {
	// The detectors keep running; only this card's score changes. Two cards
	// built in the same process must not see each other's exemptions.
	text := "The migration ran and moved 41 rows.\n\n" +
		"Two rows did not move and the reason is still open."
	exempt := &Card{Gate: []string{"dangling_end off — deliberate"}}
	_ = exempt.Check(text, 0.35)

	if got := (&Card{}).Check(text, 0.35); !hasRule(got, "dangling_end") {
		t.Errorf("a plain card stopped seeing dangling_end: %v", ids(got))
	}
}

func TestExemptAppliesToTheLoopLane(t *testing.T) {
	// loop_ask is added by CheckLane after Check has already run, so it is the
	// rule an exemption would miss if Allow only ran in one place.
	text := "Rebuilt the index and it is green.\n\nRan `go test ./...`, 41 tests.\n\n" +
		"Should I also drop the old column?"

	plain := &Card{}
	if got := plain.CheckLane(text, 0.35, LaneLoop); !hasRule(got, "loop_ask") {
		t.Fatalf("fixture trips %v, wanted loop_ask", ids(got))
	}
	exempt := &Card{Gate: []string{"loop_ask off — this loop writes for a reader who does return"}}
	if got := exempt.CheckLane(text, 0.35, LaneLoop); hasRule(got, "loop_ask") {
		t.Errorf("loop_ask survived the exemption: %v", ids(got))
	}
}

func TestMalformedGateLeavesTheRuleOn(t *testing.T) {
	text := "The migration ran and moved 41 rows.\n\n" +
		"Two rows did not move and the reason is still open."
	c := &Card{Gate: []string{"dangling_ends off — typo"}}
	if got := c.Check(text, 0.35); !hasRule(got, "dangling_end") {
		t.Errorf("a typo'd exemption turned the rule off: %v", ids(got))
	}
}

func hasRule(v []Violation, id string) bool {
	for _, x := range v {
		if x.RuleID == id {
			return true
		}
	}
	return false
}

func ids(v []Violation) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.RuleID)
	}
	return out
}
