# cope

A cope is the upper half of a foundry mould, the half carrying the shape being cast into, and this one ships a single opinionated card that is live the moment it is installed, scoring every reply on two separate questions — what the sentences sound like, and whether the reply ends somewhere its reader can act from — with the card itself sitting on disk as a file that can be edited or swapped once that matters to you.

Every file in [demo/README.md](demo/README.md) is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed; read two of them against each other and you will understand what a card does faster than the rest of this page can explain it, and [demo/README.claude-maximal.md](demo/README.claude-maximal.md) — rendered from a card that instructs every tic this model is measured to have, and deliberately hard going for exactly that reason — makes the point in a single glance.

## The two things a reply can get wrong

Ask what went wrong with a bad reply and the answer will be about wording, because wording is the part you can quote back. So let me tell you about the part you cannot. Voicing is what the sentences sound like — register, rhythm, diction, where flair is licensed — and it lives entirely in the card, which means swapping the card swaps every word of it; that is also the half with a measured result behind it. Structure is where the decision sits and how the reply ends, and it is compiled into the binary, so it is the same set of rules whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing problem. A reply that names an open problem in its closing paragraph and then simply stops, leaving "continue" nothing to refer to, is a structural one, and it can happen inside prose you would be glad to read aloud. Keeping the two apart is the only way to say which of them a reply actually failed.

A card now reaches into the structural half, in two directions and no further. Writing `@gate <rule_id> off — <why>`, one per line in the card header, declines a built-in rule the card disagrees with; that exists because a card whose `VOICE` block asked for something a built-in rule catches was being marked down for obeying itself. Writing `@shape <id>: <selector> <predicate> — <why>`, also one per line, states a structural rule of the card's own that the gate then enforces; that exists because a card's commitment about how a reply ends had nowhere to be checked. The vocabulary those assertions are written in counts words and sentences and asks whether a block poses a question, and nothing beyond that, so a check like `clause_symmetry` is not expressible in it and was never meant to be. A card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it.

## Why the obvious fix does not hold

You have edited a global `CLAUDE.md`, and you have watched it not stick, and you very likely concluded that the model is ignoring you. Let me offer a duller explanation. That file is not the system prompt. It arrives as a single message attached to the first turn, and everything written afterwards buries it a little deeper, until by turn forty it is something that happened rather than something in force. An output style sits inside the system prompt itself, which the harness keeps re-reminding the model of as the conversation runs. Take one card, change not a word of it, move it from the first place to the second: that is most of why cope works, it was measured, and `MEASUREMENTS.md` has the run. The instruction was never being ignored. It was being outvoted.

Put the card in the right place and the phrasing still drifts, which is the second and smaller problem. A global `CLAUDE.md` banning the "not A, it's B" flip is read on every single turn, and the flip appeared twice in the session that built this tool, while the ban was the subject under discussion. Name a surface form and the move relocates into a variant of it.

The third complaint is not about phrasing at all. A reply that mentions something unresolved and then stops costs you a whole round trip asking what you were supposed to decide, and no list of banned phrases was ever going to catch it, because nothing in it is a phrase.

So take the flip as an anecdote about one rule and nothing more. What the claim actually rests on is the blind discrimination test: a reader is shown a voice's own description of itself and two replies, and asked which of the two was written under it. `MEASUREMENTS.md` has the rate and the caveats. The only instrument that can see a voice being swapped is the one where the voice is the sole thing swapped.

## Installing it

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

Those two commands do the whole install. `cope-gate --setup` emits the loaded card as a Claude Code output style and wires the two hooks into `~/.claude/settings.json` with absolute paths, then prints the one step left. Before it touches that file it backs it up; it adds only what is missing, so a second run changes nothing; it leaves every other key alone, including other tools' hooks sitting on the same events; and it refuses outright to modify a settings file that does not parse. If you would rather look before it writes, `cope-gate --setup --dry-run` prints what it would change and writes nothing.

Then the step it printed for you. Under `/config`, choose Output style, and look for the entry named `claude_voice` — the shipped card's id, and not the word cope. This matters more than it looks: someone who has just installed a thing called cope will scan that menu for the word cope, fail to find it, select nothing, and walk away believing the install failed. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the route; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

