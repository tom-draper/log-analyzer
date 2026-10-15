package parse

import (
	"errors"
	"io"
	"os"
)

var errFollowReset = errors.New("follow source reset")

// FollowFile streams single-line records and follows the path across replacement
// and truncation. Line numbers count completed lines across all file generations.
func FollowFile(path string, config *Config, emit func(Extraction) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := &fileFollower{file: file, path: path}
	defer func() { reader.file.Close() }()
	return ParseStream(reader, config, true, func(record Extraction) error {
		record.SourceFile = path
		return emit(record)
	})
}

type fileFollower struct {
	file   *os.File
	path   string
	offset int64
}

func (r *fileFollower) Read(p []byte) (int, error) {
	info, err := r.file.Stat()
	if err != nil {
		return 0, err
	}
	if info.Size() < r.offset {
		if _, err := r.file.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
		r.offset = 0
		return 0, errFollowReset
	}

	n, err := r.file.Read(p)
	r.offset += int64(n)
	if err != io.EOF {
		return n, err
	}

	current, statErr := os.Stat(r.path)
	if os.IsNotExist(statErr) {
		return n, io.EOF
	}
	if statErr != nil {
		return n, statErr
	}
	if os.SameFile(info, current) {
		return n, io.EOF
	}
	file, openErr := os.Open(r.path)
	if os.IsNotExist(openErr) {
		return n, io.EOF
	}
	if openErr != nil {
		return n, openErr
	}
	r.file.Close()
	r.file = file
	r.offset = 0
	return n, errFollowReset
}
