# PRISM Control Integration

This guide covers integrating Structured Changelog with PRISM Control for roadmap traceability — linking changelog entries to roadmap items (RMIs) and strategic initiatives.

## Overview

PRISM Control is a roadmap management system that organizes work into:

- **Initiatives** — strategic goals spanning multiple projects (e.g., `INIT-STREAMING-001`)
- **Roadmap Items (RMIs)** — discrete deliverables tracked per-repository (e.g., `RMI-MYREPO-042`)

Structured Changelog supports optional `rmi` and `initiative` fields on changelog entries, enabling:

- Traceability from releases back to roadmap planning
- Automated initiative progress tracking across repositories
- Changelog-driven roadmap status updates

## JSON Schema

Add `rmi` and/or `initiative` to any changelog entry:

```json
{
  "releases": [
    {
      "version": "v1.2.0",
      "date": "2026-07-26",
      "added": [
        {
          "description": "Add streaming support via ConverseStream API",
          "rmi": "RMI-MYREPO-042",
          "initiative": "INIT-STREAMING-001",
          "commit": "abc123"
        },
        {
          "description": "Add bearer token authentication",
          "rmi": "RMI-MYREPO-043"
        }
      ]
    }
  ]
}
```

Both fields are optional and independent — use one, both, or neither as needed.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `rmi` | string | Roadmap Item ID (format: `RMI-<REPOSLUG>-<NNN>`) |
| `initiative` | string | Initiative ID (format: `INIT-<SLUG>-<NNN>`) |

## Programmatic Usage

### Building Entries

Use the builder methods for fluent entry construction:

```go
import "github.com/grokify/structured-changelog/changelog"

entry := changelog.NewEntry("Add streaming support").
    WithRMI("RMI-MYREPO-042").
    WithInitiative("INIT-STREAMING-001").
    WithCommit("abc123").
    WithAuthor("@developer")
```

### Reading Entries

Access the fields directly from parsed changelogs:

```go
cl, err := changelog.LoadFile("CHANGELOG.json")
if err != nil {
    return err
}

for _, release := range cl.Releases {
    for _, entry := range release.Added {
        if entry.RMI != "" {
            fmt.Printf("Feature %s implements %s\n", entry.Description, entry.RMI)
        }
        if entry.Initiative != "" {
            fmt.Printf("  Part of initiative: %s\n", entry.Initiative)
        }
    }
}
```

### Filtering by RMI or Initiative

Find all entries for a specific roadmap item:

```go
func findEntriesByRMI(cl *changelog.Changelog, rmi string) []changelog.Entry {
    var results []changelog.Entry
    for _, release := range cl.Releases {
        for _, entry := range release.AllEntries() {
            if entry.RMI == rmi {
                results = append(results, entry)
            }
        }
    }
    return results
}
```

## Git Commit Integration

Following the convention in CLAUDE.md, commits implementing roadmap items carry a git trailer:

```
feat(api): add streaming support

Implements server-sent events for real-time updates.

Refs: RMI-MYREPO-042
```

The `parse-commits` command can extract these references for changelog generation:

```bash
schangelog parse-commits --since=v1.1.0
```

When generating changelog entries from commits with `Refs:` trailers, include the RMI in the entry's `rmi` field.

## Use Cases

### Roadmap Progress Tracking

Query changelogs across a portfolio to track initiative completion:

```go
// Load multiple project changelogs
portfolio := []*changelog.Changelog{project1, project2, project3}

// Count completed RMIs per initiative
progress := make(map[string]int)
for _, cl := range portfolio {
    for _, release := range cl.Releases {
        for _, entry := range release.AllEntries() {
            if entry.Initiative != "" {
                progress[entry.Initiative]++
            }
        }
    }
}
```

### Release Notes Generation

Generate initiative-focused release notes:

```go
// Group entries by initiative for quarterly reports
byInitiative := make(map[string][]changelog.Entry)
for _, entry := range release.AllEntries() {
    if entry.Initiative != "" {
        byInitiative[entry.Initiative] = append(
            byInitiative[entry.Initiative], entry)
    }
}
```

### Audit Trail

The combination of `rmi`, `initiative`, and `commit` fields provides a complete audit trail:

- **What** was delivered (entry description)
- **Why** it was done (initiative context)
- **Where** it was planned (RMI reference)
- **When** it shipped (release date)
- **How** it was implemented (commit hash)

## Best Practices

1. **Use RMI for discrete deliverables** — Each RMI should map to one or a few changelog entries
2. **Use Initiative for strategic context** — Group related work across multiple releases
3. **Include both when applicable** — An entry can reference both its specific RMI and the broader initiative
4. **Omit when not planned** — Bug fixes and dependency updates typically don't need roadmap references
5. **Keep IDs stable** — Don't change RMI/Initiative IDs after they're referenced in changelogs

## Links

- [JSON IR Specification](../specification/spec.md)
- [Release Notes Guide](release-notes-guide.md)
- [Full Changelog](../changelog.md)
