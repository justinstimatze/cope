package main

// --setup does the install a reader would otherwise do by hand: emit the output
// style, wire the three hooks, and say the one thing left to click.
//
// It edits ~/.claude/settings.json, which is the most invasive thing this
// binary does. That file holds API keys, other tools' hooks and permission
// rules, so the whole design here is about touching as little of it as
// possible: parse it, add what is missing, leave every other key byte-identical
// in meaning, and back it up before writing. A settings file that fails to
// parse is left alone entirely — a half-repaired one is worse than an untouched
// one, and Claude Code ignores an unparseable settings file wholesale, which is
// how "my hooks stopped working" usually happens.
//
// The hook commands are written as absolute paths. The README has to warn that
// a bare command assumes go install's target is on PATH in whatever environment
// the hook runs in, and that a hook which silently does nothing usually wants
// the absolute path instead. --setup knows where it is running from, so it can
// remove that failure rather than document it.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/justinstimatze/cope/internal/scan"
)

// preToolMatcher is the tool-name pattern the PreToolUse entry registers.
//
// Built from the same map runPreTool checks, so the two cannot drift: a tool
// added to one and forgotten in the other would either be scored without being
// matched, or matched and then silently dropped. Sorted because a map has no
// order and a settings file that changes on every --setup looks like a hook
// rewriting itself.
func preToolMatcher() string {
	names := make([]string, 0, len(preToolTools))
	for name := range preToolTools {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, "|")
}

// setupHooks are the events worth wiring now that the card arrives as an output
// style. SessionStart is deliberately absent: that was the card's old delivery
// and it is the weakest slot available.
//
// PreToolUse is the one entry that is not about the conversation. It covers the
// prose that leaves it — see cmd/cope-gate/pretool.go — and it is not async,
// because a PreToolUse hook that returns after its tool call has run has
// nothing left to add context to.
func setupHooks(self string) map[string]any {
	return map[string]any{
		"Stop": []any{map[string]any{
			"matcher": "",
			"hooks": []any{map[string]any{
				"type": "command", "command": self, "timeout": 10, "async": true,
			}},
		}},
		"UserPromptSubmit": []any{map[string]any{
			"hooks": []any{map[string]any{
				"type": "command", "command": self + " --refresher",
			}},
		}},
		"PreToolUse": []any{map[string]any{
			"matcher": preToolMatcher(),
			"hooks": []any{map[string]any{
				"type": "command", "command": self + " --pretool", "timeout": 10,
			}},
		}},
	}
}

// hasCopeHook reports whether an event's existing entries already run this
// binary, so a second --setup adds nothing. Matching on the basename rather
// than the full command means a hook wired by hand with a different path, or
// with extra flags, still counts as present — adding a duplicate would double
// every score.
func hasCopeHook(existing any) bool {
	b, err := json.Marshal(existing)
	if err != nil {
		return false
	}
	return strings.Contains(string(b), "cope-gate")
}

func runSetup(c *scan.Card, dryRun bool) int {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: cannot resolve my own path: %v\n", err)
		return 1
	}
	if resolved, err := filepath.EvalSymlinks(self); err == nil {
		self = resolved
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: cannot resolve a home directory: %v\n", err)
		return 1
	}
	settingsPath := filepath.Join(home, ".claude", "settings.json")

	// 1. the style, which is how the card actually arrives
	if !dryRun {
		if rc := runOutputStyle(c, "", true); rc != 0 {
			return rc
		}
	} else {
		fmt.Printf("would write %s\n", filepath.Join(stylesDir(), styleName(c)+".md"))
	}

	// 2. the hooks
	settings := map[string]any{}
	raw, readErr := os.ReadFile(settingsPath)
	switch {
	case readErr == nil:
		if err := json.Unmarshal(raw, &settings); err != nil {
			fmt.Fprintf(os.Stderr,
				"cope-gate: %s does not parse (%v)\n"+
					"  leaving it alone — Claude Code ignores an unparseable settings file entirely,\n"+
					"  so repairing it halfway would be worse than not touching it\n", settingsPath, err)
			return 1
		}
	case !os.IsNotExist(readErr):
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", readErr)
		return 1
	}

	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
	}
	var added []string
	for event, block := range setupHooks(self) {
		if hasCopeHook(hooks[event]) {
			continue
		}
		// Append rather than replace: another tool's hook on the same event is
		// somebody else's and must survive.
		if cur, ok := hooks[event].([]any); ok {
			hooks[event] = append(cur, block.([]any)...)
		} else {
			hooks[event] = block
		}
		added = append(added, event)
	}
	// Sorted because setupHooks returns a map, and an unsorted line would name
	// the same three events in a different order on every run.
	sort.Strings(added)
	settings["hooks"] = hooks

	if len(added) == 0 {
		fmt.Println("hooks already wired, nothing to add")
	} else if dryRun {
		fmt.Printf("would add hooks to %s: %s\n", settingsPath, strings.Join(added, ", "))
	} else {
		if err := writeSettings(settingsPath, settings, raw); err != nil {
			fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
			return 1
		}
		fmt.Printf("wired %s in %s\n", strings.Join(added, " and "), settingsPath)
	}

	fmt.Println()
	fmt.Println("The hooks are live now — Claude Code re-reads them per prompt, so sessions")
	fmt.Println("already open pick them up without restarting (measured against 2.1.220).")
	fmt.Printf("One step left: /config -> Output style -> %s\n", styleName(c))
	fmt.Println("The style is the part read once at session start, so it applies at the")
	fmt.Println("next session or after /clear.")
	return 0
}

// writeSettings backs the file up, then replaces it atomically.
//
// The backup is the point. Everything else here is careful, and careful is not
// the same as reversible: this is somebody's live configuration and the only
// acceptable failure mode is one they can undo without knowing what cope did.
func writeSettings(path string, settings map[string]any, original []byte) error {
	if len(original) > 0 {
		backup := fmt.Sprintf("%s.before-cope-%s", path, time.Now().Format("20060102-150405"))
		if err := os.WriteFile(backup, original, 0o600); err != nil {
			return fmt.Errorf("backing up %s: %w", path, err)
		}
		fmt.Printf("backed up to %s\n", backup)
	}
	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".settings-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(out); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// 0600: this file carries API keys.
	if err := os.Chmod(tmp.Name(), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
