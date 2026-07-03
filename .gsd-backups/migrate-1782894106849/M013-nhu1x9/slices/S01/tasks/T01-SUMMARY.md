---
id: T01
parent: S01
milestone: M013-nhu1x9
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - .gitignore
key_decisions:
  - Track the native tokenizer manifest at `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`.
  - Use ignored local artifact path `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` for future reproducible local setup; S03 temp evidence remains `/tmp/fd-daulet-tokenizers-probe/libtokenizers.a`.
  - Add explicit `.gitignore` coverage for static native libraries (`*.a`) so `libtokenizers.a` cannot be accidentally committed.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:36:27.834Z
blocker_discovered: false
---

# T01: Defined the native tokenizer artifact contract: tracked manifest, ignored local `.gsd/runtime` binary, and explicit static-library ignore rule.

**Defined the native tokenizer artifact contract: tracked manifest, ignored local `.gsd/runtime` binary, and explicit static-library ignore rule.**

## What Happened

Inspected the existing ONNX artifact manifest and `.gitignore`. The native tokenizer artifact should follow the same pattern as ONNX: track JSON manifest and docs, keep large/native binary under `.gsd/runtime/`, and validate checksum before use. The local temp artifact currently exists at `/tmp/fd-daulet-tokenizers-probe/libtokenizers.a`, size 49,776,794 bytes, SHA256 `e6862b31745bb7d07980fcee70e49cd3b4318097609180f5d2d3fb394f305d50`. The future canonical local path should be `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a`.

## Verification

Read the existing ONNX manifest and `.gitignore`, and checked the local temp `libtokenizers.a` size/SHA256. The task summary names tracked manifest path, ignored local path, and binary exclusion rule.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — ONNX manifest pattern identified | 0ms |
| 2 | `read .gitignore` | 0 | ✅ pass — `.gsd/runtime/` ignored; `*.a` not yet ignored | 0ms |
| 3 | `ls -l /tmp/fd-daulet-tokenizers-probe/libtokenizers.a && sha256sum ...` | 0 | ✅ pass — size 49776794; sha256 e6862b31745bb7d07980fcee70e49cd3b4318097609180f5d2d3fb394f305d50 | 0ms |

## Deviations

None.

## Known Issues

Current `.gitignore` ignores `.so` and `.dylib`, but not `.a`; S01 T02 should add `*.a`.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `.gitignore`
