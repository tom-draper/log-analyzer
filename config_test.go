package loganalyzer

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
		{
			Pattern: "<timestamp> :: thread - message",
			Tokens: []Token{
				{Value: "<timestamp>", Timestamp: true},
				{Value: "thread"},
				{Value: "message"},
			},
		},
		{
			Pattern: "<timestamp> :: message",
			Tokens: []Token{
				{Value: "<timestamp>", Timestamp: true},
				{Value: "message"},
			},
		},
		{
			Pattern: "message",
			Tokens: []Token{
				{Value: "message"},
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("config got %+v did not match expected %+v", got, expected)
	}
}
