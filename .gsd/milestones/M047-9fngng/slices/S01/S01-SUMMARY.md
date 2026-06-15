---
id: S01
parent: M047-9fngng
milestone: M047-9fngng
provides:
  - R036 validated.
  - Issue #6 findings #15 and #25 closed.
  - Input artifact for M047 saved.
requires:
  []
affects:
  []
key_files:
  - documents/issue-6-current-m047.md
  - api/main.go
  - api/main_env_test.go
  - api/handlers/errors.go
  - api/handlers/errors_test.go
  - benchmark-results/m047-s01-contract-cleanup.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Remove un-emitted error codes instead of marking them reserved because they are not current public behavior.
patterns_established:
  - Every public registered error code must have a non-test API emitter outside the registry file.
observability_surfaces:
  - benchmark-results/m047-s01-contract-cleanup.md records red/green evidence and residual findings.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T08:12:24.017Z
blocker_discovered: false
---

# S01: Contract cleanup baseline

**Validated and fixed issue #6 small contract findings for env integer parsing and un-emitted error codes.**

## What Happened

S01 persisted the issue #6 input artifact and used TDD to pin two contract bugs. The red tests showed `getEnvInt` did not safely fallback for overflowing numeric strings and the error registry exposed three un-emitted public codes. The implementation replaced the manual integer parser with `strconv.Atoi` plus fallback for invalid/overflowing/negative values, removed `dimensions_required`, `dimensions_mismatch`, and `request_timeout` from the public error registry, and added a static test that every registered error code has a non-test API emitter outside the registry file. R036 was validated with full tests and artifact proof.

## Verification

Red evidence: `cd api && go test ./...` failed with overflow and three un-emitted code failures. Green evidence: `cd api && gofmt -w main.go main_env_test.go handlers/errors.go handlers/errors_test.go && go test ./...` passed with 283 tests. Static proof `60cf4abe-6f44-4527-8b7a-1017cbd03e71` passed. Artifact completeness `af863d74-2a3b-46a9-a1fb-ff81d10915b7` passed. UAT PASS was saved with evidence `823eb0b8-5e8e-40ad-b30c-124dd1beafa1`, `d0bf2547-6034-433c-9014-ffb9d61fa0a8`, and `a4460fc3-c2dc-44e4-8da9-6716496b3b98`.

## Requirements Advanced

None.

## Requirements Validated

- R036 — S01 full tests and static proof validate safe env parsing and emitted error-code registry policy.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Chose removal of un-emitted codes rather than reserved metadata because no current runtime behavior emits those codes and issue #6 called out misleading public contract rows.

## Known Limitations

Issue #6 findings #13, #32, #11, and #14 are intentionally left for S02-S04.

## Follow-ups

Proceed to S02 graceful listener error path.

## Files Created/Modified

- `api/main.go` — Safe `getEnvInt` parsing via `strconv.Atoi`.
- `api/main_env_test.go` — Overflow test for env integer parsing.
- `api/handlers/errors.go` — Removed un-emitted registered error codes.
- `api/handlers/errors_test.go` — Added registry emitter test and updated envelope cases.
- `benchmark-results/m047-s01-contract-cleanup.md` — S01 evidence artifact.
- `documents/issue-6-current-m047.md` — Issue #6 input artifact.
- `.gsd/REQUIREMENTS.md` — R036 validated.
