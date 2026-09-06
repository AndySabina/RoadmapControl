# PR02B Addendum — Conditional Typed Contract and Canonical JSON Proposal

**Approved contract proposal:** authorized product approval is recorded on the
approved and assigned [issue #38](https://github.com/AndySabina/RoadmapControl/issues/38). It approves
this bounded contract and the conditional use of pinned
`github.com/gowebpki/jcs` v1.0.1 through a private, fixed aggregate canonicalizer.
It does not implement, install, publish, merge, or create a pull request.

## Review path and cohesive chain

Review the contract, evidence and limitations, then the split. The chain is
`main ← tracker #31 ← docs #37 ← this docs PR (not created) → B1 issue #39 → B2 issue #40 → C`.
Tracker #31 remains unchanged, draft, and unmerged; docs #37 is published as an
open draft and remains unmerged. This approved [issue #38](https://github.com/AndySabina/RoadmapControl/issues/38)
addendum is the current documentation delta, has no PR yet, and is not part of PR #37.

| Topic | Proposed decision |
| --- | --- |
| First module | `policy.yaml` only |
| Fields | `kind: policy`; `additional_tracker_types: string[]` |
| Schema identity | exact logical identifier below; never fetched |
| Canonical bytes | RFC 8785 JCS through the constrained private wrapper |
| Deferred | generic loading, resolver, hashes, CLI: PR02C |

## Approved closed policy contract

The manifest schema URI remains an absolute logical identifier, not a network
endpoint:

```text
https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json
```

A valid module is exactly:

```yaml
kind: policy
additional_tracker_types:
  - research
  - compliance
```

`kind` is required and must be the YAML string `policy`.
`additional_tracker_types` is required and must be a YAML sequence of strings.
The field set is closed: reject missing fields, `null`, wrong node types, unknown
or duplicate fields, and any other kind. Empty strings, duplicate strings,
built-in tracker-type names, and an empty sequence are structurally valid; their
meaning, uniqueness, collisions, authority, and lifecycle are deferred to PR03.

`ValidateYAML` in [manifest.go](../../../internal/domain/roadmap/manifest.go)
already enforces one UTF-8 document, a 1 MiB limit, depth limit, string mapping
keys, and no duplicate keys or aliases. PR02B1 adds a new decoder/constructor
that requires exact string equality between the declared schema URI and the
logical identifier above, and enforces this aggregate's path and kind
restrictions; it must not change existing manifest parsing or filesystem
behavior. The only accepted mapping is `policy.yaml → policy`, with exactly one
declared module. Alternate paths, additional modules, alternate schema URIs,
and future kinds need separately approved extensions.

A valid empty `additional_tracker_types` must normalize to non-nil `[]`, not
`null`, in the aggregate and canonical JSON.

## Approved canonical JSON boundary and limitation

The decoded aggregate is self-identifying; its illustrative JCS bytes are:

```json
{"modules":[{"kind":"policy","path":"policy.yaml","policy":{"additional_tracker_types":["research","compliance"]}}],"schema":"https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json"}
```

JCS orders object names by UTF-16 code units; arrays retain declared order. It
does not normalize Unicode, so composed and decomposed text are separately valid
and have distinct bytes. The wrapper must use fixed typed structs only: no maps,
`any`, `json.RawMessage`, custom marshalers, raw JSON, or numeric API. Validate
**every** string with `utf8.ValidString` before `json.Marshal`; otherwise Go may
repair invalid bytes. Only successfully marshaled bytes may reach `Transform`.
A literal U+FFFD is valid. Re-evaluate this boundary before adding numeric fields,
raw JSON, or custom marshaling.

The candidate is not represented as fixing all input defects. Both evaluated JCS
libraries accept invalid UTF-8, repair malformed JSON surrogate-escape pairs to U+FFFD,
and accept malformed numeric forms including `+1`, `01`, `1.`, and `1 2`.
The typed wrapper and pre-marshal validation restrict this proposal's exposure;
they do not make the libraries generally safe parsers.

## Conditional candidate recommendation and evidence

The approved conditional candidate is `github.com/gowebpki/jcs` v1.0.1 pinned at
`1a4242a66e1a8e03d7458324d0bc95c327527cbb`, Apache-2.0, production standard
library only, no CGO, and `go.mod` Go 1.15. Its upstream test graph includes
`testify` v1.7.0 only for tests. Pinned production source was byte-identical in
the bounded review. Source: <https://github.com/gowebpki/jcs/tree/1a4242a66e1a8e03d7458324d0bc95c327527cbb>.

The alternative `github.com/deszhou/jcs` v1.0.0 at
`e3a84bdb40cfbdb929af946e7842ba5c56d3c603` is rejected as a replacement: it has
the same raw-parser defects, requires newer Go 1.24, and showed no correctness
advantage. Source: <https://github.com/deszhou/jcs/tree/e3a84bdb40cfbdb929af946e7842ba5c56d3c603>.

Bounded offline scratch evidence used Go 1.26.1, not the project Go 1.27.1.
Pinned production comparison; RFC serialization, UTF-16 ordering, eight number
vectors, aggregate, array-order, escaping, and Unicode-distinctness checks passed.
A representative private-wrapper run passed five top-level tests, fifteen
subtests, and twenty checks: invalid `0xff` and surrogate-encoded UTF-8 were
rejected in schema, path, kind, and each tracker-type string;
literal U+FFFD, CJK, emoji, and HTML were valid; empty arrays were non-nil `[]`;
ordered arrays and valid UTF-8 JSON output were verified. An initial scratch
syntax error was corrected before the passing rerun.

This is not full proof: full RFC conformance, fuzzing, the full upstream suite,
and repository integration were not run offline. No production tests have been
added. Future review must reproduce relevant checks with the project toolchain;
no local machine paths are repository evidence.

## Approved split, surfaces, and forecast

1. **Docs issue [#38](https://github.com/AndySabina/RoadmapControl/issues/38):** this approved addendum; its PR is not created.
2. **PR02B1 [approved and assigned issue #39](https://github.com/AndySabina/RoadmapControl/issues/39):** typed policy decoder/constructor, tests, and matching format docs; **~180–300 changed lines**.
3. **PR02B2 [approved and assigned issue #40](https://github.com/AndySabina/RoadmapControl/issues/40):** private canonical wrapper, JCS dependency, tests, and matching format docs; **~180–350 changed lines**.
4. **PR02C:** generic typed loading, local schema allowlist/resolution, relative bundled `$ref` handling, schema validation, hash integration, and CLI integration.

Issues #39 and #40 additionally authorize bounded slice-status updates to this
change's `tasks.md` and `apply-progress.md`. The owner's separate authorization
permits a local committed RED checkpoint before GREEN. The failed intermediate
head is evidence only and must never itself be pushed or published; its RED
commit may appear in a passed candidate's ancestry only when the passed
candidate HEAD is the ref pushed or published. After GREEN, TRIANGULATE, and
REFACTOR, the completed candidate requires focused and cumulative verification,
then Judgment Day final verification, before its delivery commit, push, or draft
PR publication. No other paths, scope expansion, or merges are authorized.

Each implementation PR is capped at 400 total changed lines including configuration,
documentation, and tests; stop and refine before exceeding it. Estimates are not
proof of actual size. If a human bundles this proposal with implementation,
recompute using the **actual new documentation line count plus the relevant
estimate**. A combined slice may exceed 400; any actual or forecast overage
requires splitting. Do not reuse a historical document line count.

The seven implementation/format/dependency paths, plus the two shared progress
paths named above, form the exact issue-scoped allowlist:

- **B1 / issue #39:** `internal/domain/roadmap/module.go`, `internal/domain/roadmap/module_test.go`, and `docs/roadmap-format.md`.
- **B2 / issue #40:** `internal/domain/roadmap/canonical.go`, `internal/domain/roadmap/canonical_test.go`, `docs/roadmap-format.md`, `go.mod`, and `go.sum`; the last two are conditional on direct candidate `jcs` v1.0.1. Checksum metadata for actual test dependencies may appear only after review of the actual diff, with no extra production dependencies.

The shared format document is one allowlisted path; B1 and B2 remain separate,
sequential work units. No fixtures, schemas, filesystem paths, or CLI paths are allowed.

Existing strict TDD remains mandatory for implementation: observe RED before code, GREEN after the
minimum change, TRIANGULATE with rejection/Unicode/order cases, then REFACTOR.
Run `GOTOOLCHAIN=go1.27.1 go test ./internal/domain/roadmap/...` for focused evidence
and `GOTOOLCHAIN=go1.27.1 go test ./...` for cumulative verification. Neither command
was run against a RoadmapControl implementation in this assessment.

Acceptance tests must cover missing/null/unknown/duplicate fields, wrong node and
item types, alternate paths, extra modules, and an otherwise-valid policy whose
schema is `urn:other` and is rejected; also cover required empty-list behavior,
invalid UTF-8 before marshaling, and exact canonical bytes. Preserve the existing
YAML security checks. Wrapper proof is not a substitute for these integration tests.

Rollback each implementation slice together with its tests, matching docs, and any
dependency changes. Reverse child dependencies first and preserve the PR02A parser
and filesystem behavior; no history rewrite or merge is authorized.

## Non-goals and final gate

This addendum does not approve schema files, generic resolver behavior, hashes,
CLI commands, arbitrary module kinds, YAML-loading changes, business policy,
tracker-type authorization, installation, publication outside the verified
delivery sequence above, or PR02C. It does not change the formal bootstrap
proposal or claim PR02 complete.

**Final gate: satisfied.** The recorded approval on issue #38 covers the exact
contract and path mapping, conditional dependency recommendation and limitations,
issue-scoped implementation, dependency and progress paths, PR02B1/PR02B2 split and line cap,
and issue workflow. B1 and B2 may proceed only within their assigned issue
boundaries, with strict TDD, focused verification, and Judgment Day after
verification. Any future out-of-scope expansion still requires human approval.

## References

- [Formal bootstrap proposal](proposal.md) (unchanged)
- [Roadmap-domain specification](specs/roadmap-domain/spec.md)
- [Domain and execution design](design/domain-and-execution.md)
- [Implementation plan and PR02 reconciliation](tasks.md)
- [Current apply progress](apply-progress.md)
- [Manifest safety parser](../../../internal/domain/roadmap/manifest.go)
- [Filesystem acquisition boundary](../../../internal/adapters/filesystem/manifest.go)
