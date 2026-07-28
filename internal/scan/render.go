package scan

import (
	"fmt"
	"regexp"
	"strings"
)

// maxNeverRules mirrors effigy.prompt.MAX_NEVER_RULES (effigy v0.3.1):
// stope phase-probe data showed LLMs reliably attend to ~7-10 same-weight
// constraints, and rules beyond that silently lose attention rather than
// erroring. card2json.py pulls ast.never_would_say straight through with no
// cap, so this repo would otherwise be exposed to the exact failure effigy's
// own history already paid for.
const maxNeverRules = 10

// inlineExampleMarker matches effigy.prompt._INLINE_EXAMPLE_MARKERS: the
// uppercase prefixes an author uses to embed a worked example inside a NEVER
// rule. Case-sensitive so prose-lowercase "wrong" is untouched.
var inlineExampleMarker = regexp.MustCompile(`\b(?:WRONG|RIGHT|NOT|YES|BAD|GOOD|EXAMPLE|EX|BEFORE|AFTER)\s*:`)

// stripInlineExamples ports effigy.prompt._strip_inline_examples: an inline
// WRONG:/RIGHT:/... example has the same priming effect as a standalone
// WRONG[] pair, except unintended — the model pattern-matches the example
// text regardless of the rule's negative label. Cutting at the marker keeps
// the constraint statement and drops the example.
func stripInlineExamples(s string) string {
	loc := inlineExampleMarker.FindStringIndex(s)
	if loc == nil {
		return s
	}
	return strings.TrimRight(s[:loc[0]], " \n\t-—–:;,.")
}

// prioritizeNever ports effigy.prompt's never-rule budget: CRITICAL-prefixed
// rules sort first, the list is capped at maxNeverRules, and inline examples
// are stripped from what survives. Applied here rather than in card2json.py
// because this is cope's equivalent of effigy's build_static_context — the
// point where the card becomes prompt text, not where it's compiled.
//
// The second return is what the cap discarded. It exists because the cap and
// the eleventh rule landed in the same commit (584c898): the rule was written
// for a real complaint, truncated away before it ever rendered, and sat in the
// card file for two days reading as covered. A rule that silently does not
// ship is worse than no rule, so the overflow is now reported rather than
// dropped on the floor.
func prioritizeNever(never []string) (kept, dropped []string) {
	var critical, regular []string
	for _, n := range never {
		if strings.HasPrefix(strings.ToUpper(n), "CRITICAL:") {
			critical = append(critical, n)
		} else {
			regular = append(regular, n)
		}
	}
	all := append(critical, regular...)
	if len(all) > maxNeverRules {
		all, dropped = all[:maxNeverRules], all[maxNeverRules:]
	}
	kept = make([]string, len(all))
	for i, n := range all {
		kept[i] = stripInlineExamples(n)
	}
	return kept, dropped
}

// NeverOverflow returns the NEVER rules the budget will discard at render
// time. Empty is the healthy state; anything else means the card file claims
// a constraint the model is never shown.
func (c *Card) NeverOverflow() []string {
	_, dropped := prioritizeNever(c.Never)
	return dropped
}

// Render writes the card as prompt text for SessionStart injection.
//
// The rewrite pairs come first among the examples because a worse-and-better
// pair of the same content teaches a move that a prohibition only names. Four
// judge panels in this repo's history agreed on that ordering by accident:
// every criterion they could measure moved on shape, none on vocabulary.
//
// Deliberately not included: the POSTPROC patterns. Showing the regexes invites
// writing around them rather than not doing the thing.
func (c *Card) Render() string {
	var b strings.Builder

	b.WriteString("<voice-card>\n")
	b.WriteString("How to write in this session. These are constraints on prose only; they do not change what you do or how carefully you do it.\n")
	if c.Theme != "" {
		fmt.Fprintf(&b, "\nAim: %s\n", c.Theme)
	}

	if c.Voice.Kernel != "" {
		fmt.Fprintf(&b, "\nDefault register:\n%s\n", c.Voice.Kernel)
	}
	if c.Voice.Peak != "" {
		fmt.Fprintf(&b, "\nWhen flair is licensed:\n%s\n", c.Voice.Peak)
	}

	if len(c.Wrong) > 0 {
		b.WriteString("\nSame content, written worse and better:\n")
		for _, p := range c.Wrong {
			fmt.Fprintf(&b, "  worse:  %s\n", p.Wrong)
			fmt.Fprintf(&b, "  better: %s\n\n", p.Right)
		}
	}

	if len(c.Never) > 0 {
		b.WriteString("Never:\n")
		kept, _ := prioritizeNever(c.Never)
		for _, n := range kept {
			fmt.Fprintf(&b, "  - %s\n", n)
		}
	}

	if len(c.Tests) > 0 {
		b.WriteString("\nApply to your own draft before sending:\n")
		for _, t := range c.Tests {
			fmt.Fprintf(&b, "  %s — %s\n", t.Name, t.Question)
			for _, f := range t.Fail {
				fmt.Fprintf(&b, "      fails: %s\n", f)
			}
			for _, p := range t.Pass {
				fmt.Fprintf(&b, "      passes: %s\n", p)
			}
			if t.Why != "" {
				fmt.Fprintf(&b, "      why: %s\n", t.Why)
			}
		}
	}

	if len(c.Quirks) > 0 {
		b.WriteString("\nHabits worth keeping:\n")
		for _, q := range c.Quirks {
			fmt.Fprintf(&b, "  - %s\n", q)
		}
	}

	if len(c.Traits) > 0 {
		fmt.Fprintf(&b, "\nIn a phrase: %s\n", strings.Join(c.Traits, ", "))
	}

	b.WriteString("</voice-card>\n")
	return b.String()
}

// refresherTest names the card test the mid-session refresher is built from.
// The refresher deliberately renders one test, not the card: the full card
// re-arrives at every SessionStart (startup, resume, compact), and the dead
// stretch between compacts needs a reminder, not a second copy. Anthropic's
// Opus 5 prompting guidance recommends exactly this shape — a short reminder
// of the wanted style near the end of a long context.
const refresherTest = "CONTINUE TEST"

// RenderRefresher returns the compact mid-session reminder, or "" when the
// card carries no test by the refresher's name. TestShippedCardHasARefresher
// pins the shipped card against that silent empty.
func (c *Card) RenderRefresher() string {
	for _, t := range c.Tests {
		if t.Name != refresherTest {
			continue
		}
		var b strings.Builder
		b.WriteString("<voice-refresher>\nMid-session reminder — the ending is a decision surface. ")
		b.WriteString(t.Question)
		if len(t.Pass) > 0 {
			b.WriteString("\nEndings that pass: ")
			// Raw, like the full card's pass lines: the card author already
			// quotes example sentences, so %q would double-wrap them.
			for i, p := range t.Pass {
				if i > 0 {
					b.WriteString("  /  ")
				}
				b.WriteString(p)
			}
			b.WriteString("\n")
		}
		b.WriteString("</voice-refresher>\n")
		return b.String()
	}
	return ""
}
