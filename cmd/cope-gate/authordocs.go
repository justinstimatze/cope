package main

// --author-docs emits a prompt for writing this repo's prose, following the
// pattern vidette states in its README: the tool is the structured input, the
// Claude session already running is the runtime. cope does not become an LLM
// client — it has no dependencies and adding an API client to render a README
// would be a poor trade.
//
// effigy does this differently: generate_readme.py calls the Anthropic API and
// writes README.md itself. That works, and it needs a key, an SDK, and a model
// choice baked into the repo. The half worth copying is the rule in its
// docstring — "The generator provides facts, not prose." Everything below is
// introspected. No sentence here ends up in the README.
//
// The step neither of them has: the same binary scores the result. --check
// reads the finished file back through the rules that gate every reply, so the
// docs are held to the card they describe.

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/justinstimatze/cope/card"
	"github.com/justinstimatze/cope/internal/scan"
)

// docsFacts is everything the prompt states as given. Every field is read from
// the card or the flag set at run time, so a fact cannot drift from the thing
// it describes the way a hand-written flag table does.
type docsFacts struct {
	Module       string      `json:"module"`
	NameOrigin   string      `json:"name_origin"`
	CardSource   string      `json:"card_source"`
	CardID       string      `json:"card_id"`
	Composition  composition `json:"card_composition"`
	RenderedCard sizes       `json:"rendered_sizes_chars"`
	FocusNote    string      `json:"refresher_size_note"`
	Flags        []flagFact  `json:"flags"`
	Axes         []axisFact  `json:"axes"`
	Lanes        []laneFact  `json:"lanes"`
	RegexRules   []ruleFact  `json:"regex_rules"`

	// ShippedRegexRules is the same list for the card cope ships. Both are
	// here for the reason both counts are: a page written under a demo card
	// was documenting that card's POSTPROC rules as cope's, so the front page
	// listed two rules nobody installs and omitted the three everybody gets.
	ShippedRegexRules []ruleFact `json:"regex_rules_in_shipped_card"`

	ShapeRules    []shapeFact      `json:"shape_rules"`
	Hooks         []hookFact       `json:"hooks"`
	SettingsJSON  string           `json:"settings_json_block"`
	StyleInstall  styleInstallFact `json:"output_style_install"`
	StateFiles    []stateFact      `json:"state_files"`
	NeverBudget   int              `json:"never_rule_budget"`
	NeverOverflow []string         `json:"never_rules_over_budget"`
	Gating        []gateFact       `json:"gated_card_items"`
	CardRules     authorFact       `json:"card_authored_rules"`
	TestNames     []string         `json:"card_test_names"`
	Layout        []layoutFact     `json:"layout"`
	Related       []relatedFact    `json:"related_projects"`
}

// gateFact describes a condition the card uses to narrow an injection, so the
// README can explain why the mid-session text differs from the SessionStart
// text without anyone hand-maintaining the list.
type gateFact struct {
	Condition string `json:"condition"`
	Items     int    `json:"card_items_gated_on_it"`
}

// authorFact describes how a card reaches the structure half, which is
// otherwise compiled in. Held as facts rather than prose because both syntaxes
// and the loaded card's own use of them come from code, and a README that got a
// syntax wrong would be worse than one that omitted it.
type authorFact struct {
	GateSyntax   string   `json:"gate_syntax"`
	GateExample  string   `json:"gate_example_card"`
	ShapeSyntax  string   `json:"shape_syntax"`
	ShapeWords   []string `json:"shape_vocabulary"`
	ShapeExample string   `json:"shape_example_card"`
	ThisCard     []string `json:"this_cards_own_rules"`
	Limit        string   `json:"what_it_still_cannot_do"`
}

type relatedFact struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	How  string `json:"how_it_relates"`
}

type composition struct {
	Never  int `json:"never"`
	Wrong  int `json:"wrong_pairs"`
	Mes    int `json:"whole_reply_exemplars"`
	Tests  int `json:"tests"`
	Quirks int `json:"quirks"`
	Traits int `json:"traits"`
	Regex  int `json:"postproc_rules"`

	// ShippedRegex is the same count for the card cope ships, which is a
	// different number whenever a demo card is loaded. Both are here because
	// the two get confused: a page written under a demo card was stating that
	// card's POSTPROC count and attributing it to the shipped one, so every
	// demo page claimed cope ships zero regex rules and the front page —
	// rendered from the maximal card — claimed one. It ships three.
	ShippedRegex int `json:"postproc_rules_in_shipped_card"`
}

type sizes struct {
	Card      int `json:"injected_card"`
	Refresher int `json:"refresher"`
}

type flagFact struct {
	Name    string `json:"flag"`
	Default string `json:"default"`
	Usage   string `json:"does"`
}

type ruleFact struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Why    string `json:"why"`
}

type hookFact struct {
	Event   string `json:"event"`
	Command string `json:"command"`
	Does    string `json:"does"`
}

// settingsJSON is the hooks block worth wiring, given to the prompt as a literal
// so the README shows a working configuration rather than one inferred from
// prose. A reader who has not written a Claude Code hook cannot guess the
// matcher, the type field, or that Stop wants async.
//
// SessionStart --inject is deliberately not in it. That hook was how the card
// used to arrive and it is the weakest slot available; the output style replaced
// it. What is left are the three jobs a static style cannot do: Stop reads the
// reply back and scores it, UserPromptSubmit restates the rules that have
// actually been firing this session, and PreToolUse scores the prose an
// external write is about to post — the one surface that never reaches the
// transcript at all.
//
// The commands are bare, which assumes the go install target is on PATH in the
// environment the hook runs in. It often is not, and the note under the block
// says so.
const settingsJSON = `{
  "hooks": {
    "Stop": [
      {
        "matcher": "",
        "hooks": [
          { "type": "command", "command": "cope-gate", "timeout": 10, "async": true }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [
          { "type": "command", "command": "cope-gate --refresher" }
        ]
      }
    ],
    "PreToolUse": [
      {
        "matcher": "mcp__linear__save_comment|mcp__linear__save_document|mcp__linear__save_issue|mcp__linear__save_project|mcp__linear__save_status_update",
        "hooks": [
          { "type": "command", "command": "cope-gate --pretool", "timeout": 10 }
        ]
      }
    ]
  }
}`

// styleInstallFact is the delivery that works, held as facts so the README
// cannot invent a command or a menu path for the one section where being wrong
// costs a reader the whole tool.
type styleInstallFact struct {
	// Block is the copyable install, given as a literal for the same reason
	// settingsJSON is: prose instructions to name the menu entry landed on
	// three renders out of seven and missed the front page, and a reader who
	// cannot find the entry has a style selected and nothing changed. A fact
	// marked reproduce-verbatim is the only thing in this prompt that survives
	// every render.
	Block string `json:"copyable_install_block"`

	Setup  string `json:"setup"`
	Emit   string `json:"emit_by_hand"`
	Select string `json:"select"`
	Timing string `json:"when_it_applies"`
	Why    string `json:"why_here_rather_than_a_hook"`
}

type stateFact struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	What string `json:"holds"`
}

type layoutFact struct {
	Path string `json:"path"`
	What string `json:"what"`
}

// shapeRules are the rules built into the binary and what each one actually
// matches. They carry their own descriptions because the first generated draft
// inferred them from the names and got one wrong: it said labelled_opening
// catches "a paragraph or list item", when paragraphs() skips list blocks
// outright. A name is not a specification. Keep these sentences honest against
// internal/scan — and against the card, since a description naming a POSTPROC
// rule outlives the rule. This block claimed the bolded label was covered by
// bold_label for four days after the card dropped bold_label.
//
// Axis is which of the two things the rule is about — see axes below. It is a
// judgement about what each rule is sensitive to and not a fact the card
// records, the same judgement ruleClass makes in abtest.go and at the same
// grain: that file's `rhetoric` and `deference` classes are voicing here, and
// `formatting`, `ending` and `handoff` are structure.
//
// Lane is empty for a rule that runs in both. internal/scan/loop.go drops the
// four interactive ending rules for a turn nobody is present to answer and adds
// the two loop ones in their place.
var shapeRules = []shapeFact{
	{"labelled_opening", "structure", "", "a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced"},
	{"clause_symmetry", "voicing", "", "comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint — the balanced two-beat"},
	{"paragraph_uniformity", "structure", "", "four or more prose paragraphs whose lengths have a coefficient of variation below --min-cv"},
	{"ask_not_last", "structure", "interactive", "a question or request for the reader sitting in an earlier block while the reply carries on past it"},
	{"dangling_end", "structure", "interactive", "an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving 'continue' nothing to refer to"},
	{"buried_decision", "structure", "interactive", "an open problem landing after the last question or offer, burying the decision point above it"},
	{"forked_end", "structure", "interactive", "two or more things to act on in the closing blocks with nothing marking which comes first, so answering 'continue' means picking one. Sentences opening on 'or', questions inside list items and table cells, and bare deference tags like 'your call' are read as continuing the decision above rather than adding another"},
	{"apology", "voicing", "", "the reply performs contrition instead of stating the correction and moving on"},
	{"self_postmortem", "voicing", "", "the reply turns to account for its own errors, which is a story the reader did not ask for"},
	{"announced_length", "voicing", "", "the reply announces its own length rather than cutting it"},
	{"cross_turn_repeat", "voicing", "", "a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history"},
	{"unverified_done", "structure", "loop", "says the work is done with nothing on the page that could have shown it — no command, no count, no file"},
	{"loop_ask", "structure", "loop", "an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself"},
}

type shapeFact struct {
	ID      string `json:"id"`
	Axis    string `json:"axis"`
	Lane    string `json:"lane,omitempty"`
	Catches string `json:"catches"`
}

// axisFact states what the two axes are and, more to the point, where each one
// lives. Voicing is in the card and swaps with it. Structure is compiled in and
// does not, which is the gap a reader should be told about rather than left to
// find: a card can state a structural commitment in prose and the gate cannot
// read it.
// Asked and Checked are separate fields because they are separate places, and
// collapsing them is the error the README is being rewritten to stop making: a
// card asks for its whole voicing in prose, and only the POSTPROC part of that
// is a check the card owns. The rest of the voicing rules are compiled in
// beside the structure rules and travel with the binary, not the card.
type axisFact struct {
	Name      string   `json:"axis"`
	What      string   `json:"what_it_is"`
	Asked     string   `json:"where_a_card_asks_for_it"`
	Checked   string   `json:"where_it_is_checked"`
	Authored  string   `json:"what_a_card_can_change_about_the_check"`
	Rules     []string `json:"rules_that_check_it"`
	Measuring string   `json:"what_has_measured_it"`
}

// laneFact is the one place structure does vary. Not by card — by who is going
// to read the turn.
type laneFact struct {
	Name    string   `json:"lane"`
	When    string   `json:"chosen_when"`
	Drops   []string `json:"rules_dropped,omitempty"`
	Adds    []string `json:"rules_added,omitempty"`
	Because string   `json:"why"`
}

