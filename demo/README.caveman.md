# cope

cope ships one opinionated card — cope names upper half of foundry mould, half carrying shape cast into — card runs from install, scores two jobs at once, how prose sounds and how reply ends, and card sits in one file you edit or swap.

Read [demo/README.md](demo/README.md) first: every file in that directory holds this same README written again from different card, same prompt, same facts, so card stands as only difference between them, and reading two against each other shows what card does faster than rest of this page explains. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes point in one glance, written from card that instructs every tic this model measures, hard going by design.

## Two jobs

Voicing covers sound of sentences: register, rhythm, diction, what paragraph does with detail, where flair earns place. Voicing lives whole inside card — `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES`, `POSTPROC`. Swap card, swap every word of it. Blind discrimination test measures this half.

Structure covers shape reader must use: where decision sits, whether ending gives `continue` something to point at, whether ask lands last, whether claim of finished work carries anything that could show it. Structure compiles into binary. Same rules fire under every card, varying only by lane.

One instance each. Sentence reaching for balanced two-beat — two clauses, near-equal length, one content word repeated across joint — fails voicing. Reply naming open problem in last paragraph, then stopping, leaves `continue` nothing to point at, and fails structure. Same reply passes one axis and fails other, which gives reason to keep axes apart.

Card reaches into structure two ways. Card declines built-in rule with `@gate <rule_id> off — <why>`, because card whose `VOICE` block asks for landing that `clause_symmetry` catches lost points for obeying itself. Card writes own rule with `@shape <id>: <selector> <predicate> — <why>`, because card promising reader can re-enter cold and read closing block alone had nowhere to state that promise as check.

`@shape` vocabulary counts words and sentences, asks whether block poses question, and stops there. Check like `clause_symmetry` stays in compiled rules. Card wanting check outside that vocabulary and outside `POSTPROC` regex still finds nowhere to put it.

## Three failures

Most readers arrive holding first failure backwards. Global `CLAUDE.md` looks like system prompt. It sits elsewhere: harness attaches it as one message on first turn, and conversation buries that message under everything written after. Output style goes into system prompt itself, and harness re-reminds model of style as conversation runs. Move one card between those two places, change no word of card, and prose changes. That move explains most of why cope works. Run and figures sit in [MEASUREMENTS.md](MEASUREMENTS.md).

Second failure: instruction alone leaves phrasing where it stood. Global `CLAUDE.md` banning "not A, it's B" flip gets read every turn. Flip still appeared twice in session that built this repo, while ban held topic of that session. Name one surface form, move slides into variant.

Third failure belongs to structure, and no instruction reaches it. Ending that hands reader nothing to answer costs whole round trip. Nobody wrote that ending as phrasing habit, so no ban on phrasing removes it.

Flip stands as anecdote about one rule. Claim rests on blind discrimination test: reader sees only voice's own description of itself, then picks which of two replies came from that voice. [MEASUREMENTS.md](MEASUREMENTS.md) holds rate and caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

First command builds `cope-gate`. Second command does whole install: emits output style, wires three hooks into `~/.claude/settings.json` with absolute paths, prints one step left. `--setup` backs settings file up first, adds only missing keys so second run changes nothing, leaves every other key alone including other tools' hooks on same events, and refuses to touch settings file that fails to parse. `--setup --dry-run` prints what it would change and writes nothing.

Then pick style under `/config` -> Output style. Entry reads `claude_voice`, shipped card's id, and never word cope — scan that menu for cope, find nothing, walk away with tool switched off.

