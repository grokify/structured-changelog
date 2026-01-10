# TOON Integration Design Doc

## Overview

This document describes the integration of TOON (Token-Oriented Object Notation) as the default output format for LLM-facing commands in sclog.

**Goal:** Reduce token usage by ~40% for LLM-assisted changelog generation workflows.

**Breaking Change:** Yes (pre-1.0, acceptable per semver)

## Background

### Current State

LLM-facing commands output JSON:

- `parse-commits` - Structured git history
- `suggest-category` - Category suggestions
- `validate --json` - Rich validation errors

The LLM guide documents ~5x token reduction vs raw git log. TOON can provide additional ~40% reduction.

### What is TOON?

TOON (Token-Oriented Object Notation) is a compact, human-readable encoding of the JSON data model optimized for LLM prompts.

Key features:

- Indentation-based (no braces)
- Tabular arrays (fields declared once, rows streamed)
- ~40% fewer tokens than JSON
- 74% vs 70% accuracy in LLM benchmarks
- Lossless JSON round-trips

Reference: [toonformat.dev](https://toonformat.dev/)

## Design

### Format Flag

Add `--format` flag to LLM-facing commands:

```bash
# Default: TOON (optimized for LLMs)
sclog parse-commits --since=v0.4.0

# Explicit JSON for tooling/debugging
sclog parse-commits --since=v0.4.0 --format=json
```

### Supported Values

| Value | Description |
|-------|-------------|
| `toon` | TOON format (default) |
| `json` | JSON with indentation |
| `json-compact` | JSON without indentation |

### Commands Affected

| Command | Default Before | Default After |
|---------|----------------|---------------|
| `parse-commits` | JSON | TOON |
| `suggest-category` | JSON | TOON |
| `validate --json` | JSON | TOON |
| `generate` | Markdown | Markdown (unchanged) |

### Example Output

#### parse-commits

**JSON (current):**

```json
{
  "repository": "github.com/example/project",
  "range": {
    "since": "v0.4.0",
    "until": "HEAD",
    "commit_count": 3
  },
  "commits": [
    {
      "hash": "abc123def456",
      "short_hash": "abc123d",
      "author": "John Doe",
      "date": "2026-01-04",
      "message": "feat(auth): add OAuth2 support",
      "type": "feat",
      "scope": "auth",
      "subject": "add OAuth2 support",
      "files_changed": 3,
      "insertions": 230,
      "deletions": 45,
      "suggested_category": "Added"
    },
    {
      "hash": "def456abc789",
      "short_hash": "def456a",
      "author": "Jane Smith",
      "date": "2026-01-03",
      "message": "fix: resolve memory leak in cache",
      "type": "fix",
      "subject": "resolve memory leak in cache",
      "files_changed": 1,
      "insertions": 12,
      "deletions": 8,
      "suggested_category": "Fixed"
    }
  ],
  "summary": {
    "by_type": {"feat": 1, "fix": 1},
    "by_suggested_category": {"Added": 1, "Fixed": 1}
  }
}
```

**TOON (new default):**

```
repository:github.com/example/project
range{since,until,commit_count}:
  v0.4.0,HEAD,3
commits[2]{hash,short_hash,author,date,type,scope,subject,files_changed,insertions,deletions,suggested_category}:
  abc123def456,abc123d,John Doe,2026-01-04,feat,auth,add OAuth2 support,3,230,45,Added
  def456abc789,def456a,Jane Smith,2026-01-03,fix,,resolve memory leak in cache,1,12,8,Fixed
summary:
  by_type{feat,fix}:
    1,1
  by_suggested_category{Added,Fixed}:
    1,1
```

#### suggest-category

**JSON (current):**

```json
{
  "input": "feat(auth): add OAuth2 support",
  "suggestions": [
    {
      "category": "Added",
      "tier": "core",
      "confidence": 0.95,
      "reasoning": "Conventional commit type 'feat' indicates new functionality"
    }
  ],
  "conventional_commit": {
    "type": "feat",
    "scope": "auth",
    "subject": "add OAuth2 support"
  }
}
```

**TOON (new default):**

```
input:feat(auth): add OAuth2 support
suggestions[1]{category,tier,confidence,reasoning}:
  Added,core,0.95,Conventional commit type 'feat' indicates new functionality
conventional_commit{type,scope,subject}:
  feat,auth,add OAuth2 support
```

#### validate (with errors)

**JSON (current):**

```json
{
  "valid": false,
  "errors": [
    {
      "code": "E001",
      "severity": "error",
      "path": "releases[0].date",
      "message": "Invalid date format",
      "actual": "January 4, 2026",
      "expected": "YYYY-MM-DD format (ISO 8601)",
      "suggestion": "Convert to ISO 8601 format: YYYY-MM-DD"
    }
  ],
  "summary": {
    "error_count": 1,
    "warning_count": 0,
    "releases_checked": 3,
    "entries_checked": 15
  }
}
```

**TOON (new default):**

```
valid:false
errors[1]{code,severity,path,message,actual,expected,suggestion}:
  E001,error,releases[0].date,Invalid date format,January 4, 2026,YYYY-MM-DD format (ISO 8601),Convert to ISO 8601 format: YYYY-MM-DD
summary{error_count,warning_count,releases_checked,entries_checked}:
  1,0,3,15
```

## Implementation

### Dependencies

Add `github.com/toon-format/toon-go` to go.mod.

### Code Changes

#### 1. Add format package

Create `format/format.go`:

```go
package format

import (
    "encoding/json"
    "github.com/toon-format/toon-go"
)

type Format string

const (
    FormatTOON        Format = "toon"
    FormatJSON        Format = "json"
    FormatJSONCompact Format = "json-compact"
)

func Marshal(v any, f Format) ([]byte, error) {
    switch f {
    case FormatTOON:
        return toon.Marshal(v)
    case FormatJSON:
        return json.MarshalIndent(v, "", "  ")
    case FormatJSONCompact:
        return json.Marshal(v)
    default:
        return toon.Marshal(v) // default to TOON
    }
}

func ParseFormat(s string) (Format, error) {
    switch s {
    case "toon", "":
        return FormatTOON, nil
    case "json":
        return FormatJSON, nil
    case "json-compact":
        return FormatJSONCompact, nil
    default:
        return "", fmt.Errorf("unknown format: %s", s)
    }
}
```

#### 2. Update commands

Add `--format` flag to each command:

```go
var outputFormat string

func init() {
    parseCommitsCmd.Flags().StringVar(&outputFormat, "format", "toon",
        "Output format: toon (default), json, json-compact")
}
```

Replace JSON marshaling:

```go
// Before
jsonBytes, err := json.MarshalIndent(result, "", "  ")

// After
f, _ := format.ParseFormat(outputFormat)
output, err := format.Marshal(result, f)
```

#### 3. Add TOON struct tags

Update structs with `toon` tags where field names differ:

```go
type Commit struct {
    Hash              string   `json:"hash" toon:"hash"`
    ShortHash         string   `json:"short_hash" toon:"short_hash"`
    // ... etc
}
```

### File Changes

| File | Change |
|------|--------|
| `go.mod` | Add toon-go dependency |
| `format/format.go` | New - format abstraction |
| `format/format_test.go` | New - format tests |
| `cmd/sclog/parse_commits.go` | Add --format flag, use format.Marshal |
| `cmd/sclog/suggest_category.go` | Add --format flag, use format.Marshal |
| `cmd/sclog/validate.go` | Add --format flag for JSON output |
| `gitlog/commit.go` | Add toon struct tags |
| `changelog/validate_rich.go` | Add toon struct tags |

### Testing

1. Unit tests for format package
2. Update existing command tests to cover both formats
3. Golden file tests comparing JSON vs TOON output
4. Round-trip tests (TOON → JSON → TOON)

## Migration

### Breaking Change Handling

Since we're pre-1.0, breaking changes are expected. Document in CHANGELOG.json:

```json
{
  "version": "0.6.0",
  "date": "2026-01-XX",
  "breaking": [
    {
      "description": "LLM commands now output TOON format by default; use --format=json for JSON",
      "breaking": true
    }
  ]
}
```

### User Migration

Scripts using JSON output add `--format=json`:

```bash
# Before
sclog parse-commits --since=v0.4.0 | jq '.commits'

# After
sclog parse-commits --since=v0.4.0 --format=json | jq '.commits'
```

## Rollout

1. Implement format package with tests
2. Add --format flag to all three commands
3. Update documentation (README, LLM guide)
4. Update CHANGELOG.json
5. Release as v0.6.0

## Decisions

1. **Flag name:** `--format` (avoids confusion with `--output` which often means output file)

2. **Remove --compact:** Early stage, no backwards compatibility needed. Use `--format=json-compact` instead.

3. **validate --json flag:** Replace with `--format`. The `--json` flag is removed.
