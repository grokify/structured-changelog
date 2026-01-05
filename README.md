# Structured Changelog

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Coverage][coverage-svg]][coverage-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

**A canonical, deterministic changelog framework with JSON IR and Markdown rendering.**

Structured Changelog provides a machine-readable JSON Intermediate Representation (IR) as the source of truth for your changelog, with deterministic Markdown generation for human readability. It supports optional security metadata (CVE, GHSA, SARIF) and SBOM information.

## Overview

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  CHANGELOG.json │ ───► │    Renderer     │ ───► │  CHANGELOG.md   │
│  (Canonical IR) │      │ (Deterministic) │      │ (Human-Readable)│
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

### Key Principles

1. **JSON IR is canonical** — Markdown is derived, not the source of truth
2. **Deterministic rendering** — Same input always produces identical output
3. **Keep a Changelog format** — Compatible with [keepachangelog.com](https://keepachangelog.com/) - [`github.com/olivierlacan/keep-a-changelog`](https://github.com/olivierlacan/keep-a-changelog)
4. **Semantic Versioning** — Follows [semver.org](https://semver.org/) conventions
5. **Extensible metadata** — Optional security (CVE/GHSA/SARIF) and SBOM fields
6. **Spec + tooling together** — Single source of truth for humans and machines

### Relationship to Keep a Changelog

[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) is the de-facto standard for human-readable changelogs. Structured Changelog builds on that foundation:

> **Keep a Changelog is the human spec; Structured Changelog is the structured implementation.**

| Concern | Keep a Changelog | Structured Changelog |
|---------|------------------|----------------------|
| Human-friendly text | ✓ | ✓ (via Markdown renderer) |
| Machine-readable format | ✗ | ✓ (JSON IR) |
| Deterministic rendering | ✗ | ✓ |
| Security metadata | ✗ | ✓ (CVE/GHSA/SARIF) |
| SBOM metadata | ✗ | ✓ |
| Tooling / APIs | ✗ | ✓ (Go packages + CLI) |

**What Keep a Changelog defines:**

- Section structure (`## [Unreleased]`, `## [1.0.0] - YYYY-MM-DD`)
- Categories (Added, Changed, Deprecated, Removed, Fixed, Security)
- Semantic versioning and date formatting

**What Structured Changelog adds:**

- JSON schema for machine parsing
- Deterministic ordering rules
- Security fields (CVE, GHSA, severity, CVSS, CWE, SARIF)
- SBOM fields (component, version, license)
- CLI and Go library for automation

The generated `CHANGELOG.md` conforms to Keep a Changelog 1.1.0 formatting conventions.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap grokify/tap
brew install structured-changelog
```

This installs the `sclog` CLI (also available as `structured-changelog`).

### Go Install

```bash
go install github.com/grokify/structured-changelog/cmd/sclog@latest
```

### Go Library

```bash
go get github.com/grokify/structured-changelog
```

## Quick Start

### Define your changelog in JSON

```json
{
  "ir_version": "1.0",
  "project": "my-project",
  "releases": [
    {
      "version": "1.0.0",
      "date": "2026-01-03",
      "added": [
        { "description": "Initial release with core features" }
      ]
    }
  ]
}
```

### Generate Markdown

```go
package main

import (
    "fmt"
    "os"

    "github.com/grokify/structured-changelog/changelog"
    "github.com/grokify/structured-changelog/renderer"
)

func main() {
    // Load changelog from JSON
    cl, err := changelog.LoadFile("CHANGELOG.json")
    if err != nil {
        panic(err)
    }

    // Render to Markdown
    md := renderer.RenderMarkdown(cl)
    fmt.Println(md)
}
```

### CLI Usage

Validate a changelog:

```bash
sclog validate CHANGELOG.json
```

Generate Markdown:

```bash
# Output to stdout
sclog generate CHANGELOG.json

# Output to file
sclog generate CHANGELOG.json -o CHANGELOG.md

# Minimal output (no references/metadata)
sclog generate CHANGELOG.json --minimal

# Full output (include commit SHAs)
sclog generate CHANGELOG.json --full
```

Show version:

```bash
sclog version
```

### LLM-Assisted Generation

The CLI includes tools optimized for AI-assisted changelog generation, reducing token usage by ~5x:

```bash
sclog parse-commits --since=v0.3.0     # Structured git history
sclog suggest-category "feat: ..."     # Category suggestions
sclog validate --json CHANGELOG.json   # Rich error output
```

See the [LLM Guide](https://grokify.github.io/structured-changelog/guides/llm-guide/) for prompts and workflows.

## JSON IR Schema

### Change Types

Structured Changelog supports 20 change types organized into 4 tiers. The **core** tier contains the standard [Keep a Changelog](https://keepachangelog.com/) categories, while higher tiers provide extended functionality.

#### Tiers

| Tier | Description |
|------|-------------|
| **core** | Standard types defined by Keep a Changelog (KACL) |
| **standard** | Commonly used by major providers and popular open source projects |
| **extended** | Change metadata for documentation, build, and acknowledgments |
| **optional** | For deployment teams and internal operational visibility |

#### Canonical Ordering

The following table shows all change types in canonical order, grouped by purpose:

```
┌─ OVERVIEW & CRITICAL ─────────────────────────────────────────┐
│  1. Highlights      standard   Release summaries/key takeaways│
│  2. Breaking        standard   Backward-incompatible changes  │
│  3. Upgrade Guide   standard   Migration instructions         │
│  4. Security        core       Vulnerabilities/CVE fixes      │
├─ CORE KACL ───────────────────────────────────────────────────┤
│  5. Added           core       New features                   │
│  6. Changed         core       Modified functionality         │
│  7. Deprecated      core       Future removal warnings        │
│  8. Removed         core       Removed features               │
│  9. Fixed           core       Bug fixes                      │
├─ QUALITY ─────────────────────────────────────────────────────┤
│ 10. Performance     standard   Speed/efficiency improvements  │
│ 11. Dependencies    standard   Dependency updates             │
├─ DEVELOPMENT ─────────────────────────────────────────────────┤
│ 12. Documentation   extended   Docs updates                   │
│ 13. Build           extended   CI/CD and tooling              │
│ 14. Tests           extended   Test additions and coverage    │
├─ OPERATIONS ──────────────────────────────────────────────────┤
│ 15. Infrastructure  optional   Deployment/hosting changes     │
│ 16. Observability   optional   Logging/metrics/tracing        │
│ 17. Compliance      optional   Regulatory updates             │
├─ INTERNAL ────────────────────────────────────────────────────┤
│ 18. Internal        optional   Refactors only (not tests)     │
├─ END MATTER ──────────────────────────────────────────────────┤
│ 19. Known Issues    extended   Caveats and limitations        │
│ 20. Contributors    extended   Acknowledgments                │
└───────────────────────────────────────────────────────────────┘
```

#### Tier-Based Filtering

Use tiers to control which change types to include:

```bash
# Validate: ensure changelog covers at least core types
sclog validate --min-tier core

# Validate: require core + standard coverage
sclog validate --min-tier standard

# Generate: output only core types (KACL-compliant)
sclog generate --max-tier core

# Generate: include everything up to extended
sclog generate --max-tier extended
```

### Optional Security Metadata

Entries can include security-specific fields:

```json
{
  "description": "Fix SQL injection vulnerability",
  "cve": "CVE-2026-12345",
  "ghsa": "GHSA-xxxx-xxxx-xxxx",
  "severity": "high"
}
```

### Optional SBOM Metadata

Track component changes with SBOM fields:

```json
{
  "description": "Update dependency",
  "component": "example-lib",
  "version": "2.0.0",
  "license": "MIT"
}
```

## CHANGELOG vs RELEASE_NOTES

Structured Changelog addresses `CHANGELOG.md` specifically. Understanding the difference between changelogs and release notes is important:

| Aspect | CHANGELOG.md | RELEASE_NOTES_vX.Y.Z.md |
|--------|--------------|-------------------------|
| **Purpose** | Cumulative history of all changes | Version-specific upgrade guide |
| **Format** | Concise bullet points | Detailed narrative with examples |
| **Audience** | Scanning/discovery | Users upgrading to specific version |
| **Content** | What changed | Why it changed + how to migrate |
| **Scope** | All versions in one file | Single version per file |
| **Examples** | Brief or none | Code samples, migration guides |
| **Structure** | Standardized (Keep a Changelog) | Flexible, project-specific |

### Content Placement Guidelines

| Content Type | CHANGELOG.json | RELEASE_NOTES |
|--------------|----------------|---------------|
| Breaking change flag | ✓ `"breaking": true` | Detailed explanation |
| Migration guide | ✗ | ✓ With code examples |
| Code examples | ✗ | ✓ Before/after blocks |
| API changes | Brief description | Full context + examples |
| Dependency updates | Version numbers | Implications + upgrade steps |

**Example: Breaking change workflow**

In `CHANGELOG.json` — flag it concisely:

```json
{
  "changed": [
    { "description": "Module renamed to go-opik", "breaking": true }
  ]
}
```

In `RELEASE_NOTES_v0.5.0.md` — provide migration details:

```markdown
## Migration

Update all imports:

​```go
// Before
import opik "github.com/example/old-name"

// After
import opik "github.com/example/go-opik"
​```
```

### When to Use Each

**CHANGELOG.md:**

- Quick reference for all historical changes
- Automated tooling (version bumps, release automation)
- Scanning to find when a feature was added/removed

**RELEASE_NOTES_vX.Y.Z.md:**

- Breaking change migration guides with before/after code
- Detailed feature explanations with usage examples
- Dependency change implications
- File-level change lists for major releases

See [docsrc/guides/release-notes-guide.md](docsrc/guides/release-notes-guide.md) for recommended release notes structure.

## Project Structure

```
structured-changelog/
├── README.md
├── LICENSE
├── go.mod
├── changelog/          # JSON IR structs + validation
│   ├── changelog.go
│   ├── entry.go
│   ├── release.go
│   └── validate.go
├── gitlog/             # Git log parsing for LLM workflows
│   ├── commit.go
│   ├── conventional.go
│   ├── category.go
│   └── parser.go
├── renderer/           # Deterministic Markdown renderer
│   ├── markdown.go
│   └── options.go
├── cmd/sclog/          # CLI tool (Cobra-based)
│   ├── main.go
│   ├── root.go
│   ├── validate.go
│   ├── generate.go
│   ├── parse_commits.go
│   └── suggest_category.go
├── schema/             # JSON Schema definitions
│   └── changelog-v1.schema.json
├── docsrc/             # Documentation source (MkDocs)
│   ├── index.md
│   ├── specification/
│   ├── guides/
│   ├── prd/
│   └── releases/
├── docs/               # Generated documentation (GitHub Pages)
└── examples/           # Example changelogs
    ├── basic/
    ├── security/
    └── full/
```

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

## License

MIT License - see [LICENSE](LICENSE) for details.

 [build-status-svg]: https://github.com/grokify/structured-changelog/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/structured-changelog/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/structured-changelog/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/structured-changelog/actions/workflows/lint.yaml
 [coverage-svg]: https://img.shields.io/badge/coverage-98.1%25-brightgreen
 [coverage-url]: https://github.com/grokify/structured-changelog
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/structured-changelog
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/structured-changelog
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/structured-changelog
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/structured-changelog
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fstructured-changelog
 [loc-svg]: https://tokei.rs/b1/github/grokify/structured-changelog
 [repo-url]: https://github.com/grokify/structured-changelog
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/structured-changelog/blob/master/LICENSE
