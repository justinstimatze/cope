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

// MaxNeverRules is the budget, for anything describing the card format rather
// than rendering it. The number was written out by hand in the prompt that
// teaches somebody to author a card, which is the one place a stale copy would
// teach the wrong budget.
func MaxNeverRules() int { return maxNeverRules }

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
// The cap applies to what this injection will actually print, not to the card
// file's total, so gating is resolved first. The budget is an attention limit
// on same-weight constraints the model is shown at once (see maxNeverRules);
// a rule held back for a habit the writer has not exhibited is not competing
// for that attention and must not consume a slot. Charging the cap before
// gating capped the library instead of the prompt, which is the opposite of
// what the effigy finding says to bound.
//
// The second return is what the cap discarded. It exists because the cap and
// the eleventh rule landed in the same commit (584c898): the rule was written
// for a real complaint, truncated away before it ever rendered, and sat in the
// card file for two days reading as covered. A rule that silently does not
// ship is worse than no rule, so the overflow is now reported rather than
// dropped on the floor.
func prioritizeNever(never []Gated, facts map[string]bool) (kept, dropped []Gated) {
	var critical, regular []Gated
	for _, n := range never {
		if !active(n.When, facts) {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(n.Text), "CRITICAL:") {
			critical = append(critical, n)
		} else {
			regular = append(regular, n)
		}
	}
	all := append(critical, regular...)
	if len(all) > maxNeverRules {
		all, dropped = all[:maxNeverRules], all[maxNeverRules:]
	}
	kept = make([]Gated, len(all))
	for i, n := range all {
		kept[i] = Gated{Text: stripInlineExamples(n.Text), When: n.When}
	}
	return kept, dropped
}

// active reports whether a gated item belongs in this injection.
//
// An empty condition is always-on. Anything else is effigy's fact-set form,
// and the item renders when the fact is present. Facts come from
// internal/state: at_* for the injection point, rule_* for what the writer
// has actually been tripping this session.
//
// A nil fact set means "render everything". No injection passes nil — both
// build an explicit set — so it survives for callers constructing a card by
// hand and for tests that want every gate open.
func active(when string, facts map[string]bool) bool {
	if when == "" {
		return true
	}
	if facts == nil {
		return true
	}
	return facts[strings.TrimPrefix(when, "fact:")]
}

// NeverOverflow returns the NEVER rules the budget will discard at render
// time. Empty is the healthy state; anything else means the card file claims
// a constraint the model is never shown.
func (c *Card) NeverOverflow() []Gated {
	// Charged against each injection separately, because no code path renders
	// their union: SessionStart prints the always-on rules and the refresher
	// prints the evidence-gated ones. Summing the two would report an overflow
	// no reader ever meets, and would cap the rule library at the size of one
	// injection — which is the bug this replaced.
	var out []Gated
	_, dropped := c.neverFor(c.sessionStartFacts(), false)
	out = append(out, dropped...)

	// The refresher's worst case: every gate SessionStart holds shut open at
	// once — a loop session where every detector has also fired.
	all := map[string]bool{}
	for _, n := range c.Never {
		name := strings.TrimPrefix(n.When, "fact:")
		if strings.HasPrefix(name, "rule_") || strings.HasPrefix(name, "mode_") {
			all[name] = true
		}
	}
	_, dropped = c.neverFor(all, true)
	return append(out, dropped...)
}

// neverFor returns the NEVER rules one injection prints, capped.
//
// gatedOnly selects the refresher's slice — the rules naming a habit this
// session has actually shown. The always-on rules arrived with the SessionStart
// card and are not repeated there, so they must not consume the refresher's
// budget either; the cap belongs to what prints.
func (c *Card) neverFor(facts map[string]bool, gatedOnly bool) (kept, dropped []Gated) {
	pool := c.Never
	if gatedOnly {
		pool = nil
		for _, n := range c.Never {
			if n.When != "" {
				pool = append(pool, n)
			}
		}
	}
	return prioritizeNever(pool, facts)
}

// speakers rewrites effigy's Ali:Chat placeholders for a card that is injected
// into a live session rather than driving an NPC. Leaving {{char}} in the
// prompt asks the reader to work out who that is.
var speakers = strings.NewReplacer("{{user}}:", "user:", "{{char}}:", "you:")

// Render writes the card as prompt text for SessionStart injection.
//
// The rewrite pairs come first among the examples because a worse-and-better
// pair of the same content teaches a move that a prohibition only names. Four
// judge panels in this repo's history agreed on that ordering by accident:
// every criterion they could measure moved on shape, none on vocabulary.
//
// Deliberately not included: the POSTPROC patterns. Showing the regexes invites
// writing around them rather than not doing the thing.
//
// Evidence gates stay closed here; every other gate is open. See
// sessionStartFacts.
func (c *Card) Render() string { return c.render(c.sessionStartFacts()) }

