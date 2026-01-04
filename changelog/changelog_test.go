package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	cl := New("test-project")

	if cl.IRVersion != IRVersion {
		t.Errorf("expected IRVersion %s, got %s", IRVersion, cl.IRVersion)
	}
	if cl.Project != "test-project" {
		t.Errorf("expected project test-project, got %s", cl.Project)
	}
	if len(cl.Releases) != 0 {
		t.Errorf("expected 0 releases, got %d", len(cl.Releases))
	}
}

func TestParse(t *testing.T) {
	jsonData := []byte(`{
		"ir_version": "1.0",
		"project": "my-project",
		"releases": [
			{
				"version": "1.0.0",
				"date": "2026-01-03",
				"added": [
					{"description": "Initial release"}
				]
			}
		]
	}`)

	cl, err := Parse(jsonData)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cl.Project != "my-project" {
		t.Errorf("expected project my-project, got %s", cl.Project)
	}
	if len(cl.Releases) != 1 {
		t.Fatalf("expected 1 release, got %d", len(cl.Releases))
	}
	if cl.Releases[0].Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", cl.Releases[0].Version)
	}
}

func TestAddRelease(t *testing.T) {
	cl := New("test")
	cl.AddRelease(NewRelease("1.0.0", "2026-01-01"))
	cl.AddRelease(NewRelease("1.1.0", "2026-01-02"))

	if len(cl.Releases) != 2 {
		t.Fatalf("expected 2 releases, got %d", len(cl.Releases))
	}
	// Newest should be first
	if cl.Releases[0].Version != "1.1.0" {
		t.Errorf("expected newest release first, got %s", cl.Releases[0].Version)
	}
}

func TestLatestRelease(t *testing.T) {
	cl := New("test")

	if cl.LatestRelease() != nil {
		t.Error("expected nil for empty changelog")
	}

	cl.AddRelease(NewRelease("1.0.0", "2026-01-01"))
	latest := cl.LatestRelease()

	if latest == nil {
		t.Fatal("expected non-nil latest release")
	}
	if latest.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", latest.Version)
	}
}

func TestPromoteUnreleased(t *testing.T) {
	cl := New("test")
	cl.Unreleased = &Release{
		Added: []Entry{{Description: "New feature"}},
	}

	err := cl.PromoteUnreleased("1.0.0", "2026-01-03")
	if err != nil {
		t.Fatalf("PromoteUnreleased failed: %v", err)
	}

	if cl.Unreleased != nil {
		t.Error("expected unreleased to be nil after promotion")
	}
	if len(cl.Releases) != 1 {
		t.Fatalf("expected 1 release, got %d", len(cl.Releases))
	}
	if cl.Releases[0].Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", cl.Releases[0].Version)
	}
	if len(cl.Releases[0].Added) != 1 {
		t.Error("expected promoted release to have 1 added entry")
	}
}

func TestJSON(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(NewRelease("1.0.0", "2026-01-03"))

	data, err := cl.JSON()
	if err != nil {
		t.Fatalf("JSON failed: %v", err)
	}

	// Parse it back
	cl2, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cl2.Project != cl.Project {
		t.Errorf("roundtrip failed: project mismatch")
	}
	if len(cl2.Releases) != len(cl.Releases) {
		t.Errorf("roundtrip failed: releases count mismatch")
	}
}

func TestLoadFile(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.json")

	content := []byte(`{
		"ir_version": "1.0",
		"project": "file-test",
		"releases": []
	}`)

	if err := os.WriteFile(tmpFile, content, 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cl, err := LoadFile(tmpFile)
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	if cl.Project != "file-test" {
		t.Errorf("expected project 'file-test', got %q", cl.Project)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/file.json")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(tmpFile, []byte("not valid json"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	_, err := LoadFile(tmpFile)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestWriteFile(t *testing.T) {
	cl := New("write-test")
	cl.AddRelease(NewRelease("1.0.0", "2026-01-04"))

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "output.json")

	if err := cl.WriteFile(tmpFile); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Read it back
	cl2, err := LoadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if cl2.Project != "write-test" {
		t.Errorf("expected project 'write-test', got %q", cl2.Project)
	}
	if len(cl2.Releases) != 1 {
		t.Errorf("expected 1 release, got %d", len(cl2.Releases))
	}
}

func TestWriteFile_InvalidPath(t *testing.T) {
	cl := New("test")
	err := cl.WriteFile("/nonexistent/directory/file.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte("not valid json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestPromoteUnreleased_Nil(t *testing.T) {
	cl := New("test")
	// Unreleased is nil
	err := cl.PromoteUnreleased("1.0.0", "2026-01-04")
	if err != nil {
		t.Errorf("expected no error for nil unreleased, got %v", err)
	}
	if len(cl.Releases) != 0 {
		t.Errorf("expected 0 releases, got %d", len(cl.Releases))
	}
}

func TestSummary_Empty(t *testing.T) {
	cl := New("test-project")
	s := cl.Summary()

	if s.Project != "test-project" {
		t.Errorf("expected project test-project, got %s", s.Project)
	}
	if s.IRVersion != IRVersion {
		t.Errorf("expected IR version %s, got %s", IRVersion, s.IRVersion)
	}
	if s.ReleaseCount != 0 {
		t.Errorf("expected 0 releases, got %d", s.ReleaseCount)
	}
	if s.HasUnreleased {
		t.Error("expected HasUnreleased to be false")
	}
	if s.LatestVersion != "" {
		t.Errorf("expected empty LatestVersion, got %s", s.LatestVersion)
	}
}

func TestSummary_WithReleases(t *testing.T) {
	cl := New("test-project")
	cl.AddRelease(Release{
		Version: "1.0.0",
		Date:    "2026-01-03",
		Added:   []Entry{{Description: "Initial release"}},
		Fixed:   []Entry{{Description: "Bug fix"}},
	})

	s := cl.Summary()

	if s.ReleaseCount != 1 {
		t.Errorf("expected 1 release, got %d", s.ReleaseCount)
	}
	if s.LatestVersion != "1.0.0" {
		t.Errorf("expected latest version 1.0.0, got %s", s.LatestVersion)
	}
	if s.LatestDate != "2026-01-03" {
		t.Errorf("expected latest date 2026-01-03, got %s", s.LatestDate)
	}
	if len(s.LatestCategories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(s.LatestCategories))
	}
}

func TestSummary_WithUnreleased(t *testing.T) {
	cl := New("test-project")
	cl.Unreleased = &Release{
		Added: []Entry{{Description: "New feature"}},
	}

	s := cl.Summary()

	if !s.HasUnreleased {
		t.Error("expected HasUnreleased to be true")
	}
	if len(s.UnreleasedCategories) != 1 {
		t.Errorf("expected 1 unreleased category, got %d", len(s.UnreleasedCategories))
	}
	if s.UnreleasedCategories[0] != "Added" {
		t.Errorf("expected Added category, got %s", s.UnreleasedCategories[0])
	}
}

func TestSummary_EmptyUnreleased(t *testing.T) {
	cl := New("test-project")
	cl.Unreleased = &Release{} // empty release

	s := cl.Summary()

	if s.HasUnreleased {
		t.Error("expected HasUnreleased to be false for empty unreleased")
	}
}
