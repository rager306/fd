---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Refactored LocalCache to a single mutex-owned map with idempotent Close.

Replace `sync.Map` plus separate size counter with a mutex-owned map. Make `currentSize()` derive from map length. Add an idempotent `Close() error` that stops the eviction loop. Keep `Get`, `Set`, and `Delete` API compatible. Ensure max-size enforcement retains the just-written key and handles unbounded mode.

## Inputs

- `api/cache/local.go`

## Expected Output

- `api/cache/local.go`

## Verification

cd api && go test ./cache && cd api && go test -race ./cache -run TestLocalCache

## Observability Impact

Lifecycle surface is explicit via Close; tests verify no race in LocalCache paths.
