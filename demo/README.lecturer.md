# cope

A cope is the upper half of a foundry mould, the half carrying the shape being cast into, and this one ships an opinionated card that is live the moment you install it, checks two separate jobs — how a reply sounds and how it is shaped — and keeps that card in a file you can edit or swap once the default stops fitting.

Every file under [demo/README.md](demo/README.md) is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed; read two of them side by side and you will see what a card does faster than the rest of this page can tell you. One of them, `demo/README.claude-maximal.md`, comes from a card instructing every tic this model has been measured to have, and it makes the point in a single glance. It is also deliberately hard going.

## Sound and shape

Two complaints get made about a model's prose, and they get made in the same breath, which is most of why neither ever gets fixed. The first is what the sentences sound like: register, rhythm, diction, what a paragraph does with a detail, where flair is licensed. That half lives in the card, entirely — VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC — and swapping the card swaps every word of it, and it is the half with a measured result standing behind it. The second is the shape of the reply as a thing you have to use: where the decision sits, whether the ending gives `continue` anything to refer to, whether a claim that the work is finished carries something that could have shown it. That half is compiled into the binary and reads the same whichever card is loaded. A sentence reaching for the balanced two-beat is a voicing fault; a reply that names an open problem in its last paragraph and then stops, leaving you nothing to answer, is a structural one. A reply can be immaculate on one and useless on the other, and that is the whole reason for keeping them apart.

A card used to be mute on the structural half. It now reaches in from both directions. A header line reading `@gate <rule_id> off — <why>` declines a built-in rule, which exists because a card whose VOICE asked for exactly the move a built-in rule catches was being marked down for obeying itself. A header line reading `@shape <id>: <selector> <predicate> — <why>` asserts a structural rule of the card's own, which exists because a card's commitment about how its replies end had nowhere to be checked. The vocabulary that second form writes in counts words and sentences and asks whether a block poses a question, and there it stops: a check living outside that vocabulary and outside a POSTPROC regex still has nowhere in a card to go.

## Why instructions slide off

You have written a global CLAUDE.md, watched it fail to hold, and concluded that the model does not follow instructions. Let me take the file's side for a minute. Nearly everyone who has edited that file believes it is the system prompt. It is not. It arrives as one message pinned to the first turn, and every turn written afterwards piles on top of it. An output style goes into the system prompt itself, and the harness re-reminds the model of that text as the conversation runs. Move one card between those two places without altering a syllable of it, and you have most of the difference cope makes. It was measured, and the run is in [MEASUREMENTS.md](MEASUREMENTS.md).

Instruction alone will not fix the phrasing either, and the evidence for that is close to hand and slightly humiliating. A global CLAUDE.md banning the not‑A‑but‑B flip was read on every single turn of the session that built this tool, while the ban was the topic of the session, and the flip surfaced twice anyway. Name a surface form and the move goes looking for a variant.

The third failure is no kind of phrasing habit, which is why no instruction was ever going to reach it. A reply can be beautifully written and still close by naming a problem and going quiet, and when it does, you spend an entire round trip asking what it wants from you. Nobody banned that. Nobody could have written the ban.

One anecdote about one rule is not a result, and the flip is an anecdote. What the claim actually rests on is the blind discrimination test: a reader is shown a voice's own description of itself alongside two replies and picks which one was written under it. The rate, and the reasons the rate does not carry further than it does, are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Installing

Installation is two commands and one menu choice, and the order matters more here than anywhere else on this page.

```sh
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

Then open `/config` and pick the cope card under Output style. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the way in, and the same thing can be set as `"outputStyle"` in `.claude/settings.local.json`.

What `--setup` did, before you go looking: it emitted the output style, wired the two hooks into `~/.claude/settings.json` with absolute paths, and printed the one step left. It backs the settings file up before touching it, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses outright to write to a settings file that does not parse — and `--setup --dry-run` prints what it would change and writes nothing at all. If you would rather no tool wrote to your settings, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` and stops there, with `COPE_CARD=<name>` in front of it to emit a different one.

A style is read once, at session start. Pick a new one, or re-emit a card you have just edited, and the conversation you are sitting in will not notice; the next session will, and so will this one after `/clear`. An output style sits at the end of the system prompt and gets re-raised by the harness as the conversation runs, which is the entire reason the card lands there.

The hooks are a separate matter, and the card no longer arrives through one. This is the hooks block of `~/.claude/settings.json`:

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

