package aggregate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/grokify/structured-changelog/changelog"
)

// Portfolio represents an aggregated collection of project changelogs.
type Portfolio struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Projects    []ProjectData `json:"projects"`
	DateRange   DateRange     `json:"dateRange"`
	GeneratedAt time.Time     `json:"generatedAt"`
}

// ProjectData holds a project's changelog data.
type ProjectData struct {
	Path      string               `json:"path"`
	Name      string               `json:"name"`
	Changelog *changelog.Changelog `json:"changelog"`
}

// DateRange represents a time range.
type DateRange struct {
	Start string `json:"start"` // YYYY-MM-DD
	End   string `json:"end"`   // YYYY-MM-DD
}

// LoadPortfolio loads changelogs for resolved projects and creates a portfolio.
func LoadPortfolio(manifest *Manifest, resolved []ResolvedProject) (*Portfolio, error) {
	portfolio := &Portfolio{
		Name:        manifest.Name,
		Description: manifest.Description,
		GeneratedAt: time.Now().UTC(),
	}

	var minDate, maxDate string

	for _, rp := range resolved {
		if !rp.IsLocal {
			continue // Skip unresolved projects
		}

		cl, err := changelog.LoadFile(rp.ChangelogPath)
		if err != nil {
			return nil, fmt.Errorf("loading changelog for %s: %w", rp.Ref.Path, err)
		}

		pd := ProjectData{
			Path:      rp.Ref.Path,
			Name:      cl.Project,
			Changelog: cl,
		}

		// Update date range
		for _, release := range cl.Releases {
			if release.Date != "" {
				if minDate == "" || release.Date < minDate {
					minDate = release.Date
				}
				if maxDate == "" || release.Date > maxDate {
					maxDate = release.Date
				}
			}
		}

		portfolio.Projects = append(portfolio.Projects, pd)
	}

	portfolio.DateRange = DateRange{
		Start: minDate,
		End:   maxDate,
	}

	return portfolio, nil
}

// LoadPortfolioFromPaths loads changelogs from explicit file paths.
func LoadPortfolioFromPaths(name string, paths []string) (*Portfolio, error) {
	portfolio := &Portfolio{
		Name:        name,
		GeneratedAt: time.Now().UTC(),
	}

	var minDate, maxDate string

	for _, path := range paths {
		cl, err := changelog.LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("loading changelog from %s: %w", path, err)
		}

		// Derive project path from file path
		projectPath := filepath.Dir(path)

		pd := ProjectData{
			Path:      projectPath,
			Name:      cl.Project,
			Changelog: cl,
		}

		// Update date range
		for _, release := range cl.Releases {
			if release.Date != "" {
				if minDate == "" || release.Date < minDate {
					minDate = release.Date
				}
				if maxDate == "" || release.Date > maxDate {
					maxDate = release.Date
				}
			}
		}

		portfolio.Projects = append(portfolio.Projects, pd)
	}

	portfolio.DateRange = DateRange{
		Start: minDate,
		End:   maxDate,
	}

	return portfolio, nil
}

// LoadPortfolioFile loads a portfolio from a JSON file.
func LoadPortfolioFile(path string) (*Portfolio, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading portfolio file: %w", err)
	}
	return ParsePortfolio(data)
}

// ParsePortfolio parses a portfolio from JSON bytes.
func ParsePortfolio(data []byte) (*Portfolio, error) {
	var p Portfolio
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing portfolio JSON: %w", err)
	}
	return &p, nil
}

// WriteFile writes the portfolio to a JSON file.
func (p *Portfolio) WriteFile(path string) error {
	data, err := p.JSON()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing portfolio file: %w", err)
	}
	return nil
}

// JSON returns the portfolio as formatted JSON bytes.
func (p *Portfolio) JSON() ([]byte, error) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling portfolio JSON: %w", err)
	}
	return data, nil
}

// Summary returns aggregate statistics about the portfolio.
func (p *Portfolio) Summary() PortfolioSummary {
	var summary PortfolioSummary
	summary.ProjectCount = len(p.Projects)
	summary.ByCategory = make(map[string]int)

	for _, pd := range p.Projects {
		if pd.Changelog == nil {
			continue
		}
		summary.ReleaseCount += len(pd.Changelog.Releases)

		for _, release := range pd.Changelog.Releases {
			cats := release.Categories()
			for _, cat := range cats {
				entries := release.GetEntries(cat.Name)
				summary.EntryCount += len(entries)
				summary.ByCategory[cat.Name] += len(entries)
			}
		}
	}

	summary.DateRange = p.DateRange
	return summary
}

// PortfolioSummary provides aggregate statistics.
type PortfolioSummary struct {
	ProjectCount int            `json:"projectCount"`
	ReleaseCount int            `json:"releaseCount"`
	EntryCount   int            `json:"entryCount"`
	ByCategory   map[string]int `json:"byCategory"`
	DateRange    DateRange      `json:"dateRange"`
}

// AllReleases returns all releases across all projects, sorted by date descending.
func (p *Portfolio) AllReleases() []ReleaseWithProject {
	var releases []ReleaseWithProject

	for _, pd := range p.Projects {
		if pd.Changelog == nil {
			continue
		}
		for i := range pd.Changelog.Releases {
			releases = append(releases, ReleaseWithProject{
				ProjectPath: pd.Path,
				ProjectName: pd.Name,
				Release:     &pd.Changelog.Releases[i],
			})
		}
	}

	// Sort by date descending
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Release.Date > releases[j].Release.Date
	})

	return releases
}

// ReleaseWithProject pairs a release with its project info.
type ReleaseWithProject struct {
	ProjectPath string             `json:"projectPath"`
	ProjectName string             `json:"projectName"`
	Release     *changelog.Release `json:"release"`
}

// FilterByDateRange filters releases to those within the given date range.
func (p *Portfolio) FilterByDateRange(start, end string) []ReleaseWithProject {
	all := p.AllReleases()
	var filtered []ReleaseWithProject

	for _, rwp := range all {
		if rwp.Release.Date == "" {
			continue
		}
		if start != "" && rwp.Release.Date < start {
			continue
		}
		if end != "" && rwp.Release.Date > end {
			continue
		}
		filtered = append(filtered, rwp)
	}

	return filtered
}

// ProjectNames returns the names of all projects.
func (p *Portfolio) ProjectNames() []string {
	names := make([]string, len(p.Projects))
	for i, pd := range p.Projects {
		names[i] = pd.Name
	}
	return names
}

// ProjectPaths returns the paths of all projects.
func (p *Portfolio) ProjectPaths() []string {
	paths := make([]string, len(p.Projects))
	for i, pd := range p.Projects {
		paths[i] = pd.Path
	}
	return paths
}

// GetProject returns a project by path, or nil if not found.
func (p *Portfolio) GetProject(path string) *ProjectData {
	for i := range p.Projects {
		if p.Projects[i].Path == path {
			return &p.Projects[i]
		}
	}
	return nil
}
