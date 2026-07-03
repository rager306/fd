# S02: Cache observability and log noise

**Goal:** Improve cache-path observability and reduce high-volume success log noise while preserving API responses and cache semantics.
**Demo:** After this, runtime logs/metrics show cache-path behavior without flooding INFO logs under load.

## Must-Haves

- Cache path exposes structured debug signal for L1 hit, L2 hit, cold load, and non-fatal Redis failures.
- Successful embedding/batch requests no longer emit high-volume INFO success logs by default.
- Errors/warnings remain visible.
- Existing API responses and cache semantics are unchanged.
- Go tests pass.

## Proof Level

- This slice proves: Go tests plus runtime smoke/log inspection

## Integration Closure

Runtime logs remain useful under benchmark load and cache path can be inspected with debug logging.

## Verification

- Adds structured debug cache-path events and removes per-request success INFO spam.

## Tasks

- [x] **T01: Inspect observability edit points** `est:small`
  Inspect cache and handler symbols, run GitNexus impact analysis, and decide exact observability edit points.
  - Files: `api/cache/tiered.go`, `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/main.go`
  - Verify: Impact analysis recorded for edited symbols.

- [x] **T02: Add cache-path observability** `est:medium`
  Implement configurable log level and cache debug/warn events without changing cache semantics.
  - Files: `api/main.go`, `api/cache/tiered.go`
  - Verify: Go cache tests pass.

- [x] **T03: Reduce handler success log noise** `est:medium`
  Remove or demote handler success INFO logs and verify successful requests do not spam INFO by default.
  - Files: `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: Go handler tests pass and runtime log smoke confirms no success INFO spam.

## Files Likely Touched

- api/cache/tiered.go
- api/handlers/embeddings.go
- api/handlers/batch.go
- api/main.go
- api/handlers/embeddings_integration_test.go
