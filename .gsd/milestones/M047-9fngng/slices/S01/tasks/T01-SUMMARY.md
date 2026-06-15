---
id: T01
parent: S01
milestone: M047-9fngng
key_files:
  - api/main_env_test.go
  - api/handlers/errors_test.go
  - documents/issue-6-current-m047.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:06:57.832Z
blocker_discovered: false
---

# T01: Pinned S01 issue #6 contract failures with red tests.

**Pinned S01 issue #6 contract failures with red tests.**

## What Happened

Added a `getEnvInt` overflow case using a 100-digit integer and a registry policy test that scans non-test Go source outside `handlers/errors.go` for emitters of every registered error code. The final red run fails exactly on oversized env integer parsing and the three issue #6 dead codes: `dimensions_required`, `dimensions_mismatch`, and `request_timeout`.

## Verification

`cd api && go test ./...` failed as expected: 283 passed, 3 failed. Failures were `TestGetEnvIntFallsBackForInvalidValues` for a 100-digit value and `TestAllErrorCodesHaveNonTestEmitters` for `CodeDimensionsRequired`, `CodeDimensionsMismatch`, and `CodeRequestTimeout`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 1 | ✅ expected red | 7500ms |

## Deviations

The initial registry test was too narrow and was refined to scan the full `api` source tree instead of only `api/handlers`.

## Known Issues

S01 remains red until T02 fixes env parsing and the registry contract.

## Files Created/Modified

- `api/main_env_test.go`
- `api/handlers/errors_test.go`
- `documents/issue-6-current-m047.md`
