# cope

cope ships one opinionated card — on once installed, scoring two jobs at once, how a reply sounds and how a reply ends — and card sits in a file, so swap it or edit it; name comes from foundry, where cope means upper half of mould, half carrying shape cast into metal.

Read [demo/README.md](demo/README.md) first. Every file there holds this README again, written from another card, same prompt and same facts, so card alone varies — read two side by side, learn more than rest of page teaches. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes point in one glance, written from card instructing every tic this model measures with, and hard going on purpose.

## Two Things

Voicing means sound of sentences: register, rhythm, diction, where flair gets license. Voicing lives in card. Swap card, swap every word of it. Voicing carries measured result — blind discrimination test.

Structure means shape of reply as thing reader uses: where decision sits, whether ending gives "continue" something to point at, whether claim of done work carries proof. Structure compiles into binary. Same rules under every card.

One instance each, so split stays concrete. Sentence reaching for balanced two-beat — two clauses, near-equal length, content word echoed across joint — fails on voicing. Reply naming open problem in last paragraph, then stopping, leaving "continue" nothing to refer to — fails on structure. Same reply passes one axis, fails other. Keep axes apart for that reason.

Card reaches into structure half, two directions. Card declines built-in rule: `@gate <rule_id> off — <why>`, one per line in card header. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its `VOICE` block asks for balanced landing and arriving close those two rules catch — card obeyed itself, gate marked it down. Card writes own structural rule: `@shape <id>: <selector> <predicate> — <why>`, one per line in card header. `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks reader to enter cold and read last block only, and no built-in rule checks that block reads that way. Boundary: `@shape` counts words and sentences and asks whether block poses question, nothing more, so card wanting check outside that vocabulary and outside a `POSTPROC` regex has nowhere to put it. [MEASUREMENTS.md](MEASUREMENTS.md) holds run and reasons its numbers carry no further.

## Problem

Three parts. Where instruction sits. What instruction says. Failure no instruction touches.

Where instruction sits comes first, because most readers hold it backwards. Reader edits global `CLAUDE.md`, watches it fade, believes that file forms system prompt. It does not. It arrives as one message pinned to first turn, and conversation buries it under everything written after. Output style sits inside system prompt, and harness re-reminds model of it as conversation runs. Move one card between those two places, change no word of it, and that move accounts for most of why cope works. Measured — see [MEASUREMENTS.md](MEASUREMENTS.md). No counts here; reader wanting run follows link.

Instruction alone leaves phrasing untouched. Global `CLAUDE.md` banning "not A, it's B" flip gets read every turn. Flip still landed twice in session that built this, while ban formed topic of session. Name surface form, move slides into variant. That covers voicing side.

Third part names no phrasing problem at all. Ending leaving reader nothing to answer costs whole round trip. No instruction bans it, because no habit of wording produces it.

Flip stays anecdote about one rule. Claim rests elsewhere: blind discrimination test, where reader sees only a voice's own description of itself and picks which of two replies came from it. [MEASUREMENTS.md](MEASUREMENTS.md) holds rate and caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

First command installs binary. Second command does whole install: emits output style, wires three hooks into `~/.claude/settings.json` with absolute paths, prints one step left. `--setup` backs settings file up first, adds only what stays missing so second run changes nothing, leaves every other key alone including other tools' hooks on same events, and refuses to touch settings file that fails to parse. `--setup --dry-run` prints what it would change, writes nothing.

Rather not let a tool edit your settings? `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits another one.

Now menu choice. Pick style under `/config` -> Output style. Entry reads `claude_voice` — shipped card's id, not word cope. Say name out loud, because reader scanning that menu for cope finds nothing. Standalone `/output-style` command went away in Claude Code v2.1.91, so `/config` carries it; same thing sets as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

Style gets read once at session start. New selection, or freshly emitted card, applies at next session or after `/clear`. Running conversation shows no change, and that means nothing broke.

Card goes into output style rather than into a hook because output style lands at end of system prompt and harness re-reminds model of it mid-conversation.

