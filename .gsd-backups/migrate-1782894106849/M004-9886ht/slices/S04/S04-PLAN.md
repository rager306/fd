# S04: Verification and closure

**Goal:** Run final verification, close M004, and create a local commit with GSD artifacts and code changes.
**Demo:** After this, the optimization milestone is validated, summarized, and locally committed.

## Must-Haves

- `docker compose config` passes.
- `cd api && go test ./... -short` passes.
- `uv run --python 3.13 --with requests --with redis python benchmark.py` smoke/evidence passes.
- GitNexus change detection is low/expected.
- Milestone is validated/completed and local commit created.

## Proof Level

- This slice proves: full verification gates

## Integration Closure

Leaves repository with committed measured optimization work and clean state.

## Verification

- Records final verification evidence and benchmark artifacts.

## Tasks

- [x] **T01: Run final verification gates** `est:medium`
  Run final verification commands: compose config, Go tests, benchmark evidence check, docker ps, GitNexus change detection.
  - Files: `benchmark.py`, `api/`, `benchmark-results/`
  - Verify: All commands pass.

- [x] **T02: Validate and complete milestone** `est:small`
  Validate and complete M004 with success criteria, slice audit, requirement coverage, and summary.
  - Files: `.gsd/milestones/M004-9886ht`
  - Verify: GSD validate and complete tools succeed.

- [x] **T03: Commit local optimization work** `est:small`
  Checkpoint GSD DB, stage relevant files, create local commit, and report status. Do not push without explicit user confirmation.
  - Files: `.gsd/gsd.db`, `.gsd/milestones/M004-9886ht`, `api/`, `benchmark.py`, `benchmark-results/`
  - Verify: git status clean except ahead marker; commit hash reported.

## Files Likely Touched

- benchmark.py
- api/
- benchmark-results/
- .gsd/milestones/M004-9886ht
- .gsd/gsd.db
