package aggregate

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
)

//go:embed rollup_rules.json
var rollupRulesJSON []byte

// RollupRules defines how raw categories map to rolled-up groups.
type RollupRules struct {
	Version string              `json:"version"`
	Rollups map[string][]string `json:"rollups"` // RollupName -> []CategoryNames
}

// LoadRollupRules loads the embedded rollup rules.
func LoadRollupRules() (*RollupRules, error) {
	var rules RollupRules
	if err := json.Unmarshal(rollupRulesJSON, &rules); err != nil {
		return nil, fmt.Errorf("parsing embedded rollup rules: %w", err)
	}
	return &rules, nil
}

// ParseRollupRules parses rollup rules from JSON bytes.
func ParseRollupRules(data []byte) (*RollupRules, error) {
	var rules RollupRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parsing rollup rules JSON: %w", err)
	}
	return &rules, nil
}

// Apply rolls up raw category counts to grouped counts.
// Any categories not in a rollup group are placed in "Other".
func (r *RollupRules) Apply(raw map[string]int) map[string]int {
	result := make(map[string]int)
	used := make(map[string]bool)

	// Sum counts for each rollup group
	for rollupName, categories := range r.Rollups {
		for _, cat := range categories {
			if count, ok := raw[cat]; ok {
				result[rollupName] += count
				used[cat] = true
			}
		}
	}

	// Put uncategorized items in "Other"
	for cat, count := range raw {
		if !used[cat] {
			result["Other"] += count
		}
	}

	return result
}

// Categories returns all category names that belong to a rollup group.
func (r *RollupRules) Categories(rollupName string) []string {
	return r.Rollups[rollupName]
}

// RollupNames returns all rollup group names.
func (r *RollupRules) RollupNames() []string {
	names := make([]string, 0, len(r.Rollups))
	for name := range r.Rollups {
		names = append(names, name)
	}
	return names
}

// FindRollup returns the rollup group name for a category.
// Returns empty string if not found.
func (r *RollupRules) FindRollup(category string) string {
	for rollupName, categories := range r.Rollups {
		if slices.Contains(categories, category) {
			return rollupName
		}
	}
	return ""
}

// DefaultRollupRules returns the embedded default rollup rules.
// Panics if the embedded rules are invalid (should never happen).
func DefaultRollupRules() *RollupRules {
	rules, err := LoadRollupRules()
	if err != nil {
		panic(fmt.Sprintf("invalid embedded rollup rules: %v", err))
	}
	return rules
}
