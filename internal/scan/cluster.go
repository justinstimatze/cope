package scan

import (
	"sort"
	"strings"
)

// Where several rules landed on the same paragraph.
//
// Every rule here fires on its own and knows nothing about the others, so a
// report is a flat list and a reader works down it hit by hit. That is the
// wrong shape for the common case. Three hits spread across three paragraphs
// are three small edits; three hits inside one paragraph are one paragraph to
// write again, and the flat list cannot tell them apart.
//
// The idea is borrowed rather than invented. Wikipedia's "Signs of AI writing"
// is explicit that a single tell proves nothing and that several together are
// the evidence — one em dash is punctuation, one em dash beside a triad beside
// a sales adjective is a passage. Nothing here changes what fires or how a hit
// is scored; the clustering is a second reading of the same violations, so a
// rule that would have been reported alone still is.

// Cluster is one paragraph that drew hits from several rules.
type Cluster struct {
	// Rules is the distinct rule ids that landed there, sorted.
	Rules []string
	// Excerpt is the paragraph, flattened to one line.
	Excerpt string
}

// Clusters returns the paragraphs where at least min distinct rules landed,
// most-hit first.
//
// A violation is placed by looking for its matched text in the reply, so a
// match that appears twice is attributed to its first occurrence and one whose
// text was normalised away is not placed at all. Both are acceptable: this
// decides what a report highlights, never what it reports.
func Clusters(text string, v []Violation, min int) []Cluster {
	if min < 2 || len(v) < min {
		return nil
	}
	type para struct {
		lo, hi int
	}
	var paras []para
	off := 0
	for _, p := range paraSplit.Split(text, -1) {
		// An empty block between two separators would match at offset zero and
		// swallow every hit into a paragraph with no text in it.
		if strings.TrimSpace(p) == "" {
			continue
		}
		i := strings.Index(text[off:], p)
		if i < 0 {
			continue
		}
		paras = append(paras, para{off + i, off + i + len(p)})
		off += i + len(p)
	}

	hits := make([]map[string]bool, len(paras))
	for _, x := range v {
		if x.Matched == "" {
			continue
		}
		at := strings.Index(text, x.Matched)
		if at < 0 {
			continue
		}
		for i, p := range paras {
			if at >= p.lo && at < p.hi {
				if hits[i] == nil {
					hits[i] = map[string]bool{}
				}
				hits[i][x.RuleID] = true
				break
			}
		}
	}

	var out []Cluster
	for i, h := range hits {
		if len(h) < min {
			continue
		}
		ids := make([]string, 0, len(h))
		for id := range h {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		out = append(out, Cluster{
			Rules:   ids,
			Excerpt: strings.TrimSpace(strings.ReplaceAll(text[paras[i].lo:paras[i].hi], "\n", " ")),
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return len(out[i].Rules) > len(out[j].Rules) })
	return out
}

// ClusterLine renders the clusters as one line each, or nothing.
//
// Three is the floor rather than two because two rules on a paragraph is
// ordinary: labelled_opening and clause_symmetry both read the same opening
// sentence and land together often enough that reporting it would be noise.
func ClusterLine(text string, v []Violation) string {
	cs := Clusters(text, v, 3)
	if len(cs) == 0 {
		return ""
	}
	var b strings.Builder
	for _, c := range cs {
		b.WriteString("  one paragraph drew " + strings.Join(c.Rules, ", ") +
			" — rewrite it rather than the sentences:\n      " + trimTo(c.Excerpt, 160) + "\n")
	}
	return b.String()
}

func trimTo(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && (s[n]&0xC0) == 0x80 {
		n--
	}
	return s[:n] + "…"
}
