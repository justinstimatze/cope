VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
EFFIGY_PATH ?= ../effigy

.PHONY: build install test rules check-rules check readme demo cards styles backfill pairs discriminate

CARDS_DIR ?= $(or $(XDG_CONFIG_HOME),$(HOME)/.config)/cope/cards

# build and install do not regenerate the card. card/rules.json is committed and
# compiled into the binary, so a clone with no effigy checkout still builds —
# which was not true while install depended on the rules target.
build:
	go build $(LDFLAGS) -o bin/cope-gate ./cmd/cope-gate

install:
	go install $(LDFLAGS) ./cmd/cope-gate

test:
	go test ./...

# rules regenerates the JSON compiled into the binary. It goes through effigy,
# which stays the authority on the notation: https://github.com/justinstimatze/effigy
# The gate reads .effigy directly at runtime (internal/effigy) and a test holds
# the two parsers to the same output, so this target is for the embedded card
# rather than for anything a reader writes.
rules:
	EFFIGY_PATH=$(EFFIGY_PATH) python3 tools/card2json.py card/claude_voice.effigy card/rules.json

# cards copies the demo cards where --card can find them by name. Nothing here
# is required to use your own: drop any .effigy in $(CARDS_DIR) and it is a
# voice, or point --rules at one anywhere on disk.
#
#   make cards && cope-gate --cards
#   COPE_CARD=lecturer cope-gate --inject
cards:
	@mkdir -p "$(CARDS_DIR)"
	cp card/demo/*.effigy "$(CARDS_DIR)/"
	@echo "installed to $(CARDS_DIR)"

# styles writes every card as a Claude Code output style, which is where a voice
# actually takes. --inject puts the card in one turn-zero message and it gets
# buried; an output style goes in the system prompt and the harness re-reminds
# the model of it. Measured 2026-08-03: through the hook a card asking for a
# bolded label on every paragraph produced prose indistinguishable from no card.
#
# Emitting a style does not turn it on. Pick one with /config, or set
# "outputStyle" in .claude/settings.local.json.
styles: build
	@./bin/cope-gate --output-style
	@for f in card/demo/*.effigy; do ./bin/cope-gate --rules "$$f" --output-style; done

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

# pairs runs the blind preference test: two replies to the same real prompt out
# of your own transcripts, one written with the card and one without, handed
# over unlabelled. It answers the question the hit counts only stand in for.
#
# BASELINE is required and is whatever prose guidance you already have without
# a card — for most people the Response style section of ~/.claude/CLAUDE.md,
# saved to its own file. The card is tested on top of it, so the run measures
# what the card adds rather than what any instruction at all adds.
#
#   make pairs BASELINE=~/.claude/response-style.md
#   make pairs BASELINE=~/.claude/response-style.md CARD=card/demo/laconic.effigy
#
# CARD points at any .effigy file, which is how you test your own card before
# adopting it. PAIRFLAGS passes anything else through, e.g. PAIRFLAGS="-seed 3".
TRANSCRIPTS ?= $(HOME)/.claude/projects/*/*.jsonl
PAIRS_OUT   ?= pairs
PAIRFLAGS   ?= -test-pairs 20 -control-pairs 10

pairs: build
	@test -n "$(BASELINE)" || { echo "pairs needs BASELINE=<file>; see the comment above this target in the Makefile" >&2; exit 2; }
	cd replay && go build -o ../bin/cope-replay .
ifdef CARD
	@tmp=$$(mktemp) && trap 'rm -f "$$tmp"' EXIT && \
	EFFIGY_PATH=$(EFFIGY_PATH) python3 tools/card2json.py $(CARD) "$$tmp" >/dev/null && \
	./bin/cope-replay -pairs -gate ./bin/cope-gate -rules "$$tmp" \
		-transcripts '$(TRANSCRIPTS)' -baseline $(BASELINE) -outdir $(PAIRS_OUT) $(PAIRFLAGS)
