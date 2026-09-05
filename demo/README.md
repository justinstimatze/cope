# Demo — the same facts, a different card each time

Every file here is `README.md` written again from a different voice card,
through the same prompt (`cope-gate --author-docs`) and the same introspected
facts. Only the card changed. Build them with `make demo`.

## Open two of them

[`README.claude-maximal.md`](README.claude-maximal.md) and
[`README.claude-voice.md`](README.claude-voice.md), side by side. That is the
whole demo, and it takes about thirty seconds. Same facts, same prompt, same
generator, same model. What you can count on the page:

<!-- side-by-side:start -->
| | `claude_maximal` | `claude_voice` |
|---|---|---|
| paragraphs opening on a bold label | 57 | 0 |
| em dashes | 82 | 49 |
| headings carrying an emoji | 14 | 0 |
| horizontal rules | 20 | 1 |
| tables for things that are not tabular | 5 | 0 |
| blockquote callouts | 3 | 0 |
| emoji in the body text | 16 | 0 |
| sentences apologising or taking the blame | 18 | 0 |
<!-- side-by-side:end -->

The repository's front page is a demo render on purpose —
[`../README.md`](../README.md) is written from `card/demo/fieldguide.effigy`,
because a reader has to be able to see that a card wrote the page, and a page in
the shipped register cannot show that.

It was the maximal card until 2026-08-04, on the argument that a page written in
the register cope exists to fix states the problem better than a paragraph about
it can. It does, and it costs too much: the people who arrive here arrived to
get away from that prose, and the front page served them a full page of it
before they had any reason to trust the joke. The maximal render is still here
and the front page still names it, so it is met by clicking rather than by
landing.

## Does it do the thing

Reading two pages tells you they differ. It does not tell you a *named* voice
arrived, which is the claim worth checking, and looking at your own output will
not settle it either — you wrote it, so you cannot be blind to it.

The test that can: a judge is shown one card's description of itself and two
replies, and picks which reply was written under it. Chance is 50%. Against
`lecturer` — ordinary English at ordinary length, differing on one axis only, so
there is no surface tell to win on — it picks correctly 38 times in 40, and 18 of
19 on the control. `MEASUREMENTS.md` has the run, the balance checks, and the
three earlier attempts that were invalid.

What that does not say is that the card is good. It says a card changes the
prose enough to be identified from its own description, which is the
precondition every other claim needs.

## The other four

`precise` compresses by vocabulary, naming mechanisms with their terms of art
and glossing none of them. `caveman` compresses by grammar, dropping articles,
copulas and auxiliaries. `lecturer` does not compress at all; it differs from
the shipped card on register alone, which is what made it usable as the rival
above.

`fieldguide` is the front page. Entries open on the name of the thing, an
alternative gets its own sentence beginning "Compare", and a claim with no
observation behind it reads "Not recorded". Its first draft failed and the
failure is written into the card's header: that draft described the job — tell
things apart — and came back indistinguishable from the shipped card, carrying
none of its own moves. A job is something a model can believe it is doing while
writing normally, so the instructions were rewritten at the level of the
sentence opening and the exact word a comparison starts with. `caveman` was the
other candidate for the slot and was dropped: there is a separate project of
that name by another author, cited on the front page, and wearing a
neighbour's register on the front door reads as taking it.

