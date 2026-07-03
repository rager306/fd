---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run final guardrails

Run final milestone guardrails: JSON/doc checks, py_compile/provisioning/verifiers, actionlint, Go tests/lint, tagged tests, Docker default build, leak checks, binary hygiene, port/background checks, GitNexus detect.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt`

## Expected Output

- `Verification evidence in task summary`

## Verification

All final checks pass.

## Observability Impact

Proves final milestone state before closure.
