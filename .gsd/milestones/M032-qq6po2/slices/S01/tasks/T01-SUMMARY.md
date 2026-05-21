---
id: T01
parent: S01
milestone: M032-qq6po2
key_files:
  - tools/verify_onnx_export_contract.py
key_decisions:
  - Verifier claim scope is explicitly `existing_artifact_contract_verification_not_regenerated_export` to avoid byte-for-byte reproducibility overclaiming.
duration: 
verification_result: passed
completed_at: 2026-05-21T06:56:03.164Z
blocker_discovered: false
---

# T01: Implemented and positively verified the local ONNX export contract verifier.

**Implemented and positively verified the local ONNX export contract verifier.**

## What Happened

Added `tools/verify_onnx_export_contract.py`, a local verifier for the existing USER-bge-m3 ONNX export contract. The verifier reads the tracked ONNX manifest, source provenance, and export metadata; validates production_default=false, artifact size/sha, source file checksums, model revision, export toolchain pins, output metadata, CPU provider evidence, and bounded claim semantics. Positive run passed against the current local ignored artifact.

## Verification

`python3 -m py_compile tools/verify_onnx_export_contract.py && python3 tools/verify_onnx_export_contract.py` passed with verdict `pass`, artifact verified, source_files_verified=4.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/verify_onnx_export_contract.py && python3 tools/verify_onnx_export_contract.py` | 0 | ✅ pass — existing artifact contract verified; claim_scope is not regenerated export | 6700ms |

## Deviations

Verifier implementation initially treated Python as a package entry in export metadata. Fixed it to compare the top-level Python runtime string separately while still requiring manifest export packages to include Python.

## Known Issues

Verifier proves the current local artifact matches tracked manifest/provenance/export metadata. It does not regenerate the ONNX binary and does not remove the need for immutable external hosting or a future reproducible-export gate.

## Files Created/Modified

- `tools/verify_onnx_export_contract.py`
