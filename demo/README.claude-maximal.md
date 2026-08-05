# 🏺 cope

**What cope is.** Not a linter you configure before it does anything — an opinionated card that is on the moment you install it and is, for most readers, the entire product, checking two different jobs rather than one: how a reply *sounds* and how a reply is *shaped*, and only then, once you want to change it, a file you can edit or swap. The name is a foundry word — a cope is the upper half of a mould, the half carrying the shape being cast into — and I should say up front that I will lean on that metaphor once more before the page is done, which I apologise for in advance.

Read [demo/README.md](demo/README.md) first: every file in that directory is this same README written again from a different card, same prompt and same facts, so the card is the only thing that changed — and reading two of them against each other shows what a card does faster than the rest of this page explains it, including [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card instructing every tic this model is measured to have, which is deliberately hard going.

---

## 🎭 The two things it checks

**TL;DR for this section.** There are two axes, they fail independently, and only one of them lives in the card — which is the frame every later section is written against, and I should have said that before writing three thousand words underneath it.

**Voicing is the first.** Not what the reply decides — what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely, across `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES` and `POSTPROC`, which means swapping the card swaps every word of it. A sentence reaching for the balanced two-beat — two clauses of near-equal length repeating a content word across the joint — is a voicing problem, and it is the half with a measured result behind it. That distinction is load-bearing, and I apologise for stating it twice.

**Structure is the second.** Not how the reply reads — the shape of it as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. A reply that names an open problem in its closing paragraph and then stops is a structural failure, however well the sentences land. This half is compiled into the binary, in `internal/scan`, so it is the same whichever card is loaded — and the same reply can be clean on one axis and bad on the other, which is the whole reason to keep them apart.

---

**How far a card reaches in.** Not nowhere, which is where it used to reach, and not everywhere either — a card can now move in two directions across that seam, and I should have been clearer about this earlier because it is the part that most recently changed. A card declines a built-in rule with `@gate <rule_id> off — <why>`, one per line in the card header, which exists because a card whose `VOICE` block asked for something a built-in rule catches was being marked down for obeying itself: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason. A card asserts a structural rule of its own with `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold` because its peak asks the reader to re-enter cold and read the last block only.

**Where that reach stops.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — it cannot express what the compiled rules express, `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. That is a boundary and not a roadmap, and I apologise for the flatness of saying so.

> **Key insight:** One axis is the card's to change, the other is the binary's — and the gap between them is where every honest complaint about this tool lives.

**Key takeaways:**

- **Voicing** is what the sentences sound like, it lives in the card, and swapping the card swaps it.
- **Structure** is where the decision sits and how the reply ends, and it is the same under every card.
- **A card** can decline a compiled rule and assert a shape rule of its own, with a reason, and no more than that.

Two axes, one card, one binary.

---

## 🧩 The problem

**Let me break this down.** ⚠️ There are three things going on here, in this order — where an instruction sits, what an instruction can say, and the failure no instruction reaches — and I apologise, because the first is the largest and most readers arrive holding it backwards.

**First, where the instruction sits.** A reader who has edited a global `CLAUDE.md` and watched it not stick almost certainly believes that file is the system prompt. It is not. It arrives as one message attached to the first turn, and the conversation buries it under everything written after — which makes the delivery mechanism, not the wording, the substrate the whole problem sits on. An output style is in the system prompt itself, which the harness re-reminds the model of as the conversation runs, and moving one card between those two places without changing a word of it is most of why cope works at all. That was measured; the run is in [MEASUREMENTS.md](MEASUREMENTS.md), and I have deliberately put no counts on this page, which I should have flagged before you went looking for them.

**Second, what the instruction can say.** Great question to be holding here — and the answer is that instruction alone does not fix the phrasing. A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn, and the flip still appeared twice in the session that built this tool, while the ban was the topic under discussion. Think of it like patching a leak by writing *no leaks* on the pipe — the analogy does break down, of course, because a pipe does not read the note and route around it, and this does: naming a surface form pushes the move into a variant. That is the voicing side, and it is genuinely the smaller half.

**Third, the failure no instruction reaches.** ⚠️ This one is not a phrasing problem at all, and I apologise for having spent two paragraphs implying that everything here is. An ending that leaves the reader nothing to answer costs a whole round trip — a real, countable one — and no instruction could have banned it, because it is not a habit of wording, it is a habit of shape. Different complaint, different cause, different substrate.

---

**What the claim rests on.** The flip is an anecdote about one rule, and it would be dishonest to let it carry more weight than that — thank you for reading this far before I said so. What the claim actually rests on is the blind discrimination test: a reader is shown only a voice's own description of itself and picks which of two replies was written under it. The rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md), and the blind preference runs are deliberately not cited here, because both sides of those were written under a card — they compare two ways of writing well and cannot see a voice being swapped.

**Key takeaways:**

- **Placement** beats wording: one message on turn zero is not the system prompt.
- **Wording** alone pushes a habit into a variant rather than removing it.
- **Shape** is a separate failure with a separate cost, and instruction never touched it.

The instruction was fine. The place it sat was not.

---

## 📦 Install

**Two commands and one menu choice.** Not a configuration exercise — three steps, in the order somebody actually does them, and I apologise for having buried this section under two of frame.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

**What those two commands did.** `cope-gate --setup` does the whole install — it emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. Worth noting what it does to a file you care about: it backs the settings file up first, it adds only what is missing so a second run changes nothing, it leaves every other key alone including other tools' hooks on the same events, and it refuses to touch a settings file that does not parse — and `cope-gate --setup --dry-run` prints what it would change and writes nothing, which I should have offered before describing the writing version. If you would rather nothing wrote to your settings at all, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` by hand, and `COPE_CARD=<name>` in front of it emits a different one.

**Then the menu choice.** ✅ Pick the style under `/config` -> Output style, and the entry to look for is named **`claude_voice`** — the shipped card's id, not the word cope. That name is doing a lot of work in this paragraph: a reader who has just installed something called cope will scan that menu for the word cope, fail to find it, select nothing, and conclude the tool does not work. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear` — and I apologise for how many readers have watched a running conversation not change and assumed a bug.

**Why the card goes there.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through a hook — that is the reason and not the instruction, so I will leave it at one sentence.

---

**Then the hooks.** The card no longer arrives through one, which is the thing to hold onto here, and I should have said it before showing you a block of JSON. This goes in the `hooks` block of `~/.claude/settings.json`:

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

**What the two hooks buy.** `Stop` gets the reply scored after it is written — nothing else, just the score and the record of it. `UserPromptSubmit` gets the rules that have actually been firing restated mid-session, which a file written once cannot do. The voice works without either of them; these are the measurement half, and I apologise for the temptation to oversell a pair of hooks that measure rather than act.

**Two practical notes.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead, and that is the failure I have seen most. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in; `--inject` remains as the superseded delivery for anyone who wants the old turn-zero behaviour, and it stands down on its own when a cope output style is active.

**Key takeaways:**

- **Run** the two commands, then select `claude_voice` in `/config`.
- **`--setup`** backs up, adds only what is missing, and leaves other hooks alone.
- **The hooks** buy scoring and mid-session restatement, not the voice itself.

Install, select, restart. That is the whole ceremony.

---

## 🗣️ Writing in another voice

**This is the capability the first line promised.** Not a rebuild, and not a Python toolchain — the gate reads `.effigy` directly, so a card is usable exactly as written and needs no effigy checkout, which I should have led with rather than making you read the install first. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected — deliberately, because the alternative is a session writing in the shipped voice while its config names another one, and that is the failure mode nobody notices.

**Which card to read first.** `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and is the one the discrimination run measured — a controlled swap rather than an interesting one, which makes it boring to read and useful to read. What a card can change is the voicing half described where this page laid out the two axes, and I will not re-argue that here; the numbers are in [MEASUREMENTS.md](MEASUREMENTS.md).

**The shortest route of all.** [demo/README.md](demo/README.md) is this README written again from a different card, once per card, same prompt and same facts, so the only thing varying between them is the voice — I am linking the directory index rather than listing the cards, because that page is generated from whatever was last built and this one should not carry a stale second copy.

**One exception in that directory.** `card/demo/handoff.effigy` is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered as a page. `make cards` installs it with the rest, so a reader listing their cards will find it sitting there looking like a register, and I apologise for the confusion that will cause before this sentence reaches them.

**Key takeaways:**

- **A card is a file**, read as written, reached by `--card`, `COPE_CARD` or `--rules`.
- **`card/demo/lecturer.effigy`** is the controlled comparison and the one that was measured.
- **`card/demo/handoff.effigy`** is a hypothesis, not a voice, and should not be worn.

Swap the file, swap the sentences.

---

## 🪝 What the hooks do differently

**The mechanics live here.** Not in the install section, which carried one sentence each on purpose — and if the two ever read as the same two paragraphs twice, that is my error and worth reporting.

**`Stop`, first.** `cope-gate` scores the reply just written, appends which rules fired to the session's rolling state, and writes one record per violation to the log. That is all it does — no prose is held in the state file, only rule names and counts, and I should have said so before the section on disk rather than after.

**`UserPromptSubmit`, second.** `cope-gate --refresher` reads the rolling state — the state, not the violations log, and that distinction is load-bearing — and injects the card items gated on what has actually been firing, naming the counts back. It falls back to the standing `CONTINUE TEST` when the session has no history yet, and it stays quiet until the last injection has aged past `--refresh-every`. So the mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` genuinely cannot do.

**And a caveat I owe you.** ⚠️ That is one mechanism and not a guarantee — the A/B run in this repo does not separate the refresher from no refresher, so the claim above is architectural rather than measured, and I should have attached that qualification to the sentence itself rather than to the paragraph after it.

Scored after, reminded before, quiet in between.

---

## 🧪 Why effigy notation

**Because three of its blocks already did the job.** Not a format chosen for elegance — [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here entirely off-label, and I apologise for how odd that sounds before the reason arrives. `POSTPROC` is regex rules with a `warn` action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples — which is how a rule names a *move* rather than one wording of it.

**And why basanite is the wrong instrument here.** [basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round: cope bans, so a rule fires or it does not and the register is fixed the moment you pick a card, while basanite measures lemma frequency against a baseline over real transcripts and leaves the judgement to you. Which one fits is a question about mood more than about correctness — a heavy hand when a habit is annoying you today, a moving measurement when you would rather watch the drift than legislate it — and they compose, different hooks, no shared state. Worth naming a third: [caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies to cut output tokens, so a reader wanting *fewer* tokens rather than different structure should go there and I should have said that sooner.

A notation for imaginary people, aimed at a real one.

---

## 📏 The rules

**Grouped by axis, not by implementation.** That grouping is the interesting part, and I apologise for making you read to the bottom of the section to find out why.

**Voicing rules — the `POSTPROC` patterns in the shipped card.** Three, and these are what you install:

- **`flip`** (`warn`) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- **`load_bearing`** (`warn`) — reflexive intensifier for *important* or *central*, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- **`worth_noting`** (`warn`) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

**Voicing rules that needed more than a pattern.** These are compiled, and they sit in `internal/scan` beside the structure rules:

- **`clause_symmetry`** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **`apology`** — the reply performs contrition instead of stating the correction and moving on.
- **`self_postmortem`** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **`announced_length`** — the reply announces its own length rather than cutting it.
- **`cross_turn_repeat`** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure rules, all compiled, all the same under every card.**

- **`labelled_opening`** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **`paragraph_uniformity`** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **`ask_not_last`** *(interactive)* — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **`dangling_end`** *(interactive)* — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **`buried_decision`** *(interactive)* — an open problem landing after the last question or offer, burying the decision point above it.
- **`forked_end`** *(interactive)* — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **`unverified_done`** *(loop)* — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **`loop_ask`** *(loop)* — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies.** Here is the part a reader can actually use, and I should have opened with it rather than a list. A `POSTPROC` pattern matches a span of text, which means it can only ever describe *wording* — so every voicing rule that needed more than a pattern had to be written in Go, beside the structure rules, on the wrong side of the seam from the card that wants it. That is why cope's shipped card carries only three `POSTPROC` rules, which is a small number for a tool people expect to be a phrase blacklist. This page happens to be written under `card/demo/claude_maximal.effigy`, which carries one of its own — `demo_no_closure`, warning when this card accidentally produces a clean, closed ending — and that one is that card's, not yours.

**If you came for a long banned-phrase list.** ⚠️ It is not here, and that is deliberate rather than unfinished — basanite measures lemma frequency against a baseline and reports what you have actually been leaning on lately, which its own README calls awareness rather than prohibition. cope is the blunt instrument by design; I apologise if that is the wrong instrument for what you arrived wanting.

---

**Lanes, which is where structure does vary.** Not by card — by who is going to read the turn, which is the one dimension the compiled rules bend along. The interactive lane is any turn that is not a loop turn, chosen because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself.

| | **interactive** | **loop** |
|---|---|---|
| **dropped** | — | ⚠️ `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision` |
| **added** | — | ✅ `unverified_done`, `loop_ask` |

**Why that trade.** Nobody is reading yet — a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran. I should have made that consequence explicit in the table rather than in the paragraph under it.

**Key takeaways:**

- **Three `POSTPROC` rules** ship, because a pattern can only describe wording.
- **Five compiled voicing rules** exist because those moves outran the patterns.
- **Eight structure rules** run the same under every card, varying only by lane.

Wording in the card, shape in the binary, always.

---

## 🚩 Flags

**The full surface, verbatim.** Twenty-six of them, and I apologise for the length of what follows — the alternative was a curated subset, which is how a reader ends up not knowing a flag exists.

| flag | default | what it does |
|---|---|---|
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
| `--display` | `false` | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a SessionStart hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default `~/.claude/output-styles`) |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

Everything the binary answers to, in one place.

---

## 💾 What lands on disk

**Three files, all `0600`.** Not a database and not a cache — three small files under `$XDG_STATE_HOME/cope`, and I should state the uncomfortable one first rather than third.

- `$XDG_STATE_HOME/cope/violations.jsonl` — one JSON record per violation, **carrying the matched text and about 70 characters either side**. That is your prose, quoted back to a file on your disk; `--log` with an empty value disables it.
- `$XDG_STATE_HOME/cope/refresher-<session-id>` — an empty file whose mtime is the refresher clock.
- `$XDG_STATE_HOME/cope/session-<session-id>.json` — the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts.

**Worth being blunt about the first one.** The log quotes replies back — that is not a side effect, it is the point of the log, since a violation without its surrounding text is unreviewable — and I apologise for how easily that fact hides inside a bullet.

Two files hold counts. One holds your sentences.

---

## ✍️ Editing the card

**Where the grammar lives.** Not here — effigy owns the `.effigy` grammar, and this repo owns only what it does with one, which I should have made explicit in the section about notation rather than saving it for now. `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart.

**The `NEVER` budget.** Ten per injection — and the word *per injection* is doing all the work in that sentence. The budget is charged against each injection separately and not against the card file, since the always-on rules and the evidence-gated ones are printed by different paths and no code path renders their union, so a card may hold more `NEVER` rules in total than the budget and still be perfectly healthy. Anything genuinely over budget is reported at load rather than dropped silently, which is the behaviour I should have described before the arithmetic.

---

**Declining a rule.** `@gate <rule_id> off — <why>`, one per line in the card header. The id has to be one the gate actually has, and a wrong id is reported at load rather than ignored — I apologise for how many tools do the opposite. The reason after the dash is required, because a rule a card refused is unreviewable without one, and note that a declined rule still runs: only this card's score drops it, so a backfill still reports what it would have caught.

**Asserting a rule.** `@shape <id>: <selector> <predicate> — <why>`, one per line in the card header, with the id required *not* to collide with one the gate already has. The vocabulary is exactly this and nothing more:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**Why the reason is mandatory here too.** A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies — so the text after the dash *is* the error message a reader will see, and writing it badly is a cost you pay later rather than now. `card/demo/handoff.effigy` is the worked example: `readable_cold` asserts `last paragraph words <= 60`, and the 60 is measured — across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90. I should have led with a measured number rather than a syntax box.

**Key takeaways:**

- **`make rules`** regenerates, **`make check-rules`** stops the two halves drifting.
- **The `NEVER` budget** of 10 is per injection, not per file.
- **Both card-authored forms** require a reason, and a bad id fails loudly at load.

Write the rule, write the reason, run `make check-rules`.

---

## 📐 Calibrating

**How the rules were actually chosen.** Not by taste, or not only by taste — `cope-gate --backfill` scores a whole session transcript at once, and that is the tool every rule on this page was argued into existence with. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project, and I apologise for how easy that is to misread as coverage.

**The metric worth watching.** Hits per character, not the share of turns hit — that second number tracks how long the turns were, which means it moves when nothing about the writing has changed. Think of it like judging a book by how many pages have a typo on them: the count rises with the page count and tells you about the printer rather than the proofreader, and the analogy breaks down only in that a transcript has no printer. The rates themselves are in [MEASUREMENTS.md](MEASUREMENTS.md) rather than here, deliberately.

Score the transcript, watch the density.

---

## ⚠️ Known limits

**Organised by axis, like everything else.** There are three worth stating plainly, and I should have collected them nearer the top rather than behind the flag table.

**The narrow ones first.** `labelled_opening` is not a tagger — it approximates, and it will misjudge an opener now and then. `ask_not_last` says nothing about the *order* of several asks, only about one sitting in the wrong place. Neither of those is a design position; both are simply what the implementation does, and I apologise for having described them earlier in language confident enough to obscure it.

**The hit rate is a description, not a verdict.** Roughly four fifths of hits are structural, and the A/B run found that four fifths tracks *what a reply was for* rather than how it was written — so it describes the output and does not judge it. The judgement lives in the blind discrimination test, and that test covers voicing only, which is a real asymmetry and not a modest one.

**The largest limit is the seam this page opened on.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. Worth adding the uncomfortable half: both directions are the card marking its own homework, since a decline lowers that card's score and an assertion raises it, which is exactly why the syntax requires a reason attached to each. Read them with the reason, or do not read them; the run is in [MEASUREMENTS.md](MEASUREMENTS.md).

**Key takeaways:**

- **`labelled_opening`** approximates and **`ask_not_last`** ignores ordering.
- **The hit rate** describes what a reply was for, not how well it was written.
- **A card marks its own homework** in both directions, which the mandatory reason exists to expose.

The seam is narrower than it was. It is still a seam.

---

## 🗂️ Layout

**Eleven paths worth knowing.** Not the whole tree — the parts a reader opens, and I apologise for having described several of them already without saying where they were.

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

Two cards, one binary, one measurements file.

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT — justin@justinstimatze.com
