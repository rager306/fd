# S01: Artifact provisioning contract

**Goal:** Create the artifact provisioning contract for ONNX 1024 and native tokenizer artifacts without tracking binaries or changing default runtime.
**Demo:** After this, artifact staging/checksum validation is documented and/or scripted for ONNX 1024 and native tokenizer without tracking binaries.

## Must-Haves

- Current Docker/CI boundaries are inspected.
- Artifact provisioning contract is documented or scripted.
- Checksums come from tracked manifests.
- No binary artifacts are tracked.
- Default build remains TEI/default and independent from ONNX artifacts.

## Proof Level

- This slice proves: Docs/scripts plus manifest validation and binary hygiene.

## Integration Closure

Defines the artifact boundary needed before Docker image or CI automation can consume ONNX 1024 safely.

## Verification

- Adds explicit checksum verification/failure guidance for artifact staging.

## Tasks

- [x] **T01: Inspect packaging boundary** `est:small`
  Inspect `api/Dockerfile`, `.github/workflows/go-quality.yml`, `.gitignore`, ONNX/native manifests, and existing docs to decide the minimal packaging contract.
  - Files: `api/Dockerfile`, `.github/workflows/go-quality.yml`, `.gitignore`, `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
  - Verify: Task summary states whether implementation is script, docs, Dockerfile change, or CI change and why.

- [x] **T02: Add artifact verification contract** `est:medium`
  Add a local artifact verification/staging contract that checks ONNX and native tokenizer manifests against local ignored files and emits actionable errors without printing secrets or raw text.
  - Files: `tools/verify_onnx_artifacts.py`, `docs/onnx-artifacts/README.md`
  - Verify: `python3 tools/verify_onnx_artifacts.py --onnx-manifest docs/onnx-artifacts/user-bge-m3-dense-fp32.json --native-tokenizer-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` passes locally.

- [x] **T03: Validate artifact contract** `est:small`
  Validate the contract, README, manifests, and tracked binary hygiene.
  - Verify: Script compile/run, manifest JSON validation, and tracked binary checks pass.

## Files Likely Touched

- api/Dockerfile
- .github/workflows/go-quality.yml
- .gitignore
- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
- tools/verify_onnx_artifacts.py
- docs/onnx-artifacts/README.md
