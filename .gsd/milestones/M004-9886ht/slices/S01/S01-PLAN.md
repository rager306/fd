# S01: Benchmark summary correctness

**Goal:** Fix benchmark throughput summary so it reports the true maximum row from measured throughput data.
**Demo:** After this, benchmark.py summary cannot contradict its throughput table.

## Must-Haves

- Throughput summary selects the true max req/s row.
- Associated concurrency is reported correctly.
- Fix is covered by deterministic verification.
- `uv run --python 3.13 --with requests --with redis python benchmark.py` summary no longer contradicts throughput table.

## Proof Level

- This slice proves: helper-level verification plus live benchmark smoke

## Integration Closure

Benchmark summary becomes reliable evidence for later optimization work.

## Verification

- Prevents misleading benchmark output.

## Tasks

- [x] **T01: Inspect benchmark summary bug** `est:small`
  Inspect benchmark.py throughput aggregation, identify the incorrect summary selection, and run GitNexus impact analysis before editing.
  - Files: `benchmark.py`
  - Verify: Root cause identified and impact analysis recorded.

- [x] **T02: Fix throughput summary aggregation** `est:small`
  Refactor the smallest benchmark.py helper or aggregation logic needed so summary max throughput is derived from measured rows and its concurrency matches.
  - Files: `benchmark.py`
  - Verify: Deterministic helper verification or benchmark dry run passes.

- [x] **T03: Verify benchmark summary fix** `est:medium`
  Run uv Python 3.13 benchmark smoke against the live stack and confirm summary max matches table max.
  - Files: `benchmark-results/fd-benchmark-m004-s01.txt`
  - Verify: `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s01.txt` and compare summary to table.

## Files Likely Touched

- benchmark.py
- benchmark-results/fd-benchmark-m004-s01.txt
