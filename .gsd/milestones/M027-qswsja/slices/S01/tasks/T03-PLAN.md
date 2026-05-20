---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify startup preflight diagnostics

Run S01 verification: targeted tests, default Go tests, lint, tagged checks, default Docker, binary hygiene, cleanup, GitNexus scope.

## Inputs

- `api/main.go`
- `api/embed/onnx_manifest.go`

## Expected Output

- `Task summary`

## Verification

All S01 checks pass.

## Observability Impact

Confirms no default regression.
