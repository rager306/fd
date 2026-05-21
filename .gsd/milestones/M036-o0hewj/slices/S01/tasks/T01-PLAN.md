---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T01: Inspect export evidence boundaries

Inspect export script/verifier/provenance metadata and identify pinned inputs and current claim boundaries for a reproducible-export workflow contract.

## Inputs

- `tools/export_user_bge_m3_dense_onnx.py`
- `tools/verify_onnx_export_contract.py`
- `.gsd/runtime/onnx/m010-s03/source-provenance.json`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`

## Expected Output

- `Task summary`

## Verification

Summarize pinned inputs and current non-regenerated verifier boundary.

## Observability Impact

Avoids inventing workflow requirements not grounded in existing exporter/verifier evidence.
