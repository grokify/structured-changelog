package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
)

var (
	mergeOutput      string
	mergeRelease     string
	mergeDedup       bool
	mergePrependOnly bool
)

var mergeCmd = &cobra.Command{
	Use:   "merge <base> [additions...]",
	Short: "Merge changelog files together",
	Long: `Merge one or more changelog files into a base changelog.

This command combines releases from multiple CHANGELOG.json files while
preserving metadata (maintainers, bots, repository) from the base file.

Use cases:
  - Merge a newly generated release into an existing changelog
  - Combine changelogs from multiple sources
  - Update an existing changelog with backfilled history

Examples:
  # Merge two changelog files
  schangelog merge base.json additions.json -o CHANGELOG.json

  # Merge multiple files
  schangelog merge base.json v2.27.json v2.26.json -o CHANGELOG.json

  # Prepend a single release file
  schangelog merge CHANGELOG.json --release new-release.json -o CHANGELOG.json

  # Merge with deduplication (skip versions that already exist in base)
  schangelog merge base.json additions.json --dedup -o CHANGELOG.json`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "Output file (default: stdout)")
	mergeCmd.Flags().StringVar(&mergeRelease, "release", "", "Single release file to prepend")
	mergeCmd.Flags().BoolVar(&mergeDedup, "dedup", false, "Skip versions that already exist in base")
	mergeCmd.Flags().BoolVar(&mergePrependOnly, "prepend-only", false, "Only add releases newer than base's latest")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	// Load base changelog
	basePath := args[0]
	base, err := changelog.LoadFile(basePath)
	if err != nil {
		return fmt.Errorf("failed to load base changelog %s: %w", basePath, err)
	}

	// Track existing versions for deduplication
	existingVersions := make(map[string]bool)
	for _, r := range base.Releases {
		existingVersions[r.Version] = true
	}

	var releasesToMerge []changelog.Release
	var duplicates []string

	// Handle --release flag for single release file
	if mergeRelease != "" {
		releaseFile, err := changelog.LoadFile(mergeRelease)
		if err != nil {
			return fmt.Errorf("failed to load release file %s: %w", mergeRelease, err)
		}
		releasesToMerge = append(releasesToMerge, releaseFile.Releases...)
	}

	// Handle additional changelog files
	for _, addPath := range args[1:] {
		addFile, err := changelog.LoadFile(addPath)
		if err != nil {
			return fmt.Errorf("failed to load changelog %s: %w", addPath, err)
		}
		releasesToMerge = append(releasesToMerge, addFile.Releases...)
	}

	// Merge releases
	for _, r := range releasesToMerge {
		if existingVersions[r.Version] {
			if mergeDedup {
				duplicates = append(duplicates, r.Version)
				continue
			}
			fmt.Fprintf(os.Stderr, "Warning: version %s exists in both base and additions\n", r.Version)
		}

		// For prepend-only mode, skip if version is older than latest
		if mergePrependOnly && len(base.Releases) > 0 {
			// Simple string comparison works for semver in most cases
			// For proper semver comparison, we'd need a semver library
			if r.Version <= base.Releases[0].Version {
				continue
			}
		}

		// Prepend the release
		base.Releases = append([]changelog.Release{r}, base.Releases...)
		existingVersions[r.Version] = true
	}

	// Report skipped duplicates
	if len(duplicates) > 0 {
		fmt.Fprintf(os.Stderr, "Skipped %d duplicate versions: %v\n", len(duplicates), duplicates)
	}

	// Marshal to JSON
	output, err := json.MarshalIndent(base, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal changelog: %w", err)
	}

	// Write output
	if mergeOutput != "" {
		if err := os.WriteFile(mergeOutput, output, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Merged changelog written to %s (%d releases)\n", mergeOutput, len(base.Releases))
	} else {
		fmt.Println(string(output))
	}

	return nil
}
