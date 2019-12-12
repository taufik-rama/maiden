package push

import (
	"log"

	"github.com/spf13/cobra"
	// "github.com/taufik-rama/maiden/internal/fixture/push/cassandra"
	"github.com/taufik-rama/maiden/internal/config"
	"github.com/taufik-rama/maiden/internal/fixture/push/elasticsearch"
	// "github.com/taufik-rama/maiden/internal/fixture/push/postgresql"
	// "github.com/taufik-rama/maiden/internal/fixture/push/redis"
)

// Handler for `fixture push` command
type Handler struct {

	// whether we'll push all the fixtures or not
	all bool
}

// Command returns `fixture push` command process
func (c *Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push the fixtures data",
		Run:   c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.all, "all", "", false, "Push all fixtures")

	c.register(cmd)

	return cmd
}

// RunCommand runs `fixture push` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		log.Fatalf("Unknown args `%s`", args[0])
		return
	}

	cfg := config.New(config.Args.ConfigFile)
	cfg.ResolveFixtures()

	active := cfg.Fixtures.ActiveFixtures()

	if len(active) == 0 {
		config.Print("No fixtures defined, nothing to do")
		return
	}

	if len(active) > 1 && !c.all {
		log.Printf("`--all` flags are needed when more than 1 fixtures are defined")
		return
	}

	for _, fixture := range active {
		switch fixture.Name() {
		case cfg.Fixtures.Cassandra.Name():
			// handler := &cassandra.Handler{}
			// handler.RunCommand(handler.Command(), args)
		case cfg.Fixtures.Elasticsearch.Name():
			handler := &elasticsearch.Handler{
				All: true,
			}
			handler.RunCommand(handler.Command(), args)
		case cfg.Fixtures.PostgreSQL.Name():
			// handler := &postgresql.Handler{}
			// handler.RunCommand(handler.Command(), args)
		case cfg.Fixtures.Redis.Name():
			// handler := &redis.Handler{}
			// handler.RunCommand(handler.Command(), args)
		default:
			log.Fatalf("Internal error: fixture `%s` is not regitered", fixture.Name())
		}
	}
}

// Register what `push` subcommand we'll use.
// `fixtures` field will be used as a list of registered
// subcommands.
func (c *Handler) register(cmd *cobra.Command) {
	// cmd.AddCommand((&cassandra.Handler{}).Command())
	cmd.AddCommand((&elasticsearch.Handler{}).Command())
	// cmd.AddCommand((&postgresql.Handler{}).Command())
	// cmd.AddCommand((&redis.Handler{}).Command())
}
