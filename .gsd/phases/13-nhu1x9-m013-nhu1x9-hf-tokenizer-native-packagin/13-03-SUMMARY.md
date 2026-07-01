---
id: S03
parent: M013-nhu1x9
milestone: M013-nhu1x9
provides:
  - Semantically equivalent tagged ONNX path for fixed probes.
  - Benchmark-ready ONNX path for future milestone.
requires:
  []
affects:
  - S04
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
key_decisions:
  - Tagged `hf_tokenizers` ONNX runtime path is semantically equivalent on fixed probes.
  - Default untagged build remains unchanged and green.
  - Performance benchmarking should target only the tagged HF tokenizer ONNX path.
patterns_established:
  - Use tokenizer interface plus build tags to preserve default behavior while enabling native opt-in paths.
  - Backend cosine comparisons must isolate Redis namespace.
  - Cosine pass unlocks benchmarking, but not production switch.
observability_surfaces:
  - `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt` records fixed-probe cosine pass and cache namespace.
  - Task summaries record tagged run command, native library path, and verification gates.
drill_down_paths:
  - .gsd/milestones/M013-nhu1x9/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S03/tasks/T03-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S03/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T03:56:23.414Z
blocker_discovered: false
---

# S03: Parity correct ONNX tokenizer integration

**S03 integrated the HF native tokenizer into tagged ONNX builds and passed fixed-probe TEI-vs-ONNX cosine equivalence.**

## What Happened

S03 integrated the parity-correct HF native tokenizer behind the `hf_tokenizers` build tag. `ONNXEmbedder` now uses a small tokenizer interface; untagged builds keep the existing `sugarme` path, while tagged builds construct the HF native tokenizer. The tagged API started locally and passed TEI-vs-ONNX cosine equivalence with isolated Redis namespace `m013-hf-tokenizer`. Default tests/lint and tagged tests pass, and the tagged server was cleaned up.

## Verification

Fresh S03 verification passed: default Go tests 78 passed, lint 0 issues, tagged embed tests 20 passed, cosine artifact PASS, no leaks, no background processes, GitNexus medium scope verified.

## Requirements Advanced

- M012-native-packaging-requirement — Moved from isolated tokenizer parity to tagged runtime ONNX cosine equivalence on fixed probes.

## Requirements Validated

None.

## New Requirements Surfaced

- Tagged Docker/CI build support should be added before broader team/CI use of the native ONNX path.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S03 succeeded and therefore did not produce a blocker artifact. It did not run throughput benchmarking; that remains out of scope until S04 final decision and a separate benchmark milestone.

## Known Limitations

The native artifact is still local-only. Tagged Docker/CI build support is not yet implemented. Fixed probes are not a large Russian/legal retrieval corpus.

## Follow-ups

S04 should validate the milestone and recommend an ONNX performance benchmark milestone for the tagged HF tokenizer path, while still requiring Docker/CI packaging and larger Russian/legal quality gates before production switch.

## Files Created/Modified

- `api/embed/onnx.go` — ONNX embedder tokenizer abstraction and resource close path.
- `api/embed/onnx_tokenizer_default.go` — Default untagged sugarme tokenizer adapter.
- `api/embed/onnx_tokenizer_hf.go` — Tagged HF native tokenizer adapter.
- `api/embed/hf_tokenizer_native.go` — Native HF tokenizer wrapper now satisfies ONNX tokenizer interface.
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt` — Tagged ONNX cosine comparison artifact.
