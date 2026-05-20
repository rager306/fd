---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Add CI-safe ONNX artifact contract checks

Add a CI-safe artifact contract check to the existing Go Quality workflow: run the verifier in allow-missing mode and fail if ONNX/native/runtime binaries are tracked. Include relevant docs/tools paths in workflow triggers.

## Inputs

- `.github/workflows/go-quality.yml`
- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`

## Expected Output

- `Updated .github/workflows/go-quality.yml`

## Verification

Workflow YAML parses; verifier allow-missing and binary hygiene checks pass locally.

## Observability Impact

Hosted CI can validate contract metadata without requiring unavailable binary artifacts.
