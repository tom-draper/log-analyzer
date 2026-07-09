package parse

import (
	"strings"
	"testing"
)

func TestReadLines(t *testing.T) {
	longLine := strings.Repeat("x", 200000) // exceeds bufio.Scanner's default 64KB

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "unix newlines", input: "a\nb\nc\n", want: []string{"a", "b", "c"}},
		{name: "windows newlines", input: "a\r\nb\r\n", want: []string{"a", "b"}},
		{name: "no trailing newline", input: "a\nb", want: []string{"a", "b"}},
		{name: "blank lines preserved", input: "a\n\nb\n", want: []string{"a", "", "b"}},
		{name: "long line", input: longLine + "\n", want: []string{longLine}},
		{name: "empty input", input: "", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLines(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines %q, want %d lines %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseReaderMatchesParse(t *testing.T) {
	config := Config{
		Patterns: []string{"[timestamp] ip rest"},
		Tokens:   []string{"timestamp", "ip", "rest"},
	}
	logText := "[11/Dec/2023:11:01:28] 220.203.23.174 GET /home\n" +
		"[11/Dec/2023:11:01:29] 89.238.65.53 POST /users\n"

	fromString, err := Parse(logText, &config)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	fromReader, err := ParseReader(strings.NewReader(logText), &config)
	if err != nil {
		t.Fatalf("ParseReader error: %v", err)
	}

	if len(fromReader) != len(fromString) {
		t.Fatalf("ParseReader produced %d extractions, Parse produced %d", len(fromReader), len(fromString))
	}
	for i := range fromReader {
		if fromReader[i].Line != fromString[i].Line || fromReader[i].Pattern != fromString[i].Pattern {
			t.Errorf("extraction %d mismatch: reader=%+v string=%+v", i, fromReader[i], fromString[i])
		}
	}
}
