---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Implement export contract verifier

Design and add a local verifier script for the ONNX export contract. It should read the tracked ONNX manifest, M010 source-provenance, and M010 export-metadata; validate artifact size/sha, source file checksums, package pins, model revision, dynamic axes/outputs, and bounded claim semantics; and print structured JSON with safe paths.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `.gsd/runtime/onnx/m010-s03/source-provenance.json`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`

## Expected Output

- `tools/verify_onnx_export_contract.py`

## Verification

py_compile and positive verifier run pass.

## Observability Impact

Adds deterministic local verifier output for ONNX export contract state.