func collectFacts(c *scan.Card, fs *flag.FlagSet) docsFacts {
	// buildVersion() is deliberately not a fact. It reports `git describe`, so
	// the first draft stated "v0.2.0-2-g0e017cf" — accurate for one commit and
	// wrong forever after. A generated README can restate a card count on the
	// next `make readme`; a version taken from the working tree is stale before
	// the commit that adds it lands.
	shipped := shippedRegexRules()
	f := docsFacts{
		Module:     "github.com/justinstimatze/cope",
		NameOrigin: "a cope is the upper half of a foundry mould — the half carrying the shape being cast into",
		CardSource: c.Source,
		CardID:     c.CardID,
		Composition: composition{
			Never: len(c.Never), Wrong: len(c.Wrong), Mes: len(c.Mes), Tests: len(c.Tests),
			Quirks: len(c.Quirks), Traits: len(c.Traits), Regex: len(c.Rules),
			ShippedRegex: len(shipped),
		},
		ShippedRegexRules: shipped,
		RenderedCard: sizes{
			Card:      len(c.Render()),
			Refresher: len(c.RenderRefresher()),
		},
		FocusNote:     "the refresher size in rendered_sizes_chars is the standing fallback, used when a session has no history. The focused form varies with how many rules are firing, so do not state a fixed size for it",
		ShapeRules:    shapeRules,
		Lanes:         lanes,
		NeverBudget:   maxNeverRulesForDocs,
		NeverOverflow: overflowText(c),
		Gating:        gatesOf(c),
		CardRules:     authoredBy(c),
		Hooks: []hookFact{
			{"SessionStart", "cope-gate --inject", "prints the card as prompt text. SUPERSEDED by the output style and off by default: it puts the card in one turn-zero message, which Claude Code buries as the conversation grows. Measured 2026-08-03, a card asking for a bolded label on every paragraph and an emoji on every heading produced, through this hook, prose indistinguishable from no card at all. It stays for anyone who wants the old delivery, and it stands down on its own when a cope output style is active"},
			{"Stop", "cope-gate", "scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log"},
			{"UserPromptSubmit", "cope-gate --refresher", "reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts. Falls back to the standing CONTINUE TEST when the session has no history yet, and stays quiet until the last injection has aged past --refresh-every"},
			{"PreToolUse", "cope-gate --pretool", fmt.Sprintf("scores the description, body or content field an external write is about to post, matched against the Linear save tools named in the settings block. Warn-only: it returns additionalContext and never a permissionDecision, so the call goes through and the model learns what the prose scored. It writes no session state, and it scores in the external lane, which drops %s", strings.Join(scan.NeedsAnswerableReader(), ", "))},
		},
		SettingsJSON: settingsJSON,
		StyleInstall: styleInstallFact{
			Block: fmt.Sprintf("go install %s/cmd/cope-gate@latest\ncope-gate --setup\n# then: /config -> Output style -> %s",
				"github.com/justinstimatze/cope", shippedStyleName()),
			Setup:  "cope-gate --setup does the whole install: emits the output style, wires the three hooks into ~/.claude/settings.json with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. --setup --dry-run prints what it would change and writes nothing",
			Emit:   "cope-gate --output-style writes the loaded card to ~/.claude/output-styles/<card>.md, and COPE_CARD=<name> in front of it emits a different one",
			Select: fmt.Sprintf("pick it under /config -> Output style. The entry to look for is named %q, which is the shipped card's id and not the word cope — say the name, because a reader scanning that menu for something called cope will not find it. The standalone /output-style command was removed in Claude Code v2.1.91, so /config is the way; the same thing can be set as \"outputStyle\": %q in .claude/settings.local.json", shippedStyleName(), shippedStyleName()),
			Timing: "a style is read once at session start, so a new selection or a re-emitted card applies at the next session or after /clear",
			Why:    "an output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through the hook",
		},
		StateFiles: []stateFact{
			{"$XDG_STATE_HOME/cope/violations.jsonl", "0600", "one JSON record per violation, carrying the matched text and about 70 characters either side"},
			{"$XDG_STATE_HOME/cope/refresher-<session-id>", "0600", "an empty file whose mtime is the refresher clock"},
			{"$XDG_STATE_HOME/cope/session-<session-id>.json", "0600", "the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts"},
		},
		Layout: []layoutFact{
			{"card/claude_voice.effigy", "the shipped card, in effigy notation"},
			{"card/rules.json", "generated from it; embedded in the binary"},
			{"card/demo/", "other cards, each written to sound like something else"},
			{"cmd/cope-gate/", "the hook binary"},
			{"internal/scan/", "the structure rules, the card's regex rules, and the card renderer"},
			{"internal/effigy/", "the .effigy reader, so a card is usable as written"},
			{"internal/transcript/", "Claude Code JSONL reader, and which lane a turn was written in"},
			{"replay/", "the blind-pairs and discrimination harnesses, and their own README"},
			{"demo/", "this README written again under each demo card"},
			{"tools/", "card compiler, effigy-backed scorer, cross-project sweep"},
			{"MEASUREMENTS.md", "what was run, on how much text, and what it said"},
		},
		Related: []relatedFact{
			{"caveman", "https://github.com/JuliusBrussee/caveman",
				"a separate project, by a different author, that compresses agent replies to cut output tokens — a third axis again. cope shapes prose, basanite tracks vocabulary, caveman shortens. Worth naming because a reader wanting fewer tokens rather than different structure should go there instead"},
			{"effigy", "https://github.com/justinstimatze/effigy",
				"a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it"},
			{"basanite", "https://github.com/justinstimatze/basanite",
				"the same problem answered the other way round, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead — lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you. Its own README calls that awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable"},
		},
	}
	for _, t := range c.Tests {
		f.TestNames = append(f.TestNames, t.Name)
	}
	for _, r := range c.Rules {
		// The patterns stay out, for the reason Render withholds them:
		// showing the regex invites writing around it.
		f.RegexRules = append(f.RegexRules, ruleFact{r.ID, r.Action, r.Why})
	}
	f.Axes = axesOf(c)
	fs.VisitAll(func(fl *flag.Flag) {
		def := redactHome(fl.DefValue)
		if def == "" {
			def = "(empty)"
		}
		f.Flags = append(f.Flags, flagFact{"--" + fl.Name, def, fl.Usage})
	})
	return f
}

