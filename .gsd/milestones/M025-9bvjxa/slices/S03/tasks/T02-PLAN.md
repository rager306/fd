---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Document manual hosted CI workflow

Update ONNX README/PROVISIONING docs to mention the manual workflow, required inputs, and non-secret URL guidance.

## Inputs

- `.github/workflows/onnx-packaging.yml`

## Expected Output

- `Updated docs`

## Verification

Docs reference workflow and warn that signed URLs must use secrets/masked values.

## Observability Impact

Future agents know how and when to run the workflow safely.
