package loganalyzer

type Config []Pattern

type Pattern struct {
	Pattern string
	Tokens  []Token // all unique tokens used in pattern
}

type Token struct {
	Value          string
	Timestamp      bool
	DateTimeFormat string
}

func ConfigFromJSON(path string) (Config, error) {
	return Config{}, nil
}
