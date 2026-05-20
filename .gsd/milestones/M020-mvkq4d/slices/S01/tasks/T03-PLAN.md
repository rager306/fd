---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Validate contract metadata

Verify metadata contract, no binary tracking, and artifact hygiene for S01 outputs.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Task summary with validation evidence`

## Verification

JSON/field check passes and tracked binary check reports zero ONNX/native binaries.

## Observability Impact

Confirms contract is machine-readable and safe to commit.
