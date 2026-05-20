# S03: Full hosted ONNX CI skeleton

**Goal:** Add a safe manual hosted ONNX CI skeleton after artifact provisioning, without making it a required fake-success job.
**Demo:** After this, hosted CI has a safe manual ONNX packaging workflow skeleton that runs only when artifact inputs are supplied, plus normal default CI remains unaffected.

## Must-Haves

- Manual workflow exists or docs explain why not.
- Workflow is not triggered by push/PR.
- Workflow requires explicit artifact sources before full ONNX packaging.
- Workflow uses provisioning helper, verifier, tagged tests, and build script.
- actionlint and default guardrails pass.

## Proof Level

- This slice proves: Workflow file, actionlint, dry-run/provision checks, default guardrails.

## Integration Closure

Connects provisioning helper and ONNX Docker build into a future hosted CI path that runs only with explicit artifact sources.

## Verification

- Manual workflow separates provisioning, verification, tagged tests, and image build phases in CI logs.

## Tasks

- [x] **T01: Add manual ONNX packaging workflow skeleton** `est:medium`
  Add `.github/workflows/onnx-packaging.yml` as a manual `workflow_dispatch` skeleton with explicit artifact inputs, provisioning, strict verification, tagged tests, and Docker image build.
  - Files: `.github/workflows/onnx-packaging.yml`
  - Verify: Workflow has workflow_dispatch only and actionlint passes.

- [x] **T02: Document manual hosted CI workflow** `est:small`
  Update ONNX README/PROVISIONING docs to mention the manual workflow, required inputs, and non-secret URL guidance.
  - Files: `docs/onnx-artifacts/README.md`, `docs/onnx-artifacts/PROVISIONING.md`
  - Verify: Docs reference workflow and warn that signed URLs must use secrets/masked values.

- [x] **T03: Verify hosted CI skeleton and milestone guardrails** `est:medium`
  Run S03 and milestone closure checks: actionlint, provisioning dry-run, verifier, default tests/lint, tagged tests, Docker default, binary hygiene, cleanup, GitNexus.
  - Verify: All closure checks pass.

## Files Likely Touched

- .github/workflows/onnx-packaging.yml
- docs/onnx-artifacts/README.md
- docs/onnx-artifacts/PROVISIONING.md
