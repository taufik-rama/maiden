package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/taufik-rama/maiden/config"
)

type postgresql struct {
	config config.Config
}

// Sets the `config` field
func (p *postgresql) setConfig(c config.Config, err error) error {
	p.config = c
	return err
}

// Remove every database defined on the fixtures directory.
func (p postgresql) remove() error {
	log.Println("Note: pushing without the database name is currently unimplemented :(")
	return nil
}

// Remove the database
func (p postgresql) removeDatabase(database string) error {

	db, err := sql.Open("postgres", (p.config.Fixtures().Postgresql().PushTo() + "?sslmode=disable"))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, database))
	if err != nil {
		return err
	}

	return nil
}

// Remove the table
func (p postgresql) removeTable(database, table string) error {

	db, err := sql.Open("postgres", (p.config.Fixtures().Postgresql().PushTo() + database + "?sslmode=disable"))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`DROP TABLE "%s"`, table))
	if err != nil {
		return err
	}

	return nil
}
