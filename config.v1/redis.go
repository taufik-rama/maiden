package config

// RedisDestination ...
var RedisDestination = "localhost:6379"

// Redis ...
type Redis struct {
	Source      string
	Destination string
}

// IsSet checks nil value
func (r *Redis) IsSet() bool {
	return r != nil
}

func (r *Redis) from(cfg fixtureWrapper) *Redis {
	r.Source = cfg.Fixtures.Redis.Source
	r.Destination = cfg.Fixtures.Redis.Destination
	return r
}

func (r *Redis) defaultValue() {
	if emptyString(r.Destination) {
		r.Destination = RedisDestination
	}
}

func (r *Redis) resolve(dir string) {
	if !emptyString(r.Source) {
		r.Source = dir + r.Source
	}
}

func (r Redis) replace(other *Redis) {
	if !emptyString(r.Source) {
		other.Source = r.Source
	}
	if !emptyString(r.Destination) {
		other.Destination = r.Destination
	}
}
