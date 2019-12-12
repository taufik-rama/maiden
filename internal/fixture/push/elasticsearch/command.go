package elasticsearch

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/internal/config"
	"github.com/taufik-rama/maiden/internal/structure"
)

// Handler for `fixture push elasticsearch` command
type Handler struct {

	// Push all fixtures
	All bool

	// Index to push
	index string

	pusher interface {

		// Push all
		push() error

		// Push index mapping given index name & contents
		pushIndex(string, io.Reader, map[string]fieldsInfo) error

		// Push documents values given index name & contents
		pushDocuments(string, io.Reader, *int, map[string]fieldsInfo) error

		setConfig(config.MaidenConfig)
	}
}

// Command returns `fixture push elasticsearch` command process
func (c *Handler) Command() *cobra.Command {

	pusher := &elasticsearch{}

	cmd := &cobra.Command{
		Use:     "elasticsearch",
		Aliases: []string{"es"},
		Short:   "Push the Elasticsearch fixtures data",
		Run:     c.RunCommand,
	}

	cmd.PersistentFlags().BoolVarP(&c.All, "all", "", false, "Push all fixtures")
	cmd.PersistentFlags().StringVarP(&c.index, "index", "i", "", "What index to use, multiple values w/ comma separation")
	cmd.PersistentFlags().BoolVarP(&pusher.useID, "use-id", "", false, "Use fixture document ID")
	cmd.PersistentFlags().BoolVarP(&pusher.useTime, "use-time", "", false, "Use fixture document time fields value")

	c.pusher = pusher

	return cmd
}

// RunCommand runs `fixture push elasticsearch` command
func (c *Handler) RunCommand(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		log.Fatalf("Unknown args `%s`", args[0])
	}

	cfg := config.New(config.Args.ConfigFile)
	cfg.ResolveFixtures()
	c.pusher.setConfig(cfg)

	if c.All {
		config.Print("No index specified, pushing current fixtures indices...")
		if err := c.pusher.push(); err != nil {
			log.Fatalln(err)
		}
		return
	}

	if len(strings.TrimSpace(c.index)) == 0 {
		log.Fatalf("Please specify the Elasticsearch index name")
	}

	for _, index := range strings.Split(c.index, ",") {
		config.Print("Pushing index `%s`", index)

		info := make(map[string]fieldsInfo)

		str := &structure.Structure{}
		str.ResolveFixtures(cfg.Fixtures)

		// For checking purpose
		var exists bool

		for _, fixture := range str.Fixture.Elasticsearch.Indices {

			if fixture.Name != index {
				continue
			}

			exists = true

			// File scope
			{
				file, err := os.Open(fixture.Mapping)
				if err != nil {
					log.Fatalln(err)
				}
				defer file.Close()

				if err := c.pusher.pushIndex(fixture.Name, file, info); err != nil {
					log.Fatalln(err)
				}
			}

			// For cross-file document ID
			var id int

			for _, document := range fixture.Documents {
				file, err := os.Open(document)
				if err != nil {
					log.Fatalln(err)
				}
				defer file.Close()

				if err := c.pusher.pushDocuments(fixture.Name, file, &id, info); err != nil {
					log.Fatalln(err)
				}
			}
		}

		if !exists {
			log.Fatalf("Index `%s` does not exists", index)
		}
	}
}
