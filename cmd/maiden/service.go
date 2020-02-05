package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
)

// Service command handler
type Service struct {
	statusAll  bool
	statusHTTP bool
	statusGRPC bool
}

// Command returns `service` command process
func (s Service) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"s"},
		Short:   "Parse available dummy services",
		Run:     s.RunCommand,
	}
	cmd.PersistentFlags().BoolVarP(&s.statusAll, "status-all", "", false, "Show all services")
	cmd.PersistentFlags().BoolVarP(&s.statusHTTP, "status-http", "", false, "Show HTTP services")
	cmd.PersistentFlags().BoolVarP(&s.statusGRPC, "status-grpc", "", false, "Show GRPC services")
	return cmd
}

// RunCommand runs `fixture` command
func (s *Service) RunCommand(cmd *cobra.Command, args []string) {
	if s.statusFlag() {
		if err := s.printStatus(); err != nil {
			cmd.PrintErrln(err)
		}
		return
	}
	if err := new(config.Service).Parse(internal.Args.ConfigFile); err != nil {
		cmd.PrintErrln(err)
	}
	cmd.Println("All services parsed, check with `--status-x` for the values")
}

func (s Service) statusFlag() bool {
	return s.statusAll || s.statusHTTP || s.statusGRPC
}

func (s Service) printStatus() error {

	if s.statusAll {
		s.statusGRPC = true
		s.statusHTTP = true
	}

	service := new(config.Service)
	if err := service.Parse(internal.Args.ConfigFile); err != nil {
		return err
	}

	if s.statusGRPC && len(service.GRPC) > 0 {
		s.printStatusGRPC(service.GRPC)
	}

	if s.statusHTTP && len(service.HTTP) > 0 {
		if s.statusGRPC && len(service.GRPC) > 0 {
			fmt.Println()
		}
		s.printStatusHTTP(service.HTTP)
	}

	return nil
}

func (s Service) printStatusGRPC(services config.ServiceGRPCList) {
	fmt.Printf("%d GRPC Service(s)\n", len(services))
	for name, detail := range services {
		fmt.Printf("  %s on :%d using %s\n", name, detail.Port, detail.Definition)
		for method := range detail.Methods {
			fmt.Printf("    Method %s: %d registered condition(s)\n", method, len(detail.Conditions[method]))
		}
	}
}

func (s Service) printStatusHTTP(services config.ServiceHTTPList) {
	fmt.Printf("%d HTTP Service(s)\n", len(services))
	for name, detail := range services {
		fmt.Printf("  %s on :%d\n", name, detail.Port)
		for endpoint := range detail.Endpoints {
			fmt.Printf("    Endpoint `%s`: %d registered condition(s)\n", endpoint, len(detail.Endpoints[endpoint]))
		}
	}
}
