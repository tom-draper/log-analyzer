package parse

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfigRequiresSingleJSONDocument(t *testing.T) {
	for _, tt := range []struct {
		name    string
		suffix  string
		wantErr bool
	}{
		{name: "end of file"},
		{name: "trailing whitespace", suffix: " \n\t\r"},
		{name: "trailing garbage", suffix: " not JSON", wantErr: true},
		{name: "second object", suffix: " {}", wantErr: true},
		{name: "second scalar", suffix: " null", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(`{"tokens":["value"],"patterns":["value"]}`+tt.suffix), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := LoadConfig(path)
			if tt.wantErr {
				if err == nil {
					t.Fatal("LoadConfig accepted trailing data")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := Config{Tokens: []string{"value"}, Patterns: []string{"value"}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("config = %+v, want %+v", got, want)
			}
		})
	}
}

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

func TestValidateConfig(t *testing.T) {
	valid := Config{Tokens: []string{"timestamp", "status"}, Patterns: []string{"timestamp status"}}
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{name: "no patterns", config: Config{}, want: "at least one pattern"},
		{name: "duplicate token", config: Config{Tokens: []string{"value", "value"}, Patterns: valid.Patterns}, want: "declared more than once"},
		{name: "reserved wildcard", config: Config{Tokens: []string{"*"}, Patterns: valid.Patterns}, want: "reserved"},
		{name: "unknown dependency", config: Config{Tokens: valid.Tokens, Patterns: valid.Patterns, Dependencies: map[string][]string{"missing": {"status"}}}, want: "not declared"},
		{name: "self dependency", config: Config{Tokens: valid.Tokens, Patterns: valid.Patterns, Dependencies: map[string][]string{"status": {"status"}}}, want: "cannot depend"},
		{name: "unknown conversion target", config: Config{Tokens: valid.Tokens, Patterns: valid.Patterns, Conversions: map[string]Conversion{"status": {Token: "missing", Multiplier: 1}}}, want: "target token"},
		{name: "zero conversion multiplier", config: Config{Tokens: valid.Tokens, Patterns: valid.Patterns, Conversions: map[string]Conversion{"status": {Token: "timestamp", Multiplier: 0}}}, want: "finite non-zero"},
	}

	if err := ValidateConfig(valid); err != nil {
		t.Fatalf("valid config returned error: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("ValidateConfig() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}