// redactHome rewrites the running user's home directory back to $HOME.
//
// --log resolves its default at registration time, so the flag set reports an
// absolute path containing the account name. That path was about to be copied
// verbatim into a published README. Facts get introspected from the machine
// doing the writing, which makes this the one place a private path can enter
// the docs by accident.
func redactHome(s string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" || home == "/" {
		return s
	}
	return strings.ReplaceAll(s, home, "$HOME")
}

// maxNeverRulesForDocs mirrors scan's unexported budget. Kept as a constant
// here rather than exported from scan, since one number in a prompt does not
// justify widening that package's surface.
const maxNeverRulesForDocs = 10

const docsSections = `Nowhere in the README state a version number for cope itself. It would be read
off the working tree and is wrong by the next commit.

The same goes for any figure a rebuild changes. Rendered sizes, character
counts, how many rules a card holds: state these ONLY where the fact is the
point of the sentence and the sentence would be empty without it. A number
copied into prose as colour is wrong at the next card edit and nobody rebuilds
the page to find out. The counts in facts are there to keep a claim honest, not
to be quoted for texture.

Every identifier goes in backticks, every time it appears — file paths, rule
ids, flag names, commands, card block names, the @gate and @shape forms. Two
conventions on one page is worse than either, and a bare rule id mid-sentence
reads as a typo the reader has to stop and rule out. A markdown link counts as
marked; a second mention in the same paragraph still needs its backticks.
Where a placeholder needs angle brackets, put the whole thing inside backticks
so nothing has to be escaped. Never write an HTML entity — no &lt;, no &gt;, no
&amp;. A reader looking at the raw file sees exactly what you typed.

Refer to another part of this page by naming its subject in your own words, not
by using its heading as a noun. A heading pressed into service as the subject
of a sentence parses badly, and the headings are yours to choose, so a sentence
depending on one is a sentence depending on a decision you have not made yet.

The short labels in this prompt — the phrases opening a clause before the
instruction that follows — are scaffolding for reading the brief. They are not
text. Do not copy one onto the page as a heading, a bold label, or the opening
of a paragraph, and do not keep their shape: a label here is a question or a
relative clause, and both read as a sentence that stopped early once the heading
around them is gone. How a paragraph opens is the card's decision and not this
prompt's.

The card's dangling_end and buried_decision rules describe how a conversational
turn should close, because a reader coming back to a terminal has to decide
whether to continue. A README has no such reader and no such decision. Write the
sections and stop; do not append an all-clear, and do not manufacture something
for the reader to act on so a rule about replies has something to match.

1. Title, one sentence saying what cope is, and facts.name_origin worked in
   wherever it sits best. One clause, not a paragraph. The sentence has to carry
   three things: that cope ships an opinionated card that is on once installed
   and is the whole product for most readers, that it checks two different jobs
   rather than one — how a reply sounds and how it is shaped — and only then
   that the card is a file, so it can be edited or swapped. Do NOT open by
   describing a card the reader writes. Most readers will never write one; the
   default is the thing being offered and authoring is how it is changed later.
   A title that says "gate" describes one of the three parts.
   Then, immediately after that sentence and before any heading, send the reader
   to demo/README.md in one short line of its own. Say that every file in that
   directory is this README written again from a different card, same prompt and
   same facts, so the card is the only thing that changed — and that reading two
   of them against each other shows what a card does faster than the rest of
   this page explains it. That is not a modest claim to bury in a later section:
   a reader who follows the link understands the tool, and a reader who does not
   is taking prose about registers on trust. Link it as demo/README.md. Mention
   the maximal one by name — demo/README.claude-maximal.md, written from a card
   that instructs every tic this model is measured to have — as the one that
   makes the point in a single glance. The warning that it is deliberately hard
   going belongs INSIDE that sentence, as a clause — not appended after it as a
   fragment, which is how it landed the first time and reads as a stub bolted on
   at the end. A reader who clicks it forewarned finds it funny; the same reader
   landing in it unwarned finds the thing they came here to get away from. Name
   it, do not recommend it.
   ONE EXCEPTION, only when the user message says this render IS the front page:
   the front page is written from the fieldguide card, so demo/README.fieldguide.md
   is this page under another path. Do not send a reader to it. Every other link
   above stands, the maximal one included.
   Write nothing anywhere on that page about the voice being deliberate, chosen,
   instructed, or a demonstration. No italic note above the title, no aside, no
   footer. The page is the joke and a joke does not introduce itself; a reader
   who wants to know why it reads that way has the link above and demo/README.md
   at the end of it.
2. The two things, from facts.axes. This is the frame every later section is
   written against, so it comes second and stays short — five or six sentences.
   Voicing is what the sentences sound like, it lives in the card, and swapping
   the card swaps it; that is the half with a measured result behind it.
   Structure is where the decision sits and how the reply ends, and it is
   compiled into the binary, so it is the same whichever card is loaded. Give
   the reader one concrete instance of each so the distinction is not abstract:
   a sentence that reaches for the balanced two-beat is a voicing problem, and
   a reply that names an open problem in its last paragraph and then stops,
   leaving "continue" nothing to refer to, is a structural one. The same reply
   can be fine on one axis and bad on the other, which is the reason to keep
   them apart. Then say how far a card reaches into the structure half, from
   facts.card_authored_rules, and be exact because this is what most recently
   changed. Two directions, both with the syntax from the facts: a card declines
   a built-in rule it disagrees with, and a card states a structural rule of its
   own that the gate then checks. Give one reason for each from the two example
   facts — the first exists because a card whose VOICE asked for something a
   built-in rule catches was marked down for obeying itself, the second because
   a card's own commitment about how a reply ends had nowhere to be checked.
   Then name the boundary from facts.card_authored_rules.what_it_still_cannot_do
   in one sentence, without apology and without calling it a roadmap. Do not put
   a rate in this section. MEASUREMENTS.md has the run and the reasons its
   numbers do not carry more than that.
3. The problem, in three parts and in this order: where an instruction sits,
   what an instruction can say, and the failure no instruction reaches. The
   first is the largest and most readers arrive holding it backwards, so it
   leads.
   Take the first one first. A reader who has edited a global CLAUDE.md and
   watched it not stick almost certainly believes that file is the system
   prompt. It is not. It
   arrives as one message attached to the first turn, and the conversation
   buries it under everything written after. An output style is in the system
   prompt itself, which the harness re-reminds the model of as the conversation
   runs. Say that moving one card between those two places, without changing a
   word of it, is most of why cope works. Say it was measured, point at
   MEASUREMENTS.md, and put no counts on this page: a reader here wants to know
   the thing works, and the one who wants the run will follow the link. Do not
   name the delivery that lost — section 4 gives it one clause, and a reader
   must not leave this page wiring it.
   Then the second. Instruction alone does not fix the phrasing: a
   global CLAUDE.md banning the "not A, it's B" flip is read every turn, and the
   flip still appeared twice in the session that built this, while the ban was
   the topic. Naming a surface form pushes the move into a variant. That is the
   voicing side.
   Then the third, which is not a phrasing problem at all. The structural side
   is a different complaint with a different cause — an ending that leaves the reader nothing
   to answer costs a whole round trip, and it is not a phrasing habit an
   instruction could have banned.
   Close by saying the flip is an anecdote about one rule and naming what the
   claim actually rests on: the blind discrimination test, where a reader shown
   only a voice's own description of itself picks which of two replies was
   written under it. Send them to MEASUREMENTS.md for the rate and the caveats.
   Do not cite the blind preference runs here — both sides of those were written
   under a card, so they compare two ways of writing well and cannot see a voice
   being swapped.
4. Install, in the order somebody actually does it, and the order matters more
   here than anywhere else on the page. Two commands and one menu choice, from
   facts.output_style_install. Open with
   output_style_install.copyable_install_block, reproduced VERBATIM in a fenced
   code block — character for character, including the comment line naming the
   menu entry, which is not decoration: a reader who has just installed
   something called cope will scan that menu for the word cope, and the entry
   is the card's id instead. An install that ends with a style selected and
   nothing changed is the one failure on this page that costs a reader the
   whole tool. Name the entry again in the prose that follows, from
   output_style_install.select. Lead with those and give them as a copyable block. Keep the prose
   in that same order afterwards — what the two commands did, and only then the
   menu choice that follows them. Naming the last step before explaining the one
   before it makes a reader who is following along stop and scroll. Say what
   setup did rather than listing what it saves the reader from doing, and include that
   it backs the settings file up, that a second run changes nothing, and that it
   leaves other tools' hooks alone — a reader letting a tool edit their settings
   deserves to know the shape of it before they run it, not after. Mention
   --dry-run in the same breath, in a clause. Then give the by-hand route from
   emit_by_hand for anyone who would rather not have their settings written to. State facts.output_style_install.when_it_applies
   plainly, because a reader who picks a style and sees no change in the running
   conversation will conclude the tool is broken. Then one sentence for why the
   card goes here rather than into a hook, from why_here_rather_than_a_hook —
   short, since it is the reason and not the instruction.
   Only then the hooks, and be exact about what they are now for: the card no
   longer arrives through one. Reproduce facts.settings_json_block VERBATIM in a
   fenced json code block as the hooks block of ~/.claude/settings.json. Copy it
   character for character, do not reformat or abbreviate it, and do not add a
   hook that is not in it — in particular do not add SessionStart, which is
   absent on purpose. Say what the two remaining hooks BUY, in one sentence
   each, and no more than that: the reply gets scored after it is written, and
   the rules that have actually been firing get restated mid-session, which a
   file written once cannot do. The mechanics — what each hook reads, what it
   writes, what it falls back to — belong to the scoring section further down
   and must not appear twice. A reader meeting the same two paragraphs again
   checks whether they have scrolled backwards. Say that the voice works without
   either of them and these are the measurement half.
   Under the block, one sentence that the commands assume go install's target
   directory is on PATH in the environment the hook runs in, and that a hook
   which silently does nothing usually wants the absolute path instead. Note
   that a clone builds with make install and needs no effigy checkout, because
   card/rules.json is committed and compiled in.
   Mention --inject in one clause here, as the superseded delivery that remains
   for anyone who wants it and stands down on its own when a cope output style
   is active. That clause is the ONLY prose mention of it on the whole page —
   the flag table lists it because the flag table lists everything, and that
   does not count, but no later section may name it again. It has appeared three
   times on this page before now. Do not explain it further and do not put it in
   a code block: a reader following this page must not end up wiring the weak
   path.
5. Writing in another voice. This is the capability the first sentence promised,
   so it comes before the machinery rather than after it. The gate reads .effigy
   directly, so a card is usable as written and needs no Python and no effigy
   checkout. A card dropped in $XDG_CONFIG_HOME/cope/cards is reached by
   --card <name> or COPE_CARD; --rules takes a path from anywhere; make cards
   installs the demo set. A name that resolves to nothing is an error and
   nothing is injected, rather than a session writing in the shipped voice while
   its config names another one. Point at card/demo/lecturer.effigy as the one
   to read first, because it differs from the shipped card on register alone and
   is what the discrimination run measured. Say that what a card can change is
   the voicing axis, and refer back by section title rather than arguing it
   again — the sections are numbered in this prompt and not in the README, so a
   sentence pointing at "section 2" points at nothing a reader can find. Do not
   restate the numbers here; MEASUREMENTS.md has them.
   Then link demo/README.md, in a sentence of its own, as the shortest way to
   see what a card does without writing one: every file under demo/ is this
   README written again from a different card, same prompt and same facts, so
   the only thing that varies between them is the voice. Link the directory
   index rather than listing the cards — that page is generated from whatever
   was last built and this one should not carry a second copy of the list.
   Last, one sentence that card/demo/handoff.effigy is the exception in that
   directory: it is a hypothesis rather than a voice, keeping the shipped card's
   handoff rules and dropping everything about prose, and it is meant to be run
   through make pairs against the full card rather than rendered. make cards
   installs it with the rest, so a reader who lists their cards will find it
   there and should know it is not a register to write in.
6. What the hooks do differently, from facts.hooks and
   facts.gated_card_items. Cover the two --setup wires, in this order and no
   others: Stop scores the reply and records which rules fired.
   UserPromptSubmit reads that record and injects only the items gated on what
   has been firing, naming the counts — so the mid-session text is chosen from
   measured output rather than fixed in advance, which is the one thing a
   pasted CLAUDE.md cannot do. Say this plainly and without overselling: it is
   one mechanism, not a guarantee, and the A/B in the repo does not separate
   the refresher from no refresher. Do NOT mention --inject or SessionStart
   here: the install section already gave it its one clause, and a superseded
   path named twice on one page starts to look like an option. Do not number
   the hooks in the heading — the count has changed once already and a heading
   that carries it goes stale silently.
   This section owns the mechanics, so it is where the detail goes: what each
   hook reads, what it writes, what it falls back to, when it stays quiet. The
   install section carries one sentence each and no more, and the two must not
   read as the same two paragraphs twice.
7. Why effigy notation, in two or three sentences from facts.related_projects,
   and why basanite is the wrong instrument for this in one more. A reader who
   knows neither project should still follow it.
8. The rules, grouped by facts.axes rather than by where they are implemented.
   Two lists: the voicing rules and the structure rules, each rule stated from
   facts.shape_rules and facts.regex_rules_in_shipped_card verbatim in
   substance — do not infer what a rule catches from its name.
   Take the POSTPROC rules from regex_rules_in_shipped_card and NOT from
   regex_rules. This section documents what a reader installs, and regex_rules
   holds the rules of whichever card THIS page is written from, which under any
   card in demo/ is a different set. Listing that set here is how the front
   page came to document two rules nobody has while omitting the three
   everybody gets. If the page is written under a card carrying POSTPROC rules
   of its own and they are worth a mention, name them as that card's and keep
   them out of the list. Then say what the grouping implies about
   the implementation, which is the part a reader can use: a POSTPROC pattern
   matches a span of text, so it can only ever describe wording, and every
   voicing rule that needed more than a pattern had to be written in Go beside
   the structure rules. That is why cope's shipped card carries only
   facts.card_composition.postproc_rules_in_shipped_card of them. Take that
   number from that field and from nowhere else: postproc_rules alongside it
   counts the card THIS page is written from, which is a different card and a
   different number on every page under demo/, and attributing that one to the
   shipped card is how every demo page came to state a count cope does not
   have. If the two differ and the difference is worth a clause, say whose is
   whose; otherwise state the shipped number and leave the other out. Use the basanite boundary in
   facts.related_projects here: a reader who expects a long list of banned
   phrases should learn that the list lives in another tool on purpose. Close
   with facts.lanes, which is the one place the structure rules do vary — not by
   card, by who is going to read the turn. State which rules the loop lane drops
   and adds and why, from the fact's own wording.
9. The flag table, from facts.flags verbatim. Do not paraphrase a default.
10. What lands on disk, from facts.state_files. State plainly that the log
   quotes replies back.
11. Editing the card: effigy owns the .effigy grammar, make rules regenerates,
   make check-rules is what CI runs so the enforced and injected rules cannot
   drift. Mention the NEVER budget of facts.never_rule_budget and that anything
   over it is reported at load rather than dropped silently. The budget is
   charged against each injection separately, not against the card file, since
   SessionStart prints the always-on rules and the refresher prints the
   evidence-gated ones and no code path renders their union — so a card may
   hold more NEVER rules in total than the budget and still be healthy. Say
   that. Do not state how many NEVER rules the shipped card has: the count is
   not in facts, it rots on the next card edit, and it invites the reader to
   compare it against the budget and conclude something is wrong.
   facts.never_rules_over_budget is the authoritative list of rules that really
   are discarded unrendered; it is empty when the card is healthy.
   Then both card-authored forms, from facts.card_authored_rules, with each
   syntax verbatim and the @shape vocabulary listed exactly as the fact gives it
   — a reader writing a card needs the selectors and predicates spelled out, and
   an approximation of them is worse than none. Say that a rule id has to be one
   the gate has for @gate and must not collide with one for @shape, and that a
   wrong id is reported at load rather than ignored. Say the reason after the
   dash is required in both, because a rule a card wrote and a rule a card
   refused are equally unreviewable without one. Note that a declined rule still
   runs and only this card's score drops it, so a backfill still reports what it
   would have caught, and that a @shape violation is reported in the card's own
   words rather than in any sentence the binary supplies.
12. Calibrating. cope-gate --backfill scores a whole session transcript at once
   and is how the rules were chosen; tools/backfill-sweep.sh runs it over the
   N largest transcripts found anywhere under ~/.claude/projects, which is not
   the same as one per project. Say that hits-per-character is the
   metric worth watching rather than the share of turns hit, because that
   second number tracks how long the turns were. Send the reader to
   MEASUREMENTS.md for the rates rather than quoting one here.
13. Known limits, and the axis split is what organises them. labelled_opening is
   not a tagger. ask_not_last says nothing about the order of several asks. The
   hit rate is roughly four fifths structure by facts.axes, and the A/B run
   found that four fifths tracks what a reply was for rather than how it was
   written, so it is a description of the output and not a judgement of it — the
   discrimination test is where the judgement lives, and it covers voicing only.
   The largest limit is the one the opening section named, and state it as it
   now stands rather than as it stood: a card can decline a built-in rule and
   write one of its own, and the vocabulary it writes in counts words and
   sentences and asks whether a block poses a question — so the compiled rules
   remain the only place a check like clause_symmetry can live, and a card
   wanting something outside both that vocabulary and a POSTPROC regex still has
   nowhere to put it. Say too that both directions are the card marking its own
   homework: a decline lowers that card's score and an assertion raises it, both
   are worth reading with the reason attached, and that is why the syntax
   requires one. Name the limit by the opening section's title rather than by a
   number the README does not print. Point at MEASUREMENTS.md.
14. Layout, from facts.layout.
15. A one-line footer saying this README was written by
    tools/generate_readme.py from the prompt cope-gate --author-docs emits, and
    checked with cope-gate --check. Do not explain it further.
16. License line: MIT, contact justin@justinstimatze.com.`

