# Roadmap manifest parser (PR02A-1)

PR02A-1 supplies only the typed, filesystem-independent parser for a `roadmap.yaml` manifest. The whole PR-02 roadmap-validation delivery remains incomplete.

## Accepted document

```yaml
schema: https://schemas.example/v1
modules:
  - trackers/example.yaml
```

`schema` must be an absolute URI; an authority is not required, so URNs and hostless `file:` URIs are valid syntax. This parser does not resolve or fetch the URI; offline schema resolution is deferred. `modules` must be a list of unique, normalized, slash-separated relative `.yaml` or `.yml` paths. Paths cannot be empty, absolute, backslash-separated, `.` or parent-traversing.

## Parser safety limits

The parser accepts exactly one UTF-8 YAML document of at most 1 MiB and depth at most 64. It rejects duplicate keys, aliases, unknown manifest fields, ambiguous non-string scalar values, malformed trailing content, and unsafe or duplicate module paths.

## Not included

This slice does not read the filesystem, discover or load modules, resolve schemas, validate module contents, produce hashes, or expose a CLI. The deferred filesystem loader is preserved intact for PR02A-2 at `/home/andyf/Projects/RoadmapControl-local-archive/pr02a-before-split-20260905/` (see `SHA256SUMS`).
