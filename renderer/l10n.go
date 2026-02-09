package renderer

import (
	"embed"
	"os"
	"strings"

	"github.com/grokify/structured-locale/messages"
)

//go:embed locales/*.json
var defaultLocales embed.FS

// defaultBundle holds embedded default translations.
var defaultBundle *messages.Bundle

func init() {
	defaultBundle = messages.NewBundle("en")

	entries, err := defaultLocales.ReadDir("locales")
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		data, err := defaultLocales.ReadFile("locales/" + e.Name())
		if err != nil {
			continue
		}

		loc := strings.TrimSuffix(e.Name(), ".json")
		_ = defaultBundle.AddLocale(loc, data)
	}
}

// getLocalizer returns a localizer for the given options.
func getLocalizer(opts Options) *messages.Localizer {
	locale := opts.Locale
	if locale == "" {
		locale = "en"
	}

	if opts.LocaleOverrides != "" {
		data, err := os.ReadFile(opts.LocaleOverrides)
		if err == nil {
			_ = defaultBundle.AddLocaleOverrides(locale, data)
		}
	}

	return defaultBundle.Localizer(locale)
}

// categoryToMessageID converts a changelog category name to a message ID.
// For example, "Added" -> "category.added", "Known Issues" -> "category.known_issues".
func categoryToMessageID(category string) string {
	// Convert to lowercase and replace spaces with underscores
	id := strings.ToLower(category)
	id = strings.ReplaceAll(id, " ", "_")
	return "category." + id
}
