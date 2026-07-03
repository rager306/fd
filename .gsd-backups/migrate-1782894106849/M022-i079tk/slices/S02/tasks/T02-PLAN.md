---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Document full ONNX image CI blocker

Record the CI boundary decision and update docs: full ONNX image CI requires an external artifact store/cache to provide ONNX model, libtokenizers.a, tokenizer JSON, and ONNX Runtime shared library before running tagged Docker build.

## Inputs

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`

## Expected Output

- `Decision record and README update`

## Verification

Decision saved; README states CI-safe check versus full image provisioning requirements.

## Observability Impact

Prevents fake CI readiness claims and gives the next agent concrete provisioning requirements.
