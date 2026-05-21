---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run Go ONNX API smoke

Start local Go ONNX API with isolated cache namespace on port 18000, verify `/health` runtime metadata and `/v1/embeddings` 1024-dimensional response, then stop the server.

## Inputs

- `api/main.go`
- `api/handlers/health.go`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`

## Verification

Health and embeddings smoke pass; server stopped; outcome has no raw input text/secrets/signed URLs.

## Observability Impact

Captures target-runtime API proof and safe metadata evidence.
