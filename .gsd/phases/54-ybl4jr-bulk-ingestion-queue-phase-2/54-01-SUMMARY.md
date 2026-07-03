---
id: S01
parent: M054-ybl4jr
milestone: M054-ybl4jr
provides:
  - async POST /v1/queue (202 + X-Request-Id)
  - GET /v1/queue/:id polling (pending/completed/failed/404)
  - Worker time-windowed batched processing (32 items / 10ms → 1 TEI call)
  - in-memory ResultStore with TTL eviction (5 min)
  - 5 queue metric series at /metrics
  - FD_QUEUE_ENABLED feature gate default false (backward compatible)
  - backpressure via non-blocking channel send + 503 + Retry-After
requires:
  []
affects:
  - Phase 1b optimization uses fd_queue_batch_size/process_duration metrics
  - Phase 2b persistence: extend ResultStore for Redis-backed interface
  - Phase 3 (int8): same workload pattern via queue, not on hot path
key_files:
  - api/queue/types.go
  - api/queue/id.go
  - api/queue/store.go
  - api/queue/worker.go
  - api/queue/worker_test.go
  - api/handlers/queue_handlers.go
  - api/observability/metrics.go
  - api/fd_v2_queue_integration_test.go
  - api/main.go
  - .env.example
  - README.md
key_decisions: []
patterns_established:
  - WorkerConfig + Normalize() pattern — env-tunable parameters with safe defaults
  - Async feature-gate pattern (FD_QUEUE_ENABLED default false)
  - Drain-on-shutdown pattern (no silent loss)
  - Time-windowed batching via timer + non-blocking channel drain
observability_surfaces:
  - 5 new Prometheus metrics at /metrics
  - README Async queue section in Operations
  - fd-api startup log: "queue enabled capacity=N" when FD_QUEUE_ENABLED=true
drill_down_paths:
  - /root/fd/.gsd/phases/54-ybl4jr-bulk-ingestion-queue-phase-2/
duration: ""
verification_result: passed
completed_at: 2026-07-03T11:28:15.492Z
blocker_discovered: false
---

# S01: Bulk ingestion queue MVP

**Bulk ingestion queue MVP из Issue #9 Phase 2: async POST /v1/queue + GET /v1/queue/:id polling, time-windowed batched worker (10ms × 32 inputs), in-memory result store с TTL, FD_QUEUE_ENABLED feature gate, 5 queue метрик, 4 integration теста, README async queue section. 11 пакетов тестов зелёные с -race.**

## What Happened

M054-ybl4jr S01 закрывает bulk ingestion queue MVP из Issue #9 Phase 2. Реализован в 4 тонких тасках:

T01 (Queue contract): api/queue/{types,id,store}.go + handlers/queue_handlers.go + main.go wiring через FD_QUEUE_ENABLED feature gate. Backpressure через select-default + 503 + Retry-After.

T02 (Worker single-pass): api/queue/worker.go с graceful shutdown drain, 3 worker unit tests.

T03 (Time-windowed batching): WorkerConfig struct + Normalize() с safe defaults. Drain до FD_QUEUE_BATCH_MAX_SIZE (32) items в FD_QUEUE_BATCH_WINDOW_MS (10ms) → 1 TEI call. Per-item result split. +2 batch теста.

T04 (Metrics + integration + docs): 5 queue Prometheus series (depth, drain, submit, batch_size, process_duration), 4 integration теста (submit-poll, invalid, backpressure, 404), README async queue section.

Изменения additive: sync /v1/embeddings unchanged, embed.Embedder interface unchanged, /health unchanged. FD_QUEUE_ENABLED default false - production без env vars не затрагивается. 11 packages green with -race.

## Verification

All 4 tasks complete. Full unit (5 worker tests) + integration (4 queue tests) test suite passes with -race. /metrics exposes 5 new queue series. README документирован. Feature gated by FD_QUEUE_ENABLED default false — backward compatible.

## Requirements Advanced

None.

## Requirements Validated

- R-new-async-bulk — 5 worker tests + 4 integration tests + 11 packages green with -race

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Results not persisted across fd-api restarts. In-memory only.

## Follow-ups

Phase 2b: Redis-backed persistence. Phase 1b: deferred until production metrics baseline.

## Files Created/Modified

- `api/queue/types.go` — New: queue types (Item, Result, Status) with sentinel errors
- `api/queue/id.go` — New: crypto/rand request ID generator
- `api/queue/store.go` — New: in-memory ResultStore with TTL eviction goroutine
- `api/queue/worker.go` — New: time-windowed batched worker with graceful shutdown drain
- `api/queue/worker_test.go` — New: 5 worker unit tests (process, error, cancel, batch, max-batch)
- `api/handlers/queue_handlers.go` — New: POST /v1/queue + GET /v1/queue/:id handlers with backpressure
- `api/fd_v2_queue_integration_test.go` — New: 4 integration tests
- `api/observability/metrics.go` — 5 new queue metrics + public methods
- `api/main.go` — FD_QUEUE_ENABLED feature gate + worker + handler registration
- `.env.example` — Queue env vars documentation
- `README.md` — Async queue section in Operations
