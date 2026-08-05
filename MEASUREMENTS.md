# Measurements

What was run, on how much text, and what it said. Dated, because hit rates move
when the rules or the model change. `README.md` describes the tool; this file is
the evidence the rules were chosen on.

## Baseline rate (2026-07-28)

Six transcripts, 14,178 turns, 38.9M characters. **1.29 hits per 1,000
characters.** Seven turns in ten carry at least one.

Compare hits per character, not the share of turns hit. The second number tracks
how long the turns were, so the same prose scores worse in a session of long
replies. A session drifting toward 2 per 1k is the signal worth acting on.

An earlier run of the same sweep reported 116,154 turns across 38.7M characters
and 14–24% of turns hit. The character count is comparable and the turn counts
are not: that run counted assistant records, and one reply is written as several
of them. Every turn-denominated figure published before 2026-07-28 is on a count
inflated by roughly 3× to 9× depending on how tool-heavy the session was.

## Shape detectors work, phrase detectors don't

The first run covered 36 assistant turns and 65,276 characters of the session
that produced the repo (2026-06):

| rule | hits | note |
|---|---|---|
| `bold_label` | 8 | all in the first three turns, none after a correction |
| `flip` | 3 | 2 genuine, 1 self-reference; misses every disguised variant |
| `significance_flag` | 0 | aimed at a phrasing that never occurs |
| `verdict_handoff` | 0 | same |

The rule aimed at a structural form caught everything it was meant to, and its
hit distribution matched the human complaint without tuning. The two
phrase-template rules caught nothing.

Measured 2026-07-28 over the six transcripts above:

| rule | hits |
|---|---|
| `bold_label` | 22,887 |
| `labelled_opening` | 18,335 |
| `ask_not_last` | 3,293 |
| `clause_symmetry` | 2,872 |
| `paragraph_uniformity` | 1,766 |
| `flip` | 881 |
| `dangling_end` | 103 |
| `buried_decision` | 67 |

`significance_flag` was cut earlier the same day on one hit in 20.1M characters,
for a move that happens constantly under phrasings the pattern never sees.
`verdict_handoff` went to basanite: 38 hits, and a verbatim phrase is that
tool's axis rather than this one.

## What granularity does to each rule (2026-07-28)

The largest transcript, 20.3M characters, scored twice: once per assistant
record, then once per turn, after the reader learned to group a reply's blocks
together. Same file, same rules.

| rule | per record | per turn |
|---|---|---|
| `bold_label` | 12,603 | 12,603 |
| `labelled_opening` | 11,479 | 11,442 |
| `ask_not_last` | 1,331 | 2,101 |
| `clause_symmetry` | 1,694 | 1,582 |
| `paragraph_uniformity` | 1,150 | 892 |
| `flip` | 458 | 461 |
| `dangling_end` | 350 | 36 |
| `buried_decision` | 18 | 44 |
| turns | 61,612 | 6,957 |

The local rules do not move, which is the check that the grouping did not
disturb what it should not: a bold label is a bold label wherever the paragraph
sits, and its count is identical to the digit.

Every rule that reasons about a whole reply moves, in the direction the defect
predicts. `dangling_end` loses 90% because a mid-reply fragment has
no ending, so nearly all of its old hits were the rule complaining that a
sentence stopped mid-turn. `ask_not_last` gains 58% because an ask followed by
more work is invisible when the fragment ends at the ask. `buried_decision`
more than doubles for the same reason from the other side. `paragraph_uniformity`
loses 22%, and on this repo's own session it lost 60% — a trailing scrap of four
short paragraphs reads as even in a way the reply it came from does not.

None of the three end-of-turn rules had ever been evaluated against a real
end of turn until this was fixed.

`bold_label` took two wrong guesses to get right. The tic is `**Label.**` with
the punctuation *inside* the bold, not `**Label**:`.

`ask_not_last` reported 925 hits before the 12-word paragraph floor came off it
and 1,331 after. The floor had been deleting short asks, which is most of them.

## Judge panels are reliable per criterion and unreliable in aggregate (2026-06)

Four panels, roughly 4.2M subagent tokens. The same nine drafts judged by two
different models under a Latin square:

- Per-criterion rank order replicated *exactly* across models. Clarity ranked
  the authorial blend first in both (25, 23). Specificity ranked
  selection-of-detail first in both (23, 23). Signature ranked the bare
  functional control first in both (25, 23), with the blend last at 12 in both
  runs.
- The Borda aggregate *flipped* between models (blend 58 against selection 57).

Any pipeline collapsing criteria into one score reports its choice of model.
That is why cope emits per-rule counts and no overall score.

Three earlier results that read as null were not null. The axes pull against
each other, so the sum cancelled. Completeness buys clarity and costs
signature. Terseness runs the other way.

## The mid-session reminder does not move the rate (2026-07-29)

The claim the tool was built on is that a reminder assembled from what just
fired does something a static card does not. It was tested twice, by two methods
that fail differently, and both came back flat.

**Live crossover, 238 turns across 7 sessions, 539,712 characters.** Alternate
refresher windows held back, every turn attributed to the arm live when it was
written.

| arm | turns | chars | hits | per 1k |
|---|---|---|---|---|
| hold | 109 | 263,330 | 216 | 0.820 |
| inject | 129 | 276,382 | 233 | 0.843 |

