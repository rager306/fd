---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Migrate cache test assertions

Migrate representative cache test assertions to Testify require/assert while preserving behavior.

## Inputs

- `api/cache/tiered_cache_test.go`

## Expected Output

- `api/cache/tiered_cache_test.go updated`

## Verification

`cd api && go test ./cache -short` passes.

## Observability Impact

Cache test failures become clearer.
