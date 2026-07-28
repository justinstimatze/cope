# Notes

A lab notebook, not documentation. Everything below is dated and describes what
was measured at the time; `README.md` describes the tool as it is now.

## Why effigy notation (2026-06)

Off-label, and the fit is better than it has any right to be.
[effigy](https://github.com/justinstimatze/effigy) is a character-card format
for game NPCs, and three of its blocks do exactly what a prose gate needs.
`POSTPROC` holds regex rules with an action of `reject`/`strip`/`warn`, applied
after generation — a post-hoc gate with a warn-only mode, already written.
`WRONG` holds a paired anti-pattern and its replacement. `TEST` holds a named
question with fail and pass examples, which is how a rule names the *move*
instead of one of its surface forms.

Only `strip` doesn't survive the transfer. A Stop hook sees the reply after it
is written and cannot edit it, so `cope-gate` reports any `strip` rule at load
time and treats it as `warn`.

## Shape detectors work, phrase detectors don't (2026-06, re-measured 2026-07-28)

The first run covered 36 assistant turns and 65,276 characters of the session
that produced the repo:

| rule | hits | note |
|---|---|---|
| `bold_label` | 8 | all in the first three turns, none after a correction |
| `flip` | 3 | 2 genuine, 1 self-reference; misses every disguised variant |
| `significance_flag` | 0 | aimed at a phrasing that never occurs |
| `verdict_handoff` | 0 | same |

The rule aimed at a structural form caught everything it was meant to, and its
hit distribution matched the human complaint without tuning. The phrase-template
rules caught almost nothing.

Re-measured on 2026-07-28 over 61,612 assistant turns and 20.1M characters, the
same shape held and settled the two doubtful rules:

| rule | hits |
|---|---|
| `bold_label` | 12,603 |
| `labelled_opening` | 11,479 |
| `clause_symmetry` | 1,694 |
| `ask_not_last` | 1,331 |
| `paragraph_uniformity` | 1,150 |
| `flip` | 458 |
| `dangling_end` | 350 |
| `verdict_handoff` | 38 |
| `buried_decision` | 18 |
| `significance_flag` | 1 |

`significance_flag` was cut on that evidence. `verdict_handoff` stays: 38 is
thin, but they are 38 real instances of a move worth naming, and the rule costs
one regex pass.

`bold_label` took two wrong guesses to get right. The tic is `**Label.**` with
the punctuation *inside* the bold, not `**Label**:`.

## LLM judge panels are reliable per-criterion and unreliable in aggregate (2026-06)

Four panels, roughly 4.2M subagent tokens. The same nine drafts judged by two
different models under a Latin square:

- Per-criterion rank order replicated *exactly* across models. Clarity ranked
  the authorial blend first in both (25, 23). Specificity ranked
  selection-of-detail first in both (23, 23). Signature ranked the bare
  functional control first in both (25, 23), with the blend last at 12 in both
  runs.
- The Borda aggregate *flipped* between models (blend 58 against selection 57).

Any pipeline collapsing criteria into one score reports its choice of model.
This is evidence for decorrelated per-axis judging and against scalar
aggregation.

Three earlier results that read as null were not null. The axes pull against
each other — completeness buys clarity and costs signature, terseness buys
signature and costs specificity — so the sum cancelled.

## Completeness and symmetry are separable (asserted, never measured)

The signature judges never penalised the blend's stated functions. Every quoted
offence was matched-clause symmetry: "a costume that fits well… a costume that
fits badly", "nothing to rank / nothing to match", "manufactured two-beat
symmetry". Both terms of a comparison can be named without balancing them. If
that holds, the clarity win and the signature loss come apart. The test would
vary symmetry alone with completeness fixed. It has not been run, and nothing
in the tool depends on the claim.

## Conditional gates are not justified (2026-06)

effigy rules can carry an `@when` condition. A blind check across three
registers — result report, disagreement, build proposal — found one register
serving all three: the same opening gesture, the same `So:` pivot verbatim in
two of three, the same em-dash appositive throughout. What varies is dictated by
the speech act rather than by voice.

The gates came out of the card then, and the plumbing that carried them came out
of `tools/card2json.py` and `internal/scan` on 2026-07-28. One rule keeps a note
that it slips hardest during disagreement, because disagreement offers fewer
nouns to point at.

## What basanite can and cannot see (2026-06)

[basanite](https://github.com/justinstimatze/basanite) measures word and phrase
drift: its detector ranks lemmas by frequency against a baseline, and its phrase
matcher compares a curated list verbatim against the surface word stream. A
template with variable slots is invisible to both — the words are individually
unremarkable and the surface string never repeats. Nothing in that pipeline
reads sentence shape.

Its own seed data says the output stays awareness and never prohibition. It is
doing its job. It is the wrong instrument for this.

## Still open

Nothing here measures whether the prose actually got better. The judge panels
could not distinguish three candidate voices in aggregate while agreeing exactly
per criterion, so aggregate scoring is not the answer, and no per-axis harness
has been built. Reading the output for a week remains the only test that has
been run.

Ground truth without manual labelling is the idea worth trying next: corrections
in a transcript are negative examples, locatable by the message preceding them,
and the author's own prose is the positive target. Scoring two corpora rather
than one needs the material restricted to comparable work, or it flags
differences that come from doing different jobs.