Stop scores the reply just written and records which rules fired. UserPromptSubmit restates, mid-session, the rules that have been firing lately, which a file written once cannot do. The voice works with neither of them; these two are the measurement half. Both commands assume that `go install`'s target directory is on `PATH` in the environment the hook actually runs in, and a hook that silently does nothing usually wants the absolute path instead. From a clone, `make install` builds it, and no effigy checkout is needed, because `card/rules.json` is committed and compiled into the binary. The old SessionStart delivery, `--inject`, is superseded and off by default, remains for anyone who wants that route, and stands down by itself when a cope output style is active.

## Writing in another voice

The gate reads `.effigy` directly, so a card is usable exactly as it was written, with no Python involved and no effigy checkout on disk. Drop a card in `$XDG_CONFIG_HOME/cope/cards` and reach it with `--card <name>` or `COPE_CARD`; `--rules` takes a path from anywhere at all; `make cards` installs the demo set for you. A name that resolves to nothing is a hard error and nothing gets injected, which is deliberate — the alternative is a session quietly writing in the shipped voice while its config names another one, and that is a bug you would chase for a week. Read `card/demo/lecturer.effigy` first, because it differs from the shipped card on register alone, and it is the card the discrimination run measured. Everything a card can move belongs to the first of the two jobs described under Sound and shape.

The shortest way to see what a card does without writing one is [demo/README.md](demo/README.md), where each file is this page again from a different card, same prompt and same facts, so the voice is the only variable.

One file in that directory is not a voice at all: `card/demo/handoff.effigy` is a hypothesis, keeping the shipped card's handoff rules and dropping everything about prose, and it is meant to go through `make pairs` against the full card rather than be rendered. `make cards` installs it alongside the rest, so you will find it when you list your cards, and now you know not to write in it.

## What the hooks add

A card in the system prompt is written once and true forever, which is fine right up until a session starts going wrong in one particular direction on one particular afternoon. Stop scores the reply as soon as it is written and appends the rules that fired to the session's rolling state, plus one record per violation to the log. UserPromptSubmit then reads that rolling state — the state file, never the violations log — and injects the card items gated on what has been firing, counts attached, falling back to the standing CONTINUE TEST when the session has no history yet and staying quiet until the last injection has aged past `--refresh-every`. So the mid-session text is chosen from measured output instead of being fixed in advance, and a pasted CLAUDE.md cannot do that at any price. Do not let me oversell it: this is one mechanism and not a guarantee, and the A/B in this repo does not separate the refresher from no refresher.

## Why a character-card notation

[effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, used here entirely off-label, and three of its blocks turn out to be precisely what a prose gate needs. POSTPROC is regex rules with a warn action applied after generation. WRONG holds an anti-pattern beside the thing that should have been written instead. TEST holds a named question with a failing example and a passing one, which is how a rule names a move rather than one wording of it.

[basanite](https://github.com/justinstimatze/basanite) answers the same problem from the opposite end, and it is the one to reach for when this tool feels too blunt. cope bans: a rule fires or it does not, the card says never, and the register is settled the moment you choose it. basanite measures instead, tracking lemma frequency against a baseline across real transcripts, reporting what you have been leaning on lately and leaving the judgement where it found it — awareness rather than prohibition, as its own README puts it. Which one suits you is a question about mood more than correctness: a heavy hand is what you want when a habit is annoying you today, and a moving measurement is what you want when you would rather watch the drift than legislate against it. They compose cleanly, on different hooks with no shared state. A third axis again is [caveman](https://github.com/JuliusBrussee/caveman), a separate project by a different author that compresses agent replies to cut output tokens, and if what you want is fewer tokens rather than a different shape, go there instead.

## The rules

Grouped by which of the two jobs they belong to, because that is more useful than grouping them by where the code happens to live.

Voicing, from the shipped card's POSTPROC block:

- `flip` (warn) — the not‑A‑but‑B flip in its common surface forms, including not‑only‑but‑also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` (warn) — the reflexive intensifier standing in for *important* or *central*, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- `worth_noting` (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.

Voicing, compiled into `internal/scan` because a pattern could not carry them:

- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length repeating a content word across the joint, which is the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, a story you did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. It is the only rule reading the window rather than the reply, so it cannot fire until a session has a history.

Structure, all of it compiled in:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks, with an ordinal counting as a label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets among the things they wanted, so an opener written as a bold label is unpoliced on purpose.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths show a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on *or*, questions inside list items and table cells, and bare deference tags like *your call* are read as continuing the decision above rather than adding a second one.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

That grouping tells you something you can use. A POSTPROC pattern matches a span of text, so it can only ever describe wording, and every voicing rule needing more than a pattern had to be written in Go beside the structure rules. Which is why the shipped card carries three POSTPROC rules and no more. If you came here expecting a long list of banned phrases, that list lives in basanite, deliberately.

The structure rules do vary in one place, and it has nothing to do with which card is loaded. It has to do with who is going to read the turn. An interactive turn is any turn that is not a loop turn, and it gets the full set, because somebody is waiting at a terminal and the ending is where they decide what happens next. A loop turn — opened by `/loop`, by `/goal`, or by the sentinel a dynamic-pacing loop sends itself — drops `ask_not_last`, `forked_end`, `dangling_end`, and `buried_decision`, and adds `unverified_done` and `loop_ask`. Nobody is reading yet. A report that correctly names what it left open and then stops would fail three of the dropped rules, and a question inside it lands in a log where the next iteration reads it as an instruction to itself. What arrives in their place is the claim check: a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | What it does |
| --- | --- | --- |
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

## What lands on disk

Three files, all mode `0600`, all under `$XDG_STATE_HOME/cope`. `violations.jsonl` holds one JSON record per violation, carrying the matched text and roughly seventy characters either side of it, which is to say the log quotes your replies back at you. `refresher-<session-id>` is an empty file whose mtime is the refresher clock. `session-<session-id>.json` is the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired across a twenty-turn window, with no prose stored in it — rule names and counts only.

## Editing the card

The `.effigy` grammar belongs to effigy, `make rules` regenerates `card/rules.json` from `card/claude_voice.effigy`, and `make check-rules` is what CI runs so that the rules being enforced and the rules being injected cannot quietly drift apart. There is a NEVER budget of ten, and anything over it is reported at load rather than dropped in silence. The budget is charged against each injection separately and not against the card file, since the SessionStart path prints the always-on rules and the refresher prints the evidence-gated ones and no code path anywhere renders their union — so a card can hold more NEVER rules in total than the budget and still be perfectly healthy. The authoritative list of rules that genuinely were discarded unrendered is `never_rules_over_budget`, and on a healthy card it is empty.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is short and exact:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

For `@gate` the id has to be one the gate already has; for `@shape` it must not collide with one. Either way a wrong id is reported at load rather than passed over. The reason after the dash is required in both directions, because a rule a card wrote and a rule a card refused are equally unreviewable without one. Two working examples are in the tree: `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end`, since its VOICE block asks for exactly the balanced landing and the arriving close those two catch, and `card/demo/handoff.effigy` asserts `readable_cold` — `last paragraph words <= 60` — because its peak asks the reader to re-enter cold and read only the last block, and no built-in rule checks whether that block can be read that way. That 60 was measured: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90. A declined rule keeps running, and only this card's score drops it, so a backfill still tells you what it would have caught. A `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores a whole session transcript in one pass, and that is how the rules were chosen in the first place. `tools/backfill-sweep.sh` runs it across the N largest transcripts found anywhere under `~/.claude/projects`, which is emphatically not one per project. The number to watch is hits per character, and not the share of turns hit, because that second figure mostly tracks how long the turns happened to be. Rates are in [MEASUREMENTS.md](MEASUREMENTS.md).

## Where it falls short

The two jobs organise the failures as neatly as they organise everything else. `labelled_opening` is a heuristic and not a tagger. `ask_not_last` says nothing whatever about the order of several asks. The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was *for* rather than how it was written, so treat it as a description of the output rather than a judgement on it — judgement lives in the discrimination test, and that test covers voicing only.

The largest limit is the one described under Sound and shape, and it is smaller now than it used to be. A card can decline a built-in rule and can assert one of its own, and the vocabulary it asserts in counts words and sentences and asks whether a block poses a question. The compiled rules therefore remain the only place a check like `clause_symmetry` can live, and a card wanting something outside both that vocabulary and a POSTPROC regex has nowhere to put it. There is a second thing to hold in mind while reading any card's score, which is that both new directions are the card marking its own homework: a decline lowers that card's score and an assertion raises it. That is exactly why the syntax demands a reason. Read the reasons. [MEASUREMENTS.md](MEASUREMENTS.md) has the rest.

## Layout

| Path | What |
| --- | --- |
| `card/claude_voice.effigy` | the shipped card, in effigy notation |
| `card/rules.json` | generated from it; embedded in the binary |
| `card/demo/` | other cards, each written to sound like something else |
| `cmd/cope-gate/` | the hook binary |
| `internal/scan/` | the structure rules, the card's regex rules, and the card renderer |
| `internal/effigy/` | the .effigy reader, so a card is usable as written |
| `internal/transcript/` | Claude Code JSONL reader, and which lane a turn was written in |
| `replay/` | the blind-pairs and discrimination harnesses, and their own README |
| `demo/` | this README written again under each demo card |
| `tools/` | card compiler, effigy-backed scorer, cross-project sweep |
| `MEASUREMENTS.md` | what was run, on how much text, and what it said |

---

This README was written by `tools/generate_readme.py` from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
