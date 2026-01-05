package gitlog

import (
	"testing"
)

func TestSuggestCategory(t *testing.T) {
	tests := []struct {
		commitType string
		expected   string
	}{
		{"feat", "Added"},
		{"fix", "Fixed"},
		{"docs", "Documentation"},
		{"style", "Internal"},
		{"refactor", "Changed"},
		{"perf", "Performance"},
		{"test", "Tests"},
		{"build", "Build"},
		{"ci", "Infrastructure"},
		{"chore", "Internal"},
		{"revert", "Fixed"},
		{"security", "Security"},
		{"deps", "Dependencies"},
		// Test case insensitivity
		{"FEAT", "Added"},
		{"Fix", "Fixed"},
	}

	for _, tt := range tests {
		t.Run(tt.commitType, func(t *testing.T) {
			result := SuggestCategory(tt.commitType)
			if result == nil {
				t.Errorf("expected suggestion for %s, got nil", tt.commitType)
				return
			}
			if result.Category != tt.expected {
				t.Errorf("expected category %s, got %s", tt.expected, result.Category)
			}
		})
	}
}

func TestSuggestCategoryUnknownType(t *testing.T) {
	result := SuggestCategory("unknown")
	if result != nil {
		t.Errorf("expected nil for unknown type, got %+v", result)
	}
}

func TestSuggestCategoryFromMessage(t *testing.T) {
	tests := []struct {
		name             string
		message          string
		expectedCategory string
		minConfidence    float64
	}{
		{
			name:             "conventional feat",
			message:          "feat: add new feature",
			expectedCategory: "Added",
			minConfidence:    0.90,
		},
		{
			name:             "conventional fix",
			message:          "fix: resolve bug",
			expectedCategory: "Fixed",
			minConfidence:    0.90,
		},
		{
			name:             "breaking with bang",
			message:          "feat!: remove old API",
			expectedCategory: "Breaking",
			minConfidence:    0.90,
		},
		{
			name:             "breaking in body",
			message:          "feat: change API\n\nBREAKING CHANGE: old method removed",
			expectedCategory: "Breaking",
			minConfidence:    0.90,
		},
		{
			name:             "non-conventional add",
			message:          "Add new feature",
			expectedCategory: "Added",
			minConfidence:    0.50,
		},
		{
			name:             "non-conventional fix",
			message:          "Fix memory leak",
			expectedCategory: "Fixed",
			minConfidence:    0.50,
		},
		{
			name:             "non-conventional remove",
			message:          "Remove deprecated method",
			expectedCategory: "Removed",
			minConfidence:    0.50,
		},
		{
			name:             "non-conventional deprecate",
			message:          "Deprecate old API",
			expectedCategory: "Deprecated",
			minConfidence:    0.60,
		},
		{
			name:             "security keywords",
			message:          "Fix security vulnerability CVE-2024-1234",
			expectedCategory: "Security",
			minConfidence:    0.60,
		},
		{
			name:             "documentation update",
			message:          "Update README documentation",
			expectedCategory: "Documentation",
			minConfidence:    0.50,
		},
		{
			name:             "dependency update",
			message:          "Bump go version to 1.22",
			expectedCategory: "Dependencies",
			minConfidence:    0.60,
		},
		{
			name:             "performance improvement",
			message:          "Optimize query performance",
			expectedCategory: "Performance",
			minConfidence:    0.50,
		},
		{
			name:             "unknown falls back to Changed",
			message:          "Update something",
			expectedCategory: "Changed",
			minConfidence:    0.20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SuggestCategoryFromMessage(tt.message)
			if result == nil {
				t.Errorf("expected suggestion, got nil")
				return
			}
			if result.Category != tt.expectedCategory {
				t.Errorf("expected category %s, got %s", tt.expectedCategory, result.Category)
			}
			if result.Confidence < tt.minConfidence {
				t.Errorf("expected confidence >= %f, got %f", tt.minConfidence, result.Confidence)
			}
		})
	}
}

func TestSuggestCategoryTiers(t *testing.T) {
	tests := []struct {
		commitType   string
		expectedTier string
	}{
		{"feat", "core"},
		{"fix", "core"},
		{"security", "core"},
		{"perf", "standard"},
		{"deps", "standard"},
		{"docs", "extended"},
		{"test", "extended"},
		{"build", "extended"},
		{"ci", "optional"},
		{"chore", "optional"},
		{"style", "optional"},
	}

	for _, tt := range tests {
		t.Run(tt.commitType, func(t *testing.T) {
			result := SuggestCategory(tt.commitType)
			if result == nil {
				t.Errorf("expected suggestion for %s, got nil", tt.commitType)
				return
			}
			if result.Tier != tt.expectedTier {
				t.Errorf("expected tier %s, got %s", tt.expectedTier, result.Tier)
			}
		})
	}
}

func TestGetCategoryMapping(t *testing.T) {
	mapping := GetCategoryMapping()

	// Should return a copy
	if len(mapping) == 0 {
		t.Error("expected non-empty mapping")
	}

	// Check some expected entries
	if feat, ok := mapping["feat"]; !ok {
		t.Error("expected 'feat' in mapping")
	} else if feat.Category != "Added" {
		t.Errorf("expected feat to map to Added, got %s", feat.Category)
	}

	// Modifying the copy shouldn't affect the original
	delete(mapping, "feat")
	mapping2 := GetCategoryMapping()
	if _, ok := mapping2["feat"]; !ok {
		t.Error("modifying returned map affected original")
	}
}

func TestCategorySuggestionHasReasoning(t *testing.T) {
	// All suggestions should have reasoning
	for _, typ := range KnownConventionalTypes {
		result := SuggestCategory(typ)
		if result == nil {
			continue // Some types might not have mappings
		}
		if result.Reasoning == "" {
			t.Errorf("suggestion for %s has no reasoning", typ)
		}
	}
}
