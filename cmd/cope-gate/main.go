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
//	PreToolUse        cope-gate --pretool      checks prose leaving the conversation
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
	"github.com/justinstimatze/cope/internal/effigy"
	"github.com/justinstimatze/cope/internal/scan"
	"github.com/justinstimatze/cope/internal/state"
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

// options is the parsed command line.
//
// Registration lives apart from main because --author-docs describes the flag
// set by walking it, and a flag registered inside main only exists once main
// runs. A test that asked what the binary accepts got the testing package's
// flags instead, which is the sort of quiet wrong answer this whole repo is
// about. Both callers now register through here.
type options struct {
	rulesPath   *string
	cardName    *string
	listCards   *bool
	logPath     *string
	backfill    *string
	check       *string
	authorDoc   *bool
	inject      *bool
	outputStyle *bool
	setup       *bool
	dryRun      *bool
	styleDir    *string
	describe    *bool
	minCV       *float64
	block       *bool
	refresher   *bool
	refEvery    *time.Duration
	showVer     *bool
	ab          *bool
	abArms      *string
	abReport    *string
	renderArm   *string
	renderFor   *string
	display     *bool
	displayIn   *bool
	pretool     *bool
	renderLane  *string
	checkLane   *string
	cardSample  *string
}

func registerFlags(fs *flag.FlagSet) *options {
	return &options{
		rulesPath:   fs.String("rules", "", "read the card from this .effigy or .json file instead of the one built into the binary"),
		cardName:    fs.String("card", "", "name of an installed card to write in, from the cards directory; also $COPE_CARD"),
		listCards:   fs.Bool("cards", false, "list the installed cards with the aim each one states, and exit"),
		logPath:     fs.String("log", defaultLogPath(), "append violations here; empty disables"),
		backfill:    fs.String("backfill", "", "score every assistant turn in this transcript and exit"),
		check:       fs.String("check", "", "score a prose file against the card and exit; - reads stdin"),
		authorDoc:   fs.Bool("author-docs", false, "print a prompt for writing this repo's docs: the card, the introspected facts, the sections"),
		inject:      fs.Bool("inject", false, "print the card as prompt text for a SessionStart hook"),
		setup:       fs.Bool("setup", false, "emit the output style and wire the hooks, then print the one step left"),
		dryRun:      fs.Bool("dry-run", false, "with --setup, print what would change and write nothing"),
		outputStyle: fs.Bool("output-style", false, "write the card to ~/.claude/output-styles as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message"),
		styleDir:    fs.String("output-style-dir", "", "directory to write the output style into (default ~/.claude/output-styles)"),
		describe:    fs.Bool("describe", false, "print the card's voice as a target to recognise: the aim and the register, without the machinery"),
		minCV:       fs.Float64("min-cv", 0.35, "flag paragraph-length coefficient of variation below this"),
		block:       fs.Bool("block", false, "exit 2 on a violation whose action is reject (default warn-only)"),
		refresher:   fs.Bool("refresher", false, "UserPromptSubmit entry: inject the compact reminder once the last injection has aged"),
		refEvery:    fs.Duration("refresh-every", 30*time.Minute, "minimum age of the last card or refresher injection before the refresher fires"),
		showVer:     fs.Bool("version", false, "print version and exit"),
		ab:          fs.Bool("ab", false, "rotate refresher windows through each variant in turn, recording which one a turn was written under"),
		abArms:      fs.String("ab-arms", "", "comma-separated variants to rotate through, implying -ab (default inject,hold; positive is the third)"),
		abReport:    fs.String("ab-report", "", "read a turn log and report how each variant did; - reads the default path"),
		renderArm:   fs.String("render-arm", "", "print the mid-session reminder one variant would inject, and exit"),
		display:     fs.Bool("display", false, "MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone"),
		pretool:     fs.Bool("pretool", false, "PreToolUse entry: score the prose an external write is about to post, warn-only"),
		displayIn:   fs.Bool("display-preview", false, "read prose on stdin and print it as --display would rewrite it"),
		renderFor:   fs.String("render-for", "", "comma-separated rule ids to render -render-arm against"),
		renderLane:  fs.String("render-lane", "", "render -render-arm as the given lane sees it: interactive (default) or loop"),
		checkLane:   fs.String("check-lane", "", "score -check in the given lane: interactive (default), loop, or external"),
		cardSample:  fs.String("card-from-sample", "", "print a prompt for writing a card from this writing sample; - reads stdin"),
	}
}

