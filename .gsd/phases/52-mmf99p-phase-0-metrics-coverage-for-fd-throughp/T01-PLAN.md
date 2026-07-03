---
estimated_steps: 14
estimated_files: 6
skills_used: []
---

# T01: Cache hit-rate by tier + L2 occupancy

Изменить/расширить metrics.go:
1. Существующий `fd_cache_hits_total{result}` -> добавить label `tier` со значениями "l1" | "l2" | "miss". Backward compat: result values остаются ("hit" для L1+L2 combined, "miss"), плюс новый лейбл-tier. В initLabelSeries проинициализировать ВСЕ комбинации (result x tier).
2. Существующий `fd_cache_entries{tier}` -> добавить "l2" series. L2 size нужно от TieredCache: добавить public `LocalSize()` (уже есть) и new `RedisSize() error` method on TieredCache через `RedisCache.Size()`. RedisCache.Size() через `redis-cli DBSIZE` (только если используется Redis). DBSIZE возвращает total keys in selected DB; namespace isolation через prefix/keys уже устроена через EMBEDDING_CACHE_VERSION в compose (см. .env), но DBSIZE для standalone namespace это O(1).
3. Новый gauge `fd_cache_memory_bytes{tier}` — приблизительная память (L1: `LocalSize() * 4096` где 4096 = sizeof(float32) * 1024 + overhead; L2: `RedisSize() * 4096`). Документировать в metric Help как "approximate, assumes 1024-dim float32 embeddings".
4. Instrument `cache/tiered.go` `GetManyIfPresent` и `GetOrLoad`:
   - L1 hit path -> `metrics.ObserveCacheResult(tier="l1", result="hit")`
   - L1 miss + L2 hit -> `metrics.ObserveCacheResult(tier="l2", result="hit")`
   - Both miss -> `metrics.ObserveCacheResult(tier="miss", result="miss")`
5. TieredCache constructor получает metrics observer через interface: добавить `observability.Recorder` interface (sub-package of observability, чтобы избежать cycle: cache → observability → cache; либо передавать callbacks). Простой подход: `TieredCache.SetMetricsObserver(func(tier, result string))` - опциональная, может быть nil.
6. main.go вызывает `tiered.SetMetricsObserver(metrics.ObserveCacheTierResult)` в wiring stage.
7. tests: существующие TestMetricsHandlerExposesPrometheusText - дополнить проверкой нового `tier` label в output; новые unit tests в cache/tiered_test.go проверяют observer callback receive correct tier value.

Backward compat: existing `fd_cache_hits_total` series сохраняется (result-only), НЕ breaking. Series с tier-extra появится после первого hit.

Files: api/observability/metrics.go, api/cache/tiered.go, api/cache/redis.go (add Size method), api/cache/local.go (verify Size already exposed), api/main.go (wire observer), tests metrics_test.go + tiered_test.go (new or extend existing).

Verify: go build ./... ; go test ./... -race -count=1 ; manually scrape localhost:8000/metrics после smoke /v1/embeddings.

## Inputs

- `api/observability/metrics.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/cache/local.go`
- `api/main.go`

## Expected Output

- `api/observability/metrics.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/main.go`
- `api/observability/metrics_test.go`

## Verification

go test ./api/observability/... ./api/cache/... -race -count=1 (existing + new metrics tests green); manual /metrics scrape returns fd_cache_hits_total{tier=...} after cache hit/miss.

## Observability Impact

fd_cache_hits_total получает tier label (additive). fd_cache_entries получает l2 series. fd_cache_memory_bytes новый gauge. Остальные metrics unchanged.
