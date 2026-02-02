package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/format"
	"github.com/grokify/structured-changelog/gitlog"
)

var (
	listTagsFormat  string
	listTagsRepoURL string
)

var listTagsCmd = &cobra.Command{
	Use:   "list-tags",
	Short: "List semver tags with dates and commit counts",
	Long: `List all semantic version tags in the repository with their dates
and the number of commits between each version.

This command is useful for understanding the release history and
planning changelog backfills.

Output formats:
  - toon (default): Token-Oriented Object Notation
  - json: Standard JSON with indentation
  - json-compact: Minified JSON

Examples:
  # List all tags (TOON format, default)
  schangelog list-tags

  # List tags with JSON output
  schangelog list-tags --format=json

  # Include repository URL in output
  schangelog list-tags --repo=github.com/owner/repo`,
	RunE: runListTags,
}

func init() {
	listTagsCmd.Flags().StringVar(&listTagsFormat, "format", "toon", "Output format: toon (default), json, json-compact")
	listTagsCmd.Flags().StringVar(&listTagsRepoURL, "repo", "", "Repository URL to include in output")
	rootCmd.AddCommand(listTagsCmd)
}

func runListTags(cmd *cobra.Command, args []string) error {
	// Get tags
	tagList, err := gitlog.GetTags()
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	// Set repository URL
	if listTagsRepoURL != "" {
		tagList.Repository = listTagsRepoURL
	} else {
		if repoURL, err := getRepositoryURL(); err == nil {
			tagList.Repository = repoURL
		}
	}

	// Parse output format
	f, err := format.Parse(listTagsFormat)
	if err != nil {
		return err
	}

	// Output in specified format
	outputBytes, err := format.Marshal(tagList, f)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Println(string(outputBytes))
	return nil
}
