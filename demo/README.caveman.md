# cope

Cope ships one opinionated card, on once installed, scoring two jobs — sound of sentences, shape of reply — and card sits in a file, so you edit it or swap it. Name comes from foundry: cope, upper half of mould, half carrying shape cast into.

Read [demo/README.md](demo/README.md). Every file there holds this page again, same prompt, same facts, card only change. Read two against each other. Faster than rest of page. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) shows it in one glance, card instructing every tic this model measures for. Hard going. Named, not recommended.

## Two things

Voicing: what sentences sound like. Register, rhythm, diction, what paragraph does with detail, where flair licensed. Voicing lives in card. Swap card, swap voicing. Blind discrimination test measures this half.

Structure: where decision sits, how reply ends. Compiled into binary. Same rules under every card.

One instance each. Sentence reaching for balanced two-beat — comma-joined clauses of near-equal length, content word repeated across joint — voicing problem. Reply naming open problem in last paragraph, then stopping, leaving "continue" nothing to refer to — structure problem. Same reply passes one axis, fails other. Keep them apart.

Card reaches into structure half two ways. `@gate <rule_id> off — <why>` declines built-in rule: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for balanced landing and arriving close those rules catch, and card obeying itself got marked down. `@shape <id>: <selector> <predicate> — <why>` states card's own structural rule: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks reader to re-enter cold and read last block only, and no built-in rule checks that. Boundary: @shape vocabulary counts words and sentences and asks whether block poses question. Nothing more. Card wanting check outside that vocabulary and outside a POSTPROC regex has nowhere to put it. Rates live in [MEASUREMENTS.md](MEASUREMENTS.md).

## Problem

Where instruction sits comes first. Most readers hold it backwards. Global CLAUDE.md reads like system prompt. Not system prompt. Arrives as one message on first turn. Conversation buries it under everything after. Output style sits in system prompt, and harness re-reminds model of it as conversation runs. Move one card between those two places, change no word, and most of cope works. Measured — see [MEASUREMENTS.md](MEASUREMENTS.md).

What instruction says comes second. Instruction alone fails on phrasing. Global CLAUDE.md banned not-A-but-B flip, read every turn. Flip appeared twice in session that built this, ban itself topic of that session. Name surface form, move slides into variant. Voicing side.

Third failure reaches no instruction. Ending leaves reader nothing to answer. Costs whole round trip. Not phrasing habit. No ban touches it.

Flip stands as anecdote about one rule. Claim rests elsewhere: blind discrimination test. Reader sees only a voice's own description of itself, picks which of two replies came from it. Rate and caveats in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick style: `/config` -> Output style.

`--setup` emits output style, wires two hooks into `~/.claude/settings.json` with absolute paths, prints one step left. Backs settings file up first. Adds only what missing, so second run changes nothing. Leaves every other key alone, other tools' hooks on same events included. Refuses to touch settings file that fails to parse. `--setup --dry-run` prints changes, writes nothing.

