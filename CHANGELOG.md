# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-01-03

### Added

- JSON IR schema for structured changelogs (v1.0)
- `changelog` package with IR types, parsing, and validation
- `renderer` package with deterministic Markdown generation
- `sclog` CLI with `validate` and `generate` subcommands
- Support for Keep a Changelog categories: added, changed, deprecated, removed, fixed, security
- Optional security metadata: CVE, GHSA, severity, CVSS, CWE
- Optional SBOM metadata: component, version, license
- JSON Schema for IR validation
- Example changelogs: basic, security, full
- Documentation: spec, security guide, SBOM guide
