package renderer

import (
	"strings"
	"testing"

	"github.com/grokify/structured-changelog/changelog"
)

func TestRenderMarkdown_Basic(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test-project",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "Initial release"}},
			},
		},
	}

	md := RenderMarkdown(cl)

	// Check header
	if !strings.Contains(md, "# Changelog") {
		t.Error("missing changelog header")
	}
	if !strings.Contains(md, "Keep a Changelog") {
		t.Error("missing Keep a Changelog reference")
	}
	if !strings.Contains(md, "Semantic Versioning") {
		t.Error("missing Semantic Versioning reference")
	}

	// Check release
	if !strings.Contains(md, "## [1.0.0] - 2026-01-03") {
		t.Error("missing release header")
	}
	if !strings.Contains(md, "### Added") {
		t.Error("missing Added section")
	}
	if !strings.Contains(md, "- Initial release") {
		t.Error("missing entry")
	}
}

func TestRenderMarkdown_Unreleased(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Unreleased: &changelog.Release{
			Added: []changelog.Entry{{Description: "Work in progress"}},
		},
	}

	md := RenderMarkdown(cl)

	if !strings.Contains(md, "## [Unreleased]") {
		t.Error("missing Unreleased section")
	}
	if !strings.Contains(md, "- Work in progress") {
		t.Error("missing unreleased entry")
	}
}

func TestRenderMarkdown_AllCategories(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version:    "1.0.0",
				Date:       "2026-01-03",
				Added:      []changelog.Entry{{Description: "Added item"}},
				Changed:    []changelog.Entry{{Description: "Changed item"}},
				Deprecated: []changelog.Entry{{Description: "Deprecated item"}},
				Removed:    []changelog.Entry{{Description: "Removed item"}},
				Fixed:      []changelog.Entry{{Description: "Fixed item"}},
				Security:   []changelog.Entry{{Description: "Security item"}},
			},
		},
	}

	md := RenderMarkdown(cl)

	categories := []string{"Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}
	for _, cat := range categories {
		if !strings.Contains(md, "### "+cat) {
			t.Errorf("missing %s section", cat)
		}
	}
}

func TestRenderMarkdown_BreakingChange(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "2.0.0",
				Date:    "2026-01-03",
				Changed: []changelog.Entry{{Description: "API change", Breaking: true}},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	if !strings.Contains(md, "**BREAKING:**") {
		t.Error("missing BREAKING marker")
	}
}

func TestRenderMarkdown_SecurityMetadata(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.1",
				Date:    "2026-01-03",
				Security: []changelog.Entry{
					{
						Description: "Fix vulnerability",
						CVE:         "CVE-2026-12345",
						Severity:    "high",
					},
				},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	if !strings.Contains(md, "CVE-2026-12345") {
		t.Error("missing CVE in output")
	}
	if !strings.Contains(md, "severity: high") {
		t.Error("missing severity in output")
	}
}

func TestRenderMarkdown_MinimalOptions(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.1",
				Date:    "2026-01-03",
				Security: []changelog.Entry{
					{
						Description: "Fix vulnerability",
						CVE:         "CVE-2026-12345",
						Severity:    "high",
					},
				},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, MinimalOptions())

	// CVE should NOT be included with minimal options
	if strings.Contains(md, "CVE-2026-12345") {
		t.Error("CVE should not be included with minimal options")
	}
}

func TestRenderMarkdown_Yanked(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Yanked:  true,
				Added:   []changelog.Entry{{Description: "Bad release"}},
			},
		},
	}

	md := RenderMarkdown(cl)

	if !strings.Contains(md, "[YANKED]") {
		t.Error("missing YANKED marker")
	}
}

func TestRenderMarkdown_Deterministic(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "Feature A"}, {Description: "Feature B"}},
			},
		},
	}

	// Render multiple times
	md1 := RenderMarkdown(cl)
	md2 := RenderMarkdown(cl)
	md3 := RenderMarkdown(cl)

	if md1 != md2 || md2 != md3 {
		t.Error("rendering is not deterministic")
	}
}

