# hashfm

A small YAML payload inside shell scripts, delimited by `# ---`. We call it a
hashfm – hash + frontmatter.

`hashfm` (hash + frontmatter) is a shell comment convention for embedding
structured, machine-readable metadata in scripts. It solves the problem Bash has
no standard for: script-level documentation that is machine-readable.

Other languages have JSDoc, Python docstrings, Rust doc comments. Bash has
nothing. `hashfm` aims to fills that gap.

---

## The Syntax

A `hashfm` is a YAML payload delimited by `# ---` inside a shell script:

```bash
#!/usr/bin/env bash
# ---
# key: value
# list:
#   - item one
#   - item two
# ---
```

Every line inside the delimiter is a shell comment (`# `). The script remains
fully executable by any shell.

### Rules

- A file may contain at most one hashfm.
- The hashfm may be preceded only by a shebang line and blank lines.
- Opening and closing `# ---` delimiters must appear on their own line.
- Content between the delimiters is YAML.
- A file with no hashfm is valid — absence carries no meaning.

---

## Public API

The `hashfm` package exposes two functions for implementations:

| Function | Purpose |
|---|---|
| `Extract(source string) (string, error)` | Reads a shell script and returns the YAML content of the first hashfm block. Returns an empty string if no block is found. |
| `LoadConfig() (map[string]interface{}, error)` | Finds and loads the config file from the current directory. Returns `nil` if no config file is found. Validates against the core schema. |

```go
import "github.com/sidisinsane/hashfm"

// Extract the hashfm block from a script.
yamlContent, err := hashfm.Extract(scriptSource)

// Load the config file.
cfg, err := hashfm.LoadConfig()
```

---

## Configuration

`hashfm` tools share a single `.hashfm` config file at the project root for
persistent settings. It is a YAML document with top-level keys as namespaces.

The core schema defines the envelope (`version`, `project`) and allows
tool-specific namespaces under the `hashfm-*` naming convention. Each tool
validates its own namespace. No tool registry is needed.

See [`CONFIG.md`](CONFIG.md) for the full specification — supported filenames,
schema, validation, and resolution order.

The config file schema lives in [`schema/hashfm-config.schema.json`](schema/hashfm-config.schema.json).

---

## Documentation

| Document | What it covers |
|----------|---------------|
| `spec.md` | Syntax, rules, and naming convention |
| [`CONFIG.md`](CONFIG.md) | `.hashfm` config file design and implementation |

---

## Implementations

| Tool | What it does |
|------|-------------|
| [hashfm-agent](https://github.com/sidisinsane/hashfm-agent) | Extracts hashfms from scripts and produces a machine-readable index |

Implementations should be named using the pattern `hashfm-*`, making their
relationship to the base convention explicit.

---

## Development

### Prerequisites

- Go 1.26.2
- [golangci-lint](https://golangci-lint.run) — for Go linting
- [lefthook](https://lefthook.dev) — for pre-commit hooks

### Setup

```bash
# Install lefthook
brew install lefthook

# Enable hooks
lefthook install
```

### Linting

Go files are linted with `golangci-lint`. Run manually:

```bash
golangci-lint run ./...
```

Hooks run automatically on commit via lefthook.
