---
id: S04
parent: M016-pdcjat
milestone: M016-pdcjat
provides:
  - Candidate model shortlist for future benchmark milestone.
requires:
  []
affects:
  - S03
key_files:
  - benchmark-results/fd-model-alternatives-m016-s04.txt
key_decisions:
  - Top candidates: corrected BGE-M3/USER-BGE-M3 long-text path, Qwen3-Embedding-0.6B, multilingual E5 large/instruct, ai-forever ru-en-RoSBERTa, Russian E5 derivatives, deepvk USER-base.
  - Rerankers are second-stage quality candidates, not embedding replacements.
  - MiniLM-style models are reference floors, not replacement candidates.
patterns_established:
  - Alternative model adoption requires legal-corpus evidence, not benchmark leaderboard claims alone.
  - Long-context support and prompt policy must be recorded for each candidate.
observability_surfaces:
  - Source-backed model research artifact.
  - Benchmark protocol for future candidate trials.
drill_down_paths:
  - .gsd/milestones/M016-pdcjat/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S04/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T05:20:37.474Z
blocker_discovered: false
---

# S04: Alternative model research

**S04 produced the alternative-model shortlist and fair benchmark protocol.**

## What Happened

S04 researched candidate models for future Russian/legal embedding benchmarks. It ranked candidates by fit with fd constraints: Russian/legal suitability, long-context handling, 1024-dimensional compatibility, CPU feasibility, ONNX/export risk, and prompt/chunking requirements. It produced a benchmark protocol that requires sanitized 44-ФЗ gate runs before any replacement decision.

## Verification

Artifact read and required-token checks passed.

## Requirements Advanced

- alternative-model-evaluation — Defined alternative-model research track and benchmark protocol.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S04 was executed early because it is independent research and does not depend on the ONNX root-cause slices.

## Known Limitations

No alternative models were downloaded or benchmarked in this slice. Research is source-based and must be verified on the 44-ФЗ gate before adoption.

## Follow-ups

Use this artifact when planning model trial milestones. Do not benchmark alternatives until the evaluator supports model-specific prompts/dimensions/chunking. Keep current long-text divergence investigation as the main path.

## Files Created/Modified

- `benchmark-results/fd-model-alternatives-m016-s04.txt` — Ranked alternative model research and benchmark protocol.
