---
id: T04
parent: S01
milestone: M013-nhu1x9
key_files:
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt
  - tools/compare_tokenizers.py
  - .gitignore
key_decisions:
  - S01 artifact contract is safe: native binary is ignored/untracked, manifest validates, and default builds pass.
  - S02 can use `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` as the canonical local artifact path when probing build tags.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:40:20.994Z
blocker_discovered: false
---

# T04: Verified the native tokenizer artifact contract: checksum passes, binary is ignored, default tests/lint pass, and no raw/native leaks are tracked.

**Verified the native tokenizer artifact contract: checksum passes, binary is ignored, default tests/lint pass, and no raw/native leaks are tracked.**

## What Happened

Ran S01 final verification. The tokenizer comparison tool compiles. Native artifact manifest and validation artifact parse and contain no raw probe text. No tracked native binaries were found. The canonical local `libtokenizers.a` is ignored. Default Go tests and pinned lint pass. GitNexus reports medium risk limited to comparator tool flows, with no API runtime changes.

## Verification

Fresh S01 verification passed: py_compile, manifest/artifact hygiene, no raw probe leaks, no tracked native binaries, default Go tests/lint, and GitNexus scope review.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py && uv run --python 3.13 --with transformers --with torch --with sentencepiece python native artifact hygiene check` | 0 | ✅ pass — native_manifest_hygiene=pass; raw_probe_text_leaks=0; tracked_native_binaries=0 | 0ms |
| 2 | `git check-ignore .gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` | 0 | ✅ pass — canonical native artifact path ignored | 0ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 78 passed in 4 packages | 15800ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 15700ms |
| 5 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ⚠️ medium — affected processes are comparator tool flows only | 0ms |

## Deviations

GitNexus reports medium risk because `tools/compare_tokenizers.py` changed its own comparator flows; API runtime flows are not affected. This is acceptable for S01 artifact tooling.

## Known Issues

Manifest source URL still points at `latest`; future CI hardening should pin a release if available. Native binary is local-only under `.gsd/runtime` and not committed.

## Files Created/Modified

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt`
- `tools/compare_tokenizers.py`
- `.gitignore`
