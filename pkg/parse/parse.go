package parse

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
)

func getParams(text string, regEx string) map[string]string {
	compRegEx := regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(text)

	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = strings.TrimSpace(match[i])
		}
	}
	return paramsMap
}

func escapeRegexCharacters(regEx string) string {
	// Replace all brackets
	regEx = strings.ReplaceAll(regEx, "(", "\\(")
	regEx = strings.ReplaceAll(regEx, ")", "\\)")
	regEx = strings.ReplaceAll(regEx, "]", "\\]")
	regEx = strings.ReplaceAll(regEx, "[", "\\[")
	return regEx
}

func tryPattern(line string, pattern string, tokens []string) map[string]any {
	var regEx string = pattern
	regEx = escapeRegexCharacters(regEx)
	for _, token := range tokens {
		// Encode token value to create temporary token ID as hex as any
		// brackets in token may break regex
		t := escapeRegexCharacters(token)
		if !strings.Contains(regEx, t) {
			continue
		}
		tokenID := hex.EncodeToString([]byte(t))
		regEx = strings.Replace(regEx, t, fmt.Sprintf("(?P<%s>.*)", tokenID), 1)
	}
	encodedParams := getParams(line, regEx)

	// Decode back to raw token value
	params := make(map[string]string)
	for tokenID, match := range encodedParams {
		token, err := hex.DecodeString(tokenID)
		if err == nil {
			params[string(token)] = match
		}
	}

	// Attempt to infer data types
	typedParams := parseDataTypes(params)

	return typedParams
}

func parseDataTypes(params map[string]string) map[string]any {
	typedParams := make(map[string]any)
	for token, match := range params {
		// Attempt to parse as datetime
		if t, err := dateparse.ParseAny(match); err == nil {
			typedParams[token] = t
		} else if value, err := strconv.ParseFloat(match, 64); strings.Contains(match, ".") && err == nil {
			typedParams[token] = value
		} else if value, err := strconv.Atoi(match); err == nil {
			typedParams[token] = value
		} else if value, err := strconv.ParseBool(match); err == nil {
			typedParams[token] = value
		} else {
			typedParams[token] = match
		}
	}
	return typedParams
}

func parseLine(line string, config Config) map[string]any {
	// Attempt to parse the line against each pattern in config, only taking the best
	best := make(map[string]any)
	for _, pattern := range config.Patterns {
		params := tryPattern(line, pattern, config.Tokens)
		if len(params) > len(best) {
			best = params
		}
	}
	if len(best) == 0 {
		log.Printf("no pattern matched line: %s\n", line)
	}
	return best
}

func Parse(logtext string, config Config) ([]map[string]any, error) {
	params := make([]map[string]any, 0)
	for _, line := range strings.Split(strings.ReplaceAll(logtext, "\r\n", "\n"), "\n") {
		p := parseLine(line, config)
		params = append(params, p)
	}
	return params, nil
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
	params := make([]map[string]any, 0)
	for _, path := range paths {
		fileParams, err := ParseFile(path, config)
		if err != nil {
			return nil, fmt.Errorf("unable to parse file at path %s: %w", path, err)
		}
		params = append(params, fileParams...)
	}
	return params, nil
}
