package parse

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Tokens []string
	Patterns []string
}

type TypedConfig struct {
	Tokens []Token
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

	byteValue, _ := io.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(byteValue), &config)

	return config, nil
}
