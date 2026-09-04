# Delivery and Release Specification

## Purpose

Define reviewable task delivery, tracker integration, verification boundaries, quota gates, and forward-only releases.

## Requirements

### Requirement: Tracker and task pull-request topology

The system MUST create one tracker branch and one draft tracker pull request for each top-level tracker. RoadmapControl MUST create normal task pull requests only from authorized user intent and MUST target them at the tracker branch. The `roadmapcontrol/` system namespace MUST be configurable and MUST NOT take ownership of existing project branches.

#### Scenario: Authorized task pull request

- GIVEN an authorized task workspace under a tracker
- WHEN a task pull request is created through RoadmapControl
- THEN it targets that tracker branch and uses the configured system namespace only for system-managed branches.

#### Scenario: Unauthorized or wrong target pull request

- GIVEN absent authorized user intent or a task pull request targeting another branch
- WHEN normal-flow validation runs
- THEN the system MUST reject the pull request from the normal task flow.

### Requirement: Review budget and human review

Every normal pull request MUST have authorized human review and MUST remain at or below 400 authored additions plus deletions, including tests, documentation, and configuration. The system MUST allow an indivisible over-budget change only with an audited owner exception.

#### Scenario: Normal reviewable change

- GIVEN a task pull request within the authored-change budget and an authorized reviewer
- WHEN review requirements are evaluated
- THEN the pull request may satisfy the normal review gate.

#### Scenario: Oversized change without exception

- GIVEN a normal pull request above 400 authored changed lines
- WHEN it seeks approval
- THEN the system MUST block it unless an audited owner exception covers indivisible work.

### Requirement: Tracker completion and promotion

The system MUST allow tracker approval only after every child pull request has been reviewed and merged. It MUST merge completed trackers into a configurable logical development branch. Production promotion and release MUST be separately governed even when development and production map to the same physical branch.

#### Scenario: Incomplete tracker

- GIVEN a tracker with an unreviewed or unmerged child pull request
- WHEN tracker approval is requested
- THEN the system MUST deny approval.

#### Scenario: Completed tracker promotion

- GIVEN all tracker child pull requests are reviewed and merged
- WHEN authorized tracker promotion is requested
- THEN the tracker may merge to the configured logical development branch and remain subject to separate release governance.

### Requirement: Boundary verification

The system MUST run focused checks at the task boundary, the full suite and tracker-outcome verification at the tracker boundary, and the full suite, build checks, and version checks at release. GitHub-native verification MUST use GitHub-hosted Actions runners and MUST document their dependencies and limits honestly.

#### Scenario: Release verification failure

- GIVEN a release candidate with a failed full-suite, build, or version check
- WHEN release eligibility is evaluated
- THEN the system MUST block release and retain failure evidence for recovery.

### Requirement: Actions quota policy

The system MUST classify Actions quota usage as healthy at 0–69%, warning at 70–84%, conservation at 85–94%, critical at 95–99%, and exhausted when no quota remains. It MUST block admission of new tasks at critical usage and MUST fail closed when exhausted. Documentation MAY state an initial expectation of approximately 30 pull requests per month but MUST NOT present it as guaranteed capacity.

#### Scenario: Critical quota

- GIVEN Actions usage at 95–99%
- WHEN a new task requests admission
- THEN the system MUST deny admission.

#### Scenario: Exhausted quota recovery

- GIVEN no Actions quota remains
- WHEN an operation requires quota-backed verification
- THEN the system MUST fail closed until available quota and required verification evidence are restored.

### Requirement: Versioning and rollback

The system MUST use SemVer by default while supporting CalVer and custom versioning policies. It MUST perform rollback through a new forward-only pull request and version, and MUST NOT rewrite published history.

#### Scenario: Released defect

- GIVEN a released version requires correction
- WHEN rollback is initiated
- THEN the system creates a forward-only reviewed correction and version rather than rewriting history.

### Requirement: Authorized production decision

Production promotion and release MUST require an explicit decision by an authorized authenticated human. Automated checks, prior tracker approval, or agent authority MUST NOT substitute for that production decision.

#### Scenario: Production release without an explicit human decision

- GIVEN a tracker is complete and all required release checks pass
- WHEN no authorized authenticated human has explicitly approved production promotion or release
- THEN the system MUST block production promotion and release.

#### Scenario: Authorized production release

- GIVEN release checks pass and an authorized authenticated human explicitly approves production promotion or release
- WHEN the release is requested
- THEN the system MAY proceed with the separately governed production action.
