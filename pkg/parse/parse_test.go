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
				"starting...",
				"time_timestamp :: threadint_thread_number - starting",
				"time_timestamp :: threadint_thread_number - exiting",
				"time_timestamp :: complete",
			},
			Tokens: []string{"time_timestamp", "int_thread_number"},
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
				"[time_timestamp] ip_address:str_dbname LOG: duration: float_elapsed ms statement: str_query",
				"[time_timestamp] ip_address:str_dbname LOG: statement: str_query",
			},
			Tokens: []string{"time_timestamp", "ip_address", "str_dbname", "float_elapsed", "str_query"},
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
	{
		config: Config{
			Patterns: []string{
				"[INFO ] [float_taken_msms] result found: page=str_url status=int_status_code size=int_bytes",
				"[INFO ] [float_taken_ss] result found: page=str_url status=int_status_code size=int_bytes",
			},
			Tokens: []string{"float_taken_ms", "float_taken_s", "str_url", "str_status_code", "str_bytes"},
		},
		logpath: logPath("test3.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[timestamp] rest",
			},
			Tokens: []string{"timestamp", "rest"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[time_timestamp] str_rest",
			},
			Tokens: []string{"time_timestamp", "str_rest"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[timestamp] *",
			},
			Tokens: []string{"timestamp"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"*",
			},
			Tokens: []string{},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"all",
			},
			Tokens: []string{"all"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"str_all",
			},
			Tokens: []string{"str_all"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"rest query",
			},
			Tokens: []string{"query"},
		},
		logpath: logPath("test2.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[time] ip rest",
			},
			Tokens: []string{"time", "ip", "rest"},
		},
		logpath: logPath("demo.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"*",
			},
			Tokens: []string{},
		},
		logpath: logPath("demo.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"***",
			},
			Tokens: []string{},
		},
		logpath: logPath("demo.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[*] *",
			},
			Tokens: []string{},
		},
		logpath: logPath("demo.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[*] * *",
			},
			Tokens: []string{},
		},
		logpath: logPath("demo.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"timestamp x y z source: message",
			},
			Tokens: []string{"timestamp", "x", "y", "z", "source", "message"},
		},
		logpath: logPath("loghub/Android_v1.log"),
	},
	{
		config: Config{
			Patterns: []string{
				"[timestamp] [type] jk2_init() Found child child_number in scoreboard slot slot_number",
				"[timestamp] [type] [client ip] message",
				"[timestamp] [type] mod_jk child workerEnv in error state state_number",
				"[timestamp] [type] workerEnv.init() ok /etc/httpd/conf/workers2.properties",
				"[timestamp] [type] other",
			},
			Tokens: []string{"timestamp", "type", "child_number", "slot_number", "ip", "message", "state_number", "other"},
		},
		logpath: logPath("loghub/Apache.log"),
	},
	{
		config:  loadConfig("multiline.json"),
		logpath: logPath("multiline.log"),
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
