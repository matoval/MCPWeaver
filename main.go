package main

import (
	"fmt"
	"os"

	"MCPWeaver/internal/cmd"
)

// Version information - will be set during build
var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"
)

func main() {
	// Set version information for the CLI
	cmd.SetVersionInfo(version, buildTime, commit)

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}