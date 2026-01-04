package changelog

import (
	"errors"
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test-project",
		Releases: []Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []Entry{{Description: "Initial release"}},
			},
		},
	}

	result := cl.Validate()
	if !result.Valid {
		t.Errorf("expected valid changelog, got errors: %v", result.Errors)
	}
}

func TestValidate_EmptyProject(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "",
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for empty project")
	}
	if !hasError(result.Errors, ErrEmptyProject) {
		t.Error("expected ErrEmptyProject")
	}
}

func TestValidate_InvalidVersion(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{Version: "invalid", Date: "2026-01-03"},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for bad version")
	}
	if !hasError(result.Errors, ErrInvalidVersion) {
		t.Error("expected ErrInvalidVersion")
	}
}

func TestValidate_InvalidDate(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{Version: "1.0.0", Date: "01-03-2026"},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for bad date")
	}
	if !hasError(result.Errors, ErrInvalidDate) {
		t.Error("expected ErrInvalidDate")
	}
}

func TestValidate_InvalidCVE(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", CVE: "invalid-cve"}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for bad CVE")
	}
	if !hasError(result.Errors, ErrInvalidCVE) {
		t.Error("expected ErrInvalidCVE")
	}
}

func TestValidate_ValidCVE(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", CVE: "CVE-2026-12345"}},
			},
		},
	}

	result := cl.Validate()
	if !result.Valid {
		t.Errorf("expected valid changelog, got errors: %v", result.Errors)
	}
}

func TestValidate_InvalidGHSA(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", GHSA: "invalid"}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for bad GHSA")
	}
	if !hasError(result.Errors, ErrInvalidGHSA) {
		t.Error("expected ErrInvalidGHSA")
	}
}

func TestValidate_ValidGHSA(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", GHSA: "GHSA-abcd-efgh-ijkl"}},
			},
		},
	}

	result := cl.Validate()
	if !result.Valid {
		t.Errorf("expected valid changelog, got errors: %v", result.Errors)
	}
}

func TestValidate_DuplicateVersion(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{Version: "1.0.0", Date: "2026-01-03"},
			{Version: "1.0.0", Date: "2026-01-02"},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for duplicate version")
	}
	if !hasError(result.Errors, ErrDuplicateVersion) {
		t.Error("expected ErrDuplicateVersion")
	}
}

func hasError(errs []ValidationError, target error) bool {
	for _, e := range errs {
		if errors.Is(e.Err, target) {
			return true
		}
	}
	return false
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ValidationError
		expected string
	}{
		{
			name:     "with field",
			err:      ValidationError{Field: "releases[0].version", Message: "invalid version"},
			expected: "releases[0].version: invalid version",
		},
		{
			name:     "without field",
			err:      ValidationError{Message: "general error"},
			expected: "general error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	ve := ValidationError{
		Field:   "test",
		Message: "test message",
		Err:     innerErr,
	}

	if unwrapped := ve.Unwrap(); unwrapped != innerErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}

	// Test with nil error
	ve2 := ValidationError{Field: "test", Message: "test"}
	if unwrapped := ve2.Unwrap(); unwrapped != nil {
		t.Errorf("Unwrap() = %v, want nil", unwrapped)
	}
}

func TestValidate_InvalidIRVersion(t *testing.T) {
	cl := &Changelog{
		IRVersion: "2.0",
		Project:   "test",
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for wrong IR version")
	}
	if !hasError(result.Errors, ErrInvalidIRVersion) {
		t.Error("expected ErrInvalidIRVersion")
	}
}

func TestValidate_EmptyDescription(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []Entry{{Description: ""}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for empty description")
	}
	if !hasError(result.Errors, ErrEmptyDescription) {
		t.Error("expected ErrEmptyDescription")
	}
}

func TestValidate_InvalidSeverity(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", Severity: "super-critical"}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for bad severity")
	}
	if !hasError(result.Errors, ErrInvalidSeverity) {
		t.Error("expected ErrInvalidSeverity")
	}
}

func TestValidate_InvalidCVSSScore(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", CVSSScore: 11.0}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for CVSS score > 10")
	}
	if !hasError(result.Errors, ErrInvalidCVSSScore) {
		t.Error("expected ErrInvalidCVSSScore")
	}
}

func TestValidate_InvalidCVSSScore_Negative(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "Fix", CVSSScore: -1.0}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for negative CVSS score")
	}
}

func TestValidate_Unreleased(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Unreleased: &Release{
			Added: []Entry{{Description: "New feature"}},
		},
	}

	result := cl.Validate()
	if !result.Valid {
		t.Errorf("expected valid changelog, got errors: %v", result.Errors)
	}
}

