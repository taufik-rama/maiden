package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/taufik-rama/maiden/internal/config"
)

type elasticsearch struct {
	config config.MaidenConfig

	// what index to push to
	index string

	// use the documents values, ignore values given by maiden
	useID   bool
	useTime bool
}

func (e *elasticsearch) setConfig(cfg config.MaidenConfig) {
	e.config = cfg
}

func (e elasticsearch) push() error {
	fmt.Println("Push all")
	return nil
}

func (e elasticsearch) pushIndex(index string, reader io.Reader, info map[string]fieldsInfo) error {

	url := e.config.Fixtures.Elasticsearch.ParsedDestination[0] + index
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		config.Print("`%s` already exists (status code %d)", url, res.StatusCode)
		return nil
	}

	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var mapping map[string]struct {
		Mappings map[string]interface{} `json:"mappings"`
		Settings struct {
			Index map[string]interface{} `json:"index"`
		} `json:"settings"`
	}
	if err := json.Unmarshal(contents, &mapping); err != nil {
		return err
	}

	// Because the root name is dynamic, we treat it as a
	// string and iterate through it
	for _, val := range mapping {

		// Some auto-generated value that should not be included
		delete(val.Settings.Index, "uuid")
		delete(val.Settings.Index, "sort")
		delete(val.Settings.Index, "merge")
		delete(val.Settings.Index, "version")
		delete(val.Settings.Index, "creation_date")
		delete(val.Settings.Index, "provided_name")

		if !e.useTime {
			if !parseInfo(info, val) {
				goto done
			}
		}

	done:

		contents, err := json.MarshalIndent(val, "", "\t")
		if err != nil {
			return err
		}

		request, err := http.NewRequest("PUT", url, bytes.NewReader(contents))
		request.Header.Add("Content-Type", "application/json")
		if err != nil {
			return err
		}

		client := http.Client{}
		res, err = client.Do(request)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		contents, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var response struct {
			Error struct {
				Reason string `json:"reason"`
			} `json:"error"`
			Status int `json:"status"`
		}

		if err := json.Unmarshal(contents, &response); err != nil {
			return err
		}

		if response.Status != 0 {
			return fmt.Errorf("elasticsearch index error on `%s`: '%s' (status code %d)", url, response.Error.Reason, response.Status)
		}
	}

	return nil
}

