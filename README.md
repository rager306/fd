# FD Embedding Service

High-performance text embedding API: two-tier cache (L1 local + L2 Redis) + TEI inference.

## Performance

| Metric | Value |
|--------|-------|
| Warm latency (1024d) | 2.6 ms |
| Cold latency (TEI) | 19 ms |
| Cache speedup | 28.8x |
| Max throughput | ~645 req/s |
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
- **Redis Stack** (port 6379) — binary cache + future HNSW vector search

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

## Development

```bash
cd api
go build ./...
go test ./... -short
```

```bash
# E2E benchmark
python3 benchmark.py
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
