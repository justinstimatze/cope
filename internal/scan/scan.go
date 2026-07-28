// Package scan applies a card's rules to a drafted response.
//
// Two rule families, because they fail differently. Regex rules come from the
// card's POSTPROC block and catch fixed surface forms. Shape rules are built in
// and catch structural habits that no regex expresses: paragraphs that all come
// out the same length, sentences whose two halves are balanced against each
// other, paragraphs that open with a label and then unpack it, and a question
// or request for the user that lands mid-reply instead of at the end.
//
// Measured on this repo's first 36-turn transcript: the one structural rule
// caught all 8 of its instances, and the two phrase-template rules caught 0
// across 65k characters. Shape is where the budget goes.
package scan

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
)

// Rule is one regex rule loaded from the card's POSTPROC block.
type Rule struct {
	ID      string `json:"id"`
	Action  string `json:"action"` // reject | strip | warn
	Pattern string `json:"pattern"`
	Why     string `json:"why"`

	re *regexp.Regexp
}

// Card is the JSON emitted by tools/card2json.py. Rules drive the gate; the
// remaining fields drive Render, which puts the card into context at
// SessionStart. Both halves come from the same .effigy source so the thing
// being checked and the thing being asked for cannot drift apart.
type Card struct {
	CardID string `json:"card_id"`
	Source string `json:"source"`
	Rules  []Rule `json:"rules"`

	Theme  string   `json:"theme"`
	Voice  Voice    `json:"voice"`
	Traits []string `json:"traits"`
	Never  []string `json:"never"`
	Quirks []string `json:"quirks"`
	Tests  []Test   `json:"tests"`
	Wrong  []Pair   `json:"wrong"`
}

// Voice is the VOICE block: kernel is the default register, peak is where
// flair is licensed.
type Voice struct {
	Kernel string `json:"kernel"`
	Peak   string `json:"peak"`
}

// Test is a reasoning check with worked examples — the mechanism for naming a
// move rather than a phrasing, so a banned form cannot leak into a variant.
type Test struct {
	Name     string   `json:"name"`
	Question string   `json:"question"`
	Fail     []string `json:"fail"`
	Pass     []string `json:"pass"`
	Why      string   `json:"why"`
}

// Pair is a WRONG/RIGHT rewrite of the same content.
type Pair struct {
	Wrong string `json:"wrong"`
	Right string `json:"right"`
}

// Violation is one hit, ready to print or log.
type Violation struct {
	RuleID  string `json:"rule_id"`
	Action  string `json:"action"`
	Why     string `json:"why"`
	Matched string `json:"matched"`
	Context string `json:"context"`
}

// LoadCard reads a rules JSON file. Used by --rules for a card the operator is
// iterating on; the shipped card is embedded and goes through ParseCard.
func LoadCard(path string) (*Card, []error, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	c, warnings, err := ParseCard(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return c, warnings, nil
}

// ParseCard decodes a card and compiles every pattern. effigy applies
// re.IGNORECASE to POSTPROC patterns, so the (?i) flag is added here to match.
//
// The second return is warnings, not failures: a pattern that does not compile
// under RE2 is reported and skipped rather than failing the whole gate, because
// one bad rule should not disable the others. Nothing here is fatal, so a card
// that is merely odd still produces a working gate.
func ParseCard(raw []byte) (*Card, []error, error) {
	var c Card
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, nil, fmt.Errorf("parse card: %w", err)
	}
	var skipped []error
	kept := c.Rules[:0]
	for _, r := range c.Rules {
		re, err := regexp.Compile("(?i)" + r.Pattern)
		if err != nil {
			skipped = append(skipped, fmt.Errorf("rule %s: %w", r.ID, err))
			continue
		}
		// effigy's POSTPROC grammar allows strip, which rewrites the text
		// before it ships. A Stop hook sees the reply after it is written and
		// cannot edit it, so strip is unimplementable here and silently
		// behaved as warn. Say so instead.
		if r.Action == "strip" {
			skipped = append(skipped, fmt.Errorf(
				"rule %s has action strip, which a post-generation gate cannot apply — treating it as warn", r.ID))
			r.Action = "warn"
		}
		r.re = re
		kept = append(kept, r)
	}
	c.Rules = kept

	// A NEVER rule past the budget is discarded at render time and never
	// reaches the model. Report it on every load rather than let the card
	// file keep claiming a constraint that does not ship.
	for _, d := range c.NeverOverflow() {
		skipped = append(skipped, fmt.Errorf(
			"NEVER rule %d+ over the %d-rule budget, so it never reaches the prompt — cut it or promote it to CRITICAL: %.60q",
			maxNeverRules+1, maxNeverRules, d))
	}
	return &c, skipped, nil
}

