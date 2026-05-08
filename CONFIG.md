# The Hashfm Config File

`hashfm` tools share a single `.hashfm` file at the project root for persistent
settings. This document describes the design of that file.

---

## Design

- **One file per project.** All tools read the same file.
- **Namespaced top-level keys.** Each tool owns its namespace following the
  `hashfm-*` naming convention and is responsible for validating that namespace.
- **Unified configuration.** Tool settings live alongside project metadata, not
  scattered across per-tool config files.

---

## Supported Filenames

The following filenames are probed in order, stopping at the first match:

| Filename | Format |
|----------|--------|
| `.hashfm` | YAML (assumed) |
| `.hashfm.yml` | YAML |
| `.hashfm.yaml` | YAML |
| `.hashfm.json` | JSON |

---

## Structure

The file is a YAML document with top-level keys as namespaces:

```yaml
version: "1.1.0"
project:
  name: "my-project"

hashfm-agent:
  generate:
    format: tsv
    recursive: true
    output: index.tsv
```

- `version` and `project` are owned by `hashfm`.
- `hashfm-*` namespaces are owned by the respective tools. They can be absent —
  the tool then uses its defaults.
- Future tools add their own namespace in the same way.

---

## Schema Ownership

The core schema (`hashfm/schema/hashfm-config.schema.json`) uses a
`patternProperties` rule to enforce the `hashfm-*` naming convention without
knowing any specific tool names:

```json
"patternProperties": {
  "^hashfm-[a-z][a-z0-9-]*$": {
    "type": "object",
    "description": "Tool-specific configuration."
  }
}
```

Each tool owns its namespace schema:

```text
hashfm/schema/hashfm-config.schema.json    — owns version, project, namespace pattern
hashfm-agent/schema/...                    — owns hashfm-agent namespace
hashfm-lint/schema/...                     — owns hashfm-lint namespace (future)
```

---

## Validation

When a tool runs, it reads the config file and validates it against the core
schema first. It then validates only its own namespace against its tool-specific
schema. A missing namespace is valid — tools fall back to defaults.

---

## Precedence

Settings resolve in this order (highest to lowest):

1. **CLI flags** — explicitly provided at runtime
2. **`.hashfm` config file** — persistent project settings
3. **Hardcoded defaults** — sensible tool defaults
