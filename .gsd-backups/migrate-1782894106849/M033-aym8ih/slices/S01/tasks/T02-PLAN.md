---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Verify wheel extraction failure modes

Run positive/negative synthetic wheel probes in temporary repo roots: positive member extraction, missing member, oversized member, and checksum mismatch.

## Inputs

- `tools/provision_onnx_artifacts.py`

## Expected Output

- `Task summary evidence`

## Verification

All synthetic probes pass/fail as expected with sanitized output.

## Observability Impact

Proves failure modes without mutating real runtime artifacts.
