// Command sclog is the Structured Changelog CLI tool.
//
// Usage:
//
//	sclog validate CHANGELOG.json
//	sclog generate CHANGELOG.json -o CHANGELOG.md
//	sclog version
package main

import (
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
