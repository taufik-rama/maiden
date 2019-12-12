package redis

import (
	"log"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config"
	"github.com/taufik-rama/maiden/config/fixtures"
	"github.com/xuyu/goredis"
)

// RedisCommand is the command name
var RedisCommand = "redis"

// Handler for `fixture push redis` command
type Handler struct {

	// use this before any sort of print log
	verbose bool
}

// Command returns `fixture push redis` command process
func (c *Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:     RedisCommand,
		Aliases: []string{"r"},
		Short:   "Push the Redis fixtures data",
		Run:     c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

// RunCommand runs `fixture push redis` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {

	cfg, err := config.Configure()
	if err != nil {
		log.Fatalln(err)
	}

	fxt, err := fixtures.Configure()
	if err != nil {
		log.Fatalln(err)
	}

	redisClient, err := goredis.Dial(&goredis.DialConfig{
		Address: cfg.Fixtures().Redis().PushTo(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	redisTTL := 24 * 60 * 60

keyvaluesloop:
	for key, val := range fxt.Redis().KeyValues {

		c.printf("Pushing key `%s`", key)

		kind := reflect.TypeOf(val).Kind()

		if kind == reflect.String {

			if err := redisClient.Set(key, val.(string), redisTTL, 0, false, true); err != nil {
				c.printf("error on key `%s`: %s", key, err)
			}
			continue

		} else if isNumber(kind) {

			val := strconv.FormatFloat(val.(float64), 'f', -1, 64)
			if err := redisClient.Set(key, val, redisTTL, 0, false, true); err != nil {
				c.printf("error on key `%s`: %s", key, err)
			}
			continue

		} else if kind == reflect.Slice {

			values := val.([]interface{})
			parsed := make([]string, len(values))
			for index, val := range values {
				kind := reflect.TypeOf(val).Kind()
				if kind == reflect.String {
					parsed[index] = val.(string)
				} else if isNumber(kind) {
					val := strconv.FormatFloat(val.(float64), 'f', -1, 64)
					parsed[index] = val
				} else {
					c.printf("Values inside key `%s` should only be of string / int type", key)
					continue keyvaluesloop
				}
			}
			if _, err := redisClient.LPush(key, parsed...); err != nil {
				c.printf("error on key `%s`: %s", key, err)
			}
			continue

		} else if kind == reflect.Map {

			keyType, val, skip := c.check(key, val)
			if skip {
				continue
			}

			if keyType == "basic" {

				kind := reflect.TypeOf(val).Kind()

				if kind == reflect.String {
					if err := redisClient.Set(key, val.(string), redisTTL, 0, false, true); err != nil {
						c.printf("error on key `%s`: %s", key, err)
					}
					continue
				} else if isNumber(kind) {
					val := strconv.FormatFloat(val.(float64), 'f', -1, 64)
					if err := redisClient.Set(key, val, redisTTL, 0, false, true); err != nil {
						c.printf("error on key `%s`: %s", key, err)
					}
					continue
				} else if kind == reflect.Slice {
					values := val.([]interface{})
					parsed := make([]string, len(values))
					for index, val := range values {
						kind := reflect.TypeOf(val).Kind()
						if kind == reflect.String {
							parsed[index] = val.(string)
						} else if isNumber(kind) {
							val := strconv.FormatFloat(val.(float64), 'f', -1, 64)
							parsed[index] = val
						} else {
							c.printf("Values inside key `%s` should only be of string / int type", key)
							continue keyvaluesloop
						}
					}
					if _, err := redisClient.LPush(key, parsed...); err != nil {
						c.printf("error on key `%s`: %s", key, err)
					}
					continue
				} else {
					c.printf("Type `basic` must have string / number / array type for its `value` for key `%s`", key)
					continue
				}

			} else if keyType == "hash" {

				keyvals := make(map[string]string)

				kind := reflect.TypeOf(val).Kind()
				if kind != reflect.Map {
					c.printf("Type `hash` must have object type for its `value` for key `%s`", key)
					continue
				}

				val := val.(map[string]interface{})
				for prop := range val {
					kind := reflect.TypeOf(val[prop]).Kind()
					if kind == reflect.String {
						keyvals[prop] = val[prop].(string)
					} else if isNumber(kind) {
						val := strconv.FormatFloat(val[prop].(float64), 'f', -1, 64)
						keyvals[prop] = val
					} else {
						c.printf("Currently, `hash` only supports string & number for `value` fields for key `%s`", key)
						continue keyvaluesloop
					}
				}

				if err := redisClient.HMSet(key, keyvals); err != nil {
					c.printf("error on key `%s`: %s", key, err)
				}
				continue

			} else if keyType == "set" {

				kind := reflect.TypeOf(val).Kind()
				if kind != reflect.Slice {
					c.printf("Type `set` must have array type for its `value` for key `%s`", key)
					continue
				}

				values := val.([]interface{})
				parsed := make([]string, len(values))
				for index, val := range values {
					kind := reflect.TypeOf(val).Kind()
					if kind == reflect.String {
						parsed[index] = val.(string)
					} else if isNumber(kind) {
						val := strconv.FormatFloat(val.(float64), 'f', -1, 64)
						parsed[index] = val
					} else {
						c.printf("Values inside key `%s` should only be of string / int type", key)
						continue keyvaluesloop
					}
				}

				if _, err := redisClient.SAdd(key, parsed...); err != nil {
					c.printf("error on key `%s`: %s", key, err)
				}
				continue

			} else if keyType == "sorted-set" {

				keyvals := make(map[string]float64)

				kind := reflect.TypeOf(val).Kind()
				if kind != reflect.Map {
					c.printf("Type `sorted-set` must have object type for its `value` for key `%s`", key)
					continue
				}

				val := val.(map[string]interface{})
				for prop := range val {
					score, ok := val[prop].(float64)
					if !ok {
						c.printf("Type `sorted-set` must have a `string -> number` key-value for key `%s`", key)
						continue keyvaluesloop
					}
					keyvals[prop] = score
				}

				if _, err := redisClient.ZAdd(key, keyvals); err != nil {
					c.printf("error on key `%s`: %s", key, err)
				}
				continue

			} else {
				c.printf("Unknown `type` `%s` for key `%s`", keyType, key)
				continue
			}

		} else {
			c.printf("Key `%s` has an unknown type `%s`", key, kind)
		}
	}
}

// SetVerbose unimplemeted
func (c *Handler) SetVerbose(v bool) {
	c.verbose = v
}

// Print according to the verbosity flag
func (c *Handler) printf(format string, v ...interface{}) {
	if c.verbose {
		log.Printf(format, v...)
	}
}

func (c Handler) check(key string, v interface{}) (string, interface{}, bool) {
	val := v.(map[string]interface{})
	keyType, ok := val["type"]
	if !ok {
		c.printf("Key `%s` does not have a `type`", key)
		return "", nil, true
	}
	if _, isstr := keyType.(string); !isstr {
		c.printf("Key `%s` must have a string value for `type`", key)
		return "", nil, true
	}
	value, ok := val["value"]
	if !ok {
		c.printf("Key `%s` does not have a `value`", key)
		return "", nil, true
	}
	return keyType.(string), value, false
}

func isNumber(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Float32 ||
		kind == reflect.Float64
}
