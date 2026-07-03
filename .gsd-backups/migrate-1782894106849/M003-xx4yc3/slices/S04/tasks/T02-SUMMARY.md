---
id: T02
parent: S04
milestone: M003-xx4yc3
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-baseline-py313.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:23:02.389Z
blocker_discovered: false
---

# T02: Benchmark baseline completed under uv with Python 3.13.12.

**Benchmark baseline completed under uv with Python 3.13.12.**

## What Happened

Ran benchmark.py against the live local stack using uv with Python 3.13.12 and dependencies supplied by uv. Benchmark completed all sections: cold/warm latency, 100 repeated cache-hit requests, concurrent throughput, and response format verification. Accepted baseline artifact is `benchmark-results/fd-benchmark-baseline-py313.txt`. Results: best cold latency 19.5ms, warm latency mean 2.00ms, repeated request mean 1.55ms with p95 2.17ms, max throughput about 742 req/s at 16 concurrency, response dimensions 1024.

## Verification

`uv run --python 3.13 --with requests --with redis python benchmark.py` completed successfully and wrote `benchmark-results/fd-benchmark-baseline-py313.txt`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests --with redis python --version && uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-baseline-py313.txt` | 0 | ✅ pass: Python 3.13.12; warm mean 2.00ms; max throughput ~742 req/s | 22500ms |

## Deviations

The first benchmark run used uv with default Python; after user clarified Python 3.13, benchmark was rerun with `uv run --python 3.13` and the Python 3.13 result is the accepted baseline. Earlier pip-target setup remains a discarded temp setup and was not used for the accepted benchmark.

## Known Issues

None.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-baseline-py313.txt`
