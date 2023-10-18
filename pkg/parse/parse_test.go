package parse

import (
	"testing"
)

type Profile struct{
	config Config
	logpath string
}

var profiles = []Profile{
	{
		config: Config{
		Tokens: []string{"<timestamp>", "thread", "message"},
		Patterns: []string{
			"<timestamp> :: thread - message",
			"<timestamp> :: message",
			"message",
		},
		},
		logpath: "../../test/data/test.log",
	},
}

func TestParseFile(t *testing.T) {
	for _, profile := range profiles {
		ParseFile(profile.logpath, profile.config)

	}
}
