package main

import (
	"fmt"
	"os"

	"internal/display"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No log paths provided.")
		// return
	}

	config, err := parse.LoadConfig("./config.json")
	if err != nil {
		fmt.Println("Failed to load log patterns from ./config.json")
		return
	}

	logPaths := make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		logPaths = append(logPaths, os.Args[i])
	}

	var lines []map[string]string
	if len(logPaths) > 1 {
		lines, err = parse.ParseFiles(logPaths, config)
	} else {
		// lines, err = parse.ParseFile(logPaths[0], config)
		lines, err = parse.ParseFile("./test/data/test.log", config)
	}
	if err != nil {
		panic(err)
	}

	display.DisplayLines(lines)
}
