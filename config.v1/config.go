package config

// InputFixture corresponds with the data structure of unmarshalled config values.
// It shouldn't be used directly and instead marshalled via it's parse method.
type InputFixture struct {
	Imports   interface{} `yaml:"imports"`
	Cassandra *struct {
		Sources []struct {
			Keyspace   string `yaml:"keyspace"`
			Definition string `yaml:"definition"`
			Files      string `yaml:"files"`
		} `yaml:"src"`
		Destination string `yaml:"dest"`
	}
	Elasticsearch *struct {
		Sources []struct {
			Index       string `yaml:"index"`
			Mapping     string `yaml:"mapping"`
			MappingType string `yaml:"mapping-type"`
			Files       string `yaml:"files"`
		} `yaml:"src"`
		Destination string `yaml:"dest"`
	} `yaml:"elasticsearch"`
	PostgreSQL *struct {
		Sources []struct {
			Database   string `yaml:"database"`
			Definition string `yaml:"definition"`
			Files      string `yaml:"files"`
		} `yaml:"src"`
		Destination string `yaml:"dest"`
	} `yaml:"postgresql"`
	Redis *struct {
		Source      string `yaml:"src"`
		Destination string `yaml:"dest"`
	} `yaml:"redis"`
}

// InputService corresponds with the data structure of unmarshalled config values.
// It shouldn't be used directly and instead marshalled via it's parse method.
type InputService struct {
	Imports interface{} `yaml:"imports"`
	Output  string      `yaml:"output"`
	GRPC    map[string]struct {
		Port       uint16 `yaml:"port"`
		Definition string `yaml:"definition"`
		Methods    map[string]struct {
			Request  string `yaml:"request"`
			Response string `yaml:"response"`
		} `yaml:"methods"`
		Conditions map[string][]struct {
			Request  interface{} `yaml:"request"`
			Response interface{} `yaml:"response"`
		} `yaml:"conditions"`
	} `yaml:"grpc"`
	HTTP map[string]struct {
		Port      uint16 `yaml:"port"`
		Endpoints map[string][]struct {
			Method   interface{} `yaml:"method"`
			Request  interface{} `yaml:"request"`
			Response interface{} `yaml:"response"`
		} `yaml:"endpoints"`
	} `yaml:"http"`
}
