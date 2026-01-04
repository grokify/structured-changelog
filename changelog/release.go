package changelog

// Release represents a single release in the changelog.
type Release struct {
	Version    string  `json:"version,omitempty"`
	Date       string  `json:"date,omitempty"`
	Yanked     bool    `json:"yanked,omitempty"`
	CompareURL string  `json:"compare_url,omitempty"`
	Added      []Entry `json:"added,omitempty"`
	Changed    []Entry `json:"changed,omitempty"`
	Deprecated []Entry `json:"deprecated,omitempty"`
	Removed    []Entry `json:"removed,omitempty"`
	Fixed      []Entry `json:"fixed,omitempty"`
	Security   []Entry `json:"security,omitempty"`
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
	return len(r.Added) == 0 &&
		len(r.Changed) == 0 &&
		len(r.Deprecated) == 0 &&
		len(r.Removed) == 0 &&
		len(r.Fixed) == 0 &&
		len(r.Security) == 0
}

// Categories returns all non-empty categories in standard order.
func (r *Release) Categories() []Category {
	var cats []Category
	if len(r.Added) > 0 {
		cats = append(cats, Category{Name: "Added", Entries: r.Added})
	}
	if len(r.Changed) > 0 {
		cats = append(cats, Category{Name: "Changed", Entries: r.Changed})
	}
	if len(r.Deprecated) > 0 {
		cats = append(cats, Category{Name: "Deprecated", Entries: r.Deprecated})
	}
	if len(r.Removed) > 0 {
		cats = append(cats, Category{Name: "Removed", Entries: r.Removed})
	}
	if len(r.Fixed) > 0 {
		cats = append(cats, Category{Name: "Fixed", Entries: r.Fixed})
	}
	if len(r.Security) > 0 {
		cats = append(cats, Category{Name: "Security", Entries: r.Security})
	}
	return cats
}

// Category represents a group of entries under a category heading.
type Category struct {
	Name    string
	Entries []Entry
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

// AddSecurity adds an entry to the Security category.
func (r *Release) AddSecurity(e Entry) {
	r.Security = append(r.Security, e)
}
