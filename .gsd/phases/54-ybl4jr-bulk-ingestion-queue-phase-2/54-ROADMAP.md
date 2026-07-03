# M054-ybl4jr: Bulk ingestion queue (Phase 2)

**Vision:** Phase 2 из Issue #9: добавить async bulk ingestion path для burst scenarios. Optional feature (FD_QUEUE_ENABLED), не заменяет sync /v1/embeddings path.

Структура:
- HTTP endpoint POST /v1/queue принимает embedding request, валидирует, записывает в in-memory bounded channel, возвращает 202 Accepted + request_id.
- Background worker pulls items из channel, формирует batches (time-windowed до FD_QUEUE_BATCH_WINDOW_MS ms), заполняет до 32 inputs cap, делает один TEI batch call. Результаты (включая ошибки) сохраняются в in-memory result store keyed by request_id.
- GET /v1/queue/{id} возвращает статус (pending|completed|failed) + результат.
- Backpressure: bounded channel возвращает 503 service_unavailable если queue full (новый error code queue_full).
- Конфигурация через env vars: FD_QUEUE_ENABLED, FD_QUEUE_MAX_SIZE, FD_QUEUE_BATCH_WINDOW_MS, FD_QUEUE_BATCH_MAX_SIZE, FD_QUEUE_WORKER_COUNT, FD_QUEUE_RESULT_TTL.
- /health queue_depth gauge добавляется в Phase 0 metrics (queue depth + drain rate per Issue #9).
- sync /v1/embeddings остаётся unchanged.
- /health, /metrics, /v1/healthcheck, /warmup endpoints незатронуты.

Out of scope: persistence beyond in-memory (no Redis queue, no file persistence) — graceful degrade means results lost on fd-api restart, documented.

## Slices

- [x] **S01: Bulk ingestion queue MVP** `risk:medium` `depends:[]`
  > After this: POST /v1/queue -d '{"inputs":[...]}' → 202 Accepted с X-Request-Id. Позже GET /v1/queue/<id> → 200 {status:completed, data:[...]}. Batch visibility в /metrics.

## Boundary Map

Not provided.
