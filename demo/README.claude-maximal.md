# 🏺 cope

**What cope is.** cope ships an opinionated card — a cope is the upper half of a foundry mould, the half carrying the shape being cast into — and that card is live the moment you install it, which makes it the whole product for most readers; what it checks is not one job but two, how a reply sounds and how a reply is shaped, and only after those two does the third thing matter, which is that the card is a file, so it can be edited or swapped. I apologise for the length of that sentence — it carries three claims and I should have found a way to make it carry them more lightly.

**Start here rather than here.** [demo/README.md](demo/README.md) holds this same README written again from each installed card in turn — same prompt, same facts, so the card is the only thing that changed between them — and reading two of them against each other shows what a card does faster than the rest of this page explains it; [demo/README.claude-maximal.md](demo/README.claude-maximal.md) is the one that makes the point in a single glance, written from a card that instructs every tic this model is measured to have, though it is deliberately hard going. Thank you for your patience with a detour offered before the page has earned it.

---

## 🧭 The two things it checks

**TL;DR for this section.** There are two axes, they fail independently, and one of them lives in the card while the other lives in the binary — and I should say up front that everything below is written against this split, so a reader who skims here will be lost twice later.

**Voicing is the first axis.** Voicing is what the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed — and it lives in the card, entirely: `VOICE`, `TRAITS`, `NEVER`, `WRONG`, `MES` and `POSTPROC`. Swapping the card swaps every word of it. A sentence reaching for the balanced two-beat, two clauses of near-equal length repeating a content word across the joint, is a voicing failure and nothing else, and it is worth noting that this is the half with a measured result behind it: the blind discrimination test. I should have said "measured" more carefully there — the run and its caveats live in [MEASUREMENTS.md](MEASUREMENTS.md), and I am deliberately putting no rate on this page.

**Structure is the second, and it is not a matter of taste.** Structure is the shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. A reply that names an open problem in its closing paragraph and then stops is a structural failure, because the reader has nothing to answer. That check is compiled into `internal/scan`, so it is the same whichever card is loaded, varying only by lane — and a card can write down how it wants a reply to end in prose that the gate has no way to hear. That last clause is the one I most want to be exact about, and I was not exact about it in an earlier draft.

**Why keep them apart at all?** Because the same reply can be clean on one axis and bad on the other, and a single score would hide it. Put differently: prose that sounds nothing like this page can still strand its reader, and prose that strands nobody can still be unreadable.

| | **voicing** | **structure** |
|---|---|---|
| **what it is** | what the sentences sound like | where the decision sits |
| **who owns it** | the card | the binary |
| **swap the card** | ✅ everything changes | ⚠️ nothing changes |

---

**Two ways a card reaches into the compiled half.** Not one, and this is the part that most recently changed, so I should be precise rather than brief. A card declines a built-in rule with `@gate <rule_id> off — <why>`, one per line in the card header — `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its `VOICE` block asks for the balanced landing and the arriving close those two rules catch, and a card marked down for obeying itself is a card being scored against somebody else's taste. A card states a rule of its own with `@shape <id>: <selector> <predicate> — <why>`, also one per line in the header — `card/demo/handoff.effigy` asserts `readable_cold`, `last paragraph words <= 60`, because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way. That 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

**The boundary, stated once and without apology for it.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — it cannot express what the compiled rules express, `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. I apologise for stating a limit this flatly in a section otherwise about capability.

**Key takeaways**

- **Voicing** is card-owned, swappable, and the half with a measured result.
- **Structure** is binary-owned and identical whichever card is loaded.
- **A card reaches in twice** — `@gate` to decline, `@shape` to assert.
- **The reach stops** where words, sentences and questions stop.

Two axes, one card, one binary, and one seam between them.

---

## 🧩 The problem, in three parts

