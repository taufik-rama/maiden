package tests

import (
	"testing"

	"github.com/taufik-rama/maiden/config.v1"
)

func configv1Cassandra(t *testing.T, want, got *config.Cassandra) {
	for key, val := range want.Sources {
		if got.Sources[key].Keyspace != val.Keyspace {
			t.Error()
		}
		if got.Sources[key].Definition != val.Definition {
			t.Error()
		}
		if got.Sources[key].Files != val.Files {
			t.Error()
		}
	}
	if want.Destination != got.Destination {
		t.Error()
	}
}

func configv1Elasticsearch(t *testing.T, want, got *config.Elasticsearch) {
	for key, val := range want.Sources {
		if got.Sources[key].Index != val.Index {
			t.Error()
		}
		if got.Sources[key].Mapping != val.Mapping {
			t.Error()
		}
		if got.Sources[key].Files != val.Files {
			t.Error()
		}
	}
	if want.Destination != got.Destination {
		t.Error()
	}
}

func configv1PostgreSQL(t *testing.T, want, got *config.PostgreSQL) {
	for key, val := range want.Sources {
		if got.Sources[key].Database != val.Database {
			t.Error()
		}
		if got.Sources[key].Definition != val.Definition {
			t.Error()
		}
		if got.Sources[key].Files != val.Files {
			t.Error()
		}
	}
	if want.Destination != got.Destination {
		t.Error()
	}
}

func configv1Redis(t *testing.T, want, got *config.Redis) {
	if want.Source != got.Source {
		t.Error()
	}
	if want.Destination != got.Destination {
		t.Error()
	}
}

func TestConfigV1(t *testing.T) {

	fixture := new(config.Fixture)
	if err := fixture.Parse("testdata/config.v1/config.yaml"); err != nil {
		t.Error(err)
	}

	{
		want := &config.Cassandra{
			Sources: config.CassandraSources{
				"keyspace_a": config.CassandraSource{
					Keyspace:   "keyspace_a",
					Definition: "testdata/config.v1/definition",
					Files:      "testdata/config.v1/files",
				},
			},
			Destination: "destination",
		}
		configv1Cassandra(t, want, fixture.Cassandra)
	}

	{
		want := &config.Elasticsearch{
			Sources: config.ElasticsearchSources{
				"index_a": config.ElasticsearchSource{
					Index:   "index_a",
					Mapping: "testdata/config.v1/mapping",
					Files:   "testdata/config.v1/files",
				},
			},
			Destination: "destination",
		}
		configv1Elasticsearch(t, want, fixture.Elasticsearch)
	}

	{
		want := &config.PostgreSQL{
			Sources: config.PostgreSQLSources{
				"database_a": config.PostgreSQLSource{
					Database:   "database_a",
					Definition: "testdata/config.v1/definition",
					Files:      "testdata/config.v1/files",
				},
			},
			Destination: "destination",
		}
		configv1PostgreSQL(t, want, fixture.PostgreSQL)
	}

	{
		want := &config.Redis{
			Source:      "testdata/config.v1/source",
			Destination: "destination",
		}
		configv1Redis(t, want, fixture.Redis)
	}
}
