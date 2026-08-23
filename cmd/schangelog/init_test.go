package main

import (
	"reflect"
	"testing"

	"github.com/grokify/structured-changelog/gitlog"
)

func TestBuildReleaseFromCommits_PopulatesRMIs(t *testing.T) {
	commits := []gitlog.Commit{
		{
			ShortHash:         "abc1234",
			Subject:           "add a feature",
			Type:              "feat",
			SuggestedCategory: "Added",
			// Unsorted on purpose; output must be sorted.
			RMIs: []string{"RMI-SCHANGELOG-102", "RMI-SCHANGELOG-101"},
		},
		{
			ShortHash:         "def5678",
			Subject:           "fix a bug",
			Type:              "fix",
			SuggestedCategory: "Fixed",
		},
	}

	rel := buildReleaseFromCommits("1.0.0", "2026-01-01", commits)

	if len(rel.Added) != 1 {
		t.Fatalf("expected 1 Added entry, got %d", len(rel.Added))
	}
	want := []string{"RMI-SCHANGELOG-101", "RMI-SCHANGELOG-102"}
	if !reflect.DeepEqual(rel.Added[0].RMIs, want) {
		t.Errorf("Added entry RMIs = %v, want sorted %v", rel.Added[0].RMIs, want)
	}

	if len(rel.Fixed) != 1 {
		t.Fatalf("expected 1 Fixed entry, got %d", len(rel.Fixed))
	}
	if rel.Fixed[0].RMIs != nil {
		t.Errorf("commit without trailers should have no RMIs, got %v", rel.Fixed[0].RMIs)
	}
}
