package main

// Emitting the card as a Claude Code output style.
//
// --inject was the original delivery: SessionStart prints the card and Claude
// Code puts it in as context. That works and it is the weakest slot available.
// The card lands once, at turn zero, as one message among the conversation, and
// then gets buried as the conversation grows. Measured on 2026-08-03: a card
// instructing a bolded label on every paragraph and an emoji on every heading
// produced, through the hook, prose indistinguishable from no card at all.
//
// An output style is the same instructions in the system prompt instead, and
// the docs name the two properties that matter — the style's text is appended
// to the END of the system prompt, and the harness re-reminds the model of it
// during the conversation. Those are the two jobs --inject and --refresher were
// doing by hand and losing at.
//
// Three arms, 27 samples, 22 of them usable — small, and the effect is large
// enough to read anyway. Counting the moves the card asks for: card alone in
// the system prompt, card beside a 27k CLAUDE.md, and that CLAUDE.md with no
// card. The middle arm kept 63-86% of the first arm's effect. The third
// produced none of it, with no emoji heading and no horizontal rule in any
// sample. So the card is not outweighed by what sits beside it in the system
// prompt; it was losing on position, which is the one thing this file changes.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/justinstimatze/cope/internal/scan"
)

// stylesDir resolves ~/.claude/output-styles, where Claude Code reads
// user-level styles. Unlike cardsDir this is not XDG: the path belongs to
// Claude Code and is stated in its documentation, so following XDG here would
// write somewhere nothing reads.
func stylesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "output-styles")
}

// styleBody is the file Claude Code loads: frontmatter, then the same rendered
// card --inject prints.
//
// keep-coding-instructions is true and not configurable. It defaults to false,
// which drops Claude Code's built-in software engineering instructions
// entirely — how it scopes a change, writes a comment, verifies work. A voice
// card says how a reply should read and has no opinion on any of that, so
// leaving the flag at its default would make swapping a voice quietly also
// swap the engineering behaviour, which is the confound this repo exists to
// keep out of its own measurements.
func styleBody(c *scan.Card) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "name: %s\n", styleName(c))
	if d := frontmatterSafe(c.Theme); d != "" {
		fmt.Fprintf(&b, "description: %s\n", d)
	}
	b.WriteString("keep-coding-instructions: true\n")
	b.WriteString("---\n\n")
	b.WriteString(c.Render())
	if !strings.HasSuffix(b.String(), "\n") {
		b.WriteString("\n")
	}
	return b.String()
}

func styleName(c *scan.Card) string {
	if c.CardID != "" {
		return c.CardID
	}
	return "cope"
}

// frontmatterSafe keeps a theme to one plain line. A card's theme is prose
// somebody wrote and can carry a colon or a newline, either of which turns the
// frontmatter into something Claude Code will not parse — and a style that
// fails to load is a voice that silently does not apply, which is the failure
// this whole file is fixing.
func frontmatterSafe(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if s == "" {
		return ""
	}
	if strings.ContainsAny(s, `:#"'`) {
		return strconv.Quote(s)
	}
	return s
}

// runOutputStyle writes the loaded card to the styles directory and says what
// to do next. It prints the path rather than editing any settings file: which
// style is active is the reader's choice and lives in their settings, and a
// tool that installs a voice should not also be the thing that turns it on.
func runOutputStyle(c *scan.Card, dir string) int {
	if dir == "" {
		dir = stylesDir()
	}
	if dir == "" {
		fmt.Fprintln(os.Stderr, "cope-gate: cannot resolve a home directory for ~/.claude/output-styles")
		return 1
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
		return 1
	}
	path := filepath.Join(dir, styleName(c)+".md")
	if err := os.WriteFile(path, []byte(styleBody(c)), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
		return 1
	}
	fmt.Printf("wrote %s\n", path)
	fmt.Printf("turn it on with /config -> Output style -> %s, or set\n", styleName(c))
	fmt.Printf("  \"outputStyle\": %q\n", styleName(c))
	fmt.Println("in .claude/settings.local.json. It applies at the next session or after /clear.")
	return 0
}

// activeStyle reports the output style selected for the working directory, by
// reading the outputStyle key out of the settings files Claude Code reads.
//
// Only that one key is read. These files carry API keys, tokens and hook
// command lines, and a tool with no business in any of that should not load
// them whole.
//
// Precedence is the narrowest file that names a style, matching Claude Code's
// own local-over-user order well enough for the only decision this feeds: it
// is used to keep --inject quiet, and being wrong in the cautious direction
// means one card injection nobody needed rather than two cards arguing.
func activeStyle() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	wd, _ := os.Getwd()
	for _, p := range []string{
		filepath.Join(wd, ".claude", "settings.local.json"),
		filepath.Join(wd, ".claude", "settings.json"),
		filepath.Join(home, ".claude", "settings.json"),
	} {
		var s struct {
			OutputStyle string `json:"outputStyle"`
		}
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if json.Unmarshal(b, &s) == nil && s.OutputStyle != "" {
			return s.OutputStyle
		}
	}
	return ""
}

// styleSupersedesInject reports whether an output style is already carrying a
// cope card, so --inject would be adding a second and probably different one.
//
// This is not hypothetical. On 2026-08-03 a session ran with the maximal card
// selected as its output style while the SessionStart hook went on injecting
// the shipped card, twice, as ordinary messages. The reply came out in the
// shipped card's register — correctly, since that was the instruction it had
// been given most recently and most often. The style was blamed for an hour.
//
// A style only counts if cope emitted it, which is what the file check is for.
// Somebody else's output style is not cope's to defer to.
func styleSupersedesInject() (string, bool) {
	name := activeStyle()
	if name == "" {
		return "", false
	}
	dir := stylesDir()
	if dir == "" {
		return "", false
	}
	if _, err := os.Stat(filepath.Join(dir, name+".md")); err != nil {
		return "", false
	}
	return name, true
}
