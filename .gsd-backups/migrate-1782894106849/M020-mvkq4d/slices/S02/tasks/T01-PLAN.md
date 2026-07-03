---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Record artifact contract decision

Record a GSD decision that the tracked ONNX manifest now represents an experimental 1024 runtime contract, not production readiness.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `.gsd/DECISIONS.md`

## Verification

Decision saved through GSD.

## Observability Impact

Prevents future agents from treating the manifest as production approval.
