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

## PR02A-2 — filesystem manifest loader GREEN

- **Delivery path / PR boundary:** approved chained PR02A-2 GREEN work unit, following PR02A-1; its boundary is read-only manifest acquisition and document safety only.
- **RED baseline:** committed RED `154a8a6b866939a14999c488d875c6b7bc0f7a6e`; `GOTOOLCHAIN=go1.27.1 go test -count=1 ./internal/adapters/filesystem/...` failed with `undefined: Load` at the four test call sites before production code was written.
- **Completed scope:** added `filesystem.Load(root)`, returning the typed `roadmap.Manifest` parsed by domain `ParseManifest`; it reads only declared nested/multiple modules and applies domain `ValidateYAML` to module documents. It checks root containment, symlink roots/files/intermediate components, regular files, and the 1 MiB bound before reads; it rejects unlisted case-insensitive YAML.
- **Checkbox reconciliation:** no persisted task checkbox changed. The sole PR-02 task intentionally remains unchecked because schema resolution, hashing, typed module contracts, and CLI integration are deferred; this completed PR02A-2 slice has no separate task line.

### Verification

| Command | Result |
| --- | --- |
| `GOTOOLCHAIN=go1.27.1 go test -count=1 ./internal/adapters/filesystem/...` | PASS after JD-A-001/JD-B-001 unlisted-YAML-symlink correction |
| `GOTOOLCHAIN=go1.27.1 go test -count=1 ./...` | PASS |
| `gofmt -w internal/adapters/filesystem/manifest.go && git diff --check` | PASS |
| `GOTOOLCHAIN=go1.27.1 go mod tidy` with pre/post `go.mod` and `go.sum` SHA-256 comparison | PASS; no dependency mutation |

Runtime evidence: N/A — this is a read-only library boundary exercised with `t.TempDir()`; no CLI is in scope. Rollback: remove `internal/adapters/filesystem/manifest.go` and the matching documentation/progress section without affecting the parser-only PR02A-1 work.

### TDD Cycle Evidence

| Task | RED | GREEN | TRIANGULATE | REFACTOR |
| --- | --- | --- | --- | --- |
| PR02A-2 loader | Committed RED above; corrections RED observed for unlisted `nested/roadmap.yaml` and unlisted `extra.yaml` symlink | Focused loader suite passes after exempting only the root manifest and rejecting the unlisted YAML symlink without following it | Existing nested/multiple, symlink-intermediate, unsafe-YAML, size, and unlisted-case tests pass | `gofmt`; behavior unchanged |

### Safety caveat and remaining work

Checks occur before `os.ReadFile`, avoiding intentional FIFO reads. This non-mutating cooperative boundary cannot guarantee safety against a hostile concurrent TOCTOU replacement; it makes no universal-race claim. Deferred: schema resolution, canonical hashing, typed module validation/contracts, and CLI.

**Pending GREEN commit:** parent retains Git authority; no commit, push, reset, or other worktree/authority mutation was performed here.

### Status and remaining persisted work

```yaml
changeName: bootstrap-roadmapcontrol
artifactStore: openspec
applyState: ready
actionContext:
  mode: repo-local
  workspaceRoot: /home/andyf/Projects/RoadmapControl-pr02a2
  allowedEditRoots: [internal/adapters/filesystem, docs, openspec/changes/bootstrap-roadmapcontrol]
  warnings: [PR-02 remains incomplete beyond PR02A-2]
```

Exact unchecked persisted completion lines remain (global, untouched):
- [ ] Each PR has a verified additions+deletions count of 400 or less before review; an indivisible overage stops for an audited owner `size:exception`.
- [ ] Each behavioral PR contains RED → GREEN → TRIANGULATE → REFACTOR evidence, focused/cumulative tests, runtime evidence or N/A, and its own docs.
- [ ] Every requirement in `specs/roadmap-domain`, `task-execution`, `delivery-release`, `github-operations`, `synchronization`, `agent-integration`, and `audit-lifecycle` is covered by at least one listed slice.
- [ ] PR-21 proves one complete validated control loop before PR-22 enables prospective dogfooding.
- [ ] No implementation, commit, push, issue, PR, GitHub permission/App change, or repository-rule mutation is performed by this planning phase.

## Key Learnings

- Validate each component below the selected root, not only the final module pathname, to reject an intermediate symlink escape.
- A pre-read regular-file and size check reduces accidental unsafe reads but is not a hostile-race guarantee.
