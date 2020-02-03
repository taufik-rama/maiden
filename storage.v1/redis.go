package storage

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
)

// Redis ...
type Redis struct{}

// Push ...
func (r Redis) Push(cfg *config.Redis) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Destination,
	})

	internal.Print("Pinging Redis")
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	internal.Print("Reading `%s`", cfg.Source)
	files, err := ioutil.ReadDir(cfg.Source)
	if err != nil {
		panic(err)
	}

	for _, file := range files {

		bytes, err := internal.Read(cfg.Source + string(os.PathSeparator) + file.Name())
		if err != nil {
			panic(err)
		}

		for index, line := range strings.Split(string(bytes), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			commands := strings.Split(line, " ")
			input := make([]interface{}, len(commands))
			for index, command := range commands {
				input[index] = command
			}
			if err := client.Do(input...).Err(); err != nil {
				log.Printf("Error on redis fixtures `%s` line %d: %s", (cfg.Source + file.Name()), (index + 1), err)
			}
		}
	}
}

// Remove ...
func (r Redis) Remove(cfg *config.Redis) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Destination,
	})

	internal.Print("Pinging Redis")
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	internal.Print("Flushing Redis")
	if _, err := client.FlushAll().Result(); err != nil {
		panic(err)
	}
}
