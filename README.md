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

### Comparison

#### Keep a Changelog (the specification)

[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) is the de-facto standard for human-readable changelogs. Structured Changelog implements this specification with a machine-readable JSON layer:

> **Keep a Changelog is the human spec; Structured Changelog is the structured implementation.**

The generated `CHANGELOG.md` conforms to Keep a Changelog 1.1.0 formatting conventions, while adding JSON schema for machine parsing, security fields (CVE/GHSA/CVSS), and SBOM metadata.

#### Changelog Generation Tools

Several tools generate changelogs from git history. Here's how Structured Changelog differs:

| Feature | Structured Changelog | [conventional-changelog] | [git-cliff] | [chyle] | [semverbot] |
|---------|---------------------|-------------------------|-------------|---------|-------------|
| **GitHub stars** | — | ~8.4k | ~5.5k | ~160 | ~144 |
| **Workflow** | LLM/Human-assisted | Fully automated | Fully automated | Fully automated | Fully automated |
| **Source of truth** | JSON IR | Git commits | Git + templates | Git commits | Git tags |
| **Output format** | JSON → Markdown | Markdown | Markdown | Flexible | Tags only |
| **Rendering** | Deterministic | Template-based | Template-based | Configurable | N/A |
| **Machine-readable** | ✓ (JSON IR) | ✗ | ✗ | ✗ | ✗ |
| **I18N/Localization** | ✓ (6 languages) | ✗ | ✗ | ✗ | ✗ |
| **Security metadata** | ✓ (CVE/GHSA/CVSS) | ✗ | ✗ | ✗ | ✗ |
| **SBOM metadata** | ✓ | ✗ | ✗ | ✗ | ✗ |
| **LLM optimization** | ✓ (TOON format) | ✗ | ✗ | ✗ | ✗ |
| **Version bumping** | ✗ | ✓ | ✗ | ✗ | ✓ |
| **Conventional commits** | ✓ | ✓ | ✓ | ✓ | ✓ |
| **Custom templates** | ✗ (deterministic) | ✓ | ✓ | ✓ | ✗ |
| **Language** | Go | Node.js | Rust | Go | Go |

#### LLM-Assisted vs. Fully Automated

The key architectural difference is *how* the changelog is populated:

```
Traditional tools (git-cliff, conventional-changelog):
  Git Commits → Template/Parser → Markdown
  └─ Runs in CI/CD, fully automated, pattern-based

Structured Changelog:
  Git Commits → LLM/Human → JSON IR → Markdown
  └─ Semantic understanding, judgment calls, better prose
```

**Why LLM-assisted?** Pure automation works well for simple cases, but changelogs benefit from intelligence:

- **Grouping**: Consolidate 5 related commits into one meaningful entry
- **Context**: Explain *why* something changed, not just *what*
- **Judgment**: Categorize ambiguous changes correctly
- **Quality**: Write user-friendly descriptions, not commit messages
- **Edge cases**: Handle non-conventional commits gracefully

The `parse-commits` command outputs token-optimized TOON format (~8x reduction vs raw git log), making it practical to feed git history to an LLM for changelog generation.

**When to use Structured Changelog:**

- You want LLM-assisted changelog generation with human-quality prose
- You need a machine-readable changelog for automation or APIs
- You want deterministic output (same input → identical output)
- You track security vulnerabilities with CVE/GHSA identifiers
- You need SBOM integration for compliance

**When to use other tools:**

- **conventional-changelog**: You're in a Node.js ecosystem and want automatic version bumping
- **git-cliff**: You want maximum template customization with Rust performance
- **chyle**: You need to enrich changelog data from external APIs (Jira, GitHub)
- **semverbot**: You primarily need automated semantic version tagging

