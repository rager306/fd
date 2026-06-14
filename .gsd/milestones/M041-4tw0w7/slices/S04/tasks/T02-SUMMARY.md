---
id: T02
parent: S04
milestone: M041-4tw0w7
key_files:
  - api/cache/lru.go
  - api/cache/lru_test.go
  - api/observability/metrics.go
  - api/observability/metrics_test.go
  - benchmark-results/m041-s04-t02-lru-metrics.txt
  - benchmark-results/m041-s04-t02-race.txt
  - benchmark-results/m041-s04-t02-go-test.txt
  - benchmark-results/m041-s04-t02-lint.txt
  - benchmark-results/m041-s04-t02-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T07:02:27.256Z
blocker_discovered: false
---

# T02: Added goroutine-safe in-memory LRU vector cache with TTL, env configuration, SHA256 keys, and cache metrics hooks.

**Added goroutine-safe in-memory LRU vector cache with TTL, env configuration, SHA256 keys, and cache metrics hooks.**

## What Happened

Implemented `api/cache/lru.go`: a mutex-protected LRU cache keyed by `SHA256(input + "|" + dimensions)`, storing copied `[]float32` values, with TTL expiration, max-size eviction, `FD_CACHE_SIZE` and `FD_CACHE_TTL_HOURS` env defaults (10000 entries, 24h), and metrics hooks for hit/miss and eviction. Added `api/cache/lru_test.go` covering Get/Put copies, dimension-aware keys, miss accounting, LRU eviction, TTL expiration, env configuration, stable key generation, and concurrent access. Extended `api/observability/metrics.go` with `fd_cache_evictions_total` and `ObserveCacheEviction`, and updated metrics tests to assert the new counter. Ran the cache concurrency test under the Go race detector.

## Verification

Fresh LRU/metrics evidence passed: `go test ./cache ./observability -run 'TestLRU|TestEmbeddingCacheKey|TestMetrics' -v` exit 0 and `go test -race ./cache -run TestLRUCacheConcurrentAccess -v` exit 0. Fresh M043 gate passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s04-t02-lru-metrics.txt, benchmark-results/m041-s04-t02-race.txt, benchmark-results/m041-s04-t02-go-test.txt, benchmark-results/m041-s04-t02-lint.txt, benchmark-results/m041-s04-t02-govulncheck.txt. GitNexus detect_changes reports LOW risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./cache ./observability -run 'TestLRU|TestEmbeddingCacheKey|TestMetrics' -v` | 0 | ✅ pass: LRU cache and metrics unit tests passed | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test -race ./cache -run TestLRUCacheConcurrentAccess -v` | 0 | ✅ pass: concurrent LRU access race check passed | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 5 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

This task implements the LRU primitive and metrics hooks but does not wire it into `/v1/embeddings`; route integration and `X-Cache` behavior are left to the remaining S04 tasks.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T02 changes are not committed yet.

## Files Created/Modified

- `api/cache/lru.go`
- `api/cache/lru_test.go`
- `api/observability/metrics.go`
- `api/observability/metrics_test.go`
- `benchmark-results/m041-s04-t02-lru-metrics.txt`
- `benchmark-results/m041-s04-t02-race.txt`
- `benchmark-results/m041-s04-t02-go-test.txt`
- `benchmark-results/m041-s04-t02-lint.txt`
- `benchmark-results/m041-s04-t02-govulncheck.txt`
