# S01: Exact binary hosting contract

**Goal:** Define exact-binary immutable hosting contract for the current ONNX model artifact.
**Demo:** After this, the repo has a concrete exact ONNX binary hosting contract: object key pattern, checksum/size, source policy, and pre-dispatch checklist.

## Must-Haves

- Contract defines exact binary checksum/size and recommended immutable key naming.
- Contract distinguishes planned key/source policy from an actual uploaded source.
- Contract names acceptable source forms and forbidden signed/plain secret URL forms.
- Workflow dispatch preconditions remain explicit.
- No production/default ONNX promotion.

## Proof Level

- This slice proves: Docs/manifest/outcome checks; no external action.

## Integration Closure

Builds on M031 source contract, M032 verifier, M034 workflow input contract.

## Verification

- Future operators can see required artifact source, key naming, and verification gates before dispatch.

## Tasks

- [x] **T01: Persist exact binary hosting contract** `est:small`
  Update ONNX manifest source_contract and provisioning docs with exact-binary hosting contract: required size/sha, planned immutable object key template, acceptable source forms, forbidden source forms, and pre-dispatch checklist. Do not add a fake URL.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, `docs/onnx-artifacts/PROVISIONING.md`
  - Verify: JSON/doc checks confirm no fake URL and blocker remains explicit.

- [x] **T02: Record exact binary hosting outcome** `est:small`
  Update ONNX artifacts README and create outcome artifact summarizing exact binary hosting contract, workflow input readiness, and remaining blockers. Record that no upload, push, or workflow dispatch occurred.
  - Files: `docs/onnx-artifacts/README.md`, `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`
  - Verify: Outcome/README checks pass and contain no raw input text, secrets, signed URLs, or production promotion claims.

- [x] **T03: Verify exact binary contract** `est:small`
  Run slice-level verification: manifest JSON validity, docs markers, provisioning dry-run/verifier/export-contract, actionlint, and GitNexus detect.
  - Verify: All slice-level checks pass and GitNexus reports expected scope.

## Files Likely Touched

- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/PROVISIONING.md
- docs/onnx-artifacts/README.md
- benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt
