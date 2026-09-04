# Verification, Rollout, and Traceability

## Strict TDD and verification strategy

The repository currently has no Go module or runnable test command. The first implementation work unit is a tooling prerequisite: create `go.mod`, pin the Go toolchain, and establish `go test ./...`. Creating the module does not implement product behavior. Immediately afterward, every behavioral unit follows strict red-green-refactor: add a failing test against a public package/use-case contract, capture the expected failure, implement the minimum behavior, and rerun focused plus cumulative tests.

Test layers:

| Layer | Coverage |
| --- | --- |
| Pure unit/table tests | State transitions, parent progress, ID allocation, queue scan, dependency closure/cycles, roles, quota boundaries |
| Property/fuzz tests | YAML normalization, duplicate/ambiguous input rejection, dependency DAG invariants, glob intersection symmetry, audit canonicalization |
| Schema/migration fixtures | Valid/invalid modular roadmaps, deterministic before/after digests, unsupported versions, no in-place writes |
| Git integration tests | Temporary bare remotes, ref races, non-fast-forward audit append, branch topology, worktree locks, signed checkpoint verification, exact numstat counting |
| GitHub adapter contract tests | `httptest` REST/GraphQL fixtures, pagination, actor IDs, permission denial, delayed consistency, ambiguous writes, App token boundaries |
| Saga/failure-injection tests | Crash after every lease/outbox/handoff step, idempotent replay, freeze scope, no double grant |
| Workflow policy tests | Parse workflow YAML; reject `pull_request_target`, mutable action refs, excess permissions, secret flow to PR jobs, event-recursion assumptions |
| Lifecycle tests | Rule inventory/conflict matrix, managed-block ownership, signed update pins, blocked uninstall, cache-free reconstruction |
| Adapter tests | Pi/OpenCode capability bounds, no token propagation, Engram/manifest handoff, crash/stale/offline restrictions |
| End-to-end test repository | One tracker/task/release flow on github.com with unprivileged checks and privileged control metadata only |

No test uses production signing/App secrets. Crypto tests generate ephemeral keys. Network tests default to fakes; live GitHub dogfood tests use a dedicated repository and explicit human authorization. Every work unit records focused test command/result, cumulative `go test ./...` result once available, runtime scenario or explicit N/A, and rollback boundary. Tests, documentation, and configuration stay with the behavior they verify and count toward the 400-line PR budget.

Because the complete release will exceed one normal review budget, implementation must be sliced by cohesive end-to-end work units rather than package/file type. The task plan must forecast authored additions plus deletions for each slice and use chained delivery or an audited size exception according to the orchestrator's delivery decision; code must not be compressed or tests/docs removed to fit.

## Phased rollout and dogfooding

1. **Bootstrap harness:** establish the Go module/test runner and workflow policy tests under ordinary repository review. No RoadmapControl authority is claimed yet.
2. **Read-only core:** ship schema assembly, domain validation, audit verification, rule inspection, and reconstruction in shadow/read-only mode. Compare reports manually; perform no governed GitHub writes.
3. **Non-destructive installation rehearsal:** run install preflight in a dedicated test repository, confirm exact permissions, and verify conflict blocking, key bootstrap, update pins, and conservative uninstall.
4. **Single low-risk task:** activate one tracker and one documentation/maintenance task with lease, sibling worktree, focused checks, task PR, checkpoint, and tracker integration. Keep manual owner observation at every privileged transition.
5. **Synchronization and handoff:** enable Issue projection/read-back and exercise Pi-to-OpenCode handoff plus crash/offline recovery. Projects remains disabled until Issue-only synchronization is stable.
6. **Optional Projects validation:** install the no-code-write private App in the test repository, verify drift/reconciliation and App loss behavior, then make it available as opt-in.
7. **Release loop:** verify full tracker checks, separate authenticated production decision, version policy, signed audit evidence, and forward-only correction path.
8. **Dogfood:** install an exact released pin into RoadmapControl itself through a reviewed PR. Only subsequent approved work runs under RoadmapControl; bootstrap history is not retroactively represented as governed.

Each phase has a stop/reconstruct/uninstall exercise before advancing. A failed phase remains in ordinary/manual governance or the last proven RoadmapControl state; rollout never relaxes controls to force dogfooding.

## Decisions and rejected alternatives

