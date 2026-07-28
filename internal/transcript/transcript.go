// Package transcript reads Claude Code session JSONL.
package transcript

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type block struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type record struct {
	Type    string  `json:"type"`
	Message message `json:"message"`
}

// assistantText pulls the assistant's prose out of one JSONL line, or returns
// "" for anything else. Tool calls and tool results are skipped: the gate
// judges what the user reads, not what the harness carries.
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

// tailChunk is how much LastAssistantText reads per backward step. One
// megabyte holds the last assistant turn in all but pathological transcripts;
// the loop covers the rest. A var so tests can shrink it and exercise the
// chunk boundary against a small fixture.
var tailChunk = 1 << 20

// maxCarry bounds the partial line held over when a single JSONL record spans
// more than one chunk. The longest line across the five largest transcripts on
// the development machine was 5.8 MB, so 64 MB leaves an order of magnitude
// while stopping a corrupt file from being pulled entirely into memory.
const maxCarry = 64 << 20

// LastAssistantText returns the most recent assistant prose turn.
//
// It scans backward from the end of the file, so the cost tracks how far back
// that turn sits rather than the size of the session. Reading forward cost 28.3
// seconds on a 2.6 GB transcript, against the 10-second timeout the README
// recommends for the Stop hook — the gate died silently on exactly the
// long-running sessions it exists for. Backward, the same file takes under a
// millisecond and reports the same turn.
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

	var carry []byte // a partial line continuing past the low end of the window
	for off := fi.Size(); off > 0; {
		n := int64(tailChunk)
		if n > off {
			n = off
		}
		off -= n

		buf := make([]byte, n, n+int64(len(carry)))
		if _, err := f.ReadAt(buf, off); err != nil {
			return "", err
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
			if t := assistantText(lines[i]); t != "" {
				return t, nil
			}
		}
		if off == 0 {
			break
		}
		carry = lines[0]
		if len(carry) > maxCarry {
			return "", fmt.Errorf("%s: a single line exceeds %d bytes", path, maxCarry)
		}
	}
	return "", nil
}

// AllAssistantText returns every assistant prose turn, oldest first. Used by
// --backfill to get a hit rate over a whole session at once.
//
// A read error returns the turns collected so far alongside it. A truncated
// scan is still a usable sample, and the caller decides — reporting nothing
// because the last line was malformed was the old behaviour and it was wrong.
func AllAssistantText(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	s := bufio.NewScanner(f)
	// Transcript lines routinely exceed the default 64 KiB when a turn carries
	// a large tool result. The cap matches maxCarry so both readers give up at
	// the same size rather than at two different ones.
	s.Buffer(make([]byte, 0, 1<<20), maxCarry)
	for s.Scan() {
		if t := assistantText(s.Bytes()); t != "" {
			out = append(out, t)
		}
	}
	return out, s.Err()
}
