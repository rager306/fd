# S01: Cache invalidation controls

**Goal:** Implement issue #8 AN-A: namespace-safe authenticated cache invalidation for solo operator/agent use.
**Demo:** After this slice, an authenticated operator can flush fd's embedding cache and observe MISS/HIT behavior changing without restarting services.

## Must-Haves

- Cache interfaces expose Delete and Flush primitives for active LocalCache/Redis/TieredCache paths.
- Redis flush is namespace-scoped and does not call FlushDB.
- Authenticated HTTP cache flush/delete route(s) exist under the existing API key middleware.
- Tests prove cache invalidation behavior and auth protection.
- Live rebuilt-container smoke can demonstrate MISS -> HIT -> flush -> MISS.

## Proof Level

- This slice proves: Red/green cache and handler tests; full Go tests for integration; runtime smoke in S03.

## Integration Closure

Use existing APIKeyAuthFromEnv middleware rather than a new admin-token model; use existing cache namespace semantics.

## Verification

- Adds explicit cache invalidation actions that agents can invoke and verify through cache headers.

## Tasks

- [x] **T01: Pinned cache invalidation gaps with red tests.** `est:medium`
  Add focused failing tests for Redis/Local/Tiered cache Delete/Flush behavior and authenticated HTTP cache invalidation route expectations. Tests should prove namespace safety and protected route behavior before implementation.
  - Files: `api/cache/local_test.go`, `api/cache/redis_test.go`, `api/cache/tiered_test.go`, `api/handlers/cache_test.go`, `api/main_test.go`
  - Verify: cd api && go test ./cache ./handlers (expected red before implementation).

- [x] **T02: Implemented namespace-safe cache invalidation primitives and HTTP routes.** `est:medium`
  Implement Delete/Flush on active cache types, keeping Redis operations namespace-scoped. Add an authenticated handler/route for cache flush and optionally key deletion if a safe key format is available. Preserve existing cache HIT/MISS behavior.
  - Files: `api/cache/local.go`, `api/cache/redis.go`, `api/cache/tiered.go`, `api/handlers/cache.go`, `api/main.go`
  - Verify: cd api && go test ./cache ./handlers && cd api && go test ./...

- [x] **T03: Recorded S01 evidence, advanced R040, and saved UAT for AN-A.** `est:small`
  Write S01 evidence artifact, run static proof for AN-A, validate R040, save UAT for cache invalidation source/tests, and complete the slice.
  - Files: `benchmark-results/m049-s01-cache-invalidation.md`, `.gsd/REQUIREMENTS.md`
  - Verify: Artifact/static UAT proves AN-A implementation; focused/full tests pass.

## Files Likely Touched

- api/cache/local_test.go
- api/cache/redis_test.go
- api/cache/tiered_test.go
- api/handlers/cache_test.go
- api/main_test.go
- api/cache/local.go
- api/cache/redis.go
- api/cache/tiered.go
- api/handlers/cache.go
- api/main.go
- benchmark-results/m049-s01-cache-invalidation.md
- .gsd/REQUIREMENTS.md
