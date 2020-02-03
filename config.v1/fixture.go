package config

import (
	"os"
	"path/filepath"
)

// Fixture ...
type Fixture struct {
	Cassandra     *Cassandra
	Elasticsearch *Elasticsearch
	PostgreSQL    *PostgreSQL
	Redis         *Redis
}

type fixtureWrapper struct {
	Fixtures *InputFixture `yaml:"fixtures"`
}

// Parse ...
func (f *Fixture) Parse(filename string) error {

	var config fixtureWrapper
	if err := parseYAML(filename, &config); err != nil {
		return err
	}
	if config.Fixtures == nil {
		return nil
	}

	dir := filepath.Dir(filename) + string(os.PathSeparator)

	new(Fixture).from(config).defaultValue().resolve(dir).replace(f)

	for _, imp := range toStringSlice(config.Fixtures.Imports) {
		if err := f.Parse(dir + imp); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fixture) from(cfg fixtureWrapper) *Fixture {

	if cfg.Fixtures == nil {
		return f
	}

	if cfg.Fixtures.Cassandra != nil {
		f.Cassandra = new(Cassandra).from(cfg)
	}

	if cfg.Fixtures.Elasticsearch != nil {
		f.Elasticsearch = new(Elasticsearch).from(cfg)
	}

	if cfg.Fixtures.PostgreSQL != nil {
		f.PostgreSQL = new(PostgreSQL).from(cfg)
	}

	if cfg.Fixtures.Redis != nil {
		f.Redis = new(Redis).from(cfg)
	}

	return f
}

func (f *Fixture) defaultValue() *Fixture {

	if f.Cassandra != nil {
		f.Cassandra.defaultValue()
	}

	if f.Elasticsearch != nil {
		f.Elasticsearch.defaultValue()
	}

	if f.PostgreSQL != nil {
		f.PostgreSQL.defaultValue()
	}

	if f.Redis != nil {
		f.Redis.defaultValue()
	}

	return f
}

func (f *Fixture) resolve(dir string) *Fixture {

	if f.Cassandra != nil {
		f.Cassandra.resolve(dir)
	}

	if f.Elasticsearch != nil {
		f.Elasticsearch.resolve(dir)
	}

	if f.PostgreSQL != nil {
		f.PostgreSQL.resolve(dir)
	}

	if f.Redis != nil {
		f.Redis.resolve(dir)
	}

	return f
}

func (f Fixture) replace(other *Fixture) {

	if other.Cassandra == nil {
		other.Cassandra = f.Cassandra
	} else if f.Cassandra != nil {
		f.Cassandra.replace(other.Cassandra)
	}

	if other.Elasticsearch == nil {
		other.Elasticsearch = f.Elasticsearch
	} else if f.Elasticsearch != nil {
		f.Elasticsearch.replace(other.Elasticsearch)
	}

	if other.PostgreSQL == nil {
		other.PostgreSQL = f.PostgreSQL
	} else if f.PostgreSQL != nil {
		f.PostgreSQL.replace(other.PostgreSQL)
	}

	if other.Redis == nil {
		other.Redis = f.Redis
	} else if f.Redis != nil {
		f.Redis.replace(other.Redis)
	}
}
