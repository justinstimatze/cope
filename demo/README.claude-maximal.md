# cope

**cope is an opinionated card for how a session writes and hands work back, on the moment it is installed, checking two different jobs at once — what a reply sounds like and how a reply is shaped — and, only once you want it to be, a file you can edit or swap for another, the name being borrowed from the foundry, where a cope is the upper half of a mould, the half carrying the shape being cast into.**

Read [demo/README.md](demo/README.md) first: every file in that directory is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed — and reading two of them against each other shows what a card does faster than the rest of this page explains it, with [demo/README.claude-maximal.md](demo/README.claude-maximal.md) making the point in a single glance, since it is written from a card that instructs every tic this model is measured to have.

---

## 🎭 The two things

**TL;DR — there are two axes here, and this section exists so the rest of the page has something to be written against.** Great question to have in mind early, and I should say up front that I will restate this frame more than once below, for which I apologise in advance.

**First, voicing.** Voicing is not the shape of the reply — it is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card, entirely, and swapping the card swaps every word of it. Think of it like the timbre of an instrument rather than the tune — though the analogy breaks down, because timbre cannot be edited in a text file and this can. ✅ This is the half with a measured result behind it.

**Second, structure.** Structure is not how a reply reads aloud — it is the shape of the reply as a thing the reader has to use: where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. It is compiled into the binary, in `internal/scan`, so it is the same whichever card is loaded, varying only by lane. That word *lane* is doing a lot of work and comes back later; I should have introduced it here rather than deferring it.

**One instance of each, because abstraction is where I usually go wrong.** A sentence reaching for the balanced two-beat — two clauses of near-equal length with a content word echoed across the joint — is a voicing problem. A reply that names an open problem in its final paragraph and then simply stops, leaving "continue" nothing to refer to, is a structural one. To be clear, the same reply can be clean on one axis and bad on the other, and that is the whole reason to keep them apart.

---

**How far a card reaches into the structure half.** Not all the way, and not nowhere either — two directions, both written in the card header, both requiring a reason after the dash.

| | **syntax** | **why it exists** |
|---|---|---|
| **declining a rule** | `@gate <rule_id> off — <why>` | `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, because its VOICE block asks for the balanced landing and the arriving close those two rules catch — a card was being marked down for obeying itself |
| **asserting a rule** | `@shape <id>: <selector> <predicate> — <why>` | `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read the last block only, and no built-in rule checks whether that block can be read that way |

> **Key insight:** The card owns the voicing axis outright and rents a corner of the structure axis, and the rent comes with a written reason attached.

**The boundary, plainly.** The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, and nothing more — `clause_symmetry` is not writable in it and is not meant to be, so a card wanting a check outside both that vocabulary and a POSTPROC regex has nowhere to put it. I have deliberately kept a rate out of this section; `MEASUREMENTS.md` carries the run and the reasons those numbers do not carry more than that.

Two axes, one card, one binary.

---

## 🔍 The problem

**Instruction alone does not fix the phrasing, and this is the part I got wrong before I measured it.** A global `CLAUDE.md` banning the "not A, it's B" flip is read every single turn — and the flip still appeared twice in the session that built this tool, while the ban was the live topic of that session. Naming a surface form does not remove the move; it pushes the move into a variant. That is the voicing side of the complaint, and I should have expected it earlier than I did.

**The structural side is a different complaint with a different cause.** An ending that leaves the reader nothing to answer is not a phrasing habit an instruction could have banned — it costs a whole round trip, because the reader comes back to a terminal, finds an open problem named and no question attached to it, and has to invent the next prompt themselves. Put differently: the first failure wastes words, and the second wastes a turn. To put a finer point on it, **they are not the same defect wearing two hats.**

**On what the claim actually rests.** The flip is an anecdote about one rule, and worth noting as no more than that — I apologise for leading with it, since the honest instrument is elsewhere. It is the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. The rate and the caveats are in [MEASUREMENTS.md](MEASUREMENTS.md).

