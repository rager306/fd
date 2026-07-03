---
id: T02
parent: S01
milestone: M047-9fngng
key_files:
  - api/main.go
  - api/main_env_test.go
  - api/handlers/errors.go
  - api/handlers/errors_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:09:53.865Z
blocker_discovered: false
---

# T02: Fixed safe env integer parsing and removed un-emitted error codes from the public registry.

**Fixed safe env integer parsing and removed un-emitted error codes from the public registry.**

## What Happened

Replaced the hand-rolled `getEnvInt` digit loop with `strconv.Atoi`, falling back for invalid, overflowing, or negative values. Removed the un-emitted registered error codes `dimensions_required`, `dimensions_mismatch`, and `request_timeout` from constants, registry, `AllErrorCodes`, and envelope shape tests. Added a static registry emitter test that scans non-test API source outside the registry file so future un-emitted codes are caught.

## Verification

`cd api && gofmt -w main.go main_env_test.go handlers/errors.go handlers/errors_test.go && go test ./...` passed with 283 tests in 9 packages. Static proof `60cf4abe-6f44-4527-8b7a-1017cbd03e71` passed for `strconv.Atoi` fallback and removal of un-emitted codes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && gofmt -w main.go main_env_test.go handlers/errors.go handlers/errors_test.go && go test ./...` | 0 | ✅ pass | 7400ms |
| 2 | `gsd_exec 60cf4abe-6f44-4527-8b7a-1017cbd03e71` | 0 | ✅ pass | 1153ms |

## Deviations

Chose removal rather than reserved metadata because the three codes had no active emitters and issue #6 requested avoiding un-emitted public contract rows.

## Known Issues

S02-S04 issue #6 findings remain open: graceful listener error path, TEI retry/fast-fail, and warmup retry.

## Files Created/Modified

- `api/main.go`
- `api/main_env_test.go`
- `api/handlers/errors.go`
- `api/handlers/errors_test.go`
