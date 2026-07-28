# cope

A voice card for Claude's prose, plus a Stop hook that checks the drafted reply
against it after generation instead of asking for compliance beforehand.

The name is the upper half of a foundry mould — the part carrying the shape
you cast into.

## The problem

Claude's phrasing is recognizable, and being recognizable is distracting: you
read the shape of the sentence instead of the argument. Instruction alone does
not fix it. A global `CLAUDE.md` banning the "it's not A, it's B" flip is read
at the top of every turn, and the flip still appeared twice in the session that
produced this repo — while the ban was the subject under discussion. Naming a
surface form pushes the move into a variant. Something has to read the output.

## Install

```
go install github.com/justinstimatze/cope/cmd/cope-gate@latest
```

The card is compiled into the binary, so there is nothing to configure and no
data file to place. Building from a clone works the same way:

```
make install
```

Then add the hooks to `~/.claude/settings.json`. `SessionStart` puts the card
into context, `Stop` checks what came out, and `UserPromptSubmit` re-states the
one rule that decays first once an injection has aged past four hours.

```json
{ "type": "command", "command": "cope-gate --inject", "timeout": 5 }
{ "type": "command", "command": "cope-gate", "async": true, "timeout": 10 }
{ "type": "command", "command": "cope-gate --refresher", "timeout": 5 }
```

Every half renders from the same `.effigy` file, so the rules being enforced and
the rules being asked for cannot drift apart — `make check-rules` runs in CI and
fails when the committed JSON no longer matches the card. The injected card is
about 9,000 characters, once per session, in the cached prefix; the refresher is
374.

The gate is warn-only: it prints violations, appends them to
`$XDG_STATE_HOME/cope/violations.jsonl`, and exits 0. `--block` makes it exit 2
on a rule whose action is `reject`; pairing that with `asyncRewake` is what
makes a violation come back to the model instead of landing in a log nobody
reads. Worth doing once the false-positive rate is known — not before, since a
rule firing on a quarter of turns would cost a regeneration each time.

The log quotes your replies back. Each record carries the matched text and about
70 characters either side, so `violations.jsonl` is conversation content: it is
written `0600` in your state directory, and `--log ""` turns it off.

## Flags

| flag | default | what it does |
|---|---|---|
| `--inject` | off | print the card as prompt text, for `SessionStart` |
| `--refresher` | off | print the compact reminder, for `UserPromptSubmit`, once the last injection has aged |
| `--refresh-every` | `4h` | how long an injection stays fresh |
| `--backfill FILE` | — | score every assistant turn in a transcript and exit |
| `--rules FILE` | built-in | use a card from disk instead of the compiled-in one |
| `--log FILE` | `$XDG_STATE_HOME/cope/violations.jsonl` | where violations are appended; empty disables |
| `--min-cv` | `0.35` | flag paragraph-length coefficient of variation below this |
| `--block` | off | exit 2 on a violation whose action is `reject` |
| `--version` | — | print the version and exit |

## Calibrate

```
cope-gate --backfill ~/.claude/projects/<project>/<session>.jsonl
tools/backfill-sweep.sh 6      # the largest transcripts across all projects
```

Baseline as of 2026-06 across 116,154 assistant turns and 38.7M characters:
**1.22 hits per 1,000 characters**, 14–24% of turns. Use hits-per-character, not
turns-hit — turns-hit tracks how long the turns were. A session where the
density drifts toward 2 per 1k is the signal worth acting on.

## Rules

Two families, because they fail differently.

Regex rules live in the card's `POSTPROC` block and are read through
`tools/card2json.py`, so [effigy](https://github.com/justinstimatze/effigy)
stays the only thing that parses `.effigy` notation.

Shape rules are built into `internal/scan` because they don't reduce to
patterns: paragraphs opening with a label the rest then unpacks, clauses
balanced against each other, paragraph lengths that are too even, a question or
request for the user that lands mid-reply instead of at the end, and an ending
that names an open problem without a question, an offer, or an explicit
all-clear.

Shape is where the budget goes. Over 20.1M characters measured on 2026-07-28,
`bold_label` and `labelled_opening` took the top two places with 12,603 and
11,479 hits, while `significance_flag` — aimed at the phrase "that's an
important distinction" — fired **once**. The move it targeted happens
constantly; the phrasing almost never does, and the rule was cut. Full table in
`NOTES.md`.

Code fences, inline code, block quotes, and double-quoted passages are blanked
before matching. Without that, the rules fire hardest on their own
documentation.

`NEVER` rules are a weaker channel than either family above. A bare prohibition
with no detector behind it only works if the model attends to it in the injected
prompt, and effigy's own tuning history (v0.3.1, v0.3.2) found that attention
degrades past roughly ten same-weight rules. `Render()` ports that budget:
`CRITICAL`-prefixed rules sort first, the list caps at ten, and inline
`WRONG:`/`RIGHT:` examples are stripped from what survives. Anything over the
cap is reported on every load rather than silently dropped, and
`tools/card2json.py` warns at compile time, since which rule gets discarded is
worth choosing on purpose. A rule worth adding without a detector should usually
be a `TEST` instead — a question plus fail/pass examples names the move, where a
`NEVER` only names one surface form of it.

## Layout

```
card/claude_voice.effigy   the card; validates against effigy v0.7.0
card/rules.json            generated; embedded in the binary
card/embed.go              the go:embed directive
cmd/cope-gate/             the hook binary
internal/scan/             regex + shape rules, and the card renderer
internal/transcript/       Claude Code JSONL reader
tools/card2json.py         card -> rules.json
tools/run_postproc.py      scores a transcript via effigy's own validator path
tools/backfill-sweep.sh    cross-project baseline
```

## Editing the card

Only this needs effigy, and only because effigy owns the `.effigy` grammar. The
tools default to a sibling checkout; `EFFIGY_PATH` overrides it.

```
git clone https://github.com/justinstimatze/effigy ../effigy
make rules
make check-rules
```

Commit `card/rules.json` alongside the `.effigy` change — the binary embeds the
JSON, and CI fails if the two disagree.

## Known limits

`labelled_opening` decides "is this a fragment or a sentence" with a 30-word
verb list and two suffix patterns. It is not a tagger and will miss and
over-fire. It is the rule most likely to need replacing.

`ask_not_last` flags an ask that sits before the end but says nothing about the
order of several asks. Telling "printed in the wrong order" from "printed in the
right order" needs to know what the asks mean, not just where they sit.

Nothing here measures whether the prose actually got better. Four LLM judge
panels across ~4.2M tokens failed to distinguish three candidate voices in
aggregate while agreeing exactly on each individual criterion — see `NOTES.md`.
Reading it for a week is the test.

## License

MIT. Contact: justin@justinstimatze.com
