# Security Metadata Guide

Structured Changelog supports rich security metadata for vulnerability disclosures.

## Overview

Security entries can include CVE/GHSA identifiers, severity ratings, CVSS scores, and links to SARIF static analysis results. This enables:

- Automated vulnerability tracking
- Integration with security scanners
- Compliance reporting
- Audit trails

## Fields

### CVE (Common Vulnerabilities and Exposures)

```json
{
  "description": "Fix SQL injection vulnerability",
  "cve": "CVE-2026-12345"
}
```

Format: `CVE-YYYY-NNNNN` where YYYY is the year and NNNNN is at least 4 digits.

### GHSA (GitHub Security Advisory)

```json
{
  "description": "Fix remote code execution",
  "ghsa": "GHSA-abcd-efgh-ijkl"
}
```

Format: `GHSA-xxxx-xxxx-xxxx` where x is lowercase alphanumeric.

### Severity

```json
{
  "description": "Fix authentication bypass",
  "severity": "critical"
}
```

Valid values:

- `critical` - Immediate action required
- `high` - High impact, fix soon
- `medium` - Moderate impact
- `low` - Low impact
- `informational` - No immediate risk

### CVSS (Common Vulnerability Scoring System)

```json
{
  "description": "Fix buffer overflow",
  "cvssScore": 8.5,
  "cvssVector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H"
}
```

Score ranges:

- 9.0-10.0: Critical
- 7.0-8.9: High
- 4.0-6.9: Medium
- 0.1-3.9: Low

### CWE (Common Weakness Enumeration)

```json
{
  "description": "Fix SQL injection",
  "cwe": "CWE-89"
}
```

Common CWEs:

- CWE-79: Cross-site Scripting (XSS)
- CWE-89: SQL Injection
- CWE-94: Code Injection
- CWE-287: Improper Authentication
- CWE-798: Hard-coded Credentials

### Version Ranges

```json
{
  "description": "Fix privilege escalation",
  "affectedVersions": "<2.0.0",
  "patchedVersions": ">=2.0.0"
}
```

### SARIF Integration

Link to static analysis results:

```json
{
  "description": "Fix issue identified by CodeQL",
  "sarifRuleId": "go/sql-injection"
}
```

## Complete Example

```json
{
  "irVersion": "1.0",
  "project": "secure-app",
  "releases": [
    {
      "version": "2.1.1",
      "date": "2026-01-03",
      "security": [
        {
          "description": "Fix SQL injection in user search endpoint",
          "cve": "CVE-2026-12345",
          "ghsa": "GHSA-abcd-efgh-ijkl",
          "severity": "high",
          "cvssScore": 8.5,
          "cvssVector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N",
          "cwe": "CWE-89",
          "affectedVersions": "<2.1.1",
          "patchedVersions": ">=2.1.1",
          "sarifRuleId": "go/sql-injection",
          "pr": "#123"
        }
      ]
    }
  ]
}
```

## Rendered Output

With `IncludeSecurityMetadata: true`:

```markdown
### Security

- Fix SQL injection in user search endpoint (CVE-2026-12345, GHSA-abcd-efgh-ijkl, severity: high)
```

## Best Practices

1. **Always include CVE/GHSA when available** - Makes tracking easier
2. **Use severity consistently** - Follow your organization's severity definitions
3. **Include affected versions** - Helps users determine if they're impacted
4. **Link to advisories** - Provide full details for those who need them
5. **Coordinate disclosure** - Don't publish before patches are available

## Resources

- [CVE Database](https://cve.mitre.org/)
- [GitHub Security Advisories](https://github.com/advisories)
- [CWE Database](https://cwe.mitre.org/)
- [CVSS Calculator](https://www.first.org/cvss/calculator/3.1)
