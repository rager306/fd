---
estimated_steps: 1
estimated_files: 6
skills_used: []
---

# T02: Removed dead LRUCache and unified duplicate cache hash/env helpers.

Delete dead `api/cache/lru.go` and its dedicated tests, replace the integration test scaffold with LocalCache or TieredCache active path, and unify `shortCacheKeyHash`/`shortNamespaceHash` into one package-local helper.

## Inputs

- `api/cache/local.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/fd_v2_cache_integration_test.go`

## Expected Output

- `api/fd_v2_cache_integration_test.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`

## Verification

cd api && go test ./cache && cd api && go test ./...

## Observability Impact

Reduces future agents' cache surface area.
