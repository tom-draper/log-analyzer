package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tom-draper/log-analyzer/pkg/analyze"
	"github.com/tom-draper/log-analyzer/pkg/parse"
	"github.com/tom-draper/log-analyzer/pkg/similarity"
)

//go:embed dashboard/dist/index.html
var indexHTML []byte

const version = "1.0.0"

const defaultFollowMaxRecords = 10_000

const usage = `Log Analyzer - turn unstructured log files into a dashboard.

Usage:
  log-analyzer [options] <log-file> [log-file...]
  log-analyzer suggest <log-file>

Commands:
  suggest               Group similar log lines to help you write config patterns

Options:
  -c, --config <path>   Path to the config file (default: ./config/config.json)
  -P, --port <port>     Port for the dashboard (default: 3000)
  -H, --host <host>     Interface for the dashboard (default: 127.0.0.1)
  -o, --output <path>   Write extracted data to a file instead of the dashboard
  -f, --format <fmt>    Export format: json, ndjson or csv (implies headless)
      --follow          Stream records as new lines are appended (ndjson or print)
      --follow-max-records <n>
                        Keep the latest n records in a live dashboard (default: 10000; 0 keeps all)
  -t, --test            Test config against the log file(s); write extraction to JSON
  -p, --print           Print extracted tokens for each line to the terminal
  -h, --help            Show this help message
  -v, --version         Show version information

Use - as the log path to read from stdin.

Example:
  log-analyzer ./logs/app.log
  log-analyzer -c ./config/config.json -P 8080 ./logs/app.log
  log-analyzer --format csv ./logs/app.log > app.csv
  log-analyzer -o data.ndjson ./logs/app.log
  log-analyzer --follow --format ndjson ./logs/app.log
  tail -f ./logs/app.log | log-analyzer --follow -
`

type options struct {
	logPaths            []string
	configPath          string
	port                string
	host                string
	output              string
	format              string
	follow              bool
	followMaxRecords    int
	followMaxRecordsSet bool
	test                bool
	print               bool
	help                bool
	version             bool
}

func main() {
	// Subcommands.
	if len(os.Args) > 1 && os.Args[1] == "suggest" {
		if err := runSuggest(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}

	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		fmt.Fprintln(os.Stderr, "\nrun `log-analyzer --help` for usage")
		os.Exit(2)
	}

	if opts.help {
		fmt.Print(usage)
		return
	}
	if opts.version {
		fmt.Printf("log-analyzer %s\n", version)
		return
	}

	if len(opts.logPaths) == 0 {
		fmt.Fprintln(os.Stderr, "error: no log file path provided")
		fmt.Fprintln(os.Stderr, "\nrun `log-analyzer --help` for usage")
		os.Exit(2)
	}

	// Defaults
	if opts.configPath == "" {
		opts.configPath = "./config/config.json"
	}
	if opts.port == "" {
		opts.port = "3000"
	}
	if opts.host == "" {
		opts.host = "127.0.0.1"
	}
	if !opts.followMaxRecordsSet {
		opts.followMaxRecords = defaultFollowMaxRecords
	}

	// Retrieve log line patterns from config file
	config, err := parse.LoadConfig(opts.configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config from %s: %s\n", opts.configPath, err)
		os.Exit(1)
	}

	// If testing config against log file(s), run test
	if opts.test {
		var extraction []parse.Extraction
		if len(opts.logPaths) == 1 {
			extraction, err = parse.ParseFileTest(opts.logPaths[0], &config)
		} else {
			extraction, err = parse.ParseFilesTest(opts.logPaths, &config)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to test log file(s): %s\n", err)
			os.Exit(1)
		}
		if _, err := parse.WriteConfigTest(extraction); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write config test output: %s\n", err)
			os.Exit(1)
		}
		return
	}

	// Follow mode tails the input as new lines arrive. With an export format,
	// output path or print flag it streams records to the terminal/file;
	// otherwise it serves a live-updating dashboard.
	if opts.follow {
		if opts.format != "" || opts.output != "" || opts.print {
			if err := runFollow(opts, &config); err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
			return
		}
		if len(opts.logPaths) != 1 {
			fmt.Fprintln(os.Stderr, "error: --follow requires exactly one input (a file path, or - for stdin)")
			os.Exit(2)
		}
		analyze.RunFollow(opts.logPaths[0], &config, opts.host, opts.port, opts.followMaxRecords, indexHTML)
		return
	}

	// Extract tokens from the log file(s), or stdin when the path is "-".
	extraction, err := parseLogs(opts.logPaths, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to parse log file(s): %s\n", err)
		os.Exit(1)
	}

	if len(extraction) == 0 {
		fmt.Fprintln(os.Stderr, "error: no lines extracted; ensure the log file path is correct")
		os.Exit(1)
	} else if !tokensExtracted(extraction) {
		fmt.Fprintf(os.Stderr, "error: no tokens extracted; ensure the patterns in %s are correct and all tokens are named\n", opts.configPath)
		os.Exit(1)
	}

	if opts.print {
		parse.DisplayLines(extraction)
	}

	// Headless export mode: write structured data instead of serving the dashboard.
	if opts.format != "" || opts.output != "" {
		if err := runExport(opts, extraction); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		return
	}

	analyze.Run(extraction, &config, opts.host, opts.port, indexHTML)
}

