---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T02: Implemented namespace-safe cache invalidation primitives and HTTP routes.

Implement Delete/Flush on active cache types, keeping Redis operations namespace-scoped. Add an authenticated handler/route for cache flush and optionally key deletion if a safe key format is available. Preserve existing cache HIT/MISS behavior.

## Inputs

- `api/cache/hash.go`
- `api/cache/keys.go`
- `api/middleware/auth.go`

## Expected Output

- `api/cache/local.go`
- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/handlers/cache.go`
- `api/main.go`

## Verification

cd api && go test ./cache ./handlers && cd api && go test ./...

## Observability Impact

Cache state can now be reset through a first-class authenticated action.
