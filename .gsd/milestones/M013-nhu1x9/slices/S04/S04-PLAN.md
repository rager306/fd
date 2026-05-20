# S04: Final benchmark readiness decision

**Goal:** Close M013 with the correct decision: the tagged HF tokenizer ONNX path is fixed-probe benchmark-ready, but not production-ready until native artifact Docker/CI packaging and larger quality gates are complete.
**Demo:** After this, M013 closes with proof that tagged ONNX is benchmark-ready on fixed probes, while production readiness remains gated by Docker/CI packaging and broader quality evaluation.

## Must-Haves

- Final recommendation says tagged ONNX path is benchmark-ready on fixed probes.
- Final recommendation says no production switch yet.
- Default TEI runtime health is verified.
- Default and tagged quality gates pass.
- Artifact/leak/native-binary checks pass.
- GitNexus scope is clean or explained.

## Proof Level

- This slice proves: Research synthesis plus fresh final verification gates.

## Integration Closure

Sets the next milestone direction for tagged ONNX performance benchmarking while preserving production safety boundaries.

## Verification

- Records final build/test/artifact/cosine state and exact commands needed by future benchmark work.

## Tasks

- [x] **T01: Synthesize final benchmark readiness decision** `est:small`
  Write final S04 research synthesis: native artifact contract, build-tag boundary, tagged cosine pass, remaining Docker/CI/quality gates, and next benchmark milestone recommendation.
  - Files: `.gsd/milestones/M013-nhu1x9/slices/S04/S04-RESEARCH.md`
  - Verify: Research artifact exists and states benchmark-ready but not production-ready.

- [x] **T02: Verify M013 final state** `est:small`
  Run final M013 verification: default tests/lint, tagged tests, health, artifact checks, GitNexus detect_changes, and no background process check.
  - Verify: All final gates pass.

## Files Likely Touched

- .gsd/milestones/M013-nhu1x9/slices/S04/S04-RESEARCH.md
