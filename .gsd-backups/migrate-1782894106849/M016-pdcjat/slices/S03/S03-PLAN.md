# S03: Remediation option assessment

**Goal:** Choose and document the ONNX legal-divergence remediation path using S01/S02 evidence, without switching production defaults.
**Demo:** After this, there is a concrete next remediation plan for ONNX quality.

## Must-Haves

- Root-cause evidence from S01/S02 is summarized.
- Remediation options are compared with tradeoffs.
- Recommended path is explicit and does not switch production/default runtime.
- Follow-up verification gates are listed: full legal retrieval rerun, cache namespace isolation, tagged/native build checks, and performance rerun only after quality passes.
- No raw legal corpus text is included.

## Proof Level

- This slice proves: Evidence-backed remediation plan and milestone validation.

## Integration Closure

Produces the next implementation gate for future ONNX work: quality-first 512-token/long-text remediation before packaging or performance claims.

## Verification

- Adds a sanitized remediation decision artifact tying diagnostics to follow-up verification gates.

## Tasks

- [x] **T01: Write remediation assessment** `est:small`
  Write a remediation assessment artifact comparing longer ONNX max sequence length, explicit chunking, and longer-sequence export/runtime options using S01/S02 evidence.
  - Files: `benchmark-results/fd-onnx-remediation-plan-m016-s03.txt`
  - Verify: Artifact exists, includes S01/S02 metrics, and contains no raw legal text.

- [x] **T02: Record remediation decision** `est:small`
  Record the GSD decision that TEI remains default and the next ONNX implementation gate is 512-token/long-text remediation plus full legal corpus quality rerun.
  - Files: `.gsd/DECISIONS.md`, `.gsd/gsd.db`
  - Verify: Decision is saved through gsd_decision_save and references the quality-first remediation path.

- [x] **T03: Validate milestone closure** `est:small`
  Validate M016 closure readiness: confirm S01/S02/S03/S04 outputs, run lightweight script checks and artifact hygiene checks, then prepare milestone completion if all slices are complete.
  - Verify: GSD milestone validation passes or records any remediation gaps.

## Files Likely Touched

- benchmark-results/fd-onnx-remediation-plan-m016-s03.txt
- .gsd/DECISIONS.md
- .gsd/gsd.db
