// Command cope-gate checks a drafted Claude Code reply against a voice card.
//
// Installed as a Stop hook, it reads the hook payload on stdin, pulls the last
// assistant turn out of the transcript, and reports violations. Warn-only by
// default: it prints and logs, and exits 0 so nothing is blocked. The point of
// starting warn-only is to learn the false-positive rate before a violation
// starts costing a regeneration.
//
// The card is compiled in, so there is nothing to install beside the binary:
//
//	SessionStart      cope-gate --inject       puts the card into context
//	Stop              cope-gate                checks what came out
//	UserPromptSubmit  cope-gate --refresher    restates the rule that decays first
//
// Backfill mode scores a whole session at once, which is how the rules were
// calibrated in the first place:
//
//	cope-gate --backfill ~/.claude/projects/<proj>/<session>.jsonl
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/justinstimatze/cope/card"
	"github.com/justinstimatze/cope/internal/scan"
	"github.com/justinstimatze/cope/internal/transcript"
)

// Overridden at release via -ldflags "-X main.version=$(git describe ...)".
var version = "dev"

func buildVersion() string {
	if version != "dev" {
		return version
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	var rev string
	var dirty bool
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}
	if rev == "" {
		return version
	}
	if len(rev) > 12 {
		rev = rev[:12]
	}
	if dirty {
		return rev + "-dirty"
	}
	return rev
}

// hookInput is the Stop hook payload. Only the transcript path is needed.
type hookInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	StopHookActive bool   `json:"stop_hook_active"`
}

func main() {
	var (
		rulesPath = flag.String("rules", "", "read the card from this JSON file instead of the one built into the binary")
		logPath   = flag.String("log", defaultLogPath(), "append violations here; empty disables")
		backfill  = flag.String("backfill", "", "score every assistant turn in this transcript and exit")
		inject    = flag.Bool("inject", false, "print the card as prompt text for a SessionStart hook")
		minCV     = flag.Float64("min-cv", 0.35, "flag paragraph-length coefficient of variation below this")
		block     = flag.Bool("block", false, "exit 2 on a violation whose action is reject (default warn-only)")
		refresher = flag.Bool("refresher", false, "UserPromptSubmit entry: inject the compact reminder once the last injection has aged")
		refEvery  = flag.Duration("refresh-every", 4*time.Hour, "minimum age of the last card or refresher injection before the refresher fires")
		showVer   = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *showVer {
		fmt.Println(buildVersion())
		return
	}

	c, warnings, err := loadCard(*rulesPath)
	if err != nil {
		// A missing or broken card must never wedge the session.
		fmt.Fprintf(os.Stderr, "cope-gate: %v (no check run)\n", err)
		return
	}
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", w)
	}

	if *inject {
		fmt.Print(c.Render())
		// A fresh full card resets the refresher clock: the reminder is for
		// the stretch after an injection has aged, and SessionStart (startup,
		// resume, compact) just was one.
		if id := readSessionID(); id != "" {
			stampRefresher(id)
		}
		return
	}

	if *refresher {
		runRefresher(c, *refEvery)
		return
	}

	if *backfill != "" {
		runBackfill(c, *backfill, *minCV)
		return
	}

	var in hookInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: decode hook input: %v\n", err)
		return
	}
	if in.StopHookActive {
		return // already re-entered once; do not loop
	}
	text, err := transcript.LastAssistantText(in.TranscriptPath)
	if err != nil || strings.TrimSpace(text) == "" {
		return
	}

	v := c.Check(text, *minCV)
	if len(v) == 0 {
		return
	}
	report(v, in.SessionID, *logPath)

	if *block {
		for _, x := range v {
			if x.Action == "reject" {
				os.Exit(2)
			}
		}
	}
}

// loadCard returns the card at path, or the one compiled into the binary when
// path is empty.
//
// The built-in card is the default because `go install` copies no data files.
// The old default guessed a checkout path under $HOME and read card/rules.json
// from it, which worked on exactly one machine: anywhere else the hook printed
// "no check run" on every turn and checked nothing.
func loadCard(path string) (*scan.Card, []error, error) {
	if path != "" {
		return scan.LoadCard(path)
	}
	return scan.ParseCard(card.RulesJSON)
}

// --- state directory --------------------------------------------------------

// stateDir resolves $XDG_STATE_HOME/cope, falling back to ~/.local/state/cope,
// and returns "" when neither is available. It does not create anything, so
// --version and --inject stay free of filesystem side effects.
//
// The refresher markers and the violations log both live here. The log in
// particular holds verbatim excerpts of the assistant's replies, so it belongs
// in per-machine state rather than in the repo checkout, where the previous
// default put it one `git add -A` away from being published.
func stateDir() string {
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "cope")
}

func defaultLogPath() string {
	dir := stateDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "violations.jsonl")
}

// --- mid-session refresher --------------------------------------------------
//
// The marker-file clock is basanite's claimInjection pattern (basanite
// cmd/basanite/main.go) with one inversion: basanite injects on first sight
// and re-injects after its interval, while the refresher stays SILENT on
// first sight — SessionStart just delivered the full card — and only fires
// once the last injection has aged past -refresh-every.

var sessionIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{8,128}$`)

// readSessionID decodes the hook payload's session_id from stdin, or returns
// "" for anything else — including an interactive TTY, so a manual
// `cope-gate --inject` in a terminal doesn't hang waiting for input.
func readSessionID() string {
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice != 0 {
		return ""
	}
	var in hookInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return ""
	}
	if !sessionIDRe.MatchString(in.SessionID) {
		return ""
	}
	return in.SessionID
}

func markerPath(sessionID string) (string, error) {
	dir := stateDir()
	if dir == "" {
		return "", errors.New("no state directory: $XDG_STATE_HOME and $HOME are both unset")
	}
	// 0700 because this directory holds the violations log, which quotes the
	// assistant's prose back verbatim.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "refresher-"+sessionID), nil
}

// runRefresher is the UserPromptSubmit entry. Every failure path returns
// silently with exit 0 — a style reminder must never take a prompt down.
func runRefresher(c *scan.Card, every time.Duration) {
	id := readSessionID()
	if id == "" {
		return
	}
	out := c.RenderRefresher()
	if out == "" {
		return
	}
	marker, err := markerPath(id)
	if err != nil {
		return
	}
	if !claimRefresh(marker, every) {
		return
	}
	pruneMarkers(filepath.Dir(marker))
	fmt.Print(out)
}

// claimRefresh reports whether this prompt should carry the refresher. First
// sight of a session stamps the clock and declines; after that, a marker
// older than every fires and restamps. O_EXCL keeps concurrent first prompts
// to a single stamp.
func claimRefresh(marker string, every time.Duration) bool {
	fi, err := os.Stat(marker)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(marker, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			f.Close()
		}
		return false
	}
	if err != nil || time.Since(fi.ModTime()) < every {
		return false
	}
	now := time.Now()
	return os.Chtimes(marker, now, now) == nil
}

// stampRefresher marks now as the last injection time for the session,
// creating the marker if SessionStart beat the first prompt to it.
func stampRefresher(sessionID string) {
	marker, err := markerPath(sessionID)
	if err != nil {
		return
	}
	if f, err := os.OpenFile(marker, os.O_CREATE|os.O_WRONLY, 0o600); err == nil {
		f.Close()
	}
	now := time.Now()
	_ = os.Chtimes(marker, now, now)
}

// pruneMarkers drops markers older than a week: sessions end without a
// goodbye, so the state dir would otherwise grow one file per session forever.
func pruneMarkers(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "refresher-") {
			continue
		}
		fi, err := e.Info()
		if err == nil && time.Since(fi.ModTime()) > 7*24*time.Hour {
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}

func report(v []scan.Violation, session, logPath string) {
	byRule := map[string]int{}
	for _, x := range v {
		byRule[x.RuleID]++
	}
	ids := make([]string, 0, len(byRule))
	for id := range byRule {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	fmt.Fprintf(&b, "cope: %d violation(s) —", len(v))
	for _, id := range ids {
		fmt.Fprintf(&b, " %s×%d", id, byRule[id])
	}
	b.WriteString("\n")
	for _, x := range v {
		fmt.Fprintf(&b, "  [%s] %s\n      %q\n", x.RuleID, x.Why, x.Matched)
	}
	fmt.Fprint(os.Stderr, b.String())

	if logPath == "" {
		return
	}
	if dir := filepath.Dir(logPath); dir != "" {
		_ = os.MkdirAll(dir, 0o700)
	}
	// 0600: every record carries a verbatim excerpt of the reply.
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	stamp := time.Now().UTC().Format(time.RFC3339)
	enc := json.NewEncoder(f)
	for _, x := range v {
		_ = enc.Encode(struct {
			Time    string `json:"time"`
			Session string `json:"session"`
			scan.Violation
		}{stamp, session, x})
	}
}

func runBackfill(c *scan.Card, path string, minCV float64) {
	turns, err := transcript.AllAssistantText(path)
	if err != nil {
		// A truncated read still leaves a usable sample. Say what was lost and
		// score the rest; only an empty result is worth failing on.
		fmt.Fprintf(os.Stderr, "cope-gate: %v (scoring the %d turn(s) read so far)\n", err, len(turns))
		if len(turns) == 0 {
			os.Exit(1)
		}
	}
	chars := 0
	for _, t := range turns {
		chars += len(t)
	}
	fmt.Printf("%d assistant turns, %d chars\n\n", len(turns), chars)

	byRule := map[string]int{}
	examples := map[string][]scan.Violation{}
	turnsHit := 0
	for _, t := range turns {
		v := c.Check(t, minCV)
		if len(v) > 0 {
			turnsHit++
		}
		for _, x := range v {
			byRule[x.RuleID]++
			if len(examples[x.RuleID]) < 4 {
				examples[x.RuleID] = append(examples[x.RuleID], x)
			}
		}
	}
	fmt.Printf("turns with at least one hit: %d/%d\n\n", turnsHit, len(turns))

	ids := make([]string, 0, len(byRule))
	for id := range byRule {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return byRule[ids[i]] > byRule[ids[j]] })
	for _, id := range ids {
		fmt.Printf("=== %s (%d hits)\n", id, byRule[id])
		for _, x := range examples[id] {
			fmt.Printf("  %q\n", x.Matched)
			fmt.Printf("    ...%s...\n", x.Context)
		}
		fmt.Println()
	}
	for _, r := range c.Rules {
		if byRule[r.ID] == 0 {
			fmt.Printf("(never fired: %s — /%s/)\n", r.ID, r.Pattern)
		}
	}
}
