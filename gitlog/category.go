package gitlog

import (
	"strings"
)

// CategorySuggestion represents a suggested changelog category with confidence.
type CategorySuggestion struct {
	Category   string  `json:"category"`
	Tier       string  `json:"tier"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning"`
}

// categoryMapping maps conventional commit types to changelog categories.
var categoryMapping = map[string]CategorySuggestion{
	"feat": {
		Category:   "Added",
		Tier:       "core",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'feat' indicates new functionality",
	},
	"fix": {
		Category:   "Fixed",
		Tier:       "core",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'fix' indicates bug fixes",
	},
	"docs": {
		Category:   "Documentation",
		Tier:       "extended",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'docs' indicates documentation changes",
	},
	"style": {
		Category:   "Internal",
		Tier:       "optional",
		Confidence: 0.90,
		Reasoning:  "Conventional commit type 'style' indicates formatting with no logic change",
	},
	"refactor": {
		Category:   "Changed",
		Tier:       "core",
		Confidence: 0.85,
		Reasoning:  "Conventional commit type 'refactor' indicates code restructuring",
	},
	"perf": {
		Category:   "Performance",
		Tier:       "standard",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'perf' indicates performance improvements",
	},
	"test": {
		Category:   "Tests",
		Tier:       "extended",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'test' indicates test additions or changes",
	},
	"build": {
		Category:   "Build",
		Tier:       "extended",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'build' indicates build system changes",
	},
	"ci": {
		Category:   "Infrastructure",
		Tier:       "optional",
		Confidence: 0.90,
		Reasoning:  "Conventional commit type 'ci' indicates CI/CD changes",
	},
	"chore": {
		Category:   "Internal",
		Tier:       "optional",
		Confidence: 0.85,
		Reasoning:  "Conventional commit type 'chore' indicates maintenance tasks",
	},
	"revert": {
		Category:   "Fixed",
		Tier:       "core",
		Confidence: 0.80,
		Reasoning:  "Reverting a commit typically indicates fixing a regression",
	},
	"security": {
		Category:   "Security",
		Tier:       "core",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'security' indicates security fixes",
	},
	"deps": {
		Category:   "Dependencies",
		Tier:       "standard",
		Confidence: 0.95,
		Reasoning:  "Conventional commit type 'deps' indicates dependency updates",
	},
}

// SuggestCategory suggests a changelog category for a commit based on its type.
func SuggestCategory(commitType string) *CategorySuggestion {
	t := strings.ToLower(commitType)
	if suggestion, ok := categoryMapping[t]; ok {
		return &suggestion
	}
	return nil
}

// SuggestCategoryFromMessage suggests a category by parsing the commit message.
func SuggestCategoryFromMessage(message string) *CategorySuggestion {
	cc := ParseConventionalCommit(message)
	if cc == nil {
		return inferCategoryFromMessage(message)
	}

	// Check for breaking change markers first
	if cc.Breaking {
		return &CategorySuggestion{
			Category:   "Breaking",
			Tier:       "standard",
			Confidence: 0.95,
			Reasoning:  "Commit marked with '!' indicates breaking change",
		}
	}

	// Check if body has BREAKING CHANGE
	lines := strings.SplitN(message, "\n", 2)
	if len(lines) > 1 && HasBreakingChangeMarker(lines[1]) {
		return &CategorySuggestion{
			Category:   "Breaking",
			Tier:       "standard",
			Confidence: 0.95,
			Reasoning:  "Commit body contains BREAKING CHANGE marker",
		}
	}

	return SuggestCategory(cc.Type)
}

// inferCategoryFromMessage attempts to infer category from non-conventional commits.
func inferCategoryFromMessage(message string) *CategorySuggestion {
	lower := strings.ToLower(message)

	// Check for common patterns
	// Note: Order matters - more specific patterns (like security) should come before generic ones (like fix)
	patterns := []struct {
		keywords   []string
		suggestion CategorySuggestion
	}{
		{
			keywords: []string{"security", "cve", "vulnerability", "exploit"},
			suggestion: CategorySuggestion{
				Category:   "Security",
				Tier:       "core",
				Confidence: 0.70,
				Reasoning:  "Message contains security-related keywords",
			},
		},
		{
			keywords: []string{"add ", "adds ", "added ", "adding ", "new ", "introduce ", "implement "},
			suggestion: CategorySuggestion{
				Category:   "Added",
				Tier:       "core",
				Confidence: 0.60,
				Reasoning:  "Message suggests new functionality",
			},
		},
		{
			keywords: []string{"fix ", "fixes ", "fixed ", "fixing ", "bug ", "resolve ", "repair "},
			suggestion: CategorySuggestion{
				Category:   "Fixed",
				Tier:       "core",
				Confidence: 0.60,
				Reasoning:  "Message suggests bug fix",
			},
		},
		{
			keywords: []string{"remove ", "removes ", "removed ", "delete ", "drop "},
			suggestion: CategorySuggestion{
				Category:   "Removed",
				Tier:       "core",
				Confidence: 0.60,
				Reasoning:  "Message suggests removal",
			},
		},
		{
			keywords: []string{"deprecate ", "deprecates ", "deprecated "},
			suggestion: CategorySuggestion{
				Category:   "Deprecated",
				Tier:       "core",
				Confidence: 0.70,
				Reasoning:  "Message indicates deprecation",
			},
		},
		{
			keywords: []string{"update readme", "update doc", "documentation"},
			suggestion: CategorySuggestion{
				Category:   "Documentation",
				Tier:       "extended",
				Confidence: 0.60,
				Reasoning:  "Message suggests documentation changes",
			},
		},
		{
			keywords: []string{"upgrade ", "bump ", "update depend", "update go.mod"},
			suggestion: CategorySuggestion{
				Category:   "Dependencies",
				Tier:       "standard",
				Confidence: 0.65,
				Reasoning:  "Message suggests dependency updates",
			},
		},
		{
			keywords: []string{"performance", "optimize", "speed up", "faster"},
			suggestion: CategorySuggestion{
				Category:   "Performance",
				Tier:       "standard",
				Confidence: 0.60,
				Reasoning:  "Message suggests performance improvement",
			},
		},
	}

	for _, p := range patterns {
		for _, kw := range p.keywords {
			if strings.Contains(lower, kw) {
				return &p.suggestion
			}
		}
	}

	// Default to Changed with low confidence
	return &CategorySuggestion{
		Category:   "Changed",
		Tier:       "core",
		Confidence: 0.30,
		Reasoning:  "Unable to determine specific category from message",
	}
}

// GetCategoryMapping returns the full category mapping for reference.
func GetCategoryMapping() map[string]CategorySuggestion {
	// Return a copy to prevent modification
	result := make(map[string]CategorySuggestion)
	for k, v := range categoryMapping {
		result[k] = v
	}
	return result
}
