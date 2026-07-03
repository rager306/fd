---
id: T01
parent: S04
milestone: M008-6hnowu
key_files:
  - api/cache/redis.go
  - api/cache/tiered.go
  - api/handlers/batch.go
  - api/handlers/embeddings.go
  - benchmark.py
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T16:29:55.423Z
blocker_discovered: false
---

# T01: Added Redis throughput research: fd should first benchmark MGET/pipelined batch cache hits and pool/round-trip metrics, not generic Redis tuning.

**Added Redis throughput research: fd should first benchmark MGET/pipelined batch cache hits and pool/round-trip metrics, not generic Redis tuning.**

## What Happened

Researched Redis throughput fundamentals using Redis official pipelining/client-side caching docs, Redis benchmark docs, go-redis performance/debugging guidance, and scale writeups. For fd, the most relevant near-term finding is that cached embedding throughput is likely dominated by request/response round trips, go-redis connection behavior, and per-request encode/decode/JSON response work rather than raw Redis server capacity. Redis pipelining can reduce RTT and syscall overhead; official docs note it can substantially increase operations per second by batching commands, but replies must be bounded to avoid memory buildup. fd currently uses single `GET` and `SET` per key in `RedisCache`, while handlers iterate inputs one-by-one through `GetOrLoad`; this makes batch/cache-hit workloads a good candidate for MGET or pipelined GET, especially for `/v1/embeddings` with multiple inputs and `/embeddings/batch`. go-redis guidance says default pool size is typically 10 connections per GOMAXPROCS and tuning rarely helps unless profiling shows pool timeouts or blocked connections; current fd explicitly sets PoolSize and MinIdleConns, so benchmark should capture pool stats before changing numbers. Existing L1 cache already implements the broad client-side caching idea, so Redis server-assisted client-side caching/tracking is not a first move unless fd becomes multi-process and needs invalidation semantics beyond TTL.

## Verification

Read official Redis pipelining, client-side caching, benchmark docs, go-redis performance guidance, and current fd Redis/cache/handler/benchmark code.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Search query: Redis 2026 performance best practices pipelining MGET client side caching latency diagnosis throughput connection pooling IO threads official docs cached binary values` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Fetched: https://redis.io/docs/latest/develop/use/pipelining/` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Fetched: https://redis.io/docs/latest/develop/use/client-side-caching/` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Fetched/read sources for go-redis pool/pipeline guidance and Redis benchmark docs.` | -1 | unknown (coerced from string) | 0ms |
| 5 | `Read: api/cache/redis.go, api/cache/tiered.go, api/handlers/batch.go, api/handlers/embeddings.go, benchmark.py` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Current benchmark mostly exercises single-key cached GET via `/v1/embeddings`; it does not isolate Redis round-trip cost, MGET/pipeline behavior, pool saturation, serialization/unmarshal cost, or batch cache-hit throughput. Existing batch handler currently loops input-by-input, so batch cache hits can perform N cache lookups and N potential Redis GETs.

## Files Created/Modified

- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/handlers/batch.go`
- `api/handlers/embeddings.go`
- `benchmark.py`
