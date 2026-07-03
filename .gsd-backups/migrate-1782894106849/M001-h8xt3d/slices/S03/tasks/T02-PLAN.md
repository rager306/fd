---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Implement LocalCache semantics

Change LocalCache.Set to overwrite value/TTL for existing keys, enforce maxSize for new keys, and add tests for overwrite and capacity behavior.

## Inputs

- `S03 T01`

## Expected Output

- `api/cache/local.go`
- `api/cache/local_test.go`

## Verification

cd api && go test ./cache -run 'TestLocalCache' -count=1

## Observability Impact

Makes configured L1 cache capacity meaningful.
