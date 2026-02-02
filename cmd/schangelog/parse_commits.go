package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-changelog/format"
	"github.com/grokify/structured-changelog/gitlog"
)

var (
	parseCommitsSince       string
	parseCommitsUntil       string
	parseCommitsLast        int
	parseCommitsPath        string
	parseCommitsNoFiles     bool
	parseCommitsNoMerges    bool
	parseCommitsFormat      string
	parseCommitsRepoURL     string
	parseCommitsChangelog   string
	parseCommitsAllVersions bool
)

var parseCommitsCmd = &cobra.Command{
	Use:   "parse-commits",
	Short: "Parse git commits into structured output for LLM consumption",
	Long: `Parse git commits into structured output optimized for LLM consumption.

This command runs git log and parses the output into a compact format
that reduces token usage when working with LLMs for changelog generation.

Output formats:
  - toon (default): Token-Oriented Object Notation, ~40% fewer tokens than JSON
  - json: Standard JSON with indentation
  - json-compact: Minified JSON

The output includes:
  - Parsed conventional commit components (type, scope, subject)
  - Suggested changelog categories based on commit type
  - File statistics (insertions, deletions, files changed)
  - Issue and PR references extracted from messages
  - Summary statistics grouped by type and category

Examples:
  # Parse commits since a tag (TOON format, default)
  schangelog parse-commits --since=v0.3.0

  # Parse commits with JSON output
  schangelog parse-commits --since=v0.3.0 --format=json

  # Parse commits between two refs
  schangelog parse-commits --since=v0.2.0 --until=v0.3.0

  # Parse last 20 commits
  schangelog parse-commits --last=20

  # Parse commits for specific path
  schangelog parse-commits --since=v0.3.0 --path=src/

  # Exclude file list from output
  schangelog parse-commits --since=v0.3.0 --no-files

  # Exclude merge commits
  schangelog parse-commits --since=v0.3.0 --no-merges

  # Mark external contributors (reads maintainers/bots from CHANGELOG.json)
  schangelog parse-commits --since=v0.3.0 --changelog=CHANGELOG.json

  # Parse all commits from the beginning of the repository to a tag
  schangelog parse-commits --until=v0.1.0

  # Parse commits for ALL version ranges at once (useful for backfilling)
  schangelog parse-commits --all-versions`,
	RunE: runParseCommits,
}

func init() {
	parseCommitsCmd.Flags().StringVar(&parseCommitsSince, "since", "", "Parse commits after this ref (tag, branch, or commit)")
	parseCommitsCmd.Flags().StringVar(&parseCommitsUntil, "until", "HEAD", "Parse commits up to this ref (default: HEAD)")
	parseCommitsCmd.Flags().IntVar(&parseCommitsLast, "last", 0, "Parse last N commits (alternative to --since)")
	parseCommitsCmd.Flags().StringVar(&parseCommitsPath, "path", "", "Only include commits touching this path")
	parseCommitsCmd.Flags().BoolVar(&parseCommitsNoFiles, "no-files", false, "Exclude file list from output")
	parseCommitsCmd.Flags().BoolVar(&parseCommitsNoMerges, "no-merges", false, "Exclude merge commits")
	parseCommitsCmd.Flags().StringVar(&parseCommitsFormat, "format", "toon", "Output format: toon (default), json, json-compact")
	parseCommitsCmd.Flags().StringVar(&parseCommitsRepoURL, "repo", "", "Repository URL to include in output")
	parseCommitsCmd.Flags().StringVar(&parseCommitsChangelog, "changelog", "", "CHANGELOG.json to read maintainers/bots for external contributor detection")
	parseCommitsCmd.Flags().BoolVar(&parseCommitsAllVersions, "all-versions", false, "Parse commits for all version ranges (outputs array of results)")
	rootCmd.AddCommand(parseCommitsCmd)
}

func runParseCommits(cmd *cobra.Command, args []string) error {
	// Handle --all-versions mode
	if parseCommitsAllVersions {
		return runParseAllVersions()
	}

	// Build git log command
	gitArgs := buildGitLogArgs()

	// Run git log
	output, err := runGitLog(gitArgs)
	if err != nil {
		return err
	}

	// Parse output
	parser := gitlog.NewParser()
	parser.IncludeFiles = !parseCommitsNoFiles

	result, err := parser.Parse(output)
	if err != nil {
		return fmt.Errorf("failed to parse git log output: %w", err)
	}

	// Set metadata
	if parseCommitsRepoURL != "" {
		result.Repository = parseCommitsRepoURL
	} else {
		// Try to get repository URL from git
		if repoURL, err := getRepositoryURL(); err == nil {
			result.Repository = repoURL
		}
	}

	result.Range.Since = parseCommitsSince
	result.Range.Until = parseCommitsUntil

	// If no-files flag, clear file lists from commits
	if parseCommitsNoFiles {
		for i := range result.Commits {
			result.Commits[i].Files = nil
		}
	}

	// Load changelog for external contributor detection
	var cl *changelog.Changelog
	if parseCommitsChangelog != "" {
		cl, err = changelog.LoadFile(parseCommitsChangelog)
		if err != nil {
			return fmt.Errorf("failed to load changelog %s: %w", parseCommitsChangelog, err)
		}
	}

	// Mark external contributors
	if cl != nil {
		for i := range result.Commits {
			c := &result.Commits[i]
			// IsExternal = true if author is NOT a team member
			c.IsExternal = !cl.IsTeamMemberByNameAndEmail(c.Author, c.AuthorEmail)
		}
	}

	// Compute contributors summary
	result.ComputeContributors()

	// Parse output format
	f, err := format.Parse(parseCommitsFormat)
	if err != nil {
		return err
	}

	// Output in specified format
	outputBytes, err := format.Marshal(result, f)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Println(string(outputBytes))
	return nil
}

