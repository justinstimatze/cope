# cope

cope ships one opinionated card — live from install, whole product for most readers — scoring two separate jobs, how reply sounds and how reply ends, and keeps that card in one file you edit or swap; foundry calls upper half of mould cope, half carrying shape cast into.

Start at [demo/README.md](demo/README.md): every file there holds this README written again from different card, same prompt and same facts, so card counts as only variable, and reading two of them against each other shows what card does faster than rest of this page explains. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes that point in one glance — card instructing every tic this model measures out.

## Two things

Voicing means sound of sentences: register, rhythm, diction, what paragraph does with detail, where flair lands. Voicing lives in card. Swap card, swap every word of it. Blind discrimination test measures this half.

Structure means shape of reply as thing reader uses: where decision sits, whether ending gives "continue" something to point at, whether ask lands last, whether claim of finished work carries anything showing it. Structure compiles into binary. Same rules under every card, varying only by lane.

Sentence reaching for balanced two-beat — two clauses of matched length, content word repeated across joint — fails voicing. Reply naming open problem in its last paragraph, then stopping, leaves "continue" nothing to point at, and fails structure. One reply passes one axis and fails other. Reason to hold them apart.

Card reaches into structure two ways. Card declines built-in rule with `@gate <rule_id> off — <why>`, one per line in card header: card/demo/lecturer.effigy declines clause_symmetry and dangling_end, because its VOICE block asks for balanced landing and arriving close those two rules catch, and card obeying itself lost score for it. Card states structural rule of its own with `@shape <id>: <selector> <predicate> — <why>`, one per line in card header: card/demo/handoff.effigy asserts readable_cold — last paragraph words <= 60 — because its peak asks reader to re-enter cold and read last block only, and no built-in rule checked that commitment. @shape vocabulary counts words and sentences and asks whether block poses question, nothing more, so card wanting check outside that vocabulary and outside POSTPROC regex still finds nowhere to put it. MEASUREMENTS.md holds run and holds reasons its numbers carry no further.

## Problem

Instruction alone fails on phrasing. Global CLAUDE.md bans "not A, it's B" flip. Session that built this read that ban every turn. Flip appeared twice anyway, while ban stayed topic of conversation. Name surface form, move slides into variant.

Structure fails from other cause. Ending that leaves reader nothing to answer costs whole round trip, and no instruction bans it, because no phrase carries it.

Flip counts as anecdote about one rule. Claim rests on blind discrimination test: reader sees only voice's own description of itself, then picks which of two replies came from that voice. MEASUREMENTS.md holds rate and caveats. Blind preference runs sit outside this claim — both arms wrote under card, so they compare two ways of writing well and never watch voice swap.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick style under `/config` -> Output style. Claude Code v2.1.91 removed standalone `/output-style` command, so `/config` carries it; `.claude/settings.local.json` takes same choice as `"outputStyle"`.

`--setup` emits output style, wires two hooks into ~/.claude/settings.json with absolute paths, prints one step left. It copies settings file to backup first. It adds only missing keys, so second run changes nothing. It leaves every other key alone, other tools' hooks on same events included. It refuses to touch settings file that fails to parse. `--setup --dry-run` prints changes and writes nothing.

Prefer no tool near your settings: `cope-gate --output-style` writes loaded card to ~/.claude/output-styles/&lt;card&gt;.md, and `COPE_CARD=<name>` in front of it emits different one.

Harness reads style once at session start. New selection, or freshly emitted card, applies at next session or after `/clear`. Running conversation shows nothing new.

Output style lands at end of system prompt, and harness re-reminds model of it during conversation. Card lands there for that reason, and never landed through hook.

Hooks handle measurement now; card no longer arrives through one. Hooks block of ~/.claude/settings.json:

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

Stop scores reply just written and appends which rules fired to session's rolling state, plus one record per violation to log. UserPromptSubmit reads that rolling state and injects card items gated on what fires lately, naming counts. Voice works without either hook. These two form measurement half, and `cope-gate --inject` remains only as superseded delivery for anyone wanting old path, standing down on its own once cope output style goes live.

