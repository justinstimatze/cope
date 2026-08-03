package scan

import (
	"fmt"
	"strings"
	"testing"

	"github.com/justinstimatze/cope/card"
)

// SessionStart carries the card in full apart from the gates it holds shut.
//
// This used to assert that gating never narrowed SessionStart at all. It does
// now, for rule_* items only, because the NEVER cap is an attention budget on
// what the model is shown at once: with every gated rule rendered up front, the
// budget was spent on habits the writer had not exhibited and the card could
// hold no more than the cap in total. Holding evidence-gated rules until their
// detector fires is what lets the rule library grow past ten.
//
// at_* gates are unaffected. They mark which items the mid-session injection
// selects, not whether the writer has done something, and the CONTINUE TEST
// rides on one — closing it here would drop the card's central test from the
// only injection that states the card in full.
//
// mode_* is held for a third reason: it says what kind of session this is, and
// at SessionStart that is not known either. A session becomes a loop when a
// loop prompt arrives, which is after this render.
func TestFullRenderKeepsPositionGatedItemsAndHoldsEvidenceOnes(t *testing.T) {
	c, _, err := ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatal(err)
	}
	full := c.Render()

	position, held := 0, 0
	// Only at_* renders here; rule_* and mode_* both wait.
	isHeld := func(when string) bool {
		return !strings.HasPrefix(strings.TrimPrefix(when, "fact:"), "at_")
	}
	for _, n := range c.Never {
		if n.When == "" {
			continue
		}
		if isHeld(n.When) {
			held++
			if strings.Contains(full, stripInlineExamples(n.Text)) {
				t.Errorf("held NEVER rule rendered at SessionStart: %.60q", n.Text)
			}
			continue
		}
		position++
		if !strings.Contains(full, stripInlineExamples(n.Text)) {
			t.Errorf("position-gated NEVER rule missing from the full card: %.60q", n.Text)
		}
	}
	for _, w := range c.Wrong {
		if w.When == "" || isHeld(w.When) {
			continue
		}
		position++
		if !strings.Contains(full, w.Right) {
			t.Errorf("position-gated pair missing from the full card: %.60q", w.Wrong)
		}
	}
	for _, tt := range c.Tests {
		if tt.When == "" || isHeld(tt.When) {
			continue
		}
		position++
		if !strings.Contains(full, tt.Name) {
			t.Errorf("position-gated test missing from the full card: %s", tt.Name)
		}
	}
	if held == 0 {
		t.Error("no held NEVER rules in the shipped card, so the holding-back half proves nothing")
	}
	if position == 0 {
		t.Error("no position-gated items in the shipped card, so the keeping half proves nothing")
	}
}

// The budget is charged against what an injection prints, not against the card
// file. A rule the facts do not activate is not competing for the model's
// attention and must not evict one that is.
func TestEvidenceGatedRulesDoNotSpendTheBudget(t *testing.T) {
	var never []Gated
	for range maxNeverRules {
		never = append(never, Gated{Text: "always-on rule"})
	}
	never = append(never, Gated{Text: "dormant rule", When: "fact:rule_never_fired"})
	c := &Card{Never: never}

	full := c.Render()
	if strings.Contains(full, "dormant rule") {
		t.Error("a rule with no evidence behind it rendered at SessionStart")
	}
	if n := strings.Count(full, "always-on rule"); n != maxNeverRules {
		t.Errorf("rendered %d always-on rules, want all %d — the dormant rule took a slot", n, maxNeverRules)
	}
	// No overflow: SessionStart prints the cap exactly and the refresher prints
	// one. The union is over the cap, but nothing renders the union.
	if over := c.NeverOverflow(); len(over) != 0 {
		t.Errorf("NeverOverflow reported %d rules over budget, want 0 — no injection prints more than the cap: %+v", len(over), over)
	}

	// The overflow it must still catch: one injection over the cap on its own.
	tooMany := &Card{Never: append(append([]Gated{}, never[:maxNeverRules]...), Gated{Text: "the eleventh always-on rule"})}
	if over := tooMany.NeverOverflow(); len(over) != 1 {
		t.Errorf("an eleventh always-on rule reported %d over budget, want 1", len(over))
	}

	// And the same on the refresher's side, where the always-on rules are not
	// printed and so must not be what pushes it over.
	var gated []Gated
	for i := range maxNeverRules + 1 {
		gated = append(gated, Gated{Text: fmt.Sprintf("gated rule %d", i), When: fmt.Sprintf("fact:rule_%d", i)})
	}
	if over := (&Card{Never: gated}).NeverOverflow(); len(over) != 1 {
		t.Errorf("an eleventh evidence-gated rule reported %d over budget, want 1", len(over))
	}
}

