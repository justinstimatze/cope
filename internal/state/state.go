// Package state keeps what cope has observed about a session between hook
// invocations.
//
// Hooks are separate processes with no shared memory, so anything cope wants to
// know across a turn boundary has to be written down. The Stop hook records
// what the reply scored; the next UserPromptSubmit reads it back and decides
// which parts of the card to re-state. That loop is the difference between a
// card and a file you pasted — the file says the same thing every time, and
// this can say what you have actually been doing.
package state

import (
	"encoding/json"
	"hash/fnv"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// LaneLoop mirrors transcript.LaneLoop, which this package cannot name without
// depending on the JSONL reader to describe a fact about a session.
const LaneLoop = "loop"

// recentTurns bounds the window the injection reasons over. A tic from turn 3
// of an eighty-turn session is not what the writer needs reminding of, and an
// unbounded tally would let it outvote what is happening now.
const recentTurns = 20

// Turn is one scored reply.
//
// Key is the transcript id of the prompt that opened the turn. It exists
// because a turn is scored twice: once at Stop, against a transcript missing
// the reply's closing block, and again at the next prompt once the file is
// complete. Without it the second pass would count as a second turn and the
// rule tally would double.
// Phrases is the turn's scaffolding shingles as 32-bit sums — see
// scan.Phrases. Fingerprints rather than text, so the window can notice that a
// construction recurred without this file becoming a copy of the conversation.
type Turn struct {
	At      time.Time `json:"at"`
	Key     string    `json:"key,omitempty"`
	Rules   []string  `json:"rules"`
	Chars   int       `json:"chars"`
	Phrases []uint32  `json:"phrases,omitempty"`
}

// Session is the rolling record for one Claude Code session.
//
// Windows, Arm, InjectChars and InjectNamed belong to the mid-session
// experiment and stay empty unless it is running. See NextWindow.
type Session struct {
	ID      string `json:"id"`
	Turns   int    `json:"turns"`
	Chars   int    `json:"chars"`
	Recent  []Turn `json:"recent"`
	Windows int    `json:"windows,omitempty"`
	Arm     string `json:"arm,omitempty"`

	// Lane is how the most recent turn was asked for — see transcript.LaneLoop.
	// It is on the session rather than the turn because the reader it decides
	// things for is the one arriving at the next prompt, and a session that has
	// entered a loop stays in one until a person types something.
	Lane string `json:"lane,omitempty"`

	// The shape of the last window's injection, carried so that a turn
	// written under it can be logged with what it actually saw.
	InjectChars int      `json:"inject_chars,omitempty"`
	InjectNamed []string `json:"inject_named,omitempty"`
}

// The sides of the experiment. Every arm carries the SessionStart card, so what
// is being compared is the sliced mid-session injection against nothing —
// which is the question, since the card alone is what a pasted file already
// does.
//
// ArmPositive is the same slice with its anti-patterns removed: the better
// halves of the wrong/better pairs and the passing examples, and none of the
// prohibitions or the worse lines. It exists to separate "a reminder helped"
// from "a reminder that quotes the bad version put the bad version in reach".
// It is shorter than ArmInject by construction, which is a second difference
// riding along with the first — see InjectChars, which is logged so the two can
// be told apart afterwards.
const (
	ArmInject   = "inject"
	ArmHold     = "hold"
	ArmPositive = "positive"
)

// DefaultArms is the pair the experiment runs unless told otherwise. Changing
// this changes what accumulated data means, so a third arm is opt-in: turns
// already collected under two arms stay comparable to turns collected later.
var DefaultArms = []string{ArmInject, ArmHold}

// NextWindow advances to the next refresher window and returns the arm it
// lands on.
//
// Alternating rather than tossing a coin. It balances exactly, it spaces the
// arms evenly through the session, and it cannot drift with time-in-session
// the way independent draws can. The session id decides which arm a session
// opens on, so the first window is not always the same one.
//
// Randomising by session instead would leave project and task as uncontrolled
// variance across an N of a few dozen. Alternating within a session holds both
// fixed. Guidance injected in one window does carry into the next, and that
// carry-over pulls the arms together, so the design can understate a real
// effect and cannot manufacture one.
//
// An empty arm list falls back to DefaultArms rather than dividing by zero: a
// miswired flag should cost the experiment, not the prompt.
func (s *Session) NextWindow(arms []string) string {
	if len(arms) == 0 {
		arms = DefaultArms
	}
	s.Windows++
	s.Arm = arms[(int(hash32(s.ID))+s.Windows)%len(arms)]
	return s.Arm
}

// RecordInjection stores the shape of what the last window injected, so a turn
// written under it can be logged with the reminder that produced it.
//
// Without this an arm is only a name. Two arms can differ in what they say, how
// long they are, and which rules they speak to, and a rate ratio between them
// cannot say which of the three did the work.
func (s *Session) RecordInjection(chars int, named []string) {
	s.InjectChars = chars
	s.InjectNamed = named
}

func hash32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func path(dir, id string) string { return filepath.Join(dir, "session-"+id+".json") }

// Load returns the stored session, or a fresh one. A corrupt or unreadable file
// yields a fresh session rather than an error: this is advisory state, and
// losing it costs one turn of personalisation, never the hook.
func Load(dir, id string) *Session {
	s := &Session{ID: id}
	raw, err := os.ReadFile(path(dir, id))
	if err != nil {
		return s
	}
	if err := json.Unmarshal(raw, s); err != nil || s.ID != id {
		return &Session{ID: id}
	}
	return s
}

// Record files one scored reply, replacing the newest entry when key says this
// is the same turn scored again, and appending otherwise.
//
// An empty key always appends. That is the honest reading: without an id there
// is no way to tell a rescore from a new turn, and counting one turn twice is
// the cheaper mistake than silently overwriting a real one.
//
// phrases is variadic because most callers have none: only the two hook paths
// that hold the reply's text can fingerprint it, and everything else — the
// tests, any caller tallying rules alone — would otherwise pass a nil to say
// so on every line.
func (s *Session) Record(key string, ruleIDs []string, chars int, now time.Time, phrases ...uint32) {
	if last, ok := s.newest(key); ok {
		s.Chars += chars - last.Chars
		last.At, last.Rules, last.Chars, last.Phrases = now, ruleIDs, chars, phrases
		return
	}
	s.Turns++
	s.Chars += chars
	s.Recent = append(s.Recent, Turn{At: now, Key: key, Rules: ruleIDs, Chars: chars, Phrases: phrases})
	if len(s.Recent) > recentTurns {
		s.Recent = s.Recent[len(s.Recent)-recentTurns:]
	}
}

// newest returns the most recent turn when it carries this key. Only the newest
// is eligible: a rescore always targets the turn that just happened, and
// rewriting an older entry would mean the window no longer reads in the order
// the replies were written.
func (s *Session) newest(key string) (*Turn, bool) {
	if key == "" || len(s.Recent) == 0 {
		return nil, false
	}
	last := &s.Recent[len(s.Recent)-1]
	if last.Key != key {
		return nil, false
	}
	return last, true
}

// Scored reports how many characters of this turn have already been scored, so
// the caller can tell whether the transcript has grown since. Absent means the
// turn has not been recorded — the Stop hook may not be installed, or this is
// the first prompt after a resume.
func (s *Session) Scored(key string) (int, bool) {
	if last, ok := s.newest(key); ok {
		return last.Chars, true
	}
	return 0, false
}

// Save writes the session. 0600 because Chars and the rule tally describe the
// contents of a private conversation, even though no prose is stored.
func (s *Session) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	raw, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path(dir, s.ID), raw, 0o600)
}

