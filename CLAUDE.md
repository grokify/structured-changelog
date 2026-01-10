# CLAUDE.md

Project-specific instructions for Claude Code.

## Changelog Generation

Use the `/changelog` slash command to generate changelogs:

```
/changelog <since-tag> [new-version]
```

Or manually:

1. Parse commits: `sclog parse-commits --since=<tag>`
2. Update CHANGELOG.json with appropriate entries
3. Validate: `sclog validate CHANGELOG.json`
4. Generate: `sclog generate CHANGELOG.json -o CHANGELOG.md`

## Development

```bash
go test ./...              # Run tests
golangci-lint run          # Lint
go build ./cmd/sclog       # Build CLI
```

## Key Packages

- `changelog/` - JSON IR structs, validation, change types
- `renderer/` - Deterministic Markdown generation
- `gitlog/` - Git log parsing, conventional commits
- `cmd/sclog/` - CLI commands
