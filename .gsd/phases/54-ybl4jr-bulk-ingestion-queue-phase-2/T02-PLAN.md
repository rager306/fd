---
estimated_steps: 17
estimated_files: 3
skills_used: []
---

# T02: Worker goroutine + single-pass TEI bridge

Реализовать worker goroutine без time-windowing (single-pass): drain one item from channel per tick, call TEI directly via Embedder, сохранить результат в store.

1. **worker.go**: `StartQueueWorker(ctx, channel, store, embedder, cache, logger, batchSize)` где `channel <-chan QueueItem`, batchSize = FD_QUEUE_BATCH_MAX_SIZE (default 32)
   - for loop over channel pulls
   - Вызов embedder.Embed(ctx, queueItems.texts) — single call per item (без batching пока)
   - Сохранение результата (embeddings или error) в store
   - Drain count в metrics (пока nil — метрики в T04)

2. **Подключение в main.go**: запуск worker goroutine
   - bounded channel size = FD_QUEUE_MAX_SIZE (default 1024) через make(chan QueueItem, size)
   - Worker context: привязан к глобальному shutdown (через ctx, cancel)

3. **Worker lifecycle**:
   - On shutdown: cancel context → for loop exits when channel drain no longer possible
   - Log: "queue worker stopped" с processed count

4. **Тесты** (api/queue/worker_test.go):
   - TestWorkerProcessesSingleItem — submit 1 item, verify result in store
   - TestWorkerStopsOnContextCancel — cancel ctx → worker stops after drain window
   - TestWorkerHandlesEmbedError — TEI returns error → store has status=failed, error present

**files**: api/queue/worker.go, api/queue/worker_test.go, api/main.go (wiring)

## Inputs

- `api/queue/types.go`
- `api/queue/store.go`
- `api/embed/embed.go`
- `api/cache/tiered.go`

## Expected Output

- `api/queue/worker.go`
- `api/queue/worker_test.go`
- `api/main.go (worker start)`

## Verification

go test ./api/queue/... -race (3 worker tests + existing store tests). go test ./api/... -race full suite.

## Observability Impact

Worker запускается на старте, но метрики nil пока (T04).
