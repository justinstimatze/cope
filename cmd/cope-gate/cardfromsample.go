package main

// Writing a card from prose somebody already wrote.
//
// Every card in this repo was written by hand, which is fine for the shipped
// one and wrong for everybody else. A card is a page of prose about register,
// rhythm, and what a paragraph does with a detail, and asking somebody to write
// that page before they have used the tool asks them to describe a voice
// instead of showing one. What they can always do is point at three paragraphs
// they already like.
//
// So this entry does the same thing --author-docs does, one level down: it
// prints a prompt rather than calling a model. cope has no API client and
// wants none — the session invoking it already is one, and a prompt on stdout
// composes with a pipe, a scrollback, and a paste into somebody else's tool.
//
// The idea is humanizer's. Its skill takes a writing sample and follows the
// sample's rhythm and quirks over its own style rules, which is the one thing
// a fixed rule list cannot do. cope's version is durable rather than per-call:
// the sample becomes a card, and the card is what every later turn is scored
// against.

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/justinstimatze/cope/internal/scan"
)

// runCardFromSample prints the prompt for turning a writing sample into a card.
func runCardFromSample(path string) {
	var (
		raw []byte
		err error
	)
	if path == "-" {
		raw, err = os.ReadFile("/dev/stdin")
	} else {
		raw, err = os.ReadFile(path)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
		os.Exit(1)
	}
	sample := strings.TrimSpace(string(raw))
	if sample == "" {
		fmt.Fprintf(os.Stderr, "cope-gate: %s is empty — the card is derived from the sample and there is nothing to derive from\n", path)
		os.Exit(1)
	}
	fmt.Print(buildCardPrompt(sample, path))
}

// cardBlocks is the effigy grammar a card uses, stated as what each block is
// for rather than as a schema. The .effigy notation is effigy's and documented
// there; what a prose card puts in each block is this repo's and is documented
// nowhere else.
const cardBlocks = `VOICE{ kernel: ... peak: ... }
    kernel is the default register: sentence shape, rhythm, what a paragraph
    does with a detail. peak is the one place flair is licensed, and saying
    where it is licensed is how a card avoids being either flat or florid
    throughout.

TRAITS[ hyphenated-adjectives, one line ]
    Five or fewer. Each names a disposition the sample demonstrates, not a
    virtue its writer would claim.

NEVER[ one prohibition per entry, --- between them ]
    At most %d reach the prompt; anything past that is dropped at render and
    reported at load. One entry may be marked CRITICAL: and is rendered first.
    An entry may be gated with @when fact:rule_<id> so it appears only for a
    writer who has actually been making that move.

QUIRKS[ ] TESTS and MES
    QUIRKS are habits, not rules. TEST[ name / question / fail / pass / why ]
    holds a named question with worked examples, which is how a card names a
    move rather than one wording of it. MES[ {{user}}: / {{char}}: ] are whole
    exchanges in the target voice — the highest-value block and the hardest to
    fake, because an exemplar that does not sound like the sample is visible
    immediately.

POSTPROC[ action / pattern / why / id ]
    Regex rules run over the finished text. Only for a phrase that is filler
    every single time it appears. A word that is sometimes exactly right does
    not belong here.

@gate <rule_id> off — <why>
    Declines a built-in rule for this card. The rules are: %s.

@shape <id>: <selector> <predicate> — <why>
    Adds a check of the card's own. Selectors: %s.`

// buildCardPrompt assembles the prompt: the grammar, the sample, and the two
// checks that say whether the result describes the sample or describes a card.
func buildCardPrompt(sample, path string) string {
	var b strings.Builder
	b.WriteString(`<cope-card-from-sample>
Write an effigy card describing the voice of the sample below, so cope can
score later writing against it.

Read the sample first and write down what it actually does before writing any
block: sentence length and how much it varies, where the detail sits in a
paragraph, what it does with a number, how it opens, how it closes, what it
refuses to do. The card is a description of that. It is not a list of good
writing habits, and an entry that would be true of any competent writer is an
entry that measures nothing.

Three rules about what goes in:

1. Derive every entry from something in the sample. If the sample never
   apologizes, do not add a rule against apologizing — it costs prompt budget
   to forbid a move this writer does not make. A card of universal advice is
   the failure mode.
2. Do not invent the writer's opinions, subject matter, or biography. The card
   describes how they write and nothing else.
3. Prefer showing to telling. One MES exchange in the voice is worth three
   TRAITS entries about it.

## The grammar

`)
	fmt.Fprintf(&b, cardBlocks,
		scan.MaxNeverRules(),
		strings.Join(sortedRuleIDs(), ", "),
		"first paragraph, last paragraph, every paragraph, some paragraph, reply; predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask")

	fmt.Fprintf(&b, "\n\n## The sample (%s, %d chars)\n\n", path, len(sample))
	b.WriteString(sample)

	b.WriteString(`

## Then check it, twice

    tools/card2json.py <card>.effigy > <card>.json
    cope-gate --rules <card>.json --check <the sample file>

The sample is by definition in the target voice, so a card that flags it
heavily is describing something else. Read every hit: each one is either a rule
you wrote too broadly, or a place the sample is genuinely inconsistent and you
should say which.

    cope-gate --rules <card>.json --describe

That prints the voice back as a target to recognise. If it does not sound like
the sample to somebody who has read both, the card is wrong and no amount of
extra entries will fix it.
</cope-card-from-sample>
`)
	return b.String()
}

// sortedRuleIDs is the built-in rule set, sorted, for the @gate line. It reads
// scan.RuleIDs rather than the descriptions in this file so a rule added to one
// and not the other cannot make the prompt lie about what is gateable.
func sortedRuleIDs() []string {
	out := make([]string, 0, len(scan.RuleIDs))
	for id := range scan.RuleIDs {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
