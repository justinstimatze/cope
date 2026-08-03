# cope

cope ships one opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — that is switched on the moment it is installed and is the whole product for almost everyone who installs it, and it grades two different jobs rather than one, the sound of the sentences and the shape of the reply, and only after that does it start to matter that the card is a file you can edit or trade for another.

Before anything else on this page: [demo/README.md](demo/README.md). Every file in that directory is this page written again from a different card — same prompt, same facts, the card the only variable — and reading two of them against each other shows what a card does faster than the rest of this page manages to explain it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is the one that makes the point in a single glance, since it was written from a card that instructs every tic this model has been measured to have.

## Two jobs, not one

A house style looks like one thing, because it arrives as one document and gets read as one instruction. cope treats it as two, and the seam between them is the reason everything below is arranged the way it is.

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. All of it lives in the card, and swapping the card swaps every word of it — that is the half with a measured result standing behind it. Structure is something else entirely: where the decision sits, whether the ending gives *continue* something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. That half is compiled into the binary, so it is identical whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing failure. A reply that names an open problem in its last paragraph and then stops, leaving *continue* nothing to answer, is a structural one, and the very same reply can be immaculate on one axis and useless on the other. Which is the whole reason for keeping them apart.

The boundary between them used to be absolute, and it is the part that has most recently moved. A card can now decline a compiled rule it disagrees with — `@gate <rule_id> off — <why>`, one per line in the card header — and that exists because a card whose VOICE block asked for precisely what a rule catches was being marked down for obeying itself; `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` on exactly those grounds. A card can also assert a structural rule of its own — `@shape <id>: <selector> <predicate> — <why>`, one per line in the card header — and that exists because a card's own commitment about how a reply ends had nowhere to be checked. `card/demo/handoff.effigy` asserts `readable_cold`, last paragraph words <= 60, since its peak asks the reader to re-enter cold and read the closing block alone, and no compiled rule looks at whether that block can be read that way. The 60 is not a preference: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

What a card writes those rules in counts words and sentences, and asks whether a block poses a question. That is the entire vocabulary.

## What instruction does not fix

You have already tried the obvious thing, because everybody does — a global `CLAUDE.md` with the house rules in it, read at the top of every single turn. Mine bans the "not A, it's B" flip by name. The flip appeared twice in the session that built this tool, while the ban was the topic under discussion. Naming a surface form does not remove the move underneath it. It relocates the move into a variant the ban cannot match.

The second failure is not a phrasing habit at all, which is why no instruction was ever going to catch it: a reply that raises an open problem in its closing paragraph and then simply stops has handed the reader nothing to answer, and *continue* becomes a guess about which of several loose ends comes first. That costs a full round trip, spent on a question the previous turn was in the best position to ask.

The flip is an anecdote about one rule. The claim rests on a different instrument — a blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate, and the reasons to hold it loosely.

## Installing

Two commands and one menu choice, in this order.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

The last step lives in the client: `/config`, then **Output style**, then the card by name. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the route now; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

What `--setup` did, in the order it did it: emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left. Before touching anything it backs the settings file up. It adds only what is missing, so a second run changes nothing at all. It leaves every other key alone, including other tools' hooks sitting on the same events, and it refuses outright to edit a settings file that does not parse — and if you would rather see the diff first, `--setup --dry-run` prints what it would change and writes nothing.

Anyone who would sooner not have a tool editing their settings has the by-hand route. `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start. So a new selection, or a card you have just re-emitted, applies at the next session or after `/clear` — not in the conversation you are sitting in.

The reason it goes there rather than into a hook is placement: an output style lands at the end of the system prompt, and the harness re-reminds the model of it as the conversation runs on.

The hooks are a separate matter, and they are no longer how the card arrives.

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

Those two buy the measurement half. `Stop` scores the reply that was just written and records which rules fired; `UserPromptSubmit` restates, mid-session, the items gated on what has actually been firing — which a file written once at the top of the conversation cannot do. The voice itself works with neither of them installed.

Both commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing almost always wants the absolute path instead. From a clone, `make install` builds it, and no effigy checkout is needed, because `card/rules.json` is committed and compiled in. The old `SessionStart` delivery survives as `--inject`, superseded, off by default, and standing down on its own whenever a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, which means a card is usable exactly as written — no Python, no effigy checkout, no build step between the file and the session. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set. A name that resolves to nothing is a hard error, and nothing is injected — the alternative being a session that writes in the shipped voice while its config confidently names another one.

`card/demo/lecturer.effigy` is the one to read first, since it differs from the shipped card on register alone, and it is what the discrimination run measured. What a card can change is the voicing axis, for the reasons set out under *Two jobs, not one*, and [MEASUREMENTS.md](MEASUREMENTS.md) holds the numbers.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README again from a different card, same prompt and same facts, the voice the only thing varying.

One file under `demo/` is not a register at all: `card/demo/handoff.effigy` is a hypothesis, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered — and since `make cards` installs it with the others, anyone listing their cards will find it there.

## What the hooks add

`Stop` fires after the reply is written, scores it, and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` then reads that state — the state, not the violations log — and injects the card items gated on what has been firing lately, naming the counts as it goes. It falls back to the standing CONTINUE TEST while a session has no history, and stays quiet until the last injection has aged past `--refresh-every`.

