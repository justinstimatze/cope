# cope

cope ships one opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — and once it is installed that card is on, checking two different jobs rather than one, what a reply sounds like and how it is shaped; the card is a file, so it can be edited or swapped later.

Read [demo/README.md](demo/README.md) first. Every file in that directory is this README written again from a different card, same prompt and the same facts, so the card is the only thing that changed between them, and holding two of them side by side shows what a card does faster than the rest of this page explains it — [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card that instructs every tic this model is measured to have, and it makes the point in a single glance. A reader who follows that link understands the tool. A reader who does not is taking prose about registers on trust.

## The two things

Ask what is wrong with a bad reply and the answer comes back about wording, because wording is the part you can quote. So let me tell you about the half that cannot be quoted. Voicing is what the sentences sound like — register, rhythm, diction, what a paragraph does with a detail — and it lives entirely in the card, so swapping the card swaps every word of it; that is the half with a measured result behind it. Structure is where the decision sits and how the reply ends, and it is compiled into the binary, so it is the same whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing problem. A reply that names an open problem in its last paragraph and then stops, leaving "continue" nothing to refer to, is a structural one, and no amount of rewriting the sentences will touch it. The same reply can be immaculate on one axis and useless on the other, which is the whole reason to keep them apart.

A card reaches into the structure half in two directions, and this is the part that most recently changed. It can decline a built-in rule, with `@gate <rule_id> off — <why>`, one per line in the card header; that exists because a card whose VOICE block asked for something a built-in rule catches was being marked down for obeying itself, which is what card/demo/lecturer.effigy does when it declines clause_symmetry and dangling_end. It can also assert a structural rule of its own, with `@shape <id>: <selector> <predicate> — <why>`, one per line in the same header; that exists because a card's own commitment about how a reply ends had nowhere to be checked, which is why card/demo/handoff.effigy asserts readable_cold — last paragraph words <= 60 — since its peak asks the reader to re-enter cold and read the last block only, and the 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90. The vocabulary a card writes in counts words and sentences and asks whether a block poses a question, and nothing more, so a check like clause_symmetry is not writable in it and is not meant to be. [MEASUREMENTS.md](MEASUREMENTS.md) has the runs and the reasons their numbers do not carry further than they do.

## The problem

Instruction alone does not fix the phrasing, and the cheapest way to learn that is to try it. A global CLAUDE.md banning the "not A, it's B" flip is read every single turn, sits in front of the model the whole time, and costs nothing to add. The flip appeared twice in the session that built this tool, while the ban was the topic of the session. Naming a surface form does not remove the move underneath it; it pushes the move into a variant that the wording of the ban does not reach.

That is the voicing side. The structural complaint is a different animal with a different cause: a reply that ends by naming an open problem and then stopping gives the reader nothing to answer, so the next turn is spent asking what to do, and a whole round trip is gone. No instruction banning a phrase could have caught that, because no phrase was involved.

The flip is an anecdote about one rule, and it is worth about as much as an anecdote. What the claim actually rests on is the blind discrimination test, where a reader is shown only a voice's own description of itself and picks which of two replies was written under it. The rate, and the caveats that go with it, are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

Two commands and one menu choice.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then pick it under `/config` -> Output style.

What `--setup` did, since a reader letting a tool edit their settings deserves the shape of it beforehand rather than after: it emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left. It backs the settings file up before writing, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses outright to touch a settings file that does not parse. `--setup --dry-run` prints what it would change and writes nothing at all.

For anyone who would rather not have their settings written to, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` on its own, and `COPE_CARD=<name>` in front of it emits a different one. The selection is made the same way either route: `/config` -> Output style, since the standalone `/output-style` command was removed in Claude Code v2.1.91, or as `"outputStyle"` in `.claude/settings.local.json`.

A style is read once at session start. A new selection, or a re-emitted card, therefore applies at the next session or after `/clear` — a reader who picks a style, watches the running conversation carry on unchanged, and concludes the tool is broken has simply caught it at the wrong moment.

The card goes into an output style rather than into a hook because an output style lands at the end of the system prompt and the harness re-reminds the model of it as the conversation runs.

The hooks are a separate matter now, and neither of them delivers the card:

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

Stop scores the reply just written, appending which rules fired to the session's rolling state and one record per violation to the log. UserPromptSubmit reads that rolling state and injects the card items gated on what has actually been firing, naming the counts, which is the one thing a file written once cannot do. The voice works without either of them. These two are the measurement half.

Both commands assume that `go install`'s target directory is on PATH in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. There is also `--inject`, the superseded delivery that puts the card in one turn-zero message, kept for anyone who wants it and standing down on its own whenever a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written — no Python, no effigy checkout, no build step between the file and the session. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session quietly writing in the shipped voice while its config names another one.

Read card/demo/lecturer.effigy first. It differs from the shipped card on register alone, which is what makes it the one the discrimination run measured. What a card can change is the voicing axis, in the sense The two things gave that word, and the structure rules stay where they are.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this page written again from a different card against the same prompt and the same facts, so the only thing varying between them is the voice.

One file in that directory is not a voice at all: card/demo/handoff.effigy is a hypothesis, keeping the shipped card's handoff rules and dropping everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it sitting there among the registers, and should know it is not one.

## What the hooks do differently

A pasted CLAUDE.md says the same thing on turn two and turn two hundred, which is its virtue and also the ceiling on it. The two hooks `--setup` wires exist to get above that ceiling. Stop scores the reply after it is written and records which rules fired into the session's rolling state. UserPromptSubmit reads that record — the state, not the violations log — and injects only the items gated on what has been firing, naming the counts, then stays quiet again until the last injection has aged past `--refresh-every`. When the session has no history yet, it falls back to the standing CONTINUE TEST.

So the mid-session text is chosen from measured output instead of being fixed in advance. That is a mechanism and not a guarantee, and it is worth saying that the A/B in this repo does not separate the refresher from no refresher. The old SessionStart `--inject` hook remains available and is off by default.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label, and three of its blocks turn out to do exactly what a prose gate needs. POSTPROC is regex rules with a warn action applied after generation. WRONG holds an anti-pattern beside its replacement, so a rule shows the fix and not only the fault. TEST holds a named question with fail and pass examples, which is how a rule names a move rather than one wording of it — and one wording is precisely what the CLAUDE.md ban above turned out to be.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and it is the one to reach for when this one is too blunt: cope bans, so a rule fires or it does not and the register is fixed the moment you pick a card, while basanite measures lemma frequency against a baseline over real transcripts and reports what you have been leaning on lately, leaving the judgement to you. Its own README calls that awareness rather than prohibition. Which one fits is a question about mood more than correctness, a heavy hand suits a habit that is annoying you today and a moving measurement suits watching the drift instead of legislating it, and since they share no state and sit on different hooks, running both is reasonable. [caveman](https://github.com/JuliusBrussee/caveman), by a different author again, compresses agent replies to cut output tokens — a third axis, and the right destination for a reader who wants fewer tokens rather than different structure.

## The rules

Grouped by the two jobs rather than by where they happen to be implemented.

Voicing:

- **clause_symmetry** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- **apology** — the reply performs contrition instead of stating the correction and moving on.
- **self_postmortem** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **announced_length** — the reply announces its own length rather than cutting it.
- **cross_turn_repeat** — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- **labelled_opening** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **paragraph_uniformity** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **ask_not_last** — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **dangling_end** — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **buried_decision** — an open problem landing after the last question or offer, burying the decision point above it.
- **forked_end** — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **unverified_done** — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **loop_ask** — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The grouping gives away something about the implementation that is worth having. A POSTPROC pattern matches a span of text, and a span of text is wording, so a pattern can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go, sitting beside the structure rules in `internal/scan`, and this card ships zero POSTPROC rules of its own. If you came here expecting a long list of banned phrases, the list lives in basanite on purpose, where it is a measurement rather than a prohibition.

There is one place the structure rules do vary, and it is not by card. It is by who is going to read the turn. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself, and it drops ask_not_last, forked_end, dangling_end and buried_decision, since nobody is reading yet: a report that correctly names what it left open and stops would fail three of those, and a question in it lands in a log where the next iteration reads it as an instruction to itself. In their place go unverified_done and loop_ask. A report claiming the work is done has to say what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | (empty) | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report the arms; - reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also $COPE_CARD |
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

| path | mode | holds |
| --- | --- | --- |
| `$XDG_STATE_HOME/cope/violations.jsonl` | 0600 | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | 0600 | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | 0600 | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

The log quotes your replies back at you. That is what the matched text and its surrounding characters are, and `--log` with an empty value turns the file off.

## Editing the card

The `.effigy` grammar belongs to effigy, `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs so that the rules being enforced and the rules being injected cannot drift apart.

The NEVER budget is 10, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately rather than against the card file, since SessionStart prints the always-on rules and the refresher prints the evidence-gated ones and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be perfectly healthy. The authoritative list of rules that really are discarded unrendered is empty when things are well.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is exactly this and no more:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

For `@gate` the rule id has to be one the gate already has; for `@shape` it must not collide with one. A wrong id in either is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports everything it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary could supply.

## Calibrating

`cope-gate --backfill` scores every assistant turn in a transcript at once, and it is how these rules were chosen in the first place. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is worth reading carefully, because the largest N is not the same set as one per project and a sweep can therefore be dominated by a single talkative repo.

Watch hits per character rather than the share of turns hit. The second number moves when the turns get longer, which is a fact about the turns and not about the prose in them. The rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

The axis split organises these, because a limit on one half says nothing about the other.

labelled_opening is not a tagger, and it is guessing at what counts as a verbless fragment. ask_not_last notices that an ask is not last and has nothing whatever to say about the order of several asks. The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written, which makes it a description of the output rather than a judgement of it — the judgement lives in the discrimination test, and that test covers voicing only.

The largest limit is the one The two things already named, and it stands differently now than it did. A card can decline a built-in rule and can write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. That leaves the compiled rules as the only place a check like clause_symmetry can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere at all to put it. Both directions are also the card marking its own homework: a decline lowers that card's score, an assertion raises it, and either is worth reading with the reason attached, which is why the syntax will not accept one without. [MEASUREMENTS.md](MEASUREMENTS.md) has the rest.

## Layout

| path | what |
| --- | --- |
| `card/claude_voice.effigy` | the shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in the binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | the hook binary |
| `internal/scan/` | the structure rules, the card's regex rules, and the card renderer |
| `internal/effigy/` | the .effigy reader, so a card is usable as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane a turn was written in |
| `replay/` | the blind-pairs and discrimination harnesses, and their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what was run, on how much text, and what it said |

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
