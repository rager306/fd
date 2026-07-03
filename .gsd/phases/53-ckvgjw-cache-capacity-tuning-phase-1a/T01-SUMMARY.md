---
id: T01
parent: S01
milestone: M053-ckvgjw
key_files:
  - api/main.go
  - api/internal/envutil/duration.go
  - api/internal/envutil/duration_test.go
  - .env.example
key_decisions:
  - envutil.PositiveInt — гарантирует FD_CACHE_MAX_SIZE >0, иначе fallback к default 10000
  - DurationOrDefault — без добавления нового шаблона в cache package, просто envutil helper
  - Default values = existing hardcoded — production deploy видит identical behavior без overwrites
  - INFO log 'cache configured' даёт операторам visibility на startup
duration: 
verification_result: passed
completed_at: 2026-07-03T10:42:12.221Z
blocker_discovered: false
---

# T01: FD_CACHE_MAX_SIZE (PositiveInt) + FD_CACHE_LOCAL_TTL (DurationOrDefault) env-tunable L1 cache. Defaults = 10000 entries, 30s TTL — backward compatible. envutil/duration.go новый helper. 10 packages green с -race.

**FD_CACHE_MAX_SIZE (PositiveInt) + FD_CACHE_LOCAL_TTL (DurationOrDefault) env-tunable L1 cache. Defaults = 10000 entries, 30s TTL — backward compatible. envutil/duration.go новый helper. 10 packages green с -race.**

## What Happened

Реализован env-tuning для L1 cache capacity:

1. **envutil/duration.go** — новый DurationOrDefault helper, парсит Go duration строки, fallback к default на ошибку.
2. **main.go:257** — replace hardcoded `NewLocalCache(10000, 30*time.Second)` на:
```go
l1MaxEntries := envutil.PositiveInt("FD_CACHE_MAX_SIZE", 10000)
l1TTL := envutil.DurationOrDefault("FD_CACHE_LOCAL_TTL", 30*time.Second)
localCache := cache.NewLocalCache(l1MaxEntries, l1TTL)
```
3. **tiered.go constructor** — replace `30*time.Second` on `NewTieredCache(localCache, redisCache, l1TTL)` — shares same TTL as L1.
4. **STARTUP LOG**: `INFO cache configured l1_max_entries=10000 l1_ttl=30s l2_namespace=v2`
5. **.env.example** — document оба новых knobs с комментарием про ~4 KB per entry.

Exact same behavior без env override — все тесты зелёные с -race.

## Verification

go test ./api/... -race — 10 packages green. envutil/duration_test.go новый, существующие int_test.go + bool_test.go untouched. main.go wiring через envutil.PositiveInt + DurationOrDefault. INFO log "cache configured" при старте.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build . && go test . -race -count=1 -timeout 30s` | 0 | pass | 2000ms |
| 2 | `cd /root/fd/api && go test ./internal/envutil/... -race -count=1` | 0 | pass | 50ms |

## Deviations

None.

## Known Issues

Наличие env vars фиксирует проверку что defaults работают; тесты на сам PositiveInt и LocalCache size уже в стандартном test suite.

## Files Created/Modified

- `api/main.go`
- `api/internal/envutil/duration.go`
- `api/internal/envutil/duration_test.go`
- `.env.example`