// buildDocsPrompt assembles the prompt: the card the prose must obey, the facts
// it may assert, and the sections to fill. The card goes in whole rather than
// summarised — it is the same text a session receives at SessionStart, so the
// documentation is written under the rules it documents.
func buildDocsPrompt(c *scan.Card, fs *flag.FlagSet) string {
	facts, err := json.MarshalIndent(collectFacts(c, fs), "", "  ")
	if err != nil {
		facts = []byte("{}")
	}
	var b strings.Builder
	b.WriteString(`<cope-author-docs>
Write README.md for cope. Every sentence is yours; the facts below are exact
and the voice below is binding.

This prompt is emitted by cope-gate itself, so it states no prose of its own.
That is deliberate: a generator that supplies phrasing writes the same README
every time, and the point of the card is that the phrasing comes from the card.

## The voice

Not a suggestion. This is the card a session is given at SessionStart, and the
gate that scores the result reads the same file.

`)
	b.WriteString(c.Render())
	b.WriteString("\n## The facts\n\nAssert these. Do not restate a number from memory.\n\n```json\n")
	b.Write(facts)
	b.WriteString("\n```\n\n## The sections\n\n")
	b.WriteString(docsSections)
	b.WriteString(`

## Then check it

    cope-gate --check README.md

That runs the card's regex rules and every shape rule over what you wrote —
the same pass a reply gets at Stop. Revise until it reports nothing, or until
what remains is a hit you can defend out loud. A README that cannot pass the
card it documents is the argument against the card.
</cope-author-docs>
`)
	return b.String()
}

