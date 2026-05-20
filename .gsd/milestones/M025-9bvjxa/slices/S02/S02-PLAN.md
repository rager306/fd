# S02: Operational diagnostics and rollout contract

**Goal:** Document startup diagnostics, failure modes, rollout steps, rollback to TEI, and observability requirements for ONNX artifacts.
**Demo:** After this, there is an explicit operational diagnostics and rollout/rollback contract for ONNX opt-in deployments.

## Must-Haves

- Operational diagnostics doc exists.
- Missing/mismatch artifact failure messages and health expectations are explicit.
- Rollout and rollback steps preserve TEI default.
- Decision records that operational gate is not yet implemented as production rollout.
- Docs pass hygiene checks.

## Proof Level

- This slice proves: Operations doc, decision, docs verification.

## Integration Closure

Defines operational readiness expectations before any ONNX deployment can be enabled.

## Verification

- Creates a runbook-style contract for actionable errors, health checks, safe logs, rollout stages, and rollback.

## Tasks

- [x] **T01: Write ONNX operations contract** `est:small`
  Write `docs/onnx-artifacts/OPERATIONS.md` covering startup preflight, failure diagnostics, health/status checks, safe logging fields, rollout stages, rollback to TEI, and production-default safeguards.
  - Files: `docs/onnx-artifacts/OPERATIONS.md`, `docs/onnx-artifacts/README.md`
  - Verify: Operations doc exists and is linked from README.

- [x] **T02: Record operational rollout decision** `est:small`
  Record GSD decision for ONNX operational rollout: opt-in staged rollout only, TEI default preserved, rollback by env/backend switch, production switch blocked until diagnostics are implemented and tested.
  - Files: `.gsd/DECISIONS.md`
  - Verify: Decision saved.

- [x] **T03: Verify operations contract** `est:small`
  Verify S02 docs: required sections present, no secrets/raw input text, binary hygiene unchanged.
  - Files: `docs/onnx-artifacts/OPERATIONS.md`
  - Verify: Section checks and binary hygiene pass.

## Files Likely Touched

- docs/onnx-artifacts/OPERATIONS.md
- docs/onnx-artifacts/README.md
- .gsd/DECISIONS.md
