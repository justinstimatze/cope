# cope

cope ships one opinionated card, on from the moment you install it, and it watches two different jobs in every reply the model writes — how the sentences sound, and how the reply is shaped for whoever has to act on it — and because that card is a file, the sounding half is yours to edit or swap; a cope is the upper half of a foundry mould, the half carrying the shape being cast into.

Start here: [demo/README.md](demo/README.md).

Every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed between them. Reading two of them against each other shows what a card does faster than the rest of this page explains it, and a reader who does not follow the link is taking prose about registers on trust. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes the point in a single glance, written from a card that instructs every tic this model is measured to have, and it is deliberately hard going.

## The two things it checks

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. All of it lives in the card, which is why swapping the card swaps every word of it, and it is the half with a measured result behind it. Structure is a different job — where the decision sits, whether the ending gives "continue" anything to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it — and that half is compiled into the binary, identical whichever card is loaded. A sentence reaching for the balanced two-beat, two clauses of near-equal length repeating a content word across the joint, is a voicing problem. A reply that names an open problem in its last paragraph and then stops is a structural one, and it can be beautifully written on the way there. The same reply can be fine on one axis and bad on the other, which is the whole reason to keep them apart.

A card reaches into the structure half in two directions, and both of them are narrow. It can decline a rule, with `@gate <rule_id> off — <why>` one per line in the card header, which exists because a card whose VOICE block asked for a landing that a built-in rule catches was being marked down for obeying itself. It can also assert one of its own, with `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply ends had nowhere to be checked. What a card writes those assertions in counts words and sentences and asks whether a block poses a question, and nothing more, so a check outside that vocabulary and outside a POSTPROC regex still has nowhere to live.

## The problem

You have probably edited a global CLAUDE.md, told it something plain about how you wanted replies written, and watched the model comply for three turns and then quietly stop. The instinctive reading of that is that the model is careless. The likelier one is that you were writing into the wrong slot. A global CLAUDE.md is not the system prompt: it arrives as one message attached to the first turn, and every turn written after it pushes it further down the transcript, until it is one old note among forty newer ones. An output style sits in the system prompt itself, where the harness re-reminds the model of it as the conversation runs. Take one card, change no word of it, and move it from the first place to the second. That move is most of why cope works, and it was measured — [MEASUREMENTS.md](MEASUREMENTS.md) has the run.

Delivery is only the first part, because instruction alone does not fix the phrasing either. A global CLAUDE.md banning the "not A, it's B" flip is read on every turn of every session. That flip appeared twice in the session that built this repo, while the ban was the topic of the session. Name a surface form and the move steps sideways into a variant of itself, which is why the card carries anti-patterns beside their replacements and named tests rather than a list of forbidden strings.

The structural complaint has a different cause and a different cost. A reply that ends by naming an open problem and offering nothing to answer buys you a whole round trip you did not need, and no instruction could have banned it, because it is not a habit of phrasing at all. It is a shape.

The flip, though, is an anecdote about one rule, and it is not what the claim rests on. What it rests on is the blind discrimination test: a reader is shown a voice's own description of itself and two replies, one written under that voice and one not, and asked to pick. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate and the caveats that come attached to it.

## Install

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

One menu choice is left, and `--setup` prints it: the style is selected under `/config` → Output style. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way in; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

What `--setup` did is worth knowing before you run it rather than after, since it edits a file you care about. It emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the remaining step. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses outright to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing.

If you would rather nothing wrote to your settings at all, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` and stops there, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start. So a new selection, or a re-emitted card, applies at the next session or after `/clear`, and a reader who picks the style and watches the running conversation for a change will see none and conclude the tool is broken. The reason the card goes into a style rather than a hook is the one the previous section gave: the style lands at the end of the system prompt, where the harness keeps re-reminding the model of it.

The hooks are the other half, and the card no longer arrives through any of them:

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

Stop scores the reply just written. UserPromptSubmit restates, mid-session, the rules that have actually been firing — something a file written once at session start cannot do. The voice works without either of them; these two are the measurement half.

Both commands assume that `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled into the binary. There is also `--inject`, the superseded SessionStart delivery, kept for anyone who wants it and standing down on its own whenever a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, which means a card is usable exactly as written, with no Python and no effigy checkout anywhere in the loop. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, which is deliberate — the alternative is a session writing happily in the shipped voice while its config names another one, and that failure is invisible for as long as you keep reading.

The card to read first is `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and it is the card the discrimination run measured. What a card can change is the voicing axis, in the sense given under "The two things it checks", and the structure rules stay where they are compiled.

Shorter than reading any card is reading two pages written from different ones: [demo/README.md](demo/README.md) is this README rendered again under each demo card, from the same prompt and the same facts.

One file there is not a voice at all. `card/demo/handoff.effigy` is a hypothesis: it keeps the shipped card's handoff rules, drops everything about prose, and is meant to be run through `make pairs` against the full card rather than rendered — `make cards` installs it with the rest, so anyone listing their cards will find it sitting among them looking like a register to write in.

## What the hooks add

Stop runs after a reply is written. It scores that reply, appends which rules fired to the session's rolling state, and writes one record per violation to the log.

UserPromptSubmit is the interesting one, and the reason is in what it reads. It reads the rolling state rather than the violations log, and injects the card items gated on what has been firing lately, naming the counts. So the mid-session text is chosen from your measured output instead of being fixed in advance, which is the one thing a pasted CLAUDE.md structurally cannot do. It falls back to the card's standing reminder when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. Claiming more for it than that would be dishonest: it is one mechanism, and the A/B in this repo does not separate the refresher from no refresher at all. The old SessionStart delivery still exists, superseded and off by default, and `--setup` does not wire it.

