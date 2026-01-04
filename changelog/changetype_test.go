package changelog

import (
	"testing"
)

func TestTierPriority(t *testing.T) {
	tests := []struct {
		tier     Tier
		priority int
	}{
		{TierCore, 0},
		{TierStandard, 1},
		{TierExtended, 2},
		{TierOptional, 3},
		{Tier("unknown"), 4}, // Unknown tiers get lowest priority
	}

	for _, tt := range tests {
		t.Run(string(tt.tier), func(t *testing.T) {
			if got := tt.tier.Priority(); got != tt.priority {
				t.Errorf("Tier(%q).Priority() = %d, want %d", tt.tier, got, tt.priority)
			}
		})
	}
}

func TestTierIsValid(t *testing.T) {
	tests := []struct {
		tier  Tier
		valid bool
	}{
		{TierCore, true},
		{TierStandard, true},
		{TierExtended, true},
		{TierOptional, true},
		{Tier("unknown"), false},
		{Tier(""), false},
		{Tier("CORE"), false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(string(tt.tier), func(t *testing.T) {
			if got := tt.tier.IsValid(); got != tt.valid {
				t.Errorf("Tier(%q).IsValid() = %v, want %v", tt.tier, got, tt.valid)
			}
		})
	}
}

func TestTierIncludesOrHigher(t *testing.T) {
	tests := []struct {
		tier    Tier
		maxTier Tier
		want    bool
	}{
		// Core tier includes only core
		{TierCore, TierCore, true},
		{TierStandard, TierCore, false},
		{TierExtended, TierCore, false},
		{TierOptional, TierCore, false},

		// Standard tier includes core and standard
		{TierCore, TierStandard, true},
		{TierStandard, TierStandard, true},
		{TierExtended, TierStandard, false},
		{TierOptional, TierStandard, false},

		// Extended tier includes core, standard, extended
		{TierCore, TierExtended, true},
		{TierStandard, TierExtended, true},
		{TierExtended, TierExtended, true},
		{TierOptional, TierExtended, false},

		// Optional tier includes all
		{TierCore, TierOptional, true},
		{TierStandard, TierOptional, true},
		{TierExtended, TierOptional, true},
		{TierOptional, TierOptional, true},
	}

	for _, tt := range tests {
		name := string(tt.tier) + "_in_" + string(tt.maxTier)
		t.Run(name, func(t *testing.T) {
			if got := tt.tier.IncludesOrHigher(tt.maxTier); got != tt.want {
				t.Errorf("Tier(%q).IncludesOrHigher(%q) = %v, want %v", tt.tier, tt.maxTier, got, tt.want)
			}
		})
	}
}

func TestDefaultRegistry(t *testing.T) {
	if DefaultRegistry == nil {
		t.Fatal("DefaultRegistry is nil")
	}

	// Should have 20 change types
	all := DefaultRegistry.All()
	if len(all) != 20 {
		t.Errorf("DefaultRegistry.All() returned %d types, want 20", len(all))
	}

	// Names should match count
	names := DefaultRegistry.Names()
	if len(names) != 20 {
		t.Errorf("DefaultRegistry.Names() returned %d names, want 20", len(names))
	}
}

func TestRegistryGet(t *testing.T) {
	tests := []struct {
		name     string
		wantTier Tier
		wantNil  bool
	}{
		{"Added", TierCore, false},
		{"Security", TierCore, false},
		{"Breaking", TierStandard, false},
		{"Documentation", TierExtended, false},
		{"Internal", TierOptional, false},
		{"NonExistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := DefaultRegistry.Get(tt.name)
			if tt.wantNil {
				if ct != nil {
					t.Errorf("Get(%q) = %v, want nil", tt.name, ct)
				}
				return
			}
			if ct == nil {
				t.Fatalf("Get(%q) = nil, want non-nil", tt.name)
			}
			if ct.Tier != tt.wantTier {
				t.Errorf("Get(%q).Tier = %q, want %q", tt.name, ct.Tier, tt.wantTier)
			}
		})
	}
}

func TestRegistryByTier(t *testing.T) {
	tests := []struct {
		tier      Tier
		wantCount int
	}{
		{TierCore, 6},     // Security, Added, Changed, Deprecated, Removed, Fixed
		{TierStandard, 5}, // Highlights, Breaking, Upgrade Guide, Performance, Dependencies
		{TierExtended, 5}, // Documentation, Build, Tests, Known Issues, Contributors
		{TierOptional, 4}, // Infrastructure, Observability, Compliance, Internal
	}

	for _, tt := range tests {
		t.Run(string(tt.tier), func(t *testing.T) {
			types := DefaultRegistry.ByTier(tt.tier)
			if len(types) != tt.wantCount {
				t.Errorf("ByTier(%q) returned %d types, want %d", tt.tier, len(types), tt.wantCount)
			}
		})
	}
}

