---
id: T02
parent: S01
milestone: M011-33b7wf
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - Manifest schema version starts at `1` and is intentionally JSON for easy Go/Python validation in later slices.
  - Manifest includes failure contract text so S02/S03 can implement missing/checksum/metadata diagnostics consistently.
  - Manifest classifies current artifact as `prototype_only` and `production_default=false`.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:57:44.423Z
blocker_discovered: false
---

# T02: Wrote the tracked ONNX artifact manifest and verified it matches the local ignored 1.43GB artifact.

**Wrote the tracked ONNX artifact manifest and verified it matches the local ignored 1.43GB artifact.**

## What Happened

Created `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, a tracked manifest for the M010 ONNX FP32 dense candidate. It records model revision/source hashes, ONNX artifact path/size/SHA256, export script and dependency pins, input/output metadata, validation artifacts, failure contract, and future production gates. The manifest marks the artifact as `prototype_only`, `production_default=false`, and `git_tracked=false`. Validation parsed the JSON, confirmed runtime metadata, and checked local artifact size/hash when present. Git status shows only the manifest is untracked while `.gsd/runtime/` remains ignored.

## Verification

Manifest JSON parsed and matched local artifact size/hash. Git ignored status confirmed `.gsd/runtime/` remains ignored while only the manifest JSON is tracked candidate content.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 manifest validation for docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — manifest_ok=true; artifact_present=True | 0ms |
| 2 | `git status --short --ignored .gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — manifest untracked; .gsd/runtime ignored | 0ms |

## Deviations

None. The manifest is tracked JSON only; the ONNX binary remains ignored.

## Known Issues

Manifest currently references local `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`; future production work must replace local-only storage with external artifact storage or a documented download/export procedure.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
