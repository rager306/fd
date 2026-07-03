# S02: CI artifact provisioning boundary

**Goal:** Make CI boundary honest by adding artifact-free checks that can run today and documenting why full ONNX image CI needs external artifact provisioning.
**Demo:** After this, CI limitations and next automation steps are explicit, with no fake CI support claims.

## Must-Haves

- CI artifact-free boundary is explicit.
- Safe CI checks cover verifier allow-missing mode and binary hygiene.
- Full ONNX image CI is documented as blocked on artifact store/provisioning.
- Default Go CI remains unaffected.
- Fresh verification passes.

## Proof Level

- This slice proves: Workflow syntax/local command checks and decision record.

## Integration Closure

Builds on S01 local image proof and defines what can/cannot run in hosted CI without committing binaries.

## Verification

- Adds CI-visible checks for artifact contract shape and binary hygiene.

## Tasks

- [x] **T01: Add CI-safe ONNX artifact contract checks** `est:small`
  Add a CI-safe artifact contract check to the existing Go Quality workflow: run the verifier in allow-missing mode and fail if ONNX/native/runtime binaries are tracked. Include relevant docs/tools paths in workflow triggers.
  - Files: `.github/workflows/go-quality.yml`
  - Verify: Workflow YAML parses; verifier allow-missing and binary hygiene checks pass locally.

- [x] **T02: Document full ONNX image CI blocker** `est:small`
  Record the CI boundary decision and update docs: full ONNX image CI requires an external artifact store/cache to provide ONNX model, libtokenizers.a, tokenizer JSON, and ONNX Runtime shared library before running tagged Docker build.
  - Files: `docs/onnx-artifacts/README.md`, `.gsd/DECISIONS.md`
  - Verify: Decision saved; README states CI-safe check versus full image provisioning requirements.

- [x] **T03: Verify CI boundary and milestone closure** `est:medium`
  Run closure verification for S02 and M022: workflow check commands, default tests/lint, tagged tests, Docker builds, verifier, binary hygiene, background/port cleanup, and GitNexus scope.
  - Verify: All closure commands pass or produce a concrete blocker.

## Files Likely Touched

- .github/workflows/go-quality.yml
- docs/onnx-artifacts/README.md
- .gsd/DECISIONS.md