func TestValidate_UnreleasedWithEmptyDescription(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Unreleased: &Release{
			Added: []Entry{{Description: ""}},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid for unreleased with empty description")
	}
}

func TestValidate_MissingVersion(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{Date: "2026-01-03"},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for missing version")
	}
}

func TestValidate_MissingDate(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{Version: "1.0.0"},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid changelog for missing date")
	}
}

func TestValidate_AllCategories(t *testing.T) {
	// Test that all 19 categories are validated
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:        "1.0.0",
				Date:           "2026-01-03",
				Highlights:     []Entry{{Description: "h"}},
				Breaking:       []Entry{{Description: "b"}},
				UpgradeGuide:   []Entry{{Description: "u"}},
				Security:       []Entry{{Description: "s"}},
				Added:          []Entry{{Description: "a"}},
				Changed:        []Entry{{Description: "c"}},
				Deprecated:     []Entry{{Description: "d"}},
				Removed:        []Entry{{Description: "r"}},
				Fixed:          []Entry{{Description: "f"}},
				Performance:    []Entry{{Description: "p"}},
				Dependencies:   []Entry{{Description: "dep"}},
				Documentation:  []Entry{{Description: "doc"}},
				Build:          []Entry{{Description: "bld"}},
				Infrastructure: []Entry{{Description: "i"}},
				Observability:  []Entry{{Description: "o"}},
				Compliance:     []Entry{{Description: "comp"}},
				Internal:       []Entry{{Description: "int"}},
				KnownIssues:    []Entry{{Description: "k"}},
				Contributors:   []Entry{{Description: "cont"}},
			},
		},
	}

	result := cl.Validate()
	if !result.Valid {
		t.Errorf("expected valid changelog with all categories, got errors: %v", result.Errors)
	}
}

func TestValidate_SecurityEntryWithEmptyDescription(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Security: []Entry{{Description: "", CVE: "CVE-2026-12345"}},
			},
		},
	}

	result := cl.Validate()
	if result.Valid {
		t.Error("expected invalid for security entry with empty description")
	}
}

func TestValidate_ValidSeverities(t *testing.T) {
	validSeverities := []string{"critical", "high", "medium", "low", "informational"}

	for _, severity := range validSeverities {
		t.Run(severity, func(t *testing.T) {
			cl := &Changelog{
				IRVersion: "1.0",
				Project:   "test",
				Releases: []Release{
					{
						Version:  "1.0.0",
						Date:     "2026-01-03",
						Security: []Entry{{Description: "Fix", Severity: severity}},
					},
				},
			}

			result := cl.Validate()
			if !result.Valid {
				t.Errorf("expected valid changelog for severity %q, got errors: %v", severity, result.Errors)
			}
		})
	}
}

func TestValidateMinTier_Valid(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []Entry{{Description: "New feature"}},
			},
		},
	}

	err := cl.ValidateMinTier(TierCore)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateMinTier_NoReleases(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases:  []Release{},
	}

	err := cl.ValidateMinTier(TierCore)
	if err != nil {
		t.Errorf("expected no error for empty releases, got %v", err)
	}
}

func TestValidateMinTier_NoEntriesAtTier(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version:  "1.0.0",
				Date:     "2026-01-03",
				Internal: []Entry{{Description: "Internal change"}}, // optional tier only
			},
		},
	}

	err := cl.ValidateMinTier(TierCore)
	if err == nil {
		t.Error("expected error for no entries at core tier")
	}
	if !errors.Is(err, ErrNoEntriesAtTier) {
		t.Errorf("expected ErrNoEntriesAtTier, got %v", err)
	}
}

func TestValidateMinTier_InvalidTier(t *testing.T) {
	cl := &Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []Entry{{Description: "New feature"}},
			},
		},
	}

	err := cl.ValidateMinTier(Tier("invalid"))
	if err == nil {
		t.Error("expected error for invalid tier")
	}
	if !errors.Is(err, ErrInvalidTier) {
		t.Errorf("expected ErrInvalidTier, got %v", err)
	}
}

func TestParseTier_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected Tier
	}{
		{"core", TierCore},
		{"CORE", TierCore},
		{"Core", TierCore},
		{"standard", TierStandard},
		{"extended", TierExtended},
		{"optional", TierOptional},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tier, err := ParseTier(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tier != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tier)
			}
		})
	}
}

func TestParseTier_Invalid(t *testing.T) {
	_, err := ParseTier("invalid")
	if err == nil {
		t.Error("expected error for invalid tier")
	}
	if !errors.Is(err, ErrInvalidTier) {
		t.Errorf("expected ErrInvalidTier, got %v", err)
	}
}
