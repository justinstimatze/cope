package main

// The MessageDisplay entry: the one place cope can change what the reader sees.
//
// Every other hook in this tool observes. A Stop hook reads the reply after it
// was written and cannot edit it, which is why effigy's `strip` action is
// downgraded to `warn` at card load — there was nothing for it to act on.
// MessageDisplay is different: it hands back `displayContent`, which replaces
// the text on screen while the transcript keeps the original. So the reader gets
// the better version and the scorer still sees what the model actually wrote.
// That split is what makes this safe to run during an experiment: it cannot
// contaminate the measurement it sits next to.
//
// The constraint that shapes everything here is that the hook fires per
// streaming batch. The docs are explicit that batch boundaries follow how the
// text streams rather than any structure in it, so a transform must survive
// being handed a fragment. Two consequences:
//
//   - A pattern that spans a boundary will not match. That is a miss, and a miss
//     is the only acceptable failure: half-applying a rewrite to a fragment
//     would corrupt the line. Every transform here is written so that a partial
//     match passes through untouched.
//   - Anything needing the whole message has to wait for the batch carrying
//     `final`, which is the last chance to change anything. The message text is
//     accumulated per message_id so that a `final`-time transform has the
//     context it needs. Nothing uses it yet; the tail rewrite is what it is for.

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type displayInput struct {
	SessionID string `json:"session_id"`
	MessageID string `json:"message_id"`
	Delta     string `json:"delta"`
	Index     int    `json:"index"`
	Final     bool   `json:"final"`
}

type displayOutput struct {
	HookSpecificOutput struct {
		HookEventName  string `json:"hookEventName"`
		DisplayContent string `json:"displayContent"`
	} `json:"hookSpecificOutput"`
}

// A display transform rewrites text the reader is about to see. Only
// information-preserving rewrites belong here.
//
// Deletion does not. The card's NEVER rules describe moves whose repair needs
// the sentence rebuilt — strip the negation out of "This isn't a tooling
// question, it's a workflow question" and what is left is not a shorter sentence
// but a broken one. Those stay warnings and are the tail rewrite's problem, not
// this table's.
type displayTransform struct {
	id   string
	re   *regexp.Regexp
	with string
	why  string
}

// The table is empty, and that is the current state rather than an oversight.
//
// Its only entry demoted a bold label to a plain one — `**Fixed:** the thing`
// to `Fixed: the thing` — which was the mechanical repair for the card's
// `bold_label` rule. That rule came out on 2026-07-30 after 52 blind pairs put
// bold among the three things the reader says decide a reply for him, so the
// transform was left enforcing a position the card no longer holds. Deleting it
// was the point; nothing has replaced it because no other card rule has a repair
// that does not require the sentence to be understood.
//
// What remains is the tier itself: the batch accumulation, the fail-open paths,
// and the `final` hook the tail rewrite needs. Note that the MessageDisplay
// entry is a `--display` flag nothing currently registers, so this table being
// empty changes no observed behaviour today.
var displayTransforms []displayTransform

// transformDelta applies every transform and reports which ones fired.
func transformDelta(s string) (string, []string) {
	var fired []string
	for _, t := range displayTransforms {
		if !t.re.MatchString(s) {
			continue
		}
		s = t.re.ReplaceAllString(s, t.with)
		fired = append(fired, t.id)
	}
	return s, fired
}

// runDisplay is the MessageDisplay hook entry.
//
// It fails open in every direction. An unreadable payload, an unwritable state
// directory, a transform that changes nothing — all of them end in the original
// text being emitted unchanged. The worst outcome available to this code path is
// that the reader sees exactly what the model wrote, which is the status quo.
func runDisplay(logPath string) {
	var in displayInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return
	}

	out, fired := transformDelta(in.Delta)

	// Accumulate the message so a `final`-time transform has something to read.
	// Written even when nothing fired, because the buffer is context for the
	// rewrite rather than a record of this batch.
	if in.MessageID != "" {
		appendMessageBuffer(in.MessageID, in.Delta, in.Final)
	}

	if len(fired) > 0 && logPath != "" {
		logDisplay(logPath, in, fired)
	}

	var o displayOutput
	o.HookSpecificOutput.HookEventName = "MessageDisplay"
	o.HookSpecificOutput.DisplayContent = out
	_ = json.NewEncoder(os.Stdout).Encode(o)
}

// appendMessageBuffer keeps the text of one message under its id, and drops it
// once the message is complete. A session that ends mid-message leaves a file
// behind; pruneMarkers already sweeps this directory on the refresher path.
func appendMessageBuffer(id, delta string, final bool) {
	dir := stateDir()
	if dir == "" {
		return
	}
	path := filepath.Join(dir, "display-"+filepath.Base(id)+".txt")
	if final {
		_ = os.Remove(path)
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(delta)
}

// logDisplay records what was rewritten, so the rate is measurable rather than
// asserted. Without it there is no way to tell a transform that never fires from
// one that fires constantly — which is the failure this repo spent a day
// documenting in MEASUREMENTS.md.
func logDisplay(path string, in displayInput, fired []string) {
	f, err := os.OpenFile(displayLogPath(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(map[string]any{
		"session": in.SessionID,
		"message": in.MessageID,
		"index":   in.Index,
		"rules":   fired,
		"chars":   len(in.Delta),
	})
}

func displayLogPath(logPath string) string {
	return filepath.Join(filepath.Dir(logPath), "display.jsonl")
}

// runDisplayPreview reads prose on stdin and prints it transformed, so the table
// can be tried on real text without wiring a hook.
func runDisplayPreview() {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
		os.Exit(1)
	}
	out, fired := transformDelta(string(raw))
	fmt.Print(out)
	if len(fired) > 0 {
		fmt.Fprintf(os.Stderr, "\ncope-gate: %s\n", strings.Join(fired, ", "))
	}
}
