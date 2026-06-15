# FD Embedding Service

FD is a same-host text embedding service for legal-domain workloads. It exposes an OpenAI-compatible embeddings API, uses HuggingFace Text Embeddings Inference (TEI) for inference, and keeps hot results in a two-tier cache: in-process L1 plus Redis L2.

The current production runtime is **TEI-only** with `deepvk/USER-bge-m3` at **1024 dimensions**. ONNX runtime work is intentionally inactive/future research; do not treat ONNX as part of the current deployment path.

## Table of Contents

- [Status](#status)
- [Quick Start](#quick-start)
- [What This Service Provides](#what-this-service-provides)
- [Architecture](#architecture)
- [API](#api)
- [Configuration](#configuration)
- [Operations](#operations)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [Performance](#performance)
- [Documentation Map](#documentation-map)
- [Limitations and Non-Goals](#limitations-and-non-goals)
- [License](#license)

## Status

| Area | Current value |
|------|---------------|
| Runtime backend | TEI |
| Public model identity | `deepvk/USER-bge-m3` |
| Embedding dimensions | `1024` |
| Cache | L1 local cache + L2 Redis binary cache |
| TEI startup mode | Cached local model snapshot path |
| Public API style | OpenAI-compatible `/v1/embeddings` |
| ONNX runtime | Not active; future research only |

The TEI container is configured to start from the cached local snapshot path:

```text
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae
```

This avoids the Hugging Face Hub repo-resolution path during startup while preserving `deepvk/USER-bge-m3` as the public model identity reported by fd.

## Quick Start

Start the local Docker stack:

```bash
docker compose up -d
```

Check container health:

```bash
docker compose ps
```

Check fd runtime metadata:

```bash
curl http://localhost:8000/health
```

Expected runtime block:

```json
{
  "runtime": {
    "backend": "tei",
    "model": "deepvk/USER-bge-m3",
    "dimensions": 1024,
    "production_default": true,
    "cache_namespace": "v2"
  }
}
```

Run a smoke embedding request:

```bash
curl -s http://localhost:8000/v1/embeddings \
  -H 'Content-Type: application/json' \
  -d '{"model":"deepvk/USER-bge-m3","input":"short non-sensitive smoke text"}'
```

A valid response contains one 1024-dimensional embedding vector in `data[0].embedding`.

## What This Service Provides

- OpenAI-compatible text embeddings endpoint for same-host consumers.
- FalkorDB-oriented batch embedding endpoint.
- Explicit runtime metadata through `/health` and `/info`.
- Readiness/liveness endpoints for process lifecycle checks.
- Two-tier embedding cache:
  - L1: in-process local cache for very hot requests.
  - L2: Redis binary cache for persistence across API restarts.
- OpenAI-style error envelopes for validation, capacity, and internal failures.
- Optional API-key auth, CORS, rate limiting, metrics, and trace endpoints.

## Architecture

```text
client
  │
  ▼
fd API (Go/Gin)
  ├─ validation / auth / rate-limit / lifecycle middleware
  ├─ L1 local cache
  ├─ L2 Redis cache
  └─ TEI client
        │
        ▼
HuggingFace TEI CPU container
  └─ local USER-bge-m3 snapshot mounted at /data
```

Runtime identity is intentionally separated from TEI's internal startup path:

- fd reports the public model as `deepvk/USER-bge-m3`.
- TEI starts from a local `/data/.../snapshots/...` directory.
- Cache namespace metadata is exposed so operators can detect cache identity drift.

The canonical same-host contract is maintained in [`docs/same-host-embedding-service-contract.md`](docs/same-host-embedding-service-contract.md). Keep new consumer-facing behavior aligned with that document instead of duplicating the full contract here.

## API

### Runtime and lifecycle

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/live` | Process liveness. |
| `GET` | `/ready` | Readiness; includes warmup and shutdown state. |
| `GET` | `/health` | Health and runtime metadata. |
| `GET` | `/v1/healthcheck` | Compatibility alias for `/health`. |
| `GET` | `/info` | Build, runtime, and lifecycle info. |
| `GET` | `/version` | Build/version metadata. |
| `GET` | `/metrics` | Prometheus-style metrics. |
| `GET` | `/v1/traces` | Recent in-process trace events when enabled. |
| `GET` / `POST` | `/warmup` | Inspect or trigger model warmup. |
| `GET` | `/openapi.json` | OpenAPI specification. |
| `GET` | `/docs` | API documentation surface. |

### `POST /v1/embeddings`

OpenAI-compatible embeddings request.

Request:

```json
{
  "model": "deepvk/USER-bge-m3",
  "input": ["first text", "second text"]
}
```

Response shape:

```json
{
  "object": "list",
  "model": "deepvk/USER-bge-m3",
  "data": [
    {
      "object": "embedding",
      "index": 0,
      "embedding": [0.01, -0.02]
    }
  ]
}
```

The request `model` field is compatibility metadata. Consumers should treat the response `model` and `/health.runtime.model` as authoritative runtime identity.

### `POST /v1/batch`

Batch endpoint for v1 clients. It uses the same lifecycle capacity gate as `/v1/embeddings`.

### `POST /embeddings/batch`

Batch embeddings endpoint used by FalkorDB-oriented local workflows.

## Configuration

The Docker stack sets safe local defaults. Override values through `docker-compose.override.yaml`, shell environment, or `api/.env` when running locally.

### Core runtime

| Variable | Default | Purpose |
|----------|---------|---------|
| `EMBEDDING_BACKEND` | `tei` | Active backend. Only `tei` is supported. Any other value fails startup. |
| `MODEL_ID` | `deepvk/USER-bge-m3` | Public model identity returned by fd. |
| `TEI_URL` | `http://tei:80` | Internal TEI endpoint used by fd. |
| `PORT` | `8000` | fd API port inside the container/process. |
| `BIND_HOST` | `0.0.0.0` | fd API bind host. |
| `LOG_LEVEL` | `info` | Structured log level. |
| `GIN_MODE` | unset / compose uses `release` | Gin runtime mode. |

### Redis and cache

| Variable | Default | Purpose |
|----------|---------|---------|
| `REDIS_HOST` | `redis:6379` | Redis L2 cache address. |
| `REDIS_POOL_SIZE` | `50` | Redis client pool size. |
| `REDIS_CACHE_TTL` | `24h` | Redis cache TTL. Must be a Go duration, for example `30m`, `24h`, `168h`. |
| `REDIS_CACHE_NO_EXPIRE` | `false` | If `true`, cached embeddings do not expire. Mutually exclusive with `REDIS_CACHE_TTL`. |
| `EMBEDDING_CACHE_VERSION` | `v2` | Cache namespace version. Change this to isolate incompatible cache populations. |
| `EMBEDDING_MODEL_ID` | empty | Optional model namespace component for cache keys. |
| `EMBEDDING_MODEL_REVISION` | empty | Optional model revision namespace component for cache keys. |
| `EMBEDDING_TOKENIZER_VERSION` | empty | Optional tokenizer namespace component for cache keys. |
| `EMBEDDING_CHUNKING_VERSION` | empty | Optional chunking namespace component for cache keys. |
| `REDIS_MAXMEMORY` | `2gb` | Redis container maxmemory setting. |
| `REDIS_MAXMEMORY_POLICY` | `allkeys-lru` | Redis eviction policy. |
| `REDIS_RDB_SAVE` | `300 1` | Redis RDB save policy. |
| `REDIS_AOF_ENABLED` | `no` | Redis append-only file setting. |

Cache namespace isolation matters. When comparing runtime backends or changing model/tokenizer/chunking behavior, set a new namespace value or deliberately flush Redis. Otherwise cached vectors can make incompatible backends appear equivalent.

### Controls and middleware

| Variable | Default | Purpose |
|----------|---------|---------|
| `FD_MAX_IN_FLIGHT` | `0` | Maximum concurrent embedding requests. `0` preserves unlimited behavior; positive values return `503 model_overloaded` when full. |
| `FD_API_KEY` | unset | Required for protected endpoints. If unset, protected endpoints fail closed with `401 unauthorized`; `/live`, `/ready`, `/health`, `/v1/healthcheck`, docs, and OpenAPI stay public. Do not log or commit secret values. |
| `FD_CORS_ORIGINS` | unset | Optional allowed CORS origins. |
| `FD_RATE_LIMIT_ENABLED` | unset / false | Enables IP/user rate limiting when configured. |
| `FD_TRACES_ENABLED` | unset / false | Enables `/v1/traces` diagnostic events. |

## Operations

### TEI startup

TEI is intentionally started with a local model directory:

```yaml
command: --model-id /data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae --max-batch-tokens 8192
```

Keep `docker-compose.yaml` and `docker-compose.override.yaml` aligned. Compose `command` replacement is total; an override command silently replaces the base command rather than merging with it.

The validated M045 proof showed this local-path mode reaches TEI health and fd readiness without relying on the failed `HF_HUB_OFFLINE=1` mitigation. Details are in [`documents/tei-startup-mitigation-m045.md`](documents/tei-startup-mitigation-m045.md) and [`benchmark-results/m045-tei-local-path-startup-proof.md`](benchmark-results/m045-tei-local-path-startup-proof.md).

### Health and readiness

Use these checks in order:

```bash
docker compose ps
curl http://localhost:8000/live
curl http://localhost:8000/ready
curl http://localhost:8000/health
```

`/health` exposes runtime metadata but is not a live inference probe. For readiness that matters to clients, also run a small `/v1/embeddings` smoke request.

### Logs

fd writes structured JSON logs to stdout. TEI writes its own startup/backend logs through the container runtime.

```bash
docker compose logs --tail=100 api
docker compose logs --tail=100 tei
```

Do not log request bodies or secret values in operational diagnostics.

### Redis exposure and persistence

The override binds Redis to localhost only:

```yaml
127.0.0.1:6379:6379
```

Redis uses `redis-stack:latest` with a bounded memory policy. The compose defaults favor local development and benchmark repeatability, not public network exposure.

## Troubleshooting

### TEI is unhealthy or stuck during startup

1. Confirm the TEI command uses the local `/data/.../snapshots/...` path, not `deepvk/USER-bge-m3` as a Hub ID.
2. Confirm both compose files are aligned:

   ```bash
   docker compose config tei
   ```

3. Inspect TEI logs:

   ```bash
   docker compose logs --tail=200 tei
   ```

4. A local missing-ONNX warning is acceptable if TEI falls through to Candle backend and becomes healthy. A long remote `Downloading onnx/model.onnx` path indicates the local-path mitigation is not active.

### fd `/health` is OK but embeddings fail

`/health` is metadata, not a live inference probe. Check:

```bash
curl http://localhost:8000/ready
curl -s http://localhost:8000/v1/embeddings \
  -H 'Content-Type: application/json' \
  -d '{"model":"deepvk/USER-bge-m3","input":"smoke"}'
```

Then inspect API and TEI logs.

### Cache results look suspicious

If a backend/model/tokenizer/chunking change is under test, isolate the Redis namespace with `EMBEDDING_CACHE_VERSION` or related namespace variables. Cross-backend cache reuse can produce false-positive equivalence.

### API returns `503 model_overloaded`

`FD_MAX_IN_FLIGHT` is set to a positive value and the lifecycle capacity gate is full. Either reduce concurrency, increase the configured capacity, or leave `FD_MAX_IN_FLIGHT=0` for unlimited current behavior.

### API returns auth or CORS errors

Check whether `FD_API_KEY` or `FD_CORS_ORIGINS` is set in `api/.env`, shell environment, or compose override. Never commit secret values.

## Development

Go module root:

```bash
cd api
```

Run the regular API test suite:

```bash
go test ./...
```

Run the CI-equivalent short suite:

```bash
go test ./... -short
```

Run lint with the pinned tool version used by the project:

```bash
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...
```

Run vulnerability analysis:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Run the Docker stack from the repository root:

```bash
docker compose up -d
```

Run the black-box integration suite from its standalone module:

```bash
cd tests/integration && go test -v .
```

Without `FD_INTEGRATION_API_KEY`, the suite verifies public diagnostics and auth fail-closed behavior, while protected happy-path checks skip. For full authenticated Docker e2e coverage, recreate or run the API with a matching local `FD_API_KEY`, then pass the same value as `FD_INTEGRATION_API_KEY` without printing it:

```bash
cd tests/integration && FD_INTEGRATION_API_KEY=<matching local key> go test -v .
```

Run the bounded mutation baseline from `api/` when checking assertion strength on critical cache, handler, and lifecycle files:

```bash
go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest \
  --exec 'go test ./cache ./handlers ./lifecycle' \
  --exec-timeout 45 \
  ./cache/hash.go ./cache/keys.go ./handlers/cache.go ./handlers/health.go ./lifecycle/state.go
```

This mutation command is an informational local/manual gate, not a mandatory CI hard gate. It currently requires the Go 1.25.11 toolchain via automatic toolchain switching and takes roughly a minute on the bounded scope.

The current M050 closeout evidence includes 295 passing Go tests, 0 lint issues, 0 reachable vulnerabilities, an authenticated Docker e2e suite with 9 passing checks, and a bounded mutation score of 1.0 on 143 mutants in scope.

## Performance

Validated local benchmark snapshots are committed under [`benchmark-results/`](benchmark-results/). Treat them as local-environment snapshots, not universal production SLOs.

| Metric | Latest documented local range / value |
|--------|----------------------------------------|
| Warm cache latency | ~2 ms |
| Cold TEI latency | 20-200 ms local range |
| Cache speedup | 40-80x local range |
| Max throughput | ~644 req/s in latest final run |
| Redis L2 after API restart | ~3 ms in latest diagnostic run |
| Storage per embedding | ~4 KB binary |

Run benchmark tooling from the host against a local Docker stack. Some older benchmark scripts expect `uv` and Python 3.13.

## Documentation Map

| Document | Purpose |
|----------|---------|
| [`docs/same-host-embedding-service-contract.md`](docs/same-host-embedding-service-contract.md) | Canonical HTTP/runtime contract for local consumers. |
| [`docs/fd-v2.md`](docs/fd-v2.md) | Broader fd v2 product/technical context. |
| [`documents/tei-startup-mitigation-m045.md`](documents/tei-startup-mitigation-m045.md) | TEI startup mitigation outcome and operator notes. |
| [`documents/tei-startup-recon-m045.md`](documents/tei-startup-recon-m045.md) | Non-destructive TEI startup/runtime recon. |
| [`benchmark-results/m045-tei-local-path-startup-proof.md`](benchmark-results/m045-tei-local-path-startup-proof.md) | Controlled proof for the local snapshot startup mode. |
| [`benchmark-results/`](benchmark-results/) | Local benchmark and validation artifacts. |
| [`api/openapi/spec.go`](api/openapi/spec.go) | Source for generated OpenAPI response. |

## Limitations and Non-Goals

- ONNX runtime is not active in the current fd service. Existing ONNX materials are historical or future research unless explicitly reactivated by a new milestone.
- `/health` is not a live inference probe; use `/ready` plus a small embedding request for operational readiness.
- The TEI local snapshot path is pinned. If the cached model revision changes, update compose/docs together and rerun the startup proof.
- The service is designed for same-host/local consumers, not public internet exposure by default.
- Benchmark numbers are environment-specific local snapshots.

## License

No license file is currently present in this repository. Add a `LICENSE` file before distributing or accepting external contributions.
