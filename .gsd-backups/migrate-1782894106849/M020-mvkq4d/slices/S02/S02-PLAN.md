# S02: Contract validation and next gate

**Goal:** Validate the metadata contract and record the next Docker/CI packaging gate.
**Demo:** After this, M020 records the packaging/CI follow-up gate and closes with fresh tests and artifact hygiene.

## Must-Haves

- Contract decision is recorded.
- Metadata validation remains passing.
- No binary artifacts are tracked.
- Fresh verification passes.
- Next gate is Docker/CI packaging and artifact provisioning.

## Proof Level

- This slice proves: Fresh verification, GSD decision, GitNexus scope check.

## Integration Closure

Closes M020 with a clean next step for operationalization.

## Verification

- Adds decision/summary artifacts so future agents know ONNX 1024 contract status and remaining blockers.

## Tasks

- [x] **T01: Record artifact contract decision** `est:small`
  Record a GSD decision that the tracked ONNX manifest now represents an experimental 1024 runtime contract, not production readiness.
  - Files: `.gsd/DECISIONS.md`, `.gsd/gsd.db`
  - Verify: Decision saved through GSD.

- [x] **T02: Validate M020 closure** `est:small`
  Run fresh metadata validation, tracked binary checks, relevant tests/lint, and GitNexus scope check before closing M020.
  - Verify: Fresh verification passes and no background processes remain.

## Files Likely Touched

- .gsd/DECISIONS.md
- .gsd/gsd.db
