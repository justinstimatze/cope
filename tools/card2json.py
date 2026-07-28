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

# effigy's @when condition gates are deliberately not carried across. A blind
# check over three registers (result report, disagreement, build proposal)
# found one register serving all three, so what varies is dictated by the
# speech act and not by voice. The gates were removed from the card; this drops
# the plumbing that outlived them.

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
    "traits": list(ast.traits),
    "never": [r.text for r in ast.never_would_say],
    "quirks": list(ast.quirks),
    "tests": [
        {
            "name": t.name,
            "question": t.question,
            "fail": list(t.fail_examples),
            "pass": list(t.pass_examples),
            "why": t.why,
        }
        for t in ast.tests
    ],
    "wrong": [
        {
            "wrong": getattr(w, "wrong", "") or getattr(w, "text", ""),
            "right": getattr(w, "right", "") or getattr(w, "correction", ""),
        }
        for w in ast.wrong_examples
    ],
}
# Keep in sync with internal/scan/render.go:maxNeverRules. Render() applies
# this cap at injection time (critical-prefixed rules first, then truncate),
# the same point effigy's own prompt.py applies MAX_NEVER_RULES. This script
# only warns, since card2json.py has no way to know which rule would matter
# most to keep.
MAX_NEVER_RULES = 10

with open(out, "w") as fh:
    json.dump(payload, fh, indent=2)
    fh.write("\n")

print(f"{len(rules)} rules -> {out}")
for r in rules:
    print(f"  {r['action']:6s} {r['id']}")

never_count = len(payload["never"])
if never_count > MAX_NEVER_RULES:
    critical = sum(1 for n in payload["never"] if n.upper().startswith("CRITICAL:"))
    dropped = never_count - max(MAX_NEVER_RULES, critical)
    print(
        f"warning: {never_count} NEVER rules, cap is {MAX_NEVER_RULES} — "
        f"the last {dropped} will not render at injection time"
    )
