---
id: T01
parent: S01
milestone: M052-mmf99p
key_files:
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/observability/metrics.go
  - api/observability/metrics_test.go
  - api/main.go
  - api/fd_v2_cache_integration_test.go
key_decisions:
  - D053: additive tier label strategy — backward compatible, existing dashboards continue
  - Observer pattern as simple CacheObserver callback in cache package — avoids circular import
  - RedisCache.Size по SCAN — approximate, документирован как O(N) costs at scrape cadence
  - approxEmbeddingBytes=4096 per entry — документирован как approximate, not billing-accurate
duration: 
verification_result: passed
completed_at: 2026-07-03T09:05:22.434Z
blocker_discovered: false
---

# T01: Cache hit-rate by tier + L2 occupancy: fd_cache_hits_total{result,tier}, fd_cache_entries{tier=l1|l2}, fd_cache_memory_bytes{tier=l1|l2}, CacheObserver pattern в TieredCache, RedisCache.Size, wiring в main.go — все 10 пакетов тестов зелёные с -race

**Cache hit-rate by tier + L2 occupancy: fd_cache_hits_total{result,tier}, fd_cache_entries{tier=l1|l2}, fd_cache_memory_bytes{tier=l1|l2}, CacheObserver pattern в TieredCache, RedisCache.Size, wiring в main.go — все 10 пакетов тестов зелёные с -race**

## What Happened

T01 завершён в один проход:

1. **RedisCache.Size(ctx)** — добавил SCAN-based подсчёт ключей namespace через Iterator, документированный как approximate.
2. **TieredCache.CacheObserver** — callback pattern без circular imports. Define в cache пакете, observer вызывается из GetManyIfPresent/GetOrLoad/GetIfPresent. Set-ор позволяет nil (no-op).
3. **GetManyIfPresent обогащён**: L1 hit → observe("l1",true), L2 hit → observe("l2",true), miss → observe("miss",false).
4. **GetOrLoad обогащён**: L1 hit/L2 hit/miss → observer calls в каждом пути.
5. **Metrics обновлён**:
   - `fd_cache_hits_total{result,tier}` с tier labels l1|l2|miss|all (all для legacy backward compat)
   - `fd_cache_entries{tier}` расширен l2
   - `fd_cache_memory_bytes{tier=l1|l2}` новый gauge, approxEmbeddingBytes=4096
   - `ObserveCacheResultWithTier(result, tier)`
   - `SetRedisCacheSizeObserver(fn)` callback
6. **Wiring в main.go**: SetCacheObserver + SetRedisCacheSizeObserver
7. **Тест-обновление**: metrics_test.go TestMetricsModelLoadedAndCacheResult обновлён на новый label; fd_v2_cache_integration_test.go обновлён на проверку tier="all".
8. **Full suite**: все 10 пакетов зелёные с -race.

## Verification

go test ./api/... -race -count=1 — все 10 пакетов зелёные (включая исправленный интеграционный тест). Build + vet чистые. Observer pattern in cache/tiered.go + RedisCache.Size + wiring в main.go (SetCacheObserver + SetRedisCacheSizeObserver).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build ./... && go vet ./...` | 0 | pass | 3200ms |
| 2 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 14000ms |

## Deviations

Интеграционный тест TestFdV2CacheMissThenHit не использует настоящий TieredCache (mock localVectorCache) — проверяю tier="all" вместо tier="l1". Это приемлемо — тест проверяет hit counter, не tier-specific behaviour.

## Known Issues

Redis Size через SCAN O(N) по namespace — если namespace достигнет миллионов ключей, SCAN на /metrics scrape может занять секунды. Rate-limit deferred до actual observation in production metrics.

## Files Created/Modified

- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/observability/metrics.go`
- `api/observability/metrics_test.go`
- `api/main.go`
- `api/fd_v2_cache_integration_test.go`
