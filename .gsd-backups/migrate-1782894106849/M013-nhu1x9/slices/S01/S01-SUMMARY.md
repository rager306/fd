---
id: S01
parent: M013-nhu1x9
milestone: M013-nhu1x9
provides:
  - Validated native tokenizer artifact contract.
  - Canonical ignored local artifact path for S02 tagged build work.
requires:
  []
affects:
  - S02
key_files:
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt
  - tools/compare_tokenizers.py
  - .gitignore
key_decisions:
  - Canonical local native artifact path: `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a`.
  - Tracked manifest path: `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`.
  - Static native libraries are ignored via `*.a`.
patterns_established:
  - Track manifests, not native binaries.
  - Validate native artifacts before build/tag work.
  - Keep default build safety as a first-class gate for experimental ONNX dependencies.
observability_surfaces:
  - Manifest failure contract for missing file, checksum mismatch, wrong platform, and default build dependency.
  - Validation artifact records expected/actual checksum and size.
  - S01 task summaries record commands and outcomes.
drill_down_paths:
  - .gsd/milestones/M013-nhu1x9/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T03:40:50.578Z
blocker_discovered: false
---

# S01: Native tokenizer artifact contract

**S01 created and verified the native HF tokenizer artifact contract without committing the binary or affecting default builds.**

## What Happened

S01 established the native HF tokenizer artifact contract. It added a tracked manifest for the linux-amd64 `libtokenizers.a` static library, added a static-library ignore rule, copied the local prototype artifact into the canonical ignored `.gsd/runtime` path, added repeatable validation in `tools/compare_tokenizers.py`, and produced validation evidence. Verification confirmed the local native artifact matches size/SHA256, no native binaries are tracked, no raw probe text is leaked, and default tests/lint still pass.

## Verification

S01 final verification passed: native_manifest_hygiene=pass, raw_probe_text_leaks=0, tracked_native_binaries=0, Go tests 78 passed, lint 0 issues, GitNexus scope limited to comparator tool flows.

## Requirements Advanced

- M012-native-packaging-requirement — Created the native artifact manifest and validation needed before opt-in HF tokenizer build integration.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The manifest uses the upstream `latest` release asset URL for now. This is acceptable for prototype packaging but should be pinned to a release tag before production/CI hardening.

## Known Limitations

Native artifact download/pinning is not automated yet. CI/Docker integration is not implemented yet. The binary remains local-only under `.gsd/runtime`.

## Follow-ups

S02 should add an opt-in build tag/package boundary that can use this artifact path without affecting default builds.

## Files Created/Modified

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` — Native HF tokenizers artifact manifest.
- `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt` — Native artifact validation evidence.
- `tools/compare_tokenizers.py` — Native artifact validation mode.
- `.gitignore` — Static native library ignore rule.
