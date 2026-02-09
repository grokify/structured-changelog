# Journal des modifications

Tous les changements notables de ce projet seront documentés dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.1.0/),
ce projet adhère au [Versionnement Sémantique](https://semver.org/lang/fr/),
les commits suivent [Conventional Commits](https://www.conventionalcommits.org/fr/v1.0.0/),
et ce changelog est généré par [Structured Changelog](https://github.com/grokify/structured-changelog).

## [Non publié]

## [2.0.0] - 2026-01-04

### Points forts

- Major release with breaking API changes and significant performance improvements
- New plugin architecture enables third-party extensions

### Ruptures

- **RUPTURE :** Renamed `Config.timeout` to `Config.requestTimeout` for clarity
- **RUPTURE :** Removed deprecated `v1` API endpoints

### Guide de mise à niveau

- Update config files to use `requestTimeout` instead of `timeout`
- Migrate API calls from `/v1/*` to `/v2/*` endpoints

### Sécurité

- Fixed authentication bypass vulnerability (CVE-2026-12345, severity: high)

### Ajouté

- Plugin architecture for third-party extensions ([#456](https://github.com/example/extended-example/pull/456))
- New `/v2/batch` endpoint for bulk operations ([#457](https://github.com/example/extended-example/pull/457))

### Modifié

- Improved error messages with actionable suggestions ([#458](https://github.com/example/extended-example/pull/458))

### Obsolète

- Legacy authentication method will be removed in v3.0

### Supprimé

- Removed `/v1/*` API endpoints (deprecated since v1.5)

### Corrigé

- Fixed memory leak in connection pool ([#123](https://github.com/example/extended-example/issues/123))
- Fixed race condition in concurrent requests ([#124](https://github.com/example/extended-example/issues/124))

### Performance

- Reduced API response latency by 40% through caching
- Optimized database queries for large datasets

### Dépendances

- Upgraded Go from 1.21 to 1.22
- Updated redis client to v9.0.0

### Documentation

- Added migration guide for v1 to v2 upgrade
- Updated API reference with new endpoints

### Build

- Added multi-platform Docker builds (amd64, arm64)
- Migrated CI from CircleCI to GitHub Actions

### Problèmes connus

- Batch endpoint limited to 1000 items per request
- ARM64 builds not yet tested on Windows

### Contributeurs

- @alice - Plugin architecture design and implementation
- @bob - Performance optimization work
- @charlie - Documentation updates

## [1.5.0] - 2025-12-01

### Ajouté

- Initial v2 API preview (experimental) ([#400](https://github.com/example/extended-example/pull/400))

### Obsolète

- v1 API endpoints deprecated, will be removed in v2.0

### Corrigé

- Fixed timeout handling in long-running requests ([#100](https://github.com/example/extended-example/issues/100))

### Performance

- Improved startup time by lazy-loading configuration

[unreleased]: https://github.com/example/extended-example/compare/2.0.0...HEAD
[2.0.0]: https://github.com/example/extended-example/compare/1.5.0...2.0.0
[1.5.0]: https://github.com/example/extended-example/releases/tag/1.5.0
