---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Write tracked ONNX artifact manifest

Create a tracked ONNX artifact manifest for the M010 FP32 dense candidate that records local expected path, size, sha256, source model revision/hash, output metadata, dependency pin, and production status. The manifest must not contain raw probe texts or secrets and must not include the ONNX binary.

## Inputs

- `.gsd/runtime/onnx/m010-s03/export-metadata.json`
- `benchmark-results/fd-onnx-fp32-m010-s03.txt`

## Expected Output

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Verification

Manifest JSON parses; path/size/sha256 match local artifact when present; git status shows ONNX binary remains ignored/untracked.

## Observability Impact

Provides a tracked checksum contract S02/S03 can consume for validation and error reporting.
