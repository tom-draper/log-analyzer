package parse

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
)

// getParams extracts all possible group values contained within the regular
// expression from the text and stores the extracted values in a returned
// (groupName => value) map.
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

// Replace all characters that have a special meaning within a regular
// expression with an escaped version of the character.
func escapeRegexCharacters(regEx string) string {
	regEx = strings.ReplaceAll(regEx, "(", "\\(")
	regEx = strings.ReplaceAll(regEx, ")", "\\)")
	regEx = strings.ReplaceAll(regEx, "]", "\\]")
	regEx = strings.ReplaceAll(regEx, "[", "\\[")
	return regEx
}

// tryPattern attempts to extract the corresponding token values described by
// the given pattern from the log text line. Any extracted values have their
// data types inferred and then converted.
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
	typedParams := inferDataTypes(params)

	return typedParams
}

// inferDataTypes infers the intended the data type of each extracted parameter
// value from a log text line and performs a data type conversion.
func inferDataTypes(params map[string]string) map[string]any {
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

// parseLine extracts token parameters from each line using the most appropriate
// pattern in the given config.
func parseLine(line string, config Config) (map[string]any, string) {
	// Attempt to parse the line against each pattern in config, only taking the best
	var patternUsed string
	best := make(map[string]any)
	for _, pattern := range config.Patterns {
		params := tryPattern(line, pattern, config.Tokens)
		if len(params) > len(best) {
			best = params
			patternUsed = pattern
		}
	}
	if len(best) == 0 {
		log.Printf("no pattern matched line: %s\n", line)
	}
	return best, patternUsed
}

func splitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}

// Parse separates the log text into lines and attempts to extract tokens
// parameters from each line using the most appropriate pattern in the given config.
func Parse(logtext string, config Config) ([]map[string]any, error) {
	params := make([]map[string]any, 0)
	for _, line := range splitLines(logtext) {
		p, _ := parseLine(line, config)
		params = append(params, p)
	}
	return params, nil
}

// ParseFile reads the log text from the given file path, separates the text
// into lines and attempts to extract tokens parameters from each line using the
// most appropriate pattern in the given config.
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

// ParseFile reads the log text from each of the given file paths, separates the
// text into lines and attempts to extract tokens parameters from each line
// using the most appropriate pattern in the given config.
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

type ExtractResult struct {
	Line       string         `json:"line"`
	LineNumber int            `json:"lineNumber"`
	Pattern    string         `json:"pattern"`
	Params     map[string]any `json:"params"`
}

// ParseTest runs Parse and displays a random sample of extracted parameters
// along with the origin lines from the log text.
func ParseTest(logtext string, config Config) []ExtractResult {
	results := make([]ExtractResult, 0)
	for i, line := range splitLines(logtext) {
		p, pattern := parseLine(line, config)
		result := ExtractResult{
			Line:       line,
			LineNumber: i,
			Pattern:    pattern,
			Params:     p,
		}
		results = append(results, result)
	}

	return results
}

// ParseTest runs ParseFile and displays a random sample of extracted parameters
// along with the origin lines from the log file.
func ParseFileTest(path string, config Config) ([]ExtractResult, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	results := ParseTest(string(body), config)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ParseTest runs ParseFiles and displays a random sample of extracted
// parameters along with the origin lines from the log files.
func ParseFilesTest(paths []string, config Config) ([]ExtractResult, error) {
	results := make([]ExtractResult, 0)

	var parsedAny bool
	for _, path := range paths {
		r, err := ParseFileTest(path, config)
		if err != nil {
			log.Printf("unable to read file at path %s: %s", path, fmt.Sprint(err))
			continue
		}
		parsedAny = true
		results = append(results, r...)
	}

	if !parsedAny {
		return nil, errors.New("unable to read log file path provided")
	}

	return results, nil
}