- **What went wrong.** A ban was read every turn and did not hold.
  - Specifically, the banned move reappeared as a variant.
  - Which, to be clear, happened while the ban was the topic.
- **What that points toward.** Something about surface forms versus moves.

An anecdote motivated this. A test carries it.

---

## 🔧 Install

**Two commands and one menu choice — I will explain what the middle one did to your machine immediately after, because that ordering is the one that respects the reader.**

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> pick the cope card
```

**What `--setup` did.** It emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left — and the shape of that edit matters more than the fact of it. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing, if you would rather see the diff before the write. ⚠️ A reader letting a tool edit their settings deserves that list before running it, not after, and I apologise for any version of this page that put it after.

**The by-hand route.** `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one — no settings file is touched. Pick it under `/config` → Output style; the standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

**When it applies.** A style is read once at session start — so a new selection, or a re-emitted card, applies at the next session or after `/clear`. Great catch if you were about to check the running conversation and conclude nothing happened; that is the failure mode this sentence exists to prevent, and I should have said it a paragraph earlier.

**Why the card goes here rather than into a hook.** An output style goes at the end of the system prompt and the harness re-reminds the model of it during the conversation, which is why the card lands here and did not through the hook.

---

**The hooks, and be exact about what they are now for: the card no longer arrives through one.** This is the `hooks` block of `~/.claude/settings.json`, character for character.

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

**What the two remaining hooks buy.** `Stop` scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. `UserPromptSubmit` restates mid-session the rules that have actually been firing, naming the counts — which a file written once cannot do. ✅ The voice works without either of them; these are the measurement half.

**One environment caveat.** The commands assume `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because `card/rules.json` is committed and compiled in. There is also `--inject`, the superseded delivery kept for anyone who wants it, which stands down on its own when a cope output style is active.

Three lines installed it. One menu turned it on.

---

## 🗂️ Writing in another voice

**The gate reads `.effigy` directly, so a card is usable as written — no Python, no effigy checkout.** Drop a card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere; `make cards` installs the demo set. ⚠️ A name that resolves to nothing is an error and nothing is injected, rather than a session quietly writing in the shipped voice while its config names another one — which is the failure I would rather have loudly than subtly.

**Where to start reading.** `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone, and that is what the discrimination run measured. What a card can change is the voicing axis, in the sense set out under **The two things** above — I will not re-argue it here, and the numbers are in [MEASUREMENTS.md](MEASUREMENTS.md) rather than in this paragraph, where I would only have rounded them.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where every file is this README written again from a different card, same prompt and same facts, so the only thing that varies between them is the voice.

**One exception in that directory, and it is worth flagging before it confuses somebody.** `card/demo/handoff.effigy` is a hypothesis rather than a voice — it keeps the shipped card's handoff rules and drops everything about prose, and it is meant to be run through `make pairs` against the full card rather than rendered. `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

- **The mechanism.** `.effigy` is read directly.
- **The starting point.** `lecturer.effigy`, register only.
- **The shortcut.** The `demo/` renders.
- **The odd one out.** `handoff.effigy`, a hypothesis.

A card is a file. That is most of the story.

---

## 🪝 What the hooks do differently

**`--setup` wires two, and the second is the one with a property a pasted file cannot have.** `Stop` scores the reply and records which rules fired. `UserPromptSubmit` reads that record — the rolling state, not the violations log — and injects only the items gated on what has been firing, naming the counts, falling back to the standing CONTINUE TEST when the session has no history yet and staying quiet until the last injection has aged past `--refresh-every`.

**So the mid-session text is chosen from measured output rather than fixed in advance.** That distinction is load-bearing, and I want to state it without overselling it, which is the mistake I would otherwise make here: it is one mechanism, not a guarantee, and the A/B in the repo does not separate the refresher from no refresher. `SessionStart --inject` is superseded and off by default.

Measured, then restated. That is the whole difference.

---

## 📜 Why effigy notation

