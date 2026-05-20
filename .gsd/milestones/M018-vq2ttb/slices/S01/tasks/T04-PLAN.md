---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Cleanup runtime

Stop tagged ONNX service and verify no stale benchmark runtime remains.

## Inputs

- None specified.

## Expected Output

- `Clean runtime state`

## Verification

Background process list shows the tagged ONNX service stopped.

## Observability Impact

Avoids stale runtime contaminating later runs.
