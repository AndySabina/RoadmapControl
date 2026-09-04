# Agent Integration Specification

## Purpose

Provide bounded Pi and OpenCode execution, capability-based Gentle AI interoperability, and recoverable cross-agent handoffs.

## Requirements

### Requirement: Supported adapter capabilities

The system MUST provide first-class Pi and OpenCode adapters and an initial OpenAI model profile while keeping the core model-agnostic. It MUST support Gentle AI optionally through capability negotiation and references to its authoritative rules, and MUST NOT copy Gentle AI governance rules into RoadmapControl.

#### Scenario: Pi or OpenCode task session

- GIVEN a valid task grant for Pi or OpenCode
- WHEN the adapter starts a session
- THEN it receives only the grant's task, worktree, path, and operation authority.

#### Scenario: Gentle AI unavailable or unsupported

- GIVEN Gentle AI capability negotiation is unavailable or lacks a required capability
- WHEN integration is requested
- THEN the system MUST not assume the capability or duplicate external rules, and MUST preserve the task boundary.

### Requirement: Transactional handoff

The system MUST make cross-agent handoff transactional: quiesce work, settle or cancel active operations, validate scope and secret hygiene, create and push a signed checkpoint commit to the task branch, persist Pi context to Engram plus a portable manifest, close the sender session and release its lock atomically, and require receiver verification.

#### Scenario: Successful handoff

- GIVEN a sender with an active task grant
- WHEN a handoff is approved
- THEN all handoff steps complete and the receiver verifies the checkpoint before it obtains usable authority.

#### Scenario: Handoff validation failure

- GIVEN active operations cannot settle, scope validation fails, secret hygiene fails, or checkpoint push fails
- WHEN handoff is attempted
- THEN the system MUST not transfer authority and MUST retain or recover the sender lock safely.

### Requirement: Owner-led recovery

The system MUST require owner recovery after a crashed session and MUST NOT use inactivity alone to transfer task authority.

#### Scenario: Receiver after sender crash

- GIVEN the sender crashes before completing handoff
- WHEN another agent requests the task
- THEN the system MUST require owner recovery and checkpoint evidence before a new grant is issued.

### Requirement: Human-authorized agent control actions

An agent MUST NOT autonomously promote work, assign work, approve or merge pull requests, promote or release production, reconcile drift, or expand scope. Each such action MUST require the applicable authenticated human authority. The only exception is an exact mechanical `in_progress` or `done` state-only pull request, which an agent MAY create or auto-merge only under prior authorization.

#### Scenario: Agent requests a governed control action

- GIVEN an agent has a valid task grant but no authenticated human authorization for a governed control action
- WHEN it attempts to promote, assign, approve, merge, release, reconcile, or expand scope
- THEN the system MUST deny the action.

#### Scenario: Authorized mechanical state-only pull request

- GIVEN prior authorization for an exact mechanical `in_progress` or `done` transition
- WHEN an agent creates or auto-merges the state-only pull request
- THEN the system MAY permit that action without separate autonomous authority for other control actions.
