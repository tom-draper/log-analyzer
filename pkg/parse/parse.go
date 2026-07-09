package parse

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/araddon/dateparse"
)

type Extraction struct {
	Params     map[string]Param `json:"params"`
	Pattern    string           `json:"pattern"`
	LineNumber int              `json:"lineNumber"`
	Line       string           `json:"line"`
}

type Param struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}

var (
	multiSpaceRegEx = regexp.MustCompile(`[ ]{2,}`)
	regexCache      sync.Map // map[string]*regexp.Regexp
)

// getParams extracts all possible group values contained within the regular
// expression from the text and stores the extracted values in a returned
// (groupName => value) map.
func getParams(text string, regEx string) map[string]string {
	var compRegEx *regexp.Regexp
	if cached, ok := regexCache.Load(regEx); ok {
		compRegEx = cached.(*regexp.Regexp)
	} else {
		compRegEx = regexp.MustCompile(regEx)
		regexCache.Store(regEx, compRegEx)
	}
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
	regEx = strings.ReplaceAll(regEx, "“", "\"")
	return regEx
}

// patternMatchesLine reports whether line matches pattern after expanding
// wildcard tokens (*) into .*  Used for zero-token patterns where rank
// calculation cannot be used.
func patternMatchesLine(line string, pattern string) bool {
	regEx := escapeRegexCharacters(pattern)
	regEx = strings.ReplaceAll(regEx, "*", ".*")
	var compRegEx *regexp.Regexp
	if cached, ok := regexCache.Load(regEx); ok {
		compRegEx = cached.(*regexp.Regexp)
	} else {
		var err error
		compRegEx, err = regexp.Compile(regEx)
		if err != nil {
			return false
		}
		regexCache.Store(regEx, compRegEx)
	}
	return compRegEx.MatchString(line)
}

// tryPattern attempts to extract the corresponding token values described by
// the given pattern from the log text line. Any extracted values have their
// data types inferred and then converted.
func tryPattern(line string, pattern string, tokens []string) map[string]Param {
	var regEx string = pattern
	regEx = escapeRegexCharacters(regEx)
	// Convert wildcard asterisk into underscore _ so we only have to deal
	// with one wildcard char
	regEx = strings.ReplaceAll(regEx, "*", "_")
	// Work on a copy of the tokens slice; the append and in-place sort below
	// would otherwise mutate the shared config.Tokens, which is unsafe when
	// lines are parsed concurrently.
	tokens = append(append([]string{}, tokens...), "_")
	// Sort tokens to try largest first and avoid matching substrings
	sort.Slice(tokens, func(i, j int) bool {
		return len(tokens[i]) > len(tokens[j])
	})
	for _, token := range tokens {
		// Encode token value to create temporary token ID as hex as any
		// brackets in token may break regex
		t := escapeRegexCharacters(token)
		if !strings.Contains(regEx, t) {
			continue
		}

		tokenID := hex.EncodeToString([]byte(t))
		regEx = strings.ReplaceAll(regEx, t, fmt.Sprintf("(?P<%s>.*)", tokenID))
	}
	encodedParams := getParams(line, regEx)

	// Decode back to raw token value
	params := make(map[string]string)
	for tokenID, match := range encodedParams {
		token, err := hex.DecodeString(tokenID)
		if err != nil {
			continue
		}
		stoken := string(token)
		// Avoid adding to final parameters if was a wildcard token
		if stoken == "" || stoken[0] == '_' {
			continue
		}
		if _, ok := params[stoken]; !ok {
			params[stoken] = match
		}
	}

	// Attempt to infer data types
	typedParams := inferDataTypes(params)

	return typedParams
}

