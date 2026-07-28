# Changelog

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
  38. Measurements in `NOTES.md`.
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
