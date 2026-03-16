package aggregate

import (
	"testing"
	"time"

	"github.com/grokify/structured-changelog/changelog"
)

func TestCalculateMetrics(t *testing.T) {
	portfolio := &Portfolio{
		Name: "Test Portfolio",
		Projects: []ProjectData{
			{
				Path: "repo1",
				Name: "Repo 1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{
							Version: "1.0.0",
							Date:    "2024-06-15",
							Added:   []changelog.Entry{{Description: "A"}, {Description: "B"}},
							Fixed:   []changelog.Entry{{Description: "C"}},
						},
						{
							Version: "0.9.0",
							Date:    "2024-03-10",
							Added:   []changelog.Entry{{Description: "D"}},
						},
					},
				},
			},
			{
				Path: "repo2",
				Name: "Repo 2",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{
							Version:  "2.0.0",
							Date:     "2024-06-20",
							Security: []changelog.Entry{{Description: "E"}},
						},
					},
				},
			},
		},
	}

	opts := MetricsOptions{
		Granularity:    GranularityMonth,
		IncludeRollups: true,
	}

	metrics, err := CalculateMetrics(portfolio, opts)
	if err != nil {
		t.Fatalf("CalculateMetrics() error: %v", err)
	}

	// Check totals
	if metrics.TotalReleases != 3 {
		t.Errorf("expected 3 releases, got %d", metrics.TotalReleases)
	}

	if metrics.TotalEntries != 5 {
		t.Errorf("expected 5 entries, got %d", metrics.TotalEntries)
	}

	// Check category counts
	if metrics.ByCategory["Added"] != 3 {
		t.Errorf("expected 3 Added, got %d", metrics.ByCategory["Added"])
	}

	if metrics.ByCategory["Fixed"] != 1 {
		t.Errorf("expected 1 Fixed, got %d", metrics.ByCategory["Fixed"])
	}

	if metrics.ByCategory["Security"] != 1 {
		t.Errorf("expected 1 Security, got %d", metrics.ByCategory["Security"])
	}

	// Check rollups
	if metrics.ByRollup["Features"] != 3 {
		t.Errorf("expected Features=3, got %d", metrics.ByRollup["Features"])
	}

	if metrics.ByRollup["Fixes"] != 2 { // Fixed + Security
		t.Errorf("expected Fixes=2, got %d", metrics.ByRollup["Fixes"])
	}

	// Check per-project metrics
	if len(metrics.ByProject) != 2 {
		t.Errorf("expected 2 projects, got %d", len(metrics.ByProject))
	}

	repo1 := metrics.ByProject["repo1"]
	if repo1.Releases != 2 {
		t.Errorf("expected repo1 releases=2, got %d", repo1.Releases)
	}

	if repo1.Entries != 4 {
		t.Errorf("expected repo1 entries=4, got %d", repo1.Entries)
	}
}

func TestCalculateMetricsWithDateFilter(t *testing.T) {
	portfolio := &Portfolio{
		Name: "Test",
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "3.0.0", Date: "2024-08-01", Added: []changelog.Entry{{Description: "A"}}},
						{Version: "2.0.0", Date: "2024-05-01", Added: []changelog.Entry{{Description: "B"}}},
						{Version: "1.0.0", Date: "2024-02-01", Added: []changelog.Entry{{Description: "C"}}},
					},
				},
			},
		},
	}

	opts := MetricsOptions{
		Granularity:    GranularityMonth,
		Since:          time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		Until:          time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		IncludeRollups: false,
	}

	metrics, err := CalculateMetrics(portfolio, opts)
	if err != nil {
		t.Fatalf("CalculateMetrics() error: %v", err)
	}

	// Should only include the 2024-05-01 release
	if metrics.TotalReleases != 1 {
		t.Errorf("expected 1 release in date range, got %d", metrics.TotalReleases)
	}

	if metrics.TotalEntries != 1 {
		t.Errorf("expected 1 entry in date range, got %d", metrics.TotalEntries)
	}
}

func TestNormalizeDate(t *testing.T) {
	tests := []struct {
		date        string
		granularity string
		expected    string
	}{
		{"2024-06-15", GranularityDay, "2024-06-15"},
		{"2024-06-15", GranularityMonth, "2024-06"},
		{"2024-06-15", GranularityWeek, "2024-W24"}, // June 15, 2024 is week 24
		{"2024-01-01", GranularityWeek, "2024-W01"}, // January 1, 2024 is week 1
		// Edge cases
		{"invalid-date", GranularityDay, "invalid-date"},    // Invalid date returns as-is
		{"2024-06-15", "invalid", "2024-06-15"},             // Invalid granularity returns as-is
		{"", GranularityDay, ""},                            // Empty date
	}

	for _, tt := range tests {
		t.Run(tt.date+"_"+tt.granularity, func(t *testing.T) {
			result := normalizeDate(tt.date, tt.granularity)
			if result != tt.expected {
				t.Errorf("normalizeDate(%q, %q) = %q, expected %q", tt.date, tt.granularity, result, tt.expected)
			}
		})
	}
}