// PriorPhrases tallies how many turns in the window used each phrase, skipping
// the turn under key.
//
// The skip is what makes a rescore safe. A turn is scored twice — once at Stop
// against an unfinished transcript, once on the next prompt against the
// complete one — and without this the second pass would find every phrase from
// the first pass sitting in the window and report the reply for repeating
// itself.
func (s *Session) PriorPhrases(key string) map[uint32]int {
	tally := map[uint32]int{}
	for _, t := range s.Recent {
		if key != "" && t.Key == key {
			continue
		}
		for _, p := range t.Phrases {
			tally[p]++
		}
	}
	return tally
}

// Count is one rule and how often it fired inside the window.
type Count struct {
	Rule string
	N    int
}

// Top returns the rules firing most in the window, busiest first. Ties break on
// the rule name so the injection is stable between turns rather than shuffling
// under the reader.
func (s *Session) Top(n int) []Count {
	tally := map[string]int{}
	for _, t := range s.Recent {
		for _, r := range t.Rules {
			tally[r]++
		}
	}
	out := make([]Count, 0, len(tally))
	for r, c := range tally {
		out = append(out, Count{r, c})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].N != out[j].N {
			return out[i].N > out[j].N
		}
		return out[i].Rule < out[j].Rule
	})
	if n > 0 && len(out) > n {
		out = out[:n]
	}
	return out
}

// Facts renders the session as the fact set the card's @when gates read.
//
// position is the injection point and renders as at_start / at_prompt. It was
// pos_ until a generated README read "pos_prompt" as "the positive-prompt
// fact", which is what an ambiguous abbreviation earns.
//
// topRules of zero or less means no rule facts at all. Top reads the same
// argument as "no limit", which is the opposite, and a caller that wanted an
// empty slice would otherwise get every rule in the window. The two meanings
// are kept apart here rather than in the reader's head.
func (s *Session) Facts(position string, topRules int) map[string]bool {
	f := map[string]bool{}
	// mode_* is the third namespace: not where the injection sits and not what
	// the writer has been doing, but what kind of session this is. It replaces
	// the position rather than joining it, because every at_prompt item is
	// written for a reader who is about to answer — the CONTINUE TEST and the
	// pair whose better half ends in a question. Carrying those into a lane
	// that must not end by asking would put the card on both sides of the same
	// question, which is the overconstraint failure it exists to avoid.
	switch {
	case s.Lane == LaneLoop:
		f["mode_loop"] = true
	case position != "":
		f["at_"+position] = true
	}

	if topRules <= 0 {
		return f
	}
	for _, c := range s.Top(topRules) {
		f["rule_"+c.Rule] = true
	}
	return f
}

// Prune drops session files untouched for a week. Sessions end without a
// goodbye, so nothing else would ever remove these.
func Prune(dir string, olderThan time.Duration, now time.Time) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if len(e.Name()) < 8 || e.Name()[:8] != "session-" {
			continue
		}
		if fi, err := e.Info(); err == nil && now.Sub(fi.ModTime()) > olderThan {
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}
