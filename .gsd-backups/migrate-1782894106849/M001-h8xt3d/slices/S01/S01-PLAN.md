# S01: Cache correctness and panic safety

**Goal:** Fix high-risk cache correctness and serialization safety issues.
**Demo:** Cache layer cannot cross-contaminate 512d and 1024d requests; short TEI vectors return errors rather than panics.

## Must-Haves

- L1 and singleflight keys include embedding dimension.
- marshalEmbedding validates vector length and returns errors.
- Redis/Tiered cache callers propagate cache serialization errors where correctness requires it.
- Tests cover 512d/1024d isolation and short-vector safety.

## Proof Level

- This slice proves: unit tests plus full Go short suite

## Integration Closure

Handlers continue using TieredCache.GetOrLoad with the same public behavior; cache internals become dimension-aware and error-safe.

## Verification

- Errors become explicit and contextual instead of panics or silent Redis write failures.

## Tasks

- [x] **T01: Assess cache blast radius** `est:small`
  Run GitNexus impact analysis for cache symbols, inspect current tests, and define the minimal cache API changes needed for dimension-aware keys and safe serialization.
  - Files: `api/cache/tiered.go`, `api/cache/redis.go`, `api/cache/*_test.go`
  - Verify: No code changes; record analysis in summary.

- [x] **T02: Implement cache correctness fixes** `est:medium`
  Update TieredCache to use dimension-aware local/singleflight keys and update Redis marshal/set path to validate vector lengths and return errors. Add focused unit tests for dimension isolation and short-vector error behavior.
  - Files: `api/cache/tiered.go`, `api/cache/redis.go`, `api/cache/local.go`, `api/cache/redis_binary_test.go`, `api/cache/tiered_test.go`
  - Verify: cd api && go test ./cache -run 'Test.*(Tiered|Binary|Redis|Local)' -count=1

- [x] **T03: Verify cache slice** `est:small`
  Run full Go tests for the slice and document results.
  - Verify: cd api && go test ./... -short

## Files Likely Touched

- api/cache/tiered.go
- api/cache/redis.go
- api/cache/*_test.go
- api/cache/local.go
- api/cache/redis_binary_test.go
- api/cache/tiered_test.go
