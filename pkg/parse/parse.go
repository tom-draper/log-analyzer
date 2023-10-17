package parse

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/araddon/dateparse"
)

func getParams(text string, regEx string) map[string]string {
	compRegEx := regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(text)

	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func tryPattern(line string, pattern Pattern) map[string]any {
	var regEx string = pattern.Pattern
	for _, token := range pattern.Tokens {
		// Encode token value to create temporary token ID as hex as any
		// brackets in token may break regex
		tokenID := hex.EncodeToString([]byte(token.Value))
		regEx = strings.Replace(regEx, token.Value, fmt.Sprintf("(?P<%s>.*)", tokenID), 1)
	}
	encodedParams := getParams(line, regEx)

	// Decode back to token value
	params := make(map[string]string)
	for tokenID, match := range encodedParams {
		tokenValue, err := hex.DecodeString(tokenID)
		if err == nil {
			params[string(tokenValue)] = match
		}
	}

	typedParams := applyDataTypes(params, pattern)

	return typedParams
}

func applyDataTypes(params map[string]string, pattern Pattern) map[string]any {
	typedParams := make(map[string]any)
	for tokenValue, match := range params {
		for _, token := range pattern.Tokens {
			if tokenValue != token.Value {
				continue
			}
			// If this token is for a timestamp value
			if token.Timestamp {
				// Attempt to parse as datetime
				t, err := dateparse.ParseAny(match)
				if err == nil {
					typedParams[tokenValue] = t
					continue
				}
			}
			typedParams[tokenValue] = match
		}
	}
	return typedParams
}

func parseLine(line string, config Config) map[string]any {
	// Attempt to parse the line against each pattern in config, only taking the best
	best := make(map[string]any)
	for _, pattern := range config {
		params := tryPattern(line, pattern)
		if len(params) > len(best) {
			best = params
		}
	}
	return best
}

func Parse(logtext string, config Config) ([]map[string]any, error) {
	extracted := make([]map[string]any, 0)
	for _, line := range strings.Split(strings.ReplaceAll(logtext, "\r\n", "\n"), "\n") {
		params := parseLine(line, config)
		extracted = append(extracted, params)
	}
	return extracted, nil
}

func ParseFile(path string, config Config) ([]map[string]any, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	params, err := Parse(string(body), config)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func ParseFiles(paths []string, config Config) ([]map[string]any, error) {
	allParams := make([]map[string]any, 0)
	for _, path := range paths {
		fileParams, err := ParseFile(path, config)
		if err != nil {
			return nil, fmt.Errorf("unable to parse file at path %s: %w", path, err)
		}
		allParams = append(allParams, fileParams...)
	}
	return allParams, nil
}
