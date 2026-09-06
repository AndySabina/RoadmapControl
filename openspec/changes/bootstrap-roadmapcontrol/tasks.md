# Implementation Tasks: Bootstrap RoadmapControl

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 6,000–8,500 across 22 review slices |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR-00 governance prerequisite → PR 1 → PR 2–7 read-only core → PR 8–12 install/audit/admission → PR 13–17 task and tracker delivery → PR 18–20 sync/adapters → PR 21–22 release/dogfood |
| Delivery strategy | ask-on-risk |
| Chain strategy | feature-branch-chain |

Decision needed before apply: No for approved PR02B1/PR02B2 issue scope; `feature-branch-chain` is already selected. Future scope expansion needs approval.
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

**Delivery gate:** every child slice is capped at **≤400 additions + deletions**, including tests, documentation, configuration, and fixtures; do not reduce scope by code-golfing or separating behavior from its tests/docs. Each PR body must name its predecessor, successor, current position, and a dependency diagram with the current slice marked `📍`.

**Strict-TDD evidence rule:** after PR 1 establishes `go test ./...`, every behavior slice records (1) a committed local RED test and observed failure before GREEN; the failed intermediate head is evidence only and must never itself be pushed or published; its RED commit may appear in a passed candidate's ancestry only when the passed candidate HEAD is the ref pushed or published; (2) minimal GREEN implementation and focused passing command, (3) a TRIANGULATE case/property/boundary test and cumulative passing command, and (4) a REFACTOR pass with behavior unchanged. The completed candidate then receives focused and cumulative verification followed by Judgment Day final verification before its delivery commit, push, or draft PR publication. Each slice also records a runtime/integration scenario and result, or `N/A` with why no runtime boundary exists.

**Bootstrap boundary:** PR-00 and PRs 1–21 use ordinary repository/Gentle AI controls. Do not represent bootstrap history as RoadmapControl-governed. Dogfooding begins only in PR 22 after one complete validated control loop and explicit installation of an exact released pin.

## PR map

```text
Ordinary bootstrap review
PR-00 → PR-01 → PR-02A-1 → PR-02A-2 → PR-02B1 → PR-02B2 → PR-02C → PR-03 → PR-04 → PR-05 → PR-06 → PR-07
                         ↓
PR-08 → PR-09 → PR-10 → PR-11 → PR-12 → PR-13 → PR-14 → PR-15
                                                     ↓
PR-16 → PR-17 → PR-18 → PR-19 → PR-20 → PR-21
                                                     ↓
PR-22 (dogfood transition; new governed work starts after merge)
```

## Work units