Rate ratio 1.028, 95% CI 0.854–1.237. The interval contains 1. It is also narrow
enough to say more than "no separation": an effect larger than about 20% in
either direction would have shown, and none did.

**Offline replay, 8 real prompts re-answered under each arm.** Independent calls
against reconstructed context, so nothing carries between them.

| arm | hits | chars | per 1k | vs hold |
|---|---|---|---|---|
| hold | 10 | 10,379 | 0.963 | — |
| inject | 7 | 8,837 | 0.792 | 0.822 (CI 0.313–2.160) |
| positive | 5 | 8,327 | 0.600 | 0.623 (CI 0.213–1.823) |

Both intervals span 1 on 22 hits. The ordering is the one a priming account
predicts and it is far too thin to rest on. Cost $1.03, 78% of cacheable input
served from cache.

Two methods, opposite failure modes — the crossover has carry-over between
adjacent windows, the replay has no session at all — and neither found an
effect. That agreement is the result.

### The drift reading was noise

Persona evaluation is two measurements: how a single reply scores, and whether
the voice holds across a session. A per-turn rate can stay flat while the voice
erodes. The second is what a mid-session reminder exists to act on, so the
report splits each arm at the midpoint of its session and reports the climb.

At 105 turns that split looked like the claim landing:

| arm | early | late | drift |
|---|---|---|---|
| hold | 0.727 | 1.116 | +0.389 |
| inject | 0.975 | 0.795 | −0.180 |

A difference-in-differences of about 0.57 per 1k, invisible in the mean, in the
direction the reminder predicts. At 238 turns it is gone:

| arm | early | late | drift |
|---|---|---|---|
| hold | 0.867 | 0.778 | −0.089 |
| inject | 0.859 | 0.830 | −0.030 |

Both arms drift slightly negative and agree within a rounding error. The metric
was written before that first reading was taken, which is the only reason the
sequence is worth anything — it was a prediction that failed rather than a fit
to a number already in hand. Recorded because a 0.57 gap on 105 turns is exactly
what an encouraging interim looks like, and it survived one day.

### Most of the metric is measuring reply genre, not voice (2026-07-29)

Chasing the null turned up a much larger effect that is not the arm. Split each
turn by its distance from the boundary of the window it was written in — distance
0 being the turn whose prompt carried the reminder — and the rate climbs
monotonically, in both arms:

| arm | d0 | d1-2 | d3+ |
|---|---|---|---|
| hold | 0.426 | 0.812 | 1.026 |
| inject | 0.607 | 0.802 | 0.969 |

2.4× under hold, which receives nothing. Three explanations were tested and all
three failed. Reply length: the two dominant rules repeat within a reply rather
than firing once, so their rate is flat across length (0.534 to 0.551 from the
shortest quartile to the longest). Turn tempo: inter-turn gap does relate to the
rate, but stratifying on it moves the arm ratio from 1.028 to 1.000. The
SessionStart card, which both arms receive: excluding each session's first window
makes the gradient steeper, not flatter.

What it is shows up when the rules are grouped by what they detect:

| class | d0 | d1-2 | d3+ | hits | share |
|---|---|---|---|---|---|
| formatting | 0.319 | 0.657 | 0.806 | 353 | 79% |
| rhetoric | 0.159 | 0.117 | 0.125 | 70 | 16% |
| ending | 0.035 | 0.033 | 0.065 | 26 | 6% |

`formatting` is `bold_label`, `labelled_opening` and `paragraph_uniformity` —
bold labels, verbless label openings, blocks of even length. Those are the
signature of a status report. `rhetoric` is `flip` and `clause_symmetry`, the two
that read as voice rather than layout, and they are flat.

A window opens after thirty minutes have elapsed, which lands it on a fresh
question or a change of direction. What follows is work, then a report on the
work, then more work — and reporting is where labels live. So distance from the
boundary is a proxy for how deep into an implementation stretch the turn sat, and
79% of the metric is sensitive to that.

Three things follow.

The null is explained rather than merely observed. If most of the rate is set by
what the turn was for, no reminder can move it: reporting three commits calls for
compact structure whatever the card says. The intervention was being measured
against a quantity it has little purchase on.

The drift reading is explained too. Late-half minus early-half will come out
positive in a session whose second half was implementation-heavy and negative in
one that ended in discussion. That is task composition, not voice erosion, and it
is why a 0.57 gap at 105 turns became 0.06 at 238.

And the two rules that do read as voice were stable the whole time. `flip` and
`clause_symmetry` show no gradient and no arm difference. They are 16% of the
hits, so a rate built on all eight rules is mostly reporting something else.
Whether that stability is the card working or those moves simply being rare is
not answerable from this data.

Reproduce with `cope-gate --ab-report`, which now prints both tables.

### Seven new detectors, all of them dead (2026-07-29)

Since four fifths of the metric turned out to be layout, the obvious next move
was more rules aimed at voice. The card's own negative examples are the natural
source: 13 of them across the TEST and WRONG blocks, each one a sentence the card
says not to write. Only 3 of the 13 were caught by the existing rules.

