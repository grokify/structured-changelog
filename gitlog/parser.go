package gitlog

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// commitDelimiter is a unique marker used to separate commits in git log output.
const commitDelimiter = "---COMMIT_DELIMITER---"

// GitLogFormat is the format string to use with git log for parsing.
// Use: git log --format="---COMMIT_DELIMITER---%n%H%n%h%n%an%n%ae%n%aI%n%s%n%b---END_BODY---" --numstat
const GitLogFormat = commitDelimiter + "%n%H%n%h%n%an%n%ae%n%aI%n%s%n%b---END_BODY---"

// numstatRegex matches numstat output lines: "123\t456\tfilename"
var numstatRegex = regexp.MustCompile(`^(\d+|-)\t(\d+|-)\t(.+)$`)

// Parser parses git log output into structured commits.
type Parser struct {
	IncludeFiles bool
}

// NewParser creates a new git log parser.
func NewParser() *Parser {
	return &Parser{
		IncludeFiles: true,
	}
}

// Parse parses git log output and returns a ParseResult.
func (p *Parser) Parse(input string) (*ParseResult, error) {
	result := NewParseResult()

	commits := strings.Split(input, commitDelimiter)
	for _, commitBlock := range commits {
		commitBlock = strings.TrimSpace(commitBlock)
		if commitBlock == "" {
			continue
		}

		commit := p.parseCommitBlock(commitBlock)
		if commit != nil {
			result.AddCommit(*commit)
		}
	}

	return result, nil
}

// parseCommitBlock parses a single commit block.
// Returns nil if the block is malformed.
func (p *Parser) parseCommitBlock(block string) *Commit {
	// Split on ---END_BODY--- to separate commit info from numstat
	parts := strings.SplitN(block, "---END_BODY---", 2)
	commitPart := strings.TrimSpace(parts[0])

	lines := strings.Split(commitPart, "\n")
	if len(lines) < 6 {
		return nil // Not enough lines for a valid commit
	}

	commit := &Commit{
		Hash:        strings.TrimSpace(lines[0]),
		ShortHash:   strings.TrimSpace(lines[1]),
		Author:      strings.TrimSpace(lines[2]),
		AuthorEmail: strings.TrimSpace(lines[3]),
		Message:     strings.TrimSpace(lines[5]),
	}

	// Parse date
	dateStr := strings.TrimSpace(lines[4])
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		commit.Date = t.Format("2006-01-02")
	} else {
		commit.Date = dateStr
	}

	// Extract body (lines after subject)
	if len(lines) > 6 {
		bodyLines := lines[6:]
		commit.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
	}

	// Set subject (first line of message or subject line)
	commit.Subject = commit.Message

	// Parse conventional commit
	fullMessage := commit.Message
	if commit.Body != "" {
		fullMessage = commit.Message + "\n" + commit.Body
	}

	if cc := ParseConventionalCommit(commit.Message); cc != nil {
		commit.Type = cc.Type
		commit.Scope = cc.Scope
		commit.Subject = cc.Subject
		commit.Breaking = cc.Breaking
	}

	// Check for breaking change in body
	if !commit.Breaking && commit.Body != "" {
		commit.Breaking = HasBreakingChangeMarker(commit.Body)
	}

	// Extract issue and PR references
	commit.Issue = ExtractIssueNumber(fullMessage)
	commit.PR = ExtractPRNumber(commit.Message)

	// Parse numstat if present (always parse for stats, optionally include file names)
	if len(parts) > 1 {
		p.parseNumstat(commit, strings.TrimSpace(parts[1]))
	}

	// Suggest category
	if suggestion := SuggestCategoryFromMessage(fullMessage); suggestion != nil {
		commit.SuggestedCategory = suggestion.Category
	}

	return commit
}

// parseNumstat parses the numstat output and updates the commit.
// Stats (insertions, deletions, files changed) are always parsed.
// File names are only included if IncludeFiles is true.
func (p *Parser) parseNumstat(commit *Commit, numstat string) {
	scanner := bufio.NewScanner(strings.NewReader(numstat))
	for scanner.Scan() {
		line := scanner.Text()
		matches := numstatRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		// Parse insertions (can be "-" for binary files)
		if matches[1] != "-" {
			if ins, err := strconv.Atoi(matches[1]); err == nil {
				commit.Insertions += ins
			}
		}

		// Parse deletions (can be "-" for binary files)
		if matches[2] != "-" {
			if del, err := strconv.Atoi(matches[2]); err == nil {
				commit.Deletions += del
			}
		}

		// Only include file names if requested
		if p.IncludeFiles {
			commit.Files = append(commit.Files, matches[3])
		}
		commit.FilesChanged++
	}
}

// ParseSimple parses a simpler git log format without numstat.
// Use with: git log --format="%H|%h|%an|%ae|%aI|%s"
func ParseSimple(input string) (*ParseResult, error) {
	result := NewParseResult()

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 6)
		if len(parts) < 6 {
			continue
		}

		commit := &Commit{
			Hash:        parts[0],
			ShortHash:   parts[1],
			Author:      parts[2],
			AuthorEmail: parts[3],
			Message:     parts[5],
			Subject:     parts[5],
		}

		// Parse date
		if t, err := time.Parse(time.RFC3339, parts[4]); err == nil {
			commit.Date = t.Format("2006-01-02")
		} else {
			commit.Date = parts[4]
		}

		// Parse conventional commit
		if cc := ParseConventionalCommit(commit.Message); cc != nil {
			commit.Type = cc.Type
			commit.Scope = cc.Scope
			commit.Subject = cc.Subject
			commit.Breaking = cc.Breaking
		}

		// Extract references
		commit.Issue = ExtractIssueNumber(commit.Message)
		commit.PR = ExtractPRNumber(commit.Message)

		// Suggest category
		if suggestion := SuggestCategoryFromMessage(commit.Message); suggestion != nil {
			commit.SuggestedCategory = suggestion.Category
		}

		result.AddCommit(*commit)
	}

	return result, scanner.Err()
}
