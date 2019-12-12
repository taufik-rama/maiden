package structure

import (
	"io/ioutil"
	"os"

	"github.com/taufik-rama/maiden/internal/config"
)

// Read the directory contents recursively
func readdir(dirname string) ([]string, error) {
	config.Print("Reading %s", dirname)
	var filenames []string
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := dirname + string(os.PathSeparator) + file.Name()
		if file.IsDir() {
			files, err := readdir(filename)
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, files...)
			continue
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}
