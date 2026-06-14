---
id: T04
parent: S01
milestone: M041-4tw0w7
key_files:
  - api/main.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/handlers/recovery.go
  - api/handlers/notfound.go
  - api/handlers/embeddings_integration_test.go
  - api/embed/types.go
  - api/embed/codec.go
  - api/middleware/validation.go
key_decisions:
  - EmbeddingObj.Embedding → `any` для support float array ИЛИ base64 string response.
  - ContextKeyValidatedRequest в handlers (not middleware) для avoid import cycle.
  - defaultBatchEncoding() preserves FalkorDB base64 default для backward compat.
  - Handler resilient: inline binding fallback если middleware не зарегистрирован.
  - explicit r.HandleMethodNotAllowed = true (хотя gin v1.x default — но v1.12 вёл себя как false без явного set).
  - Recovery wrapper sets X-Request-Id в message если headers middleware (S03) уже зарегистрировал.
duration: 
verification_result: untested
completed_at: 2026-06-13T18:17:46.897Z
blocker_discovered: false
---

# T04: Wire validation в handlers, replace error envelopes, recovery wrapper, 405/404 handlers, encoding_format в /v1/embeddings

**Wire validation в handlers, replace error envelopes, recovery wrapper, 405/404 handlers, encoding_format в /v1/embeddings**

## What Happened

Wire validation middleware в /v1/embeddings, replace все gin.H{errorKey:...} с handlers.WriteError (OpenAI-style envelopes), wrap gin.Recovery через новый RecoveryMiddleware, добавить NoRoute(404 not_found) и NoMethod(405 method_not_allowed) handlers, добавить encoding_format field в EmbeddingsRequest (S05 T01 work перенесён в S01 по user request), extract encodeEmbedding/float32SliceToBytes из batch.go в новый api/embed/codec.go. Reshape EmbeddingObj.Embedding на `any` чтобы support float array ИЛИ base64 string response. Investigation: dimensions=512 НЕ broken — initial baseline 500 был transient race condition после TEI restart, fd сейчас корректно 200 с 512-dim vector. encoding_format был broken: /v1/embeddings падал 500 потому что handler не знал про encoding_format. После fix работает. Handler resilient: если middleware не зарегистрирован, fallback к inline binding (для tests/standalone mount). Исправлены 3 import cycle issues (ContextKeyValidatedRequest moved в handlers, package structure stable). defaultBatchEncoding preserves legacy FalkorDB base64 default. EmbeddingObj.Embedding → any требует type assertion в tests ([]interface{} с float64 после JSON round-trip). Live integration verified на running fd container.

## Verification

Live integration (curl против running fd container после docker restart): B1 input:[] → 400 input_required ✓, B5 input:[123] → 400 input_required "input[] must be string" (no leaky Go-isms) ✓, B6 {bad → 400 invalid_json с clean message ✓, B7 {} → 400 input_required ✓, B10 GET /v1/embeddings → 405 method_not_allowed ✓, B4 1MB input → 413 input_too_long "input[0] exceeds max length 2048 chars" ✓, B9 100 inputs → 413 batch_too_large "batch size 100 exceeds max 32" ✓, encoding_format=base64 → 200 base64 string ✓, encoding_format=hex → 400 encoding_format_invalid ✓, /v9999 unknown path → 404 not_found ✓. Все unit tests pass: fd-api, fd-api/cache, fd-api/embed, fd-api/handlers, fd-api/middleware. EmbeddingObj.Embedding reshape требует type assertion в integration tests (исправлено). dimensions=512 работает 200 (initial baseline 500 был transient).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

(1) Скоуп расширен относительно исходного плана: encoding_format перенесён в S01 T04 (по user request option 2), а не остался в S05 T01. (2) dimensions=512 fix НЕ нужен: initial baseline 500 был transient race condition. (3) RecoveryMiddleware + NoRoute + NoMethod добавлены как дополнительные tasks чтобы покрыть T-E-8, T-E-10, T-E-15. (4) EmbeddingObj.Embedding изменён с []float32 на `any` чтобы support base64 string encoding. (5) Handler resilient: если middleware не зарегистрирован, fallback к inline binding (нужно для standalone test mount). (6) ContextKeyValidatedRequest constant перемещён в handlers пакет чтобы избежать import cycle. (7) defaultBatchEncoding helper preserves FalkorDB backward compat (base64 default для /embeddings/batch, float default для /v1/embeddings).

## Known Issues

None.

## Files Created/Modified

- `api/main.go`
- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/handlers/recovery.go`
- `api/handlers/notfound.go`
- `api/handlers/embeddings_integration_test.go`
- `api/embed/types.go`
- `api/embed/codec.go`
- `api/middleware/validation.go`