Prefer settings file untouched? `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits different card. Wire hooks by hand after that.

Claude Code reads style once at session start. New selection or re-emitted card applies at next session, or after `/clear`.

Output style sits at end of system prompt, where harness keeps reminding model of it. Card lands there for that reason.

Hooks block of `~/.claude/settings.json`:

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

`Stop` buys score on every reply after writing. `UserPromptSubmit` buys mid-session restatement of rules that keep firing, which file written once cannot do. Voice works without both hooks; hooks carry measurement half.

Commands above assume `go install` target directory on `PATH` inside environment where hook runs. Hook doing nothing at all usually wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since `card/rules.json` ships committed and compiled in. `--inject` survives as superseded delivery for anyone wanting old path, and stands down on its own when cope output style runs.

## Writing in another voice

Gate reads `.effigy` directly, so card works as written — no Python, no effigy checkout. Drop card in `$XDG_CONFIG_HOME/cope/cards`, reach it with `--card <name>` or `COPE_CARD`. `--rules` takes path from anywhere. `make cards` installs demo set. Name resolving to nothing raises error and injects nothing, rather than letting session write in shipped voice while config names another.

Read `card/demo/lecturer.effigy` first. It differs from shipped card on register alone, and discrimination run measured exactly that pair. Card changes voicing, described under two jobs above; structure holds still. [MEASUREMENTS.md](MEASUREMENTS.md) carries numbers.

Shortest route to seeing what card does without writing one: [demo/README.md](demo/README.md), where every file holds this page rewritten from different card, same prompt, same facts.

`card/demo/handoff.effigy` breaks pattern in that directory: hypothesis rather than voice, keeping shipped card's handoff rules and dropping everything about prose, meant for `make pairs` against full card rather than for rendering. `make cards` installs it beside rest, so anyone listing cards finds it there and should skip it as register.

## What hooks add

`Stop` scores reply just written, appends which rules fired to session's rolling state, and appends one record per violation to log.

`UserPromptSubmit` reads that rolling state — never violations log — and injects card items gated on what keeps firing, naming counts. Session with no history yet falls back to standing `CONTINUE TEST`. Hook stays quiet until last injection ages past `--refresh-every`.

Mid-session text therefore comes from measured output rather than from fixed paste. That names one mechanism, and nothing more: A/B run in this repo does not separate refresher from no refresher.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) writes character cards for game NPCs. cope uses it off-label, since three effigy blocks do what prose gate needs: `POSTPROC` holds regex rules with warn action applied after generation, `WRONG` holds anti-pattern beside its replacement, and `TEST` holds named question with fail and pass examples — which lets rule name move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers same problem from other end, and suits anyone finding this one too blunt. cope bans: rule fires or stays quiet, card says never, register locks once you pick card. basanite measures instead — lemma frequency against baseline over real transcripts — so it reports what you lean on lately and leaves judgement with you. Its own README calls that awareness rather than prohibition. Choice tracks mood more than correctness: heavy hand suits habit annoying you today, moving measurement suits watching drift. Both compose, different hooks, no shared state, so run both.

[humanizer](https://github.com/blader/humanizer) rewrites AI-sounding prose against 35 patterns drawn from Wikipedia's "Signs of AI writing", page WikiProject AI Cleanup maintains. Call humanizer on text, get rewrite back. cope fires at hook, scores what already exists, edits nothing. humanizer holds wider pattern list; anyone wanting rewrite rather than score goes there. Two tools disagree on formatting by design: cope's `bold_label` rule banned humanizer's bold mini-headings until 52 blind pairs put bold and bullets among three things deciding reply for this repo's reader, so cope deleted rule rather than tuning it.

[caveman](https://github.com/JuliusBrussee/caveman) compresses agent replies to cut output tokens — fourth axis, separate project, different author. cope shapes prose, basanite tracks vocabulary, humanizer rewrites, caveman shortens. Anyone wanting fewer tokens rather than different structure goes there instead.

## Rules

Voicing rules:

- `flip` — not-A-but-B flip in common surface forms, including not-only-but-also and inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude top of it among every writer in study, human or machine.
- `load_bearing` — reflexive intensifier standing in for important or central, at 25.6 per 1k heaviest measured lean in this register. Say what thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating content word across joint: balanced two-beat.
- `apology` — reply performs contrition instead of stating correction and moving on.
- `self_postmortem` — reply turns to account for its own errors, story reader never asked for.
- `announced_length` — reply announces own length rather than cutting it.
- `cross_turn_repeat` — turn of phrase this reply shares with several earlier replies in same session. Only rule reading window rather than reply, so it stays silent until session grows history.
- `repeated_opening` — three or more sentences in one reply opening on same two words. `cross_turn_repeat` reads session window for construction reused across turns; this one reads reply against itself, and two openings pass, since two makes rhythm.
- `fragment_run` — three consecutive sentences of five words or fewer with no finite verb in any. One fragment marks emphasis, and this repo's own register runs full of them; run of three gives staccato blind judges read as generated. Neither clipped demo card trips it, so neither declines it; card wanting that run says so with `@gate`.

Structure rules:

- `labelled_opening` — prose paragraph opening on short verbless fragment that rest unpacks; ordinal counts as label. Skips list blocks, skips paragraphs under twelve words, skips bolded form — card dropped `bold_label` in July 2026 after blind readers named bold and bullets as things they wanted, so bold opener stays deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths hold coefficient of variation below `--min-cv`.
- `ask_not_last` — question or request for reader sitting in earlier block while reply carries on past it.
- `dangling_end` — open problem named in closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to point at.
- `buried_decision` — open problem landing after last question or offer, burying decision point above it.
- `forked_end` — two or more things to act on in closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding another.
- `unverified_done` — says work finished with nothing on page that could show it: no command, no count, no file.
- `loop_ask` — unattended run ends by asking, so answer lands in log and next iteration reads question as instruction to itself.
- `echoed_heading` — heading of two or more content words whose first sentence below repeats every one, spending line to say what heading said.

Grouping tells implementers something usable. `POSTPROC` pattern matches span of text, so pattern describes wording and nothing beyond wording. Every voicing rule needing more than pattern went into Go beside structure rules. Shipped card therefore carries 3 `POSTPROC` rules. Anyone expecting long list of banned phrases should read basanite instead, where that list lives on purpose.

Structure rules vary in one way only, by lane — by who reads turn, never by card.

Interactive lane covers every turn that no loop opened, since somebody waits at terminal and ending decides what happens next.

Loop lane covers turns opened by `/loop` or `/goal`, or by sentinel that dynamic-pacing loop sends itself. Loop drops `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end`, and adds `unverified_done` and `loop_ask`. Nobody reads yet: report that correctly names what it left open and stops would fail three dropped rules, and question inside it lands in log where next iteration reads question as instruction to itself. Claim check replaces them — report saying work finished has to say what ran.

External lane covers prose that `--pretool` scores before external write posts it. External drops same four and swaps nothing in. Ticket has reader and no ending that reader can answer, read days later by somebody outside session, which names condition every surviving rule was written for.

Reports cluster hits. `Stop`, `--check` and `--pretool` all print clusters. Three or more distinct rules landing on one paragraph gives breadth. One rule landing on one paragraph three or more times gives density. Both conditions holding, report prints breadth, since naming three rules tells reader to rewrite block rather than hunt one construction. Every rule fires alone and knows nothing of others, so flat report leaves reader working down hit by hit; three hits across three paragraphs make three small edits, three hits inside one paragraph make one paragraph to write again. Density half rests on measurement: `--check` over 107 tracked documents produced 114 `flip` hits, seven worth changing, and every one of seven showed as three in paragraph rather than as anything about form. Three sets floor on both conditions, since two rules on paragraph reads ordinary and two hits of one rule reads as coincidence reader spots unaided. Clustering changes nothing about what fires or how score works.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which variant wrote each turn |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; `positive` third) |
| `--ab-report` | (empty) | read turn log, report how each variant did; `-` reads default path |
| `--author-docs` | `false` | print prompt for writing this repo's docs: card, introspected facts, sections |
| `--backfill` | (empty) | score every assistant turn in this transcript, exit |
| `--block` | `false` | exit 2 on violation carrying action `reject` (default warn-only) |
| `--card` | (empty) | name of installed card to write in, from cards directory; also `$COPE_CARD` |
| `--card-from-sample` | (empty) | print prompt for writing card from this writing sample; `-` reads stdin |
| `--cards` | `false` | list installed cards with aim each one states, exit |
| `--check` | (empty) | score prose file against card, exit; `-` reads stdin |
| `--check-lane` | (empty) | score `--check` in given lane: `interactive` (default), `loop`, or `external` |
| `--describe` | `false` | print card's voice as target to recognise: aim and register, without machinery |
| `--display` | `false` | `MessageDisplay` entry: rewrite what reader sees, leaving transcript alone |
| `--display-preview` | `false` | read prose on stdin, print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change, write nothing |
| `--inject` | `false` | print card as prompt text for `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write card to `~/.claude/output-styles` as Claude Code output style, which puts card in system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write output style into (default `~/.claude/output-styles`) |
| `--pretool` | `false` | `PreToolUse` entry: score prose that external write posts next, warn-only |
| `--refresh-every` | `30m0s` | minimum age of last card or refresher injection before refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject compact reminder once last injection ages |
| `--render-arm` | (empty) | print mid-session reminder one variant would inject, exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as given lane sees it: `interactive` (default) or `loop` |
| `--rules` | (empty) | read card from this `.effigy` or `.json` file instead of one built into binary |
| `--setup` | `false` | emit output style, wire hooks, print one step left |
| `--version` | `false` | print version, exit |

