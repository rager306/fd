---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T01: Inspect packaging boundary

Inspect `api/Dockerfile`, `.github/workflows/go-quality.yml`, `.gitignore`, ONNX/native manifests, and existing docs to decide the minimal packaging contract.

## Inputs

- `api/Dockerfile`
- `.github/workflows/go-quality.yml`
- `.gitignore`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Task summary with packaging design choice`

## Verification

Task summary states whether implementation is script, docs, Dockerfile change, or CI change and why.

## Observability Impact

Prevents changing default Docker/CI path accidentally.