Seven candidates were written for the gaps — four regex templates in POSTPROC
(`metaphor_qualifier`, `grant_withdraw`, `straw_dismissal`, a comparative arm on
`flip`) and three shape rules in Go (`stub_colon`, `ordinal_disorder`,
`announced_ask`). Together they took the card's own examples from 3/13 to 10/13.

Then they were run over real prose. All seven failed:

| rule | on 16.3M characters | verdict |
|---|---|---|
| `metaphor_qualifier` | 0 matches | matches its own example and nothing else |
| `flip` comparative arm | 0 matches | same |
| `grant_withdraw` | 1 match | same, and that one is arguably fine |
| `straw_dismissal` | 2 matches | both false — "because that's when the mine is empty" |
| `ordinal_disorder` | 0 matches | steps out of order do not occur in the corpus |
| `announced_ask` | 0 matches | its gate needs a reply with no ask anywhere, which is rare |
| `stub_colon` | 62 matches | ~5 in 6 are list headers: `The two remaining flags:\n1.` |

Reverted. The lesson is structural rather than a matter of writing better
patterns. Aim a regex at a semantic failure and there are two settings: narrow
enough to be precise, in which case it matches only the example it was written
from, or loose enough to fire, in which case it is matching something else. There
is no middle. `flip` and `bold_label` work because the not-A-but-B template and a
markup label are surface forms; "true of anything" and "the style is the entire
content" are not.

Which is also why the working rules are 79% layout. That was never an oversight
in rule selection. Layout is the part of voice that has a surface form, so it is
the part a pattern can reach.

