#!/usr/bin/env bash
# Score the largest recent transcripts across projects to see whether the
# gate's hit rate is a property of one conversation or of everything written.
set -uo pipefail

GATE="${GATE:-$HOME/go/bin/cope-gate}"
ROOT="${ROOT:-$HOME/.claude/projects}"
N="${1:-6}"
# SKIP excludes transcripts whose path matches, for leaving out the session
# doing the measuring. Empty by default: a hardcoded id silently dropped an
# arbitrary transcript on every machine but the one it was written on.
SKIP="${SKIP:-}"
HOME_PREFIX="$(echo "$HOME" | tr '/' '-')-"

find "$ROOT" -name '*.jsonl' -size +200k -printf '%s\t%p\n' \
  | sort -rn \
  | { if [ -n "$SKIP" ]; then grep -v "$SKIP"; else cat; fi; } \
  | head -n "$N" \
  | cut -f2 \
  | while read -r f; do
      proj=$(basename "$(dirname "$f")")
      echo "=== ${proj#"$HOME_PREFIX"} / $(basename "$f" | cut -c1-8)"
      "$GATE" --backfill "$f" 2>/dev/null \
        | grep -E '^(turns with|[0-9]+ assistant|=== |\(never fired)'
      echo
    done
