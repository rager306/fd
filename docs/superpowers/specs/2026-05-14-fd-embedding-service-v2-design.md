# FD Embedding Service v2 — Design Specification

**Date:** 2026-05-14
**Status:** Draft — pending user approval

---

## Overview

OpenAI API compatible embedding service: TEI (Rust) + Redis Stack + Go API. Optimized for maximum performance on AMD EPYC 12-core VPS.

**Goals:**
1. Fix existing bugs (embedding=null, failing tests)
2. Add full OpenAI API compatibility (`/v1/embeddings`)
3. Optimize build and deployment process
4. Add structured logging

---

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│  Go API     │────▶│    TEI      │
│ (OpenAI SDK)│     │  :8000      │     │  :80        │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                           ▼
                     ┌─────────────┐
                     │   Redis     │
                     │  Stack      │
                     │  :6379      │
                     └─────────────┘
```

- API accessible only inside Docker network (no external internet exposure)
- TEI: internal port 80, exposed externally on **30080** (host) for model downloads only
- Redis: internal port 6379, not exposed externally
- External clients connect via internal Docker network

---

## Endpoints

### Public (OpenAI Compatible)

**POST /v1/embeddings**
- Synchronous — returns embedding immediately
- Full OpenAI API compliance

Request:
```json
{
  "model": "deepvk/USER-bge-m3",
  "input": "текст для эмбеддинга"
}
```

Batch request:
```json
{
  "model": "deepvk/USER-bge-m3",
  "input": ["текст 1", "текст 2", "текст 3"]
}
```

Response:
```json
{
  "object": "list",
  "data": [{
    "object": "embedding",
    "embedding": [0.123, -0.456, ...],
    "index": 0
  }],
  "model": "deepvk/USER-bge-m3",
  "usage": {
    "prompt_tokens": 10,
    "total_tokens": 10
  }
}
```

**GET /health**
```json
{"status": "ok", "time": "2026-05-14T15:00:00Z"}
```

### Internal (for batch/queue optimizations, future)

**POST /embed** — async, returns 202 + UUID
**GET /embed/:id** — get result by UUID

Internal endpoints use Redis worker queue for batch processing.

---

## Configuration

All configuration via `.env` file:

```
TEI_URL=http://tei:80
REDIS_HOST=redis:6379
REDIS_QUEUE=embed:queue
CACHE_PREFIX=embed:cache:
MODEL_ID=deepvk/USER-bge-m3
OTEL_SERVICE_NAME=fd-api
GIN_MODE=release
LOG_LEVEL=info
BIND_HOST=0.0.0.0
```

Docker Compose loads `.env` via `env_file`.

Build-time vs Runtime:
- Build args (ARG) — only for compile-time constants
- Environment variables (ENV) — for runtime configuration

---

## Docker Optimization

### Multi-stage Dockerfile

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /api

# Stage 2: Runtime
FROM scratch
COPY --from=builder /app/api /api
COPY --from=builder /app/.env /.env
ENTRYPOINT ["/api"]
```

### BuildKit Optimizations

- `DOCKER_BUILDKIT=1` — enabled by default in 2026
- Layer caching: `--cache-from` / `--cache-to`
- Cache mount for go mod: `--mount=type=cache,target=/go/pkg/mod`

### Build Commands

```bash
# Fast rebuild (changed files only)
docker build --target builder -t fd-api:builder .
docker build --target runtime -t fd-api:latest .

# With cache mount (persists between builds)
docker build \
  --mount=type=cache,target=/go/pkg/mod \
  -t fd-api:latest .
```

---

## Logging

**Library:** Go `log/slog` (standard library, Go 1.21+) — zero-dependency

**Format:** JSON (machine-parseable)
```json
{"time":"2026-05-14T15:00:00Z","level":"INFO","msg":"embed request","text_len":142,"cached":false,"duration_ms":47}
```

**Levels:** DEBUG, INFO, WARN, ERROR

**Configuration:** `LOG_LEVEL` env variable

---

## Testing

### Unit Tests
- Table-driven approach
- `testify/assert`
- Run: `go test ./...`

### Integration Tests (testcontainers)
- Redis via testcontainers
- Run: `go test -tags=integration ./...`

### End-to-End (real services)
- Real TEI + Redis in Docker network
- Run: `go test -tags=integration_real ./...`

### Load/Performance
- `go test -bench=. -benchmem`
- k6 or hey for HTTP load testing
- Metrics: latency p50/p95/p99, throughput RPS

---

## Project Structure

```
/root/fd/
├── api/
│   ├── main.go
│   ├── handlers/     # HTTP handlers
│   ├── embed/         # embedding logic
│   ├── cache/         # Redis cache layer
│   └── internal/      # internal packages
├── tests/
│   ├── integration/   # testcontainers tests
│   └── e2e/          # real services tests
├── docs/
│   └── specs/         # design specs
├── Dockerfile
├── docker-compose.yaml
├── docker-compose.override.yaml
├── .env
├── .dockerignore
└── README.md
```

---

## Future (HNSW)

- HNSW index in Redis for internal optimizations
- NOT exposed as external API
- Purpose: deduplication, batch efficiency, memory optimization
- Requires further research before implementation

---

## Scope for Phase 1

1. Rewrite API with OpenAI compatible `/v1/embeddings` endpoint (from scratch)
2. Implement synchronous embedding flow (no more async/queue for v1)
3. Multi-stage Dockerfile with BuildKit optimization
4. Structured logging with slog
5. Environment-based configuration (.env)
6. Comprehensive tests (unit + integration + e2e + benchmarks)

---

## Verification Commands

```bash
# Health check
curl http://api:8000/health

# OpenAI-compatible endpoint (from inside Docker network)
curl -X POST http://api:8000/v1/embeddings \
  -H 'Content-Type: application/json' \
  -d '{"model":"deepvk/USER-bge-m3","input":"тест"}'

# Run tests
go test ./...
go test -tags=integration ./...
go test -tags=integration_real ./...
go test -bench=. -benchmem
```
