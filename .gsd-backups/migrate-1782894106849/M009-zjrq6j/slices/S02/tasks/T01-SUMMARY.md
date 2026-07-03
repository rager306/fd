---
id: T01
parent: S02
milestone: M009-zjrq6j
key_files:
  - api/cache/redis.go
  - api/main.go
  - benchmark.py
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:34:27.856Z
blocker_discovered: false
---

# T01: Designed S02 as an opt-in model-aware namespace and retention config that preserves current default Redis keys and 24h TTL.

**Designed S02 as an opt-in model-aware namespace and retention config that preserves current default Redis keys and 24h TTL.**

## What Happened

Designed the cache config surface. Add `RedisCacheOptions` and `RedisCacheNamespace` in the cache package. Defaults preserve current behavior: prefix `embed:cache:`, namespace `v2`, TTL `24h`, no-expire disabled. Env parsing should support `REDIS_CACHE_TTL`, `REDIS_CACHE_NO_EXPIRE`, `EMBEDDING_CACHE_VERSION`, `EMBEDDING_MODEL_ID`, `EMBEDDING_MODEL_REVISION`, `EMBEDDING_TOKENIZER_VERSION`, and `EMBEDDING_CHUNKING_VERSION`. If namespace fields are set, the Redis key namespace should include short hashes of correctness-affecting values to avoid raw model/tokenizer strings in keys. If `REDIS_CACHE_NO_EXPIRE=true` and `REDIS_CACHE_TTL` is also explicitly set, validation should reject the config. Invalid durations and booleans should also fail fast from `main`. Benchmark snapshot allowlist should include the new env fields plus `MODEL_ID`/`REDIS_POOL_SIZE` for runtime comparability.

## Verification

Read current Redis cache, tiered cache, main wiring, tests, and reference scan for RedisCache call sites.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Read: api/cache/redis.go` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: api/cache/tiered.go` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read: api/main.go` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read: api/cache/redis_test.go` | -1 | unknown (coerced from string) | 0ms |
| 5 | `rg references: NewRedisCache in main and tests; RedisCache in tiered cache and tests` | -1 | unknown (coerced from string) | 0ms |

## Deviations

GitNexus impact on Go symbols could not resolve `RedisCache`/`NewRedisCache` despite fresh `gitnexus analyze --force`; reference scan was used to enumerate direct call sites, and `gitnexus_detect_changes` will still be run before commit.

## Known Issues

Default key preservation and explicit model-aware namespace are in tension. The design preserves current default key namespace `v2` unless new namespace env vars are set; explicit model/revision/tokenizer/chunking fields are opt-in to avoid invalidating existing caches unexpectedly.

## Files Created/Modified

- `api/cache/redis.go`
- `api/main.go`
- `benchmark.py`
