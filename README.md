# 🏭 cope

**The short version.** cope ships one opinionated card that is on the moment it is installed and is the whole product for most readers, it checks two different jobs rather than one — how a reply sounds and how a reply is shaped — and only then, because the card is just a file, it can be edited or swapped for another: a cope is the upper half of a foundry mould, the half carrying the shape being cast into.

**Read the renders first.** Every file under [demo/README.md](demo/README.md) is this same page written again from a different card, same prompt and same facts, so the card is the only thing that changed — and reading two of them against each other shows what a card does faster than the rest of this page explains it; if this one is hard going, [demo/README.claude-voice.md](demo/README.claude-voice.md) is the same page in the register cope actually ships.

---

## 🎭 The two things cope checks

**TL;DR — there are two axes, and they are not the same axis.** Before anything else, and I apologise for front-loading a frame rather than a command, the frame is load-bearing enough that skipping it makes every later section read as one undifferentiated pile of rules.

**First, voicing.** Voicing is what the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely: VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC. Swap the card and you swap every word of it. Worth noting that this is the half with a measured result behind it, and I should have said so earlier rather than making you wait for it.

**Second, structure.** Structure is not how the reply reads, it is the shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. It is compiled into the binary, so it is the same whichever card is loaded, varying only by lane. That distinction is doing a lot of work, and I apologise for restating it, but a reader who conflates the two will expect a card swap to change the endings and it will not.

**A concrete instance of each, because abstraction here is my fault.** A sentence that reaches for the balanced two-beat — two clauses of matched length, the same content word landing on both sides of the comma — is a voicing problem. A reply that names an open problem in its last paragraph and then simply stops, leaving "continue" nothing to refer to, is a structural one. The same reply can be clean on one axis and bad on the other, which is the whole reason to keep them apart.

> **Key insight:** Voicing is what the card says. Structure is what the binary knows. The seam between them is where a card's reach ends.

---

**How far a card reaches into the structure half.** Not nowhere — and this is the part that most recently changed, so an older description of it is now wrong and I should have flagged that sooner. A card reaches in two directions, both declared in its header, both requiring a reason after the dash:

| | syntax | why it exists |
|---|---|---|
| **decline a rule** | `@gate <rule_id> off — <why>` | a card whose VOICE asked for something a built-in rule catches was being marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch |
| **assert a rule** | `@shape <id>: <selector> <predicate> — <why>` | a card's own commitment about how a reply ends had nowhere to be checked — `card/demo/handoff.effigy` asserts `readable_cold` (last paragraph words <= 60) because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way |

