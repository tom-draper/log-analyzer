package loganalyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config []Pattern

type Pattern struct {
	Pattern string  `json:"pattern"`
	Tokens  []Token `json:"tokens"` // all unique tokens used in pattern
}

type Token struct {
	Value          string `json:"value"`
	Timestamp      bool   `json:"timestamp,omitempty"`
	DateTimeFormat string `json:"dateTimeFormat,omitempty"`
}

func LoadConfig(path string) (Config, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(byteValue), &config)

	fmt.Println(config)
	return config, nil
}
