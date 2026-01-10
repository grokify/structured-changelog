package changelog

// Release represents a single release in the changelog.
type Release struct {
	Version    string `json:"version,omitempty"`
	Date       string `json:"date,omitempty"`
	Yanked     bool   `json:"yanked,omitempty"`
	CompareURL string `json:"compare_url,omitempty"`

	// Overview & Critical (standard tier, except Security which is core)
	Highlights   []Entry `json:"highlights,omitempty"`
	Breaking     []Entry `json:"breaking,omitempty"`
	UpgradeGuide []Entry `json:"upgrade_guide,omitempty"`
	Security     []Entry `json:"security,omitempty"`

	// Core KACL types (core tier)
	Added      []Entry `json:"added,omitempty"`
	Changed    []Entry `json:"changed,omitempty"`
	Deprecated []Entry `json:"deprecated,omitempty"`
	Removed    []Entry `json:"removed,omitempty"`
	Fixed      []Entry `json:"fixed,omitempty"`

	// Quality (standard tier)
	Performance  []Entry `json:"performance,omitempty"`
	Dependencies []Entry `json:"dependencies,omitempty"`

	// Development (extended tier)
	Documentation []Entry `json:"documentation,omitempty"`
	Build         []Entry `json:"build,omitempty"`
	Tests         []Entry `json:"tests,omitempty"`

	// Operations (optional tier)
	Infrastructure []Entry `json:"infrastructure,omitempty"`
	Observability  []Entry `json:"observability,omitempty"`
	Compliance     []Entry `json:"compliance,omitempty"`

	// Internal (optional tier)
	Internal []Entry `json:"internal,omitempty"`

	// End Matter (extended tier)
	KnownIssues  []Entry `json:"known_issues,omitempty"`
	Contributors []Entry `json:"contributors,omitempty"`
}

// NewRelease creates a new release with the given version and date.
func NewRelease(version, date string) Release {
	return Release{
		Version: version,
		Date:    date,
	}
}

// IsEmpty returns true if the release has no entries.
func (r *Release) IsEmpty() bool {
	return len(r.Highlights) == 0 &&
		len(r.Breaking) == 0 &&
		len(r.UpgradeGuide) == 0 &&
		len(r.Security) == 0 &&
		len(r.Added) == 0 &&
		len(r.Changed) == 0 &&
		len(r.Deprecated) == 0 &&
		len(r.Removed) == 0 &&
		len(r.Fixed) == 0 &&
		len(r.Performance) == 0 &&
		len(r.Dependencies) == 0 &&
		len(r.Documentation) == 0 &&
		len(r.Build) == 0 &&
		len(r.Tests) == 0 &&
		len(r.Infrastructure) == 0 &&
		len(r.Observability) == 0 &&
		len(r.Compliance) == 0 &&
		len(r.Internal) == 0 &&
		len(r.KnownIssues) == 0 &&
		len(r.Contributors) == 0
}

// IsMaintenanceOnly returns true if the release contains only maintenance-type
// changes (dependencies, documentation, build, tests, internal) and no
// user-facing changes (added, changed, fixed, removed, security, etc.).
func (r *Release) IsMaintenanceOnly() bool {
	// Must have at least one entry to be considered maintenance
	if r.IsEmpty() {
		return false
	}

	// User-facing categories - if any have entries, not maintenance-only
	hasUserFacing := len(r.Highlights) > 0 ||
		len(r.Breaking) > 0 ||
		len(r.UpgradeGuide) > 0 ||
		len(r.Security) > 0 ||
		len(r.Added) > 0 ||
		len(r.Changed) > 0 ||
		len(r.Deprecated) > 0 ||
		len(r.Removed) > 0 ||
		len(r.Fixed) > 0 ||
		len(r.Performance) > 0 ||
		len(r.KnownIssues) > 0

	return !hasUserFacing
}

// Categories returns all non-empty categories in canonical order.
func (r *Release) Categories() []Category {
	return r.CategoriesFiltered(TierOptional)
}

// CategoriesFiltered returns non-empty categories up to the specified tier.
func (r *Release) CategoriesFiltered(maxTier Tier) []Category {
	var cats []Category

	// Canonical order matching CHANGE_TYPES.json
	categoryMap := r.categoryMap()
	for _, name := range DefaultRegistry.NamesUpToTier(maxTier) {
		if entries, ok := categoryMap[name]; ok && len(entries) > 0 {
			cats = append(cats, Category{Name: name, Entries: entries})
		}
	}
	return cats
}

