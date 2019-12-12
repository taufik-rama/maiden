package structure

// FixtureStructure ...
type FixtureStructure struct {
	Elasticsearch ElasticsearchFixture
}

// ElasticsearchFixture ...
type ElasticsearchFixture struct {
	Indices []ElasticsearchIndex
}

// ElasticsearchIndex ...
type ElasticsearchIndex struct {
	Name      string
	Mapping   string
	Documents []string
}
