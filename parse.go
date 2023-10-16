package loganalyzer

import (
	"log"
	"os"
)

func Parse(logtext string, config Config) error {
	log.Println(logtext)
	return nil
}

func ParseFile(path string, config Config) error {
	body, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
		return err
	}
	return Parse(string(body), config)
}
