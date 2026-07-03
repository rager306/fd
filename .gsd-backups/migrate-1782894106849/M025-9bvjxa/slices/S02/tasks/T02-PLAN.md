---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Record operational rollout decision

Record GSD decision for ONNX operational rollout: opt-in staged rollout only, TEI default preserved, rollback by env/backend switch, production switch blocked until diagnostics are implemented and tested.

## Inputs

- `docs/onnx-artifacts/OPERATIONS.md`

## Expected Output

- `Decision record`

## Verification

Decision saved.

## Observability Impact

Prevents operational contract from being misread as production approval.
