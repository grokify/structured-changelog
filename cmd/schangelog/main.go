// Command schangelog is the Structured Changelog CLI tool.
//
// Usage:
//
//	schangelog validate CHANGELOG.json
//	schangelog generate CHANGELOG.json -o CHANGELOG.md
//	schangelog version
package main

import (
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
