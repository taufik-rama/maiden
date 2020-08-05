package generator

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/internal"
	"gopkg.in/yaml.v2"
)

var (

	// Cassandra service name
	Cassandra = "cassandra"

	// Elasticsearch service name
	Elasticsearch = "elasticsearch"

	// PostgreSQL service name
	PostgreSQL = "postgres"

	// Redis service name
	Redis = "redis"

	// Images name
	Images = map[string]string{
		Cassandra:     "cassandra:3",
		Elasticsearch: "elasticsearch:6.7.1",
		PostgreSQL:    "postgres:9.6-alpine",
		Redis:         "redis:5.0-alpine",
	}

	// Ports list
	Ports = map[string][]string{
		Cassandra:     {"9042:9042"},
		Elasticsearch: {"9200:9200", "9300:9300"},
		PostgreSQL:    {"5432:5432"},
		Redis:         {"6379:6379"},
	}
)

// DockerCompose config type
type DockerCompose struct {
	Output string
	Images []string
}

// GenerateCommand is the handler for cobra command-line
func (d DockerCompose) GenerateCommand(*cobra.Command, []string) {

	type service struct {
		Image       string   `yaml:"image,omitempty"`
		Ports       []string `yaml:"ports,omitempty"`
		Environment []string `yaml:"environment,omitempty"`
	}
	type services map[string]service

	var contents struct {
		Version  string   `yaml:"version"`
		Services services `yaml:"services"`
	}

	contents.Version = "3"
	contents.Services = make(services)

	for _, name := range d.Images {

		internal.Print("Parsing `%s` docker service", name)

		s := service{
			Image: Images[name],
			Ports: Ports[name],
		}

		if name == Elasticsearch {
			s.Environment = []string{
				"discovery.type=single-node",
				"cluster.routing.allocation.disk.watermark.low=95%",
				"cluster.routing.allocation.disk.watermark.high=95%",
			}
		}

		contents.Services[name] = s
	}

	bytes, err := yaml.Marshal(contents)
	if err != nil {
		panic(err)
	}

	internal.Print("Writing docker-compose to `%s`", d.Output)
	file, err := os.Create(d.Output)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		panic(err)
	}
}
