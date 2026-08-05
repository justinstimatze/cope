# cope

cope ships one opinionated card that is on from the moment it is installed — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — and it scores every reply on two separate jobs, how the sentences sound and how the reply is shaped; the card itself is a file, so it can be edited or swapped for another.

Every file under [demo/README.md](demo/README.md) is this README written again from a different card, the same prompt and the same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it — [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have, makes the point in a single glance, and is deliberately hard going.

## Voicing and structure

Voicing. What the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. It lives in the card, all of it, and swapping the card swaps every word of it. A sentence reaching for the balanced two-beat, two clauses of near-equal length repeating a content word across the joint, is a voicing fault. This is the half with a measured result behind it.

Structure. The shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives `continue` something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Compiled into the binary, in `internal/scan`, and the same whichever card is loaded. A reply that names an open problem in its last paragraph and then stops is a structural fault: `continue` has nothing to refer to. The same reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

`@gate <rule_id> off — <why>`, one per line in the card header. A card declines a built-in rule it disagrees with. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its `VOICE` block asks for the balanced landing and the arriving close those two rules catch, and a card was being marked down for obeying itself.

`@shape <id>: <selector> <predicate> — <why>`, one per line in the same header. A card states a structural rule of its own and the gate checks it. `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60`. Its peak asks the reader to re-enter cold on the closing block alone, and no built-in rule sees whether that can be done. The 60 is measured: 33 words at the median and 56 at p90, across 43,155 assistant replies.

The boundary. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex has nowhere to put it. `MEASUREMENTS.md` has the run and the reasons its numbers do not carry more than that.

## Three failures

Delivery. A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after it. An output style sits in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works, and `MEASUREMENTS.md` has the run.

Phrasing. Instruction alone does not fix it. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

Shape. A different complaint with a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no ban on a phrase reaches it.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it, and `MEASUREMENTS.md` has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`cope-gate --setup`. Does the whole install: emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. `cope-gate --setup --dry-run` prints what it would change and writes nothing.

By hand, for anyone who would rather not have their settings written to: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

The menu entry. Pick it under `/config` -> `Output style`. The entry to look for is named `claude_voice`, which is the shipped card's id and not the word cope, so a reader scanning that menu for something called cope will not find it. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`.

An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through the hook. Compare `--inject`, the superseded delivery, kept for anyone who wants the old arrangement and standing down on its own when a cope output style is active.

The hooks block of `~/.claude/settings.json`, which `--setup` writes:

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

`Stop` buys a score on the reply just written. `UserPromptSubmit` buys a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. The voice works without either; these two are the measurement half.

The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in.

## Another card

The reader. The gate reads `.effigy` directly, so a card is usable as written, with no Python and no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

`card/demo/lecturer.effigy`. The one to read first — it differs from the shipped card on register alone, which is what the discrimination run measured. A card changes voicing, the first of the two axes above, and `MEASUREMENTS.md` has the numbers.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README written again from a different card, same prompt and same facts, so the voice is the only thing that varies between them.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and can leave it alone as a register.

## The hooks

`Stop`. Scores the reply just written, then records the rules that fired in the session's rolling state and one line per violation in the log.

`UserPromptSubmit`. Reads the rolling state — the violations log is a separate file and this hook does not touch it — and injects the card items gated on what has been firing, naming the counts. Falls back to the standing `CONTINUE TEST` when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`.

Mid-session text chosen from measured output, which a pasted `CLAUDE.md` cannot do. One mechanism. The A/B in the repo does not separate the refresher from no refresher.

## Effigy notation

`effigy`, at https://github.com/justinstimatze/effigy. A character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

Compare `basanite`, at https://github.com/justinstimatze/basanite, which measures lemma frequency against a baseline over real transcripts and leaves the judgement to the reader — awareness rather than prohibition, and the wrong instrument for a card that has to say never. Compare `caveman`, at https://github.com/JuliusBrussee/caveman, a separate project by a different author that compresses agent replies to cut output tokens: a third axis, and the one to reach for when the want is fewer tokens rather than different shape.

## The rules

Voicing:

- `flip` (`POSTPROC`, warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (`POSTPROC`, warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` (`POSTPROC`, warn) — announces attention in advance instead of letting the observation earn it, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries three of them. A reader expecting a long list of banned phrases will find that list in `basanite`, on purpose.

Lanes. The structure rules vary in one way, and it is not by card. The interactive lane runs on any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane runs when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | Does |
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
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: `interactive` (default) or `loop` |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## On disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600`. One JSON record per violation, carrying the matched text and about 70 characters either side. This file quotes replies back.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600`. An empty file whose mtime is the refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600`. The rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. Rule names and counts only; no prose.

## Editing the card

`card/claude_voice.effigy`. The shipped card, in effigy notation, and effigy owns the `.effigy` grammar. `make rules` regenerates, and `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The `NEVER` budget. Ten. Anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file: the card injection carries the always-on rules, the refresher carries the evidence-gated ones, and no code path renders their union, so a card can hold more `NEVER` rules in total than the budget and still be healthy.

Both card-authored forms, one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`.

A rule id has to be one the gate already has for `@gate`, and must not collide with one for `@shape`. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs. Only this card's score drops it, so a backfill reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibration

`cope-gate --backfill`. Scores every assistant turn in a transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project. Hits per character is the number to watch; the share of turns hit tracks how long the turns were. `MEASUREMENTS.md` has the rates.

## Known limits

`labelled_opening` is not a tagger.

`ask_not_last` says nothing about the order of several asks.

The hit rate. Roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — a description of the output and not a judgement of it. The discrimination test is where the judgement lives, and it covers voicing only.

The largest limit is the one the section on voicing and structure named. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere to put it.

Both directions are the card marking its own homework: a decline lowers that card's score and an assertion raises it. Both read with the reason attached, which is why the syntax requires one, and `MEASUREMENTS.md` has the runs.

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
