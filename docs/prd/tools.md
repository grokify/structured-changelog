# PRD: CLI Tools for LLM-Assisted Changelog Generation

## Overview

This document describes CLI tools designed to improve the efficiency and consistency of LLM-assisted changelog generation. These tools reduce token usage, provide consistent parsing, and offer structured guidance for categorizing changes.

## Problem Statement

When using LLMs (like Claude) to analyze git history and generate changelogs:

1. **Token inefficiency**: Raw `git log` output is verbose, consuming 5x more tokens than necessary
2. **Inconsistent parsing**: LLMs may interpret the same commit differently across sessions
3. **Category ambiguity**: With 20 change types across 4 tiers, selecting the right category requires domain knowledge
4. **Validation feedback**: Current validation errors lack actionable guidance for fixes

## Proposed Tools

### 1. `schangelog parse-commits`

Parse git commits into structured JSON optimized for LLM consumption.

#### Use Cases

- Generate changelog entries for a new release
- Analyze commits between two versions
- Review changes since last tag

#### Input

```bash
# Parse commits since a tag
schangelog parse-commits --since=v0.3.0

# Parse commits between two refs
schangelog parse-commits --since=v0.2.0 --until=v0.3.0

# Parse last N commits
schangelog parse-commits --last=20

# Filter by path
schangelog parse-commits --since=v0.3.0 --path=src/
```

#### Output

Compact JSON optimized for token efficiency:

```json
{
  "repository": "github.com/grokify/structured-changelog",
  "range": {
    "since": "v0.3.0",
    "until": "HEAD",
    "commit_count": 12
  },
  "commits": [
    {
      "hash": "a1b2c3d",
      "short_hash": "a1b2c3d",
      "author": "John Wang",
      "date": "2026-01-04",
      "message": "feat(auth): add OAuth2 support",
      "type": "feat",
      "scope": "auth",
      "subject": "add OAuth2 support",
      "body": "Implements OAuth2 flow with PKCE.",
      "breaking": false,
      "issue": 123,
      "pr": null,
      "files_changed": 3,
      "insertions": 230,
      "deletions": 10,
      "files": [
        "src/auth/oauth.go",
        "src/auth/oauth_test.go"
      ],
      "suggested_category": "Added"
    }
  ],
  "summary": {
    "by_type": {
      "feat": 5,
      "fix": 3,
      "docs": 2,
      "test": 2
    },
    "by_suggested_category": {
      "Added": 5,
      "Fixed": 3,
      "Documentation": 2,
      "Tests": 2
    }
  }
}
```

#### Token Efficiency Analysis

| Format | Tokens per Commit | 50 Commits |
|--------|-------------------|------------|
| Raw `git log --stat` | ~150 | ~7,500 |
| Parsed JSON | ~30 | ~1,500 |
| **Savings** | **5x** | **6,000 tokens** |

#### Conventional Commit Parsing

Automatically extracts structured data from conventional commits:

| Pattern | Extracted Fields |
|---------|------------------|
| `feat(scope): message` | type=feat, scope=scope |
| `fix!: message` | type=fix, breaking=true |
| `BREAKING CHANGE:` in body | breaking=true |
| `Closes #123` | issue=123 |
| `Fixes #456` | issue=456 |
| `(#789)` in subject | pr=789 |

### 2. `schangelog suggest-category`

Suggest appropriate changelog categories for commit messages or descriptions.

#### Use Cases

- Help LLMs select the correct category from 20 options
- Validate category choices against tier requirements
- Provide reasoning for category selection

#### Input

```bash
# Suggest category for a commit message
schangelog suggest-category "feat(auth): add OAuth2 support"

# Suggest with context
schangelog suggest-category --context="security-related" "add input validation"

# Batch mode from stdin
echo "fix memory leak\nadd dark mode\nupdate README" | schangelog suggest-category --batch
```

#### Output

```json
{
  "input": "feat(auth): add OAuth2 support",
  "suggestions": [
    {
      "category": "Added",
      "tier": "core",
      "confidence": 0.95,
      "reasoning": "Conventional commit type 'feat' indicates new functionality"
    },
    {
      "category": "Security",
      "tier": "core",
      "confidence": 0.60,
      "reasoning": "Auth-related changes may have security implications"
    }
  ],
  "conventional_commit": {
    "type": "feat",
    "scope": "auth",
    "subject": "add OAuth2 support"
  }
}
```

#### Category Mapping Rules

