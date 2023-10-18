package parse

import (
	"internal/display"
	"testing"
)

type Profile struct {
	config  Config
	logpath string
}

var profiles = []Profile{
	// {
	// 	config: Config{
	// 		Tokens: []string{"<timestamp>", "thread", "message"},
	// 		Patterns: []string{
	// 			"<timestamp> :: thread - message",
	// 			"<timestamp> :: message",
	// 			"message",
	// 		},
	// 	},
	// 	logpath: "../../test/data/test.log",
	// },
	{
		config: Config{
			Tokens: []string{"timestamp", "status", "message"},
			Patterns: []string{
				"[timestamp] [status] message",
			},
		},
		logpath: "../../test/data/logs/Apache.log",
	},
}

func TestParseFile(t *testing.T) {
	for _, profile := range profiles {
		lines, err := ParseFile(profile.logpath, profile.config)
		if err != nil {
			t.Error(err)
		}
		display.DisplayLines(lines)
	}
}
