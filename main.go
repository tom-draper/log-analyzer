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
		fmt.Println("No log paths provided.")
		return
	}

	// Retrieve log line patterns from config file
	config, err := parse.LoadConfig("./config.json")
	if err != nil {
		fmt.Println("Failed to load log patterns from ./config.json")
		return
	}

	logPaths, test := getCommandLineArgs()

	// If testing config against log file(s), run test
	if test {
		if len(logPaths) == 1 {
			parse.ParseFileTest(logPaths[0], config)
		} else {
			parse.ParseFilesTest(logPaths, config)
		}
		return
	}

	// Extract tokens from log files
	var lines []map[string]any
	if len(logPaths) == 1 {
		lines, err = parse.ParseFile(logPaths[0], config)
	} else {
		lines, err = parse.ParseFiles(logPaths, config)
	}
	if err != nil {
		panic(err)
	}

	display.DisplayLines(lines)
	server.Start(lines)
}

func getCommandLineArgs() ([]string, bool) {
	// Get log file paths from command-line arguments
	test := false
	logPaths := make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--test" {
			test = true
			continue
		}
		logPaths = append(logPaths, os.Args[i])
	}
	return logPaths, test
}
