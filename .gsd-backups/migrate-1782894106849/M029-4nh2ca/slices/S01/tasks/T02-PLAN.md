---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Harden archive member extraction

Add archive member safeguards for native tokenizer extraction: require regular file, validate declared member size against expected size/hard cap before copying, and enforce copy byte limit.

## Inputs

- `tools/provision_onnx_artifacts.py`

## Expected Output

- `Updated archive materialization safety`

## Verification

Local probes cover oversized/non-regular archive member rejection.

## Observability Impact

Prevents runner exhaustion before checksum failure.