## Why effigy notation

A prose gate needs three things written down that most config formats have no slot for, and [effigy](https://github.com/justinstimatze/effigy) — a character-card notation for game NPCs, used here entirely off-label — happens to have all three. POSTPROC is regex rules with a warn action applied after generation. WRONG holds an anti-pattern beside its replacement, so a card can show the same content written worse and better instead of describing the difference. TEST holds a named question with a failing and a passing example, which is how a rule names a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and it is the one to reach for when this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead, tracking lemma frequency against a baseline over real transcripts, so it reports what you have been leaning on lately and leaves the judgement with you — its own README calls that awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose, on different hooks with no shared state, and running both is reasonable. A reader who wants fewer output tokens rather than a different structure wants neither, and should look at [caveman](https://github.com/JuliusBrussee/caveman), which is a separate project by a different author that compresses agent replies.

## The rules

Grouped by the two jobs rather than by where they happen to be implemented, because the grouping is what tells you which ones a card can touch.

Voicing, from the shipped card's POSTPROC patterns:

- **flip** (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- **load_bearing** (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- **worth_noting** (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, where a pattern was not enough:

- **clause_symmetry** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- **apology** — the reply performs contrition instead of stating the correction and moving on.
- **self_postmortem** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **announced_length** — the reply announces its own length rather than cutting it.
- **cross_turn_repeat** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, all of it compiled in:

- **labelled_opening** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks, with an ordinal counting as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **paragraph_uniformity** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **ask_not_last** — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **dangling_end** — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **buried_decision** — an open problem landing after the last question or offer, burying the decision point above it.
- **forked_end** — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **unverified_done** — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **loop_ask** — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

That split says something usable about the implementation. A POSTPROC pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a span had to be written in Go beside the structure rules. Which is why the shipped card carries three POSTPROC rules and no more. If you came here expecting a long list of banned phrases, that list lives in another tool on purpose, and the previous section names it.

The structure rules do vary in exactly one way, and it has nothing to do with which card is loaded. It is about who is going to read the turn. The interactive lane runs for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane runs when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. There, `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision` are dropped and `unverified_done` and `loop_ask` take their place, because nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check. A report saying the work is done has to say what it ran.

## Flags

| Flag | Default | What it does |
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

Three files, all `0600`:

- `$XDG_STATE_HOME/cope/violations.jsonl` — one JSON record per violation, carrying the matched text and about 70 characters either side.
- `$XDG_STATE_HOME/cope/refresher-<session-id>` — an empty file whose mtime is the refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

The middle two hold nothing you would mind. The first one quotes your replies back at you.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so the rules being enforced and the rules being injected cannot drift apart quietly.

There is a NEVER budget of 10, and anything over it is reported at load rather than dropped silently — a rule that vanished without a word would be worse than one that never existed. The budget is charged against each injection separately and not against the card file, because the SessionStart text prints the always-on rules, the refresher prints the evidence-gated ones, and no code path anywhere renders their union. A card can therefore hold more NEVER rules in total than the budget and still be perfectly healthy.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is exactly this and no more:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

For `@gate` the id has to be one the gate already has; for `@shape` it must not collide with one, and a wrong id in either is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

A declined rule is not switched off, either. It still runs, and only this card's score drops it, so a backfill still reports everything it would have caught. And a `@shape` violation comes back in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and it is how these rules were chosen in the first place. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same thing as one per project, and reading it as one per project will flatter or damn whichever project happens to have the long sessions.

The number worth watching is hits per character. The share of turns hit looks like the more honest metric and is not, because it mostly tracks how long the turns were. [MEASUREMENTS.md](MEASUREMENTS.md) has the rates.

## Known limits

The axis split organises the limits as neatly as it organises the rules.

On the structural side, `labelled_opening` is not a tagger: it counts a short verbless opener and an ordinal, and cannot tell a label from a fragment that had every right to be there. `ask_not_last` says nothing whatever about the order of several asks, only that one of them is stranded above the rest of the reply. And the hit rate is roughly four fifths structure, which sounds like a verdict on how the model writes until you look at the A/B run: that four fifths tracks what a reply was for. It is a description of the output, not a judgement of it. The judgement lives in the discrimination test, and that test covers voicing only.

The largest limit is the one "The two things it checks" already named, and it has moved recently enough to be worth stating as it now stands rather than as it stood. A card can decline a built-in rule and can write one of its own. The vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

There is a second thing to hold about both directions, and it is not a defect so much as a condition of the design. A card that declines a rule lowers its own score, and a card that asserts one raises it. Both are the card marking its own homework, both are only readable with the reason attached, and that is precisely why the syntax will not accept one without it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rest.

## Layout

- `card/claude_voice.effigy` — the shipped card, in effigy notation
- `card/rules.json` — generated from it; embedded in the binary
- `card/demo/` — other cards, each written to sound like something else
- `cmd/cope-gate/` — the hook binary
- `internal/scan/` — the structure rules, the card's regex rules, and the card renderer
- `internal/effigy/` — the `.effigy` reader, so a card is usable as written
- `internal/transcript/` — Claude Code JSONL reader, and which lane a turn was written in
- `replay/` — the blind-pairs and discrimination harnesses, and their own README
- `demo/` — this README written again under each demo card
- `tools/` — card compiler, effigy-backed scorer, cross-project sweep
- `MEASUREMENTS.md` — what was run, on how much text, and what it said

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
