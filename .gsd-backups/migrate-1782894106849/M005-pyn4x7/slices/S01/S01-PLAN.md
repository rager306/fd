# S01: Benchmark README update

**Goal:** Update README benchmark/runtime instructions to match the validated uv Python 3.13 workflow and benchmark behavior.
**Demo:** After this, README has accurate uv Python 3.13 benchmark instructions and describes benchmark side effects.

## Must-Haves

- README documents `uv run --python 3.13 --with requests --with redis python benchmark.py`.
- README states benchmark assumes local Docker stack and Redis bound on localhost.
- README warns benchmark section 5 restarts API when Docker Compose is available.
- Commands in README match current files/config.

## Proof Level

- This slice proves: docs inspection plus command/config verification

## Integration Closure

Future maintainers can reproduce benchmark and understand its runtime side effects.

## Verification

- Docs explain benchmark artifacts, Redis localhost dependency, and API restart diagnostic.

## Tasks

- [x] **T01: Inspect current docs and runtime commands** `est:small`
  Inspect README, compose files, and benchmark.py to identify stale or missing benchmark/runtime instructions.
  - Files: `README.md`, `docker-compose.yaml`, `docker-compose.override.yaml`, `benchmark.py`
  - Verify: Findings recorded with exact README sections to update.

- [x] **T02: Update benchmark README docs** `est:small`
  Update README benchmark section with uv Python 3.13 command, local stack prerequisites, generated artifacts, and API restart side effect.
  - Files: `README.md`
  - Verify: README contains current command and side-effect warning.

- [x] **T03: Verify README benchmark docs** `est:small`
  Verify README command snippets against current compose config and benchmark.py syntax without rerunning full benchmark unless needed.
  - Files: `README.md`, `benchmark.py`
  - Verify: docker compose config and Python compile check pass.

## Files Likely Touched

- README.md
- docker-compose.yaml
- docker-compose.override.yaml
- benchmark.py
