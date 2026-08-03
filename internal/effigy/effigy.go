// Package effigy parses the subset of effigy notation that cope reads and
// emits the card JSON the gate already consumes.
//
// tools/card2json.py remains the build-time path for the embedded card, and
// CI still regenerates card/rules.json through it, so effigy stays the single
// authority on the notation. This package exists for the other direction: at
// runtime a reader should be able to point cope at a .effigy file and have it
// work, without a Python interpreter and a checkout of the effigy repo sitting
// behind every card they write. A voicing framework whose second voice needs a
// build step has one voice.
//
// It is deliberately a translator rather than a second reader of the format.
// Output goes through scan.ParseCard exactly as the generated JSON does, so
// pattern compilation, the strip-is-unimplementable warning, and every other
// check happen in one place. Where the two paths could disagree, a test reads
// every card in the repo through both and compares.
//
// Blocks effigy defines and cope has no use for (ARC, GOALS, RELS, SECRETS and
// the rest) are skipped the way the reference parser skips an unknown keyword,
// so a card carrying them still loads.
package effigy

import (
	"encoding/json"
	"fmt"
	"strings"
)

type card struct {
	CardID string   `json:"card_id"`
	Source string   `json:"source"`
	Rules  []rule   `json:"rules"`
	Voice  voice    `json:"voice"`
	Theme  string   `json:"theme"`
	Traits []string `json:"traits"`
	Never  []gated  `json:"never"`
	Quirks []string `json:"quirks"`
	Tests  []test   `json:"tests"`
	Wrong  []pair   `json:"wrong"`
	Mes    []string `json:"mes"`

	// Gate carries @gate lines through verbatim. effigy's own parser drops
	// unrecognised @keys, so this is the one place cope reads a header the
	// notation does not define; the meaning of the string is scan's business,
	// not this package's.
	Gate []string `json:"gate,omitempty"`

	// Shape carries @shape lines the same way, for the card's own structure
	// rules. Same reasoning: the string's meaning is scan's business.
	Shape []string `json:"shape,omitempty"`
}

type rule struct {
	ID      string `json:"id"`
	Action  string `json:"action"`
	Pattern string `json:"pattern"`
	Why     string `json:"why"`
}

type voice struct {
	Kernel string `json:"kernel"`
	Peak   string `json:"peak"`
}

type gated struct {
	Text string `json:"text"`
	When string `json:"when"`
}

type test struct {
	Name     string   `json:"name"`
	Question string   `json:"question"`
	Fail     []string `json:"fail"`
	Pass     []string `json:"pass"`
	Why      string   `json:"why"`
	When     string   `json:"when"`
}

type pair struct {
	Wrong string `json:"wrong"`
	Right string `json:"right"`
	Why   string `json:"why"`
	When  string `json:"when"`
}

// ToJSON parses effigy notation and returns the card JSON scan.ParseCard
// reads. source is recorded as the card's origin and is not opened.
func ToJSON(src []byte, source string) ([]byte, error) {
	c, err := parse(src, source)
	if err != nil {
		return nil, err
	}
	return json.Marshal(c)
}

// braced lists the blocks effigy delimits with {} rather than []. Everything
// else it knows is bracketed, and anything in neither list is unknown.
var braced = map[string]bool{
	"VOICE": true, "ARC": true, "GOALS": true, "RELS": true,
	"SCHED": true, "DM": true, "BEHAVIORS": true,
}

func parse(src []byte, source string) (*card, error) {
	c := &card{
		Source: source,
		Rules:  []rule{}, Traits: []string{}, Never: []gated{},
		Quirks: []string{}, Tests: []test{}, Wrong: []pair{}, Mes: []string{},
	}
	p := &parser{src: src}
	for {
		p.skipSpaceAndNewlines()
		if p.done() {
			return c, nil
		}
		switch {
		case p.src[p.pos] == '#':
			p.skipLine()
		case p.src[p.pos] == '@':
			p.header(c)
		default:
			word := p.peekWord()
			if word == "" {
				p.skipLine()
				continue
			}
			if !braced[word] && !bracketed[word] {
				// Unknown keyword. The reference parser skips the line and
				// carries on rather than failing, which leaves the block's
				// body to be walked as loose lines and discarded the same way.
				p.skipLine()
				continue
			}
			p.pos += len(word)
			open, close := byte('['), byte(']')
			if braced[word] {
				open, close = '{', '}'
			}
			body, err := p.readBlock(open, close)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", word, err)
			}
			switch word {
			case "VOICE":
				kv := kvBlock(body)
				c.Voice = voice{Kernel: kv["kernel"], Peak: kv["peak"]}
			case "TRAITS":
				c.Traits = parseTraits(body)
			case "NEVER":
				c.Never = parseNever(body)
			case "QUIRKS":
				c.Quirks = parseLines(body)
			case "TEST":
				c.Tests = parseTests(body)
			case "WRONG":
				c.Wrong = parseWrong(body)
			case "MES":
				c.Mes = parseMes(body)
			case "POSTPROC":
				c.Rules = parsePostproc(body)
			}
		}
	}
}