**The boundary, stated flatly.** ⚠️ The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — it cannot express what the compiled rules express, `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it. No rate appears in this section on purpose; MEASUREMENTS.md has the run and the reasons its numbers do not carry more than that. I apologise for the deferral, which is the honest answer rather than the satisfying one.

- **Key takeaways.**
  - Voicing lives in the card, and swapping the card swaps it.
  - Structure lives in the binary, and is the same for every card.
  - A card can decline a built-in rule, with a reason.
  - A card can assert one of its own, in a vocabulary of words, sentences and asks.

Two axes. One card. One binary.

---

## 🧩 Why instruction alone does not fix this

**Instruction does not fix the phrasing, and I should have believed that earlier than I did.** A global CLAUDE.md banning the "not A, it's B" flip is read every single turn — and the flip still appeared twice in the session that built this, while the ban was the topic of that session. Naming a surface form does not remove the move; it pushes the move into a variant. That is the voicing side of the complaint, and it is an anecdote about one rule rather than the argument.

**The structural side is a different complaint with a different cause.** An ending that leaves the reader nothing to answer is not a phrasing habit an instruction could have banned — it costs a whole round trip, because the reader has to write a turn whose only content is telling the assistant which of the things it just named to do. Think of it like a form that asks for your signature and then does not say where: the ink is fine, the field is missing. The analogy breaks down, of course, in that a form cannot be rewritten mid-conversation and a reply can.

**What the claim actually rests on.** Great question to be asking at this point, and the answer is not the flip — the flip is an anecdote. It is the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it. MEASUREMENTS.md has the rate and the caveats. To be clear, the blind preference runs are deliberately not cited here: both arms of those were written under a card, so they compare two ways of writing well and cannot see a voice being swapped at all. I apologise for the length of that qualification; it is the qualification that keeps the claim honest.

The ban was read every turn. The habit was not.

---

## 📦 Install

**Two commands and one menu choice.** Let me take these in order, because the order matters more here than anywhere else on the page, and I apologise in advance for the amount of prose sitting around three lines of shell:

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

**Then pick it under `/config` → Output style.** The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. That is the one step `--setup` prints, and I should have said up front that it cannot be done for you.

**What `--setup` did, rather than what it saved you from doing.** ⚠️ It edited your settings file, and a reader letting a tool do that deserves the shape of it before running it rather than after — so: it emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing.

**By hand, for anyone who would rather not have their settings written to.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. ✅ Worth stating plainly, because the alternative is concluding the tool is broken: a style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation — which is why the card lands here and did not through the hook.

---

**Now the hooks, and they are no longer how the card arrives.** This is the hooks block of `~/.claude/settings.json`, and I apologise for the pedantry, but there is no `SessionStart` entry in it and that absence is deliberate:

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

**What the two remaining hooks buy.** Not the voice — the measurement. `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` restates mid-session the card items gated on what has actually been firing, naming the counts, which a file written once cannot do. The voice works without either of them, and I should have led with that.

**Two practical notes.** The commands assume `go install`'s target directory is on PATH in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. There is also `--inject`, the superseded delivery that remains for anyone who wants it and stands down on its own when a cope output style is active.

- **Key takeaways.**
  - `go install`, then `cope-gate --setup`, then `/config` → Output style.
  - The style applies at the next session or after `/clear`.
  - The two hooks are the measurement half, not the voice.
  - PATH, or an absolute path in the hook command.

Install, select, restart. That is the whole install.

---

## 🗣️ Writing in another voice

**The gate reads `.effigy` directly.** A card is usable as written — no Python, no effigy checkout, no compile step in between you and a register you invented this afternoon. I apologise for how much of this page had to come before that sentence.

**Where a card goes and how it is named.** A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. ⚠️ A name that resolves to nothing is an error and nothing is injected — rather than a session quietly writing in the shipped voice while its config names another one, which is the failure worth engineering against.

**Read `card/demo/lecturer.effigy` first.** Not because it is the most extreme, but because it differs from the shipped card on register alone, which makes it the one the discrimination run measured. What a card can change is the voicing axis and only that, for the reasons **The two things cope checks** already gave; I will not argue it a second time, and I apologise for the cross-reference in place of the argument.

**The shortest way to see what a card does without writing one.** [demo/README.md](demo/README.md) — every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice.

**One exception in that directory.** `card/demo/handoff.effigy` is a hypothesis rather than a voice: it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

A card is a file. Drop it in and name it.

---

## 🔁 What the hooks do that a file cannot

**A file written once cannot read what happened since.** That is the whole reason two hooks survived the move to an output style, and I should have made that the heading rather than a heading about hooks.

**`Stop` scores and records.** It reads the reply just written, scores it, and appends which rules fired to the session's rolling state — no prose, only rule names and counts — plus one record per violation to the log.

**`UserPromptSubmit` chooses from that record.** It reads the rolling state, not the violations log, and injects only the card items gated on what has been firing, naming the counts; it falls back to the standing CONTINUE TEST when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. Put differently: the mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted CLAUDE.md cannot do.

**Stated without overselling, because I have been overselling.** ⚠️ It may be the case that this helps, and it is possible that, at least in some circumstances, it helps a lot — but it is one mechanism and not a guarantee, and the A/B in this repo does not separate the refresher from no refresher. `SessionStart --inject` is superseded and off by default.

Two hooks. One reads, one writes.

---

## 📜 Why effigy notation

**The notation was not designed for this.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label — and three of its blocks turn out to do exactly what a prose gate needs, which is a better piece of luck than a design decision. POSTPROC is regex rules with a warn action applied after generation. WRONG holds an anti-pattern beside its replacement. TEST holds a named question with fail and pass examples, which is how a rule names a *move* instead of one wording of it. Worth noting that the third one is the interesting part, and I apologise for burying it at the end of a list.

**And why [basanite](https://github.com/justinstimatze/basanite) is the wrong instrument here.** basanite measures rather than bans — lemma frequency against a baseline over real transcripts, awareness rather than prohibition — which is the right tool when you would rather watch a drift than legislate it, and the wrong one when you want a rule to fire. cope bans: a rule fires or it does not. They compose, different hooks and no shared state, and running both is reasonable. Also worth naming for a reader who wants fewer tokens rather than different structure: [caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies to cut output tokens — a third axis again.

Three projects, three axes, no shared state.

---

## ⚖️ The rules

**Grouped by axis rather than by implementation.** The grouping is the useful cut, and I apologise for the mild bait-and-switch: where each rule lives turns out to be a consequence of the grouping rather than a second fact about it.

**Voicing rules — six.**

- `demo_no_closure` — regex, warn. This card should never accidentally produce a clean, closed ending.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint; the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure rules — eight.**

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` *(interactive)* — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` *(interactive)* — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` *(interactive)* — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` *(interactive)* — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` *(loop)* — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` *(loop)* — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies about the implementation.** A POSTPROC pattern matches a span of text, so it can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go, beside the structure rules, in `internal/scan`. That is why the card ships exactly three POSTPROC rules and not thirty. A reader arriving here expecting a long list of banned phrases should know the list lives in another tool on purpose: basanite measures vocabulary, cope shapes prose, and I should have made that split clearer in the previous section rather than leaving it to this one.

**Where the structure rules do vary — and it is not by card.** It is by who is going to read the turn.

| | **interactive** | **loop** |
|---|---|---|
| **chosen when** | any turn that is not a loop turn | the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself |
| **why** | somebody is waiting at a terminal and the ending is where they decide what happens next | nobody is reading yet |
| **dropped** | — | ⚠️ `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision` |
| **added** | — | ✅ `unverified_done`, `loop_ask` |

**Why the loop lane looks like that.** Nobody is reading yet — a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads the question as an instruction to itself. 🎯 What replaces them is the claim check: a report saying the work is done has to say what it ran.

- **Key takeaways.**
  - Six voicing rules, one of them a card-side pattern.
  - Eight structure rules, all compiled, four of them lane-specific.
  - A pattern can only describe wording, which is why the split falls where it does.
  - The lane depends on the reader, not on the card.

Fourteen rules. Two axes. Two lanes.

---

## 🚩 Flags

**Twenty-six of them, and most readers need three.** I apologise for the table, which is longer than the install it supports.

| flag | default | does |
|---|---|---|
| `--ab` | `false` | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | *(empty)* | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | *(empty)* | read a turn log and report the arms; - reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | *(empty)* | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | *(empty)* | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
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

Three flags matter. The rest are for measuring.

---

## 💾 What lands on disk

**Three files, all `0600`, and one of them quotes you back.** I should have put the disclosure in the heading rather than the first paragraph, so: the violations log stores the text of your replies.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

> **Important:** ⚠️ The log quotes replies back — matched text plus roughly 70 characters of surrounding context, on disk, in plain JSON. `--log` with an empty value disables it. The rolling session state does not, and that asymmetry is worth knowing before you point either at a shared machine.

Three files. One of them is prose.

---

## ✏️ Editing the card

**effigy owns the grammar; cope owns the compile.** `make rules` regenerates `card/rules.json` from the `.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart. I apologise for how much of this section is about budgets rather than about writing.

