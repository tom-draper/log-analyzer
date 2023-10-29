package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tom-draper/log-analyzer/pkg/analyze"
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
		extraction, err := parse.ParseFile(profile.logpath, &profile.config)
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

		fmt.Println(extraction)
		parse.DisplayTestLines(extraction)
	}
}

func TestLogAnalyzer(t *testing.T) {
	config, err := parse.LoadConfig("./data/config.json")
	if err != nil {
		t.Error(err)
	}

	extraction, err := parse.ParseFile("./data/test.log", &config)
	if err != nil {
		t.Error(err)
	}

	parse.DisplayLines(extraction)
	analyze.Run(extraction, &config)
}
