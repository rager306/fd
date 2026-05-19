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
  -d '{"input":["юридическая справка"]}'
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

## Operational Notes

### Redis exposure and host tuning

- The base Compose file does not publish Redis to the host. The local override binds Redis to `127.0.0.1:6379` so `benchmark.py` can run from the host without exposing Redis on all interfaces.
- Do not change the override to `6379:6379` unless you also add proper network access controls. Runtime validation previously observed unsolicited external Redis traffic when Redis was exposed on `0.0.0.0:6379`.
- Redis may log `Memory overcommit must be enabled`. That is a host-level deployment note, not an application bug. For hosts that rely on Redis persistence/background saves, set `vm.overcommit_memory=1` through the host's normal sysctl management.

### TEI backend artifacts

- The Compose command includes `--dtype fp16`, but the current `deepvk/USER-bge-m3` runtime has been observed falling back when ONNX artifacts are unavailable.
- Current local benchmarks are valid for the measured Candle/CPU fallback runtime. Treat ONNX export as a future measured optimization, not a correctness requirement.
- If ONNX artifacts are introduced later, run an A/B benchmark with the same `uv --python 3.13` command and compare cold p50/p95, memory, and response shape before making it the default.

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
