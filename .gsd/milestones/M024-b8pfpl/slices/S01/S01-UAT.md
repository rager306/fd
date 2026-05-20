# S01: Packaged ONNX performance benchmark — UAT

**Milestone:** M024-b8pfpl
**Written:** 2026-05-20T11:27:58.058Z

# S01 UAT — Packaged ONNX performance benchmark

## Checks

- [x] Packaged ONNX container healthy on port 18000.
- [x] Smoke embedding returns 1024 dimensions.
- [x] Restart command targets `fd-onnx-m024-bench`.
- [x] Benchmark ran with `uv run --python 3.13`.
- [x] Artifact contains packaged runtime labels and cache namespace.
- [x] Raw synthetic benchmark input leak check passed.
- [x] Container stopped; port 18000 clean.

## Key Metrics

- Best cold latency: `7.6ms`
- Warm latency mean: `2.03ms`
- Max throughput: `~864 req/s`
- Redis L2 restart: `3.36ms`
- Batch L1 p95: `8.91ms`
- Batch L2 p95: `5.47ms`
- Chunk reuse warm p95: `5.24ms`

## Result

Pass. Packaged ONNX Docker remains locally performance-viable.

