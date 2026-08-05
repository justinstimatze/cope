# cope

One opinionated card, on from the moment it is installed and the whole product for most readers, checking two jobs rather than one — how a reply sounds and how it is shaped — and written as a plain file, so it can be edited or swapped; a cope is the upper half of a foundry mould, the half carrying the shape being cast into.

[demo/README.md](demo/README.md) is this page written again from each other card, same prompt and same facts, so the card is the only thing that changed; reading two of them against each other shows what a card does faster than the rest of this page explains it. [demo/README.claude-maximal.md](demo/README.claude-maximal.md) comes from a card instructing every tic this model is measured to have, and it makes the point in a single glance. Deliberately hard going.

## Voicing and structure

Voicing. What the sentences sound like — register, rhythm, diction, what a paragraph does with a detail, and where flair is licensed. It lives in the card and nowhere else: VOICE, TRAITS, NEVER, WRONG, MES, POSTPROC. Swap the card and every word of it swaps. A sentence reaching for the balanced two-beat is a voicing fault, and this is the half with a measured result behind it.

Structure. The shape of the reply as a thing the reader has to use — where the decision sits, whether the ending gives "continue" something to refer to, whether an ask is last, whether a claim that the work is done carries anything that could have shown it. Compiled into the binary, in internal/scan, and the same rules whichever card is loaded, varying only by lane. A reply that names an open problem in its closing paragraph and then stops is a structural fault, whatever it sounds like. The same reply can be clean on one axis and bad on the other, which is the reason to keep them apart.

`@gate <rule_id> off — <why>`, one per line in the card header, declines a built-in rule. card/demo/lecturer.effigy declines clause_symmetry and dangling_end, because its VOICE block asks for the balanced landing and the arriving close that those two rules catch, and the card was scored down for obeying itself.

`@shape <id>: <selector> <predicate> — <why>`, also one per line in the header, states a structural rule of the card's own for the gate to check. card/demo/handoff.effigy asserts readable_cold — last paragraph words <= 60 — because its peak asks the reader to re-enter cold and read the last block only, and that commitment had nowhere to be checked.

The boundary. The @shape vocabulary counts words and sentences and asks whether a block poses a question, and nothing more; clause_symmetry is not writable in it and is not meant to be. [MEASUREMENTS.md](MEASUREMENTS.md) holds the runs and the reasons their numbers do not carry more than they do.

## The problem

Delivery, first, and most readers arrive holding it backwards. A global CLAUDE.md is not the system prompt. It arrives as one message attached to the first turn, and the conversation buries it under everything written after. Compare an output style, which sits in the system prompt itself, and which the harness re-reminds the model of as the conversation runs. Moving one card between those two places, without changing a word of it, is most of why cope works. Measured; the run is in [MEASUREMENTS.md](MEASUREMENTS.md).

A phrasing ban, second. It does not hold on its own: a global CLAUDE.md forbidding the "not A, it's B" flip is read every turn, and the flip still appeared twice in the session that built this, with the ban as the topic. Naming a surface form pushes the move into a variant.

The structural complaint, third, and no phrasing problem at all. An ending that leaves the reader nothing to answer costs a whole round trip. Different cause, and not a habit of wording an instruction could have banned.

