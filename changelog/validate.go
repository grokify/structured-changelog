package changelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validation errors.
var (
	ErrEmptyProject      = errors.New("project name is required")
	ErrInvalidIRVersion  = errors.New("invalid or unsupported IR version")
	ErrInvalidVersion    = errors.New("invalid semantic version")
	ErrInvalidDate       = errors.New("invalid date format (expected YYYY-MM-DD)")
	ErrEmptyDescription  = errors.New("entry description is required")
	ErrInvalidCVE        = errors.New("invalid CVE format")
	ErrInvalidGHSA       = errors.New("invalid GHSA format")
	ErrInvalidSeverity   = errors.New("invalid severity level")
	ErrInvalidCVSSScore  = errors.New("CVSS score must be between 0 and 10")
	ErrDuplicateVersion  = errors.New("duplicate version found")
	ErrUnsortedReleases  = errors.New("releases are not in reverse chronological order")
	ErrInvalidVersioning = errors.New("invalid versioning scheme")
	ErrInvalidCommitConv = errors.New("invalid commit convention")
)

var validVersioningSchemes = map[string]bool{
	"":               true, // empty is valid (defaults to semver)
	VersioningSemVer: true,
	VersioningCalVer: true,
	VersioningCustom: true,
	VersioningNone:   true,
}

var validCommitConventions = map[string]bool{
	"":                           true, // empty is valid (defaults to none)
	CommitConventionConventional: true,
	CommitConventionNone:         true,
}

var (
	// semverRegex matches semantic versions with optional v prefix (e.g., "1.0.0" or "v1.0.0")
	semverRegex = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	dateRegex   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	cveRegex    = regexp.MustCompile(`^CVE-\d{4}-\d{4,}$`)
	ghsaRegex   = regexp.MustCompile(`^GHSA-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}$`)
)

var validSeverities = map[string]bool{
	"critical":      true,
	"high":          true,
	"medium":        true,
	"low":           true,
	"informational": true,
}

// ValidationError contains details about a validation failure.
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ValidationResult holds the results of changelog validation.
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// Validate validates the changelog structure and content.
func (c *Changelog) Validate() ValidationResult {
	result := ValidationResult{Valid: true}

	// Check required fields
	if c.Project == "" {
		result.addError("project", "project name is required", ErrEmptyProject)
	}

	if c.IRVersion != IRVersion {
		result.addError("ir_version", fmt.Sprintf("expected %s, got %s", IRVersion, c.IRVersion), ErrInvalidIRVersion)
	}

	// Validate versioning scheme
	if !validVersioningSchemes[c.Versioning] {
		result.addError("versioning", fmt.Sprintf("invalid versioning scheme: %s (must be one of semver, calver, custom, none)", c.Versioning), ErrInvalidVersioning)
	}

	// Validate commit convention
	if !validCommitConventions[c.CommitConvention] {
		result.addError("commit_convention", fmt.Sprintf("invalid commit convention: %s (must be one of conventional, none)", c.CommitConvention), ErrInvalidCommitConv)
	}

	// Validate unreleased section
	if c.Unreleased != nil {
		c.validateRelease(c.Unreleased, "unreleased", &result, true)
	}

	// Validate releases
	versions := make(map[string]bool)
	for i, release := range c.Releases {
		field := fmt.Sprintf("releases[%d]", i)
		c.validateRelease(&release, field, &result, false)

		// Check for duplicate versions
		if release.Version != "" {
			if versions[release.Version] {
				result.addError(field+".version", "duplicate version: "+release.Version, ErrDuplicateVersion)
			}
			versions[release.Version] = true
		}
	}

	return result
}

