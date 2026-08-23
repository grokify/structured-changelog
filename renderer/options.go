package renderer

import (
	"errors"

	"github.com/grokify/structured-changelog/changelog"
)

// Options controls how the Markdown is rendered.
type Options struct {
	// IncludeReferences includes issue/PR links in entries.
	IncludeReferences bool

	// IncludeCommits includes commit SHAs in references.
	IncludeCommits bool

	// LinkReferences creates hyperlinks for issues, PRs, and commits
	// when a repository URL is available. Requires IncludeReferences.
	LinkReferences bool

	// IncludeAuthors appends author attribution for external contributors.
	// Authors listed in Changelog.Maintainers or known bots are excluded.
	IncludeAuthors bool

	// IncludeSecurityMetadata includes CVE/GHSA/severity in security entries.
	IncludeSecurityMetadata bool

	// MarkBreakingChanges prefixes breaking changes with **BREAKING:**.
	MarkBreakingChanges bool

	// IncludeCompareLinks adds version comparison links at the bottom.
	IncludeCompareLinks bool

	// IncludeUnreleasedLink adds an [Unreleased] link comparing latest version to HEAD.
	// This lets users see what's been merged since the last release.
	IncludeUnreleasedLink bool

	// CompactMaintenanceReleases groups consecutive maintenance-only releases
	// (those with only dependencies, documentation, build, tests, internal changes)
	// into a single compact section like "## Versions 0.71.1 - 0.71.10 (Maintenance)".
	CompactMaintenanceReleases bool

	// MaxTier filters change types to include only those at or above this tier.
	// Default is TierOptional (include all).
	MaxTier changelog.Tier

	// Locale specifies the BCP 47 locale tag for output (e.g., "en", "fr", "de").
	// Default is "en" (English).
	Locale string

	// LocaleOverrides specifies a path to a JSON file with locale message overrides.
	// Only the messages specified in this file will be replaced; others use defaults.
	LocaleOverrides string

	// NotableOnly when true, only includes releases that are considered "notable"
	// according to the NotabilityPolicy. Non-notable releases (maintenance-only)
	// are excluded from the output entirely.
	NotableOnly bool

	// NotabilityPolicy defines which categories make a release notable.
	// If nil and NotableOnly is true, uses DefaultNotabilityPolicy().
	NotabilityPolicy *changelog.NotabilityPolicy

	// ExcludeCategories lists category names to omit entirely from the output
	// (e.g. "Dependencies"). Excluded categories produce no section at all.
	// This applies on top of MaxTier filtering.
	ExcludeCategories []string

	// CollapseCategories lists category names to render as a single summary
	// line (an entry count) instead of listing every entry. This preserves the
	// signal that maintenance occurred without the noise of a long list; the
	// full detail remains in the JSON source. ExcludeCategories takes
	// precedence: a category in both is omitted, not collapsed.
	CollapseCategories []string
}

// DefaultOptions returns the default rendering options.
// Includes commit links and reference linking when repository URL is available.
// By default, only notable releases are included (NotableOnly: true).
func DefaultOptions() Options {
	return Options{
		IncludeReferences:          true,
		IncludeCommits:             true,
		LinkReferences:             true,
		IncludeAuthors:             true,
		IncludeSecurityMetadata:    true,
		MarkBreakingChanges:        true,
		IncludeCompareLinks:        true,
		IncludeUnreleasedLink:      true,
		CompactMaintenanceReleases: true,
		MaxTier:                    changelog.TierOptional,
		Locale:                     "en",
		NotableOnly:                true,
		NotabilityPolicy:           changelog.DefaultNotabilityPolicy(),
		// Dependency churn is the dominant source of changelog noise for human
		// readers, so collapse it to a one-line count by default. The full list
		// stays in the JSON and is one --full (or --expand-categories) away.
		CollapseCategories: []string{changelog.CategoryDependencies},
	}
}

