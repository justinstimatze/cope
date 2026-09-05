# cope

A cope is the upper half of a foundry mould — the half carrying the shape being cast into — and this one ships a single opinionated card that is live from the moment you install it, scoring two unrelated things about a reply rather than one, how the sentences sound and how the reply is shaped for the person who has to act on it, with the card itself kept as an ordinary file you can edit or swap when the default is not the voice you wanted.

Every file under [`demo/README.md`](demo/README.md) is this page written again from a different card, same prompt and same facts, so the card is the only thing that changed; reading two of them against each other shows what a card does faster than the rest of this page explains it, and [`demo/README.claude-maximal.md`](demo/README.claude-maximal.md) — written from a card that instructs every tic this model is measured to have, and deliberately hard going for that reason — makes the point in a single glance.

## Two jobs in one reply

Voicing is what the sentences sound like: register, rhythm, diction, where a flourish is licensed. It lives entirely in the card, so swapping the card swaps every word of it, and it is the half with a measured result standing behind it. Structure is a different complaint altogether — where the decision sits, how the ending leaves the reader placed — and it is compiled into the binary, the same rules whichever card you load. A sentence reaching for the balanced two-beat, two clauses of equal weight with a content word echoed across the joint, is a voicing problem. A reply that names an open problem in its last paragraph and then stops, leaving `continue` nothing to refer to, is a structural one, and it can happen in the most beautiful prose you have ever been handed. The same reply can be flawless on one axis and useless on the other. That is the whole reason for keeping them apart.

A card does reach into the structural half, in two directions, and both are recent. The first is refusal: `@gate <rule_id> off — <why>`, one per line in the card header, because a card whose `VOICE` block asked for something a built-in rule catches was being marked down for obeying itself. `card/demo/lecturer.effigy` declines `clause_symmetry` and `dangling_end` for exactly that reason. The second is assertion: `@shape <id>: <selector> <predicate> — <why>`, which exists because a card's own commitment about how a reply ends had nowhere to be checked — `card/demo/handoff.effigy` asserts `readable_cold`, requiring its last paragraph to stay inside sixty words, since its whole peak is a reader re-entering cold and reading that block alone. What a card writes in that syntax counts words and sentences and asks whether a block poses a question, and nothing beyond that; a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to live. The run behind the voicing claim, and the reasons its numbers do not carry further, are in [`MEASUREMENTS.md`](MEASUREMENTS.md).

## Where the instruction sits

You have almost certainly edited a global `CLAUDE.md`, watched it hold for three turns and then quietly stop mattering, and concluded that the model is simply forgetful. Here is the part nobody tells you: that file is not the system prompt. It arrives as one message attached to the first turn of the conversation, and everything written afterwards piles on top of it, until turn forty is being written by something that read your instructions once, a long time ago, in a room it has since left. An output style goes somewhere else entirely — into the system prompt itself, which the harness re-reminds the model of as the conversation runs. Move one card between those two places, changing not a syllable of its contents, and you get most of what `cope` does. It was measured, and [`MEASUREMENTS.md`](MEASUREMENTS.md) has the run.

Placement is not the whole of it, though, and the second problem is the one instruction cannot reach by trying harder. A global `CLAUDE.md` banning the "not A, it's B" flip is read on every single turn. The flip appeared twice in the session that built this tool, while the ban was the topic under discussion. Naming a surface form does not remove the move; it pushes the move into the next variant of itself.

The third failure is not a phrasing habit at all, which is why no phrasing instruction was ever going to catch it. An ending that leaves you nothing to answer costs a whole round trip — you type something to restart a conversation that should never have paused — and no ban on a wording would have prevented it, because nothing in the wording was wrong.

