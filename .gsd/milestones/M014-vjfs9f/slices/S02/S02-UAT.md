# S02: TEI baseline benchmark — UAT

**Milestone:** M014-vjfs9f
**Written:** 2026-05-20T04:20:32.936Z

# S02 UAT — TEI baseline benchmark

## Checks

- [x] Docker Compose default stack healthy: `fd_api`, `fd_redis`, `fd_tei`.
- [x] API `/health` returned ok.
- [x] No tagged ONNX background server running.
- [x] TEI benchmark command exited 0.
- [x] Artifact exists: `benchmark-results/fd-benchmark-m014-tei-baseline.txt`.
- [x] Artifact includes `snapshot_version: 2` and `runtime_label: tei-default`.
- [x] Artifact includes cache, throughput, batch, and chunk reuse sections.
- [x] Raw fixed-probe text leak check returned 0.

## Key TEI Results

- Best cold latency: 59.0ms.
- Warm latency mean: 2.25ms.
- Max throughput: ~750 req/s at 16 concurrent.
- Redis L2 after API restart: 2.82ms.
- Batch L1 p95: 4.16ms.
- Batch L2 p95: 5.51ms.
- Chunk reuse warm p95: 7.04ms.

## UAT Result

Pass. This artifact is ready as the TEI control for S03/S04.

