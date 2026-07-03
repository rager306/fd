---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Persist reproducible export contract

Update manifest and provisioning docs with planned reproducible-export workflow contract: pinned inputs, expected command/workflow shape, acceptance gates, and explicit non-proof status.

## Inputs

- `Task T01 findings`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`

## Expected Output

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`

## Verification

JSON/docs checks confirm planned_not_proven status and required gates.

## Observability Impact

Records no-upload path next to exact-binary source contract.
