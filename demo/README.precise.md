# cope

cope does two jobs: the card file you write fixes how replies sound, down to the register, which changes the moment you point it at a different file; the compiled gate checks how a reply is structured, the same way every time — a cope is the upper half of a foundry mould, the half carrying the shape being cast into.

Start at [demo/README.md](demo/README.md): every file in that directory is this README written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card that instructs every tic this model is measured to have; it makes the point in a single glance.

## The two things

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. The card holds all of it — VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC. Point the gate at a different card and the register goes with it; that is the half with a measured result behind it. Structure is the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives `continue` something to refer to, whether a claim that the work is done carries anything that could have shown it. Those rules are compiled into the binary and fire the same whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing fault; a reply that names an open problem in its last paragraph and then stops, leaving `continue` nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

A card reaches into the structure half in two directions, both declared in the card header, one per line.

`@gate <rule_id> off — <why>` declines a built-in rule. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch; without the decline, a card whose VOICE asked for something a built-in rule catches was marked down for obeying itself.

`@shape <id>: <selector> <predicate> — <why>` asserts a rule the gate then checks. `card/demo/handoff.effigy` asserts `readable_cold` — last paragraph words <= 60 — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90. Before `@shape`, a card's own commitment about how a reply ends had nowhere to be checked.

The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more. [MEASUREMENTS.md](MEASUREMENTS.md) has the runs and the reasons their numbers do not carry more than that.

## The problem

Instruction alone does not fix the phrasing. A global CLAUDE.md banning the "not A, it's B" flip is read every turn. The flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

The structural complaint has a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and it is not a phrasing habit an instruction could have banned.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick the style under `/config` → Output style.

`--setup` emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse; `--setup --dry-run` prints what it would change and writes nothing.

For anyone who would rather not have their settings written to: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

A style is read once at session start. A new selection or a re-emitted card applies at the next session or after `/clear`. An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here.

The hooks are the measurement half. The voice works without either of them.

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

Stop scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. UserPromptSubmit restates mid-session the card items gated on what has actually been firing, naming the counts, which a file written once cannot do.

Those commands assume `go install`'s target directory is on PATH in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. `cope-gate --inject` remains as the superseded delivery for anyone who wants it, off by default, standing down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written, with no Python and no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, and it is what the discrimination run measured. What a card changes is the voicing axis, as **The two things** sets out; the numbers are in [MEASUREMENTS.md](MEASUREMENTS.md).

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README written again from a different card, same prompt and same facts.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered — `make cards` installs it with the rest, so a reader listing their cards will find it there.

## What the hooks do differently

Stop scores the reply and records which rules fired. UserPromptSubmit reads that record — not the violations log — and injects only the items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. That is one mechanism and not a guarantee; the A/B in this repo does not separate the refresher from no refresher. It falls back to the standing CONTINUE TEST when the session has no history, and stays quiet until the last injection has aged past `--refresh-every`. SessionStart `--inject` is superseded and off by default: measured 2026-08-03, a card asking for a bolded label on every paragraph and an emoji on every heading produced, through that hook, prose indistinguishable from no card at all.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and is the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures lemma frequency against a baseline over real transcripts, reporting what you have been leaning on lately and leaving the judgement to you — awareness rather than prohibition, as its own README puts it. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable. [caveman](https://github.com/JuliusBrussee/caveman), by a different author, is a third axis again: it compresses agent replies to cut output tokens, so a reader wanting fewer tokens rather than different structure should go there instead.

## The rules

Voicing:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat
- `apology` — the reply performs contrition instead of stating the correction and moving on
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for
- `announced_length` — the reply announces its own length rather than cutting it
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself

A POSTPROC pattern matches a span of text, so it can only ever describe wording. Every voicing rule that needed more than a pattern had to be written in Go beside the structure rules, which is why the shipped card carries 0 POSTPROC rules. A reader who came here expecting a long list of banned phrases should know that the list lives in basanite on purpose.

The structure rules vary in one way, and it is not by card. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. It drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`: nobody is reading yet, a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | (empty) | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report the arms; - reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also $COPE_CARD |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; - reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as --display would rewrite it |
| `--dry-run` | `false` | with --setup, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to ~/.claude/output-styles as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default ~/.claude/output-styles) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render -render-arm against |
| `--render-lane` | (empty) | render -render-arm as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this .effigy or .json file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | 0600 | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | 0600 | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | 0600 | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

The log quotes replies back.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced and the injected rules cannot drift.

The NEVER budget is 10, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file: SessionStart prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be healthy. The authoritative list of rules that really are discarded unrendered is the over-budget report, which is empty when the card is healthy.

Card-authored rules take two forms, one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

For `@gate` the id has to be one the gate has; for `@shape` it must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores every assistant turn in a transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project. Watch hits-per-character rather than the share of turns hit; that second number tracks how long the turns were. Rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks. The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written, so it describes the output and does not judge it; the judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one **The two things** named. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are the card marking its own homework: a decline lowers that card's score, an assertion raises it, and both are worth reading with the reason attached, which is why the syntax requires one. [MEASUREMENTS.md](MEASUREMENTS.md).

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
