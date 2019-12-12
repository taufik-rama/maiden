package cassandra

import (
	"log"

	"github.com/spf13/cobra"
)

// CassandraCommand is the command name
var CassandraCommand = "cassandra"

// Handler for `fixture rm cassandra` command
type Handler struct {

	// use this before any sort of print log
	verbose bool
}

// Command returns `fixture rm cassandra` command process
func (c *Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   CassandraCommand,
		Short: "Unimplemented",
		Run:   c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

// RunCommand runs `fixture rm cassandra` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {
	log.Println("Cassandra fixture is unimplemented")
}