`card/demo/laconic.effigy` is installed by `make cards` and no longer rendered.
It is the Spartan register, after
[Laconia](https://en.wikipedia.org/wiki/Laconic_phrase) — answer what was
asked, never explain the answer, stop — and it was the original discrimination
rival until `lecturer` replaced it, because a four-word answer lets a judge
score full marks without forming a view about voice. A page written from it
demonstrated brevity, which its own first line already says.

<!-- sizes:start -->
| Card | Output | Characters |
|---|---|---|
| `card/claude_voice.effigy` | [`README.claude-voice.md`](README.claude-voice.md) | 26,172 |
| `card/demo/claude_maximal.effigy` | [`README.claude-maximal.md`](README.claude-maximal.md) | 39,718 |
| `card/demo/precise.effigy` | [`README.precise.md`](README.precise.md) | 24,615 |
| `card/demo/caveman.effigy` | [`README.caveman.md`](README.caveman.md) | 23,700 |
| `card/demo/lecturer.effigy` | [`README.lecturer.md`](README.lecturer.md) | 30,153 |
| `card/demo/fieldguide.effigy` | [`README.fieldguide.md`](README.fieldguide.md) | 26,072 |
<!-- sizes:end -->

`card/demo/handoff.effigy` is in the directory and not in the table. It is a
hypothesis rather than a voice — the shipped card's handoff rules with the prose
half deleted — and it is meant to be run through `make pairs` against the full
card rather than rendered.

---

Everything below is about the scorer rather than the prose. It is the
development instrument: useful for deciding what a rule should catch, and not
the reason to install anything. The short version is that the gate is blind to
most of what you just looked at, and the rest of this page is the argument for
why that is worth knowing.

## Where the joke card comes from

`claude_maximal` instructs every tic this model has been measured to have, and
the sourcing is in its header: a first draft from
[claudisms.ai](https://claudisms.ai/), then checked against
[basanite](https://github.com/justinstimatze/basanite)'s measured tic list,
which is not folklore — basanite measures lemma frequency over real transcripts
and has audited its list against roughly 4M tokens. The two disagree. "At the
end of the day", "you're absolutely correct", "liminal" and six more were cut
from the measured list for occurring nowhere outside writeups about how Claude
supposedly talks. Several are on the public banlist. The card instructs only
what survived, which is the same trade the rest of this repo makes.

## The gate is blind to the markup and not to the words

`claude_maximal` scores several times the shipped card. That is a recent
change and the way it changed is the useful part.

<!-- scores:start -->
| Page | scored by the shipped card | scored by its own card |
|---|---|---|
| [`README.claude-voice.md`](README.claude-voice.md) | 9 | the same card |
| [`README.claude-maximal.md`](README.claude-maximal.md) | 28 | 16 |
| [`README.precise.md`](README.precise.md) | 12 | 10 |
| [`README.caveman.md`](README.caveman.md) | 19 | 19 |
| [`README.lecturer.md`](README.lecturer.md) | 11 | 8 |
| [`README.fieldguide.md`](README.fieldguide.md) | 23 | 6 |
<!-- scores:end -->

Read that table in two columns rather than one. The left one runs the shipped
card's rules over every page, which is what a fixed ruler says. The right one
runs each page against the card it was written from, which is what the card's
own word was worth. `fieldguide` is where they part hardest: it declines
`labelled_opening` in its header, because every entry opens on the name of the
thing, and a declined rule still runs and still counts on the left. Two
numbers, one page, and the gap between them is the whole reason this section
exists.

For most of this card's life it scored **clean**, then 2, while instructing every
habit it was built from. Everything it did was markup — a bolded label, an emoji,
a horizontal rule, a table where a sentence would have done — and none of that is
on either axis cope reads. The demo page said so, and used it as the argument for
reading a score as a report rather than a verdict.

Then it was told to apologise in every paragraph, to set up each claim by first
denying its opposite, and to reach for `load-bearing` and `worth noting` by name.
None of that is markup. It is the prose, and the gate found all of it:

<!-- maximal-hits:start -->
| rule | hits on `claude_maximal` |
|---|---|
| `flip` | 8 |
| `clause_symmetry` | 8 |
| `repeated_opening` | 7 |
| `load_bearing` | 2 |
| `worth_noting` | 2 |
| `apology` | 1 |
<!-- maximal-hits:end -->

So the exclusions below are real and they are narrower than they look. A card can
be unmistakably this model on the page and score clean, as long as what makes it
unmistakable is formatting. The moment the same card reaches for the model's
sentences rather than its markup, the gate has it — and every bolded label,
emoji heading and gratuitous table in the first table on this page still passes
through untouched:

- **Bold labels.** `labelled_opening` skips the bolded form. The card carried a
  `bold_label` rule and dropped it, after a reader working through 52 blind
  pairs named bold and bullets among the things deciding which reply he
  preferred.
- **Bulleted lists.** `paragraphs()` skips list blocks outright, because a list
  is legitimately uniform and would poison the uniformity test.
- **The one-line closes.** `paragraphs()` skips anything under twelve words. The
  card asks for a closing line under twelve words. It lands under the floor.
- **Em dashes, emoji, horizontal rules.** No rule at all. cope reads sentence
  and paragraph shape; markup is not on either axis.
- **The lexical tics.** `load-bearing` and `worth noting` are now POSTPROC rules
  on the shipped card, added after measuring them at 25.6 and 6.5 per 1k across
  10,652 replies. The rest of that family stays with basanite, which reports a
  rate rather than banning a word that is sometimes right.

Each of those was a defensible call made for a stated reason on a stated date.
Together they describe a gate that cannot see this model's markup, however thick
it is laid on. [`README.claude-maximal.md`](README.claude-maximal.md) is the longest page
here and it is exactly that, and it costs the card nothing: every one of its
hits comes from a sentence rather than from anything you can see at a glance.

## Most of the gate never fires on a document

<!-- totals:start -->
Across 6 renders and 170,430 characters, 8 of the 19 rules fired at all: `labelled_opening` 36, `clause_symmetry` 29, `repeated_opening` 19, `flip` 12, `load_bearing` 2, `worth_noting` 2, `apology` 1, `ask_not_last` 1.
<!-- totals:end -->

Two of them carry it and the rest are a long tail. Nothing else in the gate
found anything at all.

That says more about the artifact than about any card. Four of the structure
rules — `ask_not_last`, `dangling_end`, `buried_decision`, `forked_end` — read
the ending of a conversational turn for whether the reader has something to
answer, and the section list forbids manufacturing one. `unverified_done` and
`loop_ask` belong to the loop lane, which a document is not in.
`cross_turn_repeat` needs a session history. What is left able to fire on a
README is almost exactly what did.

Sixteen counts thirteen compiled into the binary plus the shipped card's
POSTPROC patterns. Which of those fire, and how often, moves from render to
render — the line above is this build's answer and the next build's will differ.
`flip` is the one that reliably shows, on the joke card, after it was told to
set up each claim by first denying its opposite. That the phrase rules barely
register here is worth holding against their measured rates in real transcripts:
a README is not the artifact they were chosen on.

`apology` is the one to notice, because it wants a reply with something to be
sorry for and a README has nothing. It fires anyway, once, against the apologising-sentence count
in the first table on this page. One in that many is a poor catch rate on a
habit the card was explicitly instructed to perform, which is worth knowing
before trusting the absence of a hit anywhere else.

So read the score as separating cards on voicing and nearly silent on the other
axis. Then read the two rules that carry it, because they separate the cards for
opposite reasons.

## Where each card's hits come from

`caveman` collects most of its hits on `labelled_opening`. Drop the copula and
the opening comes out verbless, every time. The detector is not wrong to find
them, and finding them tells you nothing you did not know from the card's first
line. It is the same flaw that made caveman useless as a discrimination rival:
identifiable from morphology in three words, so anything measuring it measures
article count. The hits are real and they are not about structure.

`fieldguide` collects the same rule, turned up higher, and unlike caveman it
says so in its header. The card requires every paragraph to open on the name of
the thing it describes. It declines the rule with `@gate` and gives the reason
on the line, which is why its two columns in the score table are so far apart.

`lecturer` collects mostly `clause_symmetry` — two clauses of matching length
with a content word carried across the comma, the climb the card's VOICE block
asks for, in as many words. `clause_symmetry` is compiled into the binary to
catch it, and for a month the card had no way to say it wanted it, so it was
marked down for doing what it was told. That is what `@gate` was built for.
`lecturer` now declines that rule and `dangling_end` in its header, with the
reason on the line, and its own column drops accordingly.

The other direction is `@shape`, where a card states a structural rule and the
gate checks it. `card/demo/handoff.effigy` uses it for the one claim its own
peak makes and no built-in rule tests — that the closing block can be read cold
— as `last paragraph words <= 60`. What is left is the vocabulary: `@shape`
counts words and sentences and asks whether a block poses a question, so a
check like `clause_symmetry` stays compiled in and a card wanting something
outside that and a POSTPROC regex still has nowhere to put it.

## Why a card written to fail scores well

The gate reads whatever `paragraphs()` returns, and that function skips list
blocks, anything under twelve words, and paragraphs opening on a bolded label.
A card whose whole manner lives in those exclusions scores low whatever it
instructs — which is what `claude_maximal` demonstrates, and why the score
separates cards on a narrower thing than it looks like it does.

The general form: cope's rules are aimed at conversational replies, where the
model chooses the shape itself — where the ask lands, what the closing line
leaves the reader to do, how long a paragraph runs before it stops. A document
written to an outline has those decisions made for it before the card arrives.

## Compression by vocabulary, and where that stopped holding

The common way to shorten agent output is to strip function words and leave
fragments. That lowers word count while leaving concept density flat, so a
precise idea still costs a clause. Naming the concept instead compresses
further. Counted with Anthropic's tokenizer on one sentence:

| version | tokens |
|---|---|
| ordinary explanation | 96 |
| telegraphic fragments | 40 |
| precise vocabulary | 31 |
| precise, one clause | 24 |

The first two are the before/after from the README of
[caveman](https://github.com/JuliusBrussee/caveman), a separate project that
compresses agent replies to cut output tokens. The second two are the `precise`
card's register applied to the same content.

At document scale it reverses, and it has now reversed twice on rebuilt files.
`README.caveman.md` is the shortest in the size table above and
`README.precise.md` runs longer. Grammar compression applies to every sentence;
vocabulary compression only pays where there is a term of art to reach for, and
most of a README is not that.

## Measuring the card properly

Run it against transcripts rather than documents, where the model is choosing
the shape:

```
cope-gate --backfill ~/.claude/projects/<project>/<session>.jsonl
```

`MEASUREMENTS.md` has what that produced, and the blind discrimination test,
which is the instrument that answers whether a named voice arrived.