// runExport writes the extracted data to the configured output (a file, or
// stdout when no output path is given) in the requested or inferred format.
func runExport(opts options, extraction []parse.Extraction) error {
	format := opts.format
	if format == "" {
		format = parse.FormatFromPath(opts.output)
	}
	if format == "" {
		return fmt.Errorf("could not determine export format; pass --format json|ndjson|csv")
	}

	out := io.Writer(os.Stdout)
	if opts.output != "" {
		file, err := os.Create(opts.output)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", opts.output, err)
		}
		defer file.Close()
		out = file
	}

	if err := parse.Export(out, format, extraction); err != nil {
		return err
	}
	if opts.output != "" {
		fmt.Fprintf(os.Stderr, "wrote %d records to %s\n", len(extraction), opts.output)
	}
	return nil
}

// maxSuggestSample caps how many lines are compared when suggesting patterns;
// similarity grouping is O(n^2) in lines, so a representative sample keeps it fast.
const maxSuggestSample = 100

// runSuggest groups similar lines from a log file (or stdin) so the user can see
// the distinct log formats present and write a config pattern for each.
func runSuggest(args []string) error {
	var path string
	for _, a := range args {
		if a == "-" || (len(a) > 0 && a[0] != '-') {
			path = a
			break
		}
	}
	if path == "" {
		return fmt.Errorf("usage: log-analyzer suggest <log-file>")
	}

	lines, err := readNonEmptyLines(path)
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		return fmt.Errorf("no lines read from %s", path)
	}

	sample := sampleLines(lines, maxSuggestSample)
	groups := similarity.FindGroupsFromLines(sample)
	printSuggestions(os.Stdout, groups, len(lines), len(sample))
	return nil
}

func readNonEmptyLines(path string) ([]string, error) {
	var r io.Reader = os.Stdin
	if path != "-" {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		r = file
	}

	br := bufio.NewReader(r)
	var lines []string
	for {
		line, err := br.ReadString('\n')
		if trimmed := strings.TrimRight(line, "\r\n"); trimmed != "" {
			lines = append(lines, trimmed)
		}
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			return lines, err
		}
	}
}

// sampleLines returns at most limit lines, evenly spaced across the input so the
// sample reflects formats found throughout the file rather than just the start.
func sampleLines(lines []string, limit int) []string {
	if len(lines) <= limit {
		return lines
	}
	sample := make([]string, 0, limit)
	step := float64(len(lines)) / float64(limit)
	for i := 0; i < limit; i++ {
		sample = append(sample, lines[int(float64(i)*step)])
	}
	return sample
}

func printSuggestions(w io.Writer, groups []similarity.Group, total, sampled int) {
	if sampled < total {
		fmt.Fprintf(w, "Sampled %d of %d lines.\n", sampled, total)
	}
	fmt.Fprintf(w, "Found %d distinct log format(s):\n\n", len(groups))

	counter := 0
	var patterns, tokens []string
	for i, g := range groups {
		fmt.Fprintf(w, "Group %d (%d line%s)\n", i+1, len(g.Members), plural(len(g.Members)))
		fmt.Fprintf(w, "  example: %s\n", representative(g.Members))
		// A pattern can only be derived when there are at least two lines to
		// compare (so the varying fields can be identified).
		if len(g.Members) >= 2 {
			if pattern, toks, ok := generatePattern(g.Members, &counter); ok {
				fmt.Fprintf(w, "  pattern: %s\n", pattern)
				patterns = append(patterns, pattern)
				tokens = append(tokens, toks...)
			}
		}
		fmt.Fprintln(w)
	}

	if len(patterns) > 0 {
		fmt.Fprint(w, "Suggested config (rename the field tokens to something meaningful):\n\n")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(parse.Config{Patterns: patterns, Tokens: tokens})
	} else {
		fmt.Fprint(w, "Add a pattern to your config for each group, replacing the variable parts\n"+
			"(timestamps, IPs, ids, ...) with tokens.\n")
	}
}

