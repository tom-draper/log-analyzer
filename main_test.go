package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/tom-draper/log-analyzer/pkg/parse"
	"github.com/tom-draper/log-analyzer/pkg/similarity"
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
			name: "stdin dash is a path, not a flag",
			args: []string{"-"},
			want: options{logPaths: []string{"-"}},
		},
		{
			name: "follow with stdin and format",
			args: []string{"--follow", "-", "-f", "ndjson"},
			want: options{logPaths: []string{"-"}, follow: true, format: "ndjson"},
		},
		{
			name: "output and format flags",
			args: []string{"-o", "out.csv", "app.log"},
			want: options{logPaths: []string{"app.log"}, output: "out.csv"},
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

func TestSampleLines(t *testing.T) {
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = string(rune('a' + i%26))
	}

	// Under the limit: returns everything unchanged.
	if got := sampleLines(lines[:5], 10); len(got) != 5 {
		t.Errorf("expected 5 lines, got %d", len(got))
	}

	// Over the limit: returns exactly limit lines, starting at the first.
	got := sampleLines(lines, 10)
	if len(got) != 10 {
		t.Fatalf("expected 10 sampled lines, got %d", len(got))
	}
	if got[0] != lines[0] {
		t.Errorf("expected sample to start at first line, got %q", got[0])
	}
}

func TestRepresentative(t *testing.T) {
	members := []string{"a longer line", "short", "medium one"}
	if got := representative(members); got != "short" {
		t.Errorf("expected shortest member, got %q", got)
	}
}

func TestPrintSuggestions(t *testing.T) {
	groups := []similarity.Group{
		{Members: []string{"aaa", "aab", "aac"}},
		{Members: []string{"zzz"}},
	}
	var buf bytes.Buffer
	printSuggestions(&buf, groups, 500, 200)
	out := buf.String()

	if !strings.Contains(out, "Sampled 200 of 500 lines.") {
		t.Errorf("expected sampling note, got: %q", out)
	}
	if !strings.Contains(out, "Found 2 distinct log format(s):") {
		t.Errorf("expected group count, got: %q", out)
	}
	if !strings.Contains(out, "Group 1 (3 lines)") || !strings.Contains(out, "Group 2 (1 line)") {
		t.Errorf("expected per-group counts with correct pluralisation, got: %q", out)
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
