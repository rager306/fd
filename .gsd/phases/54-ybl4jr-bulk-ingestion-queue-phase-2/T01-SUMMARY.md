---
id: T01
parent: S01
milestone: M054-ybl4jr
key_files:
  - api/queue/types.go
  - api/queue/id.go
  - api/queue/store.go
  - api/handlers/queue_handlers.go
  - api/main.go
key_decisions:
  - FD_QUEUE_ENABLED default false — backward compatible; production без env override видит 404 на /v1/queue
  - Backpressure через select default — non-blocking send, immediate 503 if channel full + Retry-After header
  - ResultStore TTL eviction с background goroutine, lock-based map для simplicity
  - Request ID = crypto/rand 64 bits — negligible collision risk
duration: 
verification_result: passed
completed_at: 2026-07-03T10:59:33.736Z
blocker_discovered: false
---

# T01: Queue core: api/queue/types.go+id.go+store.go (sync.Mutex map с TTL eviction), api/handlers/queue_handlers.go (POST /v1/queue + GET /v1/queue/:id), main.go wiring через FD_QUEUE_ENABLED. 10 пакетов тестов зелёные. Worker пока заглушка.

**Queue core: api/queue/types.go+id.go+store.go (sync.Mutex map с TTL eviction), api/handlers/queue_handlers.go (POST /v1/queue + GET /v1/queue/:id), main.go wiring через FD_QUEUE_ENABLED. 10 пакетов тестов зелёные. Worker пока заглушка.**

## What Happened

T01 закрывает основу queue:

1. **api/queue/types.go** — Status (pending|completed|failed), Result, Item типы, ErrQueueDisabled/ErrQueueFull sentinel errors
2. **api/queue/id.go** — NewRequestID() через crypto/rand, 16 hex chars (64 bits entropy, negligible collisions)
3. **api/queue/store.go** — ResultStore in-memory sync.Mutex-guarded map с TTL eviction goroutine (5min default), Save/Get/Size/Close methods
4. **api/handlers/queue_handlers.go** — Submit (POST /v1/queue) валидирует body, генерирует request_id, non-blocking send в channel с backpressure (default = 503 queue_full + Retry-After: 5). Poll (GET /v1/queue/:id) возвращает pending/completed/failed по ResultStore.
5. **main.go wiring** — FD_QUEUE_ENABLED feature gate (default false), bounded channel FD_QUEUE_MAX_SIZE (default 1024), defer Close для ResultStore. Worker НЕ запущен — TODO T02.

T01 не делает ещё worker, поэтому submit кладёт items в channel, но Poll всегда вернёт pending. Это ожидаемо для этого таска — contract-only.

## Verification

go build ./api/... (exit 0). go test ./api/... -race — 10 packages green. FD_QUEUE_ENABLED=false (default) - feature gated off. Worker TODO T02.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build ./...` | 0 | pass | 2500ms |
| 2 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13500ms |

## Deviations

None.

## Known Issues

Worker не запущен (T02), поэтому submit видит pending до timeout. Это expected для contract-only T01.

## Files Created/Modified

- `api/queue/types.go`
- `api/queue/id.go`
- `api/queue/store.go`
- `api/handlers/queue_handlers.go`
- `api/main.go`
