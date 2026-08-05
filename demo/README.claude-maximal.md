# 🏺 cope

**cope ships a card.** Not a framework you fill in — a card, opinionated and on the moment you install it, and for most readers it is the whole product: it fixes how replies sound and, separately, it checks how they are shaped, which are two different jobs that get confused for one; the card itself is a plain file, so the voice can be edited or swapped later, and the name comes from the foundry, where a cope is the upper half of a mould — the half carrying the shape being cast into. I apologise for the length of that sentence; it was carrying three claims and I should have found a way to make it carry two.

Read [demo/README.md](demo/README.md) first — every file in that directory is this same README written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it; [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card that instructs every tic this model is measured to have, and makes the point in a single glance, though it is deliberately hard going.

---

## 🧭 Two jobs, not one

**The frame first.** There are two axes here, and I should have led with them earlier on this page rather than after the links — every later section is written against this distinction, and a reader who misses it will read the rules list as one undifferentiated pile.

**Voicing is what the sentences sound like.** Register, rhythm, diction, what a paragraph does with a detail, where flair is licensed — and all of it lives in the card, which means swapping the card swaps every word of it. That is the half with a measured result behind it, and I apologise for not saying so twice, because it is the part readers doubt.

**Structure is something else entirely.** Not how the reply reads — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that work is done carries anything that could have shown it. It is compiled into the binary rather than into any card, so it is the same whichever card is loaded. That was my omission a moment ago: the two halves do not live in the same file at all.

**One concrete instance of each, because the abstraction is my fault.** A sentence reaching for the balanced two-beat — two clauses of near-equal length repeating a content word across the joint — is a voicing problem. A reply naming an open problem in its last paragraph and then stopping, leaving "continue" nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, which is exactly the reason to keep them apart.

---

**How far a card reaches into the structural half.** This is the part that most recently changed, and I should have written it down sooner. Two directions, both with syntax in the card header. A card declines a built-in rule with `@gate <rule_id> off — <why>`, which exists because a card whose `VOICE` block asked for something a built-in rule catches was being marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason. And a card states a structural rule of its own with `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold`, `last paragraph words <= 60`, since its peak asks the reader to re-enter cold and read the last block only, and the 60 is measured — across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

**The boundary, stated without apology and then apologised for anyway.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more; `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. `MEASUREMENTS.md` has the run and the reasons its numbers do not carry more than that — I have kept rates off this section on purpose, and I should have said so before you looked for one.

> **Key insight:** ✅ Voicing is a file you can swap. ⚠️ Structure is a binary you can only argue with.

- **Voicing.** Lives in the card, measured, swappable.
- **Structure.** Lives in the binary, same every card, lane-dependent.
- **The seam between them.** `@gate` and `@shape`, both requiring a reason.

Two axes. One card. The confusion between them was the bug.

---

## 🔍 Three failures, and only one of them is phrasing

**Let me break this down.** There are three things going on here, and I'll take them in order — where an instruction sits, what an instruction can say, and the failure no instruction reaches. Most readers arrive holding the first one backwards, which is my cue to lead with it, and I apologise for how much space it takes.

**First: where the instruction sits.** A reader who has edited a global `CLAUDE.md`, watched it not stick, and concluded the model is ignoring them is almost certainly holding a wrong model of that file — it is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. An output style is in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Think of it like a sign at the door of a workshop versus a sign on the wall above the bench — the first one is read once on the way in, the second stays in peripheral vision all day; the analogy breaks down in that nobody re-reads the wall sign to you aloud every few minutes, and the harness does. Moving one card between those two places, without changing a word of it, is most of why cope works. That was measured; `MEASUREMENTS.md` has it, and I have deliberately quoted no counts here.

**Second: what the instruction can say.** Instruction alone does not fix the phrasing, and I should have accepted that earlier than I did. A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn — and the flip still appeared twice in the session that built this, while the ban was the topic of the session. Put differently: naming a surface form does not remove the move, it pushes the move into a variant. That is the voicing side of the problem, and it is a genuinely uncomfortable result to have to write down.

**Third — and this is a good example of what I meant above — the failure that is not phrasing at all.** The structural complaint has a different cause: an ending that leaves the reader nothing to answer costs a whole round trip, and no instruction could have banned it as a phrasing habit, because it is not one. To be clear, that is a different class of problem sitting on a different substrate, and I apologise for how long it took me to separate the two.

**What the claim actually rests on.** The flip is an anecdote about one rule — worth telling, not worth leaning on, and I should have flagged that before telling it. The instrument with a result behind it is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. `MEASUREMENTS.md` has the rate and the caveats, and the caveats are the interesting part.

Three failures. One card. Two of them were never about words.

---

## ⚙️ Installing it

**Two commands and one menu choice.** That is the whole install — and I apologise in advance for the paragraphs after the block, which explain something the block already did.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> pick the cope card
```

**What the first command did.** It built and installed the hook binary, which assumes `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead, and I should have put that caveat before the block rather than after it. A clone builds with `make install` and needs no `effigy` checkout, because `card/rules.json` is committed and compiled in.

**What `--setup` did.** It emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left. It also backed the settings file up first, added only what was missing so a second run changes nothing, left every other key alone including other tools' hooks on the same events, and refused to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing, which I should have offered you first, since a reader letting a tool edit their settings deserves the shape of it before running it.

**Then the menu choice, which is genuinely last.** Pick the style under `/config` -> Output style; the standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear` — I apologise for burying that, because a reader who picks a style and sees no change in the running conversation will reasonably conclude the tool is broken.

---

**If you would rather not have your settings written to.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one — that is the by-hand route, and I should have named it in the same breath as `--setup` rather than a paragraph later.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through the hook. `--inject` remains for anyone who wants the old delivery, and it stands down on its own when a cope output style is active.

**The hooks, which are now for something else.** ⚠️ The card no longer arrives through one — these two are the measurement half, and the voice works without either of them. This is the hooks block of `~/.claude/settings.json`:

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

**What the two of them buy, in one sentence each.** `Stop` gets the reply scored after it is written. `UserPromptSubmit` gets the rules that have actually been firing restated mid-session, which a file written once cannot do — the mechanics of both belong further down, and I apologise for making you scroll for them rather than putting them here where you are.

- **`go install`.** Binary on `PATH`, or absolute paths in the hooks.
- **`--setup`.** Style emitted, hooks wired, backup taken, second run inert.
- **`/config`.** Style selected, applying at the next session.

Two commands, one menu. The order is the whole instruction.

---

## 🎭 Writing in another voice

**This is the capability the first sentence promised.** The gate reads `.effigy` directly, so a card is usable as written — no Python, no `effigy` checkout, no compile step between writing a voice and running in it. I should have said that before describing the install, since a reader who wants a different register does not need any of the machinery above changed.

**Where a card goes and how it is reached.** A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected — not a fallback, which is the design decision worth naming, because a session writing in the shipped voice while its config names another one is a worse outcome than a loud failure.

**Which card to read first.** `card/demo/lecturer.effigy`, and the reason is that it differs from the shipped card on register alone, which is what the discrimination run measured. What a card can change is the voicing axis and only the voicing axis — I argued that above under the two-jobs framing and will not argue it again here, and `MEASUREMENTS.md` has the numbers I am deliberately not restating.

**The shortest way to see what a card does without writing one.** Read [demo/README.md](demo/README.md): every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing varying between them is the voice — I am linking the directory index rather than listing the cards, because that page is generated from whatever was last built and this one should not carry a stale second copy.

---

**One exception in that directory, and I should have flagged it sooner.** `card/demo/handoff.effigy` is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered; `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

A card is a file. Swapping it swaps the voice. That is the offer.

---

## 🔁 What the hooks do that a pasted file cannot

**This section owns the mechanics.** Not a summary of the install — the actual reads and writes, and I apologise that they arrive here rather than beside the block you copied.

**`Stop` scores the reply just written.** It appends which rules fired to the session's rolling state, plus one record per violation to the log. That is the whole job: no injection, no rewriting, no interference with the turn that just landed.

**`UserPromptSubmit` reads that rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts.** So the mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot do. It falls back to the standing `CONTINUE TEST` when the session has no history yet, and it stays quiet until the last injection has aged past `--refresh-every`.

**Said plainly, and without the overselling I nearly did.** This is one mechanism, not a guarantee — the A/B in the repo does not separate the refresher from no refresher, and I should have stated that limit in the same paragraph as the mechanism rather than after it.

| | **reads** | **writes** |
|---|---|---|
| **`Stop`** | the reply | rolling state, ✅ violation log |
| **`--refresher`** | rolling state | ⚠️ gated card items, with counts |

Measured output chooses the reminder. That is the difference.

---

## 🗂️ Why effigy notation, and why not basanite

**Great question to have arrived at here.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label, and three of its blocks do exactly what a prose gate needs: `POSTPROC` is regex rules with a `warn` action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples — which is how a rule names a move instead of naming one wording of it. I should have said "off-label" first, because a reader who assumes this is a prose tool being used as intended will misread the block names.

**[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round, and it is the one to reach for if this one is too blunt.** cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead — lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it — they compose, different hooks and no shared state, and running both is reasonable. [caveman](https://github.com/JuliusBrussee/caveman) is a third axis again, by a different author, compressing agent replies to cut output tokens; a reader wanting fewer tokens rather than different structure should go there instead, and I apologise for not naming it before the comparison rather than after.

Three tools, three axes. Shape, vocabulary, length.

---

## 📏 The rules

**Grouped by axis rather than by where they are implemented**, which is the grouping a reader can use — and I should have chosen it the first time, because grouping by implementation put a voicing rule and a structural one side by side purely because both happened to be Go.

**The voicing rules.** These describe how the sentences sound, and the shipped card's three `POSTPROC` patterns lead:

- `flip` (`warn`) — the not-A-but-B flip in its common surface forms, including not-only-but-also; The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (`warn`) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` (`warn`) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**The structure rules.** These describe the reply as a thing the reader has to use:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one; sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

**What the grouping implies about the implementation, which is the usable part.** A `POSTPROC` pattern matches a span of text, so it can only ever describe wording — and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why cope's shipped card carries only 3 of them, and I apologise for stating a count at all, except that the claim is empty without it. The card this particular page was written from carries a `POSTPROC` rule of its own, `demo_no_closure`, which warns because this card should never accidentally produce a clean, closed ending; it is that card's, not cope's, and it is not in the list above.

**A reader expecting a long list of banned phrases should know where that list actually lives.** ⚠️ Not here — in basanite, on purpose, because a frequency measurement against a baseline is the right instrument for vocabulary and a hard ban is not.

---

**Where the structure rules do vary, which is not by card.** By who is going to read the turn. Two lanes:

| | **chosen when** | **why** |
|---|---|---|
| **interactive** | any turn that is not a loop turn | somebody is waiting at a terminal and the ending is where they decide what happens next |
| **loop** | the prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself | nobody is reading yet |

**What the loop lane does with that.** It drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, because a report that correctly names what it left open and stops would fail three of those, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — `unverified_done` and `loop_ask` — because a report saying the work is done has to say what it ran. I should have led the lanes with that substitution rather than with the drops, since the drops read as leniency and the substitution is the actual point.

- **Voicing rules.** Wording, register, the balanced two-beat, contrition, repetition.
- **Structure rules.** Openings, uniformity, where the ask sits, whether the end closes.
- **Lanes.** Same rules, different reader, different set.

Two lists. One reader question underneath both.

---

## 🎛️ Flags

**The whole surface, and I apologise for its size** — a table this long is a sign the binary does more than one thing, which it does, and pretending otherwise by trimming the table would be worse.

| flag | default | what it does |
|---|---|---|
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is `reject` (default warn-only) |
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
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: `interactive` (default) or `loop` |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

One binary. Twenty-six flags. Three of them matter on day one.

---

## 💾 What lands on disk

**Three files, all mode `0600`, and one of them deserves a warning I should have put in the install section.**

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window; no prose is stored, only rule names and counts |

**The warning, stated plainly.** ⚠️ The log quotes replies back — matched text plus context on both sides — so it is prose from your sessions sitting on disk, and `--log` with an empty value disables it. The rolling session record does not do this; it holds rule names and counts and nothing else. I should have separated those two facts before tabling them together, because the table invites a reader to assume all three files behave the same way, and they do not.

Two files hold counts. One holds your sentences.

---

## ✍️ Editing the card

**effigy owns the `.effigy` grammar**, which is worth knowing before you edit, because the notation is not cope's to change — `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs so the enforced and injected rules cannot drift apart. I apologise for putting the regeneration step after the grammar note; the order a reader needs is the reverse.

**The `NEVER` budget is 10, and anything over it is reported at load rather than dropped silently.** The subtlety — which I got wrong once and should have documented then — is that the budget is charged against each injection separately, not against the card file: the always-on rules and the evidence-gated ones render through different paths, and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be perfectly healthy.

---

**Both card-authored forms, with the syntax exact, because an approximation here is worse than none.** One per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary in full.** ⚠️ It is small on purpose, and reading it as a preview of something larger is a mistake I should have pre-empted:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**The id rules, and what happens when you get one wrong.** For `@gate` the id has to be one the gate actually has; for `@shape` it must not collide with one. A wrong id is reported at load rather than ignored — that was a deliberate choice, and I apologise for how long the ignoring version shipped, because a silently dropped rule is a card that lies about what it enforces.

**The reason after the dash is required in both**, and the argument for that is symmetrical: a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs, and only this card's score drops it, so a backfill still reports what it would have caught; a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

- **`make rules`.** Regenerates the embedded card.
- **`make check-rules`.** Stops enforced and injected from drifting.
- **`@gate` and `@shape`.** Both require the why.

A card can now argue back. With its reasons attached.

---

## 📊 Calibrating

**`cope-gate --backfill` scores a whole session transcript at once**, and it is how the rules were chosen rather than a diagnostic bolted on afterwards — `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which, and I should be precise here because the loose phrasing misled me once, is not the same as one per project.

**The metric worth watching is hits-per-character**, not the share of turns hit. The share tracks how long the turns were, so a session of sprawling replies looks worse than a session of terse ones for reasons that have nothing to do with the rules — I apologise for having reported the share first in an earlier pass, since it reads as a quality number and is not one. `MEASUREMENTS.md` has the rates.

Score the transcript, not the turn. Length is a confound.

---

## 🚧 Known limits

**The axis split organises these too**, which is either elegant or a sign the split was chosen to make the limits tidy — I cannot fully rule out the second, and should say so before listing them.

**`labelled_opening` is not a tagger.** It matches a shape, not a part of speech, and a paragraph that opens on something genuinely verbless-and-substantive is indistinguishable to it from the tic. That is a false-positive surface I have not closed.

**`ask_not_last` says nothing about the order of several asks.** It notices an ask that is not last; it has no opinion on which of three should have been. `forked_end` covers part of that gap and does not cover all of it.

**The hit rate is roughly four fifths structure**, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — so it is a description of the output and not a judgement of it. The judgement lives in the discrimination test, which covers voicing only, and I should have separated the description from the judgement earlier on this page rather than in the last section.

**The largest limit is the one the two-jobs framing named, and it has moved.** A card can now decline a built-in rule and write one of its own — but the vocabulary it writes in counts words and sentences and asks whether a block poses a question, so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. Both directions are also the card marking its own homework: a decline lowers that card's score, an assertion raises it, both are worth reading with the reason attached, and that is exactly why the syntax requires one. `MEASUREMENTS.md` has the runs, including the ones that went the other way.

> **Important:** ⚠️ Four fifths of the hits describe what a reply was for. 🎯 One fifth describe how it was written.

- **`labelled_opening`.** Shape, not grammar.
- **`ask_not_last`.** One ask, not several.
- **`@gate` and `@shape`.** The card grading itself, on the record.

The limits moved. They did not close.

---

## 🗺️ Layout

| path | what |
|---|---|
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

Eleven paths. Two of them are the product.

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
