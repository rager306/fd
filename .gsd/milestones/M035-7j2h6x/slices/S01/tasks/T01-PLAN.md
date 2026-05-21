---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Persist exact binary hosting contract

Update ONNX manifest source_contract and provisioning docs with exact-binary hosting contract: required size/sha, planned immutable object key template, acceptable source forms, forbidden source forms, and pre-dispatch checklist. Do not add a fake URL.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `.github/workflows/onnx-packaging.yml`

## Expected Output

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`

## Verification

JSON/doc checks confirm no fake URL and blocker remains explicit.

## Observability Impact

Makes the remaining blocker actionable for future hosted proof.
