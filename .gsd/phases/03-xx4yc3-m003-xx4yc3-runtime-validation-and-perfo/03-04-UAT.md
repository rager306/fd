# S04: Benchmark and resource baseline — UAT

**Milestone:** M003-xx4yc3
**Written:** 2026-05-19T08:24:26.585Z

# UAT: S04 Benchmark and resource baseline

## Verification performed

- `uv run --python 3.13 --with requests --with redis python --version` — Python 3.13.12.
- `uv run --python 3.13 --with requests --with redis python benchmark.py` — completed all sections.
- Artifact: `benchmark-results/fd-benchmark-baseline-py313.txt`.
- Artifact: `benchmark-results/fd-runtime-stats-logs.txt`.

## Key results

- Best cold latency: 19.5ms.
- Warm latency mean: 2.00ms.
- Repeated cached mean: 1.55ms; p95: 2.17ms; p99: 4.72ms.
- Max throughput: ~742 req/s at 16 concurrent.
- Docker stats: API ~9.6MiB, Redis ~6.1MiB, TEI ~1.7GiB.

## Observations

- API logs successful embeddings at INFO for every request, which is noisy during throughput tests.
- TEI logs show historical ONNX fallback to Candle CPU and `max_batch_requests=4` backend constraint.
- Redis host memory overcommit warning remains a host tuning item.

