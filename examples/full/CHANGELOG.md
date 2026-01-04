# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Prometheus metrics exporter

## [0.5.0] - 2026-01-03

### Added

- `AnnotationManager` interface for span/trace annotations (#42)
- `DatasetManager` new methods: `GetDatasetByID`, `DeleteDataset` (#43)
- Prompt model/provider options: `WithPromptModel`, `WithPromptProvider`
- OmniLLM hook auto-creates traces when none exists in context
- Trace context helpers (`contextWithTrace`, `traceFromContext`)
- `llmops/metrics` package with evaluation metrics (#45)

### Changed

- **BREAKING:** Provider adapters moved to standalone SDKs

### Removed

- **BREAKING:** `llmops/opik` adapter (moved to go-opik)
- **BREAKING:** `llmops/phoenix` adapter (moved to go-phoenix)
- **BREAKING:** `sdk/phoenix` package (use go-phoenix directly)

## [0.4.0] - 2025-12-01

### Added

- Initial provider abstraction layer
- Opik integration
- Phoenix integration
- Langfuse integration
