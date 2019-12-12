package config

// Args used to store the global command flags
var Args struct {

	// ConfigFile for custom maiden config directory
	ConfigFile string

	// Verbose determines the verbosity of the command output
	Verbose bool

	// Fixture is the args related to `service` command
	Service struct {

		// PreferImports tell maiden to replace the values for import value (if any)
		PreferImports bool
	}
}
