package aggregate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewResolver(t *testing.T) {
	r := NewResolver()

	if len(r.SearchPaths) == 0 {
		t.Error("expected default search paths, got none")
	}
}

func TestResolverResolve(t *testing.T) {
	// Create a temp directory with a CHANGELOG.json
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "github.com", "test", "repo")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	changelogPath := filepath.Join(projectDir, "CHANGELOG.json")
	if err := os.WriteFile(changelogPath, []byte(`{"ir_version":"1.0","project":"test"}`), 0600); err != nil {
		t.Fatalf("failed to create changelog: %v", err)
	}

	// Create resolver with temp directory as search path
	r := NewResolverWithPaths([]string{dir})

	refs := []ProjectRef{
		{Path: "github.com/test/repo"},
		{Path: "github.com/other/notfound"},
	}

	resolved, err := r.Resolve(refs)
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}

	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved projects, got %d", len(resolved))
	}

	// First should be found
	if !resolved[0].IsLocal {
		t.Error("expected first project to be local")
	}
	if resolved[0].ChangelogPath != changelogPath {
		t.Errorf("expected changelog path %q, got %q", changelogPath, resolved[0].ChangelogPath)
	}

	// Second should need approval
	if resolved[1].IsLocal {
		t.Error("expected second project to not be local")
	}
	if !resolved[1].NeedsApproval {
		t.Error("expected second project to need approval")
	}
}

func TestResolverWithLocalPathOverride(t *testing.T) {
	// Create a temp directory with a CHANGELOG.json
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "CHANGELOG.json"), []byte(`{"ir_version":"1.0","project":"test"}`), 0600); err != nil {
		t.Fatalf("failed to create changelog: %v", err)
	}

	r := NewResolver()

	refs := []ProjectRef{
		{Path: "github.com/some/repo", LocalPath: dir},
	}

	resolved, err := r.Resolve(refs)
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}

	if !resolved[0].IsLocal {
		t.Error("expected project with LocalPath override to be local")
	}
}

func TestResolverSummary(t *testing.T) {
	r := NewResolver()

	resolved := []ResolvedProject{
		{Ref: ProjectRef{Path: "github.com/a/b"}, IsLocal: true},
		{Ref: ProjectRef{Path: "github.com/c/d"}, IsLocal: true},
		{Ref: ProjectRef{Path: "github.com/e/f"}, NeedsApproval: true},
	}

	summary := r.Summary(resolved)

	if summary.TotalCount != 3 {
		t.Errorf("expected total 3, got %d", summary.TotalCount)
	}
	if summary.LocalCount != 2 {
		t.Errorf("expected local 2, got %d", summary.LocalCount)
	}
	if summary.RemoteCount != 1 {
		t.Errorf("expected remote 1, got %d", summary.RemoteCount)
	}
}

func TestFilterLocal(t *testing.T) {
	resolved := []ResolvedProject{
		{Ref: ProjectRef{Path: "a"}, IsLocal: true},
		{Ref: ProjectRef{Path: "b"}, NeedsApproval: true},
		{Ref: ProjectRef{Path: "c"}, IsLocal: true},
	}

	local := FilterLocal(resolved)

	if len(local) != 2 {
		t.Errorf("expected 2 local, got %d", len(local))
	}
}

func TestFilterRemote(t *testing.T) {
	resolved := []ResolvedProject{
		{Ref: ProjectRef{Path: "a"}, IsLocal: true},
		{Ref: ProjectRef{Path: "b"}, NeedsApproval: true},
		{Ref: ProjectRef{Path: "c"}, NeedsApproval: true},
	}

	remote := FilterRemote(resolved)

	if len(remote) != 2 {
		t.Errorf("expected 2 remote, got %d", len(remote))
	}
}

func TestParseProjectPath(t *testing.T) {
	tests := []struct {
		path    string
		host    string
		owner   string
		repo    string
		subpath string
		wantErr bool
	}{
		{
			path:  "github.com/org/repo",
			host:  "github.com",
			owner: "org",
			repo:  "repo",
		},
		{
			path:    "github.com/org/repo/sdk/go",
			host:    "github.com",
			owner:   "org",
			repo:    "repo",
			subpath: "sdk/go",
		},
		{
			path:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			host, owner, repo, subpath, err := ParseProjectPath(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if host != tt.host {
				t.Errorf("host: expected %q, got %q", tt.host, host)
			}
			if owner != tt.owner {
				t.Errorf("owner: expected %q, got %q", tt.owner, owner)
			}
			if repo != tt.repo {
				t.Errorf("repo: expected %q, got %q", tt.repo, repo)
			}
			if subpath != tt.subpath {
				t.Errorf("subpath: expected %q, got %q", tt.subpath, subpath)
			}
		})
	}
}

func TestGitHubRepoURL(t *testing.T) {
	tests := []struct {
		path    string
		want    string
		wantErr bool
	}{
		{
			path: "github.com/org/repo",
			want: "https://github.com/org/repo",
		},
		{
			path: "github.com/org/repo/subdir",
			want: "https://github.com/org/repo",
		},
		{
			path:    "gitlab.com/org/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := GitHubRepoURL(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
