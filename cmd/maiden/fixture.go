package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
)

// Fixture command handler
type Fixture struct {
	statusAll           bool
	statusCassandra     bool
	statusElasticsearch bool
	statusPostgreSQL    bool
	statusRedis         bool
}

// Command returns `fixture` command process
func (f Fixture) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "fixture",
		Aliases: []string{"f"},
		Short:   "Interact with the available fixtures",
		Run:     f.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&f.statusAll, "status-all", "", false, "Show all the fixtures")
	cmd.PersistentFlags().BoolVarP(&f.statusCassandra, "status-cassandra", "", false, "Show cassandra fixtures")
	cmd.PersistentFlags().BoolVarP(&f.statusElasticsearch, "status-elasticsearch", "", false, "Show elasticsearch fixtures")
	cmd.PersistentFlags().BoolVarP(&f.statusPostgreSQL, "status-postgresql", "", false, "Show postgresql fixtures")
	cmd.PersistentFlags().BoolVarP(&f.statusRedis, "status-redis", "", false, "Show redis fixtures")

	cmd.AddCommand((FixturePush{}).Command())
	cmd.AddCommand((FixtureRemove{}).Command())
	cmd.AddCommand((FixtureGenerate{}).Command())

	return cmd
}

// RunCommand runs `fixture` command
func (f *Fixture) RunCommand(cmd *cobra.Command, args []string) {

	if f.statusFlag() {
		f.printStatus()
		return
	}

	if err := new(config.Fixture).Parse(internal.Args.ConfigFile); err != nil {
		panic(err)
	}
	cmd.Println("All fixtures parsed, check with `--status-x` for the values")
}

func (f Fixture) statusFlag() bool {
	return f.statusAll || f.statusCassandra || f.statusElasticsearch || f.statusPostgreSQL || f.statusRedis
}

func (f Fixture) printStatus() {

	if f.statusAll {
		f.statusCassandra = true
		f.statusElasticsearch = true
		f.statusPostgreSQL = true
		f.statusRedis = true
	}

	fixture := new(config.Fixture)
	if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
		panic(err)
	}

	if f.statusCassandra && fixture.Cassandra != nil {
		f.printStatusCassandra(fixture.Cassandra)
	}

	if f.statusElasticsearch && fixture.Elasticsearch != nil {
		if f.statusCassandra && fixture.Cassandra != nil {
			fmt.Println()
		}
		f.printStatusElasticsearch(fixture.Elasticsearch)
	}

	if f.statusPostgreSQL && fixture.PostgreSQL != nil {
		if f.statusElasticsearch && fixture.Elasticsearch != nil {
			fmt.Println()
		}
		f.printStatusPostgreSQL(fixture.PostgreSQL)
	}

	if f.statusRedis && fixture.Redis != nil {
		if f.statusPostgreSQL && fixture.PostgreSQL != nil {
			fmt.Println()
		}
		f.printStatusRedis(fixture.Redis)
	}
}

func (f Fixture) printStatusCassandra(cfg *config.Cassandra) {
	fmt.Println("Cassandra -", cfg.Destination)
	fmt.Println("  Sources:")
	for _, source := range cfg.Sources {
		fmt.Println("    - Keyspace:", source.Keyspace)
		fmt.Println("      Definition:", source.Definition)
		fmt.Println("      Files:", source.Files)
	}
}

func (f Fixture) printStatusElasticsearch(cfg *config.Elasticsearch) {
	fmt.Println("Elasticsearch -", cfg.Destination)
	fmt.Println("  Sources:")
	for _, source := range cfg.Sources {
		fmt.Println("    - Index:", source.Index)
		fmt.Println("      Mapping:", source.Mapping)
		fmt.Println("      Files:", source.Files)
	}
}

func (f Fixture) printStatusPostgreSQL(cfg *config.PostgreSQL) {
	fmt.Println("PostgreSQL -", cfg.Destination)
	fmt.Println("  Sources:")
	for _, source := range cfg.Sources {
		fmt.Println("    - Database:", source.Database)
		fmt.Println("      Definition:", source.Definition)
		fmt.Println("      Files:", source.Files)
	}
}

func (f Fixture) printStatusRedis(cfg *config.Redis) {
	fmt.Println("Redis -", cfg.Destination)
	fmt.Println("  Sources:", cfg.Source)
}
