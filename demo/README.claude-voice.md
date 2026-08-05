# cope

cope ships an opinionated card that a session writes in from the moment it is installed, plus a gate that scores what comes back on two separate counts — how the sentences sound and how the reply is shaped — with the card kept as a plain file you can edit or swap for another. (A cope is the upper half of a foundry mould, the half carrying the shape being cast into.)

Before anything else, [demo/README.md](demo/README.md): every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card that instructs every tic this model is measured to have, and it makes the point in one glance, though it is deliberately hard going.

## Two jobs

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. All of it lives in the card. Swap the card and every word of it changes. Structure is the shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. That half is compiled into `internal/scan` and runs the same whichever card is loaded.

A sentence reaching for the balanced two-beat — two clauses of near-equal length joined by a comma, repeating a word across the joint — is a voicing problem. A reply that names an open problem in its closing paragraph and then stops is a structural one, because the reader has nothing to answer. The same reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

A card reaches into the structure half in two directions, both written in the card header. `@gate <rule_id> off — <why>` declines a built-in rule: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its `VOICE` block asks for the balanced landing and the arriving close those two rules catch, and without the decline a card was marked down for obeying itself. `@shape <id>: <selector> <predicate> — <why>` states a structural rule of the card's own: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing further, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it.

## Where an instruction sits

A reader who has edited a global `CLAUDE.md` and watched it stop mattering usually believes that file is the system prompt. It is not. It arrives as one message attached to the first turn, and the conversation buries it under everything written afterwards. An output style goes in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works; that was measured, and the run is in [MEASUREMENTS.md](MEASUREMENTS.md).

Placement is only the first part. Instruction alone does not fix the phrasing: a global `CLAUDE.md` banning the "not A, it's B" flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic of conversation. Naming a surface form pushes the move into a variant of itself.

The structural complaint is a different thing with a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and it is not a phrasing habit an instruction could have banned.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`cope-gate --setup` emits the output style and wires the two hooks into `~/.claude/settings.json` with absolute paths, then prints the one step left. It backs the settings file up before writing, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. `--setup --dry-run` prints what it would change and writes nothing.

Then pick the style under `/config` -> Output style. The entry to look for is named `claude_voice`, which is the shipped card's id; a reader scanning that menu for something called cope will not find it. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way in, and the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

If you would rather not have your settings written to, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` on its own, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. Nothing changes in the conversation you are currently in.

The card goes here rather than into a hook because an output style sits at the end of the system prompt and the harness re-reminds the model of it during the conversation.

The hooks are separate, and the card no longer arrives through one. This is the `hooks` block of `~/.claude/settings.json`:

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

`Stop` buys you a score on every reply after it is written. `UserPromptSubmit` buys you a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. The voice works without either of them; these are the measurement half. There is also `--inject`, the superseded delivery that printed the card as a `SessionStart` prompt, kept for anyone who wants it and standing down on its own when a cope output style is active.

Both commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled into the binary.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written, with no Python and no effigy checkout. Drop a card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session quietly writing in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, and it is the card the discrimination run measured against. What a card can change is the voicing half described under Two jobs above.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README written again from a different card against the same prompt and the same facts, so the voice is the only variable.

`card/demo/handoff.effigy` is the exception in that directory. It is a hypothesis rather than a voice: it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so anyone listing their cards will find it there and should know it is not a register to write in.

## What the hooks add

`Stop` scores the reply that was just written, appends which rules fired to the session's rolling state, and writes one record per violation to the log.

`UserPromptSubmit` reads that rolling state — the state file, not the violations log — and injects the card items gated on what has been firing, naming the counts. When the session has no history yet it falls back to the standing `CONTINUE TEST` reminder, and it stays quiet until the last injection has aged past `--refresh-every`.

So the mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot do. It is a mechanism and not a guarantee: the A/B in this repo does not separate the refresher from no refresher.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs. `POSTPROC` is regex rules with a `warn` action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and it is what to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures lemma frequency against a baseline over real transcripts, reports what you have been leaning on lately, and leaves the judgement to you. They compose — different hooks, no shared state — and running both is reasonable.

[caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses agent replies to cut output tokens. Anyone wanting fewer tokens rather than different structure should go there instead.

## The rules

Grouped by which of the two jobs they check, rather than by where they happen to be implemented.

**Voicing.** Three of these are `POSTPROC` patterns in the shipped card:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

The rest needed more than a pattern and are compiled in:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure.** All compiled, all the same whichever card is loaded:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The grouping tells you something you can use. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries three `POSTPROC` rules and not thirty. Anyone expecting a long list of banned phrases should know the list lives in basanite on purpose.

The structure rules vary in exactly one way, and it is not by card. It is by who is going to read the turn. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. That lane drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, because nobody is reading yet: a report that correctly names what it left open and stops would fail three of them, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — `unverified_done` and `loop_ask` — since a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | What it does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | `(empty)` | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
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
| `--render-lane` | `(empty)` | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | `(empty)` | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

| Path | Mode | Holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | `0600` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | `0600` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | `0600` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

The log quotes your replies back at you. That is what makes a violation reviewable, and it is worth knowing before you leave it running.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart.

The `NEVER` budget is 10. Anything over it is reported at load rather than dropped in silence. The budget is charged against each injection separately and not against the card file — one path prints the always-on rules and the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more `NEVER` rules in total than the budget and still be healthy.

Both card-authored forms go one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is exactly this and no more:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

For `@gate` the id has to be one the gate already has; for `@shape` it must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

A declined rule still runs. Only this card's score drops it, so a backfill still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and that is how the rules were chosen in the first place. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Watch hits per character rather than the share of turns hit. The second number mostly tracks how long the turns were. [MEASUREMENTS.md](MEASUREMENTS.md) has the rates.

## Known limits

`labelled_opening` is not a tagger. It matches the shape of an opener and cannot tell you whether the label was a good one.

`ask_not_last` says nothing about the order of several asks. It only knows whether a reply carried on past one.

The hit rate is roughly four fifths structure, and the A/B run showed that four fifths tracks what a reply was for rather than how it was written. It describes the output; it does not judge it. The judgement lives in the discrimination test, and that test covers voicing only.

The largest limit is the one named under Two jobs above, and here is where it now stands. A card can decline a built-in rule and can write a structural rule of its own, in a vocabulary that counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it.

Both directions are the card marking its own homework. A decline lowers that card's score and an assertion raises it. Both are worth reading with the reason attached, which is why the syntax requires one. [MEASUREMENTS.md](MEASUREMENTS.md) has the runs.

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

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
