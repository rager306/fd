# S01: Cache invalidation controls — UAT

**Milestone:** M049-7dn2gp
**Written:** 2026-06-15T12:54:56.954Z

# S01: Cache invalidation controls — UAT

**Milestone:** M049-7dn2gp
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S01 validates backend cache invalidation source and tests. Runtime HIT->flush->MISS verification is explicitly deferred to S03 after aggregate container rebuild.

## Preconditions

- `benchmark-results/m049-s01-cache-invalidation.md` exists.
- Focused and full Go tests have passed.

## Smoke Test

Verify Redis namespace safety, auth route placement, and evidence artifact completeness.

## Test Cases

### 1. Redis namespace safety

1. Inspect `api/cache/redis.go`.
2. Expected: `FlushNamespace` exists, uses namespace pattern, and does not call `FlushDB`.

### 2. Route auth placement

1. Inspect `api/main.go`.
2. Expected: `POST /v1/cache/flush` and `POST /v1/cache/delete` are registered after `APIKeyAuthFromEnv` middleware.

### 3. Evidence artifact

1. Inspect `benchmark-results/m049-s01-cache-invalidation.md`.
2. Expected: artifact records LocalCache, RedisCache, TieredCache, route, and test proof.

## Requirements Proved By This UAT

- R040 advanced for cache invalidation implementation.

## Not Proven By This UAT

- Live container cache header behavior. This is planned for S03 runtime UAT.

## Notes for Tester

Evidence IDs: `94ea4377-4e0a-4327-a167-76d5bcf0404c`, `6d55b34d-006b-431c-8045-cb8e5f639981`, `8707655c-51b1-452e-9af3-1efd9ba08dda`.
