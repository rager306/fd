# M048 S03 API Polish and Closure Evidence

Captured: 2026-06-15

## Scope

S03 covers GitHub issue #7 findings:

- #24 `handleBindError` emits a malformed message for non-string array-element type errors.
- #31 `openapi.m()` silently swallows non-string keys.

S03 also runs final M048 gates.

## Red Evidence

Command:

```bash
cd api && go test ./middleware ./openapi
```

Expected red result after adding tests:

```text
middleware: TestValidationInvalidJSONNonString observed malformed message `input[] must be string, got array`.
openapi: TestMapHelperPanicsOnNonStringKey observed no panic.
```

## Fix

- `handleBindError` now checks `json.UnmarshalTypeError.Field == ""` and emits:
  - `input must be an array of strings, got <Value>`
- `openapi.m()` now panics on a non-string key instead of silently continuing.
- Regression tests cover both behaviors.

## Green Evidence

Commands:

```bash
cd api && go test ./middleware ./openapi
cd api && go test ./...
```

Results:

```text
go test ./middleware ./openapi: 53 passed in 2 packages
go test ./...: 281 passed in 10 packages
```

Static proof:

```text
gsd_exec 50f7f673-a2db-4367-bb1b-aad08226a683
PASS M048 S03 validation and OpenAPI helper source invariants
```

## Final Gates

```bash
cd api && go test ./...
cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...
```

Results:

```text
go test ./...: 281 passed in 10 packages
golangci-lint: 0 issues
govulncheck: 0 reachable vulnerabilities
```

## Requirement Outcome

- R039 validated for issue #7 findings #24 and #31.

## Residual Issue #7 Findings

None in M048 scope. Full closure matrix is in `benchmark-results/m048-issue-7-closure.md`.
