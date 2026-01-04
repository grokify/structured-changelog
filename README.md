# Structured Changelog

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

[![Go Reference](https://pkg.go.dev/badge/github.com/grokify/structured-changelog.svg)](https://pkg.go.dev/github.com/grokify/structured-changelog)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

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
3. **Keep a Changelog format** — Compatible with [keepachangelog.com](https://keepachangelog.com/)
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

Install the CLI:

```bash
go install github.com/grokify/structured-changelog/cmd/sclog@latest
```

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

## JSON IR Schema

The JSON IR supports all [Keep a Changelog](https://keepachangelog.com/) categories:

| Category     | Description                                      |
|--------------|--------------------------------------------------|
| `added`      | New features                                     |
| `changed`    | Changes in existing functionality                |
| `deprecated` | Features that will be removed in future releases |
| `removed`    | Features removed in this release                 |
| `fixed`      | Bug fixes                                        |
| `security`   | Security vulnerability fixes                     |

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

See [docs/release-notes-guide.md](docs/release-notes-guide.md) for recommended release notes structure.

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
├── renderer/           # Deterministic Markdown renderer
│   ├── markdown.go
│   └── options.go
├── cmd/sclog/          # CLI tool (Cobra-based)
│   ├── main.go
│   ├── root.go
│   ├── validate.go
│   └── generate.go
├── schema/             # JSON Schema definitions
│   └── changelog-v1.schema.json
├── docs/               # Specification documentation
│   ├── spec.md
│   ├── security.md
│   └── sbom.md
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
