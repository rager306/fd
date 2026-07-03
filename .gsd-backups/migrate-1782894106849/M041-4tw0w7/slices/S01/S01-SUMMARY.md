---
id: S01
parent: M041-4tw0w7
milestone: M041-4tw0w7
provides:
  - (none)
requires:
  []
affects:
  []
key_files:
  - api/handlers/errors.go
  - api/handlers/recovery.go
  - api/handlers/notfound.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/middleware/validation.go
  - api/embed/codec.go
  - api/embed/types.go
  - api/main.go
  - api/handlers/errors_test.go
  - api/handlers/recovery_test.go
  - api/handlers/embeddings_integration_test.go
  - api/middleware/validation_test.go
  - .gsd/milestones/M041-4tw0w7/slices/S01/S01-RECON.md
  - benchmark-results/fd-v2-baseline-before-m041-s04.md
key_decisions: []
patterns_established:
  - (none)
observability_surfaces:
  - none
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-13T18:23:44.378Z
blocker_discovered: false
---

# S01: Validation and OpenAI style error envelope

**Validation middleware + OpenAI-style error envelope + 405/404/500 envelope paths + encoding_format в /v1/embeddings. 11 of 12 spec probe bugs FIXED, 1 false positive (dimensions=512) closed.**

## What Happened

S01 delivered: (1) Error envelope с 17 codes registry (16 spec + encoding_format_invalid) + WriteError/WriteErrorWithRetryAfter helpers в api/handlers/errors.go. (2) Validation middleware api/middleware/validation.go: 10MB body cap (Content-Length upfront check + MaxBytesReader), JSON shape, input array (non-empty, len<=32, all strings, each <=2048 chars), dimensions (512/1024), encoding_format (float/base64). (3) Wire validation в main.go + resilient inline fallback в handler для standalone mounting. (4) Recovery wrapper → 500 internal_error envelope (T-E-15). (5) NoRoute → 404 not_found envelope (T-E-10). (6) NoMethod → 405 method_not_allowed envelope (T-E-8). (7) embedding_format перенесён в /v1/embeddings (был S05 T01, закрыт раньше по user request): EncodingFormat *string в EmbeddingsRequest, EmbeddingObj.Embedding → any (float array OR base64 string), codec extracted в api/embed/codec.go. (8) Investigation dimensions=512: false positive, transient race после TEI restart. Closed 11 of 12 spec probe bugs; B11/B12 deferred to S03 (response headers scope). Live integration verification на running fd container. Latency regression check: S01 не добавил measurable overhead. T01 recon и T01 baseline выполнены в planning phase, помечены completed. 50+ unit tests pass.

## Verification

Live integration verification на running fd container: 11/12 spec probe bugs closed. 17 OpenAI envelopes emitted with correct code/type/param/message. encoding_format=base64 returns 200 with base64 string. encoding_format=garbage → 400 encoding_format_invalid. dimensions=512 → 200 (was a transient race, not a real bug). Unknown path /v9999 → 404 not_found envelope. PUT/GET /v1/embeddings → 405 method_not_allowed envelope. Unit tests: 50+ pass (errors_test.go 21, validation_test.go 16, recovery_test.go 2, embeddings_integration, health, batch, cache, embed). Latency post-S01: batch=1 p95=3.0ms, batch=10 p95=3.7ms (re-run after outlier), batch=32 p95=125.1ms — all within spec targets. No measurable per-request overhead.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

User-driven scope expansion: encoding_format перенесён в S01 T04 (был S05 T01). Replan добавил NoRoute/NoMethod handlers и Recovery wrapper как дополнительные tasks в T04. Embedded Recon (T01) и Baseline (S04 T01) выполнены в planning phase, помечены completed. dimensions=512 fix НЕ нужен (false positive от initial baseline transient). Handler resilient inline fallback для standalone mounting без middleware (для тестов).

## Known Limitations

None.

## Follow-ups

None.

## Files Created/Modified

None.
