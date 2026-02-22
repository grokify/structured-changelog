package renderer

import (
	"errors"
	"testing"

	"github.com/grokify/structured-changelog/changelog"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if !opts.IncludeReferences {
		t.Error("expected IncludeReferences to be true")
	}
	if !opts.IncludeCommits {
		t.Error("expected IncludeCommits to be true")
	}
	if !opts.LinkReferences {
		t.Error("expected LinkReferences to be true")
	}
	if !opts.IncludeSecurityMetadata {
		t.Error("expected IncludeSecurityMetadata to be true")
	}
	if opts.MaxTier != changelog.TierOptional {
		t.Errorf("expected MaxTier to be optional, got %s", opts.MaxTier)
	}
}

func TestMinimalOptions(t *testing.T) {
	opts := MinimalOptions()

	if opts.IncludeReferences {
		t.Error("expected IncludeReferences to be false")
	}
	if opts.IncludeSecurityMetadata {
		t.Error("expected IncludeSecurityMetadata to be false")
	}
	if opts.MaxTier != changelog.TierCore {
		t.Errorf("expected MaxTier to be core, got %s", opts.MaxTier)
	}
}

func TestFullOptions(t *testing.T) {
	opts := FullOptions()

	if !opts.IncludeReferences {
		t.Error("expected IncludeReferences to be true")
	}
	if !opts.IncludeCommits {
		t.Error("expected IncludeCommits to be true")
	}
	if opts.MaxTier != changelog.TierOptional {
		t.Errorf("expected MaxTier to be optional, got %s", opts.MaxTier)
	}
}

func TestCoreOptions(t *testing.T) {
	opts := CoreOptions()

	if opts.MaxTier != changelog.TierCore {
		t.Errorf("expected MaxTier to be core, got %s", opts.MaxTier)
	}
}

func TestStandardOptions(t *testing.T) {
	opts := StandardOptions()

	if opts.MaxTier != changelog.TierStandard {
		t.Errorf("expected MaxTier to be standard, got %s", opts.MaxTier)
	}
}

func TestWithMaxTier(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxTier != changelog.TierOptional {
		t.Errorf("expected initial MaxTier to be optional, got %s", opts.MaxTier)
	}

	newOpts := opts.WithMaxTier(changelog.TierCore)
	if newOpts.MaxTier != changelog.TierCore {
		t.Errorf("expected new MaxTier to be core, got %s", newOpts.MaxTier)
	}

	// Original should be unchanged
	if opts.MaxTier != changelog.TierOptional {
		t.Errorf("expected original MaxTier to still be optional, got %s", opts.MaxTier)
	}
}

func TestOptionsFromPreset_Valid(t *testing.T) {
	tests := []struct {
		preset       string
		expectedTier changelog.Tier
	}{
		{"", changelog.TierOptional}, // default
		{"default", changelog.TierOptional},
		{"minimal", changelog.TierCore},
		{"full", changelog.TierOptional},
		{"core", changelog.TierCore},
		{"standard", changelog.TierStandard},
	}

	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			opts, err := OptionsFromPreset(tt.preset)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if opts.MaxTier != tt.expectedTier {
				t.Errorf("expected MaxTier %s, got %s", tt.expectedTier, opts.MaxTier)
			}
		})
	}
}

func TestOptionsFromPreset_Invalid(t *testing.T) {
	_, err := OptionsFromPreset("invalid")
	if err == nil {
		t.Error("expected error for invalid preset")
	}
	if !errors.Is(err, ErrInvalidPreset) {
		t.Errorf("expected ErrInvalidPreset, got %v", err)
	}
}

