package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

// FixtureConfig correspond with `fixture` field on main config file
type FixtureConfig struct {

	// Imports should be one of: string, array of string
	Imports interface{} `yaml:"imports"`

	Cassandra     FixtureCassandra     `yaml:"cassandra"`
	Elasticsearch FixtureElasticsearch `yaml:"elasticsearch"`
	PostgreSQL    FixturePostgreSQL    `yaml:"postgresql"`
	Redis         FixtureRedis         `yaml:"redis"`
}

func (f *FixtureConfig) resolve(filename string) error {

	bytes, err := read(filename)
	if err != nil {
		return err
	}

	var config MaidenConfig
	yaml.Unmarshal(bytes, &config)
	if config.Fixtures == nil {
		return nil
	}

	if err := config.Fixtures.Cassandra.resolve(filename); err != nil {
		return err
	} else if err := config.Fixtures.Elasticsearch.resolve(filename); err != nil {
		return err
	} else if err := config.Fixtures.PostgreSQL.resolve(filename); err != nil {
		return err
	} else if err := config.Fixtures.Redis.resolve(filename); err != nil {
		return err
	}

	f.append(config.Fixtures)

	if config.Fixtures.Imports != nil {

		imports := config.Fixtures.Imports

		if err := basicStringTypeCheck(imports); err != nil {
			return fmt.Errorf("%s fixture imports: %s", filename, err)
		}

		// Resolve the next imports statement relative to current directory
		dir := filepath.Dir(filename) + string(os.PathSeparator)

		if reflect.TypeOf(imports).Kind() == reflect.String {
			value := imports.(string)
			if err := f.resolve(dir + value); err != nil {
				return err
			}
		} else {
			for _, value := range imports.([]interface{}) {
				if err := f.resolve(dir + value.(string)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (f *FixtureConfig) append(other *FixtureConfig) {
	f.Cassandra.ParsedSource = append(f.Cassandra.ParsedSource, other.Cassandra.ParsedSource...)
	f.Cassandra.ParsedDestination = append(f.Cassandra.ParsedDestination, other.Cassandra.ParsedDestination...)
	f.Elasticsearch.ParsedSource = append(f.Elasticsearch.ParsedSource, other.Elasticsearch.ParsedSource...)
	f.Elasticsearch.ParsedDestination = append(f.Elasticsearch.ParsedDestination, other.Elasticsearch.ParsedDestination...)
	f.PostgreSQL.ParsedSource = append(f.PostgreSQL.ParsedSource, other.PostgreSQL.ParsedSource...)
	f.PostgreSQL.ParsedDestination = append(f.PostgreSQL.ParsedDestination, other.PostgreSQL.ParsedDestination...)
	f.Redis.ParsedSource = append(f.Redis.ParsedSource, other.Redis.ParsedSource...)
	f.Redis.ParsedDestination = append(f.Redis.ParsedDestination, other.Redis.ParsedDestination...)
}
