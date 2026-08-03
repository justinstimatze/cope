package scan

import (
	"fmt"
	"strings"
)

// RuleIDs is every rule id the gate can emit, and the authority on what a card
// is allowed to name in an @gate line.
//
// It lives here rather than beside the descriptions in cmd/cope-gate because a
// card is validated at load, long before anything wants to describe a rule.
// TestRuleIDsMatchTheDocumentedSet keeps the two from drifting.
var RuleIDs = map[string]bool{
	"labelled_opening":     true,
	"clause_symmetry":      true,
	"paragraph_uniformity": true,
	"ask_not_last":         true,
	"dangling_end":         true,
	"buried_decision":      true,
	"forked_end":           true,
	"apology":              true,
	"self_postmortem":      true,
	"announced_length":     true,
	"cross_turn_repeat":    true,
	"unverified_done":      true,
	"loop_ask":             true,
}

// GateDirective is one @gate line: a built-in rule this card is exempt from,
// and the card's own reason for the exemption.
//
// Exemption rather than deletion is the whole point. The rule still runs, and
// its hits are dropped for this card only, so a card that turns off
// dangling_end does not change what any other card is measured against.
type GateDirective struct {
	Rule string `json:"rule"`
	Why  string `json:"why,omitempty"`
}

// gateSyntax is what an @gate line has to look like:
//
//	@gate dangling_end off — VOICE holds the conclusion back on purpose
//
// The em dash is optional and so is everything after it. Only "off" is
// accepted: rules are on by default, so there is nothing for "on" to mean, and
// silently accepting it would let a card think it had enabled something.
const gateSyntax = `@gate <rule_id> off — <why>`

// Gates parses the card's @gate lines. The error names every problem rather
// than the first, because a card is edited by hand and one load should report
// all of what needs fixing.
//
// A card naming a rule that does not exist is an error and not a warning. The
// likely cause is a typo in a rule id, and a silently ignored exemption reads
// exactly like a working one: the card says it is exempt, the gate marks it
// down anyway, and nothing explains why.
func (c *Card) Gates() ([]GateDirective, error) {
	var out []GateDirective
	var bad []string
	seen := map[string]bool{}
	for _, line := range c.Gate {
		d, err := parseGate(line)
		if err != nil {
			if f := strings.Fields(line); len(f) > 0 && c.gateOnOwnShapeRule(f[0]) {
				bad = append(bad, fmt.Sprintf(
					"%q is this card's own @shape rule, not a built-in; delete the @shape line instead of gating it", f[0]))
				continue
			}
			bad = append(bad, err.Error())
			continue
		}
		if seen[d.Rule] {
			bad = append(bad, fmt.Sprintf("%s named twice", d.Rule))
			continue
		}
		seen[d.Rule] = true
		out = append(out, d)
	}
	if len(bad) > 0 {
		return out, fmt.Errorf("@gate: %s\n  syntax: %s", strings.Join(bad, "; "), gateSyntax)
	}
	return out, nil
}

func parseGate(line string) (GateDirective, error) {
	// Split the reason off first so a rule id is never confused with a word
	// from the prose. Both dash forms, because a card is typed by a person.
	body, why := line, ""
	for _, sep := range []string{" — ", " - ", " -- "} {
		if i := strings.Index(line, sep); i >= 0 {
			body, why = line[:i], strings.TrimSpace(line[i+len(sep):])
			break
		}
	}
	f := strings.Fields(body)
	switch {
	case len(f) == 0:
		return GateDirective{}, fmt.Errorf("empty line")
	case len(f) == 1:
		return GateDirective{}, fmt.Errorf("%q names a rule and no state, want %q off", f[0], f[0])
	case len(f) > 2:
		return GateDirective{}, fmt.Errorf("%q has %d words before the reason, want two", body, len(f))
	case f[1] != "off":
		return GateDirective{}, fmt.Errorf("%q is not a state; only \"off\" is accepted, since rules are on by default", f[1])
	case !RuleIDs[f[0]]:
		return GateDirective{}, fmt.Errorf("%q is not a built-in rule", f[0])
	}
	return GateDirective{Rule: f[0], Why: why}, nil
}

// gateOnOwnShapeRule is the confusable case: @gate names a rule this card
// defined in @shape. It parses as an unknown built-in and the bare "not a
// built-in rule" would send the author hunting for a typo in a correctly
// spelled id.
//
// It stays an error rather than becoming a feature. A card that wrote the rule
// can delete the line; there is no inheritance here for an exemption to be
// useful against, and two lines that cancel out is a worse way to say nothing
// than saying nothing.
func (c *Card) gateOnOwnShapeRule(id string) bool {
	for _, line := range c.Shape {
		if r, err := parseShape(line); err == nil && r.ID == id {
			return true
		}
	}
	return false
}

// exempt is the set of rule ids this card has turned off. A malformed line is
// skipped here and reported by Gates at load, so a card that fails validation
// is scored as though the broken line were absent rather than not scored.
func (c *Card) exempt() map[string]bool {
	if len(c.Gate) == 0 {
		return nil
	}
	out := map[string]bool{}
	for _, line := range c.Gate {
		if d, err := parseGate(line); err == nil {
			out[d.Rule] = true
		}
	}
	return out
}

// Allow drops violations from rules this card is exempt from.
//
// Filtering the results rather than skipping the detectors keeps every rule
// running against every reply, which matters for the backfill: a card can turn
// a rule off for scoring and the same corpus still says how often it would
// have fired.
func (c *Card) Allow(v []Violation) []Violation {
	ex := c.exempt()
	if len(ex) == 0 {
		return v
	}
	kept := v[:0]
	for _, x := range v {
		if !ex[x.RuleID] {
			kept = append(kept, x)
		}
	}
	return kept
}
