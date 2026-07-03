---
id: T03
parent: S01
milestone: M013-nhu1x9
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
key_decisions:
  - `validate-native-artifact` mode validates manifest shape, local artifact existence, size, SHA256, and `git_tracked=false`.
  - The validator is implemented in the existing tokenizer comparison tool to keep tokenizer/manifest checks in one place.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:39:22.413Z
blocker_discovered: false
---

# T03: Added and ran repeatable native tokenizer artifact validation; local `libtokenizers.a` matches the manifest and remains untracked.

**Added and ran repeatable native tokenizer artifact validation; local `libtokenizers.a` matches the manifest and remains untracked.**

## What Happened

Extended `tools/compare_tokenizers.py` with `validate-native-artifact` mode. The mode reads the native tokenizer artifact manifest, validates required artifact metadata, verifies the local `libtokenizers.a` size and SHA256, confirms the manifest declares `git_tracked=false`, and renders a sanitized validation artifact. Running it against `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` passed. A tracked-file scan confirmed no `.a` files are tracked.

## Verification

Python compile passed. Native artifact verifier exited 0 and wrote PASS. `git ls-files '*.a'` returned no tracked static libraries.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py` | 0 | ✅ pass | 0ms |
| 2 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode validate-native-artifact --native-artifact-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json --output benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt` | 0 | ✅ pass — size_ok=true; sha256_ok=true; git_tracked_false=true | 0ms |
| 3 | `git ls-files '*.a' '.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a'` | 0 | ✅ pass — no tracked static native libraries | 0ms |

## Deviations

Added a tracked validation evidence artifact at `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt` in addition to the manifest, so future agents can inspect the exact validation result.

## Known Issues

The validator checks local file integrity but does not download the upstream archive. Download automation remains future S02/S03 packaging work.

## Files Created/Modified

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
