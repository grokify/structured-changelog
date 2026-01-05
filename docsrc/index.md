# Structured Changelog

**A canonical, deterministic changelog framework with JSON IR and Markdown rendering.**

Structured Changelog provides a machine-readable JSON Intermediate Representation (IR) as the source of truth for your changelog, with deterministic Markdown generation for human readability. It supports optional security metadata (CVE, GHSA, SARIF) and SBOM information.

## Overview

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  CHANGELOG.json │ ───► │    Renderer     │ ───► │  CHANGELOG.md   │
│  (Canonical IR) │      │ (Deterministic) │      │ (Human-Readable)│
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

## Key Principles

1. **JSON IR is canonical** — Markdown is derived, not the source of truth
2. **Deterministic rendering** — Same input always produces identical output
3. **Keep a Changelog format** — Compatible with [keepachangelog.com](https://keepachangelog.com/)
4. **Semantic Versioning** — Follows [semver.org](https://semver.org/) conventions
5. **Extensible metadata** — Optional security (CVE/GHSA/SARIF) and SBOM fields
6. **Spec + tooling together** — Single source of truth for humans and machines

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

    "github.com/grokify/structured-changelog/changelog"
    "github.com/grokify/structured-changelog/renderer"
)

func main() {
    cl, err := changelog.LoadFile("CHANGELOG.json")
    if err != nil {
        panic(err)
    }
    md := renderer.RenderMarkdown(cl)
    fmt.Println(md)
}
```

### CLI Usage

```bash
# Validate a changelog
sclog validate CHANGELOG.json

# Generate Markdown
sclog generate CHANGELOG.json

# Output to file
sclog generate CHANGELOG.json -o CHANGELOG.md
```

## Documentation

- [JSON IR Specification](specification/spec.md) — Full schema documentation
- [Security Metadata](specification/security.md) — CVE, GHSA, SARIF fields
- [SBOM Metadata](specification/sbom.md) — Component tracking
- [Release Notes Guide](guides/release-notes-guide.md) — CHANGELOG vs RELEASE_NOTES

## Links

- [GitHub Repository](https://github.com/grokify/structured-changelog)
- [Go Package Documentation](https://pkg.go.dev/github.com/grokify/structured-changelog)
- [Keep a Changelog](https://keepachangelog.com/)