// inferDataTypes infers the intended the data type of each extracted parameter
// value from a log text line and performs a data type conversion.
func inferDataTypes(params map[string]string) map[string]Param {
	typedParams := make(map[string]Param)
	for token, match := range params {
		// Check for explicit type in token name
		if strings.HasPrefix(token, "int_") {
			if value, err := strconv.Atoi(match); err == nil {
				typedParams[token] = Param{Value: value, Type: "int"}
			}
			continue
		} else if strings.HasPrefix(token, "float_") {
			if value, err := strconv.ParseFloat(match, 64); strings.Contains(match, ".") && err == nil {
				typedParams[token] = Param{Value: value, Type: "float"}
			}
			continue
		} else if strings.HasPrefix(token, "time_") {
			if value, err := dateparse.ParseAny(match); err == nil {
				typedParams[token] = Param{Value: value, Type: "time"}
			}
			continue
		} else if strings.HasPrefix(token, "bool_") {
			if value, err := strconv.ParseBool(match); err == nil {
				typedParams[token] = Param{Value: value, Type: "bool"}
			}
			continue
		} else if strings.HasPrefix(token, "ip_") {
			if value := net.ParseIP(match); value != nil {
				typedParams[token] = Param{Value: value, Type: "ip"}
			}
			continue
		} else if strings.HasPrefix(token, "str_") {
			typedParams[token] = Param{Value: match, Type: "str"}
			continue
		}

		// Attempt to parse as datetime
		if value, err := dateparse.ParseAny(match); err == nil {
			typedParams[token] = Param{Value: value, Type: "time"}
		} else if value := net.ParseIP(match); value != nil {
			typedParams[token] = Param{Value: value, Type: "ip"}
		} else if value, err := strconv.ParseFloat(match, 64); strings.Contains(match, ".") && err == nil {
			typedParams[token] = Param{Value: value, Type: "float"}
		} else if value, err := strconv.Atoi(match); err == nil {
			typedParams[token] = Param{Value: value, Type: "int"}
		} else if value, err := strconv.ParseBool(match); err == nil {
			typedParams[token] = Param{Value: value, Type: "bool"}
		} else {
			typedParams[token] = Param{Value: match, Type: "str"}
		}
	}
	return typedParams
}

func tokenCounts(pattern string, tokens []string) int {
	var count int
	for _, token := range tokens {
		if strings.Contains(pattern, token) {
			count++
		}
	}
	return count
}

func calcExtractionRank(params map[string]Param, patternTokenCounts int) float64 {
	if patternTokenCounts == 0 {
		return 0
	}
	paramCount := float64(len(params))
	// Penalise parameter count in favour of proportion of tokens in pattern used
	rank := (paramCount * -0.05) + ((paramCount / float64(patternTokenCounts)) * 1)
	return rank
}

type PatternRank struct {
	rank   float64
	params map[string]Param
}

// parseLine extracts token parameters from each line using the most appropriate
// pattern in the given config.
func parseLine(lines []string, index int, config *Config) (map[string]Param, string) {
	// Attempt to parse the line against each pattern in config, only taking the best
	var patternUsed string
	best := PatternRank{
		rank:   0.0,
		params: make(map[string]Param),
	}
	line := lines[index]
	for _, pattern := range config.Patterns {
		lineCount := patternLineCount(pattern)
		// If not enough lines remaining for multi-line pattern
		if lineCount > 1 && index+lineCount > len(lines) {
			continue
		}

		var snippet string
		if lineCount == 1 {
			snippet = line
		} else {
			snippet = strings.Join(lines[index:min(len(lines)+1, index+lineCount)], "\n")
		}
		// If pattern contains no tokens is a perfect match for line,
		// ensure usage of this pattern is recorded even if rank would be zero
		// due to lack of tokens to base rank on
		patternTokenCounts := tokenCounts(pattern, config.Tokens)
		if snippet == pattern && patternTokenCounts == 0 {
			patternUsed = pattern // Force line to take this pattern
			break
		}

		// Measure how well this pattern fits the line(s)
		var lineRank PatternRank
		if lineCount == 1 {
			lineRank = parseSingleLine(line, pattern, config)
		} else {
			lineRank = parseMultiLine(lines, index, pattern, config)
		}
		// Record if this pattern is better than others seen so far
		if lineRank.rank > best.rank {
			best.rank = lineRank.rank
			best.params = lineRank.params
			patternUsed = pattern
		}
	}

	return best.params, patternUsed
}

func parseSingleLine(line string, pattern string, config *Config) PatternRank {
	lineRank := PatternRank{
		rank:   0.0,
		params: make(map[string]Param),
	}

	patternTokenCounts := tokenCounts(pattern, config.Tokens)

	// For zero-token patterns the rank formula always returns 0, so use a
	// direct regex match instead.  A small positive rank signals "matched".
	if patternTokenCounts == 0 {
		if patternMatchesLine(line, pattern) {
			lineRank.rank = 0.5
		}
		return lineRank
	}

	params := tryPattern(line, pattern, config.Tokens)
	rank := calcExtractionRank(params, patternTokenCounts)
	if rank > lineRank.rank {
		lineRank.rank = rank
		lineRank.params = params
	}

	// Try pattern again after eliminating multi-spaces and tab characters
	if multiSpaceRegEx.MatchString(line) {
		singleSpaceLine := multiSpaceRegEx.ReplaceAllString(strings.ReplaceAll(line, "\t", " "), " ")
		singleSpacePattern := multiSpaceRegEx.ReplaceAllString(strings.ReplaceAll(pattern, "\t", " "), " ")
		params = tryPattern(singleSpaceLine, singleSpacePattern, config.Tokens)
		rank := calcExtractionRank(params, patternTokenCounts)
		if rank > lineRank.rank {
			lineRank.rank = rank
			lineRank.params = params
		}
	}

	return lineRank
}

