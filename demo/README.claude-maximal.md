# 🪣 cope

**What this is.** cope ships an opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — that is on the moment it is installed and is the whole product for most readers, that checks two different jobs rather than one, how a reply sounds and how a reply is shaped, and that is, only after both of those, a file you can edit or swap.

**Start here, honestly.** Read [demo/README.md](demo/README.md) — every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it; [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card that instructs every tic this model is measured to have, and it makes the point in a single glance.

---

## ⚖️ The two things

**The frame first, and I should have led with it.** The real question is not whether a reply is good, it's which of two jobs it failed at — and those two jobs sit on different substrates. **Voicing** is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely — VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC — so swapping the card swaps every word of it, and that is the half with a measured result behind it. **Structure** is the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Structure is compiled into the binary, so it is the same whichever card is loaded, varying only by lane.

**Two instances, because abstraction is my mistake to avoid here.** A sentence reaching for the balanced two-beat — two clauses of near-equal length repeating a content word across the joint — is a voicing problem, and the card is where it is answered. A reply naming an open problem in its last paragraph and then stopping, leaving "continue" nothing to refer to, is a structural one, and no register change fixes it. To be clear: the same reply can be fine on one axis and bad on the other, which is the whole reason to keep them apart.

|  | **voicing** | **structure** |
|---|---|---|
| **lives in** | the card | the binary |
| **swappable** | ✅ entirely | ⚠️ no |
| **varies by** | which card you picked | which lane the turn was written in |

**How far a card reaches into the other half.** Not zero — and this is the part that most recently changed, which I should have flagged before the table rather than after it. Two directions, both stated in the card header, one per line. A card declines a built-in rule with `@gate <rule_id> off — <why>`, which exists because a card whose VOICE asked for something a built-in rule catches was being marked down for obeying itself: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch. A card states a structural rule of its own with `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and the 60 is measured, since across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

**The boundary, without apology.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — so a card wanting a check outside both that vocabulary and a POSTPROC regex has nowhere to put it.

Two jobs, two homes, one seam between them.

---

## 🧨 The problem

**Instruction alone does not fix the phrasing.** Not because instruction is ignored — because naming a surface form pushes the move into a variant. A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn, and the flip still appeared twice in the session that built this, while the ban was the topic of the session. I should have caught both at the time and did not, which is itself the evidence: the ban was in context, the ban was the subject, and the substrate underneath the ban went on producing the move in a new wording.

**The structural complaint is a different complaint with a different cause.** An ending that leaves the reader nothing to answer is not a phrasing habit an instruction could have banned — it costs a whole round trip, and it costs it whatever register the sentences were in. Put differently: the flip wastes a reader's patience, and a dangling ending wastes a reader's turn. That is the second of the two axes, and I apologise for taking two paragraphs to arrive at a distinction the section above already drew.

**What the claim actually rests on.** The flip is an anecdote about one rule and should be read as one — my mistake if the paragraph above sounded like a result. The instrument with a result behind it is the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it. See [MEASUREMENTS.md](MEASUREMENTS.md) for the rate and the caveats.

- **Key takeaways.**
  - The flip survived a standing ban that was, at the time, the topic — naming a form moves the form.
  - A dangling ending is a different failure on a different axis, and instruction was never going to reach it.
  - The measured claim is the discrimination test, not the anecdote; the rates live in MEASUREMENTS.md.

One anecdote, one round trip, one instrument.

---

## 🔧 Install

**Two commands and one menu choice.** Not a configuration exercise — three steps, and the order is load-bearing, which I should have said before the block rather than inside it.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> pick the cope card
```

**What `--setup` did, rather than what it saved you.** It emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left — and because a reader letting a tool edit their settings deserves the shape of it up front, here is the shape: it backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse, with `--setup --dry-run` printing what it would change and writing nothing. I should have put that sentence above the code block and did not.

**The by-hand route, for anyone who would rather not have their settings written to.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one. Pick it under `/config -> Output style`; the standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

> **Important:** A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. ⚠️ A reader who picks a style and watches the running conversation for a change will conclude the tool is broken, and that would be my fault for not saying so here.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation — which is why the card lands here and did not through the hook.

---

**Then the hooks, which are no longer how the card arrives.** Worth noting, and I apologise for the repetition of a point the previous paragraph implied: these two do the measuring, not the delivering. This is the `hooks` block of `~/.claude/settings.json`.

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

**What the two buy you.** `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log — and `UserPromptSubmit` reads that rolling state and injects the card items gated on what has been firing, naming the counts, which a file written once cannot do. To be clear: the voice works without either of them, and these are the measurement half.

**One operational note, and I should have volunteered it sooner.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook which silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in, and `cope-gate --inject` remains as the superseded delivery for anyone who wants it, standing down on its own when a cope output style is active.

- **Key takeaways.**
  - `go install`, then `cope-gate --setup`, then pick the style under `/config`.
  - `--setup` backs up, adds only what is missing, and refuses a settings file that does not parse.
  - The style applies at the next session; the two hooks are measurement, not delivery.

Three steps, one menu, no restart of the tool itself.

---

## 🎭 Writing in another voice

**A card is usable as written.** Not compiled first — read directly: the gate reads `.effigy`, so there is no Python step and no effigy checkout in the loop, and I apologise for having buried that under the install section for as long as it was buried. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`, `--rules` takes a path from anywhere, and `make cards` installs the demo set.

**One failure mode worth naming out loud.** A name that resolves to nothing is an error and nothing is injected — rather than a session quietly writing in the shipped voice while its config names another one. That distinction is load-bearing, and it is the kind of thing that should be stated before somebody trusts a typo.

**Which card to read first.** `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and is what the discrimination run measured — and what a card can change is the voicing axis, as the section called *The two things* argued at length and I will not argue again here. The numbers are in [MEASUREMENTS.md](MEASUREMENTS.md); quoting one here would be my second mistake in a row.

**The shortest route to seeing what a card does without writing one.** Read [demo/README.md](demo/README.md): every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice — the index there is generated from whatever was last built, which is why this page does not carry a second copy of the list.

**One exception in that directory.** `card/demo/handoff.effigy` is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered; `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

- **Key takeaways.**
  - `.effigy` is read directly: `--card`, `COPE_CARD`, or `--rules <path>`.
  - An unresolvable card name errors rather than silently falling back.
  - `lecturer.effigy` is the one to read; `handoff.effigy` is not a register.

A card is a file, and the file is the register.

---

## 🔁 What the hooks do differently

**A pasted instruction file cannot see what happened.** That is the whole difference — not volume, not placement, but feedback. `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` reads that rolling state — the state, not the violations log — and injects the card items gated on what has been firing, naming the counts, falling back to the standing CONTINUE TEST when the session has no history yet and staying quiet until the last injection has aged past `--refresh-every`. I should have said plainly at the top that the state holds rule names and counts and no prose.

**Stated plainly, without overselling it.** It is one mechanism, not a guarantee — the A/B in the repo does not separate the refresher from no refresher, and I would rather concede that here than let the paragraph above imply otherwise. `SessionStart --inject` is superseded and off by default: measured 2026-08-03, a card asking for a bolded label on every paragraph and an emoji on every heading produced, through that hook, prose indistinguishable from no card at all.

Chosen from what fired, not fixed in advance.

---

## 🧩 Why effigy notation, and why not basanite

**Off-label, and deliberately so.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, and three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples — which is how a rule names a move instead of one wording of it, the TEST here being `NOTICING TEST`. The card as rendered runs 17,418 characters of prompt text, which is a number I should have offered earlier for anyone budgeting context.

**[basanite](https://github.com/justinstimatze/basanite) answers the same problem the other way round.** cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it — whereas basanite measures, lemma frequency against a baseline over real transcripts, reporting what you have actually been leaning on lately and leaving the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it — they compose, different hooks and no shared state, and running both is reasonable.

**A third axis, named because it is not this one.** [caveman](https://github.com/JuliusBrussee/caveman) is a separate project, by a different author, that compresses agent replies to cut output tokens — cope shapes prose, basanite tracks vocabulary, caveman shortens, and a reader wanting fewer tokens rather than different structure should go there instead. I apologise for not naming it in the opening section, where a reader with that goal would have saved themselves this page.

Three tools, three axes, one overlap: none.

---

## 📋 The rules

**Grouped by job, not by implementation.** The real question is not which file a rule lives in, it's which of the two things it is watching — and I should have organised an earlier draft of this section that way rather than by mechanism.

**The voicing rules.** Six of them, and the card is where each is answered.

- `demo_no_closure` (regex, action `warn`) — this card should never accidentally produce a clean, closed ending.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**The structure rules.** Eight, compiled in, the same whichever card is loaded.

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one; sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies about the implementation.** Here is the part a reader can actually use, and I apologise for making it wait until after two lists. A POSTPROC pattern matches a span of text, so it can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go beside the structure rules, and that is why the card ships exactly one POSTPROC rule. Think of it like a spell-checker that can flag a word but not a habit; the analogy breaks down, of course, because `cross_turn_repeat` reads a whole session and no spell-checker reads yesterday.

> **Note:** A reader arriving here expecting a long list of banned phrases should know the list lives in another tool on purpose. ⚠️ That is basanite's job — lemma frequency against a baseline — and duplicating it here would have been my error to make.

**Where the structure rules do vary.** Not by card — by who is going to read the turn. The **interactive** lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The **loop** lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself, and it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision` and adds `unverified_done` and `loop_ask` — because nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran.

- **Key takeaways.**
  - Six voicing rules, eight structure rules, grouped by job rather than by file.
  - A regex can describe wording only, so one POSTPROC rule ships and the rest are Go.
  - The lanes swap four rules for two, on the question of whether anyone is reading.

Same rules, different reader, different ending.

---

## 🎛️ Flags

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

**One note under the table, which I should have put above it.** The defaults above are the defaults the binary prints — not a paraphrase, and not a rounding of `30m0s` into "half an hour".

Twenty-six flags, three of which most readers touch.

---

## 💾 What lands on disk

**Three files, all `0600`, all under `$XDG_STATE_HOME/cope`.** Worth stating plainly rather than leaving in a table nobody reads — and to be clear, the first one quotes your replies back at you.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window — no prose is stored, only rule names and counts |

**The distinction between the first and the third is load-bearing.** The log holds excerpts of what you wrote; the session state holds rule names and counts and nothing else — and `--log` with an empty value disables the first entirely, which I should have mentioned in the install section for anyone who would rather it never existed.

One file has your prose in it. ✅ It can be turned off.

---

## ✏️ Editing the card

**Who owns what.** effigy owns the `.effigy` grammar — not cope, which is worth being clear about before anyone files a syntax bug in the wrong repo. `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift.

**The NEVER budget, and the thing about it that surprises people.** The budget is 10, and anything over it is reported at load rather than dropped silently — but the budget is charged against each injection separately, not against the card file, since `SessionStart` prints the always-on rules and the refresher prints the evidence-gated ones and no code path renders their union. So a card may hold more NEVER rules in total than the budget and still be perfectly healthy, and `never_rules_over_budget` is the authoritative list of rules that really are discarded unrendered — empty when the card is healthy. I should have made that distinction the first sentence of the paragraph rather than the third.

---

**The two card-authored forms, verbatim.** One per line in the card header, and the reason after the dash is required in both — because a rule a card wrote and a rule a card refused are equally unreviewable without one.

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary, spelled out, because an approximation of it is worse than none.**

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**Two constraints on ids, and one on scoring.** A rule id has to be one the gate has for `@gate` and must not collide with one for `@shape` — a wrong id is reported at load rather than ignored, which is my preference and, I would argue, the safer default. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

- **Key takeaways.**
  - effigy owns the grammar; `make rules` regenerates and `make check-rules` prevents drift.
  - The NEVER budget of 10 is per injection, not per card file, and over-budget is reported at load.
  - Both card-authored forms need a reason, a valid id, and — for `@gate` — the honesty of a lowered score.

A card can argue with the gate, in writing, with its reasons attached.

---

## 📐 Calibrating

**One command scores a whole session.** `cope-gate --backfill` scores every assistant turn in a transcript at once, and it is how the rules were chosen rather than a diagnostic added afterwards — `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same as one per project, and I should have named that difference before somebody read a sweep as a per-project survey.

**Which number to watch.** Hits-per-character, not the share of turns hit — because that second number tracks how long the turns were, which is a fact about verbosity dressed up as a fact about quality. The rates are in [MEASUREMENTS.md](MEASUREMENTS.md), and quoting one in this paragraph would give it a permanence it has not earned.

A rate per character, over a corpus, not a score per turn.

---

## 🚧 Known limits

**The axis split is what organises the limits, too.** Two small ones and one large, and I apologise for the fact that the large one is the same limit the opening section named — restating it is reinforcement rather than repetition, but it is still a restatement.

**The two small ones.** `labelled_opening` is not a tagger — it recognises a shape, not a part of speech, and it will be wrong on sentences a linguist would classify correctly. `ask_not_last` says nothing about the order of several asks; it notices that one was left upstream and stops there.

**What the hit rate is and is not.** Roughly four fifths of the hits are structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — so the hit rate is a description of the output and not a judgement of it, which is a sentence I should have written into the calibration section as well. The judgement lives in the discrimination test, and that test covers voicing only.

**The largest limit, as it now stands rather than as it stood.** A card can decline a built-in rule and can write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are the card marking its own homework: a decline lowers that card's score and an assertion raises it, both are worth reading with the reason attached, and that is precisely why the syntax requires one — see *The two things* above, and [MEASUREMENTS.md](MEASUREMENTS.md) for what any of it has been measured to do.

- **Key takeaways.**
  - `labelled_opening` recognises a shape, not a grammatical role.
  - `ask_not_last` notices a stranded ask, not the ordering of several.
  - The hit rate describes; the discrimination test judges, and only voicing.
  - A card can decline and assert, and both are self-graded, which is why reasons are mandatory.

Two small limits, one structural one, all of them writable down.

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

**One reading of that table.** Eleven entries, and the split down the middle is the same split the whole page has been about — `card/` is voicing and `internal/scan/` is structure, with `replay/` and `MEASUREMENTS.md` being where either one gets checked.

Eleven paths, two axes, one seam.

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT — justin@justinstimatze.com
