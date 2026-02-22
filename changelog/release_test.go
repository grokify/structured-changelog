package changelog

import (
	"testing"
)

func TestNewRelease(t *testing.T) {
	r := NewRelease("1.0.0", "2026-01-04")
	if r.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", r.Version)
	}
	if r.Date != "2026-01-04" {
		t.Errorf("expected date '2026-01-04', got %q", r.Date)
	}
}

func TestReleaseIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		release  Release
		expected bool
	}{
		{"empty release", Release{}, true},
		{"with added", Release{Added: []Entry{{Description: "test"}}}, false},
		{"with changed", Release{Changed: []Entry{{Description: "test"}}}, false},
		{"with deprecated", Release{Deprecated: []Entry{{Description: "test"}}}, false},
		{"with removed", Release{Removed: []Entry{{Description: "test"}}}, false},
		{"with fixed", Release{Fixed: []Entry{{Description: "test"}}}, false},
		{"with security", Release{Security: []Entry{{Description: "test"}}}, false},
		{"with highlights", Release{Highlights: []Entry{{Description: "test"}}}, false},
		{"with breaking", Release{Breaking: []Entry{{Description: "test"}}}, false},
		{"with upgrade_guide", Release{UpgradeGuide: []Entry{{Description: "test"}}}, false},
		{"with performance", Release{Performance: []Entry{{Description: "test"}}}, false},
		{"with dependencies", Release{Dependencies: []Entry{{Description: "test"}}}, false},
		{"with documentation", Release{Documentation: []Entry{{Description: "test"}}}, false},
		{"with build", Release{Build: []Entry{{Description: "test"}}}, false},
		{"with infrastructure", Release{Infrastructure: []Entry{{Description: "test"}}}, false},
		{"with observability", Release{Observability: []Entry{{Description: "test"}}}, false},
		{"with compliance", Release{Compliance: []Entry{{Description: "test"}}}, false},
		{"with internal", Release{Internal: []Entry{{Description: "test"}}}, false},
		{"with known_issues", Release{KnownIssues: []Entry{{Description: "test"}}}, false},
		{"with contributors", Release{Contributors: []Entry{{Description: "test"}}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.release.IsEmpty(); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReleaseCategories(t *testing.T) {
	r := Release{
		Added:   []Entry{{Description: "added"}},
		Fixed:   []Entry{{Description: "fixed"}},
		Changed: []Entry{{Description: "changed"}},
	}

	cats := r.Categories()
	if len(cats) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(cats))
	}

	// Verify canonical order: Added comes before Changed, Changed before Fixed
	expectedOrder := []string{"Added", "Changed", "Fixed"}
	for i, cat := range cats {
		if cat.Name != expectedOrder[i] {
			t.Errorf("category %d: expected %q, got %q", i, expectedOrder[i], cat.Name)
		}
	}
}

func TestReleaseCategoriesFiltered(t *testing.T) {
	r := Release{
		Security:      []Entry{{Description: "security"}},      // core
		Added:         []Entry{{Description: "added"}},         // core
		Performance:   []Entry{{Description: "performance"}},   // standard
		Documentation: []Entry{{Description: "documentation"}}, // extended
		Internal:      []Entry{{Description: "internal"}},      // optional
	}

	tests := []struct {
		tier      Tier
		wantCount int
	}{
		{TierCore, 2},     // Security, Added
		{TierStandard, 3}, // Security, Added, Performance
		{TierExtended, 4}, // Security, Added, Performance, Documentation
		{TierOptional, 5}, // All
	}

	for _, tt := range tests {
		t.Run(string(tt.tier), func(t *testing.T) {
			cats := r.CategoriesFiltered(tt.tier)
			if len(cats) != tt.wantCount {
				t.Errorf("CategoriesFiltered(%q) returned %d categories, want %d", tt.tier, len(cats), tt.wantCount)
			}
		})
	}
}

func TestReleaseCategoryMap(t *testing.T) {
	r := Release{
		Added:    []Entry{{Description: "added1"}, {Description: "added2"}},
		Security: []Entry{{Description: "security1"}},
	}

	m := r.categoryMap()

	if len(m["Added"]) != 2 {
		t.Errorf("expected 2 Added entries, got %d", len(m["Added"]))
	}
	if len(m["Security"]) != 1 {
		t.Errorf("expected 1 Security entry, got %d", len(m["Security"]))
	}
	if len(m["Fixed"]) != 0 {
		t.Errorf("expected 0 Fixed entries, got %d", len(m["Fixed"]))
	}
}

func TestReleaseGetEntries(t *testing.T) {
	r := Release{
		Added: []Entry{{Description: "entry1"}, {Description: "entry2"}},
	}

	entries := r.GetEntries("Added")
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	entries = r.GetEntries("NonExistent")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for non-existent category, got %d", len(entries))
	}
}

