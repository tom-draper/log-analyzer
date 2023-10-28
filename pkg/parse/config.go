package parse

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Tokens       []string              `json:"tokens"`
	Patterns     []string              `json:"patterns"`
	Dependencies map[string][]string   `json:"dependencies,omitempty"`
	Conversions  map[string]Conversion `json:"conversions,omitempty"`
}

type TypedConfig struct {
	Tokens       []Token               `json:"tokens"`
	Patterns     []string              `json:"patterns"`
	Dependencies map[string][]string   `json:"dependencies,omitempty"`
	Conversions  map[string]Conversion `json:"conversions,omitempty"`
}

type Token struct {
	Value string `json:"value"`
	Type  string `json:"type,omitempty"` // defaults to string
}

type Conversion struct {
	Token      string  `json:"token"`
	Multiplier float64 `json:"multiplier"`
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
