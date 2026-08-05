# cope

An opinionated voice card, live once installed, checking two jobs at once — how a reply sounds and how it is shaped — and held in a file you can edit or swap; a cope is the upper half of a foundry mould, the half carrying the shape being cast into.

Start at [demo/README.md](demo/README.md): every file there is this README written again from a different card, same prompt and same facts, so the card is the only variable, and reading two against each other shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card instructing every tic this model is measured to have, makes the point in one glance, though it is deliberately hard going.

## Two axes

Voicing is what the sentences sound like: register, rhythm, diction, where flair is licensed. It lives in the card. Swap the card and every word of it swaps. That is the half with a measured result behind it.

Structure is the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to. It is compiled into the binary, identical whichever card is loaded.

A sentence reaching for the balanced two-beat is a voicing fault. A reply naming an open problem in its last paragraph and then stopping is a structural one. The same reply can be clean on one axis and bad on the other, which is why they are kept apart.

A card reaches into the structure half in two directions. `@gate <rule_id> off — <why>` declines a built-in rule, which exists because a card whose `VOICE` block asked for something a built-in rule catches was marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` on those grounds. `@shape <id>: <selector> <predicate> — <why>` asserts a structural rule the gate then checks, which exists because a card's own commitment about how a reply ends had nowhere to be checked — `card/demo/handoff.effigy` asserts `readable_cold`, `last paragraph words <= 60`, measured at 33 words median and 56 at p90 across 43,155 assistant replies. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing further; `clause_symmetry` is not writable in it. `MEASUREMENTS.md` holds the run.

## The problem

Where an instruction sits governs whether it survives. A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style sits in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without editing a word of it, is most of why cope works. Measured; see `MEASUREMENTS.md`.

Instruction alone does not fix phrasing. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

The structural complaint has a different cause. An ending that leaves the reader nothing to answer costs a whole round trip. No phrasing ban reaches it.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. Rate and caveats in `MEASUREMENTS.md`.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick the style under `/config` -> Output style.

`--setup` emits the output style and wires the two hooks into `~/.claude/settings.json` with absolute paths. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse; `--setup --dry-run` prints what it would change and writes nothing. To skip the settings edit, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. It goes here because an output style lands at the end of the system prompt and the harness re-reminds the model of it mid-conversation.

The hooks block of `~/.claude/settings.json`:

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

`Stop` buys a score on every reply after it is written. `UserPromptSubmit` buys a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. The voice works without either; these are the measurement half. The card no longer arrives through a hook, though `--inject` remains as the superseded delivery for anyone who wants it, standing down on its own when a cope output style is active.

Both commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout — `card/rules.json` is committed and compiled in.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written: no Python, no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another. Read `card/demo/lecturer.effigy` first: it differs from the shipped card on register alone, and it is what the discrimination run measured. What a card changes is voicing, as the two axes above set out. Numbers in `MEASUREMENTS.md`.

[demo/README.md](demo/README.md) is the shortest way to see what a card does without writing one — every file under `demo/` is this README again from a different card, same prompt, same facts.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## What the hooks do differently

`Stop` scores the reply just written, appends which rules fired to the session's rolling state, and writes one record per violation to the log.

`UserPromptSubmit` reads that rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts. It falls back to the standing `CONTINUE TEST` when the session has no history, and stays quiet until the last injection has aged past `--refresh-every`.

So the mid-session text is chosen from measured output rather than fixed in advance. One mechanism, not a guarantee: the A/B in the repo does not separate the refresher from no refresher.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way and is what to reach for when this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures — lemma frequency against a baseline over real transcripts — and leaves the judgement to you. Its own README calls that awareness rather than prohibition. They compose: different hooks, no shared state. Running both is reasonable. [caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies to cut output tokens — a third axis, and the right destination for a reader wanting fewer tokens rather than different structure.

## The rules

Voicing:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule reading the window rather than the reply, so it cannot fire until a session has history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so a bold opener is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A `POSTPROC` pattern matches a span of text, so it can only ever describe wording. Every voicing rule needing more than a pattern had to be written in Go beside the structure rules, which is why cope's shipped card carries three `POSTPROC` rules. A reader expecting a long list of banned phrases should know that list lives in basanite on purpose.

The structure rules vary in one dimension, and it is not the card. `interactive` is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. `loop` is chosen when the prompt opening the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. It drops `ask_not_last`, `forked_end`, `dangling_end`, and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | Does |
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

## On disk

| Path | Mode | Holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | 0600 | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | 0600 | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | 0600 | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose, only rule names and counts |

The log quotes replies back.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json`; `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. It is charged against each injection separately, not against the card file: the card-as-prompt path prints the always-on rules and the refresher prints the evidence-gated ones, and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be healthy. Rules genuinely discarded unrendered are listed at load; that list is empty when the card is healthy.

Card-authored rules, one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

`@shape` selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`.

A `@gate` id has to be one the gate has; a `@shape` id must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not one per project. Watch hits-per-character rather than share of turns hit; the second tracks how long the turns were. Rates in `MEASUREMENTS.md`.

## Known limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — a description of the output, not a judgement of it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one the two axes above set out. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it.

Both directions are the card marking its own homework: a decline lowers that card's score, an assertion raises it. Both are worth reading with the reason attached, which is why the syntax requires one. `MEASUREMENTS.md`.

## Layout

| Path | What |
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

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
