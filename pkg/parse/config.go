package parse

import (
	"encoding/json"
	"io"
	"os"
)

type Config []Pattern

type Pattern struct {
	Pattern string   `json:"pattern"`
	Tokens  []string `json:"tokens"` // all unique tokens used in pattern
}

type TypedConfig []TypedPattern

type TypedPattern struct {
	Pattern string  `json:"pattern"`
	Tokens  []Token `json:"tokens"` // all unique tokens used in pattern
}

type Token struct {
	Value string `json:"value"`
	Type  string `json:"type,omitempty"` // defaults to string
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

	return config, nil
}
