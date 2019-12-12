package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FixtureRedis ...
type FixtureRedis struct {

	// Source should be one of: string, array of string
	Source       interface{} `yaml:"src"`
	ParsedSource []string    `yaml:"-"`

	// Destination should be one of: string, array of string
	Destination       interface{} `yaml:"dest"`
	ParsedDestination []string    `yaml:"-"`
}

// Name refer to the string fixtures handler
func (f FixtureRedis) Name() string {
	return "redis"
}

func (f *FixtureRedis) resolve(filename string) error {
	if err := f.check(filename); err != nil {
		return err
	}
	f.parse(filename)
	return nil
}

// check the redis fixture configs
func (f FixtureRedis) check(filename string) error {

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

func (f FixtureRedis) checkSource(filename string) error {
	if err := basicStringTypeCheck(f.Source); err != nil {
		return fmt.Errorf("%s Redis fixture source: %s", filename, err)
	}
	return nil
}

func (f FixtureRedis) checkDestination(filename string) error {
	if err := basicStringTypeCheck(f.Destination); err != nil {
		return fmt.Errorf("%s Redis fixture destination: %s", filename, err)
	}
	return nil
}

// parse the config into usable data types
func (f *FixtureRedis) parse(filename string) {

	if f.Source != nil {
		f.parseSource(filename)
	}

	if f.Destination != nil {
		f.parseDestination()
	}
}

func (f *FixtureRedis) parseSource(filename string) {
	dir := filepath.Dir(filename) + string(os.PathSeparator)
	if source, ok := f.Source.(string); ok {
		f.ParsedSource = []string{dir + source}
	} else {
		sources := f.Source.([]interface{})
		parsedSource := make([]string, len(sources))
		for index, source := range sources {
			parsedSource[index] = dir + source.(string)
		}
		f.ParsedSource = parsedSource
	}
}

func (f *FixtureRedis) parseDestination() {
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
