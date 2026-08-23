package changelog

import "testing"

func TestEntryWithRMIs(t *testing.T) {
	e := NewEntry("add thing").WithRMIs("RMI-SCHANGELOG-101", "RMI-SCHANGELOG-102")
	if len(e.RMIs) != 2 || e.RMIs[0] != "RMI-SCHANGELOG-101" {
		t.Errorf("WithRMIs did not set RMIs: %v", e.RMIs)
	}
}

// validChangelogWithEntry builds a minimal valid changelog whose single release
// carries one Added entry, so ValidateRich runs without unrelated errors.
func validChangelogWithEntry(entry Entry) *Changelog {
	return &Changelog{
		IRVersion:        IRVersion,
		Project:          "test",
		Versioning:       VersioningSemVer,
		CommitConvention: CommitConventionConventional,
		Releases: []Release{
			{Version: "1.0.0", Date: "2026-01-01", Added: []Entry{entry}},
		},
	}
}

func countWarnings(res RichValidationResult, code ErrorCode) int {
	n := 0
	for _, w := range res.Warnings {
		if w.Code == code {
			n++
		}
	}
	return n
}

func TestValidateRich_MalformedRMI(t *testing.T) {
	entry := NewEntry("a sufficiently long description").WithRMIs("RMI-SCHANGELOG-101", "not-an-rmi")
	res := validChangelogWithEntry(entry).ValidateRich()

	if got := countWarnings(res, WarnCodeMalformedRMI); got != 1 {
		t.Errorf("expected 1 malformed-RMI warning, got %d (warnings: %+v)", got, res.Warnings)
	}
}

func TestValidateRich_WellFormedRMIsNoWarning(t *testing.T) {
	entry := NewEntry("a sufficiently long description").WithRMIs("RMI-SCHANGELOG-101", "RMI-FOO-9")
	res := validChangelogWithEntry(entry).ValidateRich()

	if got := countWarnings(res, WarnCodeMalformedRMI); got != 0 {
		t.Errorf("well-formed RMIs should not warn, got %d", got)
	}
	if got := countWarnings(res, WarnCodeRMIDisagreement); got != 0 {
		t.Errorf("no singular rmi set, should not report disagreement, got %d", got)
	}
}

func TestValidateRich_RMIDisagreement(t *testing.T) {
	// Singular rmi not present in rmis -> disagreement warning.
	entry := NewEntry("a sufficiently long description")
	entry.RMI = "RMI-SCHANGELOG-999"
	entry = entry.WithRMIs("RMI-SCHANGELOG-101")
	res := validChangelogWithEntry(entry).ValidateRich()

	if got := countWarnings(res, WarnCodeRMIDisagreement); got != 1 {
		t.Errorf("expected 1 disagreement warning, got %d", got)
	}
}

func TestValidateRich_RMIConsistentNoWarning(t *testing.T) {
	// Singular rmi present in rmis -> no disagreement.
	entry := NewEntry("a sufficiently long description")
	entry.RMI = "RMI-SCHANGELOG-101"
	entry = entry.WithRMIs("RMI-SCHANGELOG-101", "RMI-SCHANGELOG-102")
	res := validChangelogWithEntry(entry).ValidateRich()

	if got := countWarnings(res, WarnCodeRMIDisagreement); got != 0 {
		t.Errorf("consistent rmi/rmis should not warn, got %d", got)
	}
}
