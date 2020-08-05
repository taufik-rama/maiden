package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
)

var cqlFormat = ".cql"
var cassandraKeyspace = `CREATE KEYSPACE IF NOT EXISTS "%s" WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};`

// Cassandra ...
type Cassandra struct{}

// Push ...
func (c Cassandra) Push(cfg *config.Cassandra) {

	keyspaces := make(map[string]struct{})
	tables := make(map[string]struct{})

	for _, source := range cfg.Sources {
		keyspaces[source.Keyspace] = struct{}{}
		files, err := ioutil.ReadDir(source.Files)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), cqlFormat) {
				tables[file.Name()] = struct{}{}
			}
		}
	}

	c.push(cfg, keyspaces, tables)
}

// PushKeyspace ...
func (c Cassandra) PushKeyspace(cfg *config.Cassandra, keyspace string) {

	keyspaces := make(map[string]struct{})
	tables := make(map[string]struct{})

	for _, keyspace := range strings.Split(keyspace, ",") {
		if source, ok := cfg.Sources[keyspace]; !ok {
			log.Printf("Keyspace `%s` not registered", keyspace)
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
				if strings.HasSuffix(file.Name(), cqlFormat) {
					tables[file.Name()] = struct{}{}
				}
			}
		}
		keyspaces[keyspace] = struct{}{}
	}

	c.push(cfg, keyspaces, tables)
}

// PushKeyspaceTable ...
func (c Cassandra) PushKeyspaceTable(cfg *config.Cassandra, keyspace, table string) {

	keyspaces := map[string]struct{}{
		keyspace: {},
	}
	tables := make(map[string]struct{})

	source, ok := cfg.Sources[keyspace]
	if !ok {
		log.Printf("Keyspace `%s` not registered", keyspace)
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
			if strings.HasSuffix(file.Name(), cqlFormat) {
				registered[file.Name()] = struct{}{}
			}
		}
		for _, table := range strings.Split(table, ",") {
			if _, ok := registered[(table + cqlFormat)]; !ok {
				log.Printf("Table `%s` on keyspace `%s` not registered", table, keyspace)
				continue
			}
			tables[(table + cqlFormat)] = struct{}{}
		}
	}

	c.push(cfg, keyspaces, tables)
}

// Remove ...
func (c Cassandra) Remove(cfg *config.Cassandra) {
	session, err := c.createSession(cfg.Destination, "")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	for _, source := range cfg.Sources {
		if err := session.Query(fmt.Sprintf(`DROP KEYSPACE "%s"`, source.Keyspace)).Exec(); err != nil {
			panic(err)
		}
	}
}

// RemoveKeyspace ...
func (c Cassandra) RemoveKeyspace(cfg *config.Cassandra, keyspace string) {

	keyspaces := make(map[string]struct{})

	for _, keyspace := range strings.Split(keyspace, ",") {
		if _, ok := cfg.Sources[keyspace]; !ok {
			log.Printf("Keyspace `%s` not registered", keyspace)
			continue
		}
		keyspaces[keyspace] = struct{}{}
	}

	session, err := c.createSession(cfg.Destination, "")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	for _, source := range cfg.Sources {
		if _, ok := keyspaces[source.Keyspace]; !ok {
			continue
		}
		if err := session.Query(fmt.Sprintf(`DROP KEYSPACE "%s"`, source.Keyspace)).Exec(); err != nil {
			panic(err)
		}
	}
}

// RemoveKeyspaceTable ...
func (c Cassandra) RemoveKeyspaceTable(cfg *config.Cassandra, keyspace, table string) {

	source, ok := cfg.Sources[keyspace]
	if !ok {
		log.Printf("Keyspace `%s` not registered", keyspace)
		return
	}

	session, err := c.createSession(cfg.Destination, source.Keyspace)
	if err != nil {
		panic(err)
	}
	defer session.Close()

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
			if strings.HasSuffix(file.Name(), cqlFormat) {
				tables[file.Name()] = struct{}{}
			}
		}
		for _, table := range strings.Split(table, ",") {
			if _, ok := tables[(table + cqlFormat)]; !ok {
				log.Printf("Table `%s` on keyspace `%s` not registered", table, keyspace)
				continue
			}
			if err := session.Query(fmt.Sprintf(`DROP TABLE "%s"`, table)).Exec(); err != nil {
				panic(err)
			}
		}
	}
}

func (c Cassandra) push(cfg *config.Cassandra, keyspaces, tables map[string]struct{}) error {

	for _, source := range cfg.Sources {

		if _, ok := keyspaces[source.Keyspace]; !ok {
			continue
		}

		if err := c.pushKeyspace(cfg.Destination, c.definition(source.Definition, source.Keyspace)); err != nil {
			if strings.Contains(err.Error(), "Cannot add existing keyspace") {
				internal.Print("Keyspace `%s` exists, did not need to be pushed", source.Keyspace)
			} else {
				panic(err)
			}
		}

		if err := c.checkSource(source.Files, source.Keyspace); err != nil {
			panic(err)
		}

		files, err := ioutil.ReadDir(source.Files)
		if err != nil {
			panic(err)
		}

		session, err := c.createSession(cfg.Destination, source.Keyspace)
		if err != nil {
			panic(err)
		}
		defer session.Close()

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
			if err := c.executeCQL(session, string(bytes)); err != nil {
				panic(err)
			}
		}

		// Records
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			if _, ok := tables[(file.Name() + cqlFormat)]; !ok {
				continue
			}
			path := source.Files + string(os.PathSeparator) + file.Name()
			if err := c.pushRecords(path, session); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func (c Cassandra) definition(definition, keyspace string) string {
	if definition != "" {
		bytes, err := internal.Read(definition)
		if err != nil {
			panic(err)
		}
		return string(bytes)
	}
	return fmt.Sprintf(cassandraKeyspace, keyspace)
}

func (c Cassandra) executeCQL(session *gocql.Session, cql string) error {
	if strings.TrimSpace(cql) == "" {
		return nil
	}

	// Fix: Cassandra only allows 1 statement per `Query()` call, split & iterate here
	for _, query := range strings.Split(cql, ";") {
		if strings.TrimSpace(query) == "" {
			continue
		} else if err := session.Query(query).Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (c Cassandra) createSession(destination, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(destination)
	cluster.Timeout = 5 * time.Second
	cluster.Keyspace = keyspace
	return cluster.CreateSession()
}

func (c Cassandra) pushKeyspace(destination, definition string) error {
	cluster := gocql.NewCluster(destination)
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	return session.Query(definition).Exec()
}

func (c Cassandra) checkSource(source, keyspace string) error {
	if info, err := os.Lstat(source); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no fixtures are found for keyspace `%s`", keyspace)
		}
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("the fixtures for keyspace `%s` must be a directory", keyspace)
	}
	return nil
}

func (c Cassandra) pushRecords(dir string, session *gocql.Session) error {
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
		return c.executeCQL(session, string(bytes))
	})
}
