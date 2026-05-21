---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Harden Python artifact path policy and diagnostics

Implement Python provisioning/verifier path-root policy and safe display helpers; update build script missing diagnostics if needed. Add deterministic local probes for allowed/rejected paths and sanitized output.

## Inputs

- `.gsd/milestones/M028-y63tog/slices/S01/S01-RESEARCH.md`
- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`
- `tools/build_onnx_image.sh`

## Expected Output

- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`
- `tools/build_onnx_image.sh`

## Verification

Python probes and provisioning/verifier checks pass.

## Observability Impact

Tool errors avoid absolute host path disclosure by default.
