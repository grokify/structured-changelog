package gitlog

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Tag represents a git tag with metadata.
type Tag struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	DateString  string    `json:"dateString"`
	CommitHash  string    `json:"commitHash"`
	CommitCount int       `json:"commitCount,omitempty"` // Commits since previous tag
	IsInitial   bool      `json:"isInitial,omitempty"`   // True if this is the first tag
}

// TagList represents a list of tags with metadata.
type TagList struct {
	Repository  string    `json:"repository,omitempty"`
	Tags        []Tag     `json:"tags"`
	TotalTags   int       `json:"totalTags"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// semverRegex matches semantic version tags like v1.0.0, v1.2.3-beta, 1.0.0
var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)`)

// GetTags returns all semver tags in the repository sorted by version.
func GetTags() (*TagList, error) {
	// Get all tags
	cmd := exec.Command("git", "tag", "--list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tagNames := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tagNames) == 0 || (len(tagNames) == 1 && tagNames[0] == "") {
		return &TagList{
			Tags:        []Tag{},
			TotalTags:   0,
			GeneratedAt: time.Now().UTC(),
		}, nil
	}

	// Filter to semver tags only
	var semverTags []string
	for _, tag := range tagNames {
		tag = strings.TrimSpace(tag)
		if tag != "" && semverRegex.MatchString(tag) {
			semverTags = append(semverTags, tag)
		}
	}

	// Sort by semver
	sort.Slice(semverTags, func(i, j int) bool {
		return compareSemver(semverTags[i], semverTags[j]) < 0
	})

	// Get metadata for each tag
	var tags []Tag
	for i, tagName := range semverTags {
		tag, err := getTagMetadata(tagName)
		if err != nil {
			continue // Skip tags we can't get metadata for
		}

		// Calculate commit count since previous tag
		if i == 0 {
			tag.IsInitial = true
			// Count commits from beginning to this tag
			count, _ := countCommits("", tagName)
			tag.CommitCount = count
		} else {
			prevTag := semverTags[i-1]
			count, _ := countCommits(prevTag, tagName)
			tag.CommitCount = count
		}

		tags = append(tags, *tag)
	}

	return &TagList{
		Tags:        tags,
		TotalTags:   len(tags),
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// getTagMetadata retrieves date and commit hash for a tag.
func getTagMetadata(tagName string) (*Tag, error) {
	// Get commit hash
	hashCmd := exec.Command("git", "rev-list", "-n", "1", tagName)
	hashOutput, err := hashCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash for tag %s: %w", tagName, err)
	}

	// Get commit date
	dateCmd := exec.Command("git", "log", "-1", "--format=%aI", tagName)
	dateOutput, err := dateCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get date for tag %s: %w", tagName, err)
	}

	dateStr := strings.TrimSpace(string(dateOutput))
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date for tag %s: %w", tagName, err)
	}

	return &Tag{
		Name:       tagName,
		Date:       date,
		DateString: date.Format("2006-01-02"),
		CommitHash: strings.TrimSpace(string(hashOutput)),
	}, nil
}

// countCommits counts commits between two refs.
// If since is empty, counts all commits up to until.
func countCommits(since, until string) (int, error) {
	var args []string
	if since == "" {
		args = []string{"rev-list", "--count", until}
	} else {
		args = []string{"rev-list", "--count", fmt.Sprintf("%s..%s", since, until)}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// compareSemver compares two semver strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareSemver(a, b string) int {
	aMatch := semverRegex.FindStringSubmatch(a)
	bMatch := semverRegex.FindStringSubmatch(b)

	if aMatch == nil || bMatch == nil {
		return strings.Compare(a, b)
	}

	for i := 1; i <= 3; i++ {
		aNum, _ := strconv.Atoi(aMatch[i])
		bNum, _ := strconv.Atoi(bMatch[i])
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
	}

	return 0
}

// GetFirstCommit returns the hash of the first commit in the repository.
func GetFirstCommit() (string, error) {
	cmd := exec.Command("git", "rev-list", "--max-parents=0", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get first commit: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no commits found")
	}

	// Return the first (oldest) root commit
	return strings.TrimSpace(lines[len(lines)-1]), nil
}

// VersionRange represents a range between two versions for parsing.
type VersionRange struct {
	Version string `json:"version"`
	Since   string `json:"since"`   // Previous version (empty for first)
	Until   string `json:"until"`   // This version
	Date    string `json:"date"`    // Release date
	Commits int    `json:"commits"` // Commit count in range
}

// GetAllVersionRanges returns all version ranges for parsing commits.
func GetAllVersionRanges() ([]VersionRange, error) {
	tagList, err := GetTags()
	if err != nil {
		return nil, err
	}

	var ranges []VersionRange
	for i, tag := range tagList.Tags {
		vr := VersionRange{
			Version: tag.Name,
			Until:   tag.Name,
			Date:    tag.DateString,
			Commits: tag.CommitCount,
		}

		if i > 0 {
			vr.Since = tagList.Tags[i-1].Name
		}

		ranges = append(ranges, vr)
	}

	return ranges, nil
}
