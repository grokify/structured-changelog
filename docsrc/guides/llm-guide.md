# LLM-Assisted Changelog Generation

This guide covers tools and prompts for using LLMs (like Claude) to generate changelogs from git history.

## Overview

The `sclog` CLI includes tools optimized for LLM-assisted workflows:

| Tool | Purpose | Token Savings |
|------|---------|---------------|
| `parse-commits` | Convert git log to structured JSON | ~5x reduction |
| `suggest-category` | Classify commits into changelog categories | Consistent mapping |
| `validate --json` | Rich error output with suggestions | Actionable fixes |

## Tools

### parse-commits

Parses git history into structured JSON optimized for LLM consumption.

```bash
# Parse commits since a tag
sclog parse-commits --since=v0.3.0

# Parse commits between versions
sclog parse-commits --since=v0.2.0 --until=v0.3.0

# Parse last N commits
sclog parse-commits --last=20

# Compact output without file list
sclog parse-commits --since=v0.3.0 --no-files --compact

# Exclude merge commits
sclog parse-commits --since=v0.3.0 --no-merges
```

**Output includes:**

- Conventional commit parsing (type, scope, subject)
- Breaking change detection (`!` or `BREAKING CHANGE:`)
- Issue/PR references extracted from messages
- File statistics (insertions, deletions, files changed)
- Suggested changelog category for each commit
- Summary statistics grouped by type and category

**Example output:**

```json
{
  "repository": "github.com/example/project",
  "range": {
    "since": "v0.3.0",
    "until": "HEAD",
    "commit_count": 5
  },
  "commits": [
    {
      "hash": "abc123d",
      "author": "John Doe",
      "date": "2026-01-04",
      "message": "feat(auth): add OAuth2 support",
      "type": "feat",
      "scope": "auth",
      "subject": "add OAuth2 support",
      "breaking": false,
      "issue": 123,
      "files_changed": 3,
      "insertions": 230,
      "suggested_category": "Added"
    }
  ],
  "summary": {
    "by_type": { "feat": 3, "fix": 2 },
    "by_suggested_category": { "Added": 3, "Fixed": 2 }
  }
}
```

### suggest-category

Suggests changelog categories for commit messages based on conventional commit types and keywords.

```bash
# Single message
sclog suggest-category "feat(auth): add OAuth2 support"

# Batch mode from stdin
git log --format="%s" v0.3.0..HEAD | sclog suggest-category --batch

# Compact JSON output
sclog suggest-category --compact "fix: resolve memory leak"
```

**Example output:**

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
      "reasoning": "Feature relates to authentication or security"
    }
  ],
  "conventional_commit": {
    "type": "feat",
    "scope": "auth",
    "subject": "add OAuth2 support"
  }
}
```

**Category mapping:**

| Commit Type | Category | Tier |
|-------------|----------|------|
| `feat` | Added | core |
| `fix` | Fixed | core |
| `docs` | Documentation | extended |
| `perf` | Performance | standard |
| `test` | Tests | extended |
| `build` | Build | extended |
| `ci` | Infrastructure | optional |
| `chore` | Internal | optional |
| `refactor` | Changed | core |
| `security` | Security | core |
| `deps` | Dependencies | standard |

### validate --json

Validates changelog with rich, actionable error messages.

```bash
# JSON output with detailed errors
sclog validate --json CHANGELOG.json

# Strict mode (warnings become errors)
sclog validate --json --strict CHANGELOG.json
```

**Example error output:**

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
      "suggestion": "Convert to ISO 8601 format: YYYY-MM-DD",
      "documentation": "https://keepachangelog.com/en/1.1.0/#how"
    }
  ],
  "warnings": [
    {
      "code": "W001",
      "severity": "warning",
      "path": "releases[0].security[0]",
      "message": "Security entry missing CVE identifier",
      "suggestion": "Add 'cve' field with format CVE-YYYY-NNNNN"
    }
  ],
  "summary": {
    "error_count": 1,
    "warning_count": 1,
    "releases_checked": 3,
    "entries_checked": 15
  }
}
```

**Error codes:**

| Code | Description |
|------|-------------|
| E001 | Invalid date format |
| E002 | Invalid version format |
| E003 | Invalid CVE format |
| E004 | Invalid GHSA format |
| E005 | Invalid severity level |
| E006 | Invalid CVSS score |
| E007 | Invalid IR version |
| E100 | Missing required field |
| E101 | Duplicate version |
| E103 | Empty description |
| W001 | Security entry missing CVE |
| W002 | Description too short |
| W003 | Tier coverage below minimum |
| W004 | Missing severity |

## Example Prompts

Use these prompts with Claude or other LLMs to generate changelogs.

### Generate a Changelog Release

```
Generate changelog entries for v0.5.0 based on commits since v0.4.0.

1. Run `sclog parse-commits --since=v0.4.0` to get structured commit data
2. Review the commits and create appropriate CHANGELOG.json entries
3. Group related commits into single entries where appropriate
4. Write descriptions that explain "why" not just "what"
5. Validate with `sclog validate --json`
```

### Review and Categorize Changes

```
Parse the commits since v0.4.0 and help me categorize them for the changelog.

Use `sclog parse-commits --since=v0.4.0` and then:
- Group related commits together
- Identify any breaking changes
- Suggest which category each change belongs to
- Flag anything that might need security review
```

### Fix Validation Errors

```
Validate CHANGELOG.json and fix any issues:

1. Run `sclog validate --json CHANGELOG.json`
2. For each error, apply the suggested fix
3. Re-validate until clean
```

### Complete Release Workflow

```
Help me prepare the v0.5.0 release:

1. Parse commits: `sclog parse-commits --since=v0.4.0`
2. Create CHANGELOG.json entries for the new version
3. Validate: `sclog validate --json CHANGELOG.json`
4. Generate markdown: `sclog generate CHANGELOG.json -o CHANGELOG.md`
5. Summarize what's in the release
```

### Batch Processing

```
I have a backlog of releases to document. For each version tag pair,
parse the commits and generate changelog entries:

- v0.1.0 to v0.2.0
- v0.2.0 to v0.3.0
- v0.3.0 to v0.4.0

Use `sclog parse-commits --since=<from> --until=<to>` for each range.
```

## Token Efficiency

Raw `git log` output is verbose. The `parse-commits` tool reduces token usage significantly:

| Format | Tokens per Commit | 50 Commits |
|--------|-------------------|------------|
| Raw `git log --stat` | ~150 | ~7,500 |
| Parsed JSON | ~30 | ~1,500 |
| **Savings** | **5x** | **6,000 tokens** |

Use `--no-files` and `--compact` for further reduction when file lists aren't needed.

## Tips

1. **Start with parse-commits** — Get structured data before asking the LLM to categorize
2. **Use batch mode** — Process multiple commits at once with `suggest-category --batch`
3. **Validate early** — Run `validate --json` to catch issues before they accumulate
4. **Trust the suggestions** — The category suggestions are accurate for conventional commits (~95%)
5. **Override when needed** — The LLM can apply judgment for ambiguous cases
