---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Benchmark ONNX cold and warm performance with cache namespace isolation

Run ONNX cold/warm benchmark and async-combo measurement while isolating Redis cache namespace (e.g., EMBEDDING_CACHE_VERSION) or flushing deliberately so cached TEI vectors cannot masquerade as ONNX output. Capture cold batch=32, warm batch=1, and ONNX+async combo where applicable. Record runtime metadata and cache namespace in the benchmark artifact. If S02 async artifacts exist by then, compare against them; otherwise document the missing dependency.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/static-analysis-recommendation.md`

## Expected Output

- `benchmark-results/fd-v2-onnx-perf-m042.md`
- `tools/verify_fd_onnx_perf.sh`

## Verification

cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && ../tools/verify_fd_onnx_perf.sh

## Observability Impact

Perf artifact includes runtime model, dimensions, backend, and cache namespace evidence.
