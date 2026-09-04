# Bootstrap RoadmapControl as a GitHub-Native Roadmap Control System

## Intent

Create the first public release of **RoadmapControl**, an Apache-2.0 Go project hosted at <https://github.com/AndySabina/RoadmapControl>. RoadmapControl will make team-approved work explicit in a repository-versioned roadmap and keep that authority traceable across GitHub Issues, Projects, branches, worktrees, pull requests, reviews, releases, and supported agent sessions.

The project will be bootstrapped before it can dogfood its own controls. Bootstrap work must therefore establish the repository, governance, schema, CLI, GitHub Actions, and adapter foundations needed for later RoadmapControl development to run under RoadmapControl itself.

## Problem and Current-State Gap

A GitHub issue backlog records possible work, but it does not reliably identify what the team has approved for execution now. Planning data and execution artifacts can consequently drift apart: an Issue may disagree with a Project field, a branch may mix tasks, an agent may operate outside its assignment, or a pull request may no longer represent the approved plan.

RoadmapControl must close that gap by:

- maintaining a canonical record of current approved work;
- preserving hierarchy and end-to-end traceability without replacing GitHub discussion history;
- identifying the next team-approved executable item through an explicit queue;
- preventing accidental task mixing through assignment, path, lease, worktree, and pull-request boundaries;
- reconciling drift without silent overwrite or history loss; and
- failing closed when required authority or evidence is unknown.

## Product Outcome

Maintainers can answer, from repository and GitHub evidence, what work is approved, what may start next, who or which agent is authorized to perform it, which files it may affect, how it will be reviewed and released, and whether every derived GitHub representation is synchronized to the same roadmap version.

Contributors and agents receive a bounded task workspace rather than broad repository authority. Owners retain explicit control over promotion, permissions, exceptional overrides, reconciliation, release, recovery, and uninstall.

## Foundational Product Decisions

| Area | Decision |
| --- | --- |
| Canonical plan | Modular YAML under `.roadmap/`, governed by versioned JSON Schema and migrations. It contains current approved work only; GitHub Issues may exist outside it. |
| Operational ordering | Human teams decide what executes. An explicit queue is the sole operational ordering mechanism; RoadmapControl does not calculate a separate automated priority. |
| Promotion authority | Only an owner or maintainer may promote work into the approved roadmap. The default installation policy is owner-only. |
| GitHub role | Issues own comments, discussion, and evidence history. Projects are derived views. Roadmap planning fields remain canonical. |
| Hosting model | GitHub-native automation plus a local CLI and agent adapters. There is no RoadmapControl SaaS/backend, telemetry, or transfer outside GitHub. |
| Safety posture | Unknown evidence fails closed globally unless isolation to a smaller affected scope is proven. Reconciliation is explicit and history-preserving. |
| Reviewability | A normal pull request is limited to 400 authored additions plus deletions, including tests, documentation, and configuration. Only an audited owner exception may admit indivisible work. |

## First Release Scope

The first release is the smallest complete governance loop that can safely approve, execute, synchronize, review, and release repository work. Implementation may be delivered as multiple reviewable work units; each normal pull request remains within the 400-line budget unless an audited exception is approved.

### 1. Canonical roadmap and domain model

- Store the roadmap as modular YAML under `.roadmap/` and validate it with versioned JSON Schema.
- Provide versioned, pull-request-driven migrations.
- Support typed top-level trackers: `feature`, `bug`, `maintenance`, `security`, and `documentation`, with owner-configurable additional types.
- Support the bounded hierarchy `tracker -> phase -> subphase -> task`; phases and subphases are optional and used only when they add planning value.
- Require a GitHub Issue association for trackers and tasks, but not for phases or subphases.
- Assign immutable, globally unique R/P/S/T identifiers.
- Permit acyclic dependencies at any hierarchy level, including inherited and cross-tracker dependencies within the same repository.
- Use equal task weight for progress. Exclude cancelled and superseded work from progress, and derive aggregate parent states by default.
- Define exactly these states: `not_started`, `in_progress`, `blocked`, `done`, `cancelled`, and `superseded`.
- Require justification for exceptional parent-state overrides. Treat terminal states as immutable.
- When work becomes blocked, preserve its prior state together with the blocking reason and resume condition.

### 2. Task admission and bounded execution

A task may start only when it is roadmap-approved, explicitly queued and eligible, dependency-ready, assigned to an authorized identity, and has both outcome criteria and repository-relative allowed paths.

