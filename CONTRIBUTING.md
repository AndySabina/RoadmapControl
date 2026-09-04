# Contributing to RoadmapControl

RoadmapControl welcomes public discussion and carefully bounded contributions. Anyone may read, fork, open issues, and comment. Opening an issue or pull request does **not** approve it as roadmap work or grant authority to implement it upstream.

## Current bootstrap stage

RoadmapControl is not implemented or dogfooding its own controls yet. During bootstrap, maintainers use approved GitHub issues, ordinary repository protections, Gentle AI capability integration where available, and the chained pull-request discipline documented below. Bootstrap history must not be represented as RoadmapControl-governed.

## Before implementing upstream

All of these conditions are required:

1. The proposed outcome is approved for the roadmap by an authorized project role.
2. The contributor is authorized for implementation and assigned to the approved issue.
3. The work has a bounded scope and an exclusive task worktree or equivalent bootstrap lease.
4. No other active task overlaps the assigned paths.
5. The planned pull request is one cohesive work unit and normally no more than 400 authored additions plus deletions, including tests, documentation, and configuration.

The initial operating policy is owner-only. Owner, maintainer, and developer roles are planned for later operation; role names never imply authority unless the repository policy grants it for the specific action.

If you are not assigned, help by refining the issue, adding bounded reproduction evidence, or discussing design implications. Do not begin upstream implementation speculatively.

## Contribution workflow

1. Confirm the issue is approved and that you are assigned and authorized.
2. Confirm the allowed paths, acceptance criteria, dependencies, and rollback boundary.
3. Work only in the exclusive task worktree/lease assigned to that issue.
4. Keep the change to one deliverable work unit. Split larger work into a reviewed chain rather than compressing code or separating tests/docs.
5. Add or update tests and documentation in the same pull request as the behavior they verify.
6. Run the focused checks documented for the work unit.
7. Open a pull request using the repository template, link the approved issue, and apply exactly one `type:*` label.
8. Wait for required automated checks and human review. Automation or issue linkage alone does not authorize merge.

Normal pull requests must stay at or below **400 additions + deletions**. An indivisible overage requires an explicit, audited owner exception; it is not a routine waiver.

## Bootstrap chained pull requests

Before RoadmapControl can govern itself, bootstrap work follows the approved implementation slices in [`tasks.md`](openspec/changes/bootstrap-roadmapcontrol/tasks.md). Each child pull request must:

- name its predecessor and successor;
- state its position in the chain and show the dependency diagram with the current slice marked;
- contain only its current work unit;
- keep tests and documentation with the change they verify; and
- be independently reviewable and reversible.

The selected chain strategy must remain consistent across the chain. See the [governance tracker](docs/governance/README.md) for the documentation prerequisite.

## Pull-request expectations

- Link the approved issue with `Closes #N`, `Fixes #N`, or `Resolves #N`.
- Apply exactly one matching `type:*` label.
- Describe the outcome, changed surfaces, test evidence, documentation impact, and rollback boundary.
- Include chain context when applicable.
- Do not add `Co-Authored-By` or other AI attribution. The accountable human contributor remains responsible for the submitted work.
- Respond to human review and do not merge without the required authorization.

## Conduct and security

Participation is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). Report suspected vulnerabilities privately according to the [Security Policy](SECURITY.md); do not disclose secrets or exploit details in a public issue.
