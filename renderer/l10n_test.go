package renderer

import (
	"strings"
	"testing"

	"github.com/grokify/structured-changelog/changelog"
)

func TestGetLocalizer(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		msgID    string
		expected string
	}{
		{"english title", "en", "changelog.title", "Changelog"},
		{"french title", "fr", "changelog.title", "Journal des modifications"},
		{"german title", "de", "changelog.title", "Änderungsprotokoll"},
		{"spanish title", "es", "changelog.title", "Registro de cambios"},
		{"japanese title", "ja", "changelog.title", "変更履歴"},
		{"chinese title", "zh", "changelog.title", "更新日志"},
		{"english category", "en", "category.added", "Added"},
		{"french category", "fr", "category.added", "Ajouté"},
		{"fallback to english", "xx", "changelog.title", "Changelog"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithLocale(tt.locale)
			l := getLocalizer(opts)
			got := l.T(tt.msgID)
			if got != tt.expected {
				t.Errorf("getLocalizer(%q).T(%q) = %q, expected %q",
					tt.locale, tt.msgID, got, tt.expected)
			}
		})
	}
}

func TestCategoryToMessageID(t *testing.T) {
	tests := []struct {
		category string
		expected string
	}{
		{"Added", "category.added"},
		{"Fixed", "category.fixed"},
		{"Known Issues", "category.known_issues"},
		{"Upgrade Guide", "category.upgrade_guide"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			got := categoryToMessageID(tt.category)
			if got != tt.expected {
				t.Errorf("categoryToMessageID(%q) = %q, expected %q",
					tt.category, got, tt.expected)
			}
		})
	}
}

func TestRenderMarkdownWithLocale(t *testing.T) {
	cl := &changelog.Changelog{
		Repository: "https://github.com/test/repo",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2024-01-15",
				Added: []changelog.Entry{
					{Description: "New feature"},
				},
			},
		},
	}

	tests := []struct {
		name       string
		locale     string
		contains   []string
		notContain []string
	}{
		{
			name:     "english output",
			locale:   "en",
			contains: []string{"# Changelog", "### Added", "All notable changes"},
		},
		{
			name:     "french output",
			locale:   "fr",
			contains: []string{"# Journal des modifications", "### Ajouté", "Tous les changements notables"},
		},
		{
			name:     "japanese output",
			locale:   "ja",
			contains: []string{"# 変更履歴", "### 追加"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithLocale(tt.locale)
			md := RenderMarkdownWithOptions(cl, opts)

			for _, s := range tt.contains {
				if !strings.Contains(md, s) {
					t.Errorf("Expected output to contain %q, but it doesn't.\nOutput:\n%s", s, md)
				}
			}
			for _, s := range tt.notContain {
				if strings.Contains(md, s) {
					t.Errorf("Expected output NOT to contain %q, but it does.\nOutput:\n%s", s, md)
				}
			}
		})
	}
}

func TestRenderBreakingWithLocale(t *testing.T) {
	cl := &changelog.Changelog{
		Releases: []changelog.Release{
			{
				Version: "2.0.0",
				Date:    "2024-01-15",
				Changed: []changelog.Entry{
					{Description: "Changed API", Breaking: true},
				},
			},
		},
	}

	tests := []struct {
		name     string
		locale   string
		contains string
	}{
		{"english breaking", "en", "**BREAKING:**"},
		{"french breaking", "fr", "**RUPTURE :**"},
		{"japanese breaking", "ja", "**破壊的変更:**"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithLocale(tt.locale)
			md := RenderMarkdownWithOptions(cl, opts)

			if !strings.Contains(md, tt.contains) {
				t.Errorf("Expected output to contain %q for locale %q.\nOutput:\n%s",
					tt.contains, tt.locale, md)
			}
		})
	}
}

func TestRenderYankedWithLocale(t *testing.T) {
	cl := &changelog.Changelog{
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2024-01-15",
				Yanked:  true,
				Added:   []changelog.Entry{{Description: "Some feature"}}, // Need entry for notable release
			},
		},
	}

	tests := []struct {
		name     string
		locale   string
		contains string
	}{
		{"english yanked", "en", "[YANKED]"},
		{"french yanked", "fr", "[RETIRÉ]"},
		{"german yanked", "de", "[ZURÜCKGEZOGEN]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithLocale(tt.locale)
			md := RenderMarkdownWithOptions(cl, opts)

			if !strings.Contains(md, tt.contains) {
				t.Errorf("Expected output to contain %q for locale %q.\nOutput:\n%s",
					tt.contains, tt.locale, md)
			}
		})
	}
}
