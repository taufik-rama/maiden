package postgresql

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config"
	"github.com/taufik-rama/maiden/config/fixtures"
)

// PostgresqlCommand is the command name
var PostgresqlCommand = "postgresql"

// Handler for `fixture rm postgresql` command
type Handler struct {

	// use this before any sort of print log
	verbose bool

	// what database to use
	database string

	// what tables to use, comma separated
	tables string

	// what method we use for removing process
	remover interface {

		// Remove all
		remove() error

		// Remove database with given name
		removeDatabase(string) error

		// Remove database with given table & database name
		removeTable(string, string) error

		setConfig(config.Config, error) error
	}
}

// Command returns `fixture rm postgresql` command process
func (c *Handler) Command() *cobra.Command {

	c.remover = &postgresql{}

	cmd := &cobra.Command{
		Use:     PostgresqlCommand,
		Aliases: []string{"pg"},
		Short:   "Remove the PostgreSQL fixtures data",
		Run:     c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().StringVarP(&c.database, "database", "d", "", "What database to use")
	cmd.PersistentFlags().StringVarP(&c.tables, "tables", "t", "", "What table(s) to use, comma separated")

	return cmd
}

// SetVerbose flags
func (c *Handler) SetVerbose(v bool) {
	c.verbose = v
}

// RunCommand runs `fixture rm postgresql` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {

	if err := c.remover.setConfig(config.Configure()); err != nil {
		log.Fatalln(err)
	}

	// Pull all the postgresql fixture if not defined
	if c.database == "" {
		c.printf("No database specified, removing current fixtures database & tables...")
		if err := c.remover.remove(); err != nil {
			log.Fatalln(err)
		}
		return
	}

	// Check database fixtures
	dir := fixtures.PostgresqlFixturesDir + config.DirectorySeparator + c.database
	if info, err := os.Lstat(dir); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Database `%s` does not exists on the fixtures", c.database)
		}
		log.Fatalln(err)
	} else if !info.IsDir() {
		log.Fatalf("Can't rm database `%s`: not a directory inside `%s`", c.database, fixtures.PostgresqlFixturesDir)
	}

	if c.tables != "" {
		for _, table := range strings.Split(c.tables, ",") {
			if err := c.remover.removeTable(c.database, table); err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		if err := c.remover.removeDatabase(c.database); err != nil {
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