var (
	fencedBlock = regexp.MustCompile("(?s)```.*?```")
	inlineCode  = regexp.MustCompile("`[^`\n]*`")
	blockQuote  = regexp.MustCompile(`(?m)^>.*$`)
	doubleQuote = regexp.MustCompile(`"[^"\n]{0,400}"`)
)

// Redact blanks out spans the gate must not judge: code fences, inline code,
// block quotes, and double-quoted passages. Without this the rules fire on
// their own documentation — the flip rule matched this repo's own description
// of the flip rule three times in the originating session.
//
// Spans are replaced with spaces rather than deleted so that byte offsets, and
// therefore the context windows, stay aligned with the original text.
func Redact(s string) string {
	for _, re := range []*regexp.Regexp{fencedBlock, inlineCode, blockQuote, doubleQuote} {
		s = re.ReplaceAllStringFunc(s, func(m string) string {
			return strings.Repeat(" ", len(m))
		})
	}
	return s
}

func contextAround(orig string, lo, hi int) string {
	start := lo - 70
	if start < 0 {
		start = 0
	}
	end := hi + 70
	if end > len(orig) {
		end = len(orig)
	}
	return strings.ReplaceAll(strings.TrimSpace(orig[start:end]), "\n", " ")
}

// Regex runs the card's POSTPROC rules over text, ignoring redacted spans.
func (c *Card) Regex(text string) []Violation {
	return c.regex(text, Redact(text))
}

// regex takes both forms: matching runs over the redacted copy, while the
// reported match and context are sliced out of the original. Redact replaces
// spans with spaces rather than deleting them precisely so the two stay
// byte-aligned.
func (c *Card) regex(text, redacted string) []Violation {
	var out []Violation
	for _, r := range c.Rules {
		for _, loc := range r.re.FindAllStringIndex(redacted, -1) {
			out = append(out, Violation{
				RuleID:  r.ID,
				Action:  r.Action,
				Why:     r.Why,
				Matched: strings.TrimSpace(text[loc[0]:loc[1]]),
				Context: contextAround(text, loc[0], loc[1]),
			})
		}
	}
	return out
}

// --- shape rules -----------------------------------------------------------

var (
	paraSplit    = regexp.MustCompile(`\n\s*\n`)
	sentenceEnd  = regexp.MustCompile(`[.!?](?:\s|$)`)
	wordRe       = regexp.MustCompile(`[A-Za-z']+`)
	listItemHead = regexp.MustCompile(`^\s*(?:[-*+]|\d+\.)\s+`)
	// ordinalHead and clauseSplit were compiled inside their rules — once per
	// call and once per sentence. Hoisting them took the same backfill from
	// 72.7s to 53.2s on its own.
	ordinalHead = regexp.MustCompile(`(?i)^(first|second|third|fourth|fifth|next|finally|lastly)[,.:]`)
	clauseSplit = regexp.MustCompile(`[;,]\s+`)
	linkish     = regexp.MustCompile(`https?://|\]\(`)
)

func words(s string) []string { return wordRe.FindAllString(s, -1) }

// paragraphs returns prose paragraphs, skipping fenced code and list blocks —
// a bulleted list is legitimately uniform and would poison the uniformity test.
//
// It takes text that has already been through Redact. Every rule needs the same
// redacted string, and re-deriving it per rule ran Redact seven times per turn.
func paragraphs(redacted string) []string {
	var out []string
	for _, p := range paraSplit.Split(redacted, -1) {
		p = strings.TrimSpace(p)
		if p == "" || listItemHead.MatchString(p) || strings.HasPrefix(p, "#") {
			continue
		}
		if len(words(p)) < 12 {
			continue
		}
		out = append(out, p)
	}
	return out
}

// LabelledOpening flags a paragraph whose first sentence is a short fragment
// that the rest of the paragraph then unpacks. An ordinal counts as a label,
// which is the form that survived the card's original prohibition.
//
// The bold form is caught by a POSTPROC regex. This catches the unbolded one.
func LabelledOpening(text string) []Violation {
	return labelledOpening(paragraphs(Redact(text)))
}