func (c *Changelog) validateRelease(r *Release, field string, result *ValidationResult, isUnreleased bool) {
	// Version and date required for releases (not unreleased)
	if !isUnreleased {
		if r.Version == "" {
			result.addError(field+".version", "version is required", ErrInvalidVersion)
		} else if !semverRegex.MatchString(r.Version) {
			result.addError(field+".version", "invalid semantic version: "+r.Version, ErrInvalidVersion)
		}

		if r.Date == "" {
			result.addError(field+".date", "date is required", ErrInvalidDate)
		} else if !dateRegex.MatchString(r.Date) {
			result.addError(field+".date", "invalid date format: "+r.Date, ErrInvalidDate)
		}
	}

	// Validate all entries in canonical order
	// Overview & Critical
	c.validateEntries(r.Highlights, field+".highlights", result)
	c.validateEntries(r.Breaking, field+".breaking", result)
	c.validateEntries(r.UpgradeGuide, field+".upgrade_guide", result)
	c.validateSecurityEntries(r.Security, field+".security", result)

	// Core KACL
	c.validateEntries(r.Added, field+".added", result)
	c.validateEntries(r.Changed, field+".changed", result)
	c.validateEntries(r.Deprecated, field+".deprecated", result)
	c.validateEntries(r.Removed, field+".removed", result)
	c.validateEntries(r.Fixed, field+".fixed", result)

	// Quality
	c.validateEntries(r.Performance, field+".performance", result)
	c.validateEntries(r.Dependencies, field+".dependencies", result)

	// Development
	c.validateEntries(r.Documentation, field+".documentation", result)
	c.validateEntries(r.Build, field+".build", result)

	// Operations
	c.validateEntries(r.Infrastructure, field+".infrastructure", result)
	c.validateEntries(r.Observability, field+".observability", result)
	c.validateEntries(r.Compliance, field+".compliance", result)

	// Internal
	c.validateEntries(r.Internal, field+".internal", result)

	// End Matter
	c.validateEntries(r.KnownIssues, field+".known_issues", result)
	c.validateEntries(r.Contributors, field+".contributors", result)
}

func (c *Changelog) validateEntries(entries []Entry, field string, result *ValidationResult) {
	for i, entry := range entries {
		entryField := fmt.Sprintf("%s[%d]", field, i)
		if entry.Description == "" {
			result.addError(entryField+".description", "description is required", ErrEmptyDescription)
		}
	}
}

func (c *Changelog) validateSecurityEntries(entries []Entry, field string, result *ValidationResult) {
	for i, entry := range entries {
		entryField := fmt.Sprintf("%s[%d]", field, i)

		if entry.Description == "" {
			result.addError(entryField+".description", "description is required", ErrEmptyDescription)
		}

		if entry.CVE != "" && !cveRegex.MatchString(entry.CVE) {
			result.addError(entryField+".cve", "invalid CVE format: "+entry.CVE, ErrInvalidCVE)
		}

		if entry.GHSA != "" && !ghsaRegex.MatchString(entry.GHSA) {
			result.addError(entryField+".ghsa", "invalid GHSA format: "+entry.GHSA, ErrInvalidGHSA)
		}

		if entry.Severity != "" && !validSeverities[entry.Severity] {
			result.addError(entryField+".severity", "invalid severity: "+entry.Severity, ErrInvalidSeverity)
		}

		if entry.CVSSScore != 0 && (entry.CVSSScore < 0 || entry.CVSSScore > 10) {
			result.addError(entryField+".cvss_score", "CVSS score must be between 0 and 10", ErrInvalidCVSSScore)
		}
	}
}

func (r *ValidationResult) addError(field, message string, err error) {
	r.Valid = false
	r.Errors = append(r.Errors, ValidationError{
		Field:   field,
		Message: message,
		Err:     err,
	})
}

// ErrNoEntriesAtTier is returned when no entries exist at or above the required tier.
var ErrNoEntriesAtTier = errors.New("no entries at or above required tier")

// ErrInvalidTier is returned when an invalid tier string is provided.
var ErrInvalidTier = errors.New("invalid tier")

// ValidateMinTier checks that the changelog has at least one entry at or above
// the specified minimum tier in the latest release.
func (c *Changelog) ValidateMinTier(minTier Tier) error {
	if !minTier.IsValid() {
		return fmt.Errorf("%w: %q (must be one of core, standard, extended, optional)", ErrInvalidTier, minTier)
	}

	if len(c.Releases) == 0 {
		return nil // No releases to validate
	}

	latest := c.Releases[0]
	cats := latest.CategoriesFiltered(minTier)
	if len(cats) == 0 {
		return fmt.Errorf("%w %q in release %s", ErrNoEntriesAtTier, minTier, latest.Version)
	}

	return nil
}

// ParseTier parses a tier string (case-insensitive) and returns the Tier.
func ParseTier(s string) (Tier, error) {
	tier := Tier(strings.ToLower(s))
	if !tier.IsValid() {
		return "", fmt.Errorf("%w: %q (must be one of core, standard, extended, optional)", ErrInvalidTier, s)
	}
	return tier, nil
}