func TestRenderMarkdown_PRReference(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "New feature", PR: "42"}},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	if !strings.Contains(md, "#42") {
		t.Error("missing PR reference")
	}
}

func TestRenderMarkdown_FullOptions(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added: []changelog.Entry{
					{Description: "New feature", PR: "42", Commit: "abc123"},
				},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, FullOptions())

	// Full options should include commit
	if !strings.Contains(md, "abc123") {
		t.Error("missing commit SHA with full options")
	}
}

func TestRenderMarkdown_CoreOptions(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version:     "1.0.0",
				Date:        "2026-01-03",
				Added:       []changelog.Entry{{Description: "Feature"}},    // core
				Performance: []changelog.Entry{{Description: "Faster"}},     // standard
				Internal:    []changelog.Entry{{Description: "Refactored"}}, // optional
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, CoreOptions())

	// Core options should only include core tier
	if !strings.Contains(md, "### Added") {
		t.Error("missing Added section (core tier)")
	}
	if strings.Contains(md, "### Performance") {
		t.Error("Performance section should not be included with core tier")
	}
	if strings.Contains(md, "### Internal") {
		t.Error("Internal section should not be included with core tier")
	}
}

func TestRenderMarkdown_StandardOptions(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version:     "1.0.0",
				Date:        "2026-01-03",
				Added:       []changelog.Entry{{Description: "Feature"}},    // core
				Performance: []changelog.Entry{{Description: "Faster"}},     // standard
				Internal:    []changelog.Entry{{Description: "Refactored"}}, // optional
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, StandardOptions())

	// Standard options should include core + standard tiers
	if !strings.Contains(md, "### Added") {
		t.Error("missing Added section (core tier)")
	}
	if !strings.Contains(md, "### Performance") {
		t.Error("missing Performance section (standard tier)")
	}
	if strings.Contains(md, "### Internal") {
		t.Error("Internal section should not be included with standard tier")
	}
}

func TestRenderMarkdown_TierFiltering(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version:       "1.0.0",
				Date:          "2026-01-03",
				Highlights:    []changelog.Entry{{Description: "Summary"}},      // standard
				Security:      []changelog.Entry{{Description: "Fix CVE"}},      // core
				Added:         []changelog.Entry{{Description: "Feature"}},      // core
				Documentation: []changelog.Entry{{Description: "Updated docs"}}, // extended
				Internal:      []changelog.Entry{{Description: "Refactored"}},   // optional
			},
		},
	}

	// Test extended tier filtering
	opts := DefaultOptions()
	opts.MaxTier = changelog.TierExtended
	md := RenderMarkdownWithOptions(cl, opts)

	if !strings.Contains(md, "### Documentation") {
		t.Error("missing Documentation section (extended tier)")
	}
	if strings.Contains(md, "### Internal") {
		t.Error("Internal should not be included with extended tier")
	}
}

func TestRenderMarkdown_ReferenceLinks(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion:  "1.0",
		Project:    "test",
		Repository: "https://github.com/example/repo",
		Releases: []changelog.Release{
			{Version: "1.1.0", Date: "2026-01-04", Added: []changelog.Entry{{Description: "New"}}},
			{Version: "1.0.0", Date: "2026-01-03", Added: []changelog.Entry{{Description: "Initial"}}},
		},
	}

	md := RenderMarkdown(cl)

	// Check for reference links
	if !strings.Contains(md, "[1.1.0]: https://github.com/example/repo/compare/v1.0.0...v1.1.0") {
		t.Error("missing compare link for 1.1.0")
	}
	if !strings.Contains(md, "[1.0.0]: https://github.com/example/repo/releases/tag/v1.0.0") {
		t.Error("missing tag link for 1.0.0")
	}
}