// runCheck scores a prose file, or stdin when path is "-". It exists so the
// repo's own documentation goes through the gate that reads every reply;
// --backfill only understands transcript JSONL.
func runCheck(c *scan.Card, path string, minCV float64, logPath string) {
	var (
		raw []byte
		err error
	)
	if path == "-" {
		raw, err = os.ReadFile("/dev/stdin")
	} else {
		raw, err = os.ReadFile(path)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "cope-gate: %v\n", err)
		os.Exit(1)
	}

	text := string(raw)
	v := c.Check(text, minCV)
	if len(v) == 0 {
		fmt.Printf("%s: clean (%d chars)\n", path, len(text))
		return
	}
	fmt.Printf("%s: %d violation(s), %d chars\n\n", path, len(v), len(text))
	for _, x := range v {
		fmt.Printf("  [%s] %s\n      %q\n", x.RuleID, x.Why, x.Matched)
		if x.Context != "" && x.Context != x.Matched {
			fmt.Printf("      ...%s...\n", x.Context)
		}
	}
	logViolations(v, "check:"+path, logPath)
}

// lanes is the whole of the lane logic as the docs need to see it. Kept beside
// shapeRules rather than derived from internal/scan, because suppressedInLoop is
// unexported and a map has no order to report.
var lanes = []laneFact{
	{
		Name:    "interactive",
		When:    "any turn that is not a loop turn",
		Because: "somebody is waiting at a terminal and the ending is where they decide what happens next",
	},
	{
		Name:  "loop",
		When:  "the prompt that opened the turn was /loop or /goal, or the sentinel a dynamic-pacing loop sends itself",
		Drops: scan.NeedsAnswerableReader(),
		Adds:  []string{"unverified_done", "loop_ask"},
		Because: "nobody is reading yet. A report that correctly names what it left open and stops would fail three " +
			"of the dropped rules, and a question in it lands in a log where the next iteration reads it as an " +
			"instruction to itself. What replaces them is the claim check: a report saying the work is done has to " +
			"say what it ran",
	},
	{
		Name:  "external",
		When:  "the PreToolUse entry scores prose an external write is about to post, rather than a reply",
		Drops: scan.NeedsAnswerableReader(),
		Because: "a ticket has a reader and no ending they can answer. It is read days later by somebody who was " +
			"not in the session, which is the condition every rule that survives the drop was written for — so " +
			"this lane swaps nothing in, where the loop lane swapped four for two",
	},
}

