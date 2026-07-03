# S03: Runtime cache validation

**Goal:** Validate real L1/L2 cache behavior and dimension scoping.
**Demo:** Redis keys/payload sizes and warm/cold behavior prove the cache works in runtime.

## Must-Haves

- Redis keys use `embed:cache:v2:*:d<dim>` suffixes.
- 1024d payload size is 4098 bytes and 512d payload size is 2050 bytes.
- Warm request is faster than cold request.
- API restart still uses Redis L2 cache for prior key.

## Proof Level

- This slice proves: Redis CLI, timing, logs

## Integration Closure

Confirms recent cache correctness fixes work outside unit tests.

## Verification

- Captures Redis key/payload evidence and cache-miss log behavior.

## Tasks

- [x] **T01: Validate Redis key and 1024d payload** `est:small`
  Flush local Redis, make a 1024d request, and verify Redis key suffix and binary payload size.
  - Verify: Redis scan and STRLEN show d1024 key and 4098 bytes.

- [x] **T02: Validate dimension-specific cache keys** `est:small`
  Request same text at 1024d and 512d, then verify separate Redis keys and expected payload sizes.
  - Verify: Redis scan and STRLEN show d1024=4098 and d512=2050.

- [x] **T03: Validate cold versus warm behavior** `est:small`
  Measure cold and warm request timing and inspect API logs for cache miss behavior.
  - Verify: Warm request faster than cold and no repeated miss log while L1 warm.

- [x] **T04: Validate Redis L2 after API restart** `est:medium`
  Restart API and validate Redis L2 hit for an already cached key.
  - Verify: Second request after API restart succeeds without TEI miss for same key.
