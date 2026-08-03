package scan

import (
	"strings"
	"testing"
)

func TestShapeParsesEveryPredicate(t *testing.T) {
	c := &Card{Shape: []string{
		"short_landing: last paragraph words <= 12 — the register lands short",
		"no_walls: every paragraph words <= 90 — a wall of text is not a paragraph",
		"one_short: some paragraph sentences <= 2 — the rhythm needs somewhere to breathe",
		"leaves_a_move: last paragraph asks — the reader has to have something to answer",
		"no_trailing_ask: last paragraph does not ask — this lane has no reader to ask",
		"capped: reply words >= 5 — a reply this short answered nothing",
	}}
	got, err := c.Shapes()
	if err != nil {
		t.Fatalf("Shapes: %v", err)
	}
	if len(got) != 6 {
		t.Fatalf("got %d rules, want 6", len(got))
	}
	if got[0].Selector != "last paragraph" || got[0].Kind != "words" || got[0].Op != "<=" || got[0].N != 12 {
		t.Errorf("short_landing = %+v", got[0])
	}
	if got[3].Kind != "asks" || got[3].Op != "is" {
		t.Errorf("leaves_a_move = %+v", got[3])
	}
	if got[4].Op != "is not" {
		t.Errorf("no_trailing_ask = %+v", got[4])
	}
}

func TestShapeRejectsWhatCannotBeChecked(t *testing.T) {
	for name, line := range map[string]string{
		"no reason":       "short_landing: last paragraph words <= 12",
		"no id":           "last paragraph words <= 12 — reason",
		"bad id":          "Short Landing: last paragraph words <= 12 — reason",
		"unknown select":  "x: third paragraph words <= 12 — reason",
		"unknown count":   "x: last paragraph clauses <= 12 — reason",
		"bad comparison":  "x: last paragraph words < 12 — reason",
		"not a number":    "x: last paragraph words <= many — reason",
		"zero":            "x: last paragraph words <= 0 — reason",
		"built-in name":   "dangling_end: last paragraph asks — shadows a compiled rule",
		"garbage predica": "x: reply banana — reason",
	} {
		c := &Card{Shape: []string{line}}
		if _, err := c.Shapes(); err == nil {
			t.Errorf("%s: %q was accepted", name, line)
		}
	}
}

func TestShapeRejectsCollisionWithAPostprocRule(t *testing.T) {
	c := &Card{
		Rules: []Rule{{ID: "flip"}},
		Shape: []string{"flip: last paragraph asks — collides with the card's own regex rule"},
	}
	if _, err := c.Shapes(); err == nil || !strings.Contains(err.Error(), "POSTPROC") {
		t.Errorf("err = %v, want a complaint about the collision", err)
	}
}

func TestShapeRejectsTheSameIDTwice(t *testing.T) {
	c := &Card{Shape: []string{
		"landing: last paragraph words <= 12 — one",
		"landing: first paragraph words <= 12 — two",
	}}
	if _, err := c.Shapes(); err == nil || !strings.Contains(err.Error(), "twice") {
		t.Errorf("err = %v, want a complaint about the repeat", err)
	}
}

func TestShapeCatchesALongLanding(t *testing.T) {
	c := &Card{Shape: []string{
		"short_landing: last paragraph words <= 12 — the VOICE block lands short",
	}}
	long := "The migration ran.\n\n" +
		"Two rows did not move and the reason is a null tenant id, which is a " +
		"data problem rather than a migration one and needs a decision before " +
		"the next run can be scheduled at all."
	v := c.Check(long, 0.35)
	if !hasRule(v, "short_landing") {
		t.Fatalf("no hit on a long landing: %v", ids(v))
	}
	for _, x := range v {
		if x.RuleID == "short_landing" && x.Why != "the VOICE block lands short" {
			t.Errorf("why = %q, want the card's own words", x.Why)
		}
	}

	short := "The migration ran.\n\nTwo rows did not move."
	if v := c.Check(short, 0.35); hasRule(v, "short_landing") {
		t.Errorf("fired on a short landing: %v", ids(v))
	}
}

func TestShapeEveryParagraphReportsEachFailure(t *testing.T) {
	c := &Card{Shape: []string{"no_walls: every paragraph words <= 5 — keep them short"}}
	text := "one two three four five six\n\nfine\n\nseven eight nine ten eleven twelve"
	v := c.Check(text, 0.35)
	n := 0
	for _, x := range v {
		if x.RuleID == "no_walls" {
			n++
		}
	}
	if n != 2 {
		t.Errorf("got %d hits, want one per failing paragraph (2)", n)
	}
}

func TestShapeSomeParagraphFiresOnlyWhenNoneQualifies(t *testing.T) {
	c := &Card{Shape: []string{"one_short: some paragraph words <= 4 — needs a beat somewhere"}}
	has := "a paragraph that runs on for a while\n\nshort one here"
	if v := c.Check(has, 0.35); hasRule(v, "one_short") {
		t.Errorf("fired when a qualifying paragraph existed: %v", ids(v))
	}
	none := "a paragraph that runs on for a while\n\nand another that also runs on"
	if v := c.Check(none, 0.35); !hasRule(v, "one_short") {
		t.Errorf("did not fire when none qualified: %v", ids(v))
	}
}

func TestShapeAsksSelectsOnQuestions(t *testing.T) {
	c := &Card{Shape: []string{"leaves_a_move: last paragraph asks — give the reader something to answer"}}
	if v := c.Check("Did the work.\n\nWant me to push it?", 0.35); hasRule(v, "leaves_a_move") {
		t.Errorf("fired on a reply that does ask: %v", ids(v))
	}
	if v := c.Check("Did the work.\n\nThe branch is green.", 0.35); !hasRule(v, "leaves_a_move") {
		t.Errorf("did not fire on a reply that does not ask: %v", ids(v))
	}
}

func TestGatingYourOwnShapeRuleSaysToDeleteItInstead(t *testing.T) {
	// @gate is for built-ins. A card gating a rule it wrote itself is two lines
	// cancelling out, and the generic "not a built-in rule" would send the
	// author looking for a typo in an id they spelled correctly.
	c := &Card{
		Shape: []string{"short_landing: last paragraph words <= 3 — lands short"},
		Gate:  []string{"short_landing off — actually this card explains itself at the end"},
	}
	_, err := c.Gates()
	if err == nil {
		t.Fatal("gating an own shape rule was accepted")
	}
	if !strings.Contains(err.Error(), "delete the @shape line") {
		t.Errorf("err = %v, want it to name the actual fix", err)
	}
}

func TestMalformedShapeLeavesNoRuleBehind(t *testing.T) {
	c := &Card{Shape: []string{"short_landing: last paragraph words <= 12"}} // no reason
	long := "One.\n\n" + strings.Repeat("word ", 40)
	if v := c.Check(long, 0.35); hasRule(v, "short_landing") {
		t.Errorf("a rule that failed validation still fired: %v", ids(v))
	}
}