func main() {
	opt := registerFlags(flag.CommandLine)
	flag.Parse()
	rulesPath, logPath, backfill := opt.rulesPath, opt.logPath, opt.backfill
	check, authorDoc, inject := opt.check, opt.authorDoc, opt.inject
	minCV, block := opt.minCV, opt.block
	refresher, refEvery, showVer := opt.refresher, opt.refEvery, opt.showVer

	if *showVer {
		fmt.Println(buildVersion())
		return
	}

	if *opt.displayIn {
		runDisplayPreview()
		return
	}
	if *opt.display {
		runDisplay(*opt.logPath)
		return
	}
	if *opt.abReport != "" {
		runABReport(*opt.abReport)
		return
	}
	if *opt.listCards {
		runListCards()
		return
	}

	c, warnings, err := loadCard(*rulesPath, *opt.cardName)
	if err != nil {
		// A missing or broken card must never wedge the session.
		fmt.Fprintf(os.Stderr, "cope-gate: %v (no check run)\n", err)
		return
	}
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", w)
	}

	if *opt.renderArm != "" {
		runRenderArm(c, *opt.renderArm, *opt.renderFor, *opt.renderLane)
		return
	}

	if *opt.pretool {
		runPreTool(c, *minCV, *logPath)
		return
	}

	if *opt.describe {
		fmt.Print(c.Describe())
		return
	}

	if *opt.setup {
		os.Exit(runSetup(c, *opt.dryRun))
	}

	if *opt.outputStyle {
		os.Exit(runOutputStyle(c, *opt.styleDir, false))
	}

	if *inject {
		// An output style already puts a card in the system prompt, where it
		// outranks anything printed here. Printing a second one is how a
		// session ends up holding two registers and following whichever was
		// repeated most recently.
		if name, ok := styleSupersedesInject(); ok {
			fmt.Fprintf(os.Stderr, "cope-gate: output style %q is active, so the card is already in the system prompt — not injecting\n", name)
			return
		}
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
		arms, err := parseArms(*opt.abArms, *opt.ab)
		if err != nil {
			// A bad arm list disables the experiment rather than the reminder.
			fmt.Fprintf(os.Stderr, "cope-gate: %v (experiment off)\n", err)
		}
		runRefresher(c, *refEvery, *minCV, *logPath, arms)
		return
	}

	if *opt.cardSample != "" {
		runCardFromSample(*opt.cardSample)
		return
	}

	if *authorDoc {
		fmt.Print(buildDocsPrompt(c, flag.CommandLine))
		return
	}

	if *check != "" {
		runCheck(c, *check, *minCV, *logPath, *opt.checkLane)
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
	// The whole turn, not the last record — and still not the whole reply,
	// because its closing block reaches the transcript after this hook has
	// run. catchUp, on the next prompt, scores the rest.
	turn, err := transcript.LastTurn(in.TranscriptPath)
	if err != nil || strings.TrimSpace(turn.Text) == "" {
		return
	}

	v, phrases := checkTurn(c, in.SessionID, turn, *minCV)
	recordTurn(in.SessionID, turn, v, len(turn.Text), phrases)
	if len(v) == 0 {
		return
	}
	report("cope:", turn.Text, v, in.SessionID, *logPath)

	if *block {
		for _, x := range v {
			if x.Action == "reject" {
				os.Exit(2)
			}
		}
	}
}

