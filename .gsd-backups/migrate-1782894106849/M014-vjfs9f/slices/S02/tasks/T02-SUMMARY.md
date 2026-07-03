---
id: T02
parent: S02
milestone: M014-vjfs9f
key_files:
  - benchmark-results/fd-benchmark-m014-tei-baseline.txt
key_decisions:
  - Use `benchmark-results/fd-benchmark-m014-tei-baseline.txt` as the TEI control artifact for S03 tagged ONNX comparison.
  - TEI baseline uses snapshot_version 2 and runtime_label `tei-default`.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:16:43.462Z
blocker_discovered: false
---

# T02: Ran the fresh TEI baseline benchmark with snapshot v2 metadata.

**Ran the fresh TEI baseline benchmark with snapshot v2 metadata.**

## What Happened

Ran the full benchmark against the default TEI API at `http://localhost:8000` using Python 3.13 via uv. The artifact includes snapshot_version 2, runtime label `tei-default`, Docker/Redis/git/environment metadata, and performance sections for cold/warm latency, repeated cache hits, concurrency, Redis L2 after API restart, cached batch endpoint, and repeated chunk reuse. The command exited 0 and wrote `benchmark-results/fd-benchmark-m014-tei-baseline.txt`.

## Verification

Benchmark command exited 0 and artifact includes snapshot_version 2 with `runtime_label=tei-default`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `BENCHMARK_RUNTIME_LABEL=tei-default BENCHMARK_API_URL=http://localhost:8000 BENCHMARK_MODEL=deepvk/USER-bge-m3 BENCHMARK_DIMENSIONS=1024 uv run --python 3.13 --with requests --with redis python benchmark.py > benchmark-results/fd-benchmark-m014-tei-baseline.txt` | 0 | ✅ pass — benchmark artifact written | 29600ms |
| 2 | `read benchmark-results/fd-benchmark-m014-tei-baseline.txt` | 0 | ✅ pass — snapshot_version=2; runtime_label=tei-default | 0ms |

## Deviations

The benchmark snapshot records git dirty=true because GSD S02 artifacts/plans and benchmark output were uncommitted during the run. This is expected inside an active slice; the commit hash is recorded and the final slice will commit the artifact.

## Known Issues

Git dirty=true in artifact due active GSD worktree state. Benchmark.py flushed Redis and may have restarted API per existing diagnostic behavior.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
