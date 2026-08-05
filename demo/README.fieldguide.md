# cope

cope. A card of prose instructions that ships on and needs no editing, checking two separate jobs — what a reply sounds like and how it is shaped — and held in a file, so it can be edited or swapped; the name is the founding term, a cope being the upper half of a foundry mould, the half carrying the shape being cast into.

Read [demo/README.md](demo/README.md). Every file in that directory is this README written again from a different card — same prompt, same facts, so the card is the only thing that changed, and two of them read against each other show what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card instructing every tic this model is measured to have; it makes the point in a single glance and is deliberately hard going.

## The two things

Voicing. What the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. It lives in the card, entirely: VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC, so swapping the card swaps every word of it. This is the half with a measured result behind it. A sentence reaching for the balanced two-beat, two clauses of near-equal length repeating a content word across the joint, is a voicing problem.

Structure. The shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Compiled into the binary, the same rules whichever card is loaded. A reply naming an open problem in its last paragraph and then stopping, leaving "continue" nothing to refer to, is a structural problem. The same reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

A card reaches into the structure half in two directions. `@gate <rule_id> off — <why>`, one per line in the card header, declines a built-in rule; it exists because a card whose VOICE asked for the balanced landing was marked down by `clause_symmetry` for obeying itself. `@shape <id>: <selector> <predicate> — <why>`, likewise one per line, states a structural rule of the card's own; it exists because `card/demo/handoff.effigy` commits to a last block a reader can enter cold, and no built-in rule checks that. The vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — a card wanting a check outside both that vocabulary and a POSTPROC regex has nowhere to put it. [MEASUREMENTS.md](MEASUREMENTS.md) holds the run and the reasons its numbers do not carry more than that.

## The problem

Where the instruction sits. A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. Compare an output style, which is in the system prompt itself and which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works. Measured; the run is in [MEASUREMENTS.md](MEASUREMENTS.md).

What an instruction can say. Instruction alone does not fix the phrasing. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant. That is the voicing side.

The failure no instruction reaches. An ending that leaves the reader nothing to answer costs a whole round trip. Compare a phrasing habit, which a ban can at least name; this one is a property of where the decision landed, and no instruction about wording reaches it.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. Rate and caveats in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

Two commands and one menu choice.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> cope
```

`--setup` emitted the output style, backed up `~/.claude/settings.json`, and wired the two hooks into it with absolute paths, adding only what was missing — a second run changes nothing, every other key is left alone including other tools' hooks on the same events, and a settings file that does not parse is refused rather than repaired. `--setup --dry-run` prints what it would change and writes nothing.

By hand instead. `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. Pick it under `/config -> Output style`; the standalone `/output-style` command was removed in Claude Code v2.1.91, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear` — not in the conversation you are holding.

The card goes here rather than into a hook because an output style sits at the end of the system prompt and the harness keeps reminding the model of it.

The hooks, which no longer deliver the card:

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

`Stop` scores the reply just written, appending which rules fired to the session's rolling state and one record per violation to the log. `UserPromptSubmit` restates mid-session the card items gated on what has actually been firing, naming the counts — which a file written once cannot do. The voice works without either. These are the measurement half.

The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, since `card/rules.json` is committed and compiled in.

`--inject` is the superseded delivery, off by default, kept for anyone who wants it and standing down on its own when a cope output style is active.

## Writing in another voice

Cards. The gate reads `.effigy` directly, so a card is usable as written — no Python, no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name resolving to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

`card/demo/lecturer.effigy` is the one to read first. Differs from the shipped card on register alone, which is what the discrimination run measured. What a card changes is the voicing axis, described under *The two things* above; the rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

[demo/README.md](demo/README.md) is the shortest way to see what a card does without writing one: every file under `demo/` is this README again from a different card, same prompt and same facts.

`card/demo/handoff.effigy` is the exception in that directory — a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## What the hooks do differently

`Stop`. Scores the reply and records which rules fired, over a rolling 20-turn window.

`UserPromptSubmit`. Reads that record — not the violations log — and injects only the items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. That is the one thing a pasted `CLAUDE.md` cannot do. One mechanism, not a guarantee: the A/B in the repo does not separate the refresher from no refresher.

`SessionStart --inject` is superseded and off by default, and `--setup` does not wire it.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples — which is how a rule names a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) is the same problem answered the other way round, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures — lemma frequency against a baseline over real transcripts, reporting what you have been leaning on lately and leaving the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood: a heavy hand when a habit is annoying you today, a moving measurement when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable.

[caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses agent replies to cut output tokens. A third axis: cope shapes prose, basanite tracks vocabulary, caveman shortens. A reader wanting fewer tokens rather than different structure should go there.

## The rules

Voicing rules, checked either by a POSTPROC pattern in the card or in Go where a pattern was not enough.

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine. Warn.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead. Warn.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k. Warn.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure rules, compiled in.

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it. Interactive lane.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to. Interactive lane.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it. Interactive lane.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another. Interactive lane.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file. Loop lane.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself. Loop lane.

What the grouping implies. A POSTPROC pattern matches a span of text, so it can only ever describe wording, and every voicing rule needing more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries three POSTPROC rules and no long list of banned phrases; the list lives in basanite, on purpose. This page is written from the `fieldguide` card, whose own POSTPROC rules — `unrecorded_hedge` and `unpersuaded` — are that card's and not part of what a reader installs.

Lanes. The one place the structure rules vary — not by card, by who is going to read the turn. Interactive is any turn that is not a loop turn, chosen because somebody is waiting at a terminal and the ending is where they decide what happens next. Loop is a turn whose opening prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | (empty) | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report the arms; - reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
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

Three files under `$XDG_STATE_HOME/cope`, all mode 0600.

`violations.jsonl` holds one JSON record per violation, carrying the matched text and about 70 characters either side. The log quotes replies back.

`refresher-<session-id>` is an empty file whose mtime is the refresher clock.

`session-<session-id>.json` is the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The NEVER budget is 10, and anything over it is reported at load rather than dropped silently. Charged against each injection separately and not against the card file: `SessionStart` prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more NEVER rules in total than the budget and still be healthy.

`@gate <rule_id> off — <why>`, one per line in the card header. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch.

`@shape <id>: <selector> <predicate> — <why>`, one per line in the card header. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. `card/demo/handoff.effigy` asserts `readable_cold — last paragraph words <= 60`, because its peak asks the reader to re-enter cold and read the last block only. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

An id for `@gate` has to be one the gate has; an id for `@shape` must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, since a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

This card declines `labelled_opening` — an entry opens on the name of the thing and unpacks it, which is the form of the guide rather than a habit picked up from somewhere.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Hits per character is the metric worth watching. Compare the share of turns hit, which tracks how long the turns were. Rates in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

`labelled_opening` is not a tagger. It matches a shape and will take a paragraph that earned its opening fragment.

`ask_not_last` says nothing about the order of several asks. It knows only that one sits above the end.

Hit composition. Roughly four fifths structure by axis, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — a description of the output, not a judgement of it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one under *The two things*. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are the card marking its own homework: a decline lowers that card's score, an assertion raises it, and both are worth reading with the reason attached — which is why the syntax requires one. [MEASUREMENTS.md](MEASUREMENTS.md).

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | the shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in the binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | the hook binary |
| `internal/scan/` | the structure rules, the card's regex rules, and the card renderer |
| `internal/effigy/` | the .effigy reader, so a card is usable as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane a turn was written in |
| `replay/` | the blind-pairs and discrimination harnesses, and their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what was run, on how much text, and what it said |

---

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