// loadCard returns the card to write in. A path given with --rules wins, then
// a name from --card or $COPE_CARD, then the one compiled into the binary.
//
// The built-in card is the default because `go install` copies no data files.
// The old default guessed a checkout path under $HOME and read card/rules.json
// from it, which worked on exactly one machine: anywhere else the hook printed
// "no check run" on every turn and checked nothing.
//
// A named card that cannot be found is an error rather than a fall back to the
// built-in one. Falling back would leave a session writing in the shipped voice
// while its config says otherwise, which is the failure that looks like the
// card not working rather than like a typo.
func loadCard(path, name string) (*scan.Card, []error, error) {
	if path == "" && name == "" {
		name = os.Getenv("COPE_CARD")
	}
	if path == "" && name != "" {
		var err error
		if path, err = resolveCard(name); err != nil {
			return nil, nil, err
		}
	}
	if path == "" {
		return scan.ParseCard(card.RulesJSON)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasSuffix(path, ".effigy") {
		if raw, err = effigy.ToJSON(raw, path); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", path, err)
		}
	}
	c, warnings, err := scan.ParseCard(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return c, warnings, nil
}

// resolveCard turns a --card value into a file. Anything with a separator or a
// known extension is taken as a path so a card can be tried before it is
// installed; a bare word is looked up in the cards directory.
func resolveCard(name string) (string, error) {
	if strings.ContainsRune(name, filepath.Separator) ||
		strings.HasSuffix(name, ".effigy") || strings.HasSuffix(name, ".json") {
		if _, err := os.Stat(name); err != nil {
			return "", fmt.Errorf("card %s: %w", name, err)
		}
		return name, nil
	}
	dir := cardsDir()
	if dir == "" {
		return "", fmt.Errorf("card %q: no cards directory (set $XDG_CONFIG_HOME or $HOME)", name)
	}
	for _, ext := range []string{".effigy", ".json"} {
		p := filepath.Join(dir, name+ext)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no card named %q in %s — run --cards to see what is installed", name, dir)
}

// cardsDir resolves $XDG_CONFIG_HOME/cope/cards, falling back to
// ~/.config/cope/cards, and returns "" when neither is available.
//
// Cards are configuration rather than state: a card is written by hand, is
// meant to be edited and kept, and is the one file a reader of this tool will
// want to find again. The violations log and the refresher markers stay in
// stateDir, which is machine-local and disposable.
func cardsDir() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "cope", "cards")
}

// runListCards prints what is installed, each with the aim its card states, so
// picking one does not mean opening five files. The built-in card is listed
// first because it is what a session gets with nothing set.
func runListCards() {
	type entry struct{ name, path string }
	found := []entry{{name: "built-in"}}
	dir := cardsDir()
	if dir != "" {
		var paths []string
		for _, ext := range []string{"*.effigy", "*.json"} {
			matches, err := filepath.Glob(filepath.Join(dir, ext))
			if err != nil {
				continue
			}
			paths = append(paths, matches...)
		}
		sort.Strings(paths)
		for _, p := range paths {
			found = append(found, entry{strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)), p})
		}
	}
	width := 0
	for _, e := range found {
		if len(e.name) > width {
			width = len(e.name)
		}
	}
	for _, e := range found {
		fmt.Printf("%-*s  %s\n", width, e.name, themeOf(e.path))
	}
	if len(found) == 1 && dir != "" {
		fmt.Fprintf(os.Stderr, "\nno cards in %s\n", dir)
	}
}

// themeOf reads one card's stated aim for the listing. A card that will not
// parse is reported in the listing rather than skipped, because a card the
// reader wrote and cannot see is worse than one shown with its error.
func themeOf(path string) string {
	c, _, err := loadCard(path, "")
	switch {
	case err != nil:
		return "(" + err.Error() + ")"
	case c.Theme == "":
		return "(no aim stated)"
	default:
		return c.Theme
	}
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

// readHookInput decodes the hook payload from stdin, or reports false for
// anything else — including an interactive TTY, so a manual `cope-gate
// --inject` in a terminal doesn't hang waiting for input.
func readHookInput() (hookInput, bool) {
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice != 0 {
		return hookInput{}, false
	}
	var in hookInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return hookInput{}, false
	}
	if !sessionIDRe.MatchString(in.SessionID) {
		return hookInput{}, false
	}
	return in, true
}

