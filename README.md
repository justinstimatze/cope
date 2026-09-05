# cope

cope ships one opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — and from the moment it is installed that card is on, scoring two things about every reply: what the sentences sound like, and whether the shape of the reply gives a reader at a terminal anything to answer. The card is a file, so it can be edited, or swapped for another.

One README, written again from each card. Every file under `demo/` is this page rewritten from a different card, same prompt and same facts, so the card is the only thing that changed — [`demo/README.md`](demo/README.md) indexes them, and reading two of them against each other shows what a card does faster than the rest of this page explains it. `demo/README.claude-maximal.md` comes from a card instructing every tic this model is measured to have, which makes the point in a single glance and is deliberately hard going.

## Voicing and structure

Voicing. What the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card and nowhere else — `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES`, `POSTPROC` — so swapping the card swaps every word of it. A sentence reaching for the balanced two-beat, two clauses of a length with a content word repeated across the joint, is a fault on this axis. This is the half with a measured result behind it.

Structure. The shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives `continue` something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Compiled into the binary, in `internal/scan`, identical whichever card is loaded. A reply that names an open problem in its closing paragraph and then stops is a fault on this axis. A reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

Two doors into the structure half, both written in the card header, both requiring a reason after the dash. `@gate <rule_id> off — <why>` declines a built-in rule: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its `VOICE` block asks for the balanced landing and the arriving close those two rules catch, and a card marked down for obeying itself measures nothing. `@shape <id>: <selector> <predicate> — <why>` states a structural rule the gate then checks: `card/demo/handoff.effigy` asserts `readable_cold`, `last paragraph words <= 60`, because its peak asks the reader to re-enter cold on the last block alone, and nothing compiled into the binary checks for that. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more. `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. `MEASUREMENTS.md` has the runs and the reasons their numbers do not carry more than they do.

## Delivery, phrasing, and shape

Delivery. A global `CLAUDE.md` is not the system prompt, though most readers who have edited one and watched it fail to stick believe it is. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. Compare an output style, which sits in the system prompt itself and which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works; `MEASUREMENTS.md` has the run.

Phrasing. A global `CLAUDE.md` banning the not-A-but-B flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic under discussion. Naming a surface form pushes the move into a variant of itself. That is the voicing side, and instruction alone does not close it.

Shape. A different complaint with a different cause — an ending that leaves the reader nothing to answer costs a whole round trip, and no instruction could have banned it as a phrase, because it is not a phrase.

The flip is an anecdote about one rule. What the claim rests on is the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it, and `MEASUREMENTS.md` has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

`cope-gate --setup`. Emits the output style, wires the three hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — `cope-gate --setup --dry-run` prints what it would change and writes nothing.

The menu entry. Named `claude_voice`, under `/config` -> `Output style`, which is the shipped card's id and not the word cope, so a reader scanning that menu for something called cope will not find it. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; `.claude/settings.local.json` takes `"outputStyle": "claude_voice"` to the same end.

`cope-gate --output-style`. The by-hand route, for anyone who would rather not have their settings written to: it writes the loaded card to `~/.claude/output-styles/<card>.md` and does nothing else, and `COPE_CARD=<name>` in front of it emits a different one.

Timing. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`.

The reason for that placement. An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through a hook.

The hooks. Three of them, none of them how the card arrives, all wired by `--setup` already — this is the block they occupy in `~/.claude/settings.json`.

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

The gain from each hook. `Stop` gets the reply scored after it is written; `cope-gate --refresher` on `UserPromptSubmit` gets the rules that have actually been firing restated mid-session, which a file written once cannot do; `cope-gate --pretool` scores the prose an external write is about to post, warn-only, so the call goes through and the model learns what it scored. The voice works without any of them, and these are the measurement half.

