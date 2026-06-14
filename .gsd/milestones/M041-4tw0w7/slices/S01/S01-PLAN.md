# S01: Validation and OpenAI style error envelope

**Goal:** Validation middleware перед model call + OpenAI-style error envelope + правильный HTTP status mapping по Section 3 каталогу. Закрывает P0 баги B4 (1MB timeout), B7 (misleading parser error), B8 (10 inputs timeout), B9 (100 inputs 500), B10 (GET /v1/embeddings должен быть 405 не 404).
**Demo:** After this, /v1/embeddings возвращает OpenAI style error envelope с machine readable code/type, oversized batch и oversized input дают 413 (не 500), malformed JSON даёт clean 400 invalid_json (не сырую Gin ошибку), и все 16 кодов из Section 3 catalog работают.

## Must-Haves

- T-E-1..T-E-7 (validation error tests) pass: input_required, input_too_long, batch_too_large, dimensions_invalid, invalid_json все возвращают правильный code/type/status
- T-E-8 pass: GET /v1/embeddings возвращает 405 (не 404)
- T-E-14 pass: 50MB body возвращает 413 payload_too_large (не silent 500)
- T-E-15 pass: forced internal error возвращает 500 internal_error с X-Request-Id в message
- 0 cases где current implementation returns 500 на oversized batch (B9 fix)
- Все 16 error codes реализованы и покрыты тестами
- Backward compat: v1 caller (только POST /v1/embeddings) продолжает работать
- golangci-lint pass

## Proof Level

- This slice proves: contract + integration

## Integration Closure

Validation middleware вызывается из всех request handlers. Error envelope используется всеми error-emitting endpoints.

## Verification

- fd_errors_total{code=...} counter появляется в /metrics после S03, но error path уже структурирован: каждый emit логирует code, type, request_id, path.

## Tasks

- [x] **T01: Recon текущего fd Go pipeline выполнен в M041 planning phase** `est:2h`
  Прочитать api/main.go, api/handlers/embeddings.go, api/handlers/batch.go, api/handlers/health.go, и любые middleware-файлы. Зафиксировать: где request входит, какие middleware уже есть, где происходит unmarshal, где вызывается model, где текущие error responses создаются. Решить где именно вставить validation middleware. Результат: .gsd/milestones/M041-4tw0w7/slices/S01/S01-RECON.md с диаграммой pipeline и точкой вставки. Спека fd v2 в /root/fd-v2.md (внешний reference, не в репо).
  - Files: `api/main.go`, `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/handlers/health.go`
  - Verify: Recon MD написан, содержит ASCII диаграмму request path от http.Server до model call, и explicit точку вставки validation middleware (до/после каких существующих middleware).

- [x] **T02: Error envelope + 17 error codes catalog с table-driven tests** `est:3h`
  api/handlers/errors.go: type ErrorResponse {Error ErrorDetail}; type ErrorDetail {Code, Type, Param, Message}. Конструктор NewError(code, httpStatus) возвращает (httpStatus, body). Все 16 кодов из Section 3 каталога: input_required, input_too_long, batch_too_large, dimensions_invalid, dimensions_required, invalid_json, unauthorized, not_found, payload_too_large, rate_limit_exceeded, internal_error, model_not_loaded, model_overloaded, shutting_down, request_timeout, dimensions_mismatch. Каждый код маппится на (type, httpStatus) одной таблицей. Константы для каждого code value (goconst compliance). Helper WriteError(c *gin.Context, code string, param string, details ...any).
  - Files: `api/handlers/errors.go`, `api/handlers/errors_test.go`
  - Verify: go test ./api/handlers/... -run TestErrorEnvelope покрывает все 16 кодов: проверяет HTTP status, code, type, message format. Все тесты pass. golangci-lint pass.

