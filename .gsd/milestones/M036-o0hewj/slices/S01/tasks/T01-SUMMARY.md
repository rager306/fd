---
id: T01
parent: S01
milestone: M036-o0hewj
key_files:
  - tools/export_user_bge_m3_dense_onnx.py
  - tools/verify_onnx_export_contract.py
  - .gsd/runtime/onnx/m010-s03/source-provenance.json
  - .gsd/runtime/onnx/m010-s03/export-metadata.json
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T09:33:10.045Z
blocker_discovered: false
---

# T01: Mapped the current ONNX export evidence boundary for the reproducible-export contract.

**Mapped the current ONNX export evidence boundary for the reproducible-export contract.**

## What Happened

Inspected the exporter, local verifier, source provenance, and export metadata. Existing evidence pins model revision `0cc6cfe48e260fb0474c753087a69369e88709ae`, source file checksums, Python `3.13.12`, `torch==2.12.0`, `transformers==4.51.3`, `onnx==1.21.0`, `onnxruntime==1.26.0`, `safetensors==0.7.0`, opset 17, export sequence length 128, CLS pooling, L2 normalization, dynamic axes, CPUExecutionProvider, and 1024-dimensional output. The verifier confirms the existing ignored artifact against this metadata but explicitly does not regenerate the ONNX binary.

## Verification

Evidence files read and GitNexus impact for verifier file was LOW with no affected processes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gitnexus_impact(target=verify_onnx_export_contract.py, direction=upstream)` | 0 | ✅ pass — LOW risk, no impacted processes | 0ms |
| 2 | `read exporter/verifier/provenance/export-metadata` | 0 | ✅ pass — pinned inputs and claim boundary identified | 0ms |

## Deviations

None.

## Known Issues

Current verifier claim_scope remains `existing_artifact_contract_verification_not_regenerated_export`; no regenerated export proof exists.

## Files Created/Modified

- `tools/export_user_bge_m3_dense_onnx.py`
- `tools/verify_onnx_export_contract.py`
- `.gsd/runtime/onnx/m010-s03/source-provenance.json`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`
