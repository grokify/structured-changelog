package aggregate

import (
	"fmt"
	"sort"
	"time"
)

// MetricsReport contains calculated metrics for a portfolio.
type MetricsReport struct {
	Portfolio   string    `json:"portfolio"`
	DateRange   DateRange `json:"dateRange"`
	Granularity string    `json:"granularity"` // day, week, month

	// Summary totals
	TotalReleases int `json:"totalReleases"`
	TotalEntries  int `json:"totalEntries"`

	// Raw category counts (all categories)
	ByCategory map[string]int `json:"byCategory"`

	// Rolled-up counts
	ByRollup map[string]int `json:"byRollup,omitempty"`

	// Per-project breakdown
	ByProject map[string]ProjectMetrics `json:"byProject"`

	// Time series (for line/bar charts)
	TimeSeries []TimePoint `json:"timeSeries"`

	// Daily activity (for heatmap)
	DailyActivity []DailyCount `json:"dailyActivity"`
}

// ProjectMetrics contains metrics for a single project.
type ProjectMetrics struct {
	Releases   int            `json:"releases"`
	Entries    int            `json:"entries"`
	ByCategory map[string]int `json:"byCategory"`
}

// TimePoint represents metrics at a point in time.
type TimePoint struct {
	Date       string         `json:"date"` // YYYY-MM-DD, YYYY-MM, or YYYY-Www
	Releases   int            `json:"releases"`
	Entries    int            `json:"entries"`
	ByCategory map[string]int `json:"byCategory,omitempty"`
	ByRollup   map[string]int `json:"byRollup,omitempty"`
}

// DailyCount represents activity for a single day.
type DailyCount struct {
	Date  string `json:"date"`  // YYYY-MM-DD
	Count int    `json:"count"` // Total activity
}

// MetricsOptions configures metrics calculation.
type MetricsOptions struct {
	Granularity    string    // day, week, month
	Since          time.Time // Filter start
	Until          time.Time // Filter end
	IncludeRollups bool      // Include rolled-up counts
}

// GranularityDay specifies daily granularity.
const GranularityDay = "day"

// GranularityWeek specifies weekly granularity.
const GranularityWeek = "week"

// GranularityMonth specifies monthly granularity.
const GranularityMonth = "month"

// DefaultMetricsOptions returns default options (daily, last 12 months, with rollups).
func DefaultMetricsOptions() MetricsOptions {
	return MetricsOptions{
		Granularity:    GranularityDay,
		Since:          time.Now().AddDate(-1, 0, 0),
		Until:          time.Now(),
		IncludeRollups: true,
	}
}

// CalculateMetrics computes metrics for a portfolio.
func CalculateMetrics(portfolio *Portfolio, opts MetricsOptions) (*MetricsReport, error) {
	if opts.Granularity == "" {
		opts.Granularity = GranularityDay
	}

	report := &MetricsReport{
		Portfolio:   portfolio.Name,
		Granularity: opts.Granularity,
		ByCategory:  make(map[string]int),
		ByRollup:    make(map[string]int),
		ByProject:   make(map[string]ProjectMetrics),
	}

	// Date filtering
	sinceStr := ""
	untilStr := ""
	if !opts.Since.IsZero() {
		sinceStr = opts.Since.Format("2006-01-02")
		report.DateRange.Start = sinceStr
	}
	if !opts.Until.IsZero() {
		untilStr = opts.Until.Format("2006-01-02")
		report.DateRange.End = untilStr
	}

	// Collect daily data for heatmap
	dailyData := make(map[string]int) // date -> count

	// Process each project
	for _, pd := range portfolio.Projects {
		if pd.Changelog == nil {
			continue
		}

		pm := ProjectMetrics{
			ByCategory: make(map[string]int),
		}

		for _, release := range pd.Changelog.Releases {
			// Date filtering
			if release.Date != "" {
				if sinceStr != "" && release.Date < sinceStr {
					continue
				}
				if untilStr != "" && release.Date > untilStr {
					continue
				}
			}

			pm.Releases++
			report.TotalReleases++

			// Count entries by category
			cats := release.Categories()
			for _, cat := range cats {
				entries := release.GetEntries(cat.Name)
				count := len(entries)
				pm.Entries += count
				pm.ByCategory[cat.Name] += count
				report.ByCategory[cat.Name] += count
				report.TotalEntries += count
			}

			// Daily activity
			if release.Date != "" {
				dailyData[release.Date]++
			}
		}

		report.ByProject[pd.Path] = pm
	}

	// Calculate rollups
	if opts.IncludeRollups {
		rules := DefaultRollupRules()
		report.ByRollup = rules.Apply(report.ByCategory)
	}

	// Generate time series
	report.TimeSeries = generateTimeSeries(portfolio, opts)

	// Generate daily activity
	report.DailyActivity = generateDailyActivity(dailyData)

	return report, nil
}

