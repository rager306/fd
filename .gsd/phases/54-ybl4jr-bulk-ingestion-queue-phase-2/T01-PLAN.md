---
estimated_steps: 14
estimated_files: 5
skills_used: []
---

# T01: Async queue endpoint + types + result store

Реализовать основу queue: типы, handler, result store, status polling endpoint.

1. **Новый пакет `api/queue`**:
   - types.go: `QueueItem` (request_ctx, texts, response_ch), `QueueResult` (status, embeddings, error, prompt_tokens), `QueueStatus` type с константами pending/completed/failed
   - store.go: `ResultStore` — in-memory sync.Map для хранения результатов keyed по request_id. TTL eviction goroutine. Методы: Save(id, result), Get(id) (*QueueResult, bool), Close() для stop eviction goroutine
   - id.go: `NewRequestID() string` — crypto/rand UUID (12 chars), minimal collision-safe
   - worker.go (stub): `StartWorker(ctx, channel, store, embedder, cache, logger, batchSize)` — пока nil, просто подключается в T02

2. **Расширение `api/handlers/queue_handlers.go`** (или queue_api.go):
   - POST /v1/queue handler: парсит request body (мини-копия validation из sync path), генерирует request_id, кладет вход в канал. Если канал full → 503 queue_full error. Возвращает 202 Accepted с X-Request-Id.
   - GET /v1/queue/:id handler: получит status из ResultStore. Если pending → 202 с status=pending. Если completed → 200 с embeddings data. Если failed → 500 с error. Если не найден → 404.

3. **Wiring в main.go**:
   - Регистрация новых endpoints с проверкой FD_QUEUE_ENABLED (envutil.BoolOrDefault, default false)
   - Создание resultStore и start queue worker

Не менять: существующий Embedder interface, cache APIs, lifecycle /health shape.

**files**: api/queue/types.go, id.go, store.go, api/handlers/queue_handlers.go, api/main.go

## Inputs

- `api/handlers/embeddings.go (validation pattern)`
- `api/main.go (wiring patterns)`
- `api/internal/envutil/bool.go (FD_QUEUE_ENABLED)`

## Expected Output

- `api/queue/types.go`
- `api/queue/id.go`
- `api/queue/store.go`
- `api/handlers/queue_handlers.go`
- `api/main.go (wiring)`

## Verification

go build ./api/... ; go test ./api/queue/... -race (unit tests for store, ID gen)

## Observability Impact

Нет production impact пока worker не подключён (T02)
