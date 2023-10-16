package test

import (
	"testing"

	loganalyzer "github.com/tom-draper/log-analyzer"
)

func TestLogAnalyzer(t *testing.T) {
	config, err := loganalyzer.LoadConfig("./data/config.json")
	if err != nil {
		panic(err)
	}
	// config := loganalyzer.Config{
	// 	{
	// 		Pattern: "<timestamp> :: thread - message",
	// 		Tokens: []loganalyzer.Token{
	// 			{Value: "<timestamp>", Timestamp: true},
	// 			{Value: "thread"},
	// 			{Value: "message"},
	// 		},
	// 	},
	// 	{
	// 		Pattern: "<timestamp> :: message",
	// 		Tokens: []loganalyzer.Token{
	// 			{Value: "<timestamp>", Timestamp: true},
	// 			{Value: "message"},
	// 		},
	// 	},
	// 	{
	// 		Pattern: "message",
	// 		Tokens: []loganalyzer.Token{
	// 			{Value: "message"},
	// 		},
	// 	},
	// }
	loganalyzer.ParseFile("./data/test.log", config)
}