Commands above assume go install target directory on PATH inside environment hook runs in; hook that silently does nothing wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since card/rules.json sits committed and compiles into binary.

## Writing in another voice

Gate reads .effigy directly, so card works as written — no Python, no effigy checkout. Drop card in $XDG_CONFIG_HOME/cope/cards, then reach it with `--card <name>` or `COPE_CARD`. `--rules` takes path from anywhere. `make cards` installs demo set. Name resolving to nothing raises error and injects nothing, rather than session writing shipped voice while its config names another one.

Read card/demo/lecturer.effigy first. It differs from shipped card on register alone, and discrimination run measured that pair. Card changes voicing, per Two things above. MEASUREMENTS.md holds numbers.

[demo/README.md](demo/README.md) shows what card does without writing one: every file under demo/ holds this README again from different card, same prompt and same facts, voice as only variable.

card/demo/handoff.effigy stands apart in that directory as hypothesis rather than voice — it keeps shipped card's handoff rules, drops everything about prose, and wants `make pairs` against full card rather than rendering — and `make cards` installs it with rest, so reader listing cards finds it and reads it as no register to write in.

## What hooks do differently

Stop scores reply and records which rules fired. UserPromptSubmit reads that record and injects only items gated on firing counts, naming those counts, so mid-session text comes from measured output rather than from list fixed in advance. Pasted CLAUDE.md cannot do that.

One mechanism, no guarantee. A/B run in this repo separates nothing between refresher and no refresher.

SessionStart `--inject` stays superseded and off by default.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) writes character cards for game NPCs, and cope uses it off-label. Three blocks do what prose gate needs: POSTPROC holds regex rules with warn action applied after generation, WRONG holds anti-pattern beside its replacement, TEST holds named question with fail example and pass example — which lets rule name move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers same problem other way round, and suits reader finding this one too blunt. cope bans: rule fires or rule stays quiet, card says never, register locks when you pick it. basanite measures instead — lemma frequency against baseline over real transcripts — so it reports what you lean on lately and leaves judgement to you. Its own README calls that awareness rather than prohibition. Fit tracks mood more than correctness: heavy hand suits habit annoying you today, moving measurement suits watching drift instead of legislating it. They compose across different hooks with no shared state, so running both stays reasonable.

