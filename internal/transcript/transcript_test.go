package transcript

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assistantLine(text string) string {
	b, err := json.Marshal(map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"role":    "assistant",
			"content": []map[string]any{{"type": "text", "text": text}},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func toolLine(payload string) string {
	b, err := json.Marshal(map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"role":    "assistant",
			"content": []map[string]any{{"type": "tool_use", "name": "Bash", "input": payload}},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func userLine(text string) string { return promptLine("u-"+text, text) }

func promptLine(uuid, text string) string {
	b, err := json.Marshal(map[string]any{
		"type":    "user",
		"uuid":    uuid,
		"message": map[string]any{"role": "user", "content": text},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

// toolResultLine is a user record that is not a prompt. A turn with six tool
// calls writes six of these, and mistaking them for turn boundaries is what cut
// replies into fragments.
func toolResultLine(payload string) string {
	b, err := json.Marshal(map[string]any{
		"type": "user",
		"uuid": "tr-" + payload,
		"message": map[string]any{
			"role":    "user",
			"content": []map[string]any{{"type": "tool_result", "content": payload}},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

// fireLine is the system record a timer writes when it wakes the session.
func fireLine(uuid string) string {
	b, err := json.Marshal(map[string]any{
		"type":    "system",
		"subtype": "scheduled_task_fire",
		"uuid":    uuid,
		"content": "Claude resuming /loop wakeup (Apr 26 6:04pm)",
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

// promptChild is the user record the wakeup delivers, chained to it by
// parentUuid. Nothing in its text says it came from a timer.
func promptChild(uuid, parent, text string) string {
	b, err := json.Marshal(map[string]any{
		"type":       "user",
		"uuid":       uuid,
		"parentUuid": parent,
		"message":    map[string]any{"role": "user", "content": text},
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func write(t *testing.T, lines []string, trailingNewline bool) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "session.jsonl")
	body := strings.Join(lines, "\n")
	if trailingNewline {
		body += "\n"
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLastAssistantTextSkipsToolAndUserRecords(t *testing.T) {
	path := write(t, []string{
		assistantLine("an earlier turn"),
		userLine("go on"),
		assistantLine("the turn the gate should judge"),
		toolLine("git status"),
		userLine("thanks"),
	}, true)

	got, err := LastAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "the turn the gate should judge" {
		t.Errorf("got %q, want the last assistant prose turn", got)
	}
}

// The backward reader must agree with a plain forward scan on every fixture,
// including when a record straddles a chunk boundary. tailChunk is shrunk to
// force many boundaries over a small file.
func TestLastAssistantTextMatchesAForwardScanAcrossChunkBoundaries(t *testing.T) {
	var lines []string
	for i := range 200 {
		lines = append(lines, assistantLine(fmt.Sprintf("turn %d %s", i, strings.Repeat("padding ", i%40))))
		lines = append(lines, toolLine(strings.Repeat("x", i*17)))
	}
	lines = append(lines, assistantLine("the final prose turn"))
	// A run of non-prose records after it, long enough to span several chunks.
	for i := range 30 {
		lines = append(lines, toolLine(strings.Repeat("y", 500+i*97)))
	}

	for _, trailing := range []bool{true, false} {
		path := write(t, lines, trailing)
		all, err := AllAssistantText(path)
		if err != nil {
			t.Fatal(err)
		}
		want := all[len(all)-1]

		for _, chunk := range []int{16, 64, 512, 4096, 1 << 20} {
			t.Run(fmt.Sprintf("trailing=%v/chunk=%d", trailing, chunk), func(t *testing.T) {
				old := tailChunk
				tailChunk = chunk
				defer func() { tailChunk = old }()

				got, err := LastAssistantText(path)
				if err != nil {
					t.Fatal(err)
				}
				if got != want {
					t.Errorf("got %q, want %q", got, want)
				}
				if got != "the final prose turn" {
					t.Errorf("got %q, want the fixture's final prose turn", got)
				}
			})
		}
	}
}

// Turn assembly accumulates across chunks, so the boundary has to be found and
// the blocks kept in order even when the prompt sits many chunks back.
func TestLastTurnAcrossChunkBoundaries(t *testing.T) {
	lines := []string{promptLine("p-old", "an earlier ask"), assistantLine("an earlier reply")}
	lines = append(lines, promptLine("p-live", "the ask that opens the turn"))
	var want []string
	for i := range 60 {
		// Trimmed, because assistantText trims and the fixture has to agree.
		text := strings.TrimSpace(fmt.Sprintf("block %d %s", i, strings.Repeat("padding ", i%25)))
		lines = append(lines, assistantLine(text), toolLine(strings.Repeat("z", 200+i*31)))
		want = append(want, text)
	}

	for _, chunk := range []int{16, 64, 512, 4096, 1 << 20} {
		t.Run(fmt.Sprintf("chunk=%d", chunk), func(t *testing.T) {
			old := tailChunk
			tailChunk = chunk
			defer func() { tailChunk = old }()

			got, err := LastTurn(write(t, lines, true))
			if err != nil {
				t.Fatal(err)
			}
			if got.Key != "p-live" {
				t.Errorf("key %q, want p-live", got.Key)
			}
			if got.Text != strings.Join(want, "\n\n") {
				t.Errorf("turn text differs at chunk %d", chunk)
			}
		})
	}
}

func TestLastAssistantTextOnEmptyAndProseFreeFiles(t *testing.T) {
	empty := write(t, nil, false)
	got, err := LastAssistantText(empty)
	if err != nil || got != "" {
		t.Errorf("empty file: got %q, %v — want empty and no error", got, err)
	}

	noProse := write(t, []string{toolLine("ls"), userLine("hi")}, true)
	got, err = LastAssistantText(noProse)
	if err != nil || got != "" {
		t.Errorf("no assistant prose: got %q, %v", got, err)
	}
}

// A truncated write or a half-flushed line is normal in a live transcript. It
// must not blind the gate to the turn above it.
func TestMalformedLinesAreSkippedNotFatal(t *testing.T) {
	path := write(t, []string{
		assistantLine("the good turn"),
		`{"type":"assistant","message":{"role":"assistant","content":[{"type":"te`,
		"",
		"not json at all",
	}, false)

	got, err := LastAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "the good turn" {
		t.Errorf("got %q, want the last well-formed turn", got)
	}
}

func TestAllAssistantTextReturnsEveryTurnOldestFirst(t *testing.T) {
	path := write(t, []string{
		assistantLine("first"),
		toolLine("noise"),
		assistantLine("second"),
		userLine("noise"),
		assistantLine("third"),
	}, true)

	got, err := AllAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"first", "second", "third"}
	if len(got) != len(want) {
		t.Fatalf("got %d turns, want %d: %q", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("turn %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

// One reply is written as several assistant records, and the gate has to see
// all of them. Reading the last record alone gets a mid-turn aside — on this
// repo's own session that meant missing 41% of the labelled openings and 32% of
// the bold labels, while inventing paragraph uniformity in the trailing scrap.
func TestLastTurnAssemblesEveryBlockSinceThePrompt(t *testing.T) {
	path := write(t, []string{
		promptLine("p-1", "an earlier ask"),
		assistantLine("an earlier reply"),
		promptLine("p-2", "fix the turn assembly"),
		assistantLine("Reading the parser first."),
		toolLine("rg LastTurn"),
		toolResultLine("three matches"),
		assistantLine("Found it."),
		toolLine("go test ./..."),
		toolResultLine("ok"),
		assistantLine("Tests pass. Want me to push?"),
	}, true)

	got, err := LastTurn(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "Reading the parser first.\n\nFound it.\n\nTests pass. Want me to push?"
	if got.Text != want {
		t.Errorf("got %q\nwant %q", got.Text, want)
	}
	if got.Key != "p-2" {
		t.Errorf("key is %q, want the id of the prompt that opened the turn", got.Key)
	}
}

// The catch-up pass rescores a turn the Stop hook saw without its ending, and
// reports only the difference. That is sound only if the earlier read is a byte
// prefix of the later one, so the property is pinned rather than assumed.
func TestATurnGrowsByAppendOnly(t *testing.T) {
	head := []string{
		promptLine("p-1", "go on"),
		assistantLine("Starting."),
		toolLine("ls"),
		assistantLine("One file."),
	}
	partial, err := LastTurn(write(t, head, true))
	if err != nil {
		t.Fatal(err)
	}
	full, err := LastTurn(write(t, append(head, assistantLine("Done. Push it?")), true))
	if err != nil {
		t.Fatal(err)
	}

	if partial.Key != full.Key {
		t.Errorf("key changed as the turn grew: %q then %q", partial.Key, full.Key)
	}
	if !strings.HasPrefix(full.Text, partial.Text) {
		t.Errorf("the finished turn does not start with what was read early:\n  early: %q\n  full:  %q", partial.Text, full.Text)
	}
	if len(full.Text) <= len(partial.Text) {
		t.Error("the finished turn should be longer than the early read")
	}
}

// Whether the prompt being submitted is already in the file when
// UserPromptSubmit fires is not something a hook gets to know. Under the order
// where it is, stopping at the first boundary would return nothing and the
// catch-up would quietly never run, so an unanswered prompt is walked past.
func TestLastTurnWalksPastAPromptWithNoReplyYet(t *testing.T) {
	lines := []string{
		promptLine("p-1", "the ask that opened the reply"),
		assistantLine("Starting."),
		toolLine("ls"),
		assistantLine("Done. Push it?"),
	}
	answered, err := LastTurn(write(t, lines, true))
	if err != nil {
		t.Fatal(err)
	}
	pending, err := LastTurn(write(t, append(lines, promptLine("p-2", "yes push")), true))
	if err != nil {
		t.Fatal(err)
	}
	if pending.Text != answered.Text || pending.Key != answered.Key {
		t.Errorf("a new prompt changed the answer:\n  before: %+v\n  after:  %+v", answered, pending)
	}
	if pending.Key != "p-1" {
		t.Errorf("key is %q, want the prompt that opened the reply", pending.Key)
	}
}

// A transcript that opens mid-turn has no boundary to find. It should still
// return the prose rather than nothing, with an empty key saying so.
func TestLastTurnWithNoPromptAbove(t *testing.T) {
	path := write(t, []string{assistantLine("orphaned prose"), toolLine("ls")}, true)
	got, err := LastTurn(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "orphaned prose" || got.Key != "" {
		t.Errorf("got %+v, want the prose with no key", got)
	}
}

func TestAllTurnsGroupsRecordsIntoTurns(t *testing.T) {
	path := write(t, []string{
		promptLine("p-1", "first ask"),
		assistantLine("looking"),
		toolLine("ls"),
		toolResultLine("a b c"),
		assistantLine("three files"),
		promptLine("p-2", "second ask"),
		assistantLine("done"),
	}, true)

	turns, err := AllTurns(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"looking\n\nthree files", "done"}
	if len(turns) != len(want) {
		t.Fatalf("got %d turns, want %d: %q", len(turns), len(want), turns)
	}
	for i := range want {
		if turns[i].Text != want[i] {
			t.Errorf("turn %d: got %q, want %q", i, turns[i].Text, want[i])
		}
		if turns[i].Lane != LaneInteractive {
			t.Errorf("turn %d: lane %q, want %q", i, turns[i].Lane, LaneInteractive)
		}
	}
	if turns[0].Key != "p-1" || turns[1].Key != "p-2" {
		t.Errorf("keys %q/%q, want the prompt that opened each turn", turns[0].Key, turns[1].Key)
	}

	// The record-level reader counts the same file as three, which is the
	// denominator --backfill used to report and the reason its per-turn rates
	// read low.
	recs, err := AllAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 3 {
		t.Errorf("record reader found %d records, want 3", len(recs))
	}
}

// Multiple text blocks in one message are joined, which is what the model
// produces when it interleaves prose with tool calls.
func TestMultipleTextBlocksAreJoined(t *testing.T) {
	b, _ := json.Marshal(map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"role": "assistant",
			"content": []map[string]any{
				{"type": "text", "text": "part one"},
				{"type": "tool_use", "name": "Read"},
				{"type": "text", "text": "part two"},
			},
		},
	})
	path := write(t, []string{string(b)}, true)

	got, err := LastAssistantText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "part one\npart two" {
		t.Errorf("got %q, want both text blocks joined", got)
	}
}

// syntheticLine is an assistant record Claude Code composed itself. The two
// markers are set by different code paths, so each is exercised alone.
func syntheticLine(text string, marker string) string {
	msg := map[string]any{
		"role":    "assistant",
		"content": []map[string]any{{"type": "text", "text": text}},
	}
	rec := map[string]any{"type": "assistant", "uuid": "syn-" + text, "message": msg}
	switch marker {
	case "model":
		msg["model"] = "<synthetic>"
	case "flag":
		rec["isApiErrorMessage"] = true
	}
	b, err := json.Marshal(rec)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// An API error banner sits in the reply position and reads as prose the model
// wrote. It is not, and scoring it charged the model for a sentence the harness
// composed — once per retry, identically, which is what a repetition rule is
// built to notice.
func TestSyntheticAssistantRecordsAreNotTurns(t *testing.T) {
	const banner = "API Error: 500 Internal server error. This is a server-side " +
		"issue, usually temporary — try again in a moment."
	path := write(t, []string{
		userLine("go"),
		syntheticLine(banner, "model"),
		syntheticLine(banner, "flag"),
		assistantLine("the reply the model actually wrote"),
	}, true)

	turn, err := LastTurn(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(turn.Text, "API Error") {
		t.Errorf("a harness-composed record was scored as the model's prose:\n%s", turn.Text)
	}
	if turn.Text != "the reply the model actually wrote" {
		t.Errorf("got %q", turn.Text)
	}

	all, err := AllTurns(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || strings.Contains(all[0].Text, "API Error") {
		t.Errorf("backfill saw %d turn(s): %q", len(all), all)
	}
}

// The lane comes off the prompt that opened the turn, and the three shapes a
// loop-lane prompt arrives in are the slash command, the scheduled fire, and
// the sentinel the dynamic variant sends itself.
func TestLaneReadsThePromptThatOpenedTheTurn(t *testing.T) {
	for name, prompt := range map[string]string{
		"slash loop": "<command-message>loop</command-message>\n" +
			"<command-name>/loop</command-name>\n<command-args>check CI</command-args>",
		"slash goal": "<command-name>/goal</command-name>\n<command-args>tests pass</command-args>",
		"sentinel":   "<<autonomous-loop-dynamic>>",
	} {
		path := write(t, []string{promptLine("p", prompt), assistantLine("did the thing")}, true)
		turn, err := LastTurn(path)
		if err != nil {
			t.Fatal(err)
		}
		if turn.Lane != LaneLoop {
			t.Errorf("%s: lane %q, want %q", name, turn.Lane, LaneLoop)
		}
	}

	// The continuing wakeups are system records with the prompt chained under
	// them, which is the shape that made the loop lane read as empty when this
	// matched on prose instead.
	path := write(t, []string{
		fireLine("fire-1"),
		promptChild("p", "fire-1", "check on the run"),
		assistantLine("checked"),
	}, true)
	turn, err := LastTurn(path)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Lane != LaneLoop {
		t.Errorf("scheduled wakeup: lane %q, want %q", turn.Lane, LaneLoop)
	}

	// A prompt that merely talks about loops is a person typing.
	for name, prompt := range map[string]string{
		"prose about loops": "can you explain what /loop does and when it fires",
		"ordinary ask":      "rebase onto main and rerun the suite",
	} {
		path := write(t, []string{promptLine("p", prompt), assistantLine("did the thing")}, true)
		turn, err := LastTurn(path)
		if err != nil {
			t.Fatal(err)
		}
		if turn.Lane != LaneInteractive {
			t.Errorf("%s: lane %q, want %q", name, turn.Lane, LaneInteractive)
		}
	}
}

// Forward and backward readers must agree, and AllTurns has to hold the lane
// across the reply the prompt opened rather than attaching it to the next one.
func TestAllTurnsCarriesLanePerTurn(t *testing.T) {
	path := write(t, []string{
		promptLine("p-1", "rebase onto main"),
		assistantLine("rebased"),
		fireLine("fire-1"),
		promptChild("p-2", "fire-1", "check on the run"),
		assistantLine("checked CI, nothing new"),
		promptLine("p-3", "thanks"),
		assistantLine("noted"),
	}, true)

	turns, err := AllTurns(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{LaneInteractive, LaneLoop, LaneInteractive}
	if len(turns) != len(want) {
		t.Fatalf("got %d turns, want %d", len(turns), len(want))
	}
	for i, w := range want {
		if turns[i].Lane != w {
			t.Errorf("turn %d (%q): lane %q, want %q", i, turns[i].Text, turns[i].Lane, w)
		}
	}

	last, err := LastTurn(path)
	if err != nil {
		t.Fatal(err)
	}
	if last.Lane != turns[len(turns)-1].Lane {
		t.Errorf("backward reader says %q, forward says %q", last.Lane, turns[len(turns)-1].Lane)
	}
}
