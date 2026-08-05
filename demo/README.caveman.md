# cope

cope ships one card — blunt, opinionated, live from moment you install it: card fixes how reply sounds, gate checks how reply ends, and card sits in one file you edit or swap. Foundry names upper mould half cope, half carrying shape metal casts into.

Read [demo/README.md](demo/README.md) first: every file there holds this page written again from another card, same prompt and same facts, card as only difference — read two against each other and card work shows faster than rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes point in one glance, written from card instructing every tic measured in this model, and hard going on purpose.

## Two jobs

Voicing — sound of sentences. Register, rhythm, diction, what paragraph does with detail, where flair earns license. Voicing lives in card, all of it: `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES`, `POSTPROC`. Swap card, swap every word. Blind discrimination test measures this half.

Structure — where decision sits, how reply ends, whether `continue` has anything to point at. Structure compiles into binary, `internal/scan`. Same rules under every card.

One fault each. Comma-joined clauses of equal length repeating one content word across joint: voicing. Reply naming open problem in last paragraph, then stopping, leaving reader nothing to answer: structure. Same reply passes one axis and fails other. Reason to keep them apart.

Card reaches into structure half two ways. Card declines built-in rule with `@gate <rule_id> off — <why>`: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its `VOICE` block asks for balanced landing and arriving close those two rules catch, so gate marked card down for obeying itself. Card asserts own rule with `@shape <id>: <selector> <predicate> — <why>`: `card/demo/handoff.effigy` asserts `readable_cold`, because its promise about closing block had nowhere to sit and no built-in rule checks it.

`@shape` vocabulary counts words, counts sentences, asks whether block poses question. Nothing further. `clause_symmetry` stays in Go.

## Why file in your home directory fades

Global `CLAUDE.md` never reaches system prompt. Reader who edited that file and watched it stop working usually believes otherwise. `CLAUDE.md` arrives as one message attached to first turn, and conversation buries it under everything written after. Output style sits inside system prompt, and harness re-reminds model of it while conversation runs. Move one card between those two places, change no word of it, and most of cope works. Run sits in [MEASUREMENTS.md](MEASUREMENTS.md).

Instruction alone fixes no phrasing. Global `CLAUDE.md` banning not-A-but-B flip gets read every turn; session building this wrote flip twice while ban stood as topic. Name surface form and move slides to variant. Voicing side.

Third fault carries no phrasing at all. Ending that leaves reader nothing to answer costs whole round trip. No instruction bans it, because fault lives in shape.

Flip counts as anecdote about one rule. Claim rests on blind discrimination test: reader sees one voice's own description of itself, then picks which of two replies came from that voice. [MEASUREMENTS.md](MEASUREMENTS.md) holds rate and caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`cope-gate --setup` emits output style, wires two hooks into `~/.claude/settings.json` with absolute paths, prints one step left. It backs settings file up first. It adds only missing pieces, so second run changes nothing. It leaves every other key alone, other tools' hooks on same events included. It refuses to touch settings file that fails to parse. `cope-gate --setup --dry-run` prints changes and writes nothing.

Then menu step: `/config` -> Output style -> `claude_voice`. Entry reads `claude_voice`, shipped card id; word cope appears nowhere in that list, so scan for `claude_voice`. Claude Code v2.1.91 removed standalone `/output-style` command, so `/config` carries it now. Same setting goes in `.claude/settings.local.json` as `"outputStyle": "claude_voice"`.

By hand, for anyone keeping settings unwritten: `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front emits another card.

Session reads style once, at start. New selection, or freshly emitted card, applies next session or after `/clear`.

Output style goes at end of system prompt and harness re-reminds model of it during conversation. Card lands there for that reason.

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

`Stop` scores reply once model writes it. `UserPromptSubmit` restates rules that keep firing, mid-session, which file written once cannot do. Voice works without both hooks. Hooks buy measurement.

Both commands assume `go install` target directory on `PATH` inside environment hook runs in; hook doing nothing silently wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since `card/rules.json` ships committed and compiled in. `--inject` survives as superseded delivery for anyone who wants old path, and stands down on its own once cope output style goes live.

## Writing in another voice

Gate reads `.effigy` straight, so card works as written: no Python, no effigy checkout. Drop card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`. `--rules` takes path from anywhere. `make cards` installs demo set. Name resolving to nothing raises error and injects nothing, which beats session writing shipped voice while config names another one.

Read `card/demo/lecturer.effigy` first. It differs from shipped card on register alone, and discrimination run measured exactly that pair. Card changes sound of sentences and leaves shape of reply where binary holds it.

