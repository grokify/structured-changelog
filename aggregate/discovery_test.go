package aggregate

import (
	"testing"
)

func TestDefaultDiscoveryOptions(t *testing.T) {
	opts := DefaultDiscoveryOptions()

	if opts.IncludeArchived {
		t.Error("expected IncludeArchived=false by default")
	}

	if opts.IncludeForks {
		t.Error("expected IncludeForks=false by default")
	}

	if opts.MaxReposPerOrg != 0 {
		t.Errorf("expected MaxReposPerOrg=0 (unlimited), got %d", opts.MaxReposPerOrg)
	}
}

func TestDiscoveryOptionsCustom(t *testing.T) {
	opts := DiscoveryOptions{
		IncludeArchived: true,
		IncludeForks:    true,
		MaxReposPerOrg:  100,
	}

	if !opts.IncludeArchived {
		t.Error("expected IncludeArchived=true")
	}

	if !opts.IncludeForks {
		t.Error("expected IncludeForks=true")
	}

	if opts.MaxReposPerOrg != 100 {
		t.Errorf("expected MaxReposPerOrg=100, got %d", opts.MaxReposPerOrg)
	}
}

func TestDiscoveryResultStructure(t *testing.T) {
	result := DiscoveryResult{
		Projects: []ProjectRef{
			{Path: "github.com/org/repo1", Discovered: true},
			{Path: "github.com/org/repo2", Discovered: true},
		},
		Statistics: DiscoveryStats{
			SourcesScanned:     2,
			ReposScanned:       50,
			ReposWithChangelog: 10,
			ChangelogsFound:    12,
		},
	}

	if len(result.Projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(result.Projects))
	}

	if result.Statistics.SourcesScanned != 2 {
		t.Errorf("expected SourcesScanned=2, got %d", result.Statistics.SourcesScanned)
	}

	if result.Statistics.ReposScanned != 50 {
		t.Errorf("expected ReposScanned=50, got %d", result.Statistics.ReposScanned)
	}

	if result.Statistics.ReposWithChangelog != 10 {
		t.Errorf("expected ReposWithChangelog=10, got %d", result.Statistics.ReposWithChangelog)
	}

	if result.Statistics.ChangelogsFound != 12 {
		t.Errorf("expected ChangelogsFound=12, got %d", result.Statistics.ChangelogsFound)
	}
}

func TestSourceTypes(t *testing.T) {
	tests := []struct {
		sourceType string
		name       string
	}{
		{SourceTypeOrg, "grokify"},
		{SourceTypeUser, "johndoe"},
	}

	for _, tt := range tests {
		source := Source{
			Type: tt.sourceType,
			Name: tt.name,
		}

		if source.Type != tt.sourceType {
			t.Errorf("expected Type=%q, got %q", tt.sourceType, source.Type)
		}

		if source.Name != tt.name {
			t.Errorf("expected Name=%q, got %q", tt.name, source.Name)
		}
	}
}

func TestProjectRefDiscoveryFlag(t *testing.T) {
	// Manual project ref
	manual := ProjectRef{
		Path:       "github.com/org/manual-repo",
		Discovered: false,
	}

	if manual.Discovered {
		t.Error("expected manual ProjectRef to have Discovered=false")
	}

	// Discovered project ref
	discovered := ProjectRef{
		Path:       "github.com/org/discovered-repo",
		Discovered: true,
	}

	if !discovered.Discovered {
		t.Error("expected discovered ProjectRef to have Discovered=true")
	}
}

func TestDiscoveryStatsZeroValues(t *testing.T) {
	stats := DiscoveryStats{}

	if stats.SourcesScanned != 0 {
		t.Errorf("expected SourcesScanned=0, got %d", stats.SourcesScanned)
	}

	if stats.ReposScanned != 0 {
		t.Errorf("expected ReposScanned=0, got %d", stats.ReposScanned)
	}

	if stats.ReposWithChangelog != 0 {
		t.Errorf("expected ReposWithChangelog=0, got %d", stats.ReposWithChangelog)
	}

	if stats.ChangelogsFound != 0 {
		t.Errorf("expected ChangelogsFound=0, got %d", stats.ChangelogsFound)
	}
}
