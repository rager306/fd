# S01: Reproducible export contract

**Goal:** Document reproducible-export contract as the no-upload alternative path.
**Demo:** After this, the repo has a reproducible-export workflow contract that names pinned inputs, expected outputs, and gates before regenerated ONNX can replace exact-binary hosting.

## Must-Haves

- Contract records pinned model revision/toolchain inputs.
- Contract states expected output identity and acceptance criteria.
- Contract requires rerunning legal quality, performance, packaging, and hosted proof gates.
- Contract explicitly says current state is planned, not proven.

## Proof Level

- This slice proves: Manifest/docs/outcome checks; no regeneration or workflow dispatch.

## Integration Closure

Builds on M032 local export verifier and M035 exact-binary hosting contract.

## Verification

- Future agents can distinguish local existing-artifact verification from regenerated-export proof.

## Tasks

- [x] **T01: Inspect export evidence boundaries** `est:small`
  Inspect export script/verifier/provenance metadata and identify pinned inputs and current claim boundaries for a reproducible-export workflow contract.
  - Files: `tools/export_user_bge_m3_dense_onnx.py`, `tools/verify_onnx_export_contract.py`, `.gsd/runtime/onnx/m010-s03/source-provenance.json`, `.gsd/runtime/onnx/m010-s03/export-metadata.json`
  - Verify: Summarize pinned inputs and current non-regenerated verifier boundary.

- [x] **T02: Persist reproducible export contract** `est:small`
  Update manifest and provisioning docs with planned reproducible-export workflow contract: pinned inputs, expected command/workflow shape, acceptance gates, and explicit non-proof status.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, `docs/onnx-artifacts/PROVISIONING.md`
  - Verify: JSON/docs checks confirm planned_not_proven status and required gates.

- [x] **T03: Record reproducible export outcome** `est:small`
  Update README and write outcome artifact for reproducible-export workflow contract; verify no overclaim/leaks/signed URLs.
  - Files: `docs/onnx-artifacts/README.md`, `benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt`
  - Verify: Outcome/README checks pass; no proof overclaim.

- [x] **T04: Verify reproducible export contract** `est:small`
  Run S01 verification: manifest JSON, contract marker checks, provisioning/export verifier, actionlint, GitNexus detect.
  - Verify: S01 checks pass and GitNexus scope is low risk.

## Files Likely Touched

- tools/export_user_bge_m3_dense_onnx.py
- tools/verify_onnx_export_contract.py
- .gsd/runtime/onnx/m010-s03/source-provenance.json
- .gsd/runtime/onnx/m010-s03/export-metadata.json
- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/PROVISIONING.md
- docs/onnx-artifacts/README.md
- benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt
