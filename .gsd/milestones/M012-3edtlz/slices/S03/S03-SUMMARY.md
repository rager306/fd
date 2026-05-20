---
id: S03
parent: M012-3edtlz
milestone: M012-3edtlz
provides:
  - Parity-passing candidate path.
  - Native packaging blocker evidence for S04 decision.
requires:
  []
affects:
  - S04
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt
key_decisions:
  - HF Rust tokenizers binding is the correct parity path.
  - Do not integrate native dependency into default builds without build-tag/package design.
  - Keep runtime ONNX backend blocked until native tokenizer packaging is solved.
patterns_established:
  - Prove native dependency parity in isolation before touching runtime dependencies.
  - Correctness can be solved while packaging remains a separate blocker.
  - Build-time native dependencies need build-tag/package design to avoid breaking default TEI builds.
observability_surfaces:
  - `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` records exact pass evidence.
  - S03 summaries record the native library/linker requirements and integration blocker.
drill_down_paths:
  - .gsd/milestones/M012-3edtlz/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S03/tasks/T03-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S03/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T02:22:17.207Z
blocker_discovered: false
---

# S03: HF tokenizer binding feasibility and parity

**S03 proved exact tokenizer parity through HF Rust tokenizers bindings, but deferred runtime integration on native packaging safety.**

## What Happened

S03 tested the recommended HF Rust tokenizers binding. A temp-module feasibility probe successfully loaded `tokenizer.json` through `github.com/daulet/tokenizers` and matched one known baseline probe. The comparator was then extended with `go-hf-binding` mode, which compared all five fixed probes against the S01 Hugging Face baseline and passed exactly. Runtime integration was intentionally deferred because the binding requires `libtokenizers.a` and CGO linker configuration; adding that import to default API code would break normal builds unless native packaging and build tags are designed first.

## Verification

HF binding artifact passed all probes; parser/leak checks passed; Go tests `78 passed`; lint `0 issues`; default health ok; GitNexus medium scope limited to comparator tool flows.

## Requirements Advanced

- M011-surfaced-tokenizer-parity — Proved an exact tokenization parity path exists via HF Rust tokenizers binding.

## Requirements Validated

None.

## New Requirements Surfaced

- Native HF tokenizers integration requires build-tag and artifact packaging design before runtime code changes.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S03 did not integrate the binding into runtime code because native library packaging/build tags are unresolved. This is a safety deviation from the 'integrate if parity passes' branch: parity passed, but packaging is not yet safe for default builds.

## Known Limitations

The passing artifact depends on a temp local prebuilt `libtokenizers.a`. CI/Docker/native artifact management is not implemented. `api/embed/onnx.go` still uses `sugarme/tokenizer` and therefore remains semantically blocked if run as-is.

## Follow-ups

S04 should recommend a next milestone for build-tag/native packaging integration if continuing pure Go ONNX work. Candidate direction: `hf_tokenizers` build tag, explicit libtokenizers artifact manifest, Docker/CI packaging, then replace ONNX tokenizer path only in tagged opt-in builds.

## Files Created/Modified

- `tools/compare_tokenizers.py` — Extended comparator with go-hf-binding mode for isolated HF Rust tokenizers binding comparison.
- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` — Passing HF Rust tokenizers binding parity artifact.
