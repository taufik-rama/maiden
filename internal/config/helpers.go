package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// Print according to the verbosity flag
func Print(format string, v ...interface{}) {
	if Args.Verbose {
		log.Printf(format, v...)
	}
}

// Checks if the value meet the basic string type (a string, an array of string).
// This turns out to be more common than you think
func basicStringTypeCheck(str interface{}) error {
	kind := reflect.TypeOf(str).Kind()
	if kind != reflect.String && kind != reflect.Slice {
		return errors.New("not a string or an array of string")
	} else if kind == reflect.Slice {
		for _, value := range str.([]interface{}) {
			kind := reflect.TypeOf(value).Kind()
			if kind != reflect.String {
				return errors.New("not a string or an array of string")
			}
		}
	}
	return nil
}

// Read the file contents
func read(filename string) ([]byte, error) {
	Print("Parsing %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

// ActiveFixtures returns what fixtures are defined by the user, based on `source` value
func (f FixtureConfig) ActiveFixtures() []interface {
	Name() string
} {
	var fixtures []interface {
		Name() string
	}
	if len(f.Cassandra.ParsedSource) > 0 {
		fixtures = append(fixtures, f.Cassandra)
	}
	if len(f.Elasticsearch.ParsedSource) > 0 {
		fixtures = append(fixtures, f.Elasticsearch)
	}
	if len(f.PostgreSQL.ParsedSource) > 0 {
		fixtures = append(fixtures, f.PostgreSQL)
	}
	if len(f.Redis.ParsedSource) > 0 {
		fixtures = append(fixtures, f.Redis)
	}
	return fixtures
}
