---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Added red tests for LocalCache close lifecycle and concurrent overwrite accounting.

Add tests proving `Close()` exists and is idempotent, concurrent overwrites of the same key count as one entry, capacity remains bounded, and cache package passes under race detector after implementation.

## Inputs

- `api/cache/local.go`
- `api/cache/local_test.go`

## Expected Output

- `api/cache/local_test.go`

## Verification

cd api && go test ./cache -run 'TestLocalCache_(CloseIsIdempotent|ConcurrentOverwriteKeepsSingleEntry)'

## Observability Impact

Tests become executable proof for issue #3 P1 #10.
