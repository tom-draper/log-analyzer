package parse

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Export writes the extracted records to w in the given format. Supported
// formats are "json" (a pretty-printed array), "ndjson" (one JSON object per
// line) and "csv" (a table with one column per token).
func Export(w io.Writer, format string, extractions []Extraction) error {
	switch format {
	case "json":
		return exportJSON(w, extractions)
	case "ndjson":
		return exportNDJSON(w, extractions)
	case "csv":
		return exportCSV(w, extractions)
	default:
		return fmt.Errorf("unknown format %q (want json, ndjson or csv)", format)
	}
}

func exportJSON(w io.Writer, extractions []Extraction) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(extractions)
}

func exportNDJSON(w io.Writer, extractions []Extraction) error {
	encoder := json.NewEncoder(w)
	for _, e := range extractions {
		if err := encoder.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

func exportCSV(w io.Writer, extractions []Extraction) error {
	// Columns are the union of all token names, sorted for stable output.
	columnSet := make(map[string]struct{})
	for _, e := range extractions {
		for token := range e.Params {
			columnSet[token] = struct{}{}
		}
	}
	columns := make([]string, 0, len(columnSet))
	for token := range columnSet {
		columns = append(columns, token)
	}
	sort.Strings(columns)

	cw := csv.NewWriter(w)
	if err := cw.Write(append([]string{"lineNumber"}, columns...)); err != nil {
		return err
	}
	for _, e := range extractions {
		row := make([]string, 0, len(columns)+1)
		row = append(row, strconv.Itoa(e.LineNumber+1))
		for _, token := range columns {
			if p, ok := e.Params[token]; ok {
				row = append(row, formatValue(p))
			} else {
				row = append(row, "")
			}
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// formatValue renders a parameter value as a string for tabular output.
func formatValue(p Param) string {
	switch v := p.Value.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	case net.IP:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FormatFromPath infers an export format from a file extension, returning an
// empty string when the extension is unrecognised.
func FormatFromPath(path string) string {
	switch {
	case strings.HasSuffix(path, ".json"):
		return "json"
	case strings.HasSuffix(path, ".ndjson"):
		return "ndjson"
	case strings.HasSuffix(path, ".csv"):
		return "csv"
	default:
		return ""
	}
}
