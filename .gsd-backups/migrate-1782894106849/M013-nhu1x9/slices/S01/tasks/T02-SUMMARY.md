---
id: T02
parent: S01
milestone: M013-nhu1x9
key_files:
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - .gitignore
key_decisions:
  - Manifest records both canonical ignored local path and S03 temp evidence path.
  - Manifest explicitly states default builds must not require the native artifact.
  - `.gitignore` now ignores `*.a` in addition to `.gsd/runtime/`, `.so`, and `.dylib`.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:37:39.739Z
blocker_discovered: false
---

# T02: Wrote the native HF tokenizer artifact manifest and static-library ignore rule, with local checksum validation passing.

**Wrote the native HF tokenizer artifact manifest and static-library ignore rule, with local checksum validation passing.**

## What Happened

Created `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` describing the HF Rust tokenizer native artifact contract. It records module/library provenance, linux-amd64 platform, source URL, local ignored path, size, SHA256, expected build tag, linker environment, validation artifacts, and failure contract. Added `*.a` to `.gitignore` to prevent static native libraries from being committed. Verified the manifest parses and the local artifact at `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` matches the expected size and checksum.

## Verification

Manifest JSON parsed, local artifact size/SHA256 matched, and `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` is ignored by git.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 manifest parse/checksum copy from /tmp to .gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | 0 | ✅ pass — native_manifest_parse=pass; native_artifact_sha256=pass | 0ms |
| 2 | `git check-ignore .gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | 0 | ✅ pass — canonical local native artifact is ignored | 0ms |
| 3 | `git status --short` | 0 | ✅ pass — native binary not listed; only manifest/gitignore/GSD docs changed | 0ms |

## Deviations

Copied the S03 temp `libtokenizers.a` into the canonical ignored `.gsd/runtime/tokenizers/linux-amd64/` path so future S02 tagged-build work can use the manifest path. The binary remains ignored and untracked.

## Known Issues

The manifest source URL uses the `latest` release asset URL; a future hardening pass should pin an exact release tag if upstream publishes stable version tags suitable for reproducible CI.

## Files Created/Modified

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `.gitignore`
