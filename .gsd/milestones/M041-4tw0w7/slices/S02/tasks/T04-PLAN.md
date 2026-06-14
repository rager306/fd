---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: Lifecycle gate middleware на /v1/embeddings

api/middleware/lifecycle.go: gin middleware который проверяет IsReady() и !IsShuttingDown() перед передачей в handler. Если !IsReady → 503 model_not_loaded + Retry-After: 5. Если IsShuttingDown → 503 shutting_down + Retry-After: 30. Также TrackRequest(start, done) для inflight tracking. Подключается в router setup после validation (S01), до embed handler.

## Inputs

- None specified.

## Expected Output

- `api/middleware/lifecycle.go`
- `api/middleware/lifecycle_test.go`

## Verification

Unit tests: до warmup → 503 model_not_loaded, Retry-After: 5. После BeginShutdown → 503 shutting_down, Retry-After: 30. inflight counter инкрементируется на request и декрементируется на response.
