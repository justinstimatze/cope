package scan

import (
	"fmt"
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

// Cluster is one paragraph that drew a concentration of hits, either from
// several rules or from one rule several times.
type Cluster struct {
	// Rules is the distinct rule ids that landed there, sorted.
	Rules []string
	// Counts is how many times each of them did.
	Counts map[string]int
	// Repeated is the rule that landed min times or more on its own, empty
	// when the paragraph qualified by breadth instead.
	Repeated string
	// Excerpt is the paragraph, flattened to one line.
	Excerpt string
}

// Clusters returns the paragraphs where at least min distinct rules landed, or
// where one rule landed at least min times. Widest first.
//
// The second condition is the one a field report asked for and the first one
// does not cover. Running --check over 107 documents produced 114 flip hits of
// which seven were worth changing, and all seven were visible as three in a
// paragraph rather than as anything about the form: "the finding that actually
// changed my prose was 'three in this paragraph', not 'here is a flip'". Three
// hits of one rule is a concentration exactly the way three different rules is,
// and reporting only breadth prints nothing for the case with the measurement
// behind it.
//
// A violation is placed by looking for its matched text in the reply, so a
// match that appears twice is attributed to its first occurrence and one whose
// text was normalised away is not placed at all. Both are acceptable: this
// decides what a report highlights, never what it reports.
func Clusters(text string, v []Violation, min int) []Cluster {
	if min < 2 || len(v) < min {
		return nil
	}
	paras := paraSpans(text)
	hits := placeHits(text, paras, v)

	var out []Cluster
	for i, h := range hits {
		ids, repeated, ok := qualify(h, min)
		if !ok {
			continue
		}
		out = append(out, Cluster{
			Rules:    ids,
			Counts:   h,
			Repeated: repeated,
			Excerpt:  strings.TrimSpace(strings.ReplaceAll(text[paras[i].lo:paras[i].hi], "\n", " ")),
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return len(out[i].Rules) > len(out[j].Rules) })
	return out
}

// span is one paragraph's byte range in the original text.
type span struct{ lo, hi int }

// paraSpans locates each paragraph in the text it came from, so a violation
// found by byte offset can be attributed to one.
func paraSpans(text string) []span {
	var out []span
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
		out = append(out, span{off + i, off + i + len(p)})
		off += i + len(p)
	}
	return out
}

// placeHits returns, per paragraph, how many times each rule landed in it.
func placeHits(text string, paras []span, v []Violation) []map[string]int {
	hits := make([]map[string]int, len(paras))
	for _, x := range v {
		i := paraAt(text, paras, x.Matched)
		if i < 0 {
			continue
		}
		if hits[i] == nil {
			hits[i] = map[string]int{}
		}
		hits[i][x.RuleID]++
	}
	return hits
}

// paraAt is the index of the paragraph containing matched, or -1.
func paraAt(text string, paras []span, matched string) int {
	if matched == "" {
		return -1
	}
	at := strings.Index(text, matched)
	if at < 0 {
		return -1
	}
	for i, p := range paras {
		if at >= p.lo && at < p.hi {
			return i
		}
	}
	return -1
}

// qualify decides whether one paragraph's hits are a cluster, and which kind.
//
// repeated is the rule that landed min times or more on its own, and is empty
// when the paragraph also has min distinct rules. Breadth wins when both hold:
// three different rules says more about a paragraph than three hits of one, and
// naming the three is what tells a reader to rewrite the block rather than hunt
// a construction.
func qualify(h map[string]int, min int) (ids []string, repeated string, ok bool) {
	for id := range h {
		ids = append(ids, id)
	}
	// Sorted before the scan, and the scan takes strictly greater, so two rules
	// tied at the top resolve to the same one on every run. Map order would
	// otherwise make the reported rule differ between two reads of one reply.
	sort.Strings(ids)
	most := 0
	for _, id := range ids {
		if h[id] > most {
			repeated, most = id, h[id]
		}
	}
	if len(ids) < min && most < min {
		return nil, "", false
	}
	if len(ids) >= min {
		return ids, "", true
	}
	return ids, repeated, true
}

// ClusterLine renders the clusters as one line each, or nothing.
//
// Three is the floor rather than two on both conditions. Two rules on a
// paragraph is ordinary — labelled_opening and clause_symmetry read the same
// opening sentence and land together often enough that reporting it would be
// noise — and two hits of one rule is a coincidence a reader can see without
// help.
func ClusterLine(text string, v []Violation) string {
	cs := Clusters(text, v, 3)
	if len(cs) == 0 {
		return ""
	}
	var b strings.Builder
	for _, c := range cs {
		b.WriteString("  " + clusterSentence(c) + "\n      " + trimTo(c.Excerpt, 160) + "\n")
	}
	return b.String()
}

// clusterSentence says which kind of concentration this is, because the two
// call for different edits. Several rules on a paragraph means the paragraph is
// wrong and the sentences are symptoms; one rule several times means the
// construction is a habit in this passage, which is what a per-instance report
// cannot show.
func clusterSentence(c Cluster) string {
	if c.Repeated != "" {
		return fmt.Sprintf("one paragraph drew %s ×%d — the density is the finding, not any one of them:",
			c.Repeated, c.Counts[c.Repeated])
	}
	return "one paragraph drew " + strings.Join(c.Rules, ", ") +
		" — rewrite it rather than the sentences:"
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
