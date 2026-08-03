#!/usr/bin/env python3
"""Rewrite the measured tables in demo/README.md from the files themselves.

Every other page under demo/ is generated. The index was not, so its tables were
hand-maintained numbers describing files the build produces — they went stale on
every regen, silently, because nothing compares a number in prose against the
thing it counts. Six separate corrections on 2026-08-03 alone.

Run by `make demo` after the renders land. Only the blocks between the marker
comments are touched; the prose around them is written by hand and stays that
way.
"""
from __future__ import annotations

import re
import subprocess
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
INDEX = REPO / "demo" / "README.md"
GATE = REPO / "bin" / "cope-gate"

# card path -> rendered page, in the order the tables list them
CARDS = [
    ("card/claude_voice.effigy", "demo/README.claude-voice.md"),
    ("card/demo/claude_maximal.effigy", "demo/README.claude-maximal.md"),
    ("card/demo/laconic.effigy", "demo/README.laconic.md"),
    ("card/demo/precise.effigy", "demo/README.precise.md"),
    ("card/demo/caveman.effigy", "demo/README.caveman.md"),
    ("card/demo/lecturer.effigy", "demo/README.lecturer.md"),
]

EMOJI = re.compile(
    "[\U0001F300-\U0001FAFF☀-➿]"
)


def counts(path: Path) -> dict[str, int]:
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines()
    return {
        "bold": sum(1 for ln in lines if re.match(r"\*\*[^*]+\*\*", ln)),
        "emdash": text.count("—"),
        "emojihead": sum(1 for ln in lines if ln.startswith("#") and EMOJI.search(ln)),
        "hr": sum(1 for ln in lines if re.fullmatch(r"-{3,}\s*", ln)),
        "tables": sum(1 for ln in lines if ln.startswith("|---")),
        "callouts": sum(1 for ln in lines if ln.startswith("> **")),
        "bodyemoji": len(re.findall(r"[✅⚠️\U0001F3AF]", text)),
        "apologies": len(
            re.findall(
                r"I apologise|my mistake|I should have|that was mine",
                text,
                re.I,
            )
        ),
        "chars": len(text.encode("utf-8")),
    }


def rules(page: str) -> list[str]:
    """Rule ids the gate reports for one page, one entry per hit."""
    out = subprocess.run(
        [str(GATE), "--check", page, "--log", ""],
        capture_output=True,
        text=True,
        cwd=REPO,
    ).stdout
    return re.findall(r"^  \[([a-z_]+)\]", out, re.M)


def block(name: str, body: str) -> str:
    return f"<!-- {name}:start -->\n{body}\n<!-- {name}:end -->"


def splice(text: str, name: str, body: str) -> str:
    pat = re.compile(
        rf"<!-- {re.escape(name)}:start -->.*?<!-- {re.escape(name)}:end -->",
        re.S,
    )
    if not pat.search(text):
        sys.exit(
            f"error: no <!-- {name}:start --> ... <!-- {name}:end --> markers in "
            f"{INDEX.relative_to(REPO)} — add them around the table this writes"
        )
    return pat.sub(lambda _: block(name, body), text)


def main() -> None:
    if not GATE.exists():
        sys.exit(f"error: {GATE} not built — run make build first")
    for _, page in CARDS:
        if not (REPO / page).exists():
            sys.exit(f"error: {page} missing — run make demo first")

    c = {page: counts(REPO / page) for _, page in CARDS}
    mx = c["demo/README.claude-maximal.md"]
    cv = c["demo/README.claude-voice.md"]

    rows = [
        ("paragraphs opening on a bold label", "bold"),
        ("em dashes", "emdash"),
        ("headings carrying an emoji", "emojihead"),
        ("horizontal rules", "hr"),
        ("tables for things that are not tabular", "tables"),
        ("blockquote callouts", "callouts"),
        ("emoji in the body text", "bodyemoji"),
        ("sentences apologising or taking the blame", "apologies"),
    ]
    side = ["| | `claude_maximal` | `claude_voice` |", "|---|---|---|"]
    side += [f"| {label} | {mx[k]} | {cv[k]} |" for label, k in rows]

    sizes = ["| Card | Output | Characters |", "|---|---|---|"]
    for card, page in CARDS:
        name = Path(page).name
        sizes.append(f"| `{card}` | [`{name}`]({name}) | {c[page]['chars']:,} |")

    hits = rules("demo/README.claude-maximal.md")
    tally: dict[str, int] = {}
    for r in hits:
        tally[r] = tally.get(r, 0) + 1
    maximal = ["| rule | hits on `claude_maximal` |", "|---|---|"]
    maximal += [
        f"| `{r}` | {n} |" for r, n in sorted(tally.items(), key=lambda kv: -kv[1])
    ]

    text = INDEX.read_text(encoding="utf-8")
    text = splice(text, "side-by-side", "\n".join(side))
    text = splice(text, "sizes", "\n".join(sizes))
    text = splice(text, "maximal-hits", "\n".join(maximal))
    INDEX.write_text(text, encoding="utf-8")

    total = sum(v["chars"] for v in c.values())
    allhits: dict[str, int] = {}
    for _, page in CARDS:
        for r in rules(page):
            allhits[r] = allhits.get(r, 0) + 1
    print(f"demo/README.md tables rewritten from {len(CARDS)} renders")
    print(f"  total {total:,} chars")
    print(
        "  rules fired: "
        + ", ".join(f"{r} {n}" for r, n in sorted(allhits.items(), key=lambda kv: -kv[1]))
    )
    print("  prose figures outside the marked blocks are still hand-written —")
    print("  check them against the line above before committing")


if __name__ == "__main__":
    main()
