# cope

cope ships one opinionated card — on from install, and the whole product for most readers — that scores how a reply sounds and how it is shaped as two separate jobs, and keeps that card in a file you can edit or swap, a cope being the upper half of a foundry mould, the half carrying the shape being cast into.

Read [demo/README.md](demo/README.md) first: every file there is this page written again from a different card, same prompt and same facts, so the card is the only thing that varies, and two of them read against each other show what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card instructing every tic this model is measured to have, makes the point in one glance and is deliberately hard going.

## The two things

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. It is the card, entirely — VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC — so swapping the card swaps every word of it. Structure is the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Structure is compiled into `internal/scan` and runs the same whichever card is loaded, varying only by lane. A sentence reaching for the balanced two-beat is a voicing fault; a reply that names an open problem in its last paragraph and then stops, leaving "continue" nothing to refer to, is a structural one. One reply can be clean on one axis and bad on the other, which is the reason to keep them apart. Voicing is the half with a measured result behind it: the blind discrimination test.

A card reaches into the structure half in two directions. `@gate <rule_id> off — <why>`, one per line in the card header, declines a built-in rule — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close that those two rules catch, and a card marked down for obeying itself measures nothing. `@shape <id>: <selector> <predicate> — <why>`, same place, asserts one — `card/demo/handoff.effigy` asserts `readable_cold`, `last paragraph words <= 60`, because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checked that; the 60 is measured, 33 words at the median across 43,155 assistant replies and 56 at p90. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, so a check outside that vocabulary and outside a POSTPROC regex has nowhere to sit. Rates and their caveats are in [MEASUREMENTS.md](MEASUREMENTS.md).

## What an instruction cannot do

A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style sits in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, word for word unchanged, is most of why cope works. That was measured — [MEASUREMENTS.md](MEASUREMENTS.md).

Instruction alone does not fix the phrasing. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip appeared twice in the session that built this while the ban was the topic. Naming a surface form pushes the move into a variant. That is the voicing side.

The structural complaint has a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no instruction could have banned it as a phrasing habit.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. Rate and caveats in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick it under `/config` → Output style.

`--setup` emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left. It backed the settings file up first, added only what was missing so a second run changes nothing, left every other key alone including other tools' hooks on the same events, and refuses a settings file that does not parse; `--setup --dry-run` prints what it would change and writes nothing.

By hand, for anyone who would rather not have their settings written to: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

A style is read once at session start. A new selection or a re-emitted card applies at the next session or after `/clear`.

An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here.

The hooks are the measurement half, and the voice works without either of them:

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

Stop scores the reply just written. UserPromptSubmit restates mid-session the rules that have actually been firing, which a file written once cannot do. SessionStart `--inject` is superseded and off by default, kept for anyone who wants the old delivery, and it stands down on its own when a cope output style is active.

Both commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, since `card/rules.json` is committed and compiled in.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written — no Python, no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`, `--rules` takes a path from anywhere, and `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one. Read `card/demo/lecturer.effigy` first: it differs from the shipped card on register alone, and it is what the discrimination run measured. What a card changes is the voicing axis, described under "The two things" above.

[demo/README.md](demo/README.md) is the shortest way to see what a card does without writing one — every file under `demo/` is this README written again from a different card, same prompt and same facts, so the voice is the only variable.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant for `make pairs` against the full card rather than for rendering. `make cards` installs it with the rest.

## What the hooks do differently

Stop scores the reply just written, appending which rules fired to the session's rolling state and one record per violation to the log. UserPromptSubmit reads that state — not the violations log — and injects only the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. A pasted `CLAUDE.md` cannot do that. One mechanism, no guarantee: the A/B run in this repo does not separate the refresher from no refresher. `--setup` does not wire SessionStart `--inject`.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round and is the one to reach for when this one is too blunt: it measures lemma frequency against a baseline over real transcripts and reports what you have been leaning on lately, leaving the judgement to you, where cope bans — a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. Its own README calls that awareness rather than prohibition. Which one fits is a question about mood: a heavy hand for a habit that is annoying you today, a moving measurement for watching drift instead of legislating it. They compose — different hooks, no shared state — and running both is reasonable.

[caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies to cut output tokens: a third axis. A reader after fewer tokens rather than different structure should go there.

## The rules

Voicing, shipped-card POSTPROC first:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A POSTPROC pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why cope's shipped card carries three of them; the `precise` card this page is written from carries none. A reader expecting a long list of banned phrases should look in basanite, where that list lives on purpose.

The structure rules vary in one place, by who is going to read the turn. Interactive is any turn that is not a loop turn, and somebody is waiting at a terminal where the ending is where they decide what happens next. Loop is a turn whose opening prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `forked_end`, `dangling_end`, and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | false | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | (empty) | comma-separated arms to rotate through, implying `-ab` (default inject,hold; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report the arms; `-` reads the default path |
| `--author-docs` | false | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | false | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | false | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | false | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | false | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | false | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | false | with `--setup`, print what would change and write nothing |
| `--inject` | false | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | 0.35 | flag paragraph-length coefficient of variation below this |
| `--output-style` | false | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | 30m0s | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | false | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | false | emit the output style and wire the hooks, then print the one step left |
| `--version` | false | print version and exit |

## What lands on disk

All three are mode 0600 under `$XDG_STATE_HOME/cope`.

- `violations.jsonl` — one JSON record per violation, carrying the matched text and about 70 characters either side. The log quotes replies back. `--log` empty disables it.
- `refresher-<session-id>` — an empty file whose mtime is the refresher clock.
- `session-<session-id>.json` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json`, and `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The NEVER budget is 10, and anything over it is reported at load rather than dropped silently. It is charged against each injection separately rather than against the card file: SessionStart prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be healthy.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

`@gate` takes an id the gate has; `@shape` takes one that does not collide with a rule id. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, since a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project. Watch hits per character; the share of turns hit tracks how long the turns were. Rates in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks. Hits are roughly four fifths structure, and the A/B run found that fraction tracks what a reply was for rather than how it was written, so it describes the output rather than judging it; the judgement is the discrimination test, and that covers voicing only.

The largest limit is the one under "The two things", as it now stands. A card can decline a built-in rule and write one of its own, in a vocabulary that counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are the card marking its own homework — a decline lowers that card's score, an assertion raises it — which is why the syntax requires the reason. See [MEASUREMENTS.md](MEASUREMENTS.md).

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