- Protected governance and control surfaces override any task path glob.
- Default work in progress is one active task.
- Active task path scopes must not overlap.
- GitHub Actions issues a signed remote lease and creates `roadmapcontrol/task/<id>` before `roadmapctl` creates an exclusive sibling worktree.
- Permit one active agent per task.
- Agent grants bind an authenticated GitHub identity to one task, one worktree, allowed paths, and allowed operations. Grants terminate on task closure or revocation.
- Physical worktree exclusivity is a cooperative safety control. It does not claim to defeat a malicious machine owner with direct filesystem or process control.

### 3. Pull-request, tracker, and release flow

- Create one tracker branch and one draft tracker pull request for each top-level tracker.
- Create normal task pull requests through RoadmapControl from authorized user intent and target them at the tracker branch.
- Require authorized human review for every normal pull request.
- Allow tracker approval only after all child pull requests have been reviewed and merged.
- Merge completed trackers into a configurable logical development branch.
- Treat production as a separate cumulative promotion and release; development and production may be configured to the same physical branch.
- Use SemVer by default while supporting CalVer and custom versioning policies.
- Perform rollback through a new forward-only pull request and version rather than rewriting published history.
- Reserve a configurable `roadmapcontrol/` namespace for system-managed branches. Existing project branches remain under project ownership.

### 4. Verification and Actions quota policy

- Run focused checks at the task boundary.
- Run the full suite and verify tracker outcomes at the tracker boundary.
- Run the full suite, build checks, and version checks at release.
- Use GitHub-hosted Actions runners for GitHub-native verification and document their dependency and limits honestly.
- Classify Actions quota usage as: healthy at 0–69%, warning at 70–84%, conservation at 85–94%, critical at 95–99%, and exhausted when no quota remains.
- Block admission of new tasks at critical usage and fail closed when exhausted.
- Design for an initial operating expectation of approximately 30 pull requests per month without presenting that estimate as a guaranteed quota or capacity.

### 5. GitHub permissions and workflow isolation

- Use `GITHUB_TOKEN` for short-lived, repository-scoped code, Issue, and pull-request operations.
- Require an optional private GitHub App only when Projects v2 synchronization is enabled, because `GITHUB_TOKEN` cannot access that surface. The App receives no code-write permission.
- Require a human to confirm every repository and the initial GitHub permissions.
- Ensure privileged control workflows never execute pull-request code or expose secrets. Keep verification workflows unprivileged.
- Pin GitHub Actions dependencies by full commit SHA and enforce a minimal allowlist.
- Retain owner, maintainer, and developer roles, with owner-only permissions as the default.
- Permit external users to read, fork, open Issues, and comment; upstream implementation still requires explicit authorization and assignment.
- Treat direct owner bypass as break-glass behavior: detect it, freeze the affected flow, compare state, and require reconciliation.
- Do not claim that GitHub Free can technically prevent a repository owner from deliberately bypassing configured controls.

### 6. Versioned synchronization and drift reconciliation

- Emit a versioned outbox entry for every roadmap merge.
- Project that version to the linked Issues and optional Project, read the result back, and record the synchronized version before enabling dependent work.
- Freeze affected synchronization when drift is detected.
- Offer three explicit responses: restore GitHub from the roadmap, incorporate the external change through a roadmap pull request, or postpone reconciliation while the affected scope remains frozen.
- Permit exact, state-only bot pull requests for mechanical `in_progress` and `done` transitions and auto-merge them only under prior authorization.
- Require an authorized, human-approved pull request for every other roadmap change.

### 7. Agent adapters, handoff, and disconnected behavior

- Provide first-class Pi and OpenCode adapters.
- Support Gentle AI optionally through capability negotiation and references to its rules rather than copying those rules into RoadmapControl.
- Provide an initial OpenAI model profile while keeping the core model-agnostic.
- Make cross-agent handoff transactional: quiesce work, settle or cancel active operations, validate scope and secret hygiene, create and push a signed checkpoint commit to the task branch, persist Pi context to Engram plus a portable manifest, close the sender session and release its lock atomically, then require receiver verification.
- Require owner recovery after a crash. Inactivity only marks a session stale and does not silently transfer authority.
- Permit offline edits only under an existing lease. While disconnected, prohibit new tasks, handoffs, scope changes, completion transitions, pull requests, and synchronization until reconnection validation succeeds.

### 8. Auditability, durability, installation, and updates

