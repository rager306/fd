---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare packaged ONNX legal gate environment

Prepare packaged legal gate environment: verify TEI baseline health, build or reuse the M022 ONNX image, start packaged ONNX on port 18000 with an isolated cache namespace, and smoke `/health` plus non-legal embedding dimensions.

## Inputs

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `Task summary with endpoint health evidence`

## Verification

TEI `/health` and ONNX `/health` pass; ONNX smoke returns 1024 dimensions; port/process state is known.

## Observability Impact

Captures runtime labels, image tag, endpoint URLs, and isolated cache namespace before quality run.
