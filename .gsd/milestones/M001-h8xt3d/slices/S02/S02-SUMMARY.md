---
id: S02
parent: M001-h8xt3d
milestone: M001-h8xt3d
provides:
  - A stricter and directly-tested HTTP handler layer for runtime hardening in S04.
requires:
  []
affects:
  - S04
key_files:
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/handlers/embeddings_integration_test.go
key_decisions:
  - Use minimal interfaces in handlers instead of copying handler logic in tests.
patterns_established:
  - Handler dependencies should be interfaces when it enables testing production code without copy-pasting handlers.
observability_surfaces:
  - Invalid batch inputs now return explicit 400 error messages.
drill_down_paths:
  - .gsd/milestones/M001-h8xt3d/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M001-h8xt3d/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M001-h8xt3d/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T06:58:36.633Z
blocker_discovered: false
---

# S02: API validation and handler tests

**S02 made batch validation strict and tests exercise production handlers.**

## What Happened

S02 fixed the validation and test-fidelity findings. Batch requests now reject unsupported dimensions and encoding formats with clear 400 errors. Production handlers now depend on small interfaces, allowing tests to exercise the real handler methods with mocks instead of maintaining a copied test handler implementation.

## Verification

Handler package tests and full short suite passed.

## Requirements Advanced

- Review remediation medium-risk API/test findings advanced with tests. — 

## Requirements Validated

- /embeddings/batch invalid dimensions return 400 — proved by TestCreateBatchEmbeddings_Validation.
- /embeddings/batch invalid encoding_format returns 400 — proved by TestCreateBatchEmbeddings_Validation.
- Production handlers are tested directly — proved by tests constructing NewEmbeddingsHandler and NewBatchHandler.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Batch base64 response tests use mocks and do not require live TEI.

## Follow-ups

Continue with S03 LocalCache semantics, then S04 Docker/config hardening.

## Files Created/Modified

- `api/handlers/embeddings.go` — Added handler dependency interfaces for production-handler tests.
- `api/handlers/batch.go` — Added strict batch dimensions and encoding_format validation.
- `api/handlers/embeddings_integration_test.go` — Rewrote tests to exercise production handlers and added batch validation coverage.
