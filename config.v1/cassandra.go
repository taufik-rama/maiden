package config

// CassandraDestination ...
var CassandraDestination = "localhost:9042"

// Cassandra ...
type Cassandra struct {
	Sources     CassandraSources
	Destination string
}

// IsSet checks nil value
func (c *Cassandra) IsSet() bool {
	return c != nil
}

func (c *Cassandra) from(cfg fixtureWrapper) *Cassandra {
	c.Sources = make(CassandraSources)
	for _, source := range cfg.Fixtures.Cassandra.Sources {
		key := source.Keyspace
		v := c.Sources[key]
		v.Keyspace = source.Keyspace
		v.Definition = source.Definition
		v.Files = source.Files
		c.Sources[key] = v
	}
	c.Destination = cfg.Fixtures.Cassandra.Destination
	return c
}

func (c *Cassandra) defaultValue() {
	if emptyString(c.Destination) {
		c.Destination = CassandraDestination
	}
}

func (c *Cassandra) resolve(dir string) {
	for key := range c.Sources {
		if !emptyString(c.Sources[key].Definition) {
			v := c.Sources[key]
			v.Definition = dir + c.Sources[key].Definition
			c.Sources[key] = v
		}
		if !emptyString(c.Sources[key].Files) {
			v := c.Sources[key]
			v.Files = dir + c.Sources[key].Files
			c.Sources[key] = v
		}
	}
}

func (c Cassandra) replace(other *Cassandra) {
	for key, val := range c.Sources {
		other.Sources[key] = val
	}
	if emptyString(c.Destination) {
		other.Destination = c.Destination
	}
}

// CassandraSources ...
type CassandraSources map[string]CassandraSource

// CassandraSource ...
type CassandraSource struct {
	Keyspace   string
	Definition string
	Files      string
}
