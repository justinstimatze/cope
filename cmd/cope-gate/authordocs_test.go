package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/justinstimatze/cope/card"
	"github.com/justinstimatze/cope/internal/scan"
)

// testFlagSet registers the binary's real flags on a throwaway set, which is
// the point: the docs prompt describes whatever registerFlags declares, so a
// new flag shows up in the README without anyone remembering to add it.
func testFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("cope-gate", flag.ContinueOnError)
	registerFlags(fs)
	return fs
}

func embeddedCard(t *testing.T) *scan.Card {
	t.Helper()
	c, _, err := scan.ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// The facts are the whole contribution of the generator, so they have to come
// from the card and the flag set rather than from a literal that ages.
func TestCollectFactsReadsTheCardNotAConstant(t *testing.T) {
	c := embeddedCard(t)
	f := collectFacts(c, testFlagSet())

	if f.Composition.Never != len(c.Never) || f.Composition.Tests != len(c.Tests) ||
		f.Composition.Wrong != len(c.Wrong) || f.Composition.Regex != len(c.Rules) {
		t.Errorf("composition disagrees with the card: %+v", f.Composition)
	}
	if f.RenderedCard.Card != len(c.Render()) {
		t.Errorf("card size %d, want %d", f.RenderedCard.Card, len(c.Render()))
	}
	if f.RenderedCard.Refresher == 0 {
		t.Error("refresher size is 0 — the CONTINUE TEST lookup broke")
	}
	if len(f.RegexRules) != len(c.Rules) {
		t.Errorf("got %d regex rules, want %d", len(f.RegexRules), len(c.Rules))
	}
	if len(f.TestNames) != len(c.Tests) {
		t.Errorf("got %d test names, want %d", len(f.TestNames), len(c.Tests))
	}
	if len(f.NeverOverflow) != 0 {
		t.Errorf("shipped card has NEVER rules over budget: %v", f.NeverOverflow)
	}
	if f.Module == "" || f.NameOrigin == "" {
		t.Error("module and name origin must be stated, not left blank")
	}
	// git describe changes every commit, so a version in a README is stale
	// the moment it is written. The first generated draft printed one.
	if blob, _ := json.Marshal(f); strings.Contains(string(blob), buildVersion()) && buildVersion() != "dev" {
		t.Errorf("facts carry the build version %q; a generated README will print it", buildVersion())
	}
	if _, err := json.Marshal(f); err != nil {
		t.Fatalf("facts must serialize for the prompt: %v", err)
	}
}

// --author-docs is emitted after flag.Parse in main, so every flag the binary
// accepts is described without anyone maintaining a table.
func TestCollectFactsPicksUpTheRealFlagSet(t *testing.T) {
	f := collectFacts(embeddedCard(t), testFlagSet())
	got := map[string]bool{}
	for _, fl := range f.Flags {
		got[fl.Name] = true
	}
	for _, want := range []string{"--check", "--author-docs", "--inject", "--refresher", "--backfill", "--block", "--log", "--rules", "--min-cv", "--version"} {
		if !got[want] {
			t.Errorf("flag %s is missing from the facts", want)
		}
	}
}

// The facts are read off the machine writing the docs, so this is the one seam
// where a private path can walk into a published README. --log resolves an
// absolute default under $HOME and did exactly that.
func TestDocsPromptCarriesNoHomeDirectory(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("no home directory to leak")
	}
	p := buildDocsPrompt(embeddedCard(t), testFlagSet())
	if strings.Contains(p, home) {
		for _, line := range strings.Split(p, "\n") {
			if strings.Contains(line, home) {
				t.Errorf("docs prompt carries the home directory: %s", strings.TrimSpace(line))
			}
		}
	}
	if !strings.Contains(p, "$HOME/.local/state/cope/violations.jsonl") &&
		!strings.Contains(p, "$XDG_STATE_HOME") {
		t.Error("the log default should still be stated, in symbolic form")
	}
}

func TestDocsPromptCarriesTheCardAndTheFacts(t *testing.T) {
	c := embeddedCard(t)
	p := buildDocsPrompt(c, testFlagSet())

	for _, want := range []string{
		"<cope-author-docs>", "</cope-author-docs>",
		"<voice-card>", c.Theme,
		"cope-gate --check README.md",
		"## The facts", "## The sections",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("docs prompt is missing %q", want)
		}
	}
	// A prompt that supplies phrasing produces the same README every time,
	// which is the failure this whole path exists to avoid.
	if strings.Contains(p, "A voice card for Claude's prose, plus a Stop hook") {
		t.Error("docs prompt contains README prose; it must contribute facts only")
	}
}

// Same invariant Render holds: showing the patterns invites writing around them
// instead of not doing the thing.
func TestDocsPromptWithholdsTheRegexPatterns(t *testing.T) {
	c := embeddedCard(t)
	p := buildDocsPrompt(c, testFlagSet())
	for _, r := range c.Rules {
		if r.Pattern == "" {
			continue
		}
		if strings.Contains(p, r.Pattern) {
			t.Errorf("docs prompt leaked the pattern for %s", r.ID)
		}
		if !strings.Contains(p, r.ID) {
			t.Errorf("docs prompt should still name rule %s and say what it is for", r.ID)
		}
	}
}

// TestRuleIDsMatchTheDocumentedSet keeps shapeRules and scan.RuleIDs from
// drifting. They are two lists of the same thing for different reasons: this
// one carries descriptions for the README, scan's validates @gate lines at card
// load. A rule in only one of them is either undocumented or unexemptable.
func TestRuleIDsMatchTheDocumentedSet(t *testing.T) {
	documented := map[string]bool{}
	for _, r := range shapeRules {
		documented[r.ID] = true
	}
	for id := range scan.RuleIDs {
		if !documented[id] {
			t.Errorf("scan.RuleIDs has %q, shapeRules does not — it has no README description", id)
		}
	}
	for id := range documented {
		if !scan.RuleIDs[id] {
			t.Errorf("shapeRules has %q, scan.RuleIDs does not — no card can name it in @gate", id)
		}
	}
}
