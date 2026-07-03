---
id: T03
parent: S01
milestone: M052-mmf99p
key_files:
  - api/handlers/embeddings.go
  - api/cache/tiered.go
  - api/observability/metrics.go
  - api/main.go
key_decisions:
  - ObserveBatchFill ratio calc: float64(n)/32.0 — целое деление на 64-битное даёт [0,1] дробный результат
  - Request duration buckets expanded additively — existing Grafana queries continue working
  - Cache lookup timer through defer в самый конец GetManyIfPresent — captures both L1+L2 latency в одном оbservation
  - LookupDurationFn через отдельный Set — не смешиваю с CacheObserver чтобы не замедлить hot path структурированием timer
duration: 
verification_result: passed
completed_at: 2026-07-03T09:09:07.379Z
blocker_discovered: false
---

# T03: Batch-fill ratio + cache lookup timer + expanded request duration buckets: observeBatchFill в handler, SetLookupDurationObserver в TieredCache, fd_request_duration_seconds 7 buckets. Все 10 тестов зелёные с -race.

**Batch-fill ratio + cache lookup timer + expanded request duration buckets: observeBatchFill в handler, SetLookupDurationObserver в TieredCache, fd_request_duration_seconds 7 buckets. Все 10 тестов зелёные с -race.**

## What Happened

T03 wraps up batch-fill ratio + cache lookup timer:

1. **Batch-fill ratio** wired in handlers/embeddings.go:
```go
if tei, ok := h.teiClient.(interface{ ObserveBatchFill(int) }); ok {
    tei.ObserveBatchFill(len(missTexts))
}
```
Не изменяет embed.Embedder interface (retro compatibility). Через optional cast на number of cache misses per TEI batch.

2. **Cache lookup timer** в TieredCache.GetManyIfPresent:
   - `started := time.Now()` + `defer tc.recordLookupDuration(time.Since(started))`
   - SetLookupDurationObserver(fn) добавлен в TieredCache + wiring in main.go до metrics callback

3. **Request duration buckets** expanded from [0.05,0.1,0.5,1.0] to [0.005,0.01,0.05,0.1,0.5,1.0,5.0]:
   - 5ms bucket → hot path resolution
   - 5s bucket → cold-miss visibility
   - Added comment: "Covers both cache-hot (1-10ms) and cold-miss (100ms-5s) paths"
   - Backward compatible — existing metrics scrapers see more buckets, no breaking change

All 10 packages green with -race.

## Verification

go test ./api/... -race -count=1 — 10 packages green. Batch fill: handler вызывает teiClient.ObserveBatchFill(len(missTexts)). Cache lookup timer в TieredCache.GetManyIfPresent через started/defer. Request duration buckets expanded from 4 to 7 (0.005–5.0).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build && go vet ./...` | 0 | pass | 2500ms |
| 2 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13000ms |

## Deviations

None.

## Known Issues

cache_lookup_duration histogram registered in T02 но таймер срабатывает только в GetManyIfPresent; GetOrLoad timer отдельно не заводится (single-key path использует тот же GetManyIfPresent через GetIfPresent прокси). Это ОК для Phase 0 — priority на batch path.

## Files Created/Modified

- `api/handlers/embeddings.go`
- `api/cache/tiered.go`
- `api/observability/metrics.go`
- `api/main.go`
