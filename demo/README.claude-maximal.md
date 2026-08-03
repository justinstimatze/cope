# cope

**cope** ships an opinionated card that is on the moment you install it and is the whole product for most readers — it scores two different jobs rather than one, how a reply *sounds* and how a reply is *shaped*, and since a cope is the upper half of a foundry mould, the half carrying the shape being cast into, that card is a file on disk you can edit or swap for another.

Start at [demo/README.md](demo/README.md): every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed — and reading two of them against each other shows what a card does faster than the rest of this page explains it, with [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have, making the point in a single glance.

---

## 🎭 The two things cope checks

**TL;DR — there are two axes, and keeping them apart is the whole frame.** I apologise for opening on a frame rather than on an instruction, but the sections below are unreadable without it. Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. Structure is the shape of the reply as a thing the reader has to *use* — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it.

**Where each one lives.** The real question is not whether both are enforced — they are — it's *who owns them*, and the ownership is load-bearing. Voicing lives in the card, entirely: VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC, which means swapping the card swaps every word of it, and that is the half with a measured result behind it. Structure is compiled into the binary, in `internal/scan`, so it is the same whichever card is loaded and varies only by lane.

**One concrete instance of each, because abstraction here is my mistake to avoid.** A sentence reaching for the balanced two-beat — two comma-joined clauses of near-equal length repeating a content word across the joint — is a voicing problem. A reply that names an open problem in its last paragraph and then stops, leaving "continue" nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, which is exactly why they are not one number.

---

**How far a card reaches into the structure half.** This is the part that most recently changed, and I should have said so earlier rather than burying it here. Two directions, both written in the card header, one per line:

- **A card declines a built-in rule** it disagrees with: `@gate <rule_id> off — <why>`.
  - `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch.
    - Put differently: a card whose VOICE asked for something a built-in rule catches was being marked down for obeying itself.
- **A card states a structural rule of its own** that the gate then checks: `@shape <id>: <selector> <predicate> — <why>`.
  - `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.
    - In other words: a card's own commitment about how a reply ends previously had nowhere to be checked.

**The boundary, stated without apology.** ⚠️ The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — it cannot express what the compiled rules express, `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it. No rate appears in this section on purpose; [MEASUREMENTS.md](MEASUREMENTS.md) has the runs and the reasons its numbers do not carry more than that.

> **Key insight:** Voicing is the card's, structure is the binary's, and the card now gets a vote in the second — a vote with a reason attached.

**Key takeaways**

- **Voicing** — the sentences. In the card. Swap the card, swap the sound.
- **Structure** — the shape. In the binary. Same rules whichever card is loaded.
- **The card's vote** — `@gate` declines, `@shape` asserts, and both require a written reason.

Two axes, two owners, one vote across the seam.

---

## 🧨 The problem

**Instruction alone does not fix the phrasing.** Great question to ask first — and I apologise, you did not ask it, I asked it for you. A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn, and the flip still appeared twice in the session that built this, while the ban was the topic of that session. Naming a surface form does not remove the move; it pushes the move into a variant.

**The structural failure is a different complaint with a different cause.** An ending that leaves the reader nothing to answer costs a whole round trip — not a phrasing habit, not something an instruction could have banned, and not a matter of taste. Think of it like a form that asks for a signature but prints no line: the words are all correct and the reader still has to come back. The analogy breaks down, of course, because a form cannot be asked to reprint itself, and a reply can.

**What the claim actually rests on.** ⚠️ To be clear, the flip above is an anecdote about one rule, and I should not have led with an anecdote at all. The evidence is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. The rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md), and the blind *preference* runs are deliberately not cited here, since both of their arms were written under a card — they compare two ways of writing well and cannot see a voice being swapped.

Instruction is read; behaviour is measured.

---

## 📦 Install

**Three steps, in the order somebody actually does them.** I apologise for announcing the count before giving it, but the order matters more here than anywhere else on the page — two commands and one menu choice.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> pick the cope card
```

**What `--setup` did, rather than what it saved you from.** It emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left — and a reader letting a tool edit their settings deserves the shape of it before running it, not after, which is on me to state up front. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing.

**The by-hand route, for anyone who would rather not have their settings written to.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one — then pick it under `/config -> Output style`. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

**When it applies, which is the thing that looks like a bug and is not.** A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear` — I apologise for the flat delivery, but a reader who picks a style, sees no change in the running conversation, and concludes the tool is broken has been failed by the docs and not by the tool.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through the hook.

---

**Now the hooks — and be clear that the card no longer arrives through one.** This is the hooks block of `~/.claude/settings.json`:

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

**What the two remaining hooks buy.** `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log — and `UserPromptSubmit` reads that rolling state, not the log, and injects the card items gated on what has actually been firing, naming the counts. ✅ The voice works without either of them; these two are the measurement half, and I should have said that before the code block rather than after it.

**Two operational notes, which are boring and load-bearing.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in — the superseded `--inject` delivery remains for anyone who wants the old turn-zero message and stands down on its own when a cope output style is active.

Install, select, and let `Stop` start keeping score.

---

## 🗣️ Writing in another voice

**The gate reads `.effigy` directly.** Great point to raise early — and again, mine, not yours, for which I apologise. A card is usable exactly as written: no Python, no effigy checkout, no compile step in the loop. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set.

**A name that resolves to nothing is an error, and nothing is injected.** Not a silent fallback — an error, because the failure mode worth designing against is a session writing in the shipped voice while its config names another one. To put a finer point on it: the loud failure is the safe one, and I should have made that the first sentence of the paragraph.

**Which card to read first.** Read `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and is what the discrimination run measured — what a card can change is the voicing side described under *The two things cope checks* above, and I will not re-argue that section here beyond noting it. The numbers stay in [MEASUREMENTS.md](MEASUREMENTS.md).

**The shortest way to see what a card does without writing one.** [demo/README.md](demo/README.md) — every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice, and that index is generated from whatever was last built rather than copied here.

**One exception in that directory, flagged so you are not surprised by it.** ⚠️ `card/demo/handoff.effigy` is a hypothesis rather than a voice: it keeps the shipped card's handoff rules, drops everything about prose, and is meant to be run through `make pairs` against the full card rather than rendered — `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

**Key takeaways**

- **`.effigy` is read as written** — no toolchain, no checkout.
- **Resolution is strict** — an unknown card name is an error, never a quiet fallback.
- **`handoff.effigy` is not a register** — it is an experiment that happens to install.

Pick a card, or write one; the notation is the interface.

---

## 🪝 What the hooks do that a pasted file cannot

**`Stop` keeps the record.** It scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log — and I apologise for restating that from the install section, though the restatement is doing work here, because the next paragraph depends on it.

**`UserPromptSubmit` spends that record.** It reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. That is the one thing a pasted `CLAUDE.md` cannot do. When the session has no history yet it falls back to the standing CONTINUE TEST, and it stays quiet until the last injection has aged past `--refresh-every`.

**Stated plainly, and undersold on purpose.** ⚠️ It is one mechanism, not a guarantee — the A/B in the repo does not separate the refresher from no refresher, so treat this as a design claim rather than a measured one, and I should have led with the caveat rather than closing on it. `SessionStart --inject` remains available, superseded and off by default.

Measured output chooses the reminder; a file cannot.

---

## 🧩 Why effigy notation

**The notation is borrowed off-label, and the borrowing is the interesting part.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, and three of its blocks do exactly what a prose gate needs — POSTPROC is regex rules with a `warn` action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a *move* instead of one wording of it. I apologise for the density of that sentence; the three-block fit is genuinely the whole reason.

**Why basanite is the wrong instrument for this job, and the right one for a neighbouring job.** [basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round: cope bans — a rule fires or it does not, the card says never, and the register is fixed the moment you pick it — while basanite measures lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it — they compose, different hooks and no shared state, and running both is reasonable. Separately, [caveman](https://github.com/JuliusBrussee/caveman) is a third axis again, by a different author, compressing agent replies to cut output tokens — worth naming because a reader wanting *fewer* tokens rather than different structure should go there instead.

Borrowed grammar, three blocks, one fit.

---

## 📋 The rules, grouped by what they are for

**Grouped by axis rather than by implementation, because the axis is what a reader can act on.** I apologise for the second table of contents in one page, but the grouping carries information the file layout does not.

**Voicing rules — what the sentences sound like:**

- **`demo_no_closure`** *(warn)* — this card should never accidentally produce a clean, closed ending.
- **`clause_symmetry`** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **`apology`** — the reply performs contrition instead of stating the correction and moving on.
- **`self_postmortem`** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **`announced_length`** — the reply announces its own length rather than cutting it.
- **`cross_turn_repeat`** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure rules — the shape of the reply as a thing to use:**

- **`labelled_opening`** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **`paragraph_uniformity`** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **`ask_not_last`** *(interactive)* — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **`dangling_end`** *(interactive)* — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **`buried_decision`** *(interactive)* — an open problem landing after the last question or offer, burying the decision point above it.
- **`forked_end`** *(interactive)* — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **`unverified_done`** *(loop)* — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **`loop_ask`** *(loop)* — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies about the implementation.** A POSTPROC pattern matches a span of text, so it can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go beside the structure rules, and that is why the shipped card carries exactly **1** POSTPROC rule. To be clear, the split is not tidy and I am not going to pretend it is: the axes are clean, the files are not.

**If you came here for a long list of banned phrases, that list lives in another tool on purpose.** ⚠️ cope bans structure and register; [basanite](https://github.com/justinstimatze/basanite) measures vocabulary against a baseline, and the phrase-frequency work belongs there — I should have said that in the section above rather than making you read this far for it.

**Lanes, the one place the structure rules do vary — not by card, by who is going to read the turn.** 🎯

| | **interactive** | **loop** |
|---|---|---|
| **chosen when** | any turn that is not a loop turn | the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself |
| **why** | somebody is waiting at a terminal and the ending is where they decide what happens next | nobody is reading yet |
| **dropped** | — | `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision` |
| **added** | — | `unverified_done`, `loop_ask` |

**Why the loop lane drops four rules and adds two.** A report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself — so what replaces them is the claim check: a report saying the work is done has to say what it ran. That seam is load-bearing, and I apologise for not naming it before the table.

**Key takeaways**

- **Voicing rules** — six, mostly compiled, one from the card's POSTPROC block.
- **Structure rules** — eight, all compiled, four of them lane-specific.
- **Lanes** — the audience decides which rules apply, not the card.

Same axes everywhere; different reader, different rules.

---

## 🎚️ Flags

| flag | default | what it does |
|---|---|---|
| `--ab` | `false` | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | *(empty)* | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | *(empty)* | read a turn log and report the arms; - reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | *(empty)* | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | *(empty)* | name of an installed card to write in, from the cards directory; also $COPE_CARD |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | *(empty)* | score a prose file against the card and exit; - reads stdin |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as --display would rewrite it |
| `--dry-run` | `false` | with --setup, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to ~/.claude/output-styles as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | *(empty)* | directory to write the output style into (default ~/.claude/output-styles) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | *(empty)* | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | *(empty)* | comma-separated rule ids to render -render-arm against |
| `--render-lane` | *(empty)* | render -render-arm as the given lane sees it: interactive (default) or loop |
| `--rules` | *(empty)* | read the card from this .effigy or .json file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

Twenty-six flags, one of which you will ever type twice.

---

## 💾 What lands on disk

**Three files, all mode `0600`, all under `$XDG_STATE_HOME/cope` — and one of them quotes you back to you.** I apologise for burying that in a heading rather than a warning: the violations log stores the matched text.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window — no prose is stored, only rule names and counts |

> **Important:** ⚠️ The log quotes replies back. The rolling session record does not — it holds rule names and counts and nothing else — and `--log` set to empty disables the log entirely.

Two files of counters, one file of quotations.

---

## ✍️ Editing the card

**effigy owns the grammar; cope owns the compile.** `card/claude_voice.effigy` is the shipped card, `make rules` regenerates `card/rules.json` from it, and `make check-rules` is what CI runs so the enforced and injected rules cannot drift — and I should have said "cannot drift" more carefully, since what CI actually guarantees is that the drift is caught rather than prevented.

**The NEVER budget is 10, and anything over it is reported at load rather than dropped silently.** ⚠️ The budget is charged against each *injection* separately, not against the card file: `SessionStart` prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more NEVER rules in total than the budget and still be perfectly healthy. `facts.never_rules_over_budget` is the authoritative list of rules that really are discarded unrendered, and it is empty when the card is healthy.

---

**Both card-authored forms, one per line in the card header:**

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary, exactly as it stands:**

- **selectors:** `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- **predicates:** `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**The three constraints on ids and reasons.** A rule id has to be one the gate already has for `@gate`, and must *not* collide with an existing one for `@shape` — a wrong id is reported at load rather than ignored, which is the loud-failure choice again. The reason after the dash is required in both forms, because a rule a card wrote and a rule a card refused are equally unreviewable without one, and I apologise for making that sound like a policy when it is really an admission that neither is trustworthy on its own.

**What a decline does and does not do.** ✅ A declined rule still runs; only this card's score drops it — so a backfill still reports what it would have caught, and nothing is hidden. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

**Key takeaways**

- **Regenerate with `make rules`; verify with `make check-rules`.**
- **The NEVER budget is per injection, not per file, and overflow is reported at load.**
- **Both `@gate` and `@shape` require a written reason, and a bad id fails loudly.**

The card is a file; treat it like source.

---

## 📐 Calibrating

**`cope-gate --backfill` scores a whole session transcript at once, and is how the rules were chosen.** Not a reporting convenience — the selection instrument. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project, and I should have made that distinction louder the first time it mattered.

**Watch hits-per-character rather than the share of turns hit.** The second number tracks how long the turns were, which makes it a length metric wearing a quality metric's clothes — think of it like judging a book by how often a word appears per page rather than per book, except the analogy breaks down, because pages are a fixed size and turns are not. The rates themselves are in [MEASUREMENTS.md](MEASUREMENTS.md), quoted nowhere on this page on purpose.

Per character, not per turn.

---

## 🚧 Known limits

**The axis split is what organises the limits, so take them in that order.** I apologise for another announced structure; it is the shortest honest way through four unrelated caveats.

**Voicing limits.** ⚠️ `labelled_opening` is not a tagger — it is a heuristic over paragraph openings, and it will be wrong about some of them. `ask_not_last` says nothing about the ordering of several asks; it only notices that one sits above continuing prose.

**Structural limits, and the number that looks like a verdict and is not.** The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks *what a reply was for* rather than how it was written — so it is a description of the output and not a judgement of it. The judgement lives in the discrimination test, and that test covers voicing only.

**The largest limit is the one named under *The two things cope checks*, stated as it now stands.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

**Both directions are the card marking its own homework.** A decline lowers that card's score and an assertion raises it; both are worth reading with the reason attached, and that is precisely why the syntax requires one. Everything quantitative sits in [MEASUREMENTS.md](MEASUREMENTS.md).

**Key takeaways**

- **`labelled_opening` is a heuristic**, not a parser.
- **Four fifths structure describes the corpus**, not the quality of it.
- **A card grades itself in both directions**, which is why reasons are mandatory.

The limits are documented because the numbers alone would flatter the tool.

---

## 🗂️ Layout

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
