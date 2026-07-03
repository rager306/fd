---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Persist target runtime validation contract

Update ONNX manifest and provisioning docs with target-runtime validation contract: Python helper boundary, required Go API/package gates for any new or regenerated artifact, and equivalent gate rule for any future Rust backend.

## Inputs

- `T01 findings`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`

## Expected Output

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`

## Verification

JSON/docs checks confirm Python boundary, Go gates, Rust gate rule, Redis namespace isolation, and no promotion claim.

## Observability Impact

Makes target-runtime acceptance policy machine-readable and operator-readable.
