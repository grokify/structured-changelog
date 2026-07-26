package changelog

import (
	"path/filepath"
	"sort"
	"time"
)

// RepoHighlights contains highlights from a single repository for a date range.
type RepoHighlights struct {
	Project    string           `json:"project"`
	Repository string           `json:"repository,omitempty"`
	Path       string           `json:"path,omitempty"`
	Releases   []ReleaseExtract `json:"releases"`
}

// ReleaseExtract is a release with its notable entries.
type ReleaseExtract struct {
	Version    string  `json:"version"`
	Date       string  `json:"date"`
	Highlights []Entry `json:"highlights,omitempty"`
	Added      []Entry `json:"added,omitempty"`
	Changed    []Entry `json:"changed,omitempty"`
	Fixed      []Entry `json:"fixed,omitempty"`
	Security   []Entry `json:"security,omitempty"`
	Breaking   []Entry `json:"breaking,omitempty"`
}

// IsEmpty returns true if the extract has no entries.
func (r ReleaseExtract) IsEmpty() bool {
	return len(r.Highlights) == 0 &&
		len(r.Added) == 0 &&
		len(r.Changed) == 0 &&
		len(r.Fixed) == 0 &&
		len(r.Security) == 0 &&
		len(r.Breaking) == 0
}

// EntryCount returns the total number of entries.
func (r ReleaseExtract) EntryCount() int {
	return len(r.Highlights) + len(r.Added) + len(r.Changed) +
		len(r.Fixed) + len(r.Security) + len(r.Breaking)
}

// ExtractHighlights extracts notable releases within a date range.
// Only releases with dates in [since, until) are included.
// Only user-facing categories are extracted: highlights, added, changed,
// fixed, security, breaking.
func (c *Changelog) ExtractHighlights(since, until time.Time) *RepoHighlights {
	result := &RepoHighlights{
		Project:    c.Project,
		Repository: c.Repository,
		Releases:   []ReleaseExtract{},
	}

	for _, rel := range c.Releases {
		if rel.Date == "" {
			continue
		}

		// Parse release date
		relDate, err := time.Parse("2006-01-02", rel.Date)
		if err != nil {
			continue
		}

		// Check date range [since, until)
		if relDate.Before(since) || !relDate.Before(until) {
			continue
		}

		extract := ReleaseExtract{
			Version:    rel.Version,
			Date:       rel.Date,
			Highlights: rel.Highlights,
			Added:      rel.Added,
			Changed:    rel.Changed,
			Fixed:      rel.Fixed,
			Security:   rel.Security,
			Breaking:   rel.Breaking,
		}

		// Skip maintenance-only releases
		if extract.IsEmpty() {
			continue
		}

		result.Releases = append(result.Releases, extract)
	}

	return result
}

// MultiRepoHighlights aggregates highlights across multiple repositories.
type MultiRepoHighlights struct {
	Since       time.Time        `json:"since"`
	Until       time.Time        `json:"until"`
	TotalCount  int              `json:"totalCount"`
	RepoCount   int              `json:"repoCount"`
	ByRepo      []RepoHighlights `json:"byRepo"`
	TopReleases []RankedRelease  `json:"topReleases,omitempty"`
}

// RankedRelease is a release with its repo context, for cross-repo ranking.
type RankedRelease struct {
	Project    string         `json:"project"`
	Repository string         `json:"repository,omitempty"`
	Release    ReleaseExtract `json:"release"`
}

// ExtractMultiRepoHighlights loads CHANGELOG.json from multiple paths and
// extracts highlights for the date range.
func ExtractMultiRepoHighlights(paths []string, since, until time.Time) (*MultiRepoHighlights, error) {
	result := &MultiRepoHighlights{
		Since:  since,
		Until:  until,
		ByRepo: []RepoHighlights{},
	}

	for _, path := range paths {
		cl, err := LoadFile(path)
		if err != nil {
			continue // Skip repos without valid CHANGELOG.json
		}

		highlights := cl.ExtractHighlights(since, until)
		highlights.Path = path

		if len(highlights.Releases) == 0 {
			continue
		}

		result.ByRepo = append(result.ByRepo, *highlights)
		result.RepoCount++

		for _, rel := range highlights.Releases {
			result.TotalCount += rel.EntryCount()
			result.TopReleases = append(result.TopReleases, RankedRelease{
				Project:    highlights.Project,
				Repository: highlights.Repository,
				Release:    rel,
			})
		}
	}

	// Sort top releases by entry count (descending)
	sort.Slice(result.TopReleases, func(i, j int) bool {
		return result.TopReleases[i].Release.EntryCount() > result.TopReleases[j].Release.EntryCount()
	})

	return result, nil
}

// DiscoverChangelogs finds CHANGELOG.json files under the given directories.
func DiscoverChangelogs(roots []string) []string {
	var paths []string
	seen := make(map[string]bool)

	for _, root := range roots {
		// Look for CHANGELOG.json directly
		direct := filepath.Join(root, "CHANGELOG.json")
		if !seen[direct] {
			seen[direct] = true
			paths = append(paths, direct)
		}
	}

	return paths
}

// TopProjects returns projects ranked by total entry count.
func (m *MultiRepoHighlights) TopProjects(limit int) []struct {
	Project string
	Count   int
} {
	counts := make(map[string]int)
	for _, repo := range m.ByRepo {
		for _, rel := range repo.Releases {
			counts[repo.Project] += rel.EntryCount()
		}
	}

	type projectCount struct {
		Project string
		Count   int
	}
	var ranked []projectCount
	for p, c := range counts {
		ranked = append(ranked, projectCount{p, c})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Count > ranked[j].Count
	})

	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}

	result := make([]struct {
		Project string
		Count   int
	}, len(ranked))
	for i, pc := range ranked {
		result[i].Project = pc.Project
		result[i].Count = pc.Count
	}
	return result
}
