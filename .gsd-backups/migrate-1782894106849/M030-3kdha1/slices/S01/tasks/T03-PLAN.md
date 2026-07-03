---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify path security remediation

Run S01 guardrails: targeted tests/probes, Python compile, provisioning/verifier behavior, default Go tests/lint, actionlint, Docker build, binary hygiene, cleanup, GitNexus scope.

## Inputs

- `api/embed/onnx_manifest.go`
- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`

## Expected Output

- `.gsd/milestones/M030-3kdha1/slices/S01/tasks/T03-SUMMARY.md`

## Verification

All S01 checks pass.

## Observability Impact

Confirms remediation scope and no default regression.
