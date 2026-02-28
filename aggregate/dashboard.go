package aggregate

import (
	"encoding/json"
	"fmt"
)

// DashboardExport contains data formatted for dashboard consumption.
type DashboardExport struct {
	// Summary metrics widget data
	Summary SummaryData `json:"summary"`

	// Heatmap data: [[date, count], ...]
	HeatmapData [][]any `json:"heatmapData"`

	// Time series for charts
	ReleaseTrend  []map[string]any `json:"releaseTrend"`
	CategoryTrend []map[string]any `json:"categoryTrend"`

	// Project breakdown table
	ProjectTable []map[string]any `json:"projectTable"`
}

// SummaryData contains high-level summary metrics.
type SummaryData struct {
	TotalReleases int            `json:"totalReleases"`
	TotalEntries  int            `json:"totalEntries"`
	ProjectCount  int            `json:"projectCount"`
	DateRange     DateRange      `json:"dateRange"`
	ByRollup      map[string]int `json:"byRollup,omitempty"`
}

// ExportDashboard creates dashboard-compatible data from metrics.
func ExportDashboard(metrics *MetricsReport) (*DashboardExport, error) {
	export := &DashboardExport{
		Summary: SummaryData{
			TotalReleases: metrics.TotalReleases,
			TotalEntries:  metrics.TotalEntries,
			ProjectCount:  len(metrics.ByProject),
			DateRange:     metrics.DateRange,
			ByRollup:      metrics.ByRollup,
		},
	}

	// Heatmap data: [[date, count], ...]
	export.HeatmapData = make([][]any, len(metrics.DailyActivity))
	for i, dc := range metrics.DailyActivity {
		export.HeatmapData[i] = []any{dc.Date, dc.Count}
	}

	// Release trend: [{date, releases, entries}, ...]
	export.ReleaseTrend = make([]map[string]any, len(metrics.TimeSeries))
	for i, tp := range metrics.TimeSeries {
		export.ReleaseTrend[i] = map[string]any{
			"date":     tp.Date,
			"releases": tp.Releases,
			"entries":  tp.Entries,
		}
	}

	// Category trend: [{date, category, count}, ...] - flattened for chart libraries
	var categoryTrend []map[string]any
	for _, tp := range metrics.TimeSeries {
		for cat, count := range tp.ByCategory {
			categoryTrend = append(categoryTrend, map[string]any{
				"date":     tp.Date,
				"category": cat,
				"count":    count,
			})
		}
	}
	export.CategoryTrend = categoryTrend

	// Project table: [{path, name, releases, entries}, ...]
	export.ProjectTable = make([]map[string]any, 0, len(metrics.ByProject))
	for path, pm := range metrics.ByProject {
		export.ProjectTable = append(export.ProjectTable, map[string]any{
			"path":     path,
			"releases": pm.Releases,
			"entries":  pm.Entries,
		})
	}

	return export, nil
}

// DashboardOptions configures dashboard generation.
type DashboardOptions struct {
	Template    string // velocity, summary, heatmap
	Title       string
	Theme       string // light, dark, auto
	DateRange   string // last-12-months, last-6-months, ytd, all
	ShowFilters bool
}

// DefaultDashboardOptions returns default dashboard options.
func DefaultDashboardOptions() DashboardOptions {
	return DashboardOptions{
		Template:    "velocity",
		Title:       "Development Velocity Dashboard",
		Theme:       "auto",
		DateRange:   "last-12-months",
		ShowFilters: true,
	}
}

// GenerateDashboardJSON creates a complete dashforge dashboard definition.
func GenerateDashboardJSON(export *DashboardExport, opts DashboardOptions) ([]byte, error) {
	if opts.Title == "" {
		opts.Title = "Development Velocity Dashboard"
	}

	dashboard := map[string]any{
		"title":  opts.Title,
		"layout": map[string]any{"type": "grid", "columns": 12},
		"theme":  map[string]any{"mode": opts.Theme},
		"dataSources": []map[string]any{
			{
				"id":   "metrics",
				"type": "inline",
				"data": export,
			},
		},
		"widgets": generateWidgets(export, opts),
	}

	if opts.ShowFilters {
		dashboard["variables"] = []map[string]any{
			{
				"id":      "dateRange",
				"name":    "Date Range",
				"type":    "daterange",
				"default": opts.DateRange,
			},
		}
	}

	return json.MarshalIndent(dashboard, "", "  ")
}