// The mid-session render is the opposite: it carries only what the facts
// activate, and says what it observed.
func TestFocusCarriesOnlyTheLiveRules(t *testing.T) {
	c := &Card{
		Never: []Gated{
			{Text: "always-on rule"},
			{Text: "bold rule", When: "fact:rule_bold_label"},
			{Text: "flip rule", When: "fact:rule_flip"},
		},
	}
	out := c.RenderFocus(
		map[string]bool{"rule_bold_label": true},
		[]Observation{{Rule: "bold_label", N: 3}},
	)
	if !strings.Contains(out, "bold rule") {
		t.Error("the live rule is missing")
	}
	if strings.Contains(out, "flip rule") {
		t.Error("a rule that has not been firing was injected anyway")
	}
	// Always-on items belong to the full card; repeating them mid-session is
	// the flat injection again.
	if strings.Contains(out, "always-on rule") {
		t.Error("ungated items should not be repeated in the focused injection")
	}
	if !strings.Contains(out, "bold_label ×3") {
		t.Errorf("the observation should be stated with its count:\n%s", out)
	}

	// Nothing live and nothing gated means no injection at all, so the caller
	// can fall back rather than print an empty or invented reminder.
	if got := c.RenderFocus(map[string]bool{}, nil); got != "" {
		t.Errorf("expected empty focus with no live rules, got:\n%s", got)
	}
}

// The positive render is the arm that tests whether naming a move is what puts
// it in reach, so it must not carry the move's name, the bad version of it, or
// a prohibition against it. Anything that leaks defeats the arm silently.
func TestPositiveRenderCarriesNoAntiPattern(t *testing.T) {
	c, _, err := ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatal(err)
	}
	facts := map[string]bool{"at_prompt": true}
	for _, id := range []string{"bold_label", "labelled_opening", "clause_symmetry"} {
		facts["rule_"+id] = true
	}
	observed := []Observation{{Rule: "bold_label", N: 3}, {Rule: "labelled_opening", N: 2}}

	paired := c.RenderFocus(facts, observed)
	positive := c.RenderFocusPositive(facts, observed)
	if positive == "" {
		t.Fatal("the positive render is empty on facts the paired render answers")
	}

	for _, banned := range []string{"worse:", "why:", "Never:", "bold_label", "labelled_opening", "×3"} {
		if strings.Contains(positive, banned) {
			t.Errorf("positive render leaked %q:\n%s", banned, positive)
		}
	}

	// Every wrong/better pair the paired render reached must still contribute
	// its better half, or the positive arm is a shorter reminder rather than a
	// positive one.
	for _, p := range c.Wrong {
		if p.When == "" || !active(p.When, facts) {
			continue
		}
		if !strings.Contains(positive, p.Right) {
			t.Errorf("positive render dropped the better line %q", p.Right)
		}
		if strings.Contains(positive, p.Wrong) {
			t.Errorf("positive render kept the worse line %q", p.Wrong)
		}
	}

	if len(positive) >= len(paired) {
		t.Errorf("positive render is %d chars against the paired %d; it should be the shorter one", len(positive), len(paired))
	}
}
