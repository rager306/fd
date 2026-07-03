---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Cleanup benchmark runtime

Stop the tagged ONNX benchmark service and verify no stale runtime remains.

## Inputs

- None specified.

## Expected Output

- `Clean runtime state`

## Verification

Background process list shows no benchmark service.

## Observability Impact

Prevents stale benchmark server from contaminating future runs.
