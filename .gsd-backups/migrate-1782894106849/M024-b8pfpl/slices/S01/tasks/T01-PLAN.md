---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare packaged ONNX benchmark target

Prepare the packaged ONNX benchmark target: ensure TEI/default stack is not modified, start `fd-api:onnx1024-m022-final` on port 18000 with cache namespace `m024-onnx-docker-benchmark`, verify `/health` and embedding dimensions, and confirm restart command strategy.

## Inputs

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`

## Expected Output

- `Task summary with endpoint and restart strategy`

## Verification

Health and smoke embedding pass; restart command targets packaged ONNX container.

## Observability Impact

Captures target image, endpoint, namespace, and restart command before measurement.
