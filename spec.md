# hashfm

A small YAML payload inside shell scripts, delimited by `# ---`. We call it a
hashfm – hash + frontmatter.

---

## What is hashfm?

`hashfm` is a minimal shell comment convention for placing structured YAML
inside script files. It defines only the delimiter syntax and content format.
Field definitions are left to implementations.

---

## The Syntax

### Delimiters

The hashfm opens and closes with:

```shell
# ---
```

### Line prefix

Every line inside the hashfm is prefixed with `# ` (hash followed by a single
space). The canonical prefix is `# `. Parsers should be tolerant of additional
spaces after the hash.

### Content

Content between the delimiters is YAML. Any valid YAML structure is permitted —
scalars, mappings, sequences, or a sequence of mappings. The specific structure
is defined by the implementation.

### Example

```bash
#!/usr/bin/env bash
# ---
# key: value
# list:
#   - item one
#   - item two
# ---
```

---

## Rules

- A file may contain at most one hashfm.
- The hashfm may be preceded only by a shebang line and blank lines.
- Opening and closing `# ---` delimiters must appear on their own line.
- Content outside the delimiters is not part of the hashfm and must not be
  parsed as such.
- A file with no hashfm is valid — absence carries no meaning beyond the hashfm
  not being present.

---

## Implementations

An implementation of `hashfm` defines:

- The fields it expects inside the hashfm
- Which fields are mandatory and which are optional
- Validation rules for field values
- What is produced from a parsed hashfm

An implementation must document its field schema separately. It must not redefine the delimiter syntax.

---

## Naming

Implementations should be named using the pattern `hashfm-*`, making their relationship to the base convention explicit.

Examples:

- `hashfm-agent` — an implementation for script discovery and indexing
