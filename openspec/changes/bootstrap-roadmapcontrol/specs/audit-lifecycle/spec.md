# Audit and Lifecycle Specification

## Purpose

Preserve durable governance evidence through installation, operation, update, recovery, and conservative removal.

## Requirements

### Requirement: Signed append-only audit trail

The system MUST maintain a signed append-only audit branch containing compact metadata and evidence protected by a hash chain. It MUST use a per-repository Ed25519 key with the public key in the repository and private key in GitHub Actions secrets. The signing key MUST NOT grant GitHub API authority, and key rotation MUST be explicit.

#### Scenario: Audited control event

- GIVEN a governed control event is recorded
- WHEN audit evidence is appended
- THEN it is hash-chain protected and signed with the repository audit key.

#### Scenario: Signing-key rotation

- GIVEN an owner authorizes key rotation
- WHEN rotation occurs
- THEN the new key is explicitly recorded while prior audit history remains preserved and verifiable.

### Requirement: Reconstructible durable state

The system MUST make durable state reconstructible from Git and GitHub. SQLite MAY be used only as a disposable local cache and MUST NOT be authoritative.

#### Scenario: Local cache loss

- GIVEN the local SQLite cache is deleted or unavailable
- WHEN state recovery is requested
- THEN the system reconstructs durable state from Git and GitHub evidence rather than treating the cache as authority.

### Requirement: Non-destructive installation

The system MUST classify existing repository rules before installation and MUST block conflicts until owner reconciliation. It MUST initialize through a non-destructive pull request, install managed non-destructive Pi and OpenCode instruction blocks, and require one unavoidable human GitHub permission confirmation. It MUST reuse `gh` authentication, Git credential managers, or OAuth device flow without storing user tokens.

#### Scenario: Conflicting repository rules

- GIVEN installation discovers conflicting existing rules
- WHEN initialization is requested
- THEN the system MUST block installation without overwriting rules until the owner reconciles the conflict.

#### Scenario: Successful initialization

- GIVEN compatible rules and confirmed GitHub permissions
- WHEN initialization is prepared
- THEN it is delivered as a non-destructive pull request with managed instruction blocks and no stored user token.

### Requirement: Signed on-demand updates

The system MUST provide signed on-demand updates with exact repository version pins and migration through pull requests. It MUST NOT silently update an installation.

#### Scenario: Update available

- GIVEN a signed newer RoadmapControl version is available
- WHEN an owner requests update
- THEN the system proposes the exact pinned update and any migration through a pull request.

### Requirement: Conservative uninstall and compatibility envelope

The system MUST preserve roadmap and audit history during uninstall and MUST block removal while active work lacks a checkpoint or removal would lose authoritative evidence. The first release MUST support one roadmap per public or private github.com Free-or-higher repository on Linux, WSL2, and macOS, with dependencies within that repository. It MUST NOT claim first-release support for GitHub Enterprise Server, native Windows, cross-repository dependencies, or multiple independent roadmaps.

#### Scenario: Unsafe uninstall

- GIVEN active work lacks a checkpoint
- WHEN uninstall is requested
- THEN the system MUST block removal and preserve existing history.

#### Scenario: Unsupported environment

- GIVEN native Windows, GitHub Enterprise Server, cross-repository dependencies, or multiple roadmaps are requested
- WHEN capability is evaluated
- THEN the system MUST report the request as outside the first-release compatibility envelope.