func TestDefaultMetricsOptions(t *testing.T) {
	opts := DefaultMetricsOptions()

	if opts.Granularity != GranularityDay {
		t.Errorf("expected Granularity=%q, got %q", GranularityDay, opts.Granularity)
	}

	if !opts.IncludeRollups {
		t.Error("expected IncludeRollups=true")
	}

	// Since should be approximately 1 year ago
	expectedSince := time.Now().AddDate(-1, 0, 0)
	diff := opts.Since.Sub(expectedSince)
	if diff < -time.Hour || diff > time.Hour {
		t.Errorf("expected Since to be approximately 1 year ago, got %v", opts.Since)
	}

	// Until should be approximately now
	expectedUntil := time.Now()
	diff = opts.Until.Sub(expectedUntil)
	if diff < -time.Hour || diff > time.Hour {
		t.Errorf("expected Until to be approximately now, got %v", opts.Until)
	}
}

func TestMetricsTimeSeries(t *testing.T) {
	portfolio := &Portfolio{
		Name: "Test",
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "2.0.0", Date: "2024-06-15", Added: []changelog.Entry{{Description: "A"}}},
						{Version: "1.0.0", Date: "2024-06-01", Added: []changelog.Entry{{Description: "B"}}},
					},
				},
			},
		},
	}

	opts := MetricsOptions{
		Granularity:    GranularityMonth,
		IncludeRollups: true,
	}

	metrics, err := CalculateMetrics(portfolio, opts)
	if err != nil {
		t.Fatalf("CalculateMetrics() error: %v", err)
	}

	// Both releases should be in 2024-06 (same month)
	if len(metrics.TimeSeries) != 1 {
		t.Fatalf("expected 1 time point, got %d", len(metrics.TimeSeries))
	}

	point := metrics.TimeSeries[0]
	if point.Date != "2024-06" {
		t.Errorf("expected date %q, got %q", "2024-06", point.Date)
	}

	if point.Releases != 2 {
		t.Errorf("expected 2 releases, got %d", point.Releases)
	}

	if point.Entries != 2 {
		t.Errorf("expected 2 entries, got %d", point.Entries)
	}
}

func TestMetricsDailyActivity(t *testing.T) {
	portfolio := &Portfolio{
		Name: "Test",
		Projects: []ProjectData{
			{
				Path: "repo1",
				Changelog: &changelog.Changelog{
					Releases: []changelog.Release{
						{Version: "2.0.0", Date: "2024-06-15"},
						{Version: "1.0.0", Date: "2024-06-15"}, // Same day
						{Version: "0.9.0", Date: "2024-06-01"},
					},
				},
			},
		},
	}

	// Use options without date filter to include all test data
	opts := MetricsOptions{
		Granularity:    GranularityDay,
		IncludeRollups: true,
	}

	metrics, err := CalculateMetrics(portfolio, opts)
	if err != nil {
		t.Fatalf("CalculateMetrics() error: %v", err)
	}

	if len(metrics.DailyActivity) != 2 {
		t.Fatalf("expected 2 daily entries, got %d", len(metrics.DailyActivity))
	}

	// Check that June 15 has count 2
	var foundJune15 bool
	for _, dc := range metrics.DailyActivity {
		if dc.Date == "2024-06-15" {
			foundJune15 = true
			if dc.Count != 2 {
				t.Errorf("expected count 2 for 2024-06-15, got %d", dc.Count)
			}
		}
	}

	if !foundJune15 {
		t.Error("expected to find 2024-06-15 in daily activity")
	}
}

func TestMetricsTopCategories(t *testing.T) {
	metrics := &MetricsReport{
		ByCategory: map[string]int{
			"Added":    10,
			"Fixed":    5,
			"Security": 3,
			"Changed":  1,
		},
	}

	top := metrics.TopCategories(2)

	if len(top) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(top))
	}

	if top[0].Category != "Added" || top[0].Count != 10 {
		t.Errorf("expected first to be Added=10, got %s=%d", top[0].Category, top[0].Count)
	}

	if top[1].Category != "Fixed" || top[1].Count != 5 {
		t.Errorf("expected second to be Fixed=5, got %s=%d", top[1].Category, top[1].Count)
	}
}

func TestMetricsTopProjects(t *testing.T) {
	metrics := &MetricsReport{
		ByProject: map[string]ProjectMetrics{
			"repo1": {Releases: 5, Entries: 20},
			"repo2": {Releases: 3, Entries: 50},
			"repo3": {Releases: 10, Entries: 15},
		},
	}

	top := metrics.TopProjects(2)

	if len(top) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(top))
	}

	// Sorted by entries descending
	if top[0].Path != "repo2" || top[0].Entries != 50 {
		t.Errorf("expected first to be repo2=50, got %s=%d", top[0].Path, top[0].Entries)
	}

	if top[1].Path != "repo1" || top[1].Entries != 20 {
		t.Errorf("expected second to be repo1=20, got %s=%d", top[1].Path, top[1].Entries)
	}
}
