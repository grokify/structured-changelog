# Release Notes Guide

This guide provides a recommended structure for `RELEASE_NOTES_vX.Y.Z.md` files that complement your structured changelog.

## Philosophy

**CHANGELOG.md** answers: *"What changed?"*
**RELEASE_NOTES** answers: *"How do I upgrade?"*

Release notes are version-specific documents that provide the context, migration guides, and code examples that don't belong in a concise changelog.

## Recommended Structure

```markdown
# Release Notes vX.Y.Z

**Release Date:** YYYY-MM-DD

## Overview
Brief summary of the release (2-3 sentences).

## Breaking Changes
Detailed explanation of breaking changes with migration guides.

## New Features
Feature descriptions with usage examples.

## Migration
Step-by-step upgrade instructions with before/after code.

## Dependencies
Dependency changes and their implications.

## Known Issues
Current limitations or issues to be aware of.

## What's Next
Roadmap items planned for future releases.
```

## Section Guidelines

### Overview

A brief, high-level summary of the release. Answer:

- What is the main theme of this release?
- Who should upgrade immediately?
- Are there breaking changes?

**Example:**

```markdown
## Overview

This release refactors the provider architecture so that LLM observability
SDKs own their own adapters. This simplifies dependency management and
allows each SDK to be updated independently.

**Breaking changes:** Import paths have changed. See Migration section.
```

### Breaking Changes

For each breaking change:

1. Explain what changed and why
2. Show the before/after difference
3. Provide migration steps

Use tables for multiple related changes:

```markdown
## Breaking Changes

### Import Paths Changed

| Before | After |
|--------|-------|
| `import "pkg/old/path"` | `import "pkg/new/path"` |
| `import "pkg/old/other"` | `import "pkg/new/other"` |

The API remains identical—only import paths change.
```

### New Features

For each significant feature:

1. Describe what it does
2. Show a usage example
3. Link to documentation if applicable

```markdown
## New Features

### Annotation Support

Added `AnnotationManager` interface for span/trace annotations:

​```go
annotation := provider.CreateAnnotation(ctx, &Annotation{
    SpanID:  spanID,
    Name:    "quality_score",
    Value:   0.95,
})
​```

See [docs/annotations.md](docs/annotations.md) for full API reference.
```

### Migration

Provide step-by-step upgrade instructions:

```markdown
## Migration

### Step 1: Update Dependencies

​```bash
go get github.com/example/pkg@v2.0.0
​```

### Step 2: Update Imports

​```go
// Before
import "github.com/example/pkg/v1"

// After
import "github.com/example/pkg/v2"
​```

### Step 3: Update API Calls

The `DoThing()` method now requires a context:

​```go
// Before
result := client.DoThing(input)

// After
result := client.DoThing(ctx, input)
​```
```

### Dependencies

List dependency changes with context:

```markdown
## Dependencies

### Added

- `github.com/spf13/cobra` v1.10.2 - CLI framework

### Updated

- `golang.org/x/crypto` v0.17.0 → v0.18.0 (security fix)

### Removed

- `github.com/old/dep` - No longer needed after refactor
```

### Known Issues

Be transparent about limitations:

```markdown
## Known Issues

- The `--watch` flag is not yet implemented (#123)
- Large files (>100MB) may cause memory issues (#124)
```

### What's Next

Give users visibility into the roadmap:

```markdown
## What's Next

Planned for v2.1.0:

- Watch mode for automatic regeneration
- YAML input support
- GitHub Actions integration
```

## Template

Use the structure above as a template, or see the examples below.

## Naming Convention

Use the format: `RELEASE_NOTES_vX.Y.Z.md`

Examples:

- `RELEASE_NOTES_v1.0.0.md`
- `RELEASE_NOTES_v2.1.0.md`
- `RELEASE_NOTES_v0.5.0-beta.1.md`

## When to Create Release Notes

Create release notes for:

- **Major versions** (v1.0.0, v2.0.0) — Always
- **Minor versions with breaking changes** — Always
- **Minor versions with significant features** — Recommended
- **Patch versions** — Optional, usually not needed

## Relationship to CHANGELOG

| Content | CHANGELOG.json | RELEASE_NOTES |
|---------|---------------|---------------|
| Entry exists | ✓ Required | ✓ Expanded |
| Breaking flag | `"breaking": true` | Full migration guide |
| Code examples | ✗ Never | ✓ Always for features |
| Dependencies | Version only | Context + implications |
| File lists | ✗ Never | ✓ For major changes |

## Automation

While CHANGELOG.md is generated from CHANGELOG.json, release notes are typically hand-written. However, you can bootstrap them:

```bash
# Generate a starting point from changelog
schangelog generate CHANGELOG.json --format=release-notes --version=1.0.0 > RELEASE_NOTES_v1.0.0.md
```

Then expand with migration guides, code examples, and context.

## Examples

See real-world examples:

- [v0.1.0 Release Notes](../releases/v0.1.0.md) — This project's initial release
- [v0.4.0 Release Notes](../releases/v0.4.0.md) — Tier-based change types
