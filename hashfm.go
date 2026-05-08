// Package hashfm provides utilities for extracting a hashfm from shell scripts.
//
// A hashfm (hash + frontmatter) is a YAML payload delimited by `# ---`
// lines inside a shell script. This package handles extraction and
// normalization by stripping the comment prefix from each line.
package hashfm

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// configFileCandidates lists the config filenames to probe, in order of precedence.
var configFileCandidates = []string{
	".hashfm",
	".hashfm.yml",
	".hashfm.yaml",
	".hashfm.json",
}

// ErrMultipleBlocks indicates that more than one hashfm was detected in the source.
var ErrMultipleBlocks = errors.New("hashfm: multiple hashfms found")

// ErrUnclosedBlock indicates that an opening `# ---` was found without a corresponding closing delimiter.
var ErrUnclosedBlock = errors.New("hashfm: unclosed hashfm")

//go:embed schema/hashfm-config.schema.json
var ConfigSchema []byte

// Extract reads the source string and returns the YAML content of the first
// hashfm it finds. Returns an empty string if no hashfm is present.
// Returns an error if the hashfm structure is malformed.
func Extract(src string) (string, error) {
	lines := strings.Split(src, "\n")

	const delim = "# ---"

	inBlock := false
	found := false
	var yamlLines []string

	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t")

		if trimmed == delim {
			if inBlock {
				// closing delimiter
				if found {
					return "", ErrMultipleBlocks
				}
				found = true
				inBlock = false
			} else {
				// opening delimiter
				inBlock = true
				yamlLines = nil
			}
			continue
		}

		if inBlock {
			// Strip the comment prefix: `# ` or `#  ` etc.
			stripped, ok := stripPrefix(line)
			if !ok {
				// Tolerate lines that are just `#` with no space
				stripped = strings.TrimPrefix(line, "#")
				stripped = strings.TrimLeft(stripped, " ")
			}
			yamlLines = append(yamlLines, stripped)
		}
	}

	if inBlock {
		return "", ErrUnclosedBlock
	}

	if !found {
		return "", nil
	}

	return strings.TrimRight(strings.Join(yamlLines, "\n"), " \t\n"), nil
}

// stripPrefix removes the `# ` prefix (with tolerance for extra spaces).
// Returns the stripped line and true, or the original line and false if the
// prefix is not present.
func stripPrefix(line string) (string, bool) {
	if !strings.HasPrefix(line, "#") {
		return line, false
	}
	rest := line[1:] // after '#'
	if len(rest) == 0 {
		return "", true
	}
	if rest[0] != ' ' {
		return line, false
	}
	// Strip the first space (canonical prefix is `# `)
	return rest[1:], true
}

// LoadConfig finds and loads a config file from the current working directory.
// It validates the file against the hashfm config schema.
// Returns nil if no config file is found.
func LoadConfig() (map[string]interface{}, error) {
	var path string
	for _, name := range configFileCandidates {
		if _, err := os.Stat(name); err == nil {
			path = name
			break
		}
	}

	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rawConfig map[string]interface{}
	if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &rawConfig); err != nil {
			return nil, err
		}
	} else {
		if err := yaml.Unmarshal(data, &rawConfig); err != nil {
			return nil, err
		}
	}

	schemaLoader := gojsonschema.NewBytesLoader(ConfigSchema)
	documentJSON, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, err
	}
	documentLoader := gojsonschema.NewBytesLoader(documentJSON)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	if !result.Valid() {
		return nil, errors.New("hashfm: config file validation failed")
	}

	return rawConfig, nil
}
