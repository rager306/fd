# M046 S05 LocalCache Correctness

Captured: 2026-06-14

## Scope

Fix issue #3 P1 #10: `LocalCache` used `sync.Map`, a separate mutex-protected size counter, and an always-running eviction goroutine. That design made overwrite/delete/expiry accounting hard to reason about under concurrency and had no lifecycle stop surface.

## Red Evidence

```bash
cd api && go test ./cache -run 'TestLocalCache_(CloseIsIdempotent|ConcurrentOverwriteKeepsSingleEntry)'
```

Initial result:

```text
cache/local_test.go: c.Close undefined (type *LocalCache has no field or method Close)
```

## Changes

- Replaced `sync.Map` plus separate `size` counter with a single mutex-owned `map[string]l1Entry`.
- `currentSize()` now derives from `len(data)` under the same lock.
- `Set` overwrites in place and enforces max size while protecting the just-written key.
- `Get` lazily expires stale entries under the same lock.
- Added `Close() error`, backed by `sync.Once`, `stopCh`, and `doneCh`, to stop the background eviction loop.
- `NewLocalCache` no longer starts a ticker when `evictTTL <= 0`; lazy expiry still works through `Get`.
- API shutdown/error paths now close the local cache along with Redis.

## Verification

### Targeted LocalCache tests

```bash
cd api && go test ./cache -run 'TestLocalCache_(CloseIsIdempotent|ConcurrentOverwriteKeepsSingleEntry|SetAndGet|TTLExpired|SetRefreshesExistingValueAndTTL|EnforcesMaxSize|Delete|NotFound)'
```

Result:

```text
ok fd-api/cache
```

### Cache package

```bash
cd api && go test ./cache
```

Result:

```text
Go test: 44 passed in 1 packages
```

### Race-enabled LocalCache tests

```bash
cd api && go test -race ./cache -run TestLocalCache
```

Result:

```text
Go test: 9 passed in 1 packages
```

### Full Go suite

```bash
cd api && go test ./...
```

Result:

```text
Go test: 281 passed in 9 packages
```

### Lint

```bash
cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
```

Result:

```text
0 issues.
```

### Vulnerability scan

```bash
cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...
```

Result:

```text
No vulnerabilities found.
Your code is affected by 0 vulnerabilities.
```

### Static LocalCache proof

Evidence: `gsd_exec:f124000a-5996-4c68-888d-1e31237c6d39`

```text
PASS LocalCache uses lock-owned map, derived size, Close, and main shutdown integration
```

## Findings Closed

- P1 #10: `LocalCache` no longer has independent size-counter drift risk and now has an explicit lifecycle stop surface.

## Still Open

- S06: residual P1 #6 and P2/P3 closure matrix.
