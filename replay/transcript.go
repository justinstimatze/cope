package main

// Reading real prompts out of Claude Code transcripts.
//
// This duplicates a little of internal/transcript on purpose. That package is
// under cope's internal/, so a second module cannot import it, and what this
// needs is a different shape anyway: the live gate wants the assistant prose
// since the last prompt, and a replay wants the prompt plus the conversation
// standing above it.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type record struct {
	Type    string `json:"type"`
	Message struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

type block struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// A Case is one prompt to replay, with the conversation that preceded it.
type Case struct {
	Prompt  string
	History []Message // oldest first, alternating user/assistant
	Reply   string    // what was actually written, for reference only
}

type Message struct {
	Role string
	Text string
}

// text pulls the human-readable text out of a record's content, which is
// either a bare string or an array of typed blocks.
//
// A content array carrying a tool_result is a tool carrier rather than
// something a person typed, and returns "" so the caller skips it.
func text(raw json.RawMessage, skipToolResults bool) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return strings.TrimSpace(s)
	}
	var blocks []block
	if json.Unmarshal(raw, &blocks) != nil {
		return ""
	}
	var b strings.Builder
	for _, bl := range blocks {
		if skipToolResults && bl.Type == "tool_result" {
			return ""
		}
		if bl.Type == "text" {
			if b.Len() > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString(strings.TrimSpace(bl.Text))
		}
	}
	return strings.TrimSpace(b.String())
}

// readCases walks a transcript forward and returns every real user prompt with
// the conversation above it, keeping at most historyTurns messages of context.
//
// Skipped: hook-injected prompts (the card and the reminders themselves, which
// would contaminate every arm), slash commands, and prompts whose reply was too
// short to score meaningfully.
func readCases(path string, historyTurns, minReply int) ([]Case, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cases []Case
	var history []Message
	var pending *Case
	var reply strings.Builder

	closeOut := func() {
		if pending == nil {
			return
		}
		pending.Reply = strings.TrimSpace(reply.String())
		history = append(history,
			Message{Role: "user", Text: pending.Prompt},
			Message{Role: "assistant", Text: pending.Reply})
		if len(history) > historyTurns {
			history = history[len(history)-historyTurns:]
		}
		if len(pending.Reply) >= minReply {
			cases = append(cases, *pending)
		}
		pending, reply = nil, strings.Builder{}
	}

	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 256<<10), 16<<20)
	for s.Scan() {
		var r record
		if json.Unmarshal(s.Bytes(), &r) != nil {
			continue
		}
		switch r.Type {
		case "user":
			t := text(r.Message.Content, true)
			if t == "" || !realPrompt(t) {
				continue
			}
			closeOut()
			ctx := make([]Message, len(history))
			copy(ctx, history)
			pending = &Case{Prompt: t, History: ctx}
		case "assistant":
			if pending == nil {
				continue
			}
			if t := text(r.Message.Content, false); t != "" {
				if reply.Len() > 0 {
					reply.WriteString("\n\n")
				}
				reply.WriteString(t)
			}
		}
	}
	closeOut()
	return cases, s.Err()
}

// realPrompt rejects the things a person did not type. Hook output is the one
// that matters: the card and the reminders arrive as user text, and replaying a
// prompt that already carries a reminder would put it in every arm.
func realPrompt(t string) bool {
	switch {
	case strings.HasPrefix(t, "<voice-card>"), strings.Contains(t, "<voice-refresher>"):
		return false
	case strings.HasPrefix(t, "<command-name>"), strings.HasPrefix(t, "<local-command"):
		return false
	case strings.HasPrefix(t, "Caveat:"), strings.HasPrefix(t, "<system-reminder>"):
		return false
	case strings.HasPrefix(t, "[Request interrupted"):
		return false
	case len(t) < 12 || len(t) > 4000:
		return false
	}
	return true
}

func (c Case) String() string {
	p := c.Prompt
	if len(p) > 60 {
		p = p[:60] + "…"
	}
	return fmt.Sprintf("%q (%d ctx msgs, reply %d ch)", p, len(c.History), len(c.Reply))
}
