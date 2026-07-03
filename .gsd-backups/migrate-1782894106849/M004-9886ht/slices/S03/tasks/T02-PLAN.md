---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Add Redis restart cache diagnostic

Implement Redis-persisted cache diagnostic section in benchmark.py. It should prime Redis, restart API if Docker Compose is available, wait for health, then measure the same text and report skip reason if restart/check fails.

## Inputs

- `T01 design`

## Expected Output

- `benchmark.py`

## Verification

Python 3.13 compile check passes.

## Observability Impact

Benchmark can prove L2 persistence after API restart.
