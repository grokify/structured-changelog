package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-changelog/format"
)

var (
	validateStrict   bool
	validateWarnings bool
	validateMinTier  string
	validateFormat   string
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

Output formats (with --format flag):
  - toon: Token-Oriented Object Notation, ~40% fewer tokens than JSON
  - json: Standard JSON with indentation
  - json-compact: Minified JSON

Tier validation:
  --min-tier     Require at least one entry in a category at or above this tier

Tiers:
  core       KACL standard types (Security, Added, Changed, Deprecated, Removed, Fixed)
  standard   Commonly used types (core + Highlights, Breaking, Upgrade Guide, Performance, Dependencies)
  extended   Extended types (standard + Documentation, Build, Known Issues, Contributors)
  optional   All types (extended + Infrastructure, Observability, Compliance, Internal)

Examples:
  schangelog validate CHANGELOG.json
  schangelog validate CHANGELOG.json --strict
  schangelog validate CHANGELOG.json --min-tier core
  schangelog validate CHANGELOG.json --format=toon`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Enable strict validation (treat warnings as errors)")
	validateCmd.Flags().BoolVar(&validateWarnings, "warnings", true, "Show warnings")
	validateCmd.Flags().StringVar(&validateMinTier, "min-tier", "", "Minimum tier to require coverage for (core, standard, extended, optional)")
	validateCmd.Flags().StringVar(&validateFormat, "format", "", "Output format: toon, json, json-compact (enables structured output)")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Load changelog
	cl, err := changelog.LoadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", inputFile, err)
	}

	// Use rich validation for structured output
	if validateFormat != "" {
		return runValidateStructured(cl, inputFile)
	}

	// Standard validation
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

func runValidateStructured(cl *changelog.Changelog, _ string) error {
	result := cl.ValidateRich()

	// Add tier validation as warning if specified
	if validateMinTier != "" {
		tier, err := changelog.ParseTier(validateMinTier)
		if err != nil {
			result.Warnings = append(result.Warnings, changelog.RichValidationError{
				Code:       changelog.WarnCodeNoTierCoverage,
				Severity:   changelog.SeverityWarning,
				Path:       "min-tier",
				Message:    fmt.Sprintf("Invalid tier %q specified", validateMinTier),
				Suggestion: "Use one of: core, standard, extended, optional",
			})
		} else if err := cl.ValidateMinTier(tier); err != nil {
			if validateStrict {
				result.Valid = false
				result.Errors = append(result.Errors, changelog.RichValidationError{
					Code:       changelog.WarnCodeNoTierCoverage,
					Severity:   changelog.SeverityError,
					Path:       "releases[0]",
					Message:    fmt.Sprintf("No entries at or above tier %q", tier),
					Suggestion: fmt.Sprintf("Add at least one entry in a %s-tier category", tier),
				})
			} else {
				result.Warnings = append(result.Warnings, changelog.RichValidationError{
					Code:       changelog.WarnCodeNoTierCoverage,
					Severity:   changelog.SeverityWarning,
					Path:       "releases[0]",
					Message:    fmt.Sprintf("No entries at or above tier %q", tier),
					Suggestion: fmt.Sprintf("Add at least one entry in a %s-tier category", tier),
				})
			}
		}
	}

	// In strict mode, treat warnings as errors
	if validateStrict && len(result.Warnings) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, result.Warnings...)
		result.Warnings = nil
	}

	// Filter warnings if disabled
	if !validateWarnings {
		result.Warnings = nil
	}

	// Update summary counts
	result.Summary.ErrorCount = len(result.Errors)
	result.Summary.WarningCount = len(result.Warnings)

	// Parse output format
	f, err := format.Parse(validateFormat)
	if err != nil {
		return err
	}

	output, err := format.Marshal(result, f)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Println(string(output))

	if !result.Valid {
		return fmt.Errorf("validation failed")
	}
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
