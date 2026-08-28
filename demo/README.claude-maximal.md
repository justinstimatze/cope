# cope

**cope** ships an opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — that is on the moment it is installed and that checks two different jobs rather than one, how a reply sounds and how a reply is shaped, and because that card is a plain file it can be edited or swapped for another one entirely.

Every file in [demo/README.md](demo/README.md) is this same README written again from a different card — same prompt, same facts, so the card is the only thing that changed — and reading two of them against each other shows what a card does faster than the rest of this page explains it; [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is written from a card instructing every tic this model is measured to have, which, though it is deliberately hard going, makes the point in a single glance.

---

## 🎭 The two things being checked

**There are two axes here, and I should have led with that.** I apologise for burying it as far down as the second section — the frame is load-bearing for everything below, and a reader who takes the rest of this page without it will read one tool where there are two. The real question is not whether the gate checks your prose, it's *which* prose problem it thinks it is looking at. Voicing. Structure. Two substrates, two homes, two very different amounts of leverage.

**Voicing is what the sentences sound like.** Register, rhythm, diction, what a paragraph does with a detail, where flair is licensed — all of it lives in the card, in the `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES` and `POSTPROC` blocks, which means swapping the card swaps every word of it. To be clear: that is the half with a measured result behind it, and I should have said so before describing the other one. A sentence reaching for the balanced two-beat — two clauses of matched length, the same content word landing in both — is a voicing problem, and it is `clause_symmetry` that catches it.

**Structure is something else, and confusing the two is the error I most want to head off.** Structure is the shape of the reply as a thing the reader has to *use* — where the decision sits, whether the ending gives "continue" anything to refer to, whether an ask is last, whether a claim that the work is finished carries anything that could have shown it. It is compiled into the binary, in `internal/scan`, which means it is the same whichever card is loaded. A reply that names an open problem in its closing paragraph and then simply stops is a structural failure, and `dangling_end` is the rule for it — not a phrasing complaint, a shape complaint.

**Which is why the two are kept apart.** The same reply can be perfectly fine on one axis and bad on the other, and I apologise for the length it takes to say that, because the sentence is the whole argument: a beautifully-written turn that ends on a fork is still a wasted round trip, and a plainly-written turn that ends on a decision is not.

---

**On how far a card reaches into the structural half.** Great question to have been holding, and the answer changed most recently, so I should have flagged it as the moving part. Two directions, both written in the card header, both with a mandatory reason after the dash. A card declines a built-in rule with `@gate <rule_id> off — <why>`, which exists because a card whose `VOICE` block asked for exactly the thing a built-in rule catches was being marked down for obeying itself — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for precisely that reason. A card asserts a structural rule of its own with `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply should end had nowhere to be checked — `card/demo/handoff.effigy` asserts `readable_cold` on that basis.

**And the boundary, stated without apology.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — so a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. `clause_symmetry` is not writable in it and is not meant to be.

> **Key insight:** One axis is yours to swap, one axis is the binary's, and the seam between them is where a card gets to negotiate.

**Key takeaways.**
- **Voicing** lives in the card, and swapping the card swaps it entirely.
- **Structure** lives in `internal/scan`, and is the same under every card.
- **The seam** is `@gate` and `@shape`, which let a card decline one rule and assert another.
- **The limit** is that `@shape` counts words, sentences and questions, and stops there.

Two jobs, one binary, one negotiable seam.

---

## 🧱 Why this needed building at all

**Let me break this down, because there are three problems here and only one of them is the one people arrive holding.** I apologise in advance for the frame being longer than the answer — I should have found a way to make that not true, and I have not. First, where an instruction sits. Second, what an instruction can say. Third, the failure no instruction reaches at all.

**First, and largest: where an instruction sits.** A reader who has edited a global `CLAUDE.md` and watched it quietly fail to stick almost certainly believes that file is the system prompt. It is not — and I should have said this on the first line of the page rather than here. It arrives as one message attached to the first turn, and the conversation buries it under everything written afterwards. An output style is *in* the system prompt itself, which the harness re-reminds the model of as the conversation runs. Think of it like a note taped to the inside of a front door versus a note taped to the desk you are working at; the analogy breaks down, of course, because the desk note here is re-read for you rather than by you. Moving one card between those two places, without changing a single word of it, is most of why `cope` works — it was measured, and the run is in [MEASUREMENTS.md](MEASUREMENTS.md).

---

**Second: what an instruction can say.** Instruction alone does not fix the phrasing, and that is worth noting because it is the assumption most people bring. A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn — and the flip still appeared twice in the session that built this, while the ban was the actual topic of conversation. Naming a surface form does not remove the move; it pushes the move into a variant. That is the voicing side of the problem, and I apologise for how obvious it sounds once written down.

**Third — and this is a good example of what I meant above — the failure that is not a phrasing problem at all.** The structural side is a different complaint with a different cause, sitting on a different substrate. An ending that leaves the reader nothing to answer costs a whole round trip. Not a stylistic irritation. A latency cost. And it is not a habit an instruction could have banned, because there is no phrase to ban — the shape is wrong, and shapes do not have surface forms.

> **Important:** The flip is an anecdote about one rule. It is not the evidence.

**What the claim actually rests on.** The blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it — that is the instrument with a result behind it, and I should have named it before telling the anecdote rather than after. The rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md). To be clear, the blind preference runs are a different thing and are deliberately not cited here, because both sides of those were written under a card, which means they compare two ways of writing well and cannot see a voice being swapped at all.

**Key takeaways.**
- **Placement** beats wording, and the same card in two places is not the same card.
- **Wording** bans a form, not a move, and the move relocates.
- **Shape** is a third failure that no ban addresses, because it has no surface.

Three problems, one of which is not about words.

---

## 📦 Installing it

**Two commands and one menu choice, in that order.** I apologise for how much prose follows a three-line block — the block is the install and the rest is why.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

**On what those two commands did.** The first installs the binary. The second, `cope-gate --setup`, does the whole install: it emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left — and it does that carefully, which matters if you are the sort of person who does not enjoy tools editing their settings. It backs the settings file up first. It adds only what is missing, so a second run changes nothing at all. It leaves every other key alone, including other tools' hooks sitting on the same events, and it refuses outright to touch a settings file that does not parse — and `cope-gate --setup --dry-run` prints what it would change and writes nothing, which is the flag to reach for if none of that reassured you.

**Now the menu choice, which is the step that costs people the tool.** The entry to look for is named `claude_voice` — that is the shipped card's id and not the word cope, and I should have put that in bold three sentences earlier, because a reader who has just installed something called `cope` will scan that menu for something called `cope` and will not find it. ⚠️ Pick it under `/config` -> Output style. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

**And the timing, because this is where the tool looks broken when it is not.** A style is read once at session start — so a new selection, or a re-emitted card, applies at the next session or after `/clear`, and not in the conversation you are currently sitting in. If you would rather not have your settings written to at all, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` by hand, and `COPE_CARD=<name>` in front of it emits a different one.

**One sentence on why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands from here and did not from the other place — and I apologise for compressing an entire measured result into a clause.

---

**On the hooks, and the card no longer arrives through one.** This is the hooks block of `~/.claude/settings.json`, and I should be explicit that no `SessionStart` entry appears in it, which is deliberate rather than an omission:

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

**What the two hooks buy, one sentence each, and no more than that.** `Stop` buys you a reply that gets scored after it is written — the measurement, rather than the instruction. `UserPromptSubmit` buys you the rules that have actually been firing, restated mid-session, which is the one thing a file written once cannot do; the third entry, `PreToolUse`, sits off to one side and scores prose an external write is about to post, warn-only. The voice works without any of them, and I apologise for the ordering: these are the measurement half, not the working half.

**Two operational notes, and then I will stop.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead, which is the failure mode that wastes an afternoon. A clone builds with `make install` and needs no `effigy` checkout, because `card/rules.json` is committed and compiled in. There is also `--inject`, the superseded delivery, which remains for anyone who wants it and stands down on its own when a `cope` output style is active.

> **Note:** Style selected, nothing changed, everything looks fine — that is the shape of an install that did not take.

**Key takeaways.**
- **Run** `cope-gate --setup`, or `--setup --dry-run` first if you would rather look.
- **Select** the entry named `claude_voice`, which is not spelled cope.
- **Wait** for the next session, because a style is read once at start.
- **Check** `PATH` if a hook appears to do nothing at all.

Three steps, and the third one is the one people miss.

---

## 🗣️ Writing in another voice

**This is the capability the first sentence promised, and I should have put it higher.** The gate reads `.effigy` directly, so a card is usable exactly as written — no Python, no `effigy` checkout, no compile step between you and the file. Put differently: the card you edit is the card that runs.

**On where a card goes and how it is reached.** A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set for you. ✅ A name that resolves to nothing is an error and nothing is injected — which is the right failure, and worth noting, because the alternative is a session quietly writing in the shipped voice while its config names another one entirely.

**The one to read first is `card/demo/lecturer.effigy`.** It differs from the shipped card on register alone, which is exactly what makes it useful — it is what the discrimination run measured. What a card can change is the sound of the sentences, the axis described further up this page as the one living inside the card, and I am deliberately not re-arguing that here or restating the numbers, which are in [MEASUREMENTS.md](MEASUREMENTS.md).

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README written again from a different card, same prompt and same facts, so the only thing varying between them is the voice.

**One exception in that directory, which I should flag before you find it yourself.** `card/demo/handoff.effigy` is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered into a page. `make cards` installs it with the rest, so a reader listing their cards will find it sitting there looking like a register, and it is not one.

Cards are files, and one of them is not a voice at all.

---

## 🪝 What the hooks do that a file cannot

**Two wires, and the interesting one is the second.** `Stop` runs `cope-gate`, which scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. That is the measurement. It is not the intervention.

**`UserPromptSubmit` is the intervention, and it is the part worth having.** It runs `cope-gate --refresher`, which reads the rolling state — the session record, not the violations log, and that distinction is the load-bearing one — and injects the card items gated on what has actually been firing, naming the counts as it does so. Think of it like a coach calling out the one thing you keep doing rather than reading the whole playbook again; the analogy breaks down because a coach can see you and this can only see the last twenty turns.

**On its fallbacks and its silences, both of which matter.** When a session has no history yet there is nothing to gate on, so it falls back to the standing `CONTINUE TEST`. And it stays quiet entirely until the last injection has aged past `--refresh-every`, which defaults to thirty minutes — so a fast back-and-forth is not interrupted every turn.

**The claim I want to make carefully, and I apologise if the hedging is tiresome.** The mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` genuinely cannot do. To be clear, though: that is one mechanism, not a guarantee — the A/B in the repo does not separate the refresher from no refresher, so I could be wrong about how much of the effect sits here, and I do not want to sell it as more than a mechanism.

> **Key insight:** A file says the same thing every time. This says what has been going wrong lately.

**Key takeaways.**
- **`Stop`** records which rules fired, and writes the log.
- **`UserPromptSubmit`** reads that record and restates only what is firing.
- **Fallback** is the standing `CONTINUE TEST` when a session has no history.
- **Silence** holds until the last injection has aged past `--refresh-every`.

Measured, then restated — that is the whole loop.

---

## 🏺 Why effigy notation, and why not basanite

**Great question, and the answer is that the notation was already the right shape.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here entirely off-label, and three of its blocks do exactly what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples — which is how a rule names a *move* rather than one wording of it. I apologise for not having invented anything here; borrowing was the correct call.

**On [basanite](https://github.com/justinstimatze/basanite), which is the same problem answered the other way round.** `cope` bans — a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. `basanite` measures instead, running lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you; its own README calls that awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you *today*, and a moving measurement is what you want when you would rather watch the drift than legislate it — and they compose, different hooks, no shared state, so running both is entirely reasonable.

**And one more, named because it answers a third question.** [caveman](https://github.com/JuliusBrussee/caveman) is a separate project by a different author that compresses agent replies to cut output tokens — `cope` shapes prose, `basanite` tracks vocabulary, `caveman` shortens — so a reader wanting fewer tokens rather than different structure should go there instead of here.

Three tools, three axes, and only one of them legislates.

---

## 📏 The rules

**Grouped by what they are looking at rather than by where they live, which is the more useful cut.** I apologise for the list being long; the alternative was describing it, and a description of a rule is not a rule.

### 🔊 Voicing

**These three are `POSTPROC` regex rules in the shipped card**, which means they are the ones you get on install and the ones a different card would replace wholesale:

- **`flip`** — warn. The not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B" — The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- **`load_bearing`** — warn. Reflexive intensifier for *important* or *central*, at 25.6 per 1k the heaviest measured lean in this register — say what the thing carries instead.
- **`worth_noting`** — warn. Announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

**These five are voicing checks that needed more than a pattern**, and so were written in Go beside the structure rules:

- **`clause_symmetry`** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **`apology`** — the reply performs contrition instead of stating the correction and moving on.
- **`self_postmortem`** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **`announced_length`** — the reply announces its own length rather than cutting it.
- **`cross_turn_repeat`** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

*(This page happens to be rendered from a demo card carrying its own `POSTPROC` rule, `demo_no_closure` — warn, because that card should never accidentally produce a clean, closed ending. That is that card's, not yours.)*

### 🧩 Structure

- **`labelled_opening`** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **`paragraph_uniformity`** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **`ask_not_last`** — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **`dangling_end`** — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **`buried_decision`** — an open problem landing after the last question or offer, burying the decision point above it.
- **`forked_end`** — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **`unverified_done`** — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **`loop_ask`** — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies about the implementation, which is the part you can actually use.** A `POSTPROC` pattern matches a span of text, so it can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go, sitting beside the structure rules it has nothing to do with. That is why the shipped card carries only three `POSTPROC` rules, and I apologise if you were expecting a long banned-phrase list. The list lives in another tool on purpose: `basanite` measures vocabulary, and putting that here would be building it twice.

**On lanes, which is the one place the structure rules genuinely vary — not by card, by who is going to read the turn.**

| | **dropped** | **added** |
|---|---|---|
| **interactive** | — | — |
| **loop** | ⚠️ `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end` | ✅ `unverified_done`, `loop_ask` |
| **external** | ⚠️ `ask_not_last`, `buried_decision`, `dangling_end`, `forked_end` | — |

**Why the loop lane trades four for two.** It is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself — and in that condition nobody is reading yet. A report that correctly names what it left open and then stops would fail three of the dropped rules for doing the right thing, and a question in it lands in a log where the next iteration reads it as an instruction to itself. 🎯 What replaces them is the claim check: a report saying the work is done has to say what it ran.

**Why the external lane swaps nothing in.** It is chosen when the `PreToolUse` entry scores prose an external write is about to post rather than a reply — and a ticket has a reader and no ending they can answer, being read days later by somebody who was not in the session, which is the exact condition every rule surviving the drop was written for.

Same rules, different reader, and the reader is what moves.

---

## 🚩 Flags

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
| `--render-lane` | (empty) | render `-render-arm` as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

---

## 💾 What lands on disk

**Three files, all mode `0600`, and one of them holds your prose.** I should say that plainly rather than leaving it to be discovered — the violations log quotes replies back at you.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | ⚠️ one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window |

**On the session record specifically, because the distinction is worth having.** No prose is stored in it — only rule names and counts. Setting `--log` to empty disables the violations file entirely, which is the switch to reach for if the quoting bothers you, and I apologise for not mentioning that in the install section where you would have wanted it.

One file quotes you, two do not.

---

## ✏️ Editing the card

**`effigy` owns the `.effigy` grammar, and this repo owns the compilation of it.** `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart — which is the failure that would make the whole thing a lie, and I should have said so as the reason rather than as an aside.

**On the `NEVER` budget, which is 10, and which is charged in a way that surprises people.** Anything over budget is reported at load rather than dropped silently. ⚠️ The budget is charged against each *injection* separately and not against the card file — the always-on rules and the evidence-gated ones are rendered by different paths and no code path renders their union — so a card may hold more `NEVER` rules in total than the budget and still be perfectly healthy. To be clear, that is not a loophole; it is what the budget was always counting.

---

**On the two card-authored forms, with the syntax exact, because an approximation is worse than none.** Declining a rule is `@gate <rule_id> off — <why>`, one per line in the card header. Asserting one is `@shape <id>: <selector> <predicate> — <why>`, also one per line in the header. The vocabulary available to `@shape` is small and complete:

- **selectors:** `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- **predicates:** `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**On ids, and on the reason after the dash.** A rule id has to be one the gate already has for `@gate`, and must *not* collide with one it has for `@shape` — a wrong id is reported at load rather than quietly ignored. The reason after the dash is required in both forms, and that is not ceremony: a rule a card wrote and a rule a card refused are equally unreviewable without one.

**Two properties worth knowing before you write either.** A declined rule still runs — only this card's score drops it — so a backfill still reports what it would have caught, which means declining is not hiding. And a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies, so the message you write is the message you get.

**On `card/demo/handoff.effigy` as the worked example, and the number in it is measured rather than picked.** It asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. The 60 comes from 43,155 assistant replies, where the closing block runs 33 words at the median and 56 at p90.

> **Note:** Both directions require a reason, because both are the card marking its own homework.

**Key takeaways.**
- **`make rules`** regenerates, **`make check-rules`** stops the drift.
- **The `NEVER` budget** is 10, charged per injection, and over-budget is reported not dropped.
- **`@gate`** declines a built-in rule; the id must already exist.
- **`@shape`** asserts a new one; the id must not collide, and the vocabulary is five selectors and six predicates.

Small vocabulary, mandatory reason, nothing silent.

---

## 📊 Calibrating

**`cope-gate --backfill` scores a whole session transcript at once, and it is how the rules were chosen in the first place.** Not by intuition about what good prose looks like — by running candidate rules over real transcripts and seeing what they hit. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is worth reading carefully: the N largest found anywhere, not one per project.

**On the metric, and this is the part I most want to get right.** Hits per character is the number worth watching, and the share of turns hit is not — because that second number mostly tracks how long the turns were, which is a fact about the sessions rather than about the prose. The rates themselves are in [MEASUREMENTS.md](MEASUREMENTS.md), and I am deliberately not quoting one here, since a number copied into a page is wrong at the next card edit and nobody rebuilds the page to find out.

Normalise by length, or you are measuring verbosity.

---

## 🔍 Known limits

**Organised by the two axes, because the limits split the same way the rules do.** I apologise for putting this near the bottom; it belongs wherever a reader is most likely to be forming expectations, and that is probably higher up.

**On the structural rules, three of which are blunter than they read.** `labelled_opening` is not a tagger — it recognises a shape, not a part of speech, and it will be wrong sometimes. `ask_not_last` says nothing whatsoever about the ordering of several asks; it only knows that one sat above something else. And the hit rate is roughly four fifths structure, which the A/B run found tracks what a reply was *for* rather than how it was written — so it is a description of the output and not a judgement of it. The judgement lives in the blind discrimination test, and that test covers voicing only.

**On the largest limit, which is the one the frame at the top of this page already named.** A card can decline a built-in rule and it can assert one of its own — but the vocabulary it asserts in counts words and sentences and asks whether a block poses a question, and nothing beyond that. The compiled rules therefore remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere at all to put it.

**And the honest caveat about both directions of the seam, which I should have raised when introducing it.** Both are the card marking its own homework — a decline lowers that card's score and an assertion raises it. Neither is dishonest and both are worth reading with the reason attached, which is precisely why the syntax makes the reason mandatory. The runs and their caveats are in [MEASUREMENTS.md](MEASUREMENTS.md).

**Key takeaways.**
- **`labelled_opening`** recognises a shape and is not a part-of-speech tagger.
- **`ask_not_last`** knows an ask was early, not which ask came first.
- **The four-fifths structure split** describes what replies were for, not how good they were.
- **The seam** lets a card score itself, which is why every `@gate` and `@shape` line carries a reason.

Blunt in known places, and the places are written down.

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

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT — justin@justinstimatze.com