[demo/README.md](demo/README.md) shows card work without writing one: every file under `demo/` holds this README again from another card, same prompt, same facts, voice as only variable.

`card/demo/handoff.effigy` breaks that pattern — hypothesis, no register. It keeps shipped card's handoff rules, drops everything about prose, and belongs in `make pairs` against full card rather than in rendering. `make cards` installs it with rest, so reader listing cards finds it and should know it names no voice.

## What hooks add

`Stop` scores reply just written, appends which rules fired to session's rolling state, appends one record per violation to log.

`UserPromptSubmit` reads that rolling state. Violations log stays untouched. Hook injects card items gated on what keeps firing and names counts. Session with no history yet gets standing `CONTINUE TEST` instead. Hook stays quiet until last injection ages past `--refresh-every`.

So measured output picks mid-session text; pasted `CLAUDE.md` cannot. One mechanism, no guarantee: A/B run in repo compares variants and never separates refresher from silence.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) writes character cards for game NPCs. cope uses it off-label. Three blocks do what prose gate needs: `POSTPROC` holds regex rules with warn action applied after generation, `WRONG` holds anti-pattern beside its replacement, `TEST` holds named question with fail example and pass example — that shape lets rule name move rather than one wording of move.

[basanite](https://github.com/justinstimatze/basanite) answers same problem other way round, and suits reader who finds cope too blunt. cope bans: rule fires or rule stays quiet, card says never, register locks when you pick it. basanite measures lemma frequency against baseline over real transcripts, reports what you lean on lately, leaves judgement to you — awareness, in its own README's word. Heavy hand fits habit annoying you today. Moving measurement fits drift you would rather watch than legislate. Both compose: different hooks, no shared state. Run both.

[caveman](https://github.com/JuliusBrussee/caveman), separate project by another author, compresses agent replies to cut output tokens. Third axis again: cope shapes prose, basanite tracks vocabulary, caveman shortens. Reader wanting fewer tokens goes there.

## Rules

Voicing. Shipped card carries three `POSTPROC` patterns:

- `flip` (warn) — not-A-but-B flip in common surface forms, not-only-but-also included. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude top of it among every writer in study, human and machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, 25.6 per 1k, heaviest measured lean in this register. Say what thing carries.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn attention, 6.5 per 1k.

Rest of voicing compiles into Go:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating content word across joint. Balanced two-beat.
- `apology` — reply performs contrition instead of stating correction and moving on.
- `self_postmortem` — reply turns to account for own errors, story reader never requested.
- `announced_length` — reply announces own length rather than cutting it.
- `cross_turn_repeat` — turn of phrase shared with several earlier replies in same session. Only rule reading window rather than reply, so it stays silent until session builds history.

Structure, all compiled:

- `labelled_opening` — prose paragraph opening on short verbless fragment that rest of paragraph unpacks; ordinal counts as label. Skips list blocks and paragraphs under twelve words, and skips bolded form: card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as things they wanted, so bold opener goes unpoliced on purpose.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths show coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — question or request for reader sitting in earlier block while reply carries on past it.
- `dangling_end` (interactive) — open problem named in closing blocks with no question, no offer, no explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` (interactive) — open problem landing after last question or offer, decision point buried above it.
- `forked_end` (interactive) — two or more things to act on in closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding second one.
- `unverified_done` (loop) — says work done with nothing on page that could show it: no command, no count, no file.
- `loop_ask` (loop) — unattended run ends by asking, so answer lands in log and next iteration reads question as instruction to itself.

Grouping tells you something usable about implementation. `POSTPROC` pattern matches span of text, so pattern describes wording and nothing else. Every voicing rule needing more than pattern went into Go beside structure rules. Hence three `POSTPROC` rules in shipped card, and no more. Card writing this page carries none of its own.

Reader expecting long list of banned phrases should reach for basanite. That list lives in another tool on purpose.

Structure rules vary once, by lane, never by card. Interactive lane covers every turn that skips loop: somebody waits at terminal and ending carries their next decision. Loop lane covers turns opened by `/loop` or `/goal`, or by sentinel that dynamic-pacing loop sends itself. Nobody reads yet. Report naming what it left open and then stopping fails `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision`, so loop lane drops those four. Question inside such report lands in log where next iteration reads it as instruction to itself. Loop lane adds `unverified_done` and `loop_ask` instead: report claiming work done has to say what it ran.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which variant wrote each turn |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `--ab` (default `inject,hold`; positive counts as third) |
| `--ab-report` | (empty) | read turn log, report how each variant did; `-` reads default path |
| `--author-docs` | `false` | print prompt for writing this repo's docs: card, introspected facts, sections |
| `--backfill` | (empty) | score every assistant turn in this transcript, exit |
| `--block` | `false` | exit 2 on violation whose action says reject (default warns only) |
| `--card` | (empty) | name of installed card to write in, from cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list installed cards with aim each one states, exit |
| `--check` | (empty) | score prose file against card, exit; `-` reads stdin |
| `--describe` | `false` | print card's voice as target to recognise: aim and register, machinery left out |
| `--display` | `false` | `MessageDisplay` entry: rewrite what reader sees, leave transcript alone |
| `--display-preview` | `false` | read prose on stdin, print it as `--display` rewrites it |
| `--dry-run` | `false` | with `--setup`, print changes, write nothing |
| `--inject` | `false` | print card as prompt text for `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write card to `~/.claude/output-styles` as Claude Code output style, which puts card in system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of last card or refresher injection before refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject compact reminder once last injection ages |
| `--render-arm` | (empty) | print mid-session reminder one variant injects, exit |
| `--render-for` | (empty) | comma-separated rule ids to render `--render-arm` against |
| `--render-lane` | (empty) | render `--render-arm` as given lane sees it: `interactive` (default) or `loop` |
| `--rules` | (empty) | read card from this `.effigy` or `.json` file instead of card built into binary |
| `--setup` | `false` | emit output style, wire hooks, print one step left |
| `--version` | `false` | print version, exit |

## What lands on disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600` — one JSON record per violation, holding matched text plus about 70 characters either side. Log quotes your replies back at you.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600` — empty file whose mtime serves as refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600` — rolling record that mid-session injection draws on: turn count, characters, which rules fired across 20-turn window. Rule names and counts only, prose never.

## Editing card

effigy owns `.effigy` grammar. `make rules` regenerates `card/rules.json`. CI runs `make check-rules`, so enforced rules and injected rules cannot drift apart.

`NEVER` budget: ten rules. Load reports anything over budget rather than dropping it in silence. Budget charges each injection separately, never card file — one injection prints always-on rules, refresher prints evidence-gated ones, and no code path renders their union. Card holding more `NEVER` rules than ten in total stays healthy.

Two card-authored forms, one per line in card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

`@shape` vocabulary:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

`@gate` takes id that gate already carries. `@shape` takes id colliding with none. Wrong id fails at load, loudly. Reason after dash: required in both forms, since rule card wrote and rule card refused read as unreviewable without one. Declined rule still runs, and only this card's score drops it, so backfill still reports what that rule would catch. `@shape` violation reports in card's own words rather than in any sentence binary supplies.

## Calibrating

`cope-gate --backfill` scores whole session transcript in one pass, and rules came out of exactly that. `tools/backfill-sweep.sh` runs it over N largest transcripts found anywhere under `~/.claude/projects`, ignoring project boundaries. Watch hits per character. Share of turns hit tracks turn length instead. Rates sit in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

`labelled_opening` tags nothing. It matches shape of opener, never part of speech.

`ask_not_last` says nothing about order among several asks.

Hit rate runs roughly four fifths structure, and A/B run found that four fifths tracks what reply was for rather than how somebody wrote it. Read it as description of output. Judgement lives in discrimination test, and that test covers sound of sentences only.

Split above between sound and shape carries largest limit. Card declines built-in rule, card writes own rule, and vocabulary it writes in counts words, counts sentences, asks whether block poses question. Compiled rules stay only home for check like `clause_symmetry`. Card wanting check outside that vocabulary and outside `POSTPROC` regex still has nowhere to put it.

Both directions leave card marking own homework: decline lowers that card's score, assertion raises it. Read both with reason attached — syntax demands reason for that reason. [MEASUREMENTS.md](MEASUREMENTS.md) carries run.

## Layout

| Path | What |
| --- | --- |
| `card/claude_voice.effigy` | shipped card, in effigy notation |
| `card/rules.json` | generated from it, embedded in binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | hook binary |
| `internal/scan/` | structure rules, card's regex rules, card renderer |
| `internal/effigy/` | `.effigy` reader, so card works as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane wrote each turn |
| `replay/` | blind-pairs and discrimination harnesses, plus their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what ran, on how much text, and what it said |

---

`tools/generate_readme.py` wrote this README from prompt `cope-gate --author-docs` emits, and `cope-gate --check` scored it.

MIT. justin@justinstimatze.com