**Announcing the count before enumerating it.** There are three problems here and they are not the same problem — where an instruction sits, what an instruction can say, and a failure no instruction reaches — and I am taking them in that order because most readers arrive holding the first one backwards. That ordering is load-bearing, and I should have led with it in an earlier draft rather than burying it.

**First, where the instruction sits.** A reader who has edited a global `CLAUDE.md` and watched it not stick almost certainly believes that file is the system prompt. It is not — it arrives as one message attached to the first turn, and the conversation buries it under everything written after it. An output style is in the system prompt itself, which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works; that was measured, and the run is in [MEASUREMENTS.md](MEASUREMENTS.md) rather than here. Think of it like a note left on a colleague's desk versus a note taped to their monitor — the words are identical and the substrate underneath them is not, though the analogy breaks down, because a desk does not silently accumulate forty turns of other paper.

---

**Second, what an instruction can say.** Instruction alone does not fix the phrasing, and I have a small, embarrassing piece of evidence for that: a global `CLAUDE.md` banning the "not A, it's B" flip is read every turn, and the flip still appeared twice in the session that built this, while the ban was the topic under discussion. That is my mistake twice over — once for writing it and once for not noticing until the gate did. Naming a surface form pushes the move into a variant, which is the voicing side of the problem and the reason `POSTPROC` patterns can only ever be part of the answer.

**Third, and this one is not a phrasing problem at all.** The structural complaint has a different cause: an ending that leaves the reader nothing to answer costs a whole round trip, and no instruction could have banned it as a habit, because it is not a habit — it is a shape. To put a finer point on it, you cannot regex your way to a reply that closes where the decision is.

**What the claim actually rests on.** The flip is an anecdote about one rule and I would rather not have it doing more work than that. What the claim rests on is the blind discrimination test, where a reader shown only a voice's own description of itself picks which of two replies was written under it; the rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md). Worth noting for completeness — the blind preference runs are not the evidence here, because both sides of those were written under a card, so they compare two ways of writing well and cannot see a voice being swapped at all.

**Key takeaways**

- **Placement** beats wording: same card, different slot, different outcome.
- **Wording** alone leaks into variants, measurably and repeatedly.
- **Shape** is a third failure that no instruction addresses.
- **The evidence** is the discrimination test, not the anecdote.

Three problems, one of which is not about words.

---

## 📦 Installing it

**Two commands and one menu choice.** That is the whole install, and I apologise in advance for how much prose sits underneath something this short.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

**What those two commands did.** `cope-gate --setup` does the whole install: it emits the output style, wires the hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `cope-gate --setup --dry-run` prints what it would change and writes nothing, which is the flag to reach for if you would rather read the diff before you own it. I should have put that reassurance ahead of the command rather than after it, since a reader letting a tool edit their settings deserves the shape of it first.

**Then the menu choice, which is the step people lose.** Pick the style under `/config` -> Output style. ⚠️ The entry to look for is named `claude_voice` — the shipped card's id, and not the word cope — so a reader scanning that menu for something called cope will not find it, and an install that ends with a style selected and nothing changed is the one failure on this page that costs you the tool entirely. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

> **Important:** a style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear`. A reader who picks the style, watches the running conversation not change, and concludes the tool is broken has been failed by this paragraph and not by the tool.

**If you would rather not have your settings written to.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` by hand, and `COPE_CARD=<name>` in front of it emits a different one. And the reason the card goes here at all rather than into a hook: an output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through the hook.

---

**Now the hooks, which are a different purchase entirely.** The card no longer arrives through one, so do not read this block as delivery — read it as measurement.

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

**What the two main hooks buy you, in one sentence each.** `Stop` buys you a score on the reply after it is written, and `UserPromptSubmit` buys you a mid-session restatement of the rules that have actually been firing, which a file written once cannot do — the mechanics of both are further down and I am not going to say them twice. The voice works without either of them; these are the measurement half, and I should be clear that they are optional rather than load-bearing.

**Three small notes, and then I will stop.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. And `cope-gate --inject`, the superseded delivery, remains for anyone who wants it and stands down on its own when a cope output style is active.