func labelledOpening(ps []string) []Violation {
	var out []Violation
	for _, p := range ps {
		loc := sentenceEnd.FindStringIndex(p)
		if loc == nil {
			continue
		}
		head := strings.TrimSpace(p[:loc[0]])
		hw := words(head)
		rest := len(words(p)) - len(hw)
		isOrdinal := ordinalHead.MatchString(head)
		// The distinguishing feature is that the opener is a FRAGMENT — a noun
		// phrase standing where a sentence belongs. Length alone is the wrong
		// test: "Basanite measures words" is short and is a real sentence.
		// A head carrying a finite verb is prose, however terse.
		if hasFiniteVerb(hw) && !isOrdinal {
			continue
		}
		short := len(hw) <= 6 || (isOrdinal && len(hw) <= 12)
		if short && rest >= 3*len(hw) {
			out = append(out, Violation{
				RuleID:  "labelled_opening",
				Action:  "warn",
				Why:     "paragraph opens with a label the rest of it then unpacks",
				Matched: head,
				Context: contextAround(p, 0, loc[0]),
			})
		}
	}
	return out
}

// ClauseSymmetry flags a sentence whose two halves are balanced against each
// other: comma- or semicolon-joined clauses of near-equal length that repeat a
// content word across the joint.
//
// Judges under two different models independently flagged this shape as
// generated prose, quoting "a costume that fits well ... a costume that fits
// badly" and "nothing to rank / nothing to match". Naming both terms of a
// comparison is what earns clarity; balancing them is what gives it away.
func ClauseSymmetry(text string) []Violation {
	return clauseSymmetry(paragraphs(Redact(text)))
}

func clauseSymmetry(ps []string) []Violation {
	var out []Violation
	for _, p := range ps {
		for _, sent := range splitSentences(p) {
			// Markdown link lists are mechanically balanced and are not prose.
			if linkish.MatchString(sent) {
				continue
			}
			parts := clauseSplit.Split(sent, -1)
			for i := 0; i+1 < len(parts); i++ {
				a, b := words(parts[i]), words(parts[i+1])
				if len(a) < 4 || len(b) < 4 {
					continue
				}
				lo, hi := float64(len(a)), float64(len(b))
				if lo > hi {
					lo, hi = hi, lo
				}
				if lo/hi < 0.80 {
					continue
				}
				if !sharesContentWord(a, b) {
					continue
				}
				out = append(out, Violation{
					RuleID:  "clause_symmetry",
					Action:  "warn",
					Why:     "balanced two-beat: near-equal clauses repeating a word across the joint",
					Matched: strings.TrimSpace(parts[i] + ", " + parts[i+1]),
					Context: strings.TrimSpace(sent),
				})
			}
		}
	}
	return out
}

// finiteVerb is a deliberately small, high-precision set: auxiliaries and
// copulas, plus the modal set. It is not a POS tagger and does not try to be.
// A head containing any of these is treated as a sentence rather than a label,
// which is enough to separate "Basanite measures words" from "The audit half".
// Third-person -s and past -ed forms are caught by suffix instead of by list.
var finiteVerb = map[string]bool{
	"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
	"has": true, "have": true, "had": true, "does": true, "do": true, "did": true,
	"can": true, "could": true, "will": true, "would": true, "shall": true,
	"should": true, "may": true, "might": true, "must": true, "gets": true,
	"got": true, "goes": true, "went": true, "comes": true, "came": true,
	"means": true, "makes": true, "made": true, "needs": true, "wants": true,
	"gives": true, "takes": true, "sits": true, "reads": true, "writes": true,
}

// verbSuffix catches regular third-person and past forms that the small list
// misses ("measures", "carried"). Nouns ending in -s are the false-positive
// risk here; that costs recall on the label rule, which is the safe direction
// for a warn-only gate.
var verbSuffix = regexp.MustCompile(`(?i)^[a-z]{4,}(es|ed)$`)

func hasFiniteVerb(ws []string) bool {
	for _, w := range ws {
		lw := strings.ToLower(w)
		if finiteVerb[lw] || verbSuffix.MatchString(lw) {
			return true
		}
	}
	return false
}

func sharesContentWord(a, b []string) bool {
	seen := map[string]bool{}
	for _, w := range a {
		if len(w) >= 5 {
			seen[strings.ToLower(w)] = true
		}
	}
	for _, w := range b {
		if len(w) >= 5 && seen[strings.ToLower(w)] {
			return true
		}
	}
	return false
}

