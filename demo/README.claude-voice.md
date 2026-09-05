# cope

A cope is the upper half of a foundry mould — the half carrying the shape being cast into — and this one ships an opinionated card that is on the moment you install it, checking both how a reply sounds and how it is shaped; the card is a file, so you can edit it or swap it for another.

Start at [demo/README.md](demo/README.md): every file in that directory is this same README written again from a different card, same prompt and same facts, so the card is the only thing that varies, and reading two of them side by side shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card that instructs every tic this model is measured to have, which makes the point in one glance and is deliberately hard going.

## Two different jobs

Voicing is what the sentences sound like: register, rhythm, diction, where flair is licensed. It lives in the card, and swapping the card swaps every word of it. A sentence reaching for the balanced two-beat — two comma-joined clauses of near-equal length repeating a word across the joint — is a voicing problem.

Structure is the shape of the reply as a thing the reader has to use: where the decision sits, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. It is compiled into the binary and is the same whichever card is loaded. A reply that names an open problem in its last paragraph and then stops, leaving `continue` nothing to refer to, is a structural problem. The same reply can be clean on one axis and bad on the other, which is why they are kept apart.

A card reaches into the structure half in two directions. It can decline a built-in rule with `@gate <rule_id> off — <why>`, one per line in the card header, because a card whose `VOICE` block asked for the balanced landing was being marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason. It can also assert a rule of its own with `@shape <id>: <selector> <predicate> — <why>`, because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — since its peak asks the reader to re-enter cold and read the last block only. The vocabulary counts words and sentences and asks whether a block poses a question, and nothing more, so a card wanting a check outside that and outside a `POSTPROC` regex still has nowhere to put it. [MEASUREMENTS.md](MEASUREMENTS.md) has the runs.

## The problem

Most readers arrive believing their global `CLAUDE.md` is the system prompt. It is not. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style goes into the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works — that was measured, and [MEASUREMENTS.md](MEASUREMENTS.md) has the run.

Instruction alone does not fix the phrasing. A global `CLAUDE.md` banning the "not A, it's B" flip is read every turn, and the flip appeared twice in the session that built this while the ban was the topic. Naming a surface form pushes the move into a variant of itself.

The structural side is a different complaint with a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no phrasing instruction could have banned it, because it is not a phrasing habit.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test: a reader is shown only a voice's own description of itself and picks which of two replies was written under it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`cope-gate --setup` emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. `cope-gate --setup --dry-run` prints what it would change and writes nothing.

Then pick the style under `/config` -> Output style. The entry is named `claude_voice`, which is the shipped card's id — a reader scanning that menu for something called cope will not find it. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

If you would rather not have your settings written to, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` by hand, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. Nothing changes in the conversation you are already in.

The card goes here rather than into a hook because an output style sits at the end of the system prompt and the harness re-reminds the model of it as the conversation runs.

The hooks are the measurement half, and the card no longer arrives through one. This is the `hooks` block of `~/.claude/settings.json`:

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

`Stop` buys you a score on every reply after it is written. `UserPromptSubmit` buys you a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. The voice works without either.

These commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled into the binary. The older `--inject` delivery is superseded and off by default, kept for anyone who wants it, and it stands down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable as written and needs no Python and no effigy checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session quietly writing in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, which is what the blind discrimination run measured. What a card changes is the voicing half described at the top of this page; the structure rules stay where they are, minus whatever the card declines with `@gate`.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README written again from a different card against the same prompt and the same facts.

`card/demo/handoff.effigy` is the exception in that directory. It is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose — and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## What the hooks do differently

`Stop` runs `cope-gate` on the reply just written. It appends which rules fired to the session's rolling state and one record per violation to the log.

`UserPromptSubmit` runs `cope-gate --refresher`, which reads that rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts. When the session has no history yet it falls back to the standing `CONTINUE TEST` reminder. It stays quiet until the last injection has aged past `--refresh-every`, which defaults to 30 minutes.

So the mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot do. That is a mechanism and not a guarantee: the A/B run in this repo does not separate the refresher from no refresher.

`PreToolUse` runs `cope-gate --pretool`, which scores the `description`, `body` or `content` field an external write is about to post, matched against the Linear save tools in the settings block. It is warn-only — it returns `additionalContext` and never a `permissionDecision`, so the call goes through and the model learns what the prose scored. It writes no session state and scores in the `external` lane.

## Why effigy notation, and what else is out there

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs. `POSTPROC` is regex rules with a `warn` action applied after generation. `WRONG` holds an anti-pattern beside its replacement. `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) is the same problem answered the other way round, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you. Which one fits is a question about mood more than correctness — a heavy hand for a habit that is annoying you today, a moving measurement if you would rather watch the drift than legislate it. They compose, on different hooks with no shared state, and running both is reasonable.

