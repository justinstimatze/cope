package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justinstimatze/cope/internal/scan"
)

func TestStyleBodyCarriesTheFrontmatterClaudeCodeNeeds(t *testing.T) {
	c := &scan.Card{CardID: "lecturer", Theme: "Arrive somewhere, and let them feel the arriving"}
	got := styleBody(c)

	if !strings.HasPrefix(got, "---\n") {
		t.Fatalf("no frontmatter opener:\n%.80s", got)
	}
	for _, want := range []string{
		"name: lecturer",
		"description: Arrive somewhere, and let them feel the arriving",
		// The default is false, which drops Claude Code's software engineering
		// instructions. A voice card has no opinion on those and silently
		// swapping them with the voice is the confound to avoid.
		"keep-coding-instructions: true",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("frontmatter missing %q", want)
		}
	}
	if !strings.Contains(got, c.Render()) {
		t.Error("body does not carry the rendered card")
	}
}

func TestStyleDescriptionSurvivesAThemeWithAColon(t *testing.T) {
	// An unquoted colon inside a YAML scalar makes the frontmatter unparseable,
	// and a style that fails to load is a voice that silently does not apply.
	c := &scan.Card{CardID: "x", Theme: "Result first: then the rest"}
	line := ""
	for _, l := range strings.Split(styleBody(c), "\n") {
		if strings.HasPrefix(l, "description:") {
			line = l
		}
	}
	if line == "" {
		t.Fatal("no description line")
	}
	if !strings.Contains(line, `"Result first: then the rest"`) {
		t.Errorf("colon not quoted: %s", line)
	}
}

func TestStyleWritesUnderTheCardID(t *testing.T) {
	dir := t.TempDir()
	c := &scan.Card{CardID: "caveman", Theme: "Few words."}
	if rc := runOutputStyle(c, dir); rc != 0 {
		t.Fatalf("rc = %d", rc)
	}
	b, err := os.ReadFile(filepath.Join(dir, "caveman.md"))
	if err != nil {
		t.Fatalf("expected caveman.md: %v", err)
	}
	if !strings.Contains(string(b), "name: caveman") {
		t.Error("wrong card written")
	}
}

func TestStyleNameFallsBackWhenTheCardHasNoID(t *testing.T) {
	if got := styleName(&scan.Card{}); got != "cope" {
		t.Errorf("styleName = %q, want cope", got)
	}
}

func TestInjectStandsDownOnlyForACopeEmittedStyle(t *testing.T) {
	// The guard exists because a session ran with the maximal card as its
	// output style while the hook injected the shipped card twice as messages,
	// and followed the messages. It must not fire on somebody else's style.
	styles := t.TempDir()
	settings := t.TempDir()

	write := func(dir, name, body string) {
		if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".claude", name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(settings, "settings.local.json", `{"outputStyle":"lecturer"}`)
	if err := os.WriteFile(filepath.Join(styles, "lecturer.md"), []byte("---\nname: lecturer\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(settings); err != nil {
		t.Fatal(err)
	}

	if got := activeStyle(); got != "lecturer" {
		t.Fatalf("activeStyle = %q, want lecturer", got)
	}
}

func TestActiveStyleIgnoresSettingsWithoutTheKey(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	// A settings file that exists and names no style must not stop the search,
	// and must never be read for anything but this one key.
	if err := os.WriteFile(filepath.Join(dir, ".claude", "settings.local.json"),
		[]byte(`{"env":{"SECRET":"do-not-read"},"permissions":{"allow":[]}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Falls through to the user-level file, which on a test box names nothing.
	if got := activeStyle(); got == "do-not-read" {
		t.Error("read a key it had no business reading")
	}
}