> **Note:** [git-chglog](https://github.com/git-chglog/git-chglog) was previously the most popular Go-based changelog generator (~2.9k stars), but the project has been archived. Its maintainers recommend [git-cliff] as the successor.

[conventional-changelog]: https://github.com/conventional-changelog/conventional-changelog
[git-cliff]: https://github.com/orhun/git-cliff
[chyle]: https://github.com/antham/chyle
[semverbot]: https://github.com/restechnica/semverbot

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap grokify/tap
brew install structured-changelog
```

This installs the `schangelog` CLI (also available as `structured-changelog`).

### Go Install

```bash
go install github.com/grokify/structured-changelog/cmd/schangelog@latest
```

### Go Library

```bash
go get github.com/grokify/structured-changelog
```

## Quick Start

### Define your changelog in JSON

```json
{
  "irVersion": "1.0",
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
schangelog validate CHANGELOG.json
```

Generate Markdown:

```bash
# Output to stdout
schangelog generate CHANGELOG.json

# Output to file
schangelog generate CHANGELOG.json -o CHANGELOG.md

# Full output (default: includes commit links and reference linking)
schangelog generate CHANGELOG.json

# Minimal output (no references/metadata/commit links)
schangelog generate CHANGELOG.json --minimal
```

Show version:

```bash
schangelog version
```

### Localized Output (I18N)

Generate changelogs in multiple languages:

```bash
# Generate in French
schangelog generate CHANGELOG.json --locale=fr -o CHANGELOG.md

# Generate in Japanese
schangelog generate CHANGELOG.json --locale=ja -o CHANGELOG.md

# Use custom translation overrides
schangelog generate CHANGELOG.json --locale=fr --locale-file=./custom-fr.json
```

Built-in locales: English (`en`), German (`de`), Spanish (`es`), French (`fr`), Japanese (`ja`), Chinese (`zh`).

See the [Localization Guide](https://grokify.github.io/structured-changelog/guides/localization/) for customization options.

### LLM-Assisted Generation

The CLI includes tools optimized for AI-assisted changelog generation. Commands default to TOON format (Token-Oriented Object Notation) for ~8x token reduction:

```bash
schangelog parse-commits --since=v0.3.0           # Structured git history (TOON)
schangelog parse-commits --since=v0.3.0 --format=json  # JSON output
schangelog parse-commits --since=v0.3.0 --changelog=CHANGELOG.json  # Mark external contributors
schangelog suggest-category "feat: ..."           # Category suggestions
schangelog validate --format=toon CHANGELOG.json  # Rich error output
```

See the [LLM Guide](https://grokify.github.io/structured-changelog/guides/llm-guide/) for prompts and workflows.

### Reference Linking

When a repository URL is provided, references (issues, PRs, commits) are automatically linked by default:

```bash
# Generate with linked references (default behavior)
schangelog generate CHANGELOG.json
```

Example output:

```markdown
- Add OAuth2 support ([#42](https://github.com/example/repo/issues/42), [`abc123d`](https://github.com/example/repo/commit/abc123def))
```

Supports GitHub and GitLab URL formats.

### Author Attribution

External contributors are automatically attributed when an `author` field is set and `maintainers` are defined:

```json
{
  "maintainers": ["grokify"],
  "releases": [{
    "added": [{ "description": "New feature", "author": "@contributor" }]
  }]
}
```

Generates: `- New feature by [@contributor](https://github.com/contributor)`

Common bots (dependabot, renovate, etc.) are auto-detected and excluded from attribution.

### Compact Maintenance Releases

Consecutive maintenance-only releases (dependencies, documentation, build) are automatically grouped:

```markdown
## Versions 0.71.1 - 0.71.10 (Maintenance)

10 releases: 8 dependency update(s), 2 documentation change(s).
```

Use `--full` to disable grouping and show all releases expanded.

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
schangelog validate --min-tier core

# Validate: require core + standard coverage
schangelog validate --min-tier standard

# Generate: output only core types (KACL-compliant)
schangelog generate --max-tier core

# Generate: include everything up to extended
schangelog generate --max-tier extended
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

See [docs/guides/release-notes-guide.md](docs/guides/release-notes-guide.md) for recommended release notes structure.

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
│   ├── parser.go
│   └── tags.go
├── renderer/           # Deterministic Markdown renderer
│   ├── markdown.go
│   └── options.go
├── cmd/schangelog/     # CLI tool (Cobra-based)
│   ├── main.go
│   ├── root.go
│   ├── validate.go
│   ├── generate.go
│   ├── parse_commits.go
│   ├── suggest_category.go
│   ├── list_tags.go
│   └── init.go
├── schema/             # JSON Schema definitions
│   └── changelog-v1.schema.json
├── docs/               # Documentation source (MkDocs)
│   ├── index.md
│   ├── changelog.md
│   ├── specification/
│   ├── guides/
│   ├── prd/
│   └── releases/
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
 [coverage-svg]: https://img.shields.io/badge/coverage-96.1%25-brightgreen
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
