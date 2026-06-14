---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: ONNX cold/warm perf benchmark + async combo

tools/verify_fd_onnx_perf.sh: build fd-onnx binary (go build -tags onnx), start with FD_BACKEND=onnx + ONNX_ARTIFACT_MANIFEST=docs/onnx-artifacts/user-bge-m3-dense-fp32.json + ONNX_RUNTIME_LIBRARY=libonnxruntime.so + ONNX_TOKENIZER_PATH=..., measure cold/warm path × batch sizes 1/8/32/64/128, with FD_ASYNC_CHUNKS=true AND false. Output benchmark-results/fd-v2-onnx-perf-m042.md с comparison table (TEI vs ONNX × async on/off).

## Inputs

- None specified.

## Expected Output

- `tools/verify_fd_onnx_perf.sh`
- `benchmark-results/fd-v2-onnx-perf-m042.md`

## Verification

ONNX mode exit 0. Artifact содержит: cold start latency (first request), warm p95 by batch, async on/off comparison. ONNX cold ≤500ms, warm ≤10ms. TEI baseline (from S02) shown for comparison.