// axesOf splits every live rule into the two things cope is about, and states
// which half a card can actually change.
//
// The regex rules are voicing by construction rather than by judgement: a
// POSTPROC pattern matches a span of text, so the most it can ever describe is
// wording. Nothing about where a decision sits in a reply reduces to a pattern,
// which is why the structure rules are Go — and being Go is exactly why they do
// not travel with the card.
func axesOf(c *scan.Card) []axisFact {
	byAxis := map[string][]string{}
	for _, r := range c.Rules {
		byAxis["voicing"] = append(byAxis["voicing"], r.ID)
	}
	for _, s := range shapeRules {
		byAxis[s.Axis] = append(byAxis[s.Axis], s.ID)
	}
	return []axisFact{
		{
			Name:      "voicing",
			What:      "what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed",
			Asked:     "the card, entirely — VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC. Swapping the card swaps every word of it",
			Checked:   "split. The card's POSTPROC patterns, plus rules compiled into internal/scan that needed more than a pattern",
			Authored:  "the POSTPROC patterns and nothing else",
			Rules:     byAxis["voicing"],
			Measuring: "the blind discrimination test, which is the instrument with a result behind it",
		},
		{
			Name:      "structure",
			What:      "the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives 'continue' something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it",
			Asked:     "nowhere the gate reads. A card can write down how it wants a reply to end and the gate has no way to hear it",
			Checked:   "the binary, in internal/scan — the same rules whichever card is loaded, varying only by lane",
			Authored:  "nothing",
			Rules:     byAxis["structure"],
			Measuring: "the hit rate, which is roughly four fifths structure and which the A/B run showed tracks what a reply was for",
		},
	}
}

