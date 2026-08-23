package gitlog

import (
	"regexp"
	"strconv"
	"strings"
)

// ConventionalCommit represents parsed components of a conventional commit message.
type ConventionalCommit struct {
	Type     string `json:"type"`
	Scope    string `json:"scope,omitempty"`
	Subject  string `json:"subject"`
	Breaking bool   `json:"breaking"`
}

// conventionalCommitRegex matches the conventional commit format:
// type(scope)!: subject
// type!: subject
// type(scope): subject
// type: subject
var conventionalCommitRegex = regexp.MustCompile(
	`^([a-zA-Z]+)(?:\(([^)]+)\))?(!)?\s*:\s*(.+)$`,
)

// issueRefRegex matches issue references like #123, Closes #123, Fixes #456
var issueRefRegex = regexp.MustCompile(`(?i)(?:closes?|fixes?|resolves?|refs?)?\s*#(\d+)`)

// prRefRegex matches PR references in subject like "(#123)" at end of line
var prRefRegex = regexp.MustCompile(`\(#(\d+)\)\s*$`)

// breakingChangeRegex matches BREAKING CHANGE: in body
var breakingChangeRegex = regexp.MustCompile(`(?i)^BREAKING[ -]CHANGE\s*:`)

// refsTrailerRegex matches `Refs:` git trailer lines (case-insensitive key).
// The value may list multiple IDs; RMI IDs are extracted from it separately.
var refsTrailerRegex = regexp.MustCompile(`(?im)^\s*Refs\s*:\s*(.+?)\s*$`)

// rmiIDRegex matches a roadmap item ID like RMI-SCHANGELOG-101.
var rmiIDRegex = regexp.MustCompile(`\bRMI-[A-Z0-9]+-\d+\b`)

// ExtractRMIs extracts roadmap item IDs (RMI-<SLUG>-<NNN>) from a commit
// message's `Refs:` git trailers. A commit may carry multiple `Refs:` trailers,
// and a trailer value may list multiple IDs. Results are deduplicated in
// first-seen order. Only IDs on `Refs:` trailer lines are returned, so
// incidental mentions of an RMI in prose are ignored. Returns nil when none
// are present.
func ExtractRMIs(message string) []string {
	var rmis []string
	seen := make(map[string]bool)
	for _, trailer := range refsTrailerRegex.FindAllStringSubmatch(message, -1) {
		for _, id := range rmiIDRegex.FindAllString(trailer[1], -1) {
			if !seen[id] {
				seen[id] = true
				rmis = append(rmis, id)
			}
		}
	}
	return rmis
}

// ParseConventionalCommit parses a commit message into conventional commit components.
// Returns nil if the message doesn't follow conventional commit format.
func ParseConventionalCommit(message string) *ConventionalCommit {
	// Get first line
	firstLine := strings.Split(message, "\n")[0]

	matches := conventionalCommitRegex.FindStringSubmatch(firstLine)
	if matches == nil {
		return nil
	}

	cc := &ConventionalCommit{
		Type:     strings.ToLower(matches[1]),
		Scope:    matches[2],
		Breaking: matches[3] == "!",
		Subject:  strings.TrimSpace(matches[4]),
	}

	return cc
}

// IsConventionalCommit returns true if the message follows conventional commit format.
func IsConventionalCommit(message string) bool {
	return ParseConventionalCommit(message) != nil
}

// ExtractIssueNumber extracts the first issue number from a commit message.
// It looks for patterns like #123, Closes #123, Fixes #456.
func ExtractIssueNumber(message string) int {
	matches := issueRefRegex.FindStringSubmatch(message)
	if matches == nil {
		return 0
	}
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}
	return num
}

// ExtractPRNumber extracts a PR number from the subject line.
// It looks for patterns like "(#123)" at the end of the subject.
func ExtractPRNumber(subject string) int {
	matches := prRefRegex.FindStringSubmatch(subject)
	if matches == nil {
		return 0
	}
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}
	return num
}

// HasBreakingChangeMarker checks if the message body contains BREAKING CHANGE:.
func HasBreakingChangeMarker(body string) bool {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if breakingChangeRegex.MatchString(line) {
			return true
		}
	}
	return false
}

// KnownConventionalTypes are the standard conventional commit types.
var KnownConventionalTypes = []string{
	"feat",
	"fix",
	"docs",
	"style",
	"refactor",
	"perf",
	"test",
	"build",
	"ci",
	"chore",
	"revert",
	"security",
	"deps",
}

// IsKnownType returns true if the type is a recognized conventional commit type.
func IsKnownType(t string) bool {
	t = strings.ToLower(t)
	for _, known := range KnownConventionalTypes {
		if t == known {
			return true
		}
	}
	return false
}
