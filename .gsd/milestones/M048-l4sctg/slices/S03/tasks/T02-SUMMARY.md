---
id: T02
parent: S03
milestone: M048-l4sctg
key_files:
  - api/middleware/validation.go
  - api/middleware/validation_test.go
  - api/openapi/spec.go
  - api/openapi/spec_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T11:28:15.625Z
blocker_discovered: false
---

# T02: Fixed validation message clarity and OpenAPI helper fail-loud behavior.

**Fixed validation message clarity and OpenAPI helper fail-loud behavior.**

## What Happened

Updated `handleBindError` so empty `json.UnmarshalTypeError.Field` produces `input must be an array of strings, got <Value>` instead of malformed `input[]...`. Updated `openapi.m()` to panic when a key is not a string rather than silently continuing. Focused middleware/openapi tests and full test suite pass.

## Verification

`cd api && go test ./middleware ./openapi` passed with 53 tests. `cd api && go test ./...` passed with 281 tests. Static proof `50f7f673-a2db-4367-bb1b-aad08226a683` passed for validation/OpenAPI source invariants.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./middleware ./openapi` | 0 | ✅ pass | 24600ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 24500ms |
| 3 | `gsd_exec 50f7f673-a2db-4367-bb1b-aad08226a683` | 0 | ✅ pass | 160ms |

## Deviations

None.

## Known Issues

Only final closure/gates remain for M048.

## Files Created/Modified

- `api/middleware/validation.go`
- `api/middleware/validation_test.go`
- `api/openapi/spec.go`
- `api/openapi/spec_test.go`
