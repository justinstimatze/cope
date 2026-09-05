# cope

An opinionated voice card, on the moment you install it, scoring two things a reply gets wrong independently — how it sounds and how it is shaped — and the card is a file, so swap it; the name is the foundry term, the upper half of the mould, the half carrying the shape being cast into.

Read [demo/README.md](demo/README.md): every file there is this page written again from a different card, same prompt and same facts, so the card is the only variable, and two of them read side by side settle what a card does faster than the rest of this page argues it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card instructing every tic this model is measured to have, and — heavy going by construction — makes the point in one glance.

## Two axes

Voicing is what the sentences sound like: register, rhythm, diction, where flair is licensed. It lives in the card. Swapping the card swaps every word of it, and this is the half with a measured result behind it. Structure is where the decision sits and how the reply ends. It is compiled into the binary, identical whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing fault; a reply naming an open problem in its last paragraph and stopping, leaving "continue" nothing to refer to, is a structural one. One reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

A card reaches into the structure half in two directions. `@gate <rule_id> off — <why>` declines a built-in rule, because a card whose `VOICE` asked for the balanced landing was marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` on those grounds. `@shape <id>: <selector> <predicate> — <why>` asserts a structural rule of the card's own, because a card's commitment about how a reply ends had nowhere to be checked; `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — against a peak that asks the reader to re-enter cold. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing further; a check outside that vocabulary and outside a `POSTPROC` regex has nowhere to live. Rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

## The problem

Where an instruction sits decides whether it survives. A global `CLAUDE.md` is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written afterwards. An output style sits in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, unchanged word for word, is most of why cope works. Measured; see [MEASUREMENTS.md](MEASUREMENTS.md).

What an instruction can say is the second limit. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip appeared twice in the session that built this, with the ban as the topic. Naming a surface form pushes the move into a variant.

The third failure is not phrasing at all. An ending that leaves the reader nothing to answer costs a round trip, and no instruction banning a habit reaches it.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. Rate and caveats in [MEASUREMENTS.md](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`--setup` emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key and every other tool's hooks on the same events alone, and refuses a settings file that does not parse. `--setup --dry-run` prints the diff and writes nothing.

Then the menu choice. `/config` -> Output style, and the entry is named `claude_voice` — the card's id, not the word cope. The standalone `/output-style` command was removed in Claude Code v2.1.91; `"outputStyle": "claude_voice"` in `.claude/settings.local.json` sets the same thing. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`.

Prefer not to have your settings written to: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

An output style goes at the end of the system prompt, where the harness re-reminds the model of it mid-conversation.

The hooks are the measurement half; the voice works without them. The card does not arrive through one.

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

`Stop` buys a score on every reply after it is written. `UserPromptSubmit` buys a mid-session restatement of the rules that have actually been firing, which a file written once cannot do. `--inject` is the superseded delivery, kept for anyone who wants it, and it stands down on its own when a cope output style is active.

The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in; a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout — `card/rules.json` is committed and compiled in.

## Writing in another voice

The gate reads `.effigy` directly. No Python, no effigy checkout. A card in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name resolving to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another.

Read `card/demo/lecturer.effigy` first: it differs from the shipped card on register alone, and it is what the discrimination run measured. A card changes the voicing axis, described above under the two axes.

[demo/README.md](demo/README.md) is the shortest route to seeing what a card does without writing one — every file under `demo/` is this README from a different card, same prompt, same facts.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered, and installed by `make cards` alongside the registers.

## What the hooks do differently

`Stop` scores the reply just written, appends which rules fired to the session's rolling state, and appends one record per violation to the log. The state file holds turn count, characters, and which rules fired over a 20-turn window — rule names and counts, no prose.

`UserPromptSubmit` reads that rolling state, not the violations log, and injects the card items gated on what has been firing, naming the counts. The mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot do. With no history yet it falls back to the standing `CONTINUE TEST`, and it stays quiet until the last injection has aged past `--refresh-every`. One mechanism, not a guarantee: the A/B in the repo does not separate the refresher from no refresher.

`--pretool` scores the `description`, `body` or `content` field an external write is about to post, matched against the Linear save tools in the settings block. Warn-only — it returns `additionalContext` and never a `permissionDecision`, so the call goes through and the model learns what the prose scored. It writes no session state and scores in the external lane.

## Why effigy notation, and what else exists

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label. Three of its blocks do what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples — which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite) is the same problem answered the other way round, and the one to reach for if this is too blunt. cope bans: a rule fires or it does not, the card says never, the register is fixed the moment you pick it. basanite measures — lemma frequency against a baseline over real transcripts, reporting what you have been leaning on lately and leaving the judgement to you; its README calls that awareness rather than prohibition. Which fits is a question about mood: a heavy hand when a habit is annoying you today, a moving measurement when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable.

