.PHONY: all build test lint coverage clean sync-check docs docs-serve help

# Default target
all: sync-check lint test build

# Version stamp: last tag + commits-ahead + short sha + dirty flag
# (e.g. v0.16.0-3-gabc1234-dirty). Falls back to the toolchain's build info
# when git is unavailable.
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null)
COMMIT  := $(shell git rev-parse HEAD 2>/dev/null)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Build the CLI (stamped with git-describe version metadata)
build:
	go build -ldflags "$(LDFLAGS)" -o bin/schangelog ./cmd/schangelog

# Run tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML coverage report: go tool cover -html=coverage.out"

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Check that CHANGE_TYPES.json files are in sync
sync-check:
	@echo "Checking CHANGE_TYPES.json sync..."
	@diff -q CHANGE_TYPES.json changelog/change_types.json > /dev/null 2>&1 || \
		(echo "ERROR: CHANGE_TYPES.json and changelog/change_types.json are out of sync!" && \
		 echo "Run 'make sync' to fix." && exit 1)
	@echo "CHANGE_TYPES.json files are in sync."

# Sync CHANGE_TYPES.json to changelog/change_types.json
sync:
	cp CHANGE_TYPES.json changelog/change_types.json
	@echo "Synced CHANGE_TYPES.json to changelog/change_types.json"

# Generate example markdown files
examples:
	./bin/schangelog generate examples/basic/CHANGELOG.json -o examples/basic/CHANGELOG.md
	./bin/schangelog generate examples/security/CHANGELOG.json -o examples/security/CHANGELOG.md
	./bin/schangelog generate examples/full/CHANGELOG.json -o examples/full/CHANGELOG.md
	./bin/schangelog generate examples/extended/CHANGELOG.json -o examples/extended/CHANGELOG.md

# Build documentation (MkDocs)
docs:
	mkdocs build

# Serve documentation locally
docs-serve:
	mkdocs serve

# Help
help:
	@echo "Available targets:"
	@echo "  all         - Run sync-check, lint, test, and build (default)"
	@echo "  build       - Build the schangelog CLI"
	@echo "  test        - Run tests"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  lint        - Run golangci-lint"
	@echo "  clean       - Remove build artifacts"
	@echo "  sync-check  - Verify CHANGE_TYPES.json files are in sync"
	@echo "  sync        - Copy CHANGE_TYPES.json to changelog/change_types.json"
	@echo "  examples    - Generate example markdown files"
	@echo "  docs        - Build documentation (MkDocs)"
	@echo "  docs-serve  - Serve documentation locally"
	@echo "  help        - Show this help message"
