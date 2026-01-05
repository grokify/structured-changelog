package gitlog

import (
	"testing"
)

func TestParseConventionalCommit(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected *ConventionalCommit
	}{
		{
			name:    "simple feat",
			message: "feat: add new feature",
			expected: &ConventionalCommit{
				Type:     "feat",
				Subject:  "add new feature",
				Breaking: false,
			},
		},
		{
			name:    "feat with scope",
			message: "feat(auth): add OAuth2 support",
			expected: &ConventionalCommit{
				Type:     "feat",
				Scope:    "auth",
				Subject:  "add OAuth2 support",
				Breaking: false,
			},
		},
		{
			name:    "breaking change with bang",
			message: "feat!: remove deprecated API",
			expected: &ConventionalCommit{
				Type:     "feat",
				Subject:  "remove deprecated API",
				Breaking: true,
			},
		},
		{
			name:    "breaking change with scope and bang",
			message: "refactor(api)!: rename endpoints",
			expected: &ConventionalCommit{
				Type:     "refactor",
				Scope:    "api",
				Subject:  "rename endpoints",
				Breaking: true,
			},
		},
		{
			name:    "fix commit",
			message: "fix: resolve race condition",
			expected: &ConventionalCommit{
				Type:     "fix",
				Subject:  "resolve race condition",
				Breaking: false,
			},
		},
		{
			name:    "docs commit",
			message: "docs(readme): update installation guide",
			expected: &ConventionalCommit{
				Type:     "docs",
				Scope:    "readme",
				Subject:  "update installation guide",
				Breaking: false,
			},
		},
		{
			name:    "uppercase type normalized to lowercase",
			message: "FEAT: add feature",
			expected: &ConventionalCommit{
				Type:     "feat",
				Subject:  "add feature",
				Breaking: false,
			},
		},
		{
			name:     "non-conventional commit",
			message:  "Update README.md",
			expected: nil,
		},
		{
			name:     "missing colon",
			message:  "feat add new feature",
			expected: nil,
		},
		{
			name:    "extra spaces around colon",
			message: "feat : add new feature",
			expected: &ConventionalCommit{
				Type:     "feat",
				Subject:  "add new feature",
				Breaking: false,
			},
		},
		{
			name:    "multiline message uses first line only",
			message: "feat: add feature\n\nThis is the body",
			expected: &ConventionalCommit{
				Type:     "feat",
				Subject:  "add feature",
				Breaking: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseConventionalCommit(tt.message)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("expected %+v, got nil", tt.expected)
				return
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Type: expected %q, got %q", tt.expected.Type, result.Type)
			}
			if result.Scope != tt.expected.Scope {
				t.Errorf("Scope: expected %q, got %q", tt.expected.Scope, result.Scope)
			}
			if result.Subject != tt.expected.Subject {
				t.Errorf("Subject: expected %q, got %q", tt.expected.Subject, result.Subject)
			}
			if result.Breaking != tt.expected.Breaking {
				t.Errorf("Breaking: expected %v, got %v", tt.expected.Breaking, result.Breaking)
			}
		})
	}
}

func TestIsConventionalCommit(t *testing.T) {
	tests := []struct {
		message  string
		expected bool
	}{
		{"feat: add feature", true},
		{"fix(scope): fix bug", true},
		{"Update README", false},
		{"feat add feature", false},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := IsConventionalCommit(tt.message)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExtractIssueNumber(t *testing.T) {
	tests := []struct {
		message  string
		expected int
	}{
		{"fix bug #123", 123},
		{"Closes #456", 456},
		{"closes #789", 789},
		{"Fixes #100", 100},
		{"fixes #200", 200},
		{"Resolves #300", 300},
		{"Refs #400", 400},
		{"No issue reference", 0},
		{"Multiple #1 and #2", 1}, // Returns first
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := ExtractIssueNumber(tt.message)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExtractPRNumber(t *testing.T) {
	tests := []struct {
		subject  string
		expected int
	}{
		{"feat: add feature (#123)", 123},
		{"fix bug (#456)", 456},
		{"Update README (#789)  ", 789},
		{"feat: add feature #123", 0}, // Not in parens at end
		{"(#100) at start", 0},        // Not at end
		{"No PR reference", 0},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			result := ExtractPRNumber(tt.subject)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestHasBreakingChangeMarker(t *testing.T) {
	tests := []struct {
		body     string
		expected bool
	}{
		{"BREAKING CHANGE: removes old API", true},
		{"BREAKING-CHANGE: removes old API", true},
		{"Breaking Change: removes old API", true},
		{"Some text\nBREAKING CHANGE: detail\nMore text", true},
		{"No breaking changes here", false},
		{"BREAKING NEWS: not a commit", false},
	}

	for _, tt := range tests {
		t.Run(tt.body[:min(20, len(tt.body))], func(t *testing.T) {
			result := HasBreakingChangeMarker(tt.body)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsKnownType(t *testing.T) {
	knownTypes := []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert", "security", "deps"}

	for _, typ := range knownTypes {
		t.Run(typ, func(t *testing.T) {
			if !IsKnownType(typ) {
				t.Errorf("%s should be a known type", typ)
			}
			// Test uppercase
			if !IsKnownType(typ) {
				t.Errorf("%s (uppercase) should be a known type", typ)
			}
		})
	}

	unknownTypes := []string{"feature", "bugfix", "misc", "other"}
	for _, typ := range unknownTypes {
		t.Run(typ, func(t *testing.T) {
			if IsKnownType(typ) {
				t.Errorf("%s should not be a known type", typ)
			}
		})
	}
}