**Key takeaways**

- **Two commands** install the binary and wire everything.
- **One menu entry**, named `claude_voice`, is the step that gets missed.
- **A style applies next session**, never to the running one.
- **The hooks are measurement**, not delivery.

Install, select, restart — and the card is in the system prompt.

---

## 🎭 Writing in another voice

**The gate reads `.effigy` directly.** Not a compiled intermediate — the source. A card is usable as written, needs no Python and no effigy checkout, and that matters more than it sounds like it should: the file you edit is the file the gate reads, so there is no build step between an idea about register and a session written in it. I apologise for making a point of the absence of a build step, which is the sort of thing only the person who removed it finds interesting.

**Where a card lives and how it is named.** A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. ✅ A name that resolves to nothing is an error and nothing is injected — the alternative, a session writing in the shipped voice while its config names another one, is a failure mode that would be invisible for hours, and I would rather the tool stop than lie about which card it is wearing.

**Which one to read first.** Read `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and is what the discrimination run measured. What a card can change is the sounds-like half described at the top of this page, and I am not going to argue that again here; the numbers are in [MEASUREMENTS.md](MEASUREMENTS.md).

**The shortest way to see what a card does without writing one.** [demo/README.md](demo/README.md) — every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice.

**One exception in that directory.** `card/demo/handoff.effigy` is a hypothesis rather than a voice: it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered — `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in. That is my omission from an earlier draft, and it would have cost somebody a confusing session.

**Key takeaways**

- **`.effigy` is read directly** — no build, no Python.
- **`--card`, `COPE_CARD`, `--rules`** are the three ways in.
- **An unresolvable name errors** rather than falling back silently.
- **`card/demo/handoff.effigy`** is a hypothesis, not a voice.

Swap the file, swap the register.

---

## 🔁 What the hooks add that a file cannot

**The mechanism, and its honest size.** `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` reads that rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts; it falls back to the standing `CONTINUE TEST` when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. So the mid-session text is chosen from measured output rather than fixed in advance, which is the one thing a pasted `CLAUDE.md` cannot do — though I want to hedge that properly rather than let it stand: it is one mechanism, not a guarantee, and the A/B in this repo does not separate the refresher from no refresher.

**A note on the shape of the state.** No prose is stored in the rolling record — turn count, characters, and which rules fired over a 20-turn window, and nothing else. It is a bit like a tally on the back of a door rather than a diary, and the analogy breaks down at the violations log, which does keep the words; that file gets its own section below and I should have mentioned the distinction sooner.

Scored after, reminded during, and both of them optional.

---

## 🧱 Why effigy notation, and what else is out there

**Three of effigy's blocks do what a prose gate needs.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label — `POSTPROC` is regex rules with a `warn` action applied after generation, `WRONG` holds an anti-pattern beside its replacement, and `TEST` holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it. That last one is the reason the notation was worth borrowing rather than inventing something, and I should credit it more plainly than a subordinate clause allows.

