---
id: T03
parent: S01
milestone: M021-4t2wpt
key_files:
  - tools/verify_onnx_artifacts.py
  - docs/onnx-artifacts/README.md
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
key_decisions:
  - S01 contract validation passed and GitNexus sees no changed indexed symbols yet because the additions are new docs/script files outside current graph impact.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:15:47.943Z
blocker_discovered: false
---

# T03: Validated the ONNX/native artifact provisioning contract and binary hygiene.

**Validated the ONNX/native artifact provisioning contract and binary hygiene.**

## What Happened

Validated the artifact contract after adding the verifier and README. The verifier compiles, both manifests parse, local ignored ONNX/native tokenizer artifacts match manifest size and checksum, no `.onnx` or `libtokenizers.a` files are tracked, and GitNexus reports no changed indexed symbols.

## Verification

Artifact contract validation, tracked binary check, and GitNexus scope check passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/verify_onnx_artifacts.py && json.tool manifests && verify_onnx_artifacts.py` | 0 | ✅ pass — artifact_contract_validation=pass | 0ms |
| 2 | `git ls-files | grep -E '(\.onnx$|libtokenizers\.a$)'` | 0 | ✅ pass — tracked_native_onnx_binaries=0 | 0ms |
| 3 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — no changed indexed symbols | 0ms |

## Deviations

None.

## Known Issues

None for S01. Future CI integration still needs artifact provisioning design.

## Files Created/Modified

- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
