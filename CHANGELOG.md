# Changelog

## v0.4.0 (2026-08-03) — the card arrives where the model reads it

- **The card is delivered as a Claude Code output style, which is where it
  should have been all along: `cope-gate --output-style`.** Every piece of
  evidence in this file below this entry came from `tools/generate_readme.py`,
  which sends the card as the entire system prompt — a condition cope was never
  deployed in. Deployed, it went through a `SessionStart` hook, one turn-zero
  message that the conversation buries. A three-arm test measured what that
  costs: 27 samples, 22 usable after five empty responses, counting bolded-label
  paragraphs. Card as system prompt, 9.4 per reply. Card plus a global
  CLAUDE.md, 6.3. CLAUDE.md alone, 1.0. So 63–86% of the effect survives
  alongside other instructions and the failure was position, not interference.
  Claude Code's own documentation says an output style's instructions are added
  to the system prompt and that the harness re-reminds the model of them during
  the conversation; a `SessionStart` message gets neither. Selection is
  `/config` → Output style, or `"outputStyle"` in settings — `/output-style` was
  deprecated in 2.1.73 and removed in 2.1.91.
- **`cope-gate --setup` does the whole install, so it is two commands rather
  than a settings file edited by hand.** It emits the output style, wires `Stop`
  and `UserPromptSubmit` into `~/.claude/settings.json` with absolute paths, and
  prints the one step left. It backs the file up first, appends rather than
  replaces so another tool's hook on the same event survives, is idempotent, and
  refuses to touch a settings file that does not parse — Claude Code ignores an
  unparseable settings file wholesale, so a half-repair is worse than nothing.
  `--dry-run` prints what it would change and writes nothing.
- **`SessionStart --inject` stands down on its own when a cope output style is
  active**, because the two together put the card in twice. The first live test
  of the style was invalid for exactly that reason.
- **The hooks are live the moment they are written.** Claude Code re-reads hook
  configuration per prompt on 2.1.220: a session running for six hours, never
  cleared, ran a hook added three minutes earlier. Only the output style is read
  once at session start. `--setup` says both, because "restart your sessions"
  is advice that costs a day and buys nothing.
- **Two phrase rules, and the short list is the finding.** Across 10,652
  assistant replies in the ten most recent transcripts, `load-bearing` runs 25.6
  per 1k and `worth noting` 6.5. The famous ones do not: "you're absolutely
  right" appeared four times, "great question" twice, "happy to help" and "let
  me break this down" never. Banning those would lengthen the card and change
  nothing, which is the conclusion basanite reached when it cut nine folklore
  entries. Single words stayed out even where they measure high — `substrate`
  and `texture` have real uses, and this block bans rather than reports.
- **`card/demo/claude_maximal.effigy` scores clean while instructing every
  measured tic**, which is the strongest argument in this repo for reading the
  demo page's gate column as something other than a verdict. It produces 36
  bolded-label paragraphs, 55 em dashes, 13 emoji headings and 15 horizontal
  rules, and every one of those falls in a documented exclusion.
- **The caveman card broke its own concreteness rule in a passing example.** Its
  kernel says find the physical thing an abstraction came from — contention
  becomes fight, latency becomes wait — and its HEDGE TEST then passed
  "Serialization slow." Now "Packing slow", which is shorter as well as
  physical. That is the test a general simple-word rule fails: plain vocabulary
  runs longer than the term it replaces, and this card's argument rests on a
  96-to-40 token drop.
- **The README taught a three-hook install that `--setup` does not wire**, and
  the count lived in the heading, which is how it went stale without anyone
  touching it.
- **A card can write a structure rule, not only decline one: `@shape <id>:
  <selector> <predicate> — <why>`.** The set of things cope can check is no
  longer whatever the binary happens to implement. `card/demo/handoff.effigy`
  asserts `readable_cold` — `last paragraph words <= 60` — for the claim its own
  peak makes and no built-in rule tests: that a reader re-entering cold can read
  the closing block alone and know what to do. The 60 is measured rather than
  chosen. Across 43,155 assistant replies in the six largest transcripts on this
  machine the closing block runs 33 words at the median, 56 at p90 and 65 at
  p95, so it fires on roughly the longest one in fourteen — the same order as
  the 10.1% `dangling_end` rate the card was built around.
  - Selectors are `first paragraph`, `last paragraph`, `every paragraph`, `some
    paragraph`, `reply`; predicates are `words`/`sentences` against `<=` or
    `>=`, plus `asks` and `does not ask`. It counts and it does not match,
    because matching is what POSTPROC is for.
  - Every selector reads `rawParagraphs`, not `paragraphs`. A card asserting its
    closing line is short is asserting it about a block the prose rules skip for
    being short, and running this against the filtered set would make the most
    useful claim in the grammar unexpressible.
  - The violation carries the card's own `why`. A rule a card wrote should
    explain itself in the card's terms.
  - Ids are validated against the built-ins and against the card's own POSTPROC
    rules, so a card rule cannot shadow something already named.
  - `@gate` on a rule the card defined in `@shape` is an error that names the
    actual fix — delete the `@shape` line — rather than the generic "not a
    built-in rule", which would send an author hunting for a typo in a correctly
    spelled id.
  - What it still cannot express: `clause_symmetry` and its kind. That stays
    compiled in, and a card wanting a check outside both this vocabulary and a
    regex has nowhere to put it. The README says so in the same section that
    used to say structure was unauthorable.

