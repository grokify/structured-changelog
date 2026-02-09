package changelog

import (
	"fmt"
	"strings"
)

// ErrorCode represents a validation error code.
type ErrorCode string

// Validation error codes.
const (
	// Format errors (E0xx)
	ErrCodeInvalidDate       ErrorCode = "E001"
	ErrCodeInvalidVersion    ErrorCode = "E002"
	ErrCodeInvalidCVE        ErrorCode = "E003"
	ErrCodeInvalidGHSA       ErrorCode = "E004"
	ErrCodeInvalidSeverity   ErrorCode = "E005"
	ErrCodeInvalidCVSSScore  ErrorCode = "E006"
	ErrCodeInvalidIRVersion  ErrorCode = "E007"
	ErrCodeInvalidVersioning ErrorCode = "E008"
	ErrCodeInvalidCommitConv ErrorCode = "E009"

	// Structure errors (E1xx)
	ErrCodeMissingField     ErrorCode = "E100"
	ErrCodeDuplicateVersion ErrorCode = "E101"
	ErrCodeUnsortedReleases ErrorCode = "E102"
	ErrCodeEmptyDescription ErrorCode = "E103"

	// Warning codes (W0xx)
	WarnCodeMissingCVE       ErrorCode = "W001"
	WarnCodeShortDescription ErrorCode = "W002"
	WarnCodeNoTierCoverage   ErrorCode = "W003"
	WarnCodeMissingSeverity  ErrorCode = "W004"
)

// Severity represents the severity of a validation issue.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// RichValidationError provides detailed, actionable validation feedback.
type RichValidationError struct {
	Code          ErrorCode `json:"code"`
	Severity      Severity  `json:"severity"`
	Path          string    `json:"path"`
	Message       string    `json:"message"`
	Actual        string    `json:"actual,omitempty"`
	Expected      string    `json:"expected,omitempty"`
	Suggestion    string    `json:"suggestion,omitempty"`
	Documentation string    `json:"documentation,omitempty"`
}

// Error implements the error interface.
func (e RichValidationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Path, e.Message)
}

// RichValidationSummary provides summary statistics.
type RichValidationSummary struct {
	ErrorCount   int `json:"errorCount"`
	WarningCount int `json:"warningCount"`
	ReleaseCount int `json:"releasesChecked"`
	EntriesCount int `json:"entriesChecked"`
}

// RichValidationResult holds comprehensive validation results.
type RichValidationResult struct {
	Valid    bool                  `json:"valid"`
	Errors   []RichValidationError `json:"errors,omitempty"`
	Warnings []RichValidationError `json:"warnings,omitempty"`
	Summary  RichValidationSummary `json:"summary"`
}

