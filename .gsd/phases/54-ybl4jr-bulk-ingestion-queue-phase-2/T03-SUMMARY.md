---
id: T03
parent: S01
milestone: M054-ybl4jr
key_files:
  - api/queue/worker.go
  - api/queue/worker_test.go
  - api/main.go
  - .env.example
key_decisions:
  - Time-windowed batching через timer + non-blocking channel drain — конкурентно-дрен-ест
  - Batch целостность: при ошибке TEI все items fail (нельзя спарсить partial results правильно)
  - WorkerConfig struct pattern — тестируемый дизайн, env kv через PositiveInt+DurationOrDefault
  - drainRemaining при shutdown помечает items как failed немедленно, no silent drop
duration: 
verification_result: passed
completed_at: 2026-07-03T11:14:46.627Z
blocker_discovered: false
---

# T03: Time-windowed batching: worker собирает batch до 32 items / 10ms window, один TEI call. WorkerConfig с env vars (FD_QUEUE_BATCH_MAX_SIZE, FD_QUEUE_BATCH_WINDOW_MS). +3 new batch-специфичных теста. 11 пакетов зелёные с -race.

**Time-windowed batching: worker собирает batch до 32 items / 10ms window, один TEI call. WorkerConfig с env vars (FD_QUEUE_BATCH_MAX_SIZE, FD_QUEUE_BATCH_WINDOW_MS). +3 new batch-специфичных теста. 11 пакетов зелёные с -race.**

## What Happened

T03 закрывает time-windowed batching:

1. **worker.go** переписан: 
   - `StartQueueWorker(cfg WorkerConfig)` — BatchMaxSize (default 32), BatchWindow (default 10ms). Env-config wire.
   - `awaitFirstItem` — blocking receive первого item в batch
   - `drainUpToMax(items, batch, maxSize, windowStart, window)` — collect additional items non-blocking до maxSize или лимита окна
   - `processBatch(batch, emb)` — конкатенирует тексты всех items, оцин TEI call, splits результаты per-item by boundary
   - whole-batch failure на TEI ошибку (все items mark failed)
   - drainRemaining (shutdown drain with ctx err)

2. **WorkerConfig struct** с Normalize() — safe defaults, тестируемо.

3. **main.go wiring**: передача WorkerConfig из env vars FD_QUEUE_BATCH_MAX_SIZE (default 32), FD_QUEUE_BATCH_WINDOW_MS (default 10ms).

4. **.env.example extended** с FD_QUEUE_BATCH_MAX_SIZE, FD_QUEUE_BATCH_WINDOW_MS.

5. **Tests обновлены** под новый Config + новые тесты:
   - TestWorkerBatchesMultipleItems — 3 items → 1 TEI call, batch_size=3
   - TestWorkerRespectsMaxBatchSize — 8 items, maxBatchSize=4 → 2 TEI calls по 4 item каждый

All 11 packages green with -race.

## Verification

go test ./api/queue/... -race — 5 tests pass (2 single-pass + 3 batch). go test ./api/... -race — 11 packages green. Env vars .env.example документированы.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./queue/... -race -count=1 -timeout 30s` | 0 | pass | 1300ms |
| 2 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13500ms |

## Deviations

None.

## Known Issues

Default batch window 10ms — acceptable для burst pipeline. Если pipeline завал single-item requests, coalescing не сработает (всё в single-pass). Это expected для default.

## Files Created/Modified

- `api/queue/worker.go`
- `api/queue/worker_test.go`
- `api/main.go`
- `.env.example`
