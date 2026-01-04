package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
)

var (
	validateStrict   bool
	validateWarnings bool
	validateMinTier  string
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a CHANGELOG.json file",
	Long: `Validate a Structured Changelog JSON file against the IR schema.

Checks for:
  - Required fields (ir_version, project)
  - Valid semantic versions
  - Valid date formats (YYYY-MM-DD)
  - Valid security metadata (CVE, GHSA, severity)
  - No duplicate versions
  - Non-empty descriptions

Tier validation:
  --min-tier     Require at least one entry in a category at or above this tier

Tiers:
  core       KACL standard types (Security, Added, Changed, Deprecated, Removed, Fixed)
  standard   Commonly used types (core + Highlights, Breaking, Upgrade Guide, Performance, Dependencies)
  extended   Extended types (standard + Documentation, Build, Known Issues, Contributors)
  optional   All types (extended + Infrastructure, Observability, Compliance, Internal)

Examples:
  sclog validate CHANGELOG.json
  sclog validate CHANGELOG.json --strict
  sclog validate CHANGELOG.json --min-tier core`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Enable strict validation (treat warnings as errors)")
	validateCmd.Flags().BoolVar(&validateWarnings, "warnings", true, "Show warnings")
	validateCmd.Flags().StringVar(&validateMinTier, "min-tier", "", "Minimum tier to require coverage for (core, standard, extended, optional)")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Load changelog
	cl, err := changelog.LoadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", inputFile, err)
	}

	// Validate
	result := cl.Validate()

	if !result.Valid {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "  ✗ %s\n", e.Error())
		}
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	// Validate min tier if specified
	if validateMinTier != "" {
		tier, err := changelog.ParseTier(validateMinTier)
		if err != nil {
			return fmt.Errorf("invalid tier %q: must be one of core, standard, extended, optional", validateMinTier)
		}
		if err := cl.ValidateMinTier(tier); err != nil {
			return fmt.Errorf("tier validation failed: %w", err)
		}
	}

	fmt.Printf("✓ %s is valid\n", inputFile)

	// Print summary
	printSummary(cl)

	return nil
}

func printSummary(cl *changelog.Changelog) {
	s := cl.Summary()

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Project: %s\n", s.Project)
	fmt.Printf("  IR Version: %s\n", s.IRVersion)
	fmt.Printf("  Releases: %d\n", s.ReleaseCount)

	if s.HasUnreleased {
		fmt.Printf("  Unreleased: yes\n")
		if len(s.UnreleasedCategories) > 0 {
			fmt.Printf("  Unreleased categories: %s\n", strings.Join(s.UnreleasedCategories, ", "))
		}
	}

	if s.LatestVersion != "" {
		fmt.Printf("  Latest: %s (%s)\n", s.LatestVersion, s.LatestDate)
		if len(s.LatestCategories) > 0 {
			fmt.Printf("  Categories: %s\n", strings.Join(s.LatestCategories, ", "))
		}
	}
}
