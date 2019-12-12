package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

// ServiceConfig correspond with `service` field on main config file
type ServiceConfig struct {

	// Imports should be one of: string, array of string
	Imports interface{} `yaml:"imports"`

	HTTP HTTPServices `yaml:"http"`
	GRPC GRPCServices `yaml:"grpc"`
}

func (s *ServiceConfig) resolve(filename string) error {

	bytes, err := read(filename)
	if err != nil {
		return err
	}

	var config MaidenConfig
	yaml.Unmarshal(bytes, &config)
	if config.Services == nil {
		return nil
	}

	s.HTTP.replace(config.Services.HTTP)

	if config.Services.Imports != nil {

		imports := config.Services.Imports

		if err := basicStringTypeCheck(imports); err != nil {
			return fmt.Errorf("%s service imports: %s", filename, err)
		}

		// Resolve the next imports statement relative to current directory
		dir := filepath.Dir(filename) + string(os.PathSeparator)

		if reflect.TypeOf(imports).Kind() == reflect.String {
			value := imports.(string)
			if err := s.resolve(dir + value); err != nil {
				return err
			}
		} else {
			for _, value := range imports.([]interface{}) {
				if err := s.resolve(dir + value.(string)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// ::Experimental:: Handle the config check ourself via `interface{}` values
// check the config value types.
// `filename` is used solely for the error log
// func (s ServiceConfig) check(filename string) error {

// 	services, ok := s.HTTP.(map[interface{}]interface{})
// 	if !ok {
// 		return fmt.Errorf("invalid HTTP service config on `%s`", filename)
// 	}

// 	for name, service := range services {

// 		details, ok := service.(map[interface{}]interface{})
// 		if !ok {
// 			return fmt.Errorf("invalid `%s` HTTP service config details on `%s`", name, filename)
// 		}

// 		if port, ok := details["port"]; ok {
// 			if _, ok := port.(uint); !ok {
// 				return fmt.Errorf("invalid `%s` HTTP service port config details on `%s`, port have to be an unsigned int", name, filename)
// 			}
// 		}

// 		if endpoints, ok := details["endpoints"]; ok {
// 			endpoints, ok := endpoints.(map[interface{}]interface{})
// 			if !ok {
// 				return fmt.Errorf("invalid `%s` HTTP service endpoints config details on `%s`", name, filename)
// 			}

// 			for url, details := range endpoints {

// 				if _, ok := url.(string); !ok {
// 					return fmt.Errorf("invalid `%s` HTTP service URL on `%s`, endpoints must be a string", name, filename)
// 				}

// 				details, ok := details.([]interface{})
// 				if !ok {
// 					return fmt.Errorf("invalid `%s`:`%s` HTTP service endpoints config details value on `%s`", name, url, filename)
// 				}

// 				for index, detail := range details {

// 					detail, ok := detail.(map[interface{}]interface{})
// 					if !ok {
// 						return fmt.Errorf("invalid `%s`:`%s`(`%d`) HTTP service endpoints config details value on `%s`", name, url, index, filename)
// 					}

// 					if method, ok := detail["method"]; ok {
// 						if err := basicStringTypeCheck(method); err != nil {
// 							return fmt.Errorf("invalid `%s`:`%s`(`%d`) HTTP service method endpoints config detail on `%s`", name, url, index, filename)
// 						}
// 					}

// 					if _, ok := endpoints["request"]; ok {

// 					}

// 					if response, ok := endpoints["response"]; ok {
// 						if response, ok := response.(map[interface{}]interface{}); ok {
// 							if contentType, ok := response["content-type"]; ok {
// 								if _, ok := contentType.(string); !ok {
// 									return fmt.Errorf("invalid `%s`:`%s`(`%d`) HTTP service response detail on `%s`", name, url, index, filename)
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }
