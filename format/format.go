// Package format provides output format abstraction for CLI commands.
package format

import (
	"encoding/json"
	"fmt"

	toon "github.com/toon-format/toon-go"
)

// Format represents an output format type.
type Format string

// Supported output formats.
const (
	TOON        Format = "toon"
	JSON        Format = "json"
	JSONCompact Format = "json-compact"
)

// Parse parses a format string into a Format type.
// Empty string defaults to TOON.
func Parse(s string) (Format, error) {
	switch s {
	case "toon", "":
		return TOON, nil
	case "json":
		return JSON, nil
	case "json-compact":
		return JSONCompact, nil
	default:
		return "", fmt.Errorf("unknown format %q: use toon, json, or json-compact", s)
	}
}

// Marshal serializes v to the specified format.
func Marshal(v any, f Format) ([]byte, error) {
	switch f {
	case TOON:
		return toon.Marshal(v)
	case JSON:
		return json.MarshalIndent(v, "", "  ")
	case JSONCompact:
		return json.Marshal(v)
	default:
		return toon.Marshal(v)
	}
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}
