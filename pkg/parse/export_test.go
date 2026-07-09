package parse

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func exportSample() []Extraction {
	ts := time.Date(2023, 12, 11, 11, 1, 28, 0, time.UTC)
	return []Extraction{
		{
			Params: map[string]Param{
				"timestamp": {Value: ts, Type: "time"},
				"status":    {Value: 200, Type: "int"},
				"endpoint":  {Value: "/home", Type: "str"},
			},
			Pattern:    "timestamp status endpoint",
			LineNumber: 0,
			Line:       "line one",
		},
		{
			// Missing "endpoint" to exercise empty CSV cells.
			Params: map[string]Param{
				"timestamp": {Value: ts, Type: "time"},
				"status":    {Value: 500, Type: "int"},
			},
			Pattern:    "timestamp status",
			LineNumber: 4,
			Line:       "line two",
		},
	}
}

func TestExportCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, "csv", exportSample()); err != nil {
		t.Fatalf("export csv: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected header + 2 rows, got %d lines: %q", len(lines), lines)
	}
	// Columns are sorted; header starts with lineNumber.
	if lines[0] != "lineNumber,endpoint,status,timestamp" {
		t.Errorf("unexpected header: %q", lines[0])
	}
	// RFC3339 timestamp and 1-indexed line number.
	if !strings.HasPrefix(lines[1], "1,/home,200,2023-12-11T11:01:28Z") {
		t.Errorf("unexpected row 1: %q", lines[1])
	}
	// Missing endpoint renders as an empty cell.
	if !strings.HasPrefix(lines[2], "5,,500,") {
		t.Errorf("unexpected row 2 (missing cell): %q", lines[2])
	}
}

func TestExportNDJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, "ndjson", exportSample()); err != nil {
		t.Fatalf("export ndjson: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 json lines, got %d", len(lines))
	}
	for i, line := range lines {
		if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
			t.Errorf("line %d is not a single JSON object: %q", i, line)
		}
	}
}

func TestExportJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, "json", exportSample()); err != nil {
		t.Fatalf("export json: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("expected a JSON array, got: %q", out[:min(40, len(out))])
	}
}

func TestExportUnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, "yaml", exportSample()); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestFormatFromPath(t *testing.T) {
	cases := map[string]string{
		"out.json":       "json",
		"a/b/out.ndjson": "ndjson",
		"data.csv":       "csv",
		"data.txt":       "",
		"noext":          "",
	}
	for path, want := range cases {
		if got := FormatFromPath(path); got != want {
			t.Errorf("FormatFromPath(%q) = %q, want %q", path, got, want)
		}
	}
}
