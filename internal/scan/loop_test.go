package scan

import (
	"strings"
	"testing"

	"github.com/justinstimatze/cope/card"
)

// The claim is "done"; the proof is a command, a count, or a file. Nothing here
// can tell a real transcript from an invented one — what it catches is the
// summary that asserts success with no anchor at all.
func TestUnverifiedDoneWantsSomethingToCheck(t *testing.T) {
	claimOnly := "Went through the remaining items and cleaned them up. " +
		"Everything passes now and the branch is in good shape for review."
	v := UnverifiedDone(claimOnly)
	if len(v) != 1 {
		t.Fatalf("a bare success claim gave %d violations, want 1: %+v", len(v), v)
	}
	if v[0].RuleID != "unverified_done" {
		t.Errorf("rule id %q", v[0].RuleID)
	}

	for name, text := range map[string]string{
		"a command": "Cleaned up the remaining items. `go test ./...` passes across " +
			"all four packages, so the branch is ready for review.",
		"a count": "Cleaned up the remaining items. 41 tests pass and the branch " +
			"is ready for somebody to look at it properly.",
		"a file": "Dropped the dead branch in render.go and the suite passes, so " +
			"the branch is ready for somebody to look at it properly.",
		"an admission": "Wrote the fix and it is not yet run. Everything passes on " +
			"the old path, which is all the evidence there currently is.",
	} {
		if v := UnverifiedDone(text); len(v) != 0 {
			t.Errorf("%s should clear the rule, got %+v", name, v)
		}
	}
}

// A report that stops without claiming anything has nothing to verify, so the
// rule must stay quiet rather than treating silence as a bad claim.
func TestUnverifiedDoneIgnoresAReportThatClaimsNothing(t *testing.T) {
	text := "Read through the three transcripts and pulled the phrases that " +
		"recur. Nothing was changed on disk and no decision is waiting."
	if v := UnverifiedDone(text); len(v) != 0 {
		t.Errorf("fired on a report making no success claim: %+v", v)
	}
}

// Nobody is there. The question goes into a log and the next iteration reads it
// as an instruction to itself.
func TestLoopAskFlagsOnlyTheEnding(t *testing.T) {
	ending := "Rebased onto main and reran the suite, four packages green.\n\n" +
		"Want me to push it or hold until you have looked?"
	v := LoopAsk(ending)
	if len(v) != 1 {
		t.Fatalf("an ending question gave %d violations, want 1: %+v", len(v), v)
	}

	// Mid-report questions are how a plan gets laid out, and a report often
	// quotes the question it is resolving.
	midway := "The question was whether the ceiling or the style was wrong.\n\n" +
		"It is the style: 1.5 hits 1024 tokens on an arithmetic word problem " +
		"while prod's p99 sits at 300. Ceiling left alone, note in NOTES.md."
	if v := LoopAsk(midway); len(v) != 0 {
		t.Errorf("fired on a question that is not the ending: %+v", v)
	}
}

// The lane decides which rules apply. Interactive is Check unchanged, so
// nothing measured before the lane existed changes meaning.
func TestCheckLaneSwapsTheEndingRules(t *testing.T) {
	c := &Card{}
	// Names an open problem, ends with no decision point: dangling_end in the
	// interactive lane, and exactly right in the loop lane.
	report := "Rebased onto main and reran the suite, four packages green.\n\n" +
		"The stale fixture in scan_test.go is still unresolved. I left it alone " +
		"because it passes and nothing in this run touched that path."

	interactive := ruleSet(c.CheckLane(report, 0.35, "interactive"))
	if !interactive["dangling_end"] {
		t.Errorf("interactive lane should still flag the open ending: %v", interactive)
	}
	loop := ruleSet(c.CheckLane(report, 0.35, LaneLoop))
	if loop["dangling_end"] {
		t.Errorf("loop lane kept a rule about a reader who is not there: %v", loop)
	}

	asking := "Rebased onto main and reran the suite, four packages green.\n\n" +
		"Want me to push it?"
	if !ruleSet(c.CheckLane(asking, 0.35, LaneLoop))["loop_ask"] {
		t.Error("loop lane should flag an ending nobody is there to answer")
	}
	if ruleSet(c.CheckLane(asking, 0.35, "interactive"))["loop_ask"] {
		t.Error("loop_ask fired in the interactive lane, where asking is the point")
	}
}

func ruleSet(v []Violation) map[string]bool {
	out := map[string]bool{}
	for _, x := range v {
		out[x.RuleID] = true
	}
	return out
}

// The suppression list has to name rules that exist, or it silently stops
// suppressing anything the day one is renamed.
func TestSuppressedRulesAreRealRuleIDs(t *testing.T) {
	c := &Card{}
	seen := map[string]bool{}
	for _, text := range []string{
		"Something is still broken here and there is quite a lot of it left " +
			"to work through before anybody can call this finished.\n\nThat is all.",
		"Want me to commit this now?\n\nI'll run the linter over the whole tree first if you would rather have that.",
		"Do you want the long version?\n\nEither way the rebase landed cleanly " +
			"and the suite is green across every package in the module.",
	} {
		for id := range ruleSet(c.CheckLane(text, 0.35, "interactive")) {
			seen[id] = true
		}
	}
	for id := range suppressedInLoop {
		if !seen[id] && !strings.HasPrefix(id, "buried") {
			t.Errorf("suppressedInLoop names %q, which no fixture here produces — "+
				"either the fixtures are wrong or the rule id is stale", id)
		}
	}
}

// The loop injection must not carry an item written for a reader who can
// answer. This is the overconstraint check: the first version of the loop card
// printed the rule against ending by asking directly above a worked example
// whose better half ends in a question, and the CONTINUE TEST under both.
func TestLoopRenderCarriesNothingThatAsks(t *testing.T) {
	c, _, err := ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatal(err)
	}
	facts := map[string]bool{"mode_loop": true, "rule_cross_turn_repeat": true}
	out := c.RenderFocus(facts, []Observation{{Rule: "cross_turn_repeat", N: 2}})
	if out == "" {
		t.Fatal("the loop lane renders nothing")
	}
	if !strings.Contains(out, "HANDOFF TEST") {
		t.Errorf("the loop lane dropped its own test:\n%s", out)
	}
	if strings.Contains(out, "CONTINUE TEST") {
		t.Errorf("the loop render carries the test for a reader who is present:\n%s", out)
	}

	// Nothing gated on a position may appear: those items all address somebody
	// about to reply.
	for _, p := range c.Wrong {
		if strings.HasPrefix(strings.TrimPrefix(p.When, "fact:"), "at_") &&
			strings.Contains(out, p.Right) {
			t.Errorf("an at_* pair reached the loop render: %.60q", p.Right)
		}
	}
	for _, tt := range c.Tests {
		if strings.HasPrefix(strings.TrimPrefix(tt.When, "fact:"), "at_") &&
			strings.Contains(out, tt.Name) {
			t.Errorf("an at_* test reached the loop render: %s", tt.Name)
		}
	}
}
