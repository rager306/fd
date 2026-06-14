# fd v2 — Требования к доработке fd embedding service

**Audience**: LLM implementing changes in the fd repository.
**Source project**: daily-archive (M062 fd prod hardening).
**fd current version**: observed 2026-06-13 at `http://127.0.0.1:8000` (local dev).
**fd upstream repo**: separate from daily-archive; not modified by daily-archive code.

---

## 0. Контекст

`fd` — это local HTTP embedding service, обслуживающий daily-archive (scientific KG pipeline).
Модель: `deepvk/USER-bge-m3`, размерности 1024 / 512 (Matryoshka), max input length 512 tokens.

Сервис вызывается daily-archive из:
- `src/arxiv_archive/embedder.py` (class `Embedder`, async, httpx)
- `scripts/m057_table_embed.py` (helper, sync, urllib)
- `scripts/m057_figure_embed.py`, `scripts/m057_fd_validate.py`, `scripts/m058_plotextractor_embed.py`

Daily-archive pipeline обрабатывает 12,362 papers (2-hop BFS), embedding 200+ figures/tables.
Throughput target: 7+ papers/min end-to-end. fd должен выдерживать concurrent batch requests.

Daily-archive НЕ модифицирует fd. Документ `/root/fd-v2.md` — спецификация требований к fd.
После реализации fd v2: daily-archive wrapper (M062 S01) использует новые endpoints, contract tests (M062 S03) валидируют соответствие.

---

## 1. Текущее состояние fd (наблюдаемое, 2026-06-13)

### 1.1 Работающие endpoints

| Method | Path | Status | Behavior |
|---|---|---|---|
| GET | `/health` | 200 | `{"status":"ok","time":"2026-06-13T16:33:15Z"}` — shallow check, не проверяет model load |
| POST | `/v1/embeddings` (OpenAI shape) | 200 | Возвращает embeddings, model `deepvk/USER-bge-m3`, dims 1024/512 |

### 1.2 OpenAI shape (единственная рабочая)

```http
POST /v1/embeddings
Content-Type: application/json

{
  "input": ["text 1", "text 2", ...],  // required, non-empty array of strings
  "dimensions": 1024                    // optional, must be 1024 or 512
}
```

Response (200):
```json
{
  "object": "list",
  "data": [
    {"object": "embedding", "embedding": [0.003, ...], "index": 0, "dimensions": 1024}
  ],
  "model": "deepvk/USER-bge-m3",
  "usage": {"prompt_tokens": 1, "total_tokens": 1}
}
```

### 1.3 НЕ работающие/отсутствующие endpoints

| Method | Path | Status | Issue |
|---|---|---|---|
| GET | `/version` | 404 | Нет version info |
| GET | `/info` или `/v1/models` | 404 | Нет списка моделей |
| GET | `/metrics` | 404 | Нет Prometheus metrics |
| GET | `/v1/healthcheck` | 404 | Нет alias для /health |
| GET | `/ready` | 404 | Нет Kubernetes readiness probe |
| GET | `/live` | 404 | Нет Kubernetes liveness probe |
| GET | `/docs` | 404 | Нет Swagger UI |
| GET | `/openapi.json` | 404 | Нет OpenAPI schema |
| POST | `/embed` (TEI shape) | 404 | Не поддерживает TEI shape |
| POST | `/v1/batch` | 404 | Нет dedicated batch endpoint |

### 1.4 Наблюдаемые баги (12 probe tests)

| # | Сценарий | Response | Status | Severity |
|---|---|---|---|---|
| B1 | `{"input":[]}` | `{"error":"input is required"}` | 400 | OK (валидно) |
| B2 | `{"input":["test"],"dimensions":99999}` | `{"error":"dimensions must be 1024 or 512"}` | 400 | OK |
| B3 | `{"input":["test"],"dimensions":0}` | `{"error":"dimensions must be 1024 or 512"}` | 400 | OK |
| **B4** | **1MB текст в input** | **TIMEOUT 10s (нет error response)** | — | **P0 silent hang** |
| B5 | `{"input":[123]}` (non-string) | `{"error":"json: cannot unmarshal array into Go value of type string"}` | 400 | P1 — leaky Go-isms |
| B6 | malformed JSON `{bad json` | `{"error":"invalid character 'b' looking for beginning of object key string"}` | 400 | P1 — leaky parser error |
| **B7** | `{}` (missing input) | `{"error":"unexpected end of JSON input"}` | 400 | **P1 — misleading** |
| **B8** | **10 inputs (warm model)** | **TIMEOUT 10s (должно быть < 1s)** | — | **P0 performance** |
| **B9** | **100 inputs** | **500 Internal Server Error (silent)** | 500 | **P0 silent failure** |
| B10 | GET `/v1/embeddings` | `404 page not found` | 404 | P1 — should be 405 |
| B11 | Response headers | **empty** (no Server, no X-*, no Cache-*, no ETag) | — | P1 — no observability |
| B12 | Successful response headers | **только `Date`, `Content-Length`** | 200 | P1 — no X-Request-Id, no X-Cache |