- [x] **T03: Validation middleware (повторное complete после replan)** `est:4h`
  api/middleware/validation.go: gin middleware который читает body через http.MaxBytesReader(w, body, 10*1024*1024) для size limit, парсит JSON, валидирует input array (non-empty, len<=32, все strings, каждый string <=2048 chars), валидирует dimensions (если указан то 512 или 1024). На любой failure — abort с правильным error envelope из T02 через WriteError. Validate BEFORE model call. Использовать struct tags binding:"required,max=32,dive,max=2048" где возможно.
  - Files: `api/middleware/validation.go`, `api/middleware/validation_test.go`
  - Verify: Unit tests: T-E-1 ({} missing input → 400 input_required), T-E-2 (input:[] → 400 input_required), T-E-3 (dimensions:99999 → 400 dimensions_invalid), T-E-4 (input:[123] → 400 invalid_request_error НЕ 'json: cannot unmarshal'), T-E-5 (malformed JSON → 400 invalid_json), T-E-6 (100 inputs → 413 batch_too_large), T-E-7 (10000 char string → 413 input_too_long), T-E-14 (50MB body → 413 payload_too_large). Все 8 test cases pass.

- [x] **T04: Wire validation в handlers, replace error envelopes, recovery wrapper, 405/404 handlers, encoding_format в /v1/embeddings** `est:6h`
  (a) api/handlers/embeddings.go: убрать существующие ad-hoc error responses (те что возвращают gin.H{errorKey:...}), использовать errors.WriteError из T02. Подключить validation middleware из T03 в router setup. 405 handler для неправильных methods на /v1/embeddings (T-E-8). Recovery middleware обёрнут чтобы возвращал OpenAI envelope с X-Request-Id (T-E-15). (b) api/handlers/batch.go: то же самое — replace error responses. (c) INVESTIGATE dimensions=512 broken: T01 baseline показал /v1/embeddings {"input":["hello"],"dimensions":512} → 500. Root cause options: (i) TEI model --dtype fp16 не поддерживает 512-dim Matryoshka head, (ii) fd handler truncation bug, (iii) TEI returns 1024 но fd expects 512. Investigate через прямой TEI запрос: curl http://127.0.0.1:30080/embed -d '{"input":"hello","dimensions":512}'. Apply fix: если TEI side, fall back to 1024 с warning (или явная 400 dimensions_mismatch); если fd side, fix handler. (d) MOVE encoding_format codec: extract encodeEmbedding и float32SliceToBytes из api/handlers/batch.go в новый api/embed/codec.go. Добавить EncodingFormat *string в api/embed.EmbeddingsRequest. Validation в T03 принимает encoding_format=float|base64. Handler использует codec для response.
  - Files: `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/main.go`, `api/embed/types.go`, `api/embed/codec.go`, `api/embed/onnx.go`, `api/embed/tei.go`
  - Verify: Integration tests против running fd: T-E-1..T-E-8, T-E-15 все pass. B4 (1MB input) → 413 input_too_long (НЕ timeout 500). B7 ({} missing input) → 400 input_required (НЕ 'unexpected end of JSON'). B9 (100 inputs) → 413 batch_too_large (НЕ 500). B10 (GET /v1/embeddings) → 405 (НЕ 404). Recovery wrapped: forced panic → 500 internal_error с X-Request-Id. dimensions=512: либо fix работает (200), либо явная 400 dimensions_mismatch с понятным message. encoding_format=base64 → 200 с base64 string в response (per T-H-5). encoding_format=float → 200 с float array (default). encoding_format=garbage → 400 dimensions_invalid.

- [x] **T05: Live integration verification: 12 probe bugs против running fd container, все pass или в correct envelope** `est:3h`
  tests/integration/fd_v2_validation_test.go: автоматизировать все 15 error path test cases (T-E-1..T-E-15) из docs/fd-v2.md Section 5.2. Также regression test для backward compat: v1 caller POST /v1/embeddings {input:[hello]} → 200, response object/data/embedding/usage/model fields присутствуют. Также test для encoding_format=base64 (T-H-5) и dimensions=512 fix.
  - Files: `tests/integration/fd_v2_validation_test.go`, `tests/integration/fd_v1_backward_compat_test.go`
  - Verify: go test ./tests/integration/... -run TestFdV2Validation -v: все 15 test cases pass. Backward compat test pass. encoding_format и dimensions=512 tests pass.

## Files Likely Touched

- api/main.go
- api/handlers/embeddings.go
- api/handlers/batch.go
- api/handlers/health.go
- api/handlers/errors.go
- api/handlers/errors_test.go
- api/middleware/validation.go
- api/middleware/validation_test.go
- api/embed/types.go
- api/embed/codec.go
- api/embed/onnx.go
- api/embed/tei.go
- tests/integration/fd_v2_validation_test.go
- tests/integration/fd_v1_backward_compat_test.go
