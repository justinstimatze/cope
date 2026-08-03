// Package transcript reads Claude Code session JSONL.
//
// The unit that matters here is the turn, not the record. One reply is written
// as several assistant records — prose, a tool call, more prose, another tool
// call — and reading the last record alone gets a mid-turn aside rather than
// the reply. Measured on this repo's own session, 204 prose records covered 63
// turns, and 108 of those records were under 200 characters.
package transcript

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type block struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type message struct {
	Role    string          `json:"role"`
	Model   string          `json:"model"`
	Content json.RawMessage `json:"content"`
}

type record struct {
	Type       string  `json:"type"`
	Subtype    string  `json:"subtype"`
	UUID       string  `json:"uuid"`
	ParentUUID string  `json:"parentUuid"`
	APIError   bool    `json:"isApiErrorMessage"`
	Message    message `json:"message"`
}

// syntheticModel is what Claude Code writes in an assistant record it composed
// itself rather than received from the model: API error banners, interrupt
// notices, and anything else the harness needs to put in the reply position.
const syntheticModel = "<synthetic>"

// assistantText pulls the assistant's prose out of one JSONL line, or returns
// "" for anything else. Tool calls and tool results are skipped: the gate
// judges what the user reads, not what the harness carries.
//
// Harness-composed records are skipped too, and finding that out cost a
// measurement. "API Error: 500 Internal server error. This is a server-side
// issue, usually temporary — try again in a moment." is stored as an assistant
// turn, verbatim, once per retry. It was the second most repeated phrase on a
// 27,962-turn backfill and read as the model saying the same thing over and
// over, which is exactly what cross_turn_repeat was built to catch — except
// nobody wrote it. Every rule in this repo was scoring those, so every rate
// ever measured here included them.
func assistantText(line []byte) string {
	if len(bytes.TrimSpace(line)) == 0 {
		return ""
	}
	var r record
	if err := json.Unmarshal(line, &r); err != nil {
		return "" // a malformed line should not blind the gate
	}
	if r.Type != "assistant" || r.Message.Role != "assistant" {
		return ""
	}
	// Both markers, because they are set by different code paths and either one
	// alone would let the other kind through.
	if r.APIError || r.Message.Model == syntheticModel {
		return ""
	}
	var blocks []block
	if err := json.Unmarshal(r.Message.Content, &blocks); err != nil {
		return ""
	}
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" && strings.TrimSpace(b.Text) != "" {
			parts = append(parts, b.Text)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

// The two ways a reply gets asked for, and the reason cope needs to tell them
// apart: they have different readers.
//
// An interactive turn is read now, by someone who can answer it, and its
// ending is a decision surface. A loop turn is read later, by someone
// reconstructing what happened while they were gone, and nothing in it can be
// answered at the moment it is written — the next thing to read it is the next
// iteration. Half the ending rules invert across that line.
const (
	LaneInteractive = "interactive"
	LaneLoop        = "loop"
)

// loopPrompt matches a loop-lane prompt by its own text: the slash command that
// opens the loop, and the sentinel the dynamic-pacing variant sends itself.
//
// /goal is here with /loop because the difference between them is the stopping
// condition, not the reader. Both write for somebody who is not present.
//
// The continuing wakeups are not here, because they are not user records at
// all — see fireUUID.
var loopPrompt = regexp.MustCompile(
	`<command-name>/(?:loop|goal)</command-name>|<<autonomous-loop`)

func laneOf(prompt string) string {
	if loopPrompt.MatchString(prompt) {
		return LaneLoop
	}
	return LaneInteractive
}

// scheduledFire is the subtype Claude Code writes when a timer wakes the
// session. Every wakeup after the first is one of these, so a loop that ran
// overnight is invisible without it — the wakeup was read as an ordinary
// prompt, and the whole loop lane came back empty on a session with 423 of
// them.
const scheduledFire = "scheduled_task_fire"

// fireUUID returns the record id of a scheduled wakeup, or "".
//
// The wakeup is a system record and the prompt it delivers is a separate user
// record underneath it, chained by parentUuid. Matching on that chain rather
// than on adjacency is what keeps this correct if an attachment or a hook
// record ever lands between the two.
func fireUUID(line []byte) string {
	if len(bytes.TrimSpace(line)) == 0 {
		return ""
	}
	var r record
	if err := json.Unmarshal(line, &r); err != nil {
		return ""
	}
	if r.Type == "system" && r.Subtype == scheduledFire {
		return r.UUID
	}
	return ""
}

// userPromptRec reports whether a line is a real prompt from the user, and
// returns its record id, its parent id, and its text. The parent is what the
// lane chains a prompt back to the wakeup that delivered it.
//
// The distinction is the whole reason turns can be assembled at all. Tool
// results are user records too, and a turn carrying six tool calls writes six
// of them; treating those as turn boundaries is what cut replies into
// fragments. A compaction summary is a user record with string content and is
// a genuine boundary, so it needs no special case.
func userPromptRec(line []byte) (id, parent, text string, ok bool) {
	if len(bytes.TrimSpace(line)) == 0 {
		return "", "", "", false
	}
	var r record
	if err := json.Unmarshal(line, &r); err != nil {
		return "", "", "", false
	}
	if r.Type != "user" || r.Message.Role != "user" {
		return "", "", "", false
	}
	var s string
	if err := json.Unmarshal(r.Message.Content, &s); err == nil {
		return r.UUID, r.ParentUUID, s, strings.TrimSpace(s) != ""
	}
	var blocks []block
	if err := json.Unmarshal(r.Message.Content, &blocks); err != nil {
		return "", "", "", false
	}
	for _, b := range blocks {
		if b.Type == "tool_result" {
			return "", "", "", false
		}
	}
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" && strings.TrimSpace(b.Text) != "" {
			parts = append(parts, b.Text)
		}
	}
	if len(parts) == 0 {
		return "", "", "", false
	}
	return r.UUID, r.ParentUUID, strings.Join(parts, "\n"), true
}

