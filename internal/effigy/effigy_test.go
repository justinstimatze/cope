package effigy

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

// TestAgreesWithCard2json is the specification. tools/card2json.py runs the
// reference parser, so anywhere the two disagree this package is wrong by
// definition — the notation is effigy's, not cope's.
//
// It runs over every card in the repo rather than a fixture, because the
// fixture that would have caught a bug is the card somebody actually wrote.
func TestAgreesWithCard2json(t *testing.T) {
	root := repoRoot(t)
	effigyRepo := os.Getenv("EFFIGY_PATH")
	if effigyRepo == "" {
		effigyRepo = filepath.Join(filepath.Dir(root), "effigy")
	}
	if _, err := os.Stat(effigyRepo); err != nil {
		t.Skipf("no effigy checkout at %s; the oracle is unavailable", effigyRepo)
	}
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("no python3; the oracle is unavailable")
	}

	cards, err := filepath.Glob(filepath.Join(root, "card", "*.effigy"))
	if err != nil {
		t.Fatal(err)
	}
	demo, err := filepath.Glob(filepath.Join(root, "card", "demo", "*.effigy"))
	if err != nil {
		t.Fatal(err)
	}
	cards = append(cards, demo...)
	if len(cards) < 2 {
		t.Fatalf("found %d cards to compare, expected the shipped card and its demos", len(cards))
	}

	for _, path := range cards {
		t.Run(filepath.Base(path), func(t *testing.T) {
			// card2json records the path it was handed, so hand it the same
			// relative path this package will be given and compare like for
			// like.
			rel, err := filepath.Rel(root, path)
			if err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(t.TempDir(), "rules.json")
			cmd := exec.Command("python3", "tools/card2json.py", rel, out)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "EFFIGY_PATH="+effigyRepo)
			if b, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("card2json.py %s: %v\n%s", rel, err, b)
			}
			want, err := os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}

			src, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			got, err := ToJSON(src, rel)
			if err != nil {
				t.Fatalf("ToJSON: %v", err)
			}

			// Compared as values, not bytes: card2json writes indented JSON
			// with non-ASCII escaped, and neither is a fact about the card.
			var a, b any
			if err := json.Unmarshal(got, &a); err != nil {
				t.Fatalf("unmarshal ours: %v", err)
			}
			if err := json.Unmarshal(want, &b); err != nil {
				t.Fatalf("unmarshal reference: %v", err)
			}
			if !reflect.DeepEqual(a, b) {
				t.Errorf("parsed card differs from card2json.py\n%s", firstDiff(a, b))
			}
		})
	}
}

func TestQuotesComeOffTheExamplesAndNotTheReasoning(t *testing.T) {
	c := mustParse(t, `@id q
WRONG[
  WRONG: "the wrong one"
  RIGHT: "the right one"
  WHY: because "quoted" reasoning has to survive
]`)
	if len(c.Wrong) != 1 {
		t.Fatalf("got %d pairs, want 1", len(c.Wrong))
	}
	w := c.Wrong[0]
	if w.Wrong != "the wrong one" || w.Right != "the right one" {
		t.Errorf("examples kept their quotes: %+v", w)
	}
	if w.Why != `because "quoted" reasoning has to survive` {
		t.Errorf("why was altered: %q", w.Why)
	}
}

func TestAlwaysOnGateBecomesTheEmptyCondition(t *testing.T) {
	c := mustParse(t, `@id g
NEVER[
  @when *
  Never one.
  ---
  @when fact:rule_flip
  Never two.
  ---
  Never three.
]`)
	want := []string{"", "fact:rule_flip", ""}
	if len(c.Never) != len(want) {
		t.Fatalf("got %d rules, want %d", len(c.Never), len(want))
	}
	for i, w := range want {
		if c.Never[i].When != w {
			t.Errorf("rule %d gate %q, want %q", i, c.Never[i].When, w)
		}
	}
}

// A pattern is the one place a card carries brackets that are not block
// delimiters, so the depth count is what keeps a character class from ending
// the block early.
func TestACharacterClassDoesNotCloseTheBlock(t *testing.T) {
	c := mustParse(t, `@id p
POSTPROC[
  action: warn
  pattern: \b(not)\s+[^.,;]{2,40},?\s+(it's)\b
  why: the flip
  id: flip
]
TRAITS[
  after-the-pattern
]`)
	if len(c.Rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(c.Rules))
	}
	if got, want := c.Rules[0].Pattern, `\b(not)\s+[^.,;]{2,40},?\s+(it's)\b`; got != want {
		t.Errorf("pattern truncated:\n got %q\nwant %q", got, want)
	}
	if len(c.Traits) != 1 || c.Traits[0] != "after-the-pattern" {
		t.Errorf("the block after POSTPROC was not reached: %v", c.Traits)
	}
}

