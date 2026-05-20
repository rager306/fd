---
id: T01
parent: S01
milestone: M022-i079tk
key_files:
  - api/Dockerfile
  - api/.dockerignore
  - tools/verify_onnx_artifacts.py
  - docs/onnx-artifacts/README.md
key_decisions:
  - Use a generated staging Docker context under `.gsd/runtime/` instead of root context or `api/` context.
  - Packaging script will verify artifacts first, then stage only required files: API source, manifests, ONNX binary, native tokenizer static lib, tokenizer JSON, and ONNX Runtime shared library.
  - Dedicated ONNX image will use explicit `onnx hf_tokenizers` build tags and will not alter the default Dockerfile.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:33:32.148Z
blocker_discovered: false
---

# T01: Chose a generated staging context for the dedicated ONNX Docker packaging proof.

**Chose a generated staging context for the dedicated ONNX Docker packaging proof.**

## What Happened

Chose the ONNX packaging strategy. Root Docker context would be too broad and risks including unrelated files; `api/` context cannot access external ignored artifacts. A generated staging context under `.gsd/runtime/` provides explicit artifact provisioning and keeps the default Dockerfile untouched.

## Verification

Confirmed the planned tracked paths do not already exist and inspected packaging boundary inputs.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test ! -e Dockerfile.onnx && test ! -e tools/build_onnx_image.sh` | 0 | ✅ pass — new_onnx_packaging_paths_ok | 0ms |
| 2 | `read api/Dockerfile api/.dockerignore tools/verify_onnx_artifacts.py docs/onnx-artifacts/README.md` | 0 | ✅ pass — packaging inputs inspected | 0ms |

## Deviations

None.

## Known Issues

Staging context may be large because it must include the ONNX binary; that is expected for local opt-in packaging proof and remains untracked under `.gsd/runtime/`.

## Files Created/Modified

- `api/Dockerfile`
- `api/.dockerignore`
- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`
