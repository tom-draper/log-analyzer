package parse

import (
	"fmt"
	"io"
	"os"
)

func DisplayLines(extraction []Extraction) {
	fprintLines(os.Stdout, extraction)
}

func DisplayTestLines(extraction []Extraction) {
	fprintTestLines(os.Stdout, extraction)
}

func fprintLines(w io.Writer, extraction []Extraction) {
	fmt.Fprintf(w, "%d lines\n", len(extraction))
	for i, e := range extraction {
		fmt.Fprintf(w, "line %d\n", i+1)
		for k, v := range e.Params {
			fmt.Fprintf(w, "	%s(%s): %v\n", k, v.Type, v.Value)
		}
	}
}

// FprintExtraction writes a single extraction in the human-readable form used
// by the print output. It is used when streaming records one at a time.
func FprintExtraction(w io.Writer, e Extraction) {
	fmt.Fprintf(w, "line %d: %s\n", e.LineNumber+1, e.Line)
	for k, v := range e.Params {
		fmt.Fprintf(w, "	%s(%s): %v\n", k, v.Type, v.Value)
	}
}

func fprintTestLines(w io.Writer, extraction []Extraction) {
	fmt.Fprintf(w, "%d lines\n", len(extraction))
	for _, e := range extraction {
		fmt.Fprintf(w, "line %d: %s\n", e.LineNumber+1, e.Line)
		for k, v := range e.Params {
			fmt.Fprintf(w, "	%s(%s): %v\n", k, v.Type, v.Value)
		}
	}
}