**[basanite](https://github.com/justinstimatze/basanite) is the same problem answered the other way round**, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead — lemma frequency against a baseline over real transcripts, so it reports what you have actually been leaning on lately and leaves the judgement to you; its own README calls that awareness rather than prohibition. Which one fits is a question about mood more than about correctness — a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose, different hooks and no shared state, and running both is reasonable.

**[humanizer](https://github.com/blader/humanizer) rewrites where cope scores.** It is a skill, by a different author, that rewrites AI-sounding prose against 35 patterns taken from Wikipedia's "Signs of AI writing", the page WikiProject AI Cleanup maintains — humanizer is called on a text and hands back a rewrite; cope fires at a hook, scores what was already written, and edits nothing. Its pattern list is the wider one, and a reader who wants a rewrite rather than a score should go there. ⚠️ The formatting patterns are where the two disagree on purpose: cope's `bold_label` rule banned humanizer's bold mini-headings until 52 blind pairs put bold and bullets among the three things that decided a reply for this repo's reader, and the rule was deleted rather than tuned. I got that one wrong first, and deleting it was the correction.

**[caveman](https://github.com/JuliusBrussee/caveman) is a fourth axis entirely.** A separate project, by a different author, that compresses agent replies to cut output tokens — cope shapes prose, basanite tracks vocabulary, humanizer rewrites, caveman shortens. Worth naming because a reader wanting fewer tokens rather than different structure should go there instead, and sending them round this page first would have wasted their afternoon.

**Key takeaways**

- **effigy** supplies `POSTPROC`, `WRONG` and `TEST`.
- **basanite** measures where cope bans, and the two compose.
- **humanizer** rewrites where cope scores, on a wider pattern list.
- **caveman** shortens where cope shapes.

Four tools, four different complaints about the same prose.

---

## 📋 The rules

**Grouped by axis, not by implementation.** That is a deliberate choice and I should defend it before the lists rather than after: where a rule lives is an accident of what a regex can express, and what a rule catches is not.

### Voicing rules that ship in the card

- **`flip`** — `warn`. The not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- **`load_bearing`** — `warn`. Reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- **`worth_noting`** — `warn`. Announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

### Voicing rules compiled into the binary

- **`clause_symmetry`** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **`apology`** — the reply performs contrition instead of stating the correction and moving on.
- **`self_postmortem`** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **`announced_length`** — the reply announces its own length rather than cutting it.
- **`cross_turn_repeat`** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.
- **`repeated_opening`** — three or more sentences in one reply opening on the same two words.
  - `cross_turn_repeat` reads the session window for a construction reused across turns; this one reads a reply against itself.
  - Two is left alone, because two is a rhythm.
- **`fragment_run`** — three consecutive sentences of five words or fewer with no finite verb in any of them.
  - One fragment is emphasis, and this repo's own register is full of them.
  - A run of three is the staccato blind judges read as generated.
  - Neither clipped demo card trips it, so neither declines it; a card that wants the run says so with `@gate`.

### Structure rules, compiled, the same under every card

- **`labelled_opening`** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **`paragraph_uniformity`** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **`ask_not_last`** — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **`dangling_end`** — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **`buried_decision`** — an open problem landing after the last question or offer, burying the decision point above it.
- **`forked_end`** — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **`unverified_done`** — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **`loop_ask`** — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.
- **`echoed_heading`** — a heading of two or more content words whose first sentence below repeats every one of them, spending a line to say what the heading already said.

---

**What that grouping implies about the implementation.** A `POSTPROC` pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why cope's shipped card carries only three of them. A reader arriving here expecting a long list of banned phrases should know the list lives in another tool on purpose — basanite measures vocabulary and reports the drift, and cope is not trying to be that. I should have set that expectation two sections earlier and did not.

**Lanes, which are the one place the structure rules vary.** Not by card — by who is going to read the turn. Three of them, and I will take them in order.

| **lane** | **chosen when** | **what changes** |
|---|---|---|
| **interactive** | any turn that is not a loop turn | nothing; the full set |
| **loop** | the prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself | drops four, adds two |
| **external** | `--pretool` scores prose an external write is about to post | drops four, adds nothing |

**The interactive lane is the default** because somebody is waiting at a terminal and the ending is where they decide what happens next. **The loop lane** drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`, and adds `unverified_done` and `loop_ask`, because nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself — what replaces them is the claim check, since a report saying the work is done has to say what it ran. **The external lane** drops the same four and swaps nothing in, because a ticket has a reader and no ending they can answer; it is read days later by somebody who was not in the session, which is the condition every rule that survives the drop was written for. That lane is what the `PreToolUse` entry uses when it scores the description, body or content field an external write is about to post, matched against the Linear save tools in the settings block — warn-only, returning `additionalContext` and never a `permissionDecision`, so the call goes through and the model learns what the prose scored.

---

**Clustering, which is what you see at the bottom of a report.** ⚠️ Every rule fires on its own and knows nothing about the others, so a report is otherwise a flat list a reader works down hit by hit — and three hits across three paragraphs are three small edits, where three inside one paragraph are one paragraph to write again. Printed by the `Stop` hook, by `--check` and by `--pretool`, there are two conditions: **breadth**, when three or more distinct rules land on one paragraph, and **density**, when one rule lands on one paragraph three or more times. 🎯 When both hold, breadth wins, since naming three rules tells a reader to rewrite the block rather than hunt one construction.

> **Key insight:** the density half has a measurement behind it. `--check` over 107 tracked documents produced 114 `flip` hits of which seven were worth changing, and every one of the seven was visible as three in a paragraph rather than as anything about the form.

**Why three and not two.** Three is the floor on both conditions because two rules on a paragraph is ordinary and two hits of one is a coincidence a reader can see unaided. Nothing about what fires or how it is scored changes — clustering is a reading aid on top of the same hits, and I should have said that first, because a reader who thought it changed the scoring would have mistrusted the whole report.

**Key takeaways**

- **Voicing rules** split between card regexes and compiled checks.
- **Structure rules** are compiled and identical under every card.
- **Lanes** vary the structure set by reader, not by card.
- **Clustering** reads several hits together and changes no score.

Grouped by what they catch, implemented where they had to be.

---

## 🚩 Flags

| flag | default | what it does |
|---|---|---|
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | `(empty)` | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
| `--ab-report` | `(empty)` | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | `(empty)` | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is `reject` (default warn-only) |
| `--card` | `(empty)` | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--card-from-sample` | `(empty)` | print a prompt for writing a card from this writing sample; `-` reads stdin |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | `(empty)` | score a prose file against the card and exit; `-` reads stdin |
| `--check-lane` | `(empty)` | score `-check` in the given lane: `interactive` (default), `loop`, or `external` |
| `--describe` | `false` | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | `false` | `MessageDisplay` entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | `false` | read prose on stdin and print it as `--display` would rewrite it |
| `--dry-run` | `false` | with `--setup`, print what would change and write nothing |
| `--inject` | `false` | print the card as prompt text for a `SessionStart` hook |
| `--log` | `$HOME/.local/state/cope/violations.jsonl` | append violations here; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--output-style` | `false` | write the card to `~/.claude/output-styles` as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | `(empty)` | directory to write the output style into (default `~/.claude/output-styles`) |
| `--pretool` | `false` | `PreToolUse` entry: score the prose an external write is about to post, warn-only |
| `--refresh-every` | `30m0s` | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | `false` | `UserPromptSubmit` entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | `(empty)` | print the mid-session reminder one variant would inject, and exit |
| `--render-for` | `(empty)` | comma-separated rule ids to render `-render-arm` against |
| `--render-lane` | `(empty)` | render `-render-arm` as the given lane sees it: `interactive` (default) or `loop` |
| `--rules` | `(empty)` | read the card from this `.effigy` or `.json` file instead of the one built into the binary |
| `--setup` | `false` | emit the output style and wire the hooks, then print the one step left |
| `--version` | `false` | print version and exit |

Everything the binary takes, with the defaults it takes them at.

---

## 💾 What lands on disk

**Three files, all mode `0600`, and one of them quotes you back.** I want that second clause to arrive before the table rather than inside it, because a tool that keeps your prose should say so plainly and I have seen enough READMEs bury it in a schema.

| **path** | **holds** |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window; no prose is stored, only rule names and counts |

> **Note:** the violations log quotes replies back — the matched span with context either side. Setting `--log` to empty disables it, and I should have named that escape hatch in the same breath as the warning.

Two files of counts and one file of your own sentences.

---

## ✏️ Editing the card

**Where the grammar lives and what regenerates.** effigy owns the `.effigy` grammar; `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`; `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart. That last one is the guard I would keep if I could keep only one, and I apologise for not making it the first sentence.

**The `NEVER` budget, and the thing about it people get wrong.** The budget is 10, and anything over it is reported at load rather than dropped silently — a rule that vanishes without a word is worse than a rule that fails loudly. ⚠️ The budget is charged against each injection separately, not against the card file: the always-on rules and the evidence-gated ones render on different paths, no code path renders their union, and so a card may hold more `NEVER` rules in total than the budget and still be perfectly healthy. `never_rules_over_budget` in the introspected facts is the authoritative list of rules that really are discarded unrendered, and it is empty when the card is healthy.

---

**The two card-authored forms, spelled out.** Not approximated — a reader writing a card needs the selectors and predicates exactly, and an approximation of them is worse than none.

- **`@gate <rule_id> off — <why>`**, one per line in the card header.
- **`@shape <id>: <selector> <predicate> — <why>`**, one per line in the card header.
  - selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
  - predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**What the loader checks and what it refuses.** A rule id has to be one the gate has for `@gate`, and must not collide with one for `@shape` — a wrong id in either is reported at load rather than ignored, which is the same principle as the budget and I should have grouped them. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

**Two consequences worth having.** A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught — the decline is a scoring decision and not a blindfold. And a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies, which is the point of letting a card write one at all.

**Key takeaways**

- **`make rules` regenerates**, `make check-rules` stops the drift.
- **The `NEVER` budget is 10**, charged per injection, loud when exceeded.
- **`@gate` declines**, `@shape` asserts, and both require a reason.
- **A declined rule still runs** and still shows up in a backfill.

Edit the file, regenerate, and let CI catch the drift.

---

## 📏 Calibrating it on your own transcripts

**How the rules were chosen, and how you would re-choose them.** `cope-gate --backfill` scores a whole session transcript at once and is how every rule on this page earned its place; `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is not the same thing as one per project and I should have named that distinction the first time somebody read a sweep as a per-project survey.

**The metric worth watching.** Hits-per-character, rather than the share of turns hit — the second number tracks how long the turns were, so a session of essays will look worse than a session of one-liners for no reason a rule could act on. It is a bit like judging a proofreader by pages marked rather than by errors found, though the analogy breaks down, since a proofreader is not also choosing the page length.

For the rates themselves, [MEASUREMENTS.md](MEASUREMENTS.md) has what was run and on how much text.

---

## ⚖️ Known limits

**Organised by the two axes, because that is where the limits fall.** Three of them, and I would rather state them here than have a reader find them by being disappointed.

**First, on the structure side, two rules are less clever than their names suggest.** `labelled_opening` is not a tagger — it matches a shape, not a part of speech, and it will miss and over-fire accordingly. `ask_not_last` says nothing about the order of several asks; it knows that one is early and the reply carried on past it, and nothing more. I should have said that in the rules list itself.

**Second, the hit rate is roughly four fifths structure**, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — so it is a description of the output and not a judgement of it. The judgement lives in the discrimination test, and that test covers voicing only. Reading the structure share as a quality score is the misreading I most expect, and I apologise for not flagging it earlier on the page.

**Third, and largest, is the boundary the opening section already named.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. Both directions are also the card marking its own homework: a decline lowers that card's score and an assertion raises it, both are worth reading with the reason attached, and that is exactly why the syntax requires one. The numbers behind all of this are in [MEASUREMENTS.md](MEASUREMENTS.md).

**Key takeaways**

- **Two structure rules** are shape-matchers, not parsers.
- **The structure share** describes the corpus, not the quality.
- **The card's reach** stops at words, sentences and questions.
- **A card scores itself** in both directions, with a stated reason.

Known, stated, and measured where measurement was possible.

---

## 🗂️ Layout

| **path** | **what** |
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

Eleven paths, and the two at the top are the whole card.

---

*This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.*

MIT. justin@justinstimatze.com
