# S03: Benchmark diagnostic modes — UAT

**Milestone:** M004-9886ht
**Written:** 2026-05-19T10:30:22.503Z

# UAT: S03 Benchmark diagnostic modes

## Verification performed

- `uv run --python 3.13 python -m py_compile benchmark.py` — passed.
- `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s03.txt` — passed.
- Parsed benchmark output:
  - throughput table max: `754.2 req/s` at concurrency `4`
  - summary: `~754 req/s` at concurrency `4`
  - Redis L2 restart: `3.10ms`
- `docker compose ps` — api, redis, tei healthy.
- `gitnexus_detect_changes(repo=fd, scope=all)` — low risk, no affected processes.

## Result

Benchmark now reports Redis L2 cache behavior after API restart. The S03 evidence run showed a `299.11ms` prime/cold request and `3.10ms` after API restart, indicating Redis L2 persistence and L1 backfill.

