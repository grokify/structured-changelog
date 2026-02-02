# PRD: Localization Support for Structured Changelog

## Overview

Add localization (L10N) support for system labels in structured-changelog, enabling changelog generation with localized section headers and UI labels while preserving user-supplied content unchanged.

## Scope

### In Scope

Localization of system-generated labels only:

1. **Change Type Names** (20 categories)
   - Core: Added, Changed, Deprecated, Removed, Fixed, Security
   - Standard: Highlights, Breaking, Upgrade Guide, Performance, Dependencies
   - Extended: Documentation, Build, Tests, Contributors
   - Optional: Infrastructure, Observability, Compliance, Internal, Known Issues

2. **Change Type Subtitles** (20 items)
   - Brief descriptions shown in help text (e.g., "for new features")

3. **Tier Labels** (4 items)
   - Core, Standard, Extended, Optional

4. **Tier Descriptions** (4 items)
   - Explanatory text for each tier

5. **Markdown Output Labels** (~10 items)
   - "Changelog" (document title)
   - "Unreleased"
   - "[YANKED]"
   - "BREAKING:" prefix
   - Date format localization (optional)

6. **CLI Labels** (future consideration)
   - Help text and error messages

### Out of Scope

- **CHANGELOG.json content**: The JSON IR remains English-only
  - JSON field names (e.g., `added`, `fixed`, `security`)
  - Change type identifiers in the registry
  - User-supplied content (entry descriptions, project names, version numbers)
- **Validation error messages**: Keep in English for debuggability
- **Right-to-left (RTL) rendering**: Display-layer concern, not generation

### Design Principle

```
CHANGELOG.json (English, canonical) → Renderer (locale-aware) → CHANGELOG.md (localized labels)
```

The JSON Intermediate Representation is the **source of truth** and remains language-agnostic. Localization applies only at the **rendering stage** when generating Markdown output. This preserves:

1. **Schema stability**: JSON structure unchanged across locales
2. **Tooling interoperability**: Parsers work regardless of output locale
3. **Round-trip safety**: JSON → Markdown → JSON conversion unaffected by locale

## Target Locales

Based on Keep a Changelog community translations, prioritized by developer population:

### Tier 1 (High Priority)

| Code | Language | Region |
|------|----------|--------|
| en | English | Default |
| zh-CN | Chinese | Simplified |
| zh-TW | Chinese | Traditional |
| ja | Japanese | Japan |
| es-ES | Spanish | Spain |
| fr | French | France |
| de | German | Germany |
| pt-BR | Portuguese | Brazil |

### Tier 2 (Medium Priority)

| Code | Language | Region |
|------|----------|--------|
| ru | Russian | Russia |
| ko | Korean | Korea |
| it | Italian | Italy |
| nl | Dutch | Netherlands |
| pl | Polish | Poland |
| tr | Turkish | Turkey |

### Tier 3 (Community Contributed)

| Code | Language | Region |
|------|----------|--------|
| ar | Arabic | — |
| cs | Czech | Czech Republic |
| da | Danish | Denmark |
| fa | Persian | Iran |
| hr | Croatian | Croatia |
| id-ID | Indonesian | Indonesia |
| ka | Georgian | Georgia |
| nb | Norwegian | Bokmål |
| sk | Slovak | Slovakia |
| sv | Swedish | Sweden |
| uk | Ukrainian | Ukraine |

## Technical Design

### Library Choice

Use `github.com/nicksnyder/go-i18n/v2` for:

- Mature, production-proven library
- JSON/YAML/TOML message file support
- Proper plural rule handling
- Used by Docker, HashiCorp, and others

### Directory Structure

```
locales/
  en.json          # English (canonical, complete)
  zh-CN.json       # Chinese Simplified
  zh-TW.json       # Chinese Traditional
  ja.json          # Japanese
  fr.json          # French
  ...
```

### Message ID Schema

Stable, semantic identifiers that remain constant across all locales:

```
# Change type names
changetype.added.name
changetype.changed.name
changetype.deprecated.name
changetype.removed.name
changetype.fixed.name
changetype.security.name
changetype.highlights.name
changetype.breaking.name
changetype.upgrade_guide.name
changetype.performance.name
changetype.dependencies.name
changetype.documentation.name
changetype.build.name
changetype.tests.name
changetype.infrastructure.name
changetype.observability.name
changetype.compliance.name
changetype.internal.name
changetype.known_issues.name
changetype.contributors.name

# Change type subtitles
changetype.added.subtitle
changetype.changed.subtitle
... (20 total)

# Tier labels
tier.core.name
tier.standard.name
tier.extended.name
tier.optional.name

# Tier descriptions
tier.core.description
tier.standard.description
tier.extended.description
tier.optional.description

# Markdown labels
markdown.title
markdown.unreleased
markdown.yanked
markdown.breaking_prefix
```

