---
id: T02
parent: S01
milestone: M054-ybl4jr
key_files:
  - api/queue/worker.go
  - api/queue/worker_test.go
  - api/main.go
key_decisions:
  - Single-pass design (no batching yet) — simpler verification first; batching в T03
  - Response channel buffered 1 — submit handler can non-blocking send
  - Drain path on shutdown: marks all queued items failed with ctx.Err (no silent drop)
  - Worker uses single goroutine; concurrency vs ordering vs complexity tradeoff
duration: 
verification_result: passed
completed_at: 2026-07-03T11:01:21.284Z
blocker_discovered: false
---

# T02: Worker goroutine: single-pass TEI calls, graceful shutdown drain, 3 worker tests (process/error/cancel). wired в main.go после queue channel creation. 11 пакетов green с -race.

**Worker goroutine: single-pass TEI calls, graceful shutdown drain, 3 worker tests (process/error/cancel). wired в main.go после queue channel creation. 11 пакетов green с -race.**

## What Happened

T02 закрывает worker goroutine:

1. **api/queue/worker.go** — StartQueueWorker запускает goroutine, draining `items <-chan Item`. Single-pass обработка (без time-windowing — T03 pending): для каждого item вызывает `emb.Embed(ctx, texts)`, writes result в ResultStore + item.Response channel.

2. **Graceful shutdown**: ctx.Done() → drain remaining items marking их StatusFailed с ctx.Err в Response channel. Pollers не блокируются навсегда.

3. **Tests (api/queue/worker_test.go)**:
   - TestWorkerProcessesSingleItem — 1 item → 1 TEI call, status=completed, embeddings count matches
   - TestWorkerHandlesEmbedError — TEI returns error → status=failed, store updated
   - TestWorkerStopsOnContextCancel — cancel ctx → result failed (ctx.Canceled), worker exits

4. **main.go wiring**: StartQueueWorker(ctx, store, items, embeddingClient, logger) запускается после создания channel. defer resultStore.Close() для cleanup.

All 11 packages green with -race.

## Verification

go test ./api/queue/... -race — 3 tests pass. go test ./api/... -race — 11 packages green. Worker wired в main.go через StartQueueWorker.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./queue/... -race -count=1 -timeout 30s` | 0 | pass | 1100ms |
| 2 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13500ms |

## Deviations

None.

## Known Issues

No batching yet (T03) — each item → 1 TEI call. Batch visibility for coalescing decision comes after T03.

## Files Created/Modified

- `api/queue/worker.go`
- `api/queue/worker_test.go`
- `api/main.go`
