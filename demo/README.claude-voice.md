# cope

cope ships one opinionated writing card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — that is on as soon as you install it and checks two different jobs, how a reply sounds and how it is shaped, and that lives in a file you can edit or swap for another when the shipped voice is not the one you want.

Every file in [demo/README.md](demo/README.md) is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed; reading two of them against each other shows what a card does faster than the rest of this page explains it. Start with [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have, which makes the point in a single glance.

## The two things it checks

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. All of it lives in the card, so swapping the card swaps every word of it, and this is the half with a measured result behind it. Structure is the shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Structure is compiled into the binary and runs the same whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing problem; a reply that names an open problem in its last paragraph and then stops is a structural one. The same reply can be clean on one axis and bad on the other, which is why they are kept apart.

A card reaches into the structure half in two directions. It can decline a built-in rule with `@gate <rule_id> off — <why>`, one per line in the card header, which exists because a card whose VOICE block asked for something a built-in rule catches was being marked down for obeying itself: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason. It can also state a rule of its own with `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold — last paragraph words <= 60` because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule covers that. This card asserts `short_close`, one or two sentences, because past that the summary has become another section. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing beyond that, so a check like `clause_symmetry` is not writable in it. [MEASUREMENTS.md](MEASUREMENTS.md) has the runs and the reasons their numbers do not carry further than they do.

## The problem

Instruction alone does not fix the phrasing. A global CLAUDE.md banning the "not A, it's B" flip is read every turn, and the flip still appeared twice in the session that built this tool, while the ban was the topic of the session. Naming a surface form pushes the move into a variant the wording of the ban does not cover, which is the argument for scoring output rather than only asking.

The structural complaint has a different cause. A reply that names an open problem in its closing paragraph and stops leaves the reader nothing to answer, so they re-read the block to find the decision hiding in it, or reply "continue" and get work nobody wanted. That costs a round trip, and no phrasing instruction could have banned it, because the sentences are fine.

The flip is an anecdote about one rule. What the voicing claim rests on is the blind discrimination test: a reader is shown a voice's own description of itself and two replies, and picks which one was written under it. The rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick the style under `/config` → Output style. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

`--setup` emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up before writing, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `cope-gate --setup --dry-run` prints what it would change and writes nothing.

If you would rather not have your settings edited, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` on its own, and `COPE_CARD=<name>` in front of it emits a different one. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`; the running conversation will not change.

The card goes here because an output style lands at the end of the system prompt and the harness re-reminds the model of it during the conversation, which it does not do for a hook message.

The hooks are the measurement half:

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

Stop scores the reply just written and records which rules fired. UserPromptSubmit reads that record and restates mid-session the items gated on what has actually been firing, naming the counts, which a file written once at session start cannot do. The voice works with neither of them installed.

Those commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, since `card/rules.json` is committed and compiled into the binary. There is also `--inject`, the superseded delivery that prints the card as one turn-zero message, kept for anyone who wants it and standing down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written, with no Python and no effigy checkout. Drop a card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, so a session never writes in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, which is what the discrimination run measured. What a card can change is the voicing axis, described above under "The two things it checks".

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md): every file under `demo/` is this README written again from a different card, same prompt and same facts, so the voice is the only variable between them.

`card/demo/handoff.effigy` is the exception in that directory. It is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose — and it is meant to be run through `make pairs` against the full card rather than rendered, though `make cards` installs it with the rest.

## What the hooks do differently

Stop scores the reply that was just written, appends which rules fired to the session's rolling state, and writes one record per violation to the log. UserPromptSubmit reads that rolling state, not the violations log, and injects the card items gated on what has been firing, naming the counts; it falls back to the standing CONTINUE TEST when the session has no history, and stays quiet until the last injection has aged past `--refresh-every`. So the mid-session text is chosen from measured output, which is the one thing a pasted CLAUDE.md cannot do. It is a mechanism and not a guarantee — the A/B in this repo does not separate the refresher from no refresher. `--inject` is superseded and off by default; `--setup` does not wire it.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label because three of its blocks already do what a prose gate needs. POSTPROC is regex rules with a `warn` action applied after generation. WRONG holds an anti-pattern beside its replacement. TEST holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and is the one to reach for if this one is too blunt: cope bans, so a rule fires or it does not and the register is fixed the moment you pick a card, while basanite measures lemma frequency against a baseline over real transcripts and reports what you have been leaning on lately, leaving the judgement to you. They compose — different hooks, no shared state — and running both is reasonable. [caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses agent replies to cut output tokens, which is where to go if the goal is fewer tokens rather than different structure.

## The rules

Voicing:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The grouping cuts across the implementation, and that is the part a reader can use. A POSTPROC pattern matches a span of text, so it can only ever describe wording; every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the card ships only 3 POSTPROC rules. Anyone expecting a long list of banned phrases should know that the list lives in basanite on purpose, where it is a measurement against a baseline rather than a prohibition.

The structure rules do vary in one way, by who is going to read the turn. The interactive lane is any turn that is not a loop turn, chosen because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. It drops `ask_not_last`, `forked_end`, `dangling_end`, and `buried_decision`, because nobody is reading yet: a report that correctly names what it left open and stops would fail three of those, and a question in it lands in a log where the next iteration reads it as an instruction to itself. In their place it adds `unverified_done` and `loop_ask`, so a report saying the work is done has to say what it ran.

## Flags

| flag | default | what it does |
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

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | 0600 | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | 0600 | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | 0600 | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

The violations log quotes your replies back. Set `--log ""` to disable it.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart.

The NEVER budget is 10, charged against each injection separately rather than against the card file: SessionStart prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be healthy. Anything genuinely over budget is reported at load rather than dropped silently, and that list is empty for the shipped card.

A card can also write structure rules of its own, in two forms, both in the card header and one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is fixed:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

A `@gate` id has to be one the gate already has, and a `@shape` id must not collide with one; a wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies. `card/demo/handoff.effigy` picked its `readable_cold` threshold from data: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

## Calibrating

`cope-gate --backfill` scores every assistant turn in a whole session transcript at once, and that is how these rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not one per project. Watch hits-per-character rather than the share of turns hit, since that second number tracks how long the turns were. The rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

`labelled_opening` is a heuristic over sentence shape and not a tagger, so it will miss a label that happens to contain a verb. `ask_not_last` says nothing about the order of several asks; it only notices one sitting above continuing prose.

Roughly four fifths of all hits are structural, and the A/B run found that four fifths tracks what a reply was for rather than how it was written, so the hit rate describes the output and does not judge it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one under "The two things it checks", as it now stands: a card can decline a built-in rule and assert one of its own, and the vocabulary it asserts in counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, so a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are the card marking its own homework — a decline lowers that card's score, an assertion raises it — which is why the syntax requires a reason attached to each, and why both are worth reading before trusting a score. [MEASUREMENTS.md](MEASUREMENTS.md) has the rest.

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

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
