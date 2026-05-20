# S02: TEI baseline benchmark

**Goal:** Run a fresh TEI+Redis baseline benchmark with sanitized config snapshot and known cache/runtime state.
**Demo:** After this, there is a fresh TEI default benchmark artifact to compare against tagged ONNX.

## Must-Haves

- TEI stack healthy before benchmark.
- Benchmark artifact includes snapshot_version 2 and runtime label `tei-default`.
- Redis cache effects are explicit.
- Artifact parser/leak checks pass.
- Baseline result is committed as evidence.

## Proof Level

- This slice proves: Executable benchmark artifact and parser checks.

## Integration Closure

Provides the control measurement for S03 tagged ONNX comparison.

## Verification

- Records TEI runtime health, cache namespace/settings, Docker state, Redis behavior, and performance output under snapshot v2.

## Tasks

- [x] **T01: Preflight TEI benchmark runtime** `est:small`
  Verify the current default Docker stack and runtime health before running the TEI baseline. Capture compose/health state and ensure no tagged ONNX server is running.
  - Verify: Docker compose ps and API health show default TEI stack healthy; no background tagged server.

- [x] **T02: Run TEI baseline benchmark** `est:medium`
  Run `benchmark.py` against default TEI API with snapshot v2 metadata and write `benchmark-results/fd-benchmark-m014-tei-baseline.txt`. Use Python 3.13 via uv and preserve existing benchmark safety behavior.
  - Files: `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
  - Verify: Benchmark command exits 0 and artifact includes snapshot_version 2.

- [x] **T03: Verify TEI benchmark artifact** `est:small`
  Parse TEI artifact for required sections, snapshot fields, Redis/cache sections, PASS/summary markers, and raw probe text absence.
  - Files: `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
  - Verify: Parser/leak checks and GitNexus detect_changes pass.

## Files Likely Touched

- benchmark-results/fd-benchmark-m014-tei-baseline.txt
