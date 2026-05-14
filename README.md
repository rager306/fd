# FD API — Embedding Service

Embedding-сервис: TEI (Rust) + Redis Stack + Go API.

## Быстрый старт

```bash
cd /root/fd
docker compose -f docker-compose.yaml -f docker-compose.override.yaml build
docker compose -f docker-compose.yaml -f docker-compose.override.yaml up
```

## Архитектура

```
POST /embed
  → check Redis cache (by text hash)
  → if miss: LPUSH to queue, return 202 + id

GET /embed/:id
  → HGET from Redis

Worker (background):
  → BRPOP from queue
  → POST to TEI
  → HSET to Redis (by text hash AND by id)
```

## Стек

| Component | Tool | Version |
|---|---|---|
| Embedding inference | TEI (Rust) | cpu-1.9 |
| Cache + vector search | Redis Stack | latest |
| API | Go + Gin | 1.22 |
| Observability | OpenTelemetry | SDK |

## Модели (из TEI docs)

| MTEB Rank | Model | Size | Model ID |
|---|---|---|---|
| 2 | Qwen3-Embedding-8B | 7.57B | Qwen/Qwen3-Embedding-8B |
| 4 | Qwen3-Embedding-0.6B | 509M | Qwen/Qwen3-Embedding-0.6B |
| 8 | EmbeddingGemma-300M | 308M | google/embeddinggemma-300m |
| 15 | gte-Qwen2-1.5B-instruct | 1.78B | Alibaba-NLP/gte-Qwen2-1.5B-instruct |
| 35 | snowflake-arctic-embed-l-v2.0 | 568M | Snowflake/snowflake-arctic-embed-l-v2.0 |
| 79 | nomic-embed-text-v1.5 | 137M | nomic-ai/nomic-embed-text-v1.5 |

Рекомендация: **Qwen3-Embedding-0.6B** (MTEB #4, 509M params, 32K context) — хороший баланс качества и размера для CPU.

## Переменные окружения

### TEI

| Variable | Default | Description |
|---|---|---|
| `MAX_BATCH_TOKENS` | `32768` | Max tokens per batch (higher = better throughput) |
| `MAX_BATCH_REQUESTS` | `64` | Max individual requests per batch |
| `MAX_CONCURRENT_REQUESTS` | `256` | Backpressure limit |
| `TOKENIZATION_WORKERS` | `12` | Tokenizer thread count |
| `DTYPE` | `float32` | CPU: float32 faster than float16 |
| `AUTO_TRUNCATE` | `true` | Auto-truncate oversized inputs |
| `OMP_NUM_THREADS` | `12` | OpenMP threads |

### API

| Variable | Default | Description |
|---|---|---|
| `TEI_URL` | `http://localhost:8080` | TEI endpoint |
| `REDIS_HOST` | `localhost:6379` | Redis address |
| `REDIS_QUEUE` | `embed:queue` | BRPOP queue |
| `REDIS_CACHE_PREFIX` | `embed:cache:` | Redis key prefix |
| `OLLAMA_MODEL` | `nomic-embed-text-v1.5` | Model ID (compatibility) |
| `GOMAXPROCS` | all cores | Go thread count |
| `GOGC` | `100` | GC trigger threshold |
| `REDIS_POOL_SIZE` | `50` | Redis connection pool |
| `GIN_MODE` | `release` | Gin release mode |

### Redis

| Variable | Default | Description |
|---|---|---|
| `maxmemory` | `16gb` | Max memory for embeddings |
| `maxmemory-policy` | `allkeys-lru` | Eviction policy |
| `hz` | `100` | Internal cron frequency |
| `appendonly` | `no` | Disable AOF for pure cache |

## Оптимизации под среду

**Среда:** AMD EPYC 12 cores, 48GB RAM, CPU-only, 10 Gbps network, no swap.

### TEI (CPU)
- `MAX_BATCH_TOKENS=32768` — увеличивает throughput при батчинге
- `DTYPE=float32` — на CPU нет fp16 SIMD ускорения, float32 быстрее
- `OMP_NUM_THREADS=12` — все ядра для OpenMP
- `TOKENIZATION_WORKERS=12` — параллельная токенизация

### Go API
- `GOMAXPROCS=12` — все 12 ядер для горутин
- `GOGC=100` — по умолчанию, при 48GB нет нужды в агрессивном GC
- `REDIS_POOL_SIZE=50` — 50 соединений достаточно для Redis
- `GIN_MODE=release` — отключает debug overhead

### Redis
- `maxmemory 16gb` — держим весь кэш в памяти (embedding ~500MB per 10K docs)
- `hz=100` — чаще обрабатываем expire/eviction
- `appendonly=no` — Redis как чистый cache, не нужен персистентный WAL

### Сеть
- 10 Gbps ethernet — пропускная способность не узкое место
- Redis и TEI на одном хосте — локальная сеть минимизирует latency

## API Endpoints

```
GET  /health           — health check
POST /embed            — async embed, returns id
GET  /embed/:id        — get result by id
GET  /embed/search?q=  — semantic search (501 not implemented)
```

## TODO

- [ ] Semantic search — Redis Stack FT.CREATE с HNSW
- [ ] OTLP exporter для метрик (Prometheus/Jaeger)
- [ ] Graceful shutdown (SIGTERM handling)
- [ ] Health check для worker-а
- [ ] go.sum сгенерировать (создаётся автоматически при `go mod tidy`)
