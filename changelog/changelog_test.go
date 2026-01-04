package changelog

import (
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
