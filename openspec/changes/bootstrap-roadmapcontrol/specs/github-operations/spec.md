# GitHub Operations Specification

## Purpose

Define GitHub authority, permissions, workflow isolation, and GitHub Free limitations.

## Requirements

### Requirement: Installation confirmation and roles

The system MUST require human confirmation for every repository and its initial GitHub permissions. It MUST support owner, maintainer, and developer roles with owner-only permissions as the default. External users MAY read, fork, open Issues, and comment, but upstream implementation MUST require explicit authorization and assignment.

#### Scenario: Initial setup confirmation

- GIVEN a repository selected for installation
- WHEN installation reaches permission setup
- THEN the system MUST present the repository and requested permissions for human confirmation before applying them.

#### Scenario: External Issue author

- GIVEN an external user opens or comments on an Issue
- WHEN they attempt upstream implementation without assignment
- THEN the system MUST deny task authority.

### Requirement: Least-privilege GitHub access

The system MUST use `GITHUB_TOKEN` for short-lived repository-scoped code, Issue, and pull-request operations. Projects v2 synchronization MUST require an optional private GitHub App because `GITHUB_TOKEN` cannot access that surface; that App MUST NOT receive code-write permission.

#### Scenario: Projects disabled

- GIVEN Projects v2 synchronization is not enabled
- WHEN normal Issue and pull-request operations run
- THEN the system MUST NOT require a Projects App.

#### Scenario: Projects enabled

- GIVEN an owner enables Projects v2 synchronization
- WHEN the integration is configured
- THEN a private App is required and its granted permissions exclude code write.

### Requirement: Workflow trust separation

The system MUST keep privileged control workflows from executing pull-request code or exposing secrets. Verification workflows MUST remain unprivileged. All GitHub Actions dependencies MUST be pinned by full commit SHA and constrained by a minimal allowlist.

#### Scenario: Pull-request verification

- GIVEN untrusted pull-request code triggers verification
- WHEN the workflow runs
- THEN it runs without privileged secrets or control-workflow authority.

#### Scenario: Unpinned dependency

- GIVEN a workflow dependency is not pinned by full commit SHA or is outside the allowlist
- WHEN workflow policy is evaluated
- THEN the system MUST reject the workflow configuration.

### Requirement: Break-glass bypass handling

The system MUST treat a direct owner bypass as break-glass behavior: detect it, freeze the affected flow, compare authoritative state, and require reconciliation. It MUST state that GitHub Free cannot technically prevent a deliberate repository-owner bypass.

#### Scenario: Owner bypass detected

- GIVEN an owner changes a controlled artifact outside the normal flow
- WHEN the bypass is detected
- THEN the affected flow is frozen pending state comparison and reconciliation.

#### Scenario: GitHub Free boundary disclosure

- GIVEN an operator asks whether GitHub Free prevents owner bypass
- WHEN the system reports its security boundary
- THEN it MUST state that prevention is not technically guaranteed.

### Requirement: GitHub-native hosting boundary

The first release MUST NOT operate a publisher-hosted backend or SaaS, collect telemetry, store publisher-managed tokens, or transfer RoadmapControl data outside GitHub.

#### Scenario: First-release operation

- GIVEN a repository uses the first release
- WHEN RoadmapControl performs supported automation or CLI operations
- THEN it MUST use only GitHub-native automation and local tooling without publisher-hosted services, telemetry, publisher token storage, or data transfer outside GitHub.