### PR-00 — Establish public repository governance before implementation
- **Status:** completed prerequisite. **Dependencies:** none. **Tracker boundary:** documentation/configuration only; ordinary bootstrap controls.
- **Allowed edit surfaces:** `LICENSE`, `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, `.github/ISSUE_TEMPLATE/**`, `.github/PULL_REQUEST_TEMPLATE.md`, `docs/governance/README.md`, and this prerequisite entry.
- **Tasks:** publish the Apache-2.0 license, honest bootstrap project status and boundaries, contribution/security/conduct policies, bounded issue forms, a pull-request gate template, and a governance-documentation chain tracker.
- **Acceptance/evidence:** every artifact is read back; issue-form YAML parses with available standard tooling; `git diff --check` passes; exact line counts and cohesive review slices are recorded, with each slice at **≤400 additions + deletions**. No product implementation or GitHub mutation occurs.
- **Rollback:** revert only the listed public governance documents; no RoadmapControl runtime authority or GitHub configuration is active.

### PR-01 — Bootstrap the Go module, test runner, and minimal unprivileged CI
- **Status:** completed in PR #30; issue #28 closed. **Estimate:** 180–280 lines. **Dependencies:** none. **Tracker boundary:** bootstrap foundation; ordinary controls only.
- **Allowed edit surfaces:** `go.mod`, `go.sum`, `cmd/roadmapctl/`, `internal/**`, `.github/workflows/verify-pr.yml`, `README.md`, `.gitignore`.
- **Tasks:**
  1. Create module `github.com/AndySabina/RoadmapControl`, pin the supported Go toolchain, and add a minimal `roadmapctl` composition root with no product command.
  2. Establish `go test ./...` with a smoke test and a pinned, least-privilege `pull_request` workflow that runs it; document local bootstrap prerequisites.
  3. Add the first workflow-policy test fixture so later workflow changes have a test home; the production behavior scope remains empty.
- **Acceptance/evidence:** `go test ./...` passes locally and in the unprivileged workflow; workflow uses no secrets, no `pull_request_target`, and all action references are full SHAs. Record the first runnable command/result and CI run URL/ID. Runtime scenario: invoke `go run ./cmd/roadmapctl --help`.
- **Rollback:** revert only the listed bootstrap files; no RoadmapControl authority, refs, secrets, or repository settings exist.

### PR-02 — Validate canonical roadmap modules and schema pins in read-only mode
- **Status:** incomplete. The original PR-02 scope is refined into PR02A-1, PR02A-2, PR02B1, PR02B2, and PR02C below; PR-03 through PR-22 are not started and cumulative checks remain unchecked. **Dependencies:** PR-01. **Tracker boundary:** read-only core.
- **Allowed edit surfaces:** `internal/domain/roadmap/**`, `internal/adapters/filesystem/**`, `internal/ports/**`, `schemas/roadmap/v1/**`, `testdata/roadmap/**`, `docs/roadmap-format.md`. Previously authorized ancillary surfaces: minimal `cmd/roadmapctl/**` wiring, `go.mod`, `go.sum`, and this change's `apply-progress.md`; each remaining slice must narrow these surfaces before implementation.
- **Tasks:** implement explicit module-manifest loading, hardened YAML decoding, canonical JSON hashing, offline relative schema resolution, and read-only `roadmapctl validate` output.
- **Acceptance/evidence:** RED malformed/unknown/missing module tests; GREEN valid modular roadmap; TRIANGULATE duplicate-key, non-UTF-8, alias-bound, and schema-pin cases; REFACTOR package boundaries. `go test ./internal/domain/roadmap/... ./internal/adapters/filesystem/...` and cumulative suite pass. Runtime: validate a valid and invalid fixture without writes.
- **Rollback:** remove loader/schema/docs/fixtures only; no `.roadmap/` installation or GitHub mutation.

#### PR-02 delivery reconciliation

- **PR02A-1:** completed in PR #34 (issue #32 closed): typed manifest parser and safety boundaries, **335 lines**, commit `e9a81d2`. Its RED evidence was reconstructed at `1ecb062`, not claimed as original authorship.
- **PR02A-2:** completed in PR #35 (issue #33 closed): contained filesystem loading, **380 lines**, commits `b80a9dc` with correction `d84c76f`. Its committed RED is `154a8a6`.
- **Integration:** PR #35 merged into #34 at `5413f91`; PR #34 then merged into tracker at `1e4c0c9`. The resulting `02602d3` tree is the exact Judgment Day-approved combined child. Each standalone child is Judgment Day-approved, not an ordinary RDD receipt. The 709-line tracker aggregate is two previously reviewed slices, not a new ≤400-line single PR.
- **Superseded candidate:** preserve failed candidate `4191234` as evidence; it was not merged. Tracker PR #31 remains draft and unmerged, issue #27 remains open, and `main` remains untouched at `4acf0f5`.
- **Next:** original PR-02 remains unchecked. Approved docs [issue #38](https://github.com/AndySabina/RoadmapControl/issues/38) splits the remaining PR02B work into todo B1, then B2.
- [x] **B1 [approved and assigned issue #39](https://github.com/AndySabina/RoadmapControl/issues/39):** `internal/domain/roadmap/module.go`, `internal/domain/roadmap/module_test.go`, and `docs/roadmap-format.md`; **~180–300 lines**. Implementation is locally verified; independent verification and Judgment Day remain pending. It is not merged or published.
- **B2 [approved and assigned issue #40](https://github.com/AndySabina/RoadmapControl/issues/40):** `internal/domain/roadmap/canonical.go`, `internal/domain/roadmap/canonical_test.go`, shared format docs, and conditional `go.mod`/`go.sum`; **~180–350 lines**.
- **Both:** strict TDD, focused Go 1.27.1 tests, verification then Judgment Day, and rollback with their tests/docs/dependency changes. PR02C remains pending for generic typed loading, offline schema resolution, hashing, and CLI integration; neither slice may substitute arbitrary opaque modules for typed contracts. No `.roadmap/` installation is authorized.

### PR-03 — Enforce hierarchy, immutable identifiers, states, and derived progress
- **Estimate:** 300–390 lines. **Dependencies:** PR-02. **Tracker boundary:** read-only core.
- **Allowed edit surfaces:** `internal/domain/roadmap/**`, `testdata/roadmap/**`, `docs/roadmap-format.md`.
- **Tasks:** add typed R/P/S/T models, registry/high-water validation, Issue-association rules, exact state machine, blocked context, terminal immutability, parent overrides, and equal-weight progress.
- **Acceptance/evidence:** RED table tests per invariant; GREEN minimal hierarchy and transitions; TRIANGULATE retired-ID, optional phase/subphase, zero included tasks, and unjustified override cases; REFACTOR deterministic diagnostics. Focused and cumulative tests pass. Runtime: `roadmapctl validate` prints aggregate/override diagnostics for fixtures.
- **Rollback:** revert the domain-model increment while retaining PR-02 loader.

### PR-04 — Validate dependencies, explicit queue, promotion roles, and protected paths
- **Estimate:** 320–400 lines. **Dependencies:** PR-03. **Tracker boundary:** read-only core.
- **Allowed edit surfaces:** `internal/domain/roadmap/**`, `internal/policy/**`, `internal/ports/**`, `testdata/**`, `docs/roadmap-format.md`.
- **Tasks:** add stable DAG validation/inherited dependencies, queue scanning with rejection reasons and no prioritization, owner-default promotion policy, protected baseline, and fail-closed allowed-path parsing/intersection.
- **Acceptance/evidence:** RED cycle, unresolved dependency, maintainer-default denial, overlap, and protected-path tests; GREEN eligible queue selection; TRIANGULATE cross-tracker/inherited edges and wildcard automata symmetry/unknown cases; REFACTOR deterministic traversal. Focused and cumulative tests pass. Runtime: read-only `plan next` fixture reports first eligible task and skipped reasons.
- **Rollback:** remove queue/policy changes without changing canonical fixture history.

### PR-05 — Verify signed append-only audit evidence and reconstruct read-only state
- **Estimate:** 320–400 lines. **Dependencies:** PR-02. **Tracker boundary:** read-only core.
- **Allowed edit surfaces:** `internal/domain/audit/**`, `internal/app/**`, `internal/adapters/git/**`, `internal/adapters/cache/**`, `internal/ports/**`, `testdata/audit/**`, `docs/audit.md`.
- **Tasks:** implement canonical audit event verification, Ed25519 signatures/key validity/rotation continuity, hash/parent chain checks, disposable SQLite cache boundary, and read-only reconstruction/freeze reporting.
- **Acceptance/evidence:** RED invalid signature/hash/sequence tests; GREEN valid chain reconstruction; TRIANGULATE rotation, deleted cache, and divergent Git evidence cases; REFACTOR port seams. Focused and cumulative tests pass. Runtime: reconstruct from a temporary bare Git fixture with no cache.
- **Rollback:** revert audit/reconstruction code and docs; no audit branch or real key is created.

### PR-06 — Inspect repository compatibility and workflow-policy boundaries
- **Estimate:** 300–390 lines. **Dependencies:** PR-01, PR-04, PR-05. **Tracker boundary:** read-only core.
- **Allowed edit surfaces:** `internal/app/**`, `internal/policy/**`, `internal/adapters/git/**`, `internal/adapters/github/**`, `internal/ports/**`, `testdata/policy/**`, `docs/installation.md`.
- **Tasks:** implement read-only install inventory and versioned conflict matrix for refs, rules, workflow ownership, managed markers, permissions, Apps, and existing roadmap data; add workflow parser/policy checks.
- **Acceptance/evidence:** RED conflict, inaccessible evidence, mutable action, broad permission, secret-in-PR, and `pull_request_target` tests; GREEN compatible report; TRIANGULATE stricter-compatible rule and check-name collision; REFACTOR normalized fact model. Focused/cumulative tests pass. Runtime: `install inspect` on fixtures performs no writes and reports conflicts.
- **Rollback:** remove inspector/policy docs only.

### PR-07 — Add CLI read-only composition and documented security boundaries
- **Estimate:** 220–320 lines. **Dependencies:** PR-02–06. **Tracker boundary:** read-only core complete.
- **Allowed edit surfaces:** `cmd/roadmapctl/**`, `internal/app/**`, `docs/**`, `README.md`, `.github/workflows/verify-pr.yml`.
- **Tasks:** expose `validate`, `plan next`, `audit verify`, `reconstruct`, and `install inspect`; document GitHub Free/local-host limits, no SaaS/telemetry/token storage, supported platforms, and Actions quota caveats.
- **Acceptance/evidence:** RED CLI argument/error tests; GREEN read-only command wiring; TRIANGULATE unsupported environment and unavailable GitHub evidence; REFACTOR help/error consistency. Cumulative suite and `go vet ./...` pass. Runtime: execute each command against fixture repositories; verify no mutation.
- **Rollback:** remove command/docs layer while preserving independently testable core.

### PR-08 — Prepare non-destructive installation plans and managed blocks
- **Estimate:** 340–400 lines. **Dependencies:** PR-06, PR-07. **Tracker boundary:** installation rehearsal.
- **Allowed edit surfaces:** `internal/domain/lifecycle/**`, `internal/app/**`, `internal/adapters/filesystem/**`, `internal/adapters/github/**`, `internal/ports/**`, `testdata/install/**`, `docs/installation.md`.
- **Tasks:** create confirmed-install intent, exact permission display, deterministic initialization plan, sentinel/digest-managed Pi/OpenCode blocks, namespace collision checks, and token-less authentication abstraction; produce a PR plan but do not merge/apply it.
- **Acceptance/evidence:** RED unconfirmed/conflicting/modified-block tests; GREEN compatible plan; TRIANGULATE optional Projects disabled and unsupported platform cases; REFACTOR plan serialization. Focused/cumulative tests pass. Runtime: dedicated fake repository plan produces only proposed paths and never calls a mutator.
- **Rollback:** remove install-planning behavior; no remote configuration has changed.

### PR-09 — Add installation key bootstrap, pinned artifacts, updates, and conservative uninstall planning
- **Estimate:** 340–400 lines. **Dependencies:** PR-05, PR-08. **Tracker boundary:** installation rehearsal.
- **Allowed edit surfaces:** `internal/domain/lifecycle/**`, `internal/domain/audit/**`, `internal/app/**`, `internal/adapters/**`, `schemas/**`, `migrations/**`, `testdata/lifecycle/**`, `docs/installation.md`.
- **Tasks:** prepare public-key/install-pin material, encrypted-secret upload port, signed update-manifest verification, migration plan output, and uninstall inventory/blocking logic; private keys must remain transient and absent from logs/artifacts/cache.
- **Acceptance/evidence:** RED absent old key, invalid release signature, active-without-checkpoint, and secret-leak tests; GREEN valid plan paths; TRIANGULATE key rotation continuity and exact pin/migration cases; REFACTOR secret-redaction boundary. Focused/cumulative tests pass. Runtime: fake secret uploader receives encrypted material only; uninstall plan preserves history.
- **Rollback:** remove lifecycle plans/fixtures; no key secret or installation is applied.

### PR-10 — Exercise one non-destructive installation rehearsal in a dedicated test repository
- **Estimate:** 220–350 lines. **Dependencies:** PR-08, PR-09. **Tracker boundary:** installation rehearsal complete.
- **Allowed edit surfaces:** `internal/app/**`, `internal/adapters/git/**`, `internal/adapters/github/**`, `test/integration/**`, `docs/installation.md`, `.github/workflows/**`.
- **Tasks:** add hermetic integration coverage for inspect → confirmation → initialization PR preparation → read-back validation, plus explicit conflict stop; keep real GitHub execution opt-in and human-confirmed.
- **Acceptance/evidence:** RED failure-injection test for each mutation/read-back boundary; GREEN hermetic successful rehearsal; TRIANGULATE existing rule/ref/marker conflicts; REFACTOR shared test harness. Cumulative suite passes. Runtime: a dedicated test repository run is **manual, opt-in, and evidence-only**; otherwise record N/A because credentials/confirmation are unavailable.
- **Rollback:** remove rehearsal harness/docs; no production repository activation.

### PR-11 — Implement signed audit append saga and control-operation idempotency
- **Estimate:** 340–400 lines. **Dependencies:** PR-05, PR-09. **Tracker boundary:** auditable control foundation.
- **Allowed edit surfaces:** `internal/domain/audit/**`, `internal/app/**`, `internal/adapters/git/**`, `internal/ports/**`, `testdata/audit/**`, `docs/audit.md`.
- **Tasks:** add append-only audit ref operations, expected-parent CAS, operation IDs, idempotent replay, failure evidence, and frozen/unknown outcomes; retain read-only verification.
- **Acceptance/evidence:** RED non-fast-forward, duplicate operation, signing failure, and ambiguous write tests; GREEN one signed append; TRIANGULATE concurrent append/retry cases; REFACTOR saga boundaries. Focused/cumulative tests pass. Runtime: temporary bare remote demonstrates one winner and a frozen loser.
- **Rollback:** revert append saga without touching already-created audit evidence; any real evidence is preserved and handled by forward correction.

### PR-12 — Add quota receipts and fail-closed admission evidence prerequisites
- **Estimate:** 280–370 lines. **Dependencies:** PR-04, PR-11. **Tracker boundary:** admission foundation.
- **Allowed edit surfaces:** `internal/domain/execution/**`, `internal/domain/delivery/**`, `internal/app/**`, `internal/adapters/github/**`, `internal/ports/**`, `testdata/quota/**`, `docs/actions-quota.md`.
- **Tasks:** model authoritative quota receipts and healthy/warning/conservation/critical/exhausted/unknown states; deny new admission when evidence is missing, critical, exhausted, or incompatible.
- **Acceptance/evidence:** RED exact 70/85/95/100 and unknown-evidence tests; GREEN healthy admission prerequisite; TRIANGULATE unmetered response and stale/partial response; REFACTOR policy table. Focused/cumulative tests pass. Runtime: fake GitHub responses prove no admission call proceeds when quota is unknown/critical.
- **Rollback:** remove quota prerequisite, retaining audit core.

### PR-13 — Issue leases, task branches, and local sibling worktrees
- **Estimate:** 360–400 lines. **Dependencies:** PR-04, PR-11, PR-12. **Tracker boundary:** single low-risk task loop.
- **Allowed edit surfaces:** `internal/domain/execution/**`, `internal/app/**`, `internal/adapters/git/**`, `internal/adapters/filesystem/**`, `internal/ports/**`, `testdata/execution/**`, `docs/task-execution.md`.
- **Tasks:** implement lease envelope verification, branch-before-worktree ordering, Git worktree registry plus common-dir lock, branch read-back, and failure/freeze compensation rules.
- **Acceptance/evidence:** RED invalid lease, existing/mismatched ref, live lock, and audit-append-failure tests; GREEN valid branch/worktree setup; TRIANGULATE expired lease and ambiguous changed branch cases; REFACTOR Git port. Focused/cumulative tests pass. Runtime: temporary bare remote and sibling worktree prove ordering/exclusivity; document cooperative-host limitation.
- **Rollback:** remove worktree created for test scenarios and revert this unit; never delete an ambiguous/changed branch automatically.

### PR-14 — Bind grants, enforce task scope, and gate offline/reconnect behavior
- **Estimate:** 340–400 lines. **Dependencies:** PR-13. **Tracker boundary:** single low-risk task loop.
- **Allowed edit surfaces:** `internal/domain/execution/**`, `internal/app/**`, `internal/policy/**`, `internal/adapters/git/**`, `internal/ports/**`, `testdata/execution/**`, `docs/task-execution.md`.
- **Tasks:** bind actor immutable ID, task, worktree, paths, operations, session key, closure/revocation; validate Git-tree paths/symlinks; permit only leased offline edits and require full reconnect validation/owner recovery.
- **Acceptance/evidence:** RED grant mismatch/revocation/out-of-scope/offline-control-action tests; GREEN scoped operation; TRIANGULATE symlink escape, stale session, and crash recovery cases; REFACTOR grant verification. Focused/cumulative tests pass. Runtime: attempted outside-scope and offline PR action are rejected in a temp fixture.
- **Rollback:** revoke test grants/release locks, then revert grant code; no authority transfer is inferred.

### PR-15 — Create tracker/task PR topology and exact review-budget gate
- **Estimate:** 360–400 lines. **Dependencies:** PR-13, PR-14. **Tracker boundary:** single low-risk task loop.
- **Allowed edit surfaces:** `internal/domain/delivery/**`, `internal/app/**`, `internal/adapters/git/**`, `internal/adapters/github/**`, `internal/ports/**`, `testdata/delivery/**`, `docs/delivery.md`.
- **Tasks:** create configurable system tracker/task refs and draft tracker PR plans, require authenticated user intent and tracker target, calculate content-bound numstat with exact base/head, and enforce review/owner-exception evidence.
- **Acceptance/evidence:** RED wrong target/no intent/401/unknown diff tests; GREEN 400-line task PR plan; TRIANGULATE pure rename, binary, stale head, and exact owner exception cases; REFACTOR candidate identity. Focused/cumulative tests pass. Runtime: fake GitHub plus bare Git fixture shows branch topology and line result.
- **Rollback:** close unmerged test PRs and remove untouched test refs, then revert topology code; preserve audited exceptions/evidence.

### PR-16 — Run unprivileged task/tracker/release verification workflows
- **Estimate:** 320–400 lines. **Dependencies:** PR-01, PR-15. **Tracker boundary:** single low-risk task loop.
- **Allowed edit surfaces:** `.github/workflows/**`, `internal/domain/delivery/**`, `internal/app/**`, `internal/policy/**`, `testdata/workflows/**`, `docs/verification.md`.
- **Tasks:** add pinned reusable focused/full/build/version workflows and content-bound result model; explicit dispatch/poll continuation; separate no-checkout status publication from candidate execution.
- **Acceptance/evidence:** RED workflow-policy and failed/mismatched result tests; GREEN exact task-bound result; TRIANGULATE fork PR, dispatch rejection, and token-event suppression cases; REFACTOR workflow fixture helpers. Focused/cumulative tests pass. Runtime: workflow parser confirms unprivileged jobs have no secrets/write token and privileged jobs do not execute PR code.
- **Rollback:** disable/revert only these workflows and result model; do not weaken existing repository policies.

### PR-17 — Gate tracker promotion, release decisions, versioning, and forward rollback
- **Estimate:** 350–400 lines. **Dependencies:** PR-11, PR-15, PR-16. **Tracker boundary:** complete validated control loop.
- **Allowed edit surfaces:** `internal/domain/delivery/**`, `internal/app/**`, `internal/adapters/git/**`, `internal/adapters/github/**`, `internal/ports/**`, `testdata/release/**`, `docs/release.md`.
- **Tasks:** require current-head human approvals and merged children for tracker promotion; model logical development/production mapping, SemVer/CalVer/custom policy, exact release checks, authenticated production decision, release audit event, and forward-only correction.
- **Acceptance/evidence:** RED incomplete tracker/no production decision/failed check tests; GREEN authorized release candidate; TRIANGULATE same physical branch, stale approval, CalVer/custom, and correction-version cases; REFACTOR release aggregate. Focused/cumulative tests pass. Runtime: end-to-end hermetic tracker → task → checks → explicit release decision produces signed evidence.
- **Rollback:** do not rewrite any release ref; use a new correction PR/version for activated behavior, otherwise revert unmerged code.

### PR-18 — Synchronize Issue projections with outbox, read-back, drift freeze, and reconciliation
- **Estimate:** 360–400 lines. **Dependencies:** PR-04, PR-11, PR-14, PR-17. **Tracker boundary:** synchronization after stable task loop.
- **Allowed edit surfaces:** `internal/domain/synchronization/**`, `internal/app/**`, `internal/adapters/github/**`, `internal/ports/**`, `testdata/synchronization/**`, `docs/synchronization.md`.
- **Tasks:** emit merge-bound outbox events, update only Issue managed sections/version markers, read back idempotently, gate admission on verified synchronization, freeze isolated/unknown scope, and expose exactly restore/incorporate/postpone paths.
- **Acceptance/evidence:** RED ambiguous write/drift/missing outbox tests; GREEN verified Issue projection; TRIANGULATE retry-after-read-back, shared-scope freeze, and state-only diff rejection; REFACTOR saga state model. Focused/cumulative tests pass. Runtime: `httptest` fixture proves comments remain unchanged and dependent work stays disabled on uncertainty.
- **Rollback:** freeze affected synchronization and use explicit reconciliation; never blindly overwrite Issue history.

### PR-19 — Add Pi and OpenCode scoped adapters plus transactional handoff
- **Estimate:** 360–400 lines. **Dependencies:** PR-14, PR-18. **Tracker boundary:** synchronization and handoff.
- **Allowed edit surfaces:** `internal/adapters/agent/pi/**`, `internal/adapters/agent/opencode/**`, `internal/domain/execution/**`, `internal/app/**`, `internal/ports/**`, `testdata/agents/**`, `docs/agents.md`.
- **Tasks:** negotiate capabilities including optional Gentle AI references, expose grants only, journal/quiesce operations, checkpoint with session key, secret-scrub portable manifest/Engram pointer, sender revocation-before-receiver activation, and owner-led crash recovery.
- **Acceptance/evidence:** RED unavailable capability, unsettled operation, secret-hygiene, crash, and receiver-without-verification tests; GREEN Pi→OpenCode handoff; TRIANGULATE revoked sender and offline reconnect mismatch; REFACTOR adapter/core separation. Focused/cumulative tests pass. Runtime: fake adapters verify no token propagation and one active authority.
- **Rollback:** revoke active test grant/retain checkpoint, then revert adapters; do not transfer authority on partial handoff.

### PR-20 — Add opt-in Projects v2 synchronization with least-privilege App validation
- **Estimate:** 300–390 lines. **Dependencies:** PR-18. **Tracker boundary:** optional integration; remains disabled by default.
- **Allowed edit surfaces:** `internal/adapters/projects/**`, `internal/domain/synchronization/**`, `internal/app/**`, `internal/ports/**`, `testdata/projects/**`, `docs/projects.md`.
- **Tasks:** implement owner-confirmed App configuration, permission-shape validation that rejects code write, Project projection/read-back, App-token isolation, and App-loss freeze behavior.
- **Acceptance/evidence:** RED missing App/code-write/ambiguous Project response tests; GREEN opt-in verified projection; TRIANGULATE App loss and Issue-only mode; REFACTOR optional-port seam. Focused/cumulative tests pass. Runtime: contract fake validates token boundary; live App test is manual opt-in with explicit owner confirmation or N/A.
- **Rollback:** disable Projects mapping and freeze/reconcile its scope; preserve Issue and roadmap evidence.

### PR-21 — Validate the complete control loop and lifecycle stop/reconstruct/uninstall exercises
- **Estimate:** 360–400 lines. **Dependencies:** PR-10–20. **Tracker boundary:** release-ready validation.
- **Allowed edit surfaces:** `test/e2e/**`, `test/integration/**`, `testdata/**`, `docs/**`, `.github/workflows/**`, narrowly scoped test-support files under `internal/**`.
- **Tasks:** compose a hermetic one-tracker/one-documentation-or-maintenance-task loop: install rehearsal, canonical approval, quota/sync/audit evidence, lease/worktree/grant, task PR/check/review simulation, tracker merge, release decision, forward correction; add stop/reconstruct/uninstall failure exercises and operator runbook.
- **Acceptance/evidence:** RED missing evidence at each control boundary; GREEN complete loop; TRIANGULATE drift, quota exhaustion, crash, and uninstall-without-checkpoint; REFACTOR shared fixtures without reducing coverage. Full `go test ./...`, build, workflow policy, and exact line-budget checks pass. Runtime: execute hermetic loop and capture object IDs/evidence; live GitHub is manual opt-in only.
- **Rollback:** retain evidence, freeze at last proven state, and use forward correction or conservative uninstall plan; never claim the repository is dogfooding yet.

### PR-22 — Install the released pin into RoadmapControl and begin dogfooding prospectively
- **Estimate:** 240–380 lines. **Dependencies:** PR-21 plus explicit human approval of the validated release and chosen chain strategy. **Tracker boundary:** dogfood transition, not retroactive bootstrap governance.
- **Allowed edit surfaces:** `.roadmap/**`, `schemas/**`, `migrations/**`, `.github/workflows/**`, managed Pi/OpenCode instruction blocks, installation documentation, and generated initialization-plan paths proven by PR-08–10.
- **Tasks:** execute the reviewed non-destructive initialization plan against RoadmapControl itself, confirm exact repository/permissions once, install the exact released pin/public audit key/pinned workflows/managed blocks, run bootstrap read-back and one post-install control-loop verification, then create the first **new** approved roadmap item for subsequent governed work.
- **Acceptance/evidence:** RED/preflight evidence records any conflict and blocks mutation; GREEN installed-file/ref/permission/key read-back; TRIANGULATE modified managed block, unavailable audit key, and namespace collision stop paths; REFACTOR only after post-install verification remains green. Record human confirmation, installation audit object, full tests/build/workflow results, and controlled runtime evidence. Existing bootstrap commits remain ordinary history.
- **Rollback:** before activation, close/revert the initialization PR only. After activation, use a reviewed forward-only correction or conservative uninstall; preserve `.roadmap/`, audit branch, public keys, Issues/PRs, and system history.

## Completion checks

- [ ] Each PR has a verified additions+deletions count of 400 or less before review; an indivisible overage stops for an audited owner `size:exception`.
- [ ] Each behavioral PR contains RED → GREEN → TRIANGULATE → REFACTOR evidence, focused/cumulative tests, runtime evidence or N/A, and its own docs.
- [ ] Every requirement in `specs/roadmap-domain`, `task-execution`, `delivery-release`, `github-operations`, `synchronization`, `agent-integration`, and `audit-lifecycle` is covered by at least one listed slice.
- [ ] PR-21 proves one complete validated control loop before PR-22 enables prospective dogfooding.
- [ ] No implementation, commit, push, issue, PR, GitHub permission/App change, or repository-rule mutation is performed by this planning phase.
