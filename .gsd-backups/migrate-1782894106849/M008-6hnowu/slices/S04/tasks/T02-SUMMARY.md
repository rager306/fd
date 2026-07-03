---
id: T02
parent: S04
milestone: M008-6hnowu
key_files:
  - api/cache/redis.go
  - api/cache/tiered.go
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:05:45.094Z
blocker_discovered: false
---

# T02: Dense binary Redis layout is good; next spike should add configurable long retention, model/version-aware keys, and explicit memory/eviction policy.

**Dense binary Redis layout is good; next spike should add configurable long retention, model/version-aware keys, and explicit memory/eviction policy.**

## What Happened

Assessed Redis embedding data layout and long-retention policy. fd already uses an efficient binary string value format: `[dim:uint16][float32*dim]`, producing about 4098 bytes for 1024d instead of JSON-sized payloads. This should remain the default layout for dense vectors. Compression is not a first move because 4KB values are moderate and compression adds CPU to every cache hit; test only if Redis memory becomes bottleneck. For research/chunking, Redis L2 should become configurable long-lived cache with TTL/no-expire modes and model-aware namespace. Correctness-affecting key fields should include model id/revision, tokenizer version/hash, embedding/cache schema version, normalization/pooling version, dimension, and chunking version when caching chunks. Redis docs support explicit `maxmemory` and eviction policy selection; for long-lived reusable embedding cache, `allkeys-lfu` is attractive when repeated chunks are frequently reused, while `allkeys-lru` is simpler and good default for power-law access. `volatile-*` policies are risky if no-expire mode is used because no-expire keys may become non-evictable. RDB persistence is a good first persistence mode for rebuildable embedding cache because it is compact and faster to restart; AOF is more durable but larger and may add write overhead, so reserve it for cases where re-embedding loss is unacceptable.

## Verification

Read Redis eviction/persistence docs, current fd Redis/cache code, and mapped env/correctness fields.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://redis.io/docs/latest/develop/reference/eviction/` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Fetched: https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read: api/cache/redis.go` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read: api/cache/tiered.go` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Current `RedisCache` has a fixed 24h TTL and cache key lacks explicit model revision/tokenizer/chunking version. That is acceptable for the current simple service but insufficient for long-running research/chunk reuse where stale vectors from changed model/chunking settings could be reused silently.

## Files Created/Modified

- `api/cache/redis.go`
- `api/cache/tiered.go`
