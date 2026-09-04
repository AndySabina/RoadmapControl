# Pull request

> Opening a pull request does not grant roadmap approval, task authority, or merge authorization.

## Approved issue

Closes #<!-- approved issue number -->

- [ ] The linked issue is approved for this work.
- [ ] I am authorized and assigned to implement it upstream.
- [ ] This change stays within the issue's approved outcome and allowed paths.

## Type

Check exactly one box and apply exactly one matching `type:*` label:

- [ ] `type:bug` — bug fix
- [ ] `type:feature` — new feature
- [ ] `type:docs` — documentation only
- [ ] `type:refactor` — code refactoring
- [ ] `type:chore` — maintenance, tooling, build, CI, style, or tests
- [ ] `type:breaking-change` — breaking change

## Summary

- What outcome does this pull request deliver?
- What is intentionally out of scope?

## Changed surfaces

| Path | Change and reason |
| --- | --- |
| `path/to/file` | Describe the bounded change. |

## Review budget

- Additions: <!-- number -->
- Deletions: <!-- number -->
- Total additions + deletions: <!-- number -->
- [ ] The normal authored total is **400 or less**, including tests, documentation, and configuration.
- [ ] Or: an authorized owner approved an audited `size:exception` bound to this exact candidate and its indivisibility rationale.

Do not omit tests, documentation, comments, or useful whitespace to fit the budget. Split by cohesive work unit instead.

## Verification and documentation

- Focused command and observed result: <!-- exact command and result -->
- Cumulative command and observed result: <!-- exact command and result, or N/A with reason -->
- Runtime scenario and observed result: <!-- scenario and result, or N/A with reason -->
- [ ] Tests for changed behavior are included in this pull request.
- [ ] Documentation for changed behavior is included in this pull request.
- [ ] Sensitive output and private data were removed from evidence.

## Chain context

Complete this section for chained work; otherwise write `Not applicable`.

- Strategy: <!-- stacked-to-main or feature-branch-chain -->
- Position: <!-- PR-X of PR-Y -->
- Predecessor: <!-- link or none -->
- Successor: <!-- link or planned slice -->
- Start state: <!-- what is true before this slice -->
- End state: <!-- what this slice makes true -->
- Rollback boundary: <!-- exact paths/behavior removable without unrelated work -->

```text
PR-A → PR-B 📍 → PR-C
```

## Governance stage

Check exactly one:

- [ ] **Bootstrap:** RoadmapControl is not active; this pull request follows approved issues, repository protections, Gentle AI capability integration where available, and bootstrap chained-PR discipline.
- [ ] **RoadmapControl-governed:** attach the task identifier, lease/grant evidence, tracker target, and required control receipts below.

Governed evidence: <!-- task ID and evidence links, or N/A during bootstrap -->

## Contributor checklist

- [ ] I used an exclusive task worktree or approved bootstrap lease and did not mix tasks.
- [ ] Required automated checks pass on this exact candidate.
- [ ] A human review is requested; automation does not replace it.
- [ ] Commit messages follow the repository's conventional format.
- [ ] The commits and pull request contain no `Co-Authored-By` trailer or other AI attribution.
