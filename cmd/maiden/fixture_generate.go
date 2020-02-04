package main

import (
	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/generator.v1"
	"github.com/taufik-rama/maiden/internal"
)

// FixtureGenerate command handler
type FixtureGenerate struct{}

// Command returns `fixture generate` command process
func (f FixtureGenerate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate the fixtures storage service docker file",
		Run:     f.RunCommand,
	}
	return cmd
}

// RunCommand runs `fixture generate` command
func (f FixtureGenerate) RunCommand(cmd *cobra.Command, args []string) {

	fixture := new(config.Fixture)
	if err := fixture.Parse(internal.Args.ConfigFile); err != nil {
		panic(err)
	}

	var images []string
	if fixture.Cassandra != nil {
		images = append(images, generator.Cassandra)
	}
	if fixture.Elasticsearch != nil {
		images = append(images, generator.Elasticsearch)
	}
	if fixture.PostgreSQL != nil {
		images = append(images, generator.PostgreSQL)
	}
	if fixture.Redis != nil {
		images = append(images, generator.Redis)
	}

	(generator.DockerCompose{
		Output: fixture.DockerComposeOutput,
		Images: images,
	}).GenerateCommand(cmd, args)
}
