package main

import (
	"reflect"
	"testing"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want options
	}{
		{
			name: "single log path",
			args: []string{"app.log"},
			want: options{logPaths: []string{"app.log"}},
		},
		{
			name: "multiple log paths",
			args: []string{"a.log", "b.log"},
			want: options{logPaths: []string{"a.log", "b.log"}},
		},
		{
			name: "config and port long flags",
			args: []string{"--config", "cfg.json", "--port", "8080", "app.log"},
			want: options{logPaths: []string{"app.log"}, configPath: "cfg.json", port: "8080"},
		},
		{
			name: "config and port short flags",
			args: []string{"-c", "cfg.json", "-P", "8080", "app.log"},
			want: options{logPaths: []string{"app.log"}, configPath: "cfg.json", port: "8080"},
		},
		{
			name: "flags after positional (intermixed)",
			args: []string{"app.log", "--test", "--port", "9000"},
			want: options{logPaths: []string{"app.log"}, port: "9000", test: true},
		},
		{
			name: "print aliases",
			args: []string{"-d", "app.log"},
			want: options{logPaths: []string{"app.log"}, print: true},
		},
		{
			name: "help",
			args: []string{"--help"},
			want: options{help: true},
		},
		{
			name: "version",
			args: []string{"-v"},
			want: options{version: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestParseArgsErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "config flag without value", args: []string{"--config"}},
		{name: "port flag without value", args: []string{"app.log", "-P"}},
		{name: "unknown flag", args: []string{"--bogus", "app.log"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseArgs(tt.args); err == nil {
				t.Errorf("parseArgs(%v) expected error, got nil", tt.args)
			}
		})
	}
}

func TestTokensExtracted(t *testing.T) {
	none := []parse.Extraction{
		{Params: map[string]parse.Param{}},
		{Params: map[string]parse.Param{}},
	}
	if tokensExtracted(none) {
		t.Error("expected false when no params extracted")
	}

	some := []parse.Extraction{
		{Params: map[string]parse.Param{}},
		{Params: map[string]parse.Param{"ip": {Value: "1.2.3.4", Type: "ip"}}},
	}
	if !tokensExtracted(some) {
		t.Error("expected true when at least one line has params")
	}
}
