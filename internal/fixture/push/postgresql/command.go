package postgresql

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	// "github.com/taufik-rama/maiden/internal/config"
)

// PostgresqlCommand is the command name
var PostgresqlCommand = "postgresql"

// Handler for `fixture push postgresql` command
type Handler struct {

	// use this before any sort of print log
	verbose bool

	// what database to use
	database string

	// what tables to use, comma separated
	tables string

	// what method we use for pushing process
	pusher interface {

		// Push all
		push() error

		// Create database (if not exists) using given name
		createDatabase(string) (*sql.DB, error)

		// Use the reader as values for sql table query
		pushTable(*sql.DB, io.Reader) error

		// Use the reader as values for sql records query
		pushRecords(*sql.DB, io.Reader) error

		setConfig(config.Config, error) error
		setFixtures(fixtures.Config, error) error
	}
}

// Command returns `fixture push postgresql` command process
func (c *Handler) Command() *cobra.Command {

	c.pusher = &postgresql{}

	cmd := &cobra.Command{
		Use:     PostgresqlCommand,
		Aliases: []string{"pg"},
		Short:   "Push the PostgreSQL fixtures data",
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

// RunCommand runs `fixture push postgresql` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {

	if err := c.check(); err != nil {
		log.Fatalln(err)
	}

	if err := c.pusher.setConfig(config.Configure()); err != nil {
		log.Fatalln(err)
	}

	if err := c.pusher.setFixtures(fixtures.Configure()); err != nil {
		log.Fatalln(err)
	}

	// Push all the postgresql fixture if not defined
	if c.database == "" {
		c.printf("No database specified, pushing current fixtures database & tables...")
		if err := c.pusher.push(); err != nil {
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
		log.Fatalf("Can't push database `%s`: not a directory inside `%s`", c.database, fixtures.PostgresqlFixturesDir)
	}

	// Check table fixtures & assign the table name set
	tables := make(map[string]struct{})
	for _, table := range strings.Split(c.tables, ",") {

		c.printf("Checking fixture for table `%s`", table)
		dir := fixtures.PostgresqlFixturesDir + config.DirectorySeparator + c.database + config.DirectorySeparator + table + ".sql"
		if info, err := os.Lstat(dir); err != nil {
			if os.IsNotExist(err) {
				c.printf("`%s` does not exists on the database fixtures", (table + ".sql"))
				continue
			}
			log.Fatalln(err)
		} else if info.IsDir() {
			c.printf("Can't push `%s` because it's not a file", dir)
			continue
		}

		c.printf("Checking fixture for table `%s` records", table)
		dir = fixtures.PostgresqlFixturesDir + config.DirectorySeparator + c.database + config.DirectorySeparator + table
		if info, err := os.Lstat(dir); err != nil {
			if os.IsNotExist(err) {
				c.printf("Records for `%s` does not exists on the database fixtures. The directory must have the same name as table", (table + ".sql"))
				continue
			}
			log.Fatalln(err)
		} else if !info.IsDir() {
			c.printf("Can't push `%s` because it's not a directory", dir)
			continue
		}

		tables[table] = struct{}{}
	}

	db, err := c.pusher.createDatabase(c.database)
	if err != nil {
		log.Fatalln(err)
	}

	// Find the postgresql fixtures reader
	f, _ := fixtures.Configure()
	for _, database := range f.Postgresql().Databases() {
		if database.Name() == c.database {
			for _, table := range database.Tables() {

				if _, ok := tables[table.Name()]; !ok {
					continue
				}

				c.printf("Pushing table `%s`", table.Name())
				if err := c.pusher.pushTable(db, table.Schema()); err != nil {
					log.Println(err)
					continue
				}

				c.printf("Pushing records for table `%s`", table.Name())
				for _, record := range table.Records() {
					if err := c.pusher.pushRecords(db, record); err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}
}

// Check the validity of command.
// For `postgresql` we need the `-d` & `-t` command to both be filled
// or empty.
func (c *Handler) check() error {
	if c.database == "" && c.tables != "" {
		return errors.New("`--database` flag is needed if you use `--tables`")
	}
	if c.database != "" && c.tables == "" {
		return errors.New("`--tables` flag is needed if you use `--database`")
	}
	return nil
}

// Print according to the verbosity flag
func (c *Handler) printf(format string, v ...interface{}) {
	if c.verbose {
		log.Printf(format, v...)
	}
}
