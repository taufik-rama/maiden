package internal

import (
	"io/ioutil"
	"log"
	"os"
)

// Print according to the verbosity flag
func Print(format string, v ...interface{}) {
	if Args.Verbose {
		log.Printf(format, v...)
	}
}

// Read the file contents
func Read(filename string) ([]byte, error) {
	Print("Reading `%s`", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