// Describe renders the card as a target to recognise rather than as an
// instruction to follow: the aim, the default register, and where flair is
// licensed. Nothing else.
//
// This is what the discrimination test shows its reader. The rest of the card —
// the prohibitions, the worked pairs, the tests — is machinery for producing
// the voice, and a reader holding that machinery is marking prose against a
// checklist rather than recognising a voice. Those are different questions and
// only one of them is whether the voicing arrived.
//
// The card describing itself is also the only honest source for the target. A
// description written separately is one more artefact that can be good or bad
// on its own, and a reader who then failed to match it would leave two readings
// standing: the voice did not land, or the description was poor. Taking the
// text out of the card collapses that to one.
func (c *Card) Describe() string {
	var b strings.Builder
	if c.Theme != "" {
		fmt.Fprintf(&b, "Aim: %s\n", c.Theme)
	}
	if c.Voice.Kernel != "" {
		fmt.Fprintf(&b, "\nDefault register:\n%s\n", c.Voice.Kernel)
	}
	if c.Voice.Peak != "" {
		fmt.Fprintf(&b, "\nWhen flair is licensed:\n%s\n", c.Voice.Peak)
	}
	return b.String()
}

// sessionStartFacts opens every gate except the evidence ones.
//
// The two gate namespaces mean different things and only one of them should
// narrow this render. A rule_* gate says the writer has been tripping that
// rule; at SessionStart nothing has been measured, so those items stay out and
// arrive through the refresher when their detector fires. That is the whole
// mechanism a pasted CLAUDE.md does not have, and it is also what keeps the
// NEVER budget spendable: the cap is an attention limit on constraints shown
// at once (see maxNeverRules), so a rule held back for a habit this writer has
// not exhibited must not occupy one of the slots.
//
// Three gate namespaces, and only one of them opens here.
//
// at_* is a position marker rather than a condition — it selects which items
// the mid-session injection carries, and the item still belongs in the complete
// card. Closing those would drop the CONTINUE TEST from the only injection that
// states the card in full, so they open.
//
// rule_* is evidence the writer has been tripping a detector, and nothing has
// been measured at SessionStart. mode_* is what kind of session this is —
// mode_loop is the unattended lane — and at SessionStart that is not known
// either: a session becomes a loop when a loop prompt arrives, which is after
// this. Both stay shut and arrive through the refresher.
func (c *Card) sessionStartFacts() map[string]bool {
	f := map[string]bool{}
	open := func(when string) {
		if when == "" {
			return
		}
		name := strings.TrimPrefix(when, "fact:")
		if strings.HasPrefix(name, "at_") {
			f[name] = true
		}
	}
	for _, n := range c.Never {
		open(n.When)
	}
	for _, w := range c.Wrong {
		open(w.When)
	}
	for _, t := range c.Tests {
		open(t.When)
	}
	return f
}

