// Package renderer provides deterministic Markdown rendering for changelogs.
package renderer

import (
	"fmt"
	"strings"

	"github.com/grokify/structured-changelog/changelog"
)

// RenderMarkdown renders a changelog to Keep a Changelog formatted Markdown.
// The output is deterministic: same input always produces identical output.
func RenderMarkdown(cl *changelog.Changelog) string {
	return RenderMarkdownWithOptions(cl, DefaultOptions())
}

// RenderMarkdownWithOptions renders a changelog with custom options.
func RenderMarkdownWithOptions(cl *changelog.Changelog, opts Options) string {
	var sb strings.Builder

	// Header
	sb.WriteString("# Changelog\n\n")
	sb.WriteString("All notable changes to this project will be documented in this file.\n\n")
	sb.WriteString("The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),\n")
	sb.WriteString("and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n")

	// Unreleased section
	if cl.Unreleased != nil && !cl.Unreleased.IsEmpty() {
		sb.WriteString("\n## [Unreleased]\n")
		renderReleaseContent(&sb, cl.Unreleased, opts)
	}

	// Releases
	for _, release := range cl.Releases {
		sb.WriteString("\n")
		renderRelease(&sb, &release, opts)
	}

	return sb.String()
}

func renderRelease(sb *strings.Builder, r *changelog.Release, opts Options) {
	// Version header
	if r.Yanked {
		fmt.Fprintf(sb, "## [%s] - %s [YANKED]\n", r.Version, r.Date)
	} else {
		fmt.Fprintf(sb, "## [%s] - %s\n", r.Version, r.Date)
	}

	renderReleaseContent(sb, r, opts)
}

func renderReleaseContent(sb *strings.Builder, r *changelog.Release, opts Options) {
	// Render categories in standard order
	for _, cat := range r.Categories() {
		fmt.Fprintf(sb, "\n### %s\n\n", cat.Name)
		for _, entry := range cat.Entries {
			renderEntry(sb, &entry, opts, cat.Name == "Security")
		}
	}
}

func renderEntry(sb *strings.Builder, e *changelog.Entry, opts Options, isSecurity bool) {
	// Build the entry line
	var parts []string

	// Description (required)
	desc := e.Description
	if e.Breaking && opts.MarkBreakingChanges {
		desc = "**BREAKING:** " + desc
	}
	parts = append(parts, desc)

	// References
	var refs []string
	if e.Issue != "" && opts.IncludeReferences {
		refs = append(refs, formatRef("issue", e.Issue))
	}
	if e.PR != "" && opts.IncludeReferences {
		refs = append(refs, formatRef("PR", e.PR))
	}
	if e.Commit != "" && opts.IncludeReferences && opts.IncludeCommits {
		refs = append(refs, formatRef("commit", e.Commit))
	}

	// Security metadata
	if isSecurity && opts.IncludeSecurityMetadata {
		if e.CVE != "" {
			refs = append(refs, e.CVE)
		}
		if e.GHSA != "" {
			refs = append(refs, e.GHSA)
		}
		if e.Severity != "" {
			refs = append(refs, fmt.Sprintf("severity: %s", e.Severity))
		}
	}

	// Combine parts
	line := strings.Join(parts, " ")
	if len(refs) > 0 {
		line += " (" + strings.Join(refs, ", ") + ")"
	}

	sb.WriteString("- " + line + "\n")
}

func formatRef(refType, value string) string {
	// If it's already a URL, just use it
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return fmt.Sprintf("[%s](%s)", refType, value)
	}
	// Otherwise, just show the reference
	if strings.HasPrefix(value, "#") {
		return value
	}
	return fmt.Sprintf("#%s", value)
}
