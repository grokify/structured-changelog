# PRISM Control Integration

This guide covers integrating Structured Changelog with PRISM Control for roadmap traceability — linking changelog entries to roadmap items (RMIs) and strategic initiatives.

## Overview

PRISM Control is a roadmap management system that organizes work into:

- **Initiatives** — strategic goals spanning multiple projects (e.g., `INIT-STREAMING-001`)
- **Roadmap Items (RMIs)** — discrete deliverables tracked per-repository (e.g., `RMI-MYREPO-042`)

Structured Changelog supports optional `rmis`, `rmi`, and `initiative` fields on changelog entries, enabling:

- Traceability from releases back to roadmap planning
- Automated initiative progress tracking across repositories
- Changelog-driven roadmap status updates

## JSON Schema

Add `rmis` (preferred) and/or `initiative` to any changelog entry. Use the
plural `rmis` array because a curated entry often summarizes several commits
spanning multiple RMIs; the legacy singular `rmi` is retained for
back-compatibility and treated as a one-element `rmis`.

```json
{
  "releases": [
    {
      "version": "v1.2.0",
      "date": "2026-07-26",
      "added": [
        {
          "description": "Add streaming support via ConverseStream API",
          "rmis": ["RMI-MYREPO-042", "RMI-MYREPO-044"],
          "commit": "abc123"
        },
        {
          "description": "Add bearer token authentication",
          "rmis": ["RMI-MYREPO-043"]
        },
        {
          "description": "Ad-hoc refactor with no roadmap item",
          "initiative": "INIT-STREAMING-001"
        }
      ]
    }
  ]
}
```

All three fields are optional. Follow the convention: **populate `rmis`**;
set `initiative` **only** for entries that have no RMI (with an RMI present the
initiative is derivable downstream, so storing both invites drift).

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `rmis` | string[] | Roadmap Item IDs (format: `RMI-<REPOSLUG>-<NNN>`); **preferred** |
| `rmi` | string | Legacy single Roadmap Item ID; treated as a one-element `rmis` |
| `initiative` | string | Initiative ID (format: `INIT-<SLUG>-<NNN>`); use only when no RMI applies |

`schangelog validate` warns (W006) on malformed IDs in `rmis` and (W007) when a
singular `rmi` disagrees with `rmis`.

## Programmatic Usage

### Building Entries

Use the builder methods for fluent entry construction:

```go
import "github.com/grokify/structured-changelog/changelog"

entry := changelog.NewEntry("Add streaming support").
    WithRMIs("RMI-MYREPO-042", "RMI-MYREPO-044").
    WithCommit("abc123").
    WithAuthor("@developer")

// Legacy singular form (still supported); prefer WithRMIs for new code:
legacy := changelog.NewEntry("Add auth").WithRMI("RMI-MYREPO-043")
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

The `parse-commits` command extracts these references automatically, carrying
the RMI IDs on each parsed commit's `rmis`:

```bash
schangelog parse-commits --since=v1.1.0
```

When you build a changelog from those commits with `schangelog init --from-tags`,
each generated entry's `rmis` is **pre-populated** from the underlying commits'
`Refs:` trailers (deduplicated and sorted) — you curate rather than transcribe.
`initiative` is never auto-populated.

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

1. **Prefer `rmis` over `rmi`** — The plural array expresses entries that summarize several commits across multiple RMIs; the singular remains for back-compat
2. **Set `initiative` only without an RMI** — With an RMI present the initiative is derivable downstream, so storing both invites drift
3. **Let trailers do the work** — Carry `Refs: RMI-<SLUG>-<NNN>` on commits; `parse-commits` and `init` populate `rmis` for you
4. **Omit when not planned** — Bug fixes and dependency updates typically don't need roadmap references
5. **Keep IDs stable** — Don't change RMI/Initiative IDs after they're referenced in changelogs

## Links

- [JSON IR Specification](../specification/spec.md)
- [Release Notes Guide](release-notes-guide.md)
- [Full Changelog](../changelog.md)
