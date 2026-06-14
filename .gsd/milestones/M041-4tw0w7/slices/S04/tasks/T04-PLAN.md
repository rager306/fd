---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T04: Confirmed real cache-miss inference performance blocker: current fd + TEI CPU misses T-P latency targets when Redis namespace is isolated.

Если baseline из T01 НЕ достигает target (50/200/1000ms p95), выполнить targeted optimization. Возможные направления: (a) batch tensor packing — отправлять весь input array одним тензором в ONNX/TEI за раз, не per-item; (b) concurrent workers — если backend медленный, N goroutines обрабатывают chunks; (c) request coalescing — multiple requests с одинаковым input идут в один model call; (d) ONNX session warmup и graph optimization. Каждое изменение проверяется отдельным benchmark. Изменения оформляются как opt-in через env (FD_BATCH_TENSOR_PACKING=true, FD_CONCURRENT_WORKERS=4, etc).

## Inputs

- None specified.

## Expected Output

- `api/embed/optimizations.go`
- `api/embed/optimizations_test.go`
- `api/embed/perf_test.go`
- `benchmark-results/fd-v2-perf-m041-s04.md`

## Verification

После optimization: повторный baseline показывает p95 latency в target. benchmark-results/fd-v2-perf-m041-s04.md содержит before/after numbers и какие optimization сработали.