// render builds the card, keeping only the gated items this fact set activates.
// nil facts keep everything, which is what SessionStart wants.
func (c *Card) render(facts map[string]bool) string {
	var b strings.Builder

	b.WriteString("<voice-card>\n")
	b.WriteString("How to write in this session, and how to hand work back. Brevity belongs to the prose and not to the work: write less about what you did, do not do less.\n")
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
			if !active(p.When, facts) {
				continue
			}
			fmt.Fprintf(&b, "  worse:  %s\n", p.Wrong)
			fmt.Fprintf(&b, "  better: %s\n", p.Right)
			if p.Why != "" {
				fmt.Fprintf(&b, "  why:    %s\n", p.Why)
			}
			b.WriteString("\n")
		}
	}

	if len(c.Mes) > 0 {
		b.WriteString("Whole replies in this voice:\n\n")
		for _, m := range c.Mes {
			for _, line := range strings.Split(strings.TrimSpace(m), "\n") {
				fmt.Fprintf(&b, "  %s\n", speakers.Replace(line))
			}
			b.WriteString("\n")
		}
	}

	if len(c.Never) > 0 {
		b.WriteString("Never:\n")
		kept, _ := c.neverFor(facts, false)
		for _, n := range kept {
			fmt.Fprintf(&b, "  - %s\n", n.Text)
		}
	}

	if len(c.Tests) > 0 {
		b.WriteString("\nApply to your own draft before sending:\n")
		for _, t := range c.Tests {
			if !active(t.When, facts) {
				continue
			}
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

// Observation is one rule and how often it has fired recently. main builds
// these from internal/state; scan stays free of that dependency so the card
// renderer can be used without a session on disk.
type Observation struct {
	Rule string
	N    int
}

// RenderFocus is the mid-session injection, narrowed to what this session has
// been getting wrong.
//
// This is the only thing cope renders that a pasted CLAUDE.md cannot: the text
// is chosen from the writer's own measured output rather than fixed by an
// author in advance. Items gated on rule_<id> appear when that rule has been
// firing, and the counts are stated so the reminder is about something that
// happened rather than something that might.
//
// Returns "" when nothing has fired and no gated item matches. The caller falls
// back to RenderRefresher, so a quiet session gets the standing reminder rather
// than an empty block or a fabricated complaint.
func (c *Card) RenderFocus(facts map[string]bool, observed []Observation) string {
	return c.renderFocus(facts, observed, false)
}

// RenderFocusPositive is the same slice with every anti-pattern taken out: the
// better half of each wrong/better pair and the passing examples, and none of
// the prohibitions, none of the worse lines, and no rule names.
//
// It is the second arm of the question the first experiment could not answer.
// A reminder against opening on a bold label has to write "bold label" down,
// and the card's pairs print the bad version next to the good one — so a null
// or a backwards result has two readings, and only one of them is about the
// reminder failing to land. This render keeps the selection (which examples
// appear is still gated on what fired) and drops the naming.
//
// The counts stay out for the same reason. "bold_label ×3" is the anti-pattern
// in the fewest possible characters.
func (c *Card) RenderFocusPositive(facts map[string]bool, observed []Observation) string {
	return c.renderFocus(facts, observed, true)
}

func (c *Card) renderFocus(facts map[string]bool, observed []Observation, positive bool) string {
	var items strings.Builder

	// Prohibitions are anti-patterns however politely they are worded, so the
	// positive render has no equivalent block and simply drops them.
	if !positive {
		// prioritizeNever has already dropped what this fact set does not
		// activate; the remaining filter keeps the always-on rules out, since
		// those arrived with the SessionStart card and restating them is what
		// a pasted file already does.
		kept, _ := c.neverFor(facts, true)
		for _, n := range kept {
			fmt.Fprintf(&items, "  - %s\n", n.Text)
		}
	}
	var pairs strings.Builder
	for _, p := range c.Wrong {
		if p.When == "" || !active(p.When, facts) {
			continue
		}
		if positive {
			// The why explains the worse line, so it goes with it.
			fmt.Fprintf(&pairs, "  %s\n", p.Right)
			continue
		}
		fmt.Fprintf(&pairs, "  worse:  %s\n  better: %s\n", p.Wrong, p.Right)
		if p.Why != "" {
			fmt.Fprintf(&pairs, "  why:    %s\n", p.Why)
		}
	}
	var tests strings.Builder
	for _, t := range c.Tests {
		if t.When == "" || !active(t.When, facts) {
			continue
		}
		if positive {
			// The question names the failure it tests for; the passes do not.
			for _, pass := range t.Pass {
				fmt.Fprintf(&tests, "  %s\n", pass)
			}
			continue
		}
		fmt.Fprintf(&tests, "  %s — %s\n", t.Name, t.Question)
		for _, pass := range t.Pass {
			fmt.Fprintf(&tests, "      passes: %s\n", pass)
		}
	}

	if items.Len() == 0 && pairs.Len() == 0 && tests.Len() == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("<voice-refresher>\n")
	if positive {
		// No counts. A rule name is the anti-pattern in the fewest possible
		// characters, and this render exists to hold that out.
		b.WriteString("Mid-session reminder — sentences in the register to write in:\n")
		b.WriteString(pairs.String())
		if tests.Len() > 0 {
			b.WriteString(tests.String())
		}
		b.WriteString("</voice-refresher>\n")
		return b.String()
	}
	if len(observed) > 0 {
		parts := make([]string, 0, len(observed))
		for _, o := range observed {
			parts = append(parts, fmt.Sprintf("%s ×%d", o.Rule, o.N))
		}
		fmt.Fprintf(&b, "Scored in your recent replies: %s.\n\n", strings.Join(parts, ", "))
	}
	if items.Len() > 0 {
		b.WriteString("Never:\n")
		b.WriteString(items.String())
	}
	if pairs.Len() > 0 {
		b.WriteString("\nSame content, written worse and better:\n")
		b.WriteString(pairs.String())
	}
	if tests.Len() > 0 {
		b.WriteString("\nApply before sending:\n")
		b.WriteString(tests.String())
	}
	b.WriteString("</voice-refresher>\n")
	return b.String()
}

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
