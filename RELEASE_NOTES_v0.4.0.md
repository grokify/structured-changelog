# Release Notes - v0.4.0

## Overview

This release introduces a **tier-based change type system** that extends Keep a Changelog with 13 additional change types organized into 4 tiers. It also includes significant library improvements for better testability and reusability.

## Highlights

### Tier-Based Change Types

Change types are now organized into tiers for flexible filtering:

| Tier | Description | Categories |
|------|-------------|------------|
| **core** | KACL standard | Security, Added, Changed, Deprecated, Removed, Fixed |
| **standard** | Common extensions | Highlights, Breaking, Upgrade Guide, Performance, Dependencies |
| **extended** | Docs & acknowledgments | Documentation, Build, Known Issues, Contributors |
| **optional** | Operations & internal | Infrastructure, Observability, Compliance, Internal |

### CLI Improvements

- `sclog generate --max-tier <tier>` - Filter output to include only categories at or above the specified tier
- `sclog validate --min-tier <tier>` - Require at least one entry in a category at or above the specified tier

### Library Functions

New library functions for programmatic use:

- `changelog.Summary()` - Get a structured summary of changelog contents
- `changelog.ValidateMinTier(tier)` - Validate tier coverage
- `changelog.ParseTier(string)` - Parse tier strings (case-insensitive)
- `renderer.OptionsFromConfig(cfg)` - Create render options from configuration

## Breaking Changes

None. This release is fully backwards compatible.

## Test Coverage

Library packages now have 98%+ unit test coverage:

- `changelog`: 98.2%
- `renderer`: 98.9%

## Installation

```bash
# Via Go
go install github.com/grokify/structured-changelog/cmd/sclog@v0.4.0

# Via Homebrew
brew install grokify/tap/structured-changelog
```

## Links

- [Full Changelog](CHANGELOG.md)
- [Change Types Reference](CHANGE_TYPES.json)
- [Documentation](https://pkg.go.dev/github.com/grokify/structured-changelog)