## What lands on disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600`, holds one JSON record per violation. Record carries matched text plus roughly 70 characters either side, so log quotes your replies back at you.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600`, holds nothing; its mtime runs refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600`, holds rolling record that mid-session injection draws from: turn count, characters, which rules fired over 20-turn window. Rule names and counts only, no prose.

## Editing card

effigy owns `.effigy` grammar. `make rules` regenerates `card/rules.json`. CI runs `make check-rules`, so enforced rules and injected rules cannot drift apart.

`NEVER` budget stands at 10. Rules over budget get reported at load rather than dropped in silence. Budget charges each injection separately rather than charging card file: `SessionStart` prints always-on rules, refresher prints evidence-gated ones, and no code path renders union of both — so card holding more `NEVER` rules in total than budget stays healthy.

Card writes rules two ways.

`@gate <rule_id> off — <why>`, one per line in card header. Id must name rule gate already has. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, since its `VOICE` block asks for balanced landing and arriving close that those two rules catch.

`@shape <id>: <selector> <predicate> — <why>`, one per line in card header. Id must collide with nothing gate already has. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — since its peak asks reader to re-enter cold and read closing block alone, and no built-in rule checks whether that block reads that way. 60 comes from measurement: across 43,155 assistant replies, closing block runs 33 words at median and 56 at p90.

