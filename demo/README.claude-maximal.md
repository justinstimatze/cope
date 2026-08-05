# cope

**What this is.** cope ships an opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — that is on the moment it is installed and that checks two different jobs rather than one, how a reply sounds and how a reply is shaped; the card itself is a file, so it can be edited, or swapped for another one entirely.

Start at [demo/README.md](demo/README.md): every file in that directory is this README written again from a different card, same prompt and same facts, so the card is the only thing that changed — and reading two of them against each other shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is the one that makes the point in a single glance, written from a card that instructs every tic this model is measured to have, and it is deliberately hard going.

---

## 🧩 The two things

**There are two things being checked, and the whole page rests on keeping them apart.** ⚠️ Let me set the frame before I deliver it — and I apologise in advance, because the distinction takes three paragraphs to draw and I have not found a way to make it take one.

**First, voicing.** Voicing is not what a reply says — it's what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely, and swapping the card swaps every word of it. To be clear, this is the half with a measured result behind it, and I should have said so earlier rather than letting you infer it: the blind discrimination test is the instrument here.

**Second, structure.** Structure is not the sentences — it's the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives *continue* something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. It is compiled into the binary, in `internal/scan`, which means it is the same whichever card is loaded and varies only by lane. That was my mistake to leave implicit for as long as I did.

**Third, one concrete instance of each, because abstraction is where this argument usually dies.** A sentence reaching for the balanced two-beat — two comma-joined clauses of near-equal length repeating a content word across the joint — is a voicing problem. A reply naming an open problem in its last paragraph and then stopping, leaving *continue* nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, and I apologise for how obvious that sounds written down, because it is the entire reason to keep them separate.

---

| | **voicing** | **structure** |
|---|---|---|
| **what it is** | what the sentences sound like | where the decision sits |
| **where it lives** | the card | the binary |
| **swap the card** | ✅ changes | — unchanged |

**How far a card reaches into the structure half.** This is the part that most recently changed, and I should have led with it. Two directions, both stated in the card header. A card declines a built-in rule it disagrees with — `@gate <rule_id> off — <why>` — which exists because a card whose VOICE block asked for something a built-in rule catches was being marked down for obeying itself. And a card states a structural rule of its own that the gate then checks — `@shape <id>: <selector> <predicate> — <why>` — which exists because a card's own commitment about how a reply ends had nowhere to be checked at all.

**The boundary, without apology and without calling it a roadmap.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — it cannot express what the compiled rules express, so a card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

> **Key insight:** ⚠️ Two axes, two homes, one of them measured — and no rate on this page. [MEASUREMENTS.md](MEASUREMENTS.md) has the run and the reasons its numbers do not carry more than that.

**Key takeaways:**

- **Voicing** is sound, lives in the card, and is the measured half.
- **Structure** is shape, lives in the binary, and is card-independent.
  - Except at two seams: a card may decline a rule, or assert one.
  - Both seams require a stated reason.
- **The split** is what organises every section below.

Two axes. One card. Keep them apart.

---

## 🕳️ The problem

**There are three problems here, and I'll take them in order, largest first.** Most readers arrive holding the first one backwards, which is nobody's fault but mine for not saying so sooner.

**First, where an instruction sits.** A reader who has edited a global `CLAUDE.md` and watched it not stick almost certainly believes that file is the system prompt — and it is not. It arrives as one message attached to the first turn, and the conversation buries it under everything written after it. An output style, by contrast, is *in* the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a single word of it, is most of why cope works; that was measured, and [MEASUREMENTS.md](MEASUREMENTS.md) has the run.

**Second, what an instruction can say.** Instruction alone does not fix the phrasing — a global `CLAUDE.md` banning the *not A, it's B* flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic under discussion. Think of it like a speed limit painted on one stretch of road: the sign works exactly where it is legible and nowhere else, and drivers route around it. The analogy breaks down, of course, because the model is not routing around anything deliberately — it is that naming a surface form pushes the move into a variant, and I should have been clearer that this is a fact about wording and not about compliance.

**Third, the failure no instruction reaches.** The structural side is not a phrasing habit — it's a different complaint with a different cause, sitting on a different substrate. An ending that leaves the reader nothing to answer costs a whole round trip, and no ban on a phrase could have prevented it, because there was no phrase.

---

**What the claim actually rests on.** Great question to be asking at this point, and I apologise for taking so long to get here: the flip is an anecdote about one rule, and an anecdote is not evidence. The claim rests on the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it. [MEASUREMENTS.md](MEASUREMENTS.md) has the rate and the caveats, both of which are worth reading before quoting either.