For anyone who would rather not have a tool write to their settings, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` and stops there, and `COPE_CARD=<name>` in front of it emits a different card instead.

A style is read once, at session start. A new selection, or a freshly re-emitted card, therefore applies at the next session or after `/clear`, and not in the conversation you are currently sitting in. The card goes into a style rather than into a hook because the system prompt is the one place the harness keeps bringing back to the model's attention.

The hooks are the other half, and the card no longer arrives through either of them:

```json
{
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
    ]
  }
}
```

`Stop` buys you a score on the reply once it has been written. `UserPromptSubmit` buys you the rules that have actually been firing in this session, restated part-way through it, which a file written once cannot do. The voice works with neither hook present; these two are the measurement half.

Both commands assume that `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants an absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled into the binary. There is also `--inject`, the superseded delivery, kept for anyone who wants the old way and standing down of its own accord whenever a cope output style is active.

## Writing in another voice

The reason most tools like this want a build step is that their card format is a source language, and the gate here does not have that problem: it reads `.effigy` directly, so a card is usable exactly as written, with no Python and no effigy checkout anywhere near it. Drop a card into `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or with `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set. A name that resolves to nothing is a hard error and nothing is injected, which is deliberate — a session quietly writing in the shipped voice while its config names another one is the failure you would never notice.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, which is precisely what the discrimination run measured, and it also shows both card-authored forms in use: it declines `clause_symmetry` and `dangling_end`, because its own `VOICE` block asks for the balanced landing and the arriving close that those two rules catch. What a card can change is the voicing axis, in the sense set out in the two things a reply can get wrong, above; the structure rules are not on offer.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this same page rendered from a different card against the same prompt and the same facts, so the voice is the only variable in the set.

One file in that directory is not a voice. `card/demo/handoff.effigy` is a hypothesis: it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered as a page. `make cards` installs it along with the rest, so anyone listing their cards will find it sitting there looking like a register it is not.

## What the hooks add

Consider what a card cannot know at the moment it is written. It cannot know that you spent the last twenty turns opening every paragraph on a label, and it will spend the same number of words on the rule you have never once broken as on the one you break constantly. That is the gap the two hooks are for.

The `Stop` hook runs after each assistant reply, scores it against the loaded card, and appends two things: which rules fired, into the session's rolling state, and one record per violation, into the log.

The `UserPromptSubmit` hook reads that rolling state — the state file, not the violations log — and injects the card items gated on what has actually been firing, naming the counts as it does so. When a session has no history yet it falls back to the standing `CONTINUE TEST`, and it stays silent entirely until the last injection has aged past `--refresh-every`. So the mid-session text is selected from measured output rather than fixed in advance, and that is one mechanism rather than a guarantee: the A/B harness in this repo does not separate having a refresher from having none.

## Why a character-card notation