Hooks come next. Card no longer arrives through one.

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
    ],
    "PreToolUse": [
      {
        "matcher": "mcp__linear__save_comment|mcp__linear__save_document|mcp__linear__save_issue|mcp__linear__save_project|mcp__linear__save_status_update",
        "hooks": [
          { "type": "command", "command": "cope-gate --pretool", "timeout": 10 }
        ]
      }
    ]
  }
}
```

That block forms hooks block of `~/.claude/settings.json`. `Stop` buys scoring of reply after it lands. `UserPromptSubmit` buys mid-session restatement of rules that keep firing, which file written once cannot do. `PreToolUse` buys score on prose an external write posts. Voice works without any of them; hooks form measurement half.

Commands above assume `go install` target directory sits on `PATH` in environment hook runs in. Hook doing nothing silently usually wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since `card/rules.json` ships committed and compiled in. `--inject` remains as superseded delivery for anyone wanting old path, and stands down on its own once a cope output style goes active.

## Writing In Another Voice

Gate reads `.effigy` directly. Card works as written — no Python, no effigy checkout. Drop card in `$XDG_CONFIG_HOME/cope/cards`, reach it with `--card <name>` or `COPE_CARD`. `--rules` takes path from anywhere. `make cards` installs demo set. Name resolving to nothing raises error and injects nothing, rather than session writing shipped voice while config names another.

Read `card/demo/lecturer.effigy` first. It parts from shipped card on register alone, and discrimination run measured exactly that pair. Card changes voicing, per split named above under two things; structure stays fixed. Numbers live in [MEASUREMENTS.md](MEASUREMENTS.md).

Shortest route to seeing what a card does, without writing one: [demo/README.md](demo/README.md). Every file under `demo/` holds this README written again from another card, same prompt, same facts, voice alone varying.

One exception in that directory: `card/demo/handoff.effigy` states a hypothesis, not a voice. It keeps shipped card's handoff rules and drops everything about prose, and it runs through `make pairs` against full card rather than rendering a page. `make cards` installs it with rest, so reader listing cards finds it there and should know it names no register to write in.

## What Hooks Do Differently

`Stop` runs `cope-gate`, scores reply just written, appends which rules fired to session's rolling state, and appends one record per violation to log.

`UserPromptSubmit` runs `cope-gate --refresher`. It reads rolling state — not violations log — and injects card items gated on what keeps firing, naming counts. Session with no history yet falls back to standing `CONTINUE TEST`. Refresher stays quiet until last injection ages past `--refresh-every`. Mid-session text gets chosen from measured output rather than fixed in advance, which pasted `CLAUDE.md` cannot do. One mechanism, no guarantee: A/B in repo separates variants of refresher window, not refresher from nothing.

`PreToolUse` runs `cope-gate --pretool`. It scores `description`, `body` or `content` field an external write posts, matched against Linear save tools named in settings block. Warn-only: it returns `additionalContext`, never a `permissionDecision`, so call goes through and model learns what prose scored. It writes no session state. It scores in external lane, dropping `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end`.

## Why Effigy Notation

[effigy](https://github.com/justinstimatze/effigy) states character cards for game NPCs. cope uses it off-label. Three of its blocks do what prose gate needs: `POSTPROC` holds regex rules with warn action applied after generation, `WRONG` holds anti-pattern beside its replacement, and `TEST` holds named question with fail and pass examples — so a rule names a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers same problem other way round, and suits reader finding this one too blunt. cope bans: rule fires or rule stays silent, card says never, register locks once you pick it. basanite measures — lemma frequency against baseline over real transcripts — and reports what you lean on lately, leaving judgement to you. Its own README calls that awareness rather than prohibition. Which one fits asks about mood more than correctness: heavy hand suits a habit annoying you today, moving measurement suits watching drift rather than legislating it. They compose — different hooks, no shared state — so run both.

Third project, different axis again: [caveman](https://github.com/JuliusBrussee/caveman), by another author, compresses agent replies to cut output tokens. cope shapes prose, basanite tracks vocabulary, caveman shortens. Reader wanting fewer tokens rather than different structure goes there.

## Rules

Voicing rules, from shipped card and from Go:

- `flip` — not-A-but-B flip in its common surface forms, including not-only-but-also and inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at top of it among every writer in study, human or machine. Action: warn.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k heaviest measured lean in this register. Say what thing carries instead. Action: warn.
- `worth_noting` — announces that something deserves attention instead of letting it earn attention, at 6.5 per 1k. Action: warn.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating a content word across joint: balanced two-beat.
- `apology` — reply performs contrition instead of stating correction and moving on.
- `self_postmortem` — reply turns to account for its own errors, story reader never asked for.
- `announced_length` — reply announces own length rather than cutting it.
- `cross_turn_repeat` — turn of phrase this reply shares with several earlier ones in same session. Only rule reading window rather than reply, so it stays silent until session holds history.

Structure rules, all compiled:

- `labelled_opening` — prose paragraph opening on short verbless fragment rest of it unpacks; ordinal counts as label. List blocks and paragraphs under twelve words skip. Bolded form skips too — card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so opener written as bold label stays unpoliced on purpose.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths carry coefficient of variation below `--min-cv`.
- `ask_not_last` — question or request for reader sitting in earlier block while reply carries on past it.
- `dangling_end` — open problem named in closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` — open problem landing after last question or offer, burying decision point above it.
- `forked_end` — two or more things to act on in closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding another.
- `unverified_done` — says work stands done with nothing on page that could have shown it: no command, no count, no file.
- `loop_ask` — unattended run ends by asking, so answer lands in a log and next iteration reads question as instruction to itself.

