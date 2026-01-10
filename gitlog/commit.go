// Package gitlog provides parsing of git log output into structured data
// optimized for LLM-assisted changelog generation.
package gitlog

import (
	"time"
)

// Commit represents a parsed git commit with structured metadata.
type Commit struct {
	Hash              string   `json:"hash"`
	ShortHash         string   `json:"short_hash"`
	Author            string   `json:"author"`
	AuthorEmail       string   `json:"author_email,omitempty"`
	Date              string   `json:"date"`
	Message           string   `json:"message"`
	Body              string   `json:"body,omitempty"`
	Type              string   `json:"type,omitempty"`
	Scope             string   `json:"scope,omitempty"`
	Subject           string   `json:"subject"`
	Breaking          bool     `json:"breaking,omitempty"`
	Issue             int      `json:"issue,omitempty"`
	PR                int      `json:"pr,omitempty"`
	FilesChanged      int      `json:"files_changed,omitempty"`
	Insertions        int      `json:"insertions,omitempty"`
	Deletions         int      `json:"deletions,omitempty"`
	Files             []string `json:"files,omitempty"`
	SuggestedCategory string   `json:"suggested_category,omitempty"`
	IsExternal        bool     `json:"is_external,omitempty"`
}

// Range represents the commit range that was parsed.
type Range struct {
	Since       string `json:"since,omitempty"`
	Until       string `json:"until,omitempty"`
	CommitCount int    `json:"commit_count"`
}

// Summary provides aggregate statistics about the parsed commits.
type Summary struct {
	ByType              map[string]int `json:"by_type,omitempty"`
	BySuggestedCategory map[string]int `json:"by_suggested_category,omitempty"`
	TotalFilesChanged   int            `json:"total_files_changed,omitempty"`
	TotalInsertions     int            `json:"total_insertions,omitempty"`
	TotalDeletions      int            `json:"total_deletions,omitempty"`
}

// Contributor represents an author with commit count.
type Contributor struct {
	Name        string `json:"name"`
	CommitCount int    `json:"commit_count"`
	IsExternal  bool   `json:"is_external,omitempty"`
}

// ParseResult is the complete output of parsing git commits.
type ParseResult struct {
	Repository   string        `json:"repository,omitempty"`
	Range        Range         `json:"range"`
	GeneratedAt  time.Time     `json:"generated_at"`
	Commits      []Commit      `json:"commits"`
	Summary      Summary       `json:"summary"`
	Contributors []Contributor `json:"contributors,omitempty"`
}

// NewParseResult creates a new ParseResult with initialized maps.
func NewParseResult() *ParseResult {
	return &ParseResult{
		GeneratedAt: time.Now().UTC(),
		Commits:     []Commit{},
		Summary: Summary{
			ByType:              make(map[string]int),
			BySuggestedCategory: make(map[string]int),
		},
	}
}

// AddCommit adds a commit and updates summary statistics.
func (pr *ParseResult) AddCommit(c Commit) {
	pr.Commits = append(pr.Commits, c)
	pr.Range.CommitCount = len(pr.Commits)

	// Update type summary
	if c.Type != "" {
		pr.Summary.ByType[c.Type]++
	}

	// Update category summary
	if c.SuggestedCategory != "" {
		pr.Summary.BySuggestedCategory[c.SuggestedCategory]++
	}

	// Update file stats
	pr.Summary.TotalFilesChanged += c.FilesChanged
	pr.Summary.TotalInsertions += c.Insertions
	pr.Summary.TotalDeletions += c.Deletions
}

// ComputeContributors builds the Contributors list from commits.
// Call this after all commits have been added and IsExternal has been set.
func (pr *ParseResult) ComputeContributors() {
	// Count commits per author
	authorCounts := make(map[string]int)
	authorExternal := make(map[string]bool)

	for i := range pr.Commits {
		c := &pr.Commits[i]
		if c.Author == "" {
			continue
		}
		authorCounts[c.Author]++
		if c.IsExternal {
			authorExternal[c.Author] = true
		}
	}

	// Build sorted contributor list (external first, then by commit count)
	var external, internal []Contributor
	for name, count := range authorCounts {
		contrib := Contributor{
			Name:        name,
			CommitCount: count,
			IsExternal:  authorExternal[name],
		}
		if contrib.IsExternal {
			external = append(external, contrib)
		} else {
			internal = append(internal, contrib)
		}
	}

	// Sort each group by commit count (descending)
	sortByCommitCount := func(contribs []Contributor) {
		for i := 0; i < len(contribs)-1; i++ {
			for j := i + 1; j < len(contribs); j++ {
				if contribs[j].CommitCount > contribs[i].CommitCount {
					contribs[i], contribs[j] = contribs[j], contribs[i]
				}
			}
		}
	}
	sortByCommitCount(external)
	sortByCommitCount(internal)

	// External contributors first
	pr.Contributors = append(external, internal...)
}
