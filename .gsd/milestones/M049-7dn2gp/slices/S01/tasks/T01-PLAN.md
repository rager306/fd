---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T01: Pinned cache invalidation gaps with red tests.

Add focused failing tests for Redis/Local/Tiered cache Delete/Flush behavior and authenticated HTTP cache invalidation route expectations. Tests should prove namespace safety and protected route behavior before implementation.

## Inputs

- `api/cache/local.go`
- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/main.go`
- `documents/issue-8-current-m049.md`

## Expected Output

- `api/cache/local_test.go`
- `api/cache/redis_test.go`
- `api/cache/tiered_test.go`
- `api/handlers/cache_test.go`

## Verification

cd api && go test ./cache ./handlers (expected red before implementation).

## Observability Impact

Red tests pin the new operator-facing invalidation contract.