// generateTimeSeries creates time series data based on granularity.
func generateTimeSeries(portfolio *Portfolio, opts MetricsOptions) []TimePoint {
	// Collect all release dates with their data
	dateData := make(map[string]*TimePoint) // normalized date -> data

	sinceStr := ""
	untilStr := ""
	if !opts.Since.IsZero() {
		sinceStr = opts.Since.Format("2006-01-02")
	}
	if !opts.Until.IsZero() {
		untilStr = opts.Until.Format("2006-01-02")
	}

	for _, pd := range portfolio.Projects {
		if pd.Changelog == nil {
			continue
		}

		for _, release := range pd.Changelog.Releases {
			if release.Date == "" {
				continue
			}

			// Date filtering
			if sinceStr != "" && release.Date < sinceStr {
				continue
			}
			if untilStr != "" && release.Date > untilStr {
				continue
			}

			normalizedDate := normalizeDate(release.Date, opts.Granularity)

			tp, exists := dateData[normalizedDate]
			if !exists {
				tp = &TimePoint{
					Date:       normalizedDate,
					ByCategory: make(map[string]int),
					ByRollup:   make(map[string]int),
				}
				dateData[normalizedDate] = tp
			}

			tp.Releases++

			cats := release.Categories()
			for _, cat := range cats {
				entries := release.GetEntries(cat.Name)
				count := len(entries)
				tp.Entries += count
				tp.ByCategory[cat.Name] += count
			}
		}
	}

	// Apply rollups to each time point
	if opts.IncludeRollups {
		rules := DefaultRollupRules()
		for _, tp := range dateData {
			tp.ByRollup = rules.Apply(tp.ByCategory)
		}
	}

	// Convert to slice and sort
	series := make([]TimePoint, 0, len(dateData))
	for _, tp := range dateData {
		series = append(series, *tp)
	}

	sort.Slice(series, func(i, j int) bool {
		return series[i].Date < series[j].Date
	})

	return series
}

// normalizeDate converts a date to the appropriate granularity format.
func normalizeDate(date string, granularity string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date // Return as-is if parsing fails
	}

	switch granularity {
	case GranularityDay:
		return date
	case GranularityWeek:
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case GranularityMonth:
		return t.Format("2006-01")
	default:
		return date
	}
}

// generateDailyActivity creates the daily activity list sorted by date.
func generateDailyActivity(dailyData map[string]int) []DailyCount {
	activity := make([]DailyCount, 0, len(dailyData))

	for date, count := range dailyData {
		activity = append(activity, DailyCount{
			Date:  date,
			Count: count,
		})
	}

	sort.Slice(activity, func(i, j int) bool {
		return activity[i].Date < activity[j].Date
	})

	return activity
}

// TopCategories returns the top N categories by count.
func (r *MetricsReport) TopCategories(n int) []CategoryCount {
	counts := make([]CategoryCount, 0, len(r.ByCategory))
	for cat, count := range r.ByCategory {
		counts = append(counts, CategoryCount{Category: cat, Count: count})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	if n > 0 && n < len(counts) {
		counts = counts[:n]
	}

	return counts
}

// CategoryCount pairs a category with its count.
type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// TopProjects returns the top N projects by entry count.
func (r *MetricsReport) TopProjects(n int) []ProjectCount {
	counts := make([]ProjectCount, 0, len(r.ByProject))
	for path, pm := range r.ByProject {
		counts = append(counts, ProjectCount{
			Path:     path,
			Releases: pm.Releases,
			Entries:  pm.Entries,
		})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Entries > counts[j].Entries
	})

	if n > 0 && n < len(counts) {
		counts = counts[:n]
	}

	return counts
}

// ProjectCount pairs a project with its counts.
type ProjectCount struct {
	Path     string `json:"path"`
	Releases int    `json:"releases"`
	Entries  int    `json:"entries"`
}
