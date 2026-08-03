#!/usr/bin/env python3
"""Emit a card's POSTPROC rules as JSON for the Go gate.

The card is authored in effigy notation and parsed by effigy, so effigy stays
the single source of truth for the format. The gate consumes the JSON, which
keeps a second .effigy parser out of the Go side.

    python3 tools/card2json.py card/claude_voice.effigy card/rules.json
"""
import json
import sys

from effigy_bootstrap import require_effigy

require_effigy("card2json.py")

from effigy.parser import parse_file  # noqa: E402

if len(sys.argv) != 3:
    sys.exit(f"usage: {sys.argv[0]} <card.effigy> <out.json>")

card, out = sys.argv[1], sys.argv[2]
ast = parse_file(card)

rules = [
    {
        "id": r.rule_id or f"postproc_{i}",
        "action": r.action,
        # effigy compiles POSTPROC patterns with re.IGNORECASE; carry that
        # into the Go side explicitly rather than relying on a default.
        "pattern": r.pattern,
        "why": r.why,
    }
    for i, r in enumerate(ast.post_processors)
]

def when_of(obj):
    """Condition gate, normalized. Empty and '*' both mean always-on.

    The gates were dropped in June after a blind check found one register
    serving three speech acts. That tested the wrong axis. cope now gates on
    where the injection happens (at_start, at_prompt) and on what the writer
    has actually been getting wrong this session (rule_bold_label), neither of
    which is a property of the speech act. effigy's condition grammar only
    parses `fact:name` for a bare symbol, so those are the fact-set form.

    QUIRKS carry no gate — effigy parses them as plain strings — so the
    positive habits are always injected.
    """
    w = (getattr(obj, "when", "") or "").strip()
    return "" if w == "*" else w



payload = {
    "card_id": ast.char_id,
    "source": card,
    "rules": rules,
    # Everything below feeds `cope-gate --inject`, which renders the card into
    # context at SessionStart. The gate reads output; this is the other half.
    "voice": {
        "kernel": getattr(ast.voice, "kernel", "") if ast.voice else "",
        "peak": getattr(ast.voice, "peak", "") if ast.voice else "",
    },
    "theme": ast.theme,
    # @gate lines: the card's exemptions from built-in structure rules. effigy
    # keeps unrecognised header keys in .extra rather than dropping them
    # (v0.7.1+); the older parser silently discarded them, which would have made
    # a card's exemptions vanish between the .effigy source and this JSON with
    # nothing to show for it. Absent rather than empty when there are none, so
    # the committed rules.json for a card without exemptions is unchanged.
    **({"gate": list(ast.extra["@gate"])} if getattr(ast, "extra", {}).get("@gate") else {}),
    # @shape lines: the card's own structure rules. Same passthrough as @gate.
    **({"shape": list(ast.extra["@shape"])} if getattr(ast, "extra", {}).get("@shape") else {}),
    "traits": list(ast.traits),
    "never": [{"text": r.text, "when": when_of(r)} for r in ast.never_would_say],
    "quirks": list(ast.quirks),
    "tests": [
        {
            "name": t.name,
            "question": t.question,
            "fail": list(t.fail_examples),
            "pass": list(t.pass_examples),
            "why": t.why,
            "when": when_of(t),
        }
        for t in ast.tests
    ],
    "wrong": [
        {
            "wrong": getattr(w, "wrong", "") or getattr(w, "text", ""),
            "right": getattr(w, "right", "") or getattr(w, "correction", ""),
            # effigy has carried a why on WRONG pairs all along
            # (notation.py: WrongExampleAST.why) and this script dropped it.
            # A pair without one teaches a swap; the why is what makes it a
            # move the reader can apply to a sentence not in the card.
            "why": getattr(w, "why", ""),
            "when": when_of(w),
        }
        for w in ast.wrong_examples
    ],
    # MES is effigy's positive-exemplar block: whole turns in the wanted voice
    # rather than sentence-level corrections. Anthropic's Opus 5 guidance puts
    # these above prohibitions, and the card carried none until 2026-07-28.
    "mes": [getattr(m, "text", "") for m in ast.mes_examples],
}
# Keep in sync with internal/scan/render.go:maxNeverRules. Render() applies
# this cap at injection time (critical-prefixed rules first, then truncate),
# the same point effigy's own prompt.py applies MAX_NEVER_RULES. This script
# only warns, since card2json.py has no way to know which rule would matter
# most to keep.
#
# The cap is per injection, not per file. SessionStart prints the always-on
# rules; the mid-session reminder prints the evidence-gated ones that fired.
# Neither is allowed past the cap and the card file itself has no limit, so the
# two are counted apart here. Counting the file's total is what used to warn
# about a card that renders fine, and a warning that cries wolf on every build
# is one nobody reads on the build where it matters.
MAX_NEVER_RULES = 10

with open(out, "w") as fh:
    json.dump(payload, fh, indent=2)
    fh.write("\n")

print(f"{len(rules)} rules -> {out}")
for r in rules:
    print(f"  {r['action']:6s} {r['id']}")

def refresher_only(n):
    """True when the rule reaches the prompt through the refresher and not before.

    Mirrors internal/scan/render.go: SessionStart opens at_* and nothing else,
    because rule_* is evidence the writer has tripped a detector and mode_* is
    what kind of session this is, neither of which is known before the session
    has run. Both arrive later.

    Checking rule_ alone put a mode_* rule in the SessionStart pool, where it
    would be charged against a budget it cannot spend. Nothing overflowed on
    it, so the cost was a warning that would have fired early on a card one
    rule larger — which is the build where the warning is the whole point.
    """
    name = n.get("when", "").removeprefix("fact:")
    return name.startswith("rule_") or name.startswith("mode_")


for label, pool in (
    ("SessionStart", [n for n in payload["never"] if not refresher_only(n)]),
    ("the refresher", [n for n in payload["never"] if refresher_only(n)]),
):
    if len(pool) > MAX_NEVER_RULES:
        critical = sum(1 for n in pool if n["text"].upper().startswith("CRITICAL:"))
        dropped = len(pool) - max(MAX_NEVER_RULES, critical)
        print(
            f"warning: {len(pool)} NEVER rules render at {label}, cap is "
            f"{MAX_NEVER_RULES} — the last {dropped} will not reach the prompt"
        )