func (e elasticsearch) pushDocuments(index string, reader io.Reader, id *int, info map[string]fieldsInfo) error {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var mapping struct {
		Hits struct {
			Hits []struct {
				ID     string                 `json:"_id"`
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(contents, &mapping); err != nil {
		return err
	}

	var buffer bytes.Buffer

	for _, hit := range mapping.Hits.Hits {

		if !e.useID {
			*id++
			if _, ok := hit.Source["id"]; ok {
				hit.Source["id"] = *id
			} else if _, ok := hit.Source["Id"]; ok {
				hit.Source["Id"] = *id
			} else if _, ok := hit.Source["ID"]; ok {
				hit.Source["ID"] = *id
			}
		} else {
			*id, _ = strconv.Atoi(hit.ID)
		}

		metadata := fmt.Sprintf(`{"index":{"_index":"%s","_type":"default","_id":"%d"}}`, index, *id)
		buffer.WriteString(metadata + "\n")

		if !e.useTime {

			source := hit.Source

			for key := range source {

				if source[key] == nil {
					continue
				}

				kind := reflect.TypeOf(source[key]).Kind()

				// Skip if not found AND is an end product field (does not have a subfield)
				if _, ok := info[key]; !ok && kind != reflect.Map {
					continue
				}

				currentTime := time.Now()

				if kind == reflect.String {

					var pushedTime time.Time

					if strings.HasPrefix(source[key].(string), "+") {

						numStr := strings.TrimLeft(source[key].(string), "+")
						if num, err := strconv.Atoi(numStr); err == nil {
							pushedTime = currentTime.Add(time.Hour * 24 * time.Duration(num))
						}

					} else if strings.HasPrefix(source[key].(string), "-") {

						numStr := strings.TrimLeft(source[key].(string), "-")
						if num, err := strconv.Atoi(numStr); err == nil {
							pushedTime = currentTime.Add(time.Hour * 24 * -time.Duration(num))
						}

					} else {

						documentTime, err := time.Parse("2006-01-02T15:04:05Z", source[key].(string))
						if err != nil {
							pushedTime = currentTime
						} else {
							if documentTime.Before(currentTime) {
								pushedTime = currentTime
							} else {
								pushedTime = documentTime
							}
						}
					}

					source[key] = pushedTime.Format("2006-01-02T15:04:05Z")

				} else if kind == reflect.Map {

					props := source[key].(map[string]interface{})

					for propKey := range props {

						if _, ok := info[propKey]; !ok {
							continue
						}

						kind := reflect.TypeOf(props[propKey]).Kind()
						if kind == reflect.String {

							var pushedTime time.Time

							if strings.HasPrefix(props[propKey].(string), "+") {

								numStr := strings.TrimLeft(props[propKey].(string), "+")
								if num, err := strconv.Atoi(numStr); err == nil {
									pushedTime = currentTime.Add(time.Hour * 24 * time.Duration(num))
								}

							} else if strings.HasPrefix(props[propKey].(string), "-") {

								numStr := strings.TrimLeft(props[propKey].(string), "+")
								if num, err := strconv.Atoi(numStr); err == nil {
									pushedTime = currentTime.Add(time.Hour * 24 * -time.Duration(num))
								}

							} else {

								documentTime, err := time.Parse("2006-01-02T15:04:05Z", props[propKey].(string))
								if err != nil {
									pushedTime = currentTime
								} else {
									if documentTime.Before(currentTime) {
										pushedTime = currentTime
									} else {
										pushedTime = documentTime
									}
								}
							}

							props[propKey] = pushedTime.Format("2006-01-02T15:04:05Z")
						}
					}

					source[key] = props
				}
			}

			hit.Source = source
		}

		b, err := json.Marshal(hit.Source)
		if err != nil {
			return err
		}
		buffer.WriteString(string(b) + "\n")
	}

	// Don't post if there's no data
	if buffer.Len() == 0 {
		return nil
	}

	res, err := http.Post(fmt.Sprintf("%s_bulk", e.config.Fixtures.Elasticsearch.ParsedDestination[0]), "application/json", strings.NewReader(buffer.String()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	contents, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	type responseitems struct {
		Index struct {
			ID     string `json:"_id"`
			Status int    `json:"status"`
			Error  struct {
				Reason string `json:"reason"`
			} `json:"error"`
		} `json:"index"`
	}

	var response struct {
		Errors bool            `json:"errors"`
		Items  []responseitems `json:"items"`
	}

	if err := json.Unmarshal(contents, &response); err != nil {
		return err
	}

	if response.Errors {
		invalids := []responseitems{}
		for _, item := range response.Items {
			if len(invalids) > 3 {
				break
			}
			if item.Index.Status < 200 || item.Index.Status > 299 {
				invalids = append(invalids, item)
			}
		}
		return fmt.Errorf("elasticsearch documents error for index `%s`, raw values: '%v'", index, invalids)
	}

	return nil
}

type fieldsInfo struct {
	isDate bool
}

func parseInfo(info map[string]fieldsInfo, val struct {
	Mappings map[string]interface{} `json:"mappings"`
	Settings struct {
		Index map[string]interface{} `json:"index"`
	} `json:"settings"`
}) bool {

	for mapping := range val.Mappings {

		kind := reflect.TypeOf(val.Mappings[mapping]).Kind()
		if kind != reflect.Map {
			log.Printf("Field `%s` is not an object: `%s`, skipping date parse", mapping, kind)
			return false
		}
		mapping := val.Mappings[mapping].(map[string]interface{})
		properties, ok := mapping["properties"]
		if !ok {
			log.Printf("Field `properties` does not exists inside `%s`, skipping date parse", mapping)
			return false
		}

		kind = reflect.TypeOf(properties).Kind()
		if kind != reflect.Map {
			log.Printf("Field `properties` is not an object: `%s` inside `%s`, skipping date parse", mapping, kind)
			return false
		}
		props := properties.(map[string]interface{})
		for property := range props {

			var prop struct {
				Type       interface{} `json:"type"`
				Fields     interface{} `json:"fields"`
				Properties interface{} `json:"properties"`
			}

			bytes, err := json.Marshal(props[property])
			if err != nil {
				log.Printf("Error while parsing field `%s`: `%s`", property, err)
			}
			json.Unmarshal(bytes, &prop)

			if prop.Fields != nil || prop.Properties != nil {

				var properties interface{}
				if prop.Fields != nil {
					properties = prop.Fields
				} else {
					properties = prop.Properties
				}

				kind := reflect.TypeOf(properties).Kind()
				if kind != reflect.Map {
					log.Printf("Field `%s` is not an object: `%s`, skipping date parse", properties, kind)
					return false
				}

				propMap := properties.(map[string]interface{})

				for propKey := range propMap {

					var prop struct {
						Type interface{} `json:"type"`
					}

					bytes, err := json.Marshal(propMap[propKey])
					if err != nil {
						log.Printf("Error while parsing field `%s`: `%s`", propKey, err)
					}
					json.Unmarshal(bytes, &prop)

					if prop.Type == "date" {
						info[propKey] = fieldsInfo{
							isDate: true,
						}
					}
				}

			} else {

				if prop.Type == "date" {
					info[property] = fieldsInfo{
						isDate: true,
					}
				}
			}
		}
	}

	return true
}
