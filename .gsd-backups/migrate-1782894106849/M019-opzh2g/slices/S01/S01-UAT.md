# S01: ONNX 1024 performance benchmark — UAT

**Milestone:** M019-opzh2g
**Written:** 2026-05-20T08:17:11.311Z

# S01 UAT — ONNX 1024 performance benchmark

## Checks

- [x] Benchmark command plan created.
- [x] Runtime helper starts ONNX 1024 and `/health` passes.
- [x] Benchmark artifact exists.
- [x] Artifact records `ONNX_MAX_SEQUENCE_LENGTH=1024`.
- [x] Artifact records isolated namespace `m019-onnx-1024-benchmark`.
- [x] Artifact hygiene passed with `raw_benchmark_text_leaks=0`.
- [x] Runtime helper stopped service.
- [x] Port 18000 is clean.
- [x] No background processes remain.

## Key result

- Best cold latency: `8.3ms`.
- Warm latency mean: `1.19ms`.
- Max throughput: `~858 req/s`.
- Redis L2 restart: `2.02ms`.
- Batch L1 p95: `4.09ms`.
- Batch L2 p95: `6.70ms`.
- Chunk reuse warm p95: `6.51ms`.

## UAT Result

Pass. S02 should decide whether to proceed to packaging/CI gates or performance tuning.

