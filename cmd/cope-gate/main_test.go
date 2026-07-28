package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/justinstimatze/cope/internal/scan"
)

// The clock's contract, inverted from basanite's: first sight of a session
// stays silent (SessionStart just delivered the full card), an aged marker
// fires, and firing restamps so the next fire waits a full interval again.
func TestClaimRefreshSilentFirstThenFiresWhenAged(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "refresher-test-session")

	if claimRefresh(marker, time.Hour) {
		t.Fatal("first sight must stamp and stay silent, not fire")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("first sight must create the marker: %v", err)
	}

	if claimRefresh(marker, time.Hour) {
		t.Fatal("a fresh marker must not fire")
	}

	old := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(marker, old, old); err != nil {
		t.Fatal(err)
	}
	if !claimRefresh(marker, time.Hour) {
		t.Fatal("an aged marker must fire")
	}
	fi, err := os.Stat(marker)
	if err != nil {
		t.Fatal(err)
	}
	if time.Since(fi.ModTime()) > time.Minute {
		t.Error("firing must restamp the marker to now")
	}
}

func TestStampRefresherResetsAnAgedClock(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	stampRefresher("stamped-session-1")
	marker := filepath.Join(dir, "cope", "refresher-stamped-session-1")
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("stamp must create the marker: %v", err)
	}

	old := time.Now().Add(-3 * time.Hour)
	os.Chtimes(marker, old, old)
	stampRefresher("stamped-session-1")
	if claimRefresh(marker, time.Hour) {
		t.Error("a stamp (SessionStart re-inject) must reset the clock so the refresher waits again")
	}
}

// The gate must work from a clone at any path and from `go install` with no
// clone at all. The old default read card/rules.json from a directory guessed
// under $HOME, so anywhere but the author's machine it printed "no check run"
// every turn and checked nothing.
func TestLoadCardFallsBackToTheEmbeddedCard(t *testing.T) {
	c, warnings, err := loadCard("")
	if err != nil {
		t.Fatalf("the built-in card must always load: %v", err)
	}
	for _, w := range warnings {
		t.Errorf("the shipped card warns on load: %v", w)
	}
	if len(c.Rules) == 0 || len(c.Never) == 0 || len(c.Tests) == 0 {
		t.Errorf("built-in card is missing content: %d rules, %d never, %d tests",
			len(c.Rules), len(c.Never), len(c.Tests))
	}
	if c.Render() == "" {
		t.Error("built-in card renders nothing to inject")
	}

	if _, _, err := loadCard(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Error("an explicitly named card that does not exist must be an error, not a silent fallback")
	}
}

// The violations log quotes the reply verbatim, so it belongs in per-machine
// state rather than in the repo checkout the old default wrote to.
func TestDefaultLogPathFollowsXDGState(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)
	if got, want := defaultLogPath(), filepath.Join(dir, "cope", "violations.jsonl"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	t.Setenv("XDG_STATE_HOME", "")
	home := t.TempDir()
	t.Setenv("HOME", home)
	if got, want := defaultLogPath(), filepath.Join(home, ".local", "state", "cope", "violations.jsonl"); got != want {
		t.Errorf("with XDG_STATE_HOME unset: got %q, want %q", got, want)
	}
}

// report creates the state directory on demand and must not leave the log
// world-readable.
func TestReportWritesTheLogPrivately(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "nested", "violations.jsonl")
	report([]scan.Violation{{RuleID: "flip", Action: "warn", Matched: "not x, it's y"}}, "session-1", logPath)

	fi, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("log was not created: %v", err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("log mode is %04o, want 0600 — it holds verbatim excerpts of the reply", perm)
	}
}

func TestSessionIDReRejectsPathAndShellShapes(t *testing.T) {
	for _, bad := range []string{"", "short", "../../../etc/passwd", "a b c d e f g h", "x;rm -rf ~"} {
		if sessionIDRe.MatchString(bad) {
			t.Errorf("accepted %q as a session id", bad)
		}
	}
	if !sessionIDRe.MatchString("c078bae6-05af-4b87-9544-52c78ce1271f") {
		t.Error("rejected a real session id shape")
	}
}
