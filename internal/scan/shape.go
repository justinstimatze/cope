package scan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Card-authored structure rules.
//
// @gate let a card decline a built-in rule. This is the other half: a card
// states a structural commitment and the gate checks it, so the set of things
// cope can check stops being fixed at whatever the binary happens to implement.
//
// The vocabulary is deliberately small — a selector naming which blocks the
// claim is about, and a predicate that is true or false of one. It is not an
// expression language and it is not going to become one. A card wanting
// something outside it should say so in POSTPROC, where a regex reaches
// anything a regex can reach; what belongs here is the shape claims a regex
// cannot make, which is why the predicates count things rather than match them.
//
// The rules it can express are weaker than the compiled ones. clause_symmetry
// is not writable in this grammar and should not be. What is writable is the
// commitment a card's own VOICE block already states in prose — hold the
// conclusion for a short sentence, never run a paragraph past this length, end
// on something the reader can answer — which is exactly the class that used to
// have nowhere to go.

// ShapeRule is one @shape line: a card's own structural check.
type ShapeRule struct {
	ID       string `json:"id"`
	Selector string `json:"selector"`
	Predicate

	Why string `json:"why"`
}

// Predicate is the test applied to whatever the selector picked.
type Predicate struct {
	Kind string `json:"kind"` // words | sentences | asks
	Op   string `json:"op"`   // <= | >= | is | is not
	N    int    `json:"n,omitempty"`
}

// selectors map the card's words to how the blocks are picked. Every one of
// them reads rawParagraphs rather than paragraphs: a card asserting that its
// closing line is short is asserting it about a block the prose rules skip for
// being short, and running this against the filtered set would make the most
// useful claim in the grammar unexpressible.
var selectors = []string{
	"first paragraph",
	"last paragraph",
	"every paragraph",
	"some paragraph",
	"reply",
}

var shapeIDRe = regexp.MustCompile(`^[a-z][a-z0-9_]{1,39}$`)

const shapeSyntax = `@shape <id>: <selector> <predicate> — <why>` +
	"\n  selectors:  first paragraph | last paragraph | every paragraph | some paragraph | reply" +
	"\n  predicates: words <= N | words >= N | sentences <= N | sentences >= N | asks | does not ask"

// Shapes parses the card's @shape lines, reporting every problem rather than
// the first. A rule that does not parse is not applied, so a silent failure
// here would be a card believing it had a check it does not have — the same
// failure @gate guards against, in the opposite direction.
func (c *Card) Shapes() ([]ShapeRule, error) {
	var out []ShapeRule
	var bad []string
	seen := map[string]bool{}
	for _, line := range c.Shape {
		r, err := parseShape(line)
		if err != nil {
			bad = append(bad, err.Error())
			continue
		}
		switch {
		case RuleIDs[r.ID]:
			bad = append(bad, fmt.Sprintf("%q is already a built-in rule; a card rule needs its own name", r.ID))
			continue
		case seen[r.ID]:
			bad = append(bad, fmt.Sprintf("%q defined twice", r.ID))
			continue
		}
		for _, existing := range c.Rules {
			if existing.ID == r.ID {
				bad = append(bad, fmt.Sprintf("%q is already a POSTPROC rule id", r.ID))
			}
		}
		seen[r.ID] = true
		out = append(out, r)
	}
	if len(bad) > 0 {
		return out, fmt.Errorf("@shape: %s\n  syntax: %s", strings.Join(bad, "; "), shapeSyntax)
	}
	return out, nil
}

func parseShape(line string) (ShapeRule, error) {
	head, why, ok := cutReason(line)
	if !ok || why == "" {
		return ShapeRule{}, fmt.Errorf("%.40q has no reason after the dash, and a check nobody can review is worse than no check", line)
	}
	id, rest, ok := strings.Cut(head, ":")
	if !ok {
		return ShapeRule{}, fmt.Errorf("%.40q has no id, and the id is what a violation is reported under", head)
	}
	id = strings.TrimSpace(id)
	if !shapeIDRe.MatchString(id) {
		return ShapeRule{}, fmt.Errorf("%q is not a usable id: lowercase, digits and underscores, 2 to 40 characters", id)
	}

	rest = strings.TrimSpace(rest)
	sel := ""
	for _, s := range selectors {
		if rest == s || strings.HasPrefix(rest, s+" ") {
			sel, rest = s, strings.TrimSpace(strings.TrimPrefix(rest, s))
			break
		}
	}
	if sel == "" {
		return ShapeRule{}, fmt.Errorf("%q does not start with a selector", rest)
	}
	p, err := parsePredicate(rest)
	if err != nil {
		return ShapeRule{}, fmt.Errorf("%s: %w", id, err)
	}
	return ShapeRule{ID: id, Selector: sel, Predicate: p, Why: why}, nil
}

