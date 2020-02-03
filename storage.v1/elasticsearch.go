package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
)

// Elasticsearch ...
type Elasticsearch struct {

	// to keep track of ID across document file
	id int
}

// Push ...
func (e Elasticsearch) Push(cfg *config.Elasticsearch) {

	for _, source := range cfg.Sources {

		internal.Print("Reading `%s`", source.Mapping)
		mapping, err := os.Open(source.Mapping)
		if err != nil {
			panic(err)
		}
		defer mapping.Close()

		request, err := http.NewRequest("PUT", (cfg.Destination + "/" + source.Index), mapping)
		if err != nil {
			panic(err)
		}
		request.Header.Add("Content-Type", "application/json")

		response, err := new(http.Client).Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		var cause struct {
			Error struct {
				Type    string `json:"type"`
				Reasong string `json:"reason"`
			} `json:"error"`
		}
		bytes, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(bytes, &cause)

		if strings.TrimSpace(cause.Error.Type) != "" {
			log.Printf("%s: [%s] %s", response.Status, cause.Error.Type, cause.Error.Reasong)
		}

		if err := e.pushDocuments(source, cfg.Destination); err != nil {
			panic(err)
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

		internal.Print("Reading `%s`", source.Mapping)
		mapping, err := os.Open(source.Mapping)
		if err != nil {
			panic(err)
		}
		defer mapping.Close()

		request, err := http.NewRequest("PUT", (cfg.Destination + "/" + source.Index), mapping)
		if err != nil {
			panic(err)
		}
		request.Header.Add("Content-Type", "application/json")

		response, err := new(http.Client).Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		var cause struct {
			Error struct {
				Type    string `json:"type"`
				Reasong string `json:"reason"`
			} `json:"error"`
		}
		bytes, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(bytes, &cause)

		if strings.TrimSpace(cause.Error.Type) != "" {
			log.Printf("%s: [%s] %s", response.Status, cause.Error.Type, cause.Error.Reasong)
		}

		if err := e.pushDocuments(source, cfg.Destination); err != nil {
			panic(err)
		}
	}
}

func (e Elasticsearch) pushDocuments(source config.ElasticsearchSource, destination string) error {

	return filepath.Walk(source.Files, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			panic(err)
		}
		if info.IsDir() {
			return nil
		}

		var documents []interface{}

		{
			bytes, err := internal.Read(path)
			if err != nil {
				panic(err)
			}
			bytes, err = e.parseScript(bytes)
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(bytes, &documents); err != nil {
				panic(err)
			}
		}

		metadata := fmt.Sprintf(`{"index": {"_index": "%s", "_id": "%%d"}}`, source.Index)

		buffer := bytes.NewBuffer(nil)

		for _, document := range documents {
			e.id++
			buffer.WriteString((fmt.Sprintf(metadata, e.id) + "\n"))
			bytes, err := json.Marshal(document)
			if err != nil {
				panic(err)
			}
			buffer.Write(bytes)
			buffer.WriteString("\n")
		}

		request, err := http.NewRequest("POST", (destination + "/" + source.Index + "/" + source.MappingType + "/_bulk"), buffer)
		if err != nil {
			panic(err)
		}
		request.Header.Add("Content-Type", "application/x-ndjson")

		response, err := new(http.Client).Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		var cause struct {
			Error struct {
				Type    string `json:"type"`
				Reasong string `json:"reason"`
			} `json:"error"`
		}
		bytes, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(bytes, &cause)

		if strings.TrimSpace(cause.Error.Type) != "" {
			log.Printf("%s: [%s] %s", response.Status, cause.Error.Type, cause.Error.Reasong)
		}

		return nil
	})
}

// Remove ...
func (e Elasticsearch) Remove(cfg *config.Elasticsearch) {

	for _, source := range cfg.Sources {

		request, err := http.NewRequest("DELETE", (cfg.Destination + "/" + source.Index), nil)
		if err != nil {
			panic(err)
		}

		response, err := new(http.Client).Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		var cause struct {
			Error struct {
				Type    string `json:"type"`
				Reasong string `json:"reason"`
			} `json:"error"`
		}
		bytes, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(bytes, &cause)

		if strings.TrimSpace(cause.Error.Type) != "" {
			log.Printf("%s: [%s] %s", response.Status, cause.Error.Type, cause.Error.Reasong)
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

		request, err := http.NewRequest("DELETE", (cfg.Destination + "/" + source.Index), nil)
		if err != nil {
			panic(err)
		}

		response, err := new(http.Client).Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		var cause struct {
			Error struct {
				Type    string `json:"type"`
				Reasong string `json:"reason"`
			} `json:"error"`
		}
		bytes, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(bytes, &cause)

		if strings.TrimSpace(cause.Error.Type) != "" {
			log.Printf("%s: [%s] %s", response.Status, cause.Error.Type, cause.Error.Reasong)
		}
	}
}

// Unmarshal the `contents` into an appropriate type & parse for script values
func (e Elasticsearch) parseScript(contents []byte) ([]byte, error) {
	var fields []interface{}
	if err := json.Unmarshal(contents, &fields); err != nil {
		return nil, fmt.Errorf("error while unmarshaling script: %s", err)
	}
	for i := range fields {
		if _, ok := fields[i].(map[string]interface{}); !ok {
			return contents, fmt.Errorf("invalid structure: must be an object on document index %d", i)
		}
		e.parseScriptFields(fields[i].(map[string]interface{}))
	}
	return json.Marshal(fields)
}

// Recursively do a `script()` call using the `fields` value, mutates `fields`
func (e Elasticsearch) parseScriptFields(fields map[string]interface{}) {

	for field := range fields {

		// `subfields` will refer to `fields`, so there's no need to copy
		// the map data
		if subfield, ok := fields[field].(map[string]interface{}); ok {
			e.parseScriptFields(subfield)
			continue
		}

		if val, ok := fields[field].(string); ok && e.isScript(val) {
			fields[field] = e.script(e.trimScript(val))
		}
	}
}

// Custom script for elasticsearch document value.
// Return value will be the final value
func (e Elasticsearch) script(val string) interface{} {
	if strings.HasPrefix(val, "date") {
		offset, _ := strconv.Atoi(strings.TrimLeft(val, "date"))
		return time.Now().Add(time.Duration(offset) * time.Hour).Format("2006-01-02T15:04:05Z")
	}
	return val
}

// Trim the enclosing character
func (e Elasticsearch) trimScript(val string) string {
	return strings.TrimRight(strings.TrimLeft(val, "$("), ")")
}

// What we'll consider a script
func (e Elasticsearch) isScript(val string) bool {
	return strings.HasPrefix(val, "$(") && strings.HasSuffix(val, ")")
}
