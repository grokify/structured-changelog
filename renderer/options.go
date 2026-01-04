package renderer

// Options controls how the Markdown is rendered.
type Options struct {
	// IncludeReferences includes issue/PR links in entries.
	IncludeReferences bool

	// IncludeCommits includes commit SHAs in references.
	IncludeCommits bool

	// IncludeSecurityMetadata includes CVE/GHSA/severity in security entries.
	IncludeSecurityMetadata bool

	// MarkBreakingChanges prefixes breaking changes with **BREAKING:**.
	MarkBreakingChanges bool

	// IncludeCompareLinks adds version comparison links at the bottom.
	IncludeCompareLinks bool
}

// DefaultOptions returns the default rendering options.
func DefaultOptions() Options {
	return Options{
		IncludeReferences:       true,
		IncludeCommits:          false,
		IncludeSecurityMetadata: true,
		MarkBreakingChanges:     true,
		IncludeCompareLinks:     true,
	}
}

// MinimalOptions returns options for minimal output.
func MinimalOptions() Options {
	return Options{
		IncludeReferences:       false,
		IncludeCommits:          false,
		IncludeSecurityMetadata: false,
		MarkBreakingChanges:     false,
		IncludeCompareLinks:     false,
	}
}

// FullOptions returns options for maximum detail.
func FullOptions() Options {
	return Options{
		IncludeReferences:       true,
		IncludeCommits:          true,
		IncludeSecurityMetadata: true,
		MarkBreakingChanges:     true,
		IncludeCompareLinks:     true,
	}
}
