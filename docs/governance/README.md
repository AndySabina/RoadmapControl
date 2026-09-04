# Governance documentation tracker

This index tracks the public governance baseline that must land before RoadmapControl implementation begins. These documents describe bootstrap policy; they do not claim that RoadmapControl enforcement is active.

| Document | Purpose |
| --- | --- |
| [`LICENSE`](../../LICENSE) | Apache License 2.0 terms. |
| [`README.md`](../../README.md) | Honest project status, scope, platform, and trust boundaries. |
| [`CONTRIBUTING.md`](../../CONTRIBUTING.md) | Approval, authorization, assignment, worktree, review, and PR rules. |
| [`SECURITY.md`](../../SECURITY.md) | Private vulnerability reporting and bootstrap support policy. |
| [`CODE_OF_CONDUCT.md`](../../CODE_OF_CONDUCT.md) | Community behavior and enforcement process. |
| [Issue forms](../../.github/ISSUE_TEMPLATE/) | Bounded bug and feature evidence without implied approval. |
| [Pull request template](../../.github/PULL_REQUEST_TEMPLATE.md) | Review gate checklist and bootstrap/governed distinction. |

## Review chain

Recommended feature-branch chain:

```text
PR-00A license-and-governance-index
  → PR-00B project-and-contribution-guide
  → PR-00C community-and-security-policy
  → PR-00D issue-and-pull-request-templates
  → PR-01 implementation bootstrap
```

Each documentation slice must remain at or below 400 additions plus deletions and be read back before review. PR-00 is complete only when every document above is present and consistent; implementation slices PR-01 through PR-22 retain their existing numbers.