### 1.5 Отсутствующие response headers (на всех responses)

- `Server: fd/<version>` — server identification
- `X-Request-Id` — caller-passed или server-generated
- `X-Model-Id` — какая модель использовалась
- `X-Dimensions` — actual dims
- `X-Cache: HIT|MISS` — cache status (если есть cache)
- `X-RateLimit-Limit/Remaining/Reset` — rate limit status
- `Retry-After` — на 429/503
- `Connection: keep-alive` — connection reuse
- `ETag`, `Cache-Control` — response caching

### 1.6 Отсутствующая OpenAPI schema

Нет `/openapi.json`, нет `/docs`. Caller не может:
- Узнать полный список endpoints
- Узнать точный request/response shape
- Узнать error response shape
- Узнать headers
- Узнать rate limits

---

## 2. Требования (по приоритету)

### 2.1 P0 — функциональные баги (MUST FIX)

**R-P0-1**: Input length validation
- Если `len(input[i]) > MAX_INPUT_LENGTH_TOKENS` (для BGE-M3 = 512 tokens, ~2048 chars), вернуть **413 Payload Too Large** с OpenAI-style error:
  ```json
  {"error": {"code": "input_too_long", "type": "invalid_request_error", "param": "input", "message": "input[0] exceeds max length 512 tokens (got ~8192 tokens)"}}
  ```
- Валидация ДО отправки в model (не silent timeout).

**R-P0-2**: Batch size validation
- Если `len(input) > MAX_BATCH_SIZE` (recommend 32, настраивается), вернуть **413**:
  ```json
  {"error": {"code": "batch_too_large", "type": "invalid_request_error", "param": "input", "message": "batch size 100 exceeds max 32; split into smaller batches"}}
  ```
- Валидация ДО отправки в model (не silent 500).

**R-P0-3**: 500 → 503 для model not loaded / overloaded
- Если model не загружен или overload, вернуть **503 Service Unavailable** + `Retry-After: 5`.
- Не 500 (это "server bug"), а 503 (это "temporary unavailable").

**R-P0-4**: Warmup + readiness
- При старте fd: pre-warm model (1 dummy inference) ДО accepting requests.
- Endpoint `GET /ready` возвращает 200 только после pre-warm done.
- Endpoint `GET /live` возвращает 200 пока process alive (cheap, без model touch).
- Endpoint `GET /health` делает DEEP check: model loaded, GPU available, warmup done.

**R-P0-5**: Graceful shutdown
- На SIGTERM: stop accepting new requests (вернуть 503), finish in-flight (макс 30s), затем exit.
- Caller получает 503 с `Retry-After: 30` если послал после SIGTERM.

**R-P0-6**: Performance baseline
- На warm model: 1 input < 50ms p95, 10 inputs < 200ms p95, 32 inputs (max batch) < 1000ms p95.
- Если не выдерживает — add GPU/CPU optimization, batch tensor packing, concurrent workers.

### 2.2 P0 — observability endpoints (MUST ADD)

**R-P0-7**: `GET /version`
```json
{
  "service": "fd",
  "version": "2.0.0",          // semver fd
  "model": "deepvk/USER-bge-m3",
  "model_version": "v1.0",     // model-specific version
  "build_hash": "abc1234",     // git commit hash
  "build_date": "2026-06-13T00:00:00Z",
  "started_at": "2026-06-13T16:30:00Z",
  "uptime_seconds": 3600
}
```

**R-P0-8**: `GET /info` или `GET /v1/models`
```json
{
  "models": [
    {
      "id": "deepvk/USER-bge-m3",
      "dimensions": [512, 1024],
      "max_input_length_tokens": 512,
      "max_batch_size": 32,
      "loaded": true,
      "warmup_done": true,
      "device": "cuda:0"
    }
  ]
}
```

