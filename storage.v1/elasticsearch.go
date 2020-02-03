package storage

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/taufik-rama/maiden/config.v1"
)

// Elasticsearch ...
type Elasticsearch struct{}

// Push ...
func (e Elasticsearch) Push(cfg *config.Elasticsearch) {

	for _, source := range cfg.Sources {

		es, err := elasticsearch.NewDefaultClient()
		if err != nil {
			panic(err)
		}

		mapping, err := os.Open(source.Mapping)
		if err != nil {
			panic(err)
		}
		defer mapping.Close()

		request := esapi.IndicesCreateRequest{
			Index:  source.Index,
			Body:   mapping,
			Pretty: true,
		}

		response, err := request.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}

		if response.IsError() {
			panic(response.String())
		}
	}
}

// PushIndex ...
func (e Elasticsearch) PushIndex(cfg *config.Elasticsearch, index string) {

	indices := make(map[string]struct{})

	for _, index := range strings.Split(index, ",") {
		if _, ok := cfg.Sources[index]; !ok {
			log.Printf("Index `%s` not registered", index)
			continue
		}
		indices[index] = struct{}{}
	}

	for _, source := range cfg.Sources {

		if _, ok := indices[source.Index]; !ok {
			continue
		}

		es, err := elasticsearch.NewDefaultClient()
		if err != nil {
			panic(err)
		}

		mapping, err := os.Open(source.Mapping)
		if err != nil {
			panic(err)
		}
		defer mapping.Close()

		request := esapi.IndicesCreateRequest{
			Index:  source.Index,
			Body:   mapping,
			Pretty: true,
		}

		response, err := request.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}

		if response.IsError() {
			panic(response.String())
		}
	}
}

// Remove ...
func (e Elasticsearch) Remove(cfg *config.Elasticsearch) {

	for _, source := range cfg.Sources {

		es, err := elasticsearch.NewDefaultClient()
		if err != nil {
			panic(err)
		}

		request := esapi.IndicesDeleteRequest{
			Index: []string{source.Index},
		}

		response, err := request.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}

		if response.IsError() {
			panic(response.String())
		}
	}
}

// RemoveIndex ...
func (e Elasticsearch) RemoveIndex(cfg *config.Elasticsearch, index string) {

	indices := make(map[string]struct{})

	for _, index := range strings.Split(index, ",") {
		if _, ok := cfg.Sources[index]; !ok {
			log.Printf("Index `%s` not registered", index)
			continue
		}
		indices[index] = struct{}{}
	}

	for _, source := range cfg.Sources {

		if _, ok := indices[source.Index]; !ok {
			continue
		}

		es, err := elasticsearch.NewDefaultClient()
		if err != nil {
			panic(err)
		}

		request := esapi.IndicesDeleteRequest{
			Index: []string{source.Index},
		}

		response, err := request.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}

		if response.IsError() {
			panic(response.String())
		}
	}
}
