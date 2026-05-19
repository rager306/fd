# S03: Benchmark diagnostic modes

**Goal:** Extend benchmark.py to produce diagnostic cache-path evidence, especially Redis-persisted cache behavior after API restart.
**Demo:** After this, benchmark.py can produce more diagnostic cache-path evidence using uv Python 3.13.

## Must-Haves

- Benchmark has a documented section for Redis-persisted cache behavior after API restart or clearly reports a skip reason.
- Output remains concise and commit-worthy.
- Python execution uses `uv run --python 3.13`.
- Existing benchmark sections still pass.

## Proof Level

- This slice proves: uv Python 3.13 benchmark run against live Docker stack

## Integration Closure

Benchmark can guide future tuning by distinguishing normal warm cache from cache persistence across API restarts.

## Verification

- Adds benchmark evidence for L2 Redis cache usefulness after process restart.

## Tasks

- [x] **T01: Design benchmark diagnostic extension** `est:small`
  Inspect benchmark.py, run GitNexus impact analysis, and design the smallest diagnostic section for L2 Redis cache after API restart.
  - Files: `benchmark.py`
  - Verify: Impact analysis recorded for benchmark.py main.

- [x] **T02: Add Redis restart cache diagnostic** `est:medium`
  Implement Redis-persisted cache diagnostic section in benchmark.py. It should prime Redis, restart API if Docker Compose is available, wait for health, then measure the same text and report skip reason if restart/check fails.
  - Files: `benchmark.py`
  - Verify: Python 3.13 compile check passes.

- [x] **T03: Verify benchmark diagnostic output** `est:medium`
  Run the updated benchmark through uv Python 3.13 and save evidence artifact.
  - Files: `benchmark-results/fd-benchmark-m004-s03.txt`
  - Verify: `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s03.txt` passes and includes L2 diagnostic result or skip reason.

## Files Likely Touched

- benchmark.py
- benchmark-results/fd-benchmark-m004-s03.txt