| Decision | Rationale | Rejected alternative |
| --- | --- | --- |
| Single Go module, internal hexagonal packages | One deployable CLI/workflow engine with deterministic reusable policy | Microservices/SaaS: violates hosting boundary and adds distributed authority |
| Explicit module manifest plus canonical JSON hash | Deterministic snapshot, no filesystem-order ambiguity | Implicit YAML discovery or raw-text hashing: unstable and easier to hide files in |
| Monotonic typed IDs with retired registry | Deterministic, reviewable, never reused | Random UUIDs: unique but poor human traceability; issue numbers: mutable/incomplete hierarchy identity |
| Queue scan without computed priority | Preserves team ordering while allowing ineligible entries to be explained | Priority scoring or automatic queue mutation: outside approved scope |
| Automata-based glob intersection, unknown fails closed | Proves overlap instead of relying on unsafe prefixes | Prefix-only checks: miss wildcard intersections; unrestricted regex: hard to decide safely |
| Audit events on a signed append-only Git branch | GitHub-native, reconstructible, no database authority | SQLite/event service as authority: violates durability/hosting constraints |
| Sagas plus read-back/idempotency markers | Honest handling of non-atomic GitHub APIs | Distributed transaction claim or blind retries: unavailable and unsafe |
| Separate privileged and unprivileged workflows | Prevents PR code from reaching control secrets | `pull_request_target` or privileged PR build: secret/code-execution risk |
| Explicit dispatch/reconciliation after token writes | Handles recursive event suppression deterministically | Depending on bot writes to retrigger workflows: GitHub does not guarantee it |
| Session key signs checkpoints; audit key stays in Actions | Verifiable handoff without exposing repository audit secret | Giving audit key to agents/PR code: collapses trust boundary |
| Exact Git numstat algorithm bound to object IDs | Reproducible 400-line gate and candidate identity | GitHub UI estimate or current worktree count: mutable/ambiguous |
| Quota `unknown` fails closed | Billing visibility is incomplete and required checks cannot be fabricated | Treat missing billing data as healthy: unsafe optimistic inference |
| Cooperative worktree/adapter controls | Matches actual local-host boundary | Claiming sandbox security against machine owner: false |
| Explicit owner break-glass reconciliation | Only honest GitHub Free posture | Claiming branch protection can stop a deliberate owner: false |
| Preserve audit/roadmap during uninstall | Keeps authoritative history and enables reconstruction | Destructive cleanup: violates durability and recovery requirements |

## Specification traceability

References below use document-local labels in specification order.

| Specification requirement | Design sections |
| --- | --- |
| Roadmap Domain RD-1 Versioned canonical roadmap | Canonical roadmap storage; Schemas and migrations; Reconstruction |
| RD-2 Pull-request-driven schema migration | Schemas and migrations; Signed on-demand updates |
| RD-3 Typed hierarchy and identifiers | Domain model; Identifiers |
| RD-4 Dependencies and executable queue | Dependencies and queue selection |
| RD-5 State, blocking, and progress | State machines |
| RD-6 Controlled roadmap promotion | Queue admission; CLI; Rule conflict/control workflow authority |
| RD-7 Planning source ownership | Canonical roadmap storage; Synchronization/outbox |
| Task Execution TE-1 Fail-closed admission | Dependencies and queue selection; Quota; Failure semantics |
| TE-2 Path and WIP isolation | Path policy; Local workspace and grant |
| TE-3 Lease, branch, worktree, grant binding | Remote lease protocol; Local workspace and grant |
| TE-4 Cooperative local exclusivity | Trust boundaries; Local workspace and grant |
| TE-5 Disconnected execution and recovery | Crash and stale sessions; Offline behavior |
| Delivery/Release DR-1 Tracker/task PR topology | Git ref model; Pull requests and verification |
| DR-2 Review budget and human review | Exact line accounting; Audit exceptions |
| DR-3 Tracker completion and promotion | Pull requests and verification |
| DR-4 Boundary verification | Actions topology; Strict TDD strategy |
| DR-5 Actions quota policy | Quota accounting and exhaustion |
| DR-6 Versioning and rollback | Pull requests/release; Signed updates |
| DR-7 Authorized production decision | Pull requests and verification; Release aggregate |
| GitHub Operations GH-1 Installation confirmation and roles | Rule conflict detection; Installation; CLI authentication |
| GH-2 Least-privilege access | Actions credentials; Optional Projects App |
| GH-3 Workflow trust separation | Trust boundaries; Actions topology; Workflow policy tests |
| GH-4 Break-glass bypass handling | Audit limitations; Failure semantics; Threat boundaries |
| GH-5 GitHub-native hosting boundary | System context; Adapters; Audit/reconstruction |
| Synchronization SY-1 Versioned projection/read-back | Synchronization/outbox/idempotency |
| SY-2 Drift freeze/reconciliation | Synchronization; Failure semantics |
| SY-3 Controlled mechanical state updates | Synchronization state-only PR validation |
| Agent Integration AI-1 Supported adapters | CLI and agent adapters |
| AI-2 Transactional handoff | Checkpoints and transactional handoff |
| AI-3 Owner-led recovery | Crash and stale sessions |
| AI-4 Human-authorized agent controls | CLI/adapters; control workflow authority; state-only exception |
| Audit/Lifecycle AL-1 Signed append-only audit | Audit branch and key lifecycle |
| AL-2 Reconstructible durable state | Reconstruction; disposable cache boundary |
| AL-3 Non-destructive installation | Rule conflict detection; Installation |
| AL-4 Signed on-demand updates | Signed on-demand updates |
| AL-5 Conservative uninstall/compatibility | Conservative uninstall; Goals/constraints; phased rollout |

## Implementation guardrails for task planning

- Start with the module/test bootstrap, then enforce strict TDD for every behavior.
- Slice implementation by demonstrable work unit with tests/docs, not by architectural layer.
- Keep normal pull requests at or below 400 authored additions plus deletions using the exact candidate-bound algorithm.
- Never place audit/App secrets in pull-request-accessible jobs, artifacts, logs, caches, manifests, or agent context.
- Never describe GitHub Free owner controls, Actions availability, API atomicity, billing visibility, or local worktree locks more strongly than the boundaries in this design.
- Any task plan that cannot map a work unit to a specification row, focused test, runtime evidence, and independent rollback boundary is incomplete.
