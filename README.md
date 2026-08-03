# 🏺 cope

**What this is.** cope ships an opinionated card that is on the moment it is installed — for most readers that card *is* the product — and it scores what comes back on two different jobs rather than one, how a reply sounds and how a reply is shaped; the card is a file, so it can be edited or swapped, which is where the name comes from, a cope being the upper half of a foundry mould, the half carrying the shape being cast into.

**Start here, not here.** Read [demo/README.md](demo/README.md) — every file in that directory is this README written again from a different card, same prompt and same facts, so the card is the only thing that changed, and reading two of them against each other shows what a card does faster than the rest of this page explains it; if this page is hard going, [demo/README.claude-voice.md](demo/README.claude-voice.md) is the same page in the register cope actually ships.

---

## 🎭 The two things

**TL;DR for this section.** There are two axes here, they fail independently, and the whole page is written against that split — I apologise for front-loading a frame before the install instructions, but the sections below make no sense without it.

**Voicing is the first axis.** Not what a reply *does* — what it sounds like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely — VOICE, TRAITS, NEVER, WRONG, MES and POSTPROC — so swapping the card swaps every word of it. To be clear, that word *entirely* is doing a lot of work: this is the half of the tool where the file you install is the whole answer, and it is also the half with a measured result behind it. A sentence that reaches for the balanced two-beat — two clauses of near-equal length with the same content word on both sides of the comma — is a voicing failure, and no amount of restructuring fixes it.

**Structure is the second axis, and it is a different complaint with a different cause.** Not how the reply reads, but the shape of it as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. That half is compiled into the binary — `internal/scan` — so it is the same whichever card is loaded, varying only by lane. A reply that names an open problem in its last paragraph and then stops, leaving "continue" nothing to refer to, is a structural failure, and it can be written in immaculate prose. I should have said this earlier: the same reply can be clean on one axis and bad on the other, and that is the entire reason to keep them apart.

| | **voicing** | **structure** |
|---|---|---|
| **lives in** | ✅ the card | ✅ the binary |
| **swaps with the card** | yes | no |
| **example failure** | the balanced two-beat | ⚠️ an open problem, then silence |

---

**How far a card reaches into the second half.** Not nowhere, and not all the way — this is the part that most recently changed, so it is worth being exact rather than gestural. A card can move in two directions, and both are written in the card header, one rule per line:

- **Declining a built-in rule.** `@gate <rule_id> off — <why>`
  - `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch.
    - Which is to say: the mechanism exists because a card whose VOICE asked for something a built-in rule catches was being marked down for obeying itself.
- **Asserting a rule of its own.** `@shape <id>: <selector> <predicate> — <why>`
  - `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way.
    - The 60 is measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.
    - Which is to say: the mechanism exists because a card's own commitment about how a reply ends had nowhere to be checked.

