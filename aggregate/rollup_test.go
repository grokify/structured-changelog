package aggregate

import (
	"testing"
)

func TestLoadRollupRules(t *testing.T) {
	rules, err := LoadRollupRules()
	if err != nil {
		t.Fatalf("LoadRollupRules() error: %v", err)
	}

	if rules.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", rules.Version)
	}

	// Check expected rollup groups exist
	expectedGroups := []string{"Features", "Fixes", "Improvements", "Breaking", "Maintenance", "Community"}
	for _, group := range expectedGroups {
		if _, ok := rules.Rollups[group]; !ok {
			t.Errorf("expected rollup group %q to exist", group)
		}
	}
}

func TestRollupRulesApply(t *testing.T) {
	rules, _ := LoadRollupRules()

	raw := map[string]int{
		"Added":           5,
		"Fixed":           3,
		"Security":        2,
		"Dependencies":    10,
		"Documentation":   4,
		"UnknownCategory": 1, // Should go to "Other"
	}

	result := rules.Apply(raw)

	tests := []struct {
		group    string
		expected int
	}{
		{"Features", 5},     // Added
		{"Fixes", 5},        // Fixed + Security
		{"Maintenance", 14}, // Dependencies + Documentation
		{"Other", 1},        // UnknownCategory
	}

	for _, tt := range tests {
		if result[tt.group] != tt.expected {
			t.Errorf("expected %s=%d, got %d", tt.group, tt.expected, result[tt.group])
		}
	}
}

func TestRollupRulesFindRollup(t *testing.T) {
	rules, _ := LoadRollupRules()

	tests := []struct {
		category string
		expected string
	}{
		{"Added", "Features"},
		{"Fixed", "Fixes"},
		{"Security", "Fixes"},
		{"Dependencies", "Maintenance"},
		{"Breaking", "Breaking"},
		{"Contributors", "Community"},
		{"Unknown", ""},
	}

	for _, tt := range tests {
		result := rules.FindRollup(tt.category)
		if result != tt.expected {
			t.Errorf("FindRollup(%q) = %q, expected %q", tt.category, result, tt.expected)
		}
	}
}

func TestRollupRulesCategories(t *testing.T) {
	rules, _ := LoadRollupRules()

	features := rules.Categories("Features")
	if len(features) != 2 {
		t.Errorf("expected 2 categories in Features, got %d", len(features))
	}

	// Check expected categories
	found := map[string]bool{}
	for _, cat := range features {
		found[cat] = true
	}

	if !found["Added"] || !found["Highlights"] {
		t.Errorf("expected Features to contain Added and Highlights, got %v", features)
	}
}

func TestRollupRulesRollupNames(t *testing.T) {
	rules, _ := LoadRollupRules()

	names := rules.RollupNames()
	if len(names) != 6 {
		t.Errorf("expected 6 rollup names, got %d", len(names))
	}
}

func TestDefaultRollupRules(t *testing.T) {
	rules := DefaultRollupRules()

	if rules == nil {
		t.Fatal("DefaultRollupRules() returned nil")
	}

	if len(rules.Rollups) == 0 {
		t.Error("DefaultRollupRules() returned empty rollups")
	}
}

func TestParseRollupRules(t *testing.T) {
	validJSON := []byte(`{
		"version": "1.0",
		"rollups": {
			"Features": ["Added", "Highlights"],
			"Fixes": ["Fixed", "Security"]
		}
	}`)

	rules, err := ParseRollupRules(validJSON)
	if err != nil {
		t.Fatalf("ParseRollupRules() error: %v", err)
	}

	if rules.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", rules.Version)
	}

	if len(rules.Rollups) != 2 {
		t.Errorf("expected 2 rollup groups, got %d", len(rules.Rollups))
	}

	if len(rules.Rollups["Features"]) != 2 {
		t.Errorf("expected 2 categories in Features, got %d", len(rules.Rollups["Features"]))
	}
}

func TestParseRollupRulesInvalid(t *testing.T) {
	invalidJSON := []byte(`{invalid json`)

	_, err := ParseRollupRules(invalidJSON)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestRollupRulesApplyEmpty(t *testing.T) {
	rules, _ := LoadRollupRules()

	// Empty input
	result := rules.Apply(map[string]int{})

	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %d entries", len(result))
	}
}

func TestRollupRulesApplyAllUnknown(t *testing.T) {
	rules, _ := LoadRollupRules()

	raw := map[string]int{
		"UnknownCategory1": 5,
		"UnknownCategory2": 3,
	}

	result := rules.Apply(raw)

	if result["Other"] != 8 {
		t.Errorf("expected Other=8, got %d", result["Other"])
	}

	// Should only have "Other" group
	if len(result) != 1 {
		t.Errorf("expected 1 result group, got %d", len(result))
	}
}
