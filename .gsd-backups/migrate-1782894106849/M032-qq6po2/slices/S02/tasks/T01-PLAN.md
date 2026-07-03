---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Document verifier in source contract

Update provisioning docs and ONNX manifest source_contract to reference `tools/verify_onnx_export_contract.py`, clarify existing-artifact verification vs regenerated export, and list next gate options.

## Inputs

- `.gsd/milestones/M032-qq6po2/slices/S01/S01-RESEARCH.md`
- `tools/verify_onnx_export_contract.py`

## Expected Output

- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Verification

Docs/manifests parse and mention proof boundary.

## Observability Impact

Makes verifier discoverable in durable docs/manifests.
