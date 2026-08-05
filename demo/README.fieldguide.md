# cope

`cope`. An opinionated voice card for Claude Code — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — on from the moment it is installed and the whole of the product for most readers, checking two separate things about a reply, how the sentences sound and how the reply is shaped, and held in a plain file that can be edited or swapped for another.

[`demo/README.md`](demo/README.md) is every file in that directory being this page written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it; [`demo/README.claude-maximal.md`](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have, makes the point in a single glance, and is deliberately hard going.

## Voicing and structure

Voicing. What the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card and nowhere else — `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES` and `POSTPROC` — so swapping the card swaps every word of it. A sentence reaching for the balanced two-beat is a voicing fault. Compare structure, which is compiled into the binary and identical whichever card is loaded. Voicing is the half with a measured result behind it, from the blind discrimination test.

Structure. The shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives `continue` something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. A reply that names an open problem in its closing paragraph and then stops is a structural fault, and that same reply can be clean on every voicing rule. Two axes rather than one because the two verdicts come apart.

Card-authored rules. Two directions, both written one per line in the card header. `@gate <rule_id> off — <why>` declines a built-in rule, and exists because a card whose `VOICE` asked for something a built-in rule catches was marked down for obeying itself; `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` on exactly that ground. `@shape <id>: <selector> <predicate> — <why>` asserts a structural rule of the card's own, and exists because a card's commitment about how a reply ends had nowhere to be checked; `card/demo/handoff.effigy` asserts `readable_cold` for that reason. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more, so a card wanting a check outside that vocabulary and outside a `POSTPROC` regex still has nowhere to put it. The discrimination run and the reasons its numbers do not carry further are in `MEASUREMENTS.md`.

## Three failures

Three parts, in this order. Where an instruction sits, then what an instruction can say, then the failure no instruction reaches.

Placement. A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. Compare an output style, which is in the system prompt itself and which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works — measured, with the run recorded in `MEASUREMENTS.md`.

Wording. Instruction alone does not fix the phrasing. A global `CLAUDE.md` banning the "not A, it's B" flip is read every turn. The flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

Endings. A different complaint with a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no instruction could have banned it as a phrasing habit.

The claim. The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it, and `MEASUREMENTS.md` carries the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then `/config` -> `Output style`, and pick the card.

`cope-gate --setup`. Emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — `cope-gate --setup --dry-run` prints what it would change and writes nothing.

By hand. `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one; that is the route for anyone who would rather not have a tool write to their settings.

The selection. Under `/config` -> `Output style`. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

A style is read once at session start. So a new selection, or a re-emitted card, applies at the next session or after `/clear`.

The reason for putting it there. An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation.

The hooks, as `~/.claude/settings.json`:

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

`Stop`. Buys a score on the reply after it is written. `UserPromptSubmit` buys a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. The voice works without either of them; these are the measurement half.

`PATH`. The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in. A hook that silently does nothing usually wants the absolute path instead.

A clone. `make install` builds it, and no effigy checkout is needed, because `card/rules.json` is committed and compiled in.

`--inject`. The superseded delivery, kept for anyone who wants it, and standing down on its own when a cope output style is active.

## Writing in another voice

The `.effigy` reader. Built into the gate, so a card is usable as written — no Python, no effigy checkout.

Where a card goes. `$XDG_CONFIG_HOME/cope/cards`, reached by `--card <name>` or by `COPE_CARD`. `--rules` takes a path from anywhere, and `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

`card/demo/lecturer.effigy`. The one to read first. It differs from the shipped card on register alone, which is what makes it the rival the discrimination run measured. What a card changes is the sounding half of the pair described at the top of this page, and `MEASUREMENTS.md` holds the numbers.

[`demo/README.md`](demo/README.md) is the shortest way to see what a card does without writing one: every file under `demo/` is this README written again from a different card, same prompt and same facts, so the voice is the only thing that varies.

`card/demo/handoff.effigy`. The exception in that directory — a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, and meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so it turns up in a card list without being a register to write in.

## Scoring and the refresher

`Stop`. Scores the reply just written, appends which rules fired to the session's rolling state, and appends one record per violation to the log.

`UserPromptSubmit`. Reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts. It falls back to the standing `CONTINUE TEST` when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`.

Gating. The mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot do. One mechanism, not a guarantee: the A/B in this repo does not separate the refresher from no refresher.

## effigy, off-label

`effigy`, at https://github.com/justinstimatze/effigy. A character-card notation for game NPCs, used here for prose. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

Compare `basanite`, at https://github.com/justinstimatze/basanite, which answers the same problem the other way round and is the one to reach for when this one is too blunt. cope bans — a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. `basanite` measures lemma frequency against a baseline over real transcripts, reports what you have been leaning on lately, and leaves the judgement to you; its own README calls that awareness rather than prohibition. Which one fits is a question about mood rather than correctness, and they compose: different hooks, no shared state.

Compare `caveman`, at https://github.com/JuliusBrussee/caveman, by a different author again — it compresses agent replies to cut output tokens. cope shapes prose, `basanite` tracks vocabulary, `caveman` shortens; a reader wanting fewer tokens rather than different structure goes there.

## The rules

Voicing, from the shipped card's `POSTPROC` block:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, compiled into `internal/scan`:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, all compiled:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The split between those two lists. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries three `POSTPROC` rules and no more. This page is written from `card/demo/fieldguide.effigy`, whose own three — `heading_headword`, `unrecorded_hedge`, `unpersuaded` — belong to that card rather than to the shipped one.

The long list of banned phrases. In `basanite`, on purpose, and measured against a baseline rather than fixed.

Lanes. The one place the structure rules vary, and by who is going to read the turn rather than by card. The `interactive` lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The `loop` lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | `(empty)` | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; `positive` is the third) |
| `--ab-report` | `(empty)` | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | `(empty)` | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | `(empty)` | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | `(empty)` | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | `MessageDisplay` entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | `(empty)` | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | `(empty)` | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | `(empty)` | render `-render-arm` as the given lane sees it: `interactive` (default) or `loop` |
| `--rules` | `(empty)` | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## State on disk

Three files, all `0600`.

`$XDG_STATE_HOME/cope/violations.jsonl` holds one JSON record per violation, carrying the matched text and about 70 characters either side. It quotes replies back at you.

`$XDG_STATE_HOME/cope/refresher-<session-id>` is an empty file whose mtime is the refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json` holds the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

`card/claude_voice.effigy`. effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json`, and `make check-rules` is what CI runs, so the enforced and the injected rules cannot drift.

The `NEVER` budget. Ten, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately rather than against the card file — the always-on rules and the evidence-gated ones are printed by different paths and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be healthy.

`@gate <rule_id> off — <why>`. One per line in the card header. The id has to be one the gate already has, and a wrong id is reported at load rather than ignored. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught.

`@shape <id>: <selector> <predicate> — <why>`. One per line in the card header. The id must not collide with a rule the gate has, and a collision is reported at load. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. A violation is reported in the card's own words rather than in any sentence the binary supplies.

The reason after the dash. Required in both forms, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

## Calibrating

`cope-gate --backfill`. Scores a whole session transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Hits per character. The metric to watch, rather than the share of turns hit — that second number tracks how long the turns were. The rates are in `MEASUREMENTS.md`.

## Known limits

`labelled_opening`. Not a tagger. It matches a shape and cannot tell an opener that earns its fragment from one that does not.

`ask_not_last`. Says nothing about the order of several asks, only that an ask sits above prose that continues.

The hit rate. Roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — a description of the output rather than a judgement of it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one the two axes at the top of this page describe. A card can decline a built-in rule and can write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live. Both directions are the card marking its own homework: a decline lowers that card's score and an assertion raises it, both read properly only with the reason attached, and that is why the syntax requires one. `MEASUREMENTS.md` has what was run.

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

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
