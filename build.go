package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var goVersion string
var latestCommit string
var buildDate string
var buildOS string

const maidenVersion = "v0.1.0-alpha"

// VersionCommand is cobra command for Maiden `version` metadata
var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print Maiden version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`Maiden command
  Version: 	%s
  Go version: 	%s
  Git commit: 	%s
  Built: 	%s
  OS/Arch: 	%s
`, maidenVersion, goVersion, latestCommit, buildDate, buildOS)
	},
}
