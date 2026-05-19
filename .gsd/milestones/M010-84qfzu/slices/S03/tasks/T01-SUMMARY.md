---
id: T01
parent: S03
milestone: M010-84qfzu
key_files:
  - .gsd/runtime/onnx/m010-s03/source-provenance.json
key_decisions:
  - Use `.gsd/runtime/onnx/m010-s03/` as local untracked ONNX artifact workspace.
  - Use `tei-models/deepvk--USER-bge-m3` as the exact local source model path for export attempts.
  - Record source provenance in `.gsd/runtime/onnx/m010-s03/source-provenance.json` rather than a tracked file because it describes local runtime artifacts.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:39:56.983Z
blocker_discovered: false
---

# T01: Prepared the ignored S03 ONNX workspace and captured exact local model provenance/hashes before export.

**Prepared the ignored S03 ONNX workspace and captured exact local model provenance/hashes before export.**

## What Happened

Prepared the local ignored ONNX workspace at `.gsd/runtime/onnx/m010-s03/` and wrote `source-provenance.json`. The provenance records the exact local model path `tei-models/deepvk--USER-bge-m3`, model revision `0cc6cfe48e260fb0474c753087a69369e88709ae`, hashes and sizes for safetensors/tokenizer/config/SentenceTransformer files, available disk space (~84.35 GiB free), current git commit/branch/status, and `production_runtime_changed=false`. Git status confirms no ONNX or model artifacts are staged; only expected tracked GSD/benchmark/tool files are uncommitted.

## Verification

Created `.gsd/runtime/onnx/m010-s03/source-provenance.json`, read it back, and confirmed git status does not show ONNX/model artifacts because the workspace is under ignored `.gsd/runtime/`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 provenance-capture snippet writing .gsd/runtime/onnx/m010-s03/source-provenance.json && git status --short` | 0 | ✅ pass | 0ms |
| 2 | `read .gsd/runtime/onnx/m010-s03/source-provenance.json` | 0 | ✅ pass | 0ms |

## Deviations

None. Workspace is under ignored `.gsd/runtime/` as planned.

## Known Issues

`.gsd/runtime/` is intentionally ignored, so source provenance must be summarized in tracked GSD task/slice artifacts if needed for durable review. Current git status already contains expected M010 artifacts and new tools directory from S02.

## Files Created/Modified

- `.gsd/runtime/onnx/m010-s03/source-provenance.json`
