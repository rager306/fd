---
estimated_steps: 1
estimated_files: 6
skills_used: []
---

# T02: Implemented batched cache peeks for `/v1/embeddings` and Redis MGET support.

Extend the embedding cache interface with `GetManyIfPresent`, implement it for TieredCache and Redis MGET, update handlers and mocks, and preserve L1/L2 behavior, ordering, dimension checks, and L1 backfill for Redis hits.

## Inputs

- `api/handlers/embeddings.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`

## Expected Output

- `api/handlers/embeddings.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`

## Verification

cd api && go test ./handlers ./cache && cd api && go test ./...

## Observability Impact

Tests prove batched peek without runtime logging changes.
