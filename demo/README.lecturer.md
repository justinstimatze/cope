# cope

A cope is the upper half of a foundry mould, the half carrying the shape being cast into, and this one ships an opinionated card that is on the moment you install it, watching two different things about a reply — what its sentences sound like, and how the whole of it is shaped — with that card sitting on disk as a file, so it can be edited or swapped once you want a different voice.

Before any of the rest: [`demo/README.md`](demo/README.md) is a directory in which every file is this same page written again from a different card, the same prompt and the same facts each time, so that the card is the only thing that changed between them — and reading two of them side by side shows you what a card does faster than the rest of this page can explain it, with [`demo/README.claude-maximal.md`](demo/README.claude-maximal.md), written from a card that instructs every tic this model has been measured to have and heavy going entirely on purpose, making the point in a single glance.

## The two things it is looking at

Ask what went wrong with a piece of writing and you will nearly always be told about the wording, because the wording is the part you can quote back at somebody. So let me tell you about the other part. Voicing is what the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, where flair is licensed — and it lives entirely in the card, in the `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES` and `POSTPROC` blocks, so swapping the card swaps every word of it; that is the half with a measured result standing behind it. Structure is the shape of the reply as a thing you have to use: where the decision sits, whether the ending gives "continue" something to refer to, whether a claim that the work is done carries anything that could have shown it. Structure is compiled into the binary, identical whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing fault; a reply that names an open problem in its last paragraph and then simply stops is a structural one. The same reply can be clean on one axis and broken on the other, which is the entire reason for keeping them apart.

A card reaches across into the structural half through two doors, and only two. It can decline a built-in rule it disagrees with — `@gate <rule_id> off — <why>`, one per line in the card header — which exists because a card whose `VOICE` block asked for precisely the thing a built-in rule catches was being marked down for obeying itself, and `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` on exactly that ground. It can also assert a structural rule of its own — `@shape <id>: <selector> <predicate> — <why>` — which exists because a card's own commitment about how a reply ends had nowhere to be checked; `card/demo/handoff.effigy` asserts `readable_cold`, a cap on the words in its final block, because its peak asks the reader to come back cold and read that block alone. Past those two doors the card stops. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, so a check like `clause_symmetry` is not writable in it, and a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere at all to put it.

## Where an instruction sits, and what it can say

You have edited a global `CLAUDE.md`. You have watched it fail to stick somewhere around turn thirty, and you have quietly concluded that the model is bad at following instructions. Here is the piece almost everyone arrives holding backwards: that file is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written afterwards, the way anything said at the beginning of a long afternoon gets buried. An output style is in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Take one card, change not a single word of it, and move it from the first place to the second, and most of what cope does has already happened. That was measured, and [`MEASUREMENTS.md`](MEASUREMENTS.md) carries the run.

Seating the instruction correctly is only the first of the three problems, though, because instruction alone does not fix the phrasing either. A global `CLAUDE.md` forbidding the reflex that denies one thing in order to assert another is read on every single turn, and that reflex appeared twice in the session that built this tool, while the ban itself was the subject under discussion. Name a surface form and the move relocates into a variant of itself. This is why the card carries examples, anti-patterns beside their replacements, and named tests, instead of a list of strings.

The third failure is not a phrasing problem at all, which is why no amount of better instruction was ever going to reach it. A reply that ends by naming something unresolved and then gives you nothing to answer costs a full round trip: you come back to the terminal, read the last paragraph, and have to compose a question that the reply should have asked on your behalf. No banned phrase covers that. It is a fact about where the decision was put.

That reflex is an anecdote about one rule. What the claim actually rests on is the blind discrimination test, where a reader is shown a voice's own description of itself and two replies, and asked which of the two was written under it. The rate, and the reasons the rate does not carry more weight than it does, are in [`MEASUREMENTS.md`](MEASUREMENTS.md).

## Installing it

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick the style under `/config` -> Output style.

Those two commands do the whole install between them. `--setup` emits the output style and wires the hooks into `~/.claude/settings.json` with absolute paths, and it is careful about it: it backs the settings file up before writing, it adds only what is missing so that a second run changes nothing at all, it leaves every other key alone including other tools' hooks on the same events, and it refuses outright to touch a settings file that does not parse. If you would rather look before it writes, `cope-gate --setup --dry-run` prints what it would change and writes nothing. If you would rather it never wrote to your settings, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` and leaves the hooks to you, and `COPE_CARD=<name>` in front of it emits a different card instead.