[caveman](https://github.com/JuliusBrussee/caveman), separate project by different author, compresses agent replies to cut output tokens. Reader wanting fewer tokens rather than different structure goes there.

## Rules

Voicing:

- **clause_symmetry** — comma- or semicolon-joined clauses of near-equal length repeating content word across joint: balanced two-beat.
- **apology** — reply performs contrition instead of stating correction and moving on.
- **self_postmortem** — reply turns to account for its own errors, story reader never asked for.
- **announced_length** — reply announces its own length rather than cutting it.
- **cross_turn_repeat** — turn of phrase this reply shares with several earlier ones in same session. Only rule reading window rather than reply, so it stays quiet until session holds history.

Structure:

- **labelled_opening** — prose paragraph opening on short verbless fragment that rest of it unpacks; ordinal counts as label. Skips list blocks and paragraphs under twelve words, and skips bolded form: card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as something they wanted, so opener written as bold label stays deliberately unpoliced.
- **paragraph_uniformity** — four or more prose paragraphs whose lengths show coefficient of variation below `--min-cv`.
- **ask_not_last** (interactive) — question or request for reader sitting in earlier block while reply carries on past it.
- **dangling_end** (interactive) — open problem named in closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to point at.
- **buried_decision** (interactive) — open problem landing after last question or offer, burying decision point above it.
- **forked_end** (interactive) — two or more things to act on in closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" read as continuing decision above rather than adding another.
- **unverified_done** (loop) — says work done with nothing on page that could have shown it: no command, no count, no file.
- **loop_ask** (loop) — unattended run ends by asking, so answer lands in log and next iteration reads question as instruction to itself.

Grouping cuts across implementation. POSTPROC pattern matches span of text, so pattern describes wording and nothing else, and every voicing rule needing more than pattern went into Go beside structure rules. Shipped card therefore holds 0 POSTPROC rules. Reader expecting long list of banned phrases finds that list in basanite, on purpose.

Structure rules vary in one place only, and card never moves them — reader of turn does. Interactive lane takes any turn outside loop, because somebody waits at terminal and ending holds their next decision. Loop lane takes turns opened by `/loop` or `/goal`, or by sentinel dynamic-pacing loop sends itself, and drops ask_not_last, forked_end, dangling_end, buried_decision while adding unverified_done and loop_ask. Nobody reads yet in that lane. Report correctly naming what it left open, then stopping, fails three dropped rules, and question inside it lands in log where next iteration reads it as instruction to itself. Claim check replaces them: report saying work done names what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | `(empty)` | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | `(empty)` | read a turn log and report the arms; - reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | `(empty)` | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | `(empty)` | name of an installed card to write in, from the cards directory; also $COPE_CARD |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | `(empty)` | score a prose file against the card and exit; - reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as --display would rewrite it |
| `--dry-run` | `false` | with --setup, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to ~/.claude/output-styles as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory to write the output style into (default ~/.claude/output-styles) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | `(empty)` | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | `(empty)` | comma-separated rule ids to render -render-arm against |
| `--render-lane` | `(empty)` | render -render-arm as the given lane sees it: interactive (default) or loop |
| `--rules` | `(empty)` | read the card from this .effigy or .json file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | 0600 | one JSON record per violation, carrying matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | 0600 | empty file whose mtime serves as refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | 0600 | rolling record mid-session injection draws from: turn count, characters, and which rules fired over 20-turn window. No prose, only rule names and counts |

Log quotes your replies back to disk.

## Editing card

effigy owns .effigy grammar. `make rules` regenerates card/rules.json. `make check-rules` runs in CI, so enforced rules and injected rules never drift apart.

NEVER budget sits at 10 rules. Anything over budget gets reported at load, never dropped in silence. Budget charges against each injection separately, not against card file: SessionStart prints always-on rules, refresher prints evidence-gated ones, and no code path renders their union — so card holds more NEVER rules in total than budget and stays healthy. Shipped card discards none unrendered.

Two card-authored forms:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

@shape vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

@gate takes rule id gate already carries. @shape takes id colliding with none of them. Wrong id gets reported at load, never ignored. Reason after dash stays required in both forms, since rule card wrote and rule card refused read as equally unreviewable without one. Declined rule still runs, and only this card's score drops it, so backfill still reports what that rule would have caught. @shape violation reports in card's own words rather than in any sentence binary supplies.

## Calibrating

`cope-gate --backfill` scores whole session transcript at once, and rules came out of it. `tools/backfill-sweep.sh` runs backfill over N largest transcripts found anywhere under ~/.claude/projects, which differs from one per project. Watch hits per character rather than share of turns hit: that second number tracks how long turns ran. MEASUREMENTS.md holds rates.

## Known limits

labelled_opening works as heuristic, not tagger. ask_not_last says nothing about order among several asks.

Hit rate runs roughly four fifths structure, and A/B run found those four fifths track what reply was for rather than how somebody wrote it — so hit rate describes output and judges nothing. Judgement lives in discrimination test, which covers voicing only.

Largest limit carries name of Two things above. Card declines built-in rule and writes rule of its own, and vocabulary it writes in counts words and sentences and asks whether block poses question. Compiled rules stay only home for check like clause_symmetry. Card wanting something outside that vocabulary and outside POSTPROC regex still finds nowhere to put it.

Both directions amount to card marking own homework. Decline lowers that card's score. Assertion raises it. Read both with reason attached, which explains why syntax demands one. MEASUREMENTS.md holds what ran.

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | hook binary |
| `internal/scan/` | structure rules, card's regex rules, card renderer |
| `internal/effigy/` | .effigy reader, so card works as written |
| `internal/transcript/` | Claude Code JSONL reader, plus which lane turn came from |
| `replay/` | blind-pairs and discrimination harnesses, and their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what ran, on how much text, and what it said |

tools/generate_readme.py wrote this README from prompt `cope-gate --author-docs` emits, and `cope-gate --check` checked it.

MIT. justin@justinstimatze.com
