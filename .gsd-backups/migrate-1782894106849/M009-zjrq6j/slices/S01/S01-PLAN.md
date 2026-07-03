# S01: Benchmark config snapshot

**Goal:** Add sanitized effective configuration snapshot output to benchmark artifacts without changing runtime behavior.
**Demo:** After this, every benchmark run reports the effective safe configuration needed to compare future results.

## Must-Haves

- Benchmark output includes a sanitized effective config section.
- Snapshot includes git commit/branch, Docker image/compose identifiers, cache/runtime env fields, Redis INFO summary, and environment artifact reference.
- Secret-like variables and raw benchmark texts are excluded.
- Existing benchmark behavior remains intact.

## Proof Level

- This slice proves: benchmark run plus parser checks for expected fields and forbidden secret-like fields

## Integration Closure

Benchmark artifacts become comparable evidence for later cache/runtime experiments.

## Verification

- Adds git, Docker, env-derived, Redis, and environment-reference diagnostics to benchmark output while redacting secrets and raw text.

## Tasks

- [x] **T01: Design benchmark config snapshot** `est:small`
  Inspect current benchmark output flow and design a small sanitized config snapshot helper. Identify secret-redaction rules, source fields, and where the section should be printed without changing measured requests.
  - Files: `benchmark.py`
  - Verify: Design summary names included fields, excluded secret patterns, and insertion point.

- [x] **T02: Implement benchmark config snapshot** `est:medium`
  Implement the sanitized effective config snapshot in `benchmark.py`. Include git metadata, Docker compose/image identifiers when available, selected env/runtime/cache settings, Redis INFO summary when available, and environment artifact reference. Redact or omit secret-like keys and never print raw benchmark input texts.
  - Files: `benchmark.py`
  - Verify: Run Python compile and a targeted snapshot parser/check.

- [x] **T03: Verify benchmark config snapshot** `est:medium`
  Run full verification for S01: Go tests, Docker compose config, benchmark run through uv Python 3.13, and a parser check proving the snapshot exists and excludes secret-like fields. Save/inspect benchmark artifact if generated.
  - Files: `benchmark.py`, `benchmark-results/`
  - Verify: `docker compose config`; `cd api && go test ./... -short`; `uv run --python 3.13 --with requests --with redis python benchmark.py`; parser check for snapshot and redaction.

## Files Likely Touched

- benchmark.py
- benchmark-results/
