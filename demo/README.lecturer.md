# cope

A cope is the upper half of a foundry mould, the half carrying the shape being cast into, and this one arrives as a single opinionated card that is switched on the moment you install it and that checks two unrelated things about a reply — what the sentences sound like, and how the whole thing is shaped — with the card itself sitting in a file, for the day you want to argue with it.

Start with [demo/README.md](demo/README.md): every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that moved between them, and reading two of them against each other shows what a card does faster than the rest of this page manages to say it. One of them is [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model has been measured to have, and it makes the point in a single glance provided you go in knowing it is deliberately hard going.

## Sound, and shape

You have read a reply that was wrong in a way you could not name. Voicing is what the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, where flair is licensed — and all of it lives in the card, so swapping the card swaps every word of it; that is the half with a measured result behind it. Structure is the shape of the reply as a thing you have to use: where the decision sits, whether the ending gives `continue` anything to refer to, whether an ask is stranded above three paragraphs of prose, whether a claim that the work is done carries anything that could have shown it. Structure is compiled into the binary and is identical whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing problem; a reply that names an open problem in its final paragraph and then stops, leaving `continue` nothing to point at, is a structural one. The same reply can be immaculate on one axis and useless on the other, which is the whole reason for keeping them apart.

The card reaches into the structural half in exactly two directions, both written in the card header. `@gate <rule_id> off — <why>`, one per line, declines a built-in rule: a card whose `VOICE` block asked for the balanced landing that `clause_symmetry` catches was being marked down for obeying itself, and that is the hole this closes. `@shape <id>: <selector> <predicate> — <why>`, one per line, states a structural rule the gate then goes and checks, because a card can commit to how a reply ends and until this there was nowhere for that commitment to be tested. The vocabulary counts words and sentences and asks whether a block poses a question, and a card wanting a check outside that vocabulary and outside a `POSTPROC` regex still has nowhere to put it.

## Where the instruction sits

You have written rules into a global `CLAUDE.md`, watched them hold for two turns and dissolve, and concluded that the model is not really reading them. It is reading them. What it is not doing is remembering them, because that file is not the system prompt whatever its name suggests: it arrives as one message attached to the first turn, and every turn after that buries it a little deeper. An output style goes into the system prompt itself, and the harness re-reminds the model of the system prompt as the conversation runs. Take one card and move it from the first place to the second without editing a syllable of it, and that move is most of why cope works. It was measured; [MEASUREMENTS.md](MEASUREMENTS.md) has the run.

Placement is most of it, and it is not all of it, because instruction alone does not fix phrasing. A global `CLAUDE.md` banning the not-A-it's-B flip is read on every single turn. That flip appeared twice in the session that built this tool, while the ban was the subject under discussion. Name a surface form and the move relocates into a variant of itself, which is the voicing side of the complaint and the reason a gate reads output instead of trusting instruction.

The other half of the complaint is not a phrasing problem at all. An ending that leaves you nothing to answer costs a whole round trip — you type continue, and the reply spends its first paragraph working out what continue meant. No instruction bans that, because the fault is in the arrangement, and arrangement has no surface form to ban.

The flip is an anecdote about one rule. What the claim actually rests on is the blind discrimination test: a reader is shown a voice's own description of itself, plus two replies, and picks which of the two was written under it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate and the caveats that come attached to it.

## Installing it

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

The first line puts the hook binary on your machine. The second does the whole install: `cope-gate --setup` emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up before touching it, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses outright to edit a settings file that does not parse — and `cope-gate --setup --dry-run` prints what it would change while writing nothing, if you would rather look first.

The menu is the step that actually turns the voice on, and it is the step most likely to go wrong. Open `/config`, go to `Output style`, and select the entry named `claude_voice`. That name is the shipped card's id, and nowhere in it is the word cope, so a reader scanning that list for something called cope will not find it, will conclude the install failed, and will be wrong. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the route; the same choice can be written as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

For anyone who would rather not have a program editing their settings, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` on its own, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start. A fresh selection, or a card you have just re-emitted, therefore applies in your next session or after `/clear`, and the conversation you are sitting in will not change — which is the expected behaviour and no kind of failure. The reason the card lives in a style at all is that a style goes at the end of the system prompt and the harness keeps re-reminding the model of it, which is why the card lands here and did not land when it was delivered by a hook. (`--inject`, that superseded delivery, remains for anyone who wants it, and stands down on its own when a cope output style is active.)

The hooks are a separate matter, the voice works without any of them, and they are the measurement half. This is the hooks block of `~/.claude/settings.json`:

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

`Stop` buys you a score on the reply after it has been written. `UserPromptSubmit` buys you a mid-session restatement of the rules that have actually been firing, which is the one thing a file written once cannot do.

Both commands assume that `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook which silently does nothing usually wants the absolute path instead. From a clone, `make install` builds it, and no effigy checkout is needed anywhere, because `card/rules.json` is committed and compiled straight into the binary.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable exactly as written, with no Python and no effigy checkout in the loop. Drop one in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected, which is deliberate, because the alternative failure is a session writing happily in the shipped voice while its config names a different one, and you would never notice.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, which is what made it the card the discrimination run measured, and it is the card this page was written from. It also declines two rules with `@gate`: `clause_symmetry`, because its `VOICE` block asks for the balanced landing that rule catches, and `dangling_end`, because arriving and stopping is the register and there is no ask to move to the end. What a card changes is the sound of the sentences and nothing whatever about their arrangement, for the reasons set out where sound and shape were separated above.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this page again from a different card, same prompt and same facts, the voice the only thing varying.

One file in that directory is not a voice at all. `card/demo/handoff.effigy` is a hypothesis: it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it sitting there, and it is not a register to write in.

## What the hooks add

`Stop` fires when a reply finishes. It scores what was just written, appends the rule ids that fired to the session's rolling state, and writes one record per violation to the log.

`UserPromptSubmit` is where the interesting part happens. It reads the rolling state — the session file, and never the violations log — then injects the card items gated on what has actually been firing, naming the counts back at the model. The mid-session text is therefore chosen from measured output rather than fixed in advance, and that is a thing a pasted `CLAUDE.md` cannot do at any placement whatsoever. Where a session has no history yet it falls back to the standing `CONTINUE TEST`, and it stays quiet altogether until the last injection has aged past `--refresh-every`. One mechanism, and no guarantee: the A/B in this repo does not separate a session with the refresher from a session without one.

`PreToolUse` is in the block above, matched against the Linear save tools. It scores the description, body or content field an external write is about to post, warn-only — it returns `additionalContext` and never a `permissionDecision`, so the call goes through and the model gets told what its prose scored. It writes no session state, and it scores in the external lane.

## Why a game notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, and it is used here entirely off-label. Three of its blocks turn out to be what a prose gate needs and nothing else: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside the replacement for it, and `TEST` holds a named question with a failing example and a passing one, which is how a rule names a move instead of naming one wording of that move. Nothing had to be invented. The shapes were already there.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem from the other end and is the one to reach for when this one is too blunt: cope bans, so a rule fires or it does not, the card says never, and the register is fixed the moment you pick it, whereas basanite measures lemma frequency against a baseline over real transcripts and reports what you have been leaning on lately, leaving the judgement where you left it — awareness, as its own README calls it. Which of the two fits is a question about mood more than about correctness: a heavy hand for a habit that is annoying you today, a moving measurement for when you would rather watch the drift than legislate it. They compose, on different hooks with no shared state, and running both is reasonable.

[caveman](https://github.com/JuliusBrussee/caveman), a separate project by a different author again, compresses agent replies to cut output tokens, which is a third axis — cope shapes prose, basanite tracks vocabulary, caveman shortens — so if what you want is fewer tokens rather than different structure, go there instead.

## The rules

Voicing, from the card's `POSTPROC` block:

- `flip` (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, from the binary, where a pattern was not enough:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure, all of it from the binary:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

That grouping tells you something usable about the implementation. A `POSTPROC` pattern matches a span of text, so it can only ever describe wording; every voicing rule that needed more than a pattern had to be written in Go, sitting beside the structure rules. Which is why the shipped card carries three `POSTPROC` rules and no more, the card this page was written from carrying none at all, and why anyone arriving here expecting a long list of banned phrases should know the list lives in another tool on purpose.

The structure rules do vary in one way, and the variable is who is going to read the turn rather than which card is loaded. Interactive is the lane for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. A loop turn — opened by `/loop` or `/goal`, or by the sentinel a dynamic-pacing loop sends itself — drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`, because nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check, `unverified_done` and `loop_ask`, because a report saying the work is done has to say what it ran. The external lane, used when the `PreToolUse` entry scores prose an external write is about to post, drops those same four and swaps nothing in, since a ticket has a reader and no ending they can answer — it is read days later by somebody who was not in the session, which is the condition every rule that survives the drop was written for.

## Flags

| flag | default | does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default inject,hold; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
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

## What lands on disk

- `$XDG_STATE_HOME/cope/violations.jsonl`, mode `0600` — one JSON record per violation, carrying the matched text and about 70 characters either side.
- `$XDG_STATE_HOME/cope/refresher-<session-id>`, mode `0600` — an empty file whose mtime is the refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json`, mode `0600` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

The log quotes your replies back at you, and it is worth knowing that before you enable it: seventy-odd characters of context sit either side of every match. `--log` with an empty value turns it off.

## Editing the card

effigy owns the `.effigy` grammar. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected ones cannot quietly drift apart.

The `NEVER` budget is ten rules, and anything over it is reported at load rather than dropped in silence. That budget is charged against each injection separately rather than against the card file: the card injection prints the always-on rules, the refresher prints the evidence-gated ones, no code path renders their union, and a card may therefore hold more `NEVER` rules in total than the budget and still be perfectly healthy.

A card can decline a rule and a card can write one. `@gate <rule_id> off — <why>`, one per line in the card header, declines: the id has to be one the gate actually has, and a wrong one is reported at load rather than ignored. A declined rule still runs, and only this card's score drops it, so a backfill still tells you what it would have caught.

`@shape <id>: <selector> <predicate> — <why>`, one per line in the card header, asserts a rule of the card's own. The id must not collide with one the gate already has, and again a bad one is reported at load. The vocabulary is small and exact:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

`card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

The reason after the dash is required in both forms. A rule a card wrote and a rule a card refused are equally unreviewable without one, and a `@shape` violation is reported back in the card's own words rather than in any sentence the binary could have supplied.

## Calibrating

`cope-gate --backfill` scores every assistant turn in a transcript in one pass, and that is how the rules in the binary were chosen rather than guessed at. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same thing as one per project and should not be read as one. Watch hits per character rather than the share of turns hit, because that second number moves with how long the turns were, and a session of long replies will look worse than a session of short ones for reasons nobody cares about. [MEASUREMENTS.md](MEASUREMENTS.md) has the rates.

## Known limits

`labelled_opening` is not a tagger; it knows that a fragment is short and verbless and knows nothing about what the fragment means. `ask_not_last` says nothing about the order of several asks, only that one of them is stranded. Roughly four fifths of all hits are structural, and the A/B run found that the four fifths tracks what a reply was for rather than how well it was written, which makes the hit rate a description of the output rather than a judgement of it — the judgement lives in the discrimination test, and that test covers the sound of the sentences only.

The largest limit is the one described where sound and shape were separated. A card can decline a built-in rule and assert one of its own, and the vocabulary it asserts in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only home a check like `clause_symmetry` can have, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere at all to put it. Both directions are also the card marking its own homework: a decline lowers that card's score, an assertion raises it, both want reading with the reason attached, and that is precisely why the syntax will not let you write either one without it. The reasons the numbers do not carry further than they do are in [MEASUREMENTS.md](MEASUREMENTS.md).

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
