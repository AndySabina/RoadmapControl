# Apply Progress: Bootstrap RoadmapControl

## PR02A-1 — parser-only split

- **Delivery path:** approved chained split; this is PR02A-1 (parser-only) and remains at or below the 400 changed-line budget.
- **Completed scope:** retained the typed, filesystem-independent manifest parser and its tests in `internal/domain/roadmap/`, plus its YAML dependency in `go.mod`/`go.sum`; added the concise parser-boundary documentation in `docs/roadmap-format.md`.
- **Deferred scope:** filesystem loading remains PR02A-2. Before removing the active filesystem files, their SHA-256 values were confirmed byte-for-byte equal to the preserved archive:
  - `internal/adapters/filesystem/manifest.go`: `26eeb27962d533183ab94a2e86e0f6a011966932db73874804a01ce23c1a3647`
  - `internal/adapters/filesystem/manifest_test.go`: `4d7c399bea5970c41c7e940274c48ca6483bd59174dce71a3427a942d0989d4c`
  - archive locator: `/home/andyf/Projects/RoadmapControl-local-archive/pr02a-before-split-20260905/` (`SHA256SUMS`).
- **No success claim for the removed loader:** PR02 remains incomplete. This slice has no filesystem API, CLI, hashing, schema resolution, module-content validation, or module-schema semantics.

## Verification

| Command | Result |
| --- | --- |
| `GOTOOLCHAIN=go1.27.1 go test -count=1 ./internal/domain/roadmap/...` | PASS after the Judgment Day correction |
| `GOTOOLCHAIN=go1.27.1 go test -count=1 ./...` | PASS after the Judgment Day correction |
| `gofmt -w internal/domain/roadmap/manifest.go internal/domain/roadmap/manifest_test.go && git diff --check` | PASS after the Judgment Day correction |
| `GOTOOLCHAIN=go1.27.1 go mod tidy` followed by SHA-256 comparison of `go.mod` and `go.sum` | PASS; no tidy mutation |
| Changed-line ledger vs `4f3e3f69bb67c352aea1f6063da0bd345c7f0682` (including docs and untracked files) | 335 additions + 0 deletions = **335** |

Runtime evidence: N/A — PR02A-1 provides a parser package only; no CLI or filesystem boundary is included.

## TDD Cycle Evidence

| Task | Test file | Layer | Safety net | RED | GREEN | TRIANGULATE | REFACTOR |
| --- | --- | --- | --- | --- | --- | --- | --- |
| PR02A-1 parser partition | `internal/domain/roadmap/manifest_test.go` | Unit | Focused and full suite passed before partition | Historical uncommitted RED observed for undefined `ParseManifest`/`Load`; no committed chronology is claimed, and it is not evidence for the later URI correction | Existing parser tests pass | Existing table cases cover valid input and unsafe YAML/path boundaries | No behavior refactor; filesystem code removed only after archive verification |

## Judgment Day correction round1

- JD-A-001/JD-B-001 reject malformed opaque URI escapes and characters without fetching; HTTPS, URN, and hostless file URIs remain valid.
- JD-A-002 rejects NUL module paths; JD-A-003/JD-B-002 reject non-string mapping keys before YAML decoding can collapse semantic duplicates.
- **RED (observed):** `GOTOOLCHAIN=go1.27.1 go test -count=1 ./internal/domain/roadmap/...` failed for the two malformed URIs, NUL path, and duplicate-key cases before the production fix.
- **GREEN/TRIANGULATE (observed):** the same focused command passed after the fix, retaining HTTPS, URN, hostless file, relative, and invalid URI boundaries.

## Task status and next work

The persisted PR-02 task remains unchecked because its complete loader/schema/CLI scope is not complete. No task checkbox was changed. PR02A-2 must restore the filesystem files from the recorded archive and independently verify its loader behavior.

## Rollback

To remove only this correction, revert `internal/domain/roadmap/manifest.go`, `internal/domain/roadmap/manifest_test.go`, and this progress file. Do not delete or alter the PR02A-2 archive; it is the deferred filesystem rollback/recovery source.

## Status consumed

```yaml
schemaName: spec-driven
changeName: bootstrap-roadmapcontrol
artifactStore: openspec
applyState: ready
actionContext:
  mode: repo-local
  workspaceRoot: /home/andyf/Projects/RoadmapControl-pr02a
  allowedEditRoots:
    - /home/andyf/Projects/RoadmapControl-pr02a/internal/domain/roadmap
    - /home/andyf/Projects/RoadmapControl-pr02a/internal/adapters/filesystem
    - /home/andyf/Projects/RoadmapControl-pr02a
    - /home/andyf/Projects/RoadmapControl-pr02a/docs
    - /home/andyf/Projects/RoadmapControl-pr02a/openspec/changes/bootstrap-roadmapcontrol
  warnings:
    - PR02 remains incomplete; filesystem delivery is deferred to PR02A-2.
```
