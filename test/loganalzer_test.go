package test

import (
	"internal/display"
	"internal/server"
	"testing"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func TestLogAnalyzer(t *testing.T) {
	config, err := parse.LoadConfig("./data/config.json")
	if err != nil {
		panic(err)
	}

	lines, err := parse.ParseFile("./data/test.log", config)
	display.DisplayLines(lines)
	server.Start(lines)
}
