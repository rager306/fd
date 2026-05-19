# S01: Benchmark summary correctness — UAT

**Milestone:** M004-9886ht
**Written:** 2026-05-19T10:16:31.892Z

# UAT: S01 Benchmark summary correctness

## Verification performed

- `uv run --python 3.13 python -m py_compile benchmark.py` — passed.
- `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s01.txt` — passed.
- Parsed `benchmark-results/fd-benchmark-m004-s01.txt`:
  - table max: `742.8 req/s` at concurrency `16`
  - summary: `~743 req/s` at concurrency `16`
- `gitnexus_detect_changes(repo=fd, scope=all)` — low risk, no affected processes.

## Result

Benchmark summary now selects the true measured max throughput row instead of printing the final loop iteration.

