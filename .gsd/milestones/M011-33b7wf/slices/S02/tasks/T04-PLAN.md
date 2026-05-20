---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Verify backend seam safety

Run S02 verification gates: Go tests, pinned lint, manifest validation against the local M010 artifact, Docker Compose config, and GitNexus detect_changes. Record whether any production behavior changed.

## Inputs

- `api/embed/onnx_manifest.go`
- `api/main.go`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Task summary with verification evidence`

## Verification

Go tests/lint/config/GitNexus pass; TEI remains default.

## Observability Impact

Confirms S02 is a validation seam only and safe for S03 to consume.
