package gitlog

import (
	"testing"
)

func TestParserParse(t *testing.T) {
	// Sample git log output using our format
	input := `---COMMIT_DELIMITER---
abc123def456789012345678901234567890abcd
abc123d
John Doe
john@example.com
2026-01-04T10:30:00-08:00
feat(auth): add OAuth2 support

Implements OAuth2 flow with PKCE.

Closes #123
---END_BODY---
10	5	src/auth/oauth.go
5	0	src/auth/oauth_test.go
---COMMIT_DELIMITER---
def456abc789012345678901234567890abcdef
def456a
Jane Smith
jane@example.com
2026-01-03T15:00:00-08:00
fix: resolve memory leak (#456)
---END_BODY---
20	10	src/memory/pool.go
`

	parser := NewParser()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(result.Commits))
	}

	// Check first commit
	c1 := result.Commits[0]
	if c1.Hash != "abc123def456789012345678901234567890abcd" {
		t.Errorf("c1.Hash: expected abc123..., got %s", c1.Hash)
	}
	if c1.ShortHash != "abc123d" {
		t.Errorf("c1.ShortHash: expected abc123d, got %s", c1.ShortHash)
	}
	if c1.Author != "John Doe" {
		t.Errorf("c1.Author: expected John Doe, got %s", c1.Author)
	}
	if c1.Date != "2026-01-04" {
		t.Errorf("c1.Date: expected 2026-01-04, got %s", c1.Date)
	}
	if c1.Type != "feat" {
		t.Errorf("c1.Type: expected feat, got %s", c1.Type)
	}
	if c1.Scope != "auth" {
		t.Errorf("c1.Scope: expected auth, got %s", c1.Scope)
	}
	if c1.Subject != "add OAuth2 support" {
		t.Errorf("c1.Subject: expected 'add OAuth2 support', got %s", c1.Subject)
	}
	if c1.Issue != 123 {
		t.Errorf("c1.Issue: expected 123, got %d", c1.Issue)
	}
	if c1.FilesChanged != 2 {
		t.Errorf("c1.FilesChanged: expected 2, got %d", c1.FilesChanged)
	}
	if c1.Insertions != 15 {
		t.Errorf("c1.Insertions: expected 15, got %d", c1.Insertions)
	}
	if c1.Deletions != 5 {
		t.Errorf("c1.Deletions: expected 5, got %d", c1.Deletions)
	}
	if c1.SuggestedCategory != "Added" {
		t.Errorf("c1.SuggestedCategory: expected Added, got %s", c1.SuggestedCategory)
	}

	// Check second commit
	c2 := result.Commits[1]
	if c2.Type != "fix" {
		t.Errorf("c2.Type: expected fix, got %s", c2.Type)
	}
	if c2.PR != 456 {
		t.Errorf("c2.PR: expected 456, got %d", c2.PR)
	}
	if c2.SuggestedCategory != "Fixed" {
		t.Errorf("c2.SuggestedCategory: expected Fixed, got %s", c2.SuggestedCategory)
	}

	// Check summary
	if result.Summary.ByType["feat"] != 1 {
		t.Errorf("expected 1 feat commit, got %d", result.Summary.ByType["feat"])
	}
	if result.Summary.ByType["fix"] != 1 {
		t.Errorf("expected 1 fix commit, got %d", result.Summary.ByType["fix"])
	}
	if result.Summary.BySuggestedCategory["Added"] != 1 {
		t.Errorf("expected 1 Added category, got %d", result.Summary.BySuggestedCategory["Added"])
	}
}

func TestParserParseBreakingChange(t *testing.T) {
	input := `---COMMIT_DELIMITER---
abc123def456789012345678901234567890abcd
abc123d
John Doe
john@example.com
2026-01-04T10:30:00-08:00
feat!: remove deprecated API
---END_BODY---
`

	parser := NewParser()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(result.Commits))
	}

	if !result.Commits[0].Breaking {
		t.Error("expected Breaking to be true")
	}
	if result.Commits[0].SuggestedCategory != "Breaking" {
		t.Errorf("expected Breaking category, got %s", result.Commits[0].SuggestedCategory)
	}
}

func TestParserParseBreakingChangeInBody(t *testing.T) {
	input := `---COMMIT_DELIMITER---
abc123def456789012345678901234567890abcd
abc123d
John Doe
john@example.com
2026-01-04T10:30:00-08:00
feat: change API

BREAKING CHANGE: removes old method signature
---END_BODY---
`

	parser := NewParser()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(result.Commits))
	}

	if !result.Commits[0].Breaking {
		t.Error("expected Breaking to be true from body marker")
	}
}

func TestParserParseNoFiles(t *testing.T) {
	input := `---COMMIT_DELIMITER---
abc123def456789012345678901234567890abcd
abc123d
John Doe
john@example.com
2026-01-04T10:30:00-08:00
feat: add feature
---END_BODY---
10	5	src/file.go
`

	parser := NewParser()
	parser.IncludeFiles = false
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Files should still be parsed for stats
	if result.Commits[0].FilesChanged != 1 {
		t.Errorf("expected 1 file changed, got %d", result.Commits[0].FilesChanged)
	}
}

func TestParserParseEmptyInput(t *testing.T) {
	parser := NewParser()
	result, err := parser.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Commits) != 0 {
		t.Errorf("expected 0 commits, got %d", len(result.Commits))
	}
}

