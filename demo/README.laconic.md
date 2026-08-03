# cope

A cope is the upper half of a foundry mould, the half carrying the shape being cast into: here it is a file you write, swapping that file swaps every word of the register, and a second set of checks compiled into the binary watches how a reply is shaped rather than how it sounds.

Read [demo/README.md](demo/README.md) first. Every file there is this README written again from a different card, same prompt and same facts, so the card is the only thing that changed, and two of them read against each other show what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) makes the point in one glance, from a card that instructs every tic this model is measured to have.

## The two things

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. It lives in the card — VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC — and swapping the card swaps all of it. That is the half with a measured result behind it.

Structure is the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. It is compiled into the binary, the same whichever card is loaded.

A sentence reaching for the balanced two-beat is a voicing problem. A reply that names an open problem in its last paragraph and then stops, leaving "continue" nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

A card reaches into the structure half in two directions. It declines a built-in rule:

```
@gate <rule_id> off — <why>
```

one per line in the card header. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch. Without the decline, a card is marked down for obeying itself.

Or it states a structural rule of its own, which the gate then checks:

```
@shape <id>: <selector> <predicate> — <why>
```

also one per line in the header. `card/demo/handoff.effigy` asserts `readable_cold` — last paragraph words <= 60 — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more. It cannot express what the compiled rules express; `clause_symmetry` is not writable in it and is not meant to be. A card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

MEASUREMENTS.md has the run and the reasons its numbers do not carry more than that.

## The problem

Instruction alone does not fix the phrasing. A global CLAUDE.md banning the "not A, it's B" flip is read every turn; the flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

The structural complaint has a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no instruction against a phrase could have banned it.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. MEASUREMENTS.md has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick it under `/config` → Output style.

`--setup` emitted the output style and wired the two hooks into `~/.claude/settings.json` with absolute paths, then printed the one step left. It backed the settings file up first, added only what was missing so a second run changes nothing, left every other key alone including other tools' hooks on the same events, and refused to touch a settings file that does not parse; `--setup --dry-run` prints what it would change and writes nothing.

To keep your settings untouched, do it by hand. `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

A style is read once at session start. A new selection or a re-emitted card applies at the next session, or after `/clear`.

An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands there and did not through a hook.

The hooks are the measurement half. The voice works without either of them.

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

Stop scores the reply just written. UserPromptSubmit restates mid-session the rules that have actually been firing, which a file written once cannot do.

The commands assume `go install`'s target directory is on PATH in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead.

A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in.

The old SessionStart delivery survives as `--inject`, superseded and off by default, standing down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written: no Python, no effigy checkout. Drop one in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, and it is what the discrimination run measured. What a card changes is the voicing axis, as **The two things** sets out.

[demo/README.md](demo/README.md) is the shortest way to see what a card does without writing one: every file under `demo/` is this README written again from a different card, same prompt and same facts, so the voice is the only thing that varies between them.

`card/demo/handoff.effigy` is the exception in that directory — a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## What the hooks do differently

Stop scores the reply and appends which rules fired to the session's rolling state, plus one record per violation to the log.

UserPromptSubmit reads that state — not the violations log — and injects the card items gated on what has been firing, naming the counts. The mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted CLAUDE.md cannot do. It falls back to the standing CONTINUE TEST when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. The shipped card gates no items on evidence today, so the fallback is what it prints.

One mechanism, not a guarantee. The A/B in the repo does not separate the refresher from no refresher.

The SessionStart injection is superseded and off by default, and `--setup` does not wire it.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) is the same problem answered the other way round, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead — lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand suits a habit that is annoying you today, a moving measurement suits watching the drift instead of legislating it. They compose, on different hooks with no shared state, and running both is reasonable.

[caveman](https://github.com/JuliusBrussee/caveman) is a third axis again, by a different author: it compresses agent replies to cut output tokens. cope shapes prose, basanite tracks vocabulary, caveman shortens. A reader wanting fewer tokens rather than different structure should go there instead.

## The rules

Voicing:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

A POSTPROC pattern matches a span of text, so it can only ever describe wording. Every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries 0 POSTPROC rules.

A reader expecting a long list of banned phrases wants basanite. The list lives in another tool on purpose.

The structure rules vary in one way, and not by card. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. It drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask`: nobody is reading yet, a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| flag | default | does |
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
| `--display-preview` | `false` | read prose on stdin and print it as --display would rewrite it |
| `--dry-run` | `false` | with --setup, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to ~/.claude/output-styles as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default ~/.claude/output-styles) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render -render-arm against |
| `--render-lane` | (empty) | render -render-arm as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this .effigy or .json file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

- `$XDG_STATE_HOME/cope/violations.jsonl`, mode 0600 — one JSON record per violation, carrying the matched text and about 70 characters either side. The log quotes your replies back.
- `$XDG_STATE_HOME/cope/refresher-<session-id>`, mode 0600 — an empty file whose mtime is the refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json`, mode 0600 — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json`; `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The NEVER budget is 10. Anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file — SessionStart prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more NEVER rules in total than the budget and still be healthy.

A card writes its own structural rules in two forms:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

For `@gate` the id has to be one the gate has. For `@shape` it must not collide with one. A wrong id is reported at load rather than ignored.

The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

A declined rule still runs, and only this card's score drops it, so a backfill still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Watch hits per character rather than the share of turns hit. That second number tracks how long the turns were. MEASUREMENTS.md has the rates.

## Known limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written. It describes the output rather than judging it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one **The two things** named, as it now stands. A card can decline a built-in rule and write one of its own, in a vocabulary that counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex has nowhere to put it.

Both directions are the card marking its own homework. A decline lowers that card's score and an assertion raises it. Read both with the reason attached, which is why the syntax requires one.

MEASUREMENTS.md.

## Layout

- `card/claude_voice.effigy` — the shipped card, in effigy notation
- `card/rules
