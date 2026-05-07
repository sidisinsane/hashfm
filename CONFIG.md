# The `.hashfm` Config File

`hashfm` tools share a single `.hashfm` file at the project root for persistent
settings. This document describes the design of that file.

---

## Design

- **One file per project.** All tools read the same file.
- **Namespaced top-level keys.** Each tool owns its namespace and is responsible
  only for validating that namespace.
- **Unified configuration.** Tool settings live alongside project metadata, not
  scattered across per-tool config files.

---

## Structure

The file is a YAML document with top-level keys as namespaces:

```yaml
version: "1.0"
project:
  name: "my-project"

hashfm-agent:
  generate:
    format: tsv
    recursive: true
```

- `version` and `project` are owned by `hashfm`.
- `hashfm-agent` is owned by the agent. It can be absent — the tool then uses
  its defaults.
- Future tools add their own namespace in the same way.

---

## Schema Ownership

Each namespace has its own schema. The schema for a tool's namespace lives in
that tool's repository, not in `hashfm`.

```text
hashfm/schema/hashfm-config.schema.json   — owns version, project, any future core keys
hashfm-agent/...                          — owns hashfm-agent namespace
hashfm-lint/...                           — owns hashfm-lint namespace (future)
```

The core schema uses `$ref` to point to each tool's schema. This keeps schemas
with the code that owns them.

---

## Validation

When `hashfm-agent` runs, it reads the `.hashfm` file and validates only its own
namespace (`hashfm-agent`) against its schema. It ignores all other namespaces.
A missing namespace is valid — tools fall back to defaults.

---

## Precedence

Settings resolve in this order (highest to lowest):

1. **CLI flags** — explicitly provided at runtime
2. **`.hashfm` config file** — persistent project settings
3. **Hardcoded defaults** — sensible tool defaults