func TestRenderMarkdown_ReferenceLinks_WithUnreleased(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion:  "1.0",
		Project:    "test",
		Repository: "https://github.com/example/repo",
		Unreleased: &changelog.Release{
			Added: []changelog.Entry{{Description: "WIP"}},
		},
		Releases: []changelog.Release{
			{Version: "1.0.0", Date: "2026-01-03", Added: []changelog.Entry{{Description: "Initial"}}},
		},
	}

	md := RenderMarkdown(cl)

	// Check for unreleased link
	if !strings.Contains(md, "[unreleased]: https://github.com/example/repo/compare/v1.0.0...HEAD") {
		t.Error("missing unreleased compare link")
	}
}

func TestRenderMarkdown_ReferenceLinks_NonGitHub(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion:  "1.0",
		Project:    "test",
		Repository: "https://gitlab.com/example/repo",
		Releases: []changelog.Release{
			{Version: "1.0.0", Date: "2026-01-03", Added: []changelog.Entry{{Description: "Initial"}}},
		},
	}

	md := RenderMarkdown(cl)

	// Non-GitHub repos should not have reference links
	if strings.Contains(md, "[1.0.0]:") {
		t.Error("non-GitHub repos should not have reference links")
	}
}

func TestRenderMarkdown_ReferenceLinks_Disabled(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion:  "1.0",
		Project:    "test",
		Repository: "https://github.com/example/repo",
		Releases: []changelog.Release{
			{Version: "1.0.0", Date: "2026-01-03", Added: []changelog.Entry{{Description: "Initial"}}},
		},
	}

	opts := DefaultOptions()
	opts.IncludeCompareLinks = false
	md := RenderMarkdownWithOptions(cl, opts)

	// Should not have reference links
	if strings.Contains(md, "[1.0.0]:") {
		t.Error("reference links should be disabled")
	}
}

func TestRenderMarkdown_IssueReference(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Fixed:   []changelog.Entry{{Description: "Bug fix", Issue: "123"}},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	if !strings.Contains(md, "#123") {
		t.Error("missing issue reference")
	}
}

func TestRenderMarkdown_URLReference(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "Feature", Issue: "https://github.com/example/repo/issues/123"}},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	if !strings.Contains(md, "[issue](https://github.com/example/repo/issues/123)") {
		t.Error("missing URL reference link")
	}
}

func TestRenderMarkdown_HashPrefixedReference(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "Feature", Issue: "#456"}},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	// Should keep the hash prefix as-is
	if !strings.Contains(md, "#456") {
		t.Error("missing hash-prefixed reference")
	}
}

func TestRenderMarkdown_GHSAReference(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Security: []changelog.Entry{
					{Description: "Fix", GHSA: "GHSA-abcd-efgh-ijkl"},
				},
			},
		},
	}

	md := RenderMarkdownWithOptions(cl, DefaultOptions())

	if !strings.Contains(md, "GHSA-abcd-efgh-ijkl") {
		t.Error("missing GHSA reference")
	}
}

func TestRenderMarkdown_NoBreakingMarker(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "2.0.0",
				Date:    "2026-01-03",
				Changed: []changelog.Entry{{Description: "API change", Breaking: true}},
			},
		},
	}

	opts := DefaultOptions()
	opts.MarkBreakingChanges = false
	md := RenderMarkdownWithOptions(cl, opts)

	if strings.Contains(md, "**BREAKING:**") {
		t.Error("BREAKING marker should not be present when disabled")
	}
}

func TestRenderMarkdown_EmptyUnreleased(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion:  "1.0",
		Project:    "test",
		Unreleased: &changelog.Release{}, // Empty unreleased
		Releases: []changelog.Release{
			{Version: "1.0.0", Date: "2026-01-03", Added: []changelog.Entry{{Description: "Initial"}}},
		},
	}

	md := RenderMarkdown(cl)

	// Empty unreleased should not be rendered
	if strings.Contains(md, "## [Unreleased]") {
		t.Error("empty unreleased section should not be rendered")
	}
}