Wrong id gets reported at load, never ignored. Reason after dash stays required in both forms, since rule card wrote and rule card refused read equally unreviewable without one. Declined rule still runs — only that card's score drops it — so `--backfill` still reports what declined rule would have caught. `@shape` violation prints in card's own words rather than in any sentence binary supplies.

## Calibrating

`cope-gate --backfill` scores whole session transcript at once. Rules got chosen that way. `tools/backfill-sweep.sh` runs `--backfill` over N largest transcripts found anywhere under `~/.claude/projects`, which differs from one transcript per project.

Watch hits per character. Share of turns hit tracks how long turns ran, so that second number moves for reasons outside prose. [MEASUREMENTS.md](MEASUREMENTS.md) holds rates.

## Known limits

Axis split organises limits too.

`labelled_opening` runs no tagger, so it flags shapes rather than parts of speech. `ask_not_last` says nothing about order among several asks.

Hit rate runs roughly four fifths structure. A/B run found that four fifths tracks what reply was for rather than how somebody wrote it, so hit rate describes output and judges nothing. Judgement lives in discrimination test, which covers voicing alone.

Largest limit sits where two jobs above already put it, and it stands smaller now than before: card declines built-in rule and writes one of its own, in vocabulary counting words and sentences and asking whether block poses question. Compiled rules stay only home for check like `clause_symmetry`. Card wanting check outside that vocabulary and outside `POSTPROC` regex still finds nowhere to put it.

Both directions leave card marking own homework. Decline lowers that card's score; assertion raises it. Read both with reason attached — which explains why syntax demands one. [MEASUREMENTS.md](MEASUREMENTS.md) holds rest.

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | hook binary |
| `internal/scan/` | structure rules, card's regex rules, card renderer |
| `internal/effigy/` | `.effigy` reader, so card works as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane wrote each turn |
| `replay/` | blind-pairs and discrimination harnesses, plus their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what ran, on how much text, and what it said |

`tools/generate_readme.py` wrote this README from prompt `cope-gate --author-docs` emits, and `cope-gate --check` checked it.

MIT. justin@justinstimatze.com
