package config

// MaidenConfig correspond with `maiden.yml` file
type MaidenConfig struct {
	Fixtures *FixtureConfig `yaml:"fixtures"`
	Services *ServiceConfig `yaml:"services"`

	// Maiden configuration file
	config string `yaml:"-"`
}

// New maiden configuration object with `file` as filename
// relative to binary call directory.
func New(file string) MaidenConfig {
	return MaidenConfig{
		Fixtures: new(FixtureConfig),
		Services: &ServiceConfig{
			HTTP: make(HTTPServices),
			GRPC: make(GRPCServices),
		},
		config: file,
	}
}

// ResolveFixtures will resolves the fixtures config statements
func (m *MaidenConfig) ResolveFixtures() error {
	return m.Fixtures.resolve(m.config)
}

// ResolveServices will resolves the services config statements
func (m *MaidenConfig) ResolveServices() error {
	return m.Services.resolve(m.config)
}