- **A card can decline a structure rule: `@gate <rule_id> off — <why>`.** This
  was the repo's named limitation, on the README's second screen for a month —
  a card whose VOICE asked for something a built-in rule catches got marked down
  for obeying itself, and `card/demo/lecturer.effigy` was the standing example.
  It now declines `clause_symmetry` and `dangling_end` in its header and scores
  2 against its own card where it scores 4 against the shipped one. Both numbers
  are real and they answer different questions.
  - The rule id must be one the gate has; a typo is reported at load rather than
    ignored, because a silently dropped exemption reads exactly like a working
    one from inside the card file.
  - The reason after the dash is part of the syntax. An exemption is a card
    marking its own homework, and one with no stated reason cannot be reviewed.
  - The detector still runs. Only the scoring card's own hits are dropped, so a
    backfill still reports what the rule would have caught, and one card's
    exemption cannot change what another card is measured against.
  - `Allow` is applied in three places, not one: `Check`, the end of
    `CheckLane` for the two loop-only rules, and around `cross_turn_repeat` in
    `checkAgainst`, which is produced outside the card and would otherwise have
    been the one rule an `@gate` line silently failed to turn off.
  - What is still one-sided: a card can turn a rule off and cannot add one. The
    set of checkable things is fixed.
- **effigy keeps unrecognised header keys instead of dropping them** (in
  `_parse_header`, plus an `extra: dict[str, list[str]]` on `CharacterAST`).
  That is the whole upstream change, and it is deliberately ignorant of what
  cope means by `@gate` — a character notation should not learn the vocabulary
  of a prose linter. Values are lists so a repeated key accumulates. Without it
  the Python path would have silently dropped every exemption between the
  `.effigy` source and the JSON `make rules` commits, while `make check-rules`
  went on passing. `tools/card2json.py` now emits `gate` when the header carries
  one, and omits the key otherwise so an unexempted card's JSON is unchanged.
- **`scan.RuleIDs` is the authority on what a card may name**, with
  `TestRuleIDsMatchTheDocumentedSet` holding it against the descriptions in
  `shapeRules`. A rule in only one list is either undocumented or unexemptable.

- **`card/demo/claude_maximal.effigy`, which instructs every measured tic and
  scores clean.** It asks for a bolded label on every paragraph, an em dash in
  every paragraph, an emoji on every heading, a horizontal rule every few
  hundred words, a "Key takeaways" list per section and a closing line under
  twelve words. The render obeys: 36 bold-label paragraphs against the shipped
  README's 0, 55 em dashes against 23, 13 emoji headings against 0, 15
  horizontal rules against 1. The gate finds nothing. Every habit lands in an
  exclusion — `labelled_opening` skips the bolded form, `paragraphs()` skips
  list blocks and anything under twelve words, markup is on neither axis, and
  the lexical tics went to basanite in July. Each exclusion was a defensible
  call with a date attached; together they are a gate that cannot see this model
  at its most recognisable, which is a better argument for reading the demo
  column as something other than a verdict than any sentence on that page.
## v0.3.0 (2026-08-03) — the docs said things the code did not

- **The card is built from measurement, not from folklore.** First draft came
  from claudisms.ai, a public CC0 banlist of AI-writing tells. Checking it
  against basanite's measured tic list — audited over roughly 4M tokens
  of real output — found the two lists disagree: "at the end of the day",
  "you're absolutely correct", "liminal" and six more were cut on 2026-08-02 for
  occurring nowhere outside writeups about how Claude supposedly talks, and
  several of those are on the public list. The card instructs only what
  survived. The note left in that file says it better: the folklore about what
  Claude says and what Claude says are different lists.
- **The README sends the reader to `demo/` in its second paragraph.** The link
  had been sitting in section five, below the install. A reader who opens two
  demo files against each other understands the tool in a glance, and a reader
  who does not is taking prose about registers on trust; that is not something
  to bury under the hooks configuration.

- **Two rules of thirteen fire on a document, and now the demo page says which
  two.** Across all six rebuilt files — 98,000 characters written six different
  ways — only `labelled_opening` and `clause_symmetry` found anything. The page
  had estimated "roughly four rules" from which ones were structurally capable
  of firing; the measured answer is smaller and does not need the estimate.
- **Why `stock_assistant` scores 2 was explained twice and wrong twice; the
  third answer comes from the gate rather than from counting by hand.** The page
  first said it wrote more prose than the other five, then said it wrote less.
  Both numbers came from a prose-share rule invented for the occasion. Measured
  through `paragraphs()`, the function every shape rule actually calls, it has
  31 visible paragraphs — the most of any file here — and still draws zero
  `labelled_opening` hits, because 19 of them open on a bold label and
  `labelled_opening` skips the bolded form. That is the `bold_label` rule the
  card dropped on 2026-07-30 after blind readers named bold and bullets among
  the things deciding which reply they preferred. The card written to fail is
  built almost entirely out of the one opener the card gave up policing.
- **Two documented references to `bold_label` outlived the rule by four days.**
  `shapeRules` in the doc generator described `labelled_opening` as deferring
  the bolded form to "the `bold_label` regex", and the related-projects note
  offered `bold_label` as an example of what stays in POSTPROC. The card has
  carried no such rule since 2026-07-30, so the generated README told a reader
  the bolded label was covered when nothing covers it. A rule description that
  names a card rule outlives the card rule; these now say the form is
  deliberately unpoliced and that a card can add it back.
- **The README claimed the shipped card overflows the NEVER budget and reports
  it.** The budget is 10 and the card holds 12, and the generated text put those
  side by side and concluded a warning fires at load. None does, and none should:
  the budget is charged against each injection separately, because SessionStart
  prints the always-on rules while the refresher prints the evidence-gated ones
  and nothing renders their union. The prompt now says so and forbids stating a
  card total, which was a count with no fact behind it that rots on the next
  card edit.
