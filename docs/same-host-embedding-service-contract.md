# Same-Host Embedding Service Contract

**Version:** 1.0
**Date:** 2026-05-21
**Scope:** Local HTTP consumers on the same host as `fd`

This document is the canonical consumer contract for local services calling `fd`'s HTTP API on the same machine. It covers endpoint shapes, status and error codes, encoding caveats, runtime identification via `/health`, timeout and retry guidance, cache namespace expectations, and explicit non-goals.

---

## 1. Endpoints

### `GET /health`

Returns a JSON object with a `status` field and optional `runtime` metadata block.

**Response shape:**

```json
{
  "status": "ok",
  "time": "2026-05-21T10:00:00Z",
  "runtime": {
    "backend": "tei",
    "model": "deepvk/USER-bge-m3",
    "dimensions": 1024,
    "production_default": true,
    "cache_namespace": "embed:cache:v2:"
  }
}
```

When `EMBEDDING_BACKEND=onnx` (opt-in), the `runtime` block is present and additionally contains:

| Field | Type | Description |
|-------|------|-------------|
| `backend` | string | `"tei"` or `"onnx"` |
| `model` | string | Model identifier (e.g. `deepvk/USER-bge-m3`) |
| `artifact_id` | string | ONNX artifact identifier (ONNX only) |
| `dimensions` | int | Embedding dimension: `1024` or `512` |
| `max_sequence_length` | int | Configured max sequence length |
| `validated_max_sequence_length` | int | Artifact-validated max sequence length |
| `production_default` | bool | `true` for the production/default backend |
| `artifact_verified` | bool | ONNX manifest SHA256 verified (ONNX only) |
| `tokenizer_verified` | bool | Tokenizer JSON size or SHA256 verified (ONNX only) |
| `runtime_library_verified` | bool | ONNX runtime `.so` SHA256 verified; requires `ONNX_RUNTIME_SHA256` env (ONNX only) |
| `provider` | string | Execution provider (ONNX only, currently `"CPUExecutionProvider"`) |
| `cache_namespace` | string | Redis key namespace prefix |

**TEI/default `/health` includes a safe `runtime` block.** For TEI, the block reports only operational identity fields: `backend`, `model`, `dimensions`, `production_default`, and `cache_namespace`. It does not expose TEI URLs, filesystem paths, tokens, signed URLs, raw input text, or secrets.

> **Readiness scope:** `/health` is an API liveness and configuration metadata surface. For TEI, it confirms the API process started, env was parsed, runtime metadata was wired, and Redis responded to a ping before the server began accepting connections. It does **not** perform a live TEI inference probe. Clients requiring end-to-end inference readiness should send a smoke embedding request and validate the response shape.

### `POST /v1/embeddings`

OpenAI-compatible single and batch embeddings endpoint.

**Request:**

```json
{
  "model": "deepvk/USER-bge-m3",
  "input": "short non-sensitive smoke text",
  "dimensions": 1024
}
```

- `model` *(optional string)*: OpenAI-compatibility metadata only. The service does **not** use this field for routing and does **not** validate it against the configured `MODEL_ID`; clients may omit it or send a placeholder. The response `model` field and `/health.runtime.model` are authoritative for the model actually served.
- `input`: String or array of strings. Single strings are automatically promoted to a single-element array.
- `dimensions` *(optional int)*: `1024` (default) or `512`. Any other value returns `400`.

**Response:**

```json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [0.0123, -0.0456, ...],
      "index": 0,
      "dimensions": 1024
    }
  ],
  "model": "deepvk/USER-bge-m3",
  "usage": {
    "prompt_tokens": 3,
    "total_tokens": 3
  }
}
```

- `data[*].dimensions` reflects the requested dimensions (`1024` or `512`).
- `model` in the response reflects the configured `MODEL_ID`, not the request `model` field.
- Token counts are estimated (character count / 4), not tokenized counts.

**Status codes:**

| Code | Meaning |
|------|---------|
| `200` | Success |
| `400` | Invalid `dimensions` or empty `input` |
| `500` | TEI/ONNX inference failure |

### `POST /embeddings/batch`

Internal batch endpoint, primarily for FalkorDB.

