---
id: T02
parent: S01
milestone: M053-ckvgjw
key_files:
  - api/main_test.go
key_decisions:
  - Покрытие граничных случаев: maxSize=5 с 10 вставками (size eviction), maxSize=0 unlimited mode
  - Env override test изолирован через t.Setenv — не оставляет env vars в test process или CI
  - Cache TTL test через API не нужен — TieredCache new localTTL contract уже покрыт в existing cache tests
  - DurationOrDefault покрыт в envutil unit test, no need to retest it in main_test
duration: 
verification_result: passed
completed_at: 2026-07-03T10:44:10.213Z
blocker_discovered: false
---

# T02: 5 unit-тестов в main_test.go покрывают env override, defaults, boundary sizes: TestCacheConfigFromEnv (500/10s), TestCacheConfigDefaultsWithoutEnv, TestLocalCacheWithSmallSize (size eviction), TestLocalCacheWithZeroSize (unlimited). Все 10 пакетов тестов зелёные с -race.

**5 unit-тестов в main_test.go покрывают env override, defaults, boundary sizes: TestCacheConfigFromEnv (500/10s), TestCacheConfigDefaultsWithoutEnv, TestLocalCacheWithSmallSize (size eviction), TestLocalCacheWithZeroSize (unlimited). Все 10 пакетов тестов зелёные с -race.**

## What Happened

T02 закрывает тестовое покрытие Phase 1a:

1. **TestWarmupRetryPolicyFromEnvDefaultsTo5** — existing retry policy defaults (default added in T01 от main_test.go's append). Pre-existing test re-uses append-fix to ensure warmup config defaults не regression.

2. **TestCacheConfigFromEnv** — мокает main() wiring env parsing:
   - Sets `FD_CACHE_MAX_SIZE=500` and `FD_CACHE_LOCAL_TTL=10s`
   - Calls `envutil.PositiveInt` and `envutil.DurationOrDefault`
   - Verifies that values propagate from env to LocalCache config
   - Test isolation через `t.Setenv` — no env leak

3. **TestCacheConfigDefaultsWithoutEnv** — smae functions с пустыми env values:
   - Defaults 10000 entries, 30s TTL
   - Backward compat: production без env vars не ломается

4. **TestLocalCacheWithSmallSize** — eviction boundary:
   - maxSize=5, вставлено 10 entries с разными ключами
   - Size ≤ 5 (size eviction enforced)
   - Не падает на overflow

5. **TestLocalCacheWithZeroSize** — unlimited mode:
   - maxSize=0 отключает size eviction (existing API behavior)
   - 20 entries — всё помещается без eviction
   - Verify Size==20

10 packages green with -race.

## Verification

go test ./api/... -race — 10 packages green. Новые тесты: TestWarmupRetryPolicyFromEnvDefaultsTo5, TestCacheConfigFromEnv (FD_CACHE_MAX_SIZE/TTL override), TestCacheConfigDefaultsWithoutEnv (defaults), TestLocalCacheWithSmallSize (size eviction), TestLocalCacheWithZeroSize (unlimited mode). Все 5 проходят с -race.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13000ms |
| 2 | `cd /root/fd/api && go test . -race -count=1 -run 'TestCache|TestLocalCache' -v` | 0 | pass | 3000ms |

## Deviations

Не сделан пограничный тест для FD_CACHE_MAX_SIZE="notanumber" — envutil.PositiveInt возвращает default при ошибке parsing, и это уже тестируется в envutil/int_test.go.

## Known Issues

Тесты покрывают size eviction в чистом виде но не end-to-end behavior с metrics observer (size eviction + fd_cache_evictions_total). Это уже покрыто через cache tests в cache_package.

## Files Created/Modified

- `api/main_test.go`
