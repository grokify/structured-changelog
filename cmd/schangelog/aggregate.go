package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/aggregate"
)

var (
	aggregateOutput  string
	aggregatePending bool
	aggregateApprove bool
)

var aggregateCmd = &cobra.Command{
	Use:   "aggregate <manifest.json>",
	Short: "Aggregate changelogs into a unified portfolio",
	Long: `Aggregate changelogs from multiple projects defined in a manifest file.

The manifest file is a JSON file listing projects to aggregate:
  {
    "name": "Portfolio Name",
    "projects": [
      {"path": "github.com/org/repo1"},
      {"path": "github.com/org/repo2/subdir"}
    ]
  }

Projects are resolved by searching:
  1. LocalPath override in the manifest
  2. ~/go/src/<path>
  3. $GOPATH/src/<path>

Projects not found locally can be fetched from GitHub with --approve.

Examples:
  # Aggregate from manifest (local only)
  schangelog portfolio aggregate manifest.json -o portfolio.json

  # Show projects needing remote fetch
  schangelog portfolio aggregate manifest.json --pending

  # Approve and fetch remote projects
  schangelog portfolio aggregate manifest.json --approve -o portfolio.json`,
	Args: cobra.ExactArgs(1),
	RunE: runAggregate,
}

func init() {
	aggregateCmd.Flags().StringVarP(&aggregateOutput, "output", "o", "", "Output file (default: stdout)")
	aggregateCmd.Flags().BoolVar(&aggregatePending, "pending", false, "Show projects needing remote fetch")
	aggregateCmd.Flags().BoolVar(&aggregateApprove, "approve", false, "Fetch remote changelogs (requires GITHUB_TOKEN)")
	portfolioCmd.AddCommand(aggregateCmd)
}

func runAggregate(cmd *cobra.Command, args []string) error {
	manifestPath := args[0]

	// Load manifest
	manifest, err := aggregate.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	// Validate manifest
	result := manifest.Validate()
	if !result.Valid {
		fmt.Fprintf(os.Stderr, "Manifest validation failed:\n")
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "  ✗ %s: %s\n", e.Field, e.Message)
		}
		return fmt.Errorf("manifest validation failed")
	}

	// Resolve projects
	resolver := aggregate.NewResolver()
	resolved, err := resolver.Resolve(manifest.Projects)
	if err != nil {
		return fmt.Errorf("resolving projects: %w", err)
	}

	// Show summary
	summary := resolver.Summary(resolved)
	fmt.Fprintf(os.Stderr, "Resolved %d/%d projects locally\n", summary.LocalCount, summary.TotalCount)

	// Handle --pending flag
	if aggregatePending {
		if summary.RemoteCount == 0 {
			fmt.Println("All projects resolved locally.")
			return nil
		}

		fmt.Printf("\nProjects needing remote fetch (%d):\n", summary.RemoteCount)
		for _, path := range summary.RemoteProjects {
			fmt.Printf("  - %s\n", path)
		}
		fmt.Printf("\nRun with --approve to fetch these changelogs from GitHub.\n")
		return nil
	}

	// Handle --approve flag for remote fetch
	if aggregateApprove && summary.RemoteCount > 0 {
		fmt.Fprintf(os.Stderr, "Fetching %d remote changelogs...\n", summary.RemoteCount)

		client, err := aggregate.NewDiscoveryClient("")
		if err != nil {
			return fmt.Errorf("creating GitHub client: %w", err)
		}

		remoteProjects := aggregate.FilterRemote(resolved)
		for i := range remoteProjects {
			rp := &remoteProjects[i]
			data, err := client.FetchRemoteChangelog(cmd.Context(), rp.Ref.Path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  ⚠ Failed to fetch %s: %v\n", rp.Ref.Path, err)
				continue
			}

			// Write to temp file for loading
			tmpFile, err := os.CreateTemp("", "changelog-*.json")
			if err != nil {
				return fmt.Errorf("creating temp file: %w", err)
			}
			tmpPath := tmpFile.Name()
			if _, err := tmpFile.Write(data); err != nil {
				tmpFile.Close()
				os.Remove(tmpPath) //nolint:gosec // tmpPath is from os.CreateTemp, not user input
				return fmt.Errorf("writing temp file: %w", err)
			}
			tmpFile.Close()

			rp.ChangelogPath = tmpPath
			rp.IsLocal = true
			rp.NeedsApproval = false
			defer os.Remove(tmpPath) //nolint:gosec // tmpPath is from os.CreateTemp, not user input

			fmt.Fprintf(os.Stderr, "  ✓ Fetched %s\n", rp.Ref.Path)
		}

		// Update resolved list with fetched projects
		for i := range resolved {
			if resolved[i].NeedsApproval {
				for _, rp := range remoteProjects {
					if rp.Ref.Path == resolved[i].Ref.Path && rp.IsLocal {
						resolved[i] = rp
						break
					}
				}
			}
		}
	}

	// Load portfolio
	localResolved := aggregate.FilterLocal(resolved)
	if len(localResolved) == 0 {
		return fmt.Errorf("no projects resolved locally")
	}

	portfolio, err := aggregate.LoadPortfolio(manifest, localResolved)
	if err != nil {
		return fmt.Errorf("loading portfolio: %w", err)
	}

	// Output
	output, err := portfolio.JSON()
	if err != nil {
		return fmt.Errorf("marshaling portfolio: %w", err)
	}

	if aggregateOutput == "" {
		fmt.Println(string(output))
	} else {
		if err := os.WriteFile(aggregateOutput, output, 0600); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote portfolio to %s\n", aggregateOutput)
	}

	// Print summary
	s := portfolio.Summary()
	fmt.Fprintf(os.Stderr, "\nPortfolio summary:\n")
	fmt.Fprintf(os.Stderr, "  Projects: %d\n", s.ProjectCount)
	fmt.Fprintf(os.Stderr, "  Releases: %d\n", s.ReleaseCount)
	fmt.Fprintf(os.Stderr, "  Entries:  %d\n", s.EntryCount)
	if s.DateRange.Start != "" {
		fmt.Fprintf(os.Stderr, "  Date range: %s to %s\n", s.DateRange.Start, s.DateRange.End)
	}

	return nil
}