**Request:**

```json
{
  "inputs": ["text1", "text2"],
  "dimensions": 1024,
  "encoding_format": "base64"
}
```

- `dimensions`: `1024` (default) or `512`.
- `encoding_format`: `"base64"` (default, binary float32 little-endian → base64) or `"float"` (JSON-encoded float array as a string element). Any other value returns `400`.

**Response:**

```json
{
  "embeddings": ["AQAAAA==", "BBBBAA=="],
  "count": 2,
  "dimensions": 1024
}
```

- `embeddings` is always a flat array of strings. With `encoding_format=base64` (default), each string is a base64-encoded binary float32 array in little-endian byte order. With `encoding_format=float`, each string is a JSON array literal (e.g. `"[0.0123,-0.0456]"`), **not** a nested JSON array.
- `dimensions` in the response reflects the requested or defaulted dimensions.
- `count` equals `len(inputs)`.

**Status codes:**

| Code | Meaning |
|------|---------|
| `200` | Success |
| `400` | Invalid `dimensions`, `encoding_format`, or empty `inputs` |
| `500` | Inference failure |

---

## 2. Runtime and Environment

### Production default: TEI

The production and default runtime is **TEI** (Text Embeddings Inference, Hugging Face) served at `TEI_URL`. TEI runs as a separate container or process and is called by the Go API over HTTP.

- Set `TEI_URL` to point at the TEI service (default: `http://tei:80`).
- Set `MODEL_ID` to the target Hugging Face model (default: `deepvk/USER-bge-m3`).
- Set `EMBEDDING_BACKEND=tei` explicitly or omit it (TEI is the default when `EMBEDDING_BACKEND` is unset).

### Opt-in: ONNX

The ONNX backend is opt-in and requires all of:

```bash
EMBEDDING_BACKEND=onnx
ONNX_ARTIFACT_MANIFEST=/path/to/user-bge-m3-dense-fp32.json
ONNX_RUNTIME_LIBRARY=/path/to/libonnxruntime.so
ONNX_TOKENIZER_PATH=/path/to/tokenizer.json
# Optional but required for runtime_library_verified=true:
ONNX_RUNTIME_SHA256=<64-char-hex>
ONNX_MAX_SEQUENCE_LENGTH=512  # or higher, up to validated_max_sequence_length
```

Startup fails with a non-zero exit code if any required ONNX env is missing or invalid. There is no fallback to TEI within a running ONNX process.

For full ONNX preflight, health, and rollback details, see [docs/onnx-artifacts/OPERATIONS.md](./onnx-artifacts/OPERATIONS.md).

---

## 3. Health Metadata Semantics

### What `/health` proves

| Backend | What `/health` confirms |
|---------|------------------------|
| TEI (default) | API process started, env parsed, Redis responded to ping. TEI reachability is **not** actively probed before the server starts. |
| ONNX (opt-in) | All of the above **plus** ONNX artifact manifest validated, tokenizer JSON verified, and optionally the runtime library SHA256 verified via `ONNX_RUNTIME_SHA256`. |

### What `/health` does NOT prove

- **Live inference health.** Neither backend probes the inference path (TEI HTTP call or ONNX `Embed` call) as part of `/health`. For TEI, this means `status: ok` can appear even if the upstream TEI service is unreachable, as long as the Go API itself started.
- **Vector correctness.** `/health` does not return or validate embedding values.
- **Cache health.** Redis connectivity is confirmed at startup via a ping, but the cache layer is not probed per request.

### Recommended client pattern

1. Call `GET /health` and verify `status: ok`.
2. Read `runtime.backend` to confirm which backend is active.
3. For ONNX, read `runtime.runtime_library_verified` — if `false`, the runtime library was not SHA256-verified; operation is still correct but the integrity boundary is weaker.
4. For full end-to-end readiness, send a smoke `POST /v1/embeddings` request with a short non-sensitive string and verify the response shape.

---

## 4. Timeout and Retry Guidance

### Request timeouts

| Endpoint | Timeout | Reason |
|----------|---------|--------|
| `/v1/embeddings` | 30 s | Per-request context deadline; catches hung TEI/ONNX calls |
| `/embeddings/batch` | 120 s | Larger payloads and multiple TEI/ONNX calls accumulate |
| `/health` | No dedicated timeout | Gin default; health is a metadata read, not inference |