// MinimalOptions returns options for minimal output.
func MinimalOptions() Options {
	return Options{
		IncludeReferences:          false,
		IncludeCommits:             false,
		LinkReferences:             false,
		IncludeAuthors:             false,
		IncludeSecurityMetadata:    false,
		MarkBreakingChanges:        false,
		IncludeCompareLinks:        false,
		IncludeUnreleasedLink:      false,
		CompactMaintenanceReleases: true,
		MaxTier:                    changelog.TierCore,
		Locale:                     "en",
		NotableOnly:                true,
		NotabilityPolicy:           changelog.DefaultNotabilityPolicy(),
	}
}

// FullOptions returns options for maximum detail.
// Includes all releases (NotableOnly: false) and shows them expanded
// instead of grouping maintenance releases.
func FullOptions() Options {
	return Options{
		IncludeReferences:          true,
		IncludeCommits:             true,
		LinkReferences:             true,
		IncludeAuthors:             true,
		IncludeSecurityMetadata:    true,
		MarkBreakingChanges:        true,
		IncludeCompareLinks:        true,
		IncludeUnreleasedLink:      true,
		CompactMaintenanceReleases: false, // Full detail shows all releases expanded
		MaxTier:                    changelog.TierOptional,
		Locale:                     "en",
		NotableOnly:                false, // Full includes all releases
	}
}

// CoreOptions returns options for KACL-compliant core output.
func CoreOptions() Options {
	return Options{
		IncludeReferences:          true,
		IncludeCommits:             false,
		LinkReferences:             false,
		IncludeAuthors:             true,
		IncludeSecurityMetadata:    true,
		MarkBreakingChanges:        true,
		IncludeCompareLinks:        true,
		IncludeUnreleasedLink:      true,
		CompactMaintenanceReleases: true,
		MaxTier:                    changelog.TierCore,
		Locale:                     "en",
		NotableOnly:                true,
		NotabilityPolicy:           changelog.DefaultNotabilityPolicy(),
	}
}

// StandardOptions returns options including standard tier types.
func StandardOptions() Options {
	return Options{
		IncludeReferences:          true,
		IncludeCommits:             false,
		LinkReferences:             false,
		IncludeAuthors:             true,
		IncludeSecurityMetadata:    true,
		MarkBreakingChanges:        true,
		IncludeCompareLinks:        true,
		IncludeUnreleasedLink:      true,
		CompactMaintenanceReleases: true,
		MaxTier:                    changelog.TierStandard,
		Locale:                     "en",
		NotableOnly:                true,
		NotabilityPolicy:           changelog.DefaultNotabilityPolicy(),
		CollapseCategories:         []string{changelog.CategoryDependencies},
	}
}

// WithMaxTier returns a copy of the options with the MaxTier field set.
func (o Options) WithMaxTier(tier changelog.Tier) Options {
	o.MaxTier = tier
	return o
}

// WithLocale returns a copy of the options with the Locale field set.
func (o Options) WithLocale(locale string) Options {
	o.Locale = locale
	return o
}

// WithLocaleOverrides returns a copy of the options with the LocaleOverrides field set.
func (o Options) WithLocaleOverrides(path string) Options {
	o.LocaleOverrides = path
	return o
}

// WithNotableOnly returns a copy of the options with NotableOnly set.
// When enabled, only releases with entries in notable categories are included.
func (o Options) WithNotableOnly(enabled bool) Options {
	o.NotableOnly = enabled
	return o
}

// WithNotabilityPolicy returns a copy of the options with a custom NotabilityPolicy.
func (o Options) WithNotabilityPolicy(policy *changelog.NotabilityPolicy) Options {
	o.NotabilityPolicy = policy
	return o
}

// WithExcludeCategories returns a copy of the options with the given category
// names appended to ExcludeCategories.
func (o Options) WithExcludeCategories(categories ...string) Options {
	o.ExcludeCategories = append(append([]string{}, o.ExcludeCategories...), categories...)
	return o
}