- **`transcript.userPrompt` had no callers and CI had been red on it for four
  commits.** It was the variant without the parent id, orphaned when the lane
  started chaining a prompt back to the wakeup that delivered it. Its doc
  comment explains why tool results are not turn boundaries, which is the
  non-obvious part of assembling turns at all, so that moved onto
  `userPromptRec` rather than going with the function.
- **The README shows the hooks block instead of describing it.** It named the
  three events and what each one did, in prose, and left a reader who had never
  written a Claude Code hook to derive the matcher, the `type` field and Stop's
  `async` from that. The generator now carries the working configuration as a
  literal and the prompt orders it reproduced verbatim, so the one step where a
  new install fails silently is a block to copy. The note under it says the bare
  commands assume `go install`'s target is on PATH in the hook's environment,
  which is the other way it fails silently.
- **`card/demo/handoff.effigy` is explained where a reader meets it.** It is a
  hypothesis rather than a voice — the shipped card's handoff rules with the
  prose half removed, to be run through `make pairs` against the full card — and
  `make demo` correctly does not render it. But `make cards` installs it, so it
  appeared in a listing of registers to write in with nothing saying it was not
  one.
- **`demo/` is rebuilt from every card and linked from the README**, so the
  shortest way to see what a card does is to read the same document in six
  registers rather than to write one. `make demo` now hands each `.effigy`
  straight to the gate instead of compiling it through `tools/card2json.py`
  first, which had left the one target whose purpose is showing a card off
  needing an effigy checkout to run.
- **The demo page stops presenting its gate column as a verdict.** On a
  document, four of the structure rules cannot fire at all — they read the
  ending of a conversational turn for whether the reader has anything to
  answer, and the section list forbids manufacturing one — so the column
  separates cards on voicing and is silent on the other axis. What it does
  separate, it separates twice for opposite reasons: `caveman` scores highest
  almost entirely on `labelled_opening`, which is its dropped copulas meeting a
  detector for verbless openings, and `lecturer` scores on `clause_symmetry`,
  which is the climb its VOICE block asks for in as many words. The second is
  the authorability gap in a file a reader can open. The card written to trip
  everything, `stock_assistant`, scores 1.
- **The vocabulary-compression claim reversed at document scale.** `precise`
  beats `caveman` on one sentence, 24 tokens against 40, and `README.caveman.md`
  is now the shortest file in `demo/`. Grammar compression applies to every
  sentence; vocabulary compression only pays where there is a term of art, and
  most of a README is not that.
- **Structure and voicing are named as two things, and the README is written
  around the split.** Voicing is what the sentences sound like and lives
  entirely in the card, so swapping the card swaps every word of it. Structure
  is where the decision sits and how the reply ends, and it is compiled into
  `internal/scan`, so it is the same whichever card is loaded. Every rule now
  declares which axis it is on, `--author-docs` emits the split as
  `facts.axes`, and the rules section groups by it rather than by where the
  rule happens to be implemented. The split was already latent in `ruleClass`
  in `abtest.go` and in the A/B finding that four fifths of the hit rate tracks
  what a reply was for; naming it is what makes the null legible, because the
  instrument that measured the refresher is four fifths structure and the card
  is mostly voicing.
- **The README says which half is authorable, because it is only one.** A card
  can state a structural commitment in its own prose — `lecturer.effigy` holds
  its conclusion back on purpose — and nothing in the gate can read it, so a
  card whose voice asks for something a built-in rule bans is marked down for
  doing what it was told. MEASUREMENTS.md carries the run behind that and the
  reasons its numbers do not carry further: the direction is there at p=0.07
  and the sharper-looking structure number is mostly reply length.
- **The generated README had been missing live surface.** The rules section
  listed six of the thirteen built-in rules, the flag table predated six flags,
  and the loop lane was absent entirely despite inverting half the ending rules
  for unattended turns. All of it came from a hand-maintained table in
  `authordocs.go` that nothing checked against `internal/scan`.
- **`make readme` says what to do when the anthropic package is missing.** A
  distro python has been externally managed since PEP 668, so the target's old
  advice — `pip install anthropic` — is refused by pip itself, and the target
  died with that line the first time the system python moved a version. It now
  prefers a venv at `~/.venvs/cope` and prints the two commands that make one.

- **The gate reads `.effigy` directly, so a card is usable as written.**
  `internal/effigy` translates the notation into the JSON `scan.ParseCard`
  already consumes — a translator rather than a second reader of the format,
  with every check still happening in one place. Until now a second voice meant
  running `tools/card2json.py` against a checkout of the effigy repo, which
  makes "write your own card" a contributor workflow rather than a user one,
  and a voicing framework whose second voice needs a build step has one voice.
  `tools/card2json.py` stays the build-time path for the embedded card and CI
  still regenerates through it, so effigy remains the authority on the notation.
- **The two parsers are held together by an oracle test, not by care.**
  `TestAgreesWithCard2json` runs every card in the repo through both and
  compares the results, because the fixture that would have caught a bug is the
  card somebody actually wrote. Compared as values rather than bytes: the Python
  side writes indented JSON with non-ASCII escaped, and neither is a fact about
  the card.
- **`--card`, `--cards` and `$COPE_CARD`**, so switching voice is a word rather
  than an edit to a hook. Names resolve in `$XDG_CONFIG_HOME/cope/cards`; a card
  is configuration, which is written by hand and meant to be kept, where the
  violations log and the refresher markers are disposable machine state.
  `make cards` installs the demos. A name that resolves to nothing is an error
  and nothing is injected — falling back to the built-in card would leave a
  session writing in the shipped voice while its config named another one, which
  reads as the card not working rather than as a typo.
