package gitlog

import (
	"sort"
	"testing"
)

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		// Equal versions
		{"equal_v_prefix", "v1.0.0", "v1.0.0", 0},
		{"equal_no_prefix", "1.0.0", "1.0.0", 0},
		{"equal_mixed_prefix", "v1.0.0", "1.0.0", 0},

		// Major version differences
		{"major_less", "v1.0.0", "v2.0.0", -1},
		{"major_greater", "v2.0.0", "v1.0.0", 1},
		{"major_double_digit", "v9.0.0", "v10.0.0", -1},

		// Minor version differences
		{"minor_less", "v1.1.0", "v1.2.0", -1},
		{"minor_greater", "v1.3.0", "v1.2.0", 1},
		{"minor_double_digit", "v1.9.0", "v1.10.0", -1},

		// Patch version differences
		{"patch_less", "v1.0.1", "v1.0.2", -1},
		{"patch_greater", "v1.0.3", "v1.0.2", 1},
		{"patch_double_digit", "v1.0.9", "v1.0.10", -1},

		// Prerelease versions (regex only captures MAJOR.MINOR.PATCH)
		{"prerelease_same_base", "v1.0.0-alpha", "v1.0.0-beta", 0},
		{"prerelease_vs_release", "v1.0.0-alpha", "v1.0.0", 0},

		// Complex versions
		{"complex_1", "v0.1.0", "v0.2.0", -1},
		{"complex_2", "v2.26.9", "v2.27.0", -1},
		{"complex_3", "v0.73.5", "v0.73.6", -1},

		// Edge cases with non-semver (falls back to string compare)
		{"non_semver_alpha", "alpha", "beta", -1},
		{"non_semver_empty", "", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareSemver(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareSemver(%q, %q) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareSemverSorting(t *testing.T) {
	// Test that sorting produces correct order
	versions := []string{
		"v2.0.0",
		"v0.1.0",
		"v1.10.0",
		"v1.9.0",
		"v1.0.0",
		"v10.0.0",
		"v0.10.0",
		"v0.9.0",
	}

	expected := []string{
		"v0.1.0",
		"v0.9.0",
		"v0.10.0",
		"v1.0.0",
		"v1.9.0",
		"v1.10.0",
		"v2.0.0",
		"v10.0.0",
	}

	sort.Slice(versions, func(i, j int) bool {
		return compareSemver(versions[i], versions[j]) < 0
	})

	for i, v := range versions {
		if v != expected[i] {
			t.Errorf("position %d: got %q, expected %q", i, v, expected[i])
		}
	}
}

func TestSemverRegex(t *testing.T) {
	tests := []struct {
		name    string
		version string
		matches bool
	}{
		{"valid_v_prefix", "v1.0.0", true},
		{"valid_no_prefix", "1.0.0", true},
		{"valid_prerelease", "v1.0.0-alpha", true},
		{"valid_prerelease_number", "v1.0.0-beta.1", true},
		{"valid_build_metadata", "v1.0.0+build123", true},
		{"valid_full", "v1.0.0-alpha.1+build.123", true},
		{"valid_double_digit", "v10.20.30", true},
		{"valid_zero", "v0.0.0", true},

		// Invalid formats
		{"invalid_too_few_parts", "v1.0", false},
		{"invalid_four_parts", "v1.0.0.0", true}, // regex matches first 3 parts
		{"invalid_alpha_only", "alpha", false},
		{"invalid_empty", "", false},
		{"invalid_vv_prefix", "vv1.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := semverRegex.MatchString(tt.version)
			if result != tt.matches {
				t.Errorf("semverRegex.MatchString(%q) = %v, expected %v", tt.version, result, tt.matches)
			}
		})
	}
}

func TestTagListStructure(t *testing.T) {
	tagList := &TagList{
		Repository: "github.com/example/repo",
		Tags: []Tag{
			{Name: "v1.0.0", DateString: "2024-01-01", CommitHash: "abc123", IsInitial: true},
			{Name: "v1.1.0", DateString: "2024-02-01", CommitHash: "def456", CommitCount: 5},
		},
		TotalTags: 2,
	}

	if tagList.TotalTags != 2 {
		t.Errorf("expected TotalTags=2, got %d", tagList.TotalTags)
	}

	if !tagList.Tags[0].IsInitial {
		t.Error("expected first tag to be initial")
	}

	if tagList.Tags[1].CommitCount != 5 {
		t.Errorf("expected CommitCount=5, got %d", tagList.Tags[1].CommitCount)
	}
}

func TestVersionRangeStructure(t *testing.T) {
	vr := VersionRange{
		Version: "v1.1.0",
		Since:   "v1.0.0",
		Until:   "v1.1.0",
		Date:    "2024-02-01",
		Commits: 10,
	}

	if vr.Version != "v1.1.0" {
		t.Errorf("expected Version=%q, got %q", "v1.1.0", vr.Version)
	}

	if vr.Since != "v1.0.0" {
		t.Errorf("expected Since=%q, got %q", "v1.0.0", vr.Since)
	}

	if vr.Commits != 10 {
		t.Errorf("expected Commits=10, got %d", vr.Commits)
	}
}
