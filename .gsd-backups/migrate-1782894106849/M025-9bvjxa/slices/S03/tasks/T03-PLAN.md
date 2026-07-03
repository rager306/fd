---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify hosted CI skeleton and milestone guardrails

Run S03 and milestone closure checks: actionlint, provisioning dry-run, verifier, default tests/lint, tagged tests, Docker default, binary hygiene, cleanup, GitNexus.

## Inputs

- `.github/workflows/onnx-packaging.yml`

## Expected Output

- `Task summary with closure evidence`

## Verification

All closure checks pass.

## Observability Impact

Proves skeleton is syntactically valid and does not regress default CI.
