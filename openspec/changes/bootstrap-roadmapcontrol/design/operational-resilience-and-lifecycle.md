# Operational Resilience and Lifecycle

## Quota accounting and Actions exhaustion

Quota evidence is a typed snapshot containing account/repository scope, source API, total allowance when available, consumed amount, remaining amount, reset period, retrieval time, and confidence. The policy computes `floor(100 * consumed / total)` only from compatible exact values:

| State | Evidence and behavior |
| --- | --- |
| healthy | 0–69%; normal operation |
| warning | 70–84%; warn and preserve all required checks |
| conservation | 85–94%; batch safe read-backs and reduce optional scheduled runs, never skip required checks |
| critical | 95–99%; deny new task admission; existing work may checkpoint but governed transitions still require their evidence |
| exhausted | Known remaining is zero or GitHub explicitly rejects for exhausted quota; fail closed for quota-backed operations |
| unknown | Billing API unavailable, permission denied, incompatible account data, stale/partial response, or denominator absent; fail closed for new admission and any operation requiring missing Actions evidence |

An authenticated GitHub response that explicitly identifies the applicable hosted-runner usage as unmetered is classified healthy; absence of a denominator alone is not such evidence. GitHub billing visibility varies by plan, account owner permissions, public-repository treatment, and API availability. RoadmapControl does not infer healthy status from an absent bill. Documentation may describe approximately 30 pull requests per month as an initial operating expectation, never as a quota or capacity guarantee. The CLI reports the source and missing fields. A fresh authoritative response is required at admission; cached quota data is informational only.

Actions exhaustion can also prevent the workflow that would record exhaustion. Therefore local and GitHub admission paths both require a recent successful quota/control receipt, and canonical state defaults to not admitting work when that receipt cannot be refreshed. Existing offline edits remain constrained by their prior valid lease, but no new task, handoff, completion, synchronization, tracker promotion, or release can substitute local evidence for unavailable required Actions. After quota returns, an owner starts reconciliation, which discovers pending branches/outbox/checks and resumes idempotently.

## Checkpoint, handoff, crash, and offline protocols

### Checkpoints and transactional handoff

Each session uses its lease-bound ephemeral SSH Ed25519 key to create Git-signed checkpoint commits. The public-key fingerprint is in the signed remote lease, so a privileged verifier can bind the checkpoint signature to the session without exposing the repository audit key. A checkpoint also includes a portable manifest digest with task/lease IDs, base and head objects, changed paths, operation journal state, verification evidence, adapter capabilities, and redacted context references.

```mermaid
sequenceDiagram
    participant S as Sender adapter
    participant C as roadmapctl
    participant G as Git/GitHub
    participant A as Audit control
    participant R as Receiver adapter

    S->>C: request handoff
    C->>S: quiesce; deny new mutations
    C->>C: settle/cancel operations; scope and secret checks
    C->>G: signed checkpoint commit and push
    G->>A: verify lease, signature, paths, head
    A->>A: append checkpoint_verified and revoke sender authority
    C->>S: atomically mark local session closed and release lock
    R->>A: request provisional receiver grant
    R->>G: verify checkpoint and portable manifest
    R->>A: attest exact checkpoint
    A->>A: append receiver_active grant
```

For Pi, the adapter persists context through the existing Engram capability and records only a digest/pointer plus a secret-scrubbed portable manifest. RoadmapControl operates no Engram backend and treats unavailable persistence as a failed Pi handoff. The portable manifest is always required so the receiver is not dependent on conversational memory alone.

The sender's local close marker and OS lock release occur in one local critical section. Remote authority changes are not falsely described as atomic with that filesystem operation. Safety comes from ordering: the sender lease is revoked before any receiver grant can become active, and the receiver is provisional/read-only until it verifies the checkpoint. A failure before remote revocation leaves the sender authoritative; a failure after revocation leaves no active agent and requires resume/recovery. It never leaves both active.

Scope validation compares Git trees, not the operation journal alone. Secret hygiene rejects known credential formats, private-key material, configured forbidden paths, and scanner uncertainty; it cannot prove arbitrary data harmless. Failed settling, scope validation, hygiene, signing, push, audit append, local close, or receiver verification prevents transfer and records/requires recovery at the last proven state.

### Crash and stale sessions

Inactivity may append a stale observation but does not revoke a lease, release authority, or activate another agent. After a crash, an owner-led recovery loads remote lease/audit state, inspects task branch and available worktree evidence, and obtains or creates a valid signed checkpoint. The owner then authorizes a recovery event that revokes the old grant. A new agent grant is unavailable until checkpoint evidence is verified. Ambiguous local state keeps the task frozen.

### Offline behavior

An adapter may continue local edits and local signed commits only while it holds an existing locally verifiable, unexpired lease and stays inside allowed paths. Offline mode prohibits remote branch creation, new task admission, handoff, scope/assignment changes, completion state, pull-request creation, synchronization, lease renewal, tracker promotion, and release.

On reconnect, the CLI verifies authenticated immutable GitHub identity, remote lease status, audit continuity, canonical roadmap version, task/tracker branch object IDs, path scope, offline commit signatures, operation journal digest, synchronization state, and quota/check requirements. Any mismatch freezes the task for owner recovery. Offline work is never pushed or promoted merely because connectivity returned.

## Rule conflict detection and managed control surfaces

Installation preflight is read-only. It inventories repository identity/default branch, existing refs in the configured namespace, Actions permissions and allowlist, workflows and triggers, rulesets/branch protection where visible, required checks, merge settings, environments, CODEOWNERS, installed Apps' declared permissions, existing `.roadmap/` data, and Pi/OpenCode managed markers. Secret values are never read.

The inventory is normalized into facts and evaluated by a versioned conflict matrix. Blocking conflicts include:

