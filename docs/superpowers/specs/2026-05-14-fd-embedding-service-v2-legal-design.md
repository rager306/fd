# Embedding Service v2 — Legal Russian Knowledge Graph

**Date:** 2026-05-14
**Status:** Draft
**Author:** Hermes / Алекс

---

## 1. Context and Problem Statement

Embedding-сервис для построения knowledge graph в юридической предметной области на русском языке. embeddings используются для:
- Нод графа (entity embeddings, 1024d)
- Рёбер графа (relationship embeddings, 512d)
- Векторный поиск (внешняя система, за пределами scope)

Внешняя система: FalkorDB (graph) + Redis (vector search). Наша задача — генерация embeddings.

**Constraints:**
- Качество и скорость — оба критичны, без компромиссов
- Юридический домен — ошибка в similarity = ошибка в графе связей (недопустимо)
- Масштаб: средний/большой (10K–500K+ нод)
- Режим: real-time + batch indexing
- Память: 46GB RAM, swap при необходимости

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Client (external)                    │
└─────────────────────┬───────────────────────────────────┘
                      │ POST /v1/embeddings
                      │ GET /health
                      ▼
┌─────────────────────────────────────────────────────────┐
│                    fd_api (Go)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌───────────────┐  │
│  │  Handlers   │  │   embed     │  │    cache      │  │
│  │ /v1/embed   │  │   types.go  │  │   redis.go    │  │
│  └─────────────┘  └─────────────┘  └───────────────┘  │
│         │                │                  │           │
│         ▼                ▼                  ▼           │
│  ┌─────────────────────────────────────────────────┐   │
│  │  TEI Client (internal)                           │   │
│  │  - Single embedding call                         │   │
│  │  - Batch embedding call                         │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────┘
                      │ :80
        ┌─────────────┴─────────────┐
        ▼                             ▼
┌───────────────┐          ┌─────────────────┐
│   fd_tei      │          │    fd_redis     │
│ TEI cpu-1.9   │          │  Redis Stack    │
│ deepvk/USER   │          │  LRU 2GB        │
│ bge-m3        │          │  hash cache     │
│ ONNX FP16     │          │                 │
└───────────────┘          └─────────────────┘
```

---

## 3. Model Configuration

### 3.1 Model
- **ID:** `deepvk/USER-bge-m3`
- **Params:** 359M
- **Dimensions:** 1024 (full), 512 (Matryoshka slice)
- **Max tokens:** 8192
- **Backend:** TEI cpu-1.9, ONNX FP16 (not Candle)

### 3.2 ONNX Conversion
```bash
# Convert to ONNX FP16 for TEI
python3 -c "
from optimum.exporters.onnx import main
import subprocess
subprocess.run(['python3', '-m', 'optimum.exporters.onnx',
  '--model', 'deepvk/USER-bge-m3',
  '--task', 'feature-extraction',
  '--dtype', 'fp16',
  '--optimize', 'O3'], check=True)
"
```

### 3.3 Dimension Strategy
| Entity type | Dimension | Storage | Use case |
|-------------|-----------|---------|----------|
| Node (entity) | 1024d | Full vector | Exact similarity search |
| Edge (relationship) | 512d | Matryoshka slice | Graph edges, comparisons |
| Cache key | 128bit | SHA256 hash | O(1) Redis lookup |

**Implementation:**
- Both 1024d and 512d generated from single model forward pass
- `vector[:512]` = edge embedding (no recomputation)
- API returns `dim` parameter: 1024 (default) or 512

---

## 4. API Design

### 4.1 Endpoints

```
GET  /health              → {"status": "ok", "time": "..."}
GET  /ready               → {"ready": true}
POST /v1/embeddings       → OpenAI-compatible
POST /embeddings/batch    → Internal batch endpoint
POST /embed               → Internal single embedding
GET  /embed/:id           → By cache key
```

### 4.2 Request/Response: POST /v1/embeddings

```json
// Request
{
  "model": "deepvk/USER-bge-m3",
  "input": "Текст для эмбеддинга",
  "dimensions": 1024,          // optional, default 1024
  "encoding_format": "base64"   // "base64" or "float"
}

// Response
{
  "object": "list",
  "model": "deepvk/USER-bge-m3",
  "usage": {"prompt_tokens": 15, "total_tokens": 15},
  "data": [
    {
      "object": "embedding",
      "index": 0,
      "embedding": "BASE64_STRING"
    }
  ]
}
```

### 4.3 Batch Endpoint: POST /embeddings/batch

```json
// Request
{
  "inputs": ["Текст 1", "Текст 2", "..."],
  "dimensions": 512,           // for edges
  "encoding_format": "base64"
}

