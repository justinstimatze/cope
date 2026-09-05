package scan

// Three rules mined from a list this repo did not write, out of four written.
//
// Every other rule here came from one of three places: this machine's own
// transcripts, a published measurement, or a blind read. All four came from
// Wikipedia's "Signs of AI writing", the page WikiProject AI Cleanup maintains,
// reached through github.com/blader/humanizer — a skill that rewrites prose
// against 35 of its patterns. That list is curated by editors reverting AI text
// on an encyclopedia, which is a far larger corpus than anything measured here
// and a completely different reader.
//
// So it is a candidate pool and not a specification, and most of it was left
// where it was. The phrase patterns belong to basanite, which measures a word
// against a baseline instead of banning it. The formatting patterns — em
// dashes, bold, emoji, title case, curly quotes — are declined on purpose:
// bold_label banned the bolded label until 52 blind pairs put bold and bullets
// among the three things that decided a reply for this repo's reader, and it
// was deleted rather than tuned. Lifting them back would re-add what a
// measurement took out.
//
// What is left is the three below, each one a shape rather than a wording:
// something about where a sentence sits, or what it repeats, or what the
// paragraph around it is doing. That is the axis a POSTPROC pattern cannot
// reach and the axis this package is for.
//
// Each was run over the local transcript corpus before shipping, and the rate
// is recorded in MEASUREMENTS.md beside the rules already there. A rule that
// fires on every third reply is a rule nobody leaves switched on, and the
// thresholds below are what the corpus said rather than what the source page
// asserts.
//
// A fourth was written and cut on the same measurement. unanchored_close read
// the closing block for a forecast that names nothing — the future looks
// bright, bodes well, sets you up nicely — which is the source page's generic
// positive ending and is also the one move in this file the card already bans
// in prose. It found 1 hit in 10,256 replies, and widening it from the closing
// sentence to the whole closing block found the same 1. The phrases are in the
// corpus — 62 occurrences in the raw text of 400 of those transcripts — but
// almost never as the close, so the rule was precise and pointless: this repo
// cut verdict_handoff at 38 hits in
// 20.1M characters and this was two orders of magnitude under that. The
// corpus is written under a card that bans the move, so the rate is a floor
// rather than an estimate — and a floor that low is still a floor.

import "strings"

// headingStop are the words a heading can share with its first sentence
// without the sentence having restated it. Short function words carry no
// topic, so counting them would make every heading its own echo.
var headingStop = map[string]bool{
	"the": true, "a": true, "an": true, "and": true, "or": true, "of": true,
	"to": true, "in": true, "on": true, "for": true, "with": true, "by": true,
	"is": true, "are": true, "was": true, "were": true, "it": true, "this": true,
	"that": true, "what": true, "how": true, "why": true, "when": true,
}

func headingTopic(s string) []string {
	var out []string
	for _, w := range words(strings.TrimLeft(s, "# ")) {
		lw := strings.ToLower(w)
		if len(lw) >= 4 && !headingStop[lw] {
			out = append(out, lw)
		}
	}
	return out
}

// EchoedHeading flags a heading whose first sentence says it again.
//
// The heading has already told the reader the topic. A first sentence that
// repeats every content word of it spends a line saying nothing, and it is the
// shape a model reaches for when it is filling a template rather than writing
// under a heading it chose.
//
// Both the heading and the body are needed, which is why this rule splits the
// text itself: paragraphs and rawParagraphs both drop headings outright, so no
// existing rule can see one.
func EchoedHeading(text string) []Violation {
	return echoedHeading(Redact(text))
}

func echoedHeading(redacted string) []Violation {
	blocks := paraSplit.Split(redacted, -1)
	var out []Violation
	for i, b := range blocks {
		if !strings.HasPrefix(strings.TrimSpace(b), "#") {
			continue
		}
		head, body, ok := headingAndBody(blocks, i)
		if !ok {
			continue
		}
		topic := headingTopic(head)
		sents := splitSentences(strings.TrimSpace(body))
		if len(topic) < 2 || len(sents) == 0 {
			continue
		}
		if !restates(topic, sents[0]) {
			continue
		}
		out = append(out, Violation{
			RuleID:  "echoed_heading",
			Action:  "warn",
			Why:     "the first sentence repeats every content word of the heading above it — let the heading do that work",
			Matched: strings.TrimSpace(sents[0]),
			Context: strings.TrimSpace(head),
		})
	}
	return out
}

// headingAndBody pairs a heading with the prose under it.
//
// A heading and its body are one block whenever a single newline separates
// them, so the rest of this block is taken first and the next block only when
// the heading stands alone. A heading followed by another heading has no body.
func headingAndBody(blocks []string, i int) (head, body string, ok bool) {
	b := strings.TrimSpace(blocks[i])
	head, body, split := strings.Cut(b, "\n")
	if split && strings.TrimSpace(body) != "" {
		return head, body, true
	}
	if i+1 >= len(blocks) {
		return "", "", false
	}
	body = blocks[i+1]
	if strings.HasPrefix(strings.TrimSpace(body), "#") {
		return "", "", false
	}
	return b, body, true
}

