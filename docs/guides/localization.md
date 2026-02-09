# Localization Guide

structured-changelog supports localized output in 6 languages. This guide explains how to use built-in translations and create custom overrides.

## Built-in Languages

| Code | Language | Example |
|------|----------|---------|
| `en` | English (default) | "Changelog", "Added", "Breaking" |
| `de` | German | "Änderungsprotokoll", "Hinzugefügt", "Breaking Changes" |
| `es` | Spanish | "Registro de cambios", "Añadido", "Cambios importantes" |
| `fr` | French | "Journal des modifications", "Ajouté", "Changements majeurs" |
| `ja` | Japanese | "変更履歴", "追加", "破壊的変更" |
| `zh` | Chinese | "变更日志", "新增", "破坏性变更" |

## CLI Usage

```bash
# Generate in French
schangelog generate CHANGELOG.json --locale=fr -o CHANGELOG.md

# Generate in Japanese
schangelog generate CHANGELOG.json --locale=ja -o CHANGELOG.md

# Use custom overrides
schangelog generate CHANGELOG.json --locale=fr --locale-file=./my-fr.json -o CHANGELOG.md
```

## Library Usage

```go
import "github.com/grokify/structured-changelog/renderer"

// Use built-in French translations
opts := renderer.DefaultOptions().WithLocale("fr")
md, err := renderer.RenderMarkdown(changelog, opts)

// Use custom overrides
opts := renderer.DefaultOptions().
    WithLocale("fr").
    WithLocaleOverrides("./my-fr.json")
```

## Customizable Message IDs

Override any of these message IDs in your custom locale file:

### Document Structure

| ID | English Default | Description |
|----|-----------------|-------------|
| `changelog.title` | "Changelog" | Document title |
| `changelog.intro` | "All notable changes..." | Intro paragraph |
| `section.unreleased` | "Unreleased" | Unreleased section header |
| `section.yanked` | "YANKED" | Yanked version marker |

### Markers

| ID | English Default | Description |
|----|-----------------|-------------|
| `marker.breaking` | "BREAKING:" | Breaking change prefix |
| `marker.maintenance` | "Maintenance" | Maintenance release label |
| `marker.versions_range` | "Versions {{.From}} - {{.To}}" | Grouped versions header |

### Category Headers

| ID | English Default |
|----|-----------------|
| `category.highlights` | "Highlights" |
| `category.breaking` | "Breaking" |
| `category.upgrade_guide` | "Upgrade Guide" |
| `category.security` | "Security" |
| `category.added` | "Added" |
| `category.changed` | "Changed" |
| `category.deprecated` | "Deprecated" |
| `category.removed` | "Removed" |
| `category.fixed` | "Fixed" |
| `category.performance` | "Performance" |
| `category.dependencies` | "Dependencies" |
| `category.documentation` | "Documentation" |
| `category.build` | "Build" |
| `category.tests` | "Tests" |
| `category.infrastructure` | "Infrastructure" |
| `category.observability` | "Observability" |
| `category.compliance` | "Compliance" |
| `category.internal` | "Internal" |
| `category.known_issues` | "Known Issues" |
| `category.contributors` | "Contributors" |

### Plural Forms (for compact maintenance summaries)

| ID | English Default |
|----|-----------------|
| `plural.dependency_updates` | `{"one": "{{.Count}} dependency update", "other": "{{.Count}} dependency updates"}` |
| `plural.documentation_changes` | `{"one": "{{.Count}} documentation change", "other": "{{.Count}} documentation changes"}` |
| `plural.build_changes` | `{"one": "{{.Count}} build change", "other": "{{.Count}} build changes"}` |
| `plural.test_changes` | `{"one": "{{.Count}} test change", "other": "{{.Count}} test changes"}` |
| `plural.other_changes` | `{"one": "{{.Count}} other change", "other": "{{.Count}} other changes"}` |
| `plural.releases` | `{"one": "{{.Count}} release", "other": "{{.Count}} releases"}` |

## Custom Override File Format

Create a JSON file with only the messages you want to override:

```json
{
  "messages": [
    {"id": "changelog.title", "translation": "Release Notes"},
    {"id": "category.added", "translation": "New Features"},
    {"id": "category.fixed", "translation": "Bug Fixes"},
    {"id": "marker.breaking", "translation": "⚠️ BREAKING:"}
  ]
}
```

Messages not included in your override file will use the built-in translations for the selected locale.

## Examples

See `examples/l10n/` for complete rendered examples in all 6 languages:

- [English](../examples/l10n/CHANGELOG.en.md)
- [German](../examples/l10n/CHANGELOG.de.md)
- [Spanish](../examples/l10n/CHANGELOG.es.md)
- [French](../examples/l10n/CHANGELOG.fr.md)
- [Japanese](../examples/l10n/CHANGELOG.ja.md)
- [Chinese](../examples/l10n/CHANGELOG.zh.md)

## What Gets Localized

**Localized (system labels):**
- Section headers ("Changelog", "Unreleased", "Added", etc.)
- Breaking change markers
- Maintenance release summaries

**Not localized (user content):**
- Entry descriptions (your changelog text)
- Version numbers
- Dates (always ISO 8601: YYYY-MM-DD)
- Project names
- Issue/PR references

This preserves the JSON IR as the canonical, language-agnostic source of truth.
