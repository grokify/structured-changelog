# 変更履歴

このプロジェクトへのすべての注目すべき変更は、このファイルに記載されます。

このフォーマットは[Keep a Changelog](https://keepachangelog.com/ja/1.1.0/)に基づいています,
このプロジェクトは[セマンティック バージョニング](https://semver.org/lang/ja/)に準拠しています,
コミットは[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/)に従っています,
そして この変更履歴は[Structured Changelog](https://github.com/grokify/structured-changelog)によって生成されています.

## [未リリース]

## [2.0.0] - 2026-01-04

### ハイライト

- Major release with breaking API changes and significant performance improvements
- New plugin architecture enables third-party extensions

### 破壊的変更

- **破壊的変更:** Renamed `Config.timeout` to `Config.requestTimeout` for clarity
- **破壊的変更:** Removed deprecated `v1` API endpoints

### アップグレードガイド

- Update config files to use `requestTimeout` instead of `timeout`
- Migrate API calls from `/v1/*` to `/v2/*` endpoints

### セキュリティ

- Fixed authentication bypass vulnerability (CVE-2026-12345, severity: high)

### 追加

- Plugin architecture for third-party extensions ([#456](https://github.com/example/extended-example/pull/456))
- New `/v2/batch` endpoint for bulk operations ([#457](https://github.com/example/extended-example/pull/457))

### 変更

- Improved error messages with actionable suggestions ([#458](https://github.com/example/extended-example/pull/458))

### 非推奨

- Legacy authentication method will be removed in v3.0

### 削除

- Removed `/v1/*` API endpoints (deprecated since v1.5)

### 修正

- Fixed memory leak in connection pool ([#123](https://github.com/example/extended-example/issues/123))
- Fixed race condition in concurrent requests ([#124](https://github.com/example/extended-example/issues/124))

### パフォーマンス

- Reduced API response latency by 40% through caching
- Optimized database queries for large datasets

### 依存関係

- Upgraded Go from 1.21 to 1.22
- Updated redis client to v9.0.0

### ドキュメント

- Added migration guide for v1 to v2 upgrade
- Updated API reference with new endpoints

### ビルド

- Added multi-platform Docker builds (amd64, arm64)
- Migrated CI from CircleCI to GitHub Actions

### 既知の問題

- Batch endpoint limited to 1000 items per request
- ARM64 builds not yet tested on Windows

### コントリビューター

- @alice - Plugin architecture design and implementation
- @bob - Performance optimization work
- @charlie - Documentation updates

## [1.5.0] - 2025-12-01

### 追加

- Initial v2 API preview (experimental) ([#400](https://github.com/example/extended-example/pull/400))

### 非推奨

- v1 API endpoints deprecated, will be removed in v2.0

### 修正

- Fixed timeout handling in long-running requests ([#100](https://github.com/example/extended-example/issues/100))

### パフォーマンス

- Improved startup time by lazy-loading configuration

[unreleased]: https://github.com/example/extended-example/compare/2.0.0...HEAD
[2.0.0]: https://github.com/example/extended-example/compare/1.5.0...2.0.0
[1.5.0]: https://github.com/example/extended-example/releases/tag/1.5.0
