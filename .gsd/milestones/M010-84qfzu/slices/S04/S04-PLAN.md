# S04: ONNX spike recommendation

**Goal:** Synthesize M010 evidence into a clear ONNX spike recommendation: whether to proceed, what remains blocked, and which gates are required before any production adapter work.
**Demo:** After this, we have a clear decision: proceed to adapter implementation, do more artifact work, or stop ONNX path.

## Must-Haves

- Recommendation clearly states proceed/stop/continue-research.
- Evidence references S01 research, S02 TEI baseline, and S03 ONNX comparison.
- Future gates are explicit: artifact distribution, performance benchmark, Russian/legal quality corpus, dependency pin, and non-default adapter.
- Production runtime remains unchanged.
- Milestone can be validated after S04.

## Proof Level

- This slice proves: Research synthesis plus verification of all S01-S03 evidence artifacts.

## Integration Closure

Produces final milestone recommendation and updates decision/requirement state if needed. This closes the spike without changing production runtime defaults.

## Verification

- Captures evidence paths, artifact hashes, dependency pinning, and required future benchmark/quality gates for downstream agents.

## Tasks

- [x] **T01: Synthesize ONNX spike recommendation** `est:small`
  Synthesize S01-S03 evidence into `S04-RESEARCH.md`: exact model provenance, successful FP32 ONNX export/load, cosine comparison results, dependency pin issue, limitations, and recommendation. Include proceed/stop criteria and future implementation gates.
  - Files: `.gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md`
  - Verify: Research artifact exists and states recommendation, limitations, required gates, and no production runtime change.

- [x] **T02: Verify M010 spike evidence** `est:small`
  Run final evidence checks for M010: Go tests, lint, Python comparator script compile checks, artifact presence checks, raw-probe leakage checks, git status scope, and GitNexus change detection. Prepare milestone validation inputs.
  - Verify: Go tests/lint pass; Python scripts compile; artifact parser checks pass; `gitnexus_detect_changes` reports expected scope.

## Files Likely Touched

- .gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md