- namespace/ref or managed-marker ownership collisions;
- existing automation that can mutate `.roadmap/`, system refs, or audit history outside the control path;
- required workflow permissions unavailable or explicitly denied;
- required-check name collisions with different producers;
- mutable/unallowlisted actions in a would-be privileged path;
- a Projects configuration without the minimum no-code-write App permission shape;
- inability to inspect evidence required to establish safe compatibility; and
- repository rules that would force RoadmapControl to overwrite or weaken an existing owner policy.

Compatible stricter project rules are preserved and incorporated into the plan. Informational platform limitations, including owner bypass on GitHub Free, are disclosed rather than represented as enforceable rules. Preflight emits a deterministic report and proposed changes; it performs no mutation. The owner resolves conflicts explicitly, reruns preflight, confirms the exact repository and permission set once, and only then may initialization be proposed.

Managed instruction/configuration blocks carry stable begin/end sentinels, installed version, and content digest. Update/uninstall edits only a block whose ownership and digest can be proven. Unknown or manually changed blocks freeze lifecycle automation rather than being overwritten.

## Installation, updates, uninstall, and reconstruction

### Installation

1. `roadmapctl install inspect` checks the compatibility envelope and produces the rule inventory/conflict report.
2. Owner resolves conflicts outside RoadmapControl and reruns inspection.
3. CLI presents repository identity, exact GitHub permissions, optional Projects permissions, branch namespace, workflow pins, and signing-key handling for one unavoidable human confirmation.
4. CLI generates the audit key material, uploads only the encrypted private secret, and prepares a non-destructive initialization branch/PR containing schemas, initial roadmap/config, public key, exact release pins, pinned workflows, and managed Pi/OpenCode blocks.
5. Human review/merge activates controls. A privileged bootstrap validation reads back installed files, refs, permissions, and key operation before appending the installation event.

Before activation, closing/reverting the PR leaves existing project refs and policies untouched.

### Signed on-demand updates

`installation.yaml` pins RoadmapControl release version, source commit, binary digest, publisher-signature identity, schemas, migrations, workflow/action SHAs, and managed-block digests. `update inspect` fetches an owner-requested release manifest from GitHub Releases, verifies its signature against the pinned RoadmapControl release trust root, and prepares an exact-version pull request with any required schema migration. No scheduled job changes the pin, and availability notices do not authorize updates.

### Conservative uninstall

Uninstall first inventories active leases, grants, task/tracker branches, pending synchronization/outbox operations, and checkpoints. It blocks if active work lacks a verified checkpoint or if any authoritative evidence would be lost. The reviewed uninstall plan revokes authority, disables managed workflows, and removes only proven managed instruction/automation blocks. It preserves `.roadmap/` history, public audit keys, the audit branch, Issues/PRs, releases, and system branch history. Ambiguous ownership of a block means leave it and report manual follow-up, not delete it.

### Reconstruction

`roadmapctl reconstruct` starts without SQLite. It verifies installed pins and schemas, assembles canonical roadmap history, walks the audit branch from genesis validating hashes/signatures/rotations, inventories system refs and GitHub PR/Issue/release/check evidence, reads projection markers, and rebuilds derived task, lease, synchronization, quota, and lifecycle views. The resulting SQLite database stores source object IDs and is discarded whenever those IDs diverge. Missing, contradictory, or unverifiable durable evidence creates the corresponding freeze; cache contents never repair authority.

## Failure semantics

| Failure | Safe state | Recovery |
| --- | --- | --- |
| Invalid/unsupported roadmap or migration | Canonical candidate rejected; dependent admission disabled | Correct through reviewed roadmap/migration PR |
| Unknown identity, role, synchronization, audit, or quota evidence | Global freeze unless dependency analysis proves a smaller scope | Restore evidence, then explicit revalidation |
| Branch created but active lease audit append failed | No usable lease; branch inspected and frozen if non-empty/ambiguous | Idempotent saga reconciliation or owner recovery |
| GitHub/Projects write returned ambiguously | Target remains verifying/frozen; no blind repeat | Read back by idempotency marker, then continue or reconcile |
| `GITHUB_TOKEN` write emitted no recursive event | Current operation remains pending | Explicit exact-SHA dispatch or scheduled/manual reconciliation |
| Projects App unavailable | Project-enabled synchronization scope frozen; no code authority affected | Restore App permission/token and read back |
| Audit key unavailable or audit chain invalid | Privileged auditable mutations blocked | Restore key or perform explicit owner recovery; never bypass signing silently |
| PR line count unknown or above 400 | Normal review gate blocked | Make a cohesive smaller PR or obtain exact-candidate audited owner exception |
| Required Actions run cannot start/quota exhausted | Admission and quota-backed transitions blocked | Wait for quota/availability, then reconcile and rerun exact required checks |
| Owner bypass detected | Affected flow frozen and authoritative states compared | Choose explicit reconciliation; acknowledge GitHub Free cannot prevent bypass |
| Local worktree lock conflict | Workspace creation/adapter activation denied | Cooperative cleanup after identity/lease verification or owner recovery |
| Handoff fails before sender revocation | Sender retains authority, remains quiesced or resumes explicitly | Fix validation/push failure and retry with same operation state |
| Handoff fails after sender revocation | No active agent | Owner recovery and verified checkpoint before a new grant |
| Offline reconnect differs from remote authority | Task frozen; no push/control transition | Owner recovery and evidence comparison |
| Cache loss/corruption | Cache discarded | Reconstruct from Git/GitHub and signed audit evidence |

Compensation is used only when its safety is provable. Otherwise the saga stops at an explicit pending/frozen state. Retrying an operation reuses its operation ID; a new intent receives a new ID.