func parseMultiLine(lines []string, index int, pattern string, config *Config) PatternRank {
	patternLines := splitLines(pattern)
	lineRanks := make([]PatternRank, len(patternLines))
	for i, patternLine := range patternLines {
		line := lines[index+i]
		lineBest := parseSingleLine(line, patternLine, config)
		lineRanks[i] = lineBest
	}

	return avgRank(lineRanks)
}

func avgRank(ranks []PatternRank) PatternRank {
	var avg PatternRank = PatternRank{
		rank:   0.0,
		params: make(map[string]Param),
	}
	for _, rank := range ranks {
		avg.rank += rank.rank
		for k, v := range rank.params {
			avg.params[k] = v
		}
	}
	avg.rank = avg.rank / float64(len(ranks))
	return avg
}

func patternLineCount(pattern string) int {
	return strings.Count(pattern, "\n") + 1
}

func splitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}

// Parse separates the log text into lines and attempts to extract tokens
// parameters from each line using the most appropriate pattern in the given config.
func Parse(logtext string, config *Config) ([]Extraction, error) {
	lines := splitLines(logtext)
	return ParseLines(lines, config)
}

// hasMultiLinePatterns reports whether any pattern in the config spans more than
// one line. Multi-line patterns require the sequential parser so that lines
// consumed by a match are skipped correctly.
func hasMultiLinePatterns(config *Config) bool {
	for _, pattern := range config.Patterns {
		if strings.Contains(pattern, "\n") {
			return true
		}
	}
	return false
}

// ParseLines attempts to extract token parameters from each line using the most
// appropriate pattern in the given config. When the config contains no
// multi-line patterns the lines are parsed concurrently; otherwise a sequential
// pass is used so lines consumed by a multi-line match are skipped correctly.
// Both paths produce identical, in-order results.
func ParseLines(lines []string, config *Config) ([]Extraction, error) {
	if hasMultiLinePatterns(config) {
		return parseLinesSequential(lines, config)
	}
	return parseLinesConcurrent(lines, config)
}

func parseLinesSequential(lines []string, config *Config) ([]Extraction, error) {
	extraction := make([]Extraction, 0, len(lines))
	var skipLines int
	var unmatched int
	for i, line := range lines {
		if skipLines > 0 {
			skipLines--
			continue
		}
		if line == "" {
			continue
		}
		params, patternUsed := parseLine(lines, i, config)
		if patternUsed == "" {
			unmatched++
		}
		extraction = append(extraction, Extraction{
			Params:     params,
			Pattern:    patternUsed,
			LineNumber: i,
			Line:       line,
		})
		// If pattern used was multi-line, skip the next lines
		skipLines = patternLineCount(patternUsed) - 1
	}
	logUnmatched(unmatched, len(extraction))
	return extraction, nil
}

func parseLinesConcurrent(lines []string, config *Config) ([]Extraction, error) {
	extraction := make([]Extraction, len(lines))
	keep := make([]bool, len(lines))

	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan int, workers)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				line := lines[i]
				if line == "" {
					continue
				}
				params, patternUsed := parseLine(lines, i, config)
				extraction[i] = Extraction{
					Params:     params,
					Pattern:    patternUsed,
					LineNumber: i,
					Line:       line,
				}
				keep[i] = true
			}
		}()
	}
	for i := range lines {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	// Compact into an in-order slice, dropping skipped (empty) lines.
	result := make([]Extraction, 0, len(lines))
	var unmatched int
	for i := range extraction {
		if !keep[i] {
			continue
		}
		if extraction[i].Pattern == "" {
			unmatched++
		}
		result = append(result, extraction[i])
	}
	logUnmatched(unmatched, len(result))
	return result, nil
}

// logUnmatched prints a single summary line rather than one message per
// unmatched line, which would otherwise flood the terminal on large log files.
func logUnmatched(unmatched, total int) {
	if unmatched > 0 {
		log.Printf("%d of %d lines did not match any pattern in the config", unmatched, total)
	}
}

// readLines streams lines from r without loading the entire input into memory
// as a single contiguous string. Trailing "\n" and "\r\n" are stripped and
// blank lines are preserved so that line numbers stay accurate. Lines of
// arbitrary length are supported.
func readLines(r io.Reader) ([]string, error) {
	br := bufio.NewReader(r)
	lines := make([]string, 0)
	for {
		line, err := br.ReadString('\n')
		if len(line) > 0 {
			lines = append(lines, strings.TrimRight(line, "\r\n"))
		}
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			return lines, err
		}
	}
}

