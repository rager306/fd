---
estimated_steps: 19
estimated_files: 4
skills_used: []
---

# T04: Queue metrics + integration tests + documentation

Добавить queue метрики, интеграционные тесты, и README документацию.

1. **Новые метрики в `api/observability/metrics.go`** (additive):
   - `fd_queue_depth` gauge — текущий размер bounded channel
   - `fd_queue_drain_count` counter — total items processed since start
   - `fd_queue_submit_total{result=accepted|rejected}` counter — запросов принято/отклонено
   - `fd_queue_batch_size` histogram — items в каждом батче worker
   - `fd_queue_process_duration_seconds` histogram — время от submit до completed
   - Additive к существующему observer pattern, не меняет existing metrics

2. **Интеграционный smoke тест** в `api/fd_v2_queue_integration_test.go` (новый файл):
   - TestQueueSubmitThenPollCompleted (submit item, poll, verify completed with embeddings)
   - TestQueueBackpressureRejects (submit cap+1, verify rejection 503)
   - TestQueueWorkerBatcherMerges (подтвердить что batched worker делает 1 TEI call на несколько items)
   - TestQueueDisabledReturnsNotFound (FD_QUEUE_ENABLED=false)
   - TestQueueShutdownGracefulDrain (submit items, cancel ctx, verify all in-flight items либо processed либо failed)

3. **README.md**: новый раздел "Async queue (Phase 2)" под Operations, описывающий endpoint /v1/queue, env vars, cap behavior.

4. **Verify**:
   - go test ./api/... -race -count=1 — все пакеты green
   - Docker smoke: curl /v1/queue submit + poll completed
   - Evidence сохранить в .gsd/runtime/M054-ybl4jr/

## Inputs

- `api/queue/worker.go (T03)`
- `api/handlers/queue_handlers.go (T01)`
- `api/observability/metrics.go (existing metrics)`

## Expected Output

- `api/observability/metrics.go`
- `api/fd_v2_queue_integration_test.go`
- `README.md`
- `.gsd/runtime/M054-ybl4jr/`

## Verification

go test ./api/... -race -count=1 — 10+ packages green. Integration test проходит submit-completed-poll workflow.

## Observability Impact

5 новых метрических series для queue monitoring.
