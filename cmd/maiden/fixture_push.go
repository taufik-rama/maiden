package main

import (
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
	"github.com/taufik-rama/maiden/storage.v1"
)

// FixturePush command handler
type FixturePush struct{}

// Command returns `push` command process
func (f FixturePush) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push the fixtures data",
		Run:   f.RunCommand,
	}
	cmd.AddCommand(f.CassandraCommand())
	cmd.AddCommand(f.ElasticsearchCommand())
	cmd.AddCommand(f.PostgreSQLCommand())
	cmd.AddCommand(f.RedisCommand())
	return cmd
}

// RunCommand runs `push` command
func (f FixturePush) RunCommand(cmd *cobra.Command, args []string) {

	fixture := new(config.Fixture)
	if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
		panic(err)
	}

	if fixture.Cassandra != nil {
		internal.Print("Pushing Cassandra")
		(storage.Cassandra{}).Push(fixture.Cassandra)
	}

	if fixture.Elasticsearch != nil {
		internal.Print("Pushing Elasticsearch")
		(storage.Elasticsearch{}).Push(fixture.Elasticsearch)
	}

	if fixture.PostgreSQL != nil {
		internal.Print("Pushing PostgreSQL")
		(storage.PostgreSQL{}).Push(fixture.PostgreSQL)
	}

	if fixture.Redis != nil {
		internal.Print("Pushing Redis")
		(storage.Redis{}).Push(fixture.Redis)
	}
}

// CassandraCommand returns `fixture push cassandra` command process
func (f FixturePush) CassandraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cassandra",
		Aliases: []string{"ca"},
		Short:   "Push the Cassandra fixtures data",
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

		// If nothing is passed, push all the data
		if keyspace == "" && table == "" {
			internal.Print("Pushing all Cassandra data")
			c.Push(fixture.Cassandra)
			return
		}

		// If only keyspace is passed, push all the table for that keyspace(s)
		if keyspace != "" && table == "" {
			internal.Print("Pushing keyspace(s) `%s`", keyspace)
			c.PushKeyspace(fixture.Cassandra, keyspace)
			return
		}

		// Specific keyspace & table
		if keyspace != "" && table != "" {
			internal.Print("Pushing keyspace `%s` table(s) `%s`", keyspace, table)
			c.PushKeyspaceTable(fixture.Cassandra, keyspace, table)
			return
		}
	}
	return cmd
}

// ElasticsearchCommand returns `fixture push elasticsearch` command process
func (f FixturePush) ElasticsearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elasticsearch",
		Aliases: []string{"es"},
		Short:   "Push the Elasticsearch fixtures data",
	}
	var index string
	cmd.PersistentFlags().StringVarP(&index, "index", "i", "", "What index/indices to use, comma separated")
	cmd.Run = func(*cobra.Command, []string) {

		fixture := new(config.Fixture)
		if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
			panic(err)
		}

		e := storage.Elasticsearch{}

		// If nothing is passed, push all the data
		if index == "" {
			e.Push(fixture.Elasticsearch)
			return
		}

		// Push specified index/indices
		e.PushIndex(fixture.Elasticsearch, index)
	}
	return cmd
}

// PostgreSQLCommand returns `fixture push postgresql` command process
func (f FixturePush) PostgreSQLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresql",
		Aliases: []string{"pg"},
		Short:   "Push the PostgreSQL fixtures data",
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

		// If nothing is passed, push all the data
		if database == "" && table == "" {
			p.Push(fixture.PostgreSQL)
			return
		}

		// If only database is passed, push all the table for that database(s)
		if database != "" && table == "" {
			p.PushDatabase(fixture.PostgreSQL, database)
			return
		}

		// Specific database & table
		p.PushDatabaseTable(fixture.PostgreSQL, database, table)
	}
	return cmd
}

// RedisCommand returns `fixture push redis` command process
func (f FixturePush) RedisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "redis",
		Aliases: []string{"re"},
		Short:   "Push the Redis fixtures data",
		Run: func(*cobra.Command, []string) {
			fixture := new(config.Fixture)
			if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
				panic(err)
			}
			(storage.Redis{}).Push(fixture.Redis)
		},
	}
	return cmd
}
