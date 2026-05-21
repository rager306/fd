---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run live Go ONNX embedder test

Run live tagged Go ONNX embedder test against the current local artifact and native HF tokenizer path.

## Inputs

- `api/embed/onnx_test.go`
- `api/embed/hf_tokenizer_native_test.go`

## Expected Output

- `Task summary`

## Verification

`TestONNXEmbedderLiveLocalArtifact` passes with ONNX/HF tokenizer tags and configured runtime library.

## Observability Impact

Proves Go embedder can load and infer current artifact outside Python runtime.
