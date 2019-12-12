package postgresql

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/taufik-rama/maiden/config"
	"github.com/taufik-rama/maiden/config/fixtures"
)

type postgresql struct {
	config   config.Config
	fixtures fixtures.Config
}

// Sets the `config` field
func (p *postgresql) setConfig(c config.Config, err error) error {
	p.config = c
	return err
}

// Sets the `fixtures` field
func (p *postgresql) setFixtures(c fixtures.Config, err error) error {
	p.fixtures = c
	return err
}

// Push every database defined on the fixtures directory.
func (p postgresql) push() error {

	for _, database := range p.fixtures.Postgresql().Databases() {
		log.Println("Creating database", database.Name())
		db, err := p.createDatabase(database.Name())
		if err != nil {
			return err
		}
		for _, table := range database.Tables() {
			log.Println("Creating table", table.Name())
			if err := p.pushTable(db, table.Schema()); err != nil {
				log.Println(err)
				continue
			}
			for _, record := range table.Records() {
				log.Println("Creating records for table", table.Name())
				if err := p.pushRecords(db, record); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}

	return nil
}

// Create PostgreSQL database, will ignore `already exists` database error.
func (p postgresql) createDatabase(database string) (*sql.DB, error) {
	db, err := sql.Open("postgres", (p.config.Fixtures().Postgresql().PushTo() + "?sslmode=disable"))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, database))
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, err
	}
	return sql.Open("postgres", (p.config.Fixtures().Postgresql().PushTo() + database + "?sslmode=disable"))
}

// Push the PostgreSQL table
func (p postgresql) pushTable(database *sql.DB, sql io.Reader) error {
	return executeSQL(database, sql)
}

// Push the PostgreSQL table records
func (p postgresql) pushRecords(database *sql.DB, sql io.Reader) error {
	return executeSQL(database, sql)
}

// Execute the `sql` queries on `database` connection
func executeSQL(database *sql.DB, sql io.Reader) error {
	bytes, err := ioutil.ReadAll(sql)
	if err != nil {
		return fmt.Errorf("reading SQL: %s", err)
	}
	_, err = database.Exec(string(bytes))
	return err
}
