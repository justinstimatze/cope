package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecordSurvivesAProcessBoundary(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()

	s := Load(dir, "sess-1")
	if s.Turns != 0 {
		t.Fatalf("a session with no file should start empty, got %d turns", s.Turns)
	}
	s.Record("", []string{"bold_label", "bold_label", "flip"}, 900, now)
	if err := s.Save(dir); err != nil {
		t.Fatal(err)
	}

	// A separate hook invocation is a separate process; the only channel is disk.
	again := Load(dir, "sess-1")
	if again.Turns != 1 || again.Chars != 900 {
		t.Errorf("reloaded %d turns / %d chars, want 1 / 900", again.Turns, again.Chars)
	}
	top := again.Top(5)
	if len(top) != 2 || top[0].Rule != "bold_label" || top[0].N != 2 {
		t.Errorf("top rules did not survive the round trip: %+v", top)
	}

	fi, err := os.Stat(filepath.Join(dir, "session-sess-1.json"))
	if err != nil {
		t.Fatal(err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("state file mode %04o, want 0600", perm)
	}
}

// A tic from turn 3 of an eighty-turn session should not outvote what is
// happening now, so the tally reads a bounded window.
func TestTheWindowForgets(t *testing.T) {
	s := &Session{ID: "x"}
	now := time.Now()
	s.Record("", []string{"ancient"}, 10, now)
	for i := range recentTurns {
		s.Record("", []string{"current"}, 10, now.Add(time.Duration(i)*time.Minute))
	}

	if s.Turns != recentTurns+1 {
		t.Errorf("lifetime turn count should keep counting: got %d", s.Turns)
	}
	if len(s.Recent) != recentTurns {
		t.Errorf("window is %d turns, want %d", len(s.Recent), recentTurns)
	}
	for _, c := range s.Top(0) {
		if c.Rule == "ancient" {
			t.Error("a rule outside the window still counts toward the injection")
		}
	}
}

func TestTopIsStableUnderTies(t *testing.T) {
	s := &Session{ID: "x"}
	now := time.Now()
	s.Record("", []string{"zebra", "alpha"}, 10, now)

	first := s.Top(0)
	for range 20 {
		got := s.Top(0)
		for i := range got {
			if got[i] != first[i] {
				t.Fatalf("tie order shuffles between calls: %+v then %+v", first, got)
			}
		}
	}
	if first[0].Rule != "alpha" {
		t.Errorf("ties should break on name for a stable injection, got %+v", first)
	}
}

func TestFactsNameThePositionAndTheLiveRules(t *testing.T) {
	s := &Session{ID: "x"}
	now := time.Now()
	s.Record("", []string{"bold_label", "bold_label"}, 10, now)
	s.Record("", []string{"bold_label", "flip"}, 10, now)
	s.Record("", []string{"clause_symmetry"}, 10, now)

	// bold_label fired 3x; clause_symmetry and flip tied at 1 and the tie
	// breaks on name, so clause_symmetry takes the second slot.
	f := s.Facts("prompt", 2)
	for _, want := range []string{"at_prompt", "rule_bold_label", "rule_clause_symmetry"} {
		if !f[want] {
			t.Errorf("facts missing %q: %v", want, f)
		}
	}
	if f["rule_flip"] {
		t.Error("a rule outside the top N leaked into the fact set")
	}
	// Top reads 0 as "no limit"; Facts must read it as "no rules", or a caller
	// asking for nothing would be handed everything.
	if got := s.Facts("", 0); len(got) != 0 {
		t.Errorf("Facts(_, 0) should be empty, got %v", got)
	}
	if got := s.Facts("prompt", 0); len(got) != 1 || !got["at_prompt"] {
		t.Errorf("Facts(pos, 0) should carry the position only, got %v", got)
	}
}

// A turn is scored twice — once at Stop against a transcript missing the
// reply's closing block, once at the next prompt with the whole thing. The
// second pass must correct the first rather than count as another turn.
func TestRescoringTheSameTurnReplacesIt(t *testing.T) {
	s := &Session{ID: "x"}
	now := time.Now()

	s.Record("turn-1", []string{"bold_label"}, 900, now)
	s.Record("turn-1", []string{"bold_label", "dangling_end"}, 1400, now.Add(time.Minute))

	if s.Turns != 1 {
		t.Errorf("a rescore counted as a new turn: %d turns", s.Turns)
	}
	if s.Chars != 1400 {
		t.Errorf("chars is %d, want 1400 — the rescore replaces, it does not add", s.Chars)
	}
	if len(s.Recent) != 1 {
		t.Fatalf("window holds %d entries, want 1", len(s.Recent))
	}
	if got := s.Top(0); len(got) != 2 {
		t.Errorf("tally should read the rescored rules: %+v", got)
	}
	if seen, ok := s.Scored("turn-1"); !ok || seen != 1400 {
		t.Errorf("Scored returned (%d, %v), want (1400, true)", seen, ok)
	}

	s.Record("turn-2", []string{"flip"}, 300, now.Add(2*time.Minute))
	if s.Turns != 2 || s.Chars != 1700 {
		t.Errorf("a new key should append: %d turns, %d chars", s.Turns, s.Chars)
	}
	// Only the newest turn can be rescored; anything else is a stale read.
	if _, ok := s.Scored("turn-1"); ok {
		t.Error("Scored answered for a turn that is no longer the newest")
	}
}

// Without an id there is no way to tell a rescore from a new turn, so an empty
// key always appends. Two replies that both score nothing must not merge.
func TestAnEmptyKeyAlwaysAppends(t *testing.T) {
	s := &Session{ID: "x"}
	now := time.Now()
	s.Record("", []string{"flip"}, 100, now)
	s.Record("", []string{"flip"}, 200, now)
	if s.Turns != 2 || s.Chars != 300 {
		t.Errorf("got %d turns / %d chars, want 2 / 300", s.Turns, s.Chars)
	}
	if _, ok := s.Scored(""); ok {
		t.Error("Scored should never claim to know an unkeyed turn")
	}
}

// The experiment alternates so the arms balance exactly and cover equal
// stretches of the session, and it opens on a side chosen by the session id so
// the first window is not always the same arm.
func TestWindowsAlternateAndBalance(t *testing.T) {
	s := &Session{ID: "session-one"}
	counts := map[string]int{}
	prev := ""
	for i := range 20 {
		arm := s.NextWindow(DefaultArms)
		if arm != ArmInject && arm != ArmHold {
			t.Fatalf("window %d landed on %q", i, arm)
		}
		if arm == prev {
			t.Errorf("window %d repeated %q; the arms should alternate", i, arm)
		}
		prev = arm
		counts[arm]++
	}
	if counts[ArmInject] != counts[ArmHold] {
		t.Errorf("arms are unbalanced after 20 windows: %v", counts)
	}
	if s.Windows != 20 {
		t.Errorf("window count is %d, want 20", s.Windows)
	}
	if s.Arm != prev {
		t.Errorf("Arm is %q but the last window was %q", s.Arm, prev)
	}

	// Sessions should not all open on the same side.
	opening := map[string]bool{}
	for _, id := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		opening[(&Session{ID: id}).NextWindow(DefaultArms)] = true
	}
	if len(opening) != 2 {
		t.Error("every session opened on the same arm; the id is not deciding anything")
	}
}

