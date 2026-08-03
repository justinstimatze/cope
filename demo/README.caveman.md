# cope

cope — foundry word for upper half of mould, half carrying shape cast into — ships one opinionated card, live from install, scoring every reply on two jobs, sound of its sentences and shape of its handback; card sits in one file, so you edit it or swap it.

Read [demo/README.md](demo/README.md) first: every file there holds this page again from another card, same prompt and same facts, card as only difference — two of them side by side teach what card does faster than rest of this page explains it, and [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes point in one glance, written from card that instructs every tic this model shows.

## Two things

Voicing means sound of sentences. Register, rhythm, diction, what paragraph does with detail, where flair earns license. Voicing lives in card — VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC. Swap card, swap every word of it. Blind discrimination test measures this half.

Structure means shape of reply as thing reader uses. Where decision sits. Whether ending gives "continue" something to point at. Whether ask lands last. Whether claim of finished work carries anything that could show it. Structure compiles into binary, into `internal/scan`, and runs same under every card.

Two examples, one per axis.

Two clauses of near-equal length, comma between them, one content word repeated across joint — voicing fault. Reply names open problem in its last paragraph, then stops, and "continue" points at nothing — structure fault. One reply passes one axis and fails other. Keep axes apart for that reason.

Card reaches two ways into structure half. Card declines built-in rule with `@gate <rule_id> off — <why>`, one line per rule in card header: `card/demo/lecturer.effigy` declines clause_symmetry and dangling_end, since its VOICE block asks for balanced landing and arriving close those two rules catch, and card lost points for obeying itself. Card states own structural rule with `@shape <id>: <selector> <predicate> — <why>`: `card/demo/handoff.effigy` asserts readable_cold — last paragraph words <= 60 — since its peak asks reader to re-enter cold and read last block alone, and no built-in rule checks that.

Boundary sits close. `@shape` counts words, counts sentences, asks whether block poses question. Nothing beyond that. Card wanting check outside that vocabulary and outside POSTPROC regex holds nowhere to put it. MEASUREMENTS.md carries runs and carries reasons those numbers stop where they stop.

## Why instruction alone fails

Global CLAUDE.md bans "not A, it's B" flip. Model reads that file every turn. Flip surfaced twice inside session that built this gate, with ban itself as topic of that session. Name surface form, move hides in variant.

Structure fault carries different cause and different cost. Reply ends, reader finds nothing to answer, one round trip burns. No instruction bans that shape, since shape holds no phrase to ban.

Flip stands as anecdote about one rule. Claim rests elsewhere: blind discrimination test. Reader sees one voice's own description of itself, sees two replies, picks which reply came from that voice. MEASUREMENTS.md holds rate and caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick style under `/config` -> Output style.

`--setup` emits output style, wires two hooks into `~/.claude/settings.json` with absolute paths, prints one step left. It backs settings file up first. It adds only missing keys, so second run changes nothing. It leaves every other key alone, other tools' hooks on same events included. It refuses to touch settings file that fails to parse. `--setup --dry-run` prints what it would change and writes nothing.

Prefer hands off your settings: `cope-gate --output-style` writes loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits different one. Select same way, under `/config` -> Output style — Claude Code v2.1.91 removed standalone `/output-style` command — or set `"outputStyle"` in `.claude/settings.local.json`.

Style loads once per session start. New selection, or freshly emitted card, applies at next session or after `/clear`. Running conversation shows nothing new.

Output style sits at end of system prompt, and harness re-reminds model of it mid-conversation. Card lands there for that reason.

Hooks come next, and card no longer arrives through one. Hooks block of `~/.claude/settings.json`:

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

Stop scores reply just written. UserPromptSubmit restates, mid-session, rules that keep firing — file written once cannot do that. Voice works without either hook. Hooks carry measurement half.

Both commands assume `go install` target directory on PATH inside environment where hook runs; hook that quietly does nothing usually wants absolute path instead. Clone builds with `make install` and needs no effigy checkout, since `card/rules.json` sits committed and compiled in. `--inject` remains as superseded delivery for anyone who wants old path, and it stands down on its own once cope output style goes active.

## Writing in another voice

Gate reads `.effigy` directly. Card stays usable as written — no Python, no effigy checkout.

Drop card in `$XDG_CONFIG_HOME/cope/cards`, then reach it with `--card <name>` or `COPE_CARD`. `--rules` takes path from anywhere. `make cards` installs demo set. Name that resolves to nothing raises error and injects nothing, rather than letting session write in shipped voice while its config names another.

Read `card/demo/lecturer.effigy` first: it differs from shipped card on register alone, and discrimination run measured exactly that pair. Card changes voicing axis, per **Two things** above.

Shortest route past prose: [demo/README.md](demo/README.md), where every file holds this page written again from another card, same prompt and same facts, voice as only variable.

One exception in that directory: `card/demo/handoff.effigy` states hypothesis, not voice, keeping shipped card's handoff rules and dropping everything about prose, and it wants `make pairs` against full card rather than render. `make cards` installs it with rest, so anyone listing installed cards meets it and should know it names no register.

## What hooks add

Stop scores reply and appends fired rule ids to session's rolling state, plus one record per violation to log.

UserPromptSubmit reads that state — not violations log — and injects card items gated on what keeps firing, naming counts. Mid-session text comes from measured output, not from fixed list. Pasted CLAUDE.md does nothing of that kind.

One mechanism, no guarantee. A/B run in this repo separates arms, not refresher from bare session.

SessionStart `--inject` stands superseded and off by default; `--setup` wires it nowhere.

## Why effigy notation

effigy names character-card notation for game NPCs, used here off-label. Three blocks do work prose gate needs: POSTPROC holds regex rules with warn action applied after generation, WRONG holds anti-pattern beside its replacement, TEST holds named question with fail example and pass example — so rule names move rather than one wording of that move.

basanite answers same problem from other side, and fits other mood. cope bans: rule fires or stays quiet, card says never, register locks once you pick it. basanite measures — lemma frequency against baseline across real transcripts — reports what you lean on lately, leaves judgement with you, and calls that awareness rather than prohibition. Heavy hand suits habit annoying you today. Moving measurement suits drift you would rather watch. Both compose: separate hooks, no shared state.

caveman ([github.com/JuliusBrussee/caveman](https://github.com/JuliusBrussee/caveman)), separate project by different author, compresses agent replies to cut output tokens. Third axis. Fewer tokens, rather than different structure — go there.

## Rules

Voicing:

- **clause_symmetry** — comma- or semicolon-joined clauses of near-equal length repeating content word across joint: balanced two-beat.
- **apology** — reply performs contrition instead of stating correction and moving on.
- **self_postmortem** — reply turns to account for its own errors, story reader never asked for.
- **announced_length** — reply announces own length rather than cutting it.
- **cross_turn_repeat** — turn of phrase this reply shares with several earlier ones in same session. Only rule reading window rather than reply, so it stays silent until session holds history.

Structure:

- **labelled_opening** — prose paragraph opening on short verbless fragment that rest unpacks; ordinal counts as label. Skips list blocks and paragraphs under twelve words. Skips bolded form too: card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as something they wanted, so bold opener stays unpoliced on purpose.
- **paragraph_uniformity** — four or more prose paragraphs whose lengths hold coefficient of variation below `--min-cv`.
- **ask_not_last** (interactive) — question or request for reader sitting in earlier block while reply carries on past it.
- **dangling_end** (interactive) — open problem named in closing blocks, no question, no offer, no explicit all-clear anywhere, so "continue" points at nothing.
- **buried_decision** (interactive) — open problem landing after last question or offer, burying decision point above it.
- **forked_end** (interactive) — two or more things to act on in closing blocks, nothing marking first, so "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, bare deference tags like "your call" read as continuation of decision above rather than second one.
- **unverified_done** (loop) — says work finished with nothing on page that could show it: no command, no count, no file.
- **loop_ask** (loop) — unattended run ends by asking, answer lands in log, next iteration reads question as instruction to itself.

Grouping tells you something usable about implementation. POSTPROC pattern matches span of text, so pattern describes wording and nothing else. Every voicing rule needing more than pattern went into Go beside structure rules. Shipped card therefore holds 0 POSTPROC rules. Long list of banned phrases lives in basanite, on purpose.

Structure rules vary in one place only — not by card, by who reads turn.

Interactive lane takes any turn that skips loop. Somebody waits at terminal, and ending carries decision about what happens next.

Loop lane takes turn opened by `/loop` or `/goal`, or by sentinel that dynamic-pacing loop sends itself. It drops ask_not_last, forked_end, dangling_end, buried_decision. It adds unverified_done and loop_ask. Nobody reads yet: report that correctly names open work and stops would fail three dropped rules, and question inside it lands in log where next iteration reads it as instruction to itself. Claim check replaces them — report saying work finished states what it ran.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through arms, record which arm wrote each turn |
| `--ab-arms` | `(empty)` | comma-separated arms to rotate through, implying `-ab` (default inject,hold; positive as third) |
| `--ab-report` | `(empty)` | read turn log, report arms; `-` reads default path |
| `--author-docs` | `false` | print prompt for writing this repo's docs: card, introspected facts, sections |
| `--backfill` | `(empty)` | score every assistant turn in this transcript, exit |
| `--block` | `false` | exit 2 on violation whose action reads reject (default warn-only) |
| `--card` | `(empty)` | name of installed card to write in, from cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list installed cards with aim each one states, exit |
| `--check` | `(empty)` | score prose file against card, exit; `-` reads stdin |
| `--describe` | `false` | print card's voice as target to recognise: aim and register, machinery dropped |
| `--display` | `false` | MessageDisplay entry: rewrite what reader sees, leave transcript alone |
| `--display-preview` | `false` | read prose on stdin, print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change, write nothing |
| `--inject` | `false` | print card as prompt text for SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write card to `~/.claude/output-styles` as Claude Code output style, which puts it in system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory to write output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of last card or refresher injection before refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject compact reminder once last injection ages |
| `--render-arm` | `(empty)` | print mid-session reminder one arm would inject, exit |
| `--render-for` | `(empty)` | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | `(empty)` | render `-render-arm` as given lane sees it: interactive (default) or loop |
| `--rules` | `(empty)` | read card from this `.effigy` or `.json` file instead of one built into binary |
| `--setup` | `false` | emit output style, wire hooks, print one step left |
| `--version` | `false` | print version, exit |

## What lands on disk

| Path | Mode | Holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | 0600 | one JSON record per violation, carrying matched text plus about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | 0600 | empty file whose mtime drives refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | 0600 | rolling record mid-session injection draws from: turn count, characters, rules fired over 20-turn window. No prose, rule names and counts only |

Log quotes your replies back at you.

## Editing card

effigy owns `.effigy` grammar. `make rules` regenerates `card/rules.json`. CI runs `make check-rules`, so enforced rules and injected rules cannot drift apart.

NEVER budget stands at 10. Anything over budget gets reported at load rather than dropped in silence. Budget charges against each injection separately, never against card file: SessionStart prints always-on rules, refresher prints evidence-gated ones, and no code path renders their union — so card may hold more NEVER rules in total than 10 and stay healthy. `never_rules_over_budget` names rules truly discarded unrendered, and it sits empty on healthy card.

Two card-authored forms.

Decline built-in rule:

```
@gate <rule_id> off — <why>
```

State own rule:

```
@shape <id>: <selector> <predicate> — <why>
```

`@shape` vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

`@gate` needs rule id gate already holds. `@shape` needs id that collides with none. Wrong id gets reported at load, never ignored. Reason after dash stays required in both forms, since rule card wrote and rule card refused both read as unreviewable without one.

Declined rule keeps running. Only that card's score drops it, so backfill still reports what it would have caught. `@shape` violation gets reported in card's own words, never in sentence binary supplies.

## Calibrating

`cope-gate --backfill` scores whole session transcript at once, and rule set came out of that. `tools/backfill-sweep.sh` runs it over N largest transcripts found anywhere under `~/.claude/projects` — largest, not one per project.

Watch hits per character. Share of turns hit tracks turn length instead.

MEASUREMENTS.md holds rates.

## Known limits

Axis split organises them.

labelled_opening does no tagging. ask_not_last says nothing about order among several asks.

Hit rate runs roughly four fifths structure. A/B run found those four fifths tracking what reply was for rather than how somebody wrote it, so hit rate describes output and judges nothing. Judgement lives in discrimination test, and that test covers voicing alone.

Largest limit stands where **Two things** named it, and stands smaller now than before. Card declines built-in rule and writes one of its own. Vocabulary it writes in counts words, counts sentences, asks whether block poses question. Compiled rules stay only home for check like clause_symmetry. Card wanting check outside that vocabulary and outside POSTPROC regex still holds nowhere to put it.

Both directions leave card marking own homework. Decline lowers that card's score. Assertion raises it. Read each one with its reason attached — syntax demands reason for that purpose. MEASUREMENTS.md.

## Layout

| Path | What |
| --- | --- |
| `card/claude_voice.effigy` | shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | hook binary |
| `internal/scan/` | structure rules, card's regex rules, card renderer |
| `internal/effigy/` | `.effigy` reader, so card stays usable as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane wrote each turn |
| `replay/` | blind-pairs and discrimination harnesses, plus own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what ran, on how much text, what it said |

---

`tools/generate_readme.py` wrote this README from prompt `cope-gate --author-docs` emits, checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
