package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "sclog",
	Short: "Structured Changelog CLI",
	Long: `sclog is a CLI tool for working with Structured Changelogs.

It uses a JSON Intermediate Representation (IR) as the canonical source
of truth and generates deterministic Markdown output following the
Keep a Changelog format.

Examples:
  sclog validate CHANGELOG.json
  sclog generate CHANGELOG.json -o CHANGELOG.md
  sclog version`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sclog %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
