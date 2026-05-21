---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run packaged ONNX smoke

Run packaged ONNX container on port 18000 with isolated namespace, verify `/health` and `/v1/embeddings`, write smoke artifact, then stop container.

## Inputs

- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`

## Expected Output

- `benchmark-results/fd-onnx-docker-smoke-m039-s01.txt`

## Verification

Packaged health/embedding smoke passes; no raw text/secrets; container stopped; port clean.

## Observability Impact

Captures packaged runtime metadata and embedding proof.
