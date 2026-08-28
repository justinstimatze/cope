# cope

An opinionated card that ships on, checking two separate things about a reply — how its sentences sound and how the reply is shaped for the reader — and held in a file you can edit or swap, the name taken from the foundry: a cope is the upper half of a mould, the half carrying the shape being cast into.

Every file in `demo/README.md` is this README written again from a different card, same prompt and same facts, so the card is the only thing that changed; reading two of them against each other shows what a card does faster than the rest of this page explains it. `demo/README.claude-maximal.md` makes the point in a single glance, written from a card that instructs every tic this model is measured to have, and it is deliberately hard going.

## The two things

Voicing. What the sentences sound like — register, rhythm, diction, where flair is licensed. It lives in the card, so swapping the card swaps every word of it, and it is the half with a measured result behind it. A sentence reaching for the balanced two-beat, two clauses of matched length joined at a comma with a content word repeated across the joint, is a voicing problem.

Structure. Where the decision sits, whether an ask is last, whether the ending gives "continue" something to refer to, whether a claim that the work is done carries anything that could have shown it. Compiled into the binary, and the same rules whichever card is loaded. A reply naming an open problem in its last paragraph and then stopping is a structural problem. The same reply can be clean on one axis and bad on the other, which is why they stay apart.

Two ways a card reaches into the structure half. `@gate <rule_id> off — <why>`, one per line in the card header, declines a built-in rule: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` because its `VOICE` block asks for the balanced landing and the arriving close those two rules catch, and a card marked down for obeying itself needed somewhere to say so. `@shape <id>: <selector> <predicate> — <why>` states a rule of the card's own: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way.

The boundary. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more. `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex has nowhere to put it. `MEASUREMENTS.md` has the run and the reasons its numbers do not carry more than that.

## The problem

The card's address. A global `CLAUDE.md` is not the system prompt, though most readers who have watched one fail to stick believe it is. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style is in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works — measured, and `MEASUREMENTS.md` has the run.

Instruction alone. A global `CLAUDE.md` banning the "not A, it's B" flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic. Naming a surface form pushes the move into a variant.

The failure no instruction reaches. Structure is a different complaint with a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and it is not a phrasing habit an instruction could have banned.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it, and `MEASUREMENTS.md` has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`go install` puts the binary on disk. `cope-gate --setup` emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left: it backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `cope-gate --setup --dry-run` prints what it would change and writes nothing.

The menu choice. Under `/config` -> Output style, the entry to look for is named `claude_voice`, which is the shipped card's id and not the word cope, so a reader scanning that menu for something called cope will not find it. The standalone `/output-style` command was removed in Claude Code v2.1.91, and the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

By hand, for anyone who would rather not have their settings written to: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

The moment a style takes effect. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands there.

The hooks, in the `hooks` block of `~/.claude/settings.json`:

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

`Stop` buys a score on the reply just written. `UserPromptSubmit` buys a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. `PreToolUse` buys the same score on prose an external write is about to post. The voice works without any of them; these are the measurement half.

The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. `--inject` remains as the superseded delivery for anyone who wants it, and stands down on its own when a cope output style is active.

## Writing in another voice

Cards as files. The gate reads `.effigy` directly, so a card is usable as written and needs no Python and no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

`card/demo/lecturer.effigy` is the one to read first, differing from the shipped card on register alone, which is what the discrimination run measured. What a card can change is the voicing axis, described above where the two axes are separated, and `MEASUREMENTS.md` has the numbers.

The shortest way to see what a card does without writing one is `demo/README.md`, where every file is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, and meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## What the hooks do differently

`Stop`, running `cope-gate`. Scores the reply just written. Which rules fired lands in the session's rolling state; each violation adds a line to the log.

`UserPromptSubmit`, running `cope-gate --refresher`. Reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. Falls back to the standing `CONTINUE TEST` when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. One mechanism, not a guarantee: the A/B run in this repo does not separate the refresher from no refresher.

`PreToolUse`, running `cope-gate --pretool`. Scores the `description`, `body` or `content` field an external write is about to post, matched against the Linear save tools named in the settings block. Warn-only — it returns `additionalContext` and never a `permissionDecision`, so the call goes through and the model learns what the prose scored. Writes no session state, and scores in the external lane.

## Why effigy notation

`effigy` is a character-card notation for game NPCs, used here off-label, at https://github.com/justinstimatze/effigy. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

`basanite`, at https://github.com/justinstimatze/basanite, is the same problem answered the other way round and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures lemma frequency against a baseline over real transcripts, reporting what you have been leaning on lately and leaving the judgement to you. They compose — different hooks, no shared state — and running both is reasonable.

Compare `caveman`, at https://github.com/JuliusBrussee/caveman, a separate project by a different author that compresses agent replies to cut output tokens. cope shapes prose, basanite tracks vocabulary, caveman shortens; a reader wanting fewer tokens rather than different structure should go there instead.

## The rules

Voicing rules, which describe wording.

- `flip` — warn. The not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — warn. Reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` — warn. Announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure rules, which describe the shape of the reply as a thing the reader has to use.

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — interactive lane. A question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — interactive lane. An open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` — interactive lane. An open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — interactive lane. Two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — loop lane. Says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — loop lane. An unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The reach of a pattern. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule needing more than a pattern had to be written in Go beside the structure rules — which is why the shipped card carries three `POSTPROC` rules and no more. A reader expecting a long list of banned phrases should know that list lives in `basanite` on purpose.

Lanes are the one place the structure rules vary, and they vary by who will read the turn rather than by card. The interactive lane is any turn that is not a loop turn, chosen because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`, and adds `unverified_done` and `loop_ask`, because nobody is reading yet — a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran. The external lane, chosen when `cope-gate --pretool` scores prose an external write is about to post, drops the same four and swaps nothing in, because a ticket has a reader and no ending they can answer.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | `MessageDisplay` entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--pretool` | `false` | `PreToolUse` entry: score the prose an external write is about to post, warn-only |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: `interactive` (default) or `loop` |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## What lands on disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600`. One JSON record per violation, carrying the matched text and about 70 characters either side — the log quotes replies back.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600`. An empty file whose mtime is the refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600`. The rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

`card/claude_voice.effigy` is the source. `effigy` owns the `.effigy` grammar; `make rules` regenerates `card/rules.json`; `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file: the always-on rules and the evidence-gated ones are printed by different paths, and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be healthy. The authoritative list of rules that really are discarded unrendered is empty when the card is healthy.

`@gate <rule_id> off — <why>`, one per line in the card header. The id has to be one the gate has, and a wrong id is reported at load rather than ignored. A declined rule still runs; only this card's score drops it, so a backfill reports what it would have caught.

`@shape <id>: <selector> <predicate> — <why>`, one per line in the card header. The id must not collide with one the gate already has. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

The reason after the dash is required in both forms, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project. Hits per character is the metric worth watching rather than the share of turns hit, because that second number tracks how long the turns were. `MEASUREMENTS.md` has the rates.

## Known limits

`labelled_opening` is not a tagger. It matches a shape, not a part of speech, and a paragraph opening on a legitimately short clause reads to it the same as a label.

`ask_not_last` says nothing about the order of several asks. It sees whether one sits above prose that carries on past it, and nothing about which of two questions should have come first.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — a description of the output and not a judgement of it. The discrimination test is where the judgement lives, and it covers voicing only.

The largest limit is the one named where the two axes are separated. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere to put it. Both directions are the card marking its own homework: a decline lowers that card's score and an assertion raises it, both are worth reading with the reason attached, and that is why the syntax requires one. `MEASUREMENTS.md` has the runs.

## Layout

| Path | What |
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

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
