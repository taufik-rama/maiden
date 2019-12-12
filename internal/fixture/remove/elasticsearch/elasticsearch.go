package elasticsearch

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/taufik-rama/maiden/config"
	"github.com/taufik-rama/maiden/config/fixtures"
)

type elasticsearch struct {
	config   config.Config
	fixtures fixtures.Config
}

// Sets the `config` field
func (e *elasticsearch) setConfig(c config.Config, err error) error {
	e.config = c
	return err
}

// Sets the `fixtures` field
func (e *elasticsearch) setFixtures(c fixtures.Config, err error) error {
	e.fixtures = c
	return err
}

// Remove every index defined on the fixtures directory.
func (e elasticsearch) remove() error {
	log.Println("Note: pushing without the index name is currently unimplemented :(")
	return nil
}

// Remove the index
func (e elasticsearch) removeIndex(index string) error {

	url := e.config.Fixtures().Elasticsearch().PushTo() + index
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("`%s` already removed (status code %d)", url, res.StatusCode)
	}

	request, err := http.NewRequest("DELETE", (e.config.Fixtures().Elasticsearch().PushTo() + index), strings.NewReader(""))
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
