---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify benchmark summary fix

Run uv Python 3.13 benchmark smoke against the live stack and confirm summary max matches table max.

## Inputs

- `benchmark.py`

## Expected Output

- `benchmark-results/fd-benchmark-m004-s01.txt`
- `S01 summary`

## Verification

`uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s01.txt` and compare summary to table.

## Observability Impact

Produces evidence artifact for the fix.
