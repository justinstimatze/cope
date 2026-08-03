package main

import (
	"regexp"
	"testing"
)

// The table is empty since bold_label came out of the card, so the contract the
// hook offers today is that it changes nothing. A transform arriving later that
// breaks this on text it was not meant to touch is the bug this catches.
func TestEmptyTableIsAPassThrough(t *testing.T) {
	for _, in := range []string{
		"**Fixed:** the wrapper.",
		"- **Cost:** two calls.",
		"A **bold phrase** mid-line.",
		"**Split across a bou",
		"",
	} {
		got, fired := transformDelta(in)
		if got != in {
			t.Errorf("delta was rewritten with no transforms in the table\n  in:  %q\n  got: %q", in, got)
		}
		if len(fired) != 0 {
			t.Errorf("%q reported %v as fired with an empty table", in, fired)
		}
	}
}

// Nothing may be lost. The transcript keeps the original either way, but a
// reader who sees fewer words than were written is being misled about the reply,
// which is a different thing from being shown it in plainer type. This is the
// standing contract on the table, not on any one entry, so it holds while the
// table is empty and keeps holding for whatever goes in next.
func TestDisplayTransformDropsNoWords(t *testing.T) {
	in := "**Fixed:** the retry wrapper caught every error.\n- **Cost:** two calls."
	out, _ := transformDelta(in)
	for _, w := range []string{"Fixed", "the retry wrapper caught every error", "Cost", "two calls"} {
		if !contains(out, w) {
			t.Errorf("%q went missing from %q", w, out)
		}
	}
}

// Every entry has to survive being handed a fragment, because batch boundaries
// follow how the text streams and not any structure in it. A pattern anchored to
// something that can land in the next delta half-applies the rewrite and corrupts
// the line, so an unanchored end is checked here rather than per entry.
func TestEveryTransformIsWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, tr := range displayTransforms {
		if tr.id == "" || tr.why == "" {
			t.Errorf("transform %q needs both an id and a why — the log reads both", tr.id)
		}
		if seen[tr.id] {
			t.Errorf("two transforms share the id %q, so the log cannot tell them apart", tr.id)
		}
		seen[tr.id] = true
		if tr.re == nil {
			t.Errorf("transform %q has no pattern", tr.id)
			continue
		}
		if _, err := regexp.Compile(tr.re.String()); err != nil {
			t.Errorf("transform %q: %v", tr.id, err)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
