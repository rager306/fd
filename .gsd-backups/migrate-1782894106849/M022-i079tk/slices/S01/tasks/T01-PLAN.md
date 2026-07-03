---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T01: Choose ONNX Docker packaging strategy

Inspect Docker context constraints and choose the dedicated ONNX packaging strategy: root-context Dockerfile, staging script, or documented blocker.

## Inputs

- `api/Dockerfile`
- `api/.dockerignore`
- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`

## Expected Output

- `Task summary with packaging strategy`

## Verification

Strategy states context, artifact inputs, build tags, and cleanup approach.

## Observability Impact

Avoids accidentally copying ignored artifacts into the default context or git.
