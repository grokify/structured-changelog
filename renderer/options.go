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
}

// DefaultOptions returns the default rendering options.
// Includes commit links and reference linking when repository URL is available.
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
	}
}

// FullOptions returns options for maximum detail.
// Same as DefaultOptions but with CompactMaintenanceReleases disabled
// to show all releases expanded instead of grouping maintenance releases.
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
	Preset          string // default, minimal, full, core, standard
	MaxTier         string // optional tier override
	Locale          string // optional BCP 47 locale tag override
	LocaleOverrides string // optional path to locale override JSON file
}

// OptionsFromConfig creates Options from a Config struct.
// It first applies the preset, then overrides MaxTier, Locale, and LocaleOverrides if specified.
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

	return opts, nil
}
