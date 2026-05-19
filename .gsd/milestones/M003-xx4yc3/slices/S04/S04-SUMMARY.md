---
id: S04
parent: M003-xx4yc3
milestone: M003-xx4yc3
provides:
  - Evidence base for S05 performance improvement recommendations.
requires:
  []
affects:
  - S05
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-baseline-py313.txt
  - benchmark-results/fd-runtime-stats-logs.txt
key_decisions:
  - Use `uv run --python 3.13 --with requests --with redis` as the benchmark execution method.
patterns_established:
  - Run Python benchmarks with `uv run --python 3.13 --with ...` to keep dependency/runtime selection explicit.
observability_surfaces:
  - benchmark-results/fd-benchmark-baseline-py313.txt
  - benchmark-results/fd-runtime-stats-logs.txt
drill_down_paths:
  - .gsd/milestones/M003-xx4yc3/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S04/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T08:24:26.585Z
blocker_discovered: false
---

# S04: Benchmark and resource baseline

**S04 established the uv Python 3.13 benchmark and resource baseline.**

## What Happened

S04 ran the benchmark baseline and collected resource/log evidence. Python dependencies were provided by uv with Python 3.13.12. Benchmark completed all sections and showed strong warm-cache performance: warm latency mean 2.00ms, cached p95 2.17ms, and max throughput around 742 req/s. Logs/stats show API and Redis are lightweight while TEI is the heavy component (~1.7GiB) and has backend batching constraints.

## Verification

Benchmark and stats/log evidence captured successfully.

## Requirements Advanced

- Benchmark and resource baseline captured. — 

## Requirements Validated

- Python 3.13.12 benchmark runtime verified. — 
- benchmark.py completed all sections. — 
- Cold/warm latency, repeated cached p95/p99, throughput, and response shape captured. — 
- Docker stats/logs captured after benchmark. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Benchmark was first run under uv default Python, then rerun under uv Python 3.13 after user clarification; Python 3.13 artifact is the accepted baseline.

## Known Limitations

Benchmark is local/dev only and mutates Redis with FLUSHALL. Stats were captured after benchmark, not continuously during peak load.

## Follow-ups

Proceed to S05: document no further defect fixes needed for benchmark path and create evidence-backed optimization plan.

## Files Created/Modified

- `benchmark-results/fd-benchmark-baseline-py313.txt` — Benchmark baseline output under uv Python 3.13.
- `benchmark-results/fd-runtime-stats-logs.txt` — Docker stats and runtime log highlights after benchmark.
