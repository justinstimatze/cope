package main

import (
	"encoding/json"
	"github.com/justinstimatze/cope/internal/scan"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readSettings(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("result does not parse: %v", err)
	}
	return m
}

func TestSetupKeepsEverythingItDidNotComeFor(t *testing.T) {
	// The file carries API keys, permission rules and other tools' hooks. The
	// whole design is that cope adds two keys and leaves the rest alone.
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".claude", "settings.json")
	orig := `{"env":{"SECRET":"keep-me"},"permissions":{"allow":["Bash(ls:*)"]},` +
		`"hooks":{"Stop":[{"matcher":"","hooks":[{"type":"command","command":"/usr/bin/other"}]}]}}`
	if err := os.WriteFile(path, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", dir)

	if rc := runSetup(&testCard, false); rc != 0 {
		t.Fatalf("rc = %d", rc)
	}
	got := readSettings(t, path)

	if env, _ := got["env"].(map[string]any); env["SECRET"] != "keep-me" {
		t.Error("the secret did not survive")
	}
	if got["permissions"] == nil {
		t.Error("permissions did not survive")
	}
	b, _ := json.Marshal(got["hooks"])
	if !strings.Contains(string(b), "/usr/bin/other") {
		t.Error("another tool's Stop hook was dropped")
	}
	if !strings.Contains(string(b), "cope-gate") {
		t.Error("cope's hook was not added")
	}
}

func TestSetupIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	if rc := runSetup(&testCard, false); rc != 0 {
		t.Fatalf("first run rc = %d", rc)
	}
	path := filepath.Join(dir, ".claude", "settings.json")
	first := readSettings(t, path)
	if rc := runSetup(&testCard, false); rc != 0 {
		t.Fatalf("second run rc = %d", rc)
	}
	second := readSettings(t, path)

	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	if string(a) != string(b) {
		t.Error("a second --setup changed the file; hooks would be doubled")
	}
}

func TestSetupRefusesAnUnparseableSettingsFile(t *testing.T) {
	// Claude Code ignores an unparseable settings file wholesale, so a
	// half-repaired one is worse than an untouched one.
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".claude", "settings.json")
	broken := []byte(`{"hooks": { BROKEN`)
	if err := os.WriteFile(path, broken, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", dir)

	if rc := runSetup(&testCard, false); rc == 0 {
		t.Error("a broken settings file was reported as success")
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(broken) {
		t.Error("the broken file was modified")
	}
}

func TestSetupWritesAnAbsoluteHookCommand(t *testing.T) {
	// A bare command assumes go install's target is on PATH in whatever
	// environment the hook runs in. It often is not, and the failure is a hook
	// that silently does nothing.
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	if rc := runSetup(&testCard, false); rc != 0 {
		t.Fatalf("rc = %d", rc)
	}
	got := readSettings(t, filepath.Join(dir, ".claude", "settings.json"))
	b, _ := json.Marshal(got["hooks"])
	var cmds []string
	for _, part := range strings.Split(string(b), `"command":"`)[1:] {
		cmds = append(cmds, part[:strings.Index(part, `"`)])
	}
	if len(cmds) == 0 {
		t.Fatal("no commands written")
	}
	for _, c := range cmds {
		if !strings.HasPrefix(c, "/") {
			t.Errorf("command %q is not absolute", c)
		}
	}
}

func TestSetupDryRunWritesNothing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	if rc := runSetup(&testCard, true); rc != 0 {
		t.Fatalf("rc = %d", rc)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude", "settings.json")); !os.IsNotExist(err) {
		t.Error("--dry-run wrote the settings file")
	}
}

// testCard is the minimum a setup run needs: an id for the style filename.
var testCard = scan.Card{CardID: "testcard", Theme: "a test theme"}
