package tests

import (
	"testing"

	"github.com/taufik-rama/maiden/internal/config"
)

func TestFixtures(t *testing.T) {

	FixturesImports(t)

	FixturesCassandra(t)

	FixturesElasticsearch(t)

	FixturesPostgreSQL(t)

	FixturesRedis(t)
}

func FixturesImports(t *testing.T) {

	{
		c, err := config.New("definitions/fixtures/imports.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveFixtures(); err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}

		if c.Fixtures.Cassandra.Source == nil || c.Fixtures.Cassandra.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.Cassandra.Source.([]interface{})
			if source[0] != "cassandra-src-1" || source[1] != "cassandra-src-2" {
				t.Error()
			}
			destination := c.Fixtures.Cassandra.Destination.([]interface{})
			if destination[0] != "cassandra-dest-1" || destination[1] != "cassandra-dest-2" {
				t.Error()
			}
		}

		if c.Fixtures.Elasticsearch.Source == nil || c.Fixtures.Elasticsearch.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.Elasticsearch.Source.([]interface{})
			if source[0] != "elasticsearch-src-1" || source[1] != "elasticsearch-src-2" {
				t.Error()
			}
			destination := c.Fixtures.Elasticsearch.Destination.([]interface{})
			if destination[0] != "elasticsearch-dest-1" || destination[1] != "elasticsearch-dest-2" {
				t.Error()
			}
		}

		if c.Fixtures.PostgreSQL.Source == nil || c.Fixtures.PostgreSQL.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.PostgreSQL.Source.([]interface{})
			if source[0] != "postgresql-src-1" || source[1] != "postgresql-src-2" {
				t.Error()
			}
			destination := c.Fixtures.PostgreSQL.Destination.([]interface{})
			if destination[0] != "postgresql-dest-1" || destination[1] != "postgresql-dest-2" {
				t.Error()
			}
		}
		if c.Fixtures.Redis.Source == nil || c.Fixtures.Redis.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.Redis.Source.([]interface{})
			if source[0] != "redis-src-1" || source[1] != "redis-src-2" {
				t.Error()
			}
			destination := c.Fixtures.Redis.Destination.([]interface{})
			if destination[0] != "redis-dest-1" || destination[1] != "redis-dest-2" {
				t.Error()
			}
		}
	}

	{
		c, err := config.New("definitions/fixtures/imports-root.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveFixtures(); err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}

		if c.Fixtures.Cassandra.Source == nil || c.Fixtures.Cassandra.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.Cassandra.Source.([]interface{})
			if source[0] != "root-cassandra-src-1" || source[1] != "root-cassandra-src-2" {
				t.Error()
			}
			destination := c.Fixtures.Cassandra.Destination.([]interface{})
			if destination[0] != "root-cassandra-dest-1" || destination[1] != "root-cassandra-dest-2" {
				t.Error()
			}
		}

		if c.Fixtures.Elasticsearch.Source == nil || c.Fixtures.Elasticsearch.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.Elasticsearch.Source.([]interface{})
			if source[0] != "root-elasticsearch-src-1" || source[1] != "root-elasticsearch-src-2" {
				t.Error()
			}
			destination := c.Fixtures.Elasticsearch.Destination.([]interface{})
			if destination[0] != "root-elasticsearch-dest-1" || destination[1] != "root-elasticsearch-dest-2" {
				t.Error()
			}
		}

		if c.Fixtures.PostgreSQL.Source == nil || c.Fixtures.PostgreSQL.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.PostgreSQL.Source.([]interface{})
			if source[0] != "root-postgresql-src-1" || source[1] != "root-postgresql-src-2" {
				t.Error()
			}
			destination := c.Fixtures.PostgreSQL.Destination.([]interface{})
			if destination[0] != "root-postgresql-dest-1" || destination[1] != "root-postgresql-dest-2" {
				t.Error()
			}
		}
		if c.Fixtures.Redis.Source == nil || c.Fixtures.Redis.Destination == nil {
			t.Error()
		}
		{
			source := c.Fixtures.Redis.Source.([]interface{})
			if source[0] != "root-redis-src-1" || source[1] != "root-redis-src-2" {
				t.Error()
			}
			destination := c.Fixtures.Redis.Destination.([]interface{})
			if destination[0] != "root-redis-dest-1" || destination[1] != "root-redis-dest-2" {
				t.Error()
			}
		}
	}
}

