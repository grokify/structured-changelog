# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.1] - 2026-01-03

### Security

- Fix SQL injection vulnerability in user search (CVE-2026-12345, GHSA-abcd-efgh-ijkl, severity: high)
- Fix XSS vulnerability in comment rendering (CVE-2026-12346, severity: medium)

### Fixed

- Improved input validation across all forms

## [2.1.0] - 2025-12-15

### Added

- Two-factor authentication support
- Audit logging for sensitive operations

### Security

- Upgrade bcrypt to address timing attack (severity: low)
