---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Pinned S03 API polish failures with red tests.

Add focused tests proving non-string array element validation should produce a well-formed message, and `openapi.m()` should panic/fail loudly on a non-string key instead of silently continuing.

## Inputs

- `api/middleware/validation.go`
- `api/openapi/spec.go`
- `documents/issue-7-current-m048.md`

## Expected Output

- `api/middleware/validation_test.go`
- `api/openapi/spec_test.go`

## Verification

cd api && go test ./middleware ./openapi (expected red before implementation).

## Observability Impact

Red tests pin user-facing message and generator safety behavior.