The flip is an anecdote about one rule, and it would be dishonest to let it carry the argument. What the argument actually rests on is the blind discrimination test: a reader is shown a voice's own description of itself, then two replies, and picks which one was written under it. The rate and the caveats are in [`MEASUREMENTS.md`](MEASUREMENTS.md).

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
# then: /config -> Output style -> claude_voice
```

The first command builds the binary. The second does the entire install: `cope-gate --setup` emits the output style, wires the three hooks into `~/.claude/settings.json` with absolute paths, and prints the one step remaining. Before it writes anything it backs the settings file up, and it adds only what is missing, so a second run changes nothing; every other key survives untouched, including other tools' hooks sitting on the same events, and a settings file that does not parse is left strictly alone rather than repaired. If you would like to see the diff before it happens, `cope-gate --setup --dry-run` prints what it would change and writes nothing.

Then the menu, and this is the step that quietly costs people the whole tool. Open `/config`, go to Output style, and select the entry named `claude_voice`. Not `cope` — the entry carries the shipped card's id, and a reader who has just installed something called `cope` will scan that list for the word `cope`, fail to find it, and walk away with a style unselected and nothing changed. The standalone `/output-style` command was removed in Claude Code v2.1.91, so `/config` is the route; the same thing can be set as `"outputStyle": "claude_voice"` in `.claude/settings.local.json`.

If you would rather not have a tool editing your settings, `cope-gate --output-style` writes the loaded card to `~/.claude/output-styles/<card>.md` by itself, and `COPE_CARD=<name>` in front of it emits a different one.

A style is read once at session start. So a new selection, or a freshly re-emitted card, applies at the next session or after `/clear` — not in the conversation you are sitting in, which will look exactly as broken as it did a minute ago. The reason the card goes here rather than into a hook is the placement argument above: the end of the system prompt is re-reminded, and one turn-zero message is not.

The hooks are separate, and the card no longer arrives through any of them:

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

`Stop` buys you a score on the reply after it has been written. `UserPromptSubmit` buys you something a file written once cannot do at all: the rules that have actually been firing, restated mid-session, in proportion to how often they fired. The voice works without either hook. These are the measurement half.

Both commands assume that `go install`'s target directory is on `PATH` in the environment the hook runs in, and a hook that appears to do nothing at all usually wants the absolute path written out instead. A clone builds with `make install` and needs no `effigy` checkout of its own, because `card/rules.json` is committed and compiled into the binary. There is also `--inject`, the superseded delivery that prints the card as prompt text, kept for anyone who wants the old behaviour and standing down on its own whenever a `cope` output style is active.

## Writing in another voice

The gate reads `.effigy` files directly, which means a card is usable exactly as you wrote it, with no Python anywhere in the loop and no `effigy` checkout on the machine. Drop a card into `$XDG_CONFIG_HOME/cope/cards` and it is reachable by `--card <name>` or by `COPE_CARD`; `--rules` takes a path from anywhere on disk; `make cards` installs the demo set in one go. A name that resolves to nothing is an error, and nothing is injected — better a loud failure than a session writing placidly in the shipped voice while its config names another one.

Read `card/demo/lecturer.effigy` first. It differs from the shipped card on register alone, which is what made it the card the discrimination run was built around, and reading it beside `card/claude_voice.effigy` shows you the exact width of what a card controls: the voicing half described near the top of this page, all of it, and none of the compiled structure rules. The numbers for that run live in [`MEASUREMENTS.md`](MEASUREMENTS.md) and not here.

Shorter route, if you do not want to write anything: [`demo/README.md`](demo/README.md), where every file is this README rendered again from a different card against the same prompt and the same facts, so the voice is the only variable in the room.

One file in that directory is not a voice. `card/demo/handoff.effigy` is a hypothesis — it keeps the shipped card's handoff rules, drops everything about prose, and is meant to be run through `make pairs` against the full card rather than rendered into a page. `make cards` installs it alongside the rest, so you will find it when you list your cards, and it is not a register to write in.

## Inside the hooks

`Stop` fires after the reply exists. It scores what was just written, appends the rule ids that fired to the session's rolling state, and writes one record per violation to the log — nothing about the reply's contents goes into the state file, only names and counts.

`UserPromptSubmit` is the interesting one, and what makes it interesting is which file it reads. Not the violations log. It reads the rolling state, and injects the card items gated on what has been firing in this session, naming the counts as it does so, so the mid-session reminder is chosen from your measured output rather than fixed in advance by whoever wrote the card. When a session has no history yet it falls back to the standing `CONTINUE TEST`, and it stays quiet entirely until the last injection has aged past `--refresh-every`. This is one mechanism and not a guarantee; the A/B in this repo does not isolate the refresher from having no refresher at all, and it would be overselling to pretend otherwise.

`PreToolUse` scores prose on its way out of the session — the `description`, `body` or `content` field an external write is about to post, matched against the Linear save tools named in the settings block. It is warn-only by construction: it returns `additionalContext` and never a `permissionDecision`, so the call goes through regardless and the model simply learns what the prose scored. It writes no session state, and it scores in the external lane, which drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end`.

