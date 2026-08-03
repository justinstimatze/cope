package scan

import (
	"testing"

	"github.com/justinstimatze/cope/card"
)

// The flip rule matched this repo's own documentation of the flip rule in the
// originating session. Redaction is what stops the gate firing hardest during
// the conversations where the rules are being discussed.
func TestRedactBlanksQuotedAndFencedSpans(t *testing.T) {
	in := "before \"that's not a decline, it's a plateau\" after"
	got := Redact(in)
	if len(got) != len(in) {
		t.Fatalf("redaction changed length: %d != %d (offsets must stay aligned)", len(got), len(in))
	}
	if want := "before"; got[:len(want)] != want {
		t.Errorf("prefix outside the quote was altered: %q", got[:len(want)])
	}
	for _, sub := range []string{"decline", "plateau"} {
		if contains(got, sub) {
			t.Errorf("quoted text %q survived redaction", sub)
		}
	}

	fenced := "text\n```\nnot a thing, it's another\n```\ntail"
	if contains(Redact(fenced), "another") {
		t.Error("fenced block survived redaction")
	}
}

// Length is the wrong test for a labelled opening. "Basanite measures words" is
// short and is a real sentence; "The audit half" is short and is a label.
func TestLabelledOpeningNeedsAFragmentNotJustABriefSentence(t *testing.T) {
	sentence := "Basanite measures words. " + filler()
	if v := LabelledOpening(sentence); len(v) != 0 {
		t.Errorf("flagged a complete sentence as a label: %+v", v)
	}

	fragment := "The audit half. " + filler()
	if v := LabelledOpening(fragment); len(v) == 0 {
		t.Error("missed a verbless label opening")
	}

	ordinal := "Second, the detectors, in two families. " + filler()
	if v := LabelledOpening(ordinal); len(v) == 0 {
		t.Error("missed the ordinal form of the label opening")
	}
}

func TestClauseSymmetrySkipsMarkdownLinkRuns(t *testing.T) {
	links := "The most common is a banned-word list, and so on " +
		"([one](https://example.com/aaa), [two](https://example.com/bbb)). " + filler()
	for _, v := range ClauseSymmetry(links) {
		if contains(v.Matched, "https://") {
			t.Errorf("flagged a markdown link run as balanced prose: %q", v.Matched)
		}
	}

	balanced := "The flip hands you the correction pre-formed, and the distinction hands you the correction pre-assigned. " + filler()
	if len(ClauseSymmetry(balanced)) == 0 {
		t.Error("missed a balanced two-beat")
	}
}

func TestParagraphUniformityNeedsFourProseParagraphs(t *testing.T) {
	if v := ParagraphUniformity(filler()+"\n\n"+filler(), 0.35); len(v) != 0 {
		t.Errorf("fired on two paragraphs: %+v", v)
	}
	even := filler() + "\n\n" + filler() + "\n\n" + filler() + "\n\n" + filler()
	if len(ParagraphUniformity(even, 0.35)) == 0 {
		t.Error("missed four paragraphs of identical length")
	}
}

func TestAskNotLastFlagsAMidReplyAskButNotATrailingOne(t *testing.T) {
	trailing := filler() + "\n\n" + "Should I go ahead and run the migration?"
	if v := AskNotLast(trailing); len(v) != 0 {
		t.Errorf("flagged an ask that was already the last paragraph: %+v", v)
	}

	buried := "Should I go ahead and run the migration? " + filler() + "\n\n" + filler()
	if v := AskNotLast(buried); len(v) == 0 {
		t.Error("missed an ask buried before two more paragraphs")
	}

	if v := AskNotLast(filler()); len(v) != 0 {
		t.Errorf("flagged a single paragraph with no ask: %+v", v)
	}
}

