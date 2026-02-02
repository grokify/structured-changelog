package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-changelog/gitlog"
)

var (
	initFromTags   bool
	initOutput     string
	initProject    string
	initRepoURL    string
	initVersioning string
	initConvention string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a CHANGELOG.json from git history",
	Long: `Initialize a CHANGELOG.json file by scanning git history.

With --from-tags, this command creates a skeleton CHANGELOG.json with
all semver tags as releases, including dates and placeholder entries
based on commit analysis.

This is useful for:
  - Starting a new structured changelog for an existing project
  - Backfilling changelog history from git tags
  - Creating a template that can be manually refined

Examples:
  # Create skeleton from git tags
  schangelog init --from-tags

  # Specify project name and output file
  schangelog init --from-tags --project=myproject -o CHANGELOG.json

  # Set versioning and commit convention
  schangelog init --from-tags --versioning=semver --convention=conventional`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initFromTags, "from-tags", false, "Generate changelog from git tags (required)")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "", "Output file (default: stdout)")
	initCmd.Flags().StringVar(&initProject, "project", "", "Project name (default: derived from repo URL)")
	initCmd.Flags().StringVar(&initRepoURL, "repo", "", "Repository URL")
	initCmd.Flags().StringVar(&initVersioning, "versioning", "semver", "Versioning scheme: semver, calver, custom, none")
	initCmd.Flags().StringVar(&initConvention, "convention", "conventional", "Commit convention: conventional, none")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	if !initFromTags {
		return fmt.Errorf("--from-tags is required (other modes not yet implemented)")
	}

	return runInitFromTags()
}

func runInitFromTags() error {
	// Get repository URL
	repoURL := initRepoURL
	if repoURL == "" {
		if url, err := getRepositoryURL(); err == nil {
			repoURL = url
		}
	}

	// Derive project name from repo URL if not specified
	projectName := initProject
	if projectName == "" && repoURL != "" {
		parts := strings.Split(repoURL, "/")
		if len(parts) > 0 {
			projectName = parts[len(parts)-1]
		}
	}

	// Get all tags
	tagList, err := gitlog.GetTags()
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	if len(tagList.Tags) == 0 {
		return fmt.Errorf("no semver tags found in repository")
	}

	// Create changelog structure
	cl := &changelog.Changelog{
		IRVersion:        "1.0",
		Project:          projectName,
		Repository:       repoURL,
		Versioning:       initVersioning,
		CommitConvention: initConvention,
		Releases:         make([]changelog.Release, 0, len(tagList.Tags)),
	}

	// Process each tag (in reverse order - newest first)
	for i := len(tagList.Tags) - 1; i >= 0; i-- {
		tag := tagList.Tags[i]

		// Determine since ref for parsing commits
		var sinceRef string
		if i > 0 {
			sinceRef = tagList.Tags[i-1].Name
		}

		// Parse commits for this version
		commits, err := parseCommitsForVersion(sinceRef, tag.Name)
		if err != nil {
			// If we can't parse commits, create minimal release entry
			cl.Releases = append(cl.Releases, changelog.Release{
				Version: tag.Name,
				Date:    tag.DateString,
			})
			continue
		}

		// Build release from commits
		release := buildReleaseFromCommits(tag.Name, tag.DateString, commits)
		cl.Releases = append(cl.Releases, release)
	}

	// Marshal to JSON
	output, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal changelog: %w", err)
	}

	// Write output
	if initOutput != "" {
		if err := os.WriteFile(initOutput, output, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Created %s with %d releases\n", initOutput, len(cl.Releases))
	} else {
		fmt.Println(string(output))
	}

	return nil
}

// parseCommitsForVersion parses commits between two refs.
func parseCommitsForVersion(since, until string) ([]gitlog.Commit, error) {
	var args []string
	if since == "" {
		args = []string{"log", "--format=" + gitlog.GitLogFormat, until}
	} else {
		args = []string{"log", "--format=" + gitlog.GitLogFormat, fmt.Sprintf("%s..%s", since, until)}
	}

	output, err := runGitLog(args)
	if err != nil {
		return nil, err
	}

	parser := gitlog.NewParser()
	parser.IncludeFiles = false

	result, err := parser.Parse(output)
	if err != nil {
		return nil, err
	}

	return result.Commits, nil
}

// buildReleaseFromCommits creates a Release from parsed commits.
func buildReleaseFromCommits(version, date string, commits []gitlog.Commit) changelog.Release {
	release := changelog.Release{
		Version: version,
		Date:    date,
	}

	// Group commits by suggested category
	for _, commit := range commits {
		entry := changelog.Entry{
			Description: commit.Subject,
			Commit:      commit.ShortHash,
		}

		if commit.Issue > 0 {
			entry.Issue = fmt.Sprintf("%d", commit.Issue)
		}
		if commit.PR > 0 {
			entry.PR = fmt.Sprintf("%d", commit.PR)
		}
		if commit.Breaking {
			entry.Breaking = true
		}

		// Add to appropriate category based on suggested category
		switch commit.SuggestedCategory {
		case "Added":
			release.Added = append(release.Added, entry)
		case "Changed":
			release.Changed = append(release.Changed, entry)
		case "Deprecated":
			release.Deprecated = append(release.Deprecated, entry)
		case "Removed":
			release.Removed = append(release.Removed, entry)
		case "Fixed":
			release.Fixed = append(release.Fixed, entry)
		case "Security":
			release.Security = append(release.Security, entry)
		case "Documentation":
			release.Documentation = append(release.Documentation, entry)
		case "Dependencies":
			release.Dependencies = append(release.Dependencies, entry)
		case "Build":
			release.Build = append(release.Build, entry)
		case "Performance":
			release.Performance = append(release.Performance, entry)
		case "Internal":
			release.Internal = append(release.Internal, entry)
		case "Infrastructure":
			release.Infrastructure = append(release.Infrastructure, entry)
		default:
			// Default to Changed if no category
			if commit.Type == "feat" {
				release.Added = append(release.Added, entry)
			} else if commit.Type == "fix" {
				release.Fixed = append(release.Fixed, entry)
			} else {
				release.Changed = append(release.Changed, entry)
			}
		}
	}

	return release
}
