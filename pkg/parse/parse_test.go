package parse

import (
	"log"
	"testing"
)

type Profile struct {
	config  Config
	logpath string
}

var profiles = []Profile{
	{
		config:  loadConfig("Apache.json"),
		logpath: logPath("Apache.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"starting...",
				"timestamp :: threadthread_number - starting",
				"timestamp :: threadthread_number - exiting",
				"timestamp :: complete",
			},
			Tokens: []string{"timestamp", "thread_number"},
		},
		logpath: logPath("test.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[timestamp] ip:dbname LOG: duration: elapsed ms statement: query",
				"[timestamp] ip:dbname LOG: statement: query",
			},
			Tokens: []string{"timestamp", "ip", "dbname", "elapsed", "query"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[INFO ] [taken_msms] result found: page=url status=status_code size=bytes",
				"[INFO ] [taken_ss] result found: page=url status=status_code size=bytes",
			},
			Tokens: []string{"taken_ms", "taken_s", "url", "status_code", "bytes"},
			Conversions: map[string]Conversion{
				"taken_ms": {
					Multiplier: 1000,
					Token:      "taken_s",
				},
			},
		},
		logpath: logPath("test3.log"),
	},
}

func loadConfig(file string) Config {
	config, err := LoadConfig(configPath(file))
	if err != nil {
		panic(err)
	}
	return config
}

func configPath(file string) string {
	return "../../tests/data/configs/" + file
}

func logPath(file string) string {
	return "../../tests/data/logs/" + file
}

func TestParseFile(t *testing.T) {
	for i, profile := range profiles {
		log.Printf("testing profile %d...", i)
		extraction, err := ParseFile(profile.logpath, &profile.config)
		if err != nil {
			t.Error(err)
		}
		if len(extraction) == 0 {
			t.Errorf("profile %d: no lines extracted", i)
		}

		for _, e := range extraction {
			if e.Pattern == "" {
				t.Errorf("profile %d: line %d failed: %s", i, e.LineNumber, e.Line)
			} else if len(e.Params) == 0 {
				t.Logf("profile %d: no params extracted from line %d: %s", i, e.LineNumber, e.Line)
			}
		}
	}
}
