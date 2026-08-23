package renderer

import (
	"strings"
	"testing"

	"github.com/grokify/structured-changelog/changelog"
)

// notableReleaseWithDeps returns a changelog whose single release is notable
// (has an Added entry) and also carries dependency and documentation entries.
func notableReleaseWithDeps() *changelog.Changelog {
	return &changelog.Changelog{
		IRVersion: "1.0",
		Project:   "test",
		Releases: []changelog.Release{
			{
				Version: "1.0.0",
				Date:    "2026-01-03",
				Added:   []changelog.Entry{{Description: "New feature"}},
				Dependencies: []changelog.Entry{
					{Description: "Bump foo from 1.0.0 to 1.1.0"},
					{Description: "Bump bar from 2.0.0 to 2.1.0"},
				},
				Documentation: []changelog.Entry{
					{Description: "Update README"},
				},
			},
		},
	}
}

func TestRenderMarkdown_CollapsesDependenciesByDefault(t *testing.T) {
	md := RenderMarkdownWithOptions(notableReleaseWithDeps(), DefaultOptions())

	// The section header is still present.
	if !strings.Contains(md, "### Dependencies") {
		t.Fatalf("expected Dependencies section header, got:\n%s", md)
	}
	// The summary line is present.
	if !strings.Contains(md, "2 dependency updates") {
		t.Errorf("expected collapsed dependency summary, got:\n%s", md)
	}
	// The individual dependency descriptions are NOT listed.
	if strings.Contains(md, "Bump foo") || strings.Contains(md, "Bump bar") {
		t.Errorf("expected individual dependency entries to be collapsed, got:\n%s", md)
	}
	// The notable Added entry is still rendered in full.
	if !strings.Contains(md, "New feature") {
		t.Errorf("expected Added entry to render, got:\n%s", md)
	}
}

func TestRenderMarkdown_FullOptionsExpandsDependencies(t *testing.T) {
	md := RenderMarkdownWithOptions(notableReleaseWithDeps(), FullOptions())

	if !strings.Contains(md, "Bump foo") || !strings.Contains(md, "Bump bar") {
		t.Errorf("full output should list every dependency entry, got:\n%s", md)
	}
	if strings.Contains(md, "2 dependency updates") {
		t.Errorf("full output should not collapse dependencies, got:\n%s", md)
	}
}

func TestRenderMarkdown_ExpandCategoriesOverridesDefault(t *testing.T) {
	opts := DefaultOptions().WithExpandCategories(changelog.CategoryDependencies)
	md := RenderMarkdownWithOptions(notableReleaseWithDeps(), opts)

	if !strings.Contains(md, "Bump foo") || !strings.Contains(md, "Bump bar") {
		t.Errorf("expand should restore full dependency list, got:\n%s", md)
	}
	if strings.Contains(md, "2 dependency updates") {
		t.Errorf("expand should suppress the collapsed summary, got:\n%s", md)
	}
}

func TestRenderMarkdown_ExcludeCategoriesOmitsSection(t *testing.T) {
	opts := DefaultOptions().WithExcludeCategories(changelog.CategoryDependencies)
	md := RenderMarkdownWithOptions(notableReleaseWithDeps(), opts)

	if strings.Contains(md, "### Dependencies") {
		t.Errorf("excluded category should produce no section, got:\n%s", md)
	}
	if strings.Contains(md, "2 dependency updates") || strings.Contains(md, "Bump foo") {
		t.Errorf("excluded category should render nothing, got:\n%s", md)
	}
}

func TestRenderMarkdown_ExcludeTakesPrecedenceOverCollapse(t *testing.T) {
	opts := DefaultOptions().
		WithCollapseCategories(changelog.CategoryDependencies).
		WithExcludeCategories(changelog.CategoryDependencies)
	md := RenderMarkdownWithOptions(notableReleaseWithDeps(), opts)

	if strings.Contains(md, "### Dependencies") {
		t.Errorf("exclude should win over collapse, got:\n%s", md)
	}
}

func TestRenderMarkdown_CollapseArbitraryCategory(t *testing.T) {
	opts := DefaultOptions().WithCollapseCategories(changelog.CategoryDocumentation)
	md := RenderMarkdownWithOptions(notableReleaseWithDeps(), opts)

	if !strings.Contains(md, "### Documentation") {
		t.Fatalf("expected Documentation section, got:\n%s", md)
	}
	// One documentation entry -> singular generic/known summary.
	if strings.Contains(md, "Update README") {
		t.Errorf("documentation entry should be collapsed, got:\n%s", md)
	}
	if !strings.Contains(md, "1 documentation change") {
		t.Errorf("expected collapsed documentation summary, got:\n%s", md)
	}
}
