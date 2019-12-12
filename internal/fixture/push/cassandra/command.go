package cassandra

import (
	"log"

	"github.com/spf13/cobra"
)

// CassandraCommand is the command name
var CassandraCommand = "cassandra"

// Handler for `fixture push cassandra` command
type Handler struct{}

// Command returns `fixture push cassandra` command process
func (c *Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   CassandraCommand,
		Short: "Unimplemented",
		Run:   c.RunCommand,
	}

	return cmd
}

// RunCommand runs `fixture push cassandra` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {
	log.Println("Cassandra fixture is unimplemented")
}
