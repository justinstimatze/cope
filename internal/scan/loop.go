package scan

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Rules for the lane where nobody is reading yet.
//
// Addy Osmani's June 2026 "Loop Engineering" is the frame: you stop prompting
// the agent and design the system that prompts it, and the reply stops being a
// message and becomes a record somebody reconstructs the session from later.
// Two sentences from it set what a rule here has to check. "A loop running
// unattended is also a loop making mistakes unattended", and — on the split
// between the agent that does the work and the one that grades it — "'done' is
// a claim and not a proof."
//
// So the loop lane inverts half the ending rules and adds one. ask_not_last,
// forked_end, dangling_end and buried_decision all measure a decision surface
// for a reader who is present; in this lane there is nobody to decide, and a
// report that correctly names what it left open and stops would fail three of
// them. What replaces them is the claim check: a report saying the work is done
// has to say what it ran.

// LaneLoop mirrors transcript.LaneLoop. The two packages do not import each
// other — scan reads prose and transcript reads JSONL — and one string in two
// places is cheaper than a package that exists to hold it.
const LaneLoop = "loop"

// LaneInteractive is the default: a turn written to somebody about to answer
// it. Named so a caller passing a lane by string has all three, rather than
// two constants and an empty string standing for the third.
const LaneInteractive = "interactive"

// ValidLane rejects a lane nobody implements. CheckLane treats every unknown
// string as interactive, which is the right default inside the binary and the
// wrong one at a flag: a caller who typed --check-lane extenal would be told
// nothing and would silently score four rules the real hook drops.
func ValidLane(lane string) error {
	switch lane {
	case "", LaneInteractive, LaneLoop, LaneExternal:
		return nil
	}
	return fmt.Errorf("unknown lane %q — one of %s, %s, %s", lane, LaneInteractive, LaneLoop, LaneExternal)
}

// needsAnswerableReader are the rules whose subject is a reader who can answer.
// They are not wrong, they are asking a question some lanes do not pose — the
// loop lane, where nobody is reading yet, and the external lane, where the
// reader is holding a ticket rather than a turn. Both drop the same four for
// the same reason, which is why the name is about the assumption rather than
// about either lane.
var needsAnswerableReader = map[string]bool{
	"ask_not_last":    true,
	"forked_end":      true,
	"dangling_end":    true,
	"buried_decision": true,
}

