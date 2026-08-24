package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type Config struct {
	Tokens       []string              `json:"tokens"`
	Patterns     []string              `json:"patterns"`
	Dependencies map[string][]string   `json:"dependencies,omitempty"`
	Conversions  map[string]Conversion `json:"conversions,omitempty"`
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
	decoder := json.NewDecoder(bytes.NewReader(byteValue))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	if err := ValidateConfig(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// ValidateConfig checks configuration references before parsing begins, so
// typos produce useful startup errors rather than silent missing dashboard data.
func ValidateConfig(config Config) error {
	if len(config.Patterns) == 0 {
		return fmt.Errorf("config must define at least one pattern")
	}

	tokens := make(map[string]struct{}, len(config.Tokens))
	for _, token := range config.Tokens {
		if strings.TrimSpace(token) == "" {
			return fmt.Errorf("token names must not be empty")
		}
		if token == "*" {
			return fmt.Errorf("token %q is reserved as the wildcard", token)
		}
		if _, exists := tokens[token]; exists {
			return fmt.Errorf("token %q is declared more than once", token)
		}
		tokens[token] = struct{}{}
	}

	patterns := make(map[string]struct{}, len(config.Patterns))
	for _, pattern := range config.Patterns {
		if strings.TrimSpace(pattern) == "" {
			return fmt.Errorf("patterns must not be empty")
		}
		if _, exists := patterns[pattern]; exists {
			return fmt.Errorf("pattern %q is declared more than once", pattern)
		}
		patterns[pattern] = struct{}{}
	}

	for token, dependencies := range config.Dependencies {
		if _, exists := tokens[token]; !exists {
			return fmt.Errorf("dependency token %q is not declared", token)
		}
		seen := make(map[string]struct{}, len(dependencies))
		for _, dependency := range dependencies {
			if _, exists := tokens[dependency]; !exists {
				return fmt.Errorf("dependency %q for token %q is not declared", dependency, token)
			}
			if dependency == token {
				return fmt.Errorf("token %q cannot depend on itself", token)
			}
			if _, exists := seen[dependency]; exists {
				return fmt.Errorf("dependency %q for token %q is declared more than once", dependency, token)
			}
			seen[dependency] = struct{}{}
		}
	}

	for source, conversion := range config.Conversions {
		if _, exists := tokens[source]; !exists {
			return fmt.Errorf("conversion source token %q is not declared", source)
		}
		if _, exists := tokens[conversion.Token]; !exists {
			return fmt.Errorf("conversion target token %q is not declared", conversion.Token)
		}
		if source == conversion.Token {
			return fmt.Errorf("conversion source and target must differ for token %q", source)
		}
		if conversion.Multiplier == 0 || math.IsNaN(conversion.Multiplier) || math.IsInf(conversion.Multiplier, 0) {
			return fmt.Errorf("conversion multiplier for token %q must be a finite non-zero number", source)
		}
	}

	return nil
}
