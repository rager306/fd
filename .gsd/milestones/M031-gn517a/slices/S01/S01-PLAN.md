# S01: Artifact source contract research

**Goal:** Research and define immutable source contracts for ONNX model, native tokenizer, tokenizer JSON, and ONNX Runtime.
**Demo:** After this, every required ONNX artifact has a source contract decision and explicit blocker/candidate status.

## Must-Haves

- Artifact inventory covers ONNX model, native tokenizer, tokenizer JSON, ONNX Runtime.
- Source status is explicit for each artifact.
- No fake/default URLs are introduced.
- Mutable/latest/signed URLs are not blessed as immutable.
- Outcome contains exact checksums/sizes for future verification.

## Proof Level

- This slice proves: Codebase/GitNexus inventory, manifest/provenance review, web/source research where needed, source-status artifact.

## Integration Closure

Builds on M025 provisioning contract and M029/M030 security hardening by defining safe source inputs for future hosted workflow proof.

## Verification

- Clarifies what future operators must provide and verify.

## Tasks

- [x] **T01: Inventory artifact source requirements** `est:small`
  Inventory current tracked manifests, local provenance, provisioning docs, and workflow requirements into an artifact source matrix with checksums/sizes and current blockers.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`, `docs/onnx-artifacts/PROVISIONING.md`, `.github/workflows/onnx-packaging.yml`
  - Verify: Inventory artifact contains all four required artifacts and exact local checksums.

- [x] **T02: Assess immutable source candidates** `est:medium`
  Research immutable source candidates for native tokenizer and ONNX Runtime and assess source strategy for ONNX model/tokenizer JSON without downloading large artifacts. Record whether each candidate is immutable, policy-compliant, and ready or blocked.
  - Files: `.gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md`
  - Verify: Research notes distinguish immutable/candidate/blocked and cite sources.

- [x] **T03: Verify source contract research** `est:small`
  Verify S01 research artifact: marker/leak checks, no fake source URLs, no raw text/secrets, and clear status for every artifact.
  - Verify: Marker/leak/source-status checks pass.

## Files Likely Touched

- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
- docs/onnx-artifacts/PROVISIONING.md
- .github/workflows/onnx-packaging.yml
- .gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md