// ValidateRich performs validation with rich, actionable error messages.
func (c *Changelog) ValidateRich() RichValidationResult {
	result := RichValidationResult{
		Valid: true,
	}

	var entriesCount int

	// Check required fields
	if c.Project == "" {
		result.addError(RichValidationError{
			Code:       ErrCodeMissingField,
			Severity:   SeverityError,
			Path:       "project",
			Message:    "Project name is required",
			Expected:   "Non-empty string",
			Suggestion: "Add a \"project\" field with your project name",
		})
	}

	if c.IRVersion != IRVersion {
		result.addError(RichValidationError{
			Code:          ErrCodeInvalidIRVersion,
			Severity:      SeverityError,
			Path:          "ir_version",
			Message:       "Invalid or unsupported IR version",
			Actual:        c.IRVersion,
			Expected:      IRVersion,
			Suggestion:    fmt.Sprintf("Set ir_version to \"%s\"", IRVersion),
			Documentation: "https://github.com/grokify/structured-changelog#ir-version",
		})
	}

	// Validate versioning scheme
	if !validVersioningSchemes[c.Versioning] {
		result.addError(RichValidationError{
			Code:       ErrCodeInvalidVersioning,
			Severity:   SeverityError,
			Path:       "versioning",
			Message:    "Invalid versioning scheme",
			Actual:     c.Versioning,
			Expected:   "One of: semver, calver, custom, none (or omit for default)",
			Suggestion: "Use \"semver\" for Semantic Versioning or \"calver\" for Calendar Versioning",
		})
	}

	// Validate commit convention
	if !validCommitConventions[c.CommitConvention] {
		result.addError(RichValidationError{
			Code:          ErrCodeInvalidCommitConv,
			Severity:      SeverityError,
			Path:          "commit_convention",
			Message:       "Invalid commit convention",
			Actual:        c.CommitConvention,
			Expected:      "One of: conventional, none (or omit for default)",
			Suggestion:    "Use \"conventional\" for Conventional Commits specification",
			Documentation: "https://www.conventionalcommits.org/",
		})
	}

	// Validate unreleased section
	if c.Unreleased != nil {
		entriesCount += c.validateReleaseRich(c.Unreleased, "unreleased", &result, true)
	}

	// Validate releases
	versions := make(map[string]bool)
	for i, release := range c.Releases {
		field := fmt.Sprintf("releases[%d]", i)
		entriesCount += c.validateReleaseRich(&release, field, &result, false)

		// Check for duplicate versions
		if release.Version != "" {
			if versions[release.Version] {
				result.addError(RichValidationError{
					Code:       ErrCodeDuplicateVersion,
					Severity:   SeverityError,
					Path:       field + ".version",
					Message:    "Duplicate version found",
					Actual:     release.Version,
					Suggestion: "Each version must be unique; remove or rename the duplicate",
				})
			}
			versions[release.Version] = true
		}
	}

	result.Summary = RichValidationSummary{
		ErrorCount:   len(result.Errors),
		WarningCount: len(result.Warnings),
		ReleaseCount: len(c.Releases),
		EntriesCount: entriesCount,
	}

	return result
}

func (c *Changelog) validateReleaseRich(r *Release, field string, result *RichValidationResult, isUnreleased bool) int {
	entriesCount := 0

	// Version and date required for releases (not unreleased)
	if !isUnreleased {
		if r.Version == "" {
			result.addError(RichValidationError{
				Code:          ErrCodeMissingField,
				Severity:      SeverityError,
				Path:          field + ".version",
				Message:       "Version is required",
				Expected:      "Semantic version (e.g., 1.0.0)",
				Suggestion:    "Add a version following SemVer 2.0.0 format",
				Documentation: "https://semver.org/",
			})
		} else if !semverRegex.MatchString(r.Version) {
			result.addError(RichValidationError{
				Code:          ErrCodeInvalidVersion,
				Severity:      SeverityError,
				Path:          field + ".version",
				Message:       "Invalid semantic version format",
				Actual:        r.Version,
				Expected:      "MAJOR.MINOR.PATCH (e.g., 1.0.0, 2.1.3-beta.1)",
				Suggestion:    suggestVersionFix(r.Version),
				Documentation: "https://semver.org/",
			})
		}

		if r.Date == "" {
			result.addError(RichValidationError{
				Code:          ErrCodeMissingField,
				Severity:      SeverityError,
				Path:          field + ".date",
				Message:       "Date is required",
				Expected:      "YYYY-MM-DD format (ISO 8601)",
				Suggestion:    "Add a date in YYYY-MM-DD format",
				Documentation: "https://keepachangelog.com/en/1.1.0/#how",
			})
		} else if !dateRegex.MatchString(r.Date) {
			result.addError(RichValidationError{
				Code:          ErrCodeInvalidDate,
				Severity:      SeverityError,
				Path:          field + ".date",
				Message:       "Invalid date format",
				Actual:        r.Date,
				Expected:      "YYYY-MM-DD format (ISO 8601)",
				Suggestion:    suggestDateFix(r.Date),
				Documentation: "https://keepachangelog.com/en/1.1.0/#how",
			})
		}
	}

	// Validate all entries
	entriesCount += c.validateEntriesRich(r.Highlights, field+".highlights", result)
	entriesCount += c.validateEntriesRich(r.Breaking, field+".breaking", result)
	entriesCount += c.validateEntriesRich(r.UpgradeGuide, field+".upgrade_guide", result)
	entriesCount += c.validateSecurityEntriesRich(r.Security, field+".security", result)
	entriesCount += c.validateEntriesRich(r.Added, field+".added", result)
	entriesCount += c.validateEntriesRich(r.Changed, field+".changed", result)
	entriesCount += c.validateEntriesRich(r.Deprecated, field+".deprecated", result)
	entriesCount += c.validateEntriesRich(r.Removed, field+".removed", result)
	entriesCount += c.validateEntriesRich(r.Fixed, field+".fixed", result)
	entriesCount += c.validateEntriesRich(r.Performance, field+".performance", result)
	entriesCount += c.validateEntriesRich(r.Dependencies, field+".dependencies", result)
	entriesCount += c.validateEntriesRich(r.Documentation, field+".documentation", result)
	entriesCount += c.validateEntriesRich(r.Build, field+".build", result)
	entriesCount += c.validateEntriesRich(r.Tests, field+".tests", result)
	entriesCount += c.validateEntriesRich(r.Infrastructure, field+".infrastructure", result)
	entriesCount += c.validateEntriesRich(r.Observability, field+".observability", result)
	entriesCount += c.validateEntriesRich(r.Compliance, field+".compliance", result)
	entriesCount += c.validateEntriesRich(r.Internal, field+".internal", result)
	entriesCount += c.validateEntriesRich(r.KnownIssues, field+".known_issues", result)
	entriesCount += c.validateEntriesRich(r.Contributors, field+".contributors", result)

	return entriesCount
}