func TestAnExchangeKeepsItsLineBreakAndABareLineGetsSpeakerMarked(t *testing.T) {
	c := mustParse(t, `@id m
MES[
{{user}}: did it pass
{{char}}: it did
---
speaking with nobody asking
]`)
	want := []string{"{{user}}: did it pass\n{{char}}: it did", "{{char}}: speaking with nobody asking"}
	if !reflect.DeepEqual(c.Mes, want) {
		t.Errorf("got %q\nwant %q", c.Mes, want)
	}
}

func TestABlockCopeIgnoresDoesNotCostTheCard(t *testing.T) {
	c := mustParse(t, `@id u
GOALS{
  survive: keep going
}
SECRETS[
  something
]
@theme still read
TRAITS[
  reached
]`)
	if c.Theme != "still read" {
		t.Errorf("theme after an unknown block: %q", c.Theme)
	}
	if len(c.Traits) != 1 || c.Traits[0] != "reached" {
		t.Errorf("traits after an unknown block: %v", c.Traits)
	}
}

func TestAWrappedKernelJoinsIntoOneLine(t *testing.T) {
	c := mustParse(t, `@id k
VOICE{
  kernel: First part
    and the continuation.
  peak: The top.
}`)
	if got, want := c.Voice.Kernel, "First part and the continuation."; got != want {
		t.Errorf("kernel %q, want %q", got, want)
	}
	if c.Voice.Peak != "The top." {
		t.Errorf("peak %q", c.Voice.Peak)
	}
}

func TestAnUnterminatedBlockIsAnError(t *testing.T) {
	if _, err := ToJSON([]byte("@id x\nTRAITS[\n  never closed\n"), "x.effigy"); err == nil {
		t.Fatal("parsed a card whose block never closes")
	}
}

func mustParse(t *testing.T, src string) *card {
	t.Helper()
	c, err := parse([]byte(src), "test.effigy")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return c
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(wd, "..", "..")
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// firstDiff names the top-level key that differs, so a failure points at a
// block instead of printing two cards.
func firstDiff(got, want any) string {
	g, ok1 := got.(map[string]any)
	w, ok2 := want.(map[string]any)
	if !ok1 || !ok2 {
		return "the two documents are not both objects"
	}
	for k, wv := range w {
		gv, present := g[k]
		if !present {
			return "missing key: " + k
		}
		if !reflect.DeepEqual(gv, wv) {
			gj, _ := json.Marshal(gv)
			wj, _ := json.Marshal(wv)
			return "key " + k + "\n got " + string(gj) + "\nwant " + string(wj)
		}
	}
	for k := range g {
		if _, present := w[k]; !present {
			return "extra key: " + k
		}
	}
	return "the objects compare unequal but no key differs"
}

// TestGateHeadersAccumulate covers the one header key cope reads that effigy's
// own notation does not define. Repeated keys append rather than overwrite,
// because a card names one exemption per line.
func TestGateHeadersAccumulate(t *testing.T) {
	c := mustParse(t, `@id g
@gate dangling_end off — first
@gate clause_symmetry off — second
@theme still read
`)
	if len(c.Gate) != 2 {
		t.Fatalf("Gate = %v, want two lines", c.Gate)
	}
	if c.Gate[0] != "dangling_end off — first" || c.Gate[1] != "clause_symmetry off — second" {
		t.Errorf("Gate = %#v", c.Gate)
	}
	if c.Theme != "still read" {
		t.Errorf("an unknown header ate a known one: theme = %q", c.Theme)
	}
}

// TestShapeHeadersAccumulate is TestGateHeadersAccumulate's sibling: the other
// header key cope reads that effigy's notation does not define.
func TestShapeHeadersAccumulate(t *testing.T) {
	c := mustParse(t, `@id s
@shape short_landing: last paragraph words <= 12 — lands short
@shape no_walls: every paragraph words <= 90 — no walls of text
`)
	if len(c.Shape) != 2 {
		t.Fatalf("Shape = %v, want two lines", c.Shape)
	}
	if c.Shape[0] != "short_landing: last paragraph words <= 12 — lands short" {
		t.Errorf("Shape[0] = %q", c.Shape[0])
	}
}