func TestReleaseAddMethods(t *testing.T) {
	r := Release{}
	e := Entry{Description: "test"}

	// Test all Add methods
	r.AddHighlights(e)
	if len(r.Highlights) != 1 {
		t.Error("AddHighlights failed")
	}

	r.AddBreaking(e)
	if len(r.Breaking) != 1 {
		t.Error("AddBreaking failed")
	}

	r.AddUpgradeGuide(e)
	if len(r.UpgradeGuide) != 1 {
		t.Error("AddUpgradeGuide failed")
	}

	r.AddSecurity(e)
	if len(r.Security) != 1 {
		t.Error("AddSecurity failed")
	}

	r.AddAdded(e)
	if len(r.Added) != 1 {
		t.Error("AddAdded failed")
	}

	r.AddChanged(e)
	if len(r.Changed) != 1 {
		t.Error("AddChanged failed")
	}

	r.AddDeprecated(e)
	if len(r.Deprecated) != 1 {
		t.Error("AddDeprecated failed")
	}

	r.AddRemoved(e)
	if len(r.Removed) != 1 {
		t.Error("AddRemoved failed")
	}

	r.AddFixed(e)
	if len(r.Fixed) != 1 {
		t.Error("AddFixed failed")
	}

	r.AddPerformance(e)
	if len(r.Performance) != 1 {
		t.Error("AddPerformance failed")
	}

	r.AddDependencies(e)
	if len(r.Dependencies) != 1 {
		t.Error("AddDependencies failed")
	}

	r.AddDocumentation(e)
	if len(r.Documentation) != 1 {
		t.Error("AddDocumentation failed")
	}

	r.AddBuild(e)
	if len(r.Build) != 1 {
		t.Error("AddBuild failed")
	}

	r.AddInfrastructure(e)
	if len(r.Infrastructure) != 1 {
		t.Error("AddInfrastructure failed")
	}

	r.AddObservability(e)
	if len(r.Observability) != 1 {
		t.Error("AddObservability failed")
	}

	r.AddCompliance(e)
	if len(r.Compliance) != 1 {
		t.Error("AddCompliance failed")
	}

	r.AddInternal(e)
	if len(r.Internal) != 1 {
		t.Error("AddInternal failed")
	}

	r.AddKnownIssues(e)
	if len(r.KnownIssues) != 1 {
		t.Error("AddKnownIssues failed")
	}

	r.AddContributors(e)
	if len(r.Contributors) != 1 {
		t.Error("AddContributors failed")
	}

	// Verify release is not empty
	if r.IsEmpty() {
		t.Error("release should not be empty after adding entries")
	}
}

func TestReleaseCategoriesCanonicalOrder(t *testing.T) {
	// Create a release with entries in all categories
	r := Release{
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
	}

	cats := r.Categories()
	if len(cats) != 19 {
		t.Fatalf("expected 19 categories, got %d", len(cats))
	}

	// Verify canonical order
	expectedOrder := []string{
		"Highlights", "Breaking", "Upgrade Guide", "Security",
		"Added", "Changed", "Deprecated", "Removed", "Fixed",
		"Performance", "Dependencies",
		"Documentation", "Build",
		"Infrastructure", "Observability", "Compliance",
		"Internal",
		"Known Issues", "Contributors",
	}

	for i, cat := range cats {
		if cat.Name != expectedOrder[i] {
			t.Errorf("category %d: expected %q, got %q", i, expectedOrder[i], cat.Name)
		}
	}
}

func TestCategoryStruct(t *testing.T) {
	cat := Category{
		Name:    "Added",
		Entries: []Entry{{Description: "entry1"}, {Description: "entry2"}},
	}

	if cat.Name != "Added" {
		t.Errorf("expected name 'Added', got %q", cat.Name)
	}
	if len(cat.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(cat.Entries))
	}
}