The flip is an anecdote about one rule. The claim rests on the blind discrimination test: a reader shown only a voice's own description of itself picks which of two replies was written under it. MEASUREMENTS.md has the rate and the caveats.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
cope-gate --setup
```

The last step is manual. /config -> Output style, then pick the card.

Setup, in full. It emits the output style, wires the two hooks into ~/.claude/settings.json with absolute paths, and prints the one step left. It backs the settings file up first, adds only what is missing so a second run changes nothing, leaves every other key alone including other tools' hooks on the same events, and refuses to touch a settings file that does not parse. `--setup --dry-run` prints what it would change and writes nothing.

By hand, for anyone who would rather not have their settings written to. `cope-gate --output-style` writes the loaded card to ~/.claude/output-styles/&lt;card&gt;.md, and `COPE_CARD=<name>` in front of it emits a different one. The standalone /output-style command was removed in Claude Code v2.1.91, so /config is the way; the same thing can be set as `"outputStyle"` in .claude/settings.local.json.

Timing. A style is read once at session start, so a new selection or a re-emitted card applies at the next session or after /clear. Placement is the reason it goes here: an output style sits at the end of the system prompt and the harness re-reminds the model of it during the conversation, which the hook never did.

The hooks, which the card no longer arrives through. This is the hooks block of ~/.claude/settings.json:

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

Stop scores the reply just written and appends which rules fired to the session's rolling state, plus one record per violation to the log. UserPromptSubmit restates mid-session the card items gated on what has actually been firing, which a file written once cannot do. The voice works without either of them; these two are the measurement half.

PATH. Both commands assume go install's target directory is on PATH in the environment the hook runs in, and a hook that silently does nothing usually wants the absolute path instead. A clone builds with `make install` and needs no effigy checkout, because card/rules.json is committed and compiled in. `--inject` is the superseded delivery, off by default, kept for anyone who wants the old route, and it stands down on its own when a cope output style is active.

## Another voice

The .effigy reader is in the gate, so a card is usable as written — no Python, no effigy checkout. A card dropped in $XDG_CONFIG_HOME/cope/cards is reached by `--card <name>` or COPE_CARD, `--rules` takes a path from anywhere, and `make cards` installs the demo set. A name that resolves to nothing is an error and nothing is injected. Compare a session writing in the shipped voice while its config names another one, which is what the error exists to prevent.

card/demo/lecturer.effigy first. It differs from the shipped card on register alone, which is what makes it the discrimination rival and what the discrimination run measured. A card changes the voicing side described under Voicing and structure above, and reaches into the structural side only through @gate and @shape.

[demo/README.md](demo/README.md) is the shortest way to see what a card does without writing one: every file under demo/ is this README written again from a different card, same prompt and same facts, so the voice is the only variable.

card/demo/handoff.effigy is the exception in that directory. A hypothesis rather than a voice, keeping the shipped card's handoff rules and dropping everything about prose, meant to be run through `make pairs` against the full card rather than rendered; `make cards` installs it with the rest, so a reader listing their cards will find it there and should know it is not a register to write in.

## Scoring, and the mid-session reminder

Stop. It reads the reply just written, records which rules fired into the session's rolling state, and appends one record per violation to the log.

UserPromptSubmit. It reads the rolling state — not the violations log — and injects the card items gated on what has been firing, naming the counts, so the mid-session text is chosen from measured output rather than fixed in advance. That is the one thing a pasted CLAUDE.md cannot do. It falls back to the standing CONTINUE TEST when the session has no history yet, and stays quiet until the last injection has aged past `--refresh-every`. One mechanism, and no guarantee: the A/B in the repo does not separate the refresher from no refresher.

SessionStart, via `--inject`, is superseded and off by default.

## Effigy notation, and basanite

Effigy, used off-label. [effigy](https://github.com/justinstimatze/effigy) is a character-card notation for game NPCs, and three of its blocks do what a prose gate needs: POSTPROC is regex rules with a warn action applied after generation, WRONG holds an anti-pattern beside its replacement, and TEST holds a named question with fail and pass examples, which is how a rule names a move instead of one wording of it.

[basanite](https://github.com/justinstimatze/basanite). The same problem answered the other way round, and the one to reach for if this one is too blunt. cope bans: a rule fires or it does not, the card says never, and the register is fixed the moment you pick it. Compare basanite, which measures lemma frequency against a baseline over real transcripts and reports what you have been leaning on lately, leaving the judgement to you — awareness rather than prohibition, in its own README's words. A heavy hand suits a habit that is annoying you today; a moving measurement suits watching the drift instead of legislating it. They compose: different hooks, no shared state.

[caveman](https://github.com/JuliusBrussee/caveman), by a different author, compresses agent replies to cut output tokens. Compare cope, which reshapes prose and does not shorten it; a reader after fewer tokens rather than different structure should go there.

## The rules

Voicing rules. Three arrive from the shipped card's POSTPROC block, the rest are compiled:

- **flip** (warn) — the not-A-but-B flip in its common surface forms, including not-only-but-also. The Economist measured this family across 55,940 sentences on 2026-07-30 and put Claude at the top of it among every writer in the study, human or machine.
- **load_bearing** (warn) — reflexive intensifier for important or central, at 25.6 per 1k the heaviest measured lean in this register. Say what the thing carries instead.
- **worth_noting** (warn) — announces that something deserves attention instead of letting it earn the attention, at 6.5 per 1k.
- **clause_symmetry** — comma- or semicolon-joined clauses of near-equal length that repeat a content word across the joint: the balanced two-beat.
- **apology** — the reply performs contrition instead of stating the correction and moving on.
- **self_postmortem** — the reply turns to account for its own errors, which is a story the reader did not ask for.
- **announced_length** — the reply announces its own length rather than cutting it.
- **cross_turn_repeat** — a turn of phrase this reply shares with several earlier ones in the same session. The only rule that reads the window rather than the reply, so it cannot fire until a session has a history.

Structure rules, all compiled, and identical under every card:

- **labelled_opening** — a prose paragraph opening on a short verbless fragment that the rest of it unpacks; an ordinal counts as the label. List blocks and paragraphs under twelve words are skipped, and so is the bolded form: the card dropped its bold_label rule in July 2026 after blind readers named bold and bullets as something they wanted, so an opener written as a bold label is deliberately unpoliced.
- **paragraph_uniformity** — four or more prose paragraphs whose lengths have a coefficient of variation below `--min-cv`.
- **ask_not_last** (interactive) — a question or request for the reader sitting in an earlier block while the reply carries on past it.
- **dangling_end** (interactive) — an open problem named in the closing blocks with no question, offer, or explicit all-clear anywhere, leaving "continue" nothing to refer to.
- **buried_decision** (interactive) — an open problem landing after the last question or offer, burying the decision point above it.
- **forked_end** (interactive) — two or more things to act on in the closing blocks with nothing marking which comes first, so answering "continue" means picking one. Sentences opening on "or", questions inside list items and table cells, and bare deference tags like "your call" are read as continuing the decision above rather than adding another.
- **unverified_done** (loop) — says the work is done with nothing on the page that could have shown it: no command, no count, no file.
- **loop_ask** (loop) — an unattended run ends by asking, so the answer lands in a log and the next iteration reads the question as an instruction to itself.

The grouping says something about the implementation. A POSTPROC pattern matches a span of text, so it can only ever describe wording, and every voicing rule that needed more than a pattern had to be written in Go beside the structure rules. That is why the shipped card carries three POSTPROC rules and no long inventory of banned phrases. Compare basanite, where the inventory lives on purpose and is a measurement rather than a ban.

Lanes are the one place the structure rules vary — not by card, by who is going to read the turn. The interactive lane covers any turn that is not a loop turn, because somebody is waiting at a terminal and the ending is where they decide what happens next. The loop lane covers a turn opened by /loop or /goal, or by the sentinel a dynamic-pacing loop sends itself. It drops ask_not_last, forked_end, dangling_end and buried_decision, and adds unverified_done and loop_ask, because nobody is reading yet: a report that correctly names what it left open and stops would fail three of the dropped rules, and a question in it lands in a log where the next iteration reads it as an instruction to itself. What replaces them is the claim check — a report saying the work is done has to say what it ran.

## Flags

| Flag | Default | Does |
| --- | --- | --- |
| `--ab` | false | rotate refresher windows through the arms and record which arm each turn was written under |
| `--ab-arms` | (empty) | comma-separated arms to rotate through, implying -ab (default inject,hold; positive is the third) |
| `--ab-report` | (empty) | read a turn log and report the arms; - reads the default path |
| `--author-docs` | false | print a prompt for writing this repo's docs: the card, the introspected facts, the sections |
| `--backfill` | (empty) | score every assistant turn in this transcript and exit |
| `--block` | false | exit 2 on a violation whose action is reject (default warn-only) |
| `--card` | (empty) | name of an installed card to write in, from the cards directory; also $COPE_CARD |
| `--cards` | false | list the installed cards with the aim each one states, and exit |
| `--check` | (empty) | score a prose file against the card and exit; - reads stdin |
| `--describe` | false | print the card's voice as a target to recognise: the aim and the register, without the machinery |
| `--display` | false | MessageDisplay entry: rewrite what the reader sees, leaving the transcript alone |
| `--display-preview` | false | read prose on stdin and print it as --display would rewrite it |
| `--dry-run` | false | with --setup, print what would change and write nothing |
| `--inject` | false | print the card as prompt text for a SessionStart hook |
| `--log` | $HOME/.local/state/cope/violations.jsonl | append violations here; empty disables |
| `--min-cv` | 0.35 | flag paragraph-length coefficient of variation below this |
| `--output-style` | false | write the card to ~/.claude/output-styles as a Claude Code output style, which puts it in the system prompt rather than in one turn-zero message |
| `--output-style-dir` | (empty) | directory to write the output style into (default ~/.claude/output-styles) |
| `--refresh-every` | 30m0s | minimum age of the last card or refresher injection before the refresher fires |
| `--refresher` | false | UserPromptSubmit entry: inject the compact reminder once the last injection has aged |
| `--render-arm` | (empty) | print the mid-session reminder one arm would inject, and exit |
| `--render-for` | (empty) | comma-separated rule ids to render -render-arm against |
| `--render-lane` | (empty) | render -render-arm as the given lane sees it: interactive (default) or loop |
| `--rules` | (empty) | read the card from this .effigy or .json file instead of the one built into the binary |
| `--setup` | false | emit the output style and wire the hooks, then print the one step left |
| `--version` | false | print version and exit |

## On disk

| Path | Mode | Holds |
| --- | --- | --- |
| $XDG_STATE_HOME/cope/violations.jsonl | 0600 | one JSON record per violation, carrying the matched text and about 70 characters either side |
| $XDG_STATE_HOME/cope/refresher-&lt;session-id&gt; | 0600 | an empty file whose mtime is the refresher clock |
| $XDG_STATE_HOME/cope/session-&lt;session-id&gt;.json | 0600 | the rolling record the mid-session injection is chosen from: turn count, characters, and which rules fired over a 20-turn window. No prose is stored, only rule names and counts |

The log quotes replies back. Matched text, plus the surrounding context, in a file on your machine at 0600.

## Editing the card

effigy owns the .effigy grammar. `make rules` regenerates card/rules.json from card/claude_voice.effigy, and `make check-rules` is what CI runs, so the enforced and the injected rules cannot drift. The shipped card renders to 5,363 characters as injected.

The NEVER budget. Ten, and anything over it is reported at load rather than dropped silently. The budget is charged against each injection separately rather than against the card file: SessionStart prints the always-on rules, the refresher prints the evidence-gated ones, and no code path renders their union — so a card may hold more NEVER rules in total than the budget and still be healthy. Nothing in the shipped card is discarded unrendered.

Two card-authored forms, both one per line in the card header:

```
@gate <rule_id> off — <why>
@shape <id>: <selector> <predicate> — <why>
```

The @shape vocabulary:

- selectors: first paragraph, last paragraph, every paragraph, some paragraph, reply
- predicates: words <= N, words >= N, sentences <= N, sentences >= N, asks, does not ask

Ids. @gate needs one the gate already has; @shape needs one that collides with nothing. A wrong id is reported at load rather than ignored. The reason after the dash is required in both, because a rule a card wrote and a rule a card refused are equally unreviewable without one. A declined rule still runs and only this card's score drops it, so a backfill still reports what it would have caught, and a @shape violation is reported in the card's own words rather than in any sentence the binary supplies.

## Calibrating

`cope-gate --backfill` scores every assistant turn in a transcript at once, and is how the rules were chosen. tools/backfill-sweep.sh runs it over the N largest transcripts found anywhere under ~/.claude/projects. Compare one per project, which it is not.

Hits per character is the number that survives comparison. Compare the share of turns hit, which tracks how long the turns were. MEASUREMENTS.md has the rates.

## Known limits

labelled_opening is not a tagger. It looks for a short verbless fragment that the rest of the paragraph unpacks, skips list blocks and paragraphs under twelve words, counts an ordinal as a label, and leaves the bolded form alone.

ask_not_last says nothing about the order of several asks. It sees one ask in the wrong place, and no more than that.

The hit rate is roughly four fifths structure, and the A/B run found that four fifths tracks what a reply was for rather than how it was written. A description of the output, then, and no judgement of it. The judgement lives in the discrimination test, which covers voicing only.

The largest limit is the one Voicing and structure names. A card declines a built-in rule and writes one of its own, and the vocabulary it writes in counts words and sentences and asks whether a block poses a question — so the compiled rules remain the only place a check like clause_symmetry can live, and a card wanting something outside both that vocabulary and a POSTPROC regex has nowhere to put it.

Both directions are the card marking its own homework. A decline lowers that card's score, an assertion raises it, and neither reads without the reason attached, which is why the syntax requires one. MEASUREMENTS.md.

## Layout

| Path | What |
| --- | --- |
| card/claude_voice.effigy | the shipped card, in effigy notation |
| card/rules.json | generated from it; embedded in the binary |
| card/demo/ | other cards, each written to sound like something else |
| cmd/cope-gate/ | the hook binary |
| internal/scan/ | the structure rules, the card's regex rules, and the card renderer |
| internal/effigy/ | the .effigy reader, so a card is usable as written |
| internal/transcript/ | Claude Code JSONL reader, and which lane a turn was written in |
| replay/ | the blind-pairs and discrimination harnesses, and their own README |
| demo/ | this README written again under each demo card |
| tools/ | card compiler, effigy-backed scorer, cross-project sweep |
| MEASUREMENTS.md | what was run, on how much text, and what it said |

---

Written by tools/generate_readme.py from the prompt `cope-gate --author-docs` emits, and checked with `cope-gate --check`.

MIT. justin@justinstimatze.com
