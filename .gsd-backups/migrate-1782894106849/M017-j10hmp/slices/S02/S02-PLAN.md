# S02: Quality outcome decision

**Goal:** Decide the next remediation path based on the measured tagged ONNX 512 legal quality outcome.
**Demo:** After this, the project has a concrete next remediation plan based on the 512-token gate outcome.

## Must-Haves

- S01 metrics are summarized.
- 512-token ONNX outcome is classified explicitly.
- Next remediation step is chosen.
- TEI remains production/default and ONNX remains experimental.
- Milestone closure verification passes.

## Proof Level

- This slice proves: Evidence-backed assessment and validation.

## Integration Closure

Closes M017 with a concrete next milestone recommendation tied to measured evidence.

## Verification

- Adds a decision artifact and decision-register entry for the 512 gate outcome.

## Tasks

- [x] **T01: Write 512 outcome assessment** `est:small`
  Write an outcome assessment for the 512-token ONNX gate, comparing M015 128-token failure, M016 Python 512 diagnostic, and M017 tagged Go 512 results.
  - Files: `benchmark-results/fd-onnx-512-outcome-m017-s02.txt`
  - Verify: Artifact exists, includes key metrics, and contains no raw legal text.

- [x] **T02: Record 512 outcome decision** `est:small`
  Record a GSD decision that 512-token ONNX is necessary but insufficient for strict legal equivalence, so the next implementation gate must add chunking or longer sequence handling.
  - Files: `.gsd/DECISIONS.md`, `.gsd/gsd.db`
  - Verify: Decision saved through GSD.

- [x] **T03: Validate M017 closure** `est:small`
  Run fresh closure verification and complete M017 if all gates pass for measurement scope.
  - Verify: Fresh verification passes and no background processes remain.

## Files Likely Touched

- benchmark-results/fd-onnx-512-outcome-m017-s02.txt
- .gsd/DECISIONS.md
- .gsd/gsd.db