func TestIsMaintenanceOnly(t *testing.T) {
	tests := []struct {
		name     string
		release  Release
		expected bool
	}{
		{
			name:     "empty release",
			release:  Release{},
			expected: false, // empty is not maintenance
		},
		{
			name:     "dependencies only",
			release:  Release{Dependencies: []Entry{{Description: "bump lib"}}},
			expected: true,
		},
		{
			name:     "documentation only",
			release:  Release{Documentation: []Entry{{Description: "update docs"}}},
			expected: true,
		},
		{
			name:     "build only",
			release:  Release{Build: []Entry{{Description: "fix ci"}}},
			expected: true,
		},
		{
			name:     "tests only",
			release:  Release{Tests: []Entry{{Description: "add tests"}}},
			expected: true,
		},
		{
			name:     "internal only",
			release:  Release{Internal: []Entry{{Description: "refactor"}}},
			expected: true,
		},
		{
			name:     "multiple maintenance types",
			release:  Release{Dependencies: []Entry{{Description: "bump"}}, Documentation: []Entry{{Description: "docs"}}},
			expected: true,
		},
		{
			name:     "has added - not maintenance",
			release:  Release{Added: []Entry{{Description: "new feature"}}, Dependencies: []Entry{{Description: "bump"}}},
			expected: false,
		},
		{
			name:     "has changed - not maintenance",
			release:  Release{Changed: []Entry{{Description: "change"}}, Documentation: []Entry{{Description: "docs"}}},
			expected: false,
		},
		{
			name:     "has fixed - not maintenance",
			release:  Release{Fixed: []Entry{{Description: "fix"}}},
			expected: false,
		},
		{
			name:     "has security - not maintenance",
			release:  Release{Security: []Entry{{Description: "security fix"}}},
			expected: false,
		},
		{
			name:     "has removed - not maintenance",
			release:  Release{Removed: []Entry{{Description: "removed"}}},
			expected: false,
		},
		{
			name:     "has deprecated - not maintenance",
			release:  Release{Deprecated: []Entry{{Description: "deprecated"}}},
			expected: false,
		},
		{
			name:     "has highlights - not maintenance",
			release:  Release{Highlights: []Entry{{Description: "highlight"}}},
			expected: false,
		},
		{
			name:     "has breaking - not maintenance",
			release:  Release{Breaking: []Entry{{Description: "breaking"}}},
			expected: false,
		},
		{
			name:     "has performance - not maintenance",
			release:  Release{Performance: []Entry{{Description: "perf"}}},
			expected: false,
		},
		{
			name:     "infrastructure only - maintenance",
			release:  Release{Infrastructure: []Entry{{Description: "infra"}}},
			expected: true,
		},
		{
			name:     "observability only - maintenance",
			release:  Release{Observability: []Entry{{Description: "metrics"}}},
			expected: true,
		},
		{
			name:     "compliance only - maintenance",
			release:  Release{Compliance: []Entry{{Description: "audit"}}},
			expected: true,
		},
		{
			name:     "contributors only - maintenance",
			release:  Release{Contributors: []Entry{{Description: "thanks"}}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.release.IsMaintenanceOnly()
			if got != tt.expected {
				t.Errorf("IsMaintenanceOnly() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAddTests(t *testing.T) {
	r := Release{}
	e := Entry{Description: "add unit tests"}

	r.AddTests(e)

	if len(r.Tests) != 1 {
		t.Errorf("expected 1 test entry, got %d", len(r.Tests))
	}
	if r.Tests[0].Description != "add unit tests" {
		t.Errorf("expected description 'add unit tests', got %q", r.Tests[0].Description)
	}
}

func TestHasCategory(t *testing.T) {
	r := Release{
		Added:        []Entry{{Description: "added"}},
		Security:     []Entry{{Description: "security"}},
		Dependencies: []Entry{{Description: "deps"}},
	}

	tests := []struct {
		category string
		want     bool
	}{
		{"Added", true},
		{"Security", true},
		{"Dependencies", true},
		{"Fixed", false},
		{"Changed", false},
		{"NonExistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			got := r.HasCategory(tt.category)
			if got != tt.want {
				t.Errorf("HasCategory(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}

func TestIsNotable(t *testing.T) {
	tests := []struct {
		name    string
		release Release
		policy  *NotabilityPolicy
		want    bool
	}{
		{
			name:    "empty release - nil policy",
			release: Release{},
			policy:  nil,
			want:    false, // Empty release is never notable
		},
		{
			name:    "empty release - default policy",
			release: Release{},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "has notable category - nil policy",
			release: Release{Added: []Entry{{Description: "new feature"}}},
			policy:  nil,
			want:    true,
		},
		{
			name:    "has notable category - default policy",
			release: Release{Added: []Entry{{Description: "new feature"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "has security - notable",
			release: Release{Security: []Entry{{Description: "fix CVE"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "has fixed - notable",
			release: Release{Fixed: []Entry{{Description: "bug fix"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "has performance - notable",
			release: Release{Performance: []Entry{{Description: "faster"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "dependencies only - not notable",
			release: Release{Dependencies: []Entry{{Description: "bump lib"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "documentation only - not notable",
			release: Release{Documentation: []Entry{{Description: "update docs"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "build only - not notable",
			release: Release{Build: []Entry{{Description: "fix ci"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "tests only - not notable",
			release: Release{Tests: []Entry{{Description: "add tests"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "internal only - not notable",
			release: Release{Internal: []Entry{{Description: "refactor"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "mixed notable and non-notable - notable",
			release: Release{Added: []Entry{{Description: "feature"}}, Dependencies: []Entry{{Description: "bump"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "infrastructure only - not notable",
			release: Release{Infrastructure: []Entry{{Description: "infra"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "contributors only - not notable",
			release: Release{Contributors: []Entry{{Description: "thanks"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    false,
		},
		{
			name:    "custom policy - only security",
			release: Release{Added: []Entry{{Description: "feature"}}},
			policy:  NewNotabilityPolicy([]string{"Security"}),
			want:    false, // Added is not in custom policy
		},
		{
			name:    "custom policy - matches",
			release: Release{Security: []Entry{{Description: "fix"}}},
			policy:  NewNotabilityPolicy([]string{"Security"}),
			want:    true,
		},
		{
			name:    "empty policy - all notable",
			release: Release{Dependencies: []Entry{{Description: "deps"}}},
			policy:  NewNotabilityPolicy([]string{}),
			want:    true, // Empty policy means all are notable
		},
		{
			name:    "highlights - notable",
			release: Release{Highlights: []Entry{{Description: "summary"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "breaking - notable",
			release: Release{Breaking: []Entry{{Description: "breaking change"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
		{
			name:    "known issues - notable",
			release: Release{KnownIssues: []Entry{{Description: "known issue"}}},
			policy:  DefaultNotabilityPolicy(),
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.release.IsNotable(tt.policy)
			if got != tt.want {
				t.Errorf("IsNotable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultNotabilityPolicy(t *testing.T) {
	policy := DefaultNotabilityPolicy()

	// Check that notable categories are included
	notableCategories := []string{
		"Highlights", "Breaking", "Upgrade Guide", "Security",
		"Added", "Changed", "Deprecated", "Removed", "Fixed",
		"Performance", "Known Issues",
	}
	for _, cat := range notableCategories {
		if !policy.IsNotable(cat) {
			t.Errorf("expected %q to be notable", cat)
		}
	}

	// Check that maintenance categories are NOT notable
	maintenanceCategories := []string{
		"Dependencies", "Documentation", "Build", "Tests",
		"Infrastructure", "Observability", "Compliance",
		"Internal", "Contributors",
	}
	for _, cat := range maintenanceCategories {
		if policy.IsNotable(cat) {
			t.Errorf("expected %q to NOT be notable", cat)
		}
	}
}

func TestNotabilityPolicy_IsNotable(t *testing.T) {
	tests := []struct {
		name     string
		policy   *NotabilityPolicy
		category string
		want     bool
	}{
		{
			name:     "nil policy - all notable",
			policy:   nil,
			category: "Dependencies",
			want:     true,
		},
		{
			name:     "empty categories - all notable",
			policy:   &NotabilityPolicy{NotableCategories: []string{}},
			category: "Dependencies",
			want:     true,
		},
		{
			name:     "category in list - notable",
			policy:   &NotabilityPolicy{NotableCategories: []string{"Added", "Fixed"}},
			category: "Added",
			want:     true,
		},
		{
			name:     "category not in list - not notable",
			policy:   &NotabilityPolicy{NotableCategories: []string{"Added", "Fixed"}},
			category: "Security",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.policy.IsNotable(tt.category)
			if got != tt.want {
				t.Errorf("IsNotable(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}
