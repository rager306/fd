# S04: Benchmark and resource baseline

**Goal:** Run benchmark.py and correlate results with container stats/logs.
**Demo:** Benchmark baseline captures latency, throughput, resources, and logs.

## Must-Haves

- Python dependencies are available or installed in isolated env.
- `benchmark.py` completes or blocker is diagnosed.
- Cold/warm latency, p95/p99, throughput, and response format evidence captured.
- Docker stats/logs collected after benchmark.

## Proof Level

- This slice proves: benchmark output plus docker stats/logs

## Integration Closure

Produces measured baseline before proposing optimizations.

## Verification

- Creates benchmark artifact and resource/log evidence.

## Tasks

- [x] **T01: Prepare benchmark environment** `est:small`
  Verify Python benchmark dependencies without installing globally.
  - Files: `benchmark.py`
  - Verify: python3 imports requests and redis.

- [x] **T02: Run benchmark baseline** `est:medium`
  Run benchmark.py against local stack and save output artifact.
  - Files: `benchmark.py`
  - Verify: python3 benchmark.py completes all sections.

- [x] **T03: Correlate benchmark with resources** `est:small`
  Capture docker stats and recent logs after benchmark, then summarize bottlenecks.
  - Verify: docker stats --no-stream and logs captured.

## Files Likely Touched

- benchmark.py
