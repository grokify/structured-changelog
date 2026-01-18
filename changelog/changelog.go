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
	TagPath          string     `json:"tag_path,omitempty"`
	Versioning       string     `json:"versioning,omitempty"`
	CommitConvention string     `json:"commit_convention,omitempty"`
	Maintainers      []string   `json:"maintainers,omitempty"`
	Bots             []string   `json:"bots,omitempty"`
	GeneratedAt      *time.Time `json:"generated_at,omitempty"`
	Unreleased       *Release   `json:"unreleased,omitempty"`
	Releases         []Release  `json:"releases,omitempty"`
}

// CommonBots is a list of well-known bot usernames that are auto-detected.
var CommonBots = []string{
	"dependabot",
	"dependabot[bot]",
	"renovate",
	"renovate[bot]",
	"github-actions",
	"github-actions[bot]",
	"semantic-release-bot",
	"greenkeeper[bot]",
	"snyk-bot",
	"imgbot[bot]",
	"allcontributors[bot]",
}

// IsTeamMember returns true if the author is a maintainer or known bot.
func (c *Changelog) IsTeamMember(author string) bool {
	return c.IsTeamMemberByNameAndEmail(author, "")
}

// IsTeamMemberByNameAndEmail returns true if the author (by name or email) is a maintainer or known bot.
// This is useful when parsing git commits where you have both author name and email.
// It checks:
// 1. If author name matches a maintainer
// 2. If email matches a maintainer entry (for emails in maintainers list)
// 3. If GitHub username from noreply email matches a maintainer
// 4. If author matches a known bot
func (c *Changelog) IsTeamMemberByNameAndEmail(author, email string) bool {
	if author == "" && email == "" {
		return true // No author means no attribution needed
	}

	normAuthor := normalizeAuthor(author)
	normEmail := normalizeAuthor(email)

	// Check maintainers against author name and email
	for _, m := range c.Maintainers {
		normM := normalizeAuthor(m)
		if normM == normAuthor {
			return true
		}
		// Also check if email matches (allows emails in maintainers list)
		if normEmail != "" && normM == normEmail {
			return true
		}
	}

	// Check if email contains a GitHub username (noreply format)
	// Format: username@users.noreply.github.com or 12345+username@users.noreply.github.com
	if email != "" {
		if username := extractGitHubUsername(email); username != "" {
			for _, m := range c.Maintainers {
				if normalizeAuthor(m) == normalizeAuthor(username) {
					return true
				}
			}
		}
	}

	// Check custom bots
	for _, b := range c.Bots {
		if normalizeAuthor(b) == normAuthor {
			return true
		}
	}

	// Check common bots
	for _, b := range CommonBots {
		if normalizeAuthor(b) == normAuthor {
			return true
		}
	}

	return false
}

// extractGitHubUsername extracts a GitHub username from a noreply email.
// Handles formats:
// - username@users.noreply.github.com
// - 12345+username@users.noreply.github.com
func extractGitHubUsername(email string) string {
	// Check for GitHub noreply format
	suffix := "@users.noreply.github.com"
	if len(email) <= len(suffix) {
		return ""
	}

	// Convert to lowercase for comparison
	emailLower := normalizeAuthor(email)
	suffixLower := normalizeAuthor(suffix)

	if !hasEmailSuffix(emailLower, suffixLower) {
		return ""
	}

	// Extract the part before @users.noreply.github.com
	local := email[:len(email)-len(suffix)]

	// Handle 12345+username format
	if idx := indexByte(local, '+'); idx >= 0 {
		local = local[idx+1:]
	}

	return local
}

// hasEmailSuffix checks if email ends with suffix.
func hasEmailSuffix(email, suffix string) bool {
	if len(email) < len(suffix) {
		return false
	}
	return email[len(email)-len(suffix):] == suffix
}

// indexByte returns the index of the first instance of c in s, or -1 if c is not present.
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// normalizeAuthor normalizes an author string for comparison.
// Removes @ prefix and converts to lowercase.
func normalizeAuthor(author string) string {
	if len(author) > 0 && author[0] == '@' {
		author = author[1:]
	}
	// Use lowercase for case-insensitive comparison
	result := make([]byte, len(author))
	for i := 0; i < len(author); i++ {
		c := author[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
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
