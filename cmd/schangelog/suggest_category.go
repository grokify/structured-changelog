package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-changelog/format"
	"github.com/grokify/structured-changelog/gitlog"
)

var (
	suggestCategoryBatch  bool
	suggestCategoryFormat string
)

// SuggestCategoryOutput is the JSON output for a single suggestion.
type SuggestCategoryOutput struct {
	Input              string                      `json:"input"`
	Suggestions        []gitlog.CategorySuggestion `json:"suggestions"`
	ConventionalCommit *gitlog.ConventionalCommit  `json:"conventional_commit,omitempty"`
}

var suggestCategoryCmd = &cobra.Command{
	Use:   "suggest-category <message>",
	Short: "Suggest changelog category for a commit message",
	Long: `Suggest appropriate changelog categories for commit messages.

This command analyzes commit messages and suggests which changelog category
they should belong to, based on conventional commit types and message content.

Output formats:
  - toon (default): Token-Oriented Object Notation, ~40% fewer tokens than JSON
  - json: Standard JSON with indentation
  - json-compact: Minified JSON

The output includes:
  - Primary suggestion with confidence score and reasoning
  - Conventional commit parsing (if applicable)
  - Alternative suggestions when relevant

Examples:
  # Suggest category for a single message (TOON format, default)
  schangelog suggest-category "feat(auth): add OAuth2 support"

  # Suggest category with JSON output
  schangelog suggest-category --format=json "feat(auth): add OAuth2 support"

  # Batch mode from stdin (one message per line)
  echo -e "feat: add feature\nfix: resolve bug" | schangelog suggest-category --batch`,
	Args: func(cmd *cobra.Command, args []string) error {
		if suggestCategoryBatch {
			return nil // No args required in batch mode
		}
		if len(args) < 1 {
			return fmt.Errorf("requires a commit message argument (or use --batch for stdin)")
		}
		return nil
	},
	RunE: runSuggestCategory,
}

func init() {
	suggestCategoryCmd.Flags().BoolVar(&suggestCategoryBatch, "batch", false, "Read messages from stdin (one per line)")
	suggestCategoryCmd.Flags().StringVar(&suggestCategoryFormat, "format", "toon", "Output format: toon (default), json, json-compact")
	rootCmd.AddCommand(suggestCategoryCmd)
}

func runSuggestCategory(cmd *cobra.Command, args []string) error {
	if suggestCategoryBatch {
		return runSuggestCategoryBatch()
	}

	message := strings.Join(args, " ")
	output := suggestForMessage(message)
	return printSuggestOutput(output)
}

func runSuggestCategoryBatch() error {
	scanner := bufio.NewScanner(os.Stdin)
	var outputs []SuggestCategoryOutput

	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		if message == "" {
			continue
		}
		outputs = append(outputs, suggestForMessage(message))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	return printSuggestOutputs(outputs)
}

func suggestForMessage(message string) SuggestCategoryOutput {
	output := SuggestCategoryOutput{
		Input:       message,
		Suggestions: []gitlog.CategorySuggestion{},
	}

	// Parse conventional commit if applicable
	if cc := gitlog.ParseConventionalCommit(message); cc != nil {
		output.ConventionalCommit = cc
	}

	// Get primary suggestion
	if suggestion := gitlog.SuggestCategoryFromMessage(message); suggestion != nil {
		output.Suggestions = append(output.Suggestions, *suggestion)

		// Add alternative suggestions for ambiguous cases
		alternatives := getAlternativeSuggestions(message, suggestion.Category)
		output.Suggestions = append(output.Suggestions, alternatives...)
	}

	return output
}

// getAlternativeSuggestions returns alternative category suggestions
// for cases where the primary suggestion might not be the only valid option.
func getAlternativeSuggestions(message string, primaryCategory string) []gitlog.CategorySuggestion {
	var alternatives []gitlog.CategorySuggestion
	lower := strings.ToLower(message)

	// Security-related changes might also be Breaking
	if primaryCategory == "Security" && containsAny(lower, []string{"breaking", "remove", "deprecat"}) {
		alternatives = append(alternatives, gitlog.CategorySuggestion{
			Category:   "Breaking",
			Tier:       "standard",
			Confidence: 0.50,
			Reasoning:  "Security fix may introduce breaking changes",
		})
	}

	// Features with auth/security scope might have security implications
	if primaryCategory == "Added" && containsAny(lower, []string{"auth", "security", "permission", "token", "credential"}) {
		alternatives = append(alternatives, gitlog.CategorySuggestion{
			Category:   "Security",
			Tier:       "core",
			Confidence: 0.60,
			Reasoning:  "Feature relates to authentication or security",
		})
	}

	// Performance changes might be considered Changed
	if primaryCategory == "Performance" {
		alternatives = append(alternatives, gitlog.CategorySuggestion{
			Category:   "Changed",
			Tier:       "core",
			Confidence: 0.40,
			Reasoning:  "Performance improvements modify existing behavior",
		})
	}

	// Refactors might have Breaking implications
	if primaryCategory == "Changed" && containsAny(lower, []string{"refactor", "restructure", "rename", "move"}) {
		alternatives = append(alternatives, gitlog.CategorySuggestion{
			Category:   "Breaking",
			Tier:       "standard",
			Confidence: 0.40,
			Reasoning:  "Refactoring may affect public API",
		})
	}

	return alternatives
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func printSuggestOutput(output SuggestCategoryOutput) error {
	return printFormatted(output)
}

func printSuggestOutputs(outputs []SuggestCategoryOutput) error {
	return printFormatted(outputs)
}

func printFormatted(v any) error {
	f, err := format.Parse(suggestCategoryFormat)
	if err != nil {
		return err
	}

	output, err := format.Marshal(v, f)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
