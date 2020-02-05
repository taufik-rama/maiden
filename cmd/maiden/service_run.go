package main

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/generator.v1"
	"github.com/taufik-rama/maiden/internal"
)

// ServiceRun command handler
type ServiceRun struct{}

// Command returns `run` command process
func (s ServiceRun) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Build & run the services (experimental)",
		Run:   s.RunCommand,
	}
	return cmd
}

// RunCommand runs `run` command
func (s ServiceRun) RunCommand(cmd *cobra.Command, args []string) {

	service := new(config.Service)
	if err := service.Parse(internal.Args.ConfigFile); err != nil {
		cmd.PrintErrln(err)
	}

	servicesWG := new(sync.WaitGroup)

	(generator.GRPC{
		Output:     service.Output,
		ConfigGRPC: service.GRPC,
		WaitGroup:  servicesWG,
	}).RunCommand(cmd, args)

	(generator.HTTP{
		Output:     service.Output,
		ConfigHTTP: service.HTTP,
		WaitGroup:  servicesWG,
	}).RunCommand(cmd, args)

	servicesWG.Wait()
}
