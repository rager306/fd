---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Record remediation decision

Record the GSD decision that TEI remains default and the next ONNX implementation gate is 512-token/long-text remediation plus full legal corpus quality rerun.

## Inputs

- `benchmark-results/fd-onnx-remediation-plan-m016-s03.txt`

## Expected Output

- `.gsd/DECISIONS.md`

## Verification

Decision is saved through gsd_decision_save and references the quality-first remediation path.

## Observability Impact

Keeps the project decision register aligned with M016 evidence.
