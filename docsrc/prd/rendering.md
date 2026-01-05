# PRD: Rendering

## Overview

Define the rendering behavior for converting `CHANGELOG.json` to human-readable output formats, with a focus on Markdown compatibility across different rendering engines.

## Markdown Output Targets

Support both **GitHub** and **Pandoc** rendering engines.

### Target Use Cases

| Engine | Output Formats | Use Case |
|--------|---------------|----------|
| GitHub | Web display | Repository CHANGELOG.md |
| Pandoc | DOCX, PDF, HTML, etc. | Documentation, releases |

### Approach

- Generate **generic Markdown** compatible with both engines
- Use only basic constructs: headings, lists, links, emphasis
- Avoid engine-specific features

### Feature Compatibility

| Feature | GitHub | Pandoc | Used |
|---------|--------|--------|------|
| ATX headings (`#`, `##`) | Yes | Yes | Yes |
| Unordered lists (`-`) | Yes | Yes | Yes |
| Inline links `[text](url)` | Yes | Yes | Yes |
| Bold (`**text**`) | Yes | Yes | Yes |
| Task lists (`- [ ]`) | Yes | No | No |
| Alerts (`> [!NOTE]`) | Yes | No | No |
| Fenced divs (`:::`) | No | Yes | No |

Current formatting uses only universally supported constructs.

### Flavor Support (Future)

If divergence becomes necessary, add a `--markdown-flavor` flag:

```bash
sclog generate --markdown-flavor github changelog.json
sclog generate --markdown-flavor pandoc changelog.json
```

**Current assessment:** Not needed — simple formatting works universally.

## Deterministic Output

Rendering must be **deterministic**: same input always produces identical output.

### Guarantees

1. **Stable ordering** — Sections appear in defined category order
2. **Consistent whitespace** — No trailing spaces, consistent blank lines
3. **No timestamps** — Output contains no generation timestamps
4. **Locale isolation** — Same locale + same input = same output

### Why This Matters

- Git diffs remain meaningful
- CI/CD pipelines produce reproducible artifacts
- Round-trip testing is reliable

## Output Structure

### Standard Markdown Layout

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- New feature description

## [1.0.0] - 2024-01-15

### Added

- Initial release feature

### Fixed

- Bug fix description

[Unreleased]: https://github.com/user/repo/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/user/repo/releases/tag/v1.0.0
```

### Section Ordering

Categories appear in tier order, then alphabetically within tier:

1. **Core**: Added, Changed, Deprecated, Fixed, Removed, Security
2. **Standard**: Breaking, Dependencies, Highlights, Performance, Upgrade Guide
3. **Extended**: Build, Contributors, Documentation, Tests
4. **Optional**: Compliance, Infrastructure, Internal, Known Issues, Observability

Empty sections are omitted.

## Integration with Localization

Rendering respects the `--locale` flag (see `PRD_LOCALIZATION.md`):

```bash
sclog generate --locale fr changelog.json
```

Localization affects only **labels** (section headers), not structure or user content.

## References

- [CommonMark Spec](https://spec.commonmark.org/)
- [GitHub Flavored Markdown](https://github.github.com/gfm/)
- [Pandoc Markdown](https://pandoc.org/MANUAL.html#pandocs-markdown)
