---
id: T02
parent: S01
milestone: M021-4t2wpt
key_files:
  - tools/verify_onnx_artifacts.py
  - docs/onnx-artifacts/README.md
key_decisions:
  - Added local verifier script rather than modifying Dockerfile/CI in this milestone.
  - README documents `--allow-missing` as contract-shape only for default CI, not ONNX readiness evidence.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:15:21.961Z
blocker_discovered: false
---

# T02: Added the ONNX/native artifact verification contract and documentation.

**Added the ONNX/native artifact verification contract and documentation.**

## What Happened

Added `tools/verify_onnx_artifacts.py` to validate local ONNX and native tokenizer artifacts against tracked manifests, including checksum, size, production_default false, artifact.git_tracked false, and git-tracked path hygiene. Added `docs/onnx-artifacts/README.md` documenting required artifacts, verification command, validated ONNX 1024 runtime env, existing evidence, and future Docker/CI gate requirements.

## Verification

Script compilation, strict verification, and allow-missing contract-shape check all passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/verify_onnx_artifacts.py` | 0 | ✅ pass — script compiles | 0ms |
| 2 | `python3 tools/verify_onnx_artifacts.py --onnx-manifest docs/onnx-artifacts/user-bge-m3-dense-fp32.json --native-tokenizer-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` | 0 | ✅ pass — verified_all_present=true | 0ms |
| 3 | `python3 tools/verify_onnx_artifacts.py ... --allow-missing` | 0 | ✅ pass — allow_missing_contract_check=pass | 0ms |

## Deviations

None.

## Known Issues

Verifier does not download artifacts. Future Docker/CI gate must implement artifact provisioning/download/cache path.

## Files Created/Modified

- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`