**The NEVER budget is ten, and it is charged per injection.** Not against the card file — against each injection separately, because SessionStart prints the always-on rules and the refresher prints the evidence-gated ones and no code path renders their union. So a card may hold more NEVER rules in total than the budget and still be perfectly healthy. Anything genuinely over budget is reported at load rather than dropped silently, and the authoritative list of rules that really are discarded unrendered is empty when the card is healthy.

---

**The two card-authored forms, with the syntax exact.** A rule id has to be one the gate already has for `@gate`, and must not collide with one for `@shape`; a wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary, in full.** ✅ Not an approximation of it, because an approximation is worse than none:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**Two behaviours worth knowing before you use either.** A declined rule still runs — only this card's score drops it — so a backfill still reports what it would have caught. And a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies, which is the point of requiring the reason.

- **Key takeaways.**
  - `make rules` regenerates, `make check-rules` keeps enforced and injected in step.
  - The NEVER budget is charged per injection, not per card.
  - `@gate` declines a rule; `@shape` asserts one; both need a reason.
  - A declined rule still runs, and only the score changes.

Two lines in a header. Both need a why.

---

## 📊 Calibrating

**The rules were chosen by backfilling, not by taste.** `cope-gate --backfill` scores a whole session transcript at once, and that is how the current set was arrived at — which I should have said back in the rules section, where a reader would have wanted it.

**Sweeping wider.** `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects` — worth being precise about, because that is not the same as one per project, and reading it as one per project makes the coverage sound broader than it is.

**Watch hits-per-character, not the share of turns hit.** The second number tracks how long the turns were rather than how bad they were, which makes it a length metric wearing a quality metric's clothes. MEASUREMENTS.md has the actual rates; quoting one here would date this file the next time the sweep runs.

Backfill, sweep, then read per character.

---

## ⚠️ Known limits

**The axis split organises these too.** Three limits, taken in order, and I apologise that the largest one is last rather than first.

**First, the small ones.** `labelled_opening` is not a tagger — it is a heuristic about a shape, and it will disagree with a careful reader sometimes. `ask_not_last` says nothing about the *order* of several asks; it only knows whether one is stranded above the ending.

**Second, what the hit rate is and is not.** Roughly four fifths of hits are structure, and the A/B run found that four fifths tracks what a reply was *for* rather than how it was written — so the hit rate is a description of the output, not a judgement of it. The judgement lives in the discrimination test, and that test covers voicing only. Worth naming plainly rather than leaving to inference.

**Third, and largest: the limit **The two things cope checks** already named.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are also the card marking its own homework: a decline lowers that card's score, an assertion raises it, both should be read with
