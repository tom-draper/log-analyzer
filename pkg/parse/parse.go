	}
	if len(best) == 0 {
		log.Printf("no pattern matched line: %s\n", line)
	}
	return best
}

func splitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}

// Parse separates the log text into lines and attempts to extract tokens
// parameters from each line using the most appropriate pattern in the given config.
func Parse(logtext string, config Config) ([]map[string]any, error) {
	params := make([]map[string]any, 0)
	for _, line := range splitLines(logtext) {
		p := parseLine(line, config)
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

func sampleLines(lines []string, params []map[string]any, n int) ([]string, []map[string]any) {
	sampledLines := make([]string, 0)
	sampledParams := make([]map[string]any, 0)
	selectedIndices := make(map[int]struct{}, 0)
	for i := 0; i < min(n, len(lines)); i++ {
		idx := randomIndex(len(lines), selectedIndices)
		sampledLines = append(sampledLines, lines[idx])
		sampledParams = append(sampledParams, params[idx])
	}
	return sampledLines, sampledParams
}

// randomIndex selects and returns a random index from a slice of a given size.
// Indicies contained within existingIndicies memory will not be selected.
func randomIndex(size int, existingIndicies map[int]struct{}) int {
	for {
		randomIndex := rand.Intn(size)
		_, exists := existingIndicies[randomIndex]
		if !exists {
			existingIndicies[randomIndex] = struct{}{}
			return randomIndex
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// randomIndicies build and returns a list of n random index integers available
// in the slice of params.
func randomIndices(params []map[string]any, n int) []int {
	indicies := make([]int, 0)
	selectedIndices := make(map[int]struct{}, 0)
	for i := 0; i < min(n, len(params)); i++ {
		idx := randomIndex(len(params), selectedIndices)
		indicies = append(indicies, idx)
	}
	return indicies
}

// Randomly samples extracted params and displays them along with the origin log text line for user evaluation.
func displayConfigTest(logtext string, params []map[string]any) {
	indicies := randomIndices(params, 5)
	sampledLines := make([]string, 0)
	sampledParams := make([]map[string]any, 0)
	lines := splitLines(logtext)
	for _, idx := range indicies {
		sampledLines = append(sampledLines, lines[idx])
		sampledParams = append(sampledParams, params[idx])
	}

	display.DisplayTestLines(sampledLines, sampledParams, indicies)
}

// ParseTest runs Parse and displays a random sample of extracted parameters
// along with the origin lines from the log text.
func ParseTest(logtext string, config Config) {
	params, err := Parse(logtext, config)
	if err != nil {
		panic(err)
	}
	displayConfigTest(logtext, params)
}

// ParseTest runs ParseFile and displays a random sample of extracted parameters
// along with the origin lines from the log file.
func ParseFileTest(path string, config Config) {
	body, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	params, err := Parse(string(body), config)
	if err != nil {
		panic(err)
	}

	displayConfigTest(string(body), params)
}

// ParseTest runs ParseFiles and displays a random sample of extracted
// parameters along with the origin lines from the log files.
func ParseFilesTest(paths []string, config Config) {
	var logtext string                  // Build string holding all log text content
	params := make([]map[string]any, 0) // Extracted params for each log text line
	for _, path := range paths {
		fileParams, err := ParseFile(path, config)
		if err != nil {
			log.Printf("unable to parse file at path %s: %s", path, fmt.Sprint(err))
			continue
		}
		body, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		logtext += "\n" + string(body)
		params = append(params, fileParams...)
	}

	displayConfigTest(logtext, params)
}