func TestOptionsFromConfig_PresetOnly(t *testing.T) {
	cfg := Config{
		Preset: "minimal",
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if opts.MaxTier != changelog.TierCore {
		t.Errorf("expected MaxTier core, got %s", opts.MaxTier)
	}
}

func TestOptionsFromConfig_WithTierOverride(t *testing.T) {
	cfg := Config{
		Preset:  "default",
		MaxTier: "core",
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if opts.MaxTier != changelog.TierCore {
		t.Errorf("expected MaxTier core, got %s", opts.MaxTier)
	}
}

func TestOptionsFromConfig_InvalidPreset(t *testing.T) {
	cfg := Config{
		Preset: "invalid",
	}

	_, err := OptionsFromConfig(cfg)
	if err == nil {
		t.Error("expected error for invalid preset")
	}
}

func TestOptionsFromConfig_InvalidTier(t *testing.T) {
	cfg := Config{
		Preset:  "default",
		MaxTier: "invalid",
	}

	_, err := OptionsFromConfig(cfg)
	if err == nil {
		t.Error("expected error for invalid tier")
	}
}

func TestOptionsFromConfig_CaseInsensitiveTier(t *testing.T) {
	cfg := Config{
		Preset:  "default",
		MaxTier: "CORE",
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if opts.MaxTier != changelog.TierCore {
		t.Errorf("expected MaxTier core, got %s", opts.MaxTier)
	}
}

func TestWithNotableOnly(t *testing.T) {
	opts := DefaultOptions()
	// Default is now NotableOnly=true
	if !opts.NotableOnly {
		t.Error("expected NotableOnly to be true by default")
	}

	newOpts := opts.WithNotableOnly(false)
	if newOpts.NotableOnly {
		t.Error("expected NotableOnly to be false after WithNotableOnly(false)")
	}

	// Original should be unchanged
	if !opts.NotableOnly {
		t.Error("expected original NotableOnly to still be true")
	}
}

func TestWithNotabilityPolicy(t *testing.T) {
	opts := DefaultOptions()
	// Default now has NotabilityPolicy set
	if opts.NotabilityPolicy == nil {
		t.Error("expected NotabilityPolicy to be set by default")
	}

	customPolicy := changelog.NewNotabilityPolicy([]string{"Security"})
	newOpts := opts.WithNotabilityPolicy(customPolicy)
	if newOpts.NotabilityPolicy == nil {
		t.Error("expected NotabilityPolicy to be set")
	}
	if newOpts.NotabilityPolicy.IsNotable("Added") {
		t.Error("expected Added to NOT be notable in custom policy")
	}
}

func TestOptionsFromConfig_DefaultNotableOnly(t *testing.T) {
	cfg := Config{
		Preset: "default",
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Default is now NotableOnly=true
	if !opts.NotableOnly {
		t.Error("expected NotableOnly to be true by default")
	}
	if opts.NotabilityPolicy == nil {
		t.Error("expected NotabilityPolicy to be set by default")
	}
}

func TestOptionsFromConfig_AllReleases(t *testing.T) {
	cfg := Config{
		Preset:      "default",
		AllReleases: true,
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if opts.NotableOnly {
		t.Error("expected NotableOnly to be false when AllReleases is true")
	}
	if opts.NotabilityPolicy != nil {
		t.Error("expected NotabilityPolicy to be nil when AllReleases is true")
	}
}

func TestOptionsFromConfig_CustomNotableCategories(t *testing.T) {
	cfg := Config{
		Preset:            "default",
		NotableCategories: []string{"Security", "Added"},
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !opts.NotableOnly {
		t.Error("expected NotableOnly to be true")
	}
	if opts.NotabilityPolicy == nil {
		t.Error("expected NotabilityPolicy to be set")
	}

	// Verify custom categories are used
	if !opts.NotabilityPolicy.IsNotable("Security") {
		t.Error("expected Security to be notable")
	}
	if !opts.NotabilityPolicy.IsNotable("Added") {
		t.Error("expected Added to be notable")
	}
	if opts.NotabilityPolicy.IsNotable("Fixed") {
		t.Error("expected Fixed to NOT be notable (not in custom list)")
	}
}

func TestOptionsFromConfig_FullPresetIncludesAllReleases(t *testing.T) {
	cfg := Config{
		Preset: "full",
	}

	opts, err := OptionsFromConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if opts.NotableOnly {
		t.Error("expected NotableOnly to be false for full preset")
	}
}
