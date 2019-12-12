package structure

import (
	"github.com/taufik-rama/maiden/internal/config"
)

// Structure ...
type Structure struct {
	Fixture FixtureStructure
}

// ResolveFixtures ...
func (s *Structure) ResolveFixtures(fixture *config.FixtureConfig) error {
	s.Fixture.Elasticsearch.Indices = make([]ElasticsearchIndex, len(fixture.Elasticsearch.ParsedSource))
	for index, src := range fixture.Elasticsearch.ParsedSource {
		files, err := readdir(src.JSONs[0])
		if err != nil {
			return err
		}
		s.Fixture.Elasticsearch.Indices[index] = ElasticsearchIndex{
			Name:      src.IndexName,
			Mapping:   src.Definition,
			Documents: files,
		}
	}
	return nil
}
