# Roadmap Domain Specification

## Purpose

Define the repository-versioned, approved-work record and its deterministic planning rules.

## Requirements

### Requirement: Versioned canonical roadmap

The system MUST store one repository roadmap as modular YAML under `.roadmap/`, validate it against a pinned versioned JSON Schema, and treat it as the canonical record of current approved work. Issues outside that record MUST NOT become approved work implicitly.

#### Scenario: Valid approved roadmap

- GIVEN one roadmap and a supported schema version
- WHEN its modular YAML is validated
- THEN the system accepts only data conforming to that schema and identifies its approved work.

#### Scenario: Invalid or unsupported roadmap data

- GIVEN malformed YAML, schema-invalid data, or an unsupported schema version
- WHEN validation is requested
- THEN the system MUST reject the affected roadmap state and MUST NOT admit dependent work.

### Requirement: Pull-request-driven schema migration

The system MUST apply schema migrations through versioned pull requests and MUST NOT silently migrate an approved roadmap.

#### Scenario: Migration is proposed

- GIVEN a roadmap using an earlier supported schema version
- WHEN a migration is prepared
- THEN the migration is versioned and reviewable in a pull request before it becomes canonical.

### Requirement: Typed hierarchy and identifiers

The system MUST support top-level `feature`, `bug`, `maintenance`, `security`, and `documentation` trackers plus owner-configured additional types. It MUST support only the bounded hierarchy `tracker -> phase -> subphase -> task`; phases and subphases MAY be omitted. Trackers and tasks MUST have associated GitHub Issues, while phases and subphases MUST NOT require one. Every R/P/S/T identifier MUST be immutable and globally unique.

#### Scenario: Minimal tracker hierarchy

- GIVEN an approved feature tracker with a task
- WHEN the hierarchy is validated
- THEN it is valid without a phase or subphase and both tracker and task require Issue associations.

#### Scenario: Invalid identifier or hierarchy

- GIVEN a duplicate identifier, an identifier change, or a node outside the bounded hierarchy
- WHEN validation is requested
- THEN the system MUST reject the change.

### Requirement: Dependencies and executable queue

The system MUST permit acyclic dependencies at every hierarchy level, including inherited and same-repository cross-tracker dependencies. The explicit queue MUST be the sole operational ordering mechanism; the system MUST NOT infer or calculate a separate priority order.

#### Scenario: Dependency-ready queued task

- GIVEN a queued task whose direct and inherited dependencies are satisfied
- WHEN eligibility is evaluated
- THEN the task is dependency-ready according to the explicit queue.

#### Scenario: Cyclic or unresolved dependency

- GIVEN a dependency cycle or an unsatisfied dependency
- WHEN eligibility is evaluated
- THEN the system MUST reject the cycle or mark the task ineligible without reordering the queue.

### Requirement: State, blocking, and progress rules

The system MUST use exactly `not_started`, `in_progress`, `blocked`, `done`, `cancelled`, and `superseded` states. Terminal states MUST be immutable. When work becomes blocked, the system MUST preserve its prior state, blocking reason, and resume condition. Progress MUST use equal task weights, exclude cancelled and superseded tasks, and derive aggregate parent states by default; an exceptional parent-state override MUST include justification.

#### Scenario: Block and resume

- GIVEN an in-progress task that cannot continue
- WHEN it is marked blocked
- THEN its prior state, reason, and resume condition are retained for a later authorized resumption.

#### Scenario: Terminal or unjustified override

- GIVEN terminal work or a requested parent-state override without justification
- WHEN a state change is attempted
- THEN the system MUST reject the state change or override.

#### Scenario: Aggregate progress

- GIVEN a parent with done, active, cancelled, and superseded tasks
- WHEN progress is calculated
- THEN each included task has equal weight and cancelled and superseded tasks are excluded.

### Requirement: Controlled roadmap promotion

Only an owner or maintainer MUST be permitted to promote an Issue or other work into the canonical roadmap. The default installation policy MUST restrict promotion authority to owners; a maintainer MAY promote work only when an owner has explicitly configured that authority.

#### Scenario: Default promotion policy

- GIVEN an installation with no configured promotion-authority override
- WHEN a maintainer requests promotion of an Issue into the roadmap
- THEN the system MUST deny the request.

#### Scenario: Configured maintainer promotion

- GIVEN an owner has explicitly configured maintainer promotion authority
- WHEN a maintainer promotes an Issue into the roadmap
- THEN the system MAY accept the promotion.

### Requirement: Planning source ownership

The canonical roadmap MUST own planning fields and current approved-work state. Linked GitHub Issues MUST retain ownership of comments, discussion, and evidence history. GitHub Projects MUST be derived representations only and MUST NOT be treated as canonical planning authority.

#### Scenario: Derived Project conflicts with roadmap planning data

- GIVEN a Project field differs from the linked canonical roadmap planning field
- WHEN the system determines planning authority
- THEN it MUST treat the roadmap field as authoritative and the Project field as derived data.

#### Scenario: Issue discussion history is available

- GIVEN a linked Issue contains comments, discussion, or evidence
- WHEN roadmap data is synchronized or updated
- THEN the system MUST preserve that Issue history as Issue-owned content.