// Asks are short by nature, and the 12-word prose floor was deleting exactly
// the block the rule needed to see: "Build both?" scored zero while the same
// ask padded to fourteen words scored one. DanglingEnd had already been given
// rawParagraphs for this reason; AskNotLast had not.
func TestAskNotLastSeesAsksBelowTheProseWordFloor(t *testing.T) {
	short := filler() + "\n\nBuild both?\n\n" + filler()
	v := AskNotLast(short)
	if len(v) != 1 {
		t.Fatalf("got %d violations for a two-word mid-reply ask, want 1: %+v", len(v), v)
	}
	if v[0].Matched != "Build both?" {
		t.Errorf("matched %q, want the short ask itself", v[0].Matched)
	}

	trailing := filler() + "\n\n" + filler() + "\n\nBuild both?"
	if v := AskNotLast(trailing); len(v) != 0 {
		t.Errorf("flagged a short ask that was already last: %+v", v)
	}
}

// Every rule reads the same redacted text, and Check derives it once. The
// guard is that sharing it changed nothing about what gets reported.
func TestCheckAgreesWithRunningTheRulesSeparately(t *testing.T) {
	c, _, err := ParseCard([]byte(`{"rules":[{"id":"flip","action":"warn","pattern":"\\bnot [a-z]+, it's\\b"}]}`))
	if err != nil {
		t.Fatal(err)
	}

	text := "It is not slow, it's blocked on the disk here.\n\n" +
		filler() + "\n\nShould I go ahead and run it?\n\n" + filler() +
		"\n\nThe migration is still unverified on staging."

	want := append(c.Regex(text), Shapes(text, 0.35)...)
	got := c.Check(text, 0.35)
	if len(got) != len(want) {
		t.Fatalf("Check found %d violations, running the rules separately found %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("violation %d differs:\n got %+v\nwant %+v", i, got[i], want[i])
		}
	}
	if len(got) == 0 {
		t.Error("the fixture triggered nothing, so this proves nothing")
	}
}