- Maintain a signed, append-only audit branch with compact metadata and evidence protected by a hash chain.
- Use a per-repository Ed25519 key: keep the public key in the repository and the private key in GitHub Actions secrets. The signing key grants no GitHub API authority.
- Support explicit signing-key rotation.
- Make durable state reconstructible from Git and GitHub. Use SQLite only as a disposable local cache.
- Classify existing repository rules before installation and block installation on conflicts until the owner reconciles them.
- Initialize non-destructively through a pull request.
- Uninstall conservatively: preserve history and block removal while active work lacks a checkpoint.
- Provide signed, on-demand updates with exact repository version pins and migration through pull requests; do not perform silent automatic updates.
- Provide agent-guided installation with one unavoidable human GitHub permission confirmation.
- Install managed, non-destructive instruction blocks for Pi and OpenCode.
- Reuse `gh` authentication, Git credential managers, or an OAuth device flow without storing user tokens.

### 9. Initial compatibility envelope

The first release supports:

- public and private repositories on github.com Free or higher;
- Linux, WSL2, and macOS;
- one roadmap per repository; and
- dependencies contained within one repository.

## Future Capabilities

The following are candidates for later releases and are not commitments of this bootstrap change:

- GitHub Enterprise Server support;
- native Windows support;
- cross-repository dependencies and coordinated releases;
- multiple independent roadmaps within one repository;
- additional first-class agent adapters and model-provider profiles;
- broader GitHub Projects capabilities beyond the optional private Projects v2 synchronization App; and
- policy refinements based on measured GitHub Actions usage after the expected initial workload is observed.

Future work must preserve the first-release authority, auditability, non-destructive migration, and fail-closed principles rather than bypass them.

## Explicit Non-Goals

This change will not:

- turn the entire GitHub Issue backlog into approved work;
- replace Issue comments, discussion, or evidence history with roadmap YAML;
- infer team priority or reorder work outside the explicit queue;
- autonomously promote, assign, approve, merge, release, reconcile, or expand scope without the required human authority;
- provide a hosted RoadmapControl service, telemetry pipeline, or off-GitHub data transfer;
- grant the optional Projects App code-write permission;
- manage or rename existing project-owned branches outside the configured system namespace;
- support GitHub Enterprise Server, native Windows, cross-repository dependencies, or multiple roadmaps per repository in the first release;
- prove worktree exclusivity against a malicious owner of the local machine;
- guarantee that GitHub Free prevents a deliberate repository-owner bypass;
- use audit signing keys as GitHub API credentials;
- treat SQLite as authoritative durable state;
- rewrite history for rollback, silently overwrite drift, or destroy audit history during uninstall; or
- copy Gentle AI governance rules into this repository instead of negotiating capabilities and referencing the authoritative source.

## Affected Areas

| Area | Expected impact |
| --- | --- |
| Repository layout | Introduces `.roadmap/`, schemas, migrations, Go CLI/module foundations, adapter configuration, and managed instruction blocks. |
| GitHub Actions | Adds privileged control workflows, unprivileged verification workflows, pinned dependencies, lease issuance, synchronization, audit signing, and quota gates. |
| GitHub configuration | Adds role/policy configuration, system branch namespace, development/production branch mapping, optional Projects App integration, and rules compatibility checks. |
| Planning workflow | Makes roadmap promotion and queue placement explicit owner/maintainer actions rather than treating all Issues as executable. |
| Contributor workflow | Requires authorization, assignment, lease, allowed paths, system-created task branch/worktree, bounded pull requests, and human review. |
| Maintainer operations | Adds reconciliation, break-glass handling, key rotation, migration, release promotion, crash recovery, update, and uninstall procedures. |
| Agent workflow | Adds scoped grants, one-agent-per-task enforcement, transactional handoff, checkpointing, stale/crash recovery, and constrained offline behavior. |
| Support and documentation | Requires honest compatibility, GitHub permissions, Actions quota, security-boundary, failure-mode, and recovery documentation. |

## Risks and Mitigations