### Retry rules

Clients should retry only on:

- **Transport errors** (`ECONNREFUSED`, `ETIMEDOUT`, HTTP 5xx) — transient networking or service startup issues.
- **HTTP 503** — service is starting or restarting.

Clients should **NOT** retry on:

- **HTTP 400** — bad request (invalid `dimensions`, empty `input`, etc.). Fix the request before retrying.
- **HTTP 500** — inference failure. Retry is unlikely to help without a service restart.
- **HTTP 4xx** — client error. Do not retry without fixing the request.

### Bounded retry example (pseudocode)

```
for attempt in [1, 2, 3]:
    try:
        response = POST /v1/embeddings with 30s timeout
        return response
    catch (transport error or HTTP 503) and attempt < 3:
        sleep(100ms * attempt)  # simple back-off
raise "embedding request failed after 3 attempts"
```

---

## 5. Cache Behavior

### Two-tier cache (L1 + L2)

| Tier | Technology | Typical latency | Notes |
|------|------------|-----------------|-------|
| L1 | sync.Map (in-process) | ~50 ns | Per-process, not shared across API instances |
| L2 | Redis binary | ~0.5–3 ms | Shared; survives API restarts; `embed:cache:<version>:` prefix |

Cache lookups are transparent to the consumer API contract: a `POST /v1/embeddings` response after a cache hit is byte-for-byte identical to a response after a cache miss. The embedding values do not change based on cache state.

### Cache namespace

Redis keys are namespaced by:

| Env var | Purpose |
|---------|---------|
| `EMBEDDING_CACHE_VERSION` | Schema/version namespace (default: `v2`) |
| `EMBEDDING_MODEL_ID` | Optional model namespace component |
| `EMBEDDING_MODEL_REVISION` | Optional revision namespace component |
| `EMBEDDING_TOKENIZER_VERSION` | Optional tokenizer namespace component |
| `EMBEDDING_CHUNKING_VERSION` | Optional chunking/splitting namespace component |

The effective namespace string appears in `/health` as `runtime.cache_namespace`. Clients sharing a Redis instance with other services should verify they use compatible or disjoint namespaces to avoid key collisions.

### Cache TTL and expiry

| Env var | Behavior |
|---------|---------|
| `REDIS_CACHE_TTL` | Expiry for L2 cache entries (Go duration, e.g. `24h`). Default: `24h`. |
| `REDIS_CACHE_NO_EXPIRE` | `true` disables expiry on L2 entries. Use for long-lived research caches. Do not set both `REDIS_CACHE_TTL` and `REDIS_CACHE_NO_EXPIRE`; the API rejects this at startup. |

### Cache contamination risk

TEI and ONNX backends can produce different embedding vectors for the same input text. If both backends write to the same Redis namespace, a subsequent request may retrieve a cached vector produced by the other backend, falsely appearing equivalent. **Always isolate `EMBEDDING_CACHE_VERSION` or namespace env vars when comparing or switching between TEI and ONNX.**

Benchmark artifacts record the effective namespace configuration to make cross-runtime isolation auditable.

---

## 6. No-Silent-Fallback Rules

The service enforces exactly one backend per process lifetime:

1. **Backend selection is startup-only.** `EMBEDDING_BACKEND` is read once at startup. Invalid or missing ONNX config causes a non-zero exit; there is no fallback to TEI within a running ONNX process.
2. **No per-request fallback.** A running instance never switches from TEI to ONNX or from ONNX to TEI mid-flight. A request either goes to the configured backend or returns an error.
3. **TEI internal fallback is not fd behavior.** TEI's own CPU/Candle fallback warnings (logged by TEI itself) are a TEI runtime concern, not an fd-level fallback. fd sends a request to the configured TEI URL and reports any TEI error as a `500`.
4. **Cache miss path is not fallback.** Fetching from L2 or L1 on a cache miss is a cache behavior detail, not a backend or model fallback.

Any change to `EMBEDDING_BACKEND` requires a service restart.

---

## 7. Status and Error Codes

