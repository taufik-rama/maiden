package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/internal"
)

// These should be filled from build metadata
var maidenVersion string
var goVersion string
var latestCommit string
var buildDate string
var buildOS string

// VersionCommand is cobra command for Maiden `version` metadata
var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print Maiden version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`Maiden command-line info
  Version: 	%s
  Go version: 	%s
  Git commit: 	%s
  Built: 	%s
  OS/Arch: 	%s
`, maidenVersion, goVersion, latestCommit, buildDate, buildOS)
	},
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "maiden",
		Short: "'maiden' is a helper tools to automate things",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Usage()
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&internal.Args.ConfigFile, "config-file", "c", "maiden.yaml", "Maiden configuration file")
	rootCmd.PersistentFlags().BoolVarP(&internal.Args.Verbose, "verbose", "v", false, "")

	rootCmd.AddCommand((Fixture{}).Command())
	rootCmd.AddCommand((Service{}).Command())
	rootCmd.AddCommand(VersionCommand)

	rootCmd.Execute()
}
