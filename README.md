# cope

One opinionated card, on from the moment it is installed and the whole of the tool for most readers, scoring two separate jobs — how a reply sounds and how it is shaped — out of a file that can be edited or swapped for another, a cope being the upper half of a foundry mould, the half carrying the shape being cast into.

Every file under [`demo/README.md`](demo/README.md) is this README written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it — [`demo/README.claude-maximal.md`](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have and hard going by design, makes the point in a single glance.

## Voicing and structure

Voicing. What the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. It lives in the card and nowhere else, so swapping the card swaps every word of it. A sentence reaching for the balanced two-beat is a fault on this axis, and the blind discrimination test is the instrument with a result behind it.

Structure. The shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives `continue` something to refer to, whether a claim that the work is done carries anything that could have shown it. Compiled into the binary, so the same rules run whichever card is loaded. A reply that names an open problem in its last paragraph and then stops is a fault on this axis. The same reply can be clean on one and bad on the other, which is the reason to keep them apart.

`@gate <rule_id> off — <why>`, one per line in the card header. A card declines a built-in rule it disagrees with. The form exists because a card whose `VOICE` block asked for something a built-in rule catches was marked down for obeying itself: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, the two rules that catch the balanced landing and the arriving close its own `VOICE` block asks for.

`@shape <id>: <selector> <predicate> — <why>`, one per line in the card header. A card states a structural rule of its own and the gate then checks it. The form exists because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — since its peak asks the reader to re-enter cold and read the last block only.

The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing beyond that. `clause_symmetry` is not writable in it, and `MEASUREMENTS.md` carries the runs behind this split along with the reasons its numbers do not carry more than they do.

## The problem

Placement. A global `CLAUDE.md` is not the system prompt, though a reader who has edited one and watched it fail to stick commonly believes otherwise. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. Compare an output style, which sits in the system prompt itself and which the harness re-reminds the model of as the conversation runs. Moving a single card between those two places, without changing a word of it, is most of why cope works; that was measured, and the run is in `MEASUREMENTS.md`.

Wording. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

Endings. A reply closing on an open problem with nothing in it to answer costs a whole round trip — a different complaint with a different cause, and not a phrasing habit an instruction could have banned.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it, and `MEASUREMENTS.md` has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

`go install` puts the binary on disk. `cope-gate --setup` does the rest of the install: it emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. `cope-gate --setup --dry-run` prints what it would change and writes nothing.

The step left is the selection: `/config` -> Output style, and pick the card there. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing goes into `.claude/settings.local.json` as `"outputStyle"`.

`cope-gate --output-style` is the by-hand route, for anyone who would rather not have their settings written to. It writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. An output style goes at the end of the system prompt, which is why the card lands there.

The two hooks, as the `hooks` block of `~/.claude/settings.json`:

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

`Stop` buys a score on the reply once it is written. `UserPromptSubmit` buys a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. The voice works without either of them; these are the measurement half.

`PATH`. The commands above assume `go install`'s target directory is on it in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, since `card/rules.json` is committed and compiled in. `--inject` remains as the superseded delivery for anyone who wants it, and stands down on its own when a cope output style is active.

## Writing in another voice

The `.effigy` reader is in the binary, so a card is usable as written, with no Python and no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

`card/demo/lecturer.effigy` is the one to read first. It differs from the shipped card on register alone, which is what makes it usable as a discrimination rival and what the discrimination run measured. What a card changes is the sounding half of the split above, and `MEASUREMENTS.md` holds the numbers.

Reading the rendered cards is the shortest way to see the effect without writing one, and [`demo/README.md`](demo/README.md) indexes them.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there.

## The hooks

`Stop`. It scores the reply just written, appends which rules fired to the session's rolling state, and writes one record per violation to the log.

`UserPromptSubmit`. It reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. Which items are eligible is set in the card, and a card that gates none of them gets the standing fallback instead. It falls back to that standing text when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`.

One mechanism, not a guarantee. The A/B in this repo does not separate the refresher from no refresher.

## Effigy notation

[`effigy`](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

Compare [`basanite`](https://github.com/justinstimatze/basanite), which measures lemma frequency against a baseline over real transcripts and leaves the judgement to the reader, where cope bans: a rule fires or it does not, and the register is fixed the moment a card is picked. That makes a card the wrong instrument for watching drift, and its own README calls the difference awareness rather than prohibition. They compose — different hooks, no shared state — and running both is reasonable. Compare also [`caveman`](https://github.com/JuliusBrussee/caveman), by a different author, which compresses agent replies to cut output tokens, and is where a reader wanting fewer tokens rather than different structure should go.

## The rules

Voicing:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form, since the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries three `POSTPROC` rules. A reader expecting a long list of banned phrases will find that list in `basanite`, on purpose.

The lanes are the one place the structure rules vary — by who is going to read the turn, not by which card is loaded. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check, since a report saying the work is done has to say what it ran.

## The flags

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
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600`. One JSON record per violation, carrying the matched text and about 70 characters either side. The log quotes replies back.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600`. An empty file whose mtime is the refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600`. The rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

`effigy` owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so the enforced and the injected rules cannot drift.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. It is charged against each injection separately and not against the card file: the always-on rules are rendered in one injection, the evidence-gated ones in another, and no code path renders their union — so a card may hold more `NEVER` rules in total than the budget and still be healthy. No rule in this card is discarded unrendered.

`@gate <rule_id> off — <why>` takes a rule id the gate already has. `@shape <id>: <selector> <predicate> — <why>` takes an id that must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

The `@shape` vocabulary:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

A declined rule still runs, and only this card's score drops it, so `--backfill` still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Hits per character is the metric to watch. Compare the share of turns hit, which tracks how long the turns were. `MEASUREMENTS.md` has the rates.

## Known limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written, so it describes the output rather than judging it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the voicing and structure split above. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere to put it.

Both directions are the card marking its own homework. A decline lowers that card's score and an assertion raises it, both are readable with the reason attached, and that is why the syntax requires one. `MEASUREMENTS.md` has the rest.

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
