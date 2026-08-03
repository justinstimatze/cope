# Blind pairs: testing a card by reading, not by counting

`cope-gate` scores prose by regex and reports hits per 1,000 characters. That
number is a proxy, and when it was audited four fifths of what it caught tracked
what a turn was *for* — a status report, a plan, an answer — rather than how it
was written. This directory holds the instrument that asks the question directly.

A run takes real prompts out of your own transcripts, generates two replies to
each under different system prompts, and hands them to you unlabelled. You pick
which reads better. The answer key stays in `pairs.json` and the page never shows
it until you press reveal.

## Running it against your own card

```
make pairs BASELINE=~/.claude/response-style.md CARD=path/to/your.effigy
```

`BASELINE` is required and is whatever prose guidance you already have without a
card — for most people the `## Response style` section of `~/.claude/CLAUDE.md`,
saved to its own file. The card is tested *on top of* the baseline, so the run
measures what the card adds rather than what any instruction at all adds. Leave
`CARD` off to test the card built into the gate.

Three conditions exist:

| arm | system prompt |
|---|---|
| `bare` | the framing alone |
| `baseline` | framing + your baseline file |
| `card` | framing + baseline + the card |

`card` vs `baseline` is the question. `card` vs `bare` is the positive control:
the widest gap available, there to show the instrument can register a hit at all.
If the control separates and the test does not, the card adds nothing over what
you already had. If neither separates, the instrument is blind and no number it
produced means anything.

A deliberately extreme card is the cheapest sanity check there is. If a caveman
card does not separate from your baseline, stop and fix the harness before
reading anything else it tells you.

## What it can and cannot resolve

Twenty pairs detects a large effect and nothing smaller. Calling a 65%
preference at 80% power takes about 85 pairs; calling 56% takes about 540. Read a
flat result as "no large effect", never as "no effect" — and if you want to
resolve something small, the cost is measured in hours of your reading, not in
API spend.

## Contamination, and the four filters that exist because of it

- **Cases come only from transcripts with no `<voice-card>` in them.** Prior
  assistant turns anchor register hard, so history written under the card would
  seed both arms with the thing being tested.
- **`-exclude` takes pairs.json files you have already read.** Re-sampling a
  prompt carries your memory of last time's verdict into the new pair, which is
  the one form of unblinding a shuffled side assignment cannot undo.
- **Sides are balanced, not coin-flipped.** A flip is unbiased in expectation and
  lopsided in practice at this n, and any left/right reading habit would turn
  that lopsidedness into an effect.
- **Set order is shuffled.** Written in set order the control sits entirely
  behind the test, so anyone who stops early has judged no control at all.

Prompts that are hook output, task notifications, or fewer than three words are
dropped, as are cases whose prompt or history matches a credential pattern. A run
carries real conversation out of every project on the machine — `-outdir` is
gitignored for that reason and should stay outside any repo you publish.

## The other question: did the voicing arrive

Preference and identity are different axes, and the page above only reads one of
them. A voice can land completely and leave a reader indifferent about whether
it was worth landing, so a flat preference result says nothing about whether the
voicing reached the prose. Reading it as though it did is how a null at n=11
became a claim that voice guidance does nothing.

`-discriminate` asks the other question. Same generation — two replies to the
same real prompt under two conditions — but the reader is shown one card's VOICE
block, the aim and the register taken out of the card itself, and asked which
reply carries it.

```
make discriminate RIVAL=card/demo/laconic.effigy
make discriminate CARD=card/demo/precise.effigy RIVAL=card/demo/laconic.effigy
```

Chance is 50%, so the arithmetic is the same exact binomial. What differs is
what a null means. Two replies can be equally good in a hundred ways, but two
that cannot be told apart under a stated target had no voicing happen to them,
whatever the card said.

| set | pairs | asks |
|---|---|---|
| control | card vs bare framing | can the card be told from no guidance at all — the widest gap there is |
| test | card vs a second card | can two deliberate voices be told apart from their descriptions |

The named target alternates between the two cards across the test set, so a
reader cannot answer everything by recognising one voice and calling its
absence, and so both cards get measured rather than one being the permanent
subject. The answer's *side* is balanced too, which is a separate thing from the
preference page's card-side balance: the target is the rival in half the pairs,
so balancing where the card sits leaves the answer free to pile up. It did —
right in five pairs of six on the first live run, free marks for anyone with a
side habit.

The description comes out of the card rather than being written for the test. A
description written separately is one more artefact that can be good or bad on
its own, and a reader who then failed to match it would leave two readings
standing: the voice did not land, or the description was poor.

### The model judge

`-judge` answers every pair with the model, on the same materials and the same
question. Reading is the bottleneck this exists for: a person calls thirty pairs
in an evening, thirty pairs resolves a large effect and nothing smaller, and a
framework meant to carry arbitrary voices needs the loop shorter than one
revision per evening.

It is not trusted on its own. It answers the pairs you also answer, and the page
reports the agreement — until that number has been looked at, the judge's rate
is a fact about a model rather than about a voice.

One thing to know if you change the call: Opus 5 emits a thinking block ahead of
its text whether or not it was asked for. At a token ceiling sized for the
one-word answer, thinking spends the whole budget and no text arrives — which
surfaced as six pairs of "unsure", indistinguishable at a glance from a judge
that genuinely could not tell two voices apart. A missing answer is now an error
and is dropped from the run; "unsure" is reserved for a call the judge made.

## Files

- `pairs.go` — sampling, generation, the self-contained HTML page
- `discriminate.go` — target balancing, the model judge, the discrimination run
- `discpage.go` — the discrimination page
- `main.go` — flags, and the older per-arm scoring replay this grew out of
- `transcript.go` — reading cases out of Claude Code session JSONL
