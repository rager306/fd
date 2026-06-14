---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: Regression suite: M041 + M042 S02 acceptance в ONNX mode

tools/verify_fd_v2_contract.py запустить с FD_BACKEND=onnx + FD_ASYNC_CHUNKS=true. Все 45 M041 test cases + все M042 S02 async tests pass в ONNX mode. Особенно: encoding_format=base64 (response uses ONNX []float32), dimensions=512 (Matryoshka truncation), validation envelope (413, 400 etc), error propagation (ONNX errors → OpenAI envelope). Любой failure — bug в ONNX impl или в adapter.

## Inputs

- None specified.

## Expected Output

- `tools/verify_fd_v2_contract.py (run in ONNX mode)`
- `tests/integration/fd_v2_onnx_test.go`

## Verification

go test ./tests/integration/... -run TestFdV2ONNXMode: все M041 + M042 S02 acceptance pass в ONNX mode. 0 regressions.
