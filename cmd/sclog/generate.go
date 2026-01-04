package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-changelog/renderer"
)

var (
	generateOutput  string
	generateMinimal bool
	generateFull    bool
)

var generateCmd = &cobra.Command{
	Use:   "generate <file>",
	Short: "Generate CHANGELOG.md from CHANGELOG.json",
	Long: `Generate a Keep a Changelog formatted Markdown file from a
Structured Changelog JSON file.

The output is deterministic: the same input always produces identical output.

Output options:
  --minimal   Exclude references and security metadata
  --full      Include all metadata including commit SHAs

Examples:
  sclog generate CHANGELOG.json
  sclog generate CHANGELOG.json -o CHANGELOG.md
  sclog generate CHANGELOG.json --minimal
  sclog generate CHANGELOG.json --full -o docs/CHANGELOG.md`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "Output file (default: stdout)")
	generateCmd.Flags().BoolVar(&generateMinimal, "minimal", false, "Use minimal output (no references/metadata)")
	generateCmd.Flags().BoolVar(&generateFull, "full", false, "Use full output (include commits)")
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Load changelog
	cl, err := changelog.LoadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", inputFile, err)
	}

	// Validate first
	result := cl.Validate()
	if !result.Valid {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "  ✗ %s\n", e.Error())
		}
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	// Select options
	opts := renderer.DefaultOptions()
	if generateMinimal {
		opts = renderer.MinimalOptions()
	} else if generateFull {
		opts = renderer.FullOptions()
	}

	// Render
	md := renderer.RenderMarkdownWithOptions(cl, opts)

	// Write output
	if generateOutput == "" {
		// Write to stdout
		fmt.Print(md)
	} else {
		if err := os.WriteFile(generateOutput, []byte(md), 0644); err != nil { //nolint:gosec // 0644 intentional for readable output
			return fmt.Errorf("failed to write %s: %w", generateOutput, err)
		}
		fmt.Fprintf(os.Stderr, "Generated %s from %s\n", generateOutput, inputFile)
	}

	return nil
}