// strip is legal in effigy's POSTPROC grammar and unimplementable in a hook
// that sees the reply after it is written. It used to behave as warn silently.
func TestStripActionIsReportedAndDowngraded(t *testing.T) {
	c, warnings, err := ParseCard([]byte(`{"rules":[{"id":"cut","action":"strip","pattern":"x"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want one naming the strip rule: %v", len(warnings), warnings)
	}
	if !contains(warnings[0].Error(), "strip") || !contains(warnings[0].Error(), "cut") {
		t.Errorf("warning does not name the rule and its action: %v", warnings[0])
	}
	if c.Rules[0].Action != "warn" {
		t.Errorf("action is %q, want it downgraded to warn", c.Rules[0].Action)
	}
}

// The fixtures here are condensed from real Opus 5 end-of-turns in the July
// 2026 publicai sessions — the corpus where "what specifically should we do
// next" appeared 14 times as a user follow-up. Each dangling fixture preceded
// such a re-ask; each clean fixture is the Opus 4.8 ending style from the same
// projects days earlier, which never drew one.
func TestDanglingEndFlagsOpenProblemsWithNoDecisionPoint(t *testing.T) {
	dangling := filler() + "\n\n" +
		"The live probe could not reach the preview host, a TimeoutError, so the " +
		"deployed config went unverified on this run. Might be transient. Worth a look."
	if v := DanglingEnd(dangling); len(v) == 0 || v[0].RuleID != "dangling_end" {
		t.Errorf("missed an ending that names an open problem with no ask or offer: %+v", v)
	}

	observed := filler() + "\n\n" +
		"One trial hit the 300s wall clock last night, and three trials reported totals " +
		"above the token cap, which I have a guess about and have not checked."
	if v := DanglingEnd(observed); len(v) == 0 {
		t.Error("missed the observes-without-acting-or-asking ending")
	}
}

func TestDanglingEndAcceptsClosureOfferAndTrailingAsk(t *testing.T) {
	closure := filler() + "\n\nLoop closed, nothing for you to act on."
	if v := DanglingEnd(closure); len(v) != 0 {
		t.Errorf("flagged an explicit all-clear: %+v", v)
	}

	offer := filler() + "\n\n" +
		"The judge rubric is still broken in checkpoint-compare.ts. Next I'd fix the " +
		"rubric before quoting either column. Say the word and I'll start."
	if v := DanglingEnd(offer); len(v) != 0 {
		t.Errorf("flagged an open problem that ends in a concrete offer: %+v", v)
	}

	shortAsk := filler() + "\n\nBuild both?"
	if v := DanglingEnd(shortAsk); len(v) != 0 {
		t.Errorf("flagged a reply whose final block is a two-word ask — the word floor bug: %+v", v)
	}

	noProblems := filler() + "\n\n" + filler()
	if v := DanglingEnd(noProblems); len(v) != 0 {
		t.Errorf("flagged a reply that names no open problem: %+v", v)
	}
}

func TestDanglingEndFlagsObservationsBuryingTheAsk(t *testing.T) {
	buried := filler() + "\n\n" +
		"Want me to log phase 1 on the ticket? The badge is honest on all three counts " +
		"now, and the Doppler key still needs removing on your end."
	v := DanglingEnd(buried)
	if len(v) == 0 || v[0].RuleID != "buried_decision" {
		t.Errorf("missed an open problem landing after the ask: %+v", v)
	}

	clean := "The badge is honest on all three counts now, and the Doppler key still " +
		"needs removing on your end.\n\nWant me to log phase 1 on the ticket?"
	if v := DanglingEnd(clean); len(v) != 0 {
		t.Errorf("flagged the corrected order, problem first then ask: %+v", v)
	}
}

// Every fixture below is condensed from a real ending in the 60-session sample
// forked_end was measured on. The two accept cases were false positives in the
// first draft, which counted a leading "or" as a second decision and read a
// question inside a list item as an ask.
func TestForkedEndFlagsTwoThingsToActOn(t *testing.T) {
	two := filler() + "\n\n" +
		"Want me to squash these two commits into one before merge? The other thing " +
		"is the search question — should I dispatch it to the team channel?"
	v := ForkedEnd(two)
	if len(v) == 0 || v[0].RuleID != "forked_end" {
		t.Fatalf("missed two independent asks in the tail: %+v", v)
	}
	if v[0].Context == v[0].Matched {
		t.Errorf("context and match should be the two different decisions, both were %q", v[0].Matched)
	}

	askAndOffer := filler() + "\n\n" +
		"Route it into the existing instance, or a fresh one?\n\n" +
		"Say the word and I'll start it."
	if v := ForkedEnd(askAndOffer); len(v) == 0 {
		t.Error("missed an ask in one block and an unrelated offer in the next")
	}
}

func TestForkedEndAcceptsOneDecisionHoweverPrinted(t *testing.T) {
	alternative := filler() + "\n\n" +
		"Want me to land the safe one-liners on this branch now?\n" +
		"Or hold and let you triage the full list first?"
	if v := ForkedEnd(alternative); len(v) != 0 {
		t.Errorf("flagged one choice printed as two sentences, the second opening on 'or': %+v", v)
	}

	enumerated := filler() + "\n\n" +
		"- Search: does free-typing the funding question trigger search?\n" +
		"- Safety: does a toxic probe decline on prod?\n\n" +
		"Want me to finish those two live, or stop here?"
	if v := ForkedEnd(enumerated); len(v) != 0 {
		t.Errorf("flagged questions enumerated in list items alongside the one real ask: %+v", v)
	}

	closureAndOffer := filler() + "\n\n" +
		"Oracle left as-is per your call, nothing for you to act on."
	if v := ForkedEnd(closureAndOffer); len(v) != 0 {
		t.Errorf("flagged an all-clear: %+v", v)
	}

	tagged := filler() + "\n\nWant me to commit the source-tree changes? Your call."
	if v := ForkedEnd(tagged); len(v) != 0 {
		t.Errorf("flagged one ask followed by a bare deference tag: %+v", v)
	}

	tagWithContent := filler() + "\n\n" +
		"I'll run the sweep next. Let me know if you would rather I check staging first."
	if v := ForkedEnd(tagWithContent); len(v) == 0 {
		t.Error("missed a real second decision that happens to open on 'let me know'")
	}

	single := filler() + "\n\nBuild both?"
	if v := ForkedEnd(single); len(v) != 0 {
		t.Errorf("flagged a single two-word ask: %+v", v)
	}

	none := filler() + "\n\n" + filler()
	if v := ForkedEnd(none); len(v) != 0 {
		t.Errorf("flagged a tail with no decision in it at all: %+v", v)
	}
}

// The ending is the last two blocks, the same window dangling_end reads. An ask
// further up is ask_not_last's to report, and counting it here would flag every
// reply that asked something mid-way and then closed cleanly.
func TestForkedEndReadsOnlyTheTail(t *testing.T) {
	early := "Want me to check the migration first?\n\n" + filler() + "\n\n" + filler() +
		"\n\nSay the word and I'll start."
	if v := ForkedEnd(early); len(v) != 0 {
		t.Errorf("flagged an ask three blocks above the ending: %+v", v)
	}
}

// filler is a prose paragraph long enough to clear the 12-word floor.
func filler() string {
	return "This sentence exists so the paragraph clears the minimum word count the scanner requires before it will consider the text at all."
}

func contains(hay, needle string) bool {
	for i := 0; i+len(needle) <= len(hay); i++ {
		if hay[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// Ported from effigy v0.3.1 (stope phase-probe finding): an inline example
// marker primes the model on the example text regardless of the rule's
// negative label, so it gets cut rather than rendered.
func TestStripInlineExamplesCutsAtTheMarker(t *testing.T) {
	got := stripInlineExamples(`Never goes coy. WRONG: "Maybe." RIGHT: "No."`)
	if got != "Never goes coy" {
		t.Errorf("got %q, want the rule text with the example and trailing period cut", got)
	}
	plain := "Never uses academic language"
	if got := stripInlineExamples(plain); got != plain {
		t.Errorf("stripped a rule with no marker: %q", got)
	}
	lower := "Never claims something is wrong: it needs a citation"
	if got := stripInlineExamples(lower); got != lower {
		t.Errorf("stripped on a lowercase 'wrong:', which should be untouched: %q", got)
	}
}

// Ported from effigy v0.3.1/v0.3.2: CRITICAL rules sort first, the list caps
// at maxNeverRules, and whatever survives has its inline examples stripped.
func TestPrioritizeNeverCapsAndSortsCriticalFirst(t *testing.T) {
	var never []Gated
	for range maxNeverRules + 2 {
		never = append(never, Gated{Text: "regular rule"})
	}
	never = append(never, Gated{Text: "CRITICAL: the one that must survive"})

	got, dropped := prioritizeNever(never, nil)
	if len(got) != maxNeverRules {
		t.Fatalf("got %d rules, want the cap of %d", len(got), maxNeverRules)
	}
	if got[0].Text != "CRITICAL: the one that must survive" {
		t.Errorf("critical rule was not sorted first: %+v", got[0])
	}
	if len(dropped) != len(never)-maxNeverRules {
		t.Errorf("got %d dropped rules, want %d — a discarded rule must be reported, not silent",
			len(dropped), len(never)-maxNeverRules)
	}
}

// The failure this guards is the one that already happened: commit 584c898
// added an eleventh NEVER rule and the ten-rule cap together, so the rule was
// truncated away by its own commit and sat in the card for two days reading as
// covered. LoadCard now reports the overflow, and this fails the build if the
// shipped card ever carries a rule that cannot reach the prompt.
func TestShippedCardShipsEveryNeverRule(t *testing.T) {
	c, skipped, err := ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatalf("parse the embedded card: %v", err)
	}
	if over := c.NeverOverflow(); len(over) != 0 {
		t.Errorf("card/rules.json has %d NEVER rule(s) past the %d-rule budget that never render: %+v",
			len(over), maxNeverRules, over)
	}
	for _, s := range skipped {
		t.Errorf("shipped card reports a skipped element: %v", s)
	}
}

// The card that gets injected and the card that gets checked must come from
// the same file, or the gate enforces rules the writer was never given.
func TestRenderCarriesTheAuthoredCard(t *testing.T) {
	c := &Card{
		Theme:  "aim",
		Voice:  Voice{Kernel: "plain sentences", Peak: "flair only when earned"},
		Never:  []Gated{{Text: "no labels"}, {Text: "no ordinals"}},
		Wrong:  []Pair{{Wrong: "worse phrasing", Right: "better phrasing", Why: "the reason it is worse"}},
		Mes:    []string{"{{user}}: a question\n{{char}}: a reply in the wanted voice"},
		Tests:  []Test{{Name: "SWAP", Question: "swappable?", Fail: []string{"generic"}, Pass: []string{"specific"}}},
		Quirks: []string{"names the file"},
		Traits: []string{"flat-self-report"},
		Rules:  []Rule{{ID: "flip", Pattern: `not .*, it's`}},
	}
	got := c.Render()
	for _, want := range []string{
		"aim", "plain sentences", "flair only when earned",
		"worse phrasing", "better phrasing", "the reason it is worse",
		"a reply in the wanted voice",
		"no labels", "no ordinals",
		"SWAP", "swappable?", "generic", "specific",
		"names the file", "flat-self-report",
	} {
		if !contains(got, want) {
			t.Errorf("rendered card is missing %q", want)
		}
	}
	// Showing the patterns invites writing around them.
	if contains(got, `not .*, it's`) {
		t.Error("rendered card leaked a POSTPROC regex")
	}
}

// Describe is what the discrimination reader holds while comparing two
// replies. It has to carry the voice and none of the machinery: a reader given
// the prohibitions and the worked pairs is marking prose against a checklist,
// which is a different task with a different answer.
func TestDescribeCarriesTheVoiceAndNothingElse(t *testing.T) {
	c := &Card{
		Theme: "the phrasing should cost the reader nothing",
		Voice: Voice{Kernel: "plain sentences that vary in length", Peak: "flair only when earned"},
		Never: []Gated{{Text: "never opens on a bold label"}},
		Wrong: []Pair{{Wrong: "worse phrasing", Right: "better phrasing", Why: "the reason"}},
		Tests: []Test{{Name: "SWAP", Question: "swappable?"}},
		Mes:   []string{"{{char}}: a reply in the wanted voice"},
		Rules: []Rule{{ID: "flip", Pattern: `not .*, it's`}},
	}
	got := c.Describe()
	for _, want := range []string{
		"the phrasing should cost the reader nothing",
		"plain sentences that vary in length",
		"flair only when earned",
	} {
		if !contains(got, want) {
			t.Errorf("description is missing %q", want)
		}
	}
	for _, leak := range []string{
		"never opens on a bold label", "worse phrasing", "better phrasing",
		"SWAP", "swappable?", "a reply in the wanted voice", `not .*, it's`,
	} {
		if contains(got, leak) {
			t.Errorf("description leaked the machinery: %q", leak)
		}
	}
}

// A card with no VOICE block describes nothing, and the harness has to be able
// to see that rather than showing a reader an empty panel.
func TestDescribeIsEmptyWithoutAVoice(t *testing.T) {
	c := &Card{Never: []Gated{{Text: "never opens on a bold label"}}}
	if got := c.Describe(); got != "" {
		t.Errorf("a card with no voice described itself as %q", got)
	}
}

func TestRenderRefresherIsCompactAndFromTheCard(t *testing.T) {
	c := &Card{Tests: []Test{{
		Name:     "CONTINUE TEST",
		Question: `Could the reader answer the final line with the single word "continue"?`,
		Pass:     []string{"Loop closed, nothing for you to act on."},
	}}}
	out := c.RenderRefresher()
	for _, want := range []string{"<voice-refresher>", "continue", "Loop closed"} {
		if !contains(out, want) {
			t.Errorf("refresher missing %q:\n%s", want, out)
		}
	}
	if contains(out, "voice-card") {
		t.Error("refresher must not be the full card")
	}
	if (&Card{}).RenderRefresher() != "" {
		t.Error("a card with no CONTINUE TEST must render an empty refresher")
	}
}

// The refresher renders from a test looked up by name, which is the one
// content-addressed seam in the pipeline: rename the test in the card and the
// refresher silently vanishes. This pins the shipped card against that.
func TestShippedCardHasARefresher(t *testing.T) {
	c, _, err := ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatalf("parse the embedded card: %v", err)
	}
	out := c.RenderRefresher()
	if out == "" {
		t.Fatal("shipped card renders an empty refresher — was the CONTINUE TEST renamed?")
	}
	if len(out) > 900 {
		t.Errorf("refresher is %d chars; past ~900 it stops being a reminder and starts being a second card", len(out))
	}
}

