// Package effigy turns an effigy card into the card JSON the gate already
// consumes.
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
// Reading the notation — which keywords open a block, which delimiters close
// one, how a body is cut into items — belongs to effigy and comes from
// effigy/go/notation. What is left here is cope's half: which blocks it wants
// and what shape they take in the gate's JSON. Output then goes through
// scan.ParseCard exactly as the generated JSON does, so pattern compilation,
// the strip-is-unimplementable warning, and every other check happen in one
// place. Where the two paths could disagree, a test reads every card in the
// repo through both and compares.
//
// Blocks effigy defines and cope has no use for (ARC, GOALS, RELS, SECRETS and
// the rest) fall through the switch below, as does any block effigy does not
// define at all.
package effigy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/justinstimatze/effigy/go/notation"
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

	// Gate carries @gate lines through verbatim. The notation does not define
	// the key, and effigy keeps unrecognised @keys rather than dropping them,
	// so it arrives here like any other header; the meaning of the string is
	// scan's business, not this package's.
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

func parse(src []byte, source string) (*card, error) {
	c := &card{
		Source: source,
		Rules:  []rule{}, Traits: []string{}, Never: []gated{},
		Quirks: []string{}, Tests: []test{}, Wrong: []pair{}, Mes: []string{},
	}

	doc, err := notation.Scan(src)
	if err != nil {
		return nil, err
	}

	for _, h := range doc.Headers {
		switch h.Key {
		case "id":
			c.CardID = h.Value
		case "theme":
			c.Theme = h.Value
		case "shape":
			c.Shape = append(c.Shape, h.Value)
		case "gate":
			// Repeated, one rule per line, so it accumulates rather than
			// overwriting. Validation happens in scan, which knows the rule ids.
			c.Gate = append(c.Gate, h.Value)
		}
		// @name and @role are read by effigy and carried nowhere: the gate
		// renders a register, not a character sheet. Other keys belong to
		// blocks cope does not use.
	}

	for _, b := range doc.Blocks {
		switch b.Name {
		case "VOICE":
			kv := notation.KV(b.Body)
			c.Voice = voice{Kernel: kv["kernel"], Peak: kv["peak"]}
		case "TRAITS":
			c.Traits = parseTraits(b.Body)
		case "NEVER":
			c.Never = parseNever(b.Body)
		case "QUIRKS":
			c.Quirks = parseLines(b.Body)
		case "TEST":
			c.Tests = parseTests(b.Body)
		case "WRONG":
			c.Wrong = parseWrong(b.Body)
		case "MES":
			c.Mes = parseMes(b.Body)
		case "POSTPROC":
			c.Rules = parsePostproc(b.Body)
		}
	}
	return c, nil
}

// --- block bodies -----------------------------------------------------------

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
	for _, item := range notation.SplitItems(body) {
		when := ""
		var text []string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			line := strings.TrimSpace(raw)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if rest, ok := notation.Directive(line, "@when"); ok {
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
	for _, item := range notation.SplitItems(body) {
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
	for _, item := range notation.SplitItems(body) {
		var p pair
		var when string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			line := strings.TrimSpace(raw)
			if line == "" {
				continue
			}
			if rest, ok := notation.Directive(line, "@when"); ok {
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
	for _, item := range notation.SplitItems(body) {
		t := test{Fail: []string{}, Pass: []string{}}
		var when string
		for _, raw := range strings.Split(strings.TrimSpace(item), "\n") {
			line := strings.TrimSpace(raw)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if rest, ok := notation.Directive(line, "@when"); ok {
				when = rest
				continue
			}
			if _, ok := notation.Directive(line, "@beat"); ok {
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

// parseMes keeps only the text. @tier, @when and @beat gate an example for
// readers that render a character; the gate shows every example it is given, so
// notation.MES reads the annotations and this drops them.
//
// The joining rule this used to carry itself is now upstream, which also fixes
// a divergence: an @key the notation does not define is content, and the copy
// here silently swallowed it.
func parseMes(body string) []string {
	out := []string{}
	for _, ex := range notation.MES(body) {
		out = append(out, ex.Text)
	}
	return out
}

func parsePostproc(body string) []rule {
	out := []rule{}
	for i, item := range notation.SplitItems(body) {
		kv := notation.KV(item)
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

// normalizeWhen folds effigy's always-on marker to the empty condition, which
// is what the gate's own gating reads as ungated.
func normalizeWhen(w string) string {
	if w = strings.TrimSpace(w); w == "*" {
		return ""
	}
	return w
}
