# S02: Reproducibility strategy documentation and closure

**Goal:** Persist reproducibility/source strategy docs and close M032 locally.
**Demo:** After this, docs/manifests/outcome describe the exact ONNX model source options and remaining blocker without overclaiming.

## Must-Haves

- Provisioning docs explain exact binary hosting vs reproducible export strategy.
- Outcome artifact records verifier evidence and blockers.
- Decision recorded.
- No production/default switch, no external state changes, no raw text/secrets/signed URLs.

## Proof Level

- This slice proves: Docs/outcome checks plus project guardrails, GitNexus detect, commit/reindex.

## Integration Closure

Updates provisioning/source-contract docs and outcome artifacts to reference the new verifier and next gate options.

## Verification

- Future hosted proof has a sharper checklist for exact-binary hosting vs reproducible export.

## Tasks

- [x] **T01: Document verifier in source contract** `est:small`
  Update provisioning docs and ONNX manifest source_contract to reference `tools/verify_onnx_export_contract.py`, clarify existing-artifact verification vs regenerated export, and list next gate options.
  - Files: `docs/onnx-artifacts/PROVISIONING.md`, `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
  - Verify: Docs/manifests parse and mention proof boundary.

- [x] **T02: Record verifier outcome and decision** `est:small`
  Create M032 outcome artifact and record a decision about exact-binary hosting vs reproducible export strategy.
  - Files: `benchmark-results/fd-onnx-export-contract-verifier-m032-s02.txt`, `.gsd/DECISIONS.md`
  - Verify: Outcome/decision safety checks pass.

- [x] **T03: Verify and close milestone** `est:medium`
  Run final project guardrails, validate and complete M032, checkpoint DB, commit locally, run GitNexus reindex/detect, and report state.
  - Verify: Verifier positive/negative checks, py_compile, Go checks, actionlint, docs leak checks, tracked binary hygiene, GitNexus detect, commit, reindex.

## Files Likely Touched

- docs/onnx-artifacts/PROVISIONING.md
- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- benchmark-results/fd-onnx-export-contract-verifier-m032-s02.txt
- .gsd/DECISIONS.md
