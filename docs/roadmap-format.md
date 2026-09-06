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

## Policy module contract (PR02B1)

`ParsePolicyRoadmap(manifest, policyYAML)` reads the only currently typed module contract. Its manifest must name exactly the immutable `RoadmapSchemaURI` value, `https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json`, and exactly one module: `policy.yaml`. Other absolute URIs remain valid manifest syntax but are not policy contracts.

The policy document is one safe YAML mapping with exactly these required fields:

```yaml
kind: policy
additional_tracker_types: [initiative]
```

`additional_tracker_types` is a sequence of strings. Empty strings, duplicate names, and built-in tracker names are structurally retained; later business-policy work decides their meaning. Unknown, duplicate, missing, null, or wrong-type fields are rejected. The returned `PolicyRoadmap` has read-only accessors (`Schema`, `ModulePath`, and `AdditionalTrackerTypes`); the slice accessor returns a defensive copy.

## Canonical policy JSON (PR02B2)

`PolicyRoadmap.CanonicalJSON()` returns RFC 8785 canonical bytes for the fixed policy aggregate. The aggregate always emits `modules` before `schema`; its only module emits `kind`, `path`, and `policy`, whose only field is `additional_tracker_types`. Array order, including an empty `[]`, is retained.

Before ordinary JSON marshaling, the method rejects a zero or invalid aggregate and validates UTF-8 for every aggregate string. It then canonicalizes only that successful typed JSON with the pinned `github.com/gowebpki/jcs v1.0.1` dependency. It exposes no raw JSON, number, map, or custom-marshaler boundary.

Canonical output does not normalize Unicode: valid literal U+FFFD, CJK, emoji, and composed/decomposed forms remain distinct. HTML characters are not HTML-escaped in the final canonical bytes. Returned bytes are independent of the private aggregate state and of later callers' mutations.

## Not included

The loader does not discover implicit modules, return catchall module content, resolve schemas, validate typed module fields beyond the policy contract, produce hashes, or expose a CLI. Schema resolution, hashing, further module typed contracts, and CLI integration remain deferred. `filesystem.Load` remains unchanged: this contract consumes the manifest and policy YAML supplied by its caller and performs no filesystem access.
