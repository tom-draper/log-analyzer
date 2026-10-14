package parse

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFollowFileReturnsEmitterError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	if err := os.WriteFile(path, []byte("before\nafter\n"), 0600); err != nil {
		t.Fatal(err)
	}
	config := Config{Tokens: []string{"str_message"}, Patterns: []string{"str_message"}}
	stop := errors.New("stop following")
	var got []Extraction
	err := FollowFile(path, &config, func(e Extraction) error {
		got = append(got, e)
		return stop
	})
	if err != stop {
		t.Fatalf("error = %v, want emitter error", err)
	}
	want := []Extraction{{Line: "before", LineNumber: 0, Pattern: "str_message", Params: map[string]Param{"str_message": {Value: "before", Type: "str"}}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("records = %+v, want %+v", got, want)
	}
}

func TestFileFollower(t *testing.T) {
	for _, operation := range []string{"append", "replace", "truncate", "temporarily missing"} {
		t.Run(operation, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "app.log")
			write := func(text string) {
				t.Helper()
				if err := os.WriteFile(path, []byte(text), 0600); err != nil {
					t.Fatal(err)
				}
			}
			write("before rotation\n")
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			reader := &fileFollower{file: file, path: path}
			defer func() { reader.file.Close() }()
			if got, err := io.ReadAll(reader); err != nil || string(got) != "before rotation\n" {
				t.Fatalf("initial read = %q, %v", got, err)
			}
			if operation == "append" {
				file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					t.Fatal(err)
				}
				_, err = file.WriteString("after\n")
				file.Close()
				if err != nil {
					t.Fatal(err)
				}
			} else {
				if operation != "truncate" {
					if err := os.Rename(path, path+".1"); err != nil {
						t.Fatal(err)
					}
				}
				if operation == "temporarily missing" {
					if got, err := io.ReadAll(reader); err != nil || len(got) != 0 {
						t.Fatalf("missing path read = %q, %v", got, err)
					}
				}
				write("after\n")
			}
			// A source reset may interrupt one read before the next file starts.
			got, err := io.ReadAll(reader)
			if err == errFollowReset && len(got) == 0 {
				got, err = io.ReadAll(reader)
			}
			if err != nil || string(got) != "after\n" {
				t.Fatalf("read after %s = %q, %v, want after\\n", operation, got, err)
			}
		})
	}
}

type resettingReader struct {
	current io.Reader
	next    io.Reader
}

func (r *resettingReader) Read(p []byte) (int, error) {
	n, err := r.current.Read(p)
	if err == io.EOF && r.next != nil {
		r.current, r.next = r.next, nil
		return n, errFollowReset
	}
	return n, err
}

func TestParseStreamDiscardsPartialLineOnReset(t *testing.T) {
	reader := &resettingReader{
		current: strings.NewReader("before\nold incomplete"),
		next:    strings.NewReader("after\n"),
	}
	config := Config{Tokens: []string{"str_message"}, Patterns: []string{"str_message"}}
	var got []Extraction
	err := ParseStream(reader, &config, false, func(e Extraction) error {
		got = append(got, e)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []Extraction{
		{Line: "before", LineNumber: 0, Pattern: "str_message", Params: map[string]Param{"str_message": {Value: "before", Type: "str"}}},
		{Line: "after", LineNumber: 1, Pattern: "str_message", Params: map[string]Param{"str_message": {Value: "after", Type: "str"}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("records = %+v, want %+v", got, want)
	}
}
