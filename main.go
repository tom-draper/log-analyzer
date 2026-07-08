package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/tom-draper/log-analyzer/pkg/analyze"
	"github.com/tom-draper/log-analyzer/pkg/parse"
)

//go:embed dashboard/dist/index.html
var indexHTML []byte

const version = "1.0.0"

const usage = `Log Analyzer - turn unstructured log files into a dashboard.

Usage:
  log-analyzer [options] <log-file> [log-file...]

Options:
  -c, --config <path>   Path to the config file (default: ./config/config.json)
  -P, --port <port>     Port for the dashboard (default: 3000)
  -t, --test            Test config against the log file(s); write extraction to JSON
  -p, --print           Print extracted tokens for each line to the terminal
  -h, --help            Show this help message
  -v, --version         Show version information

Example:
  log-analyzer ./logs/app.log
  log-analyzer -c ./config/config.json -P 8080 ./logs/app.log
`

type options struct {
	logPaths   []string
	configPath string
	port       string
	test       bool
	print      bool
	help       bool
	version    bool
}

func main() {
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

	// Retrieve log line patterns from config file
	config, err := parse.LoadConfig(opts.configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load config from %s: %s\n", opts.configPath, err)
		os.Exit(1)
	}

	// If testing config against log file(s), run test
	if opts.test {
		if len(opts.logPaths) == 1 {
			parse.ParseFileTest(opts.logPaths[0], &config)
		} else {
			parse.ParseFilesTest(opts.logPaths, &config)
		}
		return
	}

	// Extract tokens from log files
	var extraction []parse.Extraction
	if len(opts.logPaths) == 1 {
		extraction, err = parse.ParseFile(opts.logPaths[0], &config)
	} else {
		extraction, err = parse.ParseFiles(opts.logPaths, &config)
	}
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
	analyze.Run(extraction, &config, opts.port, indexHTML)
}

func tokensExtracted(extraction []parse.Extraction) bool {
	for _, p := range extraction {
		if len(p.Params) > 0 {
			return true
		}
	}
	return false
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
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return opts, fmt.Errorf("unknown flag: %s", arg)
			}
			opts.logPaths = append(opts.logPaths, arg)
		}
	}
	return opts, nil
}
