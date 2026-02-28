package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/aggregate"
)

var (
	metricsGranularity string
	metricsSince       string
	metricsUntil       string
	metricsOutput      string
	metricsDashboard   bool
	metricsTemplate    string
	metricsTitle       string
)

var metricsCmd = &cobra.Command{
	Use:   "metrics <portfolio.json>",
	Short: "Generate metrics and dashboard data from a portfolio",
	Long: `Generate metrics and dashboard data from an aggregated portfolio.

Metrics include:
  - Total releases and entries
  - Breakdown by category and rolled-up groups
  - Per-project statistics
  - Time series data (for charts)
  - Daily activity data (for heatmaps)

Granularity options:
  - day:   Daily data points
  - week:  Weekly aggregation (ISO week)
  - month: Monthly aggregation

Examples:
  # Generate metrics with day granularity
  schangelog portfolio metrics portfolio.json --granularity day -o metrics.json

  # Filter by date range
  schangelog portfolio metrics portfolio.json --since 2024-01-01 --until 2024-12-31

  # Export dashboard-ready data
  schangelog portfolio metrics portfolio.json --dashboard -o dashboard-data.json

  # Generate complete dashforge dashboard
  schangelog portfolio metrics portfolio.json --dashboard --template velocity -o dashboard.json`,
	Args: cobra.ExactArgs(1),
	RunE: runMetrics,
}

func init() {
	metricsCmd.Flags().StringVar(&metricsGranularity, "granularity", "day", "Time granularity: day, week, month")
	metricsCmd.Flags().StringVar(&metricsSince, "since", "", "Start date (YYYY-MM-DD)")
	metricsCmd.Flags().StringVar(&metricsUntil, "until", "", "End date (YYYY-MM-DD)")
	metricsCmd.Flags().StringVarP(&metricsOutput, "output", "o", "", "Output file (default: stdout)")
	metricsCmd.Flags().BoolVar(&metricsDashboard, "dashboard", false, "Generate dashboard-ready data")
	metricsCmd.Flags().StringVar(&metricsTemplate, "template", "velocity", "Dashboard template: velocity, summary")
	metricsCmd.Flags().StringVar(&metricsTitle, "title", "", "Dashboard title (default: portfolio name)")
	portfolioCmd.AddCommand(metricsCmd)
}

func runMetrics(cmd *cobra.Command, args []string) error {
	portfolioPath := args[0]

	// Load portfolio
	portfolio, err := aggregate.LoadPortfolioFile(portfolioPath)
	if err != nil {
		return fmt.Errorf("loading portfolio: %w", err)
	}

	// Build options
	opts := aggregate.MetricsOptions{
		Granularity:    metricsGranularity,
		IncludeRollups: true,
	}

	// Parse date filters
	if metricsSince != "" {
		t, err := time.Parse("2006-01-02", metricsSince)
		if err != nil {
			return fmt.Errorf("invalid --since date: %w", err)
		}
		opts.Since = t
	}

	if metricsUntil != "" {
		t, err := time.Parse("2006-01-02", metricsUntil)
		if err != nil {
			return fmt.Errorf("invalid --until date: %w", err)
		}
		opts.Until = t
	}

	// Calculate metrics
	metrics, err := aggregate.CalculateMetrics(portfolio, opts)
	if err != nil {
		return fmt.Errorf("calculating metrics: %w", err)
	}

	var output []byte

	if metricsDashboard {
		// Export dashboard data
		export, err := aggregate.ExportDashboard(metrics)
		if err != nil {
			return fmt.Errorf("exporting dashboard: %w", err)
		}

		// If template specified, generate full dashboard
		if metricsTemplate != "" {
			dashOpts := aggregate.DefaultDashboardOptions()
			dashOpts.Template = metricsTemplate
			if metricsTitle != "" {
				dashOpts.Title = metricsTitle
			} else {
				dashOpts.Title = portfolio.Name + " - Development Velocity"
			}

			output, err = aggregate.GenerateDashboardJSON(export, dashOpts)
			if err != nil {
				return fmt.Errorf("generating dashboard: %w", err)
			}
		} else {
			output, err = export.JSON()
			if err != nil {
				return fmt.Errorf("marshaling export: %w", err)
			}
		}
	} else {
		// Output raw metrics
		output, err = json.MarshalIndent(metrics, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling metrics: %w", err)
		}
	}

	// Write output
	if metricsOutput == "" {
		fmt.Println(string(output))
	} else {
		if err := os.WriteFile(metricsOutput, output, 0600); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote output to %s\n", metricsOutput)
	}

	// Print summary to stderr
	fmt.Fprintf(os.Stderr, "\nMetrics summary:\n")
	fmt.Fprintf(os.Stderr, "  Releases: %d\n", metrics.TotalReleases)
	fmt.Fprintf(os.Stderr, "  Entries:  %d\n", metrics.TotalEntries)
	fmt.Fprintf(os.Stderr, "  Projects: %d\n", len(metrics.ByProject))

	if len(metrics.ByRollup) > 0 {
		fmt.Fprintf(os.Stderr, "  By rollup:\n")
		for name, count := range metrics.ByRollup {
			fmt.Fprintf(os.Stderr, "    %s: %d\n", name, count)
		}
	}

	return nil
}