PATH. The commands above assume `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no `effigy` checkout, because `card/rules.json` is committed and compiled in. `--inject` remains for anyone who wants the superseded turn-zero delivery, and stands down on its own when a cope output style is active.

## Another voice

The gate reads `.effigy` directly, so a card is usable as written — no Python, no `effigy` checkout. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, rather than a session writing in the shipped voice while its config names another one.

The one to read first is `card/demo/lecturer.effigy`, which differs from the shipped card on register alone, which is what makes it usable as a discrimination rival and what the discrimination run measured. What a card changes is the voicing axis described at the top of this page; `MEASUREMENTS.md` has the numbers.

The shortest way to see what a card does without writing one is [`demo/README.md`](demo/README.md), where every file is this README written again from a different card, same prompt and same facts, so the voice is the only variable between them.

`card/demo/handoff.effigy` is the exception in that directory: a hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there.

## The scoring loop

`Stop`. `cope-gate` scores the reply just written, adds which rules fired to the session's rolling state, and writes one record per violation to the log. No prose goes into the rolling state — rule names and counts only.

`UserPromptSubmit`. `cope-gate --refresher` reads that rolling state, not the violations log, and injects the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. It falls back to the standing `CONTINUE TEST` when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. One mechanism, not a guarantee: the A/B in this repo does not separate the refresher from no refresher.

## Effigy notation, and what else exists

Three blocks do what a prose gate needs, which is why the cards are written in a character-card notation for game NPCs, used here off-label: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it. The notation lives at [`effigy`](https://github.com/justinstimatze/effigy).

[`basanite`](https://github.com/justinstimatze/basanite). One problem answered the other way round, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. Compare `basanite`, which measures instead — lemma frequency against a baseline over real transcripts, reporting what you have actually been leaning on lately and leaving the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable.

[`humanizer`](https://github.com/blader/humanizer). A skill, by a different author, that rewrites AI-sounding prose against 35 patterns taken from Wikipedia's "Signs of AI writing", the page WikiProject AI Cleanup maintains. It is called on a text and hands back a rewrite; cope fires at a hook, scores what was already written, and edits nothing. Its pattern list is the wider one, and a reader who wants a rewrite rather than a score should go there. The formatting patterns are where the two disagree on purpose: cope's `bold_label` rule banned `humanizer`'s bold mini-headings until 52 blind pairs put bold and bullets among the three things that decided a reply for this repo's reader, and the rule was deleted rather than tuned.

[`caveman`](https://github.com/JuliusBrussee/caveman). A separate project, by a different author, that compresses agent replies to cut output tokens — a fourth axis. cope shapes prose, `basanite` tracks vocabulary, `humanizer` rewrites, `caveman` shortens. A reader wanting fewer tokens rather than different structure should go there instead.

## The rules

Voicing, in the shipped card as `POSTPROC` patterns:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, compiled in because a pattern could not carry them:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.
- `repeated_opening` — three or more sentences in one reply opening on the same two words. Compare `cross_turn_repeat`, which reads the session window for a construction reused across turns; this one reads a reply against itself, and two is left alone because two is a rhythm.
- `fragment_run` — three consecutive sentences of five words or fewer with no finite verb in any of them. One fragment is emphasis and this repo's own register is full of them; a run of three is the staccato blind judges read as generated. Neither clipped demo card trips it, so neither declines it; a card that wants the run says so with `@gate`.

Structure, compiled in and the same under every card:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.
- `echoed_heading` — a heading of two or more content words whose first sentence below repeats every one of them, spending a line to say what the heading already said.

The grouping says something usable about the implementation. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules — which is why the shipped card carries three `POSTPROC` rules and no more. This page is written from `card/demo/fieldguide.effigy`, whose own three — `heading_headword`, `unrecorded_hedge`, `unpersuaded` — belong to that card and are not part of what an install gets. A reader expecting a long list of banned phrases should know the list lives in `basanite` on purpose.

Lanes are where the structure rules vary, and they vary by who is going to read the turn rather than by card. The interactive lane is any turn that is not a loop turn, chosen because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself; it drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`, and adds `unverified_done` and `loop_ask`, because nobody is reading yet — a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. The external lane is chosen when `cope-gate --pretool` scores prose an external write is about to post; it drops the same four and swaps nothing in, because a ticket has a reader and no ending they can answer, read days later by somebody who was not in the session, which is the condition every surviving rule was written for.

Clustering, printed by the `Stop` hook, by `--check` and by `--pretool`. Breadth is three or more distinct rules landing on one paragraph. Density is one rule landing on one paragraph three or more times. When both hold, breadth is reported, since naming three rules tells a reader to rewrite the block rather than hunt one construction. Every rule fires on its own and knows nothing about the others, so a report is otherwise a flat list a reader works down hit by hit: three hits across three paragraphs are three small edits, and three inside one paragraph are one paragraph to write again. The density half has a measurement behind it. `--check` over 107 tracked documents produced 114 `flip` hits, seven of them worth changing. Each of those seven showed as three in a paragraph rather than as anything about the form. Three is the floor on both conditions, because two rules on a paragraph is ordinary and two hits of one is a coincidence a reader can see unaided. Nothing about what fires or how it is scored changes.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; `positive` is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit `2` on a violation whose action is reject (default warn-only) |
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

## On disk

`$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600`. One JSON record per violation, carrying the matched text and about 70 characters either side — the log quotes replies back at you, and that is what it is for.

`$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600`. An empty file whose mtime is the refresher clock.

`$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600`. The rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

## Editing the card

`effigy` owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart.

The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately and not against the card file, since the always-on rules and the evidence-gated ones go out on different paths and no code path renders their union — so a card may hold more `NEVER` rules in total than the budget and still be healthy. Nothing in the shipped card is currently discarded unrendered.

`@gate <rule_id> off — <why>`, one per line in the card header. The id has to be one the gate already has, and a wrong id is reported at load rather than ignored. A declined rule still runs and only this card's score drops it, so `--backfill` still reports what it would have caught.

`@shape <id>: <selector> <predicate> — <why>`, one per line in the card header. The id must not collide with a rule the gate already has, and a collision is reported at load. The selectors are `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. The predicates are `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`. A violation is reported in the card's own words rather than in any sentence the binary supplies.

The reason after the dash is required in both forms, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once and is how the rules were chosen. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project. Hits per character is the metric to watch: the share of turns hit tracks how long the turns were. `MEASUREMENTS.md` has the rates.

## Known limits

`labelled_opening` is not a tagger, and misreads openers a parser would resolve. `ask_not_last` says nothing about the order of several asks once they are all present.

The hit rate is roughly four fifths structure. The A/B run found that four fifths tracks what a reply was for rather than how it was written, which makes it a description of the output and not a judgement of it — the judgement lives in the discrimination test, and that covers voicing only.

The largest limit is the one the top of this page names as the boundary of card-authored rules. A card declines a built-in rule and writes one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. The compiled rules remain the only place a check like `clause_symmetry` can live, so a card wanting something outside both that vocabulary and a `POSTPROC` regex has nowhere to put it. Both directions are the card marking its own homework — a decline lowers that card's score and an assertion raises it — which is why the syntax requires a reason and why the reason is the part to read. `MEASUREMENTS.md` has the caveats.

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

Written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
