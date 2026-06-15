---
id: S04
parent: M047-9fngng
milestone: M047-9fngng
provides:
  - R034 validated.
  - Issue #6 full closure matrix.
  - Final gate evidence for milestone validation.
requires:
  - slice: S02
    provides: Graceful listener fatal-error path.
  - slice: S03
    provides: TEI retry and fast-fail behavior.
affects:
  []
key_files:
  - api/main.go
  - api/main_test.go
  - api/embed/tei.go
  - api/embed/tei_test.go
  - api/handlers/errors_test.go
  - benchmark-results/m047-s04-warmup-retry-closure.md
  - benchmark-results/m047-issue-6-closure.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Use bounded three-attempt warmup retry with injectable policy for deterministic tests.
  - Treat all issue #6 findings as fixed within M047 scope.
patterns_established:
  - Warmup and dependency retries should be deterministic in tests via injectable policies, while production defaults stay bounded.
observability_surfaces:
  - Warmup attempt/final failure logs include attempt counts, max attempts, errors, and latency; closure artifacts record all evidence.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T08:38:12.101Z
blocker_discovered: false
---

# S04: Warmup retry and closure

**Warmup now retries bounded failures and M047 issue #6 closure is fully evidenced.**

## What Happened

S04 addressed issue #6 #14 and closed the milestone. Red tests first proved warmup retry policy did not exist. The implementation added `warmupRetryPolicy` and `startModelWarmupWithPolicy`, with default three-attempt retry/backoff behavior, per-attempt error recording/logging, terminal failure logging, and later success readiness recovery through `MarkWarmupDone`. S04 also wrote the full issue #6 closure matrix covering #11, #14, #13, #32, #25, and #15, validated R034, and ran final gates. Lint findings discovered during final gates were fixed before completion.

## Verification

Red evidence: `go test ./...` failed with undefined warmup retry policy/helper. Green evidence: `cd api && gofmt -w main.go main_test.go && go test ./...` passed with 290 tests. Static proof `7ee9815e-9837-40f9-8430-8ef343422cdf` passed. Final gates passed: `go test ./...` 290 passed, golangci-lint 0 issues, govulncheck 0 reachable vulnerabilities. Closure completeness `3e33970d-90f3-40f5-9955-7fb27633019e` passed. UAT PASS was saved with evidence `19a0bc3e-9186-40ba-98ce-c456535bae8d`, `81114e5d-c0a5-423c-8d73-0e40fceed740`, `83eeafc9-ed6a-46ab-a427-7863a5f3fb51`, and `0644ab40-ee81-4eb7-9481-22ab68329ab0`.

## Requirements Advanced

None.

## Requirements Validated

- R034 — S04 deterministic warmup retry tests and static proof validate bounded retry and readiness recovery.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Final lint discovered test/source cleanup issues; fixed them before final completion.

## Known Limitations

Live TEI restart testing was not run; deterministic retry/circuit tests cover code behavior without requiring external service disruption.

## Follow-ups

Optional: with explicit confirmation, comment on or close GitHub issue #6 after push.

## Files Created/Modified

- `api/main.go` — Added warmup retry policy/helper and final lint cleanup.
- `api/main_test.go` — Added warmup retry tests and lint cleanup.
- `api/embed/tei.go` — Named return cleanup for lint after S03.
- `api/handlers/errors_test.go` — Test-only gosec annotation for repo-local source scan.
- `benchmark-results/m047-s04-warmup-retry-closure.md` — S04 evidence artifact.
- `benchmark-results/m047-issue-6-closure.md` — Full issue #6 closure matrix.
- `.gsd/REQUIREMENTS.md` — R034 validated.