// NeedsAnswerableReader is that set, sorted, for anything describing a lane
// rather than running one. The docs prompt used to carry the list written out
// by hand in two places and the count in a third, which is three copies of a
// fact that lives here.
func NeedsAnswerableReader() []string {
	out := make([]string, 0, len(needsAnswerableReader))
	for id := range needsAnswerableReader {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

var (
	// A claim that work is finished, working, or correct.
	doneClaim = regexp.MustCompile(`(?i)\b(?:` +
		`(?:tests?|suite|build|lint|vet|checks?|ci|everything|it|that|this|all of it)\s+` +
		`(?:now\s+)?(?:passes|pass|passed|builds|compiles|works|is green|are green|succeeded)|` +
		`(?:is|are|it's|its)\s+(?:now\s+)?(?:fixed|working|done|green|resolved|verified|passing)|` +
		`(?:i|i've|i have)\s+(?:verified|confirmed)\b|` +
		`no (?:errors|failures|regressions)\b` +
		`)`)

	// Evidence that something was actually executed: a command, an exit status,
	// a count, a path. Deliberately generous — the rule is looking for the
	// absence of any anchor at all, not for a well-formed citation. A report
	// that names a file and a number has done the thing this asks for even if
	// it never prints a command line.
	ranEvidence = regexp.MustCompile(`(?i)(?:` +
		"`[^`\n]+`|" + // inline code, which Redact blanks — so this runs on the raw text
		"```|" +
		`\bgo (?:test|build|vet|run)\b|\b(?:make|npm|pnpm|yarn|cargo|pytest|python3?|bash|git)\s+[a-z]|` +
		`\bexit (?:code )?\d|\b\d+\s*/\s*\d+\b|\b\d+\s+(?:tests?|checks?|files?|cases?)\b|` +
		`\bok\s+github\.com|\bPASS\b|\bFAIL\b|` +
		`[\w./-]+\.(?:go|py|ts|tsx|js|md|json|ya?ml|sh|rs|java|rb):\d+|` +
		`[\w./-]+\.(?:go|py|ts|tsx|js|md|json|ya?ml|sh|rs|java|rb)\b` +
		`)`)

	// The report saying it did not run the thing. An honest "written but not
	// executed" is the behaviour the card asks for and must not be flagged as
	// an unbacked claim.
	admitsUnrun = regexp.MustCompile(`(?i)\b(?:not (?:yet )?(?:run|executed|tested|verified)|` +
		`hasn'?t (?:been )?(?:run|executed|tested|verified)|` +
		`without running|didn'?t run|did not run|unverified|untested|` +
		`i have not (?:run|verified)|nothing was run)\b`)
)

// UnverifiedDone flags a loop report claiming the work is finished with nothing
// on the page that could have shown it.
//
// The check is deliberately whole-reply rather than per-sentence. A report that
// says "tests pass" in its first line and prints `go test ./...` in its last has
// done what this asks; requiring the evidence to sit beside the claim would be a
// formatting rule wearing a verification rule's name.
//
// It is warn-only and it will miss a lie. Naming a command is not running one,
// and nothing here can tell a real transcript from an invented one — that is
// the reader's job, and Osmani is explicit that it stays theirs. What this
// catches is the cheaper and far more common failure: a summary that asserts
// success because the work felt finished, with no anchor anyone could check.
func UnverifiedDone(text string) []Violation {
	return unverifiedDone(text, rawParagraphs(Redact(text)))
}

func unverifiedDone(raw string, ps []string) []Violation {
	if ranEvidence.MatchString(raw) || admitsUnrun.MatchString(raw) {
		return nil
	}
	for _, p := range ps {
		for _, s := range splitSentences(p) {
			s = strings.TrimSpace(s)
			if doneClaim.MatchString(s) {
				return []Violation{{
					RuleID:  "unverified_done",
					Action:  "warn",
					Why:     "says the work is done with nothing on the page that could have shown it — no command, no count, no file",
					Matched: s,
				}}
			}
		}
	}
	return nil
}

// LoopAsk flags a loop report that ends by asking.
//
// Nobody is there. The question goes into a log, and the next thing to read it
// is the next iteration, which will take it as an instruction to itself. An
// unattended run that needs a decision should say what it stopped at and why,
// so the decision is waiting when somebody arrives rather than phrased as a
// turn in a conversation that is not happening.
//
// Only the last block, and only a real question. A loop report often quotes the
// question it is resolving, and mid-report rhetorical questions are how a plan
// gets laid out.
func LoopAsk(text string) []Violation {
	return loopAsk(rawParagraphs(Redact(text)))
}

func loopAsk(ps []string) []Violation {
	if len(ps) == 0 {
		return nil
	}
	last := ps[len(ps)-1]
	for _, line := range strings.Split(last, "\n") {
		if listItem.MatchString(line) {
			continue
		}
		for _, s := range splitSentences(line) {
			s = strings.TrimSpace(s)
			if s == "" || !askPunct.MatchString(s) && !askPhrase.MatchString(s) {
				continue
			}
			return []Violation{{
				RuleID:  "loop_ask",
				Action:  "warn",
				Why:     "an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself",
				Matched: s,
			}}
		}
	}
	return nil
}

// CheckLane runs the card over one piece of prose with the lane's rules.
//
// Interactive is Check unchanged, so nothing measured before this existed
// changes meaning. The loop lane drops the four rules that assume a reader who
// can answer and adds the two that assume one who cannot. The external lane
// drops the same four and adds nothing — see LaneExternal.
func (c *Card) CheckLane(text string, minCV float64, lane string) []Violation {
	v := c.Check(text, minCV)
	if lane != LaneLoop && lane != LaneExternal {
		return v
	}
	kept := v[:0]
	for _, x := range v {
		if !needsAnswerableReader[x.RuleID] {
			kept = append(kept, x)
		}
	}
	if lane == LaneExternal {
		// Allow runs again because the filter above can only remove, and a card
		// that gated one of the four has already had it dropped twice over.
		return c.Allow(kept)
	}
	raw := rawParagraphs(Redact(text))
	kept = append(kept, unverifiedDone(text, raw)...)
	// Allow runs again over the two loop-only rules, which Check never saw.
	return c.Allow(append(kept, loopAsk(raw)...))
}