func (c *Changelog) validateEntriesRich(entries []Entry, field string, result *RichValidationResult) int {
	for i, entry := range entries {
		entryField := fmt.Sprintf("%s[%d]", field, i)
		if entry.Description == "" {
			result.addError(RichValidationError{
				Code:       ErrCodeEmptyDescription,
				Severity:   SeverityError,
				Path:       entryField + ".description",
				Message:    "Entry description is required",
				Expected:   "Non-empty description of the change",
				Suggestion: "Add a description explaining what changed and why",
			})
		} else if len(entry.Description) < 10 {
			result.addWarning(RichValidationError{
				Code:       WarnCodeShortDescription,
				Severity:   SeverityWarning,
				Path:       entryField + ".description",
				Message:    "Description is very short",
				Actual:     entry.Description,
				Suggestion: "Consider providing more detail about the change",
			})
		}
	}
	return len(entries)
}

func (c *Changelog) validateSecurityEntriesRich(entries []Entry, field string, result *RichValidationResult) int {
	for i, entry := range entries {
		entryField := fmt.Sprintf("%s[%d]", field, i)

		if entry.Description == "" {
			result.addError(RichValidationError{
				Code:       ErrCodeEmptyDescription,
				Severity:   SeverityError,
				Path:       entryField + ".description",
				Message:    "Entry description is required",
				Expected:   "Non-empty description of the security issue",
				Suggestion: "Add a description explaining the security issue and fix",
			})
		}

		if entry.CVE != "" && !cveRegex.MatchString(entry.CVE) {
			result.addError(RichValidationError{
				Code:          ErrCodeInvalidCVE,
				Severity:      SeverityError,
				Path:          entryField + ".cve",
				Message:       "Invalid CVE format",
				Actual:        entry.CVE,
				Expected:      "CVE-YYYY-NNNNN (e.g., CVE-2024-12345)",
				Suggestion:    suggestCVEFix(entry.CVE),
				Documentation: "https://cve.mitre.org/cve/identifiers/",
			})
		}

		if entry.GHSA != "" && !ghsaRegex.MatchString(entry.GHSA) {
			result.addError(RichValidationError{
				Code:          ErrCodeInvalidGHSA,
				Severity:      SeverityError,
				Path:          entryField + ".ghsa",
				Message:       "Invalid GHSA format",
				Actual:        entry.GHSA,
				Expected:      "GHSA-xxxx-xxxx-xxxx (e.g., GHSA-abcd-1234-efgh)",
				Documentation: "https://github.com/advisories",
			})
		}

		if entry.Severity != "" && !validSeverities[entry.Severity] {
			result.addError(RichValidationError{
				Code:       ErrCodeInvalidSeverity,
				Severity:   SeverityError,
				Path:       entryField + ".severity",
				Message:    "Invalid severity level",
				Actual:     entry.Severity,
				Expected:   "One of: critical, high, medium, low, informational",
				Suggestion: suggestSeverityFix(entry.Severity),
			})
		}

		if entry.CVSSScore != 0 && (entry.CVSSScore < 0 || entry.CVSSScore > 10) {
			result.addError(RichValidationError{
				Code:          ErrCodeInvalidCVSSScore,
				Severity:      SeverityError,
				Path:          entryField + ".cvss_score",
				Message:       "CVSS score out of range",
				Actual:        fmt.Sprintf("%.1f", entry.CVSSScore),
				Expected:      "Value between 0.0 and 10.0",
				Documentation: "https://www.first.org/cvss/",
			})
		}

		// Warnings for missing but recommended fields
		if entry.CVE == "" && entry.GHSA == "" {
			result.addWarning(RichValidationError{
				Code:       WarnCodeMissingCVE,
				Severity:   SeverityWarning,
				Path:       entryField,
				Message:    "Security entry missing CVE or GHSA identifier",
				Suggestion: "Add 'cve' field with format CVE-YYYY-NNNNN or 'ghsa' field",
			})
		}

		if entry.Severity == "" {
			result.addWarning(RichValidationError{
				Code:       WarnCodeMissingSeverity,
				Severity:   SeverityWarning,
				Path:       entryField,
				Message:    "Security entry missing severity level",
				Suggestion: "Add 'severity' field (critical, high, medium, low, or informational)",
			})
		}
	}
	return len(entries)
}