else
	./bin/cope-replay -pairs -gate ./bin/cope-gate \
		-transcripts '$(TRANSCRIPTS)' -baseline $(BASELINE) -outdir $(PAIRS_OUT) $(PAIRFLAGS)
endif

# discriminate runs the blind discrimination test. Same generation as pairs —
# two replies to the same real prompt under two different cards — but the reader
# is shown one card's VOICE block and asked which reply carries it, instead of
# which one reads better.
#
# The two questions come apart, and only this one can answer "make it sound like
# this instead". A voice can land completely and leave a reader indifferent
# about whether it was worth landing, so a flat preference result says nothing
# about whether the voicing arrived. Chance here is 50%: above it the voice is
# on the page, at it nothing reached the prose whatever the card said.
#
#   make discriminate RIVAL=card/demo/laconic.effigy
#   make discriminate CARD=card/demo/precise.effigy RIVAL=card/demo/laconic.effigy
#   make discriminate RIVAL=card/demo/laconic.effigy DISCFLAGS="-judge -test-pairs 40"
#
# CARD defaults to the card built into the gate; RIVAL is the second voice. The
# control set is card vs bare framing, the widest gap available — if that does
# not separate the instrument is blind and the test set is unreadable.
#
# -judge answers every pair with the model on the same materials, so revising a
# card costs a run rather than an evening of reading. The page reports its
# agreement with your own calls, which is the only thing that makes its rate a
# fact about the voice rather than about a model.
DISC_OUT  ?= discriminate
DISCFLAGS ?= -test-pairs 20 -control-pairs 10

discriminate: build
	@test -n "$(RIVAL)" || { echo "discriminate needs RIVAL=<file.effigy>, the second voice to be told apart from; see the comment above this target" >&2; exit 2; }
	cd replay && go build -o ../bin/cope-replay .
	@set -e; \
	rules=""; rival=""; \
	trap 'rm -f "$$rules" "$$rival"' EXIT; \
	rival=$$(mktemp); \
	EFFIGY_PATH=$(EFFIGY_PATH) python3 tools/card2json.py $(RIVAL) "$$rival" >/dev/null; \
	if [ -n "$(CARD)" ]; then \
		rules=$$(mktemp); \
		EFFIGY_PATH=$(EFFIGY_PATH) python3 tools/card2json.py $(CARD) "$$rules" >/dev/null; \
	fi; \
	./bin/cope-replay -pairs -discriminate -gate ./bin/cope-gate \
		$${rules:+-rules $$rules} -rival "$$rival" \
		-transcripts '$(TRANSCRIPTS)' -outdir $(DISC_OUT) $(DISCFLAGS)

# check scores prose against the card — the same pass a reply gets at Stop.
# The repo's own docs go through it; F defaults to the README.
check: build
	./bin/cope-gate --check $(or $(F),README.md) --log ""

# PYTHON is the interpreter the doc targets use. A distro python has been
# externally managed since PEP 668, so `pip install anthropic` against it is
# refused and this target died with an unhelpful "pip install anthropic" the
# first time the system python moved up a version. A venv is the answer the
# error message itself gives, so prefer one if it is there and say what to do
# when it is not.
VENV ?= $(HOME)/.venvs/cope
PYTHON ?= $(if $(wildcard $(VENV)/bin/python),$(VENV)/bin/python,python3)

