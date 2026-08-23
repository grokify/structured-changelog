package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-changelog/renderer"
)

var (
	generateOutput             string
	generateMinimal            bool
	generateFull               bool
	generateMaxTier            string
	generateLocale             string
	generateLocaleFile         string
	generateAllReleases        bool
	generateNotableCategories  string
	generateExcludeCategories  string
	generateCollapseCategories string
	generateExpandCategories   string
)

var generateCmd = &cobra.Command{
	Use:   "generate <file>",
	Short: "Generate CHANGELOG.md from CHANGELOG.json",
	Long: `Generate a Keep a Changelog formatted Markdown file from a
Structured Changelog JSON file.

The output is deterministic: the same input always produces identical output.

By default, only notable releases are included (those with user-facing changes).
Use --all-releases to include maintenance-only releases.

By default, the Dependencies category is collapsed to a one-line summary
(e.g. "17 dependency updates") to keep the human-facing Markdown readable;
the full list remains in the JSON source. Use --full or
--expand-categories Dependencies to render every entry.

Output options:
  --minimal              Exclude references and security metadata (implies --max-tier core)
  --full                 Include all metadata and all releases (implies --all-releases; expands all categories)
  --max-tier             Filter change types by tier (core, standard, extended, optional)
  --locale               Output locale for localized strings (e.g., en, fr, de, es, ja, zh)
  --locale-file          Path to JSON file with locale message overrides
  --all-releases         Include all releases (overrides default notable-only behavior)
  --notable-categories   Custom notable categories (comma-separated)
  --exclude-categories   Categories to omit entirely (comma-separated)
  --collapse-categories  Categories to render as a one-line summary (comma-separated)
  --expand-categories    Categories to force-expand, overriding collapse defaults (comma-separated)

Tiers:
  core       KACL standard types (Security, Added, Changed, Deprecated, Removed, Fixed)
  standard   Commonly used types (core + Highlights, Breaking, Upgrade Guide, Performance, Dependencies)
  extended   Extended types (standard + Documentation, Build, Known Issues, Contributors)
  optional   All types (extended + Infrastructure, Observability, Compliance, Internal)

Notable Categories (default):
  Highlights, Breaking, Upgrade Guide, Security, Added, Changed, Deprecated,
  Removed, Fixed, Performance, Known Issues

Non-notable (maintenance) categories:
  Dependencies, Documentation, Build, Tests, Infrastructure, Observability,
  Compliance, Internal, Contributors

Examples:
  schangelog generate CHANGELOG.json
  schangelog generate CHANGELOG.json -o CHANGELOG.md
  schangelog generate CHANGELOG.json --minimal
  schangelog generate CHANGELOG.json --max-tier standard
  schangelog generate CHANGELOG.json --full -o docs/CHANGELOG.md
  schangelog generate CHANGELOG.json --locale=fr
  schangelog generate CHANGELOG.json --all-releases
  schangelog generate CHANGELOG.json --notable-categories "Security,Added,Fixed"
  schangelog generate CHANGELOG.json --exclude-categories "Dependencies,Build"
  schangelog generate CHANGELOG.json --expand-categories "Dependencies"`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "Output file (default: stdout)")
	generateCmd.Flags().BoolVar(&generateMinimal, "minimal", false, "Use minimal output (no references/metadata, core tier only)")
	generateCmd.Flags().BoolVar(&generateFull, "full", false, "Use full output (include commits and all releases)")
	generateCmd.Flags().StringVar(&generateMaxTier, "max-tier", "", "Maximum tier to include (core, standard, extended, optional)")
	generateCmd.Flags().StringVar(&generateLocale, "locale", "", "Output locale (e.g., en, fr, de, es, ja, zh)")
	generateCmd.Flags().StringVar(&generateLocaleFile, "locale-file", "", "Path to locale override JSON file")
	generateCmd.Flags().BoolVar(&generateAllReleases, "all-releases", false, "Include all releases (overrides default notable-only)")
	generateCmd.Flags().StringVar(&generateNotableCategories, "notable-categories", "", "Custom notable categories (comma-separated)")
	generateCmd.Flags().StringVar(&generateExcludeCategories, "exclude-categories", "", "Categories to omit entirely from output (comma-separated)")
	generateCmd.Flags().StringVar(&generateCollapseCategories, "collapse-categories", "", "Categories to render as a one-line summary instead of a full list (comma-separated)")
	generateCmd.Flags().StringVar(&generateExpandCategories, "expand-categories", "", "Categories to force-expand, overriding collapse defaults such as Dependencies (comma-separated)")
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

	// Select options using library function
	preset := "default"
	if generateMinimal {
		preset = "minimal"
	} else if generateFull {
		preset = "full"
	}

	opts, err := renderer.OptionsFromConfig(renderer.Config{
		Preset:             preset,
		MaxTier:            generateMaxTier,
		Locale:             generateLocale,
		LocaleOverrides:    generateLocaleFile,
		AllReleases:        generateAllReleases,
		NotableCategories:  splitCSV(generateNotableCategories),
		ExcludeCategories:  splitCSV(generateExcludeCategories),
		CollapseCategories: splitCSV(generateCollapseCategories),
		ExpandCategories:   splitCSV(generateExpandCategories),
	})
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
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

// splitCSV splits a comma-separated flag value into a trimmed, non-empty slice.
// It returns nil for an empty input so callers can treat "unset" distinctly.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
