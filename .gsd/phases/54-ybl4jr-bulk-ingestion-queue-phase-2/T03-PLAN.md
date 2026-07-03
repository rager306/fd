---
estimated_steps: 17
estimated_files: 3
skills_used: []
---

# T03: Time-windowed batching + backpressure

Добавить time-windowed batching в worker: вместо single-item process, собирать items из канала в batch до FD_QUEUE_BATCH_MAX_SIZE (default 32) или timeout FD_QUEUE_BATCH_WINDOW_MS (default 10ms). Один TEI call на batch.

1. **Worker upgrade**: 
   - `StartQueueWorker` теперь использует time.Ticker (window_interval) для triggered batch drain
   - На tick: drain up to maxBatchSize items из канала (try-without-block), call `embedder.Embed(ctx, batch_texts)` one TEI batch call
   - Split results back to individual items (by index match), save to store
   - Обработка ошибок: если TEI batch call fails → все items в batch marked as failed

2. **Backpressure**:
   - `POST /v1/queue` handler проверяет: если `len(channel) >= queueCap` (full) → 503 queue_full error
   - Retry-After header: 5 секунд (впоследствии через drain-rate auto в metrics)

3. **Graceful shutdown**: 
   - Worker захватывает ctx.Done и при выходе сбрасывает все pending items в канале как failed с сообщением "queue worker shutdown"
   - Опционально: при shutdown, wait up to 5 секунд для drain оставшихся items (через sync.WaitGate или канал) перед marking failed

4. **Tests (расширение worker_test.go)**:
   - TestWorkerBatchesItems — 4 items + tick after 10ms → 1 TEI batch with 4 inputs (mock TEI captures batch size)
   - TestWorkerBatchMaxWindow — 20 items submit fast → worker drains до 16 items в первую волну (maxBatchSize=16 default) + 4 во вторую, проверяет что каждая волна имеет ровно batch_size как observed
   - TestWorkerBackpressureChannelFull — submit 1025 items сразу → первые 1024 accepted, последний → ошибка 503

**files**: api/queue/worker.go (upgrade), api/handlers/queue_handlers.go (backpressure logic), api/queue/worker_test.go (extend)

## Inputs

- `api/queue/worker.go (T02)`
- `api/handlers/queue_handlers.go (T01)`

## Expected Output

- `api/queue/worker.go`
- `api/handlers/queue_handlers.go`
- `api/queue/worker_test.go`

## Verification

go test ./api/queue/... -race (extended worker tests). go test ./api/... -race full suite green.

## Observability Impact

Backpressure surge на channel full.
