# cope

Cope ships one opinionated card, on from install, and checks two jobs at once — how a reply sounds, how a reply sits on page — and card stays a file, so edit or swap it. Name comes from foundry: cope means upper half of mould, half carrying shape cast into.

Read [demo/README.md](demo/README.md). Every file there holds this same README, same prompt, same facts, written again from another card, so card alone changed — set two beside each other and card's work shows faster than rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) comes from a card instructing every tic this model measures for, and lands the point in one glance, though reading it hurts.

## Two jobs

Voicing means sound of sentences: register, rhythm, diction, where flair gets licence. Voicing lives in card. Swap card, swap every word of it. Blind discrimination test measures this half.

Structure means shape of reply as thing a reader uses: where decision sits, how reply ends, whether an ask comes last, whether a done-claim carries proof. Structure compiles into binary, same rules under every card.

Sentence reaching for balanced two-beat — two clauses, near-equal length, one content word repeated across joint — fails voicing. Reply naming an open problem in last paragraph, then stopping, leaves "continue" nothing to point at, and fails structure. One reply passes one axis, fails other. Keep axes apart for that reason.

Card reaches into structure two ways. Card declines a built-in rule: `@gate <rule_id> off — <why>`, one per line in card header. That form exists because a card whose `VOICE` asked for a landing `clause_symmetry` catches got marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for that reason. Card also states a structure rule of its own: `@shape <id>: <selector> <predicate> — <why>`, one per line in card header. That form exists because a card's own promise about endings had nowhere to be checked — `card/demo/handoff.effigy` asserts `readable_cold`, `last paragraph words <= 60`, since its peak asks a reader to re-enter cold and read closing block alone. Sixty comes from measurement: across 43,155 assistant replies closing block runs 33 words at median, 56 at p90.

Boundary: `@shape` counts words, counts sentences, asks whether a block poses a question. Nothing more. `clause_symmetry` stays unwritable there. Card wanting a check outside that vocabulary and outside a `POSTPROC` regex has nowhere to put it.

## Problem

Three parts. Where instruction sits. What instruction says. Failure no instruction touches.

Where instruction sits comes first, since most readers hold it backwards. Reader edits global `CLAUDE.md`, watches it slide, and believes that file sits in system prompt. It does not. It arrives as one message pinned to first turn, and conversation buries it under everything written after. An output style sits in system prompt itself, and harness re-reminds model of it as conversation runs. Move one card between those two places, change no word, and most of cope's effect comes from that move alone. Measured — see [MEASUREMENTS.md](MEASUREMENTS.md).

Instruction alone leaves phrasing alone. Global `CLAUDE.md` banning not-A-but-B flip gets read every turn. Flip still appeared twice in session building this, with ban as topic. Name a surface form, move pushes into a variant. That covers voicing.

Third part holds no phrasing at all. Ending leaving a reader nothing to answer costs a whole round trip. No instruction bans that, since no habit of wording produces it.

Flip stays an anecdote about one rule. Claim rests elsewhere: blind discrimination test, where a reader sees only a voice's own description of itself, then picks which of two replies came from it. [MEASUREMENTS.md](MEASUREMENTS.md) holds rate and caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick style under `/config` -> Output style.

First command builds binary. `--setup` does whole install: writes output style, wires two hooks into `~/.claude/settings.json` with absolute paths, prints one step left. It backs settings file up first, adds only what stays missing so second run changes nothing, leaves every other key alone including other tools' hooks on same events, and refuses a settings file that will not parse. `--setup --dry-run` prints changes, writes nothing. Then menu choice above.