// WithCollapseCategories returns a copy of the options with the given category
// names appended to CollapseCategories.
func (o Options) WithCollapseCategories(categories ...string) Options {
	o.CollapseCategories = append(append([]string{}, o.CollapseCategories...), categories...)
	return o
}

// WithExpandCategories returns a copy of the options with the given category
// names removed from CollapseCategories, forcing them to render in full. This
// is how a caller overrides a preset's default collapse (e.g. Dependencies).
func (o Options) WithExpandCategories(categories ...string) Options {
	if len(categories) == 0 || len(o.CollapseCategories) == 0 {
		return o
	}
	expand := make(map[string]bool, len(categories))
	for _, c := range categories {
		expand[c] = true
	}
	kept := make([]string, 0, len(o.CollapseCategories))
	for _, c := range o.CollapseCategories {
		if !expand[c] {
			kept = append(kept, c)
		}
	}
	o.CollapseCategories = kept
	return o
}

// OptionsFromPreset returns options for the given preset name.
// Valid presets are: default, minimal, full, core, standard.
func OptionsFromPreset(preset string) (Options, error) {
	switch preset {
	case "", "default":
		return DefaultOptions(), nil
	case "minimal":
		return MinimalOptions(), nil
	case "full":
		return FullOptions(), nil
	case "core":
		return CoreOptions(), nil
	case "standard":
		return StandardOptions(), nil
	default:
		return Options{}, ErrInvalidPreset
	}
}

// ErrInvalidPreset is returned when an invalid options preset name is provided.
var ErrInvalidPreset = errors.New("invalid preset")

// Config holds configuration for rendering options.
type Config struct {
	Preset             string   // default, minimal, full, core, standard
	MaxTier            string   // optional tier override
	Locale             string   // optional BCP 47 locale tag override
	LocaleOverrides    string   // optional path to locale override JSON file
	AllReleases        bool     // include all releases (overrides default notable-only)
	NotableCategories  []string // custom notable categories (uses default if empty)
	ExcludeCategories  []string // categories to omit entirely (appended to preset)
	CollapseCategories []string // categories to render as a summary line (appended to preset)
	ExpandCategories   []string // categories to force-expand, overriding preset collapse defaults
}

// OptionsFromConfig creates Options from a Config struct.
// It first applies the preset, then overrides MaxTier, Locale, LocaleOverrides,
// and notability settings if specified.
func OptionsFromConfig(cfg Config) (Options, error) {
	opts, err := OptionsFromPreset(cfg.Preset)
	if err != nil {
		return Options{}, err
	}

	if cfg.MaxTier != "" {
		tier, err := changelog.ParseTier(cfg.MaxTier)
		if err != nil {
			return Options{}, err
		}
		opts = opts.WithMaxTier(tier)
	}

	if cfg.Locale != "" {
		opts = opts.WithLocale(cfg.Locale)
	}

	if cfg.LocaleOverrides != "" {
		opts = opts.WithLocaleOverrides(cfg.LocaleOverrides)
	}

	// AllReleases overrides the default notable-only behavior
	if cfg.AllReleases {
		opts = opts.WithNotableOnly(false)
		opts.NotabilityPolicy = nil
	} else if len(cfg.NotableCategories) > 0 {
		// Custom notable categories (only applies when not AllReleases)
		opts = opts.WithNotabilityPolicy(changelog.NewNotabilityPolicy(cfg.NotableCategories))
	}

	// Category-level render filters, applied on top of the preset. Exclude and
	// collapse are appended to preset defaults; expand removes categories from
	// the (possibly preset-seeded) collapse set. Expand is applied last so it
	// can override both the preset default and an explicit --collapse-categories.
	if len(cfg.ExcludeCategories) > 0 {
		opts = opts.WithExcludeCategories(cfg.ExcludeCategories...)
	}
	if len(cfg.CollapseCategories) > 0 {
		opts = opts.WithCollapseCategories(cfg.CollapseCategories...)
	}
	if len(cfg.ExpandCategories) > 0 {
		opts = opts.WithExpandCategories(cfg.ExpandCategories...)
	}

	return opts, nil
}
