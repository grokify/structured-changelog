package aggregate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolvedProject represents a project with resolved local paths.
type ResolvedProject struct {
	Ref           ProjectRef // Original project reference
	LocalPath     string     // Absolute path to project directory (empty if not found)
	ChangelogPath string     // Full path to CHANGELOG.json (empty if not found)
	IsLocal       bool       // true if found locally
	NeedsApproval bool       // true if remote fetch required
}

// Resolver resolves Go-style paths to local filesystem locations.
type Resolver struct {
	SearchPaths []string // Directories to search (e.g., ~/go/src)
}

// NewResolver creates a resolver with default search paths.
func NewResolver() *Resolver {
	return &Resolver{
		SearchPaths: DefaultSearchPaths(),
	}
}

// NewResolverWithPaths creates a resolver with custom search paths.
func NewResolverWithPaths(paths []string) *Resolver {
	return &Resolver{
		SearchPaths: paths,
	}
}

// DefaultSearchPaths returns the default directories to search for projects.
func DefaultSearchPaths() []string {
	paths := []string{}

	// Get home directory (cross-platform via os.UserHomeDir)
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		paths = append(paths, filepath.Join(home, "go", "src"))
	}

	// Add $GOPATH/src if set and different from ~/go
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		gopathSrc := filepath.Join(gopath, "src")
		// Only add if not already in the list
		if len(paths) == 0 || paths[0] != gopathSrc {
			paths = append(paths, gopathSrc)
		}
	}

	return paths
}

// Resolve attempts to resolve all project references to local paths.
// Projects that cannot be found locally are marked as NeedsApproval.
func (r *Resolver) Resolve(refs []ProjectRef) ([]ResolvedProject, error) {
	results := make([]ResolvedProject, len(refs))

	for i, ref := range refs {
		resolved := r.resolveOne(ref)
		results[i] = resolved
	}

	return results, nil
}

// ResolveOne resolves a single project reference.
func (r *Resolver) ResolveOne(ref ProjectRef) ResolvedProject {
	return r.resolveOne(ref)
}

func (r *Resolver) resolveOne(ref ProjectRef) ResolvedProject {
	result := ResolvedProject{
		Ref: ref,
	}

	// 1. Check explicit LocalPath override
	if ref.LocalPath != "" {
		if r.checkLocalPath(ref.LocalPath, &result) {
			return result
		}
	}

	// 2. Search through search paths
	for _, searchPath := range r.SearchPaths {
		localPath := filepath.Join(searchPath, ref.Path)
		if r.checkLocalPath(localPath, &result) {
			return result
		}
	}

	// 3. Not found locally - mark for remote fetch
	result.NeedsApproval = true
	return result
}

// checkLocalPath checks if a path contains a valid CHANGELOG.json.
// If found, populates the result and returns true.
func (r *Resolver) checkLocalPath(localPath string, result *ResolvedProject) bool {
	changelogPath := filepath.Join(localPath, "CHANGELOG.json")

	info, err := os.Stat(changelogPath)
	if err != nil || info.IsDir() {
		return false
	}

	result.LocalPath = localPath
	result.ChangelogPath = changelogPath
	result.IsLocal = true
	result.NeedsApproval = false
	return true
}

// Summary returns a summary of resolution results.
func (r *Resolver) Summary(resolved []ResolvedProject) ResolutionSummary {
	var summary ResolutionSummary

	for _, rp := range resolved {
		if rp.IsLocal {
			summary.LocalCount++
			summary.LocalProjects = append(summary.LocalProjects, rp.Ref.Path)
		} else {
			summary.RemoteCount++
			summary.RemoteProjects = append(summary.RemoteProjects, rp.Ref.Path)
		}
	}

	summary.TotalCount = summary.LocalCount + summary.RemoteCount
	return summary
}

// ResolutionSummary provides statistics about resolved projects.
type ResolutionSummary struct {
	TotalCount     int      `json:"totalCount"`
	LocalCount     int      `json:"localCount"`
	RemoteCount    int      `json:"remoteCount"`
	LocalProjects  []string `json:"localProjects,omitempty"`
	RemoteProjects []string `json:"remoteProjects,omitempty"`
}

// FilterLocal returns only projects that were found locally.
func FilterLocal(resolved []ResolvedProject) []ResolvedProject {
	var local []ResolvedProject
	for _, rp := range resolved {
		if rp.IsLocal {
			local = append(local, rp)
		}
	}
	return local
}

// FilterRemote returns only projects that need remote fetch.
func FilterRemote(resolved []ResolvedProject) []ResolvedProject {
	var remote []ResolvedProject
	for _, rp := range resolved {
		if rp.NeedsApproval {
			remote = append(remote, rp)
		}
	}
	return remote
}

// ParseProjectPath parses a project path into its components.
// Example: "github.com/org/repo/subdir" -> host, owner, repo, subpath
func ParseProjectPath(path string) (host, owner, repo, subpath string, err error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return "", "", "", "", fmt.Errorf("invalid project path: %s (expected host/owner/repo)", path)
	}

	host = parts[0]
	owner = parts[1]
	repo = parts[2]

	if len(parts) > 3 {
		subpath = strings.Join(parts[3:], "/")
	}

	return host, owner, repo, subpath, nil
}

// GitHubRepoURL returns the GitHub URL for a project path.
func GitHubRepoURL(path string) (string, error) {
	host, owner, repo, _, err := ParseProjectPath(path)
	if err != nil {
		return "", err
	}

	if host != "github.com" {
		return "", fmt.Errorf("unsupported host: %s (only github.com supported)", host)
	}

	return fmt.Sprintf("https://github.com/%s/%s", owner, repo), nil
}
