---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Document runtime wheel extraction

Update provisioning docs to explain ONNX Runtime `.whl`/`.zip` member extraction, manifest source_contract fields consumed, direct-file fallback, and remaining hosted proof blockers.

## Inputs

- `tools/provision_onnx_artifacts.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `docs/onnx-artifacts/PROVISIONING.md`

## Verification

Docs contain wheel extraction behavior and proof boundaries.

## Observability Impact

Makes new provisioning capability discoverable.
