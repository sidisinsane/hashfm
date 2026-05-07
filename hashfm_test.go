package hashfm_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sidisinsane/hashfm"
)

const testdir = "testdata"

func readFile(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(testdir, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}

func TestExtract_SingleCommand(t *testing.T) {
	got, err := hashfm.Extract(readFile(t, "valid-single.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty YAML")
	}
	for _, field := range []string{"description:", "usage:", "exits:"} {
		if !strings.Contains(got, field) {
			t.Errorf("expected %q in output", field)
		}
	}
}

func TestExtract_MultiCommand(t *testing.T) {
	got, err := hashfm.Extract(readFile(t, "valid-multi.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty YAML")
	}
	if !strings.HasPrefix(strings.TrimSpace(got), "-") {
		t.Errorf("expected multi-command YAML to start with '-', got:\n%s", got)
	}
}

func TestExtract_NoBlock(t *testing.T) {
	got, err := hashfm.Extract(readFile(t, "no-block.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got: %q", got)
	}
}

func TestExtract_UnclosedBlock(t *testing.T) {
	_, err := hashfm.Extract(readFile(t, "invalid-unclosed-block.sh"))
	if err != hashfm.ErrUnclosedBlock {
		t.Errorf("expected ErrUnclosedBlock, got: %v", err)
	}
}

func TestExtract_MultipleBlocks(t *testing.T) {
	_, err := hashfm.Extract(readFile(t, "invalid-multiple-blocks.sh"))
	if err != hashfm.ErrMultipleBlocks {
		t.Errorf("expected ErrMultipleBlocks, got: %v", err)
	}
}

func TestExtract_ExtraSpacesAfterHash(t *testing.T) {
	got, err := hashfm.Extract(readFile(t, "valid-extra-spaces.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "key: value") {
		t.Errorf("expected 'key: value', got: %q", got)
	}
}

func TestExtract_TrailingWhitespace(t *testing.T) {
	got, err := hashfm.Extract(readFile(t, "valid-trailing-whitespace.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "key: value") {
		t.Errorf("expected 'key: value', got: %q", got)
	}
}