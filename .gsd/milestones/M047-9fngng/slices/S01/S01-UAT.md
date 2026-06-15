# S01: Contract cleanup baseline — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15T08:12:24.017Z

# S01: Contract cleanup baseline — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S01 changes backend parsing and error-registry contracts. The user-visible outcome is code behavior plus a closure artifact; no browser surface is involved.

## Preconditions

- `documents/issue-6-current-m047.md` exists.
- `benchmark-results/m047-s01-contract-cleanup.md` exists.

## Smoke Test

Verify safe env parsing and removal of un-emitted error codes.

## Test Cases

### 1. Safe env integer parsing

1. Inspect `api/main.go`.
2. **Expected:** `getEnvInt` uses `strconv.Atoi`, falls back for parse errors and negative values, and does not use manual multiply accumulation.

### 2. Dead error codes removed

1. Inspect `api/handlers/errors.go`.
2. **Expected:** `CodeDimensionsRequired`, `CodeDimensionsMismatch`, and `CodeRequestTimeout` are not registered.

### 3. Evidence artifact complete

1. Inspect `benchmark-results/m047-s01-contract-cleanup.md`.
2. **Expected:** artifact covers issue #6 findings #15 and #25, red evidence, green evidence, and R036 validation.

## Edge Cases

- Very large numeric env value falls back to default.
- Negative numeric env value falls back to default.
- Future registered error codes without non-test emitters fail `TestAllErrorCodesHaveNonTestEmitters`.

## Failure Signals

- `getEnvInt` contains `n = n*10` again.
- Removed codes reappear in `errorCodeRegistry` without emitters.
- `go test ./...` fails.

## Requirements Proved By This UAT

- R036: safe env integer parsing and emitted error-code registry policy.

## Not Proven By This UAT

- TEI retry/fast-fail, warmup retry, and graceful listener shutdown; these are downstream M047 slices.

## Notes for Tester

UAT evidence: `823eb0b8-5e8e-40ad-b30c-124dd1beafa1`, `d0bf2547-6034-433c-9014-ffb9d61fa0a8`, `a4460fc3-c2dc-44e4-8da9-6716496b3b98`.
