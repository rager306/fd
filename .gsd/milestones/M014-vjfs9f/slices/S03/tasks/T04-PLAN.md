---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Verify tagged ONNX artifact and cleanup

Verify tagged ONNX artifact for required sections, snapshot v3 metadata, raw text hygiene, correctness gate reference, and cleanup the tagged server.

## Inputs

- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`

## Expected Output

- `Task summary with artifact and cleanup evidence`

## Verification

Parser/leak checks pass, health/process cleanup confirmed, GitNexus detect_changes pass.

## Observability Impact

Ensures ONNX evidence is safe to compare and no benchmark server is left running.
