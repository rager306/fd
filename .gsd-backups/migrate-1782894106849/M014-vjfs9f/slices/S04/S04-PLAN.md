# S04: Benchmark synthesis and decision

**Goal:** Compare TEI and tagged ONNX artifacts, validate milestone, and recommend next step.
**Demo:** After this, the project has a clear data-backed recommendation for ONNX performance work.

## Must-Haves

- Comparison states faster/slower by scenario.
- Caveats include fixed-probe vs corpus quality and native packaging gaps.
- Final tests/lint/tagged tests pass.
- GitNexus and artifact hygiene checks pass.
- No production switch or push occurs.

## Proof Level

- This slice proves: Synthesis plus final verification gates.

## Integration Closure

Closes the benchmark milestone and determines whether to tune ONNX, package for CI/Docker, or stop.

## Verification

- Summarizes deltas, caveats, and next operational evidence needed.

## Tasks

- [x] **T01: Compare benchmark artifacts** `est:small`
  Parse TEI and tagged ONNX benchmark artifacts and write a concise comparison artifact with deltas and caveats.
  - Files: `benchmark-results/fd-benchmark-m014-comparison.txt`
  - Verify: Comparison artifact exists and includes scenario deltas plus caveats.

- [x] **T02: Record benchmark recommendation** `est:small`
  Record the benchmark recommendation in GSD decisions and slice summary: tagged ONNX is faster for cold/direct model work, but production switch remains blocked by packaging and quality gates.
  - Files: `.gsd/DECISIONS.md`
  - Verify: Decision saved and summary references comparison artifact.

- [x] **T03: Run final M014 verification gates** `est:medium`
  Run final verification gates: Go tests, pinned lint, tagged tests, artifact hygiene, tracked binary checks, GitNexus detect, milestone validation.
  - Verify: All commands pass and no tagged server/background process remains.

## Files Likely Touched

- benchmark-results/fd-benchmark-m014-comparison.txt
- .gsd/DECISIONS.md
