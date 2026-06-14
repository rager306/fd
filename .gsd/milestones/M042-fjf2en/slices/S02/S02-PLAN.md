# S02: Async parallel chunked TEI calls in handler

**Goal:** Async parallel chunked TEI calls в handler. Bounded concurrency 4 (matches TEI max_batch_requests). Env FD_ASYNC_CHUNKS=true включает. Error aggregation, X-Concurrent-Chunks header, metrics.
**Demo:** After this, FD_ASYNC_CHUNKS=true enables parallel chunking: cold path for batch=128 падает с 25s до ≤10s, batch=32 с 6s до ≤4s. Cache hit path не regressed. New X-Concurrent-Chunks header в response для observability.

## Must-Haves

- tools/verify_fd_async_perf.sh существует, exit 0. FD_ASYNC_CHUNKS=true: cold path batch=32 ≤4s, batch=128 ≤10s (was 25s sequential). FD_ASYNC_CHUNKS=false (default): cold path не regressed. Cache hit path не regressed (≤5ms per request в обоих режимах). X-Concurrent-Chunks response header присутствует в async mode. All M041 acceptance tests (45 test cases + 10 behavior scenarios) pass в обоих режимах. benchmark-results/fd-v2-async-perf-m042.md с before/after numbers.

## Proof Level

- This slice proves: runtime + integration + benchmark

## Integration Closure

Builds on M041 S01 chunked handler. Adds errgroup with semaphore(4). Modifies main.go to read FD_ASYNC_CHUNKS env. Cache logic (GetIfPresent/Set) unchanged. Headers/metrics from M041 S03 extended.

## Verification

- New X-Concurrent-Chunks header per response (only in async mode). New metrics: fd_async_chunks_total counter, fd_async_chunk_duration_seconds histogram.

## Tasks

- [ ] **T01: Реализовать async chunked orchestrator** `est:3h`
  api/embed/async.go: AsyncChunkedEmbed(ctx, teiClient, texts, dims) ([][]float32, error) с bounded concurrency semaphore (max 4, matches TEI max_batch_requests). Использует errgroup (golang.org/x/sync) или sync.WaitGroup + atomic errors. Cache logic не меняется — handler всё ещё делает GetIfPresent per text, но для miss'ов шлёт несколько chunks в parallel. Returns concatenated [][]float32. На любой chunk error — return wrapped error, no partial result.
  - Files: `api/embed/async.go`, `api/embed/async_test.go`, `api/go.mod`
  - Verify: Unit tests: (a) 4 chunks of 8, concurrency limit 4 → все chunks запускаются; (b) 1 chunk fails → wrapped error, no partial result; (c) all chunks success → concatenated result. Race detector clean.

- [ ] **T02: Wire FD_ASYNC_CHUNKS env в handler + main.go** `est:1h`
  api/handlers/embeddings.go: в CreateEmbedding, если FD_ASYNC_CHUNKS=true → использовать async orchestrator (chunks of 32 sent in parallel), иначе current sequential loop. api/main.go: parse FD_ASYNC_CHUNKS env на startup, log mode. Sync mode default (off) для backward compat.
  - Files: `api/handlers/embeddings.go`, `api/main.go`
  - Verify: FD_ASYNC_CHUNKS=true → 4 parallel TEI calls per request (verify via TEI logs overlapping timestamps). FD_ASYNC_CHUNKS=false (default) → sequential (no regression vs M041). Integration test asserts both modes pass.

- [ ] **T03: X-Concurrent-Chunks header + metrics** `est:1h`
  api/middleware/headers.go (M041 S03 deliverable) расширить: в async mode добавить response header X-Concurrent-Chunks: N (number of chunks sent in parallel for this request). api/observability/metrics.go (M041 S03): добавить fd_async_chunks_total counter (incremented per chunk in flight), fd_async_chunk_duration_seconds histogram (per chunk inference time). Sync mode — headers/metrics absent (no overhead).
  - Files: `api/middleware/headers.go`, `api/observability/metrics.go`
  - Verify: curl -I -X POST /v1/embeddings с FD_ASYNC_CHUNKS=true показывает X-Concurrent-Chunks: 4 для batch=128. /metrics показывает fd_async_chunks_total counter incrementing.

- [ ] **T04: Perf benchmark: async vs sync cold/warm path** `est:2h`
  tools/verify_fd_async_perf.sh: прогон cold path measurements × batch sizes × async on/off. Output benchmark-results/fd-v2-async-perf-m042.md с table (batch 1/10/32/64/128 cold, sync vs async) и conclusion (improvement factor, where it falls short of 1000ms target). Также test concurrent scenario: 4 parallel fd calls × batch=32 cold, sync vs async (per M041 T-P-5 spec).
  - Files: `tools/verify_fd_async_perf.sh`, `benchmark-results/fd-v2-async-perf-m042.md`
  - Verify: tools/verify_fd_async_perf.sh exit 0. Artifact содержит: cold path table (batch × async on/off), concurrent test results, conclusion с comparison vs M041 baseline (25s → ≤10s for batch=128 cold).

- [ ] **T05: Regression suite: M041 acceptance в async mode** `est:1h`
  tools/verify_fd_v2_contract.py (M041 deliverable, deferred) запустить с FD_ASYNC_CHUNKS=true. Все 45 test cases должны pass. Любой failure — bug regression в S02. Проверить особенно: encoding_format=base64 с batch=128 (4 chunks, 128 embeddings in base64), dimensions=512, validation envelope (413, 400 etc) — всё должно работать identical к sync mode.
  - Files: `tools/verify_fd_v2_contract.py`, `tests/integration/fd_v2_async_test.go`
  - Verify: go test ./tests/integration/... -run TestFdV2AsyncMode: все M041 acceptance pass в async mode. 0 regressions.

## Files Likely Touched

- api/embed/async.go
- api/embed/async_test.go
- api/go.mod
- api/handlers/embeddings.go
- api/main.go
- api/middleware/headers.go
- api/observability/metrics.go
- tools/verify_fd_async_perf.sh
- benchmark-results/fd-v2-async-perf-m042.md
- tools/verify_fd_v2_contract.py
- tests/integration/fd_v2_async_test.go
