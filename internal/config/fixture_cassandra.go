package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FixtureCassandra ...
type FixtureCassandra struct {

	// Source should be one of: array of object
	Source interface{} `yaml:"src"`

	// Destination should be one of: string, array of string
	Destination interface{} `yaml:"dest"`

	ParsedSource []struct {
		Definition string
		DBName     string
		SQLs       []string
	} `yaml:"-"`
	ParsedDestination []string `yaml:"-"`
}

// Name refer to the string fixtures handler
func (f FixtureCassandra) Name() string {
	return "cassandra"
}

func (f *FixtureCassandra) resolve(filename string) error {
	if err := f.check(filename); err != nil {
		return err
	}
	f.parse(filename)
	return nil
}

// check the cassandra fixture configs
func (f FixtureCassandra) check(filename string) error {

	if f.Source != nil {
		if err := f.checkSource(filename); err != nil {
			return err
		}
	}

	if f.Destination != nil {
		if err := f.checkDestination(filename); err != nil {
			return err
		}
	}

	return nil
}

func (f FixtureCassandra) checkSource(filename string) error {

	sources, ok := f.Source.([]interface{})
	if !ok {
		return fmt.Errorf("invalid Cassandra fixture source value on `%s`", filename)
	}

	for index, source := range sources {

		castedSource, ok := source.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("invalid Cassandra fixture source value on `%s` index %d", filename, index)
		}

		// Special case `definition` is used
		var definition string
		if _, ok := castedSource[".definition"]; ok {
			definition = ".definition"
		} else {
			definition = "definition"
		}

		// only 1 other key can be defined when using `definition` method
		if val, ok := castedSource[definition]; ok {
			if len(castedSource) > 2 {
				return fmt.Errorf("ambiguous Cassandra source name on `%s` index %d", filename, index)
			} else if _, ok := val.(string); !ok {
				return fmt.Errorf("source definition for Cassandra is not a valid string on `%s` index %d", filename, index)
			}
		} else if len(castedSource) > 1 {

			// If the definition is not defined, there can only be 1 key
			return fmt.Errorf("ambiguous Cassandra source name on `%s` index %d", filename, index)
		}

		for key, value := range castedSource {
			if key == definition {
				continue
			}
			if _, ok := key.(string); !ok {
				return fmt.Errorf("invalid Cassandra fixture source key on `%s` index %d", filename, index)
			}
			if value == nil {
				return fmt.Errorf("invalid Cassandra fixture source nil value on `%s` index %d", filename, index)
			}
			if _, ok := value.(string); !ok {
				return fmt.Errorf("invalid Cassandra fixture source value on `%s` index %d (must be a string)", filename, index)
			}
		}
	}

	return nil
}

func (f FixtureCassandra) checkDestination(filename string) error {
	if err := basicStringTypeCheck(f.Destination); err != nil {
		return fmt.Errorf("%s Cassandra fixture destination: %s", filename, err)
	}
	return nil
}

// parse the config into usable data types
func (f *FixtureCassandra) parse(filename string) {

	if f.Source != nil {
		f.parseSource(filename)
	}

	if f.Destination != nil {
		f.parseDestination()
	}
}

func (f *FixtureCassandra) parseSource(filename string) {

	dir := filepath.Dir(filename) + string(os.PathSeparator)
	sources := f.Source.([]interface{})
	f.ParsedSource = make([]struct {
		Definition string
		DBName     string
		SQLs       []string
	}, len(sources))

	for index, source := range sources {

		// The individual fixture object
		castedSource := source.(map[interface{}]interface{})

		// The resulting object that will be used
		parsedSource := struct {
			Definition string
			DBName     string
			SQLs       []string
		}{}

		// Special case `definition` is used
		if val, ok := castedSource[".definition"]; ok {
			parsedSource.Definition = dir + val.(string)
			delete(castedSource, ".definition")
		} else if val, ok := castedSource["definition"]; ok {
			parsedSource.Definition = dir + val.(string)
			delete(castedSource, "definition")
		}

		for key, val := range castedSource {
			parsedSource.DBName = key.(string)
			if value, ok := val.(string); ok {
				parsedSource.SQLs = []string{dir + value}
			} else {
				values := val.([]interface{})
				sqls := make([]string, len(values))
				for index, val := range values {
					sqls[index] = dir + val.(string)
				}
				parsedSource.SQLs = sqls
			}
		}

		f.ParsedSource[index] = parsedSource
	}
}

func (f *FixtureCassandra) parseDestination() {
	if destination, ok := f.Destination.(string); ok {
		if !strings.HasSuffix(destination, "/") {
			destination += "/"
		}
		f.ParsedDestination = []string{destination}
	} else {
		destinations := f.Destination.([]interface{})
		parsedDestination := make([]string, len(destinations))
		for index, destination := range destinations {
			destination := destination.(string)
			if !strings.HasSuffix(destination, "/") {
				destination += "/"
			}
			parsedDestination[index] = destination
		}
		f.ParsedDestination = parsedDestination
	}
}
