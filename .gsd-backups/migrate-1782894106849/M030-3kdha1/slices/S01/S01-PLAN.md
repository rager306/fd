# S01: Artifact path policy and safe diagnostics

**Goal:** Implement path-root policy and safe path diagnostics for ONNX artifact tooling/startup surfaces.
**Demo:** After this, manifest/provisioning/verifier artifact paths are root-constrained and diagnostics avoid leaking absolute host paths by default.

## Must-Haves

- Existing approved artifact roots continue to work.
- Repo-external manifest artifact paths are rejected in Go startup validation and Python verifier/provisioning where applicable.
- Error/report display uses safe repo-relative or basename/logical artifact path.
- Local probes cover rejection and safe display.
- Guardrails pass.

## Proof Level

- This slice proves: Unit tests/probes plus provisioning/verifier guardrails and Go tests.

## Integration Closure

Closes remaining M028 low security findings before artifact source selection or hosted workflow proof.

## Verification

- Errors remain actionable via artifact labels and repo-relative/safe paths while reducing host path disclosure.

## Tasks

- [x] **T01: Harden Go manifest artifact path policy** `est:medium`
  Implement Go ONNX manifest path policy and safe path display: restrict `artifact.local_path` to approved repo-relative roots; adjust errors to avoid absolute paths where possible; add tests for allowed and rejected paths.
  - Files: `api/embed/onnx_manifest.go`, `api/embed/onnx_manifest_test.go`
  - Verify: Targeted manifest tests pass.

- [x] **T02: Harden Python artifact path policy and diagnostics** `est:medium`
  Implement Python provisioning/verifier path-root policy and safe display helpers; update build script missing diagnostics if needed. Add deterministic local probes for allowed/rejected paths and sanitized output.
  - Files: `tools/provision_onnx_artifacts.py`, `tools/verify_onnx_artifacts.py`, `tools/build_onnx_image.sh`
  - Verify: Python probes and provisioning/verifier checks pass.

- [x] **T03: Verify path security remediation** `est:medium`
  Run S01 guardrails: targeted tests/probes, Python compile, provisioning/verifier behavior, default Go tests/lint, actionlint, Docker build, binary hygiene, cleanup, GitNexus scope.
  - Verify: All S01 checks pass.

## Files Likely Touched

- api/embed/onnx_manifest.go
- api/embed/onnx_manifest_test.go
- tools/provision_onnx_artifacts.py
- tools/verify_onnx_artifacts.py
- tools/build_onnx_image.sh
