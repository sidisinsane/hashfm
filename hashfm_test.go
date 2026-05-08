package hashfm_test

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/sidisinsane/hashfm"
)

//go:embed testdata/scripts/* testdata/config/*
var testFixtures embed.FS

const scriptsTestdataDir = "testdata/scripts"
const configTestdataDir = "testdata/config"

func readScriptFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := testFixtures.ReadFile(filepath.Join(scriptsTestdataDir, name))
	if err != nil {
		t.Fatalf("read script %s: %v", name, err)
	}
	return string(data)
}

func readConfigFixture(name string) ([]byte, error) {
	return testFixtures.ReadFile(filepath.Join(configTestdataDir, name))
}

func TestExtract_SingleCommand(t *testing.T) {
	got, err := hashfm.Extract(readScriptFixture(t, "valid-single.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty YAML")
	}
	for _, field := range []string{"description:", "usage:", "exits:"} {
		if !contains(got, field) {
			t.Errorf("expected %q in output", field)
		}
	}
}

func TestExtract_MultiCommand(t *testing.T) {
	got, err := hashfm.Extract(readScriptFixture(t, "valid-multi.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty YAML")
	}
	if !hasPrefix(got, "-") {
		t.Errorf("expected multi-command YAML to start with '-', got:\n%s", got)
	}
}

func TestExtract_NoBlock(t *testing.T) {
	got, err := hashfm.Extract(readScriptFixture(t, "no-block.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got: %q", got)
	}
}

func TestExtract_UnclosedBlock(t *testing.T) {
	_, err := hashfm.Extract(readScriptFixture(t, "invalid-unclosed-block.sh"))
	if err != hashfm.ErrUnclosedBlock {
		t.Errorf("expected ErrUnclosedBlock, got: %v", err)
	}
}

func TestExtract_MultipleBlocks(t *testing.T) {
	_, err := hashfm.Extract(readScriptFixture(t, "invalid-multiple-blocks.sh"))
	if err != hashfm.ErrMultipleBlocks {
		t.Errorf("expected ErrMultipleBlocks, got: %v", err)
	}
}

func TestExtract_ExtraSpacesAfterHash(t *testing.T) {
	got, err := hashfm.Extract(readScriptFixture(t, "valid-extra-spaces.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(got, "key: value") {
		t.Errorf("expected 'key: value', got: %q", got)
	}
}

func TestExtract_TrailingWhitespace(t *testing.T) {
	got, err := hashfm.Extract(readScriptFixture(t, "valid-trailing-whitespace.sh"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(got, "key: value") {
		t.Errorf("expected 'key: value', got: %q", got)
	}
	if got != "key: value" {
		t.Errorf("expected trailing whitespace stripped, got: %q", got)
	}
}

func TestLoadConfig_Discovery(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer os.Chdir(origCwd)

	tests := []struct {
		name       string
		files      []string
		wantConfig bool
	}{
		{
			name:       "no config returns nil",
			files:      nil,
			wantConfig: false,
		},
		{
			name:       ".hashfm takes precedence",
			files:      []string{".hashfm", ".hashfm.yml"},
			wantConfig: true,
		},
		{
			name:       ".hashfm.yml discovered",
			files:      []string{".hashfm.yml"},
			wantConfig: true,
		},
		{
			name:       ".hashfm.json discovered",
			files:      []string{".hashfm.json"},
			wantConfig: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Chdir: %v", err)
			}

			for _, f := range tc.files {
				data, err := readConfigFixture(f)
				if err != nil {
					t.Fatalf("read fixture %s: %v", f, err)
				}
        
				if err := os.WriteFile(f, data, 0644); err != nil {
					t.Fatalf("write %s: %v", f, err)
				}
			}

			cfg, err := hashfm.LoadConfig()
			if err != nil {
				t.Fatalf("LoadConfig: %v", err)
			}

			if tc.wantConfig {
				if cfg == nil {
					t.Fatal("expected config, got nil")
				}
				project, ok := cfg["project"].(map[string]interface{})
				if !ok {
					t.Fatal("expected project field")
				}
				if _, ok := project["name"]; !ok {
					t.Fatal("expected project.name field")
				}
			} else {
				if cfg != nil {
					t.Errorf("expected nil, got %v", cfg)
				}
			}
		})
	}
}

func TestLoadConfig_ValidSchema(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer os.Chdir(origCwd)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	data, err := readConfigFixture("valid-minimal.yaml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(".hashfm", data, 0644); err != nil {
		t.Fatalf("write .hashfm: %v", err)
	}

	cfg, err := hashfm.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
}

func TestLoadConfig_InvalidSchema(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer os.Chdir(origCwd)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	data, err := readConfigFixture("invalid-no-project-name.yaml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(".hashfm", data, 0644); err != nil {
		t.Fatalf("write .hashfm: %v", err)
	}

	_, err = hashfm.LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid schema, got nil")
	}
}

// String helpers to avoid importing strings in tests
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr) != -1
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && searchString(s, prefix) == 0
}

func searchString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
