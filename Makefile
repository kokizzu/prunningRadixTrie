GO ?= go
GOVULNCHECK ?= govulncheck
TMPDIR ?= /tmp
GOCACHE ?= $(TMPDIR)/prunningRadixTrie-gocache
CMD ?=

.PHONY: test bench verify-dependency-security vulncheck run

test:
	$(GO) test ./...

bench:
	$(GO) test -bench=.

vulncheck:
	mkdir -p "$(GOCACHE)"
	GOCACHE="$(GOCACHE)" $(GOVULNCHECK) ./...

verify-dependency-security: vulncheck

run:
	@test -n "$(CMD)" || (echo "usage: make run CMD='go test ./...'" >&2; exit 2)
	$(CMD)
