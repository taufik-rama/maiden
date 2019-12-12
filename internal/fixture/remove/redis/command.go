package redis

import (
	"log"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"
)

// RedisCommand is the command name
var RedisCommand = "redis"

// Handler for `fixture rm redis` command
type Handler struct {

	// use this before any sort of print log
	verbose bool
}

// Command returns `fixture rm redis` command process
func (c *Handler) Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   RedisCommand,
		Short: "Unimplemented",
		Run:   c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

// RunCommand runs `fixture rm redis` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Fatalln(err)
	}

	if _, err := client.FlushAll().Result(); err != nil {
		log.Fatalln(err)
	}
}
