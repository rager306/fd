---
id: T01
parent: S03
milestone: M048-l4sctg
key_files:
  - api/middleware/validation_test.go
  - api/openapi/spec_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T11:26:10.009Z
blocker_discovered: false
---

# T01: Pinned S03 API polish failures with red tests.

**Pinned S03 API polish failures with red tests.**

## What Happened

Added focused tests for issue #7 #24 and #31. The validation test expects a clean message for non-string array input instead of the malformed `input[]` prefix. The OpenAPI test expects `m()` to panic on non-string keys instead of silently dropping the pair.

## Verification

`cd api && go test ./middleware ./openapi` failed as expected: `TestValidationInvalidJSONNonString` observed malformed message `input[] must be string, got array`, and `TestMapHelperPanicsOnNonStringKey` observed no panic.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./middleware ./openapi` | 1 | ✅ expected red | 9000ms |

## Deviations

Adjusted the expected type text from `number` to `array` to match Go's actual `json.UnmarshalTypeError.Value` on this code path while preserving the issue's clean-message intent.

## Known Issues

S03 remains red until validation and OpenAPI helper behavior are fixed.

## Files Created/Modified

- `api/middleware/validation_test.go`
- `api/openapi/spec_test.go`
