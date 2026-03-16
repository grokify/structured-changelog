package aggregate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/structured-changelog/changelog"
)

func TestLoadPortfolioFromPaths(t *testing.T) {
	dir := t.TempDir()

	// Create test changelog
	cl := &changelog.Changelog{
		IRVersion: changelog.IRVersion,
		Project:   "test-project",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2024-01-15",
				Added:   []changelog.Entry{{Description: "Feature 1"}},
			},
		},
	}

	changelogPath := filepath.Join(dir, "CHANGELOG.json")
	if err := cl.WriteFile(changelogPath); err != nil {
		t.Fatalf("failed to write changelog: %v", err)
	}

	// Load portfolio
	portfolio, err := LoadPortfolioFromPaths("Test Portfolio", []string{changelogPath})
	if err != nil {
		t.Fatalf("LoadPortfolioFromPaths() error: %v", err)
	}

	if portfolio.Name != "Test Portfolio" {
		t.Errorf("expected name %q, got %q", "Test Portfolio", portfolio.Name)
	}

	if len(portfolio.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(portfolio.Projects))
	}

	if portfolio.Projects[0].Name != "test-project" {
		t.Errorf("expected project name %q, got %q", "test-project", portfolio.Projects[0].Name)
	}

	if portfolio.DateRange.Start != "2024-01-15" {
		t.Errorf("expected date range start %q, got %q", "2024-01-15", portfolio.DateRange.Start)
	}
}

func TestPortfolioSummary(t *testing.T) {
	portfolio := &Portfolio{
		Name: "Test",
		Projects: []ProjectData{
			{
				Path: "test/repo",
				Name: "Test Repo",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{
							Version: "1.0.0",
							Date:    "2024-01-15",
							Added:   []changelog.Entry{{Description: "A"}, {Description: "B"}},
							Fixed:   []changelog.Entry{{Description: "C"}},
						},
						{
							Version: "0.9.0",
							Date:    "2024-01-01",
							Added:   []changelog.Entry{{Description: "D"}},
						},
					},
				},
			},
		},
	}

	summary := portfolio.Summary()

	if summary.ProjectCount != 1 {
		t.Errorf("expected 1 project, got %d", summary.ProjectCount)
	}

	if summary.ReleaseCount != 2 {
		t.Errorf("expected 2 releases, got %d", summary.ReleaseCount)
	}

	if summary.EntryCount != 4 {
		t.Errorf("expected 4 entries, got %d", summary.EntryCount)
	}

	if summary.ByCategory["Added"] != 3 {
		t.Errorf("expected 3 Added entries, got %d", summary.ByCategory["Added"])
	}

	if summary.ByCategory["Fixed"] != 1 {
		t.Errorf("expected 1 Fixed entry, got %d", summary.ByCategory["Fixed"])
	}
}

func TestPortfolioAllReleases(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "1.0.0", Date: "2024-01-15"},
						{Version: "0.9.0", Date: "2024-01-01"},
					},
				},
			},
			{
				Path: "repo2",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "2.0.0", Date: "2024-01-10"},
					},
				},
			},
		},
	}

	releases := portfolio.AllReleases()

	if len(releases) != 3 {
		t.Fatalf("expected 3 releases, got %d", len(releases))
	}

	// Should be sorted by date descending
	expectedOrder := []string{"2024-01-15", "2024-01-10", "2024-01-01"}
	for i, r := range releases {
		if r.Release.Date != expectedOrder[i] {
			t.Errorf("release[%d] expected date %q, got %q", i, expectedOrder[i], r.Release.Date)
		}
	}
}

func TestPortfolioFilterByDateRange(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "1.0.0", Date: "2024-03-15"},
						{Version: "0.9.0", Date: "2024-01-15"},
						{Version: "0.8.0", Date: "2023-12-01"},
					},
				},
			},
		},
	}

	filtered := portfolio.FilterByDateRange("2024-01-01", "2024-02-28")

	if len(filtered) != 1 {
		t.Errorf("expected 1 release in range, got %d", len(filtered))
	}

	if filtered[0].Release.Version != "0.9.0" {
		t.Errorf("expected version 0.9.0, got %s", filtered[0].Release.Version)
	}
}

func TestPortfolioWriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portfolio.json")

	portfolio := &Portfolio{
		Name: "Test Portfolio",
		Projects: []ProjectData{
			{
				Path: "test/repo",
				Name: "Test Repo",
				Changelog: &changelog.Changelog{
					IRVersion: changelog.IRVersion,
					Project:   "Test Repo",
				},
			},
		},
		DateRange: DateRange{Start: "2024-01-01", End: "2024-12-31"},
	}

	// Write
	if err := portfolio.WriteFile(path); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}

	// Load
	loaded, err := LoadPortfolioFile(path)
	if err != nil {
		t.Fatalf("LoadPortfolioFile error: %v", err)
	}

	if loaded.Name != portfolio.Name {
		t.Errorf("expected name %q, got %q", portfolio.Name, loaded.Name)
	}

	if len(loaded.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(loaded.Projects))
	}
}

