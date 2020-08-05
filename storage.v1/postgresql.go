package storage

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
)

var sqlFormat = ".sql"
var postgresqlDatabase = `CREATE DATABASE "%s"`

// PostgreSQL ...
type PostgreSQL struct{}

// Push ...
func (p PostgreSQL) Push(cfg *config.PostgreSQL) {

	databases := make(map[string]struct{})
	tables := make(map[string]struct{})

	for _, source := range cfg.Sources {
		databases[source.Database] = struct{}{}
		files, err := ioutil.ReadDir(source.Files)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), sqlFormat) {
				tables[file.Name()] = struct{}{}
			}
		}
	}

	p.push(cfg, databases, tables)
}

// PushDatabase ...
func (p PostgreSQL) PushDatabase(cfg *config.PostgreSQL, database string) {

	databases := make(map[string]struct{})
	tables := make(map[string]struct{})

	for _, database := range strings.Split(database, ",") {
		if source, ok := cfg.Sources[database]; !ok {
			log.Printf("Database `%s` not registered", database)
			continue
		} else {
			files, err := ioutil.ReadDir(source.Files)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if strings.HasSuffix(file.Name(), sqlFormat) {
					tables[file.Name()] = struct{}{}
				}
			}
		}
		databases[database] = struct{}{}
	}

	p.push(cfg, databases, tables)
}

// PushDatabaseTable ...
func (p PostgreSQL) PushDatabaseTable(cfg *config.PostgreSQL, database, table string) {

	databases := map[string]struct{}{
		database: {},
	}
	tables := make(map[string]struct{})

	source, ok := cfg.Sources[database]
	if !ok {
		log.Printf("Database `%s` not registered", database)
		return
	}

	{
		files, err := ioutil.ReadDir(source.Files)
		if err != nil {
			panic(err)
		}

		registered := make(map[string]struct{})
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), sqlFormat) {
				registered[file.Name()] = struct{}{}
			}
		}
		for _, table := range strings.Split(table, ",") {
			if _, ok := registered[(table + sqlFormat)]; !ok {
				log.Printf("Table `%s` on database `%s` not registered", table, database)
				continue
			}
			tables[(table + sqlFormat)] = struct{}{}
		}
	}

	p.push(cfg, databases, tables)
}

// Remove ...
func (p PostgreSQL) Remove(cfg *config.PostgreSQL) {

	db, err := sql.Open("postgres", (cfg.Destination + "/?sslmode=disable"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, source := range cfg.Sources {
		if err := p.executeSQL(db, fmt.Sprintf(`DROP DATABASE "%s"`, source.Database)); err != nil {
			panic(err)
		}
	}
}

// RemoveDatabase ...
func (p PostgreSQL) RemoveDatabase(cfg *config.PostgreSQL, database string) {

	databases := make(map[string]struct{})

	for _, database := range strings.Split(database, ",") {
		if _, ok := cfg.Sources[database]; !ok {
			log.Printf("error %s not registered", database)
			continue
		}
		databases[database] = struct{}{}
	}

	db, err := sql.Open("postgres", (cfg.Destination + "/?sslmode=disable"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, source := range cfg.Sources {
		if _, ok := databases[source.Database]; !ok {
			continue
		}
		if err := p.executeSQL(db, fmt.Sprintf(`DROP DATABASE "%s"`, source.Database)); err != nil {
			panic(err)
		}
	}
}

// RemoveDatabaseTable ...
func (p PostgreSQL) RemoveDatabaseTable(cfg *config.PostgreSQL, database, table string) {

	source, ok := cfg.Sources[database]
	if !ok {
		log.Printf("Database `%s` not registered", database)
		return
	}

	db, err := sql.Open("postgres", (cfg.Destination + "/" + database + "?sslmode=disable"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	{
		files, err := ioutil.ReadDir(source.Files)
		if err != nil {
			panic(err)
		}

		tables := make(map[string]struct{})
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), sqlFormat) {
				tables[file.Name()] = struct{}{}
			}
		}
		for _, table := range strings.Split(table, ",") {
			if _, ok := tables[(table + sqlFormat)]; !ok {
				log.Printf("Table `%s` on database `%s` not registered", table, database)
				continue
			}
			if err := p.executeSQL(db, fmt.Sprintf(`DROP TABLE "%s"`, table)); err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					internal.Print("Table `%s` does not exists, did not need to be removed", table)
				} else {
					panic(err)
				}
			}
		}
	}
}

func (p PostgreSQL) push(cfg *config.PostgreSQL, databases, tables map[string]struct{}) {

	for _, source := range cfg.Sources {

		if _, ok := databases[source.Database]; !ok {
			continue
		}

		internal.Print("Creating database `%s`", source.Database)
		if err := p.createDatabase(source.Database, source.Definition, cfg.Destination); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				internal.Print("Database `%s` exists, did not need to be pushed", source.Database)
			} else {
				panic(err)
			}
		}

		db, err := sql.Open("postgres", (cfg.Destination + "/" + source.Database + "?sslmode=disable"))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if info, err := os.Lstat(source.Files); err != nil {
			if os.IsNotExist(err) {
				internal.Print("No fixtures are found for database `%s`", source.Database)
				continue
			}
			panic(err)
		} else if !info.IsDir() {
			internal.Print("The fixtures for database `%s` must be a directory", source.Database)
			continue
		}

		internal.Print("Pushing table & records...")
		files, err := ioutil.ReadDir(source.Files)
		if err != nil {
			panic(err)
		}

		// Tables
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if _, ok := tables[file.Name()]; !ok {
				continue
			}
			path := source.Files + string(os.PathSeparator) + file.Name()
			bytes, err := internal.Read(path)
			if err != nil {
				panic(err)
			}
			if err := p.executeSQL(db, string(bytes)); err != nil {
				panic(err)
			}
		}

		// Records
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			if _, ok := tables[(file.Name() + ".sql")]; !ok {
				continue
			}
			path := source.Files + string(os.PathSeparator) + file.Name()
			if err := p.pushRecords(db, path); err != nil {
				panic(err)
			}
		}
	}
}

func (p PostgreSQL) pushRecords(db *sql.DB, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		bytes, err := internal.Read(path)
		if err != nil {
			panic(err)
		}
		return p.executeSQL(db, string(bytes))
	})
}

func (p PostgreSQL) createDatabase(database, definition, destination string) error {
	db, err := sql.Open("postgres", (destination + "/?sslmode=disable"))
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(p.definition(database, definition)); err != nil {
		return err
	}
	return nil
}

func (p PostgreSQL) definition(database, definition string) string {
	if definition != "" {
		bytes, err := internal.Read(definition)
		if err != nil {
			panic(err)
		}
		return string(bytes)
	}
	return fmt.Sprintf(postgresqlDatabase, database)
}

// Execute the `sql` queries on `database` connection
func (p PostgreSQL) executeSQL(database *sql.DB, sql string) error {
	if strings.TrimSpace(sql) == "" {
		return nil
	}
	_, err := database.Exec(sql)
	return err
}
