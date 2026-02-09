# 更新日志

此项目的所有重要更改都将记录在此文件中。

格式基于[如何维护更新日志](https://keepachangelog.com/zh-CN/1.1.0/),
本项目遵循[语义化版本](https://semver.org/lang/zh-CN/),
并且 此变更日志由[Structured Changelog](https://github.com/grokify/structured-changelog)生成.

## [未发布]

## [2.0.0] - 2026-01-04

### 亮点

- Major release with breaking API changes and significant performance improvements
- New plugin architecture enables third-party extensions

### 破坏性变更

- **破坏性变更:** Renamed `Config.timeout` to `Config.requestTimeout` for clarity
- **破坏性变更:** Removed deprecated `v1` API endpoints

### 升级指南

- Update config files to use `requestTimeout` instead of `timeout`
- Migrate API calls from `/v1/*` to `/v2/*` endpoints

### 安全

- Fixed authentication bypass vulnerability (CVE-2026-12345, severity: high)

### 新增

- Plugin architecture for third-party extensions ([#456](https://github.com/example/extended-example/pull/456))
- New `/v2/batch` endpoint for bulk operations ([#457](https://github.com/example/extended-example/pull/457))

### 变更

- Improved error messages with actionable suggestions ([#458](https://github.com/example/extended-example/pull/458))

### 弃用

- Legacy authentication method will be removed in v3.0

### 移除

- Removed `/v1/*` API endpoints (deprecated since v1.5)

### 修复

- Fixed memory leak in connection pool ([#123](https://github.com/example/extended-example/issues/123))
- Fixed race condition in concurrent requests ([#124](https://github.com/example/extended-example/issues/124))

### 性能

- Reduced API response latency by 40% through caching
- Optimized database queries for large datasets

### 依赖

- Upgraded Go from 1.21 to 1.22
- Updated redis client to v9.0.0

### 文档

- Added migration guide for v1 to v2 upgrade
- Updated API reference with new endpoints

### 构建

- Added multi-platform Docker builds (amd64, arm64)
- Migrated CI from CircleCI to GitHub Actions

### 已知问题

- Batch endpoint limited to 1000 items per request
- ARM64 builds not yet tested on Windows

### 贡献者

- @alice - Plugin architecture design and implementation
- @bob - Performance optimization work
- @charlie - Documentation updates

## [1.5.0] - 2025-12-01

### 新增

- Initial v2 API preview (experimental) ([#400](https://github.com/example/extended-example/pull/400))

### 弃用

- v1 API endpoints deprecated, will be removed in v2.0

### 修复

- Fixed timeout handling in long-running requests ([#100](https://github.com/example/extended-example/issues/100))

### 性能

- Improved startup time by lazy-loading configuration

[unreleased]: https://github.com/example/extended-example/compare/2.0.0...HEAD
[2.0.0]: https://github.com/example/extended-example/compare/1.5.0...2.0.0
[1.5.0]: https://github.com/example/extended-example/releases/tag/1.5.0
