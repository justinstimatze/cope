VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
EFFIGY_PATH ?= ../effigy

.PHONY: build install test rules check-rules backfill

# build and install do not regenerate the card. card/rules.json is committed and
# compiled into the binary, so a clone with no effigy checkout still builds —
# which was not true while install depended on the rules target.
build:
	go build $(LDFLAGS) -o bin/cope-gate ./cmd/cope-gate

install:
	go install $(LDFLAGS) ./cmd/cope-gate

test:
	go test ./...

# rules is for editing the card. It needs effigy, which stays the only thing
# that parses .effigy notation: https://github.com/justinstimatze/effigy
rules:
	EFFIGY_PATH=$(EFFIGY_PATH) python3 tools/card2json.py card/claude_voice.effigy card/rules.json

# check-rules fails when the committed JSON no longer matches the card it came
# from. CI runs it, which is what makes "the rules being enforced and the rules
# being injected cannot drift apart" a fact rather than a convention.
check-rules:
	@tmp=$$(mktemp) && trap 'rm -f "$$tmp"' EXIT && \
	EFFIGY_PATH=$(EFFIGY_PATH) python3 tools/card2json.py card/claude_voice.effigy "$$tmp" >/dev/null && \
	if diff -u card/rules.json "$$tmp"; then \
		echo "card/rules.json is up to date"; \
	else \
		echo "card/rules.json is stale — run 'make rules' and commit the result" >&2; \
		exit 1; \
	fi

backfill: build
	./bin/cope-gate --backfill $(T)
