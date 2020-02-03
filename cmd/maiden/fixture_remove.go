package main

import (
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
	"github.com/taufik-rama/maiden/storage.v1"
)

// FixtureRemove command handler
type FixtureRemove struct{}

// Command returns `fixture remove` command process
func (f FixtureRemove) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove the fixtures data",
		Run:     f.RunCommand,
	}
	cmd.AddCommand(f.CassandraCommand())
	cmd.AddCommand(f.ElasticsearchCommand())
	cmd.AddCommand(f.PostgreSQLCommand())
	cmd.AddCommand(f.RedisCommand())
	return cmd
}

// RunCommand runs `fixture remove` command
func (f FixtureRemove) RunCommand(cmd *cobra.Command, args []string) {

	fixture := new(config.Fixture)
	if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
		panic(err)
	}

	internal.Print("Removing Cassandra")
	(storage.Cassandra{}).Remove(fixture.Cassandra)

	internal.Print("Removing Elasticsearch")
	(storage.Elasticsearch{}).Remove(fixture.Elasticsearch)

	internal.Print("Removing PostgreSQL")
	(storage.PostgreSQL{}).Remove(fixture.PostgreSQL)

	internal.Print("Removing Redis")
	(storage.Redis{}).Remove(fixture.Redis)
}

// CassandraCommand returns `fixture remove cassandra` command process
func (f FixtureRemove) CassandraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cassandra",
		Aliases: []string{"ca"},
		Short:   "Remove the Cassandra fixtures data",
	}
	var keyspace, table string
	cmd.PersistentFlags().StringVarP(&keyspace, "keyspace", "d", "", "What keyspace(s) to use")
	cmd.PersistentFlags().StringVarP(&table, "table", "t", "", "What table(s) to use, comma separated")
	cmd.Run = func(*cobra.Command, []string) {

		fixture := new(config.Fixture)
		if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
			panic(err)
		}

		c := storage.Cassandra{}

		// If nothing is passed, remove all the data
		if keyspace == "" && table == "" {
			c.Remove(fixture.Cassandra)
			return
		}

		// If only keyspace is passed, remove all the table for that keyspace(s)
		if keyspace != "" && table == "" {
			c.RemoveKeyspace(fixture.Cassandra, keyspace)
			return
		}

		// Specific keyspace & table
		if keyspace != "" && table != "" {
			c.RemoveKeyspaceTable(fixture.Cassandra, keyspace, table)
			return
		}
	}
	return cmd
}

// ElasticsearchCommand returns `fixture remove elasticsearch` command process
func (f FixtureRemove) ElasticsearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elasticsearch",
		Aliases: []string{"es"},
		Short:   "Remove the Elasticsearch fixtures data",
	}
	var index string
	cmd.PersistentFlags().StringVarP(&index, "index", "i", "", "What index/indices to use, comma separated")
	cmd.Run = func(*cobra.Command, []string) {

		fixture := new(config.Fixture)
		if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
			panic(err)
		}

		e := storage.Elasticsearch{}

		// If nothing is passed, remove all the data
		if index == "" {
			e.Remove(fixture.Elasticsearch)
			return
		}

		// Remove specified index/indices
		e.RemoveIndex(fixture.Elasticsearch, index)
	}
	return cmd
}

// PostgreSQLCommand returns `fixture remove postgresql` command process
func (f FixtureRemove) PostgreSQLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresql",
		Aliases: []string{"pg"},
		Short:   "Remove the PostgreSQL fixtures data",
	}
	var database, table string
	cmd.PersistentFlags().StringVarP(&database, "database", "d", "", "What database(s) to use, comma separated")
	cmd.PersistentFlags().StringVarP(&table, "table", "t", "", "What table(s) to use, comma separated")
	cmd.Run = func(*cobra.Command, []string) {

		fixture := new(config.Fixture)
		if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
			panic(err)
		}

		p := storage.PostgreSQL{}

		// If nothing is passed, remove all the data
		if database == "" && table == "" {
			p.Remove(fixture.PostgreSQL)
			return
		}

		// If only database is passed, remove all the table for that database(s)
		if database != "" && table == "" {
			p.RemoveDatabase(fixture.PostgreSQL, database)
			return
		}

		// Specific database & table
		if database != "" && table != "" {
			p.RemoveDatabaseTable(fixture.PostgreSQL, database, table)
			return
		}
	}
	return cmd
}

// RedisCommand returns `fixture remove redis` command process
func (f FixtureRemove) RedisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "redis",
		Aliases: []string{"re"},
		Short:   "Remove the Redis fixtures data",
		Run: func(*cobra.Command, []string) {
			fixture := new(config.Fixture)
			if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
				panic(err)
			}
			(storage.Redis{}).Remove(fixture.Redis)
		},
	}
	return cmd
}