// categoryMap returns a map of category name to entries.
func (r *Release) categoryMap() map[string][]Entry {
	return map[string][]Entry{
		"Highlights":     r.Highlights,
		"Breaking":       r.Breaking,
		"Upgrade Guide":  r.UpgradeGuide,
		"Security":       r.Security,
		"Added":          r.Added,
		"Changed":        r.Changed,
		"Deprecated":     r.Deprecated,
		"Removed":        r.Removed,
		"Fixed":          r.Fixed,
		"Performance":    r.Performance,
		"Dependencies":   r.Dependencies,
		"Documentation":  r.Documentation,
		"Build":          r.Build,
		"Tests":          r.Tests,
		"Infrastructure": r.Infrastructure,
		"Observability":  r.Observability,
		"Compliance":     r.Compliance,
		"Internal":       r.Internal,
		"Known Issues":   r.KnownIssues,
		"Contributors":   r.Contributors,
	}
}

// GetEntries returns entries for a category by name.
func (r *Release) GetEntries(categoryName string) []Entry {
	return r.categoryMap()[categoryName]
}

// Category represents a group of entries under a category heading.
type Category struct {
	Name    string
	Entries []Entry
}

// AddHighlights adds an entry to the Highlights category.
func (r *Release) AddHighlights(e Entry) {
	r.Highlights = append(r.Highlights, e)
}

// AddBreaking adds an entry to the Breaking category.
func (r *Release) AddBreaking(e Entry) {
	r.Breaking = append(r.Breaking, e)
}

// AddUpgradeGuide adds an entry to the Upgrade Guide category.
func (r *Release) AddUpgradeGuide(e Entry) {
	r.UpgradeGuide = append(r.UpgradeGuide, e)
}

// AddSecurity adds an entry to the Security category.
func (r *Release) AddSecurity(e Entry) {
	r.Security = append(r.Security, e)
}

// AddAdded adds an entry to the Added category.
func (r *Release) AddAdded(e Entry) {
	r.Added = append(r.Added, e)
}

// AddChanged adds an entry to the Changed category.
func (r *Release) AddChanged(e Entry) {
	r.Changed = append(r.Changed, e)
}

// AddDeprecated adds an entry to the Deprecated category.
func (r *Release) AddDeprecated(e Entry) {
	r.Deprecated = append(r.Deprecated, e)
}

// AddRemoved adds an entry to the Removed category.
func (r *Release) AddRemoved(e Entry) {
	r.Removed = append(r.Removed, e)
}

// AddFixed adds an entry to the Fixed category.
func (r *Release) AddFixed(e Entry) {
	r.Fixed = append(r.Fixed, e)
}

// AddPerformance adds an entry to the Performance category.
func (r *Release) AddPerformance(e Entry) {
	r.Performance = append(r.Performance, e)
}

// AddDependencies adds an entry to the Dependencies category.
func (r *Release) AddDependencies(e Entry) {
	r.Dependencies = append(r.Dependencies, e)
}

// AddDocumentation adds an entry to the Documentation category.
func (r *Release) AddDocumentation(e Entry) {
	r.Documentation = append(r.Documentation, e)
}

// AddBuild adds an entry to the Build category.
func (r *Release) AddBuild(e Entry) {
	r.Build = append(r.Build, e)
}

// AddTests adds an entry to the Tests category.
func (r *Release) AddTests(e Entry) {
	r.Tests = append(r.Tests, e)
}

// AddInfrastructure adds an entry to the Infrastructure category.
func (r *Release) AddInfrastructure(e Entry) {
	r.Infrastructure = append(r.Infrastructure, e)
}

// AddObservability adds an entry to the Observability category.
func (r *Release) AddObservability(e Entry) {
	r.Observability = append(r.Observability, e)
}

// AddCompliance adds an entry to the Compliance category.
func (r *Release) AddCompliance(e Entry) {
	r.Compliance = append(r.Compliance, e)
}

// AddInternal adds an entry to the Internal category.
func (r *Release) AddInternal(e Entry) {
	r.Internal = append(r.Internal, e)
}

// AddKnownIssues adds an entry to the Known Issues category.
func (r *Release) AddKnownIssues(e Entry) {
	r.KnownIssues = append(r.KnownIssues, e)
}

// AddContributors adds an entry to the Contributors category.
func (r *Release) AddContributors(e Entry) {
	r.Contributors = append(r.Contributors, e)
}