func TestParserParseBinaryFiles(t *testing.T) {
	input := `---COMMIT_DELIMITER---
abc123def456789012345678901234567890abcd
abc123d
John Doe
john@example.com
2026-01-04T10:30:00-08:00
feat: add image
---END_BODY---
-	-	image.png
10	5	src/file.go
`

	parser := NewParser()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Binary files have "-" for insertions/deletions
	c := result.Commits[0]
	if c.FilesChanged != 2 {
		t.Errorf("expected 2 files changed, got %d", c.FilesChanged)
	}
	if c.Insertions != 10 {
		t.Errorf("expected 10 insertions (from non-binary), got %d", c.Insertions)
	}
	if c.Deletions != 5 {
		t.Errorf("expected 5 deletions (from non-binary), got %d", c.Deletions)
	}
}

func TestParseSimple(t *testing.T) {
	input := `abc123def456789012345678901234567890abcd|abc123d|John Doe|john@example.com|2026-01-04T10:30:00-08:00|feat(auth): add OAuth2 support
def456abc789012345678901234567890abcdef|def456a|Jane Smith|jane@example.com|2026-01-03T15:00:00-08:00|fix: resolve bug (#123)
`

	result, err := ParseSimple(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(result.Commits))
	}

	c1 := result.Commits[0]
	if c1.Type != "feat" {
		t.Errorf("c1.Type: expected feat, got %s", c1.Type)
	}
	if c1.Scope != "auth" {
		t.Errorf("c1.Scope: expected auth, got %s", c1.Scope)
	}

	c2 := result.Commits[1]
	if c2.PR != 123 {
		t.Errorf("c2.PR: expected 123, got %d", c2.PR)
	}
}

func TestNewParseResult(t *testing.T) {
	result := NewParseResult()

	if result.Commits == nil {
		t.Error("Commits should be initialized")
	}
	if result.Summary.ByType == nil {
		t.Error("ByType should be initialized")
	}
	if result.Summary.BySuggestedCategory == nil {
		t.Error("BySuggestedCategory should be initialized")
	}
	if result.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should be set")
	}
}

func TestParseResultAddCommit(t *testing.T) {
	result := NewParseResult()

	commit := Commit{
		Hash:              "abc123",
		Type:              "feat",
		SuggestedCategory: "Added",
		FilesChanged:      3,
		Insertions:        100,
		Deletions:         50,
	}

	result.AddCommit(commit)

	if result.Range.CommitCount != 1 {
		t.Errorf("expected CommitCount 1, got %d", result.Range.CommitCount)
	}
	if result.Summary.ByType["feat"] != 1 {
		t.Errorf("expected ByType[feat] = 1, got %d", result.Summary.ByType["feat"])
	}
	if result.Summary.BySuggestedCategory["Added"] != 1 {
		t.Errorf("expected BySuggestedCategory[Added] = 1, got %d", result.Summary.BySuggestedCategory["Added"])
	}
	if result.Summary.TotalFilesChanged != 3 {
		t.Errorf("expected TotalFilesChanged 3, got %d", result.Summary.TotalFilesChanged)
	}
	if result.Summary.TotalInsertions != 100 {
		t.Errorf("expected TotalInsertions 100, got %d", result.Summary.TotalInsertions)
	}
	if result.Summary.TotalDeletions != 50 {
		t.Errorf("expected TotalDeletions 50, got %d", result.Summary.TotalDeletions)
	}
}

func TestComputeContributors(t *testing.T) {
	result := NewParseResult()

	// Add commits from different authors
	result.Commits = []Commit{
		{Author: "Alice", IsExternal: true},
		{Author: "Alice", IsExternal: true},
		{Author: "Alice", IsExternal: true},
		{Author: "Bob", IsExternal: false},
		{Author: "Bob", IsExternal: false},
		{Author: "Charlie", IsExternal: true},
		{Author: "", IsExternal: false}, // empty author should be skipped
	}

	result.ComputeContributors()

	if len(result.Contributors) != 3 {
		t.Fatalf("expected 3 contributors, got %d", len(result.Contributors))
	}

	// External contributors should come first, sorted by commit count
	// Alice (3 commits, external), Charlie (1 commit, external), Bob (2 commits, internal)
	if result.Contributors[0].Name != "Alice" {
		t.Errorf("expected first contributor to be Alice, got %s", result.Contributors[0].Name)
	}
	if result.Contributors[0].CommitCount != 3 {
		t.Errorf("expected Alice to have 3 commits, got %d", result.Contributors[0].CommitCount)
	}
	if !result.Contributors[0].IsExternal {
		t.Error("expected Alice to be external")
	}

	if result.Contributors[1].Name != "Charlie" {
		t.Errorf("expected second contributor to be Charlie, got %s", result.Contributors[1].Name)
	}
	if !result.Contributors[1].IsExternal {
		t.Error("expected Charlie to be external")
	}

	if result.Contributors[2].Name != "Bob" {
		t.Errorf("expected third contributor to be Bob, got %s", result.Contributors[2].Name)
	}
	if result.Contributors[2].IsExternal {
		t.Error("expected Bob to be internal")
	}
}

func TestComputeContributorsEmpty(t *testing.T) {
	result := NewParseResult()
	result.ComputeContributors()

	if len(result.Contributors) != 0 {
		t.Errorf("expected 0 contributors, got %d", len(result.Contributors))
	}
}

func TestComputeContributorsAllInternal(t *testing.T) {
	result := NewParseResult()
	result.Commits = []Commit{
		{Author: "Maintainer1", IsExternal: false},
		{Author: "Maintainer2", IsExternal: false},
	}

	result.ComputeContributors()

	if len(result.Contributors) != 2 {
		t.Fatalf("expected 2 contributors, got %d", len(result.Contributors))
	}

	// All internal, so just sorted by commit count
	for _, c := range result.Contributors {
		if c.IsExternal {
			t.Errorf("expected all contributors to be internal, got external: %s", c.Name)
		}
	}
}
