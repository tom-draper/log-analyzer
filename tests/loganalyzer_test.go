package test

import (
	"internal/display"
	"internal/server"
	"os"
	"strings"
	"testing"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

type Profile struct {
	config  parse.Config
	logpath string
}

var profiles = []Profile{
	{
		config: parse.Config{
			Tokens: []string{"timestamp", "type", "unknown", "message"},
			Patterns: []string{
				"timestamp, type unknown message",
			},
		},
		logpath: "./data/logs/loghub/Windows.log",
	},
	// {
	// 	config: parse.Config{
	// 		Tokens: []string{"timestamp", "status", "func", "message"},
	// 		Patterns: []string{
	// 			"[timestamp] [status] message",
	// 			"[timestamp] [status] func() message",
	// 		},
	// 	},
	// 	logpath: "./data/logs/loghub/Apache.log",
	// },
}

func splitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}

func TestParse(t *testing.T) {
	for _, profile := range profiles {
		params, err := parse.ParseFile(profile.logpath, profile.config)
		if err != nil {
			t.Error(err)
		}

		body, err := os.ReadFile(profile.logpath)
		if err != nil {
			t.Error(err)
		}

		lines := splitLines(string(body))
		indicies := make([]int, 0)
		for i := range lines {
			indicies = append(indicies, i)
		}

		display.DisplayTestLines(lines, params, indicies)
	}
}

func TestLogAnalyzer(t *testing.T) {
	config, err := parse.LoadConfig("./data/config.json")
	if err != nil {
		t.Error(err)
	}

	params, err := parse.ParseFile("./data/test.log", config)
	if err != nil {
		t.Error(err)
	}

	display.DisplayLines(params)
	server.Start(params)
}