One thing to know before you go looking for the effect. A style is read once at session start, so a new selection, or a card you have just re-emitted, applies at the next session or after `/clear` — a reader who picks the style, watches the current conversation carry on exactly as before, and concludes the tool is broken has simply been fast. The standalone `/output-style` command was removed in Claude Code v2.1.91, which is why the path above goes through `/config`; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. The reason the card goes here rather than into a hook is the one already given: an output style lands at the end of the system prompt, and the harness keeps reminding the model of it. There is still a `--inject` flag, the superseded delivery that puts the card in a turn-zero message, kept for anyone who wants it and standing down on its own whenever a cope output style is active.

The hooks are a separate matter, and the card no longer arrives through either of them. This is the `hooks` block of `~/.claude/settings.json`:

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

What `Stop` buys you is that the reply gets scored after it has been written. What `UserPromptSubmit` buys you is that the rules which have actually been firing get restated mid-conversation, which a file written once and read once cannot do. The voice works with neither of them installed. These two are the measurement half.

Both commands assume that `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing is usually asking for the absolute path instead. From a clone, `make install` builds it, and no effigy checkout is needed for that: `card/rules.json` is committed and compiled straight into the binary.

## Writing in another voice

The gate reads `.effigy` files directly, which means a card is usable exactly as it was written, with no Python anywhere and no effigy checkout to keep in sync. Drop a card into `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or with `COPE_CARD`; point `--rules` at a path and it will read a card from anywhere on disk; run `make cards` and the demo set installs itself. A name that resolves to nothing is treated as an error and nothing is injected at all, rather than the far worse outcome, which is a session quietly writing in the shipped voice while its configuration names a different one.

The card to read first is `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and is the one the discrimination run actually measured. What a card changes is the voicing axis, in the sense set out where the two axes were separated above; the numbers for how well that change carries are in [`MEASUREMENTS.md`](MEASUREMENTS.md) and not here.

If you would rather see what a card does than write one, [`demo/README.md`](demo/README.md) is the short way: every file under it is this README written again from a different card, same prompt and same facts, so the voice is the only variable in the set.

One file in that directory is not a voice at all. `card/demo/handoff.effigy` is a hypothesis — it keeps the shipped card's rules about how a handoff ends and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered into a page. `make cards` installs it alongside the rest, so anyone listing their installed cards will find it there and should know it is not a register to write in.

## What the hooks add that a file cannot

Here is the thing a card on disk can never do, however well it is written and however well it is seated: it cannot know what you have been doing wrong this afternoon. It is the same text at turn two and at turn two hundred.

The `Stop` hook scores the reply that was just written, appends which rules fired to that session's rolling state, and writes one record per violation to the log. The `UserPromptSubmit` hook then reads that rolling state — the state file, not the violations log — and injects the card items gated on what has been firing, naming the counts as it does so. When a session has no history to read yet, it falls back to a standing reminder, and it stays quiet entirely until the last injection has aged past `--refresh-every`. So the mid-session text is chosen from measured output rather than fixed in advance. That is one mechanism and not a guarantee, and the A/B harness in this repository does not separate having the refresher from not having it.

## Why a character-card notation

Reaching for a game-NPC format to hold a prose specification looks like a joke until you look at what the format already has. Three of [effigy](https://github.com/justinstimatze/effigy)'s blocks do exactly what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern directly beside its replacement, and `TEST` holds a named question with a failing example and a passing one. That last block is the important one, because it is how a rule can name a move rather than one wording of the move.

If cope reads as too blunt an instrument — and it is blunt, since a rule either fires or it does not, and the register is fixed the moment you pick a card — then [basanite](https://github.com/justinstimatze/basanite) answers the same problem from the other end, measuring lemma frequency against a baseline across real transcripts and reporting what you have been leaning on lately without legislating anything. Its own README calls that awareness rather than prohibition. They compose, on different hooks with no shared state, and running both is perfectly reasonable. If what you want is neither of those but simply fewer output tokens, [caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies and is the right place to go instead.

## The rules

Grouped by which of the two jobs they are doing, rather than by where they happen to be implemented.

Voicing, from the shipped card's `POSTPROC` block:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, where a pattern was not enough and the rule had to be written in Go:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, which is the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. It is the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, all compiled in:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it then unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as things they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so that answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The grouping tells you something usable about the implementation. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule needing more than a pattern had to be written in Go next to the structure rules. That is why the shipped card carries three `POSTPROC` rules and no long blacklist. If you came expecting a list of banned phrases, the list lives in another tool on purpose, and that tool is basanite.

There is one place the structure rules vary, and it is not by card. It is by who is going to read the turn. The interactive lane is chosen for any turn that is not a loop turn, on the grounds that somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself, and there it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, adding `unverified_done` and `loop_ask` in their place. Nobody is reading yet, so a report that correctly names what it left open and then stops would fail three of the dropped rules for no benefit, and a question inside it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is a claim check: a report saying the work is done has to say what it ran.

## The flags

| Flag | Default | What it does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | `(empty)` | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
| `--ab-report` | `(empty)` | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | `(empty)` | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | `(empty)` | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | `(empty)` | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | `MessageDisplay` entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | `(empty)` | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | `(empty)` | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | `(empty)` | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | `(empty)` | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

Three files, all written `0600`, all under `$XDG_STATE_HOME/cope`. `violations.jsonl` holds one JSON record per violation, carrying the matched text and about seventy characters either side of it, which is to say plainly that the log quotes your replies back to you. `refresher-<session-id>` is an empty file whose mtime is the refresher clock and nothing else. `session-<session-id>.json` is the rolling record the mid-session injection gets chosen from: turn count, characters, and which rules fired over a twenty-turn window, with no prose in it at all — only rule names and counts.

## Editing the card

The `.effigy` grammar belongs to effigy, not to cope. On this side, `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so that the rules being enforced and the rules being injected cannot quietly drift apart.

