---
estimated_steps: 13
estimated_files: 2
skills_used: []
---

# T02: Unit tests — env parsing + cache boundary sizes

Добавить тесты для env parsing и boundary case behavior:

1. **TestNewLocalCacheBoundarySizes**: Вызвать `cache.NewLocalCache(0, 30*time.Second)` — размер 0 означает size eviction disabled (существующее поведение, проверить что не падает при large inserts). `-1` — NewLocalCache не возвращает ошибку (не проверяет), LocalCache допускает любые значения.

2. **TestMainCacheCapacityFromEnv**: Смоделировать main.go wiring: mock `os.Setenv("FD_CACHE_MAX_SIZE", "500")`, вызвать envutil.PositiveInt, передать в NewLocalCache → LocalSize() должно быть ≤ 500 при заполнении.

3. **TestMainCacheTTLFromEnv**: Тест DurationOrDefault парсинг:
   - "60s" → 1 min ✓
   - "" → default ✓
   - "30" → invalid, fallback to default ✓
   - "notaduration" → fallback to default ✓
   - "0s" → accepted (LocalCache handles zero TTL ✓)

4. **TestLocalCacheRejectsZeroSize** (edge): FD_CACHE_MAX_SIZE=0 → envutil.PositiveInt возвращает fallback (10000). Verify.

5. **TestLocalCacheHighLoadUnderEriction**: Заполнить 10 entries в LocalCache maxSize=5, проверить fd_cache_evictions_total growth (через metrics hook). Verify entries остаются ≤5.

Для boundary проверок использовать существующий TestCacheHandler + TestFdV2CacheMissThenHit patterns из api/fd_v2_cache_integration_test.go.

**files**: api/main_test.go (new tests), api/internal/envutil/duration_test.go (новый).

## Inputs

- `api/main.go (env vars)`
- `api/internal/envutil/duration.go (T01)`
- `api/cache/local.go (API)`

## Expected Output

- `api/main_test.go (extended)`
- `api/internal/envutil/duration_test.go`

## Verification

go test ./api/... ./api/internal/envutil/... -race -count=1 — все 10 пакетов зелёные; новые тесты все проходят.

## Observability Impact

No production impact — tests are unit-only.