| Commit Type | Primary Category | Tier | Notes |
|-------------|------------------|------|-------|
| `feat` | Added | core | New functionality |
| `fix` | Fixed | core | Bug fixes |
| `docs` | Documentation | extended | Documentation only |
| `style` | Internal | optional | Formatting, no logic change |
| `refactor` | Changed | core | Code restructuring |
| `perf` | Performance | standard | Performance improvements |
| `test` | Tests | extended | Test additions/changes |
| `build` | Build | extended | Build system changes |
| `ci` | Infrastructure | optional | CI/CD changes |
| `chore` | Internal | optional | Maintenance tasks |
| `security` | Security | core | Security fixes |
| `deps` | Dependencies | standard | Dependency updates |
| `breaking` | Breaking | standard | Breaking changes |

#### Benefits

1. **Consistency**: Same input always produces same suggestion
2. **Explainability**: LLM receives reasoning to validate or override
3. **Tier awareness**: Helps meet minimum tier coverage requirements

### 3. `schangelog validate` (Enhanced)

Enhanced validation with rich, actionable error messages.

#### Current State

```
Error: invalid date format for version 1.0.0
```

#### Enhanced Output

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
      "suggestion": "2026-01-04",
      "documentation": "https://keepachangelog.com/en/1.1.0/#how"
    }
  ],
  "warnings": [
    {
      "code": "W001",
      "severity": "warning",
      "path": "releases[0].security[0]",
      "message": "Security entry missing CVE identifier",
      "suggestion": "Add 'cve' field with format CVE-YYYY-NNNNN",
      "documentation": "https://cve.mitre.org/cve/identifiers/"
    }
  ],
  "summary": {
    "error_count": 1,
    "warning_count": 1,
    "releases_checked": 5,
    "entries_checked": 23
  }
}
```

#### Error Codes

| Code | Category | Description |
|------|----------|-------------|
| E001 | Format | Invalid date format |
| E002 | Format | Invalid version format (not semver) |
| E003 | Format | Invalid CVE format |
| E004 | Format | Invalid GHSA format |
| E005 | Structure | Missing required field |
| E006 | Structure | Duplicate version |
| E007 | Order | Releases not in descending order |
| E008 | Reference | Invalid issue/PR reference |
| W001 | Completeness | Security entry missing CVE |
| W002 | Completeness | Entry missing description |
| W003 | Coverage | Tier coverage below minimum |
| W004 | Style | Entry description too short |

#### Benefits

1. **Actionable**: Each error includes a suggestion for fixing
2. **Machine-readable**: JSON output for programmatic processing
3. **Educational**: Links to documentation for context
4. **Gradual strictness**: Separate errors from warnings

## Implementation Plan

### Phase 1: `parse-commits` Command

1. Create `gitlog` package in library
   - Commit struct with all fields
   - Parser for git log output
   - Conventional commit parser
   - Category suggester (basic)

2. Create CLI command
   - Flag handling (--since, --until, --last, --path)
   - Git command execution
   - JSON output formatting

3. Unit tests
   - Conventional commit parsing
   - Various git log formats
   - Edge cases (merge commits, empty messages)

### Phase 2: `suggest-category` Command

1. Extend `changelog` package
   - Category matching rules
   - Confidence scoring
   - Reasoning generation

2. Create CLI command
   - Single and batch modes
   - Context flag for hints

3. Unit tests
   - All commit type mappings
   - Confidence scoring
   - Edge cases

### Phase 3: Enhanced `validate` Command

1. Extend `changelog` package
   - Structured error types
   - Error code registry
   - Suggestion generator

2. Update CLI command
   - JSON output mode
   - Enhanced error formatting
   - Documentation links

3. Unit tests
   - All error codes
   - Suggestion generation
   - Various malformed inputs

## Success Metrics

| Metric | Target |
|--------|--------|
| Token reduction | 5x for commit parsing |
| Category suggestion accuracy | >90% for conventional commits |
| Unit test coverage | >95% for new library code |
| Validation error actionability | 100% errors have suggestions |

## Future Considerations

### MCP Server

Once CLI tools are stable, consider wrapping them in an MCP server for tighter LLM integration:

```
Tools:
- parse_commits(since, until, path) -> CommitList
- suggest_category(message, context) -> Suggestions
- validate_changelog(json) -> ValidationResult
- render_markdown(json, options) -> Markdown
```

Benefits of MCP vs CLI:

- Richer type information
- Stateful interactions
- Direct tool integration without shell

### Interactive Mode

Future `schangelog` could offer an interactive mode for LLM-assisted changelog creation:

```bash
schangelog generate --interactive --since=v0.3.0
```

This would:

1. Parse commits
2. Present each to LLM for categorization
3. Generate entries with LLM assistance
4. Validate and output final changelog