// readSessionID is readHookInput for the callers that only need the id.
func readSessionID() string {
	in, ok := readHookInput()
	if !ok {
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

// recordTurn appends this reply's score to the session's rolling state, so the
// next UserPromptSubmit can say what has actually been happening. Every failure
// is swallowed: this is advisory, and losing it costs one turn of
// personalisation rather than the hook.
func recordTurn(sessionID string, t transcript.Turn, v []scan.Violation, chars int, phrases []uint32) {
	if !sessionIDRe.MatchString(sessionID) {
		return
	}
	dir := stateDir()
	if dir == "" {
		return
	}
	s := state.Load(dir, sessionID)
	s.Lane = t.Lane
	s.Record(t.Key, ruleIDs(v), chars, time.Now(), phrases...)
	_ = s.Save(dir)
	logTurn(dir, s, t, ruleIDs(v), chars)
	state.Prune(dir, 7*24*time.Hour, time.Now())
}

// checkTurn scores one reply against the card and against what the session has
// already said, returning the violations and the phrases to carry forward.
//
// Every rule but one reads the reply alone, which is what Card.Check does.
// cross_turn_repeat is the exception and the reason cope is not a linter: it
// needs the window, the window is on disk, and so this is the only rule that
// goes quiet when the state directory is missing rather than reporting on a
// session it cannot see.
func checkTurn(c *scan.Card, sessionID string, t transcript.Turn, minCV float64) ([]scan.Violation, []uint32) {
	dir := stateDir()
	if dir == "" || !sessionIDRe.MatchString(sessionID) {
		return c.CheckLane(t.Text, minCV, t.Lane), scan.Sums(scan.Phrases(t.Text))
	}
	return checkAgainst(c, state.Load(dir, sessionID), t, minCV)
}

// checkAgainst is checkTurn for a caller already holding the session.
func checkAgainst(c *scan.Card, s *state.Session, t transcript.Turn, minCV float64) ([]scan.Violation, []uint32) {
	ps := scan.Phrases(t.Text)
	// Allow wraps the whole append: cross_turn_repeat is produced out here
	// rather than inside CheckLane, so a card exempting it would otherwise be
	// the one @gate line that silently did nothing.
	v := c.Allow(append(c.CheckLane(t.Text, minCV, t.Lane), scan.CrossTurnRepeat(ps, s.PriorPhrases(t.Key))...))
	return v, scan.Sums(ps)
}

func turnLogPath(dir string) string { return filepath.Join(dir, "turns.jsonl") }

// logTurn appends one scored turn to the experiment's record, and does nothing
// while the experiment is not running — Arm is empty until --ab claims a
// window.
//
// Rule names, counts and a character total, with no prose. violations.jsonl
// quotes the reply back; this file cannot, which is what makes it the one worth
// keeping past the seven-day prune.
//
// Appended rather than rewritten, so a turn scored at Stop and again on the
// next prompt appears twice. runABReport keeps the last record for a key, which
// is the pass that saw the whole reply.
func logTurn(dir string, s *state.Session, t transcript.Turn, rules []string, chars int) {
	if s.Arm == "" {
		return
	}
	f, err := os.OpenFile(turnLogPath(dir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(turnRecord{
		Time:        time.Now().UTC().Format(time.RFC3339),
		Session:     s.ID,
		Key:         t.Key,
		Lane:        t.Lane,
		Arm:         s.Arm,
		Chars:       chars,
		Rules:       rules,
		Version:     buildVersion(),
		InjectChars: s.InjectChars,
		InjectNamed: s.InjectNamed,
	})
}

// turnRecord is one scored turn, tagged with the arm that was live when it was
// written and with the shape of what that arm injected.
//
// InjectChars and InjectNamed are what let a difference between arms be
// attributed. The positive arm is shorter than the paired one and names no
// rules, so a gap between them has three candidate causes, and only a record of
// the reminder itself can start to separate them.
//
// Version is which binary scored the turn. The hook names an absolute path and
// `go install` replaces the file underneath it, so a run can span several builds
// with nothing in the data to say so — the rules that scored the first half are
// then not the rules that scored the second, and the report would average across
// the change without noticing. This happened: an arm sat behind a flag for a day
// while the installed binary predated it, and comparing file mtimes was the only
// way that surfaced.
type turnRecord struct {
	Time        string   `json:"time"`
	Session     string   `json:"session"`
	Key         string   `json:"key"`
	Lane        string   `json:"lane,omitempty"`
	Arm         string   `json:"arm"`
	Chars       int      `json:"chars"`
	Rules       []string `json:"rules"`
	Version     string   `json:"version,omitempty"`
	InjectChars int      `json:"inject_chars,omitempty"`
	InjectNamed []string `json:"inject_named,omitempty"`
}

func ruleIDs(v []scan.Violation) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.RuleID)
	}
	return out
}

// catchUp rescores the turn the Stop hook could only see part of.
//
// Claude Code appends a reply's closing block to the transcript after the Stop
// hook has already returned, so every Stop scores a turn that is missing its
// ending — which is where dangling_end, buried_decision and the last ask live.
// By the time the next prompt arrives the file is complete.
//
// Only what the Stop pass could not have seen is reported. The text it did see
// is a byte prefix of the finished turn, so rescoring that prefix reproduces
// exactly what was already printed and logged, and the difference is the rest.
//
// Every failure returns silently: this runs on the prompt path, where a style
// tool must never cost the user a turn.
func catchUp(c *scan.Card, in hookInput, minCV float64, logPath string) {
	dir := stateDir()
	if in.TranscriptPath == "" || dir == "" {
		return
	}
	turn, err := transcript.LastTurn(in.TranscriptPath)
	if err != nil || strings.TrimSpace(turn.Text) == "" {
		return
	}
	s := state.Load(dir, in.SessionID)
	seen, scored := s.Scored(turn.Key)
	if scored && seen >= len(turn.Text) {
		return // nothing arrived after the Stop hook read it
	}

	full, phrases := checkAgainst(c, s, turn, minCV)
	fresh := full
	if scored && seen > 0 && seen <= len(turn.Text) {
		head := turn
		head.Text = turn.Text[:seen]
		prefix, _ := checkAgainst(c, s, head, minCV)
		fresh = unreported(prefix, full)
	}
	s.Lane = turn.Lane
	s.Record(turn.Key, ruleIDs(full), len(turn.Text), time.Now(), phrases...)
	_ = s.Save(dir)
	logTurn(dir, s, turn, ruleIDs(full), len(turn.Text))
	if len(fresh) > 0 {
		report("cope — the close of that reply:", turn.Text, fresh, in.SessionID, logPath)
	}
}

// unreported returns the violations in full that prev does not already account
// for, as a multiset so a repeated tic is not collapsed to one. Violation is
// all strings, so it is comparable and needs no key function.
func unreported(prev, full []scan.Violation) []scan.Violation {
	seen := make(map[scan.Violation]int, len(prev))
	for _, x := range prev {
		seen[x]++
	}
	var out []scan.Violation
	for _, x := range full {
		if seen[x] > 0 {
			seen[x]--
			continue
		}
		out = append(out, x)
	}
	return out
}

// focusedRefresher renders the card narrowed to what this session has been
// getting wrong. Empty when there is no history yet, which is the common case
// early in a session and is why the caller keeps the standing reminder as a
// fallback.
// It returns the rules the reminder spoke to alongside the text, because an arm
// is only a name unless what it injected is written down beside the turns it
// produced.
func focusedRefresher(c *scan.Card, sessionID, arm string) (text string, named []string) {
	dir := stateDir()
	if dir == "" {
		return "", nil
	}
	s := state.Load(dir, sessionID)
	top := s.Top(focusRules)
	if len(top) == 0 {
		return "", nil
	}
	observed := make([]scan.Observation, 0, len(top))
	for _, t := range top {
		observed = append(observed, scan.Observation{Rule: t.Rule, N: t.N})
		named = append(named, t.Rule)
	}
	facts := s.Facts("prompt", focusRules)
	if arm == state.ArmPositive {
		return c.RenderFocusPositive(facts, observed), named
	}
	return c.RenderFocus(facts, observed), named
}

// focusRules bounds how many rules the mid-session injection speaks to. The
// point is to name the one or two habits currently in play; a list of six is
// the flat card again, with worse manners.
const focusRules = 2

// runRefresher is the UserPromptSubmit entry. Every failure path returns
// silently with exit 0 — a style reminder must never take a prompt down.
//
// It scores before it injects. catchUp finishes the previous turn, which the
// Stop hook necessarily saw without its ending, and that has to happen on every
// prompt rather than on the ~4% that carry a refresher: the state it repairs is
// what the injection is then sliced from.
// An empty arms list is the experiment switched off, which renders and prints
// exactly as it did before the experiment existed.
func runRefresher(c *scan.Card, every time.Duration, minCV float64, logPath string, arms []string) {
	in, ok := readHookInput()
	if !ok {
		return
	}
	id := in.SessionID
	catchUp(c, in, minCV, logPath)

	marker, err := markerPath(id)
	if err != nil {
		return
	}
	if !claimRefresh(marker, every) {
		return
	}
	pruneMarkers(filepath.Dir(marker))

	// The window is claimed before anything is rendered, so a held-back window
	// costs the same interval as an injected one and every arm covers an equal
	// stretch of the session. Which arm is live also decides which renderer
	// runs, so it has to be settled first.
	arm := state.ArmInject
	if len(arms) > 0 {
		arm = nextArm(id, arms)
	}
	if arm == state.ArmHold {
		recordInjection(id, 0, nil)
		return
	}

	out, named := focusedRefresher(c, id, arm)
	if out == "" {
		// The standing reminder is the same text in every arm, so a session
		// with no history yet contributes nothing that distinguishes them.
		out = c.RenderRefresher()
		named = nil
	}
	if out == "" {
		return
	}
	recordInjection(id, len(out), named)
	fmt.Print(out)
}

// runRenderArm prints what one arm would inject against a given set of fired
// rules, and exits. It exists so the arms can be driven from outside the hook —
// by an offline replay that generates prose under each arm, and by a reader who
// wants to see what the experiment is actually putting in front of the writer.
//
// ArmHold renders nothing and says so on stderr, which is the honest answer
// rather than an empty file that reads like a bug.
func runRenderArm(c *scan.Card, arm, rules, lane string) {
	if arm == state.ArmHold {
		fmt.Fprintln(os.Stderr, "cope-gate: hold injects nothing")
		return
	}
	if arm != state.ArmInject && arm != state.ArmPositive {
		fmt.Fprintf(os.Stderr, "cope-gate: unknown arm %q\n", arm)
		os.Exit(1)
	}

	// Through state.Facts rather than a literal map, so this prints what the
	// live refresher would print. Built by hand it set at_prompt and mode_loop
	// together, which no session ever does — and the loop render then carried
	// the CONTINUE TEST and a pair whose better half ends in a question, beside
	// the rule against ending a loop by asking.
	facts := (&state.Session{Lane: lane}).Facts("prompt", 0)
	var observed []scan.Observation
	for _, id := range strings.Split(rules, ",") {
		if id = strings.TrimSpace(id); id == "" {
			continue
		}
		facts["rule_"+id] = true
		// The count is a stand-in: this path has no session to read one from,
		// and only the paired render prints it at all.
		observed = append(observed, scan.Observation{Rule: id, N: 1})
	}

	out := c.RenderFocus(facts, observed)
	if arm == state.ArmPositive {
		out = c.RenderFocusPositive(facts, observed)
	}
	if out == "" {
		// Same fallback the live refresher uses, so what this prints is what a
		// session with no history would actually have seen.
		out = c.RenderRefresher()
	}
	fmt.Print(out)
}

// parseArms turns the flag pair into the arm rotation. -ab alone is the
// two-arm default; -ab-arms names the set and implies -ab.
//
// An unknown arm is an error rather than a silently dropped one: a typo that
// quietly halved the experiment would look exactly like data.
func parseArms(spec string, ab bool) ([]string, error) {
	if spec == "" {
		if ab {
			return state.DefaultArms, nil
		}
		return nil, nil
	}
	known := map[string]bool{state.ArmInject: true, state.ArmHold: true, state.ArmPositive: true}
	var arms []string
	for _, name := range strings.Split(spec, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if !known[name] {
			return nil, fmt.Errorf("unknown arm %q", name)
		}
		arms = append(arms, name)
	}
	if len(arms) < 2 {
		return nil, fmt.Errorf("an experiment needs at least two arms, got %d", len(arms))
	}
	return arms, nil
}

// nextArm advances the experiment and returns which side this window is on.
// Failure to persist returns ArmInject: the untested state of the tool is the
// one that ships, so a lost write costs a data point rather than the reminder.
func nextArm(sessionID string, arms []string) string {
	dir := stateDir()
	if dir == "" {
		return state.ArmInject
	}
	s := state.Load(dir, sessionID)
	arm := s.NextWindow(arms)
	if err := s.Save(dir); err != nil {
		return state.ArmInject
	}
	return arm
}

// recordInjection stores what this window put in front of the writer, so the
// turns logged under it carry the reminder's size and subject rather than only
// its arm name. Silent on failure, like every other write on this path.
func recordInjection(sessionID string, chars int, named []string) {
	dir := stateDir()
	if dir == "" {
		return
	}
	s := state.Load(dir, sessionID)
	s.RecordInjection(chars, named)
	_ = s.Save(dir)
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

// report prints the violations and appends them to the log. lead names which
// pass found them, because the Stop hook and the catch-up on the next prompt
// both report and the reader should know which reply is being described.
func report(lead, text string, v []scan.Violation, session, logPath string) {
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
	fmt.Fprintf(&b, "%s %d violation(s) —", lead, len(v))
	for _, id := range ids {
		fmt.Fprintf(&b, " %s×%d", id, byRule[id])
	}
	b.WriteString("\n")
	for _, x := range v {
		fmt.Fprintf(&b, "  [%s] %s\n      %q\n", x.RuleID, x.Why, x.Matched)
	}
	b.WriteString(scan.ClusterLine(text, v))
	fmt.Fprint(os.Stderr, b.String())
	logViolations(v, session, logPath)
}

// logViolations appends the log records without printing. --check does its own
// printing to stdout, and calling report there put the same violations on both
// streams.
func logViolations(v []scan.Violation, session, logPath string) {
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
	turns, err := transcript.AllTurns(path)
	if err != nil {
		// A truncated read still leaves a usable sample. Say what was lost and
		// score the rest; only an empty result is worth failing on.
		fmt.Fprintf(os.Stderr, "cope-gate: %v (scoring the %d turn(s) read so far)\n", err, len(turns))
		if len(turns) == 0 {
			os.Exit(1)
		}
	}
	chars, loopTurns := 0, 0
	for _, t := range turns {
		chars += len(t.Text)
		if t.Lane == scan.LaneLoop {
			loopTurns++
		}
	}
	fmt.Printf("%d assistant turns, %d chars (%d in the loop lane)\n\n", len(turns), chars, loopTurns)

	byRule := map[string]int{}
	byLaneRule := map[string]map[string]int{}
	byLane := map[string]int{}
	charsByLane := map[string]int{}
	hitByLane := map[string]int{}
	examples := map[string][]scan.Violation{}
	turnsHit := 0
	// A rolling window, so the cross-turn rule is scored the way it runs: each
	// reply read against the twenty before it, in the order they were written.
	// Without this it is the one rule a backfill cannot report on, which is the
	// one rule that most needs a corpus before it ships.
	var window state.Session
	for _, t := range turns {
		t.Key = "" // every backfilled turn is a distinct entry in the window
		v, phrases := checkAgainst(c, &window, t, minCV)
		window.Record("", nil, len(t.Text), time.Time{}, phrases...)
		byLane[t.Lane]++
		charsByLane[t.Lane] += len(t.Text)
		if len(v) > 0 {
			turnsHit++
			hitByLane[t.Lane]++
		}
		for _, x := range v {
			byRule[x.RuleID]++
			if byLaneRule[t.Lane] == nil {
				byLaneRule[t.Lane] = map[string]int{}
			}
			byLaneRule[t.Lane][x.RuleID]++
			if len(examples[x.RuleID]) < 4 {
				examples[x.RuleID] = append(examples[x.RuleID], x)
			}
		}
	}
	fmt.Printf("turns with at least one hit: %d/%d\n", turnsHit, len(turns))
	for _, lane := range []string{transcript.LaneInteractive, transcript.LaneLoop} {
		if byLane[lane] > 0 {
			// Mean length beside the rate, because reply genre is the standing
			// confound here: most rules need prose to fire at all, so a lane of
			// short replies reads as a lane of clean ones.
			fmt.Printf("  %-12s %d/%d (%.1f%%), mean %d chars\n", lane, hitByLane[lane], byLane[lane],
				100*float64(hitByLane[lane])/float64(byLane[lane]),
				charsByLane[lane]/byLane[lane])
			rules := make([]string, 0, len(byLaneRule[lane]))
			for id := range byLaneRule[lane] {
				rules = append(rules, id)
			}
			sort.Slice(rules, func(i, j int) bool {
				return byLaneRule[lane][rules[i]] > byLaneRule[lane][rules[j]]
			})
			for _, id := range rules {
				fmt.Printf("      %-22s %d\n", id, byLaneRule[lane][id])
			}
		}
	}
	fmt.Println()

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
