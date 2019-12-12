package fixture

import "reflect"

const stats = `Fixtures status
  Cassandra
    source: %s
    destination: %s
  Elasticsearch
    source: %s
    destination: %s
  PostgreSQL
    source: %s
    destination: %s
  Redis
    source: %s
    destination: %s
`

func format(val interface{}) string {
	var stat string
	kind := reflect.TypeOf(val).Kind()
	if kind == reflect.String {
		stat = val.(string)
	} else {
		values := val.([]interface{})
		prefix := "      "
		for _, value := range values {
			stat += "\n"
			stat += (prefix + "- " + value.(string))
		}
	}
	return stat
}
