package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

// GRPCServices ...
type GRPCServices map[string]GRPCService

// Resolve the services imports statement
func (g GRPCServices) resolve(rootdir string) error {

	for name := range g {

		// Take the existing parsed value because we need to
		// check if the value need to be updated or not
		service := g[name]

		if service.Imports == nil {
			continue
		}

		if err := basicStringTypeCheck(service.Imports); err != nil {
			return fmt.Errorf("%s service imports: %s", rootdir, err)
		}

		kind := reflect.TypeOf(service.Imports).Kind()
		if kind == reflect.String {
			value := service.Imports.(string)
			if err := service.readImport(name, rootdir, value); err != nil {
				return err
			}
		} else {
			for _, value := range service.Imports.([]interface{}) {
				if err := service.readImport(name, rootdir, value.(string)); err != nil {
					return err
				}
			}
		}

		// Remove the imports statements and update the config values
		service.Imports = nil
		g[name] = service
	}

	return nil
}

// GRPCService ...
type GRPCService struct {

	// Imports should be one of: string, array of string
	Imports interface{} `yaml:"imports"`

	Port       uint16         `yaml:"port"`
	Definition string         `yaml:"definition"`
	Methods    GRPCMethods    `yaml:"methods"`
	Conditions GRPCConditions `yaml:"conditions"`
}

// GRPCMethods ...
type GRPCMethods map[string]struct {
	Request  string `yaml:"request"`
	Response string `yaml:"response"`
}

// GRPCConditions ...
type GRPCConditions map[string][]struct {

	// Request should be one of:
	Request interface{} `yaml:"request"`

	// Response should be one of:
	Response interface{} `yaml:"response"`
}

// GRPCConditions ...
// type GRPCConditions map[string]interface{}

func (g *GRPCService) readImport(serviceName, rootdir, dir string) error {

	filename := rootdir + string(os.PathSeparator) + dir

	Print("Reading %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}

	config := MaidenConfig{}
	yaml.Unmarshal(contents, &config)
	if config.Services == nil {
		return nil
	}

	if err := g.check(); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}

	services := config.Services.GRPC
	for name := range services {
		if name != serviceName {
			continue
		}
		Print("Importing service %s from %s", serviceName, filename)
		g.replaceIfEmpty(services[name])
	}

	return nil
}

// check the config value types
func (g GRPCService) check() error {
	return nil
}

// Replace self field values with `other` if the field is empty (nil, 0, "", etc.).
// Needed because we prioritize the parent / self attribute
func (g *GRPCService) replaceIfEmpty(other GRPCService) {
	if g.Port == 0 {
		g.Port = other.Port
	}
	if g.Definition == "" {
		g.Definition = other.Definition
	}
	if g.Methods == nil {
		g.Methods = make(GRPCMethods)
	}
	if g.Conditions == nil {
		g.Conditions = make(GRPCConditions, 0)
	}
	for method := range other.Methods {
		if _, ok := g.Methods[method]; !ok {
			g.Methods[method] = other.Methods[method]
		} else {
			Print("Method `%s` is already defined", method)
		}
	}
	for condition := range other.Conditions {
		if _, ok := g.Conditions[condition]; !ok {
			g.Conditions[condition] = other.Conditions[condition]
		} else {
			g.Conditions[condition] = append(g.Conditions[condition], other.Conditions[condition]...)
		}
	}
}
