# Structured Changelog Specification v1.0

This document defines the Structured Changelog Intermediate Representation (IR) format.

## Overview

The Structured Changelog IR is a JSON format that serves as the canonical source of truth for changelogs. It is designed to be:

- **Machine-readable**: Easily parsed and processed by tools
- **Deterministic**: Same input always produces identical Markdown output
- **Extensible**: Supports optional metadata for security and SBOM
- **Compatible**: Follows Keep a Changelog and Semantic Versioning conventions

## JSON IR Structure

### Root Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `ir_version` | string | Yes | IR schema version (currently "1.0") |
| `project` | string | Yes | Project name |
| `repository` | string | No | Repository URL |
| `versioning` | string | No | Versioning scheme (see below) |
| `commit_convention` | string | No | Commit message convention (see below) |
| `maintainers` | string[] | No | Team members excluded from author attribution |
| `bots` | string[] | No | Custom bots excluded from author attribution |
| `generated_at` | datetime | No | ISO 8601 timestamp of generation |
| `unreleased` | Release | No | Unreleased changes |
| `releases` | Release[] | No | Array of releases (reverse chronological) |

#### Versioning Schemes

The `versioning` field controls what versioning reference appears in the generated header:

| Value | Description | Header Text |
|-------|-------------|-------------|
| `semver` | Semantic Versioning (default) | "adheres to Semantic Versioning" |
| `calver` | Calendar Versioning | "uses Calendar Versioning" |
| `custom` | Custom versioning | No versioning line |
| `none` | No versioning scheme | No versioning line |

#### Commit Conventions

The `commit_convention` field adds a reference to the commit message convention:

| Value | Description | Header Text |
|-------|-------------|-------------|
| `conventional` | Conventional Commits | "commits follow Conventional Commits" |
| `none` | No convention (default) | No convention line |

### Release Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | string | Yes* | Semantic version string |
| `date` | string | Yes* | Release date (YYYY-MM-DD) |
| `yanked` | boolean | No | Whether the release was retracted |
| `compare_url` | string | No | URL to diff with previous version |
| `added` | Entry[] | No | New features |
| `changed` | Entry[] | No | Changes to existing features |
| `deprecated` | Entry[] | No | Features to be removed |
| `removed` | Entry[] | No | Removed features |
| `fixed` | Entry[] | No | Bug fixes |
| `security` | Entry[] | No | Security fixes |

*Required for releases, not for unreleased section.

### Entry Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `description` | string | Yes | Description of the change |
| `issue` | string | No | Issue reference (number or URL) |
| `pr` | string | No | Pull request reference (number or URL) |
| `commit` | string | No | Commit SHA (full or short) |
| `author` | string | No | Author of the change |
| `breaking` | boolean | No | Breaking change flag |

#### Reference Linking

When a `repository` URL is provided and the renderer is configured with `LinkReferences: true` (enabled in `FullOptions`), references are automatically linked:

| Reference | GitHub URL | GitLab URL |
|-----------|------------|------------|
| `issue: "42"` | `/issues/42` | `/-/issues/42` |
| `pr: "43"` | `/pull/43` | `/-/merge_requests/43` |
| `commit: "abc123..."` | `/commit/abc123...` | `/-/commit/abc123...` |

Commits are displayed as short hashes (7 characters) in the output, but the full SHA is used in the link URL.

#### Author Attribution

When an entry includes an `author` field and the renderer is configured with `IncludeAuthors: true` (enabled by default), external contributors are automatically attributed:

```markdown
- New feature by [@contributor](https://github.com/contributor)
```

Authors listed in the root `maintainers` array are excluded from attribution. Common bots (dependabot, renovate, github-actions, etc.) are auto-detected and excluded. Custom bots can be specified in the root `bots` array.

If the description already contains inline attribution matching the author (e.g., "Feature from @user"), the inline attribution is automatically stripped to avoid duplication.

### Security Metadata (Optional)

Security entries may include additional fields:

| Field | Type | Description |
|-------|------|-------------|
| `cve` | string | CVE identifier (e.g., CVE-2026-12345) |
| `ghsa` | string | GitHub Security Advisory ID |
| `severity` | string | critical, high, medium, low, informational |
| `cvss_score` | number | CVSS score (0.0-10.0) |
| `cvss_vector` | string | CVSS vector string |
| `cwe` | string | CWE identifier (e.g., CWE-89) |
| `affected_versions` | string | Version range affected |
| `patched_versions` | string | Version range with fix |
| `sarif_rule_id` | string | SARIF rule ID for linking |

### SBOM Metadata (Optional)

Entries may include SBOM (Software Bill of Materials) fields:

| Field | Type | Description |
|-------|------|-------------|
| `component` | string | Component/dependency name |
| `component_version` | string | Component version |
| `license` | string | SPDX license identifier |

## Breaking Changes

The `breaking` field marks entries that introduce breaking changes. When rendered with `MarkBreakingChanges: true`, these entries are prefixed with `**BREAKING:**`.

### Usage

Use `"breaking": true` on entries in the `changed` or `removed` categories:

```json
{
  "changed": [
    { "description": "Rename DoThing() to DoThingV2()", "breaking": true },
    { "description": "Improve performance of DoOther()" }
  ],
  "removed": [
    { "description": "Remove deprecated LegacyAPI", "breaking": true }
  ]
}
```

### Rendered Output

```markdown
### Changed

- **BREAKING:** Rename DoThing() to DoThingV2()
- Improve performance of DoOther()

### Removed

- **BREAKING:** Remove deprecated LegacyAPI
```

### Migration Guides

Breaking changes in CHANGELOG.json should be concise. Detailed migration guides with code examples belong in `RELEASE_NOTES_vX.Y.Z.md`.

| Content | CHANGELOG.json | RELEASE_NOTES |
|---------|---------------|---------------|
| Breaking flag | ✓ `"breaking": true` | Reference only |
| What changed | ✓ Brief description | ✓ Full context |
| Code examples | ✗ Never | ✓ Before/after |
| Migration steps | ✗ Never | ✓ Step-by-step |

See the [Release Notes Guide](../guides/release-notes-guide.md) for recommended release notes structure.

## Category Order

Categories MUST be rendered in this order:

1. Added
2. Changed
3. Deprecated
4. Removed
5. Fixed
6. Security

This matches the Keep a Changelog specification.

## Version Ordering

Releases MUST be ordered in reverse chronological order (newest first).

## Deterministic Rendering

The renderer MUST produce identical output for identical input. This means:

- Category order is fixed (see above)
- Entry order within categories matches array order
- No randomization or timestamp-based formatting
- Consistent whitespace and newlines

## Validation Rules

1. `ir_version` must be "1.0"
2. `project` must be non-empty
3. Release `version` must be valid semver
4. Release `date` must be YYYY-MM-DD format
5. Entry `description` must be non-empty
6. `cve` must match pattern `CVE-\d{4}-\d{4,}`
7. `ghsa` must match pattern `GHSA-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}`
8. `severity` must be one of: critical, high, medium, low, informational
9. `cvss_score` must be between 0 and 10
10. No duplicate versions allowed

## Example

```json
{
  "ir_version": "1.0",
  "project": "my-project",
  "releases": [
    {
      "version": "1.0.0",
      "date": "2026-01-03",
      "added": [
        { "description": "Initial release" }
      ]
    }
  ]
}
```

## References

- [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
- [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- [SARIF Specification](https://sarifweb.azurewebsites.net/)
- [SPDX License List](https://spdx.org/licenses/)
