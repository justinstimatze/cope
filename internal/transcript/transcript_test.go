package transcript

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assistantLine(text string) string {
	b, err := json.Marshal(map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"role":    "assistant",
			"content": []map[string]any{{"type": "text", "text": text}},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func toolLine(payload string) string {
	b, err := json.Marshal(map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"role":    "assistant",
			"content": []map[string]any{{"type": "tool_use", "name": "Bash", "input": payload}},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func userLine(text string) string {
	b, err := json.Marshal(map[string]any{
		"type":    "user",
		"message": map[string]any{"role": "user", "content": text},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func write(t *testing.T, lines []string, trailingNewline bool) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "session.jsonl")
	body := strings.Join(lines, "\n")
	if trailingNewline {
		body += "\n"
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLastAssistantTextSkipsToolAndUserRecords(t *testing.T) {
	path := write(t, []string{
		assistantLine("an earlier turn"),
		userLine("go on"),
		assistantLine("the turn the gate should judge"),
		toolLine("git status"),
		userLine("thanks"),
	}, true)

	got, err := LastAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "the turn the gate should judge" {
		t.Errorf("got %q, want the last assistant prose turn", got)
	}
}

// The backward reader must agree with a plain forward scan on every fixture,
// including when a record straddles a chunk boundary. tailChunk is shrunk to
// force many boundaries over a small file.
func TestLastAssistantTextMatchesAForwardScanAcrossChunkBoundaries(t *testing.T) {
	var lines []string
	for i := range 200 {
		lines = append(lines, assistantLine(fmt.Sprintf("turn %d %s", i, strings.Repeat("padding ", i%40))))
		lines = append(lines, toolLine(strings.Repeat("x", i*17)))
	}
	lines = append(lines, assistantLine("the final prose turn"))
	// A run of non-prose records after it, long enough to span several chunks.
	for i := range 30 {
		lines = append(lines, toolLine(strings.Repeat("y", 500+i*97)))
	}

	for _, trailing := range []bool{true, false} {
		path := write(t, lines, trailing)
		all, err := AllAssistantText(path)
		if err != nil {
			t.Fatal(err)
		}
		want := all[len(all)-1]

		for _, chunk := range []int{16, 64, 512, 4096, 1 << 20} {
			t.Run(fmt.Sprintf("trailing=%v/chunk=%d", trailing, chunk), func(t *testing.T) {
				old := tailChunk
				tailChunk = chunk
				defer func() { tailChunk = old }()

				got, err := LastAssistantText(path)
				if err != nil {
					t.Fatal(err)
				}
				if got != want {
					t.Errorf("got %q, want %q", got, want)
				}
				if got != "the final prose turn" {
					t.Errorf("got %q, want the fixture's final prose turn", got)
				}
			})
		}
	}
}

func TestLastAssistantTextOnEmptyAndProseFreeFiles(t *testing.T) {
	empty := write(t, nil, false)
	got, err := LastAssistantText(empty)
	if err != nil || got != "" {
		t.Errorf("empty file: got %q, %v — want empty and no error", got, err)
	}

	noProse := write(t, []string{toolLine("ls"), userLine("hi")}, true)
	got, err = LastAssistantText(noProse)
	if err != nil || got != "" {
		t.Errorf("no assistant prose: got %q, %v", got, err)
	}
}

// A truncated write or a half-flushed line is normal in a live transcript. It
// must not blind the gate to the turn above it.
func TestMalformedLinesAreSkippedNotFatal(t *testing.T) {
	path := write(t, []string{
		assistantLine("the good turn"),
		`{"type":"assistant","message":{"role":"assistant","content":[{"type":"te`,
		"",
		"not json at all",
	}, false)

	got, err := LastAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "the good turn" {
		t.Errorf("got %q, want the last well-formed turn", got)
	}
}

func TestAllAssistantTextReturnsEveryTurnOldestFirst(t *testing.T) {
	path := write(t, []string{
		assistantLine("first"),
		toolLine("noise"),
		assistantLine("second"),
		userLine("noise"),
		assistantLine("third"),
	}, true)

	got, err := AllAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"first", "second", "third"}
	if len(got) != len(want) {
		t.Fatalf("got %d turns, want %d: %q", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("turn %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

// Multiple text blocks in one message are joined, which is what the model
// produces when it interleaves prose with tool calls.
func TestMultipleTextBlocksAreJoined(t *testing.T) {
	b, _ := json.Marshal(map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"role": "assistant",
			"content": []map[string]any{
				{"type": "text", "text": "part one"},
				{"type": "tool_use", "name": "Read"},
				{"type": "text", "text": "part two"},
			},
		},
	})
	path := write(t, []string{string(b)}, true)

	got, err := LastAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "part one\npart two" {
		t.Errorf("got %q, want both text blocks joined", got)
	}
}
