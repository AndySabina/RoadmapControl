# RoadmapControl

RoadmapControl is a planned GitHub-native control system for keeping approved roadmap work traceable from planning through execution, review, and release.

> **Bootstrap status:** implementation has started with a minimal Go composition root and unprivileged pull-request verification. The `roadmapctl` executable currently provides bootstrap help only—there are no product commands, roadmap enforcement, task leases, or agent controls.

## The problem

A GitHub issue backlog records possible work, not necessarily work approved for execution. Issues, Projects, branches, worktrees, pull requests, and agent sessions can drift apart. RoadmapControl is designed to make the current approved plan explicit and to deny control-plane transitions when required authority or evidence is unknown.

## Planned core concepts

| Concept | Planned responsibility |
| --- | --- |
| Canonical roadmap | Versioned, modular YAML under `.roadmap/` containing approved current work. |
| Explicit queue | Human-selected execution order; RoadmapControl will not invent priorities. |
| Bounded task | Approved issue, authorized assignee, outcome criteria, allowed paths, and dependency readiness. |
| Exclusive workspace | Signed lease, system task branch, sibling worktree, and one active agent grant. This is a cooperative control, not protection from a malicious machine owner. |
| Traceable delivery | Task pull requests feed a tracker branch; tracker and release boundaries require progressively broader verification. |
| Reconciliation | Detected drift freezes the affected scope until an authorized, history-preserving resolution is selected. |
| Audit evidence | Planned signed, append-only evidence reconstructible from Git and GitHub. |

The initial policy is owner-only. The model is designed to add owner, maintainer, and developer roles later without weakening explicit authorization.

## GitHub Actions dependency and limits

The planned first release depends on GitHub-hosted Actions runners for GitHub-native verification and control workflows. Action dependencies will be allowlisted and pinned to full commit SHAs. Pull-request verification will be unprivileged; privileged workflows will not execute pull-request code or expose control secrets.

GitHub plan limits, service availability, permissions, and Actions quota can prevent required evidence from being produced. RoadmapControl is designed to warn as quota usage rises, block new task admission at critical usage, and fail closed when quota is exhausted or required evidence is unavailable. The planning target of about 30 pull requests per month is an operating estimate, not guaranteed capacity or a GitHub quota.

GitHub Free cannot technically prevent a repository owner from deliberately bypassing configured controls. The planned response is detection, freeze, comparison, and explicit reconciliation—not a claim of absolute prevention.

## Initial compatibility target

- Public and private repositories on github.com Free or higher
- Linux, WSL2, and macOS
- Go implementation and CLI
- GitHub Actions automation
- One roadmap per repository
- Dependencies within one repository

GitHub Enterprise Server, native Windows, cross-repository dependencies, and multiple roadmaps per repository are not in the initial scope.

## Hosting and privacy boundary

The first release is planned without a RoadmapControl SaaS or publisher-operated backend. It will not collect telemetry, store publisher-managed tokens, or transfer RoadmapControl data outside GitHub. Local tooling will reuse established authentication mechanisms rather than persist user tokens.

## Local bootstrap prerequisites

Install Go 1.27.1, then run the local checks from the repository root:

```sh
go test ./...
go run ./cmd/roadmapctl --help
```

## Documentation

- [Contributing policy](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Governance documentation tracker](docs/governance/README.md)
- [Approved bootstrap proposal](openspec/changes/bootstrap-roadmapcontrol/proposal.md)
- [Implementation plan](openspec/changes/bootstrap-roadmapcontrol/tasks.md)

## License

Licensed under the [Apache License 2.0](LICENSE).