**The boundary, stated flat.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — it cannot express what the compiled rules express, `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it. I have deliberately left a rate out of this section; [MEASUREMENTS.md](MEASUREMENTS.md) has the runs and the reasons their numbers do not carry more than they do.

> **Key insight:** One axis is a file you swap. The other is a binary you don't.

**Key takeaways.**
- Voicing is sound, it lives in the card, and swapping the card swaps it.
- Structure is shape, it lives in the binary, and it is the same under every card.
- A card can decline a compiled rule and assert one of its own — with a reason, both times.

Two axes. One of them is yours.

---

## 🧩 The problem

**The short version, which I will then undercut.** Instruction alone does not fix the phrasing, and the evidence for that is embarrassing rather than theoretical — a global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn, and the flip still appeared twice in the session that built this, while the ban was the topic of the session. That is my mistake in the most literal sense available, and it is also the finding.

**Why naming a form does not remove the move.** Great question to ask of any instruction-based fix — the answer is that naming a surface form pushes the move into a variant. Ban the flip and you get the flip with a semicolon, the flip with "rather than", the flip split across two sentences. The instruction is read, the instruction is obeyed at the level of the string, and the underlying gesture survives intact. To be clear, the ban was not ignored; it was routed around.

**The structural failure is not the same complaint.** An ending that leaves the reader nothing to answer costs a whole round trip — the reader types "continue", the reply has to guess what "continue" meant, and a turn is spent recovering the decision that should have been on the page. Worth noting: that is not a phrasing habit an instruction could have banned, because nothing in the sentence is wrong. The seam is between "the prose is good" and "the reply is usable", and I should have separated those two things much earlier than I did.

---

**What the claim actually rests on.** The flip is an anecdote about one rule, and I apologise for having led with an anecdote — it is vivid and it is not the evidence. The evidence is the blind discrimination test: a reader is shown only a voice's own description of itself, then picks which of two replies was written under it. That instrument has a result behind it, and the rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md) rather than here, where they would rot.

An anecdote opened this section. A measurement is what holds it up.

---

## 📦 Install

**Two commands and one menu choice.** That is the whole install — and the order matters more here than anywhere else on this page, so I will give it as a block first and explain afterwards rather than the other way round.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

**Then pick the style.** Under `/config` → **Output style**, select the card. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way; the same thing can be set as `"outputStyle"` in `.claude/settings.local.json` if you would rather commit it than click it.

**What `--setup` did, since it edited your settings.** Not a silent rewrite — a narrow one, and you deserve the shape of it before you run it rather than after. It emits the output style, wires the two hooks into `~/.claude/settings.json` with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing, which is the version to run first if any of the above made you pause.

**If you would rather nobody wrote to your settings.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` and stops there; `COPE_CARD=<name>` in front of it emits a different one. ⚠️ A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after `/clear` — a reader who picks a style, sees no change in the running conversation, and concludes the tool is broken has been failed by this paragraph, and I should have put that sentence higher.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation — which is why the card lands from here and did not from the hook.

---

**Now the hooks, which are no longer how the card arrives.** This is the block `--setup` writes into `~/.claude/settings.json`, reproduced exactly; `SessionStart` is absent on purpose, and adding it back is not an improvement.

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

**What the two remaining hooks buy.** `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` restates mid-session the card items gated on what has actually been firing, naming the counts — which a file written once cannot do. To be clear, and this is the part I most risk overselling: the voice works without either of them, and these two are the measurement half.

**One PATH caveat, because it is the failure that looks like nothing.** Both commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in — a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout at all, because `card/rules.json` is committed and compiled in. There is also `--inject`, the superseded delivery that remains for anyone who wants it and stands down on its own when a cope output style is active, and I will say nothing further about it.

**Key takeaways.**
- Install, `--setup`, then pick the style under `/config`. ✅
- A style is read at session start, so it applies next session or after `/clear`. ⚠️
- The two hooks are measurement, not delivery. 🎯

Three steps. The card is the one that matters.

---

## 🗣️ Writing in another voice

**The capability the first sentence promised.** The gate reads `.effigy` directly, so a card is usable as written — no Python, no effigy checkout, no build step between writing a register and running under it. I apologise for how much of this page came before this section; this is the part most readers who stay will want.

**Where a card goes and how it is named.** A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. Worth noting: a name that resolves to nothing is an error and nothing is injected, rather than a session quietly writing in the shipped voice while its config names another one — a silent fallback there would be the worst available behaviour, since the failure would be invisible in exactly the output you were trying to change.

