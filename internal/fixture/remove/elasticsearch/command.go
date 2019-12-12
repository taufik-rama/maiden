package elasticsearch

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config"
	"github.com/taufik-rama/maiden/config/fixtures"
)

// ElasticsearchCommand is the command name
var ElasticsearchCommand = "elasticsearch"

// Handler for `fixture rm elasticsearch` command
type Handler struct {

	// use this before any sort of print log
	verbose bool

	// what indices to remove
	indices string

	// what method we use for removing process
	remover interface {

		// Remove all
		remove() error

		// Remove index with given name
		removeIndex(string) error

		setConfig(config.Config, error) error
		setFixtures(fixtures.Config, error) error
	}
}

// Command returns `fixture rm elasticsearch` command process
func (c *Handler) Command() *cobra.Command {

	c.remover = &elasticsearch{}

	cmd := &cobra.Command{
		Use:     ElasticsearchCommand,
		Aliases: []string{"es"},
		Short:   "Remove the Elasticsearch fixtures data",
		Run:     c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().StringVarP(&c.indices, "indices", "i", "", "What indices to use, comma separated")

	return cmd
}

// SetVerbose flags
func (c *Handler) SetVerbose(v bool) {
	c.verbose = v
}

// RunCommand runs `fixture rm elasticsearch` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {

	c.remover.setConfig(config.Configure())
	c.remover.setFixtures(fixtures.Configure())

	// Push all the elasticsearch fixture if not defined
	if c.indices == "" {
		c.printf("No index specified, removing current fixtures indices...")
		if err := c.remover.remove(); err != nil {
			log.Fatalln(err)
		}
		return
	}

	// Check index fixtures & assign the index name set
	indices := make(map[string]struct{})
	for _, index := range strings.Split(c.indices, ",") {

		c.printf("Checking fixture for index `%s`", index)
		dir := fixtures.ElasticsearchFixturesDir + config.DirectorySeparator + index + config.DirectorySeparator + index + ".json"
		if info, err := os.Lstat(dir); err != nil {
			if os.IsNotExist(err) {
				c.printf("`%s` does not exists on the index fixtures", (index + ".json"))
				continue
			}
			log.Fatalln(err)
		} else if info.IsDir() {
			c.printf("Can't push `%s` because it's not a file", dir)
			continue
		}

		indices[index] = struct{}{}
	}

	for index := range indices {
		if err := c.remover.removeIndex(index); err != nil {
			log.Fatalln(err)
		}
	}
}

// Print according to the verbosity flag
func (c *Handler) printf(format string, v ...interface{}) {
	if c.verbose {
		log.Printf(format, v...)
	}
}
