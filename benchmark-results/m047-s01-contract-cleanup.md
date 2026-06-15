# M047 S01 Contract Cleanup Evidence

Captured: 2026-06-15

## Scope

S01 covers GitHub issue #6 small contract findings:

- #15 `getEnvInt` overflow silently disables capacity gate.
- #25 registered error codes with no non-test emitters.

Input issue artifact: `documents/issue-6-current-m047.md`.

## Red Evidence

Command:

```bash
cd api && go test ./...
```

Expected red result after adding tests:

```text
283 passed, 3 failed in 9 packages
TestGetEnvIntFallsBackForInvalidValues failed for a 100-digit integer.
TestAllErrorCodesHaveNonTestEmitters failed for:
- dimensions_required / CodeDimensionsRequired
- dimensions_mismatch / CodeDimensionsMismatch
- request_timeout / CodeRequestTimeout
```

## Fix

- `api/main.go`: `getEnvInt` now uses `strconv.Atoi` and falls back for invalid, overflowing, or negative values.
- `api/handlers/errors.go`: removed the three un-emitted public error codes from constants, registry, and `AllErrorCodes()`:
  - `dimensions_required`
  - `dimensions_mismatch`
  - `request_timeout`
- `api/handlers/errors_test.go`: added `TestAllErrorCodesHaveNonTestEmitters`, which scans non-test API source outside the registry file and fails if a registered code has no emitter.
- `api/main_env_test.go`: added overflow coverage with a 100-digit integer.

## Green Evidence

Command:

```bash
cd api && gofmt -w main.go main_env_test.go handlers/errors.go handlers/errors_test.go && go test ./...
```

Result:

```text
283 passed in 9 packages
```

Static proof:

```text
gsd_exec 60cf4abe-6f44-4527-8b7a-1017cbd03e71
PASS M047 S01 safe env parsing and emitted error registry proof
```

## Requirement Outcome

- R036 validated for env integer parsing and emitted error registry policy.

## Residual Issue #6 Findings

Deferred to downstream M047 slices:

- #13 graceful listener error path.
- #32 `errors.Is(http.ErrServerClosed)`.
- #11 TEI retry/backoff/fast-fail.
- #14 warmup retry/backoff.
