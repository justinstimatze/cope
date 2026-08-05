# cope

cope ships an opinionated card that is on the moment you install it, and scores what comes back on two counts: how the sentences sound, and how the reply is shaped for whoever reads it next. The card is a file — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — so you can edit it or write another.

Start at [demo/README.md](demo/README.md): every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them side by side shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have and is hard going as a result, makes the point in a single glance.

## Sound and shape

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. It lives in the card, and swapping the card swaps every word of it. This is the half with a measured result behind it, in MEASUREMENTS.md.

Structure is where the decision sits and how the reply ends, and it is compiled into the binary, so it is the same whichever card is loaded. A sentence reaching for the balanced two-beat — matched clauses either side of a comma, echoing a word across the joint — is a voicing problem. A reply that names an open problem in its last paragraph and then stops, leaving `continue` nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, which is why they are kept apart.

A card reaches into the structure half in two directions. `@gate <rule_id> off — <why>` declines a built-in rule, which exists because a card whose `VOICE` block asked for something a built-in rule catches was marked down for obeying itself: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason. `@shape <id>: <selector> <predicate> — <why>` states a structural rule of the card's own, which exists because a card's commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — since its peak asks the reader to re-enter cold and read the last block only. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and a card wanting a check outside that and outside a `POSTPROC` regex has nowhere to put it.

## Why a file in your home directory doesn't stick

If you have edited a global `CLAUDE.md` and watched it fail to take, you probably believe that file is the system prompt. It is not. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style goes in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works — that was measured, and MEASUREMENTS.md has the run.

Instruction alone does not fix the phrasing either. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant of itself.

The structural complaint has a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no instruction about phrasing was ever going to catch it.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test: a reader is shown a voice's own description of itself and picks which of two replies was written under it. The rate and the caveats are in MEASUREMENTS.md.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

`--setup` emits the output style and wires the two hooks into `~/.claude/settings.json` with absolute paths, then prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, and leaves every other key alone, including other tools' hooks on the same events. It refuses to touch a settings file that does not parse. `--setup --dry-run` prints what it would change and writes nothing.

The step left is the menu: pick the style under `/config` -> Output style. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

If you would rather nothing wrote to your settings, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` and stops there, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. The conversation you are in now will not change. The card goes here because an output style sits at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why it lands and why it did not land through a hook.

The hooks are the measurement half, and the voice works without either of them. `Stop` scores the reply once it is written. `UserPromptSubmit` restates the rules that have actually been firing, which a file written once cannot do.

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

These commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled into the binary. `--inject` remains for anyone who wants the superseded delivery, and stands down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written, with no Python and no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`, `--rules` takes a path from anywhere, and `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, so a session cannot quietly write in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, and it is what the discrimination run measured. What a card changes is the sound of the sentences, the first of the two axes set out above; the structure rules stay where they are, minus whatever `@gate` and `@shape` reach.

[demo/README.md](demo/README.md) is the shortest way to see the effect without writing a card, since every file under `demo/` is this README written again from a different one, same prompt and same facts.

`card/demo/handoff.effigy` is the exception in that directory. It is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose — and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so it turns up in your card list, and it is not a register to write in.

## What the hooks add

`Stop` scores the reply just written, appends which rules fired to the session's rolling state, and writes one record per violation to the log.

`UserPromptSubmit` reads that rolling state — the state, not the violations log — and injects the card items gated on what has been firing, naming the counts. An item gated on `rule_labelled_opening` appears once that rule has been firing and not before. When the session has no history yet, it falls back to the standing `CONTINUE TEST`. It stays quiet until the last injection has aged past `--refresh-every`.

So the mid-session text is chosen from measured output rather than fixed in advance, which a pasted `CLAUDE.md` cannot do. That is one mechanism and not a guarantee: the A/B in this repo does not separate the refresher from no refresher.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a `warn` action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and it is the one to reach for if this is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures lemma frequency against a baseline over real transcripts, reports what you have been leaning on lately, and leaves the judgement to you — awareness rather than prohibition, as its own README puts it. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose, on different hooks with no shared state, and running both is reasonable.

[caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses agent replies to cut output tokens. cope shapes prose, basanite tracks vocabulary, caveman shortens; a reader who wants fewer tokens rather than different structure should go there.

## The rules

Voicing, the half a card owns:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, the same whichever card is loaded:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it, no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A `POSTPROC` pattern matches a span of text, so it can only ever describe wording. Every voicing rule that needed more than a pattern had to be written in Go beside the structure rules, which is why cope's shipped card carries three `POSTPROC` rules and no more. Anyone expecting a long list of banned phrases should know that the list lives in basanite on purpose.

The structure rules vary in exactly one way, and by who is going to read the turn rather than by which card is loaded. The `interactive` lane covers any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The `loop` lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. It drops `ask_not_last`, `forked_end`, `dangling_end`, and `buried_decision`, since nobody is reading yet: a report that correctly names what it left open and stops would fail three of those, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — `unverified_done` and `loop_ask` — because a report saying the work is done has to say what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode 0600, holds one JSON record per violation carrying the matched text and about 70 characters either side. The log quotes your replies back at you. `--log` with an empty value turns it off.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode 0600, is an empty file whose mtime is the refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode 0600, holds the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so what is enforced and what is injected cannot drift apart.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file: the always-on rules and the evidence-gated ones are printed by different paths, no code path renders their union, so a card can hold more `NEVER` rules in total than the budget and still be healthy.

Two forms let a card speak to the structure rules, one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`

predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

The id in `@gate` has to be one the gate already has, and the id in `@shape` must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in a sentence the binary supplies. The shipped card asserts one of its own: `short_close`, the close is one or two sentences, past which the summary has become another section.

The `60` in the handoff card's `readable_cold` is measured, not picked. Across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and it is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Watch hits per character rather than the share of turns hit, because that second number tracks how long the turns were. MEASUREMENTS.md has the rates.

## Known limits

`labelled_opening` is not a tagger, so it decides what counts as a verbless opener by heuristic and will misjudge some. `ask_not_last` says nothing about the order of several asks once they are all in the last block.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written. That makes it a description of the output and not a judgement of it. The judgement lives in the discrimination test, and that covers voicing only.

The largest limit sits in the divide between the two axes described at the top. A card can decline a built-in rule and assert one of its own, and the vocabulary it asserts in counts words and sentences and asks whether a block poses a question. `clause_symmetry` is not writable in that vocabulary and is not meant to be, so the compiled rules remain the only home for a check of that kind, and a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere to put it.

Both directions are the card marking its own homework. A decline lowers that card's score and an assertion raises it, so read either with the reason attached, which is why the syntax requires one. MEASUREMENTS.md has the rest.

## Layout

- `card/claude_voice.effigy` — the shipped card, in effigy notation
- `card/rules.json` — generated from it; embedded in the binary
- `card/demo/` — other cards, each written to sound like something else
- `cmd/cope-gate/` — the hook binary
- `internal/scan/` — the structure rules, the card's regex rules, and the card renderer
- `internal/effigy/` — the `.effigy` reader, so a card is usable as written
- `internal/transcript/` — Claude Code JSONL reader, and which lane a turn was written in
- `replay/` — the blind-pairs and discrimination harnesses, and their own README
- `demo/` — this README written again under each demo card
- `tools/` — card compiler, effigy-backed scorer, cross-project sweep
- `MEASUREMENTS.md` — what was run, on how much text, and what it said

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