[humanizer](https://github.com/blader/humanizer) is a skill by a different author that rewrites AI-sounding prose against 35 patterns taken from Wikipedia's "Signs of AI writing", the page WikiProject AI Cleanup maintains. It is called on a text and hands back a rewrite; cope fires at a hook, scores what was already written, and edits nothing. Its pattern list is the wider one, and a reader wanting a rewrite rather than a score should go there. The formatting patterns are where the two disagree on purpose: cope's `bold_label` rule banned humanizer's bold mini-headings until 52 blind pairs put bold and bullets among the three things that decided a reply for this repo's reader, and the rule was deleted rather than tuned.

[caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses agent replies to cut output tokens — a fourth axis. cope shapes prose, basanite tracks vocabulary, humanizer rewrites, caveman shortens. A reader wanting fewer tokens rather than different structure should go there instead.

## The rules

Voicing rules, the ones about what the sentences sound like:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine. `warn`.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead. `warn`.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k. `warn`.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.
- `repeated_opening` — three or more sentences in one reply opening on the same two words. Two is left alone because two is a rhythm.
- `fragment_run` — three consecutive sentences of five words or fewer with no finite verb in any of them. One fragment is emphasis and this repo's own register is full of them; a run of three is the staccato blind judges read as generated. A card that wants the run says so with `@gate`.

Structure rules, the ones about the shape of the reply:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.
- `echoed_heading` — a heading of two or more content words whose first sentence below repeats every one of them, spending a line to say what the heading already said.

The grouping tells you something about the implementation. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries only three `POSTPROC` rules. If you came expecting a long list of banned phrases, that list lives in basanite on purpose, where it is a measurement rather than a ban.

The structure rules vary in one way, and it is not by card. It is by who is going to read the turn. The `interactive` lane is any turn that is not a loop turn, and it runs everything, because somebody is waiting at a terminal and the ending is where they decide what happens next. The `loop` lane fires when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet, so a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran. The `external` lane, used when `--pretool` scores prose an external write is about to post, drops the same four and swaps nothing in — a ticket has a reader and no ending they can answer, and it is read days later by somebody who was not in the session, which is the condition every surviving rule was written for.

At the bottom of a report from the `Stop` hook, `--check` or `--pretool`, hits get clustered. Three or more distinct rules landing on one paragraph is breadth; one rule landing on one paragraph three or more times is density; when both hold you get breadth, since naming three rules tells a reader to rewrite the block rather than hunt one construction. Every rule fires on its own and knows nothing about the others, so a report is otherwise a flat list a reader works down hit by hit — but three hits across three paragraphs are three small edits, and three inside one paragraph are one paragraph to write again. The density half has a measurement behind it: `--check` over 107 tracked documents produced 114 `flip` hits of which seven were worth changing, and every one of the seven was visible as three in a paragraph rather than as anything about the form. Three is the floor on both conditions because two rules on a paragraph is ordinary and two hits of one is a coincidence a reader can see unaided. Nothing about what fires or how it is scored changes.

## Flags

| flag | default | what it does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; `positive` is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is `reject` (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--card-from-sample` | (empty) | print a prompt for writing a card from this writing sample; `-` reads stdin |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--check-lane` | (empty) | score `-check` in the given lane: `interactive` (default), `loop`, or `external` |
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

Three files, all mode `0600`, under `$XDG_STATE_HOME/cope`:

- `violations.jsonl` — one JSON record per violation, carrying the matched text and about 70 characters either side. The log quotes your replies back at you.
- `refresher-<session-id>` — an empty file whose mtime is the refresher clock.
- `session-<session-id>.json` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart.

The `NEVER` budget is 10. Anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file: the always-on rules and the evidence-gated ones go out through different paths and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be healthy.

A card can decline a built-in rule:

```
@gate <rule_id> off — <why>
```

One per line in the card header. The id has to be one the gate actually has, and a wrong id is reported at load rather than ignored. The reason after the dash is required. A declined rule still runs — only this card's score drops it — so a backfill still reports what it would have caught.

A card can also assert a rule of its own:

```
@shape <id>: <selector> <predicate> — <why>
```

Selectors are `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates are `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. The id must not collide with a rule the gate already has, and a colliding id is reported at load. The reason is required here too, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies. This card asserts `short_close` — the close is one or two sentences; past that the summary has become another section.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and it is how these rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project.

Watch hits-per-character rather than the share of turns hit. The second number tracks how long the turns were. [MEASUREMENTS.md](MEASUREMENTS.md) has the rates.

## Known limits

`labelled_opening` is not a tagger. It matches a shape, and a short verbless opener that genuinely earns its place will still be flagged.

`ask_not_last` says nothing about the order of several asks. It only knows that one is sitting above prose that carries on past it.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written. It describes the output; it does not judge it. The judgement lives in the blind discrimination test, and that covers voicing only.

The largest limit is the one the section on the two jobs named. A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. Both directions are the card marking its own homework — a decline lowers that card's score, an assertion raises it — which is why the syntax requires a reason attached to each. [MEASUREMENTS.md](MEASUREMENTS.md) has what was run.

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

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
