---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Preflight confirmed tagged ONNX artifacts, checksums, and tagged tests are ready for the benchmark.

Validate local ONNX/native/ORT artifact availability, confirm tagged test still passes, and record intended env vars including isolated cache namespace.

## Inputs

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`

## Expected Output

- `Task summary with preflight evidence`

## Verification

Artifact checks, tagged Go test, and namespace/env summary pass.

## Observability Impact

Ensures tagged benchmark will use verified artifacts and isolated cache namespace.
