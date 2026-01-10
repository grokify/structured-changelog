// Package changelog provides the JSON IR (Intermediate Representation) types
// for structured changelogs following the Keep a Changelog format.
package changelog

import (
	"encoding/json"
	"os"
	"time"
)

// IRVersion is the current version of the IR schema.
const IRVersion = "1.0"

// Versioning scheme constants.
const (
	VersioningSemVer = "semver" // Semantic Versioning (default)
	VersioningCalVer = "calver" // Calendar Versioning
	VersioningCustom = "custom" // Custom versioning scheme
	VersioningNone   = "none"   // No specific versioning scheme
)

// Commit convention constants.
const (
	CommitConventionConventional = "conventional" // Conventional Commits
	CommitConventionNone         = "none"         // No specific convention (default)
)

// Changelog represents the root of a structured changelog.
type Changelog struct {
	IRVersion        string     `json:"ir_version"`
	Project          string     `json:"project"`
	Repository       string     `json:"repository,omitempty"`
	Versioning       string     `json:"versioning,omitempty"`
	CommitConvention string     `json:"commit_convention,omitempty"`
	GeneratedAt      *time.Time `json:"generated_at,omitempty"`
	Unreleased       *Release   `json:"unreleased,omitempty"`
	Releases         []Release  `json:"releases,omitempty"`
}

// New creates a new Changelog with the current IR version.
func New(project string) *Changelog {
	return &Changelog{
		IRVersion: IRVersion,
		Project:   project,
		Releases:  []Release{},
	}
}

// LoadFile loads a Changelog from a JSON file.
func LoadFile(path string) (*Changelog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse parses a Changelog from JSON bytes.
func Parse(data []byte) (*Changelog, error) {
	var cl Changelog
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, err
	}
	return &cl, nil
}

// JSON returns the changelog as formatted JSON bytes.
func (c *Changelog) JSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// WriteFile writes the changelog to a JSON file.
func (c *Changelog) WriteFile(path string) error {
	data, err := c.JSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// AddRelease adds a new release to the changelog.
// Releases are maintained in reverse chronological order.
func (c *Changelog) AddRelease(r Release) {
	c.Releases = append([]Release{r}, c.Releases...)
}

// LatestRelease returns the most recent release, or nil if none exist.
func (c *Changelog) LatestRelease() *Release {
	if len(c.Releases) == 0 {
		return nil
	}
	return &c.Releases[0]
}

// PromoteUnreleased moves unreleased changes to a new release.
func (c *Changelog) PromoteUnreleased(version, date string) error {
	if c.Unreleased == nil {
		return nil
	}
	release := *c.Unreleased
	release.Version = version
	release.Date = date
	c.AddRelease(release)
	c.Unreleased = nil
	return nil
}

// Summary contains a summary of a changelog's contents.
type Summary struct {
	Project              string
	IRVersion            string
	ReleaseCount         int
	HasUnreleased        bool
	UnreleasedCategories []string
	LatestVersion        string
	LatestDate           string
	LatestCategories     []string
}

// Summary returns a summary of the changelog's contents.
func (c *Changelog) Summary() Summary {
	s := Summary{
		Project:      c.Project,
		IRVersion:    c.IRVersion,
		ReleaseCount: len(c.Releases),
	}

	// Check unreleased section
	if c.Unreleased != nil && !c.Unreleased.IsEmpty() {
		s.HasUnreleased = true
		for _, cat := range c.Unreleased.Categories() {
			s.UnreleasedCategories = append(s.UnreleasedCategories, cat.Name)
		}
	}

	// Get latest release info
	if len(c.Releases) > 0 {
		latest := c.Releases[0]
		s.LatestVersion = latest.Version
		s.LatestDate = latest.Date
		for _, cat := range latest.Categories() {
			s.LatestCategories = append(s.LatestCategories, cat.Name)
		}
	}

	return s
}
