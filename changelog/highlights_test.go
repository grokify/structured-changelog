package changelog

import (
	"testing"
	"time"
)

func TestExtractHighlights(t *testing.T) {
	cl := &Changelog{
		Project:    "test-project",
		Repository: "github.com/test/project",
		Releases: []Release{
			{
				Version: "v1.2.0",
				Date:    "2026-05-15",
				Highlights: []Entry{
					{Description: "Major feature X"},
				},
				Added: []Entry{
					{Description: "Added feature Y"},
				},
			},
			{
				Version: "v1.1.0",
				Date:    "2026-04-10",
				Fixed: []Entry{
					{Description: "Fixed bug Z"},
				},
			},
			{
				Version: "v1.0.0",
				Date:    "2026-03-01", // Outside Q2
				Added: []Entry{
					{Description: "Initial release"},
				},
			},
		},
	}

	// Q2 2026: April 1 - July 1
	since := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	result := cl.ExtractHighlights(since, until)

	if result.Project != "test-project" {
		t.Errorf("expected project test-project, got %s", result.Project)
	}

	if len(result.Releases) != 2 {
		t.Fatalf("expected 2 releases in Q2, got %d", len(result.Releases))
	}

	// Should be in original order (reverse chronological)
	if result.Releases[0].Version != "v1.2.0" {
		t.Errorf("expected first release v1.2.0, got %s", result.Releases[0].Version)
	}
	if result.Releases[1].Version != "v1.1.0" {
		t.Errorf("expected second release v1.1.0, got %s", result.Releases[1].Version)
	}

	// v1.2.0 should have highlights and added
	if len(result.Releases[0].Highlights) != 1 {
		t.Errorf("expected 1 highlight in v1.2.0, got %d", len(result.Releases[0].Highlights))
	}
	if len(result.Releases[0].Added) != 1 {
		t.Errorf("expected 1 added in v1.2.0, got %d", len(result.Releases[0].Added))
	}

	// v1.1.0 should have fixed only
	if len(result.Releases[1].Fixed) != 1 {
		t.Errorf("expected 1 fixed in v1.1.0, got %d", len(result.Releases[1].Fixed))
	}
}

func TestExtractHighlightsMaintenanceOnly(t *testing.T) {
	cl := &Changelog{
		Project: "test-project",
		Releases: []Release{
			{
				Version: "v1.0.1",
				Date:    "2026-05-15",
				// Only maintenance categories - should be skipped
				Dependencies: []Entry{
					{Description: "Bump deps"},
				},
			},
		},
	}

	since := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	result := cl.ExtractHighlights(since, until)

	if len(result.Releases) != 0 {
		t.Errorf("expected 0 releases (maintenance-only skipped), got %d", len(result.Releases))
	}
}

func TestReleaseExtractEntryCount(t *testing.T) {
	rel := ReleaseExtract{
		Highlights: []Entry{{Description: "h1"}, {Description: "h2"}},
		Added:      []Entry{{Description: "a1"}},
		Fixed:      []Entry{{Description: "f1"}, {Description: "f2"}, {Description: "f3"}},
	}

	if rel.EntryCount() != 6 {
		t.Errorf("expected 6 entries, got %d", rel.EntryCount())
	}
}

func TestTopProjects(t *testing.T) {
	m := &MultiRepoHighlights{
		ByRepo: []RepoHighlights{
			{
				Project: "project-a",
				Releases: []ReleaseExtract{
					{Added: []Entry{{}, {}, {}}}, // 3
				},
			},
			{
				Project: "project-b",
				Releases: []ReleaseExtract{
					{Added: []Entry{{}, {}, {}, {}, {}}}, // 5
				},
			},
			{
				Project: "project-c",
				Releases: []ReleaseExtract{
					{Added: []Entry{{}}}, // 1
				},
			},
		},
	}

	top := m.TopProjects(2)

	if len(top) != 2 {
		t.Fatalf("expected 2 top projects, got %d", len(top))
	}

	if top[0].Project != "project-b" || top[0].Count != 5 {
		t.Errorf("expected project-b with 5, got %s with %d", top[0].Project, top[0].Count)
	}

	if top[1].Project != "project-a" || top[1].Count != 3 {
		t.Errorf("expected project-a with 3, got %s with %d", top[1].Project, top[1].Count)
	}
}
