package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/justinstimatze/cope/internal/scan"
)

func TestCardPromptCarriesTheSampleWhole(t *testing.T) {
	sample := "The suite runs in forty seconds. It ran in twenty-five before the parser change, " +
		"and nothing in this branch explains the difference yet."
	got := buildCardPrompt(sample, "notes.md")
	if !strings.Contains(got, sample) {
		t.Error("the sample is not in the prompt — a card derived from a summary of a voice is a card about a summary")
	}
	if !strings.Contains(got, "notes.md") {
		t.Error("the prompt does not name where the sample came from")
	}
}

// The @gate line lists what a card may decline. A rule missing from it is a
// rule the author is never told they can turn off, which is how a card ends up
// working around a rule instead of declining it.
func TestCardPromptNamesEveryGateableRule(t *testing.T) {
	got := buildCardPrompt("some prose", "-")
	for id := range scan.RuleIDs {
		if !strings.Contains(got, id) {
			t.Errorf("prompt does not name rule %q", id)
		}
	}
}

func TestCardPromptStatesTheRealBudget(t *testing.T) {
	got := buildCardPrompt("some prose", "-")
	if !strings.Contains(got, fmt.Sprintf("At most %d reach the prompt", scan.MaxNeverRules())) {
		t.Errorf("prompt does not carry the NEVER budget of %d", scan.MaxNeverRules())
	}
}