// A markdown table is balanced by construction — every row repeats the
// column's vocabulary at matching length. Checking this repo's README turned
// one six-row layout table into two clause_symmetry hits.
func TestTablesAreNotProse(t *testing.T) {
	table := "| Path | What |\n| --- | --- |\n" +
		"| `card/claude_voice.effigy` | the card, in effigy notation |\n" +
		"| `card/rules.json` | generated from it; embedded in the binary |\n" +
		"| `cmd/cope-gate/` | the hook binary |\n" +
		"| `internal/scan/` | regex and shape rules, and the card renderer |\n"

	if v := ClauseSymmetry(table); len(v) != 0 {
		t.Errorf("flagged table rows as balanced prose: %+v", v)
	}
	if v := ParagraphUniformity(table+"\n\n"+table, 0.35); len(v) != 0 {
		t.Errorf("counted table blocks as prose paragraphs: %+v", v)
	}
	// Prose next to a table still gets read.
	mixed := table + "\n\n" + "The flip hands you the correction pre-formed, and the distinction hands you the correction pre-assigned. " + filler()
	if len(ClauseSymmetry(mixed)) == 0 {
		t.Error("a table in the document silenced the rule for the prose around it")
	}
}

// MES is effigy's positive-exemplar block and the card carried none until
// 2026-07-28. Anthropic's Opus 5 guidance puts positive examples of the wanted
// style above instructions about what not to do, so these render before the
// NEVER list rather than after it.
func TestRenderShowsWholeRepliesBeforeTheProhibitions(t *testing.T) {
	c := &Card{
		Never: []Gated{{Text: "no labels"}},
		Mes:   []string{"{{user}}: did it work?\n{{char}}: Yes, 41 rows moved."},
	}
	got := c.Render()

	// Ali:Chat placeholders are for driving an NPC; a card injected into a live
	// session should not ask the reader who {{char}} is.
	for _, leak := range []string{"{{char}}", "{{user}}"} {
		if contains(got, leak) {
			t.Errorf("rendered card leaked the placeholder %s", leak)
		}
	}
	for _, want := range []string{"user: did it work?", "you: Yes, 41 rows moved."} {
		if !contains(got, want) {
			t.Errorf("rendered card is missing %q", want)
		}
	}
	mes, never := indexOf(got, "Whole replies in this voice"), indexOf(got, "Never:")
	if mes < 0 || never < 0 {
		t.Fatalf("missing a section: mes=%d never=%d", mes, never)
	}
	if mes > never {
		t.Error("prohibitions render before the positive examples; the ordering is the point")
	}
}

