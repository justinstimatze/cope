package scan

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

// The one complaint in the field roundup that a per-reply check cannot see.
// Zvi Mowshowitz's July 2026 Opus 5 collection puts it plainly — "can we give
// it a rest with the text walls with the same turns of phrase over and over" —
// and every rule above this line reads one reply in isolation, where a phrase
// used for the fourth time looks exactly like a phrase used for the first.
//
// Reading across turns is what cope has that a linter does not: internal/state
// already carries a window of what the session did. This adds the phrases to it.
//
// Fingerprints, never prose. state.Session is written to disk and holds rule
// names, counts and character totals precisely because none of those reconstruct
// the conversation. A 32-bit sum of five words keeps that property: it is enough
// to notice a recurrence and not enough to read back what recurred. The quoted
// text in a violation comes from the turn in hand, which the gate is already
// holding.

// shingleLen is how many consecutive words make a phrase.
//
// Four is too short to be a turn of phrase — "is not the same" is a fragment
// that lands inside a hundred different sentences. Six almost never repeats
// verbatim, so the rule would report nothing and read as working. Five is where
// a construction is long enough to be recognisably the same one and short
// enough to actually recur.
const shingleLen = 5

// scaffoldFloor is how many of a shingle's words must be common ones.
//
// All five would restrict this to pure connective tissue ("that is the whole
// point") and miss the family the complaint is actually about, where a
// construction carries one content word through: "the mechanism, not the
// noise", "X is the whole game". Four admits exactly one content word.
//
// Below four the rule stops reading turns of phrase and starts reading subject
// matter: three content words in a five-word span is a topic, and a
// conversation about one topic repeats its topic by design. That is the
// distinction this constant buys, and it is the whole reason the rule can run
// on a focused session without firing on every paragraph.
const scaffoldFloor = 4

// maxPhrases bounds what one turn contributes to the stored window. A long
// reply can produce several hundred shingles, and twenty of those turns is a
// session file measured in megabytes. Truncation costs recall on the tail of a
// long turn and nothing else.
const maxPhrases = 400

// Phrase is one candidate turn of phrase: the fingerprint that gets stored and
// the words as they appeared, which do not.
type Phrase struct {
	Sum  uint32
	Text string
}

// Phrases returns every scaffolding shingle in text, one entry per distinct
// fingerprint. A phrase repeated inside a single reply appears once: the rule
// this feeds counts turns, and a turn that says something twice is still one
// turn saying it.
func Phrases(text string) []Phrase {
	return phrases(paragraphs(Redact(text)))
}

func phrases(ps []string) []Phrase {
	var out []Phrase
	seen := map[uint32]bool{}
	for _, p := range ps {
		// Per sentence, so a shingle cannot straddle a full stop. A span
		// running from the end of one sentence into the start of the next is
		// not a construction anybody wrote; it is an artefact of the window.
		for _, sent := range splitSentences(p) {
			ws := words(sent)
			for i := 0; i+shingleLen <= len(ws); i++ {
				span := ws[i : i+shingleLen]
				if scaffoldCount(span) < scaffoldFloor {
					continue
				}
				sum, text := fingerprint(span)
				if seen[sum] {
					continue
				}
				seen[sum] = true
				out = append(out, Phrase{Sum: sum, Text: text})
				if len(out) >= maxPhrases {
					return out
				}
			}
		}
	}
	return out
}

func fingerprint(span []string) (uint32, string) {
	lower := make([]string, len(span))
	for i, w := range span {
		lower[i] = strings.ToLower(w)
	}
	joined := strings.Join(lower, " ")
	h := fnv.New32a()
	_, _ = h.Write([]byte(joined))
	return h.Sum32(), joined
}

func scaffoldCount(span []string) int {
	n := 0
	for _, w := range span {
		if scaffold[strings.ToLower(w)] {
			n++
		}
	}
	return n
}

// Repeated is one phrase this turn shares with earlier ones.
type Repeated struct {
	Text  string
	Turns int // earlier turns that used it, not counting this one
}

// repeatTurns is how many earlier turns must already carry a phrase before it
// counts. Two, so a violation means the reply is the third airing — which is
// what "over and over" describes, and what keeps a phrase used twice in a long
// session from being called a habit.
const repeatTurns = 2

// Sums reduces phrases to what gets stored.
func Sums(ps []Phrase) []uint32 {
	out := make([]uint32, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Sum)
	}
	return out
}

// CrossTurnRepeat flags a construction this reply shares with several earlier
// ones. prior maps a fingerprint to how many earlier turns used it.
//
// It takes the phrases rather than the text because the caller needs them
// anyway — the same slice is what gets recorded for the next turn to read.
//
// One violation per reply, naming the phrase repeated across the most turns.
// The failure is a session that keeps reaching for the same shapes, not any
// single span, and listing six of them would be the flat injection again.
func CrossTurnRepeat(ps []Phrase, prior map[uint32]int) []Violation {
	var hits []Repeated
	for _, p := range ps {
		if n := prior[p.Sum]; n >= repeatTurns {
			hits = append(hits, Repeated{Text: p.Text, Turns: n})
		}
	}
	if len(hits) == 0 {
		return nil
	}
	// Busiest first, ties on the text so the report is stable between runs
	// rather than following map order.
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Turns != hits[j].Turns {
			return hits[i].Turns > hits[j].Turns
		}
		return hits[i].Text < hits[j].Text
	})
	top := hits[0]
	ctx := ""
	if len(hits) > 1 {
		ctx = "also repeated: " + hits[1].Text
	}
	return []Violation{{
		RuleID:  "cross_turn_repeat",
		Action:  "warn",
		Why:     fmt.Sprintf("a turn of phrase this session has already used in %d earlier turns", top.Turns),
		Matched: top.Text,
		Context: ctx,
	}}
}

// scaffold is the connective tissue of English prose: the words a construction
// is built out of rather than the words it is about.
//
// It is a frequency list, not a stoplist. Nothing here is banned — membership
// only decides whether a five-word span is read as a shape or as a subject. The
// working definition is a word that carries no topic on its own, so a span made
// of these plus at most one other word is a construction the writer reached
// for rather than a thing the conversation is about.
var scaffold = boolSet(strings.Fields(`
a about above across after again against all almost along already also although
always am among an and another any anything anyway are around as at away back
be because been before being below best better between both but by came can
cannot come comes could couldn't did didn't do does doesn't doing done don't
down during each either else enough even ever every everything except far few
first for from further get gets getting give given go goes going gone got had
hadn't half has hasn't have haven't having he her here hers herself him himself
his how however i if in indeed instead into is isn't it it's its itself just
keep kept last least left less let like likely little long look looks made make
makes making many may maybe me mean means meant merely might mine more most
much must my myself near nearly need needs neither never new next no nobody
none nor not nothing now of off often on once one only onto or other others
otherwise ought our ours out over own part particular past perhaps place point
possible put quite rather real really right said same say says see seem seems
seen sense set several shall she should shouldn't simply since small so some
something sometimes soon sort still such sure take taken takes than that that's
the their theirs them themselves then there there's these they thing things
think this those though thought three through thus time to together too took
toward true try trying turn two under until up upon us use used using very want
wants was wasn't way we well went were weren't what what's whatever when where
whether which while who whole whom whose why will with within without won't
work worse worst would wouldn't yes yet you your yours yourself
`))

func boolSet(ws []string) map[string]bool {
	m := make(map[string]bool, len(ws))
	for _, w := range ws {
		m[w] = true
	}
	return m
}
