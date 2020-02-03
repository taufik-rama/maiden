package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Service ...
type Service struct {
	Output string
	GRPC   ServiceGRPCList
	HTTP   ServiceHTTPList
}

type serviceWrapper struct {
	Services *InputService `yaml:"services"`
}

// Parse ...
func (s *Service) Parse(filename string) error {

	var config serviceWrapper
	if err := parseYAML(filename, &config); err != nil {
		return err
	}
	if config.Services == nil {
		return nil
	}

	dir := filepath.Dir(filename) + string(os.PathSeparator)

	new(Service).from(config).replace(s)

	for _, imp := range toStringSlice(config.Services.Imports) {
		if err := s.Parse(dir + imp); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) from(cfg serviceWrapper) *Service {

	if cfg.Services == nil {
		return s
	}

	if cfg.Services.GRPC != nil {
		s.GRPC = new(ServiceGRPCList).from(cfg)
	}

	if cfg.Services.HTTP != nil {
		s.HTTP = new(ServiceHTTPList).from(cfg)
	}

	return s
}

func (s Service) replace(other *Service) {

	if strings.TrimSpace(s.Output) != "" {
		other.Output = s.Output
	}

	if other.GRPC == nil {
		other.GRPC = s.GRPC
	} else if s.GRPC != nil {
		s.GRPC.replace(other.GRPC)
	}

	if other.HTTP == nil {
		other.HTTP = s.HTTP
	} else if s.HTTP != nil {
		s.HTTP.replace(other.HTTP)
	}
}