**R-P0-9**: `GET /metrics` (Prometheus text format)
```
# HELP fd_requests_total Total embedding requests
# TYPE fd_requests_total counter
fd_requests_total{status="success"} 1234
fd_requests_total{status="error"} 5
fd_requests_total{status="timeout"} 2

# HELP fd_request_duration_seconds Request latency
# TYPE fd_request_duration_seconds histogram
fd_request_duration_seconds_bucket{le="0.05"} 800
fd_request_duration_seconds_bucket{le="0.1"} 1000
fd_request_duration_seconds_bucket{le="0.5"} 1200
fd_request_duration_seconds_bucket{le="1.0"} 1230
fd_request_duration_seconds_bucket{le="+Inf"} 1239

# HELP fd_batch_size Request batch size
# TYPE fd_batch_size histogram
fd_batch_size_bucket{le="1"} 500
fd_batch_size_bucket{le="10"} 1000
fd_batch_size_bucket{le="32"} 1230
fd_batch_size_bucket{le="+Inf"} 1239

# HELP fd_cache_hits_total Cache hits/misses
# TYPE fd_cache_hits_total counter
fd_cache_hits_total{result="hit"} 800
fd_cache_hits_total{result="miss"} 439

# HELP fd_errors_total Errors by type
# TYPE fd_errors_total counter
fd_errors_total{code="input_too_long"} 3
fd_errors_total{code="batch_too_large"} 1
fd_errors_total{code="model_overloaded"} 1

# HELP fd_model_loaded Model load status (1=loaded, 0=not)
# TYPE fd_model_loaded gauge
fd_model_loaded 1
```

**R-P0-10**: `GET /v1/healthcheck` (alias для /health)

### 2.3 P0 — response headers (MUST ADD)

Каждый response должен иметь:

**R-P0-11**: `X-Request-Id`
- Если caller передал `X-Request-Id` header, echo его.
- Иначе server генерирует UUIDv4.

**R-P0-12**: `Server: fd/<version>` (где version из R-P0-7)

**R-P0-13**: `X-Model-Id: <model_id>` (на /v1/embeddings responses)

**R-P0-14**: `X-Dimensions: <actual_dims>` (на /v1/embeddings responses)

**R-P0-15**: `X-Cache: HIT|MISS` (если есть cache, см. R-P1-3)

**R-P0-16**: На 429/503 — `Retry-After: <seconds>` header

**R-P0-17**: `Connection: keep-alive` (по умолчанию)

### 2.4 P0 — error format (MUST CHANGE)

**R-P0-18**: OpenAI-style error envelope (вместо текущего `{"error": "..."}`):
```json
{
  "error": {
    "code": "input_too_long",
    "type": "invalid_request_error",
    "param": "input",
    "message": "input[0] exceeds max length 512 tokens (got ~8192 tokens)"
  }
}
```

`type` enum:
- `invalid_request_error` — 400 (caller bug)
- `authentication_error` — 401
- `permission_error` — 403
- `not_found_error` — 404
- `rate_limit_error` — 429
- `overloaded_error` — 503
- `internal_error` — 500 (server bug)

