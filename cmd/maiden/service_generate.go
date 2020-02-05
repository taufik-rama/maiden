package main

import (
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/generator.v1"
	"github.com/taufik-rama/maiden/internal"
)

// ServiceGenerate command handler
type ServiceGenerate struct{}

// Command returns `generate` command process
func (s ServiceGenerate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate the services",
		Run:     s.RunCommand,
	}
	return cmd
}

// RunCommand runs `generate` command
func (s *ServiceGenerate) RunCommand(cmd *cobra.Command, args []string) {

	service := new(config.Service)
	if err := service.Parse(internal.Args.ConfigFile); err != nil {
		cmd.PrintErrln(err)
	}

	if activeServices(service) == 0 {
		internal.Print("No services defined, nothing to do")
		return
	}

	if len(service.GRPC) != 0 {
		(generator.GRPC{
			Output:     service.Output,
			ConfigGRPC: service.GRPC,
		}).GenerateCommand(cmd, args)
	}

	if len(service.HTTP) != 0 {
		(generator.HTTP{
			Output:     service.Output,
			ConfigHTTP: service.HTTP,
		}).GenerateCommand(cmd, args)
	}
}

func activeServices(service *config.Service) int {
	var active int
	active += len(service.GRPC)
	active += len(service.HTTP)
	return active
}
