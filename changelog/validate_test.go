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
