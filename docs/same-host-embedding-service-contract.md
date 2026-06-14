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

The current fd build is TEI-only. ONNX-specific runtime metadata fields are reserved historical/schema fields and must not appear for the TEI runtime. ONNX is not a current operator-selectable backend.

---

## 2. Runtime and Environment

### Production default: TEI

The production and default runtime is **TEI** (Text Embeddings Inference, Hugging Face) served at `TEI_URL`. TEI runs as a separate container or process and is called by the Go API over HTTP.

- Set `TEI_URL` to point at the TEI service (default: `http://tei:80`).
- Set `MODEL_ID` to the target Hugging Face model (default: `deepvk/USER-bge-m3`).
- Set `TEI_URL` and `MODEL_ID` for the active runtime.
- `EMBEDDING_BACKEND` may be omitted or set to `tei`; any other value is rejected. `EMBEDDING_BACKEND=onnx` is intentionally disabled in the current product scope.

### Historical: ONNX research

Earlier milestones evaluated ONNX as a local performance candidate, but ONNX did not pass the packaging and operational-readiness gates required for the current product path. ONNX artifacts and benchmark files are retained only as historical research evidence. Reintroducing ONNX requires a separate milestone with artifact distribution, tokenizer/runtime-library packaging, Docker/CI, cache-isolation, quality, and performance gates.

---

## 3. Health Metadata Semantics

### What `/health` proves

| Runtime | What `/health` confirms |
|---------|------------------------|
| TEI (current/default) | API process started, env parsed, Redis responded to ping. TEI reachability is **not** actively probed before the server starts. |

### What `/health` does NOT prove

- **Live inference health.** The health endpoint does not probe the TEI inference path as part of `/health`. For TEI, this means `status: ok` can appear even if the upstream TEI service is unreachable, as long as the Go API itself started.
- **Vector correctness.** `/health` does not return or validate embedding values.
- **Cache health.** Redis connectivity is confirmed at startup via a ping, but the cache layer is not probed per request.

### Recommended client pattern

1. Call `GET /health` and verify `status: ok`.
2. Read `runtime.backend` and expect `tei` for the current product build.
3. For full end-to-end readiness, send a smoke `POST /v1/embeddings` request with a short non-sensitive string and verify the response shape.

---

## 4. Timeout and Retry Guidance

### Request timeouts

| Endpoint | Timeout | Reason |
|----------|---------|--------|
| `/v1/embeddings` | 30 s | Per-request context deadline; catches hung TEI calls |
| `/embeddings/batch` | 120 s | Larger payloads and multiple TEI calls accumulate |
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

Current product builds use TEI only, so normal operation does not mix runtime backends. Historical TEI/ONNX research showed that different runtimes can produce different vectors for the same input text; any future backend comparison must isolate `EMBEDDING_CACHE_VERSION` or namespace env vars so stale cached vectors cannot contaminate evidence.

Benchmark artifacts record the effective namespace configuration to make research isolation auditable.

---

## 6. No-Silent-Fallback Rules

The service enforces a TEI-only current runtime:

1. **Backend selection is effectively fixed to TEI.** `EMBEDDING_BACKEND` may be omitted or set to `tei`; any other value is rejected at startup.
2. **No per-request backend fallback.** A running instance does not switch runtimes mid-flight. A request either goes to TEI or returns an error.
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
| `500` | `{"error": "embedding generation failed"}` | TEI inference error; cache error |
| `503` | — | Not currently emitted by the service; reserved for future restart or startup probing |

### `embedding generation failed` (500)

This error means the TEI inference call or cache path returned an error. It does not distinguish between:

- TEI service unreachable.
- TEI runtime/backend error.
- Tokenizer/model mismatch inside TEI.
- Out-of-memory condition.

Check TEI service logs and Go API logs for the underlying error message.

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
| Backend detection | `GET /health` → `runtime.backend` (`"tei"` for current product builds) |
| Full readiness | Smoke `POST /v1/embeddings` → verify `data[0].dimensions` in response |
| Request timeout | 30 s for single/batch embeddings; 10 s minimum for health |
| Retry on | Transport errors, HTTP 503 |
| Do not retry on | HTTP 400, HTTP 500 |
| Cache compatibility | Verify `runtime.cache_namespace` matches your expected namespace |
| Future backend research | Use a separate milestone and isolated `EMBEDDING_CACHE_VERSION`; current product builds are TEI-only |
| Logging | Set `LOG_LEVEL=debug` to see cache hit/miss events (short key hashes only, no raw text) |

---

## 11. Related Documents

- [README.md](../README.md) — Quick start, configuration reference, and Docker Compose setup
- [M040 S04 runtime recommendation](../benchmark-results/fd-runtime-recommendation-m040-s04.md) — Final TEI-vs-ONNX same-host runtime recommendation and evidence envelope
- [docs/onnx-artifacts/OPERATIONS.md](./onnx-artifacts/OPERATIONS.md) — ONNX-specific preflight, health, and rollback contract
- [docs/onnx-artifacts/PROVISIONING.md](./onnx-artifacts/PROVISIONING.md) — ONNX artifact build and provisioning
- [benchmark.py](../benchmark.py) — Benchmark harness; uses `BENCHMARK_API_RESTART_COMMAND` for restart/cache proof
