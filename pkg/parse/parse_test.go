package parse

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	config := TypedConfig{
		{
			Pattern: "<timestamp> :: thread - message",
			Tokens: []Token{
				{Value: "<timestamp>", Type: "date"},
				{Value: "thread"},
				{Value: "message"},
			},
		},
		{
			Pattern: "<timestamp> :: message",
			Tokens: []Token{
				{Value: "<timestamp>", Type: "date"},
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
	fmt.Println(config)
	// Parse("Test1", config)
}