| Risk | Mitigation in scope |
| --- | --- |
| Bootstrap complexity causes RoadmapControl to depend on controls that do not exist yet. | Establish capabilities through non-destructive bootstrap pull requests, then transition subsequent development to dogfooding only after the control loop is operational. |
| Roadmap, Issues, and Projects diverge under asynchronous GitHub operations. | Use versioned outbox entries, write/read-back verification, synchronized-version recording, scoped freezes, and explicit reconciliation choices. |
| Privileged Actions expose secrets to untrusted pull-request code. | Separate privileged control workflows from unprivileged verification; privileged workflows never run PR code; pin actions by SHA and minimize the allowlist. |
| GitHub permission limits or Projects v2 constraints surprise users. | Require explicit initial permission confirmation, use short-lived `GITHUB_TOKEN` where possible, make the no-code-write Projects App optional, and document platform limits. |
| Actions quota exhaustion blocks normal operation. | Surface tiered quota states, conserve before exhaustion, block new tasks at critical usage, and fail closed when exhausted. |
| Path scopes or worktree locks are mistaken for hostile-host isolation. | Document them as cooperative controls and avoid claiming protection from a malicious machine owner. |
| Repository owners bypass controls, especially on GitHub Free. | Treat bypass as detectable break-glass behavior with freeze and reconciliation; explicitly document that deliberate owner bypass cannot be prevented on GitHub Free. |
| Signing-key compromise or rotation breaks audit trust. | Separate signing from API authority, keep the private key in Actions secrets, publish the public key, hash-chain records, and provide explicit rotation. |
| Oversized changes reduce review quality. | Limit normal PRs to 400 authored changed lines, keep tests/docs/config in the count, split by coherent work unit, and require an audited owner exception for indivisible work. |
| Fail-closed behavior creates avoidable global outages. | Scope freezes only when isolation is proven; otherwise preserve the safer global failure while giving owners explicit recovery and reconciliation paths. |
| Conservative uninstall or crash recovery surprises operators. | Document checkpoints and prerequisites, preserve history, and block unsafe removal or authority transfer rather than guessing. |

## Rollout and Rollback

### Rollout

1. Bootstrap the repository and foundational test harness through reviewable pull requests.
2. Classify existing repository rules, permissions, branches, and workflows.
3. Stop and require owner reconciliation for conflicts.
4. Present the exact GitHub permissions for human confirmation.
5. Install schemas, roadmap configuration, pinned workflows, audit key material, CLI/adapters, and managed instruction blocks through a non-destructive initialization pull request.
6. Validate one complete tracker-to-task-to-release flow before declaring RoadmapControl ready to dogfood its own roadmap.
7. Pin the installed RoadmapControl version; offer later signed updates and migrations only on demand through pull requests.

### Rollback

- Before activation, close or revert the initialization pull request without altering existing project-owned branches or history.
- After activation, correct released behavior through a new forward-only pull request and version.
- For uninstall, checkpoint active work first, remove managed automation and instruction blocks conservatively through reviewed changes, and retain roadmap and audit history.
- Block uninstall when active work has not been checkpointed or when removal would lose authoritative evidence.

## Success Criteria

The bootstrap change is successful when all of the following are demonstrably true within the initial compatibility envelope:

- A maintainer can install RoadmapControl non-destructively after one explicit GitHub permission confirmation, and installation blocks rather than overwrites conflicting repository rules.
- Modular roadmap YAML validates against a pinned schema version, and schema migration occurs through a pull request.
- The roadmap distinguishes approved current work from unrelated Issues and exposes one explicit operational queue without generating an independent priority.
- Domain validation enforces identifiers, hierarchy, exact statuses, dependency acyclicity and inheritance, terminal-state immutability, blocked-state context, and equal-weight progress rules.
- An authorized, dependency-ready task with outcome criteria and allowed paths can receive a signed lease, system task branch, and exclusive sibling worktree; an ineligible, overlapping, unauthorized, or out-of-scope task fails closed.
- A task pull request targets its tracker branch, stays within the normal 400-line budget or carries an audited owner exception, and cannot satisfy the normal flow without authorized human review.
- A tracker cannot be approved before all child pull requests are reviewed and merged, and a completed tracker can be promoted from the logical development branch into a separately governed release.
- Focused task checks, tracker-boundary full verification, and release build/version verification run in unprivileged GitHub-hosted workflows, while privileged workflows do not execute pull-request code.
- Every roadmap merge produces a versioned outbox record; dependent work remains disabled until GitHub projection is read back and the synchronized version is recorded.
- Detected drift freezes the affected synchronization scope and preserves all three explicit reconciliation paths without silently overwriting either side.
- Audit state is signed, append-only, hash-chained, and reconstructible from Git/GitHub without relying on SQLite.
- Pi and OpenCode can execute under scoped grants and complete a verifiable transactional handoff; crash, stale-session, and offline restrictions preserve authority boundaries.
- Quota state is visible and enforces the defined healthy, warning, conservation, critical, and exhausted behavior.
- Documentation states the real security boundary: local exclusivity cannot withstand a malicious machine owner, and GitHub Free cannot prevent deliberate owner bypass.
- No publisher-operated backend, telemetry, token storage, or transfer outside GitHub is introduced.

## Proposal Boundary

This proposal authorizes planning of the first RoadmapControl release. It does not authorize implementation, commits, pushes, releases, GitHub permission changes, App installation, or repository-rule changes. Those actions require subsequent specification, design, task planning, and explicit execution under the established SDD and repository controls.
