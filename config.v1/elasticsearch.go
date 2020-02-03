package config

// ElasticsearchDestination ...
var ElasticsearchDestination = "http://localhost:9200"

// Elasticsearch ...
type Elasticsearch struct {
	Sources     ElasticsearchSources
	Destination string
}

// IsSet checks nil value
func (e *Elasticsearch) IsSet() bool {
	return e != nil
}

func (e *Elasticsearch) from(cfg fixtureWrapper) *Elasticsearch {
	e.Sources = make(ElasticsearchSources)
	for _, source := range cfg.Fixtures.Elasticsearch.Sources {
		key := source.Index
		v := e.Sources[key]
		v.Index = source.Index
		v.Mapping = source.Mapping
		v.Files = source.Files
		e.Sources[key] = v
	}
	e.Destination = cfg.Fixtures.Elasticsearch.Destination
	return e
}

func (e *Elasticsearch) defaultValue() {
	if emptyString(e.Destination) {
		e.Destination = ElasticsearchDestination
	}
}

func (e *Elasticsearch) resolve(dir string) {
	for key := range e.Sources {
		if !emptyString(e.Sources[key].Mapping) {
			v := e.Sources[key]
			v.Mapping = dir + e.Sources[key].Mapping
			e.Sources[key] = v
		}
		if !emptyString(e.Sources[key].Files) {
			v := e.Sources[key]
			v.Files = dir + e.Sources[key].Files
			e.Sources[key] = v
		}
	}
}

func (e Elasticsearch) replace(other *Elasticsearch) {
	for key, val := range e.Sources {
		other.Sources[key] = val
	}
	if !emptyString(e.Destination) {
		other.Destination = e.Destination
	}
}

// ElasticsearchSources ...
type ElasticsearchSources map[string]ElasticsearchSource

// ElasticsearchSource ...
type ElasticsearchSource struct {
	Index   string
	Mapping string
	Files   string
}
