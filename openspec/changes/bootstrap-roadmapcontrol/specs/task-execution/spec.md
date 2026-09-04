# Task Execution Specification

## Purpose

Constrain execution to explicitly approved, isolated, and recoverable task workspaces.

## Requirements

### Requirement: Fail-closed task admission

The system MUST admit a task only when it is roadmap-approved, explicitly queued and eligible, dependency-ready, assigned to an authorized identity, and defines outcome criteria and repository-relative allowed paths. Unknown required authority or evidence MUST fail closed globally unless isolation to a smaller affected scope is proven.

#### Scenario: Eligible task starts

- GIVEN a queued, dependency-ready approved task with an authorized assignee, outcome criteria, and allowed paths
- WHEN admission is requested
- THEN the system MAY admit that task.

#### Scenario: Missing admission evidence

- GIVEN a task with missing authorization, criteria, path scope, queue eligibility, or required evidence
- WHEN admission is requested
- THEN the system MUST deny admission and preserve the task as not started.

### Requirement: Path and work-in-progress isolation

The system MUST default to one active task, MUST reject overlapping active-task path scopes, and MUST let protected governance and control surfaces override any task path glob.

#### Scenario: Overlapping scope

- GIVEN an active task whose allowed paths overlap a second task
- WHEN the second task is admitted
- THEN the system MUST deny admission.

#### Scenario: Protected surface requested

- GIVEN a task path glob that includes a protected control surface
- WHEN access is evaluated
- THEN the protected-surface policy MUST take precedence over the glob.

### Requirement: Lease, branch, worktree, and grant binding

GitHub Actions MUST issue a signed remote lease and create `roadmapcontrol/task/<id>` before `roadmapctl` creates an exclusive sibling worktree. The system MUST permit one active agent per task. An agent grant MUST bind an authenticated GitHub identity to one task, one worktree, allowed paths, and allowed operations, and MUST terminate on task closure or revocation.

#### Scenario: Authorized workspace creation

- GIVEN an admitted task and a valid signed lease
- WHEN workspace setup is requested
- THEN the task branch is created before an exclusive sibling worktree is created and the grant is bound to that workspace.

#### Scenario: Invalid, revoked, or closed grant

- GIVEN an expired, revoked, mismatched, or closed-task grant
- WHEN an agent attempts an operation
- THEN the system MUST reject it and MUST NOT extend authority.

### Requirement: Cooperative local exclusivity

The system MUST describe physical worktree exclusivity as a cooperative safety control and MUST NOT claim protection against a malicious local machine owner with direct filesystem or process control.

#### Scenario: Local-control boundary disclosure

- GIVEN documentation or an operator query about worktree isolation
- WHEN the boundary is presented
- THEN it states that malicious host-owner resistance is outside this control.

### Requirement: Disconnected execution and recovery

The system MUST permit offline edits only under an existing lease. While disconnected, it MUST prohibit new tasks, handoffs, scope changes, completion transitions, pull requests, and synchronization until reconnection validation succeeds. A crash MUST require owner recovery; inactivity MAY mark a session stale but MUST NOT silently transfer authority.

#### Scenario: Offline edit under lease

- GIVEN an active valid lease and lost connectivity
- WHEN the agent edits within its existing allowed paths
- THEN the edit is permitted locally but control-plane transitions remain unavailable.

#### Scenario: Stale or crashed agent

- GIVEN a stale or crashed agent session
- WHEN another identity requests its task
- THEN the system MUST retain the original authority until owner-led recovery completes.
