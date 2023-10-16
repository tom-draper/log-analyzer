package loganalyzer

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
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

func tryPattern(line string, pattern Pattern) map[string]string {
	var regEx string = pattern.Pattern
	for _, token := range pattern.Tokens {
		// encode token value to create temporary token ID as hex as any
		// brackets in token may break regex
		tokenID := hex.EncodeToString([]byte(token.Value))
		regEx = strings.Replace(regEx, token.Value, fmt.Sprintf("(?P<%s>.*)", tokenID), 1)
	}
	encodedParams := getParams(line, regEx)

	// decode back to token value
	params := make(map[string]string)
	for tokenID, match := range encodedParams {
		tokenValue, err := hex.DecodeString(tokenID)
		if err == nil {
			params[string(tokenValue)] = match
		}
	}

	return params
}

func Parse(logtext string, config Config) error {
	result := make([]map[string]string, 0)
	for _, line := range strings.Split(logtext, "\n") {
		best := make(map[string]string)
		for _, pattern := range config {
			params := tryPattern(line, pattern)
			if len(params) > len(best) {
				best = params
			}
		}
		result = append(result, best)
	}
	log.Println(result)
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
