---
id: T02
parent: S01
milestone: M025-9bvjxa
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Provisioning helper has no default artifact URLs; required ONNX/native sources must be explicit.
  - Dry-run exits 0 and reports missing source blockers; non-dry-run fails fast if required sources are missing.
  - Native tokenizer source may be a tar archive containing `libtokenizers.a`; checksum verification remains mandatory after extraction.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:39:27.355Z
blocker_discovered: false
---

# T02: Implemented the checksum-verified ONNX artifact provisioning helper.

**Implemented the checksum-verified ONNX artifact provisioning helper.**

## What Happened

Added `tools/provision_onnx_artifacts.py`. It reads the existing ONNX and native tokenizer manifests, supports dry-run planning, requires explicit sources for required artifacts in non-dry-run mode, copies/downloads sources, extracts `libtokenizers.a` from archives when needed, and verifies size/sha before accepting artifacts. Dry-run and missing-source failure behavior were verified, and the existing strict local verifier still passes.

## Verification

Python compile, dry-run blocker output, missing-source failure, and strict local verifier passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py && python3 tools/provision_onnx_artifacts.py --dry-run ...` | 0 | ✅ pass — provision_dry_run_blockers=pass | 13900ms |
| 2 | `python3 tools/provision_onnx_artifacts.py ... without required sources` | 0 | ✅ pass — provision_missing_source_failure=pass | 13800ms |
| 3 | `python3 tools/verify_onnx_artifacts.py --onnx-manifest ... --native-tokenizer-manifest ...` | 0 | ✅ pass — strict_local_verifier=pass | 13800ms |

## Deviations

None.

## Known Issues

The helper supports ONNX/native artifacts and optional tokenizer JSON/ONNX Runtime source, but the project still needs immutable external sources for hosted CI.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
