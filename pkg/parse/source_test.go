package parse

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFilesPreservesSource(t *testing.T) {
	config := Config{Tokens: []string{"str_message"}, Patterns: []string{"str_message"}}
	paths := []string{filepath.Join(t.TempDir(), "first.log"), filepath.Join(t.TempDir(), "second.log")}
	for _, path := range paths {
		if err := os.WriteFile(path, []byte("hello\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	records, err := ParseFiles(paths, &config)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}
	for i, record := range records {
		if record.SourceFile != paths[i] || record.LineNumber != 0 || record.Line != "hello" {
			t.Errorf("record %d = %+v, want source %q and local line 0", i, record, paths[i])
		}
	}
}