// restates reports whether every content word of the heading appears in the
// sentence below it.
func restates(topic []string, first string) bool {
	lower := strings.ToLower(first)
	for _, w := range topic {
		if !strings.Contains(lower, w) {
			return false
		}
	}
	return true
}

// RepeatedOpening flags several sentences in one reply starting the same way.
//
// cross_turn_repeat reads the session window for a construction reused across
// turns; nothing read one reply against itself. Three sentences opening on the
// same two words is the within-reply form, and it reads as a template being
// filled even when each sentence is individually fine.
//
// Deliberate anaphora exists and is not this. The threshold is what separates
// them: two is a rhythm, and the corpus rate at three is what decided the
// number rather than the source page, which gives none.
func RepeatedOpening(text string) []Violation {
	return repeatedOpening(paragraphs(Redact(text)))
}

func repeatedOpening(ps []string) []Violation {
	type opener struct {
		count int
		first string
	}
	seen := map[string]*opener{}
	var order []string
	for _, p := range ps {
		for _, s := range rawSentences(p) {
			ws := words(s)
			if len(ws) < 4 {
				continue
			}
			key, ok := openingKey(s)
			if !ok {
				continue
			}
			o, ok := seen[key]
			if !ok {
				o = &opener{first: strings.TrimSpace(s)}
				seen[key] = o
				order = append(order, key)
			}
			o.count++
		}
	}
	var out []Violation
	for _, key := range order {
		if seen[key].count < 3 {
			continue
		}
		out = append(out, Violation{
			RuleID:  "repeated_opening",
			Action:  "warn",
			Why:     "three or more sentences open on the same two words, which reads as a template rather than a rhythm",
			Matched: seen[key].first,
			Context: key,
		})
	}
	return out
}

// rawSentences splits like splitSentences and does not trim.
//
// The whitespace is the evidence openingKey needs: a trimmed sentence that
// began with a code span is indistinguishable from one that began with its
// second word, because Redact left blanks and TrimSpace removed them.
func rawSentences(p string) []string {
	locs := sentenceEnd.FindAllStringIndex(p, -1)
	var out []string
	prev := 0
	for _, l := range locs {
		out = append(out, p[prev:l[1]])
		prev = l[1]
	}
	if prev < len(p) {
		out = append(out, p[prev:])
	}
	return out
}

// openingKey is a sentence's first two words, lowercased, when they really are
// the first two words.
//
// Redact blanks code spans in place, so a sentence written as "`cope-gate
// --refresher` runs on every prompt" arrives here as a run of spaces followed
// by "runs on every prompt", and taking words[0] and words[1] silently reads
// across the hole. The key becomes a pair that is not adjacent in the reply and
// the quoted match carries a stretch of blanks — which is how demo/README.precise.md
// reported three sentences opening on "runs which". A sentence whose opening is
// a code span has no prose opening to compare, so it is skipped.
func openingKey(s string) (string, bool) {
	locs := wordRe.FindAllStringIndex(s, 2)
	if len(locs) < 2 {
		return "", false
	}
	if strings.Contains(s[:locs[0][0]], "  ") {
		return "", false
	}
	if locs[1][0]-locs[0][1] > 2 {
		return "", false
	}
	return strings.ToLower(s[locs[0][0]:locs[0][1]] + " " + s[locs[1][0]:locs[1][1]]), true
}

// FragmentRun flags three or more fragments in a row.
//
// One fragment is emphasis and this repo's own house style is full of them.
// Three consecutive is the staccato that judges read as generated, and the
// run length is what the source page itself says to require: flag dramatic
// fragments only when several appear together.
//
// A card wanting the clipped register declines it with @gate. None of the demo
// cards need to: run against the READMEs rendered under precise and caveman,
// the two terse ones, the rule found nothing, so the run length is above the
// register rather than inside it.
func FragmentRun(text string) []Violation {
	return fragmentRun(paragraphs(Redact(text)))
}

func fragmentRun(ps []string) []Violation {
	var out []Violation
	for _, p := range ps {
		if start, ok := firstFragmentRun(p); ok {
			out = append(out, Violation{
				RuleID:  "fragment_run",
				Action:  "warn",
				Why:     "three fragments in a row — one is emphasis, a run of them is the staccato that reads as generated",
				Matched: start,
				Context: strings.TrimSpace(strings.ReplaceAll(p, "\n", " ")),
			})
		}
	}
	return out
}

// firstFragmentRun returns the first of three consecutive fragments in a
// paragraph, if there are three.
func firstFragmentRun(p string) (string, bool) {
	run, start := 0, ""
	for _, s := range splitSentences(p) {
		ws := words(s)
		if len(ws) == 0 {
			continue
		}
		if !isFragment(ws) {
			run = 0
			continue
		}
		if run == 0 {
			start = strings.TrimSpace(s)
		}
		if run++; run == 3 {
			return start, true
		}
	}
	return "", false
}

// isFragment is a noun phrase standing where a sentence belongs: short, and
// carrying no finite verb.
func isFragment(ws []string) bool {
	return len(ws) <= 5 && !hasFiniteVerb(ws)
}