func TestShippedCardCarriesExemplarsAndReasons(t *testing.T) {
	c, _, err := ParseCard(card.RulesJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Mes) == 0 {
		t.Error("shipped card has no MES entries — did card2json.py stop emitting them?")
	}
	for i, p := range c.Wrong {
		if p.Why == "" {
			t.Errorf("wrong pair %d has no why; effigy parses one and the pipeline used to drop it", i)
		}
	}
}

func indexOf(hay, needle string) int {
	for i := 0; i+len(needle) <= len(hay); i++ {
		if hay[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

// The three field rules. Each fires on the move Zvi's Opus 5 roundup names and
// stays quiet on the sentence it is most likely to be confused with, because
// that neighbour is behaviour this repo asks for elsewhere: correct plainly,
// explain the bug, say the reply is long when the reader asked how long it is.
func TestApologyFlagsContritionNotCorrection(t *testing.T) {
	for _, s := range []string{
		"I apologize — the count was 38, not 42.",
		"Sorry about that, the path was wrong.",
		"My apologies for the confusion.",
		"I'm really sorry, that was my read and it was wrong.",
	} {
		if got := Apology(s); len(got) == 0 {
			t.Errorf("no apology flagged in %q", s)
		}
	}
	for _, s := range []string{
		"The count was 38, not 42. I misread the column header.",
		"That was wrong. The correct path is internal/scan/scan.go.",
		"The build is in a sorry state and the cache is stale.",
		"You said sorry was banned, so the rule now covers it.",
	} {
		if got := Apology(s); len(got) > 0 {
			t.Errorf("plain correction flagged as an apology in %q: %+v", s, got)
		}
	}
}

func TestSelfPostmortemNeedsAFirstPersonSubject(t *testing.T) {
	for _, s := range []string{
		"I made two errors this session. And I owe you an explanation as to why.",
		"Let me explain what I got wrong in the first pass.",
		"In the interest of full transparency, the first attempt never compiled.",
		"An honest disclosure: the fixture was never run.",
	} {
		if got := SelfPostmortem(s); len(got) == 0 {
			t.Errorf("no post-mortem flagged in %q", s)
		}
	}
	for _, s := range []string{
		"Let me explain why the build went wrong: the cache key omitted the model.",
		"The test made an assertion about the old contract, which is why it failed.",
		"Two errors surfaced in CI, both in the fixture loader.",
		"Here is where the parser went wrong on line 40.",
		// Both flagged on a real backfill by an earlier draft. The first is the
		// plain correction the card asks for; the second is a caveat about file
		// permissions that has nothing to do with the writer.
		"My mistake was locating the safety in the sandbox when it lives in the schema.",
		"One caveat for honesty: the directory is world-visible at 664.",
		// The backfill's only self_postmortem match, and it was describing
		// somebody else's remark rather than announcing one.
		"That's an honest disclosure but it's a real scoping question.",
	} {
		if got := SelfPostmortem(s); len(got) > 0 {
			t.Errorf("ordinary debugging prose flagged as a post-mortem in %q: %+v", s, got)
		}
	}
}

func TestAnnouncedLengthFlagsTheWarningNotTheLength(t *testing.T) {
	for _, s := range []string{
		"Bear with me, there are three moving parts.",
		"You might want to grab a cup of tea because this takes a while to lay out.",
		"This is going to be a bit long.",
		"There is a lot to unpack in the gating change.",
	} {
		if got := AnnouncedLength(s); len(got) == 0 {
			t.Errorf("no announcement flagged in %q", s)
		}
	}
	for _, s := range []string{
		"The diff is long; most of it is the generated card.",
		"This is a lot of state for one struct.",
		"There is a lot going on in the refresher path.",
		"The card is 14k characters and the handoff variant is 5k.",
		// Flagged by an earlier draft: an offer to talk the reader through
		// operator steps is the help being asked for, not a preamble.
		"I committed locally but did not push — say the word and I'll walk you through it.",
		// All three of this rule's backfill matches were this shape: "dense"
		// describing a scene, not the length of the reply describing it.
		"The humming station, seal fractures, bare earth circle — this is dense with leads.",
	} {
		if got := AnnouncedLength(s); len(got) > 0 {
			t.Errorf("a plain statement of size flagged as an announcement in %q: %+v", s, got)
		}
	}
}
