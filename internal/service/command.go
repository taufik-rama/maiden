package service

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/internal/config"
	// "github.com/taufik-rama/maiden/internal/service/generate"
	// "github.com/taufik-rama/maiden/internal/service/run"
)

// Handler for `service` command
type Handler struct{}

// Command returns `service` command process
func (h Handler) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"s"},
		Short:   "Generate pre-defined dummy services",
		Run: func(cmd *cobra.Command, args []string) {
			c := config.New(config.Args.ConfigFile)
			if err := c.ResolveServices(); err != nil {
				log.Fatalln(err)
			}
			for name, service := range c.Services.HTTP {
				fmt.Println(">>> name", name)
				fmt.Println(">>> details")
				fmt.Println(service.Port)
				fmt.Println(service.Endpoints)
			}
			for name, service := range c.Services.GRPC {
				fmt.Println(">>> name", name)
				fmt.Println(">>> details")
				fmt.Println(service.Port)
				fmt.Println(service.Definition)
				fmt.Println(service.Methods)
				fmt.Println(service.Conditions)
			}
			// fmt.Printf(
			// 	stats,
			// 	formatHTTP(c.Services.HTTP),
			// 	formatGRPC(c.Services.GRPC),
			// )
		},
	}
	// cmd.AddCommand(generate.Command())
	// cmd.AddCommand(run.Command())

	cmd.PersistentFlags().BoolVarP(&config.Args.Service.PreferImports, "prefer-imports", "", false, "prioritize the imports value")

	return cmd
}