// generateWidgets creates the widget configuration for the dashboard.
func generateWidgets(export *DashboardExport, _ DashboardOptions) []map[string]any {
	yPos := 0 // Track vertical position

	// Core metric widgets (always shown)
	widgets := []map[string]any{
		{
			"id":           "total-releases",
			"type":         "metric",
			"title":        "Total Releases",
			"position":     map[string]int{"x": 0, "y": yPos, "w": 3, "h": 2},
			"dataSourceId": "metrics",
			"config":       map[string]any{"field": "summary.totalReleases"},
		},
		{
			"id":           "total-entries",
			"type":         "metric",
			"title":        "Total Changes",
			"position":     map[string]int{"x": 3, "y": yPos, "w": 3, "h": 2},
			"dataSourceId": "metrics",
			"config":       map[string]any{"field": "summary.totalEntries"},
		},
		{
			"id":           "project-count",
			"type":         "metric",
			"title":        "Projects",
			"position":     map[string]int{"x": 6, "y": yPos, "w": 3, "h": 2},
			"dataSourceId": "metrics",
			"config":       map[string]any{"field": "summary.projectCount"},
		},
	}
	yPos += 2

	// Charts (always shown)
	widgets = append(widgets, map[string]any{
		"id":           "releases-by-month",
		"type":         "chart",
		"title":        "Releases Over Time",
		"position":     map[string]int{"x": 0, "y": yPos, "w": 6, "h": 4},
		"dataSourceId": "metrics",
		"config": map[string]any{
			"dataPath": "releaseTrend",
			"mark":     map[string]any{"type": "bar"},
			"encodings": map[string]any{
				"x": map[string]any{"field": "date", "type": "temporal"},
				"y": map[string]any{"field": "releases", "type": "quantitative"},
			},
		},
	})
	widgets = append(widgets, map[string]any{
		"id":           "category-breakdown",
		"type":         "chart",
		"title":        "Changes by Rollup",
		"position":     map[string]int{"x": 6, "y": yPos, "w": 6, "h": 4},
		"dataSourceId": "metrics",
		"config": map[string]any{
			"dataPath": "summary.byRollup",
			"mark":     map[string]any{"type": "pie"},
			"encodings": map[string]any{
				"theta": map[string]any{"field": "value", "type": "quantitative"},
				"color": map[string]any{"field": "key", "type": "nominal"},
			},
		},
	})
	yPos += 4

	// Activity heatmap (only if data exists)
	if len(export.HeatmapData) > 0 {
		widgets = append(widgets, map[string]any{
			"id":           "activity-heatmap",
			"type":         "chart",
			"title":        "Activity Heatmap",
			"position":     map[string]int{"x": 0, "y": yPos, "w": 12, "h": 4},
			"dataSourceId": "metrics",
			"config": map[string]any{
				"dataPath": "heatmapData",
				"mark":     map[string]any{"type": "heatmap"},
				"encodings": map[string]any{
					"x":    map[string]any{"field": "0", "type": "temporal"},
					"heat": map[string]any{"field": "1", "type": "quantitative"},
				},
				"style": map[string]any{
					"calendar": map[string]any{
						"range":    "auto",
						"cellSize": []any{"auto", 13},
					},
					"colors": []string{"#ebedf0", "#9be9a8", "#40c463", "#30a14e", "#216e39"},
				},
			},
		})
		yPos += 4
	}

	// Project table (only if projects exist)
	if len(export.ProjectTable) > 0 {
		widgets = append(widgets, map[string]any{
			"id":           "project-table",
			"type":         "table",
			"title":        "Projects",
			"position":     map[string]int{"x": 0, "y": yPos, "w": 12, "h": 6},
			"dataSourceId": "metrics",
			"config": map[string]any{
				"dataPath": "projectTable",
				"columns": []map[string]any{
					{"field": "path", "title": "Project"},
					{"field": "releases", "title": "Releases"},
					{"field": "entries", "title": "Changes"},
				},
				"sortable": true,
			},
		})
	}

	return widgets
}

// JSON returns the export as formatted JSON bytes.
func (e *DashboardExport) JSON() ([]byte, error) {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling dashboard export JSON: %w", err)
	}
	return data, nil
}
