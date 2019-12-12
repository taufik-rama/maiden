package fixture

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/internal/config"
	"github.com/taufik-rama/maiden/internal/fixture/push"
)

// Handler for `fixture` command
type Handler struct{}

// Command returns `fixture` command process
func (h Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "fixture",
		Aliases: []string{"f"},
		Short:   "Interact with the available fixtures",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.New(config.Args.ConfigFile)
			cfg.ResolveFixtures()
			fmt.Println(cfg.Fixtures.Cassandra)
			fmt.Println(cfg.Fixtures.Elasticsearch)
			fmt.Println(cfg.Fixtures.PostgreSQL)
			fmt.Println(cfg.Fixtures.Redis)
		},
	}

	h.register(cmd)

	return cmd
}

func (h Handler) register(cmd *cobra.Command) {
	cmd.AddCommand((&push.Handler{}).Command())
	// cmd.AddCommand((&remove.Handler{}).Command())
}
