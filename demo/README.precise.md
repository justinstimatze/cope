# cope

An opinionated prose card, on once installed, checking two things a reply can get wrong — how it sounds and how it is shaped — and shipped as a file you can edit or swap; a cope is the upper half of a foundry mould, the half carrying the shape being cast into.

Start at [demo/README.md](demo/README.md): every file there is this README written again from a different card, same prompt and same facts, so the card is the only variable, and reading two against each other shows what a card does faster than this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card instructing every tic this model is measured to have, makes the point in one glance, though it is deliberately hard going.

## Two axes

Voicing is what the sentences sound like: register, rhythm, diction, where flair is licensed. It lives in the card, and swapping the card swaps every word of it. Structure is the shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether a claim that the work is done carries anything that could have shown it. Structure is compiled into `internal/scan`, identical whichever card is loaded.

A clause reaching for the balanced two-beat is a voicing fault. A reply naming an open problem in its last paragraph and stopping is a structural one. The same reply can be clean on one axis and bad on the other, which is why they are kept apart. Voicing is the axis with a measured result behind it.

A card reaches into the structure axis in two directions. `@gate <rule_id> off — <why>`, one per line in the card header, declines a built-in rule: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` because its `VOICE` block asks for the balanced landing and the arriving close those rules catch, and a card marked down for obeying itself is a broken instrument. `@shape <id>: <selector> <predicate> — <why>` asserts a rule of the card's own: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checked that commitment. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more; `clause_symmetry` is not writable in it and is not meant to be. See [MEASUREMENTS.md](MEASUREMENTS.md).

## The problem

Where an instruction sits. A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style is in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works. Measured — see [MEASUREMENTS.md](MEASUREMENTS.md).

What an instruction can say. A global `CLAUDE.md` banning the "not A, it's B" flip is read every turn, and the flip appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

The failure no instruction reaches. An ending that leaves the reader nothing to answer costs a whole round trip. That is not a phrasing habit a ban could have named.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. Rate and caveats in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`cope-gate --setup` emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. `cope-gate --setup --dry-run` prints what it would change and writes nothing. To skip the settings edit: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

Then select it under `/config` -> Output style. The entry is named `claude_voice` — the shipped card's id, not the word cope. A reader scanning that menu for cope will not find it. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`.

An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here.

Then the hooks. The card no longer arrives through one; these are the measurement half, and the voice works without either.

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
}
```

`Stop` scores the reply just written. `UserPromptSubmit` restates the rules that have actually been firing, which a file written once cannot do. `PreToolUse` scores prose an external write is about to post.

These commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout — `card/rules.json` is committed and compiled in. `--inject` remains as the superseded delivery for anyone who wants it, and stands down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written — no Python, no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name resolving to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another.

Read `card/demo/lecturer.effigy` first: it differs from the shipped card on register alone, and it is what the discrimination run measured. A card changes the voicing axis and nothing else, per the axis split above. Numbers in [MEASUREMENTS.md](MEASUREMENTS.md).

The shortest route to seeing what a card does without writing one is [demo/README.md](demo/README.md) — every file under `demo/` is this README written again from a different card, same prompt and same facts.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## What the hooks add

`Stop` runs `cope-gate`, which scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` runs `cope-gate --refresher`, which reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts. It falls back to the standing `CONTINUE TEST` when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. So the mid-session text is chosen from measured output rather than fixed in advance. One mechanism, not a guarantee: the A/B in the repo does not separate the refresher from no refresher.

`PreToolUse` runs `cope-gate --pretool`, which scores the `description`, `body` or `content` field an external write is about to post, matched against the Linear save tools named in the settings block. Warn-only — it returns `additionalContext` and never a `permissionDecision`, so the call goes through and the model learns what the prose scored. It writes no session state, and it scores in the external lane.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and is the one to reach for if this one is too blunt: it measures lemma frequency against a baseline over real transcripts and leaves the judgement to you, where cope bans. Which fits is a question about mood — a heavy hand for a habit annoying you today, a moving measurement to watch drift rather than legislate it. They compose: different hooks, no shared state. [caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies to cut output tokens — a third axis, and the right destination for a reader wanting fewer tokens rather than different structure.

## The rules

Voicing, in the shipped card's `POSTPROC` block:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, compiled:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, all compiled:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A `POSTPROC` pattern matches a span of text, so it can only ever describe wording. Every voicing rule that needed more than a pattern was written in Go beside the structure rules, which is why the shipped card carries three `POSTPROC` rules. A reader expecting a long list of banned phrases wants basanite: the list lives in another tool on purpose.

Lanes are the one place the structure rules vary, and they vary by who will read the turn rather than by card. The interactive lane is any turn that is not a loop turn: somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet — a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran. The external lane drops the same four and adds nothing: a ticket has a reader and no ending they can answer, read days later by somebody who was not in the session, which is the condition every surviving rule was written for.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; `positive` is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | `MessageDisplay` entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--pretool` | `false` | `PreToolUse` entry: score the prose an external write is about to post, warn-only |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: `interactive` (default) or `loop` |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## On disk

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | `0600` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | `0600` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | `0600` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

The log quotes replies back.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json`; `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately, not against the card file: one path prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more `NEVER` rules in total than the budget and still be healthy. Nothing in the shipped card is currently discarded unrendered.

Both card-authored forms go in the card header, one per line:

    @gate <rule_id> off — <why>
    @shape <id>: <selector> <predicate> — <why>

The `@shape` vocabulary:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

A rule id has to be one the gate has for `@gate`, and must not collide with one for `@shape`. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

The `60` in `readable_cold` is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not one per project. Watch hits-per-character rather than the share of turns hit — that second number tracks how long the turns were. Rates in [MEASUREMENTS.md](MEASUREMENTS.md).

## Limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written. It is a description of the output, not a judgement of it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one the axis split above names, as it now stands: a card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. Both directions are the card marking its own homework — a decline lowers that card's score, an assertion raises it — which is why the syntax requires a reason. See [MEASUREMENTS.md](MEASUREMENTS.md).

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

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