func FixturesCassandra(t *testing.T) {

	{
		c, err := config.New("definitions/fixtures/cassandra/array.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.Cassandra.Source == nil || c.Fixtures.Cassandra.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.Cassandra.Source.([]interface{})
		if source[0] != "cassandra-src-1" || source[1] != "cassandra-src-2" {
			t.Error()
		}
		destination := c.Fixtures.Cassandra.Destination.([]interface{})
		if destination[0] != "cassandra-dest-1" || destination[1] != "cassandra-dest-2" {
			t.Error()
		}
	}

	{
		c, err := config.New("definitions/fixtures/cassandra/string.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.Cassandra.Source == nil || c.Fixtures.Cassandra.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.Cassandra.Source.(string)
		if source != "cassandra-src" {
			t.Error()
		}
		destination := c.Fixtures.Cassandra.Destination.(string)
		if destination != "cassandra-dest" {
			t.Error()
		}
	}
}

func FixturesElasticsearch(t *testing.T) {

	{
		c, err := config.New("definitions/fixtures/elasticsearch/array.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.Elasticsearch.Source == nil || c.Fixtures.Elasticsearch.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.Elasticsearch.Source.([]interface{})
		if source[0] != "elasticsearch-src-1" || source[1] != "elasticsearch-src-2" {
			t.Error()
		}
		destination := c.Fixtures.Elasticsearch.Destination.([]interface{})
		if destination[0] != "elasticsearch-dest-1" || destination[1] != "elasticsearch-dest-2" {
			t.Error()
		}
	}

	{
		c, err := config.New("definitions/fixtures/elasticsearch/string.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.Elasticsearch.Source == nil || c.Fixtures.Elasticsearch.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.Elasticsearch.Source.(string)
		if source != "elasticsearch-src" {
			t.Error()
		}
		destination := c.Fixtures.Elasticsearch.Destination.(string)
		if destination != "elasticsearch-dest" {
			t.Error()
		}
	}
}

func FixturesPostgreSQL(t *testing.T) {

	{
		c, err := config.New("definitions/fixtures/postgresql/array.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.PostgreSQL.Source == nil || c.Fixtures.PostgreSQL.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.PostgreSQL.Source.([]interface{})
		if source[0] != "postgresql-src-1" || source[1] != "postgresql-src-2" {
			t.Error()
		}
		destination := c.Fixtures.PostgreSQL.Destination.([]interface{})
		if destination[0] != "postgresql-dest-1" || destination[1] != "postgresql-dest-2" {
			t.Error()
		}
	}

	{
		c, err := config.New("definitions/fixtures/postgresql/string.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.PostgreSQL.Source == nil || c.Fixtures.PostgreSQL.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.PostgreSQL.Source.(string)
		if source != "postgresql-src" {
			t.Error()
		}
		destination := c.Fixtures.PostgreSQL.Destination.(string)
		if destination != "postgresql-dest" {
			t.Error()
		}
	}
}

func FixturesRedis(t *testing.T) {

	{
		c, err := config.New("definitions/fixtures/redis/array.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.Redis.Source == nil || c.Fixtures.Redis.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.Redis.Source.([]interface{})
		if source[0] != "redis-src-1" || source[1] != "redis-src-2" {
			t.Error()
		}
		destination := c.Fixtures.Redis.Destination.([]interface{})
		if destination[0] != "redis-dest-1" || destination[1] != "redis-dest-2" {
			t.Error()
		}
	}

	{
		c, err := config.New("definitions/fixtures/redis/string.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures == nil || c.Services != nil {
			t.Error()
		}
		if c.Fixtures.Redis.Source == nil || c.Fixtures.Redis.Destination == nil {
			t.Error()
		}
		source := c.Fixtures.Redis.Source.(string)
		if source != "redis-src" {
			t.Error()
		}
		destination := c.Fixtures.Redis.Destination.(string)
		if destination != "redis-dest" {
			t.Error()
		}
	}
}
