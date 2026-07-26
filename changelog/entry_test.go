package changelog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewEntry(t *testing.T) {
	e := NewEntry("Test description")
	if e.Description != "Test description" {
		t.Errorf("expected description 'Test description', got %q", e.Description)
	}
}

func TestEntryWithIssue(t *testing.T) {
	e := NewEntry("Fix bug").WithIssue("#123")
	if e.Issue != "#123" {
		t.Errorf("expected issue '#123', got %q", e.Issue)
	}
	if e.Description != "Fix bug" {
		t.Errorf("expected description preserved, got %q", e.Description)
	}
}

func TestEntryWithPR(t *testing.T) {
	e := NewEntry("Add feature").WithPR("#456")
	if e.PR != "#456" {
		t.Errorf("expected PR '#456', got %q", e.PR)
	}
}

func TestEntryWithCommit(t *testing.T) {
	e := NewEntry("Refactor").WithCommit("abc123")
	if e.Commit != "abc123" {
		t.Errorf("expected commit 'abc123', got %q", e.Commit)
	}
}

func TestEntryWithAuthor(t *testing.T) {
	e := NewEntry("Update docs").WithAuthor("@alice")
	if e.Author != "@alice" {
		t.Errorf("expected author '@alice', got %q", e.Author)
	}
}

func TestEntryWithBreaking(t *testing.T) {
	e := NewEntry("API change").WithBreaking()
	if !e.Breaking {
		t.Error("expected Breaking to be true")
	}
}

func TestEntryWithCVE(t *testing.T) {
	e := NewEntry("Security fix").WithCVE("CVE-2026-12345")
	if e.CVE != "CVE-2026-12345" {
		t.Errorf("expected CVE 'CVE-2026-12345', got %q", e.CVE)
	}
}

func TestEntryWithGHSA(t *testing.T) {
	e := NewEntry("Security fix").WithGHSA("GHSA-xxxx-xxxx-xxxx")
	if e.GHSA != "GHSA-xxxx-xxxx-xxxx" {
		t.Errorf("expected GHSA 'GHSA-xxxx-xxxx-xxxx', got %q", e.GHSA)
	}
}

func TestEntryWithSeverity(t *testing.T) {
	e := NewEntry("Vulnerability").WithSeverity("high")
	if e.Severity != "high" {
		t.Errorf("expected severity 'high', got %q", e.Severity)
	}
}

func TestEntryWithCVSS(t *testing.T) {
	e := NewEntry("Vulnerability").WithCVSS(7.5, "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N")
	if e.CVSSScore != 7.5 {
		t.Errorf("expected CVSS score 7.5, got %f", e.CVSSScore)
	}
	if e.CVSSVector != "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N" {
		t.Errorf("expected CVSS vector, got %q", e.CVSSVector)
	}
}

func TestEntryWithCWE(t *testing.T) {
	e := NewEntry("SQL Injection").WithCWE("CWE-89")
	if e.CWE != "CWE-89" {
		t.Errorf("expected CWE 'CWE-89', got %q", e.CWE)
	}
}

func TestEntryWithComponent(t *testing.T) {
	e := NewEntry("Update dependency").WithComponent("redis", "7.0.0", "BSD-3-Clause")
	if e.Component != "redis" {
		t.Errorf("expected component 'redis', got %q", e.Component)
	}
	if e.ComponentVersion != "7.0.0" {
		t.Errorf("expected version '7.0.0', got %q", e.ComponentVersion)
	}
	if e.License != "BSD-3-Clause" {
		t.Errorf("expected license 'BSD-3-Clause', got %q", e.License)
	}
}

func TestEntryIsSecurityEntry(t *testing.T) {
	tests := []struct {
		name     string
		entry    Entry
		expected bool
	}{
		{"empty", Entry{Description: "Test"}, false},
		{"with CVE", Entry{Description: "Test", CVE: "CVE-2026-12345"}, true},
		{"with GHSA", Entry{Description: "Test", GHSA: "GHSA-xxxx-xxxx-xxxx"}, true},
		{"with severity", Entry{Description: "Test", Severity: "high"}, true},
		{"with all", Entry{Description: "Test", CVE: "CVE-2026-12345", GHSA: "GHSA-xxxx-xxxx-xxxx", Severity: "critical"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.IsSecurityEntry(); got != tt.expected {
				t.Errorf("IsSecurityEntry() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEntryWithRMI(t *testing.T) {
	e := NewEntry("Add widget support").WithRMI("RMI-MYREPO-042")
	if e.RMI != "RMI-MYREPO-042" {
		t.Errorf("expected RMI 'RMI-MYREPO-042', got %q", e.RMI)
	}
	if e.Description != "Add widget support" {
		t.Errorf("expected description preserved, got %q", e.Description)
	}
}

func TestEntryWithInitiative(t *testing.T) {
	e := NewEntry("Phase 2 work").WithInitiative("INIT-PRISMCONTROL-001")
	if e.Initiative != "INIT-PRISMCONTROL-001" {
		t.Errorf("expected initiative 'INIT-PRISMCONTROL-001', got %q", e.Initiative)
	}
	if e.Description != "Phase 2 work" {
		t.Errorf("expected description preserved, got %q", e.Description)
	}
}

func TestEntryRMIJSONRoundTrip(t *testing.T) {
	original := NewEntry("Add streaming support").
		WithRMI("RMI-MYREPO-042").
		WithInitiative("INIT-X-001").
		WithCommit("abc123").
		WithAuthor("@dev")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Entry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Description != original.Description {
		t.Errorf("description: got %q, want %q", decoded.Description, original.Description)
	}
	if decoded.RMI != original.RMI {
		t.Errorf("rmi: got %q, want %q", decoded.RMI, original.RMI)
	}
	if decoded.Initiative != original.Initiative {
		t.Errorf("initiative: got %q, want %q", decoded.Initiative, original.Initiative)
	}
	if decoded.Commit != original.Commit {
		t.Errorf("commit: got %q, want %q", decoded.Commit, original.Commit)
	}
	if decoded.Author != original.Author {
		t.Errorf("author: got %q, want %q", decoded.Author, original.Author)
	}

	// Verify omitempty: empty RMI/Initiative fields should not appear in JSON
	minimal := NewEntry("No PRISM fields")
	minData, err := json.Marshal(minimal)
	if err != nil {
		t.Fatal(err)
	}
	minJSON := string(minData)
	if strings.Contains(minJSON, `"rmi"`) {
		t.Error("empty RMI should be omitted from JSON")
	}
	if strings.Contains(minJSON, `"initiative"`) {
		t.Error("empty Initiative should be omitted from JSON")
	}
}

func TestEntryChaining(t *testing.T) {
	e := NewEntry("Complex entry").
		WithIssue("#100").
		WithPR("#101").
		WithCommit("abc123").
		WithAuthor("@dev").
		WithBreaking()

	if e.Description != "Complex entry" {
		t.Errorf("description not preserved through chaining")
	}
	if e.Issue != "#100" {
		t.Errorf("issue not set correctly")
	}
	if e.PR != "#101" {
		t.Errorf("PR not set correctly")
	}
	if e.Commit != "abc123" {
		t.Errorf("commit not set correctly")
	}
	if e.Author != "@dev" {
		t.Errorf("author not set correctly")
	}
	if !e.Breaking {
		t.Errorf("breaking not set correctly")
	}
}
