# S03: API polish and closure

**Goal:** Resolve issue #7 findings #24 and #31, write issue #7 closure matrix, run final gates, validate R039, and close M048.
**Demo:** Validation messages and OpenAPI helper failures are clear, then issue #7 closes with final gates and closure matrix.

## Must-Haves

- Red tests reproduce malformed array-element validation message and silent OpenAPI helper key-drop behavior.
- Validation message for non-string array elements is well-formed.
- `openapi.m()` misuse fails loudly with a regression test.
- Issue #7 closure matrix covers #19, #24, #26, #27, #28, #29, #30, #31.
- Final tests, lint, govulncheck, artifact UAT, milestone validation, and completion pass.

## Proof Level

- This slice proves: Focused middleware/openapi tests, full Go tests, lint, govulncheck, artifact UAT.

## Integration Closure

OpenAI-compatible error envelope remains canonical; only message clarity and generator misuse behavior change.

## Verification

- Issue #7 closure artifact records all P3 cleanup proof.

## Tasks

- [x] **T01: Pinned S03 API polish failures with red tests.** `est:small`
  Add focused tests proving non-string array element validation should produce a well-formed message, and `openapi.m()` should panic/fail loudly on a non-string key instead of silently continuing.
  - Files: `api/middleware/validation_test.go`, `api/openapi/spec_test.go`
  - Verify: cd api && go test ./middleware ./openapi (expected red before implementation).

- [x] **T02: Fixed validation message clarity and OpenAPI helper fail-loud behavior.** `est:small`
  Guard empty `json.UnmarshalTypeError.Field` in validation and change `openapi.m()` non-string key behavior from silent continue to panic/fail-loud. Keep existing valid OpenAPI spec tests green.
  - Files: `api/middleware/validation.go`, `api/middleware/validation_test.go`, `api/openapi/spec.go`, `api/openapi/spec_test.go`
  - Verify: cd api && go test ./middleware ./openapi && cd api && go test ./...

- [x] **T03: Ran final M048 gates, wrote issue #7 closure matrix, and validated R039.** `est:medium`
  Write S03 and issue #7 closure artifacts, validate R039, run full tests/lint/govulncheck, save UAT, validate milestone, complete milestone, and commit final changes.
  - Files: `benchmark-results/m048-s03-api-polish-closure.md`, `benchmark-results/m048-issue-7-closure.md`, `.gsd/REQUIREMENTS.md`
  - Verify: cd api && go test ./...; golangci-lint; govulncheck; artifact UAT; milestone validation.

## Files Likely Touched

- api/middleware/validation_test.go
- api/openapi/spec_test.go
- api/middleware/validation.go
- api/openapi/spec.go
- benchmark-results/m048-s03-api-polish-closure.md
- benchmark-results/m048-issue-7-closure.md
- .gsd/REQUIREMENTS.md
