package changelog

// Entry represents a single changelog entry.
type Entry struct {
	Description string `json:"description"`
	Issue       string `json:"issue,omitempty"`
	PR          string `json:"pr,omitempty"`
	Commit      string `json:"commit,omitempty"`
	Author      string `json:"author,omitempty"`
	Breaking    bool   `json:"breaking,omitempty"`

	// SBOM metadata
	Component        string `json:"component,omitempty"`
	ComponentVersion string `json:"componentVersion,omitempty"`
	License          string `json:"license,omitempty"`

	// Security metadata
	CVE              string  `json:"cve,omitempty"`
	GHSA             string  `json:"ghsa,omitempty"`
	Severity         string  `json:"severity,omitempty"`
	CVSSScore        float64 `json:"cvssScore,omitempty"`
	CVSSVector       string  `json:"cvssVector,omitempty"`
	CWE              string  `json:"cwe,omitempty"`
	AffectedVersions string  `json:"affectedVersions,omitempty"`
	PatchedVersions  string  `json:"patchedVersions,omitempty"`
	SARIFRuleID      string  `json:"sarifRuleId,omitempty"`
}

// NewEntry creates a new entry with the given description.
func NewEntry(description string) Entry {
	return Entry{Description: description}
}

// WithIssue sets the issue reference.
func (e Entry) WithIssue(issue string) Entry {
	e.Issue = issue
	return e
}

// WithPR sets the pull request reference.
func (e Entry) WithPR(pr string) Entry {
	e.PR = pr
	return e
}

// WithCommit sets the commit SHA.
func (e Entry) WithCommit(commit string) Entry {
	e.Commit = commit
	return e
}

// WithAuthor sets the author.
func (e Entry) WithAuthor(author string) Entry {
	e.Author = author
	return e
}

// WithBreaking marks the entry as a breaking change.
func (e Entry) WithBreaking() Entry {
	e.Breaking = true
	return e
}

// WithCVE sets CVE identifier for security entries.
func (e Entry) WithCVE(cve string) Entry {
	e.CVE = cve
	return e
}

// WithGHSA sets GitHub Security Advisory identifier.
func (e Entry) WithGHSA(ghsa string) Entry {
	e.GHSA = ghsa
	return e
}

// WithSeverity sets the severity level.
func (e Entry) WithSeverity(severity string) Entry {
	e.Severity = severity
	return e
}

// WithCVSS sets the CVSS score and vector.
func (e Entry) WithCVSS(score float64, vector string) Entry {
	e.CVSSScore = score
	e.CVSSVector = vector
	return e
}

// WithCWE sets the CWE identifier.
func (e Entry) WithCWE(cwe string) Entry {
	e.CWE = cwe
	return e
}

// WithComponent sets SBOM component information.
func (e Entry) WithComponent(name, version, license string) Entry {
	e.Component = name
	e.ComponentVersion = version
	e.License = license
	return e
}

// IsSecurityEntry returns true if the entry has security metadata.
func (e Entry) IsSecurityEntry() bool {
	return e.CVE != "" || e.GHSA != "" || e.Severity != ""
}