There is a NEVER budget of 10, and anything over it is reported at load rather than dropped silently, which matters more than it sounds like it does. The budget is charged against each injection separately and not against the card file: the always-on rules and the evidence-gated ones are printed by different paths and no code renders their union, so a card may perfectly well hold more NEVER rules in total than the budget and still be entirely healthy. The list of rules that really are being discarded unrendered is its own field, and on a healthy card it is empty.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is exactly this and no more. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. For `@gate` the id has to be one the gate already has; for `@shape` it must not collide with one, and either way a wrong id is reported at load rather than ignored. The reason after the dash is required in both forms, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs — only this card's score drops it — so a backfill will still tell you what it would have caught, and a `@shape` violation is reported back in the card's own words rather than in any sentence the binary could have supplied.

One number in there is worth knowing where it came from. `card/demo/handoff.effigy` caps its last paragraph at 60 words, and that 60 was measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

## Calibrating it against your own transcripts

You do not have to take the rules on faith, and the rules themselves were not chosen on faith either. `cope-gate --backfill` scores a whole session transcript at once, which is how they were picked in the first place, and `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects` — the largest N, note, and not one per project.

When you read the output, watch hits per character rather than the share of turns hit. The share of turns hit mostly tracks how long the turns were, which is a fact about your prompts. The rates themselves are in [`MEASUREMENTS.md`](MEASUREMENTS.md).

## Known limits

`labelled_opening` is not a tagger, and it will occasionally take an ordinary short opening for a label. `ask_not_last` says nothing whatsoever about the order of several asks; it only knows that one is not last.

Roughly four fifths of all hits are structural, and the A/B run found that this four fifths tracks what a reply was *for* rather than how well it was written. Take the hit rate as a description of the output and not as a judgement of it. The judgement lives in the discrimination test, and the discrimination test covers voicing only.

The largest limit is the one already named where the two axes were separated. A card can decline a built-in rule and can assert a structural rule of its own, and the vocabulary it writes that rule in counts words, counts sentences, and asks whether a block poses a question. The compiled rules therefore remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere at all to put it. Both directions are also the card marking its own homework: a decline lowers that card's score, an assertion raises it, and both need reading with their reason attached. That is precisely why the syntax refuses to accept one without a reason.

## Layout

| Path | What is in it |
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

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