`code` enum (canonical, machine-readable):
- `input_too_long`
- `batch_too_large`
- `input_required`
- `dimensions_invalid` (not 512/1024)
- `dimensions_required`
- `dimensions_mismatch` (model doesn't support requested)
- `model_not_loaded`
- `model_overloaded`
- `rate_limit_exceeded`
- `request_timeout`
- `payload_too_large`
- `internal_error`

**R-P0-19**: HTTP status code mapping (machine-readable):
- 200 success
- 400 invalid_request_error (caller bug)
- 401 authentication_error
- 403 permission_error
- 404 not_found_error (path)
- 405 method_not_allowed
- 413 payload_too_large (input_too_long, batch_too_large)
- 429 rate_limit_error
- 500 internal_error (server bug)
- 503 overloaded_error (model not loaded, overloaded, shutting down)
- 504 gateway_timeout (request_timeout)

### 2.5 P1 — health checks (SHOULD ADD)

**R-P1-1**: `GET /health` — deep check
- Проверяет: model loaded, GPU available, warmup done, last inference < 60s ago.
- Response:
```json
{
  "status": "ok",  // "ok" | "degraded" | "down"
  "time": "2026-06-13T16:33:15Z",
  "model_loaded": true,
  "warmup_done": true,
  "device": "cuda:0",
  "last_inference_at": "2026-06-13T16:33:00Z",
  "in_flight_requests": 3
}
```
- 200 если `status="ok"`, 503 если `status="degraded"` или `"down"`.

**R-P1-2**: `GET /warmup` — warmup status
- `{"status": "warming_up"|"ready", "progress": 0.5}`

**R-P1-3**: `POST /warmup` — trigger warmup on demand
- Если не warm — загружает model, делает 1 dummy inference.
- 200 если уже warm, 202 если warming.

### 2.6 P1 — features (SHOULD ADD)

**R-P1-4**: Cache с X-Cache header
- LRU cache на (input_text, dimensions) → embedding.
- Cache size: 10000 (настраивается).
- TTL: 24h (настраивается).
- HIT на cache → skip model inference, return cached.
- `X-Cache: HIT|MISS` в response.
- Metrics: `fd_cache_hits_total{result="hit"|"miss"}`.

**R-P1-5**: `encoding_format` option (OpenAI v2 compat)
- `encoding_format: "float"` (default) — array of floats
- `encoding_format: "base64"` — base64-encoded float32 array (экономит ~30% bandwidth)

**R-P1-6**: `user` field (OpenAI v2 compat)
- `user: "caller-id-123"` — для abuse tracking и per-user rate limits.

**R-P1-7**: `priority` option
- `priority: "low"|"normal"|"high"` — caller помечает приоритет (для routing/queue).

**R-P1-8**: API key auth (env var `FD_API_KEY`)
- Если env var set, требует header `Authorization: Bearer <key>`.
- 401 если отсутствует/wrong.

**R-P1-9**: CORS headers (для web clients)
- `Access-Control-Allow-Origin: *` (или из config)
- `Access-Control-Allow-Methods: GET, POST, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Request-Id`

### 2.7 P2 — nice-to-have (COULD ADD)

**R-P2-1**: OpenAPI schema at `/openapi.json`
- Полная OpenAPI 3.1 spec.
- Swagger UI at `/docs`.

**R-P2-2**: ETag + Cache-Control
- `ETag: "<hash-of-input>"`
- `Cache-Control: public, max-age=86400`

**R-P2-3**: `/v1/batch` separate endpoint
- Для caller-ов которые хотят explicit batch semantics.
- Принимает `{"batches": [[...], [...], ...]}`.
- Возвращает `{"batches": [[...], [...], ...]}`.

**R-P2-4**: Streaming response (SSE)
- Для large responses > 1MB.

**R-P2-5**: Rate limiting
- Per-IP: 100 req/min
- Per-user (если R-P1-6 используется): 1000 req/min
- Headers: `X-RateLimit-Limit/Remaining/Reset`

**R-P2-6**: `/v1/traces` (debugging)
- Recent N requests с latency, status, model_id, request_id.
- Useful для debugging caller-side issues.

---

## 3. Error catalog (machine-readable)

| HTTP | code | type | Когда | Message template |
|---|---|---|---|---|
| 400 | `input_required` | `invalid_request_error` | `input` field missing | `input is required` |
| 400 | `input_too_long` | `invalid_request_error` | input[i] > 512 tokens | `input[0] exceeds max length 512 tokens (got ~{N} tokens)` |
| 400 | `batch_too_large` | `invalid_request_error` | `len(input)` > 32 | `batch size {N} exceeds max 32; split into smaller batches` |
| 400 | `dimensions_invalid` | `invalid_request_error` | `dimensions` not 512/1024 | `dimensions must be 1024 or 512, got {N}` |
| 400 | `dimensions_required` | `invalid_request_error` | `dimensions` missing AND required | `dimensions is required` |
| 400 | `invalid_json` | `invalid_request_error` | body malformed | `invalid JSON: {parser_error}` |
| 401 | `unauthorized` | `authentication_error` | `Authorization` missing/wrong (если R-P1-8) | `missing or invalid API key` |
| 404 | `not_found` | `not_found_error` | unknown path | `path {path} not found` |
| 405 | `method_not_allowed` | (no error envelope) | wrong HTTP method | (HTTP status only, no body) |
| 413 | `payload_too_large` | `invalid_request_error` | body size > MAX_BODY_SIZE (e.g., 10MB) | `request body {N} bytes exceeds max {MAX_BODY_SIZE} bytes` |
| 429 | `rate_limit_exceeded` | `rate_limit_error` | rate limit hit | `rate limit exceeded; retry after {seconds}s` |
| 500 | `internal_error` | `internal_error` | unexpected server bug | `internal server error; request_id={X-Request-Id}` |
| 503 | `model_not_loaded` | `overloaded_error` | model not loaded yet | `model not loaded; retry after {seconds}s` |
| 503 | `model_overloaded` | `overloaded_error` | model overloaded | `model overloaded; retry after {seconds}s` |
| 503 | `shutting_down` | `overloaded_error` | SIGTERM received | `service shutting down; retry after {seconds}s` |
| 504 | `request_timeout` | `overloaded_error` | inference > REQUEST_TIMEOUT | `request timed out after {seconds}s` |

---

## 4. OpenAPI 3.1 spec (sketch)

```yaml
openapi: 3.1.0
info:
  title: fd Embedding Service
  version: 2.0.0
  description: Local HTTP embedding service for daily-archive scientific KG pipeline.
servers:
  - url: http://127.0.0.1:8000
    description: Local dev
paths:
  /health:
    get:
      summary: Deep health check (model loaded, warmup done, GPU available)
      responses:
        '200': {description: ok}
        '503': {description: degraded or down}
  /live:
    get:
      summary: Liveness probe (cheap, no model touch)
      responses:
        '200': {description: alive}
  /ready:
    get:
      summary: Readiness probe (returns 200 only after warmup done)
      responses:
        '200': {description: ready}
        '503': {description: not ready}
  /warmup:
    get:
      summary: Warmup status
      responses:
        '200': {description: ok}
    post:
      summary: Trigger warmup on demand
      responses:
        '200': {description: already warm}
        '202': {description: warming up}
  /version:
    get:
      summary: Service version info
      responses:
        '200': {description: ok}
  /info:
    get:
      summary: Model info (loaded models, dimensions, limits)
      responses:
        '200': {description: ok}
  /metrics:
    get:
      summary: Prometheus metrics
      responses:
        '200':
          description: ok
          content:
            text/plain: {}
  /v1/embeddings:
    post:
      summary: Generate embeddings
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [input]
              properties:
                input:
                  type: array
                  minItems: 1
                  maxItems: 32
                  items:
                    type: string
                    maxLength: 2048
                dimensions:
                  type: integer
                  enum: [512, 1024]
                  default: 1024
                encoding_format:
                  type: string
                  enum: [float, base64]
                  default: float
                user:
                  type: string
                priority:
                  type: string
                  enum: [low, normal, high]
                  default: normal
      responses:
        '200': {description: ok}
        '400': {description: invalid_request_error}
        '413': {description: payload_too_large}
        '429': {description: rate_limit_error}
        '503': {description: overloaded_error}
        '504': {description: request_timeout}
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
security:
  - bearerAuth: []
```

---

## 5. Test cases (validation against actual fd)

Каждый test имеет: input, expected HTTP status, expected error code, expected headers.

### 5.1 Happy path (10 tests)

```
T-H-1: POST /v1/embeddings {"input":["hello"]} → 200, dimensions=1024 in response
T-H-2: POST /v1/embeddings {"input":["hello"],"dimensions":512} → 200, dimensions=512
T-H-3: POST /v1/embeddings {"input":["a","b","c"]} → 200, 3 embeddings
T-H-4: POST /v1/embeddings {"input":["a"]*32} → 200, 32 embeddings (max batch)
T-H-5: POST /v1/embeddings {"input":["a"],"encoding_format":"base64"} → 200, base64 string
T-H-6: POST /v1/embeddings {"input":["a"],"priority":"high"} → 200
T-H-7: GET /health → 200, body contains model_loaded: true
T-H-8: GET /live → 200
T-H-9: GET /ready → 200 (after warmup)
T-H-10: GET /version → 200, body contains version field
```

### 5.2 Error path (15 tests)

```
T-E-1: POST /v1/embeddings {} → 400, code=input_required
T-E-2: POST /v1/embeddings {"input":[]} → 400, code=input_required
T-E-3: POST /v1/embeddings {"input":["a"],"dimensions":99999} → 400, code=dimensions_invalid
T-E-4: POST /v1/embeddings {"input":[123]} → 400, code=invalid_request_error, NOT "json: cannot unmarshal"
T-E-5: POST /v1/embeddings {malformed → 400, code=invalid_json
T-E-6: POST /v1/embeddings {"input":["a"]*100} → 413, code=batch_too_large
T-E-7: POST /v1/embeddings {"input":["x"*10000]} → 413, code=input_too_long
T-E-8: GET /v1/embeddings → 405, NOT 404
T-E-9: POST /v1/embeddings {"input":["a"]} (no auth, если R-P1-8) → 401, code=unauthorized
T-E-10: GET /v9999 → 404, code=not_found
T-E-11: POST /v1/embeddings (during shutdown) → 503, code=shutting_down, Retry-After: 30
T-E-12: POST /v1/embeddings (model not loaded) → 503, code=model_not_loaded, Retry-After: 5
T-E-13: POST /v1/embeddings (rate limit hit) → 429, code=rate_limit_exceeded, Retry-After: 60
T-E-14: POST /v1/embeddings (oversized body 50MB) → 413, code=payload_too_large
T-E-15: 500 scenario (force internal error) → 500, code=internal_error, X-Request-Id in body
```

### 5.3 Headers (10 tests)

```
T-HDR-1: любой response → Server: fd/2.0.0
T-HDR-2: caller passes X-Request-Id: my-id → response echoes my-id
T-HDR-3: caller doesn't pass X-Request-Id → response has X-Request-Id: <uuid>
T-HDR-4: POST /v1/embeddings response → X-Model-Id: deepvk/USER-bge-m3
T-HDR-5: POST /v1/embeddings response → X-Dimensions: 1024 (or 512)
T-HDR-6: repeat same input → X-Cache: HIT (если R-P1-4)
T-HDR-7: first request → X-Cache: MISS
T-HDR-8: 429/503 response → Retry-After: <seconds>
T-HDR-9: любой response → Connection: keep-alive
T-HDR-10: cache hit response → ETag: <hash>
```

### 5.4 Performance (5 tests)

```
T-P-1: 1 input (warm) → p95 < 50ms
T-P-2: 10 inputs (warm) → p95 < 200ms
T-P-3: 32 inputs (warm, max batch) → p95 < 1000ms
T-P-4: 100 sequential requests → 0 errors, 0 timeouts
T-P-5: concurrent 4 callers × 8 inputs each → all succeed, total time < 2s
```

### 5.5 Endpoints existence (5 tests)

```
T-E-1: GET /version → 200 (NOT 404)
T-E-2: GET /info → 200 (NOT 404)
T-E-3: GET /metrics → 200 (NOT 404), Content-Type: text/plain
T-E-4: GET /openapi.json → 200 (NOT 404)
T-E-5: GET /docs → 200 (NOT 404)
```

---

## 6. Behavior scenarios (как fd должен вести себя)

### 6.1 Startup sequence

```
1. Process starts
2. Initialize HTTP server (bind port 8000)
3. /live returns 200 (process alive)
4. /ready returns 503 (warmup not done)
5. /health returns 503, status="down", warmup_done=false
6. Load model into GPU (async, may take 5-30s)
7. Run 1 dummy inference (warmup)
8. /ready returns 200 (warmup done)
9. /health returns 200, status="ok", warmup_done=true
10. Accept normal traffic
```

### 6.2 Normal request

```
1. Caller: POST /v1/embeddings {"input":["text"],"dimensions":1024,"user":"x","priority":"normal"}
2. Server: validate input (non-empty, strings only, length OK, batch OK, dims OK)
3. Server: validate auth (если R-P1-8)
4. Server: check model loaded + warmup done (если нет → 503)
5. Server: check rate limit (если R-P2-5, если превышен → 429)
6. Server: check cache (если R-P1-4, если hit → return cached + X-Cache: HIT)
7. Server: run inference (timeout REQUEST_TIMEOUT seconds, e.g., 30s)
8. Server: return 200 + embeddings + X-Cache: MISS + metrics update
9. Server: log with X-Request-Id
```

### 6.3 Failure scenarios

**F-1: Model not loaded (startup race)**
```
Caller: POST /v1/embeddings
Server: model not loaded
Response: 503, code=model_not_loaded, Retry-After: 5
Caller: waits 5s, retries
Server: model loaded
Response: 200
```

**F-2: Model overloaded (concurrent > capacity)**
```
Server: queue full or worker pool exhausted
Response: 503, code=model_overloaded, Retry-After: 5
Caller: retries with backoff
```

**F-3: Request timeout (long inference)**
```
Server: inference > REQUEST_TIMEOUT
Server: cancel inference, return 504, code=request_timeout, Retry-After: 1
```

**F-4: Cache miss, then hit**
```
Caller: POST /v1/embeddings {"input":["hello"]}
Response: 200, X-Cache: MISS
Caller: POST /v1/embeddings {"input":["hello"]}
Response: 200, X-Cache: HIT, latency < 5ms
```

**F-5: Graceful shutdown**
```
Operator: kill -TERM <pid>
Server: receives SIGTERM
Server: stop accepting new requests (return 503 with code=shutting_down)
Server: wait for in-flight requests to complete (max 30s)
Server: exit 0
```

**F-6: Invalid input (caller bug)**
```
Caller: POST /v1/embeddings {"input":[123]}
Server: validate, type check fails
Response: 400, code=invalid_request_error, param=input, message="input[0] must be string"
Caller: fixes input, retries
```

**F-7: Batch too large (caller bug)**
```
Caller: POST /v1/embeddings {"input":["a"]*100}
Server: validate, batch > 32
Response: 413, code=batch_too_large, message="batch size 100 exceeds max 32"
Caller: splits into chunks of 32, retries
```

**F-8: Input too long (caller bug)**
```
Caller: POST /v1/embeddings {"input":["x"*10000]}
Server: validate, input[0] > 512 tokens (~2048 chars)
Response: 413, code=input_too_long, message="input[0] exceeds max length 512 tokens (got ~2500 tokens)"
Caller: truncates, retries
```

**F-9: Server bug (unexpected panic)**
```
Server: panic during inference
Server: recover, log with stack trace
Response: 500, code=internal_error, X-Request-Id in message
Server: increment fd_errors_total{code=internal_error}
```

**F-10: Caller passes X-Request-Id**
```
Caller: POST /v1/embeddings, X-Request-Id: my-trace-123
Server: echo X-Request-Id: my-trace-123 in response
Server: log: [request_id=my-trace-123] ...
```

---

## 7. Implementation guidance (architecture hints)

### 7.1 Tech stack (likely Go, based on error messages)

```go
// Server skeleton
type FDServer struct {
    model         *Model          // loaded model
    cache         *LRUCache       // LRU cache
    rateLimiter   *RateLimiter    // per-IP/user rate limit
    metrics       *Metrics        // Prometheus counters/histograms
    shuttingDown  atomic.Bool     // graceful shutdown flag
    warmupDone    atomic.Bool     // warmup complete flag
    inflight      sync.WaitGroup  // in-flight requests
}

// Handler chain
http.Handler → Middleware(X-Request-Id, Metrics, Auth, RateLimit, Validation) → Endpoint handler
```

### 7.2 Critical concurrency primitives

- `sync.RWMutex` для cache reads/writes
- `atomic.Bool` для `shuttingDown`, `warmupDone`
- `sync.WaitGroup` для tracking in-flight requests при shutdown
- `context.Context` для request timeout cancellation

### 7.3 Model loading pattern

```go
// Startup
go func() {
    s.model = loadModel(config)
    s.warmupDone.Store(true)  // signal /ready
}()

// Per request
if !s.warmupDone.Load() {
    return 503, "model_not_loaded"
}
```

### 7.4 Shutdown pattern

```go
// On SIGTERM
s.shuttingDown.Store(true)
s.inflight.Wait()  // wait for in-flight (with timeout)
httpServer.Shutdown(ctx)
```

### 7.5 Validation pattern (BEFORE model call)

```go
func validateRequest(req EmbedRequest) error {
    if len(req.Input) == 0 {
        return ErrInputRequired
    }
    if len(req.Input) > MaxBatchSize {
        return ErrBatchTooLarge
    }
    for i, t := range req.Input {
        if !isString(t) {
            return ErrInputNotString(i)
        }
        if len(t) > MaxInputLength {
            return ErrInputTooLong(i, len(t))
        }
    }
    if req.Dimensions != 512 && req.Dimensions != 1024 {
        return ErrDimensionsInvalid
    }
    return nil
}
```

### 7.6 Error response builder

```go
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}
type ErrorDetail struct {
    Code    string `json:"code"`
    Type    string `json:"type"`
    Param   string `json:"param,omitempty"`
    Message string `json:"message"`
}
```

### 7.7 Metrics middleware

```go
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w, r) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        status := strconv.Itoa(getStatusCode(w))
        
        fdRequestsTotal.WithLabelValues(status).Inc()
        fdRequestDuration.Observe(duration.Seconds())
    })
}
```

---

## 8. Out of scope (NOT part of fd v2)

- ❌ Replace BGE-M3 with another model (separate milestone)
- ❌ Multi-model support (F-P2-3 batch endpoint is OK, but multi-model = separate)
- ❌ Persistent storage (cache is in-memory only)
- ❌ Distributed tracing (R-P2-6 traces endpoint is in-memory, not distributed)
- ❌ Authentication beyond API key (no OAuth, no JWT)
- ❌ Multi-tenant isolation
- ❌ GPU multi-tenancy
- ❌ Auto-scaling

---

## 9. Migration path (от fd v1 к fd v2)

### 9.1 Backward compatibility

- v1 caller (только POST /v1/embeddings с OpenAI shape) должен продолжать работать.
- v2 только ДОБАВЛЯЕТ endpoints, headers, validation, errors.
- v2 МЕНЯЕТ error format (caller должен обновиться, чтобы парсить `{"error":{"code":...}}`).
- v2 ADD headers (caller может игнорировать).

### 9.2 Rollout strategy

1. **Week 1**: реализовать R-P0-1..R-P0-6 (P0 functional), deploy как fd v2.0.0-beta
2. **Week 2**: реализовать R-P0-7..R-P0-19 (P0 observability + headers + errors), deploy v2.0.0
3. **Week 3**: реализовать R-P1-* (P1 health + features), deploy v2.1.0
4. **Week 4+**: реализовать R-P2-* (P2 nice-to-have), по запросу

### 9.3 Daily-archive integration

- M062 S01: wrapper обновляется, retries/circuit breaker/graceful degradation
- M062 S02: ADR-019 binding contract
- M062 S03: contract tests validate fd v2 (запускаются daily-archive, не fd)
- После deploy fd v2.0.0: daily-archive wrapper использует новые headers (X-Request-Id, X-Cache, etc.)

---

## 10. LLM Reading Notes

- **Read this once before implementing fd v2.** This document is the single source of truth for fd v2 requirements.
- **Section 1 (current state)**: complete list of bugs and missing endpoints. Do not re-discover — trust the probe results.
- **Section 2 (P0/P1/P2)**: prioritized requirements. R-P0-* are MUST, R-P1-* are SHOULD, R-P2-* are COULD.
- **Section 3 (error catalog)**: machine-readable error contract. Use exactly these codes, types, and HTTP statuses.
- **Section 4 (OpenAPI)**: implement endpoints matching this spec.
- **Section 5 (test cases)**: 45 test cases total. Use as acceptance criteria.
- **Section 6 (behavior scenarios)**: 10 scenarios showing how fd should respond. Implement recovery logic.
- **Section 7 (implementation)**: Go-flavored hints (since fd is likely Go based on error messages). Adapt to actual stack.
- **Section 8 (out of scope)**: explicit non-goals. Do not add.
- **Section 9 (migration)**: backward compatibility strategy. v1 callers must work unchanged in v2.0.0.
- **Section 10 (this section)**: navigation guide.

**Tone**: precise, machine-actionable, no marketing prose. All requirements testable. All errors machine-readable.

**Gotchas**:
- Do not invent endpoints. Implement exactly R-P0-7..R-P0-10, R-P1-1..R-P1-3.
- Do not invent error codes. Use exactly section 3 catalog.
- Do not change request shape (OpenAI compat is binding).
- Do not add authentication without R-P1-8 (env var FD_API_KEY).
- Cache is optional (R-P1-4), but if implemented, must follow R-P0-15 (X-Cache header).
- Health check must be DEEP (R-P1-1), not just process alive.
- Warmup is mandatory (R-P0-4) — do not skip pre-warm.

**Cross-references**:
- Daily-archive wrapper spec: M062 S01 (src/arxiv_archive/embedder.py + new tests)
- Daily-archive contract tests: M062 S03 (scripts/test_fd_contract.py)
- Daily-archive ADR: ADR-019 (M034 template, 14 sections, LLM Reading Notes)
