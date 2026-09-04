# Synchronization Specification

## Purpose

Project canonical roadmap versions to GitHub representations without silent loss of authority or history.

## Requirements

### Requirement: Versioned projection and read-back

The system MUST emit a versioned outbox entry for every roadmap merge. It MUST project that version to linked Issues and, when enabled, the optional Project; read the result back; and record the synchronized version before enabling dependent work.

#### Scenario: Successful projection

- GIVEN a merged roadmap version with linked Issues
- WHEN synchronization runs
- THEN an outbox entry is emitted, each enabled representation is read back, and the synchronized version is recorded before dependent work is enabled.

#### Scenario: Projection uncertainty

- GIVEN projection, read-back, or version recording fails or returns unknown evidence
- WHEN dependent work is evaluated
- THEN the system MUST keep the affected scope disabled and preserve the failure evidence.

### Requirement: Drift freeze and explicit reconciliation

When drift is detected, the system MUST freeze affected synchronization. It MUST offer only these explicit responses: restore GitHub from the roadmap, incorporate the external change through a roadmap pull request, or postpone reconciliation while the affected scope remains frozen. It MUST NOT silently overwrite either representation or discard history.

#### Scenario: External Issue drift

- GIVEN a linked Issue no longer reflects its recorded synchronized version
- WHEN drift is detected
- THEN the affected scope is frozen and the three reconciliation responses are available.

#### Scenario: Postponed reconciliation

- GIVEN an authorized owner postpones reconciliation
- WHEN the decision is recorded
- THEN the affected scope MUST remain frozen until an explicit reconciliation completes.

### Requirement: Controlled mechanical state updates

The system MAY create exact state-only bot pull requests for mechanical `in_progress` and `done` transitions and MAY auto-merge them only under prior authorization. Every other roadmap change MUST require an authorized human-approved pull request.

#### Scenario: Authorized state-only update

- GIVEN prior authorization and an exact mechanical transition to `in_progress`
- WHEN the bot proposes the state update
- THEN it may use a state-only pull request and auto-merge only under that authorization.

#### Scenario: Non-mechanical roadmap change

- GIVEN a queue, scope, dependency, or other non-mechanical roadmap change
- WHEN it is proposed
- THEN the system MUST require an authorized human-approved pull request.
