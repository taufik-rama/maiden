package config

// PostgreSQLDestination ...
var PostgreSQLDestination = "postgres://postgres@localhost:5432"

// PostgreSQL ...
type PostgreSQL struct {
	Sources     PostgreSQLSources
	Destination string
}

// IsSet checks nil value
func (p *PostgreSQL) IsSet() bool {
	return p != nil
}

func (p *PostgreSQL) from(cfg fixtureWrapper) *PostgreSQL {
	p.Sources = make(PostgreSQLSources)
	for _, source := range cfg.Fixtures.PostgreSQL.Sources {
		key := source.Database
		v := p.Sources[key]
		v.Database = source.Database
		v.Definition = source.Definition
		v.Files = source.Files
		p.Sources[key] = v
	}
	p.Destination = cfg.Fixtures.PostgreSQL.Destination
	return p
}

func (p *PostgreSQL) defaultValue() {
	if emptyString(p.Destination) {
		p.Destination = PostgreSQLDestination
	}
}

func (p *PostgreSQL) resolve(dir string) {
	for key := range p.Sources {
		if !emptyString(p.Sources[key].Definition) {
			v := p.Sources[key]
			v.Definition = dir + p.Sources[key].Definition
			p.Sources[key] = v
		}
		if !emptyString(p.Sources[key].Files) {
			v := p.Sources[key]
			v.Files = dir + p.Sources[key].Files
			p.Sources[key] = v
		}
	}
}

func (p PostgreSQL) replace(other *PostgreSQL) {
	for key, val := range p.Sources {
		other.Sources[key] = val
	}
	if !emptyString(p.Destination) {
		other.Destination = p.Destination
	}
}

// PostgreSQLSources ...
type PostgreSQLSources map[string]PostgreSQLSource

// PostgreSQLSource ...
type PostgreSQLSource struct {
	Database   string
	Definition string
	Files      string
}