Grouping tells reader something usable about implementation. `POSTPROC` pattern matches a span of text, so it describes wording and nothing else. Every voicing rule needing more than a pattern got written in Go beside structure rules. Shipped cope card therefore carries 3 `POSTPROC` rules. Reader expecting long list of banned phrases finds that list in basanite instead, by design: cope bans, basanite measures.

Structure rules vary in one place — not by card, by who reads turn. Lanes:

- interactive lane — every turn that runs outside a loop. Somebody waits at terminal and ending carries their next decision.
- loop lane — prompt opening turn was `/loop` or `/goal`, or sentinel a dynamic-pacing loop sends itself. Drops `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end`. Adds `unverified_done`, `loop_ask`. Nobody reads yet. Report correctly naming what it left open and stopping would fail three of dropped rules, and question in it lands in log where next iteration reads it as instruction to itself. Claim check replaces them: report saying work stands done has to say what it ran.
- external lane — `PreToolUse` entry scoring prose an external write posts, rather than a reply. Drops `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end`, swaps nothing in. Ticket holds a reader and no ending they answer. Somebody reads it days later, outside session, and that condition describes every rule surviving drop.

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
| `--pretool` | `false` | PreToolUse entry: score the prose an external write is about to post, warn-only |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What Lands On Disk

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | `0600` | one JSON record per violation, carrying matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | `0600` | empty file whose mtime forms refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | `0600` | rolling record mid-session injection gets chosen from: turn count, characters, which rules fired over 20-turn window. No prose, only rule names and counts |

Log quotes replies back. Plainly: matched text plus surrounding context lands on disk in `violations.jsonl`. Session state stores no prose.

## Editing The Card

effigy owns `.effigy` grammar. `make rules` regenerates `card/rules.json`. `make check-rules` runs in CI, so enforced rules and injected rules cannot drift apart.

`NEVER` budget stands at 10. Anything over budget gets reported at load rather than dropped in silence. Budget charges against each injection separately, not against card file: one path prints always-on rules, refresher prints evidence-gated ones, and no code path renders their union. So card may hold more `NEVER` rules in total than budget and still stay healthy. Rules that really do get discarded unrendered appear in a list of their own, empty under a healthy card.

Two card-authored forms.

`@gate <rule_id> off — <why>`, one per line in card header. `@shape <id>: <selector> <predicate> — <why>`, one per line in card header.

`@shape` vocabulary, exactly:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

`@gate` takes rule id gate already holds. `@shape` takes id colliding with none. Wrong id gets reported at load rather than ignored. Reason after dash stays required in both forms, since rule a card wrote and rule a card refused read as unreviewable without one. Declined rule still runs; only this card's score drops it, so backfill still reports what it would have caught. `@shape` violation reports in card's own words rather than in any sentence binary supplies.

`card/demo/handoff.effigy` shows the assertion form with a measured number: `readable_cold`, `last paragraph words <= 60`, and 60 comes from 43,155 assistant replies where closing block runs 33 words at median and 56 at p90.

## Calibrating

`cope-gate --backfill` scores whole session transcript at once, and that command chose the rules. `tools/backfill-sweep.sh` runs it over N largest transcripts found anywhere under `~/.claude/projects` — not one per project.

Watch hits per character. Share of turns hit tracks how long turns ran, so it measures length more than prose. [MEASUREMENTS.md](MEASUREMENTS.md) holds rates.

## Known Limits

`labelled_opening` runs no tagger. It matches a fragment shape, so it misfires on prose a tagger would clear.

`ask_not_last` says nothing about order among several asks. It sees one ask sitting early and reply continuing past it.

Hit rate runs roughly four fifths structure. A/B run found that four fifths tracks what a reply was for rather than how it got written, so it describes output and judges nothing. Discrimination test holds judgement, and covers voicing only.

Largest limit sits where two things above named it. Card declines a built-in rule and writes one of its own, and vocabulary it writes in counts words and sentences and asks whether block poses question. Compiled rules remain only home for a check like `clause_symmetry`. Card wanting something outside that vocabulary and outside a `POSTPROC` regex still finds nowhere to put it.

Both directions amount to card marking own homework. Decline lowers that card's score; assertion raises it. Read both with reason attached — which is why syntax demands one.

[MEASUREMENTS.md](MEASUREMENTS.md) holds runs and caveats.

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
