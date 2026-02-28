package aggregate

import (
	"encoding/json"
	"testing"
)

func TestExportDashboard(t *testing.T) {
	metrics := &MetricsReport{
		Portfolio:     "Test Portfolio",
		TotalReleases: 10,
		TotalEntries:  50,
		DateRange:     DateRange{Start: "2024-01-01", End: "2024-12-31"},
		ByCategory: map[string]int{
			"Added": 30,
			"Fixed": 20,
		},
		ByRollup: map[string]int{
			"Features": 30,
			"Fixes":    20,
		},
		ByProject: map[string]ProjectMetrics{
			"repo1": {Releases: 5, Entries: 25},
			"repo2": {Releases: 5, Entries: 25},
		},
		TimeSeries: []TimePoint{
			{Date: "2024-01", Releases: 3, Entries: 15},
			{Date: "2024-02", Releases: 4, Entries: 20},
		},
		DailyActivity: []DailyCount{
			{Date: "2024-01-15", Count: 2},
			{Date: "2024-02-20", Count: 3},
		},
	}

	export, err := ExportDashboard(metrics)
	if err != nil {
		t.Fatalf("ExportDashboard() error: %v", err)
	}

	// Check summary
	if export.Summary.TotalReleases != 10 {
		t.Errorf("expected TotalReleases=10, got %d", export.Summary.TotalReleases)
	}

	if export.Summary.TotalEntries != 50 {
		t.Errorf("expected TotalEntries=50, got %d", export.Summary.TotalEntries)
	}

	if export.Summary.ProjectCount != 2 {
		t.Errorf("expected ProjectCount=2, got %d", export.Summary.ProjectCount)
	}

	// Check heatmap data
	if len(export.HeatmapData) != 2 {
		t.Errorf("expected 2 heatmap entries, got %d", len(export.HeatmapData))
	}

	if export.HeatmapData[0][0] != "2024-01-15" {
		t.Errorf("expected first heatmap date %q, got %v", "2024-01-15", export.HeatmapData[0][0])
	}

	// Check release trend
	if len(export.ReleaseTrend) != 2 {
		t.Errorf("expected 2 trend entries, got %d", len(export.ReleaseTrend))
	}

	// Check project table
	if len(export.ProjectTable) != 2 {
		t.Errorf("expected 2 project entries, got %d", len(export.ProjectTable))
	}
}

func TestGenerateDashboardJSON(t *testing.T) {
	export := &DashboardExport{
		Summary: SummaryData{
			TotalReleases: 10,
			TotalEntries:  50,
			ProjectCount:  5,
		},
		HeatmapData: [][]any{
			{"2024-01-15", 2},
		},
		ReleaseTrend: []map[string]any{
			{"date": "2024-01", "releases": 3},
		},
		ProjectTable: []map[string]any{
			{"path": "repo1", "releases": 5},
		},
	}

	opts := DashboardOptions{
		Title: "Test Dashboard",
		Theme: "dark",
	}

	data, err := GenerateDashboardJSON(export, opts)
	if err != nil {
		t.Fatalf("GenerateDashboardJSON() error: %v", err)
	}

	// Parse and check structure
	var dashboard map[string]any
	if err := json.Unmarshal(data, &dashboard); err != nil {
		t.Fatalf("failed to parse dashboard JSON: %v", err)
	}

	if dashboard["title"] != "Test Dashboard" {
		t.Errorf("expected title %q, got %v", "Test Dashboard", dashboard["title"])
	}

	// Check theme
	theme, ok := dashboard["theme"].(map[string]any)
	if !ok {
		t.Fatal("expected theme to be a map")
	}
	if theme["mode"] != "dark" {
		t.Errorf("expected theme mode %q, got %v", "dark", theme["mode"])
	}

	// Check widgets exist
	widgets, ok := dashboard["widgets"].([]any)
	if !ok {
		t.Fatal("expected widgets to be an array")
	}

	if len(widgets) < 4 { // At least metrics, charts, heatmap, table
		t.Errorf("expected at least 4 widgets, got %d", len(widgets))
	}

	// Check data sources
	dataSources, ok := dashboard["dataSources"].([]any)
	if !ok {
		t.Fatal("expected dataSources to be an array")
	}

	if len(dataSources) != 1 {
		t.Errorf("expected 1 data source, got %d", len(dataSources))
	}
}

func TestDashboardExportJSON(t *testing.T) {
	export := &DashboardExport{
		Summary: SummaryData{
			TotalReleases: 10,
		},
	}

	data, err := export.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	// Should be valid JSON
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Check summary
	summary, ok := parsed["summary"].(map[string]any)
	if !ok {
		t.Fatal("expected summary to be a map")
	}

	if summary["totalReleases"] != float64(10) {
		t.Errorf("expected totalReleases=10, got %v", summary["totalReleases"])
	}
}

func TestDefaultDashboardOptions(t *testing.T) {
	opts := DefaultDashboardOptions()

	if opts.Template != "velocity" {
		t.Errorf("expected template %q, got %q", "velocity", opts.Template)
	}

	if opts.Theme != "auto" {
		t.Errorf("expected theme %q, got %q", "auto", opts.Theme)
	}

	if !opts.ShowFilters {
		t.Error("expected ShowFilters to be true")
	}
}

func TestGenerateWidgetsConditional(t *testing.T) {
	// Test with empty heatmap data
	export := &DashboardExport{
		Summary:      SummaryData{ProjectCount: 0},
		HeatmapData:  [][]any{},
		ProjectTable: []map[string]any{},
	}

	widgets := generateWidgets(export, DefaultDashboardOptions())

	// Should have fewer widgets without heatmap and project table
	hasHeatmap := false
	hasProjectTable := false

	for _, w := range widgets {
		wmap := w
		if wmap["id"] == "activity-heatmap" {
			hasHeatmap = true
		}
		if wmap["id"] == "project-table" {
			hasProjectTable = true
		}
	}

	if hasHeatmap {
		t.Error("expected no heatmap widget with empty data")
	}

	if hasProjectTable {
		t.Error("expected no project table widget with empty data")
	}

	// Test with data
	exportWithData := &DashboardExport{
		Summary:      SummaryData{ProjectCount: 1},
		HeatmapData:  [][]any{{"2024-01-01", 1}},
		ProjectTable: []map[string]any{{"path": "repo"}},
	}

	widgetsWithData := generateWidgets(exportWithData, DefaultDashboardOptions())

	hasHeatmap = false
	hasProjectTable = false

	for _, w := range widgetsWithData {
		wmap := w
		if wmap["id"] == "activity-heatmap" {
			hasHeatmap = true
		}
		if wmap["id"] == "project-table" {
			hasProjectTable = true
		}
	}

	if !hasHeatmap {
		t.Error("expected heatmap widget with data")
	}

	if !hasProjectTable {
		t.Error("expected project table widget with data")
	}
}
