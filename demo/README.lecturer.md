# cope

cope ships one opinionated card — a cope being the upper half of a foundry mould, the half carrying the shape being cast into — and once it is installed that card is on, checking two different jobs rather than one, how a reply sounds and how a reply is shaped, with the card itself being only a file you can edit or swap for another.

The short way in is [demo/README.md](demo/README.md), where every file is this same page written again from a different card, the same prompt and the same facts behind each one, so the card is the only thing that varies between them; two of them read against each other will show you what a card does faster than the rest of this page can explain it. The one that makes the point in a single glance is [demo/README.claude-maximal.md](demo/README.claude-maximal.md), written from a card that instructs every tic this model is measured to have. A reader who takes the link understands the tool. A reader who does not is taking prose about registers on trust.

## The two things it checks

Voicing is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. All of it lives in the card, and swapping the card swaps every word of it, which is also why voicing is the half with a measured result behind it. Structure is something else entirely — the shape of the reply as a thing you have to use, meaning where the decision sits, whether the ending gives "continue" anything to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. That half is compiled into the binary and is the same whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing problem; a reply that names an open problem in its closing paragraph and then stops, leaving you nothing to answer, is a structural one. The same reply can be clean on one axis and useless on the other, and that is the whole reason for keeping them apart.

A card does reach into the structure half, though not far, and this is the part that has most recently changed. It reaches in two directions. A card can decline a built-in rule with `@gate <rule_id> off — <why>`, one per line in the card header, which exists because a card whose VOICE block asked for something a built-in rule catches was being marked down for obeying itself: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason, since its register asks for the balanced landing and the arriving close. A card can also assert a rule of its own with `@shape <id>: <selector> <predicate> — <why>`, again one per line in the header, which exists because a card's own commitment about how a reply ends had nowhere to be checked: `card/demo/handoff.effigy` asserts `readable_cold` — last paragraph words <= 60 — because its peak asks the reader to re-enter cold and read the closing block alone, and no built-in rule asks whether that block can be read that way. The 60 is not a taste; across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

What a card cannot do is write the checks that needed more than counting. The `@shape` vocabulary counts words and sentences and asks whether a block poses a question, so `clause_symmetry` is not writable in it and was never meant to be, and a card wanting a check outside both that vocabulary and a POSTPROC regex still has nowhere to put it.

## Why instruction is not enough

Tell a model never to write a particular sentence and what it hears is the sentence, not the move underneath it. On the machine this was built on there is a global CLAUDE.md banning the not-A-it's-B flip in as many words, and that file is read at the top of every single turn. The flip appeared twice during the session that built cope, while the ban was the topic under discussion. Naming a surface form does not remove the move; it relocates it into a variant. The instruction was obeyed and the habit survived.

The structural complaint has a different cause and cannot be reached the same way. When a reply names an open problem in its last paragraph and stops, nothing in the wording is wrong, and no ban on phrasing could have caught it — what it costs you is a whole round trip, you typing the question the reply should have asked on your behalf.

The flip is an anecdote about one rule, and it is worth being clear about what the claim actually rests on. It rests on the blind discrimination test, where a reader is shown a voice's own description of itself and two replies, and picks which of the two was written under it. MEASUREMENTS.md has the rate, and more usefully the caveats.

## Install

Two commands, then one menu choice.

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Letting a tool edit your settings is a thing you should know the shape of beforehand rather than afterwards, so here is the shape of it. `--setup` emits the output style and wires the two hooks into `~/.claude/settings.json` with absolute paths, then prints the one step left. It backs the settings file up before touching it, it adds only what is missing so a second run changes nothing at all, it leaves every other key alone including other tools' hooks on the same events, and it refuses outright to touch a settings file that does not parse — and `--setup --dry-run` prints what it would change while writing nothing.

Anyone who would rather not have their settings written to can emit the style by hand: `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md`, and `COPE_CARD=<name>` in front of it emits a different one.

