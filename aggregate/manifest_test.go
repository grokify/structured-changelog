package aggregate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewManifest(t *testing.T) {
	m := NewManifest("Test Portfolio")

	if m.Name != "Test Portfolio" {
		t.Errorf("expected name %q, got %q", "Test Portfolio", m.Name)
	}

	if len(m.Projects) != 0 {
		t.Errorf("expected empty projects, got %d", len(m.Projects))
	}
}

func TestManifestAddProject(t *testing.T) {
	m := NewManifest("Test")

	// Add first project
	added := m.AddProject(ProjectRef{Path: "github.com/org/repo1"})
	if !added {
		t.Error("expected first project to be added")
	}

	// Add duplicate
	added = m.AddProject(ProjectRef{Path: "github.com/org/repo1"})
	if added {
		t.Error("expected duplicate to not be added")
	}

	// Add different project
	added = m.AddProject(ProjectRef{Path: "github.com/org/repo2"})
	if !added {
		t.Error("expected second project to be added")
	}

	if len(m.Projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(m.Projects))
	}
}

func TestManifestAddSource(t *testing.T) {
	m := NewManifest("Test")

	// Add first source
	added := m.AddSource(Source{Type: SourceTypeOrg, Name: "myorg"})
	if !added {
		t.Error("expected first source to be added")
	}

	// Add duplicate
	added = m.AddSource(Source{Type: SourceTypeOrg, Name: "myorg"})
	if added {
		t.Error("expected duplicate to not be added")
	}

	// Add different source
	added = m.AddSource(Source{Type: SourceTypeUser, Name: "myuser"})
	if !added {
		t.Error("expected second source to be added")
	}

	if len(m.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(m.Sources))
	}
}

func TestManifestValidate(t *testing.T) {
	tests := []struct {
		name     string
		manifest *Manifest
		valid    bool
		errField string
	}{
		{
			name:     "valid manifest",
			manifest: &Manifest{Name: "Test", Projects: []ProjectRef{{Path: "github.com/org/repo"}}},
			valid:    true,
		},
		{
			name:     "missing name",
			manifest: &Manifest{Projects: []ProjectRef{{Path: "github.com/org/repo"}}},
			valid:    false,
			errField: "name",
		},
		{
			name:     "invalid source type",
			manifest: &Manifest{Name: "Test", Sources: []Source{{Type: "invalid", Name: "test"}}},
			valid:    false,
			errField: "sources[0].type",
		},
		{
			name:     "missing source name",
			manifest: &Manifest{Name: "Test", Sources: []Source{{Type: SourceTypeOrg, Name: ""}}},
			valid:    false,
			errField: "sources[0].name",
		},
		{
			name:     "missing project path",
			manifest: &Manifest{Name: "Test", Projects: []ProjectRef{{Path: ""}}},
			valid:    false,
			errField: "projects[0].path",
		},
		{
			name: "duplicate project paths",
			manifest: &Manifest{
				Name: "Test",
				Projects: []ProjectRef{
					{Path: "github.com/org/repo"},
					{Path: "github.com/org/repo"},
				},
			},
			valid:    false,
			errField: "projects[1].path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.manifest.Validate()
			if result.Valid != tt.valid {
				t.Errorf("expected valid=%v, got %v", tt.valid, result.Valid)
			}

			if !tt.valid && tt.errField != "" {
				found := false
				for _, err := range result.Errors {
					if err.Field == tt.errField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error on field %q, got errors: %+v", tt.errField, result.Errors)
				}
			}
		})
	}
}

func TestManifestJSON(t *testing.T) {
	m := &Manifest{
		Name:        "Test Portfolio",
		Description: "A test portfolio",
		Sources: []Source{
			{Type: SourceTypeOrg, Name: "myorg"},
		},
		Projects: []ProjectRef{
			{Path: "github.com/org/repo", Discovered: true},
		},
	}

	data, err := m.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	// Parse back
	var m2 Manifest
	if err := json.Unmarshal(data, &m2); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if m2.Name != m.Name {
		t.Errorf("expected name %q, got %q", m.Name, m2.Name)
	}

	if len(m2.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(m2.Projects))
	}
}

func TestManifestWriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")

	m := &Manifest{
		Name: "Test Portfolio",
		Projects: []ProjectRef{
			{Path: "github.com/org/repo"},
		},
	}

	// Write
	if err := m.WriteFile(path); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}

	// Load
	loaded, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest error: %v", err)
	}

	if loaded.Name != m.Name {
		t.Errorf("expected name %q, got %q", m.Name, loaded.Name)
	}
}