// tailChunk is how much the backward reader takes per step. One megabyte holds
// the last turn in all but pathological transcripts; the loop covers the rest.
// A var so tests can shrink it and exercise the chunk boundary against a small
// fixture.
var tailChunk = 1 << 20

// maxCarry bounds the partial line held over when a single JSONL record spans
// more than one chunk. The longest line across the five largest transcripts on
// the development machine was 5.8 MB, so 64 MB leaves an order of magnitude
// while stopping a corrupt file from being pulled entirely into memory.
const maxCarry = 64 << 20

// maxTurnBlocks bounds how far back turn assembly walks when it finds no
// boundary. A transcript whose user records are all unreadable would otherwise
// pull the whole session in as one turn.
const maxTurnBlocks = 512

// maxFireLookback bounds how far past the boundary prompt the backward reader
// looks for the wakeup that delivered it. The two are adjacent in every
// transcript checked, so this only has to survive an attachment or a hook
// record landing between them.
const maxFireLookback = 32

// turnSep joins the prose blocks of one turn. A blank line, because that is how
// the reader sees them and because every shape rule splits on paragraphs.
//
// It also makes the turn grow by append only, which the catch-up pass depends
// on: the text the Stop hook scored is a byte prefix of the finished turn, so
// rescoring that prefix reproduces exactly what was already reported.
const turnSep = "\n\n"

// eachLineBackward calls fn with each line of f, last first, and stops when fn
// returns false or the start of the file is reached.
//
// Backward because the cost then tracks how far back the turn sits rather than
// the size of the session. Reading forward cost 28.3 seconds on a 2.6 GB
// transcript, against the 10-second timeout the README recommends for the Stop
// hook — the gate died silently on exactly the long-running sessions it exists
// for. Backward, the same file takes under a millisecond.
func eachLineBackward(f *os.File, size int64, fn func(line []byte) bool) error {
	var carry []byte // a partial line continuing past the low end of the window
	for off := size; off > 0; {
		n := int64(tailChunk)
		if n > off {
			n = off
		}
		off -= n

		buf := make([]byte, n, n+int64(len(carry)))
		if _, err := f.ReadAt(buf, off); err != nil {
			return err
		}
		buf = append(buf, carry...)

		lines := bytes.Split(buf, []byte{'\n'})
		// lines[0] is a whole line only when the window reaches the start of
		// the file. Anywhere else it is a fragment to carry down.
		first := 1
		if off == 0 {
			first = 0
		}
		for i := len(lines) - 1; i >= first; i-- {
			if !fn(lines[i]) {
				return nil
			}
		}
		if off == 0 {
			return nil
		}
		carry = lines[0]
		if len(carry) > maxCarry {
			return fmt.Errorf("a single line exceeds %d bytes", maxCarry)
		}
	}
	return nil
}

// Turn is one reply: every assistant prose block written since the user's last
// real prompt, in the order it was written.
//
// Key is the id of the prompt that opened the turn. It exists so a second read
// of the same turn can be recognised as a rescore rather than counted as a new
// turn, which is what the Stop hook's blind spot requires — see LastTurn.
// Lane is which reader the reply was written for — see LaneInteractive. It is
// read off the prompt that opened the turn, so a turn with no prompt above it
// reads as interactive: that is the conservative default, since the loop lane
// is where rules get suppressed.
type Turn struct {
	Key  string
	Text string
	Lane string
}

