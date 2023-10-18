package parse

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	got, err := LoadConfig("./test/data/config.json")
	if err != nil {
		t.Error("could not load config file")
	}
	expected := Config{
		Tokens: []string{"<timestamp>", "thread", "message"},
		Patterns: []string{
			"<timestamp> :: thread - message",
			"<timestamp> :: message",
			"message",
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("config got %+v did not match expected %+v", got, expected)
	}
}
