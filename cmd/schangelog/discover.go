package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/aggregate"
)

var (
	discoverOrgs     []string
	discoverUsers    []string
	discoverManifest string
	discoverOutput   string
	discoverUpdate   bool
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover projects with CHANGELOG.json in GitHub orgs/users",
	Long: `Scan GitHub organizations and users for repositories containing CHANGELOG.json.

Requires GITHUB_TOKEN environment variable for authentication.

Discovery scans all non-archived, non-fork repositories in the specified
organizations and users, looking for CHANGELOG.json files (including nested
paths in monorepos).

Examples:
  # Discover from orgs and users, create new manifest
  schangelog portfolio discover --org agentplexus --org fleet-ops --user grokify -o manifest.json

  # Update existing manifest with new discoveries
  schangelog portfolio discover --manifest manifest.json --update

  # Add new sources to existing manifest
  schangelog portfolio discover --manifest manifest.json --org new-org --update`,
	RunE: runDiscover,
}

func init() {
	discoverCmd.Flags().StringArrayVar(&discoverOrgs, "org", nil, "GitHub organization to scan (can be specified multiple times)")
	discoverCmd.Flags().StringArrayVar(&discoverUsers, "user", nil, "GitHub user to scan (can be specified multiple times)")
	discoverCmd.Flags().StringVar(&discoverManifest, "manifest", "", "Existing manifest to update")
	discoverCmd.Flags().StringVarP(&discoverOutput, "output", "o", "", "Output file (default: stdout or manifest file if --update)")
	discoverCmd.Flags().BoolVar(&discoverUpdate, "update", false, "Update existing manifest in place")
	portfolioCmd.AddCommand(discoverCmd)
}

func runDiscover(cmd *cobra.Command, args []string) error {
	var manifest *aggregate.Manifest
	var err error

	// Load or create manifest
	if discoverManifest != "" {
		manifest, err = aggregate.LoadManifest(discoverManifest)
		if err != nil {
			return fmt.Errorf("loading manifest: %w", err)
		}
	} else {
		manifest = aggregate.NewManifest("Discovered Projects")
	}

	// Build sources list
	var sources []aggregate.Source
	for _, org := range discoverOrgs {
		sources = append(sources, aggregate.Source{Type: aggregate.SourceTypeOrg, Name: org})
		manifest.AddSource(aggregate.Source{Type: aggregate.SourceTypeOrg, Name: org})
	}
	for _, user := range discoverUsers {
		sources = append(sources, aggregate.Source{Type: aggregate.SourceTypeUser, Name: user})
		manifest.AddSource(aggregate.Source{Type: aggregate.SourceTypeUser, Name: user})
	}

	// If no new sources specified, use manifest sources
	if len(sources) == 0 && len(manifest.Sources) > 0 {
		sources = manifest.Sources
	}

	if len(sources) == 0 {
		return fmt.Errorf("no sources specified (use --org or --user, or provide a manifest with sources)")
	}

	// Create discovery client
	client, err := aggregate.NewDiscoveryClient("")
	if err != nil {
		return err
	}

	// Discover projects
	fmt.Fprintf(os.Stderr, "Scanning %d source(s)...\n", len(sources))
	discovered, err := client.DiscoverProjects(cmd.Context(), sources)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	// Add discovered projects to manifest
	added := 0
	for _, proj := range discovered {
		if manifest.AddProject(proj) {
			added++
		}
	}

	fmt.Fprintf(os.Stderr, "Discovered %d projects, %d new\n", len(discovered), added)

	// Update generated timestamp
	now := time.Now().UTC()
	manifest.Generated = &now

	// Output
	output, err := manifest.JSON()
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	// Determine output destination
	outputPath := discoverOutput
	if outputPath == "" && discoverUpdate && discoverManifest != "" {
		outputPath = discoverManifest
	}

	if outputPath == "" {
		fmt.Println(string(output))
	} else {
		if err := os.WriteFile(outputPath, output, 0600); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote manifest to %s\n", outputPath)
	}

	// Print summary
	fmt.Fprintf(os.Stderr, "\nManifest summary:\n")
	fmt.Fprintf(os.Stderr, "  Sources: %d\n", len(manifest.Sources))
	fmt.Fprintf(os.Stderr, "  Projects: %d\n", len(manifest.Projects))

	return nil
}