Rather not let a tool edit settings? `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it writes another. Standalone `/output-style` command went away in Claude Code v2.1.91, so `/config` remains route; same thing sets as `"outputStyle"` in `.claude/settings.local.json`.

Session reads a style once at start. New selection, or re-emitted card, applies at next session or after `/clear`. Running conversation shows no change, and that means nothing broke.

Card goes here rather than into a hook because an output style lands at end of system prompt and harness keeps reminding model of it.

Now hooks. Card no longer arrives through one. Hooks block of `~/.claude/settings.json`:

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

`Stop` buys a score on reply just written. `UserPromptSubmit` buys a mid-session restatement of rules actually firing, which a file written once cannot give. Voice runs without either. These two form measurement half.

Commands assume `go install` target directory sits on `PATH` inside environment hook runs in. Hook doing nothing, quietly, usually wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since `card/rules.json` sits committed and compiles in. `--inject` remains as superseded delivery for anyone wanting old route, and stands down on its own once a cope output style goes active.

## Another voice

Gate reads `.effigy` straight, so a card works as written — no Python, no effigy checkout. Drop a card in `$XDG_CONFIG_HOME/cope/cards`, reach it with `--card <name>` or `COPE_CARD`. `--rules` takes a path from anywhere. `make cards` installs demo set. Name resolving to nothing raises an error and injects nothing, rather than a session writing shipped voice while its config names another.

Read `card/demo/lecturer.effigy` first. It differs from shipped card on register alone, and discrimination run measured exactly that pair. Card changes voicing axis, per two jobs above. [MEASUREMENTS.md](MEASUREMENTS.md) holds numbers.

Want card's work without writing one? Read [demo/README.md](demo/README.md) — every file under `demo/` holds this README again from another card, same prompt, same facts, voice as sole variable.

`card/demo/handoff.effigy` stands apart in that directory: hypothesis, not voice. It keeps shipped card's handoff rules, drops everything about prose, and wants `make pairs` against full card rather than a render. `make cards` installs it with rest, so a reader listing cards finds it and should know it names no register to write in.

## What hooks add

`Stop` scores reply and appends fired rule names to session's rolling state, plus one record per violation to log.

`UserPromptSubmit` reads that rolling state — not violations log — and injects card items gated on what fires, counts named. No history yet, so it falls back to standing `CONTINUE TEST`. It stays quiet until last injection ages past `--refresh-every`. Mid-session text comes from measured output rather than fixed ahead, and a pasted `CLAUDE.md` cannot do that. One mechanism, no guarantee: A/B in repo does not separate refresher from no refresher.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) names character cards for game NPCs. Cope uses it off-label, since three blocks do what a prose gate needs: `POSTPROC` holds regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples — so a rule names a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers same problem other way round, and suits a reader finding cope too blunt. Cope bans: rule fires or does not, card says never, register locks once picked. Basanite measures: lemma frequency against a baseline over real transcripts, reporting what you lean on lately, judgement left to you. Its README calls that awareness rather than prohibition. Mood decides, not correctness — heavy hand suits a habit annoying you today, moving measurement suits watching drift rather than legislating it. They compose: different hooks, no shared state. Run both.

[caveman](https://github.com/JuliusBrussee/caveman), separate project, different author, compresses agent replies to cut output tokens — third axis again. Cope shapes prose, basanite tracks vocabulary, caveman shortens. Reader wanting fewer tokens rather than different structure goes there.

## Rules

Voicing rules.

- `flip` — not-A-but-B flip in common surface forms, not-only-but-also included. Action `warn`. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude atop it among every writer in study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k heaviest measured lean in this register. Action `warn`. Say what thing carries instead.
- `worth_noting` — announces that something deserves attention rather than letting it earn attention, at 6.5 per 1k. Action `warn`.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating a content word across joint. Balanced two-beat.
- `apology` — reply performs contrition rather than stating correction and moving on.
- `self_postmortem` — reply turns to account for its own errors, a story reader never asked for.
- `announced_length` — reply announces own length rather than cutting it.
- `cross_turn_repeat` — turn of phrase this reply shares with several earlier ones in same session. Only rule reading window rather than reply, so it cannot fire until session holds history.

Structure rules.

- `labelled_opening` — prose paragraph opening on a short verbless fragment rest of it unpacks; an ordinal counts as label. Skips list blocks and paragraphs under twelve words, and skips bolded form too — card dropped `bold_label` in July 2026 after blind readers named bold and bullets as things they wanted, so a bold opener stays unpoliced on purpose.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths carry a coefficient of variation below `--min-cv`.
- `ask_not_last` — question or request for reader sitting in an earlier block while reply carries on past it. Interactive lane.
- `dangling_end` — open problem named in closing blocks, no question, no offer, no explicit all-clear anywhere, so "continue" points at nothing. Interactive lane.
- `buried_decision` — open problem landing after last question or offer, burying decision point above it. Interactive lane.
- `forked_end` — two or more things to act on in closing blocks, nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding another. Interactive lane.
- `unverified_done` — says work finished with nothing on page that could show it: no command, no count, no file. Loop lane.
- `loop_ask` — unattended run ends by asking, so answer lands in a log and next iteration reads question as instruction to itself. Loop lane.

Grouping tells a reader something usable about implementation. A `POSTPROC` pattern matches a span of text, so it describes wording alone. Every voicing rule needing more than a pattern got written in Go beside structure rules. Shipped card therefore carries 3 `POSTPROC` rules. Reader expecting a long list of banned phrases wants basanite, which holds that list on purpose.

Structure rules vary in one place, and card never causes it. Lane does. Interactive lane covers any turn that stays outside a loop: somebody waits at a terminal and ending decides what happens next. Loop lane covers a turn opened by `/loop`, `/goal`, or sentinel a dynamic-pacing loop sends itself. Loop drops `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision`, and adds `unverified_done`, `loop_ask`. Nobody reads yet. Report naming what it left open, then stopping, fails three of dropped rules, and a question inside it lands in a log where next iteration reads it as instruction to itself. Claim check replaces them: report saying work finished has to say what ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default inject,hold; positive is the third) |
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

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | `0600` | one JSON record per violation, carrying matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | `0600` | empty file whose mtime holds refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | `0600` | rolling record mid-session injection gets chosen from: turn count, characters, which rules fired over a 20-turn window. No prose, only rule names and counts |

Log quotes replies back.

## Editing card

Effigy owns `.effigy` grammar. `make rules` regenerates `card/rules.json`. `make check-rules` runs in CI, so enforced and injected rules cannot drift.

`NEVER` budget sits at 10. Anything over budget gets reported at load rather than dropped quietly. Budget charges against each injection separately, not against card file: `SessionStart` prints always-on rules, refresher prints evidence-gated ones, no code path renders their union. So a card holds more `NEVER` rules in total than budget and stays healthy.

Card-authored rules, two forms.

`@gate <rule_id> off — <why>`, one per line in card header.

`@shape <id>: <selector> <predicate> — <why>`, one per line in card header. Vocabulary:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

`@gate` needs a rule id gate already holds. `@shape` needs an id colliding with none. Wrong id gets reported at load rather than ignored. Reason after dash stays required in both forms, since a rule a card wrote and a rule a card refused both read as unreviewable without one. Declined rule still runs; this card's score alone drops it, so a backfill still reports what it would have caught. `@shape` violation gets reported in card's own words rather than any sentence binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and rules got chosen that way. `tools/backfill-sweep.sh` runs it over N largest transcripts found anywhere under `~/.claude/projects` — largest, not one per project.

Watch hits per character, not share of turns hit. Second number tracks turn length. [MEASUREMENTS.md](MEASUREMENTS.md) holds rates.

## Limits

Axis split organises limits.

`labelled_opening` tags nothing. `ask_not_last` says nothing about order among several asks.

Hit rate runs roughly four fifths structure. A/B run found that four fifths tracks what a reply was for rather than how it got written, so hit rate describes output and judges nothing. Judgement lives in discrimination test, and that test covers voicing alone.

Largest limit sits where two jobs above named it, and current shape reads so: a card declines a built-in rule and writes one of its own, and vocabulary it writes in counts words, counts sentences, asks whether a block poses a question. Compiled rules stay only home for a check like `clause_symmetry`. Card wanting something outside that vocabulary and outside a `POSTPROC` regex still has nowhere to put it.

Both directions leave card marking own homework. Decline lowers that card's score. Assertion raises it. Read both with reason attached — hence syntax demanding one. [MEASUREMENTS.md](MEASUREMENTS.md).

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | hook binary |
| `internal/scan/` | structure rules, card's regex rules, card renderer |
| `internal/effigy/` | `.effigy` reader, so a card works as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane a turn came from |
| `replay/` | blind-pairs and discrimination harnesses, plus own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what ran, on how much text, what it said |

---

`tools/generate_readme.py` wrote this README from prompt `cope-gate --author-docs` emits, checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
