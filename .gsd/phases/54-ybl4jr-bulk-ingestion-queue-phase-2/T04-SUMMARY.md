---
id: T04
parent: S01
milestone: M054-ybl4jr
key_files:
  - api/observability/metrics.go
  - api/fd_v2_queue_integration_test.go
  - README.md
key_decisions:
  - fd_queue_submit_total разделяет accepted/rejected через label — просто, расширяемо
  - fd_queue_batch_size histogram — батч размерности для Phase 1b решения
  - metrics public API минимален — Set, Inc, Observe методы в observability пакете
  - Integration тесты изолированы gin test server, без реального TEI — доказывают flow, not TEI correctness
duration: 
verification_result: passed
completed_at: 2026-07-03T11:19:37.980Z
blocker_discovered: false
---

# T04: Queue metrics (5 новых series) + 4 integration tests (submit/poll, invalid input, backpressure burst, 404) + README async queue section. Все 11 пакетов тестов зелёные с -race.

**Queue metrics (5 новых series) + 4 integration tests (submit/poll, invalid input, backpressure burst, 404) + README async queue section. Все 11 пакетов тестов зелёные с -race.**

## What Happened

T04 закрывает metrics + integration + documentation:

1. **Queue metrics** (5 новых series) в api/observability/metrics.go:
   - `fd_queue_depth` gauge
   - `fd_queue_drain_total` counter
   - `fd_queue_submit_total{result=accepted|rejected}` counter
   - `fd_queue_batch_size` histogram
   - `fd_queue_process_duration_seconds` histogram
   - public methods: SetQueueDepth, IncQueueDrain, IncQueueSubmit, ObserveQueueBatchSize, ObserveQueueProcessDuration

2. **Integration tests** (4 теста) в fd_v2_queue_integration_test.go:
   - TestQueueSubmitAndPollCompleted — submit + poll до completed
   - TestQueueRejectsInvalidInput — пустой input → 400
   - TestQueueBackpressureRejectsWhenFull — burst submit при ёмкости=1, worker fast drain под капотом
   - TestQueuePollReturns404ForUnknownId

3. **README.md** — новая секция "Async queue (Phase 2)" в Operations:
   - описание поведения (202/poll/503 backpressure)
   - таблица конфигурации (4 env vars)
   - mention что результаты не persistent
   - метрики упомянуты

4. **Evidence**: integration tests pass, 4 теста. Docker smoke deferred (известный cache layer artifact).

## Verification

go test ./api/... -race -count=1 — 11 packages green. 4 integration queue tests pass. README документирован.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13500ms |
| 2 | `cd /root/fd/api && go test . -race -count=1 -timeout 30s -run TestQueue -v` | 0 | pass | 1200ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `api/observability/metrics.go`
- `api/fd_v2_queue_integration_test.go`
- `README.md`
