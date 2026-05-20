# S02: Docker CI boundary validation

**Goal:** Validate Docker/CI packaging boundaries and record the next artifact provisioning gate.
**Demo:** After this, local Docker/CI checks prove the default path is not broken and the ONNX packaging path has a clear next gate or local proof.

## Must-Haves

- Default Docker build path is validated or a concrete environment blocker is recorded.
- Artifact verifier passes locally.
- Tagged native tests pass.
- No binary artifacts are tracked.
- Decision records next gate as Docker/CI artifact provisioning.
- TEI remains production/default.

## Proof Level

- This slice proves: Docker build if available, tests/lint/tagged checks, binary hygiene, GSD decision.

## Integration Closure

Closes M021 with proof that default Docker/CI remains unaffected and opt-in ONNX packaging has a clear next implementation step.

## Verification

- Captures default Docker build/test status and artifact contract decision.

## Tasks

- [x] **T01: Validate default Docker boundary** `est:medium`
  Run the default Docker build to prove the TEI/default path still does not require ONNX/native artifacts.
  - Verify: `docker build -f api/Dockerfile -t fd-api:m021-default api` exits 0, or records a concrete Docker environment blocker.

- [x] **T02: Record packaging contract decision** `est:small`
  Record a GSD decision that M021 added artifact verification contract and the next gate is Docker/CI provisioning, with no production/default switch.
  - Files: `.gsd/DECISIONS.md`, `.gsd/gsd.db`
  - Verify: Decision saved through GSD.

- [x] **T03: Validate M021 closure** `est:small`
  Run fresh closure verification: artifact verifier, default Go tests, pinned lint, tagged tests, tracked binary check, Docker status, and GitNexus scope.
  - Verify: Fresh verification passes and no background processes remain.

## Files Likely Touched

- .gsd/DECISIONS.md
- .gsd/gsd.db