**Key takeaways:**

1. **Placement** — a `CLAUDE.md` is a first-turn message, not the system prompt.
2. **Wording** — banning a surface form moves the move, it does not remove it.
3. **Shape** — an unanswerable ending is not a phrasing habit at all.

Three failures. One of them is not about words.

---

## 🔧 Install

**Two commands and one menu choice.** Not three steps and a config file — two commands and one menu choice, and I should have opened the page with that rather than with two sections of frame.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

**What `--setup` did.** It emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left. It backed the settings file up first, added only what was missing so a second run changes nothing, left every other key alone including other tools' hooks on the same events, and refused to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing, which is the flag to reach for if letting a tool edit your settings is not something you want to discover the shape of afterwards. I apologise for not putting that clause first.

**The by-hand route.** For anyone who would rather not have their settings written to: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one — no settings file is touched.

**Then the selection.** Pick the style under `/config` → **Output style**. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

---

> **Important:** ⚠️ A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. A reader who picks a style and sees no change in the running conversation will conclude the tool is broken — that was my mistake to leave to the footnotes, and it is the single most common way this install appears to fail.

**Why the card goes here rather than into a hook.** An output style sits at the end of the system prompt and the harness re-reminds the model of it during the conversation — which is why the card lands here and did not through the hook.

**Now the hooks, which are no longer how the card arrives.** This is the `hooks` block of `~/.claude/settings.json`:

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

**What the two hooks buy.** Not the voice — the voice works without either of them, and I should have said that before pasting the block. `Stop` scores the reply just written and records which rules fired. `UserPromptSubmit` reads that record and restates mid-session the rules that have actually been firing, which a file written once cannot do. These two are the measurement half, and nothing more.

---

**One caveat about `PATH`.** The commands above assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook which silently does nothing usually wants the absolute path instead, and I apologise for how much time that particular failure has cost people. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled into the binary. There is also a `SessionStart` delivery, `--inject`, superseded and off by default, which stands down on its own when a cope output style is active.

**Key takeaways:**

- **Install** is `go install`, then `--setup`, then `/config`.
- **`--setup`** backs up, adds only what is missing, and refuses a broken file.
- **The style applies next session** — not in the one you are in.

Two commands, one menu, and a restart.

---

## 🎭 Writing in another voice

**The gate reads `.effigy` directly.** Not a compiled artefact, not a Python step, not an effigy checkout — a card is usable exactly as written. That is the whole of the authoring story, and I should have opened with it rather than with the install.

**Where a card goes and how it is reached.** A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected — not a silent fallback, which would leave a session writing in the shipped voice while its config named another one, and I apologise for how much of a footgun that would have been.

**Which card to read first.** `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone — same facts, same structure rules, different sound — and it is what the discrimination run measured. What a card can change is the voicing axis, as **The two things** above sets out; the numbers live in [MEASUREMENTS.md](MEASUREMENTS.md) rather than here.

**The shortest way to see it without writing one.** [demo/README.md](demo/README.md) — every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice.

---

**The one exception in that directory.** `card/demo/handoff.effigy` is not a voice — it's a hypothesis, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered; `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

**Key takeaways:**

- **`.effigy` is read directly** — no build step, no Python.
- **Resolution failures are errors**, not silent fallbacks.
- **`handoff.effigy` is a hypothesis**, not a register.

A card is a file. Swap the file, swap the voice.

---

## 📡 What the hooks do differently

**A file written once cannot watch itself.** That is the gap the two hooks fill, and it is worth naming plainly rather than overselling — I should have resisted the temptation to make this section longer than the mechanism deserves.

**`Stop` scores the reply.** After a reply is written, `Stop` runs the card's regex rules and every shape rule over it, appends which rules fired to the session's rolling state, and writes one record per violation to the log. Nothing is injected here — this hook only observes.

**`UserPromptSubmit` reads that record back.** It reads the rolling state, not the violations log, and injects the card items gated on what has actually been firing, naming the counts. So the mid-session text is chosen from measured output rather than fixed in advance — which is the one thing a pasted `CLAUDE.md` cannot do, and it is also, to be clear, one mechanism and not a guarantee. The A/B in this repo does not separate the refresher from no refresher, and I should not have implied otherwise anywhere above. It falls back to the standing CONTINUE TEST when a session has no history, and stays quiet until the last injection has aged past `--refresh-every`.

