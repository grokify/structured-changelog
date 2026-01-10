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

	// MaxTier filters change types to include only those at or above this tier.
	// Default is TierOptional (include all).
	MaxTier changelog.Tier
}

// DefaultOptions returns the default rendering options.
func DefaultOptions() Options {
	return Options{
		IncludeReferences:       true,
		IncludeCommits:          false,
		LinkReferences:          false,
		IncludeAuthors:          true,
		IncludeSecurityMetadata: true,
		MarkBreakingChanges:     true,
		IncludeCompareLinks:     true,
		MaxTier:                 changelog.TierOptional,
	}
}

// MinimalOptions returns options for minimal output.
func MinimalOptions() Options {
	return Options{
		IncludeReferences:       false,
		IncludeCommits:          false,
		LinkReferences:          false,
		IncludeAuthors:          false,
		IncludeSecurityMetadata: false,
		MarkBreakingChanges:     false,
		IncludeCompareLinks:     false,
		MaxTier:                 changelog.TierCore,
	}
}

// FullOptions returns options for maximum detail.
func FullOptions() Options {
	return Options{
		IncludeReferences:       true,
		IncludeCommits:          true,
		LinkReferences:          true,
		IncludeAuthors:          true,
		IncludeSecurityMetadata: true,
		MarkBreakingChanges:     true,
		IncludeCompareLinks:     true,
		MaxTier:                 changelog.TierOptional,
	}
}

// CoreOptions returns options for KACL-compliant core output.
func CoreOptions() Options {
	return Options{
		IncludeReferences:       true,
		IncludeCommits:          false,
		LinkReferences:          false,
		IncludeAuthors:          true,
		IncludeSecurityMetadata: true,
		MarkBreakingChanges:     true,
		IncludeCompareLinks:     true,
		MaxTier:                 changelog.TierCore,
	}
}

// StandardOptions returns options including standard tier types.
func StandardOptions() Options {
	return Options{
		IncludeReferences:       true,
		IncludeCommits:          false,
		LinkReferences:          false,
		IncludeAuthors:          true,
		IncludeSecurityMetadata: true,
		MarkBreakingChanges:     true,
		IncludeCompareLinks:     true,
		MaxTier:                 changelog.TierStandard,
	}
}

// WithMaxTier returns a copy of the options with the MaxTier field set.
func (o Options) WithMaxTier(tier changelog.Tier) Options {
	o.MaxTier = tier
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
	Preset  string // default, minimal, full, core, standard
	MaxTier string // optional tier override
}

// OptionsFromConfig creates Options from a Config struct.
// It first applies the preset, then overrides MaxTier if specified.
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

	return opts, nil
}
