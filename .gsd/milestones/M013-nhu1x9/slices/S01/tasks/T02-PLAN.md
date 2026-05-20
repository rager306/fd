---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Write native tokenizer artifact manifest

Create the tracked native tokenizer artifact manifest and setup notes. Use the local S03 prebuilt linux-amd64 `libtokenizers.a` evidence to record checksum, size, source URL, architecture, and expected local path.

## Inputs

- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`

## Expected Output

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Verification

Manifest JSON parses and local artifact checksum/size match when artifact exists.

## Observability Impact

Gives future build-tag integration a reproducible native artifact contract.
