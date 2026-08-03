#!/usr/bin/env python3
"""Write README.md by giving the card to a model and the result back to the gate.

The prompt comes from `cope-gate --author-docs`: the rendered voice card, the
facts introspected from the card and the flag set, and the section list. This
script contributes no prose. That rule is effigy's, from its own
generate_readme.py — "The generator provides facts, not prose" — and it is the
reason the README reads like the card rather than like a template.

What cope adds is the second half. `cope-gate --check` scores the draft with the
same rules that score a reply at Stop, and any violations go back to the model
as a revision turn. The documentation is held to the card it documents, and the
loop runs until the gate is quiet or the round budget is spent.

    pip install anthropic
    python3 tools/generate_readme.py --rounds 3

ANTHROPIC_API_KEY comes from the environment or from a .env beside this repo.
"""

import argparse
import os
import shutil
import subprocess
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
DEFAULT_MODEL = "claude-opus-5"
DEFAULT_TARGET = REPO / "README.md"


def load_env():
    """Read KEY=value lines from .env without overriding a real environment."""
    env_path = REPO / ".env"
    if not env_path.exists():
        return
    for line in env_path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        os.environ.setdefault(k.strip(), v.strip().strip("'\""))


def gate():
    """Locate cope-gate: a local build first, then $PATH."""
    local = REPO / "bin" / "cope-gate"
    if local.exists():
        return str(local)
    found = shutil.which("cope-gate")
    if not found:
        sys.exit("error: no cope-gate on $PATH and no bin/cope-gate — run `make build`")
    return found


def run_gate(*args):
    proc = subprocess.run([gate(), *args], capture_output=True, text=True)
    return proc.stdout, proc.returncode


def check(target):
    """Score the file. Returns the gate's own report, or '' when clean."""
    out, _ = run_gate("--check", str(target), "--log", "")
    return "" if ": clean" in out else out


def response_text(resp):
    """Join the text blocks of a reply.

    Not content[0]: a model with extended thinking on puts a ThinkingBlock
    first, and indexing blind raises AttributeError halfway through a paid
    call. Selecting by type works whether or not thinking is enabled.
    """
    parts = [b.text for b in resp.content if getattr(b, "type", "") == "text"]
    if not parts:
        sys.exit("error: the model returned no text block")
    return "\n".join(parts).strip()


def main():
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--model", default=DEFAULT_MODEL)
    ap.add_argument("--target", default=str(DEFAULT_TARGET), type=Path)
    ap.add_argument("--rules", default="",
                    help="compile the prompt from this rules JSON instead of the built-in card")
    ap.add_argument("--rounds", type=int, default=3,
                    help="revision rounds after the first draft (default 3). Use 0 for a demo "
                         "card, where the point is to see what the first draft scores")
    ap.add_argument("--dry-run", action="store_true",
                    help="print the prompt and its size, call nothing")
    args = ap.parse_args()

    gate_args = ["--author-docs"]
    if args.rules:
        gate_args = ["--rules", args.rules] + gate_args
    system, rc = run_gate(*gate_args)
    if rc != 0 or not system.strip():
        sys.exit("error: cope-gate --author-docs produced nothing")

    if args.dry_run:
        print(system)
        print(f"\n--- system prompt: {len(system)} chars ---", file=sys.stderr)
        return

    load_env()
    if not os.environ.get("ANTHROPIC_API_KEY"):
        sys.exit("error: set ANTHROPIC_API_KEY, or put it in a .env at the repo root")

    try:
        import anthropic
    except ImportError:
        sys.exit("pip install anthropic")

    client = anthropic.Anthropic()
    # The card and the facts are identical on every round, so the prefix is
    # marked once and each revision pays ~10% on it instead of full price.
    system_blocks = [{
        "type": "text",
        "text": system,
        "cache_control": {"type": "ephemeral"},
    }]
    # The prompt needs to know whether this render IS the front page, because
    # the instruction to point at the maximal demo becomes a link to the page
    # you are standing on the moment the front page is written from that card.
    front_page = args.target.resolve() == DEFAULT_TARGET.resolve()
    where = ("This render IS the repository front page — the file a visitor sees "
             "first on GitHub." if front_page else
             f"This render is a demo copy written to {args.target.name}, one of "
             "several of the same page in different voices. It is not the front page.")
    messages = [{
        "role": "user",
        "content": f"Write README.md now. {where} Output raw markdown only — no "
                   "wrapping code fence, no preamble, no explanation of what you wrote.",
    }]

    draft = ""
    for round_no in range(args.rounds + 1):
        resp = client.messages.create(
            model=args.model,
            max_tokens=16384,
            system=system_blocks,
            messages=messages,
        )
        draft = response_text(resp)
        u = resp.usage
        cached = getattr(u, "cache_read_input_tokens", 0) or 0
        print(f"round {round_no}: {len(draft)} chars "
              f"(in {u.input_tokens} + {cached} cached, out {u.output_tokens})",
              file=sys.stderr)

        args.target.write_text(draft + "\n")
        report = check(args.target)
        if not report:
            print(f"gate is clean after round {round_no} → {args.target}", file=sys.stderr)
            return
        print(report, file=sys.stderr)
        if round_no == args.rounds:
            break
        messages += [
            {"role": "assistant", "content": draft},
            {"role": "user", "content":
                "cope-gate scored that draft against the card and reported this:\n\n"
                f"{report}\n"
                "Rewrite the offending sentences. Keep every fact and every section. "
                "Do not delete a section to silence a rule. Output the whole README "
                "again, raw markdown only."},
        ]

    if args.rounds == 0:
        print(f"{args.target} written unrevised — the score above is the demo", file=sys.stderr)
        return
    print(f"round budget spent; {args.target} still has hits above", file=sys.stderr)
    sys.exit(1)


if __name__ == "__main__":
    main()