// Response
{
  "embeddings": ["BASE64_1", "BASE64_2", "..."],
  "count": 3,
  "dimensions": 512
}
```

---

## 5. Caching Strategy

### 5.1 Cache Key
```go
key = SHA256(text)[:16]  // 128-bit, collision-resistant enough for cache
```

### 5.2 Redis Structure
```
Key:   embed:cache:{sha256_prefix}
Value: JSON { "embedding": "base64", "dim": 1024, "created": "timestamp" }
TTL:   24h
```

### 5.3 Cache Flow
1. Compute SHA256(text)
2. Redis GET — O(1) lookup, ~1-2ms
3. HIT → return cached
4. MISS → call TEI → store → return

### 5.4 Cache Invalidation
- TTL-based (24h)
- Manual purge by pattern: `embed:cache:*`

---

## 6. Docker Configuration

### 6.1 Services
```yaml
services:
  tei:
    image: ghcr.io/huggingface/text-embeddings-inference:cpu-1.9
    command: --model-id deepvk/USER-bge-m3 --dtype fp16 --max-batch-tokens 8192
    # No memory limit — 46GB available
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:80/health"]
      interval: 5s
      timeout: 5s
      retries: 20
      start_period: 180s

  api:
    build: ./api
    ports: ["8000:8000"]
    environment:
      - TEI_URL=http://tei:80
      - REDIS_ADDR=redis:6379
      - PORT=8000
    depends_on:
      tei:   { condition: service_healthy }
      redis: { condition: service_healthy }
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]

  redis:
    image: redis/redis-stack:latest
    command: redis-server --maxmemory 2gb --maxmemory-policy allkeys-lru --protected-mode no
    deploy:
      resources:
        limits:
          cpus: "2"
          memory: 2G
```

### 6.2 Memory Allocation
- TEI: model ~2GB FP16, activations scale with batch
- Redis: 2GB LRU (cache only)
- API: <100MB (stateless Go binary)

---

## 7. Expected Performance

### 7.1 Latency (measured after ONNX FP16 switch)

| Path | Latency | Notes |
|------|---------|-------|
| Cold (first time) | ~12-15ms | FP16 vs 19ms FP32 |
| Cache hit | ~1.5-2ms | Redis O(1) |
| Batch 32 texts | ~200-400ms | TEI dynamic batching |

### 7.2 Throughput
| Concurrency | Requests/sec |
|-------------|--------------|
| 1 | ~70-80 |
| 4 | ~250-300 |
| 8 | ~400-500 |

### 7.3 Memory
| Component | RAM |
|----------|-----|
| TEI (model) | ~2GB FP16 |
| TEI (activations, batch 64) | ~4-6GB |
| Redis cache (2GB limit) | ~2GB |
| API | <100MB |
| **Total** | ~8-10GB (with batch 64) |

---

## 8. Quality Assurance

### 8.1 Model Quality (RU Legal)
- deepvk/USER-bge-m3: SOTA для русского (ruMTEB benchmark)
- FP16: ~0.1-0.5% loss vs FP32 (negligible)
- 512d (Matryoshka): loss устойчив к truncation для relationship comparison

### 8.2 Verification Tests
- Unit tests: cache, handlers, types
- Integration tests: Redis, TEI connectivity
- E2E: curl /v1/embeddings with Russian legal text
- Quality: manual cosine similarity check on known legal entity pairs

---

## 9. Out of Scope (for this spec)

- FalkorDB integration
- Redis vector search / HNSW index
- Model fine-tuning
- Authentication / API keys
- Metrics / monitoring (Prometheus)
- CI/CD

---

## 10. Open Questions (deferred)

1. Edge embeddings: concat vs separate text → resolved externally
2. Batch size tuning → empiric after deployment
3. Cache warmup strategy for known FAQs
4. FalkorDB vector index configuration

---

## 11. Files to Change

```
/root/fd/
  docker-compose.yaml          # Add ONNX FP16, batch config
  api/
    internal/embed/types.go    # Add dimensions param
    internal/embed/tei.go      # Add batch endpoint
    internal/handlers/embeddings.go  # Handle dim param
    internal/cache/redis.go    # Store dim in cache value
    cmd/api/main.go            # Pass dim config
  scripts/
    convert_onnx_fp16.sh       # ONNX conversion script (NEW)
```
