---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Record packaging contract decision

Record a GSD decision that M021 added artifact verification contract and the next gate is Docker/CI provisioning, with no production/default switch.

## Inputs

- `docs/onnx-artifacts/README.md`
- `tools/verify_onnx_artifacts.py`

## Expected Output

- `.gsd/DECISIONS.md`

## Verification

Decision saved through GSD.

## Observability Impact

Prevents future agents from treating contract docs as packaged runtime readiness.