func TestPortfolioGetProject(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{Path: "repo1", Name: "Repo 1"},
			{Path: "repo2", Name: "Repo 2"},
		},
	}

	// Found
	p := portfolio.GetProject("repo1")
	if p == nil {
		t.Fatal("expected to find repo1")
	}
	if p.Name != "Repo 1" {
		t.Errorf("expected name %q, got %q", "Repo 1", p.Name)
	}

	// Not found
	p = portfolio.GetProject("nonexistent")
	if p != nil {
		t.Error("expected nil for nonexistent project")
	}
}

func TestPortfolioProjectNames(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{Path: "repo1", Name: "Project A"},
			{Path: "repo2", Name: "Project B"},
		},
	}

	names := portfolio.ProjectNames()

	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}

	if names[0] != "Project A" || names[1] != "Project B" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestPortfolioProjectPaths(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{Path: "github.com/org/repo1", Name: "Repo 1"},
			{Path: "github.com/org/repo2", Name: "Repo 2"},
			{Path: "github.com/other/repo3", Name: "Repo 3"},
		},
	}

	paths := portfolio.ProjectPaths()

	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(paths))
	}

	expectedPaths := []string{
		"github.com/org/repo1",
		"github.com/org/repo2",
		"github.com/other/repo3",
	}

	for i, path := range paths {
		if path != expectedPaths[i] {
			t.Errorf("path[%d]: expected %q, got %q", i, expectedPaths[i], path)
		}
	}
}

func TestPortfolioJSON(t *testing.T) {
	portfolio := &Portfolio{
		Name: "Test Portfolio",
		Projects: []ProjectData{
			{
				Path: "test/repo",
				Name: "Test Repo",
				Changelog: &changelog.Changelog{
					IRVersion: changelog.IRVersion,
					Project:   "Test Repo",
				},
			},
		},
	}

	data, err := portfolio.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty JSON output")
	}

	// Verify it's valid JSON by parsing it back
	var parsed Portfolio
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if parsed.Name != portfolio.Name {
		t.Errorf("expected name %q, got %q", portfolio.Name, parsed.Name)
	}
}

func TestParsePortfolio(t *testing.T) {
	jsonData := []byte(`{
		"name": "Parsed Portfolio",
		"description": "Test description",
		"projects": [
			{"path": "repo1", "name": "Repo 1"}
		],
		"dateRange": {"start": "2024-01-01", "end": "2024-12-31"}
	}`)

	portfolio, err := ParsePortfolio(jsonData)
	if err != nil {
		t.Fatalf("ParsePortfolio() error: %v", err)
	}

	if portfolio.Name != "Parsed Portfolio" {
		t.Errorf("expected name %q, got %q", "Parsed Portfolio", portfolio.Name)
	}

	if portfolio.Description != "Test description" {
		t.Errorf("expected description %q, got %q", "Test description", portfolio.Description)
	}

	if len(portfolio.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(portfolio.Projects))
	}
}

func TestParsePortfolioInvalid(t *testing.T) {
	invalidJSON := []byte(`{invalid json`)

	_, err := ParsePortfolio(invalidJSON)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoadPortfolioFile404(t *testing.T) {
	_, err := LoadPortfolioFile("/nonexistent/path/portfolio.json")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestPortfolioEmptyProjects(t *testing.T) {
	portfolio := &Portfolio{
		Name:     "Empty Portfolio",
		Projects: []ProjectData{},
	}

	paths := portfolio.ProjectPaths()
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}

	names := portfolio.ProjectNames()
	if len(names) != 0 {
		t.Errorf("expected 0 names, got %d", len(names))
	}

	summary := portfolio.Summary()
	if summary.ProjectCount != 0 {
		t.Errorf("expected ProjectCount=0, got %d", summary.ProjectCount)
	}
}

func TestPortfolioAllReleasesEmpty(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{},
				},
			},
		},
	}

	releases := portfolio.AllReleases()
	if len(releases) != 0 {
		t.Errorf("expected 0 releases, got %d", len(releases))
	}
}

func TestPortfolioFilterByDateRangeEmpty(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "1.0.0", Date: "2024-06-15"},
					},
				},
			},
		},
	}

	// Filter with range that excludes all releases
	filtered := portfolio.FilterByDateRange("2025-01-01", "2025-12-31")
	if len(filtered) != 0 {
		t.Errorf("expected 0 releases in range, got %d", len(filtered))
	}
}

func TestPortfolioSummaryEmpty(t *testing.T) {
	portfolio := &Portfolio{
		Projects: []ProjectData{},
	}

	summary := portfolio.Summary()

	if summary.ProjectCount != 0 {
		t.Errorf("expected ProjectCount=0, got %d", summary.ProjectCount)
	}

	if summary.ReleaseCount != 0 {
		t.Errorf("expected ReleaseCount=0, got %d", summary.ReleaseCount)
	}

	if summary.EntryCount != 0 {
		t.Errorf("expected EntryCount=0, got %d", summary.EntryCount)
	}
}
