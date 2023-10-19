package parse

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Tokens   []string
	Patterns []string
}

type TypedConfig struct {
	Tokens   []Token
	Patterns []string
}

type Token struct {
	Value string `json:"value"`
	Type  string `json:"type,omitempty"` // defaults to string
}

func LoadConfig(path string) (Config, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal([]byte(byteValue), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
