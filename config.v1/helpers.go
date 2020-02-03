package config

import (
	"strings"

	"github.com/taufik-rama/maiden/internal"
	"gopkg.in/yaml.v2"
)

// Casts `v` into string slice, empty slice on invalid value
func toStringSlice(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	if _, ok := v.(string); ok {
		return []string{v.(string)}
	} else if vals, ok := v.([]interface{}); ok {
		slice := make([]string, len(vals))
		for i, val := range vals {
			if val, ok := val.(string); ok {
				slice[i] = val
			}
		}
		return slice
	}
	return []string{}
}

// Parse `filename` with YAML encoding
func parseYAML(filename string, v interface{}) error {
	contents, err := internal.Read(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(contents, v)
}

// Checks for trimmed empty string
func emptyString(s string) bool {
	return strings.TrimSpace(string(s)) == ""
}
