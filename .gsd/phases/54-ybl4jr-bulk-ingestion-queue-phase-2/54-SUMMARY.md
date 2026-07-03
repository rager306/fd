---
id: M054-ybl4jr
title: "Bulk ingestion queue (Phase 2)"
status: complete
completed_at: 2026-07-03T11:29:38.307Z
key_decisions:
  - Time-windowed batching via timer + non-blocking channel drain
  - WorkerConfig + Normalize() pattern
  - Public API minimal (Set, Inc, Observe)
  - FD_QUEUE_ENABLED default false — feature backward compatibility
  - Drain-on-shutdown pattern (no silent loss)
  - Whole-batch failure semantics
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
lessons_learned:
  - WorkerConfig + Normalize() pattern reusable for future async workers in fd
  - Test setup helper signals avoid "declared and not used" through `_ = ...` blocks cleanly
  - gin test server end-to-end integration tests по sequence: submit → poll → verify lifecycle testable without real TEI
---

# M054-ybl4jr: Bulk ingestion queue (Phase 2)

**Bulk ingestion queue MVP из Issue #9 Phase 2: async POST /v1/queue + GET /v1/queue/:id polling, time-windowed batched worker (10ms × 32), in-memory result store с TTL, FD_QUEUE_ENABLED feature gate, 5 queue метрик, 4 integration теста. 11 пакетов тестов зелёные с -race.**

## What Happened

M054-ybl4jr закрыт. Bulk ingestion queue MVP из Issue #9 Phase 2 реализован в 4 тонких тасках с полным покрытием unit + integration тестов с -race. 11 пакетов тестов зелёные. Feature gated default false — backward compatible. READ MD и .env.example документированы.

## Success Criteria Results

15/15 success criteria PASS. 5 worker tests + 4 integration tests + 11 packages green.

## Definition of Done Results

All 4 tasks + 15 success criteria met. Feature gated default false. Sync /v1/embeddings unchanged. /health shape unchanged. embed.Embedder interface unchanged.

## Requirement Outcomes

R-new-async-bulk validated. R046 (continuity from M051) still validated.

## Deviations

Два этапа test setup артефакта (unused store var) — фиксил в течение T04.

## Follow-ups

Phase 2b: Redis-backed persistence. Phase 1b: deferred.
