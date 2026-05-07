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
# name: deploy
# description: Deploy the application to staging
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

## The `.hashfm` Config File

Separately, `hashfm` tools share a single `.hashfm` config file at the project
root. Each tool owns a top-level namespace:

```yaml
version: "1.0"
project:
  name: "my-project"

hashfm-agent:
  generate:
    format: tsv
    recursive: true
```

This keeps configuration unified — one file, all tools.

---

## Block vs Config

| | hashfm | Config |
|---|---|---|
| **Location** | Inside script files | Project root (`.hashfm`) |
| **Owned by** | The script itself | Each tool owns its namespace |
| **Purpose** | Script metadata | Tool settings |

---

## Implementations

| Tool | What it does |
|------|-------------|
| [hashfm-agent](https://github.com/sidisinsane/hashfm-agent) | Extracts hashfms from scripts and produces a machine-readable index |

---

## Go Package

```bash
go get github.com/sidisinsane/hashfm
```

```go
import "github.com/sidisinsane/hashfm"

// Extract reads a shell script source and returns its hashfm.
yamlContent, err := hashfm.Extract(scriptSource)
```

`hashfm.Extract` returns the cleaned YAML content from the first `# ---` block,
or an empty string if no hashfm is found.

---

## Schema

The hashfm syntax itself has no required fields — field definitions are left to
implementations. The `.hashfm` config file schema is in `schema/`.

---

## Development

### Prerequisites

- Go 1.26+
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