## Why a notation meant for game characters

The card is written in [`effigy`](https://github.com/justinstimatze/effigy), a character-card notation for game NPCs, used here well off-label. Three of its blocks turn out to be exactly what a prose gate needs: `POSTPROC` is regex rules with a warn action applied after generation, `WRONG` holds an anti-pattern beside the thing that should have been written instead, and `TEST` holds a named question with a failing and a passing example — which is how a rule comes to name a move rather than one particular wording of it.

[`basanite`](https://github.com/justinstimatze/basanite) is the same problem answered the other way round, and the one to reach for when this one feels too blunt in the hand. `cope` bans: a rule fires or it does not, the card says never, and the register is settled the moment you choose it. `basanite` measures instead — lemma frequency against a baseline over real transcripts, reporting what you have actually been leaning on lately and leaving the judgement with you, which its own README calls awareness rather than prohibition. Which one fits is a question about mood far more than about correctness. A heavy hand is what you want when a habit has been annoying you all week; a moving measurement is what you want when you would rather watch the drift than legislate against it. They compose cleanly — different hooks, no shared state — and running both is entirely reasonable.

[`humanizer`](https://github.com/blader/humanizer) is a skill, by a different author, that rewrites AI-sounding prose against 35 patterns drawn from Wikipedia's "Signs of AI writing", the page WikiProject AI Cleanup maintains. It is called on a text and hands back a rewrite; `cope` fires at a hook, scores what was already written, and edits nothing. Its pattern list is the wider of the two, and if you want a rewrite rather than a score you want that one. The formatting patterns are where the two disagree on purpose: `cope` had a `bold_label` rule that banned `humanizer`'s bold mini-headings, until 52 blind pairs put bold and bullets among the three things deciding a reply for this repo's reader, and the rule was deleted rather than tuned.

[`caveman`](https://github.com/JuliusBrussee/caveman) is a separate project, by a different author, that compresses agent replies to cut output tokens — a fourth axis entirely. `cope` shapes prose, `basanite` tracks vocabulary, `humanizer` rewrites, `caveman` shortens. It is worth naming here because a reader who came for fewer tokens rather than different structure should go there instead of reading on.

## The rules, and which job each one is doing

Voicing, from the card and from the compiled set together:

- `flip` — the not-A-but-B flip in its common surface forms, including not-only-but-also and the inverted "A, not B". The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- `load_bearing` — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register; say what the thing carries instead.
- `worth_noting` — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- `clause_symmetry` — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint, the balanced two-beat.
- `apology` — the reply performs contrition instead of stating the correction and moving on.
- `self_postmortem` — the reply turns to account for its own errors, which is a story the reader did not ask for.
- `announced_length` — the reply announces its own length rather than cutting it.
- `cross_turn_repeat` — a turn of phrase this reply shares with several earlier ones in the same session. The only rule reading the window rather than the reply, so it cannot fire until a session has a history.
- `repeated_opening` — three or more sentences in one reply opening on the same two words. Where `cross_turn_repeat` reads the session window for a construction reused across turns, this one reads a reply against itself, and two is left alone because two is a rhythm.
- `fragment_run` — three consecutive sentences of five words or fewer with no finite verb in any of them. One fragment is emphasis and this repo's own register is full of them; a run of three is the staccato blind judges read as generated. Neither clipped demo card trips it, so neither declines it, and a card that wants the run says so with `@gate`.

Structure, all of it compiled and none of it card-supplied:

- `labelled_opening` — a prose paragraph opening on a short verbless fragment that the rest of it unpacks, with an ordinal counting as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its `bold_label` rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- `paragraph_uniformity` — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- `ask_not_last` — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- `dangling_end` — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving `continue` nothing to refer to.
- `buried_decision` — an open problem landing after the last question or offer, burying the decision point above it.
- `forked_end` — two or more things to act on in the closing blocks with nothing marking which comes first, so answering `continue` means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- `unverified_done` — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- `loop_ask` — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.
- `echoed_heading` — a heading of two or more content words whose first sentence below repeats every one of them, spending a line to say what the heading already said.

Look at where the line falls in that pair of lists and you can read the implementation off it. A `POSTPROC` pattern matches a span of text, which means it can only ever describe a wording; every voicing rule that needed more than a wording — a rule about a whole reply, or about a session, or about the shape of three sentences in a row — had to be written in Go, sitting beside the structure rules in `internal/scan`. That is why the shipped card carries three `POSTPROC` rules and not thirty. If you arrived expecting a long banned-phrase list, the list is real and it lives in `basanite`, on purpose, because a frequency measurement over your own transcripts is a better instrument for that job than a regex someone else wrote.

The structure rules do vary, though not by card. They vary by who is going to read the turn. The interactive lane is the default and applies to any turn that is not a loop turn, on the grounds that somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane is chosen when the prompt opening the turn was `/loop` or `/goal`, or the sentinel a dynamic-pacing loop sends itself, and it drops `ask_not_last`, `buried_decision`, `dangling_end` and `forked_end` while adding `unverified_done` and `loop_ask` — because nobody is reading yet, a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check: a report saying the work is done has to say what it ran. The external lane, used when `--pretool` scores prose an external write is about to post, drops the same four and swaps nothing in, because a ticket has a reader and no ending they can answer — it is read days later by somebody who was not in the session, which is the condition every surviving rule was written for in the first place.

At the bottom of a report you will sometimes see hits clustered rather than listed flat, and the `Stop` hook, `--check` and `--pretool` all do this. Breadth clustering happens when three or more distinct rules land on one paragraph; density clustering when one rule lands on one paragraph three or more times; when both hold, breadth wins, since naming three rules tells you to rewrite the block rather than hunt one construction through it. The reason is that every rule fires alone and knows nothing whatsoever about the others, so a report is a flat list you work down hit by hit — and three hits across three paragraphs are three small edits, while three hits inside one paragraph are one paragraph to write again. The density half has a measurement behind it: `--check` over 107 tracked documents produced 114 `flip` hits, of which seven were worth changing, and every one of the seven was visible as three in a paragraph rather than as anything about the form itself. Three is the floor on both conditions because two rules on a paragraph is ordinary and two hits of one rule is a coincidence you can see unaided. Nothing about what fires, or how it is scored, changes.

## Flags

| Flag | Default | What it does |
| --- | --- | --- |
| `--ab` | `false` | rotate refresher windows through each variant in turn, recording which one a turn was written under |
| `--ab-arms` | (empty) | comma-separated variants to rotate through, implying `-ab` (default `inject,hold`; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report how each variant did; `-` reads the default path |
| `--author-docs` | `false` | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | `false` | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also `$COPE_CARD` |
| `--card-from-sample` | (empty) | print a prompt for writing a card from this writing sample; `-` reads stdin |
| `--cards` | `false` | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; `-` reads stdin |
| `--check-lane` | (empty) | score `-check` in the given lane: interactive (default), loop, or external |
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

## What lands on disk

Three files, all mode `0600`, all under `$XDG_STATE_HOME/cope`.

`violations.jsonl` holds one JSON record per violation, carrying the matched text and roughly seventy characters either side of it. Read that sentence twice before you install this on anything sensitive: the log quotes your replies back to you, on disk, in the clear. `refresher-<session-id>` is an empty file whose mtime is the refresher clock and nothing else. `session-<session-id>.json` is the rolling record the mid-session injection is chosen from — turn count, characters, and which rules fired over a twenty-turn window — and it stores no prose at all, only rule names and counts.

## Editing the card

The `.effigy` grammar belongs to [`effigy`](https://github.com/justinstimatze/effigy). Edit `card/claude_voice.effigy`, run `make rules` to regenerate `card/rules.json`, and CI runs `make check-rules` so that what is enforced and what is injected cannot drift apart while nobody is looking.

There is a budget of 10 on `NEVER` rules, and the important thing about it is what it is charged against. Not the card file — each injection separately, because the always-on rules and the evidence-gated ones are printed by different code paths and nothing anywhere renders their union. A card may therefore hold more `NEVER` rules in total than the budget and be in perfect health. Anything genuinely over budget is reported at load rather than dropped in silence, and the authoritative list of rules discarded unrendered is empty when the card is well.

Both card-authored forms go in the card header, one per line:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The `@shape` vocabulary is small and exact:

- selectors: `first paragraph`, `last paragraph`, `every paragraph`, `some paragraph`, `reply`
- predicates: `words <= N`, `words >= N`, `sentences <= N`, `sentences >= N`, `asks`, `does not ask`

For `@gate` the id has to be one the gate actually has; for `@shape` it must not collide with one. Either way a wrong id is reported at load and not quietly ignored. The reason after the dash is required in both directions, because a rule a card invented and a rule a card refused are equally unreviewable without one. Note also that a declined rule still runs — only this card's score drops it — so `--backfill` will still tell you what it would have caught, and a `@shape` violation is reported in the card's own words rather than in any sentence the binary supplies.

`card/demo/handoff.effigy` shows the assertion form doing real work. Its `readable_cold` shape caps the last paragraph at sixty words, and the sixty is not a taste: across 43,155 assistant replies the closing block runs 33 words at the median and 56 at p90.

## Calibrating

`cope-gate --backfill` scores an entire session transcript in one pass, and that is how the rules in here were chosen in the first place. `tools/backfill-sweep.sh` runs it across the N largest transcripts found anywhere under `~/.claude/projects` — largest overall, which is emphatically not one per project, and the difference will bite you if you assume otherwise. Watch hits-per-character rather than the share of turns hit. The share of turns hit mostly tells you how long your turns were. Rates are in [`MEASUREMENTS.md`](MEASUREMENTS.md).

## Known limits

`labelled_opening` is not a part-of-speech tagger and will sometimes disagree with one. `ask_not_last` has nothing to say about the ordering of several asks relative to each other, only about an ask that is not last.

The hit rate is roughly four fifths structure, and the A/B run found that the four fifths tracks what a reply was for rather than how well it was written — which makes it a description of the output and not a verdict on it. The verdict lives in the discrimination test, and that test covers voicing only.

The largest limit is the split described near the top of this page, in the form it now takes rather than the form it took. A card can decline a built-in rule and can write one of its own, and the vocabulary it writes in counts words, counts sentences, and asks whether a block poses a question. The compiled rules therefore remain the only place something like `clause_symmetry` can live, and a card wanting a check outside both that vocabulary and a `POSTPROC` regex still has nowhere to put it. Both directions are the card marking its own homework: a decline lowers that card's score, an assertion raises it, and each is worth reading with the reason attached. That is precisely why the syntax insists on one. The rest of the caveats are in [`MEASUREMENTS.md`](MEASUREMENTS.md).

## Layout

| Path | What |
| --- | --- |
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

MIT. justin@justinstimatze.com