func (r *RichValidationResult) addError(err RichValidationError) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

func (r *RichValidationResult) addWarning(err RichValidationError) {
	r.Warnings = append(r.Warnings, err)
}

// Suggestion helper functions

func suggestVersionFix(version string) string {
	// Try to provide a helpful suggestion based on common mistakes
	v := strings.TrimPrefix(version, "v")
	if semverRegex.MatchString(v) {
		return fmt.Sprintf("Remove the 'v' prefix: %q", v)
	}
	parts := strings.Split(v, ".")
	if len(parts) == 2 {
		return fmt.Sprintf("Add patch version: %q", v+".0")
	}
	return "Use format MAJOR.MINOR.PATCH (e.g., 1.0.0)"
}

func suggestDateFix(date string) string {
	// Common date format fixes
	date = strings.TrimSpace(date)

	// Try to detect common formats and suggest fixes
	if strings.Contains(date, "/") {
		parts := strings.Split(date, "/")
		if len(parts) == 3 {
			// Could be MM/DD/YYYY or DD/MM/YYYY
			return "Use ISO 8601 format: YYYY-MM-DD"
		}
	}

	// Month name format
	months := []string{"january", "february", "march", "april", "may", "june",
		"july", "august", "september", "october", "november", "december",
		"jan", "feb", "mar", "apr", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	lower := strings.ToLower(date)
	for _, m := range months {
		if strings.Contains(lower, m) {
			return "Convert to ISO 8601 format: YYYY-MM-DD"
		}
	}

	return "Use format YYYY-MM-DD (e.g., 2024-01-15)"
}

func suggestCVEFix(cve string) string {
	upper := strings.ToUpper(cve)
	if !strings.HasPrefix(upper, "CVE-") {
		return fmt.Sprintf("Add CVE prefix: CVE-%s", strings.TrimPrefix(upper, "CVE"))
	}
	return "Use format CVE-YYYY-NNNNN (e.g., CVE-2024-12345)"
}

func suggestSeverityFix(severity string) string {
	lower := strings.ToLower(severity)
	// Check for close matches
	suggestions := map[string]string{
		"crit":      "critical",
		"hi":        "high",
		"med":       "medium",
		"lo":        "low",
		"info":      "informational",
		"none":      "informational",
		"moderate":  "medium",
		"important": "high",
	}
	if suggestion, ok := suggestions[lower]; ok {
		return fmt.Sprintf("Did you mean %q?", suggestion)
	}
	return "Use one of: critical, high, medium, low, informational"
}