func TestRenderMarkdown_CommitReference(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "Feature", Commit: "abc123def"}},
			},
		},
	}

	// Default options should NOT include commits
	md := RenderMarkdownWithOptions(cl, DefaultOptions())
	if strings.Contains(md, "abc123def") {
		t.Error("commits should not be included with default options")
	}

	// Full options should include commits
	md = RenderMarkdownWithOptions(cl, FullOptions())
	if !strings.Contains(md, "abc123def") {
		t.Error("commits should be included with full options")
	}
}

func TestRenderMarkdown_AllExtendedCategories(t *testing.T) {
	cl := &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version:        "1.0.0",
				Date:           "2026-01-03",
				Highlights:     []changelog.Entry{{Description: "h"}},
				Breaking:       []changelog.Entry{{Description: "b"}},
				UpgradeGuide:   []changelog.Entry{{Description: "u"}},
				Security:       []changelog.Entry{{Description: "s"}},
				Added:          []changelog.Entry{{Description: "a"}},
				Changed:        []changelog.Entry{{Description: "c"}},
				Deprecated:     []changelog.Entry{{Description: "d"}},
				Removed:        []changelog.Entry{{Description: "r"}},
				Fixed:          []changelog.Entry{{Description: "f"}},
				Performance:    []changelog.Entry{{Description: "p"}},
				Dependencies:   []changelog.Entry{{Description: "dep"}},
				Documentation:  []changelog.Entry{{Description: "doc"}},
				Build:          []changelog.Entry{{Description: "bld"}},
				Infrastructure: []changelog.Entry{{Description: "i"}},
				Observability:  []changelog.Entry{{Description: "o"}},
				Compliance:     []changelog.Entry{{Description: "comp"}},
				Internal:       []changelog.Entry{{Description: "int"}},
				KnownIssues:    []changelog.Entry{{Description: "k"}},
				Contributors:   []changelog.Entry{{Description: "cont"}},
			},
		},
	}

	md := RenderMarkdown(cl)

	// Check all 19 categories are present
	categories := []string{
		"Highlights", "Breaking", "Upgrade Guide", "Security",
		"Added", "Changed", "Deprecated", "Removed", "Fixed",
		"Performance", "Dependencies",
		"Documentation", "Build",
		"Infrastructure", "Observability", "Compliance",
		"Internal",
		"Known Issues", "Contributors",
	}

	for _, cat := range categories {
		if !strings.Contains(md, "### "+cat) {
			t.Errorf("missing %s section", cat)
		}
	}
}

func TestOptions_DefaultMaxTier(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxTier != changelog.TierOptional {
		t.Errorf("DefaultOptions MaxTier = %q, want %q", opts.MaxTier, changelog.TierOptional)
	}
}

func TestOptions_MinimalMaxTier(t *testing.T) {
	opts := MinimalOptions()
	if opts.MaxTier != changelog.TierCore {
		t.Errorf("MinimalOptions MaxTier = %q, want %q", opts.MaxTier, changelog.TierCore)
	}
}

func TestOptions_FullMaxTier(t *testing.T) {
	opts := FullOptions()
	if opts.MaxTier != changelog.TierOptional {
		t.Errorf("FullOptions MaxTier = %q, want %q", opts.MaxTier, changelog.TierOptional)
	}
}

func TestOptions_CoreMaxTier(t *testing.T) {
	opts := CoreOptions()
	if opts.MaxTier != changelog.TierCore {
		t.Errorf("CoreOptions MaxTier = %q, want %q", opts.MaxTier, changelog.TierCore)
	}
}

func TestOptions_StandardMaxTier(t *testing.T) {
	opts := StandardOptions()
	if opts.MaxTier != changelog.TierStandard {
		t.Errorf("StandardOptions MaxTier = %q, want %q", opts.MaxTier, changelog.TierStandard)
	}
}
