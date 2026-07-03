---
estimated_steps: 9
estimated_files: 3
skills_used: []
---

# T01: Add envurable FD_CACHE_MAX_SIZE and FD_CACHE_LOCAL_TTL

Изменить main.go: заменить hardcoded `cache.NewLocalCache(10000, 30*time.Second)` на env-читаемые значения.

1. В main.go добавить парсинг FD_CACHE_MAX_SIZE через `envutil.PositiveInt("FD_CACHE_MAX_SIZE", 10000)` — PositiveInt гарантирует >0, иначе возвращает fallback.
2. Добавить парсинг FD_CACHE_LOCAL_TTL через `envutil.DurationOrDefault` (новый helper в api/internal/envutil/duration.go):
   - `DurationOrDefault(key string, default d time.Duration) time.Duration` — парсит Go duration строку, fallback к default при ошибке.
3. Логировать значения: `logger.Info("cache config", "l1_max_entries", maxSize, "l1_ttl", localTTL.String())`
4. Обновить .env.example документацией для FD_CACHE_MAX_SIZE и FD_CACHE_LOCAL_TTL с defaults в комментарии.

**Важно**: обе переменные должны опциональными — деплой без .env изменений получает identical behavior (10000, 30s).

Не менять: main.go существующий wiring после NewLocalCache, TieredCache конструктор, LocalCache.API. Все остальные изменения additive.

**files**: api/main.go (lines 257 + add env parsing), api/internal/envutil/duration.go (новый), .env.example (документация).

## Inputs

- `api/main.go`
- `api/internal/envutil/int.go (pattern reference)`
- `.env.example`

## Expected Output

- `api/internal/envutil/duration.go`
- `api/main.go (cache config wiring)`
- `api/.env.example`

## Verification

go build ./... && go test ./api/internal/envutil/... -race -count=1 (existing + new tests). Verify default behavior: без env values, main.go uses defaults 10000 and 30s.

## Observability Impact

INFO log "cache config" с l1_max_entries и l1_ttl