// segment is one unit of a log line: a bracketed group "[...]", a quoted group
// "\"...\"", or a bare whitespace-delimited run. Treating bracketed and quoted
// groups as single units keeps variable-length fields (such as quoted requests
// or user agents) aligned across lines.
type segment struct {
	text string
	kind byte // '[', '"', or 'b' for bare
}

func tokenizeLine(line string) []segment {
	var segs []segment
	for i := 0; i < len(line); {
		if line[i] == ' ' || line[i] == '\t' {
			i++
			continue
		}
		switch line[i] {
		case '[':
			if j := strings.IndexByte(line[i:], ']'); j != -1 {
				segs = append(segs, segment{line[i : i+j+1], '['})
				i += j + 1
				continue
			}
		case '"':
			if j := strings.IndexByte(line[i+1:], '"'); j != -1 {
				segs = append(segs, segment{line[i : i+1+j+1], '"'})
				i += 1 + j + 1
				continue
			}
		}
		// Bare run until the next whitespace.
		end := i
		for end < len(line) && line[end] != ' ' && line[end] != '\t' {
			end++
		}
		segs = append(segs, segment{line[i:end], 'b'})
		i = end
	}
	return segs
}

// generatePattern derives a config pattern from a group of similar lines by
// aligning them field-by-field: fields that are identical across every line
// become literals, and fields that vary become tokens (field1, field2, ...,
// numbered via counter so names stay unique across groups). Bracket and quote
// delimiters are kept as literals around their content, and their content is
// recursively sub-split into tokens too (so a quoted request becomes
// "method endpoint version" rather than a single token). It reports ok=false
// when the lines do not share the same field structure.
func generatePattern(lines []string, counter *int) (string, []string, bool) {
	tokenized := make([][]segment, len(lines))
	for i, line := range lines {
		tokenized[i] = tokenizeLine(line)
	}

	fields := len(tokenized[0])
	for _, segs := range tokenized {
		if len(segs) != fields {
			return "", nil, false
		}
	}

	var parts, tokens []string
	for pos := 0; pos < fields; pos++ {
		texts := make([]string, len(tokenized))
		for i, segs := range tokenized {
			texts[i] = segs[pos].text
		}
		if allEqual(texts) {
			parts = append(parts, texts[0])
			continue
		}
		part, toks := patternForSegment(texts, tokenized[0][pos].kind, counter)
		parts = append(parts, part)
		tokens = append(tokens, toks...)
	}
	return strings.Join(parts, " "), tokens, true
}

// patternForSegment produces the pattern fragment for a single varying field.
// Bracketed and quoted fields keep their delimiters and have their inner
// whitespace-separated content aligned into sub-tokens when every line splits
// into the same number of parts; otherwise the whole field becomes one token.
func patternForSegment(texts []string, kind byte, counter *int) (string, []string) {
	if kind == '[' || kind == '"' {
		open, close := "[", "]"
		if kind == '"' {
			open, close = "\"", "\""
		}
		rows := make([][]string, len(texts))
		for i, t := range texts {
			rows[i] = strings.Fields(t[1 : len(t)-1]) // strip delimiters
		}
		if len(rows[0]) >= 1 && sameRowLength(rows) {
			inner, toks := alignFields(rows, counter)
			return open + inner + close, toks
		}
		name := nextToken(counter)
		return open + name + close, []string{name}
	}
	name := nextToken(counter)
	return name, []string{name}
}

// alignFields aligns equal-length rows of tokens, emitting a literal where every
// row agrees and a token where they differ.
func alignFields(rows [][]string, counter *int) (string, []string) {
	var parts, tokens []string
	for pos := 0; pos < len(rows[0]); pos++ {
		column := make([]string, len(rows))
		for i, r := range rows {
			column[i] = r[pos]
		}
		if allEqual(column) {
			parts = append(parts, column[0])
			continue
		}
		name := nextToken(counter)
		parts = append(parts, name)
		tokens = append(tokens, name)
	}
	return strings.Join(parts, " "), tokens
}

func nextToken(counter *int) string {
	*counter++
	return fmt.Sprintf("field%d", *counter)
}

