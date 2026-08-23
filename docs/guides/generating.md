# Generating Changelogs

This guide covers `schangelog generate`, which renders a human-readable `CHANGELOG.md` from the canonical `CHANGELOG.json`. Output is **deterministic** — the same input always produces identical output — and highly configurable through presets, change-type tiers, release filtering, and per-category controls.

## Basic Usage

```bash
# Write to stdout
schangelog generate CHANGELOG.json

# Write to a file
schangelog generate CHANGELOG.json -o CHANGELOG.md
```

`CHANGELOG.json` is always the source of truth; the generated Markdown is derived and should not be hand-edited.

## Presets

A preset selects a coherent bundle of rendering options. Pick one with `--minimal` / `--full`, or use the default.

| Preset | Flag | What it produces |
|--------|------|------------------|
| **default** | *(none)* | Notable releases only, commit/reference links, Dependencies collapsed, maintenance releases compacted |
| **minimal** | `--minimal` | Core-tier types only, no references/metadata/commit links |
| **full** | `--full` | Every release and every category expanded, all metadata, nothing collapsed or compacted |

Two additional presets are available programmatically (`core`, `standard`) via the library; see [Programmatic Rendering](#programmatic-rendering).

## Change-Type Tiers

Every change type belongs to a tier. `--max-tier` includes only types at or above the given tier, letting you scale detail up or down.

| Tier | Includes |
|------|----------|
| **core** | Security, Added, Changed, Deprecated, Removed, Fixed (the Keep a Changelog set) |
| **standard** | core + Highlights, Breaking, Upgrade Guide, Performance, Dependencies |
| **extended** | standard + Documentation, Build, Tests, Known Issues, Contributors |
| **optional** | extended + Infrastructure, Observability, Compliance, Internal |

```bash
# KACL-compliant core output only
schangelog generate CHANGELOG.json --max-tier core

# Everything up to extended
schangelog generate CHANGELOG.json --max-tier extended
```

## Notable Releases (Default)

By default, only **notable** releases are rendered. Maintenance-only releases (those with just dependencies, documentation, build, tests, etc.) are excluded to keep the changelog focused on user-facing change.

```bash
# Default: notable releases only
schangelog generate CHANGELOG.json -o CHANGELOG.md

# Include every release
schangelog generate CHANGELOG.json --all-releases

# Redefine what counts as notable
schangelog generate CHANGELOG.json --notable-categories "Security,Added,Fixed"
```

- **Notable categories:** Highlights, Breaking, Upgrade Guide, Security, Added, Changed, Deprecated, Removed, Fixed, Performance, Known Issues
- **Non-notable (maintenance):** Dependencies, Documentation, Build, Tests, Infrastructure, Observability, Compliance, Internal, Contributors

## Compact Maintenance Releases

With `--all-releases`, consecutive maintenance-only releases are grouped into a single compact section instead of listing each one:

```markdown
## Versions 0.71.1 - 0.71.10 (Maintenance)

10 releases: 8 dependency updates, 2 documentation changes.
```

Use `--full` to include every release fully expanded (disables both notable-only filtering and grouping).

## Category Filtering

Within a rendered release, individual categories can be **collapsed** to a one-line summary or **excluded** entirely. This keeps the human-facing Markdown readable while the full detail remains in the JSON source of truth.

### Dependencies Are Collapsed by Default

Dependency churn is the dominant source of changelog noise for readers, so the `default` and `standard` presets render the Dependencies category as a one-line summary:

```markdown
### Dependencies

- _17 dependency updates_
```

!!! note "Behavior since v0.16.0"
    Collapsing Dependencies by default is a change to generated output. Projects that carry `dependencies` entries will see those sections collapse on the next `generate`. The full list is unchanged in `CHANGELOG.json` and is one flag away. The `core` and `minimal` presets are unaffected — they already drop Dependencies via tier filtering.

### Controls

```bash
# Restore the full dependency list (either works)
schangelog generate CHANGELOG.json --expand-categories Dependencies
schangelog generate CHANGELOG.json --full

# Collapse additional noisy categories to a one-line count
schangelog generate CHANGELOG.json --collapse-categories "Build,Tests"

# Omit categories entirely (no section at all)
schangelog generate CHANGELOG.json --exclude-categories "Dependencies,Build"
```

| Flag | Effect |
|------|--------|
| `--collapse-categories` | Render the named categories as a one-line entry count |
| `--exclude-categories` | Omit the named categories entirely (no section) |
| `--expand-categories` | Force-expand the named categories, overriding collapse defaults |

**Precedence:** `--exclude-categories` wins over `--collapse-categories` for the same category (the section is omitted, not summarized). `--expand-categories` is applied last, so it overrides both a preset's collapse default and an explicit `--collapse-categories`.

## Reference Linking

For GitHub and GitLab repositories (via the changelog's `repository` field), commit SHAs, issues, and PRs are rendered as links, and version comparison links are appended at the bottom. Commit links are included by default and suppressed by `--minimal`.

## Localized Output

Rendered section headings and summary phrases can be localized:

```bash
# Built-in locales: en, fr, de, es, ja, zh
schangelog generate CHANGELOG.json --locale=fr -o CHANGELOG.md

# Override specific messages with a custom locale file
schangelog generate CHANGELOG.json --locale=fr --locale-file=./custom-fr.json
```

See the [Localization guide](localization.md) for details.

## Programmatic Rendering

The same controls are available on `renderer.Options` when generating from Go:

```go
import (
    "github.com/grokify/structured-changelog/changelog"
    "github.com/grokify/structured-changelog/renderer"
)

cl, err := changelog.LoadFile("CHANGELOG.json")
if err != nil {
    return err
}

opts := renderer.DefaultOptions().
    WithCollapseCategories("Build").   // one-line summaries
    WithExcludeCategories("Tests").    // omit entirely
    WithExpandCategories("Dependencies") // override the collapse default

md := renderer.RenderMarkdownWithOptions(cl, opts)
```

Preset constructors (`DefaultOptions`, `MinimalOptions`, `FullOptions`, `CoreOptions`, `StandardOptions`) and `OptionsFromConfig` mirror the CLI presets and flags.

## Determinism

Generation is deterministic: identical `CHANGELOG.json` plus identical options always yields byte-identical Markdown. This makes the output safe to commit and to diff in CI — regenerate and fail on drift if the committed `CHANGELOG.md` no longer matches its source.

## Links

- [LLM-Assisted Generation](llm-guide.md)
- [Release Notes Guide](release-notes-guide.md)
- [Localization](localization.md)
- [Full Changelog](../changelog.md)
