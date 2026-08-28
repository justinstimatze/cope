package main

// The PreToolUse entry: the one place cope sees prose that is not a reply.
//
// Every other hook here reads the conversation. Stop scores the turn just
// written, the refresher restates what has been firing, MessageDisplay rewrites
// what the reader sees — all of them are about the text between one person and
// this session. A ticket description posted through an MCP tool never touches
// that surface. It is written inside a turn, goes out mid-turn, and is read
// days later by somebody who was not here. Until this file existed cope could
// not see it at all: internal/scan reaches prose through the Stop transcript,
// --check on a file, and --backfill on a JSONL, and a tool_input is none of the
// three.
//
// Three decisions shape what follows, and each is a decision not to do the
// obvious thing.
//
// It does not gate. The output carries additionalContext and no
// permissionDecision, so the write proceeds and the model learns what the prose
// scored. cope's own posture already settles this: POSTPROC rules carry
// warn or reject, --block is opt-in and off by default, and it is off at Stop —
// where the surface is a chat reply nobody else reads. Gating an external write
// harder than cope gates its own output would invert the tool. If a rule ever
// earns a gate it earns it by being action: reject under an explicit --block,
// which is vocabulary that exists.
//
// The honest cost of that: additionalContext without a deny arrives after the
// call. The model learns the description scored badly once it is already
// posted. For Linear that is survivable — save_issue and save_comment are
// editable and the model holds the id — but this catches the prose one beat
// late rather than stopping it, and nobody should expect otherwise.
//
// It writes no state. The obvious design is a per-session seen-set so a tic is
// named once and not on every ticket, and that design is a known bug: a single
// set shared across destinations lets a local write spend the one warning
// before an external write ever asks for it. cope's version would be worse than
// a shared set. state.Session.Record replaces the newest entry when the key
// matches, and the key is the transcript id of the prompt that opened the turn
// — so a mid-turn scan carrying that key would overwrite the turn's real reply
// score, and passing an empty key appends a phantom turn into a 20-turn window
// instead. Downstream, Top feeds the refresher, so ticket prose would decide
// which card items get restated to somebody writing chat replies, and
// PriorPhrases would call a construction repeated across two audiences. None of
// that is worth a deduplicated warning, so this path reads no session file and
// writes none. Per-write feedback repeats. That is the correct amount.
//
// It scores in the external lane, not the interactive one. See
// scan.LaneExternal: four rules measure an ending for a reader who can answer
// "continue", and a ticket has no such reader.

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/justinstimatze/cope/internal/scan"
)

// preToolTools are the tool names this entry scores.
//
// The list is Linear's write surface and matches what the settings matcher
// registers. It is a list rather than a prefix test because mcp__linear__ also
// carries the read tools, and scoring a search result would be scoring
// somebody else's prose.
var preToolTools = map[string]bool{
	"mcp__linear__save_issue":         true,
	"mcp__linear__save_comment":       true,
	"mcp__linear__save_project":       true,
	"mcp__linear__save_document":      true,
	"mcp__linear__save_status_update": true,
}

// proseFields are the tool_input keys worth scoring, in report order.
//
// Titles are deliberately absent. A title is one line with no paragraph
// structure, so every shape rule abstains on it and the voicing rules would be
// reading a label as prose. What is left is the body of the thing, under
// whichever of the three names the tool gave it.
var proseFields = []string{"description", "body", "content"}

type preToolInput struct {
	SessionID string         `json:"session_id"`
	ToolName  string         `json:"tool_name"`
	ToolInput map[string]any `json:"tool_input"`
}

type preToolOutput struct {
	HookSpecificOutput struct {
		HookEventName     string `json:"hookEventName"`
		AdditionalContext string `json:"additionalContext"`
	} `json:"hookSpecificOutput"`
}

// runPreTool is the PreToolUse hook entry.
//
// It fails open in every direction. An unreadable payload, an unknown tool, a
// field that is not a string, a clean score — all of them end in nothing on
// stdout and exit 0, which lets the write through untouched. A style scorer
// must never be the reason a ticket does not get filed.
func runPreTool(c *scan.Card, minCV float64, logPath string) {
	var in preToolInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return
	}
	if !preToolTools[in.ToolName] {
		return
	}

	var b strings.Builder
	for _, field := range proseFields {
		text, ok := in.ToolInput[field].(string)
		if !ok || strings.TrimSpace(text) == "" {
			continue
		}
		v := c.CheckLane(text, minCV, scan.LaneExternal)
		if len(v) == 0 {
			continue
		}
		writePreToolReport(&b, field, v)
		logViolations(v, in.ToolName+":"+field, logPath)
	}
	if b.Len() == 0 {
		return
	}

	var out preToolOutput
	out.HookSpecificOutput.HookEventName = "PreToolUse"
	out.HookSpecificOutput.AdditionalContext = preToolPreamble(in.ToolName) + b.String()
	_ = json.NewEncoder(os.Stdout).Encode(out)
}

// preToolPreamble says what the reader is looking at and what it can do about
// it. Without the second half the context is a complaint: the write has already
// gone out by the time this is read, so the actionable move is an edit, and the
// note has to say so or it reads as a rule that arrived too late to obey.
func preToolPreamble(tool string) string {
	return fmt.Sprintf("cope scored the prose this %s call is about to post, in the external lane "+
		"— the four ending rules that assume a reader who can answer are off. This is a warning, "+
		"not a block: the call proceeds. If a hit is right, edit the item rather than reposting it.\n\n",
		tool)
}

// writePreToolReport renders one field's violations, tallied by rule the way
// the Stop hook's report does, so the same hit reads the same in both places.
func writePreToolReport(b *strings.Builder, field string, v []scan.Violation) {
	byRule := map[string]int{}
	for _, x := range v {
		byRule[x.RuleID]++
	}
	ids := make([]string, 0, len(byRule))
	for id := range byRule {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	fmt.Fprintf(b, "%s: %d violation(s) —", field, len(v))
	for _, id := range ids {
		fmt.Fprintf(b, " %s×%d", id, byRule[id])
	}
	b.WriteString("\n")
	for _, x := range v {
		fmt.Fprintf(b, "  [%s] %s\n      %q\n", x.RuleID, x.Why, x.Matched)
	}
}
