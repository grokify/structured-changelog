package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
)

var (
	validateStrict   bool
	validateWarnings bool
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

Examples:
  sclog validate CHANGELOG.json
  sclog validate CHANGELOG.json --strict`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Enable strict validation (treat warnings as errors)")
	validateCmd.Flags().BoolVar(&validateWarnings, "warnings", true, "Show warnings")
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

	fmt.Printf("✓ %s is valid\n", inputFile)

	// Print summary
	printSummary(cl)

	return nil
}

func printSummary(cl *changelog.Changelog) {
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Project: %s\n", cl.Project)
	fmt.Printf("  IR Version: %s\n", cl.IRVersion)
	fmt.Printf("  Releases: %d\n", len(cl.Releases))

	if cl.Unreleased != nil && !cl.Unreleased.IsEmpty() {
		fmt.Printf("  Unreleased: yes\n")
	}

	if len(cl.Releases) > 0 {
		latest := cl.Releases[0]
		fmt.Printf("  Latest: %s (%s)\n", latest.Version, latest.Date)
	}
}
