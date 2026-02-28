package aggregate

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v84/github"
	"github.com/grokify/gogithub/auth"
	"github.com/grokify/gogithub/repo"
)

// DiscoveryClient scans GitHub orgs/users for repos with changelogs.
type DiscoveryClient struct {
	gh *github.Client
}

// NewDiscoveryClient creates a discovery client with the given token.
// If token is empty, checks GITHUB_TOKEN environment variable.
func NewDiscoveryClient(token string) (*DiscoveryClient, error) {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required for discovery")
	}

	client := auth.NewGitHubClient(context.Background(), token)

	return &DiscoveryClient{
		gh: client,
	}, nil
}

// DiscoverProjects scans sources for repos containing CHANGELOG.json.
// Returns discovered ProjectRefs not already in the manifest.
func (d *DiscoveryClient) DiscoverProjects(ctx context.Context, sources []Source) ([]ProjectRef, error) {
	var projects []ProjectRef

	for _, source := range sources {
		discovered, err := d.discoverSource(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("discovering %s %s: %w", source.Type, source.Name, err)
		}
		projects = append(projects, discovered...)
	}

	return projects, nil
}

func (d *DiscoveryClient) discoverSource(ctx context.Context, source Source) ([]ProjectRef, error) {
	var repos []*github.Repository
	var err error

	switch source.Type {
	case SourceTypeOrg:
		repos, err = repo.ListOrgRepos(ctx, d.gh, source.Name)
	case SourceTypeUser:
		repos, err = repo.ListUserRepos(ctx, d.gh, source.Name)
	default:
		return nil, fmt.Errorf("unsupported source type: %s", source.Type)
	}

	if err != nil {
		return nil, err
	}

	var projects []ProjectRef
	for _, r := range repos {
		if r.GetArchived() || r.GetFork() {
			continue
		}

		paths, err := d.FindChangelogPaths(ctx, r.GetOwner().GetLogin(), r.GetName())
		if err != nil {
			// Log but continue - repo might not have a changelog
			continue
		}

		for _, path := range paths {
			projects = append(projects, ProjectRef{
				Path:       path,
				Discovered: true,
			})
		}
	}

	return projects, nil
}

// FindChangelogPaths searches a repo for CHANGELOG.json files.
// Returns paths like "github.com/org/repo" or "github.com/org/repo/subdir".
func (d *DiscoveryClient) FindChangelogPaths(ctx context.Context, owner, repoName string) ([]string, error) {
	var paths []string

	// First check root using gogithub's FileExists
	exists, err := repo.FileExists(ctx, d.gh, owner, repoName, "CHANGELOG.json", nil)
	if err == nil && exists {
		paths = append(paths, fmt.Sprintf("github.com/%s/%s", owner, repoName))
	}

	// Search for nested changelogs (common in monorepos)
	query := fmt.Sprintf("filename:CHANGELOG.json repo:%s/%s", owner, repoName)
	result, _, err := d.gh.Search.Code(ctx, query, nil)
	if err != nil {
		// Search may fail for various reasons, return what we have
		return paths, nil
	}

	for _, codeResult := range result.CodeResults {
		filePath := codeResult.GetPath()
		if filePath == "CHANGELOG.json" {
			continue // Already handled root
		}

		// Extract directory path
		dir := strings.TrimSuffix(filePath, "/CHANGELOG.json")
		dir = strings.TrimSuffix(dir, "CHANGELOG.json")
		if dir != "" {
			dir = strings.TrimSuffix(dir, "/")
			fullPath := fmt.Sprintf("github.com/%s/%s/%s", owner, repoName, dir)
			paths = append(paths, fullPath)
		}
	}

	return paths, nil
}

// FetchRemoteChangelog fetches a CHANGELOG.json from a GitHub repository.
func (d *DiscoveryClient) FetchRemoteChangelog(ctx context.Context, projectPath string) ([]byte, error) {
	host, owner, repoName, subpath, err := ParseProjectPath(projectPath)
	if err != nil {
		return nil, err
	}

	if host != "github.com" {
		return nil, fmt.Errorf("unsupported host: %s (only github.com supported)", host)
	}

	changelogPath := "CHANGELOG.json"
	if subpath != "" {
		changelogPath = subpath + "/CHANGELOG.json"
	}

	// Use gogithub's GetFileContent
	content, err := repo.GetFileContent(ctx, d.gh, owner, repoName, changelogPath, nil)
	if err != nil {
		return nil, fmt.Errorf("fetching changelog: %w", err)
	}

	return content, nil
}

// DiscoveryOptions configures discovery behavior.
type DiscoveryOptions struct {
	IncludeArchived bool
	IncludeForks    bool
	MaxReposPerOrg  int // 0 = unlimited
}

// DefaultDiscoveryOptions returns default discovery options.
func DefaultDiscoveryOptions() DiscoveryOptions {
	return DiscoveryOptions{
		IncludeArchived: false,
		IncludeForks:    false,
		MaxReposPerOrg:  0,
	}
}

// DiscoveryResult contains discovery results and statistics.
type DiscoveryResult struct {
	Projects   []ProjectRef   `json:"projects"`
	Statistics DiscoveryStats `json:"statistics"`
}

// DiscoveryStats provides discovery statistics.
type DiscoveryStats struct {
	SourcesScanned     int `json:"sourcesScanned"`
	ReposScanned       int `json:"reposScanned"`
	ReposWithChangelog int `json:"reposWithChangelog"`
	ChangelogsFound    int `json:"changelogsFound"`
}
