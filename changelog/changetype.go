package changelog

import (
	"embed"
	"encoding/json"
	"fmt"
	"slices"
)

//go:embed change_types.json
var changeTypesFS embed.FS

// Tier represents the priority level of a change type.
type Tier string

// Tier constants in order of priority (highest to lowest).
const (
	TierCore     Tier = "core"
	TierStandard Tier = "standard"
	TierExtended Tier = "extended"
	TierOptional Tier = "optional"
)

// TierOrder defines the canonical ordering of tiers from highest to lowest priority.
var TierOrder = []Tier{TierCore, TierStandard, TierExtended, TierOptional}

// TierDescriptions provides human-readable descriptions for each tier.
var TierDescriptions = map[Tier]string{
	TierCore:     "Standard types defined by Keep a Changelog (KACL)",
	TierStandard: "Commonly used by major providers and popular open source projects",
	TierExtended: "Change metadata for documentation, build, and acknowledgments",
	TierOptional: "For deployment teams and internal operational visibility",
}

// Priority returns the numeric priority of a tier (lower is higher priority).
func (t Tier) Priority() int {
	for i, tier := range TierOrder {
		if tier == t {
			return i
		}
	}
	return len(TierOrder) // Unknown tiers have lowest priority
}

// IsValid returns true if the tier is a recognized value.
func (t Tier) IsValid() bool {
	return slices.Contains(TierOrder, t)
}

// IncludesOrHigher returns true if this tier should be included when filtering
// at the given maximum tier level.
func (t Tier) IncludesOrHigher(maxTier Tier) bool {
	return t.Priority() <= maxTier.Priority()
}

// ChangeType represents a single change type definition.
type ChangeType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subtitle    string `json:"subtitle"`
	Tier        Tier   `json:"tier"`
}

// ChangeTypeRegistry holds all change type definitions.
type ChangeTypeRegistry struct {
	types    []ChangeType
	byName   map[string]*ChangeType
	byTier   map[Tier][]ChangeType
	nameList []string
}

// DefaultRegistry is the global registry loaded from embedded change_types.json.
var DefaultRegistry *ChangeTypeRegistry

func init() {
	var err error
	DefaultRegistry, err = LoadEmbeddedChangeTypes()
	if err != nil {
		panic(fmt.Sprintf("failed to load embedded change types: %v", err))
	}
}

// LoadEmbeddedChangeTypes loads change types from the embedded JSON file.
func LoadEmbeddedChangeTypes() (*ChangeTypeRegistry, error) {
	data, err := changeTypesFS.ReadFile("change_types.json")
	if err != nil {
		return nil, fmt.Errorf("reading embedded change_types.json: %w", err)
	}
	return ParseChangeTypes(data)
}

// ParseChangeTypes parses change type definitions from JSON bytes.
func ParseChangeTypes(data []byte) (*ChangeTypeRegistry, error) {
	var types []ChangeType
	if err := json.Unmarshal(data, &types); err != nil {
		return nil, fmt.Errorf("parsing change types JSON: %w", err)
	}

	registry := &ChangeTypeRegistry{
		types:  types,
		byName: make(map[string]*ChangeType),
		byTier: make(map[Tier][]ChangeType),
	}

	for i := range types {
		ct := &types[i]
		registry.byName[ct.Name] = ct
		registry.byTier[ct.Tier] = append(registry.byTier[ct.Tier], *ct)
		registry.nameList = append(registry.nameList, ct.Name)
	}

	return registry, nil
}

// All returns all change types in canonical order.
func (r *ChangeTypeRegistry) All() []ChangeType {
	return r.types
}

// Names returns all change type names in canonical order.
func (r *ChangeTypeRegistry) Names() []string {
	return r.nameList
}

// Get returns a change type by name, or nil if not found.
func (r *ChangeTypeRegistry) Get(name string) *ChangeType {
	return r.byName[name]
}

// ByTier returns all change types for a given tier.
func (r *ChangeTypeRegistry) ByTier(tier Tier) []ChangeType {
	return r.byTier[tier]
}

// FilterByMaxTier returns change types up to and including the given tier.
func (r *ChangeTypeRegistry) FilterByMaxTier(maxTier Tier) []ChangeType {
	var result []ChangeType
	for _, ct := range r.types {
		if ct.Tier.IncludesOrHigher(maxTier) {
			result = append(result, ct)
		}
	}
	return result
}

// NamesUpToTier returns change type names up to and including the given tier.
func (r *ChangeTypeRegistry) NamesUpToTier(maxTier Tier) []string {
	var result []string
	for _, ct := range r.types {
		if ct.Tier.IncludesOrHigher(maxTier) {
			result = append(result, ct.Name)
		}
	}
	return result
}

// IsValidName returns true if the name is a recognized change type.
func (r *ChangeTypeRegistry) IsValidName(name string) bool {
	_, ok := r.byName[name]
	return ok
}

// CoreTypes returns the names of all core (KACL) change types.
func (r *ChangeTypeRegistry) CoreTypes() []string {
	var result []string
	for _, ct := range r.byTier[TierCore] {
		result = append(result, ct.Name)
	}
	return result
}
