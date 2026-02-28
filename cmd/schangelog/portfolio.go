package main

import (
	"github.com/spf13/cobra"
)

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "Multi-project changelog aggregation and metrics",
	Long: `Commands for working with changelogs across multiple projects.

These commands enable you to:
  - Discover projects with CHANGELOG.json in GitHub orgs/users
  - Aggregate changelogs into a unified portfolio
  - Generate metrics and dashboard data

Workflow:
  1. Discover projects:  schangelog portfolio discover --org myorg -o manifest.json
  2. Aggregate:          schangelog portfolio aggregate manifest.json -o portfolio.json
  3. Generate metrics:   schangelog portfolio metrics portfolio.json -o metrics.json

Examples:
  schangelog portfolio discover --org myorg --user myuser -o manifest.json
  schangelog portfolio aggregate manifest.json -o portfolio.json
  schangelog portfolio metrics portfolio.json --dashboard -o dashboard.json`,
}

func init() {
	rootCmd.AddCommand(portfolioCmd)
}
