---
description: Generate or update CHANGELOG.json and CHANGELOG.md from git history
allowed-tools: Bash, Read, Write, Edit, Glob, Grep
argument-hint: <since-tag> [version]
---

# Structured Changelog Generator

Generate a structured changelog from git commits using the sclog CLI.

## Arguments

- `$ARGUMENTS` - First argument is the since-tag (e.g., v0.4.0), optional second is the new version

## Workflow

1. **Parse commits** from the specified tag to HEAD
2. **Review and categorize** changes into appropriate sections
3. **Create or update** CHANGELOG.json with new entries
4. **Validate** the changelog structure
5. **Generate** CHANGELOG.md from the JSON

## Instructions

Run `sclog parse-commits --since=<tag>` to get structured commit data, then:

- Group related commits into single changelog entries
- Use descriptions that explain "why" not just "what"
- Mark breaking changes with `"breaking": true`
- Use appropriate categories: Added, Changed, Deprecated, Removed, Fixed, Security
- For extended types: Highlights, Performance, Dependencies, Documentation, Build, Tests

After creating/updating CHANGELOG.json:

1. Validate: `sclog validate CHANGELOG.json`
2. Generate: `sclog generate CHANGELOG.json -o CHANGELOG.md`

## JSON Entry Format

```json
{
  "version": "X.Y.Z",
  "date": "YYYY-MM-DD",
  "added": [
    { "description": "New feature description", "issue": "123" }
  ],
  "changed": [
    { "description": "Modified behavior", "breaking": true }
  ],
  "fixed": [
    { "description": "Bug fix description", "pr": "456" }
  ]
}
```

## Reference

See the [LLM Guide](https://grokify.github.io/structured-changelog/guides/llm-guide/) for detailed workflows.