// LastTurn returns the most recent turn that has a reply in it.
//
// It reads what is on disk, and at Stop time that is not the whole reply.
// Claude Code appends a reply's closing block after the Stop hook has already
// run, so the hook sees the turn without its ending — which is where
// dangling_end, buried_decision and the last ask all live. Across six live
// sessions on 2026-07-28 every one of the 17 recorded turns was scored against
// text that ended early, and in four of them the text belonged to an earlier
// turn entirely. The next UserPromptSubmit reads a finished file and scores the
// rest; Key is what lets the two passes agree they are the same turn.
//
// A prompt with no reply under it is walked past rather than returned empty.
// Whether the prompt being submitted is already in the file when
// UserPromptSubmit fires is not something a hook gets to know, and under the
// order where it is, stopping at the first boundary would hand back nothing
// and the catch-up would quietly never run.
func LastTurn(path string) (Turn, error) {
	f, err := os.Open(path)
	if err != nil {
		return Turn{}, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return Turn{}, err
	}

	var parts []string // newest first until reversed below
	var key, lane string
	// Set once the boundary prompt is found, to keep walking a little further
	// for the wakeup that delivered it. Reading backward, the system record
	// sits below the prompt, so the lane is not known when the prompt is.
	wantFire, budget := "", maxFireLookback
	err = eachLineBackward(f, fi.Size(), func(line []byte) bool {
		if wantFire != "" {
			if fireUUID(line) == wantFire {
				lane = LaneLoop
				return false
			}
			budget--
			return budget > 0
		}
		if id, parent, prompt, ok := userPromptRec(line); ok {
			if len(parts) == 0 {
				return true // a prompt still waiting for its reply
			}
			key, lane = id, laneOf(prompt)
			if lane == LaneLoop || parent == "" {
				return false
			}
			wantFire = parent
			return true
		}
		if t := assistantText(line); t != "" {
			parts = append(parts, t)
			return len(parts) < maxTurnBlocks
		}
		return true
	})
	if err != nil {
		return Turn{}, fmt.Errorf("%s: %w", path, err)
	}
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	if lane == "" {
		lane = LaneInteractive
	}
	return Turn{Key: key, Text: strings.Join(parts, turnSep), Lane: lane}, nil
}

// LastAssistantText returns the most recent assistant prose record.
//
// One record, not one turn: this is the raw primitive, kept because the
// backward reader is worth testing at record granularity. Nothing scores
// against it — the gate reads LastTurn, for the reasons in this file's header.
func LastAssistantText(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	var out string
	err = eachLineBackward(f, fi.Size(), func(line []byte) bool {
		if t := assistantText(line); t != "" {
			out = t
			return false
		}
		return true
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return out, nil
}

// scanner returns a line scanner sized for transcript records, which routinely
// exceed the default 64 KiB when a turn carries a large tool result. The cap
// matches maxCarry so both readers give up at the same size rather than at two
// different ones.
func scanner(f *os.File) *bufio.Scanner {
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 1<<20), maxCarry)
	return s
}

// AllTurns returns every turn in the file, oldest first. Used by --backfill to
// get a rate over a whole session at once.
//
// A read error returns the turns collected so far alongside it. A truncated
// scan is still a usable sample, and the caller decides — reporting nothing
// because the last line was malformed was the old behaviour and it was wrong.
//
// Reading forward, the prompt arrives before the reply it opens, so the lane
// has to be held until the turn under it closes.
func AllTurns(path string) ([]Turn, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Turn
	var buf []string
	key, lane := "", LaneInteractive
	lastFire := ""
	flush := func() {
		if len(buf) > 0 {
			out = append(out, Turn{Key: key, Text: strings.Join(buf, turnSep), Lane: lane})
			buf = nil
		}
	}
	s := scanner(f)
	for s.Scan() {
		line := s.Bytes()
		if t := assistantText(line); t != "" {
			buf = append(buf, t)
			continue
		}
		if id := fireUUID(line); id != "" {
			lastFire = id
			continue
		}
		if id, parent, prompt, ok := userPromptRec(line); ok {
			flush() // the blocks above belong to the previous prompt
			key, lane = id, laneOf(prompt)
			if parent != "" && parent == lastFire {
				lane = LaneLoop
			}
		}
	}
	flush()
	return out, s.Err()
}

// AllAssistantText returns every assistant prose record, oldest first. Record
// granularity, so it counts a turn's mid-reply asides separately — which is why
// --backfill reads AllTurns instead. Kept as the forward reference the backward
// reader's tests check themselves against.
func AllAssistantText(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	s := scanner(f)
	for s.Scan() {
		if t := assistantText(s.Bytes()); t != "" {
			out = append(out, t)
		}
	}
	return out, s.Err()
}
