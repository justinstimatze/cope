#!/usr/bin/env python3
"""Run a card's POSTPROC rules over the assistant turns of a Claude Code transcript.

Uses effigy's own validators path (validators_from_ast + validate), not a
re-implementation, so what gets tested is the code the hook would call.
"""
import json
import sys
from collections import Counter, defaultdict

from effigy_bootstrap import require_effigy

require_effigy("run_postproc.py")

from effigy.parser import parse_file  # noqa: E402
from effigy.validators import validate, validators_from_ast  # noqa: E402

CARD = sys.argv[1]
TRANSCRIPT = sys.argv[2]

ast = parse_file(CARD)
validators = validators_from_ast(ast)
print(f"loaded {len(validators)} POSTPROC rules: {[v.rule_id for v in validators]}\n")

# Assistant text blocks only. Tool calls, tool results, and user turns are
# excluded — the hook would only ever see the drafted reply.
turns = []
with open(TRANSCRIPT) as fh:
    for line in fh:
        try:
            rec = json.loads(line)
        except json.JSONDecodeError:
            continue
        msg = rec.get("message") or {}
        if rec.get("type") != "assistant" or msg.get("role") != "assistant":
            continue
        content = msg.get("content")
        if not isinstance(content, list):
            continue
        text = "\n".join(
            b.get("text", "") for b in content
            if isinstance(b, dict) and b.get("type") == "text"
        ).strip()
        if text:
            turns.append(text)

print(f"{len(turns)} assistant text turns, {sum(len(t) for t in turns):,} chars\n")

hits = Counter()
examples = defaultdict(list)
turns_with_any = 0
for i, text in enumerate(turns):
    violations = validate(text, ast, validators)
    if violations:
        turns_with_any += 1
    for v in violations:
        hits[v.rule_id] += 1
        start, end = v.span
        ctx = text[max(0, start - 70):min(len(text), end + 70)].replace("\n", " ")
        examples[v.rule_id].append((i, v.matched_text, ctx))

print(f"turns with at least one hit: {turns_with_any}/{len(turns)}\n")
for v in validators:
    rid = v.rule_id
    print(f"=== {rid}  ({hits[rid]} hits)  /{v.pattern}/")
    for turn_i, matched, ctx in examples[rid][:6]:
        print(f"  turn {turn_i}: MATCH {matched!r}")
        print(f"    ...{ctx}...")
    if len(examples[rid]) > 6:
        print(f"  (+{len(examples[rid]) - 6} more)")
    print()
