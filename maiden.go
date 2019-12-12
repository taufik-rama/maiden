package main

import (
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/internal/config"
	"github.com/taufik-rama/maiden/internal/fixture"
	"github.com/taufik-rama/maiden/internal/service"
)

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

	rootCmd.PersistentFlags().StringVarP(&config.Args.ConfigFile, "config-file", "c", "maiden.yaml", "Maiden configuration file")
	rootCmd.PersistentFlags().BoolVarP(&config.Args.Verbose, "verbose", "v", false, "")

	rootCmd.AddCommand((fixture.Handler{}).Command())
	rootCmd.AddCommand((service.Handler{}).Command())
	rootCmd.AddCommand(VersionCommand)

	rootCmd.Execute()
}