The point of that arrangement is narrow and worth stating without decoration: the mid-session text is selected from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot manage. It is a mechanism, not a guarantee, and the A/B in this repo does not separate the refresher from no refresher.

`SessionStart --inject` remains, superseded and off by default, for anyone who wants the old delivery back.

## Why a character-card notation

[effigy](https://github.com/justinstimatze/effigy) is a card notation for game NPCs, used here entirely off-label, and three of its blocks turn out to be the three things a prose gate needs. POSTPROC is regex rules with a warn action applied after generation. WRONG holds an anti-pattern next to its replacement, so the card can show the swap rather than describe it. TEST holds a named question with a failing and a passing example, which is how a rule names a *move* instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem from the other end, and it is the one to reach for when this one is too blunt: cope bans, a rule fires or it does not, the card says never, and the register is settled the moment you pick a card — whereas basanite measures lemma frequency against a baseline over real transcripts and reports what you have actually been leaning on, leaving the judgement with you. Its own README calls that awareness rather than prohibition. Which of them fits is a question about mood more than correctness, they share no state and sit on different hooks, and running both is perfectly reasonable. A third project, [caveman](https://github.com/JuliusBrussee/caveman) by a different author, compresses agent replies to cut output tokens; anybody wanting fewer tokens rather than different structure should be over there instead.

## The rules

Grouped by job rather than by where they happen to be implemented.

**Voicing**

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure**

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks, an ordinal counting as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving *continue* nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering *continue* means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

That grouping tells you something the implementation would otherwise hide. A POSTPROC pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules — which is why this card ships no POSTPROC rules at all. Anyone arriving in expectation of a long list of banned phrases should know that the list lives in basanite, on purpose.

The structure rules do vary in one place, and it is not by card. It is by who is going to read the turn. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself, and it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision` while adding `unverified_done` and `loop_ask`. Nobody is reading yet in that lane: a report that correctly names what it left open and then stops would fail three of the dropped rules, and a question inside it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

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
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

Three files, all mode `0600`, all under `$XDG_STATE_HOME/cope`.

`violations.jsonl` holds one JSON record per violation, carrying the matched text and roughly 70 characters either side of it — which is to say the log quotes your replies back at you, and that is worth knowing before it accumulates. `refresher-<session-id>` is an empty file whose mtime is the refresher clock. `session-<session-id>.json` is the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired across a 20-turn window. No prose is stored in that last one, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot quietly drift apart.

The NEVER budget is 10, and anything over it is reported at load rather than dropped in silence. That budget is charged against each injection separately rather than against the card file, because SessionStart prints the always-on rules, the refresher prints the evidence-gated ones, and no code path anywhere renders their union — so a card may hold more NEVER rules in total than the budget and still be entirely healthy. The authoritative list of rules that really are being discarded unrendered is `never_rules_over_budget`, and on a healthy card it is empty.

Both card-authored forms live in the card header, one per line.

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is exactly this and nothing beyond it:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

For `@gate` the id has to be one the gate actually has; for `@shape` it must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both directions, since a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs — only this card's score drops it, so a backfill still reports everything it would have caught — and a `@shape` violation comes back in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores every assistant turn in a whole session transcript at once, and it is how these rules were chosen in the first place. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same thing as one per project, and mistaking the two will skew whatever you conclude.

The number to watch is hits per character, not the share of turns hit, because that second number mostly tracks how long the turns were. [MEASUREMENTS.md](MEASUREMENTS.md) has the rates.

## Known limits

The axis split organises these, so they divide the way everything else on this page has.

`labelled_opening` is not a tagger, and it will mistake things a tagger would not. `ask_not_last` has nothing to say about the ordering of several asks once they are all in the right region. The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was *for* rather than how it was written — which makes it a description of the output rather than a judgement of it. The judgement lives in the discrimination test, and that test covers voicing only.

The largest limit is the one *Two jobs, not one* already named, and it is worth stating as it now stands rather than as it stood. A card can decline a compiled rule and assert one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are also the card marking its own homework: a decline lowers that card's score, an assertion raises it, and both are only readable with the reason attached. Which is why the syntax insists on one. [MEASUREMENTS.md](MEASUREMENTS.md) has the rest.

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | the shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in the binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | the hook binary |
| `internal/scan/` | the structure rules, the card's regex rules, and the card renderer |
| `internal/effigy/` | the `.effigy` reader, so a card is usable as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane a turn was written in |
| `replay/` | the blind-pairs and discrimination harnesses, and their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what was run, on how much text, and what it said |

---

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