// A third arm has to divide the windows as evenly as two did, or the arm that
// draws the short straw carries a wider interval for the whole run.
func TestAThirdArmStillBalances(t *testing.T) {
	arms := []string{ArmInject, ArmHold, ArmPositive}
	s := &Session{ID: "session-one"}
	counts := map[string]int{}
	for range 30 {
		counts[s.NextWindow(arms)]++
	}
	for _, arm := range arms {
		if counts[arm] != 10 {
			t.Errorf("30 windows over 3 arms gave %v, want 10 each", counts)
			break
		}
	}

	// Consecutive windows must land on different arms, or a session with few
	// windows could spend its whole run on one of them.
	seen := map[string]bool{}
	prev := ""
	fresh := &Session{ID: "session-two"}
	for i := range 6 {
		arm := fresh.NextWindow(arms)
		if arm == prev {
			t.Errorf("window %d repeated %q", i, arm)
		}
		prev = arm
		seen[arm] = true
	}
	if len(seen) != 3 {
		t.Errorf("six windows over three arms reached only %v", seen)
	}
}

// A miswired flag must cost the experiment and not the prompt, so an empty arm
// list falls back rather than dividing by zero.
func TestNoArmsFallsBackInsteadOfPanicking(t *testing.T) {
	s := &Session{ID: "x"}
	arm := s.NextWindow(nil)
	if arm != ArmInject && arm != ArmHold {
		t.Errorf("empty arm list gave %q, want one of the defaults", arm)
	}
}

