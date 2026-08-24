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
			name: "follow record limit",
			args: []string{"--follow", "--follow-max-records", "500", "app.log"},
			want: options{logPaths: []string{"app.log"}, follow: true, followMaxRecords: 500, followMaxRecordsSet: true},
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
		{name: "negative follow record limit", args: []string{"--follow-max-records", "-1", "app.log"}},
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

func TestTokenizeLine(t *testing.T) {
	segs := tokenizeLine(`[11/Dec:01] 1.2.3.4 "GET /a HTTP/1.1" 200`)
	if len(segs) != 4 {
		t.Fatalf("expected 4 segments, got %d: %+v", len(segs), segs)
	}
	if segs[0].kind != '[' || segs[0].text != "[11/Dec:01]" {
		t.Errorf("segment 0 = %+v, want bracketed timestamp", segs[0])
	}
	if segs[2].kind != '"' || segs[2].text != `"GET /a HTTP/1.1"` {
		t.Errorf("segment 2 = %+v, want quoted request as one unit", segs[2])
	}
}

func TestGeneratePattern(t *testing.T) {
	lines := []string{
		`[11/Dec/2023:11:01:28] 220.203.23.174 "GET /a HTTP/1.1" 200 182`,
		`[11/Dec/2023:11:01:29] 89.238.65.53 "POST /b HTTP/1.1" 201 145`,
	}
	counter := 0
	pattern, tokens, ok := generatePattern(lines, &counter)
	if !ok {
		t.Fatal("expected generation to succeed for aligned lines")
	}
	// timestamp, ip, and (method, endpoint) inside the quoted request, plus
	// status and bytes vary -> 6 tokens; the HTTP version stays literal.
	if len(tokens) != 6 {
		t.Errorf("expected 6 tokens, got %d: %v (pattern %q)", len(tokens), tokens, pattern)
	}
	if !strings.Contains(pattern, "[field1]") {
		t.Errorf("expected bracket kept as literal around token: %q", pattern)
	}
	if !strings.Contains(pattern, "HTTP/1.1") {
		t.Errorf("expected constant HTTP version kept literal inside the quotes: %q", pattern)
	}
	// The quoted request is sub-split rather than a single token.
	if strings.Contains(pattern, `"field3"`) {
		t.Errorf("expected quoted request to be sub-split, not a single token: %q", pattern)
	}
}

func TestGeneratePatternMisaligned(t *testing.T) {
	counter := 0
	if _, _, ok := generatePattern([]string{"a b c", "a b"}, &counter); ok {
		t.Error("expected generation to fail for lines with different field counts")
	}
}

// TestGeneratedPatternParses is the important round-trip: a generated pattern,
// fed back through the real parser, must actually match its source lines and
// extract every token.
func TestGeneratedPatternParses(t *testing.T) {
	lines := []string{
		`[t1] 1.1.1.1 "GET /a HTTP/1.1" 200 10`,
		`[t2] 2.2.2.2 "POST /b HTTP/1.1" 404 20`,
	}
	counter := 0
	pattern, tokens, ok := generatePattern(lines, &counter)
	if !ok {
		t.Fatal("expected generation to succeed")
	}

	config := parse.Config{Patterns: []string{pattern}, Tokens: tokens}
	extraction, err := parse.ParseLines(lines, &config)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(extraction) != len(lines) {
		t.Fatalf("expected %d extractions, got %d", len(lines), len(extraction))
	}
	for i, e := range extraction {
		if e.Pattern == "" {
			t.Errorf("line %d did not match the generated pattern %q", i, pattern)
		}
		if len(e.Params) != len(tokens) {
			t.Errorf("line %d extracted %d params, want %d: %+v", i, len(e.Params), len(tokens), e.Params)
		}
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
