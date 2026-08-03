// Package card carries the shipped voice card, compiled into the binary.
//
// The gate is installed with `go install`, which copies no data files, so a
// card read from disk at a path guessed from $HOME only ever worked on the
// machine that wrote it. Embedding removes the runtime file dependency
// entirely: a clone anywhere, or `go install github.com/justinstimatze/cope/
// cmd/cope-gate@latest` with no clone at all, gets a working gate.
//
// rules.json is generated from claude_voice.effigy by tools/card2json.py and
// committed so this embed has something to read. CI regenerates it and fails
// on any difference, which is what makes README's claim true — the rules being
// enforced and the rules being injected cannot drift apart.
package card

import _ "embed"

// RulesJSON is the compiled card. Callers pass it to scan.ParseCard.
//
//go:embed rules.json
var RulesJSON []byte
