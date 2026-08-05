# 🏺 cope

**What cope is.** cope ships an opinionated card — a *cope*, in the foundry sense, being the upper half of a mould, the half carrying the shape being cast into — which is on the moment it is installed and is the whole product for most readers, and which checks two different jobs rather than one, how a reply *sounds* and how a reply is *shaped*, and only after that is it a file on disk you can edit or swap for another. That sentence is carrying more than one sentence should, and I apologise for the density of it — I should have found a way to split it without losing the ordering, since the ordering is the part that matters.

**Before anything else, go and look.** Every file in [demo/README.md](demo/README.md) is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed — read two of them against each other and you will understand what a card does faster than the rest of this page explains it, and I should concede that plainly rather than burying it in a later section. One of them is [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have; it is deliberately hard going.

---

## 🎭 The two things

**The frame, up front.** There are two axes here, and every later section on this page is written against them — I should have said so before the install instructions in earlier drafts, and I apologise for the confusion that caused.

**Voicing is the first.** Voicing is not the shape of the reply — it is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely, across VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC, so swapping the card swaps every word of it. A sentence reaching for the balanced two-beat — two clauses of near-equal length repeating a content word across the joint — is a voicing problem, and it is the half with a measured result behind it.

**Structure is the second, and it does not move.** Structure is not phrasing — it is the shape of the reply as a thing the reader has to *use*: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. A reply that names an open problem in its closing paragraph and then stops is a structural failure, and it is compiled into the binary, so it is the same whichever card is loaded. The same reply can be clean on one axis and bad on the other, which is precisely why they are kept apart — and I should have made that separation explicit sooner.

---

**How far a card reaches into the structure half.** This is the part that most recently changed, so let me be exact rather than gestural — and I apologise if earlier phrasing implied more reach than exists. A card can move in two directions, both from its header, one rule per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**Declining, first.** `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch — a card whose voice asked for something a built-in rule catches was being marked down for obeying itself, which was my mistake to leave unaddressed. **Asserting, second.** `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way; the 60 is measured, at 33 words median and 56 at p90 across 43,155 assistant replies. The boundary is simple: the `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more.

> **Key insight:** A card owns its voice outright and rents a narrow corner of its structure.

- **Voicing** lives in the card, swaps with the card, and is the measured half.
- **Structure** lives in the binary and does not vary by card, only by lane.
- **A card can decline a compiled rule and assert a small one of its own** — with a reason, always.

Two axes, one file, one binary, and a seam between them.

---

## 🧱 The problem

**Three things are going on here, and I will take them in order.** Where an instruction sits, what an instruction can say, and the failure no instruction reaches — the first is the largest, most readers arrive holding it backwards, and I should have led with it long before I did.

**First, where it sits.** A reader who has edited a global `CLAUDE.md` and watched it not stick almost certainly believes that file is the system prompt. It is not — it arrives as one message attached to the first turn, and the conversation buries it under everything written after it. An output style is *in* the system prompt itself, which the harness re-reminds the model of as the conversation runs. Think of it like a note taped to the inside of the front door versus a note written on the door itself: the first is read once on the way in and forgotten by the kitchen; the second is passed every time you move through the house. The analogy breaks down here, of course, because nobody re-reads a door. Moving one card between those two places, without changing a word of it, is most of why cope works — it was measured, and [MEASUREMENTS.md](MEASUREMENTS.md) has the run.

**Second, what an instruction can say.** Instruction alone does not fix the phrasing, and I should have caught this earlier than I did. A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn — and the flip still appeared twice in the session that built this, while the ban *was the topic*. Naming a surface form does not remove the move; it pushes the move into a variant. That is the voicing side of the complaint.

**Third, the failure no instruction reaches.** The structural side is not a phrasing habit at all — it is a different complaint with a different cause. An ending that leaves the reader nothing to answer costs a whole round trip, and no ban on a wording could have prevented it, which is my mistake to have conflated in earlier framings of this section.

---

**What the claim actually rests on.** The flip is an anecdote about one rule and I should not have let it stand in for evidence — the claim rests on the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it. The rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md). The blind *preference* runs are not the evidence for this and are not cited here: both sides of those were written under a card, so they compare two ways of writing well and cannot see a voice being swapped.

- **Placement** is most of the effect, and it was measured.
- **Instruction** shifts a habit into a variant rather than removing it.
- **Structure** is a separate failure that instruction never addressed.

The card works because of where it sits, not because it asks nicely.

---

## 📦 Install

**Two commands and one menu choice.** That is the whole of it, and I apologise for the length of the explanation that follows — the doing is short and the describing is not.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

**What `--setup` did, stated rather than sold.** It emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left — and before any of that it backed the settings file up, added only what was missing so a second run changes nothing, left every other key alone including other tools' hooks on the same events, and refused to touch a settings file that does not parse. A reader letting a tool edit their settings deserves the shape of it before running it and not after, which is my mistake for having buried it previously; `--setup --dry-run` prints what it would change and writes nothing. If you would rather nothing wrote to your settings at all, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

**Then the menu choice, which is where people stall.** Pick it under `/config` → **Output style** — the standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. ⚠️ A style is read once at session start, so a new selection or a re-emitted card applies at the *next* session or after `/clear`; a reader who picks a style and sees no change in the running conversation will conclude the tool is broken, and I should have said so before the command rather than after it.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation — which is why the card lands here and did not through the hook.

---

**Now the hooks, which are no longer how the card arrives.** This is the hooks block of `~/.claude/settings.json`:

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

**What the two buy you.** `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log; `UserPromptSubmit` restates mid-session the card items gated on what has actually been firing, naming the counts — which a file written once cannot do. ✅ The voice works without either of them; these are the measurement half, and I should have made that division of labour clearer at the top.

**Two operational notes I owe you.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. There is also `--inject`, the superseded SessionStart delivery, which remains for anyone who wants it and stands down on its own when a cope output style is active.

| | **what it is for** | **needed for the voice?** |
|---|---|---|
| **output style** | delivering the card | ✅ yes |
| **`Stop` hook** | scoring what was written | ⚠️ measurement only |
| **`UserPromptSubmit` hook** | evidence-gated reminders | ⚠️ measurement only |

- **Install** is two commands and a menu choice.
- **`--setup`** backs up, adds only what is missing, and touches nothing else.
- **A style applies at the next session**, which is the single most common confusion.

Install it, pick it, then start a new session.

---

## 🎙️ Writing in another voice

**The gate reads `.effigy` directly.** A card is usable exactly as written — no Python, no effigy checkout, no compile step in the loop — and I apologise for how long this took to be true. Drop a card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set.

**A name that resolves to nothing is an error.** Not a fallback — an error, and nothing is injected. ⚠️ The alternative would be a session quietly writing in the shipped voice while its config named another one, which is the failure mode that would waste the most of your time, and it was my mistake to have ever allowed it.

**Read `card/demo/lecturer.effigy` first.** It differs from the shipped card on register alone, which is exactly what makes it the useful one to read, and it is what the discrimination run measured. What a card can change is the voicing axis, as **The two things** set out above — the numbers are in [MEASUREMENTS.md](MEASUREMENTS.md) and I will not restate them here, having already restated too much.

**The shortest route is not writing a card at all.** [demo/README.md](demo/README.md) is this README written again from each card, same prompt and same facts, so the only thing varying between them is the voice.

**One exception in that directory, so you are not surprised by it.** `card/demo/handoff.effigy` is a hypothesis rather than a voice — it keeps the shipped card's handoff rules, drops everything about prose, and is meant to be run through `make pairs` against the full card rather than rendered; `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

- **Cards are `.effigy` files** and need no toolchain to use.
- **A bad card name fails loudly**, which is the behaviour you want.
- **`demo/` shows the effect** without asking you to write anything.

Pick a card, or read two, and the axis becomes concrete.

---

## 🔁 What the hooks do differently

**A file written once cannot notice anything.** That is the whole of the distinction, and I should have compressed the previous three sections into that sentence.

**`Stop`, first.** It scores the reply just written and appends which rules fired to the session's rolling state, alongside one record per violation to the log — no prose in the state file, only rule names and counts.

**`UserPromptSubmit`, second.** It reads that rolling record — not the violations log — and injects only the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. It falls back to the standing CONTINUE TEST when the session has no history yet, and it stays quiet until the last injection has aged past `--refresh-every`. That is the one thing a pasted `CLAUDE.md` cannot do.

**A calibration I owe you before you over-read that.** It is one mechanism and not a guarantee — the A/B in the repo does not separate the refresher from no refresher, and I should have volunteered that limitation in the same breath as the claim rather than after it. `--inject` on SessionStart is superseded and off by default.

- **`Stop`** records what fired.
- **`UserPromptSubmit`** reminds you of what fired, with counts.
- **Neither is required** for the voice itself.

The hooks are how the card learns what you actually do.

---

## 🗂️ Why effigy notation

**Three of its blocks do what a prose gate needs.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here entirely off-label — POSTPROC is regex rules with a `warn` action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a *move* instead of one wording of it. I should concede that borrowing an NPC format for this looks like a joke and is not one.

**And why [basanite](https://github.com/justinstimatze/basanite) is the wrong instrument here.** cope bans — a rule fires or it does not, the card says never, and the register is fixed the moment you pick it — where basanite measures lemma frequency against a baseline over real transcripts and leaves the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable. If what you want is fewer tokens rather than different structure, [caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses replies, and it is a third axis again.

Different tools for prose, for vocabulary, and for length.

---

## 📏 The rules

**Grouped by axis rather than by implementation**, because the axis is what a reader can act on — and I apologise for the earlier versions of this page that grouped them the other way.

**Voicing rules.**

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for *important* or *central*, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure rules.**

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies about the implementation.** This is the usable part, and I should have led the section with it. A POSTPROC pattern matches a span of text, so it can only ever describe *wording* — every voicing rule that needed more than a pattern had to be written in Go, beside the structure rules, in `internal/scan`. That is why cope's shipped card carries only **three** POSTPROC rules.

**If you arrived expecting a long list of banned phrases**, that list lives in another tool on purpose: basanite measures vocabulary against a baseline and reports what you have been leaning on, where cope bans a short, defended handful. ⚠️ A reader who wants breadth here will be disappointed, and I would rather disappoint them in this paragraph than after installation.

---

**The one place structure rules vary, and it is not by card.** It is by who is going to read the turn.

| | **interactive** | **loop** |
|---|---|---|
| **chosen when** | any turn that is not a loop turn | the prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself |
| **why** | somebody is waiting at a terminal and the ending is where they decide what happens next | nobody is reading yet |

**What the loop lane drops and adds.** It drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask` — because a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran.

> **Note:** The lane is decided by the transcript, not by a flag you remember to set.

- **Voicing rules** are the card's, and three of them are patterns.
- **Structure rules** are the binary's, and the same everywhere.
- **The lane** is the only axis on which structure moves.

Who reads the turn decides which ending counts as finished.

---

## 🚩 Flags

| flag | default | does |
|---|---|---|
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

Twenty-six flags, and most sessions use none of them.

---

## 💾 What lands on disk

**Three files, all mode `0600`.** I should say the important part before the table rather than after it, since it is the one that would change a reader's mind: ⚠️ the violations log quotes your replies back, carrying the matched text and about 70 characters either side.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window — no prose, only rule names and counts |

**`--log` with an empty value disables the quoting one entirely**, which is the lever to reach for if that trade does not suit you — and I should have offered the lever in the same paragraph as the warning rather than a line below it.

One file quotes you, one is a clock, one is counts.

---

## ✏️ Editing the card

**effigy owns the `.effigy` grammar.** `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart — a drift that would be invisible from either side alone, and I apologise for how long the repo went without that check.

**The NEVER budget is 10, and it is charged per injection.** ⚠️ Not against the card file — against each injection separately, since SessionStart prints the always-on rules and the refresher prints the evidence-gated ones and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be perfectly healthy. Anything genuinely over budget is reported at load rather than dropped silently, which was my mistake in an earlier version and is fixed.

---

**Both card-authored forms, with their syntax exactly as the gate reads it.**

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary, spelled out**, because an approximation of it is worse than none:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**The id rules, which are the ones people trip on.** For `@gate` the id has to be one the gate already has; for `@shape` it must not collide with one — and a wrong id is reported at load rather than ignored, which I should have made true from the start. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

**Two consequences worth having.** A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught; and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

- **`make rules` and `make check-rules`** keep injection and enforcement aligned.
- **The NEVER budget is per injection**, not per file.
- **Both `@gate` and `@shape` demand a reason**, and neither hides a bad id.

A card can argue with the gate, in writing, on the record.

---

## 📐 Calibrating

**`cope-gate --backfill` scores a whole session transcript at once**, and it is how the rules were chosen rather than a diagnostic bolted on afterwards — I should have said that before describing the sweep, since it is the reason the sweep exists.

**`tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`** — which is not the same as one per project, and reading it as one per project will skew what you conclude.

**Watch hits-per-character, not the share of turns hit.** ⚠️ The second number tracks how long the turns were more than it tracks anything about the writing, which is a good way to congratulate yourself for terseness. The rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

Normalise by length or the number is measuring length.

---

## ⚠️ Known limits

**The axis split is what organises these**, which is itself the admission — a limit that crosses both axes would be harder to see, and I have not gone looking for one.

**`labelled_opening` is not a tagger.** It approximates a shape, and it will misjudge sentences that a person would read correctly at a glance. **`ask_not_last` says nothing about the order of several asks.** It notices one sitting early and carries on.

**The hit rate is roughly four fifths structure**, and the A/B run found that four fifths tracks what a reply was *for* rather than how it was *written* — so the hit rate is a description of the output and not a judgement of it. The judgement lives in the discrimination test, and that covers voicing only. I should have separated those two claims much earlier on this page.

---

**The largest limit is the one The two things named.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

**Both directions are the card marking its own homework**, and I would rather say that plainly than have it discovered. A decline lowers that card's score and an assertion raises it; both are worth reading with the reason attached, and that is exactly why the syntax requires one. [MEASUREMENTS.md](MEASUREMENTS.md) has the runs and their caveats.

- **The rules approximate**, and two of them approximate visibly.
- **The hit rate describes output**, and the discrimination test judges voice.
- **A card grades itself** at both edges of what it may change.

The limits are in the notation, and the notation is on the page.

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

Eleven entries, and the card is the first of them.

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT — justin@justinstimatze.com
