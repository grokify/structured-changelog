# CLAUDE.md

Project-specific instructions for Claude Code.

## Changelog Generation

Use the `/changelog` slash command to generate changelogs:

```
/changelog <since-tag> [new-version]
```

Or manually:

1. Parse commits: `schangelog parse-commits --since=<tag>`
2. Update CHANGELOG.json with appropriate entries
3. Validate: `schangelog validate CHANGELOG.json`
4. Generate: `schangelog generate CHANGELOG.json -o CHANGELOG.md`

## Development

```bash
go test ./...              # Run tests
golangci-lint run          # Lint
go build ./cmd/schangelog       # Build CLI
```

## Git Log Alternatives

**Always prefer `schangelog` commands over raw git log** for commit analysis. The schangelog CLI outputs TOON format by default, which is ~8x more token-efficient than raw git log.

| Instead of | Use |
|------------|-----|
| `git log --oneline --reverse <tag>..HEAD` | `schangelog parse-commits --since=<tag>` |
| `git log --format="%H %s" --reverse` | `schangelog parse-commits --since=<tag>` |
| `git log --oneline \| head -N` | `schangelog parse-commits --last=N` |
| `git log <from>..<to>` | `schangelog parse-commits --since=<from> --until=<to>` |

Benefits of schangelog:

- Pre-parsed conventional commits (type, scope, subject)
- Suggested changelog categories for each commit
- Summary statistics by type and category
- TOON format optimized for LLM consumption

Use `--format=json` when you need to debug or pipe to jq.

## Key Packages

- `changelog/` - JSON IR structs, validation, change types
- `renderer/` - Deterministic Markdown generation
- `gitlog/` - Git log parsing, conventional commits
- `cmd/schangelog/` - CLI commands
