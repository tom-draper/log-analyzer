package main

import (
	"fmt"
	"os"

	"github.com/tom-draper/log-analyzer/pkg/analyze"
	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no log file paths provided\nprovide log file path(s) as command-line arguments\n\nexample:\n./main ./logs/postgres-main.log ./logs/postgres-1.log")
		return
	}

	logPaths, configPath, test, printLines := getCommandLineArgs()
	// Default to ./config.json
	if configPath == "" {
		configPath = "./config/config.json"
	}

	// Retrieve log line patterns from config file
	config, err := parse.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("failed to load patterns from %s\n", configPath)
		return
	}

	// If testing config against log file(s), run test
	if test {
		if len(logPaths) == 1 {
			parse.ParseFileTest(logPaths[0], &config)
		} else {
			parse.ParseFilesTest(logPaths, &config)
		}
		return
	}

	// Extract tokens from log files
	var extraction []parse.Extraction
	if len(logPaths) == 1 {
		extraction, err = parse.ParseFile(logPaths[0], &config)
		if err != nil {
			fmt.Printf("unable to parse log file: %s\n", fmt.Sprint(err))
		}
	} else {
		extraction, err = parse.ParseFiles(logPaths, &config)
		if err != nil {
			fmt.Printf("unable to parse log files: %s\n", fmt.Sprint(err))
		}
	}

	if len(extraction) == 0 {
		fmt.Println("no lines extracted\nensure log file path is correct")
		return
	} else if !tokensExtracted(extraction) {
		fmt.Println("no tokens extracted\nensure patterns in `config/config.json` are correct and all tokens are named")
		return
	}

	if printLines {
		parse.DisplayLines(extraction)
	}
	analyze.Run(extraction, &config)
}

func tokensExtracted(extraction []parse.Extraction) bool {
	for _, p := range extraction {
		if len(p.Params) > 0 {
			return true
		}
	}
	return false
}

func getCommandLineArgs() (logPaths []string, configPath string, test bool, print bool) {
	// Get log file paths from command-line arguments
	logPaths = make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-t" || arg == "--test" {
			test = true
			continue
		} else if arg == "-p" || arg == "--print" || arg == "-d" || arg == "--display" {
			print = true
			continue
		} else if arg == "-c" || arg == "--config" {
			// Skip as path will be recorded next iteration
			continue
		} else if i > 1 && (os.Args[i-1] == "-c" || os.Args[i-1] == "--config") {
			configPath = os.Args[i]
			continue
		}
		logPaths = append(logPaths, arg)
	}
	return logPaths, configPath, test, print
}
