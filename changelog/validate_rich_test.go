package changelog

import (
	"testing"
)

func TestValidateRich_Valid(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2024-01-15",
		Added: []Entry{
			{Description: "Added a new feature"},
		},
	})

	result := cl.ValidateRich()

	if !result.Valid {
		t.Errorf("expected valid result, got errors: %v", result.Errors)
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}
}

func TestValidateRich_MissingProject(t *testing.T) {
	cl := &Changelog{
		IRVersion: IRVersion,
		Project:   "",
	}

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeMissingField && err.Path == "project" {
			found = true
			if err.Suggestion == "" {
				t.Error("expected suggestion for missing project")
			}
		}
	}
	if !found {
		t.Error("expected missing project error")
	}
}

func TestValidateRich_InvalidIRVersion(t *testing.T) {
	cl := &Changelog{
		IRVersion: "0.5",
		Project:   "test",
	}

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeInvalidIRVersion {
			found = true
			if err.Actual != "0.5" {
				t.Errorf("expected actual '0.5', got %q", err.Actual)
			}
			if err.Expected != IRVersion {
				t.Errorf("expected expected %q, got %q", IRVersion, err.Expected)
			}
		}
	}
	if !found {
		t.Error("expected invalid IR version error")
	}
}

func TestValidateRich_InvalidVersion(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "v1.0",
		Date:    "2024-01-15",
	})

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeInvalidVersion {
			found = true
			if err.Suggestion == "" {
				t.Error("expected suggestion for invalid version")
			}
		}
	}
	if !found {
		t.Error("expected invalid version error")
	}
}

func TestValidateRich_InvalidDate(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "January 15, 2024",
	})

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeInvalidDate {
			found = true
			if err.Documentation == "" {
				t.Error("expected documentation link for date error")
			}
		}
	}
	if !found {
		t.Error("expected invalid date error")
	}
}

func TestValidateRich_InvalidCVE(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2024-01-15",
		Security: []Entry{
			{Description: "Security fix", CVE: "CVE2024-1234"},
		},
	})

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeInvalidCVE {
			found = true
			if err.Suggestion == "" {
				t.Error("expected suggestion for invalid CVE")
			}
		}
	}
	if !found {
		t.Error("expected invalid CVE error")
	}
}

func TestValidateRich_InvalidSeverity(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2024-01-15",
		Security: []Entry{
			{Description: "Security fix", Severity: "super-high"},
		},
	})

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeInvalidSeverity {
			found = true
		}
	}
	if !found {
		t.Error("expected invalid severity error")
	}
}

func TestValidateRich_Warnings(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2024-01-15",
		Security: []Entry{
			{Description: "Security fix without CVE or severity"},
		},
	})

	result := cl.ValidateRich()

	// Should be valid but have warnings
	if !result.Valid {
		t.Error("expected valid result with warnings")
	}

	if len(result.Warnings) < 2 {
		t.Errorf("expected at least 2 warnings (missing CVE and severity), got %d", len(result.Warnings))
	}
}

func TestValidateRich_ShortDescription(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2024-01-15",
		Added: []Entry{
			{Description: "Short"},
		},
	})

	result := cl.ValidateRich()

	// Should be valid but have warning
	if !result.Valid {
		t.Error("expected valid result")
	}

	found := false
	for _, warn := range result.Warnings {
		if warn.Code == WarnCodeShortDescription {
			found = true
		}
	}
	if !found {
		t.Error("expected short description warning")
	}
}

func TestValidateRich_Summary(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2024-01-15",
		Added:   []Entry{{Description: "Feature one"}, {Description: "Feature two"}},
		Fixed:   []Entry{{Description: "Bug fix"}},
	})
	cl.AddRelease(Release{
		Version: "0.9.0",
		Date:    "2024-01-01",
		Added:   []Entry{{Description: "Initial feature"}},
	})

	result := cl.ValidateRich()

	if result.Summary.ReleaseCount != 2 {
		t.Errorf("expected 2 releases, got %d", result.Summary.ReleaseCount)
	}
	if result.Summary.EntriesCount != 4 {
		t.Errorf("expected 4 entries, got %d", result.Summary.EntriesCount)
	}
}

func TestValidateRich_DuplicateVersion(t *testing.T) {
	cl := New("test-project")
	cl.Releases = []Release{
		{Version: "1.0.0", Date: "2024-01-15"},
		{Version: "1.0.0", Date: "2024-01-10"},
	}

	result := cl.ValidateRich()

	if result.Valid {
		t.Error("expected invalid result for duplicate versions")
	}

	found := false
	for _, err := range result.Errors {
		if err.Code == ErrCodeDuplicateVersion {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate version error")
	}
}

func TestSuggestVersionFix(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"v1.0.0", "Remove the 'v' prefix"},
		{"1.0", "Add patch version"},
		{"abc", "MAJOR.MINOR.PATCH"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := suggestVersionFix(tt.input)
			if result == "" {
				t.Error("expected non-empty suggestion")
			}
		})
	}
}

func TestSuggestDateFix(t *testing.T) {
	tests := []struct {
		input    string
		notEmpty bool
	}{
		{"01/15/2024", true},
		{"January 15, 2024", true},
		{"2024-1-15", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := suggestDateFix(tt.input)
			if tt.notEmpty && result == "" {
				t.Error("expected non-empty suggestion")
			}
		})
	}
}

func TestSuggestSeverityFix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"crit", "critical"},
		{"hi", "high"},
		{"moderate", "medium"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := suggestSeverityFix(tt.input)
			if tt.expected != "" && result == "" {
				t.Errorf("expected suggestion containing %q", tt.expected)
			}
		})
	}
}

func TestRichValidationError_Error(t *testing.T) {
	err := RichValidationError{
		Code:    ErrCodeInvalidDate,
		Path:    "releases[0].date",
		Message: "Invalid date format",
	}

	expected := "[E001] releases[0].date: Invalid date format"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
