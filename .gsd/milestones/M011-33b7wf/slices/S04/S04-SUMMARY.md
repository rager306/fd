---
id: S04
parent: M011-33b7wf
milestone: M011-33b7wf
provides:
  - Final M011 recommendation.
  - Next milestone direction: tokenizer parity.
  - Safety verification evidence for default runtime.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M011-33b7wf/slices/S04/S04-RESEARCH.md
  - benchmark-results/fd-go-onnx-m011-s03.txt
key_decisions:
  - Close M011 as a blocked prototype, not a production-ready ONNX backend.
  - Do not benchmark ONNX speed until tokenization parity passes.
  - Future backend comparisons must isolate Redis cache namespace.
patterns_established:
  - Evidence-backed blocker closure is preferable to invalid performance claims.
  - Backend comparison must isolate cache namespace.
  - Tokenizer parity is a prerequisite for embedding runtime equivalence.
observability_surfaces:
  - S04 research artifact records blocker, cache namespace pitfall, future gates, and verification commands.
  - T02 summary records fresh Go/lint/Compose/health/manifest/GitNexus verification evidence.
drill_down_paths:
  - .gsd/milestones/M011-33b7wf/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S04/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T01:52:36.478Z
blocker_discovered: false
---

# S04: Blocker synthesis and recommendation

**S04 closed M011 with an honest blocker recommendation and verified TEI default safety.**

## What Happened

S04 synthesized the M011 outcome and verified the project safety boundary. It documented that M011 successfully added manifest validation, opt-in runtime selection, and a real Go ONNX load/run path, but blocked semantic equivalence on tokenizer mismatch. It also verified the default TEI stack remains healthy and quality gates pass. The recommendation is to stop M011 here and plan the next milestone around tokenizer parity before any ONNX throughput benchmark or production-readiness claim.

## Verification

Fresh S04 verification passed: Go tests `78 passed in 4 packages`; pinned GolangCI-Lint `0 issues`; Compose config and health passed; manifest/comparison/tracked-artifact checks passed; GitNexus reported low risk and no affected processes.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- Tokenizer parity must be validated before ONNX backend performance benchmarking or production recommendation.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S04 intentionally did not run an ONNX throughput benchmark because S03 failed semantic equivalence. This changes the milestone outcome from 'prototype performance evidence' to 'blocked prototype with evidence'.

## Known Limitations

Tokenizer parity is unresolved. Stable ONNX Runtime shared library packaging is unresolved. Larger Russian/legal corpus validation remains future work.

## Follow-ups

Plan M012 for tokenizer parity: build Python HF token baseline, compare Go tokenizer candidates, require token-level equality before rerunning cosine and performance benchmarks.

## Files Created/Modified

- `.gsd/milestones/M011-33b7wf/slices/S04/S04-RESEARCH.md` — Final M011 blocker synthesis and recommendation.
- `benchmark-results/fd-go-onnx-m011-s03.txt` — Failed isolated-cache Go ONNX API comparison artifact.
- `api/main.go` — Opt-in ONNX backend code and default TEI-safe wiring verified by S04.
