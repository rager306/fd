# FD Embedding Service

[![Go Quality](https://github.com/rager306/fd/actions/workflows/go-quality.yml/badge.svg?branch=master)](https://github.com/rager306/fd/actions/workflows/go-quality.yml)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.2-blue)](./CHANGELOG.md)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-3.0-4baaaa.svg)](./CODE_OF_CONDUCT.md)

FD is a same-host text embedding service for Russian and legal-domain workloads. It exposes an OpenAI-compatible `POST /v1/embeddings` API, calls Hugging Face Text Embeddings Inference (TEI) for inference, and caches results in two tiers: an in-process L1 and a Redis L2.

The runtime is **TEI-only** with `deepvk/USER-bge-m3`. Vectors are 1024-dimensional; callers can request `dimensions: 512` to get the first 512 components. ONNX is not part of the current deployment path; the remaining ONNX files are historical research artifacts.

| Area | Current value |
|------|---------------|
| Inference backend | TEI (`ghcr.io/huggingface/text-embeddings-inference:cpu-1.9.4`) |
| Model | `deepvk/USER-bge-m3` |
| Dimensions | `1024` default, `512` on request |
| Cache | L1 exact LRU in process + L2 Redis binary cache |
| Cross-request coalescing | Opt-in, `FD_COALESCE_ENABLED` (default `false`) |
| Async bulk queue | Opt-in, `FD_QUEUE_ENABLED` (default `false`) |
| Auth | Bearer `FD_API_KEY`, fail-closed |
| Observability | Prometheus at `/metrics`, trace ring at `/v1/traces` |

[CONTRIBUTING.md](./CONTRIBUTING.md) covers setup and the gates a change must pass; [SECURITY.md](./SECURITY.md) explains private vulnerability reporting.

## Contents

- [Quick Start](#quick-start)
- [What It Provides](#what-it-provides)
- [Architecture](#architecture)
- [API](#api)
- [Configuration](#configuration)
- [Operations](#operations)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [Performance](#performance)
- [Repository Layout](#repository-layout)
- [Documentation Map](#documentation-map)
- [Limitations and Non-Goals](#limitations-and-non-goals)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

## Quick Start

Prerequisites: Docker with Compose v2. On a new host, populate the TEI model volume first (skip this if the volume already has the snapshot):

```bash
docker compose run --rm tei \
  --model-id deepvk/USER-bge-m3 \
  --revision 0cc6cfe48e260fb0474c753087a69369e88709ae \
  --max-batch-tokens 8192
# Wait for the model download to finish, then press Ctrl-C.
```

Then start the stack:

```bash
cp .env.example .env
# Edit .env and set FD_API_KEY to a local random value, for example:
#   openssl rand -hex 32
docker compose up -d
```

`docker-compose.override.yaml` loads `.env` when it exists, and `FD_API_KEY` fails closed: with an empty or missing key, every protected endpoint returns `401`. Compose pins `GIN_MODE=release`.

Wait until the API is healthy. TEI loads the BERT model on CPU at first start and can take several minutes; fd retries warmup in the background and reports `503 degraded` until TEI answers.

```bash
until curl -fsS http://localhost:8000/ready >/dev/null; do sleep 5; done
```

Export the same key and run a smoke request:

```bash
export FD_API_KEY="$(grep -E '^FD_API_KEY=' .env | cut -d= -f2-)"

curl -s http://localhost:8000/v1/embeddings \
  -H "Authorization: Bearer $FD_API_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"model":"deepvk/USER-bge-m3","input":"short non-sensitive smoke text"}'
```

A valid response has `object: "list"`, one `data[0].embedding` array of 1024 floats, and `X-Cache: MISS` on the first call.

The compose files pin TEI to a snapshot inside the `tei_data` volume:

```text
/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae
```

The volume must contain that snapshot; the population step above downloads it. See [documents/tei-startup-mitigation-m045.md](documents/tei-startup-mitigation-m045.md) for the inventory and rollback.

## What It Provides

- OpenAI-compatible embeddings, with `input` as a string or an array of strings.
- `dimensions` 1024 (default) or 512, and `encoding_format` `float` (default) or `base64`.
- Two-tier caching with namespace isolation and cache-control endpoints.
- Optional coalescing of concurrent inputs and an optional bounded async queue.
- Bearer auth (fail-closed), CORS, and per-IP/per-user rate limiting.

## Architecture

```text
client
  |
  v
fd API (Go, Gin)
  |-- headers, CORS, API-key auth, per-IP rate limit
  |-- route validation and lifecycle gate
  |-- L1 cache (exact LRU: entry and byte budgets, TTL with jitter)
  |-- L2 cache (Redis, binary blobs, namespace + TTL)
  |-- optional coalescer (FD_COALESCE_ENABLED)
  |-- optional queue worker (FD_QUEUE_ENABLED)
  v
TEI (CPU) -- deepvk/USER-bge-m3, local snapshot under /data
```

Request flow for `POST /v1/embeddings`:

1. Middleware sets `Server` and `X-Request-Id`, applies CORS, authentication, and the per-IP limit; the route then validates shape and limits.
2. The lifecycle gate rejects warming-up or shutting-down requests and enforces `FD_MAX_IN_FLIGHT` when set.
3. Inputs are looked up in L1 by `(text, dimensions)`, then in L2 with one Redis `MGET` over `sha256(text)` keys for the remaining misses.
4. Misses go to TEI in chunks of at most 32 texts; the optional coalescer can merge concurrent calls. The response is shaped and new vectors are written to L1 and L2. `X-Cache` is `HIT` only when every input came from cache.

L1 is an exact LRU with independent entry-count and payload-byte budgets. The byte budget counts `len(key) + len(value)` per entry, not Go heap. Each entry lives for `FD_CACHE_LOCAL_TTL` with ±10% jitter. L2 stores a 2-byte dimension prefix plus `4 * dimensions` bytes of little-endian float32. The cache namespace is part of every L2 key. Changing the model, tokenizer, or chunking behavior without changing the namespace can serve incompatible vectors from the old population.

## API

### Authentication and public routes

Authentication uses `Authorization: Bearer <FD_API_KEY>` and compares SHA-256 digests in constant time. When `FD_API_KEY` is empty, protected routes reject every request with `401 unauthorized`; they never fall open. `OPTIONS` bypasses auth so CORS preflight completes.

Public routes: `GET /live`, `/ready`, `/health`, `/v1/healthcheck`, `/docs` and `/docs/*`, `/openapi.json`. Everything else is protected, including `/version`, `/info`, `/warmup`, `/metrics`, `/v1/traces`, `/v1/cache/*`, `/v1/embeddings`, both batch endpoints, and `/v1/queue`.

### Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/live`, `/ready`, `/health`, `/v1/healthcheck` | Liveness, readiness, deep health, and its alias. |
| `GET` | `/info`, `/version` | Model metadata; build hash, date, and uptime. |
| `GET` | `/metrics` | Prometheus text. Protected. |
| `GET` | `/v1/traces` | Ring buffer of the last 100 requests. |
| `GET` / `POST` | `/warmup` | Warmup status, or start warmup in the background. |
| `POST` | `/v1/embeddings` | OpenAI-compatible embeddings. |
| `POST` | `/v1/batch` | Up to 100 groups of up to 32 inputs each. |
| `POST` | `/embeddings/batch` | Legacy batch shape for an in-product client; base64 by default; max 128 inputs. |
| `POST` | `/v1/cache/flush` | Drop L1 and the L2 namespace. |
| `POST` | `/v1/cache/delete` | Delete one input or an array; `dimensions` defaults to 1024. |
| `POST` | `/v1/queue` | Async submission; `202`, or `503` with `Retry-After: 5` when full. |
| `GET` | `/v1/queue/:id` | Poll: `202` pending, `200` done, `500` failed, `404` unknown or expired. |
| `GET` | `/openapi.json`, `/docs` | OpenAPI 3.1 document and Swagger UI. |

Unknown paths return `404 not_found`. A known path with the wrong method returns `405 method_not_allowed` plus an `Allow` header.

### `POST /v1/embeddings`

```json
{
  "model": "deepvk/USER-bge-m3",
  "input": ["first text", "second text"],
  "dimensions": 1024,
  "encoding_format": "float",
  "user": "consumer-1",
  "priority": "normal"
}
```

Only `input` is required. It may be a string or an array of up to 128 strings. `dimensions` must be `512` or `1024`. `encoding_format` must be `float` or `base64`; base64 is the raw little-endian float32 bytes. `user` feeds per-user rate limiting. `priority` accepts `low`, `normal`, or `high` and is validated but not used for scheduling.

`model` is compatibility metadata: requests are not routed by it. The response `model` and `/health.runtime.model` are authoritative.

```json
{
  "object": "list",
  "model": "deepvk/USER-bge-m3",
  "data": [{"object": "embedding", "index": 0, "dimensions": 1024, "embedding": [0.01, -0.02]}],
  "usage": {"prompt_tokens": 8, "total_tokens": 8}
}
```

`usage.prompt_tokens` is an estimate of input characters divided by 4, not a tokenizer count. With `dimensions: 512`, fd returns the first 512 components of the 1024-dimensional vector; it is a truncation, not a separately trained 512-dimensional space.

### Limits

| Limit | Value | Notes |
|-------|-------|-------|
| Request body | 10 MiB (10,485,760 bytes) | A larger declared `Content-Length` is rejected before reading. |
| Inputs per `/v1/embeddings` | 128 | Each 32-text chunk is one TEI call. |
| Total input characters | 1,200,000 | Sum across all inputs. |
| Per-input tokens (`/v1/embeddings`) | 8192 | Model context limit; TEI truncates longer inputs (auto-truncate is on). |
| `dimensions` | 512 or 1024 | Default 1024. |
| `encoding_format` | `float` or `base64` | Default `float`. |
| Batch endpoints | 100 groups × 32 inputs, or 128 inputs | 2048 characters per input. |

### Response headers

| Header | When | Value |
|--------|------|-------|
| `Server`, `X-Content-Type-Options` | All responses | `fd/<version>` (`fd/dev` locally); `nosniff`. |
| `X-Request-Id` | All responses | Echoed request value, or a generated UUID. |
| `X-Model-Id`, `X-Dimensions`, `X-Cache` | `/v1/embeddings` | Model ID, `1024`/`512`, `HIT`/`MISS`. |
| `ETag`, `Cache-Control` | `/v1/embeddings`, `/info` | `private, max-age=86400`; `If-None-Match` returns `304`. |
| `X-RateLimit-Limit/Remaining/Reset` | Rate-limited routes | `Reset` is seconds until the bucket is full. |
| `Retry-After` | `429`, warmup/overload `503`, queue full, shutdown | Seconds to wait. |
| `Allow` | `405` | Methods registered for the path. |

### Errors

Errors use one envelope; `param` is omitted when no field applies:

```json
{
  "error": {"code": "input_required", "type": "invalid_request_error",
            "param": "input", "message": "input is required (non-empty array of strings)"}
}
```

| Status | Codes |
|--------|-------|
| `400` | `invalid_json`, `input_required`, `dimensions_invalid`, `encoding_format_invalid`, `priority_invalid` |
| `401` | `unauthorized` |
| `404`, `405` | `not_found`, `method_not_allowed` |
| `413` | `batch_too_large`, `payload_too_large`, `input_too_long` |
| `429` | `rate_limit_exceeded` |
| `500` | `internal_error` |
| `503` | `model_not_loaded`, `model_overloaded`, `shutting_down` |

Rate-limited responses carry `Retry-After` with the seconds until the next token refills; `X-RateLimit-Reset` reports the seconds until the bucket is full.

### Batch, cache, and queue behaviour

`POST /v1/batch` takes `{"batches":[["a","b"],["c"]]}` and returns `{"batches":[[[...]],[[...]]]}` in group order. `POST /embeddings/batch` takes `{"inputs":["a"],"dimensions":1024,"encoding_format":"base64"}` and returns `{"embeddings":["..."],"count":1,"dimensions":1024}`; `encoding_format: "float"` returns arrays instead.

`POST /v1/cache/flush` ignores its body and returns the number of L2 keys removed. `POST /v1/cache/delete` accepts `{"input":"text"}` or `{"input":["a","b"]}` and reports how many inputs it processed, not how many keys existed.

The queue is enabled only with `FD_QUEUE_ENABLED=true`. Submit accepts the same body as `/v1/embeddings`, returns `202` with `X-Request-Id`, and stores results in memory for about five minutes. Results are lost on restart; items still queued at shutdown are marked failed so pollers do not wait forever.

## Configuration

Variables are read at startup. `.env.example` lists the local compose values; the real `.env` is gitignored.

### Runtime and HTTP server

| Variable | Default | Purpose |
|----------|---------|---------|
| `EMBEDDING_BACKEND` | `tei` | Only `tei` is supported; any other value fails startup. |
| `MODEL_ID` | `deepvk/USER-bge-m3` | Model identity in responses and `/health`. |
| `TEI_URL` | `http://tei:80` | TEI base URL. |
| `PORT` | `8000` | API listen port. |
| `BIND_HOST` | `0.0.0.0` | API bind address. |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error`; JSON logs to stdout. |
| `GIN_MODE` | unset | Compose pins `release`; `debug` enables Gin debug output. |
| `FD_HTTP_READ_TIMEOUT` | `35s` | Read deadline: 30s handler budget + 5s body margin. |
| `FD_HTTP_WRITE_TIMEOUT` | `75s` | Write deadline: 35s read + 30s handler + 10s response margin. |
| `FD_HTTP_IDLE_TIMEOUT` | `75s` | Keep-alive idle timeout. |

Zero, negative, or unparseable values for any of the three deadlines fall back to the default. `ReadHeaderTimeout` is fixed at 10s. fd does not set `GOMAXPROCS`: the runtime sizes itself from the container CPU limit, and the resolved value is logged as `gomaxprocs`. The TEI HTTP client has a 30s timeout.

### Cache

| Variable | Default | Purpose |
|----------|---------|---------|
| `REDIS_HOST` | `redis:6379` | L2 Redis address. |
| `REDIS_POOL_SIZE` | `50` | Redis connection pool size. |
| `REDIS_CACHE_TTL` | `24h` | L2 TTL; Go duration, must be positive. |
| `REDIS_CACHE_NO_EXPIRE` | `false` | Disable L2 expiry. Mutually exclusive with a set `REDIS_CACHE_TTL`. |
| `REDIS_CACHE_KEY_FASTPATH` | `true` | String literals for 1024/512 dims; key format unchanged. |
| `EMBEDDING_CACHE_VERSION` | `v2` | Primary namespace component. |
| `EMBEDDING_MODEL_ID`, `EMBEDDING_MODEL_REVISION`, `EMBEDDING_TOKENIZER_VERSION`, `EMBEDDING_CHUNKING_VERSION` | empty | Optional namespace components, hashed into the key. |
| `FD_CACHE_MAX_SIZE` | `200000` | L1 entry budget; `0` disables it. |
| `FD_CACHE_MAX_BYTES` | `536870912` | L1 payload budget (512 MiB); `0` disables it. |
| `FD_CACHE_LOCAL_TTL` | `30s` | L1 entry TTL with ±10% jitter. |
| `REDIS_MAXMEMORY`, `REDIS_MAXMEMORY_POLICY`, `REDIS_RDB_SAVE`, `REDIS_AOF_ENABLED` | `2gb`, `allkeys-lru`, `300 1`, `no` | Compose-only Redis settings. |

### Auth, CORS, and capacity

| Variable | Default | Purpose |
|----------|---------|---------|
| `FD_API_KEY` | unset | Bearer key for protected routes. Unset means fail-closed `401`. |
| `FD_CORS_ORIGINS` | unset, effective `*` | Comma-separated allowlist. Empty or `*` allows any origin. |
| `FD_MAX_IN_FLIGHT` | `0` | Concurrent embedding requests; `0` unlimited. |
| `FD_RATE_LIMIT_ENABLED` | `false` | Enables the token-bucket limiters. The value must be `true`. |
| `FD_RATE_LIMIT_IP_RPM` | `100` | Per-IP requests per minute. |
| `FD_RATE_LIMIT_USER_RPM` | `1000` | Per-`user` requests per minute. |

CORS allows `GET, POST, OPTIONS` and the `Content-Type`, `Authorization`, and `X-Request-Id` request headers; preflight returns `204`. The per-IP limiter applies to every route when enabled. The per-user limiter applies to the embedding and batch routes, and only when the request carries `user`.

### Coalescing, queue, warmup, observability

| Variable | Default | Purpose |
|----------|---------|---------|
| `FD_COALESCE_ENABLED` | `false` | Merge concurrent calls into fewer TEI requests. |
| `FD_COALESCE_WINDOW_MS` | `5ms` | Collection window before a batch is flushed. |
| `FD_QUEUE_ENABLED` | `false` | Registers `/v1/queue` and starts the worker. |
| `FD_QUEUE_MAX_SIZE` | `1024` | Bounded in-memory channel capacity. |
| `FD_QUEUE_BATCH_MAX_SIZE` | `32` | Maximum items per worker batch. |
| `FD_QUEUE_BATCH_WINDOW_MS` | `10ms` | Time the worker waits to fill a batch. |
| `FD_WARMUP_START_MAX_ATTEMPTS` | `5` | Startup attempts before the recovery loop takes over. |
| `FD_WARMUP_START_BACKOFF_SEC` | `5` | Fixed wait between startup attempts. |
| `FD_WARMUP_RECOVERY_INTERVAL_SEC` | `30` | Background retry interval. |
| `FD_WARMUP_RECOVERY_ENABLED` | `true` | Set `false` for manual-restart-only behavior. |
| `FD_TRACES_ENABLED` | `true` | Trace ring for `/v1/traces`; set `false` to disable. |

The coalescer flushes when it has collected 32 texts or the window expires; callers block until their batch returns. Shutdown cancels the in-flight batch and fails pending jobs instead of dropping them.

## Operations

### TEI startup

TEI starts from the local snapshot path shown in [Quick Start](#quick-start), not from a Hub repo ID. This avoids the Hub resolution path but makes the `tei_data` volume a prerequisite. Keep `docker-compose.yaml` and `docker-compose.override.yaml` aligned: a compose `command` replacement is total, not merged. A local ORT/ONNX warning at TEI startup is expected because the ONNX files are intentionally absent; TEI falls back to its Candle backend. A long remote `Downloading onnx/model.onnx` phase means the local path is missing or wrong.

### Health, readiness, and warmup recovery

`/health` returns lifecycle state, the runtime block, in-flight counts, and cheap TEI and Redis reachability probes with a 750ms budget. It does not run inference. `/health` is `503` with `status: degraded` while warming up and `status: down` during shutdown; `last_error` keeps the latest failure visible. `/ready` is `200` only after warmup.

The `runtime` block reports `backend`, `model`, `dimensions`, `production_default`, and `cache_namespace`:

```bash
curl -s http://localhost:8000/health
```

Startup warmup makes up to `FD_WARMUP_START_MAX_ATTEMPTS` attempts with `FD_WARMUP_START_BACKOFF_SEC` between them. A background loop also retries every `FD_WARMUP_RECOVERY_INTERVAL_SEC` until TEI answers, the process shuts down, or the loop is disabled. No operator action is needed.

```bash
curl -s http://localhost:8000/ready
docker compose logs --follow api | grep -E 'warmup recovery|model warmup'
```

### Logs, rebuild, and Redis

```bash
docker compose logs --tail=100 api
docker compose logs --tail=100 tei

docker compose build --no-cache api
docker compose up -d api
```

Logs are structured JSON on stdout; request bodies and secrets are never logged. `docker compose up -d --force-recreate api` re-reads `.env`: if the file is missing or `FD_API_KEY` is empty, protected routes start returning `401`, so verify the key after a recreate. Never commit the key. The override publishes Redis only on `127.0.0.1:6379`; L2 survives API restarts in the `redis_data` volume.

### Metrics

`GET /metrics` returns Prometheus text and is protected by the API key. The series are:

- HTTP: `fd_requests_total{status="success|error|timeout"}`, `fd_request_duration_seconds`, `fd_batch_size`, `fd_errors_total{code}`.
- Lifecycle: `fd_model_loaded`, `fd_in_flight_requests`, `fd_in_flight_capacity`.
- Cache: `fd_cache_hits_total{result="hit|miss",tier="l1|l2|miss|all"}`, `fd_cache_entries{tier}`, `fd_cache_memory_bytes{tier}`, `fd_cache_lookup_duration_seconds`, `fd_cache_evictions_total` (capacity evictions only, not TTL expiry).
- TEI: `fd_tei_request_duration_seconds`, `fd_tei_requests_in_flight`, `fd_tei_errors_total{reason="timeout|http_error|circuit_open|model_mismatch|transport"}`, `fd_tei_batch_fill_ratio`.
- Queue: `fd_queue_depth`, `fd_queue_drain_total`, `fd_queue_submit_total{result}`, `fd_queue_batch_size`, `fd_queue_process_duration_seconds`. These collectors are registered, but this tree does not increment them yet.

`fd_cache_entries` and `fd_cache_memory_bytes` are sampled at scrape time. L2 uses `SCAN` over the namespace, so scrape cost grows with key count.

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `401`, `api key is not configured` or a bearer error | Key empty, missing, or different | Set `FD_API_KEY` in `.env`, recreate api, send `Authorization: Bearer <key>`. |
| `413 payload_too_large` | Body above 10 MiB, or total input characters above 1,200,000 | Split the request. |
| `413 batch_too_large` | More than 128 inputs, 100 groups, or 32 per group | Split into smaller batches. |
| `413 input_too_long` | A batch input above 2048 characters | Shorten or split the input. |
| `429 rate_limit_exceeded` | IP or user bucket empty | Read `Retry-After`, then retry after that many seconds. |
| `503 model_overloaded` | `FD_MAX_IN_FLIGHT` reached, or queue full (`Retry-After: 5`) | Lower concurrency or raise the limit. |
| `503 model_not_loaded` | Warmup not complete | Wait for `/ready`; inspect warmup logs. |
| `/health` stays `503 degraded` | TEI still loading or unreachable | Check `docker compose logs tei`; recovery retries automatically. |
| Embeddings look wrong after a model change | Old vectors in the L2 namespace | Change `EMBEDDING_CACHE_VERSION` or flush the namespace. |
| `/metrics` returns `401` | It is a protected route | Add the bearer key. |

## Development

The module lives in `api/` and declares `go 1.26.10`; CI reads the toolchain from `api/go.mod`.

```bash
cd api
go test ./... -short -count=1     # the CI suite
go test ./... -race -count=1      # full suite
go vet ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 \
  run --config ../.golangci.yml ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Use `-race` for changes that touch `cache`, `embed` (the coalescer), or `queue`. All packages pass with these commands on the pinned toolchain; the lint configuration enables 18 linters in fail mode. TEI on CPU can take minutes for the first load, so do not restart the API repeatedly during that window.

```bash
cd tests/integration
FD_BASE_URL=http://localhost:8000 \
FD_INTEGRATION_API_KEY="<same key as FD_API_KEY>" \
go test -v .
```

Without `FD_INTEGRATION_API_KEY`, the suite checks public diagnostics and auth fail-closed behavior and skips protected happy paths. `docker compose config` validates both compose files after an edit. Host-side checks live in `tools/`: `verify_fd_v2_contract.py`, `verify_fd_v2_perf.sh`, and `verify_warmup_recovery.sh`; `benchmark.py` reproduces the latency/cache snapshot.

### CI and review

- `Go Quality` runs on every pull request to `master` and on pushes to `master`. Its job `Go tests and lint` is a required check: `go test ./... -short -count=1`, `go vet`, golangci-lint v2.12.2, and `govulncheck`. There are no path filters, so a documentation-only PR still reports the check.
- `Release notes` runs on a `v*` tag and generates the GitHub release notes with `cliff.toml`; `CHANGELOG.md` stays authoritative.
- CodeRabbit is configured in `.coderabbit.yaml`. The repository is public with fewer than 10 stars, so automatic review is not triggered; comment `@coderabbitai review` on a pull request to request one.

## Performance

These are local snapshots from committed runs, not production SLOs. They move with CPU, TEI version, batch shape, and cache temperature.

| Measurement | Value | Source |
|-------------|------:|--------|
| Cache-hot single input, p95 | 2.008 ms | `benchmark-results/fd-v2-validation-m041.md` |
| Cache-hot 10 inputs, p95 | 5.783 ms | `benchmark-results/fd-v2-validation-m041.md` |
| TEI warm latency, mean | 2.25 ms | `documents/te-perf-root-cause-m042.md` |
| TEI best cold latency | 59.0 ms | `documents/te-perf-root-cause-m042.md` |
| TEI throughput, controlled run | about 750 req/s | `documents/te-perf-root-cause-m042.md` |
| Direct TEI cold p50, batch=1 / 8 / 32 | 199 ms / 1.35 s / 5.34 s | `documents/te-perf-root-cause-m042.md` |
| L2 blob per 1024-dim vector | 4098 bytes | `api/cache/redis.go` |

Cold latency grows steeply with batch size and concurrent load because TEI schedules and queues internally. Cache-hot latency is the number that matters for repeated same-host traffic; do not compare a warm-cache p95 with a cold TEI p95. `benchmark.py` reproduces a snapshot from the host against a local stack.

## Repository Layout

```text
api/                     Go service: cache/, embed/, handlers/, lifecycle/,
                         middleware/, observability/, openapi/, queue/
tests/integration/       black-box Go module for a running stack
tools/                   contract, perf, and model verification scripts
docs/                    consumer contract and design notes
documents/               investigation and closeout reports
benchmark-results/       committed local benchmark snapshots
.github/                 workflows, issue forms, PR template, CODEOWNERS, dependabot
benchmark.py             host-side latency and cache benchmark
docker-compose.yaml      tei + api + redis
docker-compose.override.yaml  local env, Redis port binding, TEI health grace
.env.example             local environment template
.golangci.yml            lint configuration
.coderabbit.yaml         CodeRabbit review configuration
AGENTS.md                rules for AI coding agents
CITATION.cff             citation metadata
cliff.toml               release-note configuration
CONTRIBUTING.md          setup, gates, PR expectations
SECURITY.md              supported versions, private reporting
CODE_OF_CONDUCT.md       Contributor Covenant 3.0
CHANGELOG.md             release history
```

Only `api/` is compiled into the service. Everything else is tooling, documentation, or evidence.

## Documentation Map

| Document | Purpose |
|----------|---------|
| [`docs/same-host-embedding-service-contract.md`](docs/same-host-embedding-service-contract.md) | Canonical HTTP and runtime contract for same-host consumers. |
| [`docs/fd-v2.md`](docs/fd-v2.md) | Requirements and error catalog behind the current API. |
| [`documents/te-perf-root-cause-m042.md`](documents/te-perf-root-cause-m042.md) | TEI queue-time and startup analysis, with direct-TEI timings. |
| [`documents/tei-startup-mitigation-m045.md`](documents/tei-startup-mitigation-m045.md) | Local snapshot startup mode, inventory, and rollback. |
| [`benchmark-results/fd-v2-validation-m041.md`](benchmark-results/fd-v2-validation-m041.md) | Contract validation run and cache-hot latency. |
| [`benchmark-results/m045-tei-local-path-startup-proof.md`](benchmark-results/m045-tei-local-path-startup-proof.md) | Controlled proof of the local snapshot startup. |
| [`api/openapi/spec.go`](api/openapi/spec.go) | Source of the served OpenAPI document. |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Setup, required gates, PR expectations. |
| [`SECURITY.md`](SECURITY.md), [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) | Private vulnerability reporting; Contributor Covenant 3.0. |
| [`CHANGELOG.md`](CHANGELOG.md) | Release history and the `[Unreleased]` section. |
| [`AGENTS.md`](AGENTS.md) | Rules for AI coding agents in this repository. |

## Limitations and Non-Goals

- The service targets same-host consumers. It has no TLS, no multi-tenant key management, and no public-internet hardening.
- ONNX is not an active runtime; historical artifacts and benchmarks are not a deployment path.
- `/health` is not an inference probe. Use `/ready` plus a small embedding request.
- The TEI snapshot path is pinned; changing the model revision needs a compose change and a fresh startup proof.
- The async queue is in-memory: results do not survive a restart and are retained for about five minutes.
- `dimensions: 512` truncates the vector; it does not renormalize or project it.
- Benchmark numbers are environment-specific snapshots. A single Redis namespace is shared by every caller of one deployment; isolation is an operator responsibility.

## Contributing

Read [CONTRIBUTING.md](./CONTRIBUTING.md) before opening a pull request. One logical change per PR; `go test ./... -short`, `go vet`, golangci-lint v2.12.2, and `govulncheck` must be clean; user-visible changes need a `CHANGELOG.md` entry under `[Unreleased]`; new configuration must be documented here and in `.env.example`. The [issue chooser](https://github.com/rager306/fd/issues/new/choose) has forms for bugs and feature requests.

## Security

Report vulnerabilities privately through [GitHub security advisories](https://github.com/rager306/fd/security/advisories/new) or by email as described in [SECURITY.md](./SECURITY.md). Do not open a public issue, and do not include real API keys or customer text in a report.

## License

[MIT License](./LICENSE). Copyright (c) 2026 rager306.

See [CHANGELOG.md](./CHANGELOG.md) and [GitHub Releases](https://github.com/rager306/fd/releases) for release history.