- **`card/demo/lecturer.effigy`, the first rival a discrimination run can learn
  anything from.** The other demos are identifiable before a reader forms any
  view about voice — laconic is four words long, caveman has no articles,
  precise is a wall of terminology — so a high score against them measures the
  presence of a gimmick. Lecturer is ordinary English at ordinary length and
  differs from the shipped card on one axis: whether the prose may call
  attention to itself. Its WRONG blocks therefore read backwards against the
  rest of the repo, and one of them is quoted verbatim from `precise.effigy`'s
  MES, where it is the right answer. It scored 95% where laconic scored 90%,
  which retires the length confound the previous run could not rule out.
- **`card/demo/caveman.effigy`**, the missing half of the argument
  `precise.effigy` makes. Both cut an answer to a quarter; caveman removes the
  words holding a sentence together and stops at 40 tokens, precise replaces a
  paraphrase with its term of art and reaches 24 while still reading aloud. It
  is the demo's smoke test — a failure is visible from across the room — and it
  is useless as a rival for exactly that reason.
- **A discrimination harness, because preference cannot answer "sound like
  this instead".** The blind pairs measure which of two replies reads better,
  and both arms of every run so far were Claude writing well — an instrument
  that structurally cannot detect a voice being swapped. `-discriminate` shows
  the reader one card's VOICE block and asks which reply carries it. Chance is
  50% and the arithmetic is the same binomial, but a null means something much
  stronger: two replies can be equally good in a hundred ways, and two that
  cannot be told apart under a stated target had no voicing happen to them.
- **`cope-gate --describe`** prints a card as a target to recognise rather than
  an instruction to follow: the aim and the register, none of the machinery. The
  card describing itself is the only honest source — a description written for
  the test is one more artefact that can be wrong on its own, and a failed match
  would then leave two readings standing.
- **`-judge` answers the same pairs with the model**, so revising a card costs a
  run rather than an evening. It is not trusted alone: it answers what the
  reader answers and the page reports the agreement, because until that number
  is looked at its rate is a fact about a model, not about a voice.
- **The answer's side is balanced, which is not the same as the card's side.**
  `balanceSides` puts the card left in half the pairs, which is the answer on a
  preference run. Here the target is the rival in half the pairs, so the answer
  landed on the right in five of six on the first live run — free marks for any
  reader with a side habit, in the one direction blinding cannot cover.
- **A truncated judge reply is missing data, not an answer of "unsure".** Opus 5
  emits a thinking block ahead of its text whether or not it was asked for, so a
  ceiling sized for the one-word answer left no room for the answer: six pairs
  came back "unsure", indistinguishable at a glance from a judge that could not
  tell two voices apart. Missing answers are now errors and are dropped.
- `make pairs` could not have built its own binary since `replay/` got its own
  `go.mod` — `go build ./replay` from the root module cannot see it. Both
  targets now build from inside the module.
- **A loop lane, because the reply stops being a message and becomes a record.**
  Addy Osmani's June 2026 "Loop Engineering" is the frame — you stop prompting
  the agent and design the system that prompts it — and the consequence for a
  voice gate is that `/loop` and `/goal` write for a reader who is asleep.
  `transcript.Turn` now carries a lane, read off the prompt that opened it, and
  `Card.CheckLane` drops the four rules whose subject is a reader who can answer
  (`ask_not_last`, `forked_end`, `dangling_end`, `buried_decision`) and adds two
  that assume one who cannot: `loop_ask` and `unverified_done`.
- **`loop_ask` earned a card rule; `unverified_done` did not.** Over 870
  loop-lane turns, `loop_ask` fired 55 times — "Want me to shorten it, change
  the tone, or adjust anything?" written into a log at an hour nobody was
  reading it. Each of those endings is well-formed by the interactive rules,
  which is what a lane-blind gate cannot see. `unverified_done` fired once, so
  it keeps its detector and loses its card entry; the card already carries an
  always-on rule against calling written work verified.
- **Loop replies are shorter, not cleaner.** Across 34,543 turns from 105
  sessions: interactive hits on 62.4% of turns at 2,757 chars mean, loop on
  21.1% at 427. Per thousand characters that is 0.23 against 0.50 — the loop
  lane is twice as dense. One session's loop replies averaged 37 characters,
  which is not a report.
- **A third gate namespace.** `at_*` is where an injection sits, `rule_*` is
  what the writer has been doing, and `mode_*` is what kind of session this is.
  SessionStart opens only `at_*`: the other two are not known that early.
  `mode_loop` replaces the position fact rather than joining it, because every
  `at_prompt` item is written for a reader about to answer — the CONTINUE TEST,
  and a pair whose better half ends in a question. Carrying those into a lane
  that must not end by asking would put the card on both sides of one question.
- `--render-lane loop` prints what the unattended lane injects, and `--backfill`
  splits its rates by lane with the mean reply length beside each, because reply
  genre is the standing confound and a lane of short replies reads as a lane of
  clean ones.
- **`transcript` reads scheduled wakeups.** Every loop iteration after the first
  arrives as a `system` record with subtype `scheduled_task_fire`, with the
  prompt it delivers chained underneath by `parentUuid`. Matching on prose
  instead returned an empty loop lane on a session with 423 of them.
- **`cross_turn_repeat`, the first rule that reads the session rather than the
  reply.** The most-cited complaint in the July 2026 roundup is repetition —
  "the same turns of phrase over and over" — and it is the one thing a
  per-reply linter cannot see, because a construction used for the fourth time
  looks exactly like one used for the first. `internal/scan/repeat.go` splits a
  reply into five-word shingles, keeps the ones built mostly of common words
  (four of five, so one content word rides along inside a frame and three make
  it subject matter rather than shape), and fingerprints each to 32 bits.
  `state.Turn` carries the sums, so the window says which constructions this
  session has already used. A phrase already in two earlier turns fires the
  rule, once per reply, naming the busiest. On 27,610 assistant turns from 80
  other-project sessions it fires on 4.96%, median 1.8% per session, and the
  worst session — 13.3% over 6,927 turns — printed nothing but the "let me look
  at the" preamble family. Roughly three in four sampled matches are
  constructions; the rest is subject matter that got through the filter.