# readme regenerates README.md from the card. Needs ANTHROPIC_API_KEY in .env
# and the anthropic package on $(PYTHON); see tools/generate_readme.py.
#
#   python3 -m venv ~/.venvs/cope && ~/.venvs/cope/bin/pip install anthropic
#
# CARD and TARGET point it at a different card, which is how demo/ is built:
#   make readme CARD=card/demo/precise.effigy TARGET=demo/README.precise.md
#
# With no CARD, this writes the FRONT PAGE, and the front page is deliberately
# a demo card rather than the shipped one: a reader has to be able to see that a
# card wrote the page, and a page in the shipped register cannot show that.
#
# It was the maximal card until 2026-08-04, on the argument that a README
# written in the register cope exists to fix states the problem better than a
# paragraph about it can. It does, and it costs too much: the people who arrive
# here arrived to get away from exactly that prose, and the front page served
# them a full page of it before they had any reason to trust the joke. The
# maximal render stays under demo/ and the front page names it, so the reader
# meets it by clicking rather than by landing.
#
# fieldguide replaced it. It is visibly card-shaped — entries open on the name
# of the thing, lookalikes get a sentence beginning "Compare", an unmeasured
# claim reads "Not recorded" — so it still shows a card at work, and none of
# those moves hurt to read. caveman was the other candidate and was dropped:
# there is a separate project of that name by another author, cited in the
# README, and wearing a neighbour's register on the front door reads as taking
# it.
#
# generate_readme.py tells the prompt which of the two it is writing, because
# the instruction to send a reader to the demo set becomes a link to the page
# they are standing on the moment the front page is one of those cards.
readme: build
	@$(PYTHON) -c 'import anthropic' 2>/dev/null || { \
	  echo "no anthropic package on $(PYTHON)"; \
	  echo "  python3 -m venv $(VENV) && $(VENV)/bin/pip install anthropic"; \
	  exit 1; }
ifdef CARD
	$(PYTHON) tools/generate_readme.py --rules $(CARD) --target $(TARGET) --rounds $(or $(ROUNDS),3)
else
	$(PYTHON) tools/generate_readme.py --rules card/demo/fieldguide.effigy --target README.md --rounds $(or $(ROUNDS),3)
endif

# demo writes the README again from every card in card/demo/ and scores them
# all, so a reader can A/B the shipped card against the others and see the
# gate's verdict on each. caveman cuts an answer by removing things;
# precise cuts it by replacing them; lecturer does not cut at all and differs
# from the shipped card only in register, which is what makes it the one usable
# as a discrimination rival. fieldguide is also the front page, so it renders
# twice, to README.md and to demo/README.fieldguide.md.
#
# fieldguide scores high in the check below and the number is not what it looks
# like. It declines labelled_opening with @gate, because every entry opens on
# the name of the thing — and the check runs without --rules, so all seven pages
# are scored against the SHIPPED card's rules, which still counts the rule this
# card refused. That is the documented behaviour: a declined rule keeps running
# and only its own card's score drops it. Under its own rules the page scores
# lowest of the set.
#
# Each card goes to the gate as .effigy, which internal/effigy reads directly.
# This ran through tools/card2json.py until 2026-08-02 and so needed an effigy
# checkout to render a demo — a build step between having a card and reading
# what it does, in the one target whose whole purpose is showing that.
# The six renders share nothing but the binary and the API key, so they run
# concurrently. Sequentially this target took about ten minutes; three at once
# measured 203s on 2026-08-03, and the wall clock is now the slowest single
# render rather than the sum. Backgrounded explicitly rather than through -j so
# it does not depend on how the caller invokes make.
demo: build
	@set -e; pids=""; \
	for spec in \
	  "card/demo/precise.effigy:demo/README.precise.md" \
	  "card/demo/caveman.effigy:demo/README.caveman.md" \
	  "card/demo/lecturer.effigy:demo/README.lecturer.md" \
	  "card/demo/fieldguide.effigy:demo/README.fieldguide.md" \
	  "card/demo/claude_maximal.effigy:demo/README.claude-maximal.md" \
	  "card/claude_voice.effigy:demo/README.claude-voice.md" ; do \
	  card=$${spec%%:*}; target=$${spec##*:}; \
	  echo "  render $$card -> $$target"; \
	  $(PYTHON) tools/generate_readme.py --rules "$$card" --target "$$target" --rounds 0 & \
	  pids="$$pids $$!"; \
	done; \
	fail=0; for p in $$pids; do wait $$p || fail=1; done; \
	test $$fail -eq 0 || { echo "at least one render failed" >&2; exit 1; }
	@echo
	$(PYTHON) tools/demo_index.py
	@echo
	@for f in README.md demo/README.*.md; do ./bin/cope-gate --check $$f --log "" | head -1; done
