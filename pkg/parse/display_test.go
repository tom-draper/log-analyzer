package parse

import (
	"bytes"
	"strings"
	"testing"
)

func sampleExtraction() []Extraction {
	return []Extraction{
		{
			Params: map[string]Param{
				"status": {Value: 200, Type: "int"},
				"method": {Value: "GET", Type: "str"},
			},
			Pattern:    "method status",
			LineNumber: 0,
			Line:       "GET 200",
		},
		{
			Params: map[string]Param{
				"message": {Value: "boot complete", Type: "str"},
			},
			Pattern:    "message",
			LineNumber: 4,
			Line:       "boot complete",
		},
	}
}

func TestDisplayLines(t *testing.T) {
	var buf bytes.Buffer
	fprintLines(&buf, sampleExtraction())
	out := buf.String()

	if !strings.HasPrefix(out, "2 lines\n") {
		t.Errorf("expected line count header, got: %q", out)
	}
	// Lines are numbered from 1 regardless of the origin LineNumber.
	if !strings.Contains(out, "line 1\n") || !strings.Contains(out, "line 2\n") {
		t.Errorf("expected sequential line headers, got: %q", out)
	}
	// Params render as token(type): value, not the Param struct.
	if !strings.Contains(out, "status(int): 200") {
		t.Errorf("expected typed param output, got: %q", out)
	}
	if strings.Contains(out, "parse.Param") || strings.Contains(out, "{200 int}") {
		t.Errorf("param should not render as a struct, got: %q", out)
	}
}

func TestDisplayTestLines(t *testing.T) {
	var buf bytes.Buffer
	fprintTestLines(&buf, sampleExtraction())
	out := buf.String()

	if !strings.HasPrefix(out, "2 lines\n") {
		t.Errorf("expected line count header, got: %q", out)
	}
	// Test output uses the origin line number (1-indexed) and the raw line.
	if !strings.Contains(out, "line 1: GET 200") {
		t.Errorf("expected origin line 1, got: %q", out)
	}
	if !strings.Contains(out, "line 5: boot complete") {
		t.Errorf("expected origin line 5, got: %q", out)
	}
	if !strings.Contains(out, "message(str): boot complete") {
		t.Errorf("expected typed param output, got: %q", out)
	}
}