func TestRegistryFilterByMaxTier(t *testing.T) {
	tests := []struct {
		maxTier   Tier
		wantCount int
	}{
		{TierCore, 6},      // Only core types
		{TierStandard, 11}, // Core + standard
		{TierExtended, 16}, // Core + standard + extended
		{TierOptional, 20}, // All types
	}

	for _, tt := range tests {
		t.Run(string(tt.maxTier), func(t *testing.T) {
			types := DefaultRegistry.FilterByMaxTier(tt.maxTier)
			if len(types) != tt.wantCount {
				t.Errorf("FilterByMaxTier(%q) returned %d types, want %d", tt.maxTier, len(types), tt.wantCount)
			}
		})
	}
}

func TestRegistryNamesUpToTier(t *testing.T) {
	tests := []struct {
		maxTier   Tier
		wantCount int
	}{
		{TierCore, 6},
		{TierStandard, 11},
		{TierExtended, 16},
		{TierOptional, 20},
	}

	for _, tt := range tests {
		t.Run(string(tt.maxTier), func(t *testing.T) {
			names := DefaultRegistry.NamesUpToTier(tt.maxTier)
			if len(names) != tt.wantCount {
				t.Errorf("NamesUpToTier(%q) returned %d names, want %d", tt.maxTier, len(names), tt.wantCount)
			}
		})
	}
}

func TestRegistryIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"Added", true},
		{"Security", true},
		{"Breaking", true},
		{"Highlights", true},
		{"Unknown", false},
		{"added", false}, // Case sensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultRegistry.IsValidName(tt.name); got != tt.valid {
				t.Errorf("IsValidName(%q) = %v, want %v", tt.name, got, tt.valid)
			}
		})
	}
}

func TestRegistryCoreTypes(t *testing.T) {
	coreTypes := DefaultRegistry.CoreTypes()

	// Should have exactly 6 core types
	if len(coreTypes) != 6 {
		t.Errorf("CoreTypes() returned %d types, want 6", len(coreTypes))
	}

	// Check expected core types are present
	expected := map[string]bool{
		"Security":   true,
		"Added":      true,
		"Changed":    true,
		"Deprecated": true,
		"Removed":    true,
		"Fixed":      true,
	}

	for _, name := range coreTypes {
		if !expected[name] {
			t.Errorf("CoreTypes() contains unexpected type %q", name)
		}
	}
}

func TestCanonicalOrder(t *testing.T) {
	// Verify the canonical order matches CHANGE_TYPES.json
	expectedOrder := []string{
		"Highlights",
		"Breaking",
		"Upgrade Guide",
		"Security",
		"Added",
		"Changed",
		"Deprecated",
		"Removed",
		"Fixed",
		"Performance",
		"Dependencies",
		"Documentation",
		"Build",
		"Tests",
		"Infrastructure",
		"Observability",
		"Compliance",
		"Internal",
		"Known Issues",
		"Contributors",
	}

	names := DefaultRegistry.Names()

	if len(names) != len(expectedOrder) {
		t.Fatalf("Names() returned %d items, want %d", len(names), len(expectedOrder))
	}

	for i, want := range expectedOrder {
		if names[i] != want {
			t.Errorf("Names()[%d] = %q, want %q", i, names[i], want)
		}
	}
}

func TestTierDescriptions(t *testing.T) {
	// Ensure all tiers have descriptions
	for _, tier := range TierOrder {
		desc, ok := TierDescriptions[tier]
		if !ok {
			t.Errorf("TierDescriptions missing entry for tier %q", tier)
		}
		if desc == "" {
			t.Errorf("TierDescriptions[%q] is empty", tier)
		}
	}
}

func TestChangeTypeFields(t *testing.T) {
	// Verify each change type has required fields populated
	for _, ct := range DefaultRegistry.All() {
		if ct.Name == "" {
			t.Error("Found change type with empty Name")
		}
		if ct.Description == "" {
			t.Errorf("Change type %q has empty Description", ct.Name)
		}
		if ct.Subtitle == "" {
			t.Errorf("Change type %q has empty Subtitle", ct.Name)
		}
		if !ct.Tier.IsValid() {
			t.Errorf("Change type %q has invalid Tier %q", ct.Name, ct.Tier)
		}
	}
}
