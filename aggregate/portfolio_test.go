package aggregate

import (
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
