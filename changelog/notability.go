package changelog

// NotabilityPolicy defines which releases are considered notable.
// A release is notable if it contains at least one entry in any of the
// specified notable categories.
type NotabilityPolicy struct {
	// NotableCategories specifies which categories make a release notable.
	// If a release has entries in ANY of these categories, it is considered notable.
	// If empty and UseDefault is false, all releases are considered notable.
	NotableCategories []string
}

// DefaultNotableCategories returns the default list of categories that make a
// release notable. These are the inverse of maintenance-only categories.
//
// Notable categories (user-facing changes):
//   - Highlights, Breaking, Upgrade Guide, Security
//   - Added, Changed, Deprecated, Removed, Fixed
//   - Performance, Known Issues
//
// Non-notable (maintenance) categories:
//   - Dependencies, Documentation, Build, Tests
//   - Infrastructure, Observability, Compliance
//   - Internal, Contributors
func DefaultNotableCategories() []string {
	return []string{
		CategoryHighlights,
		CategoryBreaking,
		CategoryUpgradeGuide,
		CategorySecurity,
		CategoryAdded,
		CategoryChanged,
		CategoryDeprecated,
		CategoryRemoved,
		CategoryFixed,
		CategoryPerformance,
		CategoryKnownIssues,
	}
}

// DefaultNotabilityPolicy returns a policy using the default notable categories.
func DefaultNotabilityPolicy() *NotabilityPolicy {
	return &NotabilityPolicy{
		NotableCategories: DefaultNotableCategories(),
	}
}

// NewNotabilityPolicy creates a policy with custom notable categories.
func NewNotabilityPolicy(categories []string) *NotabilityPolicy {
	return &NotabilityPolicy{
		NotableCategories: categories,
	}
}

// IsNotable returns true if the given category is considered notable by this policy.
func (p *NotabilityPolicy) IsNotable(categoryName string) bool {
	if p == nil || len(p.NotableCategories) == 0 {
		return true // No policy or empty = all notable
	}
	for _, cat := range p.NotableCategories {
		if cat == categoryName {
			return true
		}
	}
	return false
}