// cutReason splits a card line on the first dash form it finds. Shared with
// @gate's parser in spirit and kept separate in fact, because the two grammars
// are free to diverge and a shared helper would quietly couple them.
func cutReason(line string) (head, why string, ok bool) {
	for _, sep := range []string{" — ", " -- ", " - "} {
		if i := strings.Index(line, sep); i >= 0 {
			return line[:i], strings.TrimSpace(line[i+len(sep):]), true
		}
	}
	return line, "", false
}

func parsePredicate(s string) (Predicate, error) {
	switch s {
	case "asks":
		return Predicate{Kind: "asks", Op: "is"}, nil
	case "does not ask":
		return Predicate{Kind: "asks", Op: "is not"}, nil
	}
	f := strings.Fields(s)
	if len(f) != 3 {
		return Predicate{}, fmt.Errorf("%q is not a predicate", s)
	}
	if f[0] != "words" && f[0] != "sentences" {
		return Predicate{}, fmt.Errorf("%q counts nothing this knows; use words or sentences", f[0])
	}
	if f[1] != "<=" && f[1] != ">=" {
		return Predicate{}, fmt.Errorf("%q is not a comparison; use <= or >=", f[1])
	}
	n, err := strconv.Atoi(f[2])
	if err != nil || n < 1 {
		return Predicate{}, fmt.Errorf("%q is not a positive count", f[2])
	}
	return Predicate{Kind: f[0], Op: f[1], N: n}, nil
}

// shapeRules runs the card's own structure rules over the reply.
func (c *Card) shapeRules(redacted string) []Violation {
	if len(c.Shape) == 0 {
		return nil
	}
	ps := rawParagraphs(redacted)
	if len(ps) == 0 {
		return nil
	}
	var out []Violation
	for _, line := range c.Shape {
		r, err := parseShape(line)
		if err != nil || RuleIDs[r.ID] {
			continue // reported by Shapes at load
		}
		out = append(out, r.apply(ps, redacted)...)
	}
	return out
}

func (r ShapeRule) apply(ps []string, whole string) []Violation {
	switch r.Selector {
	case "reply":
		if r.holds(whole) {
			return nil
		}
		return []Violation{r.violation(firstLine(whole))}
	case "first paragraph":
		if r.holds(ps[0]) {
			return nil
		}
		return []Violation{r.violation(ps[0])}
	case "last paragraph":
		last := ps[len(ps)-1]
		if r.holds(last) {
			return nil
		}
		return []Violation{r.violation(last)}
	case "some paragraph":
		for _, p := range ps {
			if r.holds(p) {
				return nil
			}
		}
		// Nothing in the reply satisfied it, so the match is the whole reply
		// rather than any one block. Quote the last, which is where a reader
		// looking for a missing closing move will look first.
		return []Violation{r.violation(ps[len(ps)-1])}
	case "every paragraph":
		var out []Violation
		for _, p := range ps {
			if !r.holds(p) {
				out = append(out, r.violation(p))
			}
		}
		return out
	}
	return nil
}

func (r ShapeRule) holds(s string) bool {
	switch r.Kind {
	case "words":
		return compare(len(words(s)), r.Op, r.N)
	case "sentences":
		return compare(len(nonEmptySentences(s)), r.Op, r.N)
	case "asks":
		asks := false
		for _, sent := range splitSentences(s) {
			if askPunct.MatchString(sent) || askPhrase.MatchString(sent) {
				asks = true
				break
			}
		}
		return asks == (r.Op == "is")
	}
	return true
}

func compare(got int, op string, want int) bool {
	if op == "<=" {
		return got <= want
	}
	return got >= want
}

func nonEmptySentences(p string) []string {
	var out []string
	for _, s := range splitSentences(p) {
		if strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

func (r ShapeRule) violation(matched string) Violation {
	return Violation{
		RuleID: r.ID,
		Action: "warn",
		// The card's own words. A rule the card wrote should explain itself in
		// the card's terms, not in a sentence this file made up about it.
		Why:     r.Why,
		Matched: strings.TrimSpace(firstLine(matched)),
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
