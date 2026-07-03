# S01: Workflow runtime source alignment

**Goal:** Align manual ONNX packaging workflow with M033 runtime wheel provisioning behavior.
**Demo:** After this, the workflow no longer blocks runtime wheel provisioning solely because `onnx_runtime_sha256` is omitted when manifest metadata supplies it.

## Must-Haves

- Workflow validation no longer requires `onnx_runtime_sha256` when `onnx_runtime_source_url` is set.
- Provisioning step passes `--onnx-runtime-sha256` only when non-empty.
- Workflow comments/descriptions clarify manifest-derived sha fallback.
- actionlint passes.

## Proof Level

- This slice proves: Workflow diff, actionlint, provisioning checks.

## Integration Closure

Connects M033 provisioning helper behavior to the manual workflow skeleton.

## Verification

- Workflow validation messages and docs reflect manifest-derived runtime checksum behavior.

## Tasks

- [x] **T01: Align workflow runtime sha handling** `est:small`
  Update `.github/workflows/onnx-packaging.yml` input descriptions, validation step, and provisioning argument assembly so runtime sha is optional when manifest metadata supplies it, while preserving checksum verification by provisioning.
  - Files: `.github/workflows/onnx-packaging.yml`
  - Verify: actionlint and text checks pass.

- [x] **T02: Verify workflow alignment** `est:small`
  Verify workflow behavior locally through actionlint and by checking provisioning still exposes manifest-derived runtime sha in dry-run.
  - Verify: actionlint, dry-run expected runtime sha, py_compile pass.

## Files Likely Touched

- .github/workflows/onnx-packaging.yml