func allEqual(values []string) bool {
	for _, v := range values {
		if v != values[0] {
			return false
		}
	}
	return true
}

func sameRowLength(rows [][]string) bool {
	for _, r := range rows {
		if len(r) != len(rows[0]) {
			return false
		}
	}
	return true
}

// representative returns the shortest member, which is usually closest to the
// group's underlying template.
func representative(members []string) string {
	rep := members[0]
	for _, m := range members[1:] {
		if len(m) < len(rep) {
			rep = m
		}
	}
	return rep
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func tokensExtracted(extraction []parse.Extraction) bool {
	for _, p := range extraction {
		if len(p.Params) > 0 {
			return true
		}
	}
	return false
}

// parseLogs parses the given log paths, reading from stdin when the single path
// is "-".
func parseLogs(paths []string, config *parse.Config) ([]parse.Extraction, error) {
	if len(paths) == 1 {
		if paths[0] == "-" {
			return parse.ParseReader(os.Stdin, config)
		}
		return parse.ParseFile(paths[0], config)
	}
	return parse.ParseFiles(paths, config)
}

// runFollow tails a single input (a file, or stdin when the path is "-") and
// emits each extracted record as new lines arrive.
func runFollow(opts options, config *parse.Config) error {
	if len(opts.logPaths) != 1 {
		return fmt.Errorf("--follow requires exactly one input (a file path, or - for stdin)")
	}

	emit, closeFn, err := streamEmitter(opts)
	if err != nil {
		return err
	}
	defer closeFn()

	path := opts.logPaths[0]
	if path == "-" {
		// A pipe delivers lines as they are written and ends at EOF, so
		// following is implicit and we stop when the pipe closes.
		return parse.ParseStream(os.Stdin, config, false, emit)
	}

	fmt.Fprintf(os.Stderr, "following %s (Ctrl-C to stop)\n", path)
	return parse.FollowFile(path, config, emit)
}

// streamEmitter builds a per-record emitter for follow mode based on the chosen
// format and output. Streaming supports ndjson and human-readable print output.
func streamEmitter(opts options) (func(parse.Extraction) error, func(), error) {
	format := opts.format
	if format == "" && opts.output != "" {
		format = parse.FormatFromPath(opts.output)
	}

	out := io.Writer(os.Stdout)
	closeFn := func() {}
	if opts.output != "" {
		file, err := os.Create(opts.output)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create output file %s: %w", opts.output, err)
		}
		out = file
		closeFn = func() { file.Close() }
	}

	switch format {
	case "", "print":
		return func(e parse.Extraction) error {
			parse.FprintExtraction(out, e)
			return nil
		}, closeFn, nil
	case "ndjson":
		encoder := json.NewEncoder(out)
		return func(e parse.Extraction) error {
			return encoder.Encode(e)
		}, closeFn, nil
	default:
		closeFn()
		return nil, nil, fmt.Errorf("--follow supports --format ndjson or plain print output, not %q", format)
	}
}

// parseArgs parses command-line arguments, allowing flags and log file paths to
// be intermixed in any order.
func parseArgs(args []string) (options, error) {
	var opts options
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.help = true
		case "-v", "--version":
			opts.version = true
		case "--follow":
			opts.follow = true
		case "--follow-max-records":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s flag requires a non-negative integer", arg)
			}
			i++
			value, err := strconv.Atoi(args[i])
			if err != nil || value < 0 {
				return opts, fmt.Errorf("%s flag requires a non-negative integer", arg)
			}
			opts.followMaxRecords = value
			opts.followMaxRecordsSet = true
		case "-t", "--test":
			opts.test = true
		case "-p", "--print", "-d", "--display":
			opts.print = true
		case "-c", "--config":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s flag requires a path argument", arg)
			}
			i++
			opts.configPath = args[i]
		case "-P", "--port":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s flag requires a port argument", arg)
			}
			i++
			opts.port = args[i]
		case "-H", "--host":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s flag requires a host argument", arg)
			}
			i++
			opts.host = args[i]
		case "-o", "--output":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s flag requires a path argument", arg)
			}
			i++
			opts.output = args[i]
		case "-f", "--format":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s flag requires a format argument", arg)
			}
			i++
			opts.format = args[i]
		default:
			// "-" is the conventional stdin marker, not a flag.
			if arg != "-" && len(arg) > 0 && arg[0] == '-' {
				return opts, fmt.Errorf("unknown flag: %s", arg)
			}
			opts.logPaths = append(opts.logPaths, arg)
		}
	}
	return opts, nil
}
