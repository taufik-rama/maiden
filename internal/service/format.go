package service

import (
	"sort"
	"strconv"

	"github.com/taufik-rama/maiden/internal/config"
)

const stats = `Services status
  HTTP: %s
  GRPC: %s
`

func formatHTTP(services config.HTTPServices) string {
	var stat string
	prefix := "    "
	var names []string
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		stat += "\n"
		stat += (prefix + "- " + name + ": " + formatHTTPService(services[name]))
	}
	return stat
}

func formatHTTPService(service config.HTTPService) string {
	var stat string
	prefix := "      "

	// Port
	stat += "\n"
	stat += (prefix + "port: " + strconv.FormatUint(uint64(service.Port), 10))

	// Endpoints
	stat += "\n"
	stat += (prefix + "endpoints: " + formatHTTPServiceEndpoints(service.Endpoints))

	return stat
}

func formatHTTPServiceEndpoints(endpoints config.HTTPEndpoints) string {
	var stat string
	prefix := "        "
	var names []string
	for endpoint := range endpoints {
		names = append(names, endpoint)
	}
	sort.Strings(names)
	for _, name := range names {
		stat += "\n"
		stat += (prefix + "- " + name)
	}
	return stat
}

func formatGRPC(services config.GRPCServices) string {
	var stat string
	prefix := "    "
	var names []string
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		stat += "\n"
		stat += (prefix + "- " + name + ": " + formatGRPCService(services[name]))
	}
	return stat
}

func formatGRPCService(service config.GRPCService) string {
	var stat string
	prefix := "      "

	// Port
	stat += "\n"
	stat += (prefix + "port: " + strconv.FormatUint(uint64(service.Port), 10))

	// Definition
	stat += "\n"
	stat += (prefix + "definition: " + service.Definition)

	// Methods
	stat += "\n"
	stat += (prefix + "methods: " + formatGRPCServiceMethods(service.Methods))

	return stat
}

func formatGRPCServiceMethods(methods config.GRPCMethods) string {
	var stat string
	prefix := "        "
	var names []string
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		stat += "\n"
		stat += (prefix + "- " + name)
	}
	return stat
}