### Message File Format

Using go-i18n JSON format:

```json
[
  {
    "id": "changetype.added.name",
    "translation": "Added"
  },
  {
    "id": "changetype.added.subtitle",
    "translation": "for new features"
  },
  {
    "id": "tier.core.name",
    "translation": "Core"
  },
  {
    "id": "markdown.title",
    "translation": "Changelog"
  }
]
```

### API Changes

#### Renderer Options

```go
type RenderOptions struct {
    Locale string // BCP 47 tag: "en", "fr", "zh-CN"
    // ... existing options
}
```

#### CLI Flag

```bash
schangelog generate --locale fr changelog.json
schangelog generate -l zh-CN changelog.json
```

#### Default Behavior

- Default locale: `en`
- Fallback chain: `zh-CN` → `zh` → `en`
- Missing translation: Fall back to English, emit warning

### Determinism Guarantee

Localized output must remain deterministic:

- Same locale + same input = identical output
- Locale affects only label strings, not structure or ordering
- No runtime locale detection; explicit locale required

## Implementation Phases

### Phase 1: Infrastructure

1. Add `go-i18n` dependency
2. Create `locales/` directory structure
3. Create canonical `en.json` with all message IDs
4. Add locale loading to `ChangeTypeRegistry`
5. Add `--locale` flag to CLI

### Phase 2: Core Locales

1. Implement localized rendering in `renderer/markdown.go`
2. Add translations for Tier 1 locales (8 languages)
3. Add locale parity tests

### Phase 3: Extended Locales

1. Add Tier 2 locale translations
2. Create contributor documentation for translators
3. Add translation coverage badge to README

### Phase 4: Community

1. Add Tier 3 locales via community contribution
2. Implement translation parity CI check
3. Document translation workflow in CONTRIBUTING.md

## Testing Strategy

### Unit Tests

- Verify all message IDs exist in English
- Verify locale fallback behavior
- Verify deterministic output per locale

### Parity Tests

```go
func TestTranslationParity(t *testing.T) {
    // All locales must have same message IDs as English
    englishIDs := getMessageIDs("en")
    for _, locale := range supportedLocales {
        localeIDs := getMessageIDs(locale)
        assert.Equal(t, englishIDs, localeIDs)
    }
}
```

### Golden File Tests

- One golden output file per locale
- Ensures rendering changes are intentional

## Success Metrics

- 100% message ID coverage for Tier 1 locales
- Zero fallback warnings for Tier 1/2 locales
- Translation parity enforced in CI

## Design Decisions

### Embedded Locale Files (`//go:embed`)

Locale files will be **embedded** in the binary rather than loaded at runtime.

**Rationale:**

| Factor | Embedded | Runtime |
|--------|----------|---------|
| Distribution | Single binary | Binary + files |
| Reliability | Guaranteed | File-not-found risk |
| Portability | Works anywhere | Path resolution issues |
| Version consistency | Translations match code | Potential skew |

**Additional considerations:**

- `change_types.json` is already embedded — consistent approach
- CLI tools benefit from single-binary deployment
- Translation updates ship with releases (acceptable for developer tooling)
- Community contributions via PR ensure review and quality control

### Date Format

Dates remain **ISO 8601** (`YYYY-MM-DD`) regardless of locale.

**Rationale:**

- Changelog dates are for version history, not human scheduling
- ISO 8601 is unambiguous and sort-friendly
- Keep a Changelog specifies this format

### Plural Forms

**Not required** — all labels are singular or uncountable (e.g., "Added", "Security", "Changelog").

### CLI Help Text

CLI help text and error messages remain **English-only**.

**Rationale:**

- **CLI users are developers** — assumed to have English proficiency
- **Changelog readers are end users** — may prefer native language
- Localization effort focused on user-facing output, not developer tooling

## References

- [go-i18n Documentation](https://github.com/nicksnyder/go-i18n)
- [BCP 47 Language Tags](https://www.rfc-editor.org/info/bcp47)
- [Keep a Changelog Translations](https://keepachangelog.com/)
- [golang.org/x/text/language](https://pkg.go.dev/golang.org/x/text/language)
