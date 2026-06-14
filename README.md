# FD Embedding Service

High-performance text embedding API: two-tier cache (L1 local + L2 Redis) + TEI inference.

## Performance

Validated local benchmark snapshots are committed under `benchmark-results/`.
Run benchmarks from the host with `uv` and Python 3.13 against the local Docker stack.

| Metric | Value |
|--------|-------|
| Warm latency (1024d) | ~2 ms |
| Cold latency (TEI) | 20-200 ms local range |
| Cache speedup | 40-80x local range |
| Max throughput | ~644 req/s in latest final run |
| Redis L2 after API restart | ~3 ms in latest diagnostic run |
| Storage per embedding | 4 KB (binary) |

## Quick Start

```bash
docker compose up -d
curl -X POST http://localhost:8000/v1/embeddings \
  -H 'Content-Type: application/json' \
  -d '{"input":["short non-sensitive smoke text"]}'
```

## API

```
GET  /health              — health check
POST /v1/embeddings       — OpenAI-compatible embeddings endpoint
POST /embeddings/batch    — batch embeddings for FalkorDB
```

### Embeddings Request

```json
{
  "input": ["text to embed"],
  "dimensions": 1024
}
```

`dimensions`: 1024 (nodes) or 512 (edges). Default: 1024.

### Batch Request

```json
{
  "inputs": ["text1", "text2"],
  "dimensions": 1024,
  "encoding_format": "base64"
}
```

`encoding_format`: `base64` (default) or `float`.

> **Full service contract:** For same-host local clients, see [docs/same-host-embedding-service-contract.md](docs/same-host-embedding-service-contract.md). It covers `/health` readiness semantics, timeout/retry guidance, cache namespace isolation, and the TEI-only current runtime contract.

## Architecture

```
Request → L1 (sync.Map, ~50ns) → L2 (Redis binary, ~0.5ms) → TEI (~70ms)
                 ↓ miss                   ↓ miss              ↓ miss
              backfill                backfill            cache + return
```

**Components:**
- **Go API** (port 8000) — HTTP handler, two-tier cache orchestration
- **TEI** (port 30080) — Rust embedding inference server (deepvk/USER-bge-m3)
- **Redis Stack** (`127.0.0.1:6379` in local override) — binary cache + future HNSW vector search

**Stack:** Go 1.25, Gin, go-redis v9, TEI cpu-1.9, Redis Stack, Docker Compose

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TEI_URL` | `http://tei:80` | TEI endpoint |
| `REDIS_HOST` | `redis:6379` | Redis address |
| `MODEL_ID` | `deepvk/USER-bge-m3` | HuggingFace model ID |
| `REDIS_POOL_SIZE` | `50` | Connection pool size |
| `BIND_HOST` | `0.0.0.0` | Bind address |
| `PORT` | `8000` | API port |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error`; `debug` enables cache-path events |
| `REDIS_CACHE_TTL` | `24h` | Redis L2 embedding TTL, parsed as a Go duration; ignored when no-expire mode is enabled |
| `REDIS_CACHE_NO_EXPIRE` | `false` | Set `true` for long-lived reusable research cache entries with no key TTL |
| `EMBEDDING_CACHE_VERSION` | `v2` | Cache schema/version namespace. Default preserves existing `v2` keys |
| `EMBEDDING_MODEL_ID` | unset | Optional model namespace component for Redis keys; set for long-lived multi-model research caches |
| `EMBEDDING_MODEL_REVISION` | unset | Optional model revision namespace component |
| `EMBEDDING_TOKENIZER_VERSION` | unset | Optional tokenizer namespace component |
| `EMBEDDING_CHUNKING_VERSION` | unset | Optional chunking/splitting namespace component |
| `REDIS_MAXMEMORY` | `2gb` | Redis server maxmemory used by Compose |
| `REDIS_MAXMEMORY_POLICY` | `allkeys-lru` | Redis eviction policy used by Compose |
| `REDIS_RDB_SAVE` | `300 1` | Redis RDB snapshot rule for rebuildable cache persistence |
| `REDIS_AOF_ENABLED` | `no` | Redis appendonly mode; keep `no` unless cache write loss is more expensive than AOF overhead |

## Operational Notes

### Redis exposure, persistence, and host tuning

- The base Compose file does not publish Redis to the host. The local override binds Redis to `127.0.0.1:6379` so `benchmark.py` can run from the host without exposing Redis on all interfaces.
- Do not change the override to `6379:6379` unless you also add proper network access controls. Runtime validation previously observed unsolicited external Redis traffic when Redis was exposed on `0.0.0.0:6379`.
- Redis L2 is a rebuildable embedding cache, not the source of truth. Compose defaults to an RDB-first persistence profile (`REDIS_RDB_SAVE="300 1"`, `REDIS_AOF_ENABLED=no`) so cache entries can survive normal Redis/container restarts while avoiding AOF write overhead.
- For repeated research/chunking workflows, set `REDIS_CACHE_TTL` to a long duration such as `168h` or set `REDIS_CACHE_NO_EXPIRE=true`. Do not set both at once; the API rejects that configuration at startup.
- For long-lived or multi-model cache reuse, set `EMBEDDING_MODEL_ID`, `EMBEDDING_MODEL_REVISION`, `EMBEDDING_TOKENIZER_VERSION`, and `EMBEDDING_CHUNKING_VERSION` as appropriate. Redis keys use short hashes of these values; benchmark artifacts record the clear effective config values for comparison.
- Redis memory behavior is controlled by `REDIS_MAXMEMORY` and `REDIS_MAXMEMORY_POLICY` in Compose. `allkeys-lru` is the default; `allkeys-lfu` may be worth benchmarking for repeated chunk reuse.
- Redis may log `Memory overcommit must be enabled`. That is a host-level deployment note, not an application bug. For hosts that rely on Redis persistence/background saves, set `vm.overcommit_memory=1` through the host's normal sysctl management.

### TEI backend artifacts

- The current supported runtime path is TEI. Historical ONNX experiments remain research artifacts only and are not an operator option in the current build.
- Current local benchmarks are valid for the measured TEI/Candle CPU runtime. Treat any future ONNX work as a separate research milestone with fresh A/B benchmarks and operational gates.

### Runtime logging

- Default `LOG_LEVEL=info` keeps startup, connection, warning, and error logs visible without logging every successful embedding request.
- Set `LOG_LEVEL=debug` when diagnosing cache behavior. Debug cache events include short key hashes and dimensions, not raw input text.

## Development

```bash
cd api
go build ./...
go test ./... -short
```

### Quality tooling

Tests use Go's standard `testing` package plus [Testify](https://github.com/stretchr/testify) `assert`/`require` helpers for clearer failures in representative cache and handler tests.

Run the full Go test suite:

```bash
cd api
go test ./... -short
```

Run the configured GolangCI-Lint gate from the Go module directory:

```bash
cd api
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 \
  run --config ../.golangci.yml ./...
