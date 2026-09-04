# Technical Design: Bootstrap RoadmapControl

## Decision summary

RoadmapControl will be a single Go module with a hexagonal core, repository-backed YAML as planning authority, Git/GitHub as durable operational storage, and SQLite only as a rebuildable cache. Privileged GitHub Actions will run a pinned RoadmapControl release from a trusted ref, never pull-request code. Unprivileged workflows will verify pull-request code without secrets. Every privileged mutation is a resumable, idempotent state machine whose durable events are appended to a signed audit branch.

The first release deliberately provides cooperative controls rather than claiming hostile-owner isolation. A repository owner on GitHub Free can bypass controls, and a local machine owner can bypass worktree locks. RoadmapControl detects evidence divergence, freezes the smallest provably isolated scope, and requires explicit reconciliation; it cannot make those platform boundaries disappear.

## How to use this design

The detailed design is organized by review concern. Read this index first, then follow the module links for the decisions and evidence relevant to a specification requirement.

| Module | Review focus |
| --- | --- |
| [Architecture and canonical storage](design/architecture-and-storage.md) | Goals, trust boundaries, package layout, canonical YAML, schemas, and migrations |
| [Domain model and execution control](design/domain-and-execution.md) | Aggregates, identifiers, state, dependencies, path policy, refs, leases, and local grants |
| [GitHub operations and delivery](design/github-and-delivery.md) | CLI/adapters, workflow trust separation, credentials, PR topology, verification, and line accounting |
| [Synchronization and audit](design/synchronization-audit-and-resilience.md) | Projection sagas, drift handling, signed audit history, and key lifecycle |
| [Operational resilience and lifecycle](design/operational-resilience-and-lifecycle.md) | Quota, handoff, crash/offline behavior, installation, updates, uninstall, reconstruction, and failures |
| [Verification, rollout, and traceability](design/verification-rollout-and-traceability.md) | Strict TDD, dogfooding phases, rejected alternatives, specification mapping, and planning guardrails |

## Architectural decision index

| Concern | Decision | Detail |
| --- | --- | --- |
| Authority | `.roadmap/` is canonical; Git and GitHub hold durable operational evidence; SQLite is disposable. | [Architecture and canonical storage](design/architecture-and-storage.md) |
| Structure | One Go module uses a hexagonal core with deterministic domain packages and external adapters. | [Architecture and canonical storage](design/architecture-and-storage.md) |
| Execution | Admission, leases, grants, path isolation, and recovery fail closed on missing or stale evidence. | [Domain model and execution control](design/domain-and-execution.md) |
| GitHub trust | Privileged trusted-code workflows and unprivileged pull-request verification remain separate. | [GitHub operations and delivery](design/github-and-delivery.md) |
| Delivery | Tracker/task PR topology and exact object-bound line accounting enforce the normal 400-line gate. | [GitHub operations and delivery](design/github-and-delivery.md) |
| Consistency | Non-atomic GitHub operations use durable sagas, idempotency markers, read-back, and scoped freezes. | [Synchronization and audit](design/synchronization-audit-and-resilience.md) |
| Audit | Signed append-only events on a Git branch provide reconstructible authority with explicit platform limits. | [Synchronization and audit](design/synchronization-audit-and-resilience.md) |
| Continuity | Quota loss, crashes, handoffs, offline work, lifecycle changes, and cache loss preserve the last proven safe state. | [Operational resilience and lifecycle](design/operational-resilience-and-lifecycle.md) |
| Verification | Implementation starts with the Go test harness, then follows strict red-green-refactor in cohesive work units. | [Verification, rollout, and traceability](design/verification-rollout-and-traceability.md) |

## Specification traceability

The complete requirement-to-design mapping is preserved in [Specification traceability](design/verification-rollout-and-traceability.md#specification-traceability). Each mapping points to the detailed sections in these modules; section names remain unchanged from the approved design.

## Scope guardrail

This modularization changes navigation only. Requirements, architectural decisions, trust claims, failure behavior, delivery limits, rollout sequencing, and rejected alternatives remain those of the approved design modules.