func splitSentences(p string) []string {
	locs := sentenceEnd.FindAllStringIndex(p, -1)
	var out []string
	prev := 0
	for _, l := range locs {
		out = append(out, strings.TrimSpace(p[prev:l[1]]))
		prev = l[1]
	}
	if prev < len(p) {
		out = append(out, strings.TrimSpace(p[prev:]))
	}
	return out
}

// ParagraphUniformity flags a response whose paragraphs are all about the same
// length. Coefficient of variation below the threshold across four or more
// prose paragraphs means the rhythm is machine-even.
func ParagraphUniformity(text string, minCV float64) []Violation {
	return paragraphUniformity(paragraphs(Redact(text)), minCV)
}

func paragraphUniformity(ps []string, minCV float64) []Violation {
	if len(ps) < 4 {
		return nil
	}
	var lens []float64
	for _, p := range ps {
		lens = append(lens, float64(len(words(p))))
	}
	var sum float64
	for _, l := range lens {
		sum += l
	}
	mean := sum / float64(len(lens))
	if mean == 0 {
		return nil
	}
	var ss float64
	for _, l := range lens {
		ss += (l - mean) * (l - mean)
	}
	cv := math.Sqrt(ss/float64(len(lens))) / mean
	if cv >= minCV {
		return nil
	}
	counts := make([]string, len(lens))
	for i, l := range lens {
		counts[i] = fmt.Sprintf("%d", int(l))
	}
	return []Violation{{
		RuleID:  "paragraph_uniformity",
		Action:  "warn",
		Why:     fmt.Sprintf("paragraph lengths too even (CV %.2f < %.2f)", cv, minCV),
		Matched: strings.Join(counts, "/") + " words",
		Context: fmt.Sprintf("%d prose paragraphs, mean %.0f words", len(lens), mean),
	}}
}

var (
	askPunct  = regexp.MustCompile(`\?\s*$`)
	askPhrase = regexp.MustCompile(`(?i)\b(let me know|should i|do you want|want me to|which (?:one|option|approach) do you|go ahead and|confirm(?:ing)?\b.{0,20}\byou|your call|need (?:from|you) to|waiting on you)\b`)
)

// AskNotLast flags a question or request for the user that sits in an earlier
// paragraph while the reply keeps going after it. The complaint it targets is
// concrete: a reader who opens a long reply from the bottom, because that is
// where the terminal leaves the cursor, and misses an ask that landed 2/3 of
// the way up.
//
// Order among multiple asks is a real form of the same complaint but is not
// checked here — telling "printed in the wrong order" apart from "printed in
// the right order" needs to know what the asks mean, not just where they sit.
//
// It splits on rawParagraphs for the same reason DanglingEnd does. Under the
// 12-word prose floor a short ask was invisible: "…\n\nBuild both?\n\n…" gave
// zero violations while the identical ask padded to fourteen words gave one.
// Asks are short by nature, so the floor was deleting the evidence.
func AskNotLast(text string) []Violation {
	return askNotLast(rawParagraphs(Redact(text)))
}

func askNotLast(ps []string) []Violation {
	if len(ps) < 2 {
		return nil
	}
	last := len(ps) - 1
	var out []Violation
	for i, p := range ps {
		if i == last {
			continue
		}
		for _, sent := range splitSentences(p) {
			if askPunct.MatchString(sent) || askPhrase.MatchString(sent) {
				out = append(out, Violation{
					RuleID:  "ask_not_last",
					Action:  "warn",
					Why:     "a question or request for the user sits mid-reply instead of at the end",
					Matched: strings.TrimSpace(sent),
					Context: strings.TrimSpace(strings.ReplaceAll(p, "\n", " ")),
				})
				break
			}
		}
	}
	return out
}

// --- end-of-turn decision surface ------------------------------------------

var (
	closurePhrase = regexp.MustCompile(`(?i)\b(nothing (?:for you )?to act on|no action (?:needed|required|from you)|loop closed|nothing (?:left |remains )?for you)\b`)
	offerPhrase   = regexp.MustCompile(`(?i)\b(next,? i'?(?:d|ll) |say the word|say ["']?continue["']?|the next step is|i'?ll (?:start|take|do|fix|run) )`)
	openProblem   = regexp.MustCompile(`(?i)\b(still (?:open|broken|failing|unresolved|unfixed|missing|unverified|sitting|gated|binding|live|needs)|hasn'?t been|haven'?t (?:been|checked|verified|measured)|never got|remains? (?:open|unverified|unfixed|broken|unresolved)|unresolved|unverified|unaudited|unfixed|needs? (?:fixing|a fix|someone|removing|checking)|would need|worth a look|should agree to|i have a guess)\b`)
)

