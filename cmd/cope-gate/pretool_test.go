package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runPreToolOn feeds one PreToolUse payload through the entry and returns what
// it wrote to stdout, which is the whole of what Claude Code reads back.
func runPreToolOn(t *testing.T, payload map[string]any, logPath string) string {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	in, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := in.Write(raw); err != nil {
		t.Fatal(err)
	}
	if _, err := in.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	out, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}

	oldIn, oldOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = in, out
	defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()
	runPreTool(&testCard, 0.35, logPath)

	if _, err := out.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The labelled opening is the fixture throughout: a noun phrase standing where
// a sentence belongs, which no lane drops.
const labelledTicket = "Row loss on cursor reset. The sync drops every row written after the " +
	"upstream cursor rewinds, because the batch commits before the cursor " +
	"advances and the next poll starts from the stale position."

func TestPreToolWarnsWithoutGating(t *testing.T) {
	got := runPreToolOn(t, map[string]any{
		"session_id": "abc",
		"tool_name":  "mcp__linear__save_issue",
		"tool_input": map[string]any{"title": "Sync drops rows", "description": labelledTicket},
	}, "")
	if got == "" {
		t.Fatal("scored nothing on prose that fails labelled_opening")
	}

	var out preToolOutput
	if err := json.Unmarshal([]byte(got), &out); err != nil {
		t.Fatalf("output does not parse: %v\n%s", err, got)
	}
	if out.HookSpecificOutput.HookEventName != "PreToolUse" {
		t.Errorf("wrong event name: %q", out.HookSpecificOutput.HookEventName)
	}
	if !strings.Contains(out.HookSpecificOutput.AdditionalContext, "labelled_opening") {
		t.Errorf("context does not name the rule that fired:\n%s", out.HookSpecificOutput.AdditionalContext)
	}
	if !strings.Contains(out.HookSpecificOutput.AdditionalContext, "description:") {
		t.Errorf("context does not name the field it scored:\n%s", out.HookSpecificOutput.AdditionalContext)
	}

	// The decision this file exists to hold: no permissionDecision anywhere in
	// the payload. A gate here would be stricter than cope is at Stop, where
	// --block is off by default and the surface is a reply nobody else reads.
	if strings.Contains(got, "permissionDecision") {
		t.Errorf("PreToolUse emitted a permission decision — this entry warns only:\n%s", got)
	}
}

func TestPreToolIgnoresToolsItDoesNotCover(t *testing.T) {
	// A read tool on the same server. Scoring a search result would be scoring
	// somebody else's prose.
	if got := runPreToolOn(t, map[string]any{
		"tool_name":  "mcp__linear__list_issues",
		"tool_input": map[string]any{"description": labelledTicket},
	}, ""); got != "" {
		t.Errorf("scored a tool outside the covered set:\n%s", got)
	}
}

func TestPreToolSaysNothingOnCleanProse(t *testing.T) {
	if got := runPreToolOn(t, map[string]any{
		"tool_name":  "mcp__linear__save_comment",
		"tool_input": map[string]any{"body": "Reverted the cursor change. The staging replica is back to a clean sync."},
	}, ""); got != "" {
		t.Errorf("added context with nothing to say:\n%s", got)
	}
}

// Every failure path ends in silence and exit 0, because the write it sits in
// front of matters more than the score.
func TestPreToolFailsOpen(t *testing.T) {
	for name, payload := range map[string]map[string]any{
		"no tool_input":     {"tool_name": "mcp__linear__save_issue"},
		"no prose field":    {"tool_name": "mcp__linear__save_issue", "tool_input": map[string]any{"title": "x"}},
		"field is not text": {"tool_name": "mcp__linear__save_issue", "tool_input": map[string]any{"description": 42}},
		"field is blank":    {"tool_name": "mcp__linear__save_issue", "tool_input": map[string]any{"description": "   "}},
		"no tool name":      {"tool_input": map[string]any{"description": labelledTicket}},
	} {
		if got := runPreToolOn(t, payload, ""); got != "" {
			t.Errorf("%s: expected silence, got:\n%s", name, got)
		}
	}
}

// The ending rules are dropped here and nowhere else. A ticket that names an
// open problem and stops is a ticket doing its job.
func TestPreToolScoresInTheExternalLane(t *testing.T) {
	dangling := "Rebased onto main and reran the suite, four packages green.\n\n" +
		"The stale fixture in scan_test.go is still unresolved. I left it alone " +
		"because it passes and nothing in this run touched that path."
	got := runPreToolOn(t, map[string]any{
		"tool_name":  "mcp__linear__save_issue",
		"tool_input": map[string]any{"description": dangling},
	}, "")
	if strings.Contains(got, "dangling_end") {
		t.Errorf("flagged an ending in prose with no reader to answer it:\n%s", got)
	}
}

func TestPreToolLogsWhatItFound(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "violations.jsonl")
	runPreToolOn(t, map[string]any{
		"tool_name":  "mcp__linear__save_issue",
		"tool_input": map[string]any{"description": labelledTicket},
	}, logPath)

	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("nothing logged: %v", err)
	}
	// The session field names the destination, not a transcript id: these
	// records have to be separable from reply scores after the fact.
	if !strings.Contains(string(b), "mcp__linear__save_issue:description") {
		t.Errorf("log record does not name the tool and field:\n%s", b)
	}
}

// The matcher and the covered set are built from one map, so a tool cannot be
// scored without being matched or matched without being scored.
func TestPreToolMatcherCoversExactlyWhatItScores(t *testing.T) {
	got := strings.Split(preToolMatcher(), "|")
	if len(got) != len(preToolTools) {
		t.Fatalf("matcher names %d tools, the entry scores %d", len(got), len(preToolTools))
	}
	for _, name := range got {
		if !preToolTools[name] {
			t.Errorf("matcher names %q, which runPreTool would ignore", name)
		}
	}
}
