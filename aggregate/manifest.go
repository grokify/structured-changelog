// Package aggregate provides multi-project changelog aggregation.
// It enables combining changelogs from multiple repositories into
// a unified portfolio view with metrics and dashboard export capabilities.
package aggregate

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Manifest lists projects to aggregate (like go.mod but for changelogs).
type Manifest struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Sources     []Source     `json:"sources,omitempty"`
	Projects    []ProjectRef `json:"projects"`
	Generated   *time.Time   `json:"generated,omitempty"`
}

// Source defines a GitHub org or user to scan for changelogs.
type Source struct {
	Type string `json:"type"` // "org" or "user"
	Name string `json:"name"` // GitHub org or username
}

// SourceTypeOrg indicates a GitHub organization source.
const SourceTypeOrg = "org"

// SourceTypeUser indicates a GitHub user source.
const SourceTypeUser = "user"

// ProjectRef references a project containing a changelog.
type ProjectRef struct {
	Path       string `json:"path"`                 // github.com/org/repo or github.com/org/repo/subdir
	LocalPath  string `json:"localPath,omitempty"`  // Override local resolution
	Discovered bool   `json:"discovered,omitempty"` // Auto-discovered vs manually added
}

// LoadManifest loads a manifest from a JSON file.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file: %w", err)
	}
	return ParseManifest(data)
}

// ParseManifest parses a manifest from JSON bytes.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest JSON: %w", err)
	}
	return &m, nil
}

// WriteFile writes the manifest to a JSON file.
func (m *Manifest) WriteFile(path string) error {
	data, err := m.JSON()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing manifest file: %w", err)
	}
	return nil
}

// JSON returns the manifest as formatted JSON bytes.
func (m *Manifest) JSON() ([]byte, error) {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling manifest JSON: %w", err)
	}
	return data, nil
}

// Validate checks the manifest for errors.
func (m *Manifest) Validate() ValidationResult {
	var result ValidationResult

	if m.Name == "" {
		result.addError("name", "manifest name is required")
	}

	// Validate sources
	for i, s := range m.Sources {
		if s.Type != SourceTypeOrg && s.Type != SourceTypeUser {
			result.addError(fmt.Sprintf("sources[%d].type", i),
				fmt.Sprintf("invalid source type %q, must be %q or %q", s.Type, SourceTypeOrg, SourceTypeUser))
		}
		if s.Name == "" {
			result.addError(fmt.Sprintf("sources[%d].name", i), "source name is required")
		}
	}

	// Validate projects
	for i, p := range m.Projects {
		if p.Path == "" {
			result.addError(fmt.Sprintf("projects[%d].path", i), "project path is required")
		}
	}

	// Check for duplicate paths
	seen := make(map[string]bool)
	for i, p := range m.Projects {
		if seen[p.Path] {
			result.addError(fmt.Sprintf("projects[%d].path", i),
				fmt.Sprintf("duplicate project path: %s", p.Path))
		}
		seen[p.Path] = true
	}

	result.Valid = len(result.Errors) == 0
	return result
}

// AddProject adds a project to the manifest if not already present.
// Returns true if added, false if already exists.
func (m *Manifest) AddProject(ref ProjectRef) bool {
	for _, p := range m.Projects {
		if p.Path == ref.Path {
			return false
		}
	}
	m.Projects = append(m.Projects, ref)
	return true
}

// AddSource adds a source to the manifest if not already present.
// Returns true if added, false if already exists.
func (m *Manifest) AddSource(source Source) bool {
	for _, s := range m.Sources {
		if s.Type == source.Type && s.Name == source.Name {
			return false
		}
	}
	m.Sources = append(m.Sources, source)
	return true
}

// ValidationResult holds manifest validation errors.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *ValidationResult) addError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// NewManifest creates a new manifest with the given name.
func NewManifest(name string) *Manifest {
	return &Manifest{
		Name:     name,
		Projects: []ProjectRef{},
	}
}
