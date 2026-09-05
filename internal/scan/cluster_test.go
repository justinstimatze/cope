package scan

import (
	"strings"
	"testing"
)

func TestClustersFindsTheParagraphThatDrewSeveralRules(t *testing.T) {
	text := "The gate reads the transcript once per turn and writes nothing back.\n\n" +
		"Row loss on cursor reset. The sync drops every row written after the cursor rewinds, " +
		"and the batch commits before the cursor advances.\n\n" +
		"Nothing else changed in this run."
	v := []Violation{
		{RuleID: "labelled_opening", Matched: "Row loss on cursor reset"},
		{RuleID: "clause_symmetry", Matched: "the batch commits before the cursor advances"},
		{RuleID: "flip", Matched: "The sync drops every row"},
		// A fourth rule on a different paragraph, which must not join them.
		{RuleID: "dangling_end", Matched: "Nothing else changed in this run."},
	}
	cs := Clusters(text, v, 3)
	if len(cs) != 1 {
		t.Fatalf("got %d clusters, want 1: %+v", len(cs), cs)
	}
	if got := strings.Join(cs[0].Rules, ","); got != "clause_symmetry,flip,labelled_opening" {
		t.Errorf("rules %q", got)
	}
	if !strings.HasPrefix(cs[0].Excerpt, "Row loss on cursor reset.") {
		t.Errorf("excerpt is the wrong paragraph: %q", cs[0].Excerpt)
	}
}

func TestClustersStaysQuietOnScatteredHits(t *testing.T) {
	text := "One paragraph here with enough words to stand on its own.\n\n" +
		"A second paragraph, also standing on its own and saying something else.\n\n" +
		"A third, closing the reply."
	v := []Violation{
		{RuleID: "labelled_opening", Matched: "One paragraph here"},
		{RuleID: "clause_symmetry", Matched: "A second paragraph"},
		{RuleID: "flip", Matched: "A third, closing"},
	}
	if cs := Clusters(text, v, 3); len(cs) != 0 {
		t.Errorf("clustered three hits on three paragraphs: %+v", cs)
	}
	if ClusterLine(text, v) != "" {
		t.Error("printed a cluster line for scattered hits")
	}
}

// A violation whose matched text is not in the reply — normalised, or quoted
// from somewhere else — is skipped rather than misplaced.
func TestClustersSkipsUnplaceableHits(t *testing.T) {
	text := "Row loss on cursor reset. The sync drops rows after the cursor rewinds."
	v := []Violation{
		{RuleID: "labelled_opening", Matched: "Row loss on cursor reset"},
		{RuleID: "clause_symmetry", Matched: "text that is not in the reply at all"},
		{RuleID: "flip", Matched: ""},
	}
	if cs := Clusters(text, v, 3); len(cs) != 0 {
		t.Errorf("placed a hit it could not find: %+v", cs)
	}
}