The three examples no rule could touch are the three the card had already sorted
into TEST blocks: SWAP ("could this sentence appear in a different conversation
with the nouns swapped"), REACTION ("am I telling the reader what to conclude"),
and EARNED FLAIR ("strip the style off — is there still something here"). Each is
a question with worked fail and pass examples, written for a reader with
judgement. The card knew which failures needed one.

### The card adds nothing a global CLAUDE.md does not (2026-07-29)

The A/B held the card constant — both arms receive the full SessionStart
injection, and only the mid-session slice varied. So the experiment that ran for
weeks never tested the card. It tested an addition to it.

The card had never been tested at all. Every transcript written without it still
carried a global CLAUDE.md with voice guidance in it, which makes the available
comparison card+CLAUDE.md against CLAUDE.md alone. That is the comparison worth
making anyway: the card's claim over a pasted file is what the whole tool rests
on.

Transcripts classify exactly, because one carrying `<voice-card>` had the hook
live. 17 do, out of 448 above 60k excluding subagents. Scored with one binary so
the rule set is constant, and measured on the rhetoric rules alone:

| group | files | turns | chars | rhetoric/1k | formatting/1k |
|---|---|---|---|---|---|
| with card | 17 | 311 | 643,039 | 0.186 (120) | 0.769 (495) |
| CLAUDE.md only | 24 | 274 | 850,270 | 0.190 (162) | 0.657 (559) |

**Rhetoric: rate ratio 0.979, 95% CI 0.773–1.240.** The interval contains 1, and
is tight enough to rule out an effect above about 25%.

Formatting runs the other way — 1.170, CI 1.037–1.320, clearing 1 — which is the
genre confound arriving on cue rather than the card making layout worse. The 17
card transcripts are recent and implementation-heavy, including this repo's own
development sessions, and implementation stretches are where status reports and
their labels live. Its size is about what the distance-from-boundary gradient
predicts, which is reasonable independent confirmation that excluding formatting
from the comparison was necessary rather than convenient.

**Why rhetoric alone, and why unpaired.** The first design paired transcripts
within project to hold genre constant, and capped at six pairs — 41 turns on the
control side, two pairs favouring each direction, no signal available at that
size. Pairing existed only to control genre, and genre only dominates because the
formatting rules are 79% of the hits. `flip` and `clause_symmetry` do not move
with genre: flat across arms, across distance from the window boundary, across
session halves. Dropping the formatting rules removes the reason for pairing and
buys 41 transcripts instead of six pairs.

**What this shows and what it does not.** It shows the rhetoric metric cannot
detect a difference. That is not the same claim as the card doing nothing, and
two hypotheses survive it: the card genuinely adds nothing over a CLAUDE.md, or
it does something these two rules cannot see. Nothing here separates them.

The design is observational and the confounds are real: different periods,
different model versions, different task mixes, no randomisation. The crossover
was an experiment; this is 41 transcripts compared after the fact.

The second hypothesis is the one that survived. Blind pairs run the following
day put the card ahead of the same CLAUDE.md at 79% over 38 pairs — see below.
The heading here still describes what these two rules can see, and what they
can see turned out not to be the thing.

### Three measurement attempts, two of them invalid

Recorded because the apparatus was wrong more often than the thing being measured
was, and each failure has a general form.

**A stale binary.** The first run scored everything with three rules that had
been reverted hours earlier. The card is embedded in the binary, so a reverted
card source and an uninstalled binary disagree in silence. The turn-log version
stamp added the same morning does not help here: it records which build wrote a
turn, not which build ran an analysis. Any script that scores prose has to print
its version and its rule list before scoring, and `cardtest2.sh` onward does.

**Genre in the control group.** The first run's control was six `agent-*.jsonl`
transcripts. Subagents write structured reports to a parent rather than
conversational replies, and they receive no SessionStart hook, so every one of
them sits in the without-card bucket and sampling by recency finds them first.
The genre confound established that morning appeared in the sampling of the
experiment testing for it.

**Pre-registering a reading the data could not deliver.** Three readings were
committed to before the paired run, without first counting how many pairs
existed. Six was the ceiling. Pre-registration is worth nothing if the design
cannot distinguish the outcomes, and the count that would have shown it took
under a minute afterward.

### What the null does and does not cover

The rate is a proxy. It counts the moves the card names in the text that was
written, not whether the writing is better, and the two can come apart. What
this rules out is the strong version — that assembling the reminder from
measured output visibly changes how often those moves appear.

Two live confounds, both of which bias toward the null and neither of which
manufactures one. Carry-over from an injected window into the held window that
follows pulls the arms together. And both arms receive the full card at
SessionStart, so what is under test is only the slicing, against a floor that
already carries every rule.

The third arm was never in this data. `positive` was built, tested, and left
behind a flag while the live experiment ran two arms for a day — the code
existed and the running binary did not have it. From 2026-07-29 the live arms
are `positive,hold`, against the same held-back reference, with the phase-1 log
kept at `turns-phase1.jsonl`.

## Blind pairs: the card beats a CLAUDE.md at 79% (2026-07-29 / 07-30)

Every section above counts moves the card names. This one asks the reader
which reply he would rather have gotten, with the answer key hidden. Real
prompts out of his own transcripts, two replies generated per prompt under
different system prompts, sides balanced rather than coin-flipped, set order
interleaved so any prefix samples both arms. Prompts already read in an earlier
run are excluded from later ones, because remembering last time's verdict is
the one form of unblinding a shuffle cannot undo. The harness is `replay/`,
documented in `replay/README.md`. Generator was `claude-opus-5` throughout.

Three runs, all `card` against `framing + the reader's CLAUDE.md Response style
section`, with `card` against `framing alone` as a positive control:

| run | card | test | control | note |
|---|---|---|---|---|
| 2026-07-29 | before the 07-30 edits | 19/34, 56%, p=0.608 | 11/17, 65%, p=0.332 | neither arm separates |
| 2026-07-30 | `bold_label` out, three QUIRKS in | 16/19, 84%, p=0.004 | not run | |
| 2026-07-30 | plus the do-not-instruct NEVER | 14/19, 74%, p=0.064 | not run | replication, all-new prompts |
| 2026-07-30 | same | not run | 10/12, 83%, p=0.039 | control re-run on the edited card |

**Pooled across the two runs of the edited card: 30/38, 79%, p=0.0005.**
Different case samples, same direction, and the gap between 84% and 74% is
inside the noise at this n. The single-run reading is what the page prints and
it is the wrong reading: 0.064 is a threshold, not a result.

The first run's flat 56% is not a fourth data point against the same card. Four
things changed between it and the two that follow: the `bold_label` ban came
out, three QUIRKS and three WRONG pairs were added, and the injection preamble
in `internal/scan/render.go` stopped telling the model the card constrained
prose only. Which of the four carries the effect is not identified, and
isolating it costs four more runs of reading.

**The mechanism is not the one the reader named.** Asked what decided the pairs
for him, he said conciseness, bullets and bold for skimmability, and a clear
closing action. Counting the replies he actually chose:

| run | arm | bullets/reply | bold/reply | chars |
|---|---|---|---|---|
| 1 | card | 0.82 | 0.44 | 1,103 |
| 1 | CLAUDE.md | 3.18 | 4.18 | 1,705 |
| 2 | card | 0.85 | 1.95 | 1,591 |
| 2 | CLAUDE.md | 1.70 | 5.70 | 2,006 |
| 3 | card | 0.42 | 0.42 | 1,046 |
| 3 | CLAUDE.md | 1.11 | 3.47 | 1,452 |

The winning side has a third to a half the bullets and a fifth the bold of the
side he rejected, in every run including the one he split 56/44. Stated
preference and revealed preference disagree, and the disagreement is not
marginal.

Length does not rescue it either. The card arm is 35% shorter in the run that
scored 56% and 28% shorter in the run that scored 74%, so shorter is not what
moved between them.

One earlier reading here was wrong and is corrected: run 2's card arm showed
bold at 1.95 against 0.44 in run 1, which read as the `bold_label` deletion
arriving on cue. Run 3's card arm is at 0.42. The 1.95 was case mix, the
deletion did not raise bold, and a two-point difference read as a trend is the
thing this file exists to stop.

**The control is what makes the test arm readable.** Card against bare framing
on the edited card is 10/12, 83%, p=0.039 — the instrument registers a hit at
the widest gap available, so a flat test arm would have meant something. It also
implies the thing nobody measured directly: 83% against nothing and 79% against
the Response style section puts that section close to nothing on this reader's
preference. A bare-versus-baseline arm would test it rather than infer it, and
has not been run.

**What this does and does not settle.** It settles that the card beats a pasted
Response style section on the only question the tool is for, which the
observational comparison above could not detect and the hit-rate proxy cannot
see. It does not identify which part of the card does it, and 38 pairs cannot
resolve anything smaller than a large effect — calling a 65% preference at 80%
power takes about 85 pairs, calling 56% takes about 540.

### A handoff-only card beats the shipped card 11 of 12 (2026-07-30)

The shipped card is about 14,000 characters and most of it legislates prose.
Every rule aimed at prose measured flat: seven candidate detectors reverted at
0–2 real matches on 16.3M characters, `flip` and `clause_symmetry` show no
gradient across any arm. What kept earning its place was the handoff — the
ending rules, and the QUIRKS and NEVER added on 07-30 that moved the pairs from
56% to 79%.

So `card/demo/handoff.effigy` keeps only that and drops the voice: no register
kernel, no flair licensing, no flip, no `clause_symmetry`, none of the SWAP,
REACTION or EARNED FLAIR tests. 5,125 characters injected against 14,011. Run
head to head with `-rival`, neither arm carrying a baseline, 12 pairs, all
prompts new against the 103 already read:

**The shipped card lost, 1 of 12, p=0.006.**

That is not the card making prose worse than no card at all: the same shipped
card beat bare framing at 83% two hours earlier. What lost is the voice axis
against the handoff axis, on a corpus where length was not the tell — 1,169
characters per reply against 1,038, trading places pair to pair.

**What is not identified.** The handoff card is not the shipped card minus its
voice half. It also carries two rules the shipped card has never had, both
written the same afternoon from preferences the reader had just stated: one
against ending on more than one actionable thing, one against pasting a document
into the reply. Some of 11/12 is the subtraction and some is those additions,
and this run cannot say how much.

### The separating arm says it was the two rules (2026-07-30)

Same handoff card, unchanged, against the shipped card with those two rules
added and nothing else changed. 12 pairs, new prompts again.

| card side | result |
|---|---|
| shipped card | 1 / 12 |
| shipped card + the two rules | 5 / 11, one tie, p=1.000 |

Two rules erased the deficit. The difference between the runs is p=0.069 by
Fisher's exact, which is suggestive rather than established, and the competing
reading is that 1/12 was an extreme draw regressing to the mean. Both readings
point the same way in practice, so the two rules are now in the shipped card and
the NEVER block holds 9 of its 10.

**The voice half is inert, not harmful.** With the handoff rules present on both
sides, 14,334 characters ties 5,125. The extra 9,209 buy nothing this instrument
can see, and they are not costing anything either — the earlier 1/12 was the
missing rules, not the present voice. Whether to carry them is a context-budget
question rather than a quality one, and at n=11 a small effect in either
direction stays invisible.

What survives all six runs is narrower than the tool: the rules that move a
reader are about how a reply hands work back, and the rules about how sentences
are built are the ones that never moved anything — flat across every arm, dead
as detectors, and now measurably worth nothing as instruction.

Two readings survive either way, and they differ in what they say to do. The
voice half may be inert, in which case it is 9,000 characters of context bought
for nothing. Or it may be actively costly — Anthropic's own guidance names
overconstraint, many directives arriving in one request, as the Opus 5 failure
mode, and 14,000 characters of prose legislation is a candidate for it.

**It also un-settles the reminder null.** The mid-session experiment above is a
null on hits per 1k, and hits per 1k is 79% formatting rules that track reply
genre. The rate cannot see what the pairs just measured, so that experiment does
not show re-injection failing — it shows the rate not moving, which is a
different claim about a metric now known not to track the reader. Whether the
third hook earns its place is open, and the pairs harness is what would answer
it: a run whose arms are the full card against the card with `--refresher`
disabled.

## Not measured

The rate is not the thing. Every number outside the blind-pairs section counts
moves the card names, not whether the writing is better, and the pairs are the
only place the second question gets asked directly. They answer it for one
reader, over 38 pairs, against one comparator.

**Whether the voicing arrives at all has never been measured, by anything on
this page.** Preference and identity are different axes. The pairs ask which
reply reads better, and both arms of every run above were Claude writing well —
so a flat result there is silent about whether a voice was swapped, and reading
it otherwise is how a null at n=11 turned into "the voice half is inert" when
the same paragraph says a small effect stays invisible at that n. The dead voice
detectors say no more: seven regexes failing to catch Claude's tics is a fact
about regex sensitivity.

The instrument for it now exists — `replay/ -discriminate`, which shows a reader
one card's VOICE block and asks which of two replies carries it. Its first real
run is below.

Two further gaps.

Completeness and symmetry look separable, and that is an assertion. The
signature judges never penalised the blend's stated functions; every quoted
offence was matched-clause symmetry — "a costume that fits well… a costume that
fits badly", "nothing to rank / nothing to match". If both terms of a comparison
can be named without being balanced, the clarity win and the signature loss come
apart. The test would vary symmetry alone with completeness fixed. It has not
been run, and nothing in the tool depends on the claim.

Ground truth without hand labelling is the idea worth trying next. Corrections
in a transcript are negative examples and can be found by the message preceding
them; the author's own prose is the positive target. Scoring two corpora rather
than one needs the material restricted to comparable work first. Otherwise what
it flags is the two sources having been written for different purposes.


## Rules taken from other people's complaints

The rules before this section were mined from this repo's own transcripts. These
three came from Zvi Mowshowitz's July 2026 Opus 5 reaction roundup, which
collects named practitioners on what makes the model unpleasant to talk to. The
same three moves recur there across otherwise unrelated complaints: apology,
narrating its own errors, and announcing that a long passage is coming.

Backfilled over 20,911 assistant turns and 52,599,321 characters from other
projects' sessions.

| rule | hits | verdict |
| --- | --- | --- |
| `apology` | 25 | every one genuine contrition; in the card |
| `self_postmortem` | 1 | the one match described somebody else's remark; detector only |
| `announced_length` | 3 | all three described a room in a game, not the reply; detector only |

The two nulls are not the same claim as "the move does not happen". This corpus
is Claude Code under a global CLAUDE.md that already bans preamble and forbids
ruminating on mistakes, and the roundup's complaints come from chat. What the
null rules out is the move appearing in this writer's setup often enough to
score. Both detectors stay so the question can be answered later by a corpus
that has not already been filtered.

Four false positives surfaced only under backfill, after the rules passed their
own fixtures — the same order of events as `forked_end`. A caveat about
world-visible file permissions matched "for honesty". An offer to talk someone
through operator steps matched "I'll walk you through". A one-sentence
correction, which is the behaviour the card asks for, matched "my mistake was".
And "dense" matched three descriptions of a location. Each is now a fixture.

## The repetition rule, and the corpus bug it found (2026-07-30)

Repetition is the most-cited complaint in the roundup and the only one a
per-reply check cannot see: a construction used for the fourth time reads
identically to one used for the first. `cross_turn_repeat` reads the window
instead. Five-word shingles, kept when four of the five words are common ones,
fingerprinted to 32 bits and stored on the turn; a phrase already in two earlier
turns fires the rule.

Backfilled over 27,610 assistant turns and 72,839,042 characters from 80
sessions in other projects, each reply scored against the twenty before it.

| | |
| --- | --- |
| turns with a hit | 1,369 (4.96%) |
| sessions with any hit | 47 of 80 |
| median session rate | 1.8% |
| median rate, sessions of 500+ turns | 1.6% |
| worst session | 918 of 6,927 (13.3%) |

For scale, `dangling_end` was 10.1% of Opus 5 endings and `forked_end` 5.4%.

The worst session is the rule working, not failing. Every phrase it printed
there is the same preamble family — "let me look at the", "let me check how
the", "let me check for any", "now let me run the" — in a 6,927-turn session
where the global CLAUDE.md has banned preamble the whole time. The next worst,
at 10.9%, gave "say the word to keep" and "so the next session doesn't".

Across the whole corpus the sampler printed 131 matches. Roughly three in four
are constructions: 43 in the `let me ...` family, 21 in `want me to ...`, 7 in
`say the word ...`, plus "nothing further to do here", "i have everything i
need", "from my last message is". The other quarter is subject matter that got
through the four-of-five filter — "berkeley is the other one", "on long runs
where the", "the go module now has". That is the rule's failure mode and it has
no fix short of knowing what a sentence means: a five-word span carrying one
content word can be a frame or a fact, and only the topic says which.

### The corpus had been scoring the harness

The second most repeated phrase in the first pass was "temporary try again in
a", 18 sampled matches across 9 sessions. It is not a phrase the model wrote.
Claude Code stores its own API error banner as an assistant record —

> API Error: 500 Internal server error. This is a server-side issue, usually
> temporary — try again in a moment. If it persists, check status.claude.com.

— once per retry, identically, in the reply position. Two markers separate
those from real replies: a top-level `isApiErrorMessage`, and `message.model`
set to `<synthetic>`. `internal/transcript` now skips both.

Every rate this file records before today included them. The correction is
small — 27,962 turns to 27,610, and `cross_turn_repeat` from 5.04% to 4.96% —
because banners are short and most rules need prose to fire. It is worth
recording anyway: nothing else in the repo would have surfaced it, and it took
a rule that reads across turns to notice a sentence that was identical every
time.

## The loop lane (2026-07-31)

Addy Osmani's June 2026 "Loop Engineering" names the shift this measures: "Loop
engineering is replacing yourself as the person who prompts the agent. You
design the system that does it instead." Boris Cherny, quoted there: "I don't
prompt Claude anymore. I have loops running that prompt Claude and figuring out
what to do. My job is to write loops".

The relevant consequence for a voice gate is that the reply stops being a
message and becomes a record. `/loop` re-runs on a cadence and `/goal` runs to a
condition; either way the reader is asleep. Two of Osmani's lines set what a
rule for that lane has to check — "a loop running unattended is also a loop
making mistakes unattended", and, on the maker/checker split, "'done' is a claim
and not a proof."

So `transcript.Turn` carries a lane, and `Card.CheckLane` drops the four rules
whose subject is a reader who can answer (`ask_not_last`, `forked_end`,
`dangling_end`, `buried_decision`) and adds two that assume one who cannot:
`loop_ask` and `unverified_done`.

Backfilled over 34,543 assistant turns from 105 sessions in other projects.

| lane | turns | turns with a hit | mean length | hits per 1k chars |
| --- | --- | --- | --- | --- |
| interactive | 33,602 | 62.4% | 2,757 chars | 0.23 |
| loop | 941 | 21.1% | 427 chars | 0.50 |

Per turn the loop lane looks three times cleaner. Per thousand characters it is
twice as dense. The per-turn number is mostly reply length: loop replies are a
sixth the size, and most rules in this repo need prose to fire at all. One
session's loop replies averaged 37 characters, which is not a report, it is
"nothing to do".

The rate is not a like-for-like comparison — the two lanes run different rule
sets by construction — so the number to read is the examples.

### What the loop lane actually does

`loop_ask` is the rule that fires, and its matches are the finding:

> Want me to shorten it, change the tone, or adjust anything?

> What's next — DESIGN.md writeup for all three tracks, start on visualization,
> or something else?

> Given the hour — want to push into Part 2 now, or stop here for tonight?

> Stopping the loop — work is committed, tests pass, awaiting your call on
> deploy timing.

Every one of those was written into a log at an hour nobody was reading it. The
question does not reach a person; it reaches the next iteration, which has to
decide what its own predecessor meant. This is the thing a lane-blind gate
cannot see, because each of those endings is well-formed by the interactive
rules — several of them are what `dangling_end` asks for.

Broken out over the 870 loop-lane turns in the sessions that use loops at all:

| rule | hits |
| --- | --- |
| `labelled_opening` | 138 |
| `loop_ask` | 55 |
| `paragraph_uniformity` | 26 |
| `clause_symmetry` | 26 |
| `cross_turn_repeat` | 11 |
| `flip` | 8 |
| `unverified_done` | 1 |

The two lane-specific rules are 21% of the hits, so the density above is not an
artefact of adding them. `labelled_opening` dominating is the reply-genre
confound this file already records: a status note is written in labels.

`unverified_done` fired once. It keeps its detector and loses its card rule,
which is the same call `self_postmortem` and `announced_length` got — and here
there is a second reason, since the card already carries an always-on "never
describes work as done, working, or verified when it was written but not
executed". A gated rule that restates an ungated one spends budget to say the
same thing twice.

The honest limit: 941 loop turns against 33,602 interactive ones is a 2.7%
sample, from eight sessions, and one of those sessions supplies 422 of them. A
rate from that is a rate from one person's loops, not from the lane.

## The voicing arrives: 90% discrimination (2026-08-02)

Every section above this one asks whether the writing is better. This asks
whether it is *different* — the reader is shown one card's VOICE block and picks
which of two replies carries it. Chance is 50%.

Sixty pairs, `claude-opus-5` generating both sides, real prompts from card-free
transcripts. The shipped card against `card/demo/laconic.effigy` for the test
set, against bare framing for the control.

| set | comparator | correct | rate | p |
|---|---|---|---|---|
| test | laconic | 36/40 | 90% | 1.9e-07 |
| control | bare framing | 18/20 | 90% | 0.0004 |
| pooled | | 54/60 | 90% | 9.7e-11 |

Split by which voice was named:

| target named | correct |
|---|---|
| laconic | 20/20 |
| the shipped card | 34/40 |

Every one of the six misses came when the card was the target. Spartan brevity
is unmistakable; "plain prose that costs the reader nothing" sits nearer to
default Claude, and is found rather than spotted.

**This is the model judge, not the reader, and the two have not been compared.**
The harness prints agreement between them precisely because a judge's rate is a
fact about a model until that number exists. Read it as: a model given a voice
description and two replies picks correctly nine times in ten. Whether a person
does, and whether either is reading register rather than length, is unmeasured.

The balance held — the answer fell on the left in exactly 20 of 40 test pairs
and 10 of 20 control, with the target split 20/20 between the two cards. That is
the fix from the pilot run, where the answer landed right in five of six.

What this does *not* say is that the card is good. It says a card changes the
prose enough to be identified from its own description, which is the
precondition every voicing claim above needed and none of them had.

## The rival with no gimmick: 95% (2026-08-02)

The section above has a hole in it, and the last paragraph names it without
following it: whether either reader was reading register rather than length.
`card/demo/laconic.effigy` answers in four words. A judge shown two replies where
one is four words long can score 20/20 on the laconic target without ever
forming a view about voice, and the run cannot tell that apart from the voice
having arrived.

So this run swaps the rival for one with no surface tell.
`card/demo/lecturer.effigy` is ordinary English at ordinary length, written for
this test, differing from the shipped card on one axis — whether the prose is
allowed to call attention to itself. The two descriptions the reader holds:

- shipped: *No lift at the end of a paragraph, no summary of what was just said.*
- lecturer: *The paragraph has been climbing the whole way, and the last four words are the top of it.*

Same generation as before, `claude-opus-5` on both sides, real prompts from
card-free transcripts.

| set | comparator | correct | rate | p |
|---|---|---|---|---|
| test | lecturer | 38/40 | 95% | 7.5e-10 |
| control | bare framing | 18/19 | 95% | 3.8e-05 |
| pooled | | 56/59 | 95% | 6.0e-14 |

Split by which voice was named:

| target named | correct |
|---|---|
| lecturer | 20/20 |
| the shipped card | 36/39 |

It went up. Removing the surface cue raised the rate from 90% to 95%, which
retires the confound rather than confirming it — whatever the judge was reading,
it was not article count or reply length. The asymmetry from the previous run
survived intact: every miss came looking for the shipped card, and the rival
scored perfectly both times.

The stronger reading is in the rival column. `lecturer.effigy` was written from
scratch an hour before the run, by a notation the gate has no detectors for,
against a register nothing in this repo was tuned toward — and 20 of 20 replies
written under it were identified from its own description. That is the framework
claim, which every measurement before this one had left untested: the previous
runs all asked whether the shipped card differs from not-the-shipped-card.

This is still the model judge, and the human agreement number the harness prints
is still unread. The case sample is also not the same one — 40 transcripts
yielded 58 usable cases where 60 were needed, so this run drew from 120. The two
runs are separate samples of the same corpus rather than a paired comparison,
and the 90-to-95 gap is inside what that could produce on its own.

One pair is missing from the 60. The judge returned no text on it and the
harness dropped it rather than scoring it as a call, which is the guard added
after the pilot run counted empty answers as unsure.

### The question the judge was asked, and the wording that changed after

Both runs above were produced under a judge instruction that said one reply was
written under the named voice "and the other was not". That is true of the
control arm and false of every test pair, where both replies were written under
a card and only one was written under the named one. It points a judge at
whichever reply reads as more voiced rather than at the voice it was given.

The bias runs one way and it is the useful way. The shipped card is the quiet
voice in both pairings, so the wording pushed against the card-as-target arm and
toward the rival — which makes 36/39 and 34/40 the harder-won numbers and the
rival's 20/20 the softer one. Neither result changes sign under it.

The instruction now says the other reply "was written under different
instructions, which may include a different voice". Runs after that wording are
not comparable to these two, and the numbers here should not be pooled with
them.

## The scorer does not swap when the card swaps (2026-08-02)

`shapes()` in `internal/scan/scan.go` takes no card. The regex rules come from
the card's POSTPROC block and the rest of the rules are Go, so loading a
different card changes what gets written and leaves what gets checked exactly
where it was. That much is a fact about the code and needs no run behind it.

What a run can say is whether it matters. The 60 discrimination pairs give
three sets of replies to the same prompts — the shipped card, `lecturer`, and a
bare framing — so scoring all of them with the same binary asks what the gate
makes of prose it was not written for.

Split by what each rule is sensitive to. Voicing is `clause_symmetry`, `flip`
and the three deference rules, normalised over prose characters because those
rules skip list blocks. Structure is `labelled_opening`,
`paragraph_uniformity`, and the four ending rules, over the whole reply.

| written under | replies | chars | voicing / 1k prose | structure / 1k |
| --- | --- | --- | --- | --- |
| bare framing | 20 | 25,659 | 0.10 | 1.13 |
| the shipped card | 60 | 59,661 | 0.20 | 0.70 |
| lecturer | 40 | 71,145 | 0.39 | 0.34 |

The two columns run in opposite directions, which is what the split predicts
and the reason to keep the two axes apart. Neither leg is a result on its own.

Voicing, shipped card against lecturer: rate ratio 1.93, p=0.07 by exact
conditional Poisson. `clause_symmetry` alone is 0.13 against 0.25, ratio 1.98,
p=0.15. Directionally this is the thing worth naming — lecturer's VOICE block
asks for a cadence in so many words, and `clause_symmetry` exists to catch that
cadence, so the card is marked down for doing what it was told — but 7 hits
against 17 does not separate and should not be quoted as though it did.

Structure looks stronger at p=0.004 and is worth less. It reverses under a
different denominator: per 1000 characters the rival is twice as good (0.34
against 0.70), and per reply the two are almost level (0.60 against 0.70).
Lecturer's replies average 1,779 characters against the shipped card's 994, and
that alone moves the number by more than the cards do. Nothing here separates a
card effect from a length effect.

The prose share is not the confound. It runs 80% bare, 92% card, 95% lecturer,
which is the right direction to explain part of the voicing gap and about a
sixth the size of it.

So: one code fact, one direction worth stating with its n attached, and one
number that is mostly reply length. The claim the README makes from this is the
code fact.

### What would settle it

The rules would have to be authorable. A card that could name which structural
rules it wants — or supply its own — turns this from a confound into an arm: run
the same 60 prompts with lecturer scored under lecturer's rules and under the
shipped card's, and the difference is the thing being argued about. Until then
every card is measured against one card's idea of how a reply should end.

## Where the card sits beats what it says (2026-08-03)

Every run above this line came from `tools/generate_readme.py`, which sends the
card as the entire system prompt. Deployed, cope did not do that: the card went
through a `SessionStart` hook — one message at turn zero, which the conversation
buries under everything written after it.

27 samples, 22 usable after five empty responses, counting bolded-label
paragraphs per reply against a card that asks for a bolded label on every
paragraph.

| where the card sat | bolded labels per reply |
| --- | --- |
| the system prompt | 9.4 |
| the system prompt, next to a global CLAUDE.md | 6.3 |
| a global CLAUDE.md, no card | 1.0 |

63–86% of the effect survives sitting next to other instructions, so what was
killing the hook delivery was position and not interference. Claude Code's
documentation says an output style's instructions are added to the system prompt
and that the harness re-reminds the model of them during the conversation; a
turn-zero message gets neither.

Run through the hook, the same card — bolded label on every paragraph, emoji on
every heading — produced prose indistinguishable from no card at all.

**What this does not say.** A bolded-label count is a compliance proxy, picked
because it is countable and unambiguous, not because markup is what a card is
for. It says an instruction arrived. It says nothing about whether the writing
got better, which is what the discrimination runs above are for.
