package main

import (
	"fmt"
	"os"

	"internal/display"
	"internal/server"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no log file paths provided\nprovide log file path(s) as command-line arguments\n\nexample:\n./main ./logs/postgres-main.log ./logs/postgres-1.log")
		return
	}

	logPaths, configPath, test := getCommandLineArgs()
	if configPath == "" {
		configPath = "./config.json"
	}

	// Retrieve log line patterns from config file
	config, err := parse.LoadConfig(configPath)
	if err != nil {
		fmt.Println("failed to load log patterns from ./config.json")
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
	var lines []map[string]any
	if len(logPaths) == 1 {
		lines, err = parse.ParseFile(logPaths[0], &config)
		if err != nil {
			fmt.Printf("unable to parse log file: %s\n", fmt.Sprint(err))
		}
	} else {
		lines, err = parse.ParseFiles(logPaths, &config)
		if err != nil {
			fmt.Printf("unable to parse log files: %s\n", fmt.Sprint(err))
		}
	}

	if len(lines) == 0 {
		fmt.Println("no lines extracted\nensure log file path is correct")
		return
	} else if !tokensExtracted(lines) {
		fmt.Println("no tokens extracted\nensure patterns in `config.json` are correct and all tokens are named")
		return
	}

	display.DisplayLines(lines)
	server.Start(lines)
}

func tokensExtracted(params []map[string]any) bool {
	for _, p := range params {
		if len(p) > 0 {
			return true
		}
	}
	return false
}

func getCommandLineArgs() ([]string, string, bool) {
	// Get log file paths from command-line arguments
	test := false
	var configPath string
	logPaths := make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--test" {
			test = true
			continue
		} else if i > 1 && (os.Args[i-1] == "-c" || os.Args[i-1] == "--config") {
			configPath = os.Args[i]
			continue
		}
		logPaths = append(logPaths, os.Args[i])
	}
	return logPaths, configPath, test
}