// The experiment's fields must stay out of a session that is not running it,
// so an ordinary state file is unchanged by the feature existing.
func TestWindowsAreAbsentUntilTheExperimentRuns(t *testing.T) {
	dir := t.TempDir()
	s := Load(dir, "quiet")
	s.Record("k", []string{"flip"}, 10, time.Now())
	if err := s.Save(dir); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "session-quiet.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"windows", "arm", "inject_chars", "inject_named"} {
		if strings.Contains(string(raw), field) {
			t.Errorf("%q is in a state file for a session not running the experiment: %s", field, raw)
		}
	}
}

func TestPruneDropsOnlyOldSessions(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	for _, id := range []string{"fresh", "stale"} {
		s := &Session{ID: id}
		s.Record("", []string{"flip"}, 10, now)
		if err := s.Save(dir); err != nil {
			t.Fatal(err)
		}
	}
	old := now.Add(-8 * 24 * time.Hour)
	if err := os.Chtimes(filepath.Join(dir, "session-stale.json"), old, old); err != nil {
		t.Fatal(err)
	}
	// An unrelated file in the same directory must survive.
	other := filepath.Join(dir, "violations.jsonl")
	if err := os.WriteFile(other, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(other, old, old); err != nil {
		t.Fatal(err)
	}

	Prune(dir, 7*24*time.Hour, now)

	if _, err := os.Stat(filepath.Join(dir, "session-fresh.json")); err != nil {
		t.Error("pruned a session inside the window")
	}
	if _, err := os.Stat(filepath.Join(dir, "session-stale.json")); !os.IsNotExist(err) {
		t.Error("stale session survived the prune")
	}
	if _, err := os.Stat(other); err != nil {
		t.Error("prune removed a file it does not own")
	}
}

// Advisory state must never take a hook down, so damage reads as "no history".
func TestCorruptStateReadsAsEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "session-bad.json"), []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if s := Load(dir, "bad"); s.Turns != 0 || s.ID != "bad" {
		t.Errorf("corrupt state should read as a fresh session, got %+v", s)
	}

	// A file whose id disagrees with the one asked for is someone else's state.
	mismatched := &Session{ID: "other"}
	mismatched.Record("", []string{"flip"}, 10, time.Now())
	if err := os.WriteFile(filepath.Join(dir, "session-mixed.json"), mustJSON(t, mismatched), 0o600); err != nil {
		t.Fatal(err)
	}
	if s := Load(dir, "mixed"); s.Turns != 0 {
		t.Error("state carrying a different session id should be discarded")
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// The tally counts turns, not airings. A phrase said three times inside one
// reply is one turn using it, and the rule it feeds asks how many turns have.
func TestPriorPhrasesCountsTurns(t *testing.T) {
	now := time.Now()
	s := &Session{ID: "s"}
	s.Record("a", nil, 10, now, 1, 2)
	s.Record("b", nil, 10, now, 2, 3)
	s.Record("c", nil, 10, now, 2)

	got := s.PriorPhrases("")
	for sum, want := range map[uint32]int{1: 1, 2: 3, 3: 1} {
		if got[sum] != want {
			t.Errorf("phrase %d counted %d times, want %d", sum, got[sum], want)
		}
	}
}

// A turn is scored twice — once at Stop against an unfinished transcript, once
// on the next prompt against the complete one. The second pass must not find
// the first pass's phrases sitting in the window and report the reply for
// repeating itself.
func TestPriorPhrasesSkipsTheTurnBeingRescored(t *testing.T) {
	now := time.Now()
	s := &Session{ID: "s"}
	s.Record("older", nil, 10, now, 7)
	s.Record("again", nil, 10, now, 7, 8)

	if got := s.PriorPhrases("again")[7]; got != 1 {
		t.Errorf("phrase 7 counted %d times while rescoring its own turn, want 1", got)
	}
	if got := s.PriorPhrases("again")[8]; got != 0 {
		t.Errorf("a phrase only this turn used counted %d times, want 0", got)
	}
	// An empty key is the backfill and the un-keyed hook: nothing to exclude.
	if got := s.PriorPhrases("")[7]; got != 2 {
		t.Errorf("with no key to skip, phrase 7 should count %d turns, got %d", 2, got)
	}
}

// A rescore replaces the turn's phrases rather than adding to them, for the
// same reason it replaces its rules: it is the same reply, read further.
func TestRescoreReplacesPhrases(t *testing.T) {
	now := time.Now()
	s := &Session{ID: "s"}
	s.Record("k", nil, 10, now, 1)
	s.Record("k", nil, 20, now, 1, 2)
	if s.Turns != 1 {
		t.Fatalf("rescore counted as %d turns", s.Turns)
	}
	if got := s.Recent[0].Phrases; len(got) != 2 {
		t.Errorf("phrases after rescore = %v, want the second pass's two", got)
	}
}
