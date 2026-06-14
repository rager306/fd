# S05: OpenAI v2 compat features OpenAPI schema and P2 enhancements

**Goal:** OpenAI v2 compat (encoding_format=user=priority), optional API key auth (FD_API_KEY), CORS, OpenAPI 3.1 schema at /openapi.json + Swagger UI at /docs, /v1/batch, rate limiting (per-IP/user), ETag+Cache-Control, /v1/traces. Закрывает R-P1-5..R-P1-9, R-P2-1..R-P2-6.
**Demo:** After this, remaining OpenAI v2 compat work is complete without reimplementing encoding_format already covered by S01; user/priority/auth/CORS/OpenAPI/rate-limit/traces work passes M043 gates and documents any new API godoc.

## Must-Haves

- T-H-5 pass: encoding_format=base64 возвращает base64 string
- T-H-6 pass: priority=high принимается
- T-E-9 pass: с FD_API_KEY=test, запрос без Authorization → 401 unauthorized
- CORS preflight (OPTIONS) → 200 с правильными Access-Control-Allow-* headers
- T-E-4/5 (Section 5.5) pass: /openapi.json 200 валидный JSON, /docs 200 HTML
- /v1/batch принимает {"batches": [[..],[..]]} и возвращает {"batches": [[..],[..]]}
- Rate limit: 101-й запрос за минуту → 429 с X-RateLimit-* headers и Retry-After
- ETag: повторный запрос с If-None-Match → 304
- /v1/traces возвращает последние 100 requests с latency, status, model_id, request_id
- tools/verify_fd_v2_contract.py exit 0 на running fd v2, 45/45 test cases pass
- benchmark-results/fd-v2-validation-m041.md final artifact
- golangci-lint pass

## Proof Level

- This slice proves: contract + integration

## Integration Closure

Использует validation из S01 (новые поля валидируются тем же middleware), error envelope из S01, headers из S03, lifecycle state из S02. Auth middleware (если FD_API_KEY) сидит после headers middleware, до validation.

## Verification

- Добавляет X-RateLimit-* headers, расширяет metrics: fd_rate_limit_exceeded_total, fd_cache_evictions_total. /v1/traces использует существующий request log ring buffer.

## Tasks

- [x] **T01: Added OpenAI-compatible `user` and `priority` request fields, priority validation, and tests for base64/user/priority behavior.** `est:3h`
  api/handlers/embeddings.go: расширить request struct: EncodingFormat *string (valid: float|base64), User *string, Priority *string (valid: low|normal|high). Validation в S01 middleware: невалидный encoding_format → 400. Base64 encoding для response: при encoding_format=base64, кодировать []float32 в base64-encoded float32 LE array.
  - Files: `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/embed/base64.go`, `api/middleware/validation.go`, `api/handlers/embeddings_test.go`
  - Verify: Unit tests: T-H-5 (encoding_format=base64 → base64 string в response), T-H-6 (priority=high принимается), user field принимается. Невалидный encoding_format → 400.

- [x] **T02: Added optional FD_API_KEY bearer auth and CORS/preflight middleware with tests and main wiring.** `est:3h`
  api/middleware/auth.go: если env FD_API_KEY задан, требует Authorization: Bearer <key> на всех endpoints кроме /live, /metrics, /docs, /openapi.json. На missing/wrong → 401 unauthorized. api/middleware/cors.go: Access-Control-Allow-Origin (из env FD_CORS_ORIGINS или * default), Access-Control-Allow-Methods: GET,POST,OPTIONS, Access-Control-Allow-Headers: Content-Type,Authorization,X-Request-Id. OPTIONS preflight → 204.
  - Files: `api/middleware/auth.go`, `api/middleware/cors.go`, `api/middleware/auth_test.go`, `api/middleware/cors_test.go`
  - Verify: Unit tests: T-E-9 (с FD_API_KEY=test, без Authorization → 401 unauthorized, с правильным Bearer → 200). OPTIONS preflight → 204 с правильными CORS headers.

- [x] **T03: Added opt-in per-IP and per-user rate limiting with X-RateLimit headers and 429 retry envelopes.** `est:3h`
  api/middleware/ratelimit.go: token bucket per IP (100 req/min default) и per user (1000 req/min default если user field задан). Env FD_RATE_LIMIT_IP_RPM, FD_RATE_LIMIT_USER_RPM для конфигурации. Headers X-RateLimit-Limit/Remaining/Reset на каждом response. На превышение → 429 rate_limit_exceeded + Retry-After: 60. Опционально через FD_RATE_LIMIT_ENABLED=true (default false для обратной совместимости).
  - Files: `api/middleware/ratelimit.go`, `api/middleware/ratelimit_test.go`
  - Verify: Unit tests: с включённым rate limit, 101-й запрос за минуту → 429 с X-RateLimit-* headers и Retry-After: 60. Per-user limit отдельно от per-IP.