// overflowText flattens the budget overflow for the facts blob.
func overflowText(c *scan.Card) []string {
	var out []string
	for _, g := range c.NeverOverflow() {
		out = append(out, g.Text)
	}
	return out
}

// exemptionsOf reports the @gate mechanism and whatever the loaded card does
// with it. Exemptors is empty for the shipped card, which claims none — the
// README should describe the mechanism either way, because a reader writing a
// card needs it and the shipped card is not the only card.
func authoredBy(c *scan.Card) authorFact {
	f := authorFact{
		GateSyntax:  "@gate <rule_id> off — <why>, one per line in the card header",
		GateExample: "card/demo/lecturer.effigy declines clause_symmetry and dangling_end, because its VOICE block asks for the balanced landing and the arriving close those two rules catch",
		ShapeSyntax: "@shape <id>: <selector> <predicate> — <why>, one per line in the card header",
		ShapeWords: []string{
			"selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply",
			"predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask",
		},
		ShapeExample: "card/demo/handoff.effigy asserts readable_cold — last paragraph words <= 60 — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90",
		Limit:        "the @shape vocabulary counts words and sentences and asks whether a block poses a question, and nothing more. It cannot express what the compiled rules express — clause_symmetry is not writable in it and is not meant to be — so a card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it",
	}
	for _, d := range mustGates(c) {
		f.ThisCard = append(f.ThisCard, "declines "+d.Rule+" — "+d.Why)
	}
	for _, r := range mustShapes(c) {
		f.ThisCard = append(f.ThisCard, "asserts "+r.ID+" — "+r.Why)
	}
	return f
}

// mustShapes drops the validation error for the same reason mustGates does.
func mustShapes(c *scan.Card) []scan.ShapeRule {
	r, _ := c.Shapes()
	return r
}

// shippedRegexRules returns the POSTPROC rules on the embedded card, whichever
// card --rules loaded. It returns nil rather than failing when the embedded
// card will not parse, because a docs prompt that refuses to print is worse
// than one field missing from it, and internal/scan's tests already fail the
// build in that case.
// shippedStyleName is the output-style entry --setup actually writes, which is
// the SHIPPED card's id and never the loaded one. The install block took the
// loaded card's name for one render and the front page — rendered from a demo
// card — told readers to select an entry --setup does not create. Same fault as
// the POSTPROC counts: a fact about the product, taken from whichever card
// happened to be loaded.
func shippedStyleName() string {
	c, _, err := scan.ParseCard(card.RulesJSON)
	if err != nil {
		return "cope"
	}
	return styleName(c)
}

func shippedRegexRules() []ruleFact {
	c, _, err := scan.ParseCard(card.RulesJSON)
	if err != nil {
		return nil
	}
	var out []ruleFact
	for _, r := range c.Rules {
		// Patterns stay out here for the reason they stay out of RegexRules:
		// showing the regex invites writing around it.
		out = append(out, ruleFact{r.ID, r.Action, r.Why})
	}
	return out
}

// mustGates drops the validation error, which ParseCard has already reported to
// the operator by the time anything asks for documentation.
func mustGates(c *scan.Card) []scan.GateDirective {
	d, _ := c.Gates()
	return d
}

// gatesOf tallies the card's @when conditions. A condition with no items is
// not reported, so the README describes gating that exists rather than gating
// the notation merely permits.
func gatesOf(c *scan.Card) []gateFact {
	tally := map[string]int{}
	for _, n := range c.Never {
		if n.When != "" {
			tally[n.When]++
		}
	}
	for _, t := range c.Tests {
		if t.When != "" {
			tally[t.When]++
		}
	}
	for _, w := range c.Wrong {
		if w.When != "" {
			tally[w.When]++
		}
	}
	keys := make([]string, 0, len(tally))
	for k := range tally {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]gateFact, 0, len(keys))
	for _, k := range keys {
		out = append(out, gateFact{k, tally[k]})
	}
	return out
}
