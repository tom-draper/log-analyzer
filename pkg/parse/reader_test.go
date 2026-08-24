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

func TestParseStream(t *testing.T) {
	config := Config{
		Patterns: []string{"[timestamp] ip rest"},
		Tokens:   []string{"timestamp", "ip", "rest"},
	}
	// Blank middle line is skipped; last line has no trailing newline.
	input := "[11/Dec/2023:11:01:28] 1.2.3.4 GET /a\n\n[11/Dec/2023:11:01:29] 5.6.7.8 POST /b"

	var got []Extraction
	err := ParseStream(strings.NewReader(input), &config, false, func(e Extraction) error {
		got = append(got, e)
		return nil
	})
	if err != nil {
		t.Fatalf("ParseStream error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d: %+v", len(got), got)
	}
	// Line numbers account for the skipped blank line at index 1.
	if got[0].LineNumber != 0 || got[1].LineNumber != 2 {
		t.Errorf("unexpected line numbers: %d, %d", got[0].LineNumber, got[1].LineNumber)
	}
	for i, e := range got {
		if e.Pattern == "" {
			t.Errorf("record %d did not match a pattern: %q", i, e.Line)
		}
		if _, ok := e.Params["ip"]; !ok {
			t.Errorf("record %d missing ip param: %+v", i, e.Params)
		}
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

func TestParseFileTestReturnsReadError(t *testing.T) {
	_, err := ParseFileTest("does-not-exist.log", &Config{})
	if err == nil {
		t.Fatal("expected an error for a missing log file")
	}
}