// bracketed lists every []-delimited block effigy knows. Blocks cope ignores
// are here so their bodies are consumed as a unit; leaving them out would send
// their contents through the top-level loop a line at a time.
var bracketed = map[string]bool{
	"TRAITS": true, "NEVER": true, "QUIRKS": true, "MES": true, "UNC": true,
	"SECRETS": true, "ERA": true, "ARRIVE": true, "DEPART": true,
	"WRONG": true, "TEST": true, "PROPS": true, "POSTPROC": true,
}

func (p *parser) header(c *card) {
	p.pos++ // @
	start := p.pos
	for !p.done() && p.src[p.pos] != ' ' && p.src[p.pos] != '\t' && p.src[p.pos] != '\n' {
		p.pos++
	}
	key := string(p.src[start:p.pos])
	p.skipSpace()
	value := strings.TrimSpace(p.readLine())
	switch key {
	case "id":
		c.CardID = value
	case "theme":
		c.Theme = value
	case "shape":
		c.Shape = append(c.Shape, value)
	case "gate":
		// Repeated, one rule per line, so it accumulates rather than
		// overwriting. Validation happens in scan, which knows the rule ids.
		c.Gate = append(c.Gate, value)
	}
	// @name and @role are read by effigy and carried nowhere: the gate renders
	// a register, not a character sheet. Other keys belong to blocks cope does
	// not use.
}

// --- block bodies -----------------------------------------------------------

// splitItems cuts a block body on --- lines, dropping blank and commented
// lines. Interior indentation survives; the per-block parsers trim.
func splitItems(body string) []string {
	var items []string
	var cur []string
	for _, line := range strings.Split(body, "\n") {
		s := strings.TrimSpace(line)
		switch {
		case s == "---":
			if len(cur) > 0 {
				items = append(items, strings.Join(cur, "\n"))
				cur = nil
			}
		case s != "" && !strings.HasPrefix(s, "#"):
			cur = append(cur, line)
		}
	}
	if len(cur) > 0 {
		items = append(items, strings.Join(cur, "\n"))
	}
	return items
}

// kvBlock reads `key: value` lines. A line with no colon continues the value
// above it, which is how a kernel can be wrapped across lines in a card.
func kvBlock(body string) map[string]string {
	out := map[string]string{}
	last := ""
	for _, line := range strings.Split(body, "\n") {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		if i := strings.Index(s, ":"); i >= 0 {
			k := strings.TrimSpace(s[:i])
			out[k] = strings.TrimSpace(s[i+1:])
			last = k
			continue
		}
		if last != "" {
			out[last] += " " + s
		}
	}
	return out
}

func parseTraits(body string) []string {
	var kept []string
	for _, line := range strings.Split(body, "\n") {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		kept = append(kept, s)
	}
	out := []string{}
	for _, t := range strings.Split(strings.Join(kept, " "), ",") {
		if t = strings.TrimSpace(t); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func parseNever(body string) []gated {
	out := []gated{}
	for _, item := range splitItems(body) {
		when := ""
		var text []string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			line := strings.TrimSpace(raw)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if rest, ok := directive(line, "@when"); ok {
				when = rest
				continue
			}
			text = append(text, line)
		}
		if t := strings.Join(text, " "); t != "" {
			out = append(out, gated{Text: t, When: normalizeWhen(when)})
		}
	}
	return out
}

func parseLines(body string) []string {
	out := []string{}
	for _, item := range splitItems(body) {
		var lines []string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			if s := strings.TrimSpace(raw); s != "" {
				lines = append(lines, s)
			}
		}
		if t := strings.Join(lines, " "); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func parseWrong(body string) []pair {
	out := []pair{}
	for _, item := range splitItems(body) {
		var p pair
		var when string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			line := strings.TrimSpace(raw)
			if line == "" {
				continue
			}
			if rest, ok := directive(line, "@when"); ok {
				when = rest
				continue
			}
			switch {
			case strings.HasPrefix(line, "WRONG:"):
				// effigy strips quotes from the two examples and not from the
				// why, so a card can quote inside its reasoning.
				p.Wrong = strings.Trim(strings.TrimSpace(line[6:]), `"`)
			case strings.HasPrefix(line, "RIGHT:"):
				p.Right = strings.Trim(strings.TrimSpace(line[6:]), `"`)
			case strings.HasPrefix(line, "WHY:"):
				p.Why = strings.TrimSpace(line[4:])
			}
		}
		if p.Wrong != "" {
			p.When = normalizeWhen(when)
			out = append(out, p)
		}
	}
	return out
}

func parseTests(body string) []test {
	out := []test{}
	for _, item := range splitItems(body) {
		t := test{Fail: []string{}, Pass: []string{}}
		var when string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			line := strings.TrimSpace(raw)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if rest, ok := directive(line, "@when"); ok {
				when = rest
				continue
			}
			if _, ok := directive(line, "@beat"); ok {
				continue
			}
			low := strings.ToLower(line)
			switch {
			case strings.HasPrefix(low, "name:"):
				t.Name = strings.TrimSpace(line[5:])
			case strings.HasPrefix(low, "question:"):
				t.Question = strings.TrimSpace(line[9:])
			case strings.HasPrefix(low, "fail:"):
				t.Fail = append(t.Fail, strings.TrimSpace(line[5:]))
			case strings.HasPrefix(low, "pass:"):
				t.Pass = append(t.Pass, strings.TrimSpace(line[5:]))
			case strings.HasPrefix(low, "why:"):
				t.Why = strings.TrimSpace(line[4:])
			}
			// dimension: is effigy's own axis label and reaches no prompt.
		}
		// A test with no question asks nothing, so effigy drops it rather than
		// rendering a heading with nothing under it.
		if t.Name != "" && t.Question != "" {
			t.When = normalizeWhen(when)
			out = append(out, t)
		}
	}
	return out
}

