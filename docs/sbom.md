# SBOM Metadata Guide

Structured Changelog supports Software Bill of Materials (SBOM) metadata for tracking component changes.

## Overview

SBOM fields help track:

- Dependency updates and their licenses
- Component version changes
- License compliance
- Supply chain visibility

## Fields

### Component

The name of the dependency or component:

```json
{
  "description": "Upgrade HTTP client library",
  "component": "github.com/go-resty/resty/v2"
}
```

### Component Version

The new version of the component:

```json
{
  "description": "Upgrade database driver",
  "component": "github.com/lib/pq",
  "component_version": "1.10.9"
}
```

### License

SPDX license identifier:

```json
{
  "description": "Add logging library",
  "component": "go.uber.org/zap",
  "component_version": "1.26.0",
  "license": "MIT"
}
```

Common SPDX identifiers:

- `MIT`
- `Apache-2.0`
- `BSD-3-Clause`
- `GPL-3.0-only`
- `LGPL-2.1-or-later`
- `MPL-2.0`

## Use Cases

### Dependency Updates

Track when dependencies are updated:

```json
{
  "version": "1.5.0",
  "date": "2026-01-03",
  "changed": [
    {
      "description": "Upgrade Go to 1.23",
      "component": "go",
      "component_version": "1.23.0"
    },
    {
      "description": "Update protobuf library",
      "component": "google.golang.org/protobuf",
      "component_version": "1.32.0",
      "license": "BSD-3-Clause"
    }
  ]
}
```

### New Dependencies

Track when new dependencies are added:

```json
{
  "added": [
    {
      "description": "Add OpenTelemetry instrumentation",
      "component": "go.opentelemetry.io/otel",
      "component_version": "1.21.0",
      "license": "Apache-2.0"
    }
  ]
}
```

### Removed Dependencies

Track when dependencies are removed:

```json
{
  "removed": [
    {
      "description": "Remove deprecated logging library",
      "component": "github.com/sirupsen/logrus"
    }
  ]
}
```

### Security Updates

Combine with security metadata:

```json
{
  "security": [
    {
      "description": "Upgrade crypto library to fix timing attack",
      "component": "golang.org/x/crypto",
      "component_version": "0.18.0",
      "license": "BSD-3-Clause",
      "severity": "medium",
      "cve": "CVE-2026-54321"
    }
  ]
}
```

## Integration with SBOM Tools

### CycloneDX

The component metadata can be used to generate CycloneDX SBOMs:

```json
{
  "component": "github.com/example/lib",
  "component_version": "1.2.3",
  "license": "MIT"
}
```

Maps to CycloneDX:

```xml
<component type="library">
  <name>github.com/example/lib</name>
  <version>1.2.3</version>
  <licenses>
    <license>
      <id>MIT</id>
    </license>
  </licenses>
</component>
```

### SPDX

Similarly maps to SPDX format for compliance reporting.

## Best Practices

1. **Use SPDX identifiers** - Ensures license compatibility checking
2. **Include version numbers** - Essential for reproducibility
3. **Track all dependency changes** - Complete audit trail
4. **Note breaking changes** - Use the `breaking` field for major upgrades
5. **Combine with security data** - Link component updates to CVEs when relevant

## Resources

- [SPDX License List](https://spdx.org/licenses/)
- [CycloneDX Specification](https://cyclonedx.org/)
- [SPDX Specification](https://spdx.dev/)
- [NTIA SBOM Guidance](https://www.ntia.gov/SBOM)
