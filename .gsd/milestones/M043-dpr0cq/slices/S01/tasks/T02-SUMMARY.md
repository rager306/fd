---
id: T02
parent: S01
milestone: M043-dpr0cq
key_files:
  - api/cache/redis.go (errors.Is, maxUint16 bounds, errors import)
  - api/cache/local.go (package comment)
  - api/embed/codec.go (package comment)
  - api/embed/onnx_manifest.go (2x //nolint:gosec G304)
  - api/handlers/batch.go (early-return, package comment)
  - api/handlers/constants.go (deleted, was empty)
  - api/handlers/errors_test.go (paramInput, paramDimensions, paramEncodingFormat consts)
  - api/main.go (package comment, defaultValue, ReadHeaderTimeout, //nolint:gosec G304)
  - api/middleware/validation.go (deleted teiSubBatchSize const + comment)
  - api/middleware/validation_test.go (paramInput const)
key_decisions:
  - (none)
duration: 
verification_result: untested
completed_at: 2026-06-14T03:51:44.746Z
blocker_discovered: false
---

# T02: Fixed all 11 genuine issues: errorlint redis.Nil, gosec G112/G115/G304, revive package-comments + var-naming + early-return, unused consts, test fixture consts

**Fixed all 11 genuine issues: errorlint redis.Nil, gosec G112/G115/G304, revive package-comments + var-naming + early-return, unused consts, test fixture consts**

## What Happened

Fixes applied: (1) api/cache/redis.go: `err == redis.Nil` → `errors.Is(err, redis.Nil)` (added errors import). (2) api/cache/redis.go: explicit `dim > maxUint16` bounds check + `//nolint:gosec` G115 with bounds-checked-above justification. (3) api/main.go: added `ReadHeaderTimeout: 10 * time.Second` to http.Server (gosec G112 Slowloris fix). (4) api/main.go, api/embed/onnx_manifest.go (2 spots): `//nolint:gosec` G304 for env-controlled operator paths (ONNX_RUNTIME_SHA256, ONNX_ARTIFACT_MANIFEST). (5) api/main.go: `default_` → `defaultValue` (revive var-naming). (6) api/handlers/batch.go: early-return refactor (revive). (7) api/handlers/constants.go: deleted (unused `errorKey`). (8) api/middleware/validation.go: deleted unused `teiSubBatchSize` const + explanatory comment. (9) Package comments: added to api/cache/local.go, api/embed/codec.go, api/handlers/batch.go, api/main.go (revive package-comments). (10) api/handlers/errors_test.go: extracted paramInput, paramDimensions, paramEncodingFormat consts (goconst). (11) api/middleware/validation_test.go: extracted paramInput const.

## Verification

go test ./api/... все 5 packages pass. golangci-lint run --config .golangci.yml ./... 0 issues, exit 0. Phase 1 report documents каждое fix с commit reference rationale.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `api/cache/redis.go (errors.Is, maxUint16 bounds, errors import)`
- `api/cache/local.go (package comment)`
- `api/embed/codec.go (package comment)`
- `api/embed/onnx_manifest.go (2x //nolint:gosec G304)`
- `api/handlers/batch.go (early-return, package comment)`
- `api/handlers/constants.go (deleted, was empty)`
- `api/handlers/errors_test.go (paramInput, paramDimensions, paramEncodingFormat consts)`
- `api/main.go (package comment, defaultValue, ReadHeaderTimeout, //nolint:gosec G304)`
- `api/middleware/validation.go (deleted teiSubBatchSize const + comment)`
- `api/middleware/validation_test.go (paramInput const)`