func buildGitLogArgs() []string {
	args := []string{
		"log",
		"--format=" + gitlog.GitLogFormat,
		"--numstat",
	}

	if parseCommitsNoMerges {
		args = append(args, "--no-merges")
	}

	if parseCommitsLast > 0 {
		args = append(args, fmt.Sprintf("-n%d", parseCommitsLast))
	} else if parseCommitsSince != "" {
		args = append(args, fmt.Sprintf("%s..%s", parseCommitsSince, parseCommitsUntil))
	} else if parseCommitsUntil != "" && parseCommitsUntil != "HEAD" {
		// If only --until is specified (no --since), get all commits up to that ref
		args = append(args, parseCommitsUntil)
	}

	if parseCommitsPath != "" {
		args = append(args, "--", parseCommitsPath)
	}

	return args
}

func runGitLog(args []string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git log failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run git log: %w", err)
	}
	return string(output), nil
}

func getRepositoryURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	url := strings.TrimSpace(string(output))

	// Convert SSH URL to HTTPS
	if strings.HasPrefix(url, "git@") {
		// git@github.com:owner/repo.git -> github.com/owner/repo
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		url = strings.TrimSuffix(url, ".git")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
		url = strings.TrimSuffix(url, ".git")
	}

	return url, nil
}

// AllVersionsResult contains parse results for all version ranges.
type AllVersionsResult struct {
	Repository  string               `json:"repository,omitempty"`
	Versions    []VersionParseResult `json:"versions"`
	TotalCount  int                  `json:"total_count"`
	GeneratedAt string               `json:"generated_at"`
}

// VersionParseResult contains parse result for a single version.
type VersionParseResult struct {
	Version     string          `json:"version"`
	Date        string          `json:"date"`
	Since       string          `json:"since,omitempty"`
	CommitCount int             `json:"commit_count"`
	Commits     []gitlog.Commit `json:"commits"`
	Summary     gitlog.Summary  `json:"summary"`
}

// runParseAllVersions parses commits for all version ranges at once.
func runParseAllVersions() error {
	// Get all version ranges
	ranges, err := gitlog.GetAllVersionRanges()
	if err != nil {
		return fmt.Errorf("failed to get version ranges: %w", err)
	}

	if len(ranges) == 0 {
		return fmt.Errorf("no semver tags found in repository")
	}

	// Get repository URL
	repoURL := parseCommitsRepoURL
	if repoURL == "" {
		if url, err := getRepositoryURL(); err == nil {
			repoURL = url
		}
	}

	// Load changelog for external contributor detection
	var cl *changelog.Changelog
	if parseCommitsChangelog != "" {
		cl, err = changelog.LoadFile(parseCommitsChangelog)
		if err != nil {
			return fmt.Errorf("failed to load changelog %s: %w", parseCommitsChangelog, err)
		}
	}

	// Parse commits for each version
	result := AllVersionsResult{
		Repository:  repoURL,
		Versions:    make([]VersionParseResult, 0, len(ranges)),
		GeneratedAt: gitlog.NewParseResult().GeneratedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	totalCommits := 0
	for _, vr := range ranges {
		// Build git args for this range
		var args []string
		if vr.Since == "" {
			args = []string{"log", "--format=" + gitlog.GitLogFormat, "--numstat", vr.Until}
		} else {
			args = []string{"log", "--format=" + gitlog.GitLogFormat, "--numstat", fmt.Sprintf("%s..%s", vr.Since, vr.Until)}
		}

		if parseCommitsNoMerges {
			args = append(args, "--no-merges")
		}

		output, err := runGitLog(args)
		if err != nil {
			// Skip versions we can't parse
			continue
		}

		parser := gitlog.NewParser()
		parser.IncludeFiles = !parseCommitsNoFiles

		parseResult, err := parser.Parse(output)
		if err != nil {
			continue
		}

		// Mark external contributors
		if cl != nil {
			for i := range parseResult.Commits {
				c := &parseResult.Commits[i]
				c.IsExternal = !cl.IsTeamMemberByNameAndEmail(c.Author, c.AuthorEmail)
			}
		}

		// Clear file lists if requested
		if parseCommitsNoFiles {
			for i := range parseResult.Commits {
				parseResult.Commits[i].Files = nil
			}
		}

		vpr := VersionParseResult{
			Version:     vr.Version,
			Date:        vr.Date,
			Since:       vr.Since,
			CommitCount: len(parseResult.Commits),
			Commits:     parseResult.Commits,
			Summary:     parseResult.Summary,
		}

		result.Versions = append(result.Versions, vpr)
		totalCommits += len(parseResult.Commits)
	}

	result.TotalCount = totalCommits

	// Parse output format
	f, err := format.Parse(parseCommitsFormat)
	if err != nil {
		return err
	}

	// Output in specified format
	outputBytes, err := format.Marshal(result, f)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Println(string(outputBytes))
	return nil
}
