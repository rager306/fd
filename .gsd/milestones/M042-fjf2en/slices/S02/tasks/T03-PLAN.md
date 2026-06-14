---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: X-Concurrent-Chunks header + metrics

api/middleware/headers.go (M041 S03 deliverable) расширить: в async mode добавить response header X-Concurrent-Chunks: N (number of chunks sent in parallel for this request). api/observability/metrics.go (M041 S03): добавить fd_async_chunks_total counter (incremented per chunk in flight), fd_async_chunk_duration_seconds histogram (per chunk inference time). Sync mode — headers/metrics absent (no overhead).

## Inputs

- None specified.

## Expected Output

- `api/middleware/headers.go (extended)`
- `api/observability/metrics.go (extended)`

## Verification

curl -I -X POST /v1/embeddings с FD_ASYNC_CHUNKS=true показывает X-Concurrent-Chunks: 4 для batch=128. /metrics показывает fd_async_chunks_total counter incrementing.
