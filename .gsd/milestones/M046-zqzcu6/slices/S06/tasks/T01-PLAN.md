---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Added red test proving `/v1/embeddings` does not yet use batched cache peeks.

Add tests proving `/v1/embeddings` chunk cache-peek uses a batched cache API rather than per-item `GetIfPresent`, and that hits/misses preserve output order. Add TieredCache/Redis unit tests where feasible for batched lookup shape.

## Inputs

- `api/handlers/embeddings.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`

## Expected Output

- `api/handlers/embeddings_integration_test.go`
- `api/cache/tiered_test.go`
- `api/cache/redis_test.go`

## Verification

cd api && go test ./handlers -run TestCreateEmbeddingUsesBatchedCachePeek && cd api && go test ./cache -run 'Test(TieredCache|RedisCache).*Many'

## Observability Impact

Red tests capture the residual P1 #6 contract.
