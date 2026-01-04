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
