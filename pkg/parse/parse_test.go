package parse

import "testing"

func TestParse(t *testing.T) {
	config := Config{
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
	Parse("Test1", config)
}
