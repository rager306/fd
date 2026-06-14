---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Реализовать async chunked orchestrator

api/embed/async.go: AsyncChunkedEmbed(ctx, teiClient, texts, dims) ([][]float32, error) с bounded concurrency semaphore (max 4, matches TEI max_batch_requests). Использует errgroup (golang.org/x/sync) или sync.WaitGroup + atomic errors. Cache logic не меняется — handler всё ещё делает GetIfPresent per text, но для miss'ов шлёт несколько chunks в parallel. Returns concatenated [][]float32. На любой chunk error — return wrapped error, no partial result.

## Inputs

- None specified.

## Expected Output

- `api/embed/async.go`
- `api/embed/async_test.go`
- `api/go.mod (errgroup dep)`

## Verification

Unit tests: (a) 4 chunks of 8, concurrency limit 4 → все chunks запускаются; (b) 1 chunk fails → wrapped error, no partial result; (c) all chunks success → concatenated result. Race detector clean.