**Which one to read first.** `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone and it is the card the discrimination run measured. What a card can change is the voicing axis, as **The two things** above set out, and I will not argue that again here — the numbers are in [MEASUREMENTS.md](MEASUREMENTS.md) rather than in this paragraph.

**The shortest route to understanding a card without writing one.** [demo/README.md](demo/README.md) — every file under `demo/` is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice.

**One exception in that directory.** `card/demo/handoff.effigy` is a hypothesis rather than a voice: it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

A card is a file. That was always the interesting part.

---

## 🪝 What the hooks do differently

**The claim, stated before I qualify it.** A pasted `CLAUDE.md` is fixed in advance; the mid-session text here is chosen from measured output. That is one difference, and it is the only one I would defend.

**`Stop` records.** It scores the reply just written and appends which rules fired to the session's rolling state — rule names and counts, no prose. Think of it like a turnstile counter at a station: it does not know who went through, only how many and through which gate. The analogy breaks down immediately, of course, since a turnstile does not read the passenger's sentences to decide which gate they used.

**`UserPromptSubmit` chooses.** It reads that record — not the violations log — and injects only the card items gated on what has been firing, naming the counts, falling back to the standing CONTINUE TEST when the session has no history yet, and staying quiet until the last injection has aged past `--refresh-every`. To be clear, and I should have led with this rather than with the mechanism: it is one mechanism, not a guarantee, and the A/B in this repo does not separate the refresher from no refresher.

**And the one that is off.** `SessionStart --inject` is superseded and off by default; `--setup` does not wire it.

Evidence-gated, not fixed in advance. That is the whole difference.

---

## 🧠 Why effigy notation

**Not a config format chosen for taste — a format that already had the three blocks.** [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here off-label, and three of its blocks do exactly what a prose gate needs: POSTPROC is regex rules with a `warn` action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples — which is how a rule names a *move* instead of one wording of it. I should have found that framing sooner than I did.

**Why [basanite](https://github.com/justinstimatze/basanite) is the wrong instrument for this, and the right one for something adjacent.** cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it — whereas basanite measures, lemma frequency against a baseline over real transcripts, reporting what you have actually been leaning on lately and leaving the judgement to you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood more than about correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable. A reader who wants *fewer tokens* rather than different structure wants neither, and should look at [caveman](https://github.com/JuliusBrussee/caveman), a separate project by a different author that compresses agent replies to cut output cost — a third axis again.

Three tools, three axes. Prohibition, awareness, compression.

---

## 📏 The rules

**Grouped by what they check, not by where they live.** That grouping is the useful one, and the implementation falls out of it afterwards rather than organising it — I apologise for asking you to hold the axis split for this long, but here is where it pays.

**Voicing rules.**
- `demo_no_closure` *(regex, `warn`)* — this card should never accidentally produce a clean, closed ending.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure rules.**
- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` *(interactive)* — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` *(interactive)* — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` *(interactive)* — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` *(interactive)* — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` *(loop)* — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` *(loop)* — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

**What the grouping implies about the implementation.** A POSTPROC pattern matches a span of text, so it can only ever describe wording — which means every voicing rule that needed more than a pattern had to be written in Go, beside the structure rules, and that is why the card ships exactly **1** POSTPROC rule rather than a page of them. If you came here expecting a long list of banned phrases, the list lives in another tool on purpose, and that tool is basanite.

---

**The one place the structure rules vary.** Not by card — by who is going to read the turn.

| | **interactive** | **loop** |
|---|---|---|
| **chosen when** | any turn that is not a loop turn | the prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself |
| **why** | somebody is waiting at a terminal and the ending is where they decide what happens next | ⚠️ nobody is reading yet |

**What the loop lane drops, and what it puts in their place.** It drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision`, and adds `unverified_done` and `loop_ask` — because a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran.

Same card, same binary. Different reader, different rules.

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

Twenty-six flags. Three of them are the install.

---

## 💾 What lands on disk

**Three files, all `0600`, all under `$XDG_STATE_HOME/cope`.** I should say the uncomfortable one first rather than third: the log quotes your replies back.

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | ⚠️ one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window — no prose, only rule names and counts |

**The distinction between those last two rows and the first.** The session file is counts; the log is text. `--log` with an empty value disables the log entirely, and a reader who is uneasy about excerpted prose sitting on disk should reach for that before reading further.

One file quotes you. Two only count.

---

## ✏️ Editing the card