- **The transcript reader was scoring Claude Code's own error banners.** "API
  Error: 500 … usually temporary — try again in a moment" is stored as an
  assistant record, once per retry, identically — and it came second in the
  first repetition backfill. `internal/transcript` now skips records marked
  `isApiErrorMessage` or carrying `model: "<synthetic>"`. Every rate in
  MEASUREMENTS.md before today included them; the correction is small (27,962
  turns to 27,610) because banners are short, but it took a rule that reads
  across turns to find it.
- Fingerprints rather than text, because `state.Session` is a file on disk and
  its whole claim is that rule names and character counts do not reconstruct a
  conversation. A 32-bit sum of five words keeps that: enough to notice a
  recurrence, not enough to read one back. The quoted phrase in a violation
  comes from the reply the gate is already holding.
- `--backfill` now carries a rolling window, so the cross-turn rule is scored
  the way it runs. Without it, the one rule that most needs a corpus before
  shipping is the one rule a backfill could not report on.
- `card2json.py` stopped warning about a card that renders fine. It counted the
  file's NEVER total against a cap that has been per-injection since the budget
  fix below, so every build warned; it now counts the two injections apart.
- **The NEVER budget was capping the rule library instead of the prompt.** The
  ten-rule cap is an attention limit on same-weight constraints shown at once,
  and it was being charged against the card file's total before gating resolved
  — so an evidence-gated rule held back for a habit the writer had never
  exhibited still consumed a slot, and the card could never hold more than ten
  rules in total. The cap now applies to what each injection prints. SessionStart
  renders the always-on rules, the refresher renders the evidence-gated ones that
  fired, and neither exceeds ten. The card file is free to grow past it.
- **SessionStart no longer renders evidence-gated rules.** `rule_*` gates say the
  writer has been tripping something; at SessionStart nothing has been measured,
  so those items now wait for their detector to fire and arrive through the
  refresher. `at_*` gates are untouched — they mark an injection point rather
  than evidence, and closing them would have dropped the CONTINUE TEST from the
  only injection that states the card in full. The card measured at 79% rendered
  all twelve items; this one renders seven, so that result stands against the
  older render.
