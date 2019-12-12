package remove

import (
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/oldcode/fixture/pull/postgresql"
	"github.com/taufik-rama/maiden/oldcode/fixture/set/elasticsearch"
)

// Handler for `fixture rm` command
type Handler struct{}

// Command returns `fixture rm` command process
func (c *Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove the fixtures data",
		Run:     c.RunCommand,
	}

	c.register(cmd)

	return cmd
}

// RunCommand runs `fixture push` command
func (c Handler) RunCommand(cmd *cobra.Command, args []string) {
	// // Run every subcommand if not specified
	// if c.verbose {
	// 	log.Println("No fixtures specified, running every fixtures...")
	// }
	// for name, fixture := range c.fixtures {
	// 	if c.verbose {
	// 		log.Printf("Running `%s` fixture\n", name)
	// 	}
	// 	fixture.SetVerbose(c.verbose)
	// 	fixture.RunCommand(cmd, args)
	// }
}

// Register what `push` subcommand we'll use.
// `fixtures` field will be used as a list of registered
// subcommands.
func (c Handler) register(cmd *cobra.Command) {
	cmd.AddCommand((&elasticsearch.Handler{}).Command())
	cmd.AddCommand((&postgresql.Handler{}).Command())
}