```

`.golangci.yml` enables Staticcheck through GolangCI-Lint, along with `go vet`, `errcheck`, `unused`, `ineffassign`, `goconst`, and `misspell`. Use this lint command before committing Go code changes.

The same quality gate is prepared for GitHub Actions in `.github/workflows/go-quality.yml`. It runs on `push`, `pull_request`, and manual dispatch for API, lint config, workflow, or README changes. Remote CI run evidence is only available after these local commits are pushed to GitHub.

### Local benchmark

Prerequisites:

- Docker Compose stack is running and healthy: `docker compose up -d --build`.
- API is reachable from the host at `http://localhost:8000`.
- Redis is reachable from the host at `127.0.0.1:6379` through `docker-compose.override.yaml`.
- Python benchmark execution uses `uv` with Python 3.13; do not use ad-hoc virtualenvs for the recorded benchmark flow.

```bash
mkdir -p benchmark-results
uv run --python 3.13 --with requests --with redis python benchmark.py \
  | tee benchmark-results/fd-benchmark-local.txt
```

Benchmark side effects:

- `benchmark.py` calls `FLUSHALL` on the local Redis instance to measure cold/warm cache behavior.
- Section 5 restarts the `api` service with `docker compose restart api` when Docker Compose is available, waits for `/health`, then measures Redis L2 cache behavior after API restart.
- Do not run this benchmark against a shared or production Redis/API environment.

Useful verification after a benchmark run:

```bash
docker compose ps
docker compose logs --tail=100 api redis tei
```

## Project Structure

```
fd/
├── api/
│   ├── cache/
│   │   ├── local.go      # L1: sync.Map cache with TTL
│   │   ├── redis.go      # L2: Redis binary storage + pool timeouts
│   │   └── tiered.go     # Two-tier cache with singleflight
│   ├── embed/
│   │   ├── tei.go        # TEI client
│   │   └── types.go      # Request/response types
│   ├── handlers/
│   │   ├── embeddings.go  # /v1/embeddings
│   │   ├── batch.go      # /embeddings/batch
│   │   └── health.go     # /health
│   └── main.go           # Entry point, DI wiring
├── docker-compose.yaml
├── docker-compose.override.yaml
└── benchmark.py
```

## Changelog

See [CHANGELOG.md](./CHANGELOG.md).
