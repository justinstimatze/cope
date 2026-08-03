package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/justinstimatze/cope/internal/scan"
	"github.com/justinstimatze/cope/internal/state"
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
	c, warnings, err := loadCard("", "")
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

	if _, _, err := loadCard(filepath.Join(t.TempDir(), "absent.json"), ""); err == nil {
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
	report("cope:", []scan.Violation{{RuleID: "flip", Action: "warn", Matched: "not x, it's y"}}, "session-1", logPath)

	fi, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("log was not created: %v", err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("log mode is %04o, want 0600 — it holds verbatim excerpts of the reply", perm)
	}
}

// The catch-up rescores a whole turn and must report only what the Stop hook
// could not have seen, counting repeats rather than collapsing them.
func TestUnreportedReturnsOnlyWhatIsNew(t *testing.T) {
	bold := func(m string) scan.Violation {
		return scan.Violation{RuleID: "bold_label", Action: "warn", Matched: m}
	}
	prev := []scan.Violation{bold("**One.**"), bold("**One.**")}
	full := []scan.Violation{
		bold("**One.**"),
		bold("**One.**"),
		bold("**One.**"),
		{RuleID: "dangling_end", Action: "warn", Matched: "still unverified"},
	}

	got := unreported(prev, full)
	if len(got) != 2 {
		t.Fatalf("got %d new violations, want 2: %+v", len(got), got)
	}
	if got[0] != bold("**One.**") {
		t.Errorf("the third repeat should survive as new: %+v", got[0])
	}
	if got[1].RuleID != "dangling_end" {
		t.Errorf("the ending violation should survive: %+v", got[1])
	}
	if n := len(unreported(full, full)); n != 0 {
		t.Errorf("an unchanged turn reported %d violations again", n)
	}
	if n := len(unreported(nil, full)); n != len(full) {
		t.Errorf("with nothing reported before, all %d should be new, got %d", len(full), n)
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

// A typo in the arm list that quietly halved the experiment would look exactly
// like data, so an unknown arm is an error and the experiment goes off.
func TestParseArms(t *testing.T) {
	if arms, err := parseArms("", false); err != nil || arms != nil {
		t.Errorf("no flags should leave the experiment off, got %v %v", arms, err)
	}
	arms, err := parseArms("", true)
	if err != nil || len(arms) != 2 {
		t.Errorf("-ab alone should give the default pair, got %v %v", arms, err)
	}
	arms, err = parseArms("inject, hold ,positive", false)
	if err != nil || len(arms) != 3 {
		t.Fatalf("-ab-arms should imply -ab and keep all three, got %v %v", arms, err)
	}
	if arms[2] != state.ArmPositive {
		t.Errorf("arm order should follow the flag, got %v", arms)
	}
	if _, err := parseArms("inject,hodl", false); err == nil {
		t.Error("a misspelled arm should be an error, not a silently dropped arm")
	}
	if _, err := parseArms("inject", false); err == nil {
		t.Error("one arm is not an experiment")
	}
}

// A card is only switchable if a name resolves to a file, and only trustworthy
// if a name that resolves to nothing says so. Silently writing in the shipped
// voice while the config names another one is the failure that reads as the
// card not working.
func TestCardNameResolvesAndAMissTellsYou(t *testing.T) {
	cfg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfg)
	dir := filepath.Join(cfg, "cope", "cards")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join("..", "..", "card", "demo", "lecturer.effigy"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "lecturer.effigy"), src, 0o644); err != nil {
		t.Fatal(err)
	}

	c, _, err := loadCard("", "lecturer")
	if err != nil {
		t.Fatalf("installed card did not load: %v", err)
	}
	if c.CardID != "lecturer" {
		t.Errorf("loaded %q, want the lecturer card", c.CardID)
	}

	if _, _, err := loadCard("", "nosuchvoice"); err == nil {
		t.Error("a named card that is not installed must be an error, not the built-in card")
	}
}

// The env var is how a session picks a voice without editing a hook, and the
// flag has to beat it or a shell export would pin every run on the machine.
func TestCopeCardEnvIsReadAndTheFlagWins(t *testing.T) {
	cfg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfg)
	dir := filepath.Join(cfg, "cope", "cards")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"lecturer", "caveman"} {
		src, err := os.ReadFile(filepath.Join("..", "..", "card", "demo", name+".effigy"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name+".effigy"), src, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("COPE_CARD", "caveman")
	c, _, err := loadCard("", "")
	if err != nil {
		t.Fatalf("$COPE_CARD did not load: %v", err)
	}
	if c.CardID != "caveman" {
		t.Errorf("$COPE_CARD gave %q, want caveman", c.CardID)
	}

	c, _, err = loadCard("", "lecturer")
	if err != nil {
		t.Fatal(err)
	}
	if c.CardID != "lecturer" {
		t.Errorf("--card lost to $COPE_CARD: got %q", c.CardID)
	}
}

// Reading .effigy directly is the whole point of the native parser: a card is
// usable as written. It has to arrive at the same card the build-time path
// produces, or the two would drift apart in a way only a run would show.
func TestEffigyLoadsDirectlyAndMatchesTheGeneratedJSON(t *testing.T) {
	fromEffigy, _, err := loadCard(filepath.Join("..", "..", "card", "claude_voice.effigy"), "")
	if err != nil {
		t.Fatalf("the shipped card did not load from notation: %v", err)
	}
	fromJSON, _, err := loadCard("", "")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := fromEffigy.Render(), fromJSON.Render(); got != want {
		t.Error("the card parsed from .effigy injects different text than the embedded JSON")
	}
	if len(fromEffigy.Rules) != len(fromJSON.Rules) {
		t.Errorf("%d rules from notation, %d from JSON", len(fromEffigy.Rules), len(fromJSON.Rules))
	}
}