By hand, for anyone guarding their settings file: `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front picks another. Claude Code v2.1.91 removed standalone `/output-style`, so `/config` remains the way; `"outputStyle"` in `.claude/settings.local.json` sets same thing.

Style reads once at session start. New selection, or freshly emitted card, applies at next session or after `/clear`. Running conversation shows nothing.

Card goes into output style because output style sits at end of system prompt and harness re-reminds model of it mid-conversation. Hook delivery got no such reminder.

Hooks come after, and card no longer arrives through one:

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

Stop scores reply just written, appends fired rules to session's rolling state, appends one record per violation to log. UserPromptSubmit reads that state and injects card items gated on what fires, counts named. Voice works without both. These two form measurement half.

Commands assume `go install` target directory on PATH inside environment hook runs in. Hook doing nothing silently wants absolute path. Clone builds with `make install`, no effigy checkout needed: `card/rules.json` sits committed and compiles in.

`--inject` remains as superseded delivery, off by default, standing down on its own when cope output style goes active.

## Writing in another voice

Gate reads `.effigy` directly. Card runs as written. No Python, no effigy checkout. Drop card in `$XDG_CONFIG_HOME/cope/cards`, reach it with `--card <name>` or `COPE_CARD`. `--rules` takes path from anywhere. `make cards` installs demo set. Name resolving to nothing errors, and nothing injects — no session writes shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. Differs from shipped card on register alone. Discrimination run measured that pair. Card changes voicing, as the two axes above set out.

[demo/README.md](demo/README.md) shows what a card does without writing one: every file under `demo/` holds this README again from a different card, same prompt, same facts, voice the only variable.

One exception in that directory: `card/demo/handoff.effigy` states a hypothesis, not a voice. Keeps shipped card's handoff rules, drops everything about prose, and runs through `make pairs` against full card rather than rendering. `make cards` installs it, so a reader listing cards finds it there. Not a register to write in.

## What hooks add

Stop scores reply after model writes it, records which rules fired over rolling window. UserPromptSubmit reads that record and injects only items gated on firing rules, naming counts. Mid-session text comes from measured output, not from advance decision. Pasted CLAUDE.md cannot do that. One mechanism, no guarantee: A/B in repo does not separate refresher from no refresher. SessionStart `--inject` stays superseded and off by default.

## Why effigy

[effigy](https://github.com/justinstimatze/effigy) states a character-card notation for game NPCs, used here off-label. Three blocks do what a prose gate needs. POSTPROC holds regex rules with warn action applied after generation. WRONG holds anti-pattern beside its replacement. TEST holds named question with fail and pass examples, so a rule names a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers same problem other way round, and suits a reader finding cope too blunt. cope bans: rule fires or fails to, card says never, register fixes the moment you pick it. basanite measures — lemma frequency against baseline over real transcripts — reports what you lean on lately, leaves judgement to you. Its README calls that awareness over prohibition. Choice tracks mood over correctness: heavy hand when a habit annoys you today, moving measurement when you would rather watch drift than legislate it. They compose, different hooks, no shared state. Run both.

Third project, separate author: [caveman](https://github.com/JuliusBrussee/caveman) compresses agent replies to cut output tokens. Third axis. cope shapes prose, basanite tracks vocabulary, caveman shortens. Reader wanting fewer tokens rather than different structure goes there.

## Rules

Voicing rules. Three POSTPROC patterns in shipped card:

- `flip` (warn) — not-A-but-B flip in common surface forms, not-only-but-also included. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude top of it among every writer in study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, 25.6 per 1k, heaviest measured lean in this register. Say what thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn attention, 6.5 per 1k.

Five voicing rules compiled in Go:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating a content word across joint. Balanced two-beat.
- `apology` — reply performs contrition instead of stating correction and moving on.
- `self_postmortem` — reply turns to account for its own errors, story reader never asked for.
- `announced_length` — reply announces own length rather than cutting it.
- `cross_turn_repeat` — turn of phrase this reply shares with several earlier ones in same session. Only rule reading window rather than reply, so it stays quiet until session holds history.

Structure rules, all compiled:

- `labelled_opening` — prose paragraph opening on short verbless fragment that rest of it unpacks; ordinal counts as label. Skips list blocks, paragraphs under twelve words, and bolded form — card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so bold opener stays deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths hold coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — question or request for reader sitting in earlier block while reply carries on past it.
- `dangling_end` (interactive) — open problem named in closing blocks, no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — open problem landing after last question or offer, burying decision point above it.
- `forked_end` (interactive) — two or more things to act on in closing blocks, nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding another.
- `unverified_done` (loop) — says work done with nothing on page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — unattended run ends by asking, so answer lands in log and next iteration reads question as instruction to itself.

Grouping tells you something usable about implementation. POSTPROC pattern matches span of text, so it describes wording and nothing else. Every voicing rule needing more than a pattern went into Go beside structure rules. Hence three POSTPROC rules in shipped card and no more.

Reader expecting long list of banned phrases wants basanite. That list lives in another tool on purpose.

Structure rules vary one way, and card plays no part: lane. Interactive lane covers any turn that is not a loop turn — somebody waits at a terminal and ending decides what happens next. Loop lane covers turns opened by `/loop` or `/goal`, or by the sentinel a dynamic-pacing loop sends itself. Loop lane drops `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody reads yet. Report naming what it left open, then stopping, fails three of dropped rules, and a question in it lands in a log where next iteration reads it as instruction to itself. Claim check replaces them: report saying work done has to say what it ran.

## Flags

| Flag | Default | Does |
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

- `$XDG_STATE_HOME/cope/violations.jsonl`, mode 0600 — one JSON record per violation, carrying matched text and about 70 characters either side. This file quotes replies back.
- `$XDG_STATE_HOME/cope/refresher-<session-id>`, mode 0600 — empty file, mtime serves as refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json`, mode 0600 — rolling record mid-session injection gets chosen from: turn count, characters, rules fired over 20-turn window. No prose stored. Rule names and counts only.

## Editing the card

effigy owns `.effigy` grammar. `make rules` regenerates `card/rules.json`. `make check-rules` runs in CI, so enforced rules and injected rules cannot drift.

NEVER budget: 10. Anything over budget gets reported at load, never dropped silently. Budget charges against each injection separately, not against card file: SessionStart prints always-on rules, refresher prints evidence-gated ones, no code path renders their union. Card holds more NEVER rules in total than 10 and stays healthy.

Two card-authored forms, one per line in card header.

`@gate <rule_id> off — <why>`

`@shape <id>: <selector> <predicate> — <why>`

@shape vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

`@gate` needs a rule id gate already holds. `@shape` needs an id colliding with none. Wrong id reports at load, never gets ignored. Reason after dash: required in both. Rule a card wrote and rule a card refused read as unreviewable without one.

Declined rule still runs. Only this card's score drops it, so backfill still reports what it would have caught. @shape violation reports in card's own words, never in a sentence binary supplies.

`handoff.effigy` sets 60 words from measurement: across 43,155 assistant replies, closing block runs 33 words at median, 56 at p90.

## Calibrating

`cope-gate --backfill` scores whole session transcript at once. Rules got chosen that way. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects` — largest N, not one per project.

Watch hits per character. Share of turns hit tracks turn length instead. Rates sit in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

Limits split by axis.

`labelled_opening` runs no tagger. `ask_not_last` says nothing about order among several asks.

Hit rate runs roughly four fifths structure. A/B run found that four fifths tracks what a reply was for rather than how somebody wrote it, so hit rate describes output and judges nothing. Judgement lives in discrimination test, and that test covers voicing only.

Largest limit sits where Two things named it. Card declines a built-in rule and writes one of its own, in a vocabulary counting words and sentences and asking whether a block poses a question. Compiled rules remain the only home for a check like `clause_symmetry`. Card wanting something outside that vocabulary and outside a POSTPROC regex still has nowhere to put it.

Both directions amount to card marking own homework. Decline lowers that card's score. Assertion raises it. Read both with reason attached — hence syntax demanding one. [MEASUREMENTS.md](MEASUREMENTS.md) holds the runs.

## Layout

| Path | What |
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

`tools/generate_readme.py` wrote this README from the prompt `cope-gate --author-docs` emits, checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