**Where the grammar lives.** effigy owns the `.effigy` grammar, `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs so the enforced and injected rules cannot drift apart — that drift is the failure mode worth engineering against, because it is invisible from either side alone.

**The NEVER budget, and what it is charged against.** The budget is **10**, and anything over it is reported at load rather than dropped silently. Worth noting, because the arithmetic surprises people: the budget is charged against each *injection* separately, not against the card file — `SessionStart` prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union, so a card may hold more NEVER rules in total than the budget and still be perfectly healthy. `never_rules_over_budget` in the introspected facts is the authoritative list of rules that really are discarded unrendered, and it is empty when the card is healthy.

---

**The two things a card can write about the rules.** Both go one per line in the card header, and both require the reason after the dash:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

**The `@shape` vocabulary, spelled out, because an approximation of it is worse than none.**
- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**Ids, and what happens when one is wrong.** For `@gate` the id has to be one the gate already has; for `@shape` it must not collide with one. Either way a wrong id is reported at load rather than ignored — silence there would let a card believe it had declined a rule that was still scoring it.

**Why the reason is mandatory in both.** A rule a card wrote and a rule a card refused are equally unreviewable without one, and I should have made that argument before giving the syntax rather than after. Two consequences follow: a declined rule still *runs* and only this card's score drops it, so a backfill still reports what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

**Key takeaways.**
- `make rules` regenerates; `make check-rules` keeps enforced and injected from drifting.
- The NEVER budget is per injection, not per card file. ⚠️
- `@gate` declines, `@shape` asserts, and both owe a reason. 🎯

A card can argue with the gate. In writing, with reasons.

---

## 📊 Calibrating

**How the rules were chosen, since it was not by taste.** `cope-gate --backfill` scores a whole session transcript at once, and that is the instrument the rules were selected against — every assistant turn in a real conversation, not a handful of examples written to make a rule look good.

**Sweeping wider.** `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects` — which, to be clear, is not the same as one per project, and reading it as one-per-project will make the sample look broader than it is. That was my mistake in an earlier reading of the output, and it is worth naming rather than quietly correcting.

**The metric to watch.** Hits-per-character, not the share of turns hit — because that second number tracks how long the turns were, so a session of longer replies looks worse without writing worse. The rates themselves are in [MEASUREMENTS.md](MEASUREMENTS.md), which records what was run, on how much text, and what it said.

Measured on real transcripts. The rates live in one file.

---

## ⚠️ Known limits

**Two rules that are narrower than their names suggest.** `labelled_opening` is not a tagger — it recognises a shape, not a part of speech, and it will miss and over-fire accordingly. `ask_not_last` says nothing about the *order* of several asks; it notices one sitting above prose that carried on past it, and stops there.

**What the hit rate is and is not.** Roughly four fifths of hits are structure, and the A/B run found that four fifths tracks what a reply was *for* rather than how it was written — so the hit rate is a description of the output, not a judgement of it. The judgement lives in the discrimination test, which covers voicing only. I apologise for how easy the aggregate number is to misread; I have not found a framing that makes it harder.

---

**The largest limit is the one **The two things** already named, restated as it now stands.** A card can decline a built-in rule and assert one of its own, and the vocabulary it asserts in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

**And both directions are the card marking its own homework.** A decline lowers that card's score; an assertion raises it. Neither is readable without the reason attached, which is precisely why the syntax requires one — and it is also why a backfill keeps scoring declined rules anyway, so the thing a card refused is still visible to anyone who goes looking. [MEASUREMENTS.md](MEASUREMENTS.md) is where that looking starts.

**Key takeaways.**
- `labelled_opening` recognises a shape, not a part of speech.
- Four fifths of hits are structure, and that tracks purpose rather than quality. ⚠️
- A card grading itself is legible only with its reasons attached.

The limits are the interesting part. They usually are.

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

Eleven paths. Two of them are the card.

---

*This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.*

MIT — justin@justinstimatze.com
