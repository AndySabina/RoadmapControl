# Roadmap manifest loading (PR02A-1 / PR02A-2)

PR02A-1 supplies the typed, filesystem-independent parser. PR02A-2 adds a read-only loader for a `roadmap.yaml` manifest and its explicitly declared modules; the whole PR-02 roadmap-validation delivery remains incomplete.

## Accepted document

```yaml
schema: https://schemas.example/v1
modules:
  - trackers/example.yaml
```

`schema` must be an absolute URI; an authority is not required, so URNs and hostless `file:` URIs are valid syntax. This parser does not resolve or fetch the URI; offline schema resolution is deferred. `modules` must be a list of unique, normalized, slash-separated relative `.yaml` or `.yml` paths. Paths cannot be empty, absolute, backslash-separated, `.` or parent-traversing.

## Parser safety limits

The parser accepts exactly one UTF-8 YAML document of at most 1 MiB and depth at most 64. It rejects duplicate keys, aliases, unknown manifest fields, ambiguous non-string scalar values, malformed trailing content, and unsafe or duplicate module paths.

## Filesystem loader boundary

`filesystem.Load(root)` accepts a non-symlink directory containing a regular, at-most-1-MiB `roadmap.yaml`. It parses the manifest, reads only its declared regular YAML modules (nested paths and multiple modules are allowed), validates each module's YAML document safety, and rejects unlisted `.yaml`/`.yml` files case-insensitively. It rejects symlink roots, manifest/modules, and any intermediate module-path component.

Loading is non-mutating and rejects unsafe files before reading them, so it does not intentionally read a FIFO. This is a cooperative boundary: a hostile concurrent filesystem change can still race check-then-read operations; it does not claim universal TOCTOU safety.

## Not included

The loader does not discover implicit modules, return catchall module content, resolve schemas, validate typed module fields, produce hashes, or expose a CLI. Schema resolution, hashing, module typed contracts, and CLI integration remain deferred.