Reaching for a game-NPC format to describe prose style looks like a joke until you look at which blocks effigy already has. [effigy](https://github.com/justinstimatze/effigy) gives three of them for free: `POSTPROC` is regex rules carrying a warn action and applied after generation, `WRONG` holds an anti-pattern beside the replacement for it, and `TEST` holds a named question with a failing example and a passing one — which is how a rule can name a *move* rather than one wording of that move. The notation is used here entirely off-label, and it fits because those three blocks are what a prose gate needed anyway.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem from the other end, and is the one to reach for when this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead — lemma frequency against a baseline over real transcripts — and reports what you have in fact been leaning on lately, leaving the judgement with you; its own README calls that awareness rather than prohibition. Which one suits is a question about mood rather than about correctness, since a heavy hand is what you want when a habit is annoying you today and a moving measurement is what you want when you would rather watch a drift than legislate against it. They share no state and sit on different hooks, so running both is entirely reasonable.

Worth naming for a different reason: [caveman](https://github.com/JuliusBrussee/caveman), by a different author again, compresses agent replies to cut output tokens. cope shapes prose, basanite tracks vocabulary, caveman shortens — a reader who wants fewer tokens rather than different structure should go there instead.

## The rules

They are grouped here by which of the two jobs they check, rather than by where they happen to be implemented, because that is the grouping that tells you what a card can move.

Voicing, from the shipped card's `POSTPROC` block:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — the reflexive intensifier standing in for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, compiled into the binary because a pattern could not express them:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, which is the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to accounting for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. It is the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, compiled and identical under every card:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it then unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as things they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

That grouping tells you something about the implementation that is genuinely usable. A `POSTPROC` pattern matches a span of text, which means it can only ever describe wording, and every voicing rule that needed more than a span had to be written in Go alongside the structure rules. Hence the shipped card carrying three `POSTPROC` rules and no more. Anyone arriving here expecting a long list of banned phrases should know that the list lives in basanite on purpose, and that measuring your own drift is a better use of a list than banning it.

The structure rules vary in exactly one way, and it is not by card. It is by who is going to read the turn. Interactive is the default and covers any turn that is not a loop turn, on the reasoning that somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet in that lane. A report that correctly names what it left open and then stops would fail three of the dropped rules, and a question inside it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran.

## Flags

| flag | default | what it does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | `(empty)` | comma-separated variants to rotate through, implying `-ab` (default inject,hold; positive is the third) |
| `--ab-report` | `(empty)` | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | `(empty)` | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | `(empty)` | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | `(empty)` | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | `(empty)` | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | `(empty)` | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | `(empty)` | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | `(empty)` | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

Three files, all written `0600`:

- `$XDG_STATE_HOME/cope/violations.jsonl` — one JSON record per violation, carrying the matched text and roughly seventy characters either side of it.
- `$XDG_STATE_HOME/cope/refresher-<session-id>` — an empty file whose mtime is the refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a twenty-turn window. No prose is stored in it, only rule names and counts.

The middle file and the last one hold no text of yours. The first one does: the violations log quotes your replies back to you, in fragments, and it is on disk until you delete it. Setting `--log` to empty turns it off.

## Editing the card

effigy owns the `.effigy` grammar, `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so what the gate enforces and what the card injects cannot drift apart unnoticed.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped in silence. The subtlety worth knowing is what the budget is charged against: each injection separately, not the card file. One path prints the always-on rules and another prints the evidence-gated ones, and no code path anywhere renders their union, so a card may hold more `NEVER` rules in total than the budget and still be perfectly healthy.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is exactly this and no more. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`.

For `@gate`, the id has to be one the gate already has; for `@shape`, it must not collide with one. Either way a wrong id is reported at load rather than quietly ignored. The reason after the dash is required in both forms, on the grounds that a rule a card invented and a rule a card refused are equally unreviewable without one. A declined rule still runs — only this card's score drops it — so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

`card/demo/handoff.effigy` shows the assertion form doing real work. It declares `readable_cold` as `last paragraph words <= 60`, because its peak asks the reader to come back cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is not a guess: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

## Calibrating against your own transcripts

Choosing rules by intuition produces a gate that fires on prose you like, which is why every rule here was chosen against real output. `cope-gate --backfill` scores an entire session transcript in one pass, and `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects` — which is not the same as one per project, and the difference matters if your work is lopsided across repos.

The number to watch is hits per character, not the share of turns hit. That second figure mostly tracks how long your turns happened to be, and it will move whenever your habits of length move, whether or not anything about your writing changed. `MEASUREMENTS.md` has the rates.

## Known limits

`labelled_opening` is not a tagger, and it will miss labels it cannot see the shape of. `ask_not_last` says nothing at all about the ordering of several asks once they are in the right place. The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was *for* rather than how well it was written, so the rate describes the output rather than judging it; the judging happens in the discrimination test, and that covers voicing only.

The largest limit is the one named in the two things a reply can get wrong, above, and it should be stated as it now stands rather than as it stood. A card can decline a built-in rule and can assert one of its own, and the vocabulary it asserts in counts words and sentences and asks whether a block poses a question. That leaves the compiled rules as the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex with nowhere to put it. Both directions are also, unavoidably, the card marking its own homework — a decline lowers that card's score and an assertion raises it — which is the entire reason the syntax will not accept either one without a reason attached. `MEASUREMENTS.md` has the rest.

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | the shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in the binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | the hook binary |
| `internal/scan/` | the structure rules, the card's regex rules, and the card renderer |
| `internal/effigy/` | the `.effigy` reader, so a card is usable as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane a turn was written in |
| `replay/` | the blind-pairs and discrimination harnesses, and their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what was run, on how much text, and what it said |

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