### Common responses

| HTTP code | Body | When |
|-----------|------|------|
| `200` | Embeddings or batch response | Success |
| `400` | `{"error": "dimensions must be 1024 or 512"}` | Invalid `dimensions`, `encoding_format`, or empty `input`/`inputs` |
| `500` | `{"error": "embedding generation failed"}` | TEI or ONNX inference error; cache error |
| `503` | — | Not currently emitted by the service; reserved for future restart or startup probing |

### `embedding generation failed` (500)

This error means the inference backend (TEI HTTP call or ONNX `Embed` call) returned an error. It does not distinguish between:

- TEI service unreachable.
- ONNX runtime error (e.g. sequence length exceeded).
- Tokenizer mismatch.
- Out-of-memory condition.

For TEI, check TEI service logs. For ONNX, check the Go API logs for the underlying error message.

---

## 8. Encoding Caveats

### `encoding_format=base64` (batch, default)

Binary little-endian float32. To decode in Python:

```python
import base64, struct

def decode_base64_embedding(b64_str: str) -> list[float]:
    raw = base64.b64decode(b64_str)
    return list(struct.unpack(f"<{len(raw)//4}f", raw))
```

### `encoding_format=float` (batch)

JSON-encoded string containing a float array literal. The `embeddings` array contains **strings**, not JSON arrays. Each element is a quoted JSON array literal:

```json
{
  "embeddings": ["[0.0123,-0.0456]"],
  "count": 1,
  "dimensions": 1024
}
```

To decode, parse each string element as JSON — do not expect a top-level JSON array.

> **Breaking change risk:** The batch response shape (`embeddings` as `[]string`) is stable. However, `encoding_format=float` returning a string-ified JSON array rather than a nested array is a known shape quirk. Clients expecting a fully nested JSON array should use `encoding_format=base64`.

---

## 9. Non-Goals

The following are explicitly **out of scope** for the `fd` same-host embedding service:

| Non-goal | Reason |
|----------|--------|
| Request-level TEI ↔ ONNX fallback | Violates reproducibility and cache isolation guarantees |
| Per-request model selection | Model is a deployment-time configuration |
| Automatic tokenizer switching | Tokenizer is bound to the backend at startup |
| Live TEI/ONNX inference health probing via `/health` | Adds complexity and latency; smoke embedding is the correct readiness probe |
| Hosted or remote deployment proof | Scope is same-host local operation |
| Open-ended alternative model evaluation | Bounded by M040 S03 evidence gates; not a general model registry |
| Library/embedded integration surface | HTTP API only |
| INT8, NUMA, or OpenVINO runtime variants | Not in current evidence envelope |

---

## 10. Operational Summary for Local Clients

| Concern | Guidance |
|---------|---------|
| Endpoint base | `http://localhost:8000` (or the host's fd service address) |
| Backend detection | `GET /health` → `runtime.backend` (`"tei"` or `"onnx"`) |
| Full readiness | Smoke `POST /v1/embeddings` → verify `data[0].dimensions` in response |
| Request timeout | 30 s for single/batch embeddings; 10 s minimum for health |
| Retry on | Transport errors, HTTP 503 |
| Do not retry on | HTTP 400, HTTP 500 |
| Cache compatibility | Verify `runtime.cache_namespace` matches your expected namespace |
| TEI ↔ ONNX switch | Restart the service with new `EMBEDDING_BACKEND`; flush Redis or use a new `EMBEDDING_CACHE_VERSION` |
| Logging | Set `LOG_LEVEL=debug` to see cache hit/miss events (short key hashes only, no raw text) |

---

## 11. Related Documents

- [README.md](../README.md) — Quick start, configuration reference, and Docker Compose setup
- [docs/onnx-artifacts/OPERATIONS.md](./onnx-artifacts/OPERATIONS.md) — ONNX-specific preflight, health, and rollback contract
- [docs/onnx-artifacts/PROVISIONING.md](./onnx-artifacts/PROVISIONING.md) — ONNX artifact build and provisioning
- [benchmark.py](../benchmark.py) — Benchmark harness; uses `BENCHMARK_API_RESTART_COMMAND` for restart/cache proof