// parseMes cannot use splitItems: a MES entry is a whole exchange and its line
// breaks are the thing that makes it one, so the lines are kept apart here and
// joined only at the end.
func parseMes(body string) []string {
	var items [][]string
	var cur []string
	flush := func() {
		if len(cur) > 0 {
			items = append(items, cur)
			cur = nil
		}
	}
	for _, line := range strings.Split(body, "\n") {
		s := strings.TrimSpace(line)
		switch {
		case s == "---":
			flush()
		case s == "" || strings.HasPrefix(s, "#"):
		case strings.HasPrefix(s, "@"):
			// @tier, @when and @beat annotate the example; none reach the card.
		default:
			cur = append(cur, s)
		}
	}
	flush()

	out := []string{}
	for _, lines := range items {
		hasUser := false
		for _, l := range lines {
			if strings.HasPrefix(l, "{{user}}:") {
				hasUser = true
				break
			}
		}
		if hasUser {
			out = append(out, strings.Join(lines, "\n"))
			continue
		}
		// A bare example is the character speaking unprompted.
		text := strings.Join(lines, " ")
		if !strings.HasPrefix(text, "{{char}}:") {
			text = "{{char}}: " + text
		}
		out = append(out, text)
	}
	return out
}

func parsePostproc(body string) []rule {
	out := []rule{}
	for i, item := range splitItems(body) {
		kv := kvBlock(item)
		action := strings.ToLower(strings.TrimSpace(kv["action"]))
		pattern := strings.TrimSpace(kv["pattern"])
		if pattern == "" || (action != "reject" && action != "strip" && action != "warn") {
			// An authoring error in one rule should not cost the card.
			continue
		}
		id := strings.TrimSpace(kv["id"])
		if id == "" {
			id = fmt.Sprintf("postproc_%d", i)
		}
		out = append(out, rule{ID: id, Action: action, Pattern: pattern, Why: strings.TrimSpace(kv["why"])})
	}
	return out
}

// directive matches an @-prefixed annotation and returns what follows it.
func directive(line, name string) (string, bool) {
	if !strings.HasPrefix(line, name) {
		return "", false
	}
	return strings.TrimSpace(line[len(name):]), true
}

// normalizeWhen folds effigy's always-on marker to the empty condition, which
// is what the gate's own gating reads as ungated.
func normalizeWhen(w string) string {
	if w = strings.TrimSpace(w); w == "*" {
		return ""
	}
	return w
}

// --- scanning ---------------------------------------------------------------

type parser struct {
	src []byte
	pos int
}

func (p *parser) done() bool { return p.pos >= len(p.src) }

func (p *parser) skipSpace() {
	for !p.done() && (p.src[p.pos] == ' ' || p.src[p.pos] == '\t') {
		p.pos++
	}
}

func (p *parser) skipSpaceAndNewlines() {
	for !p.done() {
		switch p.src[p.pos] {
		case ' ', '\t', '\n', '\r':
			p.pos++
		default:
			return
		}
	}
}

func (p *parser) skipLine() {
	for !p.done() && p.src[p.pos] != '\n' {
		p.pos++
	}
	if !p.done() {
		p.pos++
	}
}

func (p *parser) readLine() string {
	start := p.pos
	for !p.done() && p.src[p.pos] != '\n' {
		p.pos++
	}
	line := string(p.src[start:p.pos])
	if !p.done() {
		p.pos++
	}
	return line
}

func (p *parser) peekWord() string {
	i := p.pos
	for i < len(p.src) && (p.src[i] == ' ' || p.src[i] == '\t') {
		i++
	}
	j := i
	for j < len(p.src) && isAlpha(p.src[j]) {
		j++
	}
	return string(p.src[i:j])
}

// readBlock consumes a delimited body, counting nested delimiters so a regex
// carrying a character class survives inside a bracketed block.
func (p *parser) readBlock(open, close byte) (string, error) {
	p.skipSpaceAndNewlines()
	if p.done() || p.src[p.pos] != open {
		return "", fmt.Errorf("expected %q", string(open))
	}
	p.pos++
	depth, start := 1, p.pos
	for !p.done() {
		switch p.src[p.pos] {
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				body := string(p.src[start:p.pos])
				p.pos++
				return body, nil
			}
		}
		p.pos++
	}
	return "", fmt.Errorf("unterminated block, no closing %q", string(close))
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