- [x] **T04: Added `/v1/batch` endpoint for multiple inner batches with validation, cache/model execution, and tests.** `est:3h`
  api/handlers/v1batch.go: POST /v1/batch принимает {"batches": [[s1, s2, ...], [s1, s2, ...], ...]} (каждый inner array ≤ 32 strings, total ≤ 100 batches). Возвращает {"batches": [[e1, e2, ...], [e1, e2, ...], ...]}. Validation аналогично /v1/embeddings. Каждый inner batch обрабатывается через тот же pipeline (validation → lifecycle → cache → model).
  - Files: `api/handlers/v1batch.go`, `api/handlers/v1batch_test.go`
  - Verify: Unit tests: 2 batches × 4 strings → 200 с 2 batches × 4 embeddings. Oversized inner batch → 413. Empty batches → 400 input_required.

- [x] **T05: Added ETag and Cache-Control middleware for `/v1/embeddings` and `/info` with If-None-Match 304 support.** `est:2h`
  api/middleware/cache_headers.go: на /v1/embeddings и /info responses вычислять ETag = SHA256(response body) и выставлять Cache-Control: public, max-age=86400. Поддержка If-None-Match: если request header matches ETag → 304 Not Modified без body.
  - Files: `api/middleware/cache_headers.go`, `api/middleware/cache_headers_test.go`
  - Verify: Unit tests: первый request → ETag: <hash>, Cache-Control: public, max-age=86400. Повторный с If-None-Match: <hash> → 304.

- [ ] **T06: /v1/traces debug endpoint** `est:2h`
  api/observability/traces.go: in-memory ring buffer (последние 100 requests) с timestamp, latency, status, model_id, request_id, path, dimensions. GET /v1/traces возвращает JSON массив. Использует request_id из headers middleware (S03). Опционально через FD_TRACES_ENABLED=true (default true).
  - Files: `api/observability/traces.go`, `api/observability/traces_test.go`
  - Verify: Unit tests: после 5 requests GET /v1/traces → 200 с 5 entries. Каждая entry содержит timestamp, latency_ms, status, model_id, request_id, path.

- [ ] **T07: OpenAPI 3.1 schema generation и Swagger UI** `est:4h`
  api/openapi/spec.go: генерирует OpenAPI 3.1 spec программно (на основе реальных routes и типов). Включает все endpoints: /health, /live, /ready, /warmup, /version, /info, /metrics, /v1/embeddings, /v1/batch, /v1/healthcheck, /v1/traces, /openapi.json, /docs. Все request/response schemas, headers, error envelope. GET /openapi.json → 200 application/json. GET /docs → 200 text/html с Swagger UI (swagger-ui-dist или CDN). Валидация через openapi-spec-validator.
  - Files: `api/openapi/spec.go`, `api/openapi/spec_test.go`, `api/handlers/openapi.go`, `api/handlers/docs.go`
  - Verify: Unit tests: /openapi.json возвращает валидный OpenAPI 3.1 (validate через openapi-spec-validator). /docs возвращает HTML с swagger-ui в body. Integration test: curl /openapi.json | openapi-spec-validator — exit 0.

- [ ] **T08: Final 45-test acceptance suite и MILESTONE-UAT** `est:3h`
  tools/verify_fd_v2_contract.py: автоматизировать ВСЕ 45 test cases. Скрипт запускает каждый test case против running fd v2, проверяет HTTP status, body shape, headers, и пишет results в JSON. Final artifact: benchmark-results/fd-v2-validation-m041.md со всеми 45 test results, p95 perf numbers, и pass/fail summary. Если хоть 1 test fail — exit 1.
  - Files: `tools/verify_fd_v2_contract.py`, `benchmark-results/fd-v2-validation-m041.md`
  - Verify: tools/verify_fd_v2_contract.py exit 0 на running fd v2, все 45 test cases pass. Artifact содержит pass/fail breakdown.

## Files Likely Touched

- api/handlers/embeddings.go
- api/handlers/batch.go
- api/embed/base64.go
- api/middleware/validation.go
- api/handlers/embeddings_test.go
- api/middleware/auth.go
- api/middleware/cors.go
- api/middleware/auth_test.go
- api/middleware/cors_test.go
- api/middleware/ratelimit.go
- api/middleware/ratelimit_test.go
- api/handlers/v1batch.go
- api/handlers/v1batch_test.go
- api/middleware/cache_headers.go
- api/middleware/cache_headers_test.go
- api/observability/traces.go
- api/observability/traces_test.go
- api/openapi/spec.go
- api/openapi/spec_test.go
- api/handlers/openapi.go
- api/handlers/docs.go
- tools/verify_fd_v2_contract.py
- benchmark-results/fd-v2-validation-m041.md
