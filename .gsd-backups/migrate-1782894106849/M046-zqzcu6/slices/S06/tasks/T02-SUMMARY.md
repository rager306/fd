---
id: T02
parent: S06
milestone: M046-zqzcu6
key_files:
  - api/handlers/embeddings.go
  - api/handlers/embeddings_integration_test.go
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/cache/tiered_test.go
  - api/cache/lru.go
  - api/fd_v2_lifecycle_integration_test.go
key_decisions:
  - Use `GetManyIfPresent` keyed by input index so duplicate texts and response ordering remain unambiguous.
duration: 
verification_result: passed
completed_at: 2026-06-15T05:48:25.885Z
blocker_discovered: false
---

# T02: Implemented batched cache peeks for `/v1/embeddings` and Redis MGET support.

**Implemented batched cache peeks for `/v1/embeddings` and Redis MGET support.**

## What Happened

Extended the embedding cache surface with `GetManyIfPresent`. `/v1/embeddings` now calls that once per bounded chunk and no longer loops through per-item `GetIfPresent` in `collectCacheMisses`. `TieredCache.GetManyIfPresent` checks L1 for all indexes, batches L2 misses through `RedisCache.GetMany` using one Redis MGET, backfills L1 for Redis hits, and returns hits keyed by original input index. `RedisCache.GetMany` decodes MGET bulk values and omits misses/malformed/dimension-mismatched entries. Compatibility methods were added for LRUCache and lifecycle test cache.

## Verification

The red handler test now passes; cache-level batched lookup tests pass; `cd api && go test ./handlers ./cache` passes; `cd api && go test ./...` passes; static proof `43c16c32-c290-499a-a42a-b8602a0ce6ee` confirms `/v1/embeddings` uses batched cache peeks and Redis MGET.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers -run TestCreateEmbeddingUsesBatchedCachePeek` | 0 | ✅ pass | 28000ms |
| 2 | `cd api && go test ./cache -run 'Test(TieredCacheGetManyIfPresentUsesLocalHitsByIndex|RedisBulkBytes)'` | 0 | ✅ pass | 9000ms |
| 3 | `cd api && go test ./handlers ./cache && go test ./...` | 0 | ✅ pass | 24000ms |
| 4 | `gsd_exec 43c16c32-c290-499a-a42a-b8602a0ce6ee` | 0 | ✅ pass | 106ms |

## Deviations

No live Redis MGET integration test was added because existing Redis tests avoid requiring a live server; static proof plus cache/handler unit tests cover the code path shape.

## Known Issues

Closure matrix and final gates remain in T03/T04.

## Files Created/Modified

- `api/handlers/embeddings.go`
- `api/handlers/embeddings_integration_test.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/cache/tiered_test.go`
- `api/cache/lru.go`
- `api/fd_v2_lifecycle_integration_test.go`