- **Three rules from outside this repo's transcripts, one of which earned its
  place.** Zvi Mowshowitz's July 2026 Opus 5 reaction roundup collects named
  practitioners on what makes the model unpleasant to talk to; three moves recur
  across otherwise unrelated complaints. Detectors for all three shipped; only
  `apology` reaches the card. On a 52.6M-character backfill across 20,911
  assistant turns from other projects: `apology` 25 hits, all genuine contrition
  ("Sorry for the wasted commit cycle", "Apologies for the cumulative miss-rate
  this session"). `self_postmortem` and `announced_length` found nothing real —
  their only matches were a security caveat, an offer to talk someone through
  operator steps, a one-sentence correction, and three uses of "dense" describing
  a room in a game. Each became a fixture. Both keep their detector and neither
  gets a card rule, because a gated rule whose detector never fires never renders,
  which is the failure 584c898 already paid for.
- The corpus is the caveat on those two nulls. It is Claude Code sessions under a
  global CLAUDE.md that already bans preamble and rumination, while the roundup's
  complaints come from chat. A null here says the move does not appear in this
  writer's setup, not that it does not appear.
- `deference` is a fourth rule class in the A/B report, holding the three
  field-sourced rules apart from the ones mined from these transcripts so the two
  groups can be seen to move independently.

- **Two NEVER rules were worth more than nine thousand characters of voice
  guidance.** A 5,125-character card carrying only the handoff rules beat the
  shipped 14,011-character card 11 of 12 head to head. Adding the two rules it
  had that the shipped card lacked — never end on more than one actionable
  thing, never paste a document into the reply — moved that to 5 of 11, a tie.
  The rival side never changed between the runs, so the two rules are what moved
  it, at p=0.069 on the difference. Both are in the card now and the NEVER block
  holds 9 of its 10. The voice half turns out to be inert rather than costly:
  with handoff rules on both sides, 14,334 characters ties 5,125.
- `-rival` runs two cards head to head instead of stacking one on a baseline,
  which is the comparison anyone iterating a card actually wants — is the new
  one better than the old one. `card/demo/handoff.effigy` is the candidate that
  won, kept as a demo rather than adopted.
- **The README claimed more than the measurements supported.** It said twice
  that nothing here measures whether the prose got better, which stopped being
  true when the blind pairs ran. It presented the refresher null as a verdict on
  the refresher, when the metric that produced it is four fifths reply genre by
  its own account. It described three hooks as equally load-bearing. And four
  counts had rotted against the card they were generated from. All corrected,
  and the pairs results are in it now.
- **The card beats a global CLAUDE.md at 79% when a reader judges it blind.**
  Two runs of the edited card against the reader's Response style section,
  16/19 and 14/19, pooled 30/38 at p=0.0005. The observational comparison below
  could not detect this and neither can the hit-rate proxy, which is the point:
  the metric counts moves, and the pairs ask which reply he would rather have
  gotten. Also the mechanism he named is not the one operating — the side he
  picks has a third the bullets and a fifth the bold of the side he rejects, in
  every run. `MEASUREMENTS.md` carries all three runs and the correction to an
  earlier reading that took case mix for a trend.
- **`bold_label` is out of the card**, rule and POSTPROC pattern both, along
  with the display transform that repaired it and the refresher item gated on
  it. Fifty-two blind pairs put bold among the three things the reader says
  decide a reply for him, and the rule had been banning it. The MessageDisplay
  transform table is empty now and says so.
- The injection preamble told the model the card constrained prose only and did
  not change what it did. That fenced the card out of exactly the handoff
  behaviour it was supposed to shape; `internal/scan/render.go` now says brevity
  belongs to the prose and not to the work.
- Three QUIRKS and a NEVER, all four about handing work back rather than about
  sentences: say how good the result has to be, name what your own proposal
  would fail at while proposing it, ask the invalidating question first, and
  never hand over a procedure you could have run yourself.
- **`make pairs` runs the blind test against any card.** `BASELINE` is the
  prose guidance you already have, `CARD` points at any `.effigy` file, and the
  card is tested on top of the baseline so the run measures what the card adds.
  The reader's own CLAUDE.md section is no longer compiled into the harness, the
  generator retries on 429/5xx instead of losing pairs to them, and
  `replay/README.md` documents what 20 pairs can and cannot resolve.
- **The card came back null against a global CLAUDE.md too.** 17 transcripts
  carrying `<voice-card>` against 24 without, one binary so the rule set is
  constant, measured on the rhetoric rules alone: 0.186 against 0.190 hits per
  1k, a rate ratio of 0.979 with a 95% interval of 0.773–1.240. Formatting runs
  the other way and clears 1, which is the genre confound rather than the card —
  the card transcripts are recent and implementation-heavy. Observational, not
  randomised. It shows the rhetoric metric cannot see a difference, which is not
  the same as the card doing nothing, and separating those needs an instrument
  that reads meaning rather than surface form.
- **The `inject`/`hold` experiment came back null.** 238 turns over 7 sessions:
  0.843 hits per 1k characters injected against 0.820 held back, a rate ratio of
  1.028 with a 95% interval of 0.854–1.237. An offline replay of the same
  prompts under each arm agreed. The interval is narrow enough to rule out an
  effect larger than about 20%. The drift split, which read +0.389 against
  −0.180 at 105 turns, was −0.089 against −0.030 by 238 — noise. `MEASUREMENTS.md`
  carries both runs and the interim that failed.
- `--ab-report` splits each turn by distance from the boundary of the window it
  was written in, and groups the rules by what they detect. This is where the
  null got its explanation. The rate climbs 2.4× with distance under `hold`,
  which receives no reminder, and the climb is almost entirely `formatting` —
  bold labels, verbless openings, even blocks, which is what a status report
  looks like. Reporting clusters in the working stretches of a session, so
  distance from the boundary stands in for how deep into an implementation run
  the turn sat. Four fifths of the metric is sensitive to that; the two rules
  that read as voice rather than layout are flat across distance and across arms.
  Reply length, turn tempo and the SessionStart card were each tested as the
  cause first and each failed.
- `turns.jsonl` records which build scored each turn, and `--ab-report` names
  every build in its header. The hook names an absolute path and `go install`
  replaces the file underneath it, so a run can span several binaries with
  nothing in the data to say the rules changed partway. Phase 1 did exactly that
  and comparing file mtimes was the only way it surfaced.
- The recommended live arms are now `positive,hold`. Swapping the third arm in
  beats adding it: three arms divide the same windows three ways and cost half
  the power of a question already too underpowered to resolve. What is given up
  is the direct `positive`/`inject` comparison, which was never the one under
  test — `hold` is the reference either way.
- A third arm, `positive`, off by default and reached with
  `--ab-arms inject,hold,positive`. It renders the same slice with every
  anti-pattern removed: the better half of each wrong/better pair and the
  passing examples, and none of the prohibitions, none of the worse lines, and
  no rule names or counts. The first arm cannot tell "the reminder did not
  land" from "the reminder put the move in reach by naming it", and a reminder
  against opening on a bold label has to write "bold label" down. Selection is
  unchanged — which examples appear still comes from what fired.
  Opt-in because a third arm divides the same windows three ways. The two-arm
  default is unchanged so turns already collected stay comparable.
- `turns.jsonl` records the length and the named rules of the reminder that was
  live when each turn was written. The positive render is about a third the
  length of the paired one, so the two arms differ in size as well as in
  content; without the length on disk a gap between them could not be
  attributed to either.
- `--ab-report` discovers its arms from the log rather than assuming two, and
  measures every arm against `hold` as the reference. It prints each arm's mean
  injection length beside its rate for the same reason.
- `--refresh-every` drops from 4h to 30m. Replaying the old clock against real
  session timelines, the refresher fired on 0.9% to 5.7% of prompts: the window
  is wall-clock and every compaction restamps it, so a long working stretch
  never reached one. The card's 11,923 characters arrive 8 to 48 times a
  session against that, which made cope about 96% a large paste by volume.
- `--ab` runs the experiment the project turns on — whether the sliced
  mid-session injection does anything a static file does not. It holds back every other
  refresher window and records which arm each turn was written under, to
  `turns.jsonl` — rule names and counts, no prose. `--ab-report` compares the
  arms as a Poisson rate ratio with a 95% interval, and says when the interval
  contains 1 rather than picking the larger number.
  Windows alternate rather than falling to a coin, and the crossover sits
  inside a session rather than across sessions, which holds project and task
  fixed instead of leaving them as variance across an N of a few dozen.
- The gate scores a turn instead of a record. One reply is written to the
  transcript as several assistant records — prose, a tool call, more prose —
  and `LastAssistantText` returned the last of them, so the gate judged a
  mid-reply aside. Scoring this repo's own session both ways: 26 labelled
  openings against 44, 19 bold labels against 28, and 15 paragraph-uniformity
  hits against 6, the last of those being the false-positive direction, since
  a trailing fragment of four short paragraphs looks even in a way the reply
  does not. `--backfill` read the same file as 204 turns where it holds 63.
- Scoring now spans two hooks, because a reply's closing block reaches the
  transcript after the Stop hook has already run. Across six live sessions on
  2026-07-28 every one of the 17 turns in state had been scored against text
  that ended early, and four against text from an earlier turn entirely. The
  Stop hook scores what is there and the next UserPromptSubmit rescores the
  finished turn, reporting only what the first pass could not have seen. The
  ending is where `dangling_end`, `buried_decision` and the last ask live, so
  none of the three had ever been evaluated against a real ending in a live
  session.
- `state.Turn` carries the id of the prompt that opened it, which is what lets
  the second pass correct the first rather than count a second turn.
- `--check` no longer prints its violations to both stdout and stderr.

- `internal/state` keeps a rolling record per session: turns, characters, and
  which rules fired, over a 20-turn window. The Stop hook writes it; the next
  UserPromptSubmit reads it. Hooks are separate processes, so disk is the only
  channel between turns.
- The mid-session injection is now chosen from that record. `@when` gates are
  back on NEVER rules, tests and worse/better pairs — not on the axis that was
  tried and dropped in June (register by speech act, which turned out to be one
  register), but on where the injection happens (`at_prompt`) and on what the
  writer has actually been tripping this session (`rule_bold_label`). A cold
  session still gets the standing reminder; three bolded labels later, the
  injection names the count and carries the rule that addresses it.
  This is the one thing cope renders that a pasted CLAUDE.md cannot, because
  the text is selected from measured output rather than fixed in advance.
- `Render()` is unchanged at SessionStart: nil facts keep every gated item, so
  the full card is never narrowed by a gate. A test fails if it is.
- Quirks stay ungated — effigy parses them as plain strings with no `@when`
  support, verified against the parser rather than assumed.

- `verdict_handoff` left the card for basanite's `known-tics.txt`. It matched a
  verbatim phrase, which is basanite's axis, and the two tools should not both
  be watching the same string. POSTPROC is down to `flip` and `bold_label`,
  both syntactic: a template with a variable slot, and markup. The README now
  says why the list is short instead of leaving it looking unfinished.
- The card gained an `MES` block — six whole replies in the wanted voice.
  effigy has had the block all along and cope used none of it, while
  Anthropic's Opus 5 guidance puts positive examples of the wanted style above
  instructions about what not to do. They render before the NEVER list, and a
  test fails if that order flips.
- Worse/better pairs carry a `WHY`. effigy parses one
  (`notation.py: WrongExampleAST.why`) and `card2json.py` was dropping it, so
  every pair taught a swap without the move behind it.
- The injected card is 11,923 characters, up from 9,022.

- `NOTES.md` became `MEASUREMENTS.md` and lost half its length. Four sections
  were history the changelog already carries or ideas never run: why effigy
  notation, why the `@when` gates went, what basanite cannot see, and a "still
  open" list overlapping the README's known limits. The effigy and basanite
  answers are facts in the docs prompt now, so the README carries them in a
  sentence each. What is left is the evidence: the baseline rate, both hit
  tables, the judge-panel result, and an honest account of what has not been
  measured.

- Markdown tables are no longer read as prose. A grid is balanced by
  construction — every row repeats the column's vocabulary at near-equal
  length — so `clause_symmetry` turned the README's own six-row layout table
  into two violations. Found by pointing `--check` at this repo's docs.
- `README.md` is now generated. `tools/generate_readme.py` runs the card
  through a model and the draft back through the gate; the committed file is
  the first one the gate passed. Shape-rule descriptions became facts after a
  draft inferred them from the rule names and got `labelled_opening` wrong,
  claiming it catches list items when `paragraphs()` skips list blocks
  outright.
- `--check FILE` scores prose against the card, the same pass a reply gets at
  Stop. `--backfill` only understood transcript JSONL, so the repo's own
  documentation had no way through the gate it describes. `-` reads stdin,
  `make check` runs it over the README.
- `--author-docs` prints a prompt for writing this repo's prose: the rendered
  card, the facts introspected from the card and the flag set, and the section
  list. It states no prose of its own — effigy's rule, from the docstring of
  its `generate_readme.py`. The flag table cannot drift from the flags because
  it is read out of the flag set at run time.
- `tools/generate_readme.py` runs the loop: prompt in, draft out, `--check`
  over the draft, violations back as a revision turn until the gate is quiet.
  The system block is marked `cache_control`, so revision rounds pay about a
  tenth on the card and facts they re-send. Needs `ANTHROPIC_API_KEY` in a
  gitignored `.env`.
- Flag registration moved out of `main` into `registerFlags`, so the docs
  prompt and its test describe the same flag set. A test asking the binary what
  it accepts was getting the testing package's flags instead.
- `--log` resolves an absolute default under `$HOME`, and the facts are read off
  whichever machine writes the docs, so flag defaults are rewritten back to
  `$HOME` before they reach the prompt. That path was headed for a published
  README verbatim.
- `.env` added to `.gitignore`.
- `tools/effigy_bootstrap.py` holds the effigy-path contract that
  `card2json.py` and `run_postproc.py` each carried a copy of. The copies had
  already drifted once — the sibling-checkout default landed in one before the
  other. adit flagged the duplicate `EFFIGY_PATH` definition, calque flagged
  the shared seam.

## v0.2.0 (2026-07-28) — works on a machine that isn't mine

A public-readiness pass. The two items at the top are the ones that made the
gate a no-op anywhere but its author's laptop.

### Fixed
- The card is compiled into the binary with `go:embed`. `--rules` defaulted to
  `card/rules.json` under a checkout path guessed from `$HOME`, so a clone
  anywhere else printed `no check run` on every turn and checked nothing.
  `go install …/cmd/cope-gate@latest` now works with no clone at all.
- `make install` no longer depends on the `rules` target, so a build no longer
  requires a Python effigy checkout. `make rules` and `make check-rules` are for
  editing the card; CI runs the latter against effigy v0.7.0, which is what
  makes the no-drift claim in the README enforceable rather than aspirational.
- `LastAssistantText` reads backward from the end of the transcript. Reading
  forward took 28.3s on a 2.6 GB session against the README's recommended 10s
  Stop-hook timeout, so the gate was being killed on exactly the long-running
  conversations it exists for. The same file now takes under a millisecond.
- `ask_not_last` sees asks below the 12-word prose floor. "Build both?" scored
  zero while the identical ask padded to fourteen words scored one — asks are
  short by nature, and the floor was deleting the evidence. Over 61,612 turns
  the rule went from 925 hits to 1,331.
- `violations.jsonl` moved to `$XDG_STATE_HOME/cope` and is written `0600`. It
  holds verbatim excerpts of replies and used to land in the repo checkout at
  `0644`. An existing checkout's file is left where it is.
- `--backfill` scores what it managed to read instead of exiting on a truncated
  or malformed transcript.
- `tools/backfill-sweep.sh` no longer defaults `SKIP` to a session id from one
  machine, which silently excluded an arbitrary transcript everywhere else.

### Changed
- `significance_flag` cut: 1 hit in 20.1M characters. `verdict_handoff` stays at
  38. Measurements in `MEASUREMENTS.md`.
- effigy's `@when` condition gates were removed from the card in June; the
  plumbing that carried them is now gone from `card2json.py`, `internal/scan`,
  and the rendered output. `never` in `rules.json` is a list of strings.
- A `strip` rule is reported at load and treated as `warn`. A hook that sees the
  reply after it is written cannot rewrite it, and the action was silently doing
  nothing.
- `HANDOFF.md` became `NOTES.md`, dated and rewritten as a lab notebook. It had
  stale counts, line-number citations into another repo, and a "next" list whose
  first three items had shipped.
- CI pins staticcheck to 2026.1, runs `go test -race`, and declares read-only
  permissions.

### Performance
- One redaction and two paragraph splits per turn, against seven and five when
  every rule derived its own; `ordinalHead` and `clauseSplit` hoisted out of the
  loops that recompiled them. A 20.1M-character backfill went from 72.7s to
  43.4s reporting the same violations.

### Added
- `internal/transcript` has tests: the backward reader is checked against a
  forward scan across five chunk sizes, with and without a trailing newline.

## v0.1.0 (2026-07-28) — first tag

### Added since the initial baseline
- Decision-surface shape rules, from measuring ~475 Opus 5 turns against
  ~1,800 Opus 4.x turns in the same projects: `dangling_end` (open problems
  named, no question, offer, or all-clear) and `buried_decision` (an open
  problem after the last ask). `ask_not_last` for the mid-reply ask.
- `--refresher` — UserPromptSubmit entry that injects the CONTINUE TEST
  (~374 chars) once the last card injection has aged past `-refresh-every`
  (4h default), on a marker-file clock in `$XDG_STATE_HOME/cope`. Silent on
  a session's first prompt; `--inject` restamps the clock.
- NEVER-rule budget ported from effigy (cap 10, CRITICAL first), with
  overflow reported at load rather than silently discarded. The card's list
  trimmed 11 → 7; the ask-placement rule became two QUIRKS, a WRONG pair,
  and the CONTINUE TEST — positive structure over prohibition, per
  Anthropic's Opus 5 prompting guidance.

### Added
- `card/claude_voice.effigy` — voice card in effigy notation. 10 NEVER rules,
  4 WRONG pairs, 3 TESTs, 4 POSTPROC rules. Validates against effigy 0.7.0.
- `cope-gate` — Stop hook binary. Reads the hook payload on stdin, pulls the
  last assistant turn from the transcript, reports violations, appends to
  `violations.jsonl`, exits 0. `--block` opts into exit 2 on `reject` rules.
- `--backfill` mode for scoring a whole session, and `tools/backfill-sweep.sh`
  for a cross-project baseline.
- Shape rules in `internal/scan`: `labelled_opening`, `clause_symmetry`,
  `paragraph_uniformity`. These do not reduce to regex.
- Redaction of code fences, inline code, block quotes, and double-quoted spans
  before matching.
- `--inject` renders the card as prompt text for a `SessionStart` hook, so the
  writer receives the same rules the gate enforces. The POSTPROC patterns are
  deliberately withheld from the render — showing the regexes invites writing
  around them. `tools/card2json.py` now emits the whole card, not just POSTPROC.

### Notes on what the measurements said
- Baseline over 116,154 assistant turns / 38.7M characters: 1.22 hits per 1k
  characters. Hits-per-character is the metric; turns-hit tracks turn length.
- `significance_flag` fired twice in 38.7M characters and `verdict_handoff` 58
  times. Both target phrasings rather than moves and should be replaced.
- `bold_label` needed two wrong guesses: the tic is `**Label.**` with the
  punctuation inside the bold, not `**Label**:`.
- `labelled_opening` first keyed on the opener being short, which flagged
  "Basanite measures words" — a real sentence. Rewritten to test for a finite
  verb instead.
- The `@when` conditional gates were removed. A blind check across three
  registers found one register serving all three; what varies is dictated by
  the speech act rather than by voice.
