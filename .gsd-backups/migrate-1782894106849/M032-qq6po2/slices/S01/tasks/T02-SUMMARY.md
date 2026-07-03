---
id: T02
parent: S01
milestone: M032-qq6po2
key_files:
  - tools/verify_onnx_export_contract.py
  - .gsd/exec/2b3f90be-849f-4c8f-af94-b87cd7aed5f1.stdout
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T06:56:42.386Z
blocker_discovered: false
---

# T02: Verified export contract verifier failure modes.

**Verified export contract verifier failure modes.**

## What Happened

Ran three negative probes using temporary tampered copies of the manifest/provenance/export metadata. The verifier failed explicitly on artifact checksum mismatch, model revision mismatch, and transformers version mismatch. Error outputs were structured JSON with labels and safe repo-relative paths.

## Verification

Tampered metadata probes failed as expected with explicit labels: `artifact_sha256`, `model_revision`, and `export_metadata`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec M032 verifier negative probes` | 0 | ✅ pass — all tampered probes failed as expected with sanitized structured errors | 3071ms |

## Deviations

None.

## Known Issues

Negative probes do not exercise every possible failure, but cover the core artifact checksum, provenance revision, and export toolchain mismatch classes.

## Files Created/Modified

- `tools/verify_onnx_export_contract.py`
- `.gsd/exec/2b3f90be-849f-4c8f-af94-b87cd7aed5f1.stdout`