// rawParagraphs splits like paragraphs but keeps short blocks and lists. The
// decision-surface rules need the real final block: "Build both?" is two words
// and is the whole point of the reply, so a word floor would delete exactly the
// evidence that clears the check.
func rawParagraphs(redacted string) []string {
	var out []string
	for _, p := range paraSplit.Split(redacted, -1) {
		p = strings.TrimSpace(p)
		if p == "" || strings.HasPrefix(p, "#") {
			continue
		}
		out = append(out, p)
	}
	return out
}

func isDecision(sent string) bool {
	return askPunct.MatchString(sent) || askPhrase.MatchString(sent) ||
		closurePhrase.MatchString(sent) || offerPhrase.MatchString(sent)
}

// DanglingEnd flags the ending Opus 5 made a habit of: the reply names open
// problems near its end and then stops without a question, an offer, or an
// explicit all-clear. Measured across 473 Opus 5 end-of-turns in July 2026,
// 10.1% ended this way, against 6.0% for Opus 4.x — and the user's follow-up
// "what specifically should we do next" appeared 14 times against zero.
//
// The reader it protects re-enters the terminal cold from another window and
// needs the final block to be answerable: continue, redirect, or decide. Two
// forms fail that. No decision point anywhere while the tail names open
// problems (dangling_end), and a decision point that observations then bury
// (buried_decision).
func DanglingEnd(text string) []Violation {
	return danglingEnd(rawParagraphs(Redact(text)))
}

func danglingEnd(ps []string) []Violation {
	if len(ps) == 0 {
		return nil
	}
	type located struct {
		text string
		para int
	}
	var sents []located
	for i, p := range ps {
		for _, s := range splitSentences(p) {
			sents = append(sents, located{s, i})
		}
	}
	lastDecision := -1
	for i, s := range sents {
		if isDecision(s.text) {
			lastDecision = i
		}
	}
	tailPara := len(ps) - 2
	if tailPara < 0 {
		tailPara = 0
	}
	var out []Violation
	if lastDecision == -1 {
		for _, s := range sents {
			if s.para >= tailPara && openProblem.MatchString(s.text) {
				out = append(out, Violation{
					RuleID:  "dangling_end",
					Action:  "warn",
					Why:     "names an open problem and ends with no question, offer, or all-clear — the reader cannot answer 'continue'",
					Matched: strings.TrimSpace(s.text),
					Context: strings.TrimSpace(strings.ReplaceAll(ps[len(ps)-1], "\n", " ")),
				})
				break
			}
		}
		return out
	}
	for i := lastDecision + 1; i < len(sents); i++ {
		if openProblem.MatchString(sents[i].text) {
			out = append(out, Violation{
				RuleID:  "buried_decision",
				Action:  "warn",
				Why:     "an open problem lands after the last question or offer, burying the decision point above it",
				Matched: strings.TrimSpace(sents[i].text),
				Context: strings.TrimSpace(sents[lastDecision].text),
			})
			break
		}
	}
	return out
}

// Shapes runs every built-in structural rule. The individual rules are
// exported for tests and for callers wanting one of them; Card.Check is the
// entry point for running everything over a real turn.
func Shapes(text string, minCV float64) []Violation {
	return shapes(Redact(text), minCV)
}

func shapes(redacted string, minCV float64) []Violation {
	prose := paragraphs(redacted)
	raw := rawParagraphs(redacted)
	var out []Violation
	out = append(out, labelledOpening(prose)...)
	out = append(out, clauseSymmetry(prose)...)
	out = append(out, paragraphUniformity(prose, minCV)...)
	out = append(out, askNotLast(raw)...)
	out = append(out, danglingEnd(raw)...)
	return out
}

// Check runs the card's regex rules and every shape rule over one turn.
//
// Redaction happens once here and the paragraph split twice, against seven
// redactions and five splits back when every rule derived its own. With the
// two loop-compiled patterns hoisted to package level, a 20.1M-character
// backfill went from 72.7s to 43.4s reporting the same violations.
func (c *Card) Check(text string, minCV float64) []Violation {
	redacted := Redact(text)
	return append(c.regex(text, redacted), shapes(redacted, minCV)...)
}