[humanizer](https://github.com/blader/humanizer) is a skill, by a different author, rewriting AI-sounding prose against 35 patterns from Wikipedia's "Signs of AI writing", the page WikiProject AI Cleanup maintains. It is called on a text and hands back a rewrite; cope fires at a hook, scores what was already written, and edits nothing. Its pattern list is the wider one, and a reader wanting a rewrite rather than a score should go there. Formatting is where the two disagree on purpose: cope's `bold_label` rule banned humanizer's bold mini-headings until 52 blind pairs put bold and bullets among the three things that decided a reply for this repo's reader, and the rule was deleted rather than tuned.

[caveman](https://github.com/JuliusBrussee/caveman) is a separate project, by a different author, compressing agent replies to cut output tokens — a fourth axis. cope shapes prose, basanite tracks vocabulary, humanizer rewrites, caveman shortens. A reader wanting fewer tokens rather than different structure should go there instead.

## The rules

Voicing, from the shipped card's `POSTPROC` block:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, compiled:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating a content word across the joint: the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule reading the window rather than the reply, so it cannot fire until a session has a history.
- `repeated_opening` — three or more sentences in one reply opening on the same two words. `cross_turn_repeat` reads the session window for a construction reused across turns; this one reads a reply against itself, and two is left alone because two is a rhythm.
- `fragment_run` — three consecutive sentences of five words or fewer with no finite verb in any of them. One fragment is emphasis and this repo's own register is full of them; a run of three is the staccato blind judges read as generated. Neither clipped demo card trips it, so neither declines it; a card wanting the run says so with `@gate`.

Structure, compiled:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it. Interactive lane.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to. Interactive lane.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it. Interactive lane.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another. Interactive lane.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file. Loop lane.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself. Loop lane.
- `echoed_heading` — a heading of two or more content words whose first sentence below repeats every one of them, spending a line to say what the heading already said.

A `POSTPROC` pattern matches a span of text, so it can only ever describe wording; every voicing rule needing more than a pattern was written in Go beside the structure rules. Hence three `POSTPROC` rules in the shipped card. A reader expecting a long list of banned phrases should read basanite, above: the list lives in another tool on purpose.

The structure rules vary by lane — not by card, by who is going to read the turn. The interactive lane is any turn that is not a loop turn: somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself. It drops `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end` and adds `unverified_done` and `loop_ask`: nobody is reading yet, a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran. The external lane, used when `--pretool` scores prose an external write is about to post, drops the same four and swaps nothing in: a ticket has a reader and no ending they can answer, read days later by somebody who was not in the session, which is the condition every surviving rule was written for.

At the bottom of a report from `Stop`, `--check` or `--pretool`, hits cluster. Three or more distinct rules on one paragraph is breadth; one rule three or more times on one paragraph is density; when both hold, breadth wins, since naming three rules tells a reader to rewrite the block rather than hunt one construction. Every rule fires on its own and knows nothing about the others, so a report is otherwise a flat list worked down hit by hit — three hits across three paragraphs are three small edits, three inside one paragraph are one paragraph to write again. The density half is measured: `--check` over 107 tracked documents produced 114 `flip` hits of which seven were worth changing, and every one of the seven was visible as three in a paragraph rather than as anything about the form. Three is the floor on both conditions, because two rules on a paragraph is ordinary and two hits of one is a coincidence a reader can see unaided. Nothing about what fires or how it is scored changes.

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
| `--card-from-sample` | (empty) | print a prompt for writing a card from this writing sample; `-` reads stdin |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--check-lane` | (empty) | score `-check` in the given lane: interactive (default), loop, or external |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--pretool` | `false` | PreToolUse entry: score the prose an external write is about to post, warn-only |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

## On disk

- `$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600` — one JSON record per violation, carrying the matched text and about 70 characters either side. The log quotes your replies back.
- `$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600` — an empty file whose mtime is the refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. Rule names and counts, no prose.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json`; `make check-rules` is what CI runs, so the enforced and injected rules cannot drift.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. It is charged against each injection separately, not against the card file: the always-on rules and the evidence-gated ones render on different paths and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be healthy.

Two card-authored forms, one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

`@shape` selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. The `60` in `card/demo/handoff.effigy` is measured — across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

A `@gate` id has to be one the gate has; a `@shape` id must not collide with one. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, since a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project. Watch hits-per-character rather than the share of turns hit; that second number tracks how long the turns were. Rates in [MEASUREMENTS.md](MEASUREMENTS.md).

## Known limits

`labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — a description of the output, not a judgement of it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one stated above under the two axes, as it now stands: a card can decline a built-in rule and assert one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere to put it. Both directions are the card marking its own homework — a decline lowers that card's score, an assertion raises it — and both are worth reading with the reason attached, which is why the syntax requires one. See [MEASUREMENTS.md](MEASUREMENTS.md).

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

---

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