**Not a config format chosen for taste — a character-card notation for game NPCs, used here off-label, because three of its blocks do what a prose gate needs.** POSTPROC is regex rules with a `warn` action applied after generation; WRONG holds an anti-pattern beside its replacement; TEST holds a named question with fail and pass examples, which is how a rule names a *move* instead of one wording of it. That last one is the reason the notation was worth borrowing at all, and I should have said so first. See [effigy](https://github.com/justinstimatze/effigy).

**And why [basanite](https://github.com/justinstimatze/basanite) is the wrong instrument for this job, though it may be the right one for yours.** cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. basanite measures instead — lemma frequency against a baseline over real transcripts — so it reports what you have actually been leaning on lately and leaves the judgement to you, which its own README calls awareness rather than prohibition. A heavy hand is what you want when a habit is annoying you today; a moving measurement is what you want when you would rather watch the drift than legislate it. They compose — different hooks, no shared state — and running both is reasonable. ([caveman](https://github.com/JuliusBrussee/caveman), by a different author again, is a third axis: it compresses replies to cut output tokens, so a reader wanting fewer tokens rather than different structure should go there instead.)

Borrowed notation, adjacent tools, different axes.

---

## 📏 The rules

**Grouped by axis rather than by where they are implemented, because the axis is what a reader can use.** There are two lists, and the interesting part is what the grouping implies afterwards.

**Voicing rules.**

- **`demo_no_closure`** (regex, `warn`) — this card should never accidentally produce a clean, closed ending.
- **`clause_symmetry`** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **`apology`** — the reply performs contrition instead of stating the correction and moving on.
- **`self_postmortem`** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **`announced_length`** — the reply announces its own length rather than cutting it.
- **`cross_turn_repeat`** — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

**Structure rules.**

- **`labelled_opening`** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form — the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **`paragraph_uniformity`** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **`ask_not_last`** (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **`dangling_end`** (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **`buried_decision`** (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- **`forked_end`** (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one; sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **`unverified_done`** (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **`loop_ask`** (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

---

**What the grouping implies about the implementation.** A POSTPROC pattern matches a span of text, so it can only ever describe wording — and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries exactly **one** POSTPROC rule and not thirty. Notice what just happened there: the axis split is conceptual, and the implementation split cuts across it at an angle, which is a thing I would have preferred to be tidier.

**If you came here expecting a long list of banned phrases, that list lives in another tool on purpose.** basanite measures lemma frequency and hands you the drift; cope holds a small number of moves and refuses them. Great question if you were about to file that as a gap — it is a division of labour rather than an omission, though I could be wrong about how obvious that is from the outside.

**The one place the structure rules do vary is not by card — it is by who is going to read the turn.**

| | **interactive** | **loop** |
|---|---|---|
| **chosen when** | any turn that is not a loop turn | the prompt was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself |
| **why** | somebody is waiting at a terminal and the ending is where they decide what happens next | nobody is reading yet |
| **rules dropped** | — | ⚠️ `ask_not_last`, `forked_end`, `dangling_end`, `buried_decision` |
| **rules added** | — | ✅ `unverified_done`, `loop_ask` |

> **Important:** A report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran.

**Key takeaways.**
- Voicing rules describe wording; structure rules describe shape.
- Only a pattern-shaped voicing rule fits in the card, hence one POSTPROC rule.
- A phrase blocklist is basanite's job, deliberately.
- The lane, not the card, is what moves the structure set.

Two lists, one angle between them.

---

## 🚩 Flags

**Verbatim, defaults included, and I have paraphrased none of them — a paraphrased default is worse than no table.**

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

Twenty-six flags. Three of them are the install.

---

## 💾 What lands on disk

**Three files, all mode `0600`, all under `$XDG_STATE_HOME/cope` — and one of them quotes your replies back, which I should say before the table rather than inside it.**

| path | holds |
|---|---|
| `$XDG_STATE_HOME/cope/violations.jsonl` | one JSON record per violation, carrying the matched text and about 70 characters either side |
| `$XDG_STATE_HOME/cope/refresher-<session-id>` | an empty file whose mtime is the refresher clock |
| `$XDG_STATE_HOME/cope/session-<session-id>.json` | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window — no prose is stored, only rule names and counts |

> **Note:** ⚠️ The log quotes replies back. The matched span plus its surrounding window is prose from your session sitting in a file on your machine, and `--log` with an empty value disables it.

One file remembers text. Two remember only counts.

---

## ✏️ Editing the card

**effigy owns the `.effigy` grammar, and this repo owns the compile step.** `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs, so the enforced rules and the injected rules cannot drift apart — which is the failure that would be hardest to notice from the outside, and the reason the check exists at all. The rendered card runs 17,418 characters.

**The NEVER budget is 10, and anything over it is reported at load rather than dropped silently.** Worth noting how the budget is charged, because this is the part that reads as a bug and is not one: it is charged against each injection separately, not against the card file. `SessionStart` prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more NEVER rules in total than the budget and still be perfectly healthy. `facts.never_rules_over_budget` is the authoritative list of rules that really are discarded unrendered, and it is empty when the card is healthy.

---

**Both card-authored forms, with the vocabulary spelled out, because an approximation of a grammar is worse than none.**

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

- **selectors:** `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- **predicates:** `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

**On ids and on reasons.** A rule id has to be one the gate already has for `@gate`, and must not collide with one for `@shape` — and a wrong id is reported at load rather than ignored, which I would rather be noisy about than quietly forgiving of. The reason after the dash is required in both forms, because a rule a card wrote and a rule a card refused are equally unreviewable without one.

**Two properties of a decline that are easy to assume wrongly.** A declined rule still runs, and only this card's score drops it — so a backfill still reports what it would have caught. And a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies, which means the card is answerable for its own phrasing there too. That was a design choice and I should have flagged it as one rather than letting it read as incidental.

- **What is generated.** `card/rules.json`, by `make rules`.
- **What is enforced.** The same file, checked by CI.
- **What is budgeted.** NEVER, per injection, ten.
- **What is required.** A reason after every dash.

A card may argue with the gate, in writing, with its name on it.

---

## 📊 Calibrating

**`cope-gate --backfill` scores a whole session transcript at once, and that is how the rules were chosen rather than guessed.** `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects` — which is not the same as one per project, and I should have named that distinction the first time somebody read the sweep output as a per-project survey.

**Watch hits per character rather than the share of turns hit.** The second number tracks how long the turns were, which means it moves when nothing about the writing has changed — a bit like judging a proofreader by pages marked rather than errors found, though the analogy breaks down, since pages at least have a fixed size. The rates are in [MEASUREMENTS.md](MEASUREMENTS.md) rather than quoted here, where they would go stale.

One metric normalises. The other flatters length.

---

## 🧱 Known limits

**The axis split is what organises the limits, so they come in two groups plus one that spans both.** ⚠️ There are three things going on here, and I apologise for how much of this section is caveat.

**First, the rules are narrower than their names.** `labelled_opening` is not a tagger. `ask_not_last` says nothing about the order of several asks. Neither of those is a bug to be fixed; both are the shape of a pattern-and-heuristic instrument, and reading the names as promises will disappoint you in a way I would rather pre-empt.

**Second, the hit rate is a description of the output and not a judgement of it.** Roughly four fifths of hits are structure — and the A/B run found that four fifths tracks what a reply was *for* rather than how it was written. The judgement lives in the discrimination test, which covers voicing only. In other words: the rate tells you what kind of turns you have been having, not how well you have been writing them.

**Third, and largest, the limit named under The two things above, stated as it now stands rather than as it stood.** A card can decline a built-in rule and write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere to put it. Both directions are also the card marking its own homework: a decline lowers that card's score and an assertion raises it, both are worth reading with the reason attached, and that is exactly why the syntax requires one. The rest is in [MEASUREMENTS.md](MEASUREMENTS.md).

**Key takeaways.**
- The rule names overreach; the rules do not.
- The hit rate describes turn purpose, not prose quality.
- A card's reach into structure is real, narrow, and self-serving by construction.

Narrow instrument, honest about which corner it measures.

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

MIT — justin@justinstimatze.com