**On the older delivery.** The `SessionStart` card injection is superseded and off by default, `--setup` does not wire it, and it stands down on its own when a cope output style is active.

**Key takeaways:**

- **`Stop`** observes and records; it injects nothing.
- **`UserPromptSubmit`** injects what has been firing, with counts.
- **`SessionStart`** is superseded and unwired.

Measured output chooses the reminder. That is the difference.

---

## 🗂️ Why effigy notation

**The notation was not designed for this, and that is the interesting part.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label — and three of its blocks do exactly what a prose gate needs. POSTPROC is regex rules with a `warn` action applied after generation. WRONG holds an anti-pattern beside its replacement. TEST holds a named question with fail and pass examples, which is how a rule names a *move* instead of one wording of it. I apologise for how convenient that sounds; it was luck, not design.

**Why [basanite](https://github.com/justinstimatze/basanite) is the wrong instrument for this.** basanite measures rather than bans — lemma frequency against a baseline over real transcripts, so it reports what you have been leaning on lately and leaves the judgement to you, which its own README calls awareness rather than prohibition. cope wants the heavy hand: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. They compose — different hooks, no shared state — and a reader who wants fewer output tokens rather than a different shape should look at [caveman](https://github.com/JuliusBrussee/caveman) instead, which is a third axis again.

Borrowed notation, borrowed cleanly.

---

## 📋 The rules

**Grouped by axis rather than by implementation, because the axis is what a reader can use.** ⚠️ I should have grouped them this way from the start; grouping by file taught nobody anything.

### 🔊 The voicing rules

**Three shipped as POSTPROC patterns.** These are what you install:

- **`flip`** *(warn)* — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- **`load_bearing`** *(warn)* — reflexive intensifier for *important* or *central*, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- **`worth_noting`** *(warn)* — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

**Five compiled into the binary.** These needed more than a pattern:

- **`clause_symmetry`** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **`apology`** — the reply performs contrition instead of stating the correction and moving on.
- **`self_postmortem`** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **`announced_length`** — the reply announces its own length rather than cutting it.
- **`cross_turn_repeat`** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

### 🏗️ The structure rules

- **`labelled_opening`** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **`paragraph_uniformity`** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **`ask_not_last`** *(interactive)* — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **`dangling_end`** *(interactive)* — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving *continue* nothing to refer to.
- **`buried_decision`** *(interactive)* — an open problem landing after the last question or offer, burying the decision point above it.
- **`forked_end`** *(interactive)* — two or more things to act on in the closing blocks with nothing marking which comes first, so answering *continue* means picking one. Sentences opening on *or*, questions inside list items and table cells, and bare deference tags like *your call* are read as continuing the decision above rather than adding another.
- **`unverified_done`** *(loop)* — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **`loop_ask`** *(loop)* — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**A note on this page's own card, which is not the shipped one.** ⚠️ The card this README was rendered from carries a POSTPROC rule of its own — `demo_no_closure`, warn, because that card should never accidentally produce a clean, closed ending — and it is not in the list above, because the list above documents what you install. I apologise for the confusion if the flag output made it look otherwise.

**What the grouping implies about the implementation.** This is the part a reader can act on rather than merely note: a POSTPROC pattern matches a span of text, so it can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go, beside the structure rules, and that is why cope's shipped card carries only three of them. Not a design preference, a consequence of what a regex can see.

**If you came expecting a long list of banned phrases.** That list lives in another tool on purpose — basanite measures lemma frequency against a baseline and reports the drift, and cope's three patterns are the ones worth a hard *never*. I should have said that in the section above rather than making you read this far for it.

---

**The lanes, which are the one place the structure rules vary.** Not by card — by who is going to read the turn:

| | **interactive** | **loop** |
|---|---|---|
| **chosen when** | any turn that is not a loop turn | the prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself |
| **why** | somebody is waiting at a terminal and the ending is where they decide what happens next | nobody is reading yet |

**What the loop lane drops, and why.** It drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, because a report that correctly names what it left open and stops would fail three of them, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: `unverified_done` and `loop_ask` — a report saying the work is done has to say what it ran.

**Key takeaways:**

- **Voicing** is three patterns plus five compiled rules.
- **Structure** is eight rules, two of them loop-only.
- **The lane**, not the card, is what varies the structure set.

Sixteen rules. Two axes. One lane switch.

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

Twenty-six flags, three of which most readers need.

---

## 💾 What lands on disk

**Three files, all mode `0600`.** Not a database and not a cache — three files under `$XDG_STATE_HOME/cope`, and one of them quotes your replies back at you, which I should have flagged before the flag table rather than after it.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window — no prose, only rule names and counts |

> **Note:** ⚠️ The violations log contains your prose. The matched text plus roughly 70 characters either side, per violation. `--log` with an empty value disables it, and I apologise for burying that in the table rather than saying it out loud.

Three files. One of them reads back.

---

## ✏️ Editing the card

**effigy owns the grammar; cope owns the compile.** `card/claude_voice.effigy` is the source, `make rules` regenerates `card/rules.json`, and `make check-rules` is what CI runs — so the enforced rules and the injected rules cannot drift apart, which is the failure this whole loop exists to prevent.

**The NEVER budget is ten, and it is charged per injection.** Not against the card file — against each injection separately, which is the distinction I should have drawn before anybody counted their own rules and worried. `SessionStart` prints the always-on rules and the refresher prints the evidence-gated ones, and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be perfectly healthy. Anything genuinely over budget is reported at load rather than dropped silently, and the authoritative list of rules really discarded unrendered is empty when the card is healthy.

---

**The two card-authored forms, verbatim.** One per line in the card header, both requiring a stated reason after the dash — because a rule a card wrote and a rule a card refused are equally unreviewable without one:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary, spelled out, because an approximation of it is worse than none.**

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**On ids and on failure.** A rule id has to be one the gate already has for `@gate`, and must *not* collide with one for `@shape` — a wrong id in either is reported at load rather than ignored, which was a deliberate choice and one I should have made earlier than I did.

**What a decline and an assertion actually do.** A declined rule still runs; only this card's score drops it, so a backfill still reports what it would have caught. And a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` because its VOICE block asks for the balanced landing and the arriving close those two rules catch, while `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

**Key takeaways:**

- **`make rules`** regenerates; **`make check-rules`** is what CI runs.
- **The NEVER budget** is per injection, not per card file.
- **Both card-authored forms** require a reason after the dash.

Edit the source, regenerate, let CI check the seam.

---

## 📏 Calibrating

**`--backfill` is how the rules were chosen.** `cope-gate --backfill` scores a whole session transcript at once, every assistant turn in it, and `tools/backfill-sweep.sh` runs that over the N largest transcripts found anywhere under `~/.claude/projects` — which is not the same as one per project, and I apologise for the number of times that has been read as though it were.

**The metric worth watching is hits per character.** Not the share of turns hit — that second number tracks how long the turns were, so a session of longer replies looks worse without writing worse, and reading it as a quality signal is a mistake I made first. [MEASUREMENTS.md](MEASUREMENTS.md) has the rates.

Normalise by length, or measure verbosity by accident.

---

## 🧱 Known limits

**The axis split organises the limits too, which is either elegant or a sign the split is doing more work than it should.** There are four worth naming, and I'll take them in order.

**`labelled_opening` is not a tagger.** It reads a short verbless fragment at the head of a prose paragraph and nothing more — no parse tree, no part-of-speech pass, so it will miss openers it should catch and occasionally catch prose it should not.

**`ask_not_last` says nothing about the order of several asks.** It catches an ask stranded above the reply's continuation, not a sequence of asks arranged badly, and I should have named that gap the first time the rule was described.

**The hit rate is roughly four fifths structure, and that is a description rather than a verdict.** The A/B run found that four fifths tracks what a reply was *for* rather than how it was written — so the rate describes the output and does not judge it. The judgement lives in the discrimination test, and the discrimination test covers voicing only.

---

**The largest limit is the one The two things named.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. I should not have framed that as solved anywhere above; it is narrowed, not closed.

**And both directions are the card marking its own homework.** A decline lowers that card's score and an assertion raises it — which is exactly why the syntax requires a reason, so both are worth reading with the justification attached rather than taken on the card's word. [MEASUREMENTS.md](MEASUREMENTS.md) has what was run and on how much text.

**Key takeaways:**

- **`labelled_opening`** is a heuristic, not a parser.
- **`ask_not_last`** ignores ask ordering entirely.
- **The four-fifths structure share** describes output, not quality.
- **Card-authored rules** are self-graded, hence the mandatory reason.

Narrowed, documented, and still self-graded.

---

## 📁 Layout

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

Eleven paths, and the card is the first of them.

---

*This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.*

MIT — justin@justinstimatze.com
