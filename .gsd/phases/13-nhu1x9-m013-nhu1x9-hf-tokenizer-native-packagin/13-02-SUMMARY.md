---
id: S02
parent: M013-nhu1x9
milestone: M013-nhu1x9
provides:
  - Opt-in native tokenizer package boundary.
  - Tagged parity proof for S03 runtime integration.
requires:
  []
affects:
  - S03
key_files:
  - api/embed/hf_tokenizer_native.go
  - api/embed/hf_tokenizer_native_test.go
  - benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt
key_decisions:
  - Use build tag `hf_tokenizers` for native tokenizer code.
  - Keep default builds free of native artifact/linker requirements.
  - Use `CGO_LDFLAGS=-L../.gsd/runtime/tokenizers/linux-amd64` for tagged tests from `api`.
patterns_established:
  - Native dependencies stay behind build tags until fully integrated.
  - Default build isolation is a required proof before runtime integration.
  - Tagged parity artifacts document native build commands and checksums.
observability_surfaces:
  - Tagged parity artifact records command, native artifact checksum, and PASS output.
  - Task summaries record default-vs-tagged verification commands.
drill_down_paths:
  - .gsd/milestones/M013-nhu1x9/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S02/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T03:47:19.200Z
blocker_discovered: false
---

# S02: Opt in build tag boundary

**S02 created a safe opt-in `hf_tokenizers` build boundary: default builds pass, tagged native parity passes.**

## What Happened

S02 implemented the build-tag boundary. Native HF tokenizer imports are isolated to `//go:build hf_tokenizers` files. Default tests/lint pass without native flags, proving default TEI builds are not affected. Tagged tests pass with the validated `libtokenizers.a`, proving the project-local native tokenizer path can match the M012 baseline. The tagged parity proof is persisted in `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`.

## Verification

Fresh verification passed: default Go tests 78 passed; lint 0 issues; tagged test 1 passed; artifact/leak/native scan passed; GitNexus low risk.

## Requirements Advanced

- M012-native-packaging-requirement — Proved native HF tokenizer dependency can be isolated from default builds while preserving token parity under tagged conditions.

## Requirements Validated

None.

## New Requirements Surfaced

- Docker/CI tagged build support remains required before production use.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Tagged test duplicates fixed probe texts in test code to avoid adding runtime data loaders. This does not affect artifacts; raw probe text leakage checks for artifacts pass.

## Known Limitations

Docker/CI do not yet run tagged builds. Runtime ONNX embedder still uses `sugarme` until S03 changes it.

## Follow-ups

S03 can now integrate the parity-correct tokenizer into the opt-in ONNX path under `hf_tokenizers` build tag and rerun isolated-cache cosine if integration succeeds.

## Files Created/Modified

- `api/embed/hf_tokenizer_native.go` — Tagged native HF tokenizer wrapper.
- `api/embed/hf_tokenizer_native_test.go` — Tagged native tokenizer parity test.
- `api/go.mod / api/go.sum` — Go dependency updates for `github.com/daulet/tokenizers`.
- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt` — Tagged parity artifact.
