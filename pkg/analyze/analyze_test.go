package analyze

import (
	"testing"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func TestRetainRecent(t *testing.T) {
	extractions := []parse.Extraction{{LineNumber: 1}, {LineNumber: 2}, {LineNumber: 3}}

	got := retainRecent(extractions, 2)
	if len(got) != 2 || got[0].LineNumber != 2 || got[1].LineNumber != 3 {
		t.Errorf("retainRecent() = %#v, want last two records", got)
	}
	if got := retainRecent(extractions, 0); len(got) != 3 {
		t.Errorf("retainRecent(..., 0) kept %d records, want all", len(got))
	}
}