// ParseReader streams log lines from r and extracts token parameters from each
// line using the most appropriate pattern in the given config. Reading from a
// stream avoids holding a second full-size copy of the log in memory, so it
// scales to large files and also supports stdin.
func ParseReader(r io.Reader, config *Config) ([]Extraction, error) {
	lines, err := readLines(r)
	if err != nil {
		return nil, err
	}
	return ParseLines(lines, config)
}

// ParseStream reads lines from r one at a time and calls emit with the
// extraction for each non-blank line as soon as the line is available. Only
// single-line patterns are considered, since streaming cannot look ahead for
// multi-line patterns.
//
// When follow is true, reaching the end of the input means "wait for more data"
// rather than "done", so a growing file can be tailed; the call then returns
// only on error or an emit failure. When follow is false it returns at EOF,
// which also lets it stream naturally from a pipe (each line is emitted as soon
// as it is written, and the pipe closing ends the stream).
func ParseStream(r io.Reader, config *Config, follow bool, emit func(Extraction) error) error {
	br := bufio.NewReader(r)
	lineNum := 0
	var partial strings.Builder

	emitLine := func() error {
		line := strings.TrimRight(partial.String(), "\r\n")
		partial.Reset()
		if line != "" {
			params, pattern := parseLine([]string{line}, 0, config)
			if err := emit(Extraction{Params: params, Pattern: pattern, LineNumber: lineNum, Line: line}); err != nil {
				return err
			}
		}
		lineNum++
		return nil
	}

	for {
		chunk, err := br.ReadString('\n')
		if err == nil {
			partial.WriteString(chunk)
			if e := emitLine(); e != nil {
				return e
			}
			continue
		}
		if err == io.EOF {
			// Keep any unterminated trailing text for the next read.
			partial.WriteString(chunk)
			if follow {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			if partial.Len() > 0 {
				if e := emitLine(); e != nil {
					return e
				}
			}
			return nil
		}
		return err
	}
}

// ParseFile streams the log text from the given file path and attempts to
// extract token parameters from each line using the most appropriate pattern in
// the given config.
func ParseFile(path string, config *Config) ([]Extraction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseReader(file, config)
}

// ParseFile reads the log text from each of the given file paths, separates the
// text into lines and attempts to extract tokens parameters from each line
// using the most appropriate pattern in the given config.
func ParseFiles(paths []string, config *Config) ([]Extraction, error) {
	extraction := make([]Extraction, 0)
	for _, path := range paths {
		ex, err := ParseFile(path, config)
		if err != nil {
			return nil, fmt.Errorf("unable to parse file at path %s: %w", path, err)
		}
		extraction = append(extraction, ex...)
	}
	return extraction, nil
}

func writeConfigTest(extractions []Extraction) (int, error) {
	//write data as buffer to json encoder
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(extractions)
	if err != nil {
		return 0, err
	}
	fileDir := fmt.Sprintf("test-%d.json", time.Now().Unix())
	file, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return 0, err
	}
	n, err := file.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}
	fmt.Println(fileDir) // Display output location
	return n, nil
}

type ExtractionDebug struct {
	Line       string         `json:"line"`
	LineNumber int            `json:"lineNumber"`
	Pattern    string         `json:"pattern"`
	Params     map[string]any `json:"params"`
}

// ParseTest runs Parse and displays a random sample of extracted parameters
// along with the origin lines from the log text.
func ParseTest(logtext string, config *Config) []Extraction {
	extractions, err := Parse(logtext, config)
	if err != nil {
		panic(err)
	}
	// Random sample...
	writeConfigTest(extractions)
	return extractions
}

// ParseTest runs ParseFile and displays a random sample of extracted parameters
// along with the origin lines from the log file.
func ParseFileTest(path string, config *Config) ([]Extraction, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	extractions := ParseTest(string(body), config)
	if err != nil {
		return nil, err
	}

	return extractions, nil
}

// ParseTest runs ParseFiles and displays a random sample of extracted
// parameters along with the origin lines from the log files.
func ParseFilesTest(paths []string, config *Config) ([]Extraction, error) {
	extractions := make([]Extraction, 0)

	var parsedAny bool
	for _, path := range paths {
		r, err := ParseFileTest(path, config)
		if err != nil {
			log.Printf("unable to read file at path %s: %s\n", path, fmt.Sprint(err))
			continue
		}
		parsedAny = true
		extractions = append(extractions, r...)
	}

	if !parsedAny {
		return nil, errors.New("unable to read log file path provided")
	}

	return extractions, nil
}