The menu choice is under `/config` -> Output style. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way in, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`. One thing to know before you look for a change and fail to find one: a style is read once at session start, so a new selection or a re-emitted card applies at the next session, or after `/clear`.

An output style goes at the end of the system prompt and the harness re-reminds the model of it as the conversation runs, which is why the card lands from there and did not land from a hook.

The hooks are for something else now, and the card no longer arrives through one of them. This is the hooks block of `~/.claude/settings.json`:

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

Stop scores the reply that was just written and records which rules fired. UserPromptSubmit reads that record and restates, mid-session, the items that have actually been firing — which is the one thing a file written once cannot do. The voice works with neither of them installed. These two are the measurement half.

Both commands assume that go install's target directory is on PATH in the environment the hook runs in, and a hook that silently does nothing almost always wants the absolute path instead. From a clone, `make install` builds it and no effigy checkout is needed, because `card/rules.json` is committed and compiled in. There is also a superseded SessionStart delivery, `--inject`, which remains for anyone who wants the old behaviour and stands down on its own when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable exactly as written, with no Python and no effigy checkout anywhere near it. A card dropped in `$XDG_CONFIG_HOME/cope/cards` is reached by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected at all, which is deliberate — the alternative is a session writing in the shipped voice while its config confidently names another one.

The one to read first is `card/demo/lecturer.effigy`, because it differs from the shipped card on register alone, and it is the card the discrimination run measured. What a card changes is the voicing half described under **The two things it checks**, and the numbers for that live in MEASUREMENTS.md rather than here.

Seeing what a card does without writing one takes about a minute at [demo/README.md](demo/README.md), where each file is this page again under a different card, same prompt, same facts, voice the only variable.

One file in that directory is not a voice at all: `card/demo/handoff.effigy` is a hypothesis, keeping the shipped card's handoff rules and dropping everything about prose, and it is meant to go through `make pairs` against the full card rather than to be rendered. `make cards` installs it with the others, so anyone listing their cards will find it sitting there among them, looking like a register it is not.

## What the hooks do that a file cannot

Stop fires after the reply exists, scores it against the loaded card, appends which rules fired to the session's rolling state, and writes one record per violation to the log. UserPromptSubmit then reads that rolling state — the state, not the violations log — and injects the card items gated on what has been firing, naming the counts as it goes, once the last injection has aged past `--refresh-every`. When a session has no history yet it falls back to the standing CONTINUE TEST. So the mid-session text is selected from measured output rather than fixed in advance, which is the thing a pasted CLAUDE.md cannot manage.

That is one mechanism and not a guarantee, and it should be read as such: the A/B in this repo does not separate the refresher from no refresher. The old SessionStart delivery via `--inject` is superseded and off by default, and `--setup` does not wire it.

## Why effigy notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here well off-label, and three of its blocks happen to be exactly what a prose gate wants: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with a failing and a passing example. That last one is why a rule in a card can name a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem from the opposite end, and it is the one to reach for when this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you choose it. basanite measures instead, comparing lemma frequency against a baseline over real transcripts, so it reports what you have been leaning on lately and leaves the judgement with you — its own README calls that awareness rather than prohibition. Which one fits is a question about mood more than about correctness, since a heavy hand is what you want when a habit is annoying you today and a moving measurement is what you want when you would rather watch the drift than legislate it. They compose, on different hooks with no shared state, and running both is reasonable.

There is a third axis again in [caveman](https://github.com/JuliusBrussee/caveman), a separate project by a different author that compresses agent replies to cut output tokens. cope shapes prose, basanite tracks vocabulary, caveman shortens. Wanting fewer tokens rather than different structure is a reason to go there instead.

## The rules

Grouping them by where they are implemented tells you about the code. Grouping them by which of the two jobs they check tells you something you can use.

Voicing:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session; the only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks, with an ordinal counting as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- `buried_decision` (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` (loop) — says the work is done with nothing on the page that could have shown it, no command, no count, no file.
- `loop_ask` (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

Read the two lists together and the implementation explains itself. A POSTPROC pattern matches a span of text, which means it can only ever describe wording, so every voicing rule that needed more than a pattern had to be written in Go alongside the structure rules. That is why the shipped card carries 0 POSTPROC rules of its own. Anyone arriving here expecting a long list of banned phrases should know that the list lives in basanite, on purpose.

The structure rules do vary in one place, and it is not by card — it is by who is going to read the turn. The interactive lane is chosen for any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt that opened the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself, and it drops `ask_not_last`, `forked_end`, `dangling_end` and `buried_decision` while adding `unverified_done` and `loop_ask`. Nobody is reading yet. A report that correctly names what it left open and then stops would fail three of the dropped rules, and a question inside it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran.

## Flags

| flag | default | does |
| --- | --- | --- |
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

## What lands on disk

Three files, all mode 0600, all under `$XDG_STATE_HOME/cope`. `violations.jsonl` holds one JSON record per violation, carrying the matched text and about 70 characters either side, which is to say the log quotes your replies back at you. `refresher-<session-id>` is an empty file whose mtime is the refresher clock. `session-<session-id>.json` is the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window, holding no prose at all, only rule names and counts.

## Editing the card

effigy owns the `.effigy` grammar, `make rules` regenerates `card/rules.json` from the card, and `make check-rules` is what CI runs so that the enforced rules and the injected rules cannot quietly drift apart.

There is a NEVER budget of 10, and anything over it is reported at load rather than dropped in silence. The budget is charged against each injection separately and not against the card file, because SessionStart prints the always-on rules while the refresher prints the evidence-gated ones and no code path renders their union — so a card may hold more NEVER rules in total than the budget allows and be perfectly healthy. The authoritative list of rules that really are discarded unrendered is `never_rules_over_budget`, and on a healthy card it is empty.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is small and worth having in front of you exactly as it stands. Selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`. Predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`.

A rule id for `@gate` has to be one the gate actually has, and an id for `@shape` must not collide with one, and a wrong id in either is reported at load rather than ignored. The reason after the dash is required in both forms, since a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still tells you what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript at once, and it is how these rules were chosen in the first place. `tools/backfill-sweep.sh` runs it over the N largest transcripts found anywhere under `~/.claude/projects`, which is emphatically not the same thing as one per project. The number worth watching as you tune is hits per character, not the share of turns hit, because that second number mostly tracks how long the turns were. MEASUREMENTS.md has the rates.

## Known limits

`labelled_opening` is not a tagger and should not be mistaken for one. `ask_not_last` says nothing whatsoever about the ordering of several asks among themselves.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written — which makes it a description of the output and not a judgement of it. The judgement lives in the discrimination test, and that test covers voicing only.

The largest limit is the one named under **The two things it checks**, and it is worth stating as it now stands rather than as it stood. A card can decline a built-in rule and can write one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question. So the compiled rules remain the only home a check like `clause_symmetry` can have, and a card wanting something outside both that vocabulary and a POSTPROC regex still has nowhere at all to put it. Both directions are also the card marking its own homework, since a decline lowers that card's score and an assertion raises it; both are worth reading with the reason attached, which is precisely why the syntax refuses to accept one without a reason. MEASUREMENTS.md is where the evidence for all of this sits.

## Layout

- `card/claude_voice.effigy` — the shipped card, in effigy notation
- `card/rules.json` — generated from it; embedded in the binary
- `card/demo/` — other cards, each written to sound like something else
- `cmd/cope-gate/` — the hook binary
- `internal/scan/` — the structure rules, the card's regex rules, and the card renderer
- `internal/effigy/` — the .effigy reader, so a card is usable as written
- `internal/transcript/` — Claude Code JSONL reader, and which lane a turn was written in
- `replay/` — the blind-pairs and discrimination harnesses, and their own README
- `demo/` — this README written again under each demo card
- `tools/` — card compiler, effigy-backed scorer, cross-project sweep
- `MEASUREMENTS.md` — what was run, on how much text, and what it said

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
