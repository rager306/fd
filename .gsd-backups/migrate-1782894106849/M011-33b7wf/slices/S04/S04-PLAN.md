# S04: Blocker synthesis and recommendation

**Goal:** Synthesize the M011 outcome honestly: artifact/config/Go ONNX runtime integration works, but semantic equivalence is blocked by tokenizer parity, so throughput benchmarking and production recommendations are deferred.
**Demo:** After this, the project has evidence that Go ONNX runtime integration loads/runs but is blocked on tokenizer parity, plus a recommendation for the next milestone.

## Must-Haves

- Recommendation clearly states ONNX Go backend is blocked on tokenizer parity.
- Evidence references S01 manifest, S02 validation seam, and S03 failed isolated-cache comparison.
- Future gates are explicit: tokenizer parity, stable ORT shared library, artifact storage, quality/performance.
- Default TEI path remains verified.
- Milestone can be validated without claiming ONNX production readiness.

## Proof Level

- This slice proves: Research synthesis plus fresh Go tests/lint/config/artifact checks/GitNexus verification.

## Integration Closure

Closes M011 with a clear next-step recommendation and verification that TEI default remains safe.

## Verification

- Captures the cache masking pitfall, tokenizer mismatch evidence, and required future diagnostic checks for backend comparisons.

## Tasks

- [x] **T01: Synthesize tokenizer parity blocker recommendation** `est:small`
  Write S04 research synthesis with final recommendation: stop M011 at blocker, do not benchmark throughput yet, and plan tokenizer parity research/remediation as the next gate. Include cache namespace masking lesson and shared-library caveat.
  - Files: `.gsd/milestones/M011-33b7wf/slices/S04/S04-RESEARCH.md`
  - Verify: Research artifact exists and states blocker, recommendation, and no production switch.

- [x] **T02: Verify M011 blocked prototype safety** `est:small`
  Run final M011 verification: Go tests, pinned lint, Docker Compose config, default TEI health/API smoke if available, manifest parser, artifact checks, no raw probe leakage, and GitNexus detect_changes. Record evidence for milestone validation.
  - Verify: All verification gates pass; final recommendation remains blocker, not production-ready.

## Files Likely Touched

- .gsd/milestones/M011-33b7wf/slices/S04/S04-RESEARCH.md
