# cope

cope ships one loud card — on from install, whole product for most readers — and scores two jobs, sound of sentences and shape of reply; card sits in one file, so edit it or swap it. Name comes from foundry work: cope means upper half of mould, half carrying shape cast into.

Read [demo/README.md](demo/README.md). Every file there holds this same page written again from other card — same prompt, same facts, card only change — so reading two of them against each other shows what card does faster than rest of this page explains. One file makes point in single glance: [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from card instructing every tic this model carries. Hard going on purpose.

## Two jobs

Voicing means sound of sentences: register, rhythm, diction, what paragraph does with detail, where flair gets licence. Voicing lives in card — VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC — so swapping card swaps every word of it, and voicing holds measured result behind it. Structure means shape of reply as thing reader uses: where decision sits, whether ending gives "continue" something to point at, whether ask lands last, whether claim of finished work carries anything showing it. Structure compiles into binary, same under every card, varying by lane alone. Sentence reaching for balanced two-beat fails on voicing. Reply naming open problem in last paragraph, then stopping, fails on structure. One reply passes one axis and fails other — reason to keep axes apart.

Card reaches into structure two ways.

Card declines built-in rule with `@gate <rule_id> off — <why>`, one per line in card header. `card/demo/lecturer.effigy` declines clause_symmetry and dangling_end: its VOICE block asks for balanced landing and arriving close, so card obeys own VOICE and gate marks card down for it. Card states own structural rule with `@shape <id>: <selector> <predicate> — <why>`, again one per line. `card/demo/handoff.effigy` asserts readable_cold — last paragraph words <= 60 — because its peak asks reader to re-enter cold and read last block alone, and no built-in rule checks whether that block reads that way. That 60 comes from measurement: across 43,155 assistant replies, closing block runs 33 words at median and 56 at p90.

Boundary sits here: `@shape` vocabulary counts words, counts sentences, asks whether block poses question, nothing further. clause_symmetry falls outside it. Card wanting check outside both that vocabulary and POSTPROC regex holds nowhere to put it. Run and reasons live in [MEASUREMENTS.md](MEASUREMENTS.md).

## Problem

Most readers arrive holding first part backwards.

Reader edits global CLAUDE.md, watches voice drift back, concludes that file sits in system prompt. Wrong place. CLAUDE.md arrives as one message attached to first turn, and conversation buries it under everything written after. Output style sits inside system prompt, and harness re-reminds model of style while conversation runs. Move one card between those two places, change no word of card, and most of cope's effect comes from that move alone. Measured — see [MEASUREMENTS.md](MEASUREMENTS.md).

Instruction alone leaves phrasing where it stood. Global CLAUDE.md bans "not A, it's B" flip. Model reads that ban every turn. Flip still appears twice in session that built this repo, ban itself as topic. Naming one surface form pushes move into variant. That covers voicing side.

Structure side carries different complaint from different cause. Ending leaving reader nothing to answer costs whole round trip, and no phrasing ban reaches it.

Flip counts as anecdote about one rule. Claim rests elsewhere: blind discrimination test, where reader sees only voice's own description of itself and picks which of two replies came from that voice. Rate and caveats sit in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick style under `/config` -> Output style.

`--setup` emits output style, wires two hooks into `~/.claude/settings.json` with absolute paths, prints one step left. It backs settings file up first. It adds only missing pieces, so second run changes nothing. It leaves every other key alone, other tools' hooks on same events included, and refuses to touch settings file failing parse. `--setup --dry-run` prints planned changes and writes nothing.

By-hand route stays open for anyone keeping settings untouched by tools: `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits other card. Select under `/config` -> Output style; standalone `/output-style` command left Claude Code at v2.1.91. Same choice sets as `"outputStyle"` in `.claude/settings.local.json`.

Style reads once at session start. New selection, or freshly emitted card, applies at next session or after `/clear`.

Output style goes at end of system prompt, and harness re-reminds model of it mid-conversation. Hook delivery lacks that.

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

Stop scores reply just written. UserPromptSubmit restates mid-session those card items actually firing, which file written once cannot do. Voice works without both hooks; hooks form measurement half.

Both commands assume go install's target directory on PATH inside environment running hook. Hook doing nothing silently usually wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since `card/rules.json` sits committed and compiled in.

`--inject` remains as superseded delivery for anyone wanting old path, and it stands down on its own once cope output style runs.

## Writing in another voice

Gate reads `.effigy` direct, so card works as written — no Python, no effigy checkout. Drop card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes path from anywhere; `make cards` installs demo set. Name resolving to nothing raises error and injects nothing, rather than letting session write shipped voice while its config names other card. Read `card/demo/lecturer.effigy` first: that card differs from shipped card on register alone, and discrimination run measured that pair. Card changes voicing axis, as "Two jobs" sets out. Numbers stay in [MEASUREMENTS.md](MEASUREMENTS.md).

Shortest look at what card does, without writing one: [demo/README.md](demo/README.md), where every file holds this page written again from other card, same prompt and same facts, voice as only variable.

`card/demo/handoff.effigy` stands as exception in that directory: hypothesis, not voice — keeps shipped card's handoff rules, drops everything about prose, and runs through `make pairs` against full card instead of rendering. `make cards` installs it beside rest, so anyone listing cards finds it and learns it names no register.

## What hooks add

Stop scores reply and appends fired rule names to session's rolling state, plus one record per violation to log. UserPromptSubmit reads that rolling state — not violations log — and injects card items gated on what fires lately, naming counts. Mid-session text comes from measured output instead of fixed text, and pasted CLAUDE.md reaches nothing like that. One mechanism, no guarantee: A/B run in this repo separates arms, not refresher from no refresher. SessionStart `--inject` stays superseded and off by default.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) holds character-card notation for game NPCs, used here off-label. Three of its blocks do work prose gate needs: POSTPROC holds regex rules with warn action applied after generation, WRONG holds anti-pattern beside its replacement, TEST holds named question with fail example and pass example — that way rule names move rather than one wording of move.

[basanite](https://github.com/justinstimatze/basanite) answers same problem other way round, and suits anyone finding cope too blunt. cope bans: rule fires or rule stays quiet, card says never, register fixes at moment of choice. basanite measures — lemma frequency against baseline over real transcripts — reports what you lean on lately, leaves judgement with you, and its own README calls that awareness rather than prohibition. Choice tracks mood, not correctness: heavy hand suits habit annoying you today, moving measurement suits watching drift. Both compose — different hooks, no shared state.

## Rules

Voicing rules:

- **flip** (warn) — not-A-but-B flip in common surface forms, not-only-but-also included. Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude top of it among every writer in study, human and machine.
- **load_bearing** (warn) — reflexive intensifier standing in for important or central. 25.6 per 1k, heaviest measured lean in this register. Say what thing carries.
- **worth_noting** (warn) — announces that something deserves attention instead of letting it earn attention. 6.5 per 1k.
- **clause_symmetry** — comma- or semicolon-joined clauses of near-equal length repeating content word across joint. Balanced two-beat.
- **apology** — reply performs contrition instead of stating correction and moving on.
- **self_postmortem** — reply turns to account for own errors, story reader never asked for.
- **announced_length** — reply announces own length rather than cutting it.
- **cross_turn_repeat** — turn of phrase this reply shares with several earlier ones in same session. Only rule reading window rather than reply, so it stays quiet until session holds history.

Structure rules:

- **labelled_opening** — prose paragraph opening on short verbless fragment that rest of paragraph unpacks; ordinal counts as label. Skips list blocks and paragraphs under twelve words. Skips bolded form too: card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as things they want, so bold opener goes unpoliced on purpose.
- **paragraph_uniformity** — four or more prose paragraphs whose lengths carry coefficient of variation below `--min-cv`.
- **ask_not_last** — question or request for reader sitting in earlier block while reply carries on past it.
- **dangling_end** — open problem named in closing blocks, no question, no offer, no explicit all-clear anywhere, so "continue" points at nothing.
- **buried_decision** — open problem landing after last question or offer, decision point buried above it.
- **forked_end** — two or more things to act on in closing blocks, nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding another.
- **unverified_done** — says work finished with nothing on page showing it: no command, no count, no file.
- **loop_ask** — unattended run ends by asking, so answer lands in log and next iteration reads question as instruction to itself.

Grouping tells reader something usable about implementation. POSTPROC pattern matches span of text, so pattern describes wording and nothing else. Every voicing rule needing more than pattern went into Go beside structure rules. That leaves shipped card three POSTPROC rules.

Anyone expecting long list of banned phrases finds that list in basanite, on purpose.

Structure rules vary in one place alone — not by card, by who reads turn:

- **interactive** — any turn falling outside loop lane. Somebody waits at terminal, and ending holds their next decision.
- **loop** — prompt opening turn reads `/loop` or `/goal`, or sentinel that dynamic-pacing loop sends itself. Drops ask_not_last, forked_end, dangling_end, buried_decision. Adds unverified_done, loop_ask. Nobody reads yet. Report correctly naming what it left open, then stopping, fails three of dropped rules, and question inside it lands in log where next iteration reads it as instruction to itself. Claim check replaces them: report saying work done must say what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotates refresher windows through arms, records which arm each turn ran under |
| `--ab-arms` | `(empty)` | comma-separated arms to rotate through, implying -ab (default inject,hold; positive comes third) |
| `--ab-report` | `(empty)` | reads turn log, reports arms; `-` reads default path |
| `--author-docs` | `false` | prints prompt for writing this repo's docs: card, introspected facts, sections |
| `--backfill` | `(empty)` | scores every assistant turn in this transcript, exits |
| `--block` | `false` | exits 2 on violation whose action reads reject (default warn-only) |
| `--card` | `(empty)` | names installed card to write in, from cards directory; reads `$COPE_CARD` too |
| `--cards` | `false` | lists installed cards with aim each one states, exits |
| `--check` | `(empty)` | scores prose file against card, exits; `-` reads stdin |
| `--describe` | `false` | prints card's voice as target to recognise: aim and register, machinery left out |
| `--display` | `false` | MessageDisplay entry: rewrites what reader sees, leaves transcript alone |
| `--display-preview` | `false` | reads prose on stdin, prints it as `--display` rewrites it |
| `--dry-run` | `false` | with `--setup`, prints planned changes, writes nothing |
| `--inject` | `false` | prints card as prompt text for SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | appends violations here; empty disables |
| `--min-cv` | `0.35` | flags paragraph-length coefficient of variation below this |
| `--output-style` | `false` | writes card to `~/.claude/output-styles` as Claude Code output style, landing card in system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory receiving output style (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of last card or refresher injection before refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: injects compact reminder once last injection ages |
| `--render-arm` | `(empty)` | prints mid-session reminder one arm injects, exits |
| `--render-for` | `(empty)` | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | `(empty)` | renders `-render-arm` as given lane sees it: interactive (default) or loop |
| `--rules` | `(empty)` | reads card from this `.effigy` or `.json` file instead of card built into binary |
| `--setup` | `false` | emits output style, wires hooks, prints one step left |
| `--version` | `false` | prints version, exits |

## What lands on disk

Three files, each mode 0600:

- `$XDG_STATE_HOME/cope/violations.jsonl` — one JSON record per violation, carrying matched text plus about 70 characters either side. Log quotes replies back.
- `$XDG_STATE_HOME/cope/refresher-<session-id>` — empty file whose mtime serves as refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json` — rolling record driving mid-session injection: turn count, characters, which rules fired over 20-turn window. No prose, rule names and counts alone.

## Editing card

effigy owns `.effigy` grammar. `make rules` regenerates. `make check-rules` runs in CI, so enforced rules and injected rules cannot drift apart.

NEVER budget holds at 10. Rules past budget get reported at load, never dropped in silence. Budget charges each injection on its own, not card file: SessionStart prints always-on rules, refresher prints evidence-gated rules, and no code path renders their union — so card holding more NEVER rules than budget still counts healthy.

Card-authored rules take two forms, one per line in card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

`@shape` vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

Rule id under `@gate` must name rule gate already holds. Rule id under `@shape` must collide with none. Wrong id gets reported at load, never ignored. Reason after dash stays required in both forms: rule card wrote and rule card refused both read as unreviewable without one. Declined rule still runs, and this card's score alone drops it, so backfill still reports what that rule catches. `@shape` violation reports in card's own words rather than in sentence binary supplies.

## Calibrating

`cope-gate --backfill` scores whole session transcript at once, and rules came out of that. `tools/backfill-sweep.sh` runs it over N largest transcripts found anywhere under `~/.claude/projects`, which differs from one per project.

Watch hits per character. Share of turns hit tracks turn length instead. Rates sit in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

Axis split organises them.

labelled_opening runs no tagger. ask_not_last says nothing about order among several asks. Hit rate runs roughly four fifths structure, and A/B run found that four fifths tracks what reply was for rather than how reply got written — description of output, no judgement of it. Judgement lives in discrimination test, and that test covers voicing alone.

Largest limit belongs to "Two jobs", as it now stands: card declines built-in rule and writes rule of its own, and vocabulary it writes in counts words, counts sentences, asks whether block poses question. Compiled rules stay only home for check like clause_symmetry. Card wanting something outside both that vocabulary and POSTPROC regex holds nowhere to put it.

Both directions leave card marking own homework. Decline lowers that card's score; assertion raises it. Read each with reason attached — syntax demands reason for that. See [MEASUREMENTS.md](MEASUREMENTS.md).

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | hook binary |
| `internal/scan/` | structure rules, card's regex rules, card renderer |
| `internal/effigy/` | `.effigy` reader, so card works as written |
| `internal/transcript/` | Claude Code JSONL reader, plus lane each turn ran in |
| `replay/` | blind-pairs and discrimination harnesses, with own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what ran, on how much text, and what it said |

---

`tools/generate_readme.py` writes this README from prompt `cope-gate --author-docs` emits, checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
